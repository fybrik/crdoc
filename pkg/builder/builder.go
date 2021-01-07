package builder

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type ModelBuilder struct {
	Model              *Model
	Strict             bool
	TemplatesDirOrFile string
	OutputFilepath     string
	Links              map[string]int
}

// Write outputs markdown to the output direcory
func (b *ModelBuilder) Write() error {
	outputFilepath := filepath.Clean(b.OutputFilepath)

	// create dirs if needed
	err := os.MkdirAll(filepath.Dir(outputFilepath), os.ModePerm)
	if err != nil {
		return err
	}

	// create the file
	f, err := os.Create(outputFilepath)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf("Error closing file: %s\n", err)
		}
	}()

	// Load and process template
	templatesFilepath := filepath.Clean(b.TemplatesDirOrFile)
	info, err := os.Stat(templatesFilepath)
	if err != nil {
		return err
	}

	var t *template.Template
	if info.IsDir() {
		t = template.Must(template.New("main.tmpl").Funcs(sprig.TxtFuncMap()).ParseGlob(path.Join(templatesFilepath, "*")))
	} else {
		t = template.Must(template.New(filepath.Base(templatesFilepath)).Funcs(sprig.TxtFuncMap()).ParseFiles(b.TemplatesDirOrFile))
	}
	return t.Execute(f, *b.Model)
}

// Add adds a CustomResourceDefinition to the model
func (b *ModelBuilder) Add(crd *apiextensions.CustomResourceDefinition) error {
	// Add chapter for each version
	for _, version := range crd.Spec.Versions {
		group := crd.Spec.Group
		gv := fmt.Sprintf("%s/%s", group, version.Name)
		kind := crd.Spec.Names.Kind

		// Find matching group/version
		groupModel := b.Model.findGroupModel(group, version.Name)
		if groupModel == nil {
			if b.Strict {
				log.Warn(fmt.Sprintf("group/version not found in TOC: %s", gv))
				continue
			}
			groupModel = &GroupModel{
				Group:   group,
				Version: version.Name,
			}
			b.Model.Groups = append(b.Model.Groups, groupModel)
		}

		// Find matching kind
		kindModel := groupModel.findKindModel(kind)
		if kindModel == nil {
			if b.Strict {
				log.Warn(fmt.Sprintf("group/version/kind not found in TOC: %s/%s", gv, kind))
				continue
			}
			kindModel = &KindModel{
				Name: kind,
			}
			groupModel.Kinds = append(groupModel.Kinds, kindModel)
		}

		// Find schema
		validation := version.Schema
		if validation == nil {
			// Fallback to resource level schema
			validation = crd.Spec.Validation
		}
		if validation == nil {
			return errors.New("missing validation field in input CRD")
		}
		schema := validation.OpenAPIV3Schema

		// Recusively add type models
		_, _ = b.addTypeModels(groupModel, kindModel, kind, schema, true)
	}

	return nil
}

func (b *ModelBuilder) createLink(name string) string {
	link := fmt.Sprintf("#%s", headingID(name))
	if b.Links == nil {
		b.Links = make(map[string]int)
	}
	if value, exists := b.Links[link]; exists {
		value += 1
		link = fmt.Sprintf("%s-%d", link, value)
	} else {
		b.Links[link] = 0
	}
	return link
}

func (b *ModelBuilder) addTypeModels(groupModel *GroupModel, kindModel *KindModel, name string, schema *apiextensions.JSONSchemaProps, isTopLevel bool) (string, *TypeModel) {
	typeName := getTypeName(schema)
	if typeName == "object" && schema.Properties != nil {
		// Create an object type model
		typeModel := &TypeModel{
			Name:        name,
			Key:         b.createLink(name),
			Description: schema.Description,
		}
		kindModel.Types = append(kindModel.Types, typeModel)

		// For each field
		for _, fieldName := range orderedPropertyKeys(schema.Required, schema.Properties, true) {
			property := schema.Properties[fieldName]

			fieldFullname := strings.Join([]string{name, fieldName}, ".")
			fieldTypename, fieldTypeModel := b.addTypeModels(groupModel, kindModel, fieldFullname, &property, false)
			var fieldTypeKey *string = nil
			if fieldTypeModel != nil {
				fieldTypeKey = &fieldTypeModel.Key
				fieldTypeModel.ParentKey = &typeModel.Key
			}

			fieldDescription := property.Description
			if fieldTypename == "enum" {
				fieldDescription = fmt.Sprintf("%s %s", fieldDescription, property.Enum)
			}

			// Create field model
			fieldModel := &FieldModel{
				Name:        fieldName,
				Type:        fieldTypename,
				TypeKey:     fieldTypeKey,
				Description: fieldDescription,
				Required:    isRequiredProperty(fieldName, schema.Required),
			}
			typeModel.Fields = append(typeModel.Fields, fieldModel)
		}
		return typeName, typeModel
	} else if typeName == "[]" {
		childTypeName, childTypeModel := b.addTypeModels(groupModel, kindModel,
			fmt.Sprintf("%s[index]", name), schema.Items.Schema, false)
		return "[]" + childTypeName, childTypeModel
	} else if typeName == "map[string]" {
		childTypeName, childTypeModel := b.addTypeModels(groupModel, kindModel,
			fmt.Sprintf("%s[key]", name), schema.AdditionalProperties.Schema, false)
		return "map[string]" + childTypeName, childTypeModel
	}
	return typeName, nil
}

func getTypeName(props *apiextensions.JSONSchemaProps) string {
	// map
	if props.Type == "object" && props.AdditionalProperties != nil {
		if props.AdditionalProperties.Schema == nil && props.AdditionalItems.Allows {
			return "map[string]string"
		}
		return "map[string]"
	}

	// array
	if props.Type == "array" {
		if props.Items == nil {
			return "[]object"
		}
		return "[]"
	}

	// enum
	if props.Type == "string" && props.Enum != nil && len(props.Enum) > 0 {
		return "enum"
	}

	// Get the value for primitive types
	if props.Format != "" && props.Type == "byte" {
		return "[]byte"
	}

	return props.Type
}

// headingID returns the ID built by hugo for a given header
func headingID(s string) string {
	result := strings.ToLower(s)
	result = strings.TrimSpace(result)
	result = regexp.MustCompile(`([^\w\- ]+)`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`(\s)`).ReplaceAllString(result, "-")
	result = regexp.MustCompile(`(\-+$)`).ReplaceAllString(result, "")

	return result
}

// orderedPropertyKeys returns the keys of m alphabetically ordered
// keys in required will be placed first
func orderedPropertyKeys(required []string, m map[string]apiextensions.JSONSchemaProps, isResource bool) []string {
	sort.Strings(required)

	if isResource {
		mkeys := make(map[string]struct{})
		for k := range m {
			mkeys[k] = struct{}{}
		}
		for _, special := range []string{"metadata", "kind", "apiVersion"} {
			if !isRequiredProperty(special, required) {
				if _, ok := mkeys[special]; ok {
					required = append([]string{special}, required...)
				}
			}
		}
	}

	keys := make([]string, len(m)-len(required))
	i := 0
	for k := range m {
		if !isRequiredProperty(k, required) {
			keys[i] = k
			i++
		}
	}
	sort.Strings(keys)
	return append(required, keys...)
}

// isRequired returns true if k is in the required array
func isRequiredProperty(k string, required []string) bool {
	for _, r := range required {
		if r == k {
			return true
		}
	}
	return false
}

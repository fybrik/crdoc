package builder

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type ModelBuilder struct {
	Model        *Model
	Strict       bool
	TemplatesDir string
	OutputDir    string
	Links        map[string]int
}

// Write outputs markdown to the output direcory
func (b *ModelBuilder) Write() error {
	filename := path.Join(b.OutputDir, "out.md")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	t := template.Must(template.New("all.tmpl").Funcs(sprig.TxtFuncMap()).ParseFiles(fmt.Sprintf("%s/all.tmpl", b.TemplatesDir)))
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
		schema := validation.OpenAPIV3Schema

		// Recusively add type models
		_ = b.addTypeModels(groupModel, kindModel, []string{}, kind, schema, true)
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

func (b *ModelBuilder) addTypeModels(groupModel *GroupModel, kindModel *KindModel,
	prefix []string, name string, schema *apiextensions.JSONSchemaProps, isTopLevel bool) *TypeModel {

	fullName := strings.Join(append(prefix, name), ".")
	typeModel := &TypeModel{
		Name:        fullName,
		Key:         b.createLink(fullName),
		Description: schema.Description,
	}
	kindModel.Types = append(kindModel.Types, typeModel)

	for _, fieldName := range orderedPropertyKeys(schema.Required, schema.Properties, true) {
		property := schema.Properties[fieldName]
		typename, typekey := getTypeNameAndKey(fieldName, property)
		fieldModel := &FieldModel{
			Name:        fieldName,
			Type:        typename,
			Description: property.Description,
			Required:    isRequiredProperty(fieldName, schema.Required),
		}
		typeModel.Fields = append(typeModel.Fields, fieldModel)

		if typekey != nil {
			var tm *TypeModel = nil
			if property.Type == "array" {
				tm = b.addTypeModels(groupModel, kindModel,
					append(prefix, name), fmt.Sprintf("%s[index]", fieldName), property.Items.Schema, false)
			} else if property.Type == "object" && property.AdditionalProperties != nil {
				tm = b.addTypeModels(groupModel, kindModel,
					append(prefix, name), fmt.Sprintf("%s[key]", fieldName), property.AdditionalProperties.Schema, false)
			} else if property.Type == "object" && property.Properties != nil {
				tm = b.addTypeModels(groupModel, kindModel,
					append(prefix, name), fieldName, &property, false)
			}
			if tm != nil {
				tm.ParentKey = &typeModel.Key
				fieldModel.TypeKey = &tm.Key
			}
		}
	}

	return typeModel
}

// getTypeNameAndKey returns the display name of a Schema type.
func getTypeNameAndKey(fieldName string, s apiextensions.JSONSchemaProps) (string, *string) {
	// Recurse if type is array
	if s.Type == "array" {
		typ, key := getTypeNameAndKey(fieldName, *s.Items.Schema)
		return fmt.Sprintf("[]%s", typ), key
	}

	// Recurse if type is map
	if s.Type == "object" && s.AdditionalProperties != nil {
		typ, key := getTypeNameAndKey(fieldName, *s.AdditionalProperties.Schema)
		return fmt.Sprintf("map[string]%s", typ), key
	}

	// Handle complex types
	if s.Type == "object" && s.Properties != nil {
		key := "object"
		return key, &key
	}

	// Get the value for primitive types
	value := s.Type
	if s.Format != "" && value == "byte" {
		value = "[]byte"
	}

	return value, nil
}

// headingID returns the ID built by hugo for a given header
func headingID(s string) string {
	result := s
	result = strings.ToLower(s)
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

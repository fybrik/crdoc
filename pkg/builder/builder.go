// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"fybrik.io/crdoc/pkg/functions"
	"github.com/Masterminds/sprig/v3"
	log "github.com/sirupsen/logrus"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// ModelBuilder is the entry point for docs generation
type ModelBuilder struct {
	Model              *Model
	Strict             bool
	TemplatesDirOrFile string
	OutputFilepath     string

	keys             map[string]int
	builtinTemplates fs.FS
}

func NewModelBuilder(model *Model, strict bool, templatesDirOrFile string, outputFilepath string, builtinTemplates embed.FS) *ModelBuilder {

	builder := &ModelBuilder{
		Model:              model,
		Strict:             strict,
		TemplatesDirOrFile: templatesDirOrFile,
		OutputFilepath:     outputFilepath,
		builtinTemplates:   builtinTemplates,
	}
	builder.keys = make(map[string]int)
	return builder
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

// Output writes markdown to the output direcory
func (b *ModelBuilder) Output() error {
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

	// Values for embedded templates
	templatesFs := b.builtinTemplates
	pattern := "templates/**.tmpl"

	dir, file := filepath.Split(b.TemplatesDirOrFile)
	if dir != "" {
		// Override to take extenal template
		templatesFs = os.DirFS(dir)
		pattern = "**.tmpl"
	}

	return template.Must(template.New(file).
		Funcs(sprig.TxtFuncMap()).
		Funcs(functions.ExportedMap).
		ParseFS(templatesFs, pattern)).
		Execute(f, *b.Model)
}

func (b *ModelBuilder) addTypeModels(groupModel *GroupModel, kindModel *KindModel, name string, schema *apiextensions.JSONSchemaProps, isTopLevel bool) (string, *TypeModel) {
	typeName := getTypeName(schema)
	if typeName == "object" && schema.Properties != nil {
		// Create an object type model
		typeModel := &TypeModel{
			Name:        name,
			Key:         b.createKey(name),
			Description: schema.Description,
			IsTopLevel:  isTopLevel,
		}
		kindModel.Types = append(kindModel.Types, typeModel)

		// For each field
		for _, fieldName := range orderedPropertyKeys(schema.Required, schema.Properties, isTopLevel) {
			property := getEnrichedProperty(schema, fieldName)

			fieldFullname := strings.Join([]string{name, fieldName}, ".")
			fieldTypename, fieldTypeModel := b.addTypeModels(groupModel, kindModel, fieldFullname, &property, false)
			var fieldTypeKey *string = nil
			if fieldTypeModel != nil {
				fieldTypeKey = &fieldTypeModel.Key
				fieldTypeModel.ParentKey = &typeModel.Key
			}

			fieldDescription := property.Description

			// Create field model
			fieldModel := &FieldModel{
				Name:        fieldName,
				Type:        fieldTypename,
				TypeKey:     fieldTypeKey,
				Description: fieldDescription,
				Required:    containsString(fieldName, schema.Required),
				Schema:      property,
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

func (b *ModelBuilder) createKey(name string) string {
	key := functions.Anchorize(name)
	if value, exists := b.keys[key]; exists {
		value++
		key = fmt.Sprintf("%s-%d", key, value)
	} else {
		b.keys[key] = 0
	}
	return key
}

func getTypeName(props *apiextensions.JSONSchemaProps) string {
	if props.XIntOrString {
		return "int or string"
	}

	if props.XEmbeddedResource {
		return "RawExtension"
	}

	if props.Type == "" && props.XPreserveUnknownFields != nil {
		return "JSON"
	}

	if props.Type == "string" && props.Enum != nil && len(props.Enum) > 0 {
		return "enum"
	}

	if props.Format != "" && props.Type == "byte" {
		return "[]byte"
	}

	// map
	if props.Type == "object" && props.AdditionalProperties != nil {
		if props.AdditionalProperties.Schema == nil && props.AdditionalProperties.Allows {
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

	return props.Type
}

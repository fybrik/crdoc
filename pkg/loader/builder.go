package loader

import (
	"fmt"
	"os"
	"path"
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
	// Links        map[string]string
}

func (b *ModelBuilder) Add(crd *apiextensions.CustomResourceDefinition) error {
	// Add chapter for each version
	for _, version := range crd.Spec.Versions {
		group := crd.Spec.Group
		gv := fmt.Sprintf("%s/%s", group, version.Name)
		kind := crd.Spec.Names.Kind

		log.Info("processing ", gv, kind)

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
			log.Info("adding group ", groupModel)
			b.Model.Groups = append(b.Model.Groups, groupModel)
		}
		if groupModel.Metadata.Title == "" {
			groupModel.Metadata.Title = gv
		}
		if groupModel.Metadata.Description == "" {
			groupModel.Metadata.Description = fmt.Sprintf("API Reference for %s", gv)
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
			log.Info("adding kind ", kindModel)
			groupModel.Kinds = append(groupModel.Kinds, kindModel)
		}
		if kindModel.Metadata.Title == "" {
			kindModel.Metadata.Title = kind
		}
		if kindModel.Metadata.Description == "" {
			kindModel.Metadata.Description = fmt.Sprintf("API Reference for %s (%s)", kind, gv)
		}

		// Find schema
		validation := version.Schema
		if validation == nil {
			// Fallback to resource level schema
			validation = crd.Spec.Validation
		}
		schema := validation.OpenAPIV3Schema

		// Recusively add type models
		b.addTypeModels(kindModel, kind, schema, true)
	}

	return nil
}

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

func (b *ModelBuilder) addTypeModels(kindModel *KindModel, name string, schema *apiextensions.JSONSchemaProps, isTopLevel bool) {
	typeModel := &TypeModel{
		Name:        name,
		Description: schema.Description,
		Link:        fmt.Sprintf("#%s", headingID(name)),
	}
	kindModel.Types = append(kindModel.Types, typeModel)
	log.Info("adding type ", typeModel)

	for _, fieldName := range orderedPropertyKeys(schema.Required, schema.Properties, true) {
		log.Info("adding field ", fieldName)
		property := schema.Properties[fieldName]
		typename, typekey := getTypeNameAndKey(fieldName, property)
		fieldModel := &FieldModel{
			Name:        fieldName,
			Description: property.Description,
			Typename:    typename,
			Typekey:     typekey, // TODO: handle as link
			Required:    isRequired(fieldName, schema.Required),
		}
		typeModel.Fields = append(typeModel.Fields, fieldModel)

		// Recurse if type is array
		if property.Type == "array" {
			b.addTypeModels(kindModel, *typekey, property.Items.Schema, false)
		} else if property.Type == "object" && property.AdditionalProperties != nil {
			b.addTypeModels(kindModel, *typekey, property.AdditionalProperties.Schema, false)
		} else if property.Type == "object" && property.Properties != nil {
			b.addTypeModels(kindModel, *typekey, &property, false)
		}
	}
}

package loader

import (
	"io/ioutil"
	"path"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

func LoadModel(filepath string) (*Model, error) {
	model := &Model{}
	if filepath != "" {
		filecontent, err := ioutil.ReadFile(filepath)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(filecontent, model); err != nil {
			return nil, err
		}
	}
	return model, nil
}

func LoadCRDs(dirpath string) ([]*apiextensions.CustomResourceDefinition, error) {
	files, err := filepath.Glob(path.Join(dirpath, "*"))
	if err != nil {
		return nil, err
	}

	resources := []*apiextensions.CustomResourceDefinition{}
	for _, filepath := range files {
		crd, err := LoadCRD(filepath)
		if err != nil {
			return nil, err
		}
		resources = append(resources, crd)
	}

	return resources, nil
}

func LoadCRD(filepath string) (*apiextensions.CustomResourceDefinition, error) {
	filecontent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return DecodeCRD(filecontent)
}

func DecodeCRD(content []byte) (*apiextensions.CustomResourceDefinition, error) {
	sch := runtime.NewScheme()
	_ = scheme.AddToScheme(sch)
	_ = apiextensions.AddToScheme(sch)
	_ = apiextensionsv1.AddToScheme(sch)
	_ = apiextensionsv1.RegisterConversions(sch)
	_ = apiextensionsv1beta1.AddToScheme(sch)
	_ = apiextensionsv1beta1.RegisterConversions(sch)

	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode
	obj, _, err := decode(content, nil, nil)
	if err != nil {
		return nil, err
	}

	crd := &apiextensions.CustomResourceDefinition{}
	err = sch.Convert(obj, crd, nil)
	if err != nil {
		return nil, err
	}
	return crd, err
}

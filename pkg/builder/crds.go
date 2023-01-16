// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

// LoadCRDs loads all CustomResourceDefinition resources from a directory (glob)
func LoadCRDs(dirpath string) ([]*apiextensions.CustomResourceDefinition, error) {
	files, err := filepath.Glob(path.Join(dirpath, "*"))
	if err != nil {
		return nil, err
	}

	// Glob ignores file system errors, so check the supplied path when there
	// are no results. When it is a file, treat it like a single result from
	// Glob. When it does not exist, return an error.
	if len(files) == 0 {
		info, err := os.Stat(dirpath)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			files = append(files, dirpath)
		}
	}

	resources := []*apiextensions.CustomResourceDefinition{}

	for _, file := range files {
		// Read file
		filecontent, err := os.ReadFile(filepath.Clean(file))
		if err != nil {
			return nil, err
		}

		// Split if multiple YAML documents are defined in the file
		fileDocuments, err := loadYAMLDocuments(filecontent)
		if err != nil {
			return nil, err
		}

		for _, document := range fileDocuments {
			crd, err := decodeCRD(document)
			if err != nil {
				return nil, err
			}
			if crd != nil {
				resources = append(resources, crd)
			}
		}
	}

	return resources, nil
}

func loadYAMLDocuments(filecontent []byte) ([][]byte, error) {
	dec := yaml.NewDecoder(bytes.NewReader(filecontent))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := yaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}

func decodeCRD(content []byte) (*apiextensions.CustomResourceDefinition, error) {
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
		// Not a Kubernetes resource
		if strings.HasPrefix(err.Error(), "Object 'Kind' is missing in ") {
			return nil, nil
		}
		// Not a kind from the above registered schemas
		if strings.HasPrefix(err.Error(), "no kind ") {
			return nil, nil
		}
		return nil, err
	}

	if obj.GetObjectKind().GroupVersionKind().Kind != "CustomResourceDefinition" {
		return nil, nil
	}

	crd := &apiextensions.CustomResourceDefinition{}
	err = sch.Convert(obj, crd, nil)
	if err != nil {
		return nil, err
	}
	return crd, err
}

package builder

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// LoadModel loads a Model from a TOC yaml
func LoadModel(file string) (*Model, error) {
	model := &Model{}
	if file != "" {
		filecontent, err := ioutil.ReadFile(filepath.Clean(file))
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(filecontent, model); err != nil {
			return nil, err
		}
	}
	return model, nil
}

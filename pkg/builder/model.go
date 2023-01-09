// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package builder

import "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

type Model struct {
	Metadata PageMetadata  `yaml:"metadata"`
	Groups   []*GroupModel `yaml:"groups"`
}

type GroupModel struct {
	// Metadata PageMetadata `yaml:"metadata"`
	Group   string       `yaml:"group"`
	Version string       `yaml:"version"`
	Kinds   []*KindModel `yaml:"kinds"`
}

type KindModel struct {
	// Metadata PageMetadata `yaml:"metadata"`
	Name  string `yaml:"name"`
	Key   string
	Types []*TypeModel
}

type PageMetadata struct {
	Title       string `yaml:"title"`
	Weight      int    `yaml:"weight"`
	Description string `yaml:"description"`
}

type TypeModel struct {
	Order       int
	Name        string
	NameConcise string
	Key         string
	Parents     []Parent
	Description string
	IsTopLevel  bool
	Headings    string
	Fields      []*FieldModel
}

type Parent struct {
	Name string
	Key  string
}

type FieldModel struct {
	Name        string
	Type        string
	TypeKey     *string
	Description string
	Required    bool
	Schema      apiextensions.JSONSchemaProps
}

func (m *Model) findGroupModel(group, version string) *GroupModel {
	for _, v := range m.Groups {
		if group == v.Group && version == v.Version {
			return v
		}
	}
	return nil
}

func (m *GroupModel) findKindModel(kind string) *KindModel {
	for _, v := range m.Kinds {
		if kind == v.Name {
			return v
		}
	}
	return nil
}

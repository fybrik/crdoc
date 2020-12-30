package loader

type Model struct {
	Metadata PageMetadata  `yaml:"metadata"`
	Groups   []*GroupModel `yaml:"groups"`
}

type GroupModel struct {
	Metadata PageMetadata `yaml:"metadata"`
	Group    string       `yaml:"group"`
	Version  string       `yaml:"version"`
	Kinds    []*KindModel `yaml:"kinds"`
}

type KindModel struct {
	Metadata PageMetadata `yaml:"metadata"`
	Name     string       `yaml:"name"`

	Types []*TypeModel
}

type PageMetadata struct {
	Title       string `yaml:"title"`
	Weight      int    `yaml:"weight"`
	Description string `yaml:"description"`
}

type TypeModel struct {
	Name        string
	Description string
	Link        string
	Fields      []*FieldModel
}

type FieldModel struct {
	Name        string
	Typename    string
	Typekey     *string
	Description string
	Required    bool
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

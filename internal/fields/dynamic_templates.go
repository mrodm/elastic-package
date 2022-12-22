package fields

import (
	"strings"

	"github.com/pkg/errors"
)

const prefixMapping = "_embedded_ecs"

type dynamicTemplate struct {
	PathMatch string `yaml:"path_match"`
	Match     string `yaml:"match"`
	Mapping   struct {
		Type   string `yaml:"type"`
		Fields struct {
			Text struct {
				Type string `yaml:"type"`
			} `yaml:"text"`
		} `yaml:""fields"`
	} `yaml:"mapping"`
}

func (d *dynamicTemplate) Name() string {
	if d.Match != "" {
		return d.Match
	}
	if d.PathMatch != "" {
		return d.PathMatch
	}

	return ""
}

func (d *dynamicTemplate) Type() string {
	return d.Mapping.Type
}

func (d *dynamicTemplate) FieldDefinition() (*FieldDefinition, error) {
	var field FieldDefinition
	field.Name = d.Name()
	field.Type = d.Type()

	return &field, nil
}

type templates struct {
	Elasticsearch struct {
		IndexTemplate struct {
			Mappings struct {
				DynamicTemplates []map[string]dynamicTemplate `yaml:"dynamic_templates"`
			} `yaml:"mappings"`
		} `yaml:"index_template"`
	} `yaml:"elasticsearch"`
}

func (t *templates) FieldDefinitions() ([]FieldDefinition, error) {
	var fields []FieldDefinition
	for _, template := range t.Elasticsearch.IndexTemplate.Mappings.DynamicTemplates {
		if len(template) > 1 {
			return nil, errors.New("template with more than one definition")
		}
		for name, v := range template {
			field, err := v.FieldDefinition()

			if strings.HasPrefix(name, prefixMapping) {
				field.External = "ecs"
			}

			if err != nil {
				return nil, err
			}
			fields = append(fields, *field)
		}
	}
	return fields, nil
}

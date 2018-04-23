package formatter

import (
	"bytes"
	"github.com/twhiston/factd/lib/common"
	"gopkg.in/yaml.v2"
)

// YAMLFormatter prints-out facts in YAML format
type YAMLFormatter struct {
}

// NewYAMLFormatter returns new YAML formatter
func NewYAMLFormatter() *YAMLFormatter {
	return &YAMLFormatter{}
}

// Format marshals a map of FactList's keyed by string to YAML
func (jf *YAMLFormatter) Format(facts map[string]common.FactList) (*bytes.Buffer, error) {
	b, err := yaml.Marshal(facts)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

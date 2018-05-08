package formatter

import (
	"bytes"
	"github.com/twhiston/factd/lib/common"
	"gopkg.in/yaml.v2"
)

// YAMLFormatter prints-out facts in YAML format
type YAMLFormatter struct {
}

// Name returns the formatter name, in a format suitable for using as a map key
func (f *YAMLFormatter) Name() string {
	return GetFormatterName(f)
}

// Format marshals a map of FactList's keyed by string to YAML
func (f *YAMLFormatter) Format(facts map[string]common.FactList) (*bytes.Buffer, error) {
	b, err := yaml.Marshal(facts)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

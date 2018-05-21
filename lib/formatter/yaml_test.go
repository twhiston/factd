package formatter

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestYAMLFormatter_Name(t *testing.T) {
	runTestFormatterName(&YAMLFormatter{}, "yaml", t)
}

func TestYAMLFormatter_Format(t *testing.T) {
	f := YAMLFormatter{}
	facts := getExampleFactList()
	formatted, err := f.Format(facts)
	if err != nil {
		t.Error(err)
	}
	if formatted.Len() == 0 {
		t.Error("formatted length should not be nil")
	}
	if !isYAML(formatted.Bytes()) {
		t.Error("could not marshall output back to YAML")
	}
}

func isYAML(b []byte) bool {
	var y interface{}
	return yaml.Unmarshal(b, &y) == nil
}

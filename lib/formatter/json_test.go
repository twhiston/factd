package formatter

import (
	"encoding/json"
	"testing"
)

func TestJSONFormatter_Name(t *testing.T) {
	runTestFormatterName(&JSONFormatter{}, "json", t)
}

func TestJSONFormatter_Format(t *testing.T) {
	f := JSONFormatter{}
	facts := getExampleFactList()
	formatted, err := f.Format(facts)
	if err != nil {
		t.Error(err)
	}
	if formatted.Len() == 0 {
		t.Error("formatted length should not be nil")
	}
	if !isJSON(formatted.Bytes()) {
		t.Error("could not marshall output back to JSON")
	}
}

func isJSON(b []byte) bool {
	var js interface{}
	return json.Unmarshal(b, &js) == nil
}

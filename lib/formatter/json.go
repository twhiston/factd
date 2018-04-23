package formatter

import (
	"bytes"
	j "encoding/json"
	"github.com/twhiston/factd/lib/common"
)

// JSONFormatter prints-out facts in JSON format
type JSONFormatter struct {
}

// NewJSONFormatter returns new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format prints-out facts in JSON format
func (jf *JSONFormatter) Format(facts map[string]common.FactList) (*bytes.Buffer, error) {
	b, err := j.MarshalIndent(facts, "", "  ")
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

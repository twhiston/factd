package formatter

import (
	"bytes"
	j "encoding/json"
	"github.com/twhiston/factd/lib/common"
)

// JSONFormatter prints-out facts in JSON format
type JSONFormatter struct {
}

// Name returns the formatter name, in a format suitable for using as a map key
func (f *JSONFormatter) Name() string {
	return GetFormatterName(f)
}

// Format prints-out facts in JSON format
func (f *JSONFormatter) Format(facts map[string]common.FactList) (*bytes.Buffer, error) {
	b, err := j.MarshalIndent(facts, "", "  ")
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

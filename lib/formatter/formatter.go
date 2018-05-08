package formatter

import (
	"bytes"

	"github.com/twhiston/factd/lib/common"
	"reflect"
	"strings"
)

// Formatter interface
type Formatter interface {
	Format(map[string]common.FactList) (*bytes.Buffer, error)
	Name() string
}

// GetFormatterName is a helper function to return the formatter name nicely formatted
func GetFormatterName(myvar interface{}) string {
	name := reflect.TypeOf(myvar).String()
	parts := strings.Split(name, ".")
	noSuffix := strings.TrimSuffix(parts[len(parts)-1], "Formatter")
	return strings.ToLower(noSuffix)
}

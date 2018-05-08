package formatter

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/common/logging"
	"io"
)

// PlainTextFormatter prints-out facts in k=>v format
type PlainTextFormatter struct {
}

// Name returns the formatter name, in a format suitable for using as a map key
func (f *PlainTextFormatter) Name() string {
	return GetFormatterName(f)
}

// Format prints-out facts in k=>v format
func (f *PlainTextFormatter) Format(facts map[string]common.FactList) (*bytes.Buffer, error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	for k, v := range facts {
		fmt.Fprintf(writer, "%v\n", k)
		for fn, f := range v {
			indentPrint(f, fn, 1, " => ", writer)
		}
	}
	logging.Fatal(writer.Flush())
	return &b, nil
}

// TODO - gross :P
func indentPrint(data interface{}, name string, amount int, divider string, writer io.Writer) {
	switch val := data.(type) {
	case common.FactList:
		fmt.Fprintf(writer, "%v%v\n", indent(name, amount), divider)
		amount++
		for k, v := range val {
			indentPrint(v, k, amount, divider, writer)
		}
	case map[string]interface{}:
		fmt.Fprintf(writer, "%v%v\n", indent(name, amount), divider)
		amount++
		for k, v := range val {
			indentPrint(v, k, amount, divider, writer)
		}
	case map[string]string:
		fmt.Fprintf(writer, "%v%v\n", indent(name, amount), divider)
		amount++
		for k, v := range val {
			indentPrint(v, k, amount, divider, writer)
		}
	default:
		fmt.Fprintf(writer, "%v%v%v\n", indent(name, amount), divider, val)
	}
}

// indent inserts prefix at the beginning of each non-empty line of s. The
// end-of-line marker is NL.
func indent(s string, count int) string {
	for i := count; i > 0; i-- {
		s = "	" + s
	}
	return s
}

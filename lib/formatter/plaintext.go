package formatter

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/common/logging"
	"io"
	"reflect"
	"strconv"
)

// PlainTextFormatter prints-out facts in k=>v format
type PlainTextFormatter struct {
	Divider string
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
		for fn, fd := range v {
			f.indentPrint(fd, fn, 1, writer)
		}
	}
	logging.Fatal(writer.Flush())
	return &b, nil
}

func (f *PlainTextFormatter) indentPrint(data interface{}, name string, amount int, writer io.Writer) {
	rt := reflect.TypeOf(data)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		fmt.Fprintf(writer, "%v%v\n", indent(name, amount), f.Divider)
		amount++
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			val := s.Index(i).Interface()
			f.indentPrint(val, strconv.Itoa(i), amount, writer)
		}
	case reflect.Map:
		fmt.Fprintf(writer, "%v%v\n", indent(name, amount), f.Divider)
		amount++
		s := reflect.ValueOf(data)
		keys := s.MapKeys()
		for i := 0; i < s.Len(); i++ {
			val := s.MapIndex(keys[i]).Interface()
			f.indentPrint(val, keys[i].String(), amount, writer)
		}
	case reflect.Struct:
		s := reflect.ValueOf(data)
		for i := 0; i < s.NumField(); i++ {
			val := s.Field(i).Interface()
			f.indentPrint(val, s.Type().Field(i).Name, amount, writer)
		}
	default:
		fmt.Fprintf(writer, "%v%v%v\n", indent(name, amount), f.Divider, data)
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

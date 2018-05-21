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

// FlattenedFormatter prints-out facts in flattened format
type FlattenedFormatter struct {
	Divider   string
	KvDivider string
	Prepend   string
}

// Name returns the formatter name, in a format suitable for using as a map key
func (ff *FlattenedFormatter) Name() string {
	return GetFormatterName(ff)
}

// Format prints-out facts in k=>v format
func (ff *FlattenedFormatter) Format(facts map[string]common.FactList) (*bytes.Buffer, error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	for k, v := range facts {
		for fn, f := range v {
			key := ff.Prepend + k + ff.Divider + fn
			ff.flatPrint(f, key, writer)
		}
	}
	logging.Fatal(writer.Flush())
	return &b, nil
}

func (ff *FlattenedFormatter) flatPrint(data interface{}, name string, writer io.Writer) {
	rt := reflect.TypeOf(data)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			val := s.Index(i).Interface()
			key := name + ff.Divider + strconv.Itoa(i)
			ff.flatPrint(val, key, writer)
		}
	case reflect.Map:
		s := reflect.ValueOf(data)
		keys := s.MapKeys()
		for i := 0; i < s.Len(); i++ {
			val := s.MapIndex(keys[i]).Interface()
			key := name + ff.Divider + keys[i].String()
			ff.flatPrint(val, key, writer)
		}
	case reflect.Struct:
		s := reflect.ValueOf(data)
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i).Interface()
			key := name + ff.Divider + s.Type().Field(i).Name
			ff.flatPrint(f, key, writer)
		}
	default:
		fmt.Fprintf(writer, "%v%v%v\n", name, ff.KvDivider, data)
	}
}

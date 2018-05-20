package formatter

import (
	"strings"
	"testing"
)

func TestFlattenedFormatter_Name(t *testing.T) {
	runTestFormatterName(&FlattenedFormatter{}, "flattened", t)
}

func TestFlattenedFormatter_Format(t *testing.T) {

	f := FlattenedFormatter{"/", " ", ""}

	facts := getExampleFactList()

	bytesBuf, err := f.Format(facts)
	if err != nil {
		t.Error(err)
	}
	if bytesBuf.Len() == 0 {
		t.Error("buffer should not be empty")
	}
	str := bytesBuf.String()
	testVal(str, "basic/text simple", t)
	testVal(str, "slice/data/2 c", t)
	testVal(str, "factlist/nested/inner map", t)
	testVal(str, "factlist/slice/1 data", t)
	testVal(str, "factlist/struct/Map/data1/2 c", t)
	testVal(str, "factlist/struct/Factlist/fact1 x", t)

}

func testVal(actual, expected string, t *testing.T) {
	if !strings.Contains(actual, expected) {
		t.Error(actual + " did not match expected: " + expected)
	}
}

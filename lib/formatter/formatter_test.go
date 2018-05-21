package formatter

import (
	"github.com/twhiston/factd/lib/common"
	"testing"
)

type testStructVals struct {
	Text1    string
	Factlist common.FactList
	Map      map[string][]string
}

func getExampleFactList() map[string]common.FactList {
	ns := testStructVals{
		Text1:    "structValue",
		Factlist: common.FactList{"fact1": "x", "fact2": "y"},
		Map:      map[string][]string{"data1": {"a", "b", "c", "d"}},
	}
	facts := map[string]common.FactList{}
	facts["basic"] = common.FactList{"text": "simple"}
	facts["slice"] = common.FactList{"data": []string{"a", "b", "c", "d"}}
	facts["map"] = common.FactList{"stringmap": map[string]string{"this": "is", "a": "test"}}
	nested := map[string]string{"inner": "map", "has": "values"}
	facts["factlist"] = common.FactList{"test": "simple", "nested": nested, "slice": []string{"more", "data"}, "struct": ns}
	return facts
}

func runTestFormatterName(f Formatter, expectedName string, t *testing.T) {
	if f.Name() != expectedName {
		t.Error("incorrect name: " + f.Name() + " expected: " + expectedName)
	}
}

type TestFormatter struct{}
type TestDifferentName struct{}

func TestGetFormatterName(t *testing.T) {

	n := GetFormatterName(TestFormatter{})
	if n != "test" {
		t.Error("TestFormatter name should return as test. Actual: " + n)
	}

	n = GetFormatterName(TestDifferentName{})
	if n != "testdifferentname" {
		t.Error("TestDifferentName name should return as testdifferentname. Actual: " + n)
	}

}

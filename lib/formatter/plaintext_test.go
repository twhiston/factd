package formatter

import "testing"

func TestPlainTextFormatter_Name(t *testing.T) {
	runTestFormatterName(&PlainTextFormatter{}, "plaintext", t)
}

func TestPlainTextFormatter_Format(t *testing.T) {

	f := PlainTextFormatter{Divider: " => "}

	facts := getExampleFactList()

	bytesBuf, err := f.Format(facts)
	if err != nil {
		t.Error(err)
	}
	if bytesBuf.Len() == 0 {
		t.Error("buffer should not be empty")
	}
	str := bytesBuf.String()
	t.Log(str)
	testVal(str, "basic\n	text => simple", t)
	testVal(str, "slice\n	data => \n		0 => a", t)
	//testVal(str, "factlist\n	test => simple", t)
	//testVal(str, "factlist/slice/1 data", t)
	//testVal(str, "factlist/struct/Map/data1/2 c", t)
	//testVal(str, "factlist/struct/Factlist/fact1 x", t)

}

package factd

import (
	"fmt"
	"strings"
	"testing"

	"bytes"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/formatter"
)

func TestFacter(t *testing.T) {
	testKey := "test"
	testValue := "test value"

	f := New(Config{})
	if f == nil {
		t.Fail()
	}
	f.AddFact("test", testKey, testValue)
	value, ok := f.Get("test", testKey)
	if !ok || strings.Compare(fmt.Sprintf("%v", value), testValue) != 0 {
		t.Fatalf("Failed to get K/V: %v:%v:%v", testKey, value, ok)
	}
	f.Delete(testKey)
	value, ok = f.Get("test", testKey)
	if ok {
		t.Fatalf("Got %v, value %v", ok, value)
	}
}

type FakeFormatter struct {
	facts map[string]common.FactList
}

func (ff *FakeFormatter) Get(group string, k string) (interface{}, bool) {
	val, ok := ff.facts[group][k]
	return val, ok
}

func NewFakeFormatter() *FakeFormatter {
	f := FakeFormatter{}
	f.facts = make(map[string]common.FactList)
	return &f
}

func (ff *FakeFormatter) Format(facts map[string]common.FactList) (*bytes.Buffer, error) {
	for k, v := range facts {
		if ff.facts[k] == nil {
			ff.facts[k] = make(map[string]interface{})
		}
		ff.facts[k] = v
	}
	var buf []byte
	return bytes.NewBuffer(buf), nil
}

func TestFacterFormatter(t *testing.T) {
	testGroup := "test"
	testKey := "test"
	testValue := "test_value"
	ff := NewFakeFormatter()
	conf := Config{
		Formatter: ff,
	}
	f := New(conf)
	if f == nil {
		t.Fatal()
	}
	f.AddFact(testGroup, testKey, testValue)
	f.Format()
	val, ok := ff.Get(testGroup, testKey)
	if !ok {
		t.Fatal()
	}
	if strings.Compare(fmt.Sprintf("%v", val), testValue) != 0 {
		t.Fatal()
	}
}

func TestNewNilConf(t *testing.T) {
	f := New(Config{})
	if f == nil {
		t.Fail()
	}
}

func TestNewConf(t *testing.T) {
	conf := Config{
		Formatter: formatter.NewPlainTextFormatter(),
	}
	f := New(conf)
	if f == nil {
		t.Fail()
	}
}

package env

import (
	"log"
	"reflect"
	"testing"
)

type testStruct struct {
	A int64  `env:"test,int"`
	B string `env:"test2"`
	C string `env:"test3,string"`
}

var TestEnv = []string{
	"test=test",
	"test2=te=st",
	"test3",
}

var TestStruct = testStruct{
	A: 1,
	B: "test",
	C: "test",
}

var testEnv = []string{
	"test=1",
	"test2=stoff",
	"test3=",
}

func TestSplitEnv(t *testing.T) {

	a, b, err := splitEnv(TestEnv[0])
	if err != nil {
		t.Errorf("splitEnv(%s) = %s, %s, %s, wanted: test, test, \"\"", TestEnv[0], a, b, err.Error())
	}

	if a != "test" || b != "test" {
		t.Errorf("splitEnv(%s) = %s, %s, wanted: test, test", TestEnv[0], a, b)
	}

	a, b, err = splitEnv(TestEnv[1])
	if err != nil {
		t.Errorf("splitEnv(%s) = %s, %s, %s, wanted: test2, te=st, \"\"", TestEnv[1], a, b, err.Error())
	}

	if a != "test2" || b != "te=st" {
		t.Errorf("splitEnv(%s) = %s, %s, wanted: test2, te=st", TestEnv[1], a, b)
	}

	a, b, err = splitEnv(TestEnv[2])
	if err == nil {
		t.Errorf("splitEnv(%s) should fail but doesn't", TestEnv[2])
	}
}

func TestParse(t *testing.T) {
	a, err := parse(TestEnv[:2])
	if err != nil {
		t.Errorf("parse(%+v) returned error: %s", TestEnv[:2], err.Error())
	}

	if a["test"] != "test" {
		t.Errorf("parse(%+v)[test] = %s, wanted: 'test'", TestEnv[:2], a["test"])
	} else if a["test2"] != "te=st" {
		t.Errorf("parse(%+v)[test2] = %s, wanted: 'te=st'", TestEnv[:2], a["test2"])
	}

	a, err = parse(TestEnv)
	if err == nil {
		t.Errorf("parse(%+v) should return error but didn't", TestEnv)
	}

}

func TestParseTag(t *testing.T) {
	x := reflect.TypeOf(TestStruct)
	f, _ := x.FieldByName("A")
	tn, _ := f.Tag.Lookup(tagName)
	name, opt := parseTag(tn)
	if name != "test" || opt[0] != "int" {
		t.Errorf("parseTag(TestStruct.a) = %s, %s. Wanted: %s, %s", name, opt, "test", "int")
	}

	f, _ = x.FieldByName("B")
	tn, _ = f.Tag.Lookup(tagName)
	name, opt = parseTag(tn)
	if name != "test2" || !reflect.DeepEqual(opt, []string{}) {
		t.Errorf("parseTag(TestStruct.b) = %s, %s. Wanted: %s, %s", name, opt, "test2", "[]")
	}

	f, _ = x.FieldByName("C")
	tn, _ = f.Tag.Lookup(tagName)
	name, opt = parseTag(tn)
	if name != "test3" || opt[0] != "string" {
		t.Errorf("parseTag(TestStruct.c) = %s, %s. Wanted: %s, %s", name, opt, "test3", "string")
	}
}

func TestUnmarshal(t *testing.T) {
	ts := &testStruct{}

	err := Unmarshal(testEnv, ts)
	if err != nil {
		t.Errorf("Unmarshal(testEnv, testStruct) failed with err: %s", err.Error())
	}
	if ts.A != 1 {
		t.Errorf("unmarshal(testEnv, testStruct) A = %d, wanted: %d", ts.A, 1)
	}
	log.Printf("%+v", ts)
}
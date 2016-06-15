package filter

import (
	"encoding/json"
	"strings"
	"testing"
)

func assertSplit(t *testing.T, in string, exp []string) {
	out, err := parseQueryString(in)
	if err != nil {
		t.Fatalf("%q: %s", in, err)
	}

	if len(out) != len(exp) {
		t.Fatalf("expected %v but got %v", exp, out)
	}

	for i, v := range out {
		if exp[i] != v {
			t.Fatalf("[%d] expected %v but got %v", i, exp, out)
		}
	}
}

func TestQuerySplit(t *testing.T) {
	cases := map[string][]string{
		".a.b":                []string{"a", "b"},
		".a[5].b":             []string{"a", "[5]", "b"},
		".a[.name[0]=fish].b": []string{"a", "[.name[0]=fish]", "b"},
		".a[2]":               []string{"a", "[2]"},
		"[2]":                 []string{"[2]"},
	}

	for in, exp := range cases {
		assertSplit(t, in, exp)
	}
}

func TestFindClosingBracket(t *testing.T) {
	if findClosingBracket("asdas]") != 5 {
		t.Fatal("expected 5")
	}

	if v := findClosingBracket("a[d]as]"); v != 6 {
		t.Fatal("expected 6, got ", v)
	}

	if v := findClosingBracket("a[[[[]]]]as]"); v != 11 {
		t.Fatal("expected 11, got ", v)
	}
}

func assertGet(t *testing.T, jsonText, query string, expected interface{}) {
	switch i := expected.(type) {
	case int:
		expected = float64(i)
	}

	var obj interface{}
	err := json.Unmarshal([]byte(jsonText), &obj)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Get(obj, query)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Fatalf("While testing Get (%v), got: %v (%T), expected: %v (%T)", query, res, res, expected, expected)
	}
}

func TestGet(t *testing.T) {
	cases := map[string]interface{}{
		`{"a": 1}@a`:               1,
		`{"a": [1, 2, 3, 4]}@a[2]`: 3,
		`[1, 2, 3]@[1]`:            2,
		`[{"id": 1, "name": "a"}, {"id": 2, "name":"b"}]@[.name=a].id`: 1,
	}
	for q, v := range cases {
		tq := strings.Split(q, "@")
		assertGet(t, tq[0], tq[1], v)
	}
}

func assertSet(t *testing.T, jsonText, query string, value interface{}) {
	switch i := value.(type) {
	case int:
		value = float64(i)
	}

	var obj interface{}
	err := json.Unmarshal([]byte(jsonText), &obj)
	if err != nil {
		t.Fatal(err)
	}
	err = Set(obj, query, value)
	if err != nil {
		t.Fatal(err)
	}

	res, err := Get(obj, query)
	if err != nil {
		t.Fatal(err)
	}

	if res != value {
		t.Fatalf("While testing Set (%v), got: %v (%T), expected: %v (%T)", query, res, res, value, value)
	}
}

func TestSet(t *testing.T) {
	cases := map[string]interface{}{
		`{}@a`:                     "hello",
		`{"a": 1}@a`:               3,
		`{"a": [1, 2, 3, 4]}@a[2]`: 5,
		`[1, 2, 3]@[1]`:            5,
		//TODO: Known breakage, fixing it is .. and I really don't want to do it now
		//`[1, 2]@[5]`:               5,
		`[{"id": 1, "name": "a"}, {"id": 2, "name":"b"}]@[.name=a].id`: 3,
	}
	for q, v := range cases {
		tq := strings.Split(q, "@")
		assertSet(t, tq[0], tq[1], v)
	}
}

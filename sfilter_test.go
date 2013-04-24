package sfilter_test

import (
	"reflect"
	"testing"

	"git.cupcake.io/picard/sfilter"
)

// Hook up gocheck into the gotest runner.
type NestedExample struct {
	A string `sfilter:"one,two"`
	B string `sfilter:"two"`
}

var example = struct {
	None   string
	One    string `sfilter:"one"`
	Two    string `sfilter:"two"`
	Tagged string `json:"tagged,omitempty" sfilter:"one,two"`

	NestedOne NestedExample  `sfilter:"one"`
	NestedTwo *NestedExample `sfilter:"two"`

	SliceOne []NestedExample  `sfilter:"one"`
	SliceTwo []*NestedExample `sfilter:"two"`
}{
	None:   "a",
	One:    "1",
	Two:    "2",
	Tagged: "t",

	NestedOne: NestedExample{A: "a", B: "b"},
	NestedTwo: &NestedExample{A: "b", B: "a"},

	SliceOne: []NestedExample{{A: "ab", B: "ba"}, {A: "b", B: "a"}},
	SliceTwo: []*NestedExample{{A: "ab", B: "ba"}, {A: "b", B: "a"}},
}

var mapTests = []struct {
	tags     []string
	expected map[string]interface{}
}{
	{[]string{"one"}, map[string]interface{}{
		"One":    "1",
		"tagged": "t",

		"NestedOne": map[string]interface{}{"A": "a"},
		"SliceOne":  []map[string]interface{}{{"A": "ab"}, {"A": "b"}},
	}},
	{[]string{"two"}, map[string]interface{}{
		"Two":    "2",
		"tagged": "t",

		"NestedTwo": map[string]interface{}{"A": "b", "B": "a"},
		"SliceTwo":  []map[string]interface{}{{"A": "ab", "B": "ba"}, {"A": "b", "B": "a"}},
	}},
	{[]string{"one", "two"}, map[string]interface{}{
		"One":    "1",
		"Two":    "2",
		"tagged": "t",

		"NestedOne": map[string]interface{}{"A": "a", "B": "b"},
		"NestedTwo": map[string]interface{}{"A": "b", "B": "a"},
		"SliceOne":  []map[string]interface{}{{"A": "ab", "B": "ba"}, {"A": "b", "B": "a"}},
		"SliceTwo":  []map[string]interface{}{{"A": "ab", "B": "ba"}, {"A": "b", "B": "a"}},
	}},
}

func TestMap(t *testing.T) {
	for i, test := range mapTests {
		actual, err := sfilter.Map(example, test.tags...)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("got %#v, want %#v, test %d: %#v", actual, test.expected, i, test.tags)
		}
	}
}

type Emptiness struct {
	Foo string `json:"foo,omitempty" sfilter:"a"`
	Bar Bar    `json:"bar,omitempty" sfilter:"a"`
}

type Bar struct {
	Baz string `json:"baz,omitempty" sfilter:"a"`
}

var omitemptyTests = []struct {
	example  Emptiness
	expected map[string]interface{}
}{
	{Emptiness{"", Bar{""}}, map[string]interface{}{}},
	{Emptiness{"a", Bar{""}}, map[string]interface{}{"foo": "a"}},
}

func TestOmitEmpty(t *testing.T) {
	for i, test := range omitemptyTests {
		actual, err := sfilter.Map(test.example, "a")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("got %#v, want %#v, test %d", actual, test.expected, i)
		}
	}
}

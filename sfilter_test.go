package sfilter_test

import (
	"testing"

	"git.cupcake.io/picard/sfilter"
	. "launchpad.net/gocheck"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

type NestedExample struct {
	A string `sfilter:"one,two"`
	B string `sfilter:"two"`
}

var example = struct {
	All    string
	One    string `sfilter:"one"`
	Two    string `sfilter:"two"`
	Tagged string `json:"tagged,omitempty"`

	NestedOne NestedExample  `sfilter:"one"`
	NestedTwo *NestedExample `sfilter:"two"`

	SliceOne []NestedExample  `sfilter:"one"`
	SliceTwo []*NestedExample `sfilter:"two"`

	Other []int
}{
	All:    "a",
	One:    "1",
	Two:    "2",
	Tagged: "t",

	NestedOne: NestedExample{A: "a", B: "b"},
	NestedTwo: &NestedExample{A: "b", B: "a"},

	SliceOne: []NestedExample{{A: "ab", B: "ba"}, {A: "b", B: "a"}},
	SliceTwo: []*NestedExample{{A: "ab", B: "ba"}, {A: "b", B: "a"}},

	Other: []int{1, 2, 3},
}

var mapTests = []struct {
	tags     []string
	expected map[string]interface{}
}{
	{[]string{"one"}, map[string]interface{}{
		"All":    "a",
		"One":    "1",
		"tagged": "t",

		"NestedOne": map[string]interface{}{"A": "a"},
		"SliceOne":  []map[string]interface{}{{"A": "ab"}, {"A": "b"}},

		"Other": []int{1, 2, 3},
	}},
	{[]string{"two"}, map[string]interface{}{
		"All":    "a",
		"Two":    "2",
		"tagged": "t",

		"NestedTwo": map[string]interface{}{"A": "b", "B": "a"},
		"SliceTwo":  []map[string]interface{}{{"A": "ab", "B": "ba"}, {"A": "b", "B": "a"}},

		"Other": []int{1, 2, 3},
	}},
	{[]string{"one", "two"}, map[string]interface{}{
		"All":    "a",
		"One":    "1",
		"Two":    "2",
		"tagged": "t",

		"NestedOne": map[string]interface{}{"A": "a", "B": "b"},
		"NestedTwo": map[string]interface{}{"A": "b", "B": "a"},
		"SliceOne":  []map[string]interface{}{{"A": "ab", "B": "ba"}, {"A": "b", "B": "a"}},
		"SliceTwo":  []map[string]interface{}{{"A": "ab", "B": "ba"}, {"A": "b", "B": "a"}},

		"Other": []int{1, 2, 3},
	}},
}

func (s *TestSuite) TestMap(c *C) {
	for _, t := range mapTests {

		actual, err := sfilter.Map(example, t.tags...)
		if err != nil {
			c.Log(err)
			c.Fail()
		}
		c.Assert(actual, DeepEquals, t.expected, Commentf("%s", t.tags))
	}
}

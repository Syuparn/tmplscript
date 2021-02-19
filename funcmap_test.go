package main

import (
	"fmt"
	"testing"
	"text/template"
)

func TestSearchFunc(t *testing.T) {
	tests := []struct {
		title    string
		funcMap  template.FuncMap
		key      string
		expected []string
	}{
		{
			"not found",
			map[string]interface{}{},
			"",
			[]string{},
		},
		{
			"empty key gets all functions",
			map[string]interface{}{
				"a": func() {},
				"b": func() {},
			},
			"",
			[]string{"a", "b"},
		},
		{
			"key gets all functions whose key starts with it",
			map[string]interface{}{
				"fail":   func() {},
				"find":   func() {},
				"finish": func() {},
				"get":    func() {},
			},
			"fi",
			[]string{"find", "finish"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf(tt.title), func(t *testing.T) {
			actual := searchFunc(tt.funcMap)(tt.key)

			if len(actual) != len(tt.expected) {
				t.Fatalf("wrong length. got %d, expected %d", len(actual), len(tt.expected))
			}

			for i, elem := range actual {
				if elem != tt.expected[i] {
					t.Errorf("actual[%d] is wrong. got %s, expected %s",
						i, elem, tt.expected[i])
				}
			}
		})
	}
}

func TestDocFunc(t *testing.T) {
	tests := []struct {
		title    string
		funcMap  template.FuncMap
		key      string
		expected string
	}{
		{
			"arity 0, return 0",
			map[string]interface{}{
				"fail": func() {},
			},
			"fail",
			"fail -> ()",
		},
		{
			"arity 1, return 1",
			map[string]interface{}{
				"itoa": func(i int) string { return fmt.Sprint(i) },
			},
			"itoa",
			"itoa int -> (string)",
		},
		{
			"arity 3, return 1",
			map[string]interface{}{
				"myTernary": func(i int, s string, f float64) bool { return false },
			},
			"myTernary",
			"myTernary int string float64 -> (bool)",
		},
		{
			"arity 1, return 3",
			map[string]interface{}{
				"myFunc": func(i int) (string, bool, error) { return "", true, nil },
			},
			"myFunc",
			"myFunc int -> (string bool error)",
		},
		{
			"variadic",
			map[string]interface{}{
				"myFunc": func(i int, strs ...string) {},
			},
			"myFunc",
			"myFunc int ...string -> ()",
		},
		{
			"not found",
			map[string]interface{}{},
			"myFunc",
			"function myFunc is not defined (or embedded)",
		},
		{
			"not function",
			map[string]interface{}{
				"val": 1,
			},
			"val",
			"val is not a function",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf(tt.title), func(t *testing.T) {
			actual := docFunc(tt.funcMap)(tt.key)

			if actual != tt.expected {
				t.Errorf("wrong output. expected `%s`, got `%s`",
					tt.expected, actual)
			}
		})
	}
}

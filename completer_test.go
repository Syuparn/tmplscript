package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCompleter(t *testing.T) {
	c := newWordCompleter()

	type completion struct {
		head        string
		completions []string
		tail        string
	}

	tests := []struct {
		title    string
		line     string
		pos      int
		expected completion
	}{
		{
			"empty",
			"",
			0,
			completion{
				head:        "",
				completions: []string{},
				tail:        "",
			},
		},
		{
			"only partials",
			"low",
			3,
			completion{
				head:        "",
				completions: []string{"lower"},
				tail:        "",
			},
		},
		{
			"cursor is not at the end",
			"low",
			1,
			completion{
				head:        "",
				completions: []string{"lower"},
				tail:        "",
			},
		},
		{
			"two tokens",
			"if low",
			6,
			completion{
				head:        "if ",
				completions: []string{"lower"},
				tail:        "",
			},
		},
		{
			"two tokens and cursor in the middle",
			"if low",
			3,
			completion{
				head:        "if ",
				completions: []string{"lower"},
				tail:        "",
			},
		},
		{
			"cursor points non-function",
			"if low",
			0,
			completion{
				head:        "",
				completions: []string{},
				tail:        " low",
			},
		},
		{
			"many tokens (cursor at the end)",
			"1 | add int6",
			12,
			completion{
				head:        "1 | add ",
				completions: []string{"int64"},
				tail:        "",
			},
		},
		{
			"many tokens (cursor in the middle)",
			"1 | ad int64",
			4,
			completion{
				head:        "1 | ",
				completions: []string{"add", "add1f", "add1", "addf", "adler32sum"},
				tail:        " int64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			head, completions, tail := c(tt.line, tt.pos)

			assert.Equal(t, tt.expected.head, head, "wrong head")
			assert.ElementsMatch(t, tt.expected.completions, completions, "wrong completions")
			assert.Equal(t, tt.expected.tail, tail, "wrong tail")
		})
	}
}

func allFuncNames() []string {
	names := []string{}
	for name := range funcMap() {
		names = append(names, name)
	}
	return names
}

package algorithm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	type input struct {
		text, pattern string
	}
	tests := map[string]struct {
		input
		want []int
	}{
		"pattern and text are zero length": {
			input: input{
				text:    "",
				pattern: "",
			},
			want: []int{},
		},
		"pattern longer than text": {
			input: input{
				text:    "abc",
				pattern: "abcd",
			},
			want: []int{},
		},
		"pattern not in text": {
			input: input{
				text:    "abc",
				pattern: "dd",
			},
			want: []int{},
		},
		"pattern at the start of text": {
			input: input{
				text:    "abc",
				pattern: "ab",
			},
			want: []int{0},
		},
		"pattern at the end of text": {
			input: input{
				text:    "abc",
				pattern: "bc",
			},
			want: []int{1},
		},
		"pattern in the middle of text": {
			input: input{
				text:    "abcdfg",
				pattern: "cd",
			},
			want: []int{2},
		},
		"pattern doesn't overlap in text": {
			input: input{
				text:    "aaabaaa",
				pattern: "aaa",
			},
			want: []int{0, 4},
		},
		"pattern overlaps in text": {
			input: input{
				text:    "aaaaa",
				pattern: "aa",
			},
			want: []int{0, 1, 2, 3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Search(tt.text, tt.pattern)
			assert.Equal(t, tt.want, got)
		})
	}
}

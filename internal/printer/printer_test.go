package printer

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccumulatePatterns(t *testing.T) {
	type input struct {
		line     string
		patterns []string
	}

	tests := map[string]struct {
		input
		want []pattern
	}{
		"no patterns": {
			input: input{
				line:     "abcd",
				patterns: []string{},
			},
			want: []pattern{},
		},
		"no text": {
			input: input{
				line:     "",
				patterns: []string{"ab", "dcf"},
			},
			want: []pattern{},
		},
		"pattern is zero length": {
			input: input{
				line:     "abcd",
				patterns: []string{""},
			},
			want: []pattern{},
		},
		"pattern are longer than text": {
			input: input{
				line:     "abcd",
				patterns: []string{"abcdefg"},
			},
			want: []pattern{},
		},
		"patterns are longer than text": {
			input: input{
				line:     "abcd",
				patterns: []string{"abcdefg", "cdefghij"},
			},
			want: []pattern{},
		},
		"pattern is in text": {
			input: input{
				line:     "abcdefg",
				patterns: []string{"abcd"},
			},
			want: []pattern{
				{0, 3},
			},
		},
		"patterns are in text": {
			input: input{
				line:     "abcdefg",
				patterns: []string{"abcd", "def", "fg"},
			},
			want: []pattern{
				{0, 3},
				{3, 5},
				{5, 6},
			},
		},
		"patters duplicate": {
			input: input{
				line:     "abcdefg",
				patterns: []string{"abcd", "abcd", "fg"},
			},
			want: []pattern{
				{0, 3},
				{0, 3},
				{5, 6},
			},
		},
		"patterns overlap": {
			input: input{
				line:     "abcdefg",
				patterns: []string{"abcd", "bcd", "defg"},
			},
			want: []pattern{
				{0, 3},
				{1, 3},
				{3, 6},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := accumulatePatterns(tt.line, tt.patterns)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergePatternRanges(t *testing.T) {
	tests := map[string]struct {
		input []pattern
		want  []pattern
	}{
		"no patters": {
			input: []pattern{},
			want:  []pattern{},
		},
		"single pattern": {
			input: []pattern{{0, 3}},
			want:  []pattern{{0, 3}},
		},
		"not overlapping patterns": {
			input: []pattern{{1, 3}, {5, 8}},
			want:  []pattern{{1, 3}, {5, 8}},
		},
		"overlapping patterns": {
			input: []pattern{{1, 3}, {2, 4}, {5, 8}, {8, 9}},
			want:  []pattern{{1, 4}, {5, 9}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := mergePatternRanges(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPrintDefault(t *testing.T) {
	type input struct {
		reader   io.Reader
		name     string
		patterns []string
	}

	tests := map[string]struct {
		want string
		input
	}{
		"no patterns": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{},
			},
			want: fmt.Sprintf(NAME_PLACEHOLDER, "r_name"),
		},
		"no text": {
			input: input{
				reader:   strings.NewReader(""),
				name:     "r_name",
				patterns: []string{"abc"},
			},
			want: fmt.Sprintf(NAME_PLACEHOLDER, "r_name"),
		},
		"patterns not in text": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"fgh", "iyg"},
			},
			want: fmt.Sprintf(NAME_PLACEHOLDER, "r_name"),
		},
		"duplicate patterns in text": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"abc", "abc"},
			},
			want: fmt.Sprintf("%s1: abcd\n", fmt.Sprintf(NAME_PLACEHOLDER, "r_name")),
		},
		"nonoverlapping patterns in text on a single line": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"abc", "d"},
			},
			want: fmt.Sprintf("%s1: abcd\n", fmt.Sprintf(NAME_PLACEHOLDER, "r_name")),
		},
		"overlapping patterns in text on a single line": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"abc", "cd"},
			},
			want: fmt.Sprintf("%s1: abcd\n", fmt.Sprintf(NAME_PLACEHOLDER, "r_name")),
		},
		"nonoverlapping patterns in text on multiple lines": {
			input: input{
				reader:   strings.NewReader("abcd\ncde\nfgh"),
				name:     "r_name",
				patterns: []string{"abc", "fg"},
			},
			want: fmt.Sprintf("%s1: abcd\n3: fgh\n", fmt.Sprintf(NAME_PLACEHOLDER, "r_name")),
		},
		"overlapping patterns in text on multiple lines": {
			input: input{
				reader:   strings.NewReader("abcd\ncde\nfgh"),
				name:     "r_name",
				patterns: []string{"abc", "fg", "h", "gh"},
			},
			want: fmt.Sprintf("%s1: abcd\n3: fgh\n", fmt.Sprintf(NAME_PLACEHOLDER, "r_name")),
		},
	}
	for name, tt := range tests {
		buf := new(bytes.Buffer)
		t.Run(name, func(t *testing.T) {
			FprintDefault(tt.reader, buf, tt.name, tt.patterns)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestPrintQuite(t *testing.T) {
	type input struct {
		reader   io.Reader
		name     string
		patterns []string
	}

	tests := map[string]struct {
		want string
		input
	}{
		"no patterns": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{},
			},
			want: "",
		},
		"no text": {
			input: input{
				reader:   strings.NewReader(""),
				name:     "r_name",
				patterns: []string{"abc"},
			},
			want: "",
		},
		"patterns not in text": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"fgh", "iyg"},
			},
			want: "",
		},
		"duplicate patterns in text": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"abc", "abc"},
			},
			want: "abcd\n",
		},
		"nonoverlapping patterns in text on a single line": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"abc", "d"},
			},
			want: "abcd\n",
		},
		"overlapping patterns in text on a single line": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"abc", "cd"},
			},
			want: "abcd\n",
		},
		"nonoverlapping patterns in text on multiple lines": {
			input: input{
				reader:   strings.NewReader("abcd\ncde\nfgh"),
				name:     "r_name",
				patterns: []string{"abc", "fg"},
			},
			want: "abcd\nfgh\n",
		},
		"overlapping patterns in text on multiple lines": {
			input: input{
				reader:   strings.NewReader("abcd\ncde\nfgh"),
				name:     "r_name",
				patterns: []string{"abc", "fg", "h", "gh"},
			},
			want: "abcd\nfgh\n",
		},
	}
	for name, tt := range tests {
		buf := new(bytes.Buffer)
		t.Run(name, func(t *testing.T) {
			FprintQuite(tt.reader, buf, tt.patterns)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestPrintNonHuman(t *testing.T) {
	type input struct {
		reader   io.Reader
		name     string
		patterns []string
	}

	tests := map[string]struct {
		want string
		input
	}{
		"no patterns": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{},
			},
			want: "",
		},
		"no text": {
			input: input{
				reader:   strings.NewReader(""),
				name:     "r_name",
				patterns: []string{"abc"},
			},
			want: "",
		},
		"patterns not in text": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"fgh", "iyg"},
			},
			want: "",
		},
		"duplicate patterns in text": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"abc", "abc"},
			},
			want: "1 0 2 0 2\n",
		},
		"nonoverlapping patterns in text on a single line": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"abc", "d"},
			},
			want: "1 0 2 3 3\n",
		},
		"overlapping patterns in text on a single line": {
			input: input{
				reader:   strings.NewReader("abcd"),
				name:     "r_name",
				patterns: []string{"abc", "cd"},
			},
			want: "1 0 2 2 3\n",
		},
		"nonoverlapping patterns in text on multiple lines": {
			input: input{
				reader:   strings.NewReader("abcd\ncde\nfgh"),
				name:     "r_name",
				patterns: []string{"abc", "fg"},
			},
			want: "1 0 2\n3 0 1\n",
		},
		"overlapping patterns in text on multiple lines": {
			input: input{
				reader:   strings.NewReader("abcd\ncde\nfgh"),
				name:     "r_name",
				patterns: []string{"abc", "fg", "h", "gh"},
			},
			want: "1 0 2\n3 0 1 1 2 2 2\n",
		},
	}
	for name, tt := range tests {
		buf := new(bytes.Buffer)
		t.Run(name, func(t *testing.T) {
			FprintNonHuman(tt.reader, buf, tt.patterns)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

package wordwrap

import (
	"testing"
)

func TestWrapString(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		output string
		lim    uint
	}{
		{
			name:   "A simple word passes through",
			input:  "foo",
			output: "foo",
			lim:    4,
		},
		{
			name:   "A single word that is too long passes through, don't break words",
			input:  "foobarbaz",
			output: "foob\narba\nz",
			lim:    4,
		},
		{
			name:   "Lines are broken at whitespace",
			input:  "foo bar baz",
			output: "foo\nbar\nbaz",
			lim:    4,
		},
		{
			name:   "Lines are broken at whitespace even if words are too long, don't break words",
			input:  "foo bars bazzes",
			output: "foo\nbars\nbazz\nes",
			lim:    4,
		},
		{
			name:   "A word that would run beyond the width is wrapped",
			input:  "fo sop",
			output: "fo\nsop",
			lim:    4,
		},
		// Whitespace that trails a line and fits the width
		// passes through, as does whitespace prefixing an
		// explicit line break. A tab counts as one character.
		{
			name:   "A tab counts as one character",
			input:  "foo\nb\tar\n baz",
			output: "foo\nb\tar\n baz",
			lim:    4,
		},
		// Trailing whitespace is removed if it doesn't fit the width.
		// Runs of whitespace on which a line is broken are removed.
		{
			name:   "Trailing whitespace is removed if it doesn't fit the width",
			input:  "foo    \nb   ar   ",
			output: "foo\nb\nar",
			lim:    4,
		},
		{
			name:   "An explicit line break at the end of the input is preserved",
			input:  "foo bar baz\n",
			output: "foo\nbar\nbaz\n",
			lim:    4,
		},
		{
			name:   "Explicit break are always preserved",
			input:  "\nfoo bar\n\n\nbaz\n",
			output: "\nfoo\nbar\n\n\nbaz\n",
			lim:    4,
		},
		{
			name:   "Ignore ansi colors (A single word that is too long passes through, don't break words)",
			input:  "\x1b[34mfoo\x1b[0m\x1b[32mbar\x1b[0m\x1b[38;5;198mbaz\x1b[0m",
			output: "\x1b[34mfoo\x1b[0m\x1b[32mb\nar\x1b[0m\x1b[38;5;198mba\nz\x1b[0m",
			lim:    4,
		},
		{
			name:   "Ignore ansi colors (Lines are broken at whitespace even if words are too long, don't break words)",
			input:  "\x1b[34mfoo\x1b[0m \x1b[38;5;198mbars bazzes\x1b[0m",
			output: "\x1b[34mfoo\x1b[0m\n\x1b[38;5;198mbars\nbazz\nes\x1b[0m",
			lim:    4,
		},
		{
			name:   "Complete example",
			input:  " This is a list: \n\n\t* foo\n\t* bar\n\n\n\t* baz  \nBAM    ",
			output: " This\nis a\nlist: \n\n\t* foo\n\t* bar\n\n\n\t* baz\nBAM",
			lim:    6,
		},
	}
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := WrapString(tc.input, tc.lim)
			if actual != tc.output {
				t.Fatalf("Case %d Input:\n\n`%q`\n\nActual Output:\n\n`%q` expected `%q`", i, tc.input, actual, tc.output)
			}
		})
	}
}

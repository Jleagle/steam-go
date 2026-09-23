package steamvdf

import (
	"errors"
	"io"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestReadText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want KeyValue
	}{
		{
			name: "key value",
			in:   `"root" { "k" "v" }`,
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "k", Value: "v"}}},
		},
		{
			name: "nested objects",
			in:   `"root" { "a" { "b" { "c" "d" } } }`,
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "a", Children: []KeyValue{{Key: "b", Children: []KeyValue{{Key: "c", Value: "d"}}}}}}},
		},
		{
			name: "empty object has nil children",
			in:   `"root" { "empty" { } }`,
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "empty"}}},
		},
		{
			name: "single backslash before a letter does not escape the closing quote",
			in:   `"root" { "path" "C:\Games" "next" "value" }`,
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "path", Value: `C:\Games`}, {Key: "next", Value: "value"}}},
		},
		{
			name: "escape sequences are kept as written",
			in:   `"root" { "n" "line1\nline2" "q" "a\"b" "b" "a\\" "next" "v" }`,
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "n", Value: `line1\nline2`}, {Key: "q", Value: `a\"b`}, {Key: "b", Value: `a\\`}, {Key: "next", Value: "v"}}},
		},
		{
			name: "string may contain braces and newlines",
			in:   "\"root\" { \"k\" \"{\n}\" }",
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "k", Value: "{\n}"}}},
		},
		{
			name: "duplicate keys keep the last value in place",
			in:   `"root" { "k" "1" "other" "x" "k" "2" }`,
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "k", Value: "2"}, {Key: "other", Value: "x"}}},
		},
		{
			name: "repeated root key keeps the last block",
			in:   `"root" { "k" "1" } "root" { "k" "2" }`,
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "k", Value: "2"}}},
		},
		{
			name: "comments and directives are skipped",
			in:   "\"root\"\n{\n\t// double slash\n\t/ single slash\n\t/* block on one line */ \"ignored\" \"x\"\n\t#base \"file\"\n\t\"k\" \"v\" // trailing\n}\n",
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "k", Value: "v"}}},
		},
		{
			name: "trailing null byte is ignored",
			in:   "\"root\" { \"k\" \"v\" }\x00",
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "k", Value: "v"}}},
		},
		{
			name: "windows line endings",
			in:   "\"root\"\r\n{\r\n\t\"k\"\t\"v\"\r\n}\r\n",
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "k", Value: "v"}}},
		},
		{name: "empty input", in: "", want: KeyValue{}},
		{name: "whitespace only", in: " \n\t\r\n", want: KeyValue{}},
		{name: "comment only", in: "// nothing here\n", want: KeyValue{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ReadBytes([]byte(tc.in))
			assert.NilError(t, err)
			assert.DeepEqual(t, got, tc.want)
		})
	}
}

func TestReadTextErrors(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantIs   error  // checked with errors.Is when non-nil
		wantText string // checked with strings.Contains when non-empty
	}{
		{name: "lone key", in: `"a"`, wantIs: io.ErrUnexpectedEOF, wantText: `key "a" on line 1 has no value`},
		{name: "key without value before closing brace", in: `"root" { "a" }`, wantText: `key "a" on line 1 has no value`},
		{name: "stray open brace", in: `{`, wantText: "unexpected '{' on line 1"},
		{name: "stray close brace", in: `}`, wantText: "unexpected '}' on line 1"},
		{name: "missing close brace", in: `"root" { "k" "v"`, wantIs: io.ErrUnexpectedEOF, wantText: `missing closing brace for "root"`},
		{name: "extra close brace", in: "\"root\" { \"k\" \"v\" }\n}", wantText: "unexpected '}' on line 2"},
		{name: "unterminated string", in: "\"root\"\n{\n\t\"k\" \"v", wantIs: io.ErrUnexpectedEOF, wantText: "unterminated string starting on line 3"},
		{name: "unquoted token", in: `root { "k" "v" }`, wantText: "unexpected character 'r' on line 1"},
		{name: "line numbers count newlines inside strings", in: "\"root\"\n{\n\t\"multi\n\nline\" \"v\"\n\t\"k\"\n}", wantText: `key "k" on line 6 has no value`},
		{name: "two distinct roots", in: `"a" { } "b" { }`, wantIs: ErrMultipleRoots},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ReadBytes([]byte(tc.in))
			assert.Assert(t, err != nil, "expected an error")
			if tc.wantIs != nil {
				assert.Assert(t, errors.Is(err, tc.wantIs), "error %q should wrap %q", err, tc.wantIs)
			}
			if tc.wantText != "" {
				assert.Assert(t, strings.Contains(err.Error(), tc.wantText), "error %q should contain %q", err, tc.wantText)
			}
		})
	}
}

// nestedText builds a document of depth nested objects: "a" { "a" { ... } }.
func nestedText(depth int) string {
	return strings.Repeat(`"a" { `, depth) + strings.Repeat("} ", depth)
}

func TestReadTextNestingDepth(t *testing.T) {
	_, err := ReadBytes([]byte(nestedText(maxNestingDepth)))
	assert.NilError(t, err)

	_, err = ReadBytes([]byte(nestedText(maxNestingDepth + 1)))
	assert.Assert(t, errors.Is(err, ErrNestingTooDeep), "error %q should wrap ErrNestingTooDeep", err)
	assert.Assert(t, strings.Contains(err.Error(), `key "a" on line 1`), "error %q should name the key and line", err)
}

package steamvdf

import (
	"errors"
	"io"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

// binHeader is the 4-byte prefix found on PICS product-info blobs. The first
// byte is non-zero so that any prefix of three bytes or more is still
// detected as binary.
var binHeader = []byte{1, 0, 0, 0}

// binDoc builds a binary document: header, entries, closing TypeEnd.
func binDoc(entries ...[]byte) []byte {
	out := append([]byte{}, binHeader...)
	for _, e := range entries {
		out = append(out, e...)
	}
	return append(out, TypeEnd)
}

// binEntry builds one scalar entry: type byte, null-terminated key, payload.
func binEntry(typ byte, key string, payload ...byte) []byte {
	out := append([]byte{typ}, key...)
	out = append(out, 0)
	return append(out, payload...)
}

// binObject builds a nested object entry containing the given entries.
func binObject(key string, entries ...[]byte) []byte {
	out := binEntry(TypeNone, key)
	for _, e := range entries {
		out = append(out, e...)
	}
	return append(out, TypeEnd)
}

func TestReadBinary(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want KeyValue
	}{
		{
			name: "string",
			in:   binDoc(binObject("root", binEntry(TypeString, "s", 'v', 0))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "s", Value: "v"}}},
		},
		{
			name: "int32 negative",
			in:   binDoc(binObject("root", binEntry(TypeInt32, "i", 0xFF, 0xFF, 0xFF, 0xFF))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "i", Value: "-1"}}},
		},
		{
			name: "pointer renders as int32",
			in:   binDoc(binObject("root", binEntry(TypePointer, "p", 0x2A, 0, 0, 0))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "p", Value: "42"}}},
		},
		{
			name: "float32",
			in:   binDoc(binObject("root", binEntry(TypeFloat32, "f", 0x00, 0x00, 0xC0, 0x3F))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "f", Value: "1.5"}}},
		},
		{
			name: "color renders as r g b a",
			in:   binDoc(binObject("root", binEntry(TypeColor, "c", 1, 2, 3, 4))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "c", Value: "1 2 3 4"}}},
		},
		{
			name: "uint64 max",
			in:   binDoc(binObject("root", binEntry(TypeUint64, "u", 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "u", Value: "18446744073709551615"}}},
		},
		{
			name: "int64 negative",
			in:   binDoc(binObject("root", binEntry(TypeInt64, "l", 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "l", Value: "-2"}}},
		},
		{
			name: "nested objects and empty object",
			in:   binDoc(binObject("root", binObject("a", binObject("b", binEntry(TypeString, "c", 'd', 0))), binObject("empty"))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "a", Children: []KeyValue{{Key: "b", Children: []KeyValue{{Key: "c", Value: "d"}}}}}, {Key: "empty"}}},
		},
		{
			name: "duplicate keys keep the last value in place",
			in:   binDoc(binObject("root", binEntry(TypeString, "k", '1', 0), binEntry(TypeString, "o", 'x', 0), binEntry(TypeString, "k", '2', 0))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "k", Value: "2"}, {Key: "o", Value: "x"}}},
		},
		{
			name: "repeated root key keeps the last block",
			in:   binDoc(binObject("root", binEntry(TypeString, "k", '1', 0)), binObject("root", binEntry(TypeString, "k", '2', 0))),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "k", Value: "2"}}},
		},
		{
			name: "empty document",
			in:   binDoc(),
			want: KeyValue{},
		},
		{
			name: "trailing bytes after the root are ignored",
			in:   append(binDoc(binObject("root", binEntry(TypeString, "s", 'v', 0))), 0xAA, 0xBB),
			want: KeyValue{Key: "root", Children: []KeyValue{{Key: "s", Value: "v"}}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Assert(t, IsBinary(tc.in), "test input should be detected as binary")
			got, err := ReadBytes(tc.in)
			assert.NilError(t, err)
			assert.DeepEqual(t, got, tc.want)
		})
	}
}

func TestReadBinaryErrors(t *testing.T) {
	tests := []struct {
		name     string
		in       []byte
		wantIs   error
		wantText string
	}{
		{name: "short header", in: []byte{1, 0, 0}, wantIs: io.ErrUnexpectedEOF, wantText: "binary header"},
		{name: "header only", in: binHeader, wantIs: io.ErrUnexpectedEOF},
		{name: "missing final end marker", in: append(append([]byte{}, binHeader...), binObject("root")...), wantIs: io.ErrUnexpectedEOF, wantText: "at byte"},
		{name: "unterminated key", in: append(append([]byte{}, binHeader...), TypeString, 'k', 'e', 'y'), wantIs: io.ErrUnexpectedEOF},
		{name: "wide string", in: binDoc(binObject("root", binEntry(TypeWideString, "w", 0, 0))), wantIs: ErrWideString},
		{name: "unknown type", in: binDoc(binObject("root", binEntry(9, "x", 0))), wantText: `unknown value type 9 for key "x"`},
		{name: "two distinct roots", in: binDoc(binObject("a"), binObject("b")), wantIs: ErrMultipleRoots},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ReadBytes(tc.in)
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

// TestReadBinaryTruncated feeds every proper prefix of at least three bytes
// to the parser: each must be detected as binary, none may panic, and all
// must return an error. Shorter prefixes hold no null byte after IsBinary
// trims the trailing one and route to the text parser instead.
func TestReadBinaryTruncated(t *testing.T) {
	doc := binDoc(binObject("root",
		binEntry(TypeString, "s", 'v', 0),
		binObject("nested", binEntry(TypeInt32, "i", 1, 0, 0, 0), binEntry(TypeUint64, "u", 1, 0, 0, 0, 0, 0, 0, 0)),
		binEntry(TypeFloat32, "f", 0, 0, 0, 0),
	))

	for n := 3; n < len(doc); n++ {
		assert.Assert(t, IsBinary(doc[:n]), "prefix of length %d should be detected as binary", n)
		_, err := ReadBytes(doc[:n])
		assert.Assert(t, err != nil, "prefix of length %d should fail", n)
	}
}

// nestedBinary builds a document of depth nested empty objects all keyed "a".
func nestedBinary(depth int) []byte {
	var body []byte
	for i := 0; i < depth; i++ {
		body = append(body, binEntry(TypeNone, "a")...)
	}
	for i := 0; i < depth; i++ {
		body = append(body, TypeEnd)
	}
	return binDoc(body)
}

func TestReadBinaryNestingDepth(t *testing.T) {
	_, err := ReadBytes(nestedBinary(maxNestingDepth))
	assert.NilError(t, err)

	_, err = ReadBytes(nestedBinary(maxNestingDepth + 1))
	assert.Assert(t, errors.Is(err, ErrNestingTooDeep), "error %q should wrap ErrNestingTooDeep", err)
}

// Package steamvdf parses Steam's KeyValue (VDF) format, in both its text and
// binary encodings, into a KeyValue tree.
package steamvdf

import (
	"bytes"
	"errors"
	"os"
)

// ErrMultipleRoots is returned when a document has more than one distinct
// top-level key.
var ErrMultipleRoots = errors.New("vdf: document has more than one top-level key")

// maxNestingDepth bounds how deeply objects may nest. Steam documents nest a
// handful of levels; the bound keeps a hostile document from exhausting the
// goroutine stack.
const maxNestingDepth = 1000

// ErrNestingTooDeep is returned when a document nests objects more than
// maxNestingDepth levels deep.
var ErrNestingTooDeep = errors.New("vdf: nesting too deep")

// ReadFile parses the VDF document at path, detecting the text or binary
// encoding automatically.
func ReadFile(path string) (KeyValue, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return KeyValue{}, err
	}
	return ReadBytes(b)
}

// ReadBytes parses a VDF document, detecting the text or binary encoding
// automatically.
func ReadBytes(b []byte) (KeyValue, error) {
	if IsBinary(b) {
		return readBinary(b)
	}
	return readText(b)
}

// IsBinary reports whether b is a binary-encoded document. Text documents may
// end with a single null byte, which is ignored.
func IsBinary(b []byte) bool {
	b = bytes.TrimSuffix(b, []byte{TypeNone})
	return bytes.Contains(b, []byte{TypeNone})
}

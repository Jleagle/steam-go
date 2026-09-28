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

// Option configures how ReadFile and ReadBytes parse a document.
type Option func(*options)

type options struct {
	keepDuplicates bool
}

// KeepDuplicates keeps every entry of a key that repeats within an object, in
// document order, instead of replacing the earlier value. Some documents, such
// as the items_game.txt files of CS2 and Dota 2, split one section across many
// same-keyed blocks; GetChildren collects them. A document that repeats its
// root key parses to a single root holding the children of every occurrence.
func KeepDuplicates() Option {
	return func(o *options) {
		o.keepDuplicates = true
	}
}

// ReadFile parses the VDF document at path, detecting the text or binary
// encoding automatically.
func ReadFile(path string, opts ...Option) (KeyValue, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return KeyValue{}, err
	}
	return ReadBytes(b, opts...)
}

// ReadBytes parses a VDF document, detecting the text or binary encoding
// automatically.
func ReadBytes(b []byte, opts ...Option) (KeyValue, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	if IsBinary(b) {
		return readBinary(b, o)
	}
	return readText(b, o)
}

// singleRoot returns the one top-level node of a parsed document. When
// duplicates are kept, repeated root keys are merged into one node.
func singleRoot(root KeyValue, o options) (KeyValue, error) {
	switch len(root.Children) {
	case 0:
		return KeyValue{}, nil
	case 1:
		return root.Children[0], nil
	}

	if !o.keepDuplicates {
		return KeyValue{}, ErrMultipleRoots
	}

	merged := KeyValue{Key: root.Children[0].Key}
	for _, child := range root.Children {
		if child.Key != merged.Key {
			return KeyValue{}, ErrMultipleRoots
		}
		merged.Children = append(merged.Children, child.Children...)
	}
	return merged, nil
}

// IsBinary reports whether b is a binary-encoded document. Text documents may
// end with a single null byte, which is ignored.
func IsBinary(b []byte) bool {
	b = bytes.TrimSuffix(b, []byte{TypeNone})
	return bytes.Contains(b, []byte{TypeNone})
}

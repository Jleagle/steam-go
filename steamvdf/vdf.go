// Package steamvdf parses Steam's KeyValue (VDF) format, in both its text and
// binary encodings, into a KeyValue tree.
package steamvdf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// ErrMultipleRoots is returned when a document has more than one distinct
// top-level key.
var ErrMultipleRoots = errors.New("vdf: document has more than one top-level key")

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
		return readBinaryReader(bytes.NewReader(b))
	}
	return readText(b)
}

func readBinaryReader(r io.Reader) (kv KeyValue, err error) {

	var d float32
	err = binary.Read(r, binary.LittleEndian, &d)
	if err != nil {
		return kv, fmt.Errorf("binary.Read failed: %w", err)
	}

	root := KeyValue{}
	err = readBinary(r, &root, nil)
	return root, err
}

// IsBinary reports whether b is a binary-encoded document. Text documents may
// end with a single null byte, which is ignored.
func IsBinary(b []byte) bool {
	b = bytes.TrimSuffix(b, []byte{TypeNone})
	return bytes.Contains(b, []byte{TypeNone})
}

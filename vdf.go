package vdf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReadBinaryFile(path string) (kv KeyValue, err error) {

	r, err := os.Open(path)
	if err != nil {
		return kv, err
	}

	return ReadBinaryReader(r)
}

func ReadBinaryReader(r io.Reader) (kv KeyValue, err error) {

	var d float32
	err = binary.Read(r, binary.LittleEndian, &d)

	root := KeyValue{}
	err = readBinary(r, &root, nil)
	return root, err
}

func ReadBinaryBytes(b []byte) (kv KeyValue, err error) {
	return ReadBinaryReader(bytes.NewReader(b))
}

func IsBinary(b []byte) bool {

	b = bytes.TrimSuffix(b, []byte{TypeNone}) // Text VDF contains a null byte suffix

	return bytes.Contains(b, []byte{TypeNone})
}

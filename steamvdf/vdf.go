package steamvdf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

func ReadFile(path string) (kv KeyValue, err error) {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return kv, err
	}

	return ReadBytes(b)
}

func ReadBytes(b []byte) (kv KeyValue, err error) {
	if IsBinary(b) {
		return readBinaryReader(bytes.NewReader(b))
	} else {
		return readText(b)
	}
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

func IsBinary(b []byte) bool {

	b = bytes.TrimSuffix(b, []byte{TypeNone}) // Text VDF contains a null byte suffix

	return bytes.Contains(b, []byte{TypeNone})
}

package steamvdf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image/color"
	"io"
	"strconv"
)

const (
	TypeNone       byte = 0
	TypeString     byte = 1
	TypeInt32      byte = 2
	TypeFloat32    byte = 3
	TypePointer    byte = 4
	TypeWideString byte = 5
	TypeColor      byte = 6
	TypeUint64     byte = 7
	TypeEnd        byte = 8
	TypeInt64      byte = 10
)

var (
	ErrWideString = errors.New("WideString not supported")
)

// Thanks to https://github.com/SteamRE/SteamKit
func readBinary(r io.Reader, current *KeyValue, parent *KeyValue) (err error) {

	for {
		var b byte
		err := binary.Read(r, binary.LittleEndian, &b)
		if err != nil {
			return err
		}

		if b == TypeEnd {
			break
		}

		current.Key, err = readString(r)
		if err != nil {
			return err
		}

		switch b {
		case TypeNone:

			var child = KeyValue{}
			err = readBinary(r, &child, current)
			if err == nil {
				break
			}

		case TypeString:

			current.Value, err = readString(r)

		case TypeColor:

			var d color.NRGBA
			err = binary.Read(r, binary.LittleEndian, &d)

		case TypeInt32, TypePointer:

			var d int32
			err := binary.Read(r, binary.LittleEndian, &d)
			if err == nil {
				current.Value = strconv.Itoa(int(d))
			}

		case TypeFloat32:

			var d float32
			err := binary.Read(r, binary.LittleEndian, &d)
			if err == nil {
				current.Value = strconv.FormatFloat(float64(d), 'f', -1, 32)
			}

		case TypeWideString:

			return ErrWideString

		case TypeUint64:

			var d uint64
			err := binary.Read(r, binary.LittleEndian, &d)
			if err == nil {
				current.Value = strconv.FormatUint(d, 10)
			}

		case TypeInt64:

			var d int64
			err := binary.Read(r, binary.LittleEndian, &d)
			if err == nil {
				current.Value = strconv.FormatInt(d, 10)
			}

		default:

			err = fmt.Errorf("vdf: unknown pack type %d", b)

		}

		if err != nil {
			return err
		}

		if parent != nil {
			parent.SetChild(*current)
		}

		current = &KeyValue{}
	}

	return nil
}

func readString(r io.Reader) (string, error) {
	c := make([]byte, 0)
	var err error
	for {
		var b byte
		err = binary.Read(r, binary.LittleEndian, &b)
		if b == byte(0x0) || err != nil {
			break
		}
		c = append(c, b)
	}
	return string(c), err
}

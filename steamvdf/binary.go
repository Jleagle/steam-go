package steamvdf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
)

// Binary-format value type markers. Each entry in an object is a type byte, a
// null-terminated key, then a payload whose layout the type determines.
const (
	TypeNone       byte = 0 // nested object; entries follow until TypeEnd
	TypeString     byte = 1
	TypeInt32      byte = 2
	TypeFloat32    byte = 3
	TypePointer    byte = 4
	TypeWideString byte = 5
	TypeColor      byte = 6
	TypeUint64     byte = 7
	TypeEnd        byte = 8 // closes the enclosing object
	TypeInt64      byte = 10
)

// ErrWideString is returned when a binary document contains a wide-string
// value, which this package does not decode.
var ErrWideString = errors.New("vdf: wide string values are not supported")

// binaryHeaderSize is the length of the uint32 that precedes the root object
// in PICS product-info blobs. Its meaning is undocumented; SteamKit skips it
// too.
const binaryHeaderSize = 4

// readBinary parses a binary-format document. Integers render in decimal,
// floats with the shortest representation that round-trips, and colours as
// "r g b a" to match Valve's text rendering. Unless o keeps duplicates,
// duplicate keys within an object replace the earlier value in place. Errors
// name the byte offset at which parsing stopped.
func readBinary(data []byte, o options) (KeyValue, error) {
	r := binaryReader{data: data, keepDuplicates: o.keepDuplicates}
	if _, err := r.take(binaryHeaderSize); err != nil {
		return KeyValue{}, fmt.Errorf("vdf: reading binary header: %w", err)
	}

	var root KeyValue
	if err := r.readObject(&root, 0); err != nil {
		return KeyValue{}, fmt.Errorf("vdf: at byte %d: %w", r.pos, err)
	}

	return singleRoot(root, o)
}

// binaryReader is a cursor over a binary document.
type binaryReader struct {
	data           []byte
	pos            int
	keepDuplicates bool
}

// take returns the next n bytes.
func (r *binaryReader) take(n int) ([]byte, error) {
	if len(r.data)-r.pos < n {
		return nil, io.ErrUnexpectedEOF
	}
	b := r.data[r.pos : r.pos+n]
	r.pos += n
	return b, nil
}

func (r *binaryReader) readByte() (byte, error) {
	b, err := r.take(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// readString returns the next null-terminated string.
func (r *binaryReader) readString() (string, error) {
	end := bytes.IndexByte(r.data[r.pos:], 0)
	if end < 0 {
		return "", io.ErrUnexpectedEOF
	}
	s := string(r.data[r.pos : r.pos+end])
	r.pos += end + 1
	return s, nil
}

// readObject reads entries into parent until the object's TypeEnd marker.
// depth is the nesting level of parent; the root object is at depth 0.
func (r *binaryReader) readObject(parent *KeyValue, depth int) error {
	setter := childSetter{kv: parent, keepDuplicates: r.keepDuplicates}
	for {
		typ, err := r.readByte()
		if err != nil {
			return err
		}
		if typ == TypeEnd {
			return nil
		}

		key, err := r.readString()
		if err != nil {
			return err
		}

		child := KeyValue{Key: key}
		if typ == TypeNone {
			if depth+1 > maxNestingDepth {
				return fmt.Errorf("%w: key %q", ErrNestingTooDeep, key)
			}
			err = r.readObject(&child, depth+1)
		} else {
			child.Value, err = r.readValue(typ, key)
		}
		if err != nil {
			return err
		}

		setter.set(child)
	}
}

// readValue decodes the payload of a scalar entry of type typ.
func (r *binaryReader) readValue(typ byte, key string) (string, error) {
	switch typ {
	case TypeString:
		return r.readString()
	case TypeInt32, TypePointer:
		b, err := r.take(4)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(int64(int32(binary.LittleEndian.Uint32(b))), 10), nil
	case TypeFloat32:
		b, err := r.take(4)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(float64(math.Float32frombits(binary.LittleEndian.Uint32(b))), 'f', -1, 32), nil
	case TypeColor:
		b, err := r.take(4)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d %d %d %d", b[0], b[1], b[2], b[3]), nil
	case TypeUint64:
		b, err := r.take(8)
		if err != nil {
			return "", err
		}
		return strconv.FormatUint(binary.LittleEndian.Uint64(b), 10), nil
	case TypeInt64:
		b, err := r.take(8)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(int64(binary.LittleEndian.Uint64(b)), 10), nil
	case TypeWideString:
		return "", ErrWideString
	default:
		return "", fmt.Errorf("vdf: unknown value type %d for key %q", typ, key)
	}
}

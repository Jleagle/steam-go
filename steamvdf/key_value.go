package steamvdf

import (
	"encoding/json"
	"slices"
	"strings"
)

// KeyValue is a node in a VDF document: a key with either a string value or
// child nodes.
type KeyValue struct {
	Key      string     `json:"key"`
	Value    string     `json:"value"`
	Children []KeyValue `json:"children"`
}

// SortChildren sorts the children by key, recursively.
func (kv *KeyValue) SortChildren() {
	slices.SortFunc(kv.Children, func(a, b KeyValue) int {
		return strings.Compare(a.Key, b.Key)
	})
	for i := range kv.Children {
		kv.Children[i].SortChildren()
	}
}

// GetChildrenAsSlice returns each child rendered with String, in order.
func (kv KeyValue) GetChildrenAsSlice() []string {
	ret := make([]string, 0, len(kv.Children))
	for _, child := range kv.Children {
		ret = append(ret, child.String())
	}
	return ret
}

// GetChildrenAsMap returns the children keyed by Key, each rendered with
// String.
func (kv KeyValue) GetChildrenAsMap() map[string]string {
	ret := make(map[string]string, len(kv.Children))
	for _, child := range kv.Children {
		ret[child.Key] = child.String()
	}
	return ret
}

// GetChild returns the child with the given key.
func (kv KeyValue) GetChild(key string) (KeyValue, bool) {
	for _, child := range kv.Children {
		if child.Key == key {
			return child, true
		}
	}
	return KeyValue{}, false
}

// GetChildren returns every child with the given key, in order. Documents
// read with KeepDuplicates can hold several.
func (kv KeyValue) GetChildren(key string) []KeyValue {
	var ret []KeyValue
	for _, child := range kv.Children {
		if child.Key == key {
			ret = append(ret, child)
		}
	}
	return ret
}

// SetChild replaces the child with the same key in place, or appends child
// when there is none.
func (kv *KeyValue) SetChild(child KeyValue) {
	for i, existing := range kv.Children {
		if existing.Key == child.Key {
			kv.Children[i] = child
			return
		}
	}
	kv.Children = append(kv.Children, child)
}

// HasChild reports whether a child with the given key exists.
func (kv KeyValue) HasChild(key string) bool {
	_, ok := kv.GetChild(key)
	return ok
}

// childIndexThreshold is the number of children at which childSetter starts
// indexing keys. Below it, the linear scan in SetChild is cheaper than a map.
const childIndexThreshold = 16

// childSetter inserts children into kv with SetChild's replace-or-append
// semantics while keeping each insert O(1) once the object is wide. The
// parsers use it so that objects with thousands of children do not rescan
// every existing child on each insert. With keepDuplicates it always appends.
type childSetter struct {
	kv             *KeyValue
	index          map[string]int // key to position in kv.Children, built lazily
	keepDuplicates bool
}

func (s *childSetter) set(child KeyValue) {
	if s.keepDuplicates {
		s.kv.Children = append(s.kv.Children, child)
		return
	}

	if s.index == nil {
		if len(s.kv.Children) < childIndexThreshold {
			s.kv.SetChild(child)
			return
		}
		s.index = make(map[string]int, len(s.kv.Children)*2)
		for i, existing := range s.kv.Children {
			s.index[existing.Key] = i
		}
	}

	if i, ok := s.index[child.Key]; ok {
		s.kv.Children[i] = child
		return
	}
	s.index[child.Key] = len(s.kv.Children)
	s.kv.Children = append(s.kv.Children, child)
}

// String renders the node. A node with a Value returns it. A node with
// children returns the JSON object of its children (as ToMapInner builds it).
// A node with neither returns a JSON object mapping its key to null.
func (kv KeyValue) String() string {
	if kv.Value != "" {
		return kv.Value
	}

	var m map[string]any
	if len(kv.Children) == 0 {
		m = map[string]any{kv.Key: nil}
	} else {
		m = toMap(kv)
	}

	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

// ToMapOuter converts the node to nested maps with this node's Key as the
// single top-level key.
func (kv KeyValue) ToMapOuter() map[string]any {
	return toMap(KeyValue{Children: []KeyValue{kv}})
}

// ToMapInner converts the node's children to nested maps.
func (kv KeyValue) ToMapInner() map[string]any {
	return toMap(kv)
}

// toMap converts children to a map: a child with a Value becomes a string, a
// child with neither Value nor children becomes nil, and a child with
// children recurses.
func toMap(kv KeyValue) map[string]any {
	m := make(map[string]any, len(kv.Children))
	for _, child := range kv.Children {
		switch {
		case child.Value != "":
			m[child.Key] = child.Value
		case len(child.Children) == 0:
			m[child.Key] = nil
		default:
			m[child.Key] = toMap(child)
		}
	}
	return m
}

// FromMap is the inverse of ToMapOuter: it converts nested maps back into a
// KeyValue tree, with children ordered by key. A map with a single key returns
// that key's node; a map with several keys returns a root node with an empty
// Key holding one child per key.
func FromMap(m map[string]any) KeyValue {
	kv := fromMap("", m)
	if len(kv.Children) == 1 {
		return kv.Children[0]
	}
	return kv
}

// fromMap converts one map entry. Strings become values, maps recurse, and
// anything else (including nil) becomes an empty node.
func fromMap(key string, v any) KeyValue {
	out := KeyValue{Key: key}

	switch v := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			out.Children = append(out.Children, fromMap(k, v[k]))
		}
	case string:
		out.Value = v
	}

	return out
}

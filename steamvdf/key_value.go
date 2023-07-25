package steamvdf

import (
	"encoding/json"
	"sort"
)

type KeyValue struct {
	Key      string     `json:"key"`
	Value    string     `json:"value"`
	Children []KeyValue `json:"children"`
}

func (kv *KeyValue) SortChildren() {

	sort.Slice(kv.Children, func(i, j int) bool {
		return kv.Children[i].Key < kv.Children[j].Key
	})

	for _, v := range kv.Children {
		v.SortChildren()
	}
}

func (kv KeyValue) GetChildrenAsSlice() (ret []string) {
	for _, v := range kv.Children {
		ret = append(ret, v.String())
	}
	return ret
}

func (kv KeyValue) GetChildrenAsMap() (ret map[string]string) {
	ret = map[string]string{}
	for _, v := range kv.Children {
		ret[v.Key] = v.String()
	}
	return ret
}

func (kv KeyValue) GetChild(key string) (child KeyValue, found bool) {
	for _, child := range kv.Children {
		if child.Key == key {
			return child, true
		}
	}
	return child, false
}

func (kv *KeyValue) SetChild(value KeyValue) {
	for k, child := range kv.Children {
		if child.Key == value.Key {
			kv.Children[k] = value
			return
		}
	}
	kv.Children = append(kv.Children, value)
}

func (kv *KeyValue) HasChild(key string) bool {
	for _, child := range kv.Children {
		if child.Key == key {
			return true
		}
	}
	return false
}

// Returns kv.Value or the children in json form
func (kv KeyValue) String() string {

	if kv.Value != "" {
		return kv.Value
	}

	if len(kv.Children) == 0 {
		b, err := json.Marshal(map[string]interface{}{kv.Key: nil})
		if err != nil || string(b) == "{}" {
			return ""
		}
		return string(b)
	}

	b, err := json.Marshal(toMap(kv))
	if err != nil || string(b) == "{}" {
		return ""
	}

	return string(b)
}

// Transforms to nested maps
// Includes top level
func (kv KeyValue) ToMapOuter() (m map[string]interface{}) {

	return toMap(KeyValue{
		Key:      "",
		Children: []KeyValue{kv},
	})
}

// Transforms to nested maps
// Does not include top level
func (kv KeyValue) ToMapInner() (m map[string]interface{}) {

	return toMap(kv)
}

func toMap(kv KeyValue) map[string]interface{} {

	m := map[string]interface{}{}

	for _, child := range kv.Children {

		if child.Value != "" {
			m[child.Key] = child.Value
		} else if len(child.Children) == 0 {
			m[child.Key] = nil
		} else {
			m[child.Key] = toMap(child)
		}
	}

	return m
}

func FromMap(m map[string]interface{}) KeyValue {

	kv := fromMap("", m)
	if len(kv.Children) > 0 {
		return kv.Children[0]
	}
	return kv
}

func fromMap(key string, m interface{}) (out KeyValue) {

	out.Key = key

	switch m := m.(type) {
	case map[string]interface{}:
		for k, v := range m {
			out.Children = append(out.Children, fromMap(k, v))
		}
		out.Value = ""
	case string:
		out.Value = m
	}

	return out
}

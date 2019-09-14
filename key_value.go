package vdf

import (
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
		ret = append(ret, v.Value)
	}
	return ret
}

func (kv KeyValue) GetChildrenAsMap() (ret map[string]string) {
	ret = map[string]string{}
	for _, v := range kv.Children {
		ret[v.Key] = v.Value
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

func (kv KeyValue) ToMap() (m map[string]interface{}) {

	return toMap(KeyValue{
		Key:      "",
		Children: []KeyValue{kv},
	})
}

func toMap(kv KeyValue) map[string]interface{} {

	m := map[string]interface{}{}

	for _, child := range kv.Children {

		if child.Value == "" {
			m[child.Key] = toMap(child)
		} else {
			m[child.Key] = child.Value
		}
	}

	return m
}

func FromMap(m map[string]interface{}) KeyValue {
	return fromMap("", m).Children[0]
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

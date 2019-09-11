package vdf

type KeyValue struct {
	Key      string     `json:"key"`
	Value    string     `json:"value"`
	Children []KeyValue `json:"children"`
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

func (kv KeyValue) ChildrenAsMap() (m map[string]interface{}) {

	m = map[string]interface{}{}

	for _, v := range kv.Children {

		if v.Value == "" {
			m[v.Key] = v.ChildrenAsMap()
		} else {
			m[v.Key] = v.Value
		}
	}

	return m
}

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

func (kv KeyValue) Map() (m map[string]interface{}) {

	return toMap(KeyValue{
		Key:      "",
		Children: []KeyValue{kv},
	})
}

func toMap(kv KeyValue) (m map[string]interface{}) {

	m = map[string]interface{}{}

	for _, child := range kv.Children {

		if child.Value == "" {
			m[child.Key] = toMap(child)
		} else {
			m[child.Key] = child.Value
		}
	}

	return m
}

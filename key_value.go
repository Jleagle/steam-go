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

func (kv KeyValue) Map() (m map[string]interface{}) {

	m = map[string]interface{}{}

	if kv.Value != "" {
		m[kv.Key] = kv.Value
	} else {
		c := map[string]interface{}{}
		for _, v := range kv.Children {
			c[v.Key] = v.Map()
		}
		m[kv.Key] = c
	}

	return m
}

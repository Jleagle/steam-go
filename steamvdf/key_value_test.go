package steamvdf

import (
	"strconv"
	"testing"

	"gotest.tools/v3/assert"
)

// TestChildSetterMatchesSetChild inserts the same sequence, with duplicates
// on both sides of the index threshold, through SetChild and childSetter and
// requires identical results.
func TestChildSetterMatchesSetChild(t *testing.T) {
	var direct, indexed KeyValue
	setter := childSetter{kv: &indexed}

	insert := func(key, value string) {
		direct.SetChild(KeyValue{Key: key, Value: value})
		setter.set(KeyValue{Key: key, Value: value})
	}

	for i := 0; i < childIndexThreshold*3; i++ {
		insert("k"+strconv.Itoa(i), "v"+strconv.Itoa(i))
		if i%5 == 0 {
			insert("k0", "replaced at "+strconv.Itoa(i)) // duplicate below and above the threshold
		}
	}
	insert("k"+strconv.Itoa(childIndexThreshold-1), "replaced entry inserted before the index existed")
	insert("k"+strconv.Itoa(childIndexThreshold*2), "replaced above threshold")

	assert.DeepEqual(t, indexed, direct)
	assert.Equal(t, len(indexed.Children), childIndexThreshold*3)
	assert.Equal(t, indexed.Children[0].Value, "replaced at "+strconv.Itoa(childIndexThreshold*3-3))
}

func TestFromMapSingleKeyReturnsThatNode(t *testing.T) {
	kv := FromMap(map[string]any{"root": map[string]any{"b": "2", "a": "1", "empty": nil, "obj": map[string]any{"z": "26"}}})
	assert.DeepEqual(t, kv, KeyValue{Key: "root", Children: []KeyValue{
		{Key: "a", Value: "1"},
		{Key: "b", Value: "2"},
		{Key: "empty"},
		{Key: "obj", Children: []KeyValue{{Key: "z", Value: "26"}}},
	}})
}

func TestFromMapMultipleKeysReturnsRoot(t *testing.T) {
	kv := FromMap(map[string]any{"b": "2", "a": "1", "c": "3"})
	assert.DeepEqual(t, kv, KeyValue{Children: []KeyValue{
		{Key: "a", Value: "1"},
		{Key: "b", Value: "2"},
		{Key: "c", Value: "3"},
	}})
}

func TestFromMapEmpty(t *testing.T) {
	assert.DeepEqual(t, FromMap(map[string]any{}), KeyValue{})
}

func TestSortChildrenRecurses(t *testing.T) {
	kv := KeyValue{Key: "root", Children: []KeyValue{
		{Key: "b", Children: []KeyValue{{Key: "y"}, {Key: "x"}}},
		{Key: "a"},
	}}
	kv.SortChildren()
	assert.DeepEqual(t, kv, KeyValue{Key: "root", Children: []KeyValue{
		{Key: "a"},
		{Key: "b", Children: []KeyValue{{Key: "x"}, {Key: "y"}}},
	}})
}

func TestString(t *testing.T) {
	assert.Equal(t, KeyValue{Key: "k", Value: "v"}.String(), "v")
	assert.Equal(t, KeyValue{Key: "k"}.String(), `{"k":null}`)
	assert.Equal(t, KeyValue{Key: "k", Children: []KeyValue{{Key: "a", Value: "1"}, {Key: "e"}, {Key: "o", Children: []KeyValue{{Key: "z", Value: "26"}}}}}.String(), `{"a":"1","e":null,"o":{"z":"26"}}`)
}

func TestGetChildrenAsMapAndSlice(t *testing.T) {
	kv := KeyValue{Key: "k", Children: []KeyValue{{Key: "a", Value: "1"}, {Key: "o", Children: []KeyValue{{Key: "z", Value: "26"}}}}}
	assert.DeepEqual(t, kv.GetChildrenAsMap(), map[string]string{"a": "1", "o": `{"z":"26"}`})
	assert.DeepEqual(t, kv.GetChildrenAsSlice(), []string{"1", `{"z":"26"}`})
}

func TestGetSetHasChild(t *testing.T) {
	var kv KeyValue
	_, ok := kv.GetChild("a")
	assert.Assert(t, !ok)
	assert.Assert(t, !kv.HasChild("a"))

	kv.SetChild(KeyValue{Key: "a", Value: "1"})
	kv.SetChild(KeyValue{Key: "b", Value: "2"})
	kv.SetChild(KeyValue{Key: "a", Value: "3"})
	assert.DeepEqual(t, kv.Children, []KeyValue{{Key: "a", Value: "3"}, {Key: "b", Value: "2"}})

	child, ok := kv.GetChild("a")
	assert.Assert(t, ok)
	assert.Equal(t, child.Value, "3")
	assert.Assert(t, kv.HasChild("b"))
}

func TestGetChildren(t *testing.T) {
	kv := KeyValue{Children: []KeyValue{{Key: "a", Value: "1"}, {Key: "b", Value: "2"}, {Key: "a", Value: "3"}}}
	assert.DeepEqual(t, kv.GetChildren("a"), []KeyValue{{Key: "a", Value: "1"}, {Key: "a", Value: "3"}})
	assert.Assert(t, kv.GetChildren("missing") == nil)
}

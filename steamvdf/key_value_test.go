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
	insert("k"+strconv.Itoa(childIndexThreshold-1), "replaced last below threshold")
	insert("k"+strconv.Itoa(childIndexThreshold*2), "replaced above threshold")

	assert.DeepEqual(t, indexed, direct)
	assert.Equal(t, len(indexed.Children), childIndexThreshold*3)
	assert.Equal(t, indexed.Children[0].Value, "replaced at "+strconv.Itoa(childIndexThreshold*3-3))
}

func TestChildSetterLeavesEmptyObjectNil(t *testing.T) {
	var kv KeyValue
	_ = childSetter{kv: &kv}
	assert.Assert(t, kv.Children == nil)
}

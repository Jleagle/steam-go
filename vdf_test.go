package vdf

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"

	"gotest.tools/assert"
)

func TestUnmarshalBinary(t *testing.T) {

	ints := regexp.MustCompile("[0-9]+")

	files, _ := filepath.Glob("testdata/package_*.bin")
	for _, file := range files {

		id := ints.FindString(file)

		kv, err := ReadBinaryFile(file)
		if err != nil {
			t.Error(err)
		}

		assert.Assert(t, kv.Key == id)
		assert.Assert(t, len(kv.Children) > 0)

		packageid, found := kv.GetChild("packageid")
		assert.Assert(t, found)
		assert.Assert(t, packageid.Value == id)

		switch id {
		case "55058":

			fmt.Println("Testing " + id)

			child, ok := kv.GetChild("billingtype")
			assert.Assert(t, ok)
			assert.Assert(t, child.Value == "10", id)

			child, ok = kv.GetChild("licensetype")
			assert.Assert(t, ok)
			assert.Assert(t, child.Value == "1", id)

			extended, ok := kv.GetChild("extended")
			assert.Assert(t, ok)

			child, ok = extended.GetChild("allowcrossregiontradingandgifting")
			assert.Assert(t, ok)
			assert.Assert(t, child.Value == "false", id)

		case "228891":

			fmt.Println("Testing " + id)

			child, ok := kv.GetChild("billingtype")
			assert.Assert(t, ok)
			assert.Assert(t, child.Value == "3", id)

			child, ok = kv.GetChild("licensetype")
			assert.Assert(t, ok)
			assert.Assert(t, child.Value == "1", id)

			//
			extended, ok := kv.GetChild("extended")
			assert.Assert(t, ok)

			child, ok = extended.GetChild("curatorconnect")
			assert.Assert(t, ok)
			assert.Assert(t, child.Value == "1", id)

			child, ok = extended.GetChild("releasestateoverride")
			assert.Assert(t, ok)
			assert.Assert(t, child.Value == "released", id)

			//
			appids, ok := kv.GetChild("appids")
			assert.Assert(t, ok)
			assert.Assert(t, len(appids.Children) == 2)

			child, ok = appids.GetChild("0")
			assert.Assert(t, ok)
			assert.Assert(t, child.Value == "680590", id)

			child, ok = appids.GetChild("1")
			assert.Assert(t, ok)
			assert.Assert(t, child.Value == "1155030", id)

		case "394363":

			fmt.Println("Testing " + id)

			extended, ok := kv.GetChild("extended")
			assert.Assert(t, ok)
			assert.Assert(t, len(extended.Children) == 7)

		}
	}
}

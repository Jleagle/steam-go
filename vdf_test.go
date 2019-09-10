package vdf

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestReadBinary(t *testing.T) {

	ints := regexp.MustCompile("[0-9]+")

	files, _ := filepath.Glob("testdata/*")
	for _, file := range files {

		id := ints.FindString(file)
		isApp := strings.Contains(file, "app")

		kv, err := ReadFile(file)
		if err != nil {
			t.Error(err)
		}

		if isApp {

			assert.Assert(t, kv.Key == "appinfo")

			appid, ok := kv.GetChild("appid")
			assert.Assert(t, ok)
			assert.Assert(t, appid.Value == id)

		} else {

			assert.Assert(t, kv.Key == id)
			assert.Assert(t, len(kv.Children) > 0)

			packageid, ok := kv.GetChild("packageid")
			assert.Assert(t, ok)
			assert.Assert(t, packageid.Value == id)
		}

		switch id {
		case "10":

			fmt.Println("Testing " + id)

			assert.Assert(t, len(kv.Children) == 6)

			common, ok := kv.GetChild("common")
			assert.Assert(t, ok)

			associations, ok := common.GetChild("associations")
			assert.Assert(t, ok)
			assert.Assert(t, len(associations.Children) == 2, id)
			assert.Assert(t, associations.Children[0].Children[1].Value == "Valve", id)
			assert.Assert(t, associations.Children[1].Children[1].Value == "Valve", id)

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
			assert.Assert(t, len(extended.Children) == 1)

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
			assert.Assert(t, len(extended.Children) == 2)

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

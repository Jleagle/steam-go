package steamvdf

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
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

		switch file {
		case "testdata/app_600760.vdf":

			fmt.Println("Testing " + file)

			str := kv.String()

			assert.Assert(t, strings.Contains(str, "\"gamedir\":null"), "kv.String()", "empty value is not null")

			var vdf struct {
				Extended struct {
					Gamedir string `json:"gamedir"`
				} `json:"extended"`
			}

			err := json.Unmarshal([]byte(str), &vdf)

			assert.Assert(t, err == nil, "json.Unmarshal()", "expected to unmarshal empty string value", "error", err)

		case "testdata/app_212200.vdf":

			fmt.Println("Testing " + file)

			assert.Assert(t, strings.Contains(kv.String(), `\\\"`))

		case "testdata/app_10.vdf":

			fmt.Println("Testing " + file)

			assert.Assert(t, len(kv.Children) == 6)

			common, ok := kv.GetChild("common")
			assert.Assert(t, ok)

			associations, ok := common.GetChild("associations")
			assert.Assert(t, ok)
			assert.Assert(t, len(associations.Children) == 2, id)
			assert.Assert(t, associations.Children[0].Children[1].Value == "Valve", id)
			assert.Assert(t, associations.Children[1].Children[1].Value == "Valve", id)

			kvMap := kv.ToMapOuter()

			appinfo := kvMap["appinfo"].(map[string]interface{})
			assert.Assert(t, appinfo["appid"] == "10", id)

			commonM := appinfo["common"].(map[string]interface{})
			assert.Assert(t, len(commonM) == 28, id)

			kvStruct := FromMap(kvMap)
			kv.SortChildren()
			kvStruct.SortChildren()
			assert.DeepEqual(t, kv, kvStruct)

		case "testdata/app_917720.vdf":

			fmt.Println("Testing " + file)

			common, ok := kv.GetChild("common")
			assert.Assert(t, ok)
			assert.Assert(t, len(common.Children) == 25, id)

			assets, ok := common.GetChild("library_assets")
			assert.Assert(t, ok)
			assert.Assert(t, len(assets.Children) == 4, id)

		case "testdata/package_55058.bin":

			fmt.Println("Testing " + file)

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

		case "testdata/package_228891.bin":

			fmt.Println("Testing " + file)

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

		case "testdata/package_394363.bin":

			fmt.Println("Testing " + file)

			extended, ok := kv.GetChild("extended")
			assert.Assert(t, ok)
			assert.Assert(t, len(extended.Children) == 7)

		}
	}
}

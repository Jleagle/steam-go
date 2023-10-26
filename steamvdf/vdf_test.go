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

const improperlyEscapedQuote = "improperly parsed escaped quote (\\\")"

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

			assert.Assert(t, kv.Key == "appinfo", "id", id)

			appid, ok := kv.GetChild("appid")
			assert.Assert(t, ok, "id", id)
			assert.Assert(t, appid.Value == id)

		} else {

			assert.Assert(t, kv.Key == id)
			assert.Assert(t, len(kv.Children) > 0)

			packageid, ok := kv.GetChild("packageid")
			assert.Assert(t, ok)
			assert.Assert(t, packageid.Value == id)
		}

		switch file {
		case "testdata/app_0.vdf":
			fmt.Println("Testing " + file)

			str := kv.String()

			assert.Assert(t, !strings.Contains(str, "ignoredKey"), "block comment should consume keys and values on its line")

			_, err := json.MarshalIndent(str, "", "    ")

			assert.Assert(t, err == nil, fmt.Sprintf("json.MarshalIndent() failed to marshal kv.String() value: %v", err))

		case "testdata/app_574720.vdf":

			fmt.Println("Testing " + file)

			var vdf struct {
				Ufs struct {
					Rootoverrides map[string]struct {
						Os             string `json:"os"`
						Oscompare      string `json:"oscompare"`
						Pathtransforms map[string]struct {
							Find    string `json:"find"`
							Replace string `json:"replace"`
						} `json:"pathtransforms"`
						Root       string `json:"root"`
						Useinstead string `json:"useinstead"`
					} `json:"rootoverrides"`
				} `json:"ufs"`
			}

			err := json.Unmarshal([]byte(kv.String()), &vdf)
			assert.Assert(t, err == nil, "json.Unmarshal()", "expected to unmarshal config value", "error", err)

			assert.Assert(t, vdf.Ufs.Rootoverrides["0"].Pathtransforms["0"].Find == `Mirage Game Studios\\Little Big Workshop\\saves`, fmt.Sprintf("got `%s`\nerror in kv.String(): %s", vdf.Ufs.Rootoverrides["0"].Pathtransforms["0"].Find, improperlyEscapedQuote))
			assert.Assert(t, vdf.Ufs.Rootoverrides["0"].Pathtransforms["0"].Replace == `Mirage Game Studios\\Little Big Workshop\\`, fmt.Sprintf("got `%s`\nerror in kv.String(): %s", vdf.Ufs.Rootoverrides["0"].Pathtransforms["0"].Replace, improperlyEscapedQuote))

			assert.Assert(t, vdf.Ufs.Rootoverrides["1"].Pathtransforms["0"].Find == `Mirage Game Studios\\Little Big Workshop\\saves`, fmt.Sprintf("got `%s`\nerror in kv.String(): %s", vdf.Ufs.Rootoverrides["1"].Pathtransforms["0"].Find, improperlyEscapedQuote))
			assert.Assert(t, vdf.Ufs.Rootoverrides["1"].Pathtransforms["0"].Replace == `unity.Mirage Game Studios.Little Big Workshop\\`, fmt.Sprintf("got `%s`\nerror in kv.String(): %s", vdf.Ufs.Rootoverrides["1"].Pathtransforms["0"].Replace, improperlyEscapedQuote))

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

			var vdf struct {
				Config struct {
					Launch map[string]struct {
						Arguments string `json:"arguments"`
						Config    struct {
							Oslist string `json:"oslist"`
						} `json:"config"`
						Description    string `json:"description"`
						DescriptionLoc struct {
							English string `json:"english"`
						} `json:"description_loc"`
						Executable string `json:"executable"`
						Type       string `json:"type"`
					} `json:"launch"`
				} `json:"config"`
			}

			err := json.Unmarshal([]byte(kv.String()), &vdf)
			assert.Assert(t, err == nil, "json.Unmarshal()", "expected to unmarshal config value", "error", err)

			assert.Assert(t, vdf.Config.Launch["0"].Arguments == `\"Client.exe code:1622 locale:USA env:Regular ver:246 logip:35.162.171.43 logport:11000 chatip:54.214.176.167 chatport:8002 setting:\\\"file://data/features.xml\\\" sn:{tracking_uid} sid:{tracking_sessionid} /P:{passport} -Steam\" --nx:title=Mabinogi --nx:serviceId=880915460`, fmt.Sprintf("got `%s`\nerror in kv.String(): %s", vdf.Config.Launch["0"].Arguments, improperlyEscapedQuote))

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

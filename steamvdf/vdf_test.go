package steamvdf

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

var digits = regexp.MustCompile(`[0-9]+`)

const improperlyEscapedQuote = `improperly parsed escaped quote (\")`

// fileChecks holds assertions specific to individual corpus files, keyed by
// file name. TestCorpus runs the matching check after the generic ones.
var fileChecks = map[string]func(t *testing.T, kv KeyValue){
	"app_0.vdf": func(t *testing.T, kv KeyValue) {
		str := kv.String()
		assert.Assert(t, !strings.Contains(str, "ignoredKey"), "block comment should consume keys and values on its line")
		assert.Assert(t, json.Valid([]byte(str)), "kv.String() is not valid JSON:\n%s", str)
	},
	"app_574720.vdf": func(t *testing.T, kv KeyValue) {
		var vdf struct {
			Ufs struct {
				Rootoverrides map[string]struct {
					Pathtransforms map[string]struct {
						Find    string `json:"find"`
						Replace string `json:"replace"`
					} `json:"pathtransforms"`
				} `json:"rootoverrides"`
			} `json:"ufs"`
		}
		assert.NilError(t, json.Unmarshal([]byte(kv.String()), &vdf))
		assert.Equal(t, vdf.Ufs.Rootoverrides["0"].Pathtransforms["0"].Find, `Mirage Game Studios\\Little Big Workshop\\saves`, improperlyEscapedQuote)
		assert.Equal(t, vdf.Ufs.Rootoverrides["0"].Pathtransforms["0"].Replace, `Mirage Game Studios\\Little Big Workshop\\`, improperlyEscapedQuote)
		assert.Equal(t, vdf.Ufs.Rootoverrides["1"].Pathtransforms["0"].Find, `Mirage Game Studios\\Little Big Workshop\\saves`, improperlyEscapedQuote)
		assert.Equal(t, vdf.Ufs.Rootoverrides["1"].Pathtransforms["0"].Replace, `unity.Mirage Game Studios.Little Big Workshop\\`, improperlyEscapedQuote)
	},
	"app_600760.vdf": func(t *testing.T, kv KeyValue) {
		str := kv.String()
		assert.Assert(t, strings.Contains(str, `"gamedir":null`), "empty value should render as null")
		var vdf struct {
			Extended struct {
				Gamedir string `json:"gamedir"`
			} `json:"extended"`
		}
		assert.NilError(t, json.Unmarshal([]byte(str), &vdf))
	},
	"app_212200.vdf": func(t *testing.T, kv KeyValue) {
		var vdf struct {
			Config struct {
				Launch map[string]struct {
					Arguments string `json:"arguments"`
				} `json:"launch"`
			} `json:"config"`
		}
		assert.NilError(t, json.Unmarshal([]byte(kv.String()), &vdf))
		assert.Equal(t, vdf.Config.Launch["0"].Arguments, `\"Client.exe code:1622 locale:USA env:Regular ver:246 logip:35.162.171.43 logport:11000 chatip:54.214.176.167 chatport:8002 setting:\\\"file://data/features.xml\\\" sn:{tracking_uid} sid:{tracking_sessionid} /P:{passport} -Steam\" --nx:title=Mabinogi --nx:serviceId=880915460`, improperlyEscapedQuote)
	},
	"app_10.vdf": func(t *testing.T, kv KeyValue) {
		assert.Equal(t, len(kv.Children), 6)

		common, ok := kv.GetChild("common")
		assert.Assert(t, ok)

		associations, ok := common.GetChild("associations")
		assert.Assert(t, ok)
		assert.Equal(t, len(associations.Children), 2)
		assert.Equal(t, associations.Children[0].Children[1].Value, "Valve")
		assert.Equal(t, associations.Children[1].Children[1].Value, "Valve")

		kvMap := kv.ToMapOuter()
		appinfo := kvMap["appinfo"].(map[string]any)
		assert.Equal(t, appinfo["appid"], "10")

		commonM := appinfo["common"].(map[string]any)
		assert.Equal(t, len(commonM), 28)

		kvStruct := FromMap(kvMap)
		kv.SortChildren()
		kvStruct.SortChildren()
		assert.DeepEqual(t, kv, kvStruct)
	},
	"app_917720.vdf": func(t *testing.T, kv KeyValue) {
		common, ok := kv.GetChild("common")
		assert.Assert(t, ok)
		assert.Equal(t, len(common.Children), 25)

		assets, ok := common.GetChild("library_assets")
		assert.Assert(t, ok)
		assert.Equal(t, len(assets.Children), 4)
	},
	"package_55058.bin": func(t *testing.T, kv KeyValue) {
		assertChildValue(t, kv, "billingtype", "10")
		assertChildValue(t, kv, "licensetype", "1")

		extended, ok := kv.GetChild("extended")
		assert.Assert(t, ok)
		assert.Equal(t, len(extended.Children), 1)
		assertChildValue(t, extended, "allowcrossregiontradingandgifting", "false")
	},
	"package_228891.bin": func(t *testing.T, kv KeyValue) {
		assertChildValue(t, kv, "billingtype", "3")
		assertChildValue(t, kv, "licensetype", "1")

		extended, ok := kv.GetChild("extended")
		assert.Assert(t, ok)
		assert.Equal(t, len(extended.Children), 2)
		assertChildValue(t, extended, "curatorconnect", "1")
		assertChildValue(t, extended, "releasestateoverride", "released")

		appids, ok := kv.GetChild("appids")
		assert.Assert(t, ok)
		assert.Equal(t, len(appids.Children), 2)
		assertChildValue(t, appids, "0", "680590")
		assertChildValue(t, appids, "1", "1155030")
	},
	"package_394363.bin": func(t *testing.T, kv KeyValue) {
		extended, ok := kv.GetChild("extended")
		assert.Assert(t, ok)
		assert.Equal(t, len(extended.Children), 7)
	},
}

// assertChildValue asserts that kv has a child with the given key and value.
func assertChildValue(t *testing.T, kv KeyValue, key, want string) {
	t.Helper()
	child, ok := kv.GetChild(key)
	assert.Assert(t, ok, "missing child %q", key)
	assert.Equal(t, child.Value, want, "child %q", key)
}

// TestCorpus parses every file under testdata, checks the root key and the
// app or package ID, and then runs the file-specific checks.
func TestCorpus(t *testing.T) {
	files, err := filepath.Glob("testdata/*")
	assert.NilError(t, err)
	assert.Assert(t, len(files) > 0, "no test data found")

	for _, file := range files {
		name := filepath.Base(file)
		t.Run(name, func(t *testing.T) {
			kv, err := ReadFile(file)
			assert.NilError(t, err)

			id := digits.FindString(name)
			if strings.HasPrefix(name, "app_") {
				assert.Equal(t, kv.Key, "appinfo")
				assertChildValue(t, kv, "appid", id)
			} else {
				assert.Equal(t, kv.Key, id)
				assert.Assert(t, len(kv.Children) > 0)
				assertChildValue(t, kv, "packageid", id) //nolint:misspell
			}

			if check, ok := fileChecks[name]; ok {
				check(t, kv)
			}
		})
	}
}

// TestFileChecksTargetExistingFiles guards against a check silently never
// running because its file was renamed or removed.
func TestFileChecksTargetExistingFiles(t *testing.T) {
	for name := range fileChecks {
		_, err := os.Stat(filepath.Join("testdata", name))
		assert.NilError(t, err, name)
	}
}

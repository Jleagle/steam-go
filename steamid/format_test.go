package steamid

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSteam2(t *testing.T) {
	tests := []struct {
		id   ID
		want string
	}{
		{meID, "STEAM_0:0:4180232"},
		{NewID(UniversePublic, AccountTypeIndividual, InstanceDesktop, 4294967295), "STEAM_0:1:2147483647"},
		{NewID(UniverseBeta, AccountTypeIndividual, InstanceDesktop, 3), "STEAM_2:1:1"},
		{NewID(UniversePublic, AccountTypeClan, InstanceAll, 275018), "103582791429796426"},
	}
	for _, tc := range tests {
		if got := tc.id.Steam2(); got != tc.want {
			t.Errorf("Steam2(%d) = %q, want %q", tc.id, got, tc.want)
		}
	}
}

func TestSteam3(t *testing.T) {
	tests := []struct {
		id   ID
		want string
	}{
		{meID, "[U:1:8360464]"},
		{NewID(UniversePublic, AccountTypeIndividual, InstanceConsole, meAccountID), "[U:1:8360464:2]"},
		{NewID(UniversePublic, AccountTypeClan, InstanceAll, 275018), "[g:1:275018]"},
		{NewID(UniversePublic, AccountTypeChat, InstanceAll, 275018), "[T:1:275018]"},
		{NewID(UniversePublic, AccountTypeChat, InstanceID(ChatInstanceClan), 275018), "[c:1:275018]"},
		{NewID(UniversePublic, AccountTypeChat, InstanceID(ChatInstanceLobby), 275018), "[L:1:275018]"},
		{NewID(UniversePublic, AccountTypeAnonGameServer, 456, 123), "[A:1:123:456]"},
		{NewID(UniversePublic, AccountTypeMultiseat, 2, 123), "[M:1:123:2]"},
		{NewID(UniversePublic, AccountTypeGameServer, InstanceConsole, 9), "[G:1:9]"},
		{NewID(UniversePublic, AccountTypeConsoleUser, InstanceAll, 5), "[i:1:5]"},
		{NewID(UniversePublic, AccountTypeInvalid, InstanceAll, 5), "[I:1:5]"},
	}
	for _, tc := range tests {
		if got := tc.id.Steam3(); got != tc.want {
			t.Errorf("Steam3(%d) = %q, want %q", tc.id, got, tc.want)
		}
	}
}

func TestSteam3RoundTrip(t *testing.T) {
	for _, in := range []string{"[U:1:8360464]", "[U:1:8360464:2]", "[g:1:275018]", "[T:1:275018]", "[c:1:275018]", "[L:1:275018]", "[A:1:123:456]", "[M:1:123:2]", "[G:1:9]"} {
		id, err := ParsePlayerID(in)
		if err != nil {
			t.Errorf("ParsePlayerID(%q): %v", in, err)
			continue
		}
		if got := id.Steam3(); got != in {
			t.Errorf("ParsePlayerID(%q).Steam3() = %q", in, got)
		}
	}
}

func TestSteam2RoundTrip(t *testing.T) {
	for _, in := range []string{"STEAM_0:0:4180232", "STEAM_0:1:2147483647", "STEAM_2:1:1"} {
		id, err := ParsePlayerID(in)
		if err != nil {
			t.Errorf("ParsePlayerID(%q): %v", in, err)
			continue
		}
		if got := id.Steam2(); got != in {
			t.Errorf("ParsePlayerID(%q).Steam2() = %q", in, got)
		}
	}

	// The STEAM_1 spelling of the public universe parses to the same ID and
	// renders back in the STEAM_0 form.
	id, err := ParsePlayerID("STEAM_1:0:4180232")
	if err != nil {
		t.Fatalf("ParsePlayerID: %v", err)
	}
	if got := id.Steam2(); got != "STEAM_0:0:4180232" {
		t.Errorf("Steam2 = %q, want STEAM_0:0:4180232", got)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	type doc struct {
		ID  ID            `json:"id"`
		IDs map[ID]string `json:"ids"`
	}
	in := doc{ID: meID, IDs: map[ID]string{meID: "me"}}

	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := `{"id":"76561197968626192","ids":{"76561197968626192":"me"}}`; string(b) != want {
		t.Errorf("Marshal = %s, want %s", b, want)
	}

	var out doc
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if out.ID != meID || out.IDs[meID] != "me" {
		t.Errorf("Unmarshal = %+v, want %+v", out, in)
	}
}

func TestUnmarshalJSONAcceptsNumbersAndNull(t *testing.T) {
	var out struct {
		ID ID `json:"id"`
	}
	if err := json.Unmarshal([]byte(`{"id":76561197968626192}`), &out); err != nil {
		t.Fatalf("Unmarshal number: %v", err)
	}
	if out.ID != meID {
		t.Errorf("Unmarshal number = %d, want %d", out.ID, meID)
	}

	out.ID = meID
	if err := json.Unmarshal([]byte(`{"id":null}`), &out); err != nil {
		t.Fatalf("Unmarshal null: %v", err)
	}
	if out.ID != meID {
		t.Errorf("Unmarshal null changed the ID to %d", out.ID)
	}
}

func TestUnmarshalJSONErrors(t *testing.T) {
	for _, in := range []string{`{"id":"abc"}`, `{"id":true}`, `{"id":"STEAM_0:0:4180232"}`, `{"id":-1}`, `{"id":1.5}`, `{"id":""}`} {
		var out struct {
			ID ID `json:"id"`
		}
		if err := json.Unmarshal([]byte(in), &out); err == nil {
			t.Errorf("Unmarshal(%s) succeeded with %d; want an error", in, out.ID)
		}
	}
}

func TestTextMarshalling(t *testing.T) {
	b, err := meID.MarshalText()
	if err != nil || string(b) != "76561197968626192" {
		t.Errorf("MarshalText = %q, %v", b, err)
	}

	var id ID
	if err := id.UnmarshalText([]byte("76561197968626192")); err != nil || id != meID {
		t.Errorf("UnmarshalText = %d, %v", id, err)
	}

	err = id.UnmarshalText([]byte("x"))
	if err == nil || !strings.HasPrefix(err.Error(), "steamid: ") {
		t.Errorf("UnmarshalText(x) error = %v; want a steamid-prefixed error", err)
	}
}

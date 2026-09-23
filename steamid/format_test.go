package steamid

import "testing"

func TestSteam2(t *testing.T) {
	tests := []struct {
		id   ID
		want string
	}{
		{meID, "STEAM_1:0:4180232"},
		{NewID(UniversePublic, AccountTypeIndividual, InstanceDesktop, 4294967295), "STEAM_1:1:2147483647"},
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
		{NewID(UniversePublic, AccountTypeGameServer, InstanceDesktop, 9), "[G:1:9]"},
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
	for _, in := range []string{"STEAM_1:0:4180232", "STEAM_1:1:2147483647", "STEAM_2:1:1"} {
		id, err := ParsePlayerID(in)
		if err != nil {
			t.Errorf("ParsePlayerID(%q): %v", in, err)
			continue
		}
		if got := id.Steam2(); got != in {
			t.Errorf("ParsePlayerID(%q).Steam2() = %q", in, got)
		}
	}
}

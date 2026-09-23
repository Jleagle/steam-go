package steamid

import (
	"errors"
	"fmt"
	"testing"
)

const (
	meAccountID AccountID = 8360464
	meID        ID        = 76561197968626192
)

func TestNewIDAndGetters(t *testing.T) {
	id := NewID(UniversePublic, AccountTypeIndividual, InstanceDesktop, meAccountID)
	if id != meID {
		t.Fatalf("NewID = %d, want %d", id, meID)
	}
	if got := id.GetUniverseID(); got != UniversePublic {
		t.Errorf("GetUniverseID = %d, want %d", got, UniversePublic)
	}
	if got := id.GetAccountType(); got != AccountTypeIndividual {
		t.Errorf("GetAccountType = %d, want %d", got, AccountTypeIndividual)
	}
	if got := id.GetInstanceID(); got != InstanceDesktop {
		t.Errorf("GetInstanceID = %d, want %d", got, InstanceDesktop)
	}
	if got := id.GetAccountID(); got != meAccountID {
		t.Errorf("GetAccountID = %d, want %d", got, meAccountID)
	}
}

// TestGettersOnValues checks that read methods work on values that have no
// address, such as a function result, and that a plain ID is a fmt.Stringer.
func TestGettersOnValues(t *testing.T) {
	if got := NewID(UniversePublic, AccountTypeIndividual, InstanceDesktop, meAccountID).GetAccountID(); got != meAccountID {
		t.Errorf("GetAccountID on function result = %d, want %d", got, meAccountID)
	}
	var s fmt.Stringer = meID
	if got := s.String(); got != "76561197968626192" {
		t.Errorf("String = %q, want %q", got, "76561197968626192")
	}
	if got := fmt.Sprint(meID); got != "76561197968626192" {
		t.Errorf("Sprint = %q, want %q", got, "76561197968626192")
	}
}

func TestSettersLeaveOtherFieldsAlone(t *testing.T) {
	var id ID
	id.SetAccountID(meAccountID)
	id.SetInstanceID(InstanceDesktop)
	id.SetAccountType(AccountTypeIndividual)
	id.SetUniverseID(UniversePublic)
	if id != meID {
		t.Fatalf("after setters id = %d, want %d", id, meID)
	}

	id.SetInstanceID(InstanceWeb)
	if id.GetInstanceID() != InstanceWeb || id.GetAccountID() != meAccountID || id.GetAccountType() != AccountTypeIndividual || id.GetUniverseID() != UniversePublic {
		t.Errorf("SetInstanceID disturbed other fields: %d", id)
	}
	id.SetAccountType(AccountTypeClan)
	if id.GetAccountType() != AccountTypeClan || id.GetInstanceID() != InstanceWeb || id.GetAccountID() != meAccountID || id.GetUniverseID() != UniversePublic {
		t.Errorf("SetAccountType disturbed other fields: %d", id)
	}
}

func TestClanChatConversion(t *testing.T) {
	clan := NewID(UniversePublic, AccountTypeClan, InstanceAll, 275018)

	chat := clan
	if err := chat.ClanToChat(); err != nil {
		t.Fatalf("ClanToChat: %v", err)
	}
	if chat.GetAccountType() != AccountTypeChat || chat.GetInstanceID() != InstanceID(ChatInstanceClan) || chat.GetAccountID() != 275018 {
		t.Errorf("ClanToChat produced %d (type %d, instance %#x)", chat, chat.GetAccountType(), chat.GetInstanceID())
	}

	back := chat
	if err := back.ChatToClan(); err != nil {
		t.Fatalf("ChatToClan: %v", err)
	}
	if back != clan {
		t.Errorf("ChatToClan = %d, want %d", back, clan)
	}

	if err := clan.ChatToClan(); !errors.Is(err, ErrInvalidChatID) {
		t.Errorf("ChatToClan on a clan ID: error = %v, want ErrInvalidChatID", err)
	}
	if err := chat.ClanToChat(); !errors.Is(err, ErrInvalidClanID) {
		t.Errorf("ClanToChat on a chat ID: error = %v, want ErrInvalidClanID", err)
	}
}

func TestParsePlayerID(t *testing.T) {
	tests := []struct {
		in   string
		want ID
	}{
		{"STEAM_0:0:4180232", meID},
		{"U:1:8360464", meID},
		{"[U:1:8360464]", meID},
		{"U:1:8360464:1", meID},
		{"[U:1:8360464:1]", meID},
		{"8360464", meID},
		{"76561197968626192", meID},
		{"STEAM_1:0:4180232", meID},
		{"  76561197968626192\n", meID},
		{"STEAM_0:1:2147483647", NewID(UniversePublic, AccountTypeIndividual, InstanceDesktop, 4294967295)},
		{"[U:1:8360464:2]", NewID(UniversePublic, AccountTypeIndividual, InstanceConsole, meAccountID)},
		{"[g:1:275018]", NewID(UniversePublic, AccountTypeClan, InstanceAll, 275018)},
		{"[T:1:275018]", NewID(UniversePublic, AccountTypeChat, InstanceAll, 275018)},
		{"[c:1:275018]", NewID(UniversePublic, AccountTypeChat, InstanceID(ChatInstanceClan), 275018)},
		{"[L:1:275018]", NewID(UniversePublic, AccountTypeChat, InstanceID(ChatInstanceLobby), 275018)},
		{"[A:1:123:456]", NewID(UniversePublic, AccountTypeAnonGameServer, 456, 123)},
		{"[A:1:0:456]", NewID(UniversePublic, AccountTypeAnonGameServer, 456, 0)},
		{"[M:1:123:2]", NewID(UniversePublic, AccountTypeMultiseat, 2, 123)},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ParsePlayerID(tc.in)
			if err != nil {
				t.Fatalf("ParsePlayerID(%q): %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("ParsePlayerID(%q) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestParseGroupID(t *testing.T) {
	tests := map[string]ID{
		"275018":             103582791429796426, // Reddit
		"3381077":            103582791432902485, // Counter-Strike: Global Offensive
		"103582791432902485": 103582791432902485,
		" 275018 ":           103582791429796426,
	}
	for in, want := range tests {
		t.Run(in, func(t *testing.T) {
			got, err := ParseGroupID(in)
			if err != nil {
				t.Fatalf("ParseGroupID(%q): %v", in, err)
			}
			if got != want {
				t.Errorf("ParseGroupID(%q) = %d, want %d", in, got, want)
			}
		})
	}
}

func TestParsePlayerIDErrors(t *testing.T) {
	inputs := []string{
		"",
		"abc",
		"STEAM_0:1:2147483648", // account bits overflow 32 bits
		"STEAM_0:0:0",          // account ID zero
		"99999999999",          // 11 digits
		"4294967296",           // 10 digits but above uint32
		"0",
		"[U:1:8360464",      // unmatched bracket
		"U:1:8360464]",      // unmatched bracket
		"[Z:1:1]",           // unknown type character
		"[U:1:0]",           // account ID zero
		"[U:1:1:4]",         // instance above web
		"[U:1:1:1048576]",   // instance above 20 bits
		"12345678901234567", // 17 digits, universe 0
		"76561197968626192x",
		"103582791429796426", // a group ID
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			id, err := ParsePlayerID(in)
			if !errors.Is(err, ErrInvalidPlayerID) {
				t.Errorf("ParsePlayerID(%q) = %d, %v; want an error wrapping ErrInvalidPlayerID", in, id, err)
			}
			if id != 0 {
				t.Errorf("ParsePlayerID(%q) returned %d alongside an error; want 0", in, id)
			}
		})
	}
}

func TestParseGroupIDErrors(t *testing.T) {
	inputs := []string{
		"",
		"abc",
		"0",
		"99999999999",
		"76561197968626192",  // a player ID
		"108086391057166922", // a chat ID
		"103582791429796426x",
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			id, err := ParseGroupID(in)
			if !errors.Is(err, ErrInvalidGroupID) {
				t.Errorf("ParseGroupID(%q) = %d, %v; want an error wrapping ErrInvalidGroupID", in, id, err)
			}
			if id != 0 {
				t.Errorf("ParseGroupID(%q) returned %d alongside an error; want 0", in, id)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want bool
	}{
		{"individual", NewID(UniversePublic, AccountTypeIndividual, InstanceDesktop, 1), true},
		{"individual web instance", NewID(UniversePublic, AccountTypeIndividual, InstanceWeb, 1), true},
		{"individual instance above web", NewID(UniversePublic, AccountTypeIndividual, 4, 1), false},
		{"individual account zero", NewID(UniversePublic, AccountTypeIndividual, InstanceDesktop, 0), false},
		{"universe invalid", NewID(UniverseInvalid, AccountTypeIndividual, InstanceDesktop, 1), false},
		{"universe max", NewID(UniverseMax, AccountTypeIndividual, InstanceDesktop, 1), false},
		{"type invalid", NewID(UniversePublic, AccountTypeInvalid, InstanceDesktop, 1), false},
		{"type max", NewID(UniversePublic, AccountTypeMax, InstanceDesktop, 1), false},
		{"clan", NewID(UniversePublic, AccountTypeClan, InstanceAll, 1), true},
		{"clan with instance", NewID(UniversePublic, AccountTypeClan, InstanceDesktop, 1), false},
		{"clan account zero", NewID(UniversePublic, AccountTypeClan, InstanceAll, 0), false},
		{"game server account zero", NewID(UniversePublic, AccountTypeGameServer, InstanceDesktop, 0), false},
		{"anon game server account zero", NewID(UniversePublic, AccountTypeAnonGameServer, 456, 0), true},
		{"chat with clan flag", NewID(UniversePublic, AccountTypeChat, InstanceID(ChatInstanceClan), 1), true},
		{"zero", 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.id.IsValid(); got != tc.want {
				t.Errorf("IsValid(%d) = %v, want %v", tc.id, got, tc.want)
			}
		})
	}
}

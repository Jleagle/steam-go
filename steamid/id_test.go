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

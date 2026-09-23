package steamid

import (
	"fmt"
	"strconv"
)

// accountTypeToCharacter is the inverse of CharacterToAccountType, used when
// rendering SteamID3 strings. Chat IDs render "c" or "L" from their instance
// flags, and types without a character render "i".
var accountTypeToCharacter = map[AccountType]string{
	AccountTypeInvalid:        "I",
	AccountTypeIndividual:     "U",
	AccountTypeMultiseat:      "M",
	AccountTypeGameServer:     "G",
	AccountTypeAnonGameServer: "A",
	AccountTypePending:        "P",
	AccountTypeContentServer:  "C",
	AccountTypeClan:           "g",
	AccountTypeChat:           "T",
	AccountTypeAnonUser:       "a",
}

// Steam2 renders an individual ID in the "STEAM_X:Y:Z" form: X is the
// universe as stored (1 for public), Y the low bit of the account ID and Z the
// remaining bits. Other account types have no SteamID2 form and render as the
// decimal SteamID64.
func (id ID) Steam2() string {
	if id.GetAccountType() != AccountTypeIndividual {
		return id.String()
	}
	account := uint32(id.GetAccountID())
	return fmt.Sprintf("STEAM_%d:%d:%d", id.GetUniverseID(), account&1, account>>1)
}

// Steam3 renders the ID in the "[U:1:8360464]" form. The instance is appended
// for multiseat and anonymous game server IDs and for individual IDs whose
// instance is not the desktop instance. Chat IDs render as "c" when the clan
// flag is set and "L" when the lobby flag is set.
func (id ID) Steam3() string {
	accountType := id.GetAccountType()
	instance := id.GetInstanceID()

	char, ok := accountTypeToCharacter[accountType]
	if !ok {
		char = "i"
	}
	if accountType == AccountTypeChat {
		switch {
		case instance&InstanceID(ChatInstanceClan) != 0:
			char = "c"
		case instance&InstanceID(ChatInstanceLobby) != 0:
			char = "L"
		}
	}

	withInstance := accountType == AccountTypeAnonGameServer ||
		accountType == AccountTypeMultiseat ||
		(accountType == AccountTypeIndividual && instance != InstanceDesktop)
	if withInstance {
		return fmt.Sprintf("[%s:%d:%d:%d]", char, id.GetUniverseID(), id.GetAccountID(), instance)
	}
	return fmt.Sprintf("[%s:%d:%d]", char, id.GetUniverseID(), id.GetAccountID())
}

// MarshalText implements encoding.TextMarshaler by rendering the decimal
// SteamID64, so encoding/json writes the ID as a string and JavaScript
// clients keep every digit.
func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and accepts the decimal
// SteamID64 that MarshalText writes.
func (id *ID) UnmarshalText(b []byte) error {
	v, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return fmt.Errorf("steamid: %w", err)
	}
	*id = ID(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler. It accepts the decimal SteamID64
// as a JSON string, as MarshalText writes it, or as a JSON number, which is
// how documents encoded before ID implemented encoding.TextMarshaler hold it.
// A JSON null leaves the ID unchanged.
func (id *ID) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	if len(b) >= 2 && b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	return id.UnmarshalText(b)
}

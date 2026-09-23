package steamid

import "fmt"

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

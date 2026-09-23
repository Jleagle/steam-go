// Package steamid parses, validates and formats Steam IDs.
package steamid

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// AccountID is the 32-bit account number within a universe.
type AccountID uint32

// InstanceID is the 20-bit instance field of a Steam ID.
type InstanceID uint32

// Steam instance IDs.
const (
	InstanceAll     InstanceID = 0
	InstanceDesktop InstanceID = 1
	InstanceConsole InstanceID = 2
	InstanceWeb     InstanceID = 3
)

// ChatInstanceID holds the flag bits that a chat ID carries in its instance
// field.
type ChatInstanceID InstanceID

// Steam chat instance flags.
const (
	ChatInstanceClan     ChatInstanceID = 0x100000 >> 1
	ChatInstanceLobby    ChatInstanceID = 0x100000 >> 2
	ChatInstanceMMSLobby ChatInstanceID = 0x100000 >> 3
)

// UniverseID identifies the Steam universe an ID belongs to.
type UniverseID uint8

// Steam universe IDs.
const (
	UniverseInvalid  UniverseID = 0
	UniversePublic   UniverseID = 1
	UniverseBeta     UniverseID = 2
	UniverseInternal UniverseID = 3
	UniverseDev      UniverseID = 4
	UniverseMax      UniverseID = 5
)

// AccountType identifies what kind of entity an ID refers to.
type AccountType uint8

// Steam account types.
const (
	AccountTypeInvalid        AccountType = 0
	AccountTypeIndividual     AccountType = 1
	AccountTypeMultiseat      AccountType = 2
	AccountTypeGameServer     AccountType = 3
	AccountTypeAnonGameServer AccountType = 4
	AccountTypePending        AccountType = 5
	AccountTypeContentServer  AccountType = 6
	AccountTypeClan           AccountType = 7
	AccountTypeChat           AccountType = 8
	AccountTypeConsoleUser    AccountType = 9
	AccountTypeAnonUser       AccountType = 10
	AccountTypeMax            AccountType = 11
)

// CharacterToAccountType maps the account type character of a SteamID3 string
// to its AccountType. "c" and "L" are chat IDs whose instance carries the
// clan or lobby flag; ParsePlayerID sets that flag.
var CharacterToAccountType = map[string]AccountType{
	"I": AccountTypeInvalid,
	"i": AccountTypeInvalid,
	"U": AccountTypeIndividual,
	"M": AccountTypeMultiseat,
	"G": AccountTypeGameServer,
	"A": AccountTypeAnonGameServer,
	"P": AccountTypePending,
	"C": AccountTypeContentServer,
	"g": AccountTypeClan,
	"T": AccountTypeChat,
	"c": AccountTypeChat,
	"L": AccountTypeChat,
	"a": AccountTypeAnonUser,
}

// ID is a 64-bit Steam ID: universe (8 bits), account type (4 bits), instance
// (20 bits) and account ID (32 bits).
type ID uint64

var (
	// ErrInvalidPlayerID is returned when a player ID is invalid.
	ErrInvalidPlayerID = errors.New("invalid player id")
	// ErrInvalidClanID is returned when a clan ID is invalid.
	ErrInvalidClanID = errors.New("invalid clan id")
	// ErrInvalidChatID is returned when a chat ID is invalid.
	ErrInvalidChatID = errors.New("invalid chat id")
	// ErrInvalidGroupID is returned when a group ID is invalid.
	ErrInvalidGroupID = errors.New("invalid group id")
)

var (
	regexpSteam2     = regexp.MustCompile(`^STEAM_([0-5]):([01]):(\d+)$`)          // universe, low bit, account ID >> 1
	regexpSteam3     = regexp.MustCompile(`^([a-zA-Z]):([0-5]):(\d+)(?::(\d+))?$`) // type character, universe, account ID, optional instance
	regexpAccountID  = regexp.MustCompile(`^\d{1,10}$`)
	regexpPlayerID64 = regexp.MustCompile(`^\d{17}$`)
	regexpGroupID64  = regexp.MustCompile(`^\d{18}$`)
)

// ParsePlayerID parses a Steam ID written as SteamID2 ("STEAM_0:0:4180232"),
// SteamID3 ("[U:1:8360464]", brackets optional, instance optional), a bare
// account ID ("8360464") or a SteamID64 ("76561197968626192"). Surrounding
// whitespace is ignored. The result satisfies IsValid; any other input
// returns an error wrapping ErrInvalidPlayerID.
func ParsePlayerID(s string) (ID, error) {
	s = strings.TrimSpace(s)

	id, err := parsePlayerID(s)
	if err != nil {
		return 0, err
	}
	if !id.IsValid() {
		return 0, fmt.Errorf("%w: %q", ErrInvalidPlayerID, s)
	}
	return id, nil
}

func parsePlayerID(s string) (ID, error) {
	if m := regexpSteam2.FindStringSubmatch(s); m != nil {
		universe := UniverseID(m[1][0] - '0')
		if universe == UniverseInvalid {
			universe = UniversePublic // the legacy "STEAM_0:" form means the public universe
		}
		lowBit := uint64(m[2][0] - '0')
		high, err := parseField(m[3], 31, ErrInvalidPlayerID)
		if err != nil {
			return 0, err
		}
		return NewID(universe, AccountTypeIndividual, InstanceDesktop, AccountID(high<<1|lowBit)), nil
	}

	if body, ok := stripBrackets(s); ok {
		if m := regexpSteam3.FindStringSubmatch(body); m != nil {
			return parseSteam3(m)
		}
	}

	if regexpAccountID.MatchString(s) {
		account, err := parseField(s, 32, ErrInvalidPlayerID)
		if err != nil {
			return 0, err
		}
		return NewID(UniversePublic, AccountTypeIndividual, InstanceDesktop, AccountID(account)), nil
	}

	if regexpPlayerID64.MatchString(s) {
		v, err := parseField(s, 64, ErrInvalidPlayerID)
		if err != nil {
			return 0, err
		}
		return ID(v), nil
	}

	return 0, fmt.Errorf("%w: %q", ErrInvalidPlayerID, s)
}

// parseSteam3 builds an ID from the submatches of regexpSteam3: type
// character, universe, account ID and optional instance. Clan, chat, clan
// chat and lobby chat IDs default to no instance; every other type defaults
// to the desktop instance.
func parseSteam3(m []string) (ID, error) {
	char := m[1]
	accountType, ok := CharacterToAccountType[char]
	if !ok {
		return 0, fmt.Errorf("%w: unknown account type character %q", ErrInvalidPlayerID, char)
	}

	universe := UniverseID(m[2][0] - '0')

	account, err := parseField(m[3], 32, ErrInvalidPlayerID)
	if err != nil {
		return 0, err
	}

	instance := InstanceDesktop
	switch char {
	case "g", "T", "c", "L":
		instance = InstanceAll
	}
	if m[4] != "" {
		v, err := parseField(m[4], 20, ErrInvalidPlayerID)
		if err != nil {
			return 0, err
		}
		instance = InstanceID(v)
	}
	switch char {
	case "c":
		instance |= InstanceID(ChatInstanceClan)
	case "L":
		instance |= InstanceID(ChatInstanceLobby)
	}

	return NewID(universe, accountType, instance, AccountID(account)), nil
}

// stripBrackets removes a matched pair of surrounding square brackets. It
// reports false when only one bracket is present.
func stripBrackets(s string) (string, bool) {
	open := strings.HasPrefix(s, "[")
	closed := strings.HasSuffix(s, "]")
	switch {
	case open && closed:
		return s[1 : len(s)-1], true
	case open || closed:
		return "", false
	default:
		return s, true
	}
}

// parseField parses a decimal field that must fit in bitSize bits, wrapping
// any failure in sentinel.
func parseField(s string, bitSize int, sentinel error) (uint64, error) {
	v, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", sentinel, err)
	}
	return v, nil
}

// ParseGroupID parses a Steam group (clan) ID written as a bare account ID
// ("275018") or as a SteamID64 ("103582791429796426"). Surrounding whitespace
// is ignored. The result is a valid clan ID; any other input returns an error
// wrapping ErrInvalidGroupID.
func ParseGroupID(s string) (ID, error) {
	s = strings.TrimSpace(s)

	var id ID
	switch {
	case regexpAccountID.MatchString(s):
		account, err := parseField(s, 32, ErrInvalidGroupID)
		if err != nil {
			return 0, err
		}
		id = NewID(UniversePublic, AccountTypeClan, InstanceAll, AccountID(account))
	case regexpGroupID64.MatchString(s):
		v, err := parseField(s, 64, ErrInvalidGroupID)
		if err != nil {
			return 0, err
		}
		id = ID(v)
	default:
		return 0, fmt.Errorf("%w: %q", ErrInvalidGroupID, s)
	}

	if id.GetAccountType() != AccountTypeClan || !id.IsValid() {
		return 0, fmt.Errorf("%w: %q", ErrInvalidGroupID, s)
	}
	return id, nil
}

// NewID builds an ID from its fields.
func NewID(universe UniverseID, accountType AccountType, instance InstanceID, accountID AccountID) ID {
	var id ID
	id.SetAccountID(accountID)
	id.SetInstanceID(instance)
	id.SetUniverseID(universe)
	id.SetAccountType(accountType)
	return id
}

// IsValid reports whether every field is within Steam's defined range,
// following SteamKit's SteamID.IsValid: the universe and account type are
// defined values, individual, clan and game server IDs have a non-zero
// account ID, individual IDs use the desktop, console or web instance, and
// clan IDs have no instance.
func (id ID) IsValid() bool {
	accountType := id.GetAccountType()
	if accountType <= AccountTypeInvalid || accountType >= AccountTypeMax {
		return false
	}

	universe := id.GetUniverseID()
	if universe <= UniverseInvalid || universe >= UniverseMax {
		return false
	}

	account, instance := id.GetAccountID(), id.GetInstanceID()
	switch accountType {
	case AccountTypeIndividual:
		return account != 0 && instance <= InstanceWeb
	case AccountTypeClan:
		return account != 0 && instance == InstanceAll
	case AccountTypeGameServer:
		return account != 0
	}
	return true
}

// Helpers
func (id ID) get(offset uint, mask uint64) uint64 {
	return (uint64(id) >> offset) & mask
}

func (id *ID) set(offset uint, mask, value uint64) {
	*id = ID((uint64(*id) & ^(mask << offset)) | (value&mask)<<offset)
}

// String returns the decimal SteamID64.
func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

// GetAccountID returns the account ID.
func (id ID) GetAccountID() AccountID {
	//nolint:gosec // account ID is 32-bit.
	return AccountID(id.get(0, 0xFFFFFFFF))
}

// SetAccountID sets the account ID.
func (id *ID) SetAccountID(accountID AccountID) {
	id.set(0, 0xFFFFFFFF, uint64(accountID))
}

// GetInstanceID returns the instance ID.
func (id ID) GetInstanceID() InstanceID {
	//nolint:gosec // instance ID is 20-bit.
	return InstanceID(id.get(32, 0xFFFFF))
}

// SetInstanceID sets the instance ID.
func (id *ID) SetInstanceID(instanceID InstanceID) {
	id.set(32, 0xFFFFF, uint64(instanceID))
}

// GetAccountType returns the account type.
func (id ID) GetAccountType() AccountType {
	//nolint:gosec // account type is 4-bit.
	return AccountType(id.get(52, 0xF))
}

// SetAccountType sets the account type.
func (id *ID) SetAccountType(accountType AccountType) {
	id.set(52, 0xF, uint64(accountType))
}

// GetUniverseID returns the universe ID.
func (id ID) GetUniverseID() UniverseID {
	//nolint:gosec // universe ID is 4-bit.
	return UniverseID(id.get(56, 0xF))
}

// SetUniverseID sets the universe ID.
func (id *ID) SetUniverseID(universeID UniverseID) {
	id.set(56, 0xF, uint64(universeID))
}

// ClanToChat converts a clan ID into the ID of its group chat.
func (id *ID) ClanToChat() error {
	if id.GetAccountType() != AccountTypeClan {
		return ErrInvalidClanID
	}

	id.SetInstanceID(InstanceID(ChatInstanceClan))
	id.SetAccountType(AccountTypeChat)

	return nil
}

// ChatToClan converts a group chat ID into the ID of its clan.
func (id *ID) ChatToClan() error {
	if id.GetAccountType() != AccountTypeChat {
		return ErrInvalidChatID
	}

	id.SetInstanceID(InstanceAll)
	id.SetAccountType(AccountTypeClan)

	return nil
}

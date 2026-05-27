// Package steamid provides tools for parsing and handling Steam IDs.
package steamid

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// AccountID is a unique identifier for a Steam account.
type AccountID uint32

// InstanceID is a unique identifier for a Steam instance.
type InstanceID uint32

// Steam instance IDs.
const (
	InstanceAll     InstanceID = 0
	InstanceDesktop InstanceID = 1
	InstanceConsole InstanceID = 2
	InstanceWeb     InstanceID = 3
)

// ChatInstanceID is a unique identifier for a Steam chat instance.
type ChatInstanceID InstanceID

// Steam chat instance IDs.
const (
	ChatInstanceClan     ChatInstanceID = 0x100000 >> 1
	ChatInstanceLobby    ChatInstanceID = 0x100000 >> 2
	ChatInstanceMMSLobby ChatInstanceID = 0x100000 >> 3
)

// UniverseID is a unique identifier for a Steam universe.
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

// AccountType is a unique identifier for a Steam account type.
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

// CharacterToAccountType maps Steam account type characters to their respective AccountType.
var CharacterToAccountType = map[string]AccountType{
	"I": AccountTypeInvalid,
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

// ID is a unique identifier for a Steam object.
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

// Regexps for parsing Steam IDs.
var (
	RegexpPlayerID1  = regexp.MustCompile(`^STEAM_([0-5]):([01]):(\d+)$`)              // Universe ID, Lowest bit, Highest bit
	RegexpPlayerID3  = regexp.MustCompile(`^\[?([a-zA-Z])\:([0-5])\:(\d+)\]?$`)        // Account type character, Universe ID, Account ID
	RegexpPlayerID3I = regexp.MustCompile(`^\[?([a-zA-Z])\:([0-5])\:(\d+)\:(\d+)\]?$`) // Account type character, Universe ID, Account ID, Instance ID
	RegexpPlayerID32 = regexp.MustCompile(`^\d{1,16}$`)                                // Account ID
	RegexpPlayerID64 = regexp.MustCompile(`^\d{17}$`)                                  // ID
)

// ParsePlayerID parses a Steam ID from a string.
func ParsePlayerID(id string) (out ID, err error) {

	id = strings.TrimSpace(id)

	switch {
	case RegexpPlayerID1.MatchString(id):

		// Get universe
		parts := RegexpPlayerID1.FindStringSubmatch(id)
		i, err := strconv.ParseInt(parts[1], 10, 8)
		if err != nil {
			return out, err
		}
		if i == int64(UniverseInvalid) {
			i = int64(UniversePublic)
		}

		// Get account
		part2, err := strconv.ParseUint(parts[2], 10, 32)
		if err != nil {
			return out, err
		}

		// AccountID
		part3, err := strconv.ParseUint(parts[3], 10, 32)
		if err != nil {
			return out, err
		}

		account := (uint32(part3) << 1) | uint32(part2)

		//nolint:gosec // UniverseID(i) is safe as i is parsed as 8-bit int and checked against constants.
		return NewID(UniverseID(i), AccountTypeIndividual, InstanceDesktop, AccountID(account)), nil

	case RegexpPlayerID3I.MatchString(id):

		parts := RegexpPlayerID3I.FindStringSubmatch(id)

		// Account Type
		if accountType, ok := CharacterToAccountType[parts[1]]; ok {

			// Universe ID
			part2, err := strconv.ParseUint(parts[2], 10, 8)
			if err != nil {
				return out, err
			}

			// Account ID
			part3, err := strconv.ParseUint(parts[3], 10, 32)
			if err != nil {
				return out, err
			}

			// Instance ID
			part4, err := strconv.ParseUint(parts[4], 10, 32)
			if err != nil {
				return out, err
			}

			//
			return NewID(UniverseID(part2), accountType, InstanceID(part4), AccountID(part3)), nil
		}

		return out, ErrInvalidPlayerID

	case RegexpPlayerID3.MatchString(id):

		parts := RegexpPlayerID3.FindStringSubmatch(id)

		// Account Type
		if accountType, ok := CharacterToAccountType[parts[1]]; ok {

			// Universe ID
			part2, err := strconv.ParseUint(parts[2], 10, 8)
			if err != nil {
				return out, err
			}

			// Account ID
			part3, err := strconv.ParseUint(parts[3], 10, 32)
			if err != nil {
				return out, err
			}

			// Instance
			instanceID := InstanceDesktop
			if accountType == AccountTypeClan || accountType == AccountTypeChat {
				instanceID = InstanceAll
			}

			//
			return NewID(UniverseID(part2), accountType, instanceID, AccountID(part3)), nil
		}

		return out, ErrInvalidPlayerID

	case RegexpPlayerID32.MatchString(id):

		i, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return out, err
		}

		return NewID(UniversePublic, AccountTypeIndividual, InstanceDesktop, AccountID(i)), nil

	case RegexpPlayerID64.MatchString(id):

		i, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return out, err
		}

		return ID(i), nil

	default:

		return out, ErrInvalidPlayerID
	}
}

// Regexps for parsing Steam group IDs.
var (
	RegexpGroupID64 = regexp.MustCompile(`^\d{18}$`)   // ID
	RegexpGroupID   = regexp.MustCompile(`^\d{1,17}$`) // Account ID
)

// ParseGroupID parses a Steam group ID from a string.
func ParseGroupID(id string) (out ID, err error) {

	id = strings.TrimSpace(id)

	switch {
	case RegexpGroupID64.MatchString(id):

		i, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return out, err
		}

		return ID(i), nil

	case RegexpGroupID.MatchString(id):

		// Account ID
		i, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return out, err
		}

		return NewID(UniversePublic, AccountTypeClan, InstanceAll, AccountID(i)), nil

	default:

		return out, ErrInvalidGroupID
	}
}

// NewID creates a new Steam ID.
func NewID(universe UniverseID, accountType AccountType, instance InstanceID, accountID AccountID) (id ID) {
	id.SetAccountID(accountID)
	id.SetInstanceID(instance)
	id.SetUniverseID(universe)
	id.SetAccountType(accountType)
	return id
}

// Helpers
func (id *ID) get(offset uint, mask uint64) uint64 {
	return (uint64(*id) >> offset) & mask
}

func (id *ID) set(offset uint, mask, value uint64) {
	*id = ID((uint64(*id) & ^(mask << offset)) | (value&mask)<<offset)
}

func (id *ID) String() string {
	return strconv.FormatUint(uint64(*id), 10)
}

// GetAccountID returns the account ID.
func (id *ID) GetAccountID() AccountID {
	//nolint:gosec // account ID is 32-bit.
	return AccountID(id.get(0, 0xFFFFFFFF))
}

// SetAccountID sets the account ID.
func (id *ID) SetAccountID(accountID AccountID) {
	id.set(0, 0xFFFFFFFF, uint64(accountID))
}

// GetInstanceID returns the instance ID.
func (id *ID) GetInstanceID() InstanceID {
	//nolint:gosec // instance ID is 20-bit.
	return InstanceID(id.get(32, 0xFFFFF))
}

// SetInstanceID sets the instance ID.
func (id *ID) SetInstanceID(instanceID InstanceID) {
	id.set(32, 0xFFFFF, uint64(instanceID))
}

// GetAccountType returns the account type.
func (id *ID) GetAccountType() AccountType {
	//nolint:gosec // account type is 4-bit.
	return AccountType(id.get(52, 0xF))
}

// SetAccountType sets the account type.
func (id *ID) SetAccountType(accountType AccountType) {
	id.set(52, 0xF, uint64(accountType))
}

// GetUniverseID returns the universe ID.
func (id *ID) GetUniverseID() UniverseID {
	//nolint:gosec // universe ID is 4-bit.
	return UniverseID(id.get(56, 0xF))
}

// SetUniverseID sets the universe ID.
func (id *ID) SetUniverseID(universeID UniverseID) {
	id.set(56, 0xF, uint64(universeID))
}

// ClanToChat converts a clan ID to a chat ID.
func (id *ID) ClanToChat() error {

	if id.GetAccountType() != AccountTypeClan {
		return ErrInvalidClanID
	}

	id.SetInstanceID(InstanceID(ChatInstanceClan))
	id.SetAccountType(AccountTypeChat)

	return nil
}

// ChatToClan converts a chat ID to a clan ID.
func (id *ID) ChatToClan() error {

	if id.GetAccountType() != AccountTypeChat {
		return ErrInvalidChatID
	}

	id.SetInstanceID(0)
	id.SetAccountType(AccountTypeClan)

	return nil
}

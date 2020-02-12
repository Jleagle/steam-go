package steamid

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type AccountID uint32

type InstanceID uint32

type ChatInstanceID InstanceID

//noinspection GoUnusedConst
const (
	ChatInstanceClan     ChatInstanceID = 0x100000 >> 1
	ChatInstanceLobby    ChatInstanceID = 0x100000 >> 2
	ChatInstanceMMSLobby ChatInstanceID = 0x100000 >> 3
)

type UniverseID uint8

//noinspection GoUnusedConst
const (
	UniverseInvalid  UniverseID = 0
	UniversePublic   UniverseID = 1
	UniverseBeta     UniverseID = 2
	UniverseInternal UniverseID = 3
	UniverseDev      UniverseID = 4
	UniverseMax      UniverseID = 5
)

type AccountType uint8

//noinspection GoUnusedConst
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

type ID uint64

//noinspection RegExpRedundantEscape
var (
	regexpID1  = regexp.MustCompile(`^STEAM_([0-5]):([01]):(\d+)$`)              // Universe ID, Lowest bit, Highest bit
	regexpID3  = regexp.MustCompile(`^\[?([a-zA-Z])\:([0-5])\:(\d+)\]?$`)        // Account type character, Universe ID, Account ID
	regexpID3I = regexp.MustCompile(`^\[?([a-zA-Z])\:([0-5])\:(\d+)\:(\d+)\]?$`) // Account type character, Universe ID, Account ID, Instance ID
	regexpID32 = regexp.MustCompile(`^\d{1,16}$`)                                // Account ID
	regexpID64 = regexp.MustCompile(`^\d{17}$`)                                  // ID
)

var (
	ErrInvalidPlayerID = errors.New("invalid player id")
	ErrInvalidClanID   = errors.New("invalid clan id")
	ErrInvalidChatID   = errors.New("invalid chat id")
)

func ParsePlayerID(id string) (out ID, err error) {

	id = strings.TrimSpace(id)

	switch {
	case regexpID1.MatchString(id):

		// Get universe
		parts := regexpID1.FindStringSubmatch(id)
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

		//
		return NewID(UniverseID(i), AccountTypeIndividual, 1, AccountID(account)), nil

	case regexpID3I.MatchString(id):

		parts := regexpID3I.FindStringSubmatch(id)

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

	case regexpID3.MatchString(id):

		parts := regexpID3.FindStringSubmatch(id)

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

			//
			return NewID(UniverseID(part2), accountType, 1, AccountID(part3)), nil
		}

		return out, ErrInvalidPlayerID

	case regexpID32.MatchString(id):

		i, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return out, err
		}

		return NewID(UniversePublic, AccountTypeIndividual, 1, AccountID(i)), nil

	case regexpID64.MatchString(id):

		i, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return out, err
		}

		return ID(i), nil

	default:

		return out, ErrInvalidPlayerID
	}
}

func NewID(universe UniverseID, accountType AccountType, instance InstanceID, accountId AccountID) (id ID) {
	id.SetAccountID(accountId)
	id.SetInstanceID(instance)
	id.SetUniverseID(universe)
	id.SetAccountType(accountType)
	return id
}

// Helpers
func (id ID) get(offset uint, mask uint64) uint64 {
	return (uint64(id) >> offset) & mask
}

func (id *ID) set(offset uint, mask, value uint64) {
	*id = ID((uint64(*id) & ^(mask << offset)) | (value&mask)<<offset)
}

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

// Account ID
func (id ID) GetAccountID() AccountID {
	return AccountID(id.get(0, 0xFFFFFFFF))
}

func (id *ID) SetAccountID(accountID AccountID) {
	id.set(0, 0xFFFFFFFF, uint64(accountID))
}

// Instance ID
func (id ID) GetInstanceID() InstanceID {
	return InstanceID(id.get(32, 0xFFFFF))
}

func (id *ID) SetInstanceID(instanceID InstanceID) {
	id.set(32, 0xFFFFF, uint64(instanceID))
}

// Account Type
func (id ID) GetAccountType() AccountType {
	return AccountType(id.get(52, 0xF))
}

func (id *ID) SetAccountType(accountType AccountType) {
	id.set(52, 0xF, uint64(accountType))
}

// Universe ID
func (id ID) GetUniverseID() UniverseID {
	return UniverseID(id.get(56, 0xF))
}

func (id *ID) SetUniverseID(universeID UniverseID) {
	id.set(56, 0xF, uint64(universeID))
}

// Clan / Chat
func (id ID) ClanToChat() error {

	if id.GetAccountType() != AccountTypeClan {
		return ErrInvalidClanID
	}

	id.SetInstanceID(InstanceID(ChatInstanceClan))
	id.SetAccountType(AccountTypeChat)

	return nil
}

func (id ID) ChatToClan() error {

	if id.GetAccountType() != AccountTypeChat {
		return ErrInvalidChatID
	}

	id.SetInstanceID(0)
	id.SetAccountType(AccountTypeClan)

	return nil
}

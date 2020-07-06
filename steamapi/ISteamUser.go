package steamapi

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/Jleagle/unmarshal-go/ctypes"
)

var ErrProfileMissing = errors.New("profile missing")
var ErrProfilePrivate = errors.New("private profile")

func (c Client) GetFriendList(playerID int64) (friends []Friend, b []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))
	options.Set("relationship", "friend")

	b, err = c.getFromAPI("ISteamUser/GetFriendList/v1", options, true)
	if err != nil {
		return friends, b, err
	}

	// Unmarhsal
	var resp FriendListResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return friends, b, err
	}

	return resp.Friendslist.Friends, b, nil
}

type FriendListResponse struct {
	Friendslist FriendsList `json:"friendslist"`
}

type FriendsList struct {
	Friends []Friend `json:"friends"`
}

type Friend struct {
	SteamID      ctypes.Int64 `json:"steamid"`
	Relationship string       `json:"relationship"`
	FriendSince  int64        `json:"friend_since"`
}

//noinspection GoUnusedConst
const (
	VanityURLProfile   = 1
	VanityURLGroup     = 2
	VanityURLGameGroup = 3
)

func (c Client) ResolveVanityURL(vanityURL string, urlType int) (info VanityURL, b []byte, err error) {

	options := url.Values{}
	options.Set("vanityurl", vanityURL)
	options.Set("url_type", strconv.Itoa(urlType))

	b, err = c.getFromAPI("ISteamUser/ResolveVanityURL/v1", options, true)
	if err != nil {
		return info, b, err
	}

	// Unmarhsal
	var resp VanityURLRepsonse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return info, b, err
	}

	if resp.Response.Success != 1 {
		return resp.Response, b, ErrProfileMissing
	}

	return resp.Response, b, nil
}

type VanityURLRepsonse struct {
	Response VanityURL
}

type VanityURL struct {
	SteamID ctypes.Int64 `json:"steamid"`
	Success int8         `json:"success"`
	Message string       `json:"message"`
}

func (c Client) GetPlayer(playerID int64) (player PlayerSummary, b []byte, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.FormatInt(playerID, 10))

	b, err = c.getFromAPI("ISteamUser/GetPlayerSummaries/v2", options, true)
	if err != nil {
		return player, b, err
	}

	// Unmarshal
	var resp PlayerResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return player, b, err
	}

	if len(resp.Response.Players) == 0 {
		return player, b, ErrProfileMissing
	}

	return resp.Response.Players[0], b, nil
}

type PlayerResponse struct {
	Response struct {
		Players []PlayerSummary `json:"players"`
	} `json:"response"`
}

type PlayerSummary struct {
	SteamID                  ctypes.Int64 `json:"steamid"`
	CommunityVisibilityState int          `json:"communityvisibilitystate"`
	ProfileState             int          `json:"profilestate"`
	PersonaName              string       `json:"personaname"`
	LastLogOff               int64        `json:"lastlogoff"`
	CommentPermission        int          `json:"commentpermission"`
	ProfileURL               string       `json:"profileurl"`
	Avatar                   string       `json:"avatar"`
	AvatarMedium             string       `json:"avatarmedium"`
	AvatarFull               string       `json:"avatarfull"`
	AvatarHash               string       `json:"avatarhash"`
	PersonaState             int          `json:"personastate"`
	RealName                 string       `json:"realname"`
	PrimaryClanID            string       `json:"primaryclanid"`
	TimeCreated              int64        `json:"timecreated"`
	PersonaStateFlags        int          `json:"personastateflags"`
	CountryCode              string       `json:"loccountrycode"`
	StateCode                string       `json:"locstatecode"`
}

func (c Client) GetPlayerBans(playerID int64) (bans GetPlayerBanResponse, b []byte, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.FormatInt(playerID, 10))

	b, err = c.getFromAPI("ISteamUser/GetPlayerBans/v1", options, true)
	if err != nil {
		return bans, b, err
	}

	// Unmarshal
	var resp GetPlayerBansResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return bans, b, err
	}

	if len(resp.Players) == 0 {
		return bans, b, ErrProfileMissing
	}

	return resp.Players[0], b, nil
}

type GetPlayerBansResponse struct {
	Players []GetPlayerBanResponse `json:"players"`
}

type GetPlayerBanResponse struct {
	SteamID          ctypes.Int64 `json:"SteamId"`
	CommunityBanned  bool         `json:"CommunityBanned"`
	VACBanned        bool         `json:"VACBanned"`
	NumberOfVACBans  int          `json:"NumberOfVACBans"`
	DaysSinceLastBan int          `json:"DaysSinceLastBan"`
	NumberOfGameBans int          `json:"NumberOfGameBans"`
	EconomyBan       string       `json:"EconomyBan"`
}

func (c Client) GetUserGroupList(playerID int64) (groups UserGroupList, b []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))

	b, err = c.getFromAPI("ISteamUser/GetUserGroupList/v1", options, true)
	// err checked below

	var resp UserGroupListResponse
	err2 := json.Unmarshal(b, &resp)
	if err2 != nil {
		return groups, b, err2
	}

	if !resp.Response.Success && strings.HasPrefix(resp.Response.Error, "Failed to get information about account") {
		return groups, b, ErrProfileMissing
	}

	if !resp.Response.Success && resp.Response.Error == "Private profile" {
		return groups, b, ErrProfilePrivate
	}

	if err != nil {
		return groups, b, err
	}

	return resp.Response, b, nil
}

type UserGroupListResponse struct {
	Response UserGroupList `json:"response"`
}

type UserGroupList struct {
	Success bool `json:"success"`
	Groups  []struct {
		GID string `json:"gid"` // Can be over 64 bit
	} `json:"groups"`
	Error string `json:"error"`
}

func (u UserGroupList) GetIDs() (ids []string) {
	for _, v := range u.Groups {
		ids = append(ids, v.GID)
	}
	return ids
}

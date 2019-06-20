package steam

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/Jleagle/unmarshal-go/ctypes"
)

var (
	ErrNoUserFound = errors.New("no user found")
)

func (s Steam) GetFriendList(playerID int64) (friends []Friend, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))
	options.Set("relationship", "friend")

	bytes, err = s.getFromAPI("ISteamUser/GetFriendList/v1", options)
	if err != nil {
		return friends, bytes, err
	}

	// Unmarhsal
	var resp FriendListResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return friends, bytes, err
	}

	return resp.Friendslist.Friends, bytes, nil
}

type FriendListResponse struct {
	Friendslist FriendsList `json:"friendslist"`
}

type FriendsList struct {
	Friends []Friend `json:"friends"`
}

type Friend struct {
	SteamID      ctypes.CInt64 `json:"steamid"`
	Relationship string        `json:"relationship"`
	FriendSince  int64         `json:"friend_since"`
}

const (
	VanityURLProfile   = 1
	VanityURLGroup     = 2
	VanityURLGameGroup = 3
)

func (s Steam) ResolveVanityURL(vanityURL string, urlType int) (info VanityURL, bytes []byte, err error) {

	options := url.Values{}
	options.Set("vanityurl", vanityURL)
	options.Set("url_type", strconv.Itoa(urlType))

	bytes, err = s.getFromAPI("ISteamUser/ResolveVanityURL/v1", options)
	if err != nil {
		return info, bytes, err
	}

	// Unmarhsal
	var resp VanityURLRepsonse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return info, bytes, err
	}

	if resp.Response.Success != 1 {
		return resp.Response, bytes, ErrNoUserFound
	}

	return resp.Response, bytes, nil
}

type VanityURLRepsonse struct {
	Response VanityURL
}

type VanityURL struct {
	SteamID ctypes.CInt64 `json:"steamid"`
	Success int8          `json:"success"`
	Message string        `json:"message"`
}

func (s Steam) GetPlayer(playerID int64) (player PlayerSummary, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.FormatInt(playerID, 10))

	bytes, err = s.getFromAPI("ISteamUser/GetPlayerSummaries/v2", options)
	if err != nil {
		return player, bytes, err
	}

	// Unmarshal
	var resp PlayerResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return player, bytes, err
	}

	if len(resp.Response.Players) == 0 {
		return player, bytes, ErrNoUserFound
	}

	return resp.Response.Players[0], bytes, nil
}

type PlayerResponse struct {
	Response struct {
		Players []PlayerSummary `json:"players"`
	} `json:"response"`
}

type PlayerSummary struct {
	SteamID                  ctypes.CInt64 `json:"steamid"`
	CommunityVisibilityState int           `json:"communityvisibilitystate"`
	ProfileState             int           `json:"profilestate"`
	PersonaName              string        `json:"personaname"`
	LastLogOff               int64         `json:"lastlogoff"`
	CommentPermission        int           `json:"commentpermission"`
	ProfileURL               string        `json:"profileurl"`
	Avatar                   string        `json:"avatar"`
	AvatarMedium             string        `json:"avatarmedium"`
	AvatarFull               string        `json:"avatarfull"`
	PersonaState             int           `json:"personastate"`
	RealName                 string        `json:"realname"`
	PrimaryClanID            ctypes.CInt   `json:"primaryclanid"`
	TimeCreated              int64         `json:"timecreated"`
	PersonaStateFlags        int           `json:"personastateflags"`
	LOCCountryCode           string        `json:"loccountrycode"`
	LOCStateCode             string        `json:"locstatecode"`
}

func (s Steam) GetPlayerBans(playerID int64) (bans GetPlayerBanResponse, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.FormatInt(playerID, 10))

	bytes, err = s.getFromAPI("ISteamUser/GetPlayerBans/v1", options)
	if err != nil {
		return bans, bytes, err
	}

	// Unmarshal
	var resp GetPlayerBansResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return bans, bytes, err
	}

	if len(resp.Players) == 0 {
		return bans, bytes, ErrNoUserFound
	}

	return resp.Players[0], bytes, nil
}

type GetPlayerBansResponse struct {
	Players []GetPlayerBanResponse `json:"players"`
}

type GetPlayerBanResponse struct {
	SteamID          ctypes.CInt64 `json:"SteamId"`
	CommunityBanned  bool          `json:"CommunityBanned"`
	VACBanned        bool          `json:"VACBanned"`
	NumberOfVACBans  int           `json:"NumberOfVACBans"`
	DaysSinceLastBan int           `json:"DaysSinceLastBan"`
	NumberOfGameBans int           `json:"NumberOfGameBans"`
	EconomyBan       string        `json:"EconomyBan"`
}

func (s Steam) GetUserGroupList(playerID int64) (groups UserGroupList, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))

	bytes, err = s.getFromAPI("ISteamUser/GetUserGroupList/v1", options)
	if err != nil {
		return groups, bytes, err
	}

	// Unmarshal
	var resp UserGroupListResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return groups, bytes, err
	}

	if len(resp.Response.Groups) == 0 {
		return resp.Response, bytes, nil
	}

	return resp.Response, bytes, nil
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

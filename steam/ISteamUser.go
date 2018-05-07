package steam

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
)

var (
	ErrNoUserFound = errors.New("no user found")
)

func (s Steam) GetFriendList(id int) (friends FriendsList, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))
	options.Set("relationship", "friend")

	bytes, err = s.getFromAPI("ISteamUser/GetFriendList/v1/", options)
	if err != nil {
		return friends, bytes, err
	}

	// Regex
	var regex *regexp.Regexp
	var str = string(bytes)

	regex = regexp.MustCompile(`"steamid":\s?"(\d+)"`)
	str = regex.ReplaceAllString(str, `"steamid": $1`)

	bytes = []byte(str)

	// Unmarhsal
	var resp FriendListResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return friends, bytes, err
	}

	return resp.Friendslist, bytes, nil
}

type FriendListResponse struct {
	Friendslist FriendsList `json:"friendslist"`
}

type FriendsList struct {
	Friends []Friend `json:"friends"`
}

type Friend struct {
	SteamID      int64  `json:"steamid"`
	Relationship string `json:"relationship"`
	FriendSince  int64  `json:"friend_since"`
}

func (s Steam) ResolveVanityURL(id string) (info VanityURL, bytes []byte, err error) {

	options := url.Values{}
	options.Set("vanityurl", id)
	options.Set("url_type", "1") // 1 (default): Individual profile, 2: Group, 3: Official game group

	bytes, err = s.getFromAPI("ISteamUser/ResolveVanityURL/v1/", options)
	if err != nil {
		return info, bytes, err
	}

	// Regex
	var regex *regexp.Regexp
	var str = string(bytes)

	regex = regexp.MustCompile(`"steamid":\s?"(\d+)"`)
	str = regex.ReplaceAllString(str, `"steamid": $1`)

	bytes = []byte(str)

	// Unmarhsal
	var resp VanityURLRepsonse
	if err := json.Unmarshal(bytes, &resp); err != nil {
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
	SteamID int64  `json:"steamid"`
	Success int8   `json:"success"`
	Message string `json:"message"`
}

func (s Steam) GetPlayer(id int) (player PlayerSummary, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.Itoa(id))

	bytes, err = s.getFromAPI("ISteamUser/GetPlayerSummaries/v2/", options)
	if err != nil {
		return player, bytes, err
	}

	// Regex
	var regex *regexp.Regexp
	var str = string(bytes)

	regex = regexp.MustCompile(`"primaryclanid":\s?"(\d+)"`)
	str = regex.ReplaceAllString(str, `"primaryclanid": $1`)

	regex = regexp.MustCompile(`"steamid":\s?"(\d+)"`)
	str = regex.ReplaceAllString(str, `"steamid": $1`)

	bytes = []byte(str)

	// Unmarshal
	var resp PlayerResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return player, bytes, err
	}

	if len(resp.Response.Players) == 0 {
		return player, bytes, ErrNoUserFound
	}

	return resp.Response.Players[0], bytes, nil
}

type PlayerResponse struct {
	Response PlayerSummaries `json:"response"`
}

type PlayerSummaries struct {
	Players []PlayerSummary `json:"players"`
}

type PlayerSummary struct {
	SteamID                  int64  `json:"steamid"`
	CommunityVisibilityState int    `json:"communityvisibilitystate"`
	ProfileState             int    `json:"profilestate"`
	PersonaName              string `json:"personaname"`
	LastLogOff               int64  `json:"lastlogoff"`
	CommentPermission        int    `json:"commentpermission"`
	ProfileURL               string `json:"profileurl"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarmedium"`
	AvatarFull               string `json:"avatarfull"`
	PersonaState             int    `json:"personastate"`
	RealName                 string `json:"realname"`
	PrimaryClanID            int    `json:"primaryclanid"`
	TimeCreated              int64  `json:"timecreated"`
	PersonaStateFlags        int    `json:"personastateflags"`
	LOCCountryCode           string `json:"loccountrycode"`
	LOCStateCode             string `json:"locstatecode"`
}

func (s Steam) GetPlayerBans(id int) (bans GetPlayerBanResponse, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.Itoa(id))

	bytes, err = s.getFromAPI("ISteamUser/GetPlayerBans/v1", options)
	if err != nil {
		return bans, bytes, err
	}

	// Regex
	var regex *regexp.Regexp
	var str = string(bytes)

	regex = regexp.MustCompile(`"SteamId":\s?"(\d+)"`)
	str = regex.ReplaceAllString(str, `"SteamId": $1`)

	bytes = []byte(str)

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
	SteamID          int64  `json:"SteamId"`
	CommunityBanned  bool   `json:"CommunityBanned"`
	VACBanned        bool   `json:"VACBanned"`
	NumberOfVACBans  int    `json:"NumberOfVACBans"`
	DaysSinceLastBan int    `json:"DaysSinceLastBan"`
	NumberOfGameBans int    `json:"NumberOfGameBans"`
	EconomyBan       string `json:"EconomyBan"`
}

func (s Steam) GetUserGroupList(id int) (groups UserGroupList, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))

	bytes, err = s.getFromAPI("ISteamUser/GetUserGroupList/v1", options)
	if err != nil {
		return groups, bytes, err
	}

	// Regex
	var regex *regexp.Regexp
	var str = string(bytes)

	regex = regexp.MustCompile(`"gid":\s?"(\d+)"`)
	str = regex.ReplaceAllString(str, `"gid": $1`)

	bytes = []byte(str)

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
	Success bool        `json:"success"`
	Groups  []UserGroup `json:"groups"`
}

func (u UserGroupList) GetIDs() (ids []int) {
	for _, v := range u.Groups {
		ids = append(ids, v.GID)
	}
	return ids
}

type UserGroup struct {
	GID int `json:"gid"`
}

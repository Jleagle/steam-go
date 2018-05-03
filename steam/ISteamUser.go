package steam

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNoUserFound = errors.New("no user found")
)

func GetFriendList(id int) (friends FriendsList, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))
	options.Set("relationship", "friend")

	bytes, err = get("ISteamUser/GetFriendList/v1/", options)
	if err != nil {
		return friends, bytes, err
	}

	if strings.Contains(string(bytes), "Internal Server Error") {
		return friends, bytes, ErrNoUserFound
	}

	// Unmarshal JSON
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
	SteamID      string `json:"steamid"`
	Relationship string `json:"relationship"`
	FriendSince  int    `json:"friend_since"`
}

func ResolveVanityURL(id string) (info VanityURL, bytes []byte, err error) {

	options := url.Values{}
	options.Set("vanityurl", id)
	options.Set("url_type", "1") // 1 (default): Individual profile, 2: Group, 3: Official game group

	bytes, err = get("ISteamUser/ResolveVanityURL/v1/", options)
	if err != nil {
		return info, bytes, err
	}

	// Unmarshal JSON
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
	SteamID string `json:"steamid"`
	Success int8   `json:"success"`
	Message string `json:"message"`
}

func GetPlayer(id int) (player PlayerSummary, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.Itoa(id))

	bytes, err = get("ISteamUser/GetPlayerSummaries/v2/", options)
	if err != nil {
		return player, bytes, err
	}

	var regex *regexp.Regexp
	var s = string(bytes)

	// Convert strings to ints
	regex = regexp.MustCompile(`"primaryclanid":\s?"(\d+)"`)
	s = regex.ReplaceAllString(s, `"primaryclanid": $1`)

	regex = regexp.MustCompile(`"steamid":\s?"(\d+)"`)
	s = regex.ReplaceAllString(s, `"steamid": $1`)

	bytes = []byte(s)

	// Unmarshal JSON
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
	SteamID                  int    `json:"steamid"`
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

func GetPlayerBans(id int) (bans GetPlayerBanResponse, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.Itoa(id))

	bytes, err = get("ISteamUser/GetPlayerBans/v1", options)
	if err != nil {
		return bans, bytes, err
	}

	// Unmarshal JSON
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
	SteamID          string `json:"SteamId"`
	CommunityBanned  bool   `json:"CommunityBanned"`
	VACBanned        bool   `json:"VACBanned"`
	NumberOfVACBans  int    `json:"NumberOfVACBans"`
	DaysSinceLastBan int    `json:"DaysSinceLastBan"`
	NumberOfGameBans int    `json:"NumberOfGameBans"`
	EconomyBan       string `json:"EconomyBan"`
}

// todo, regex the ids into ints
func GetUserGroupList(id int) (groups UserGroupList, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))

	bytes, err = get("ISteamUser/GetUserGroupList/v1", options)
	if err != nil {
		return groups, bytes, err
	}

	// Unmarshal JSON
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
		gidx, _ := strconv.Atoi(v.GID)
		ids = append(ids, gidx)
	}
	return ids
}

type UserGroup struct {
	GID string `json:"gid"`
}

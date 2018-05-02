package steam

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/steam-authority/steam-authority/logger"
)

var (
	ErrNoUserFound = errors.New("no user found")
)

func GetFriendList(id int) (friends []GetFriendListFriend, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))
	options.Set("relationship", "friend")

	bytes, err := get("ISteamUser/GetFriendList/v1/", options)
	if err != nil {
		return friends, err
	}

	if strings.Contains(string(bytes), "Internal Server Error") {
		return friends, ErrNoUserFound
	}

	// Unmarshal JSON
	var response = new(GetFriendListBody)
	if err := json.Unmarshal(bytes, &response); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return friends, err
	}

	return response.Friendslist.Friends, nil
}

type GetFriendListBody struct {
	Friendslist struct {
		Friends []GetFriendListFriend `json:"friends"`
	} `json:"friendslist"`
}

type GetFriendListFriend struct {
	SteamID      string `json:"steamid"`
	Relationship string `json:"relationship"`
	FriendSince  int    `json:"friend_since"`
}

func ResolveVanityURL(id string) (resp ResolveVanityURLBody, err error) {

	options := url.Values{}
	options.Set("vanityurl", id)
	options.Set("url_type", "1") // 1 (default): Individual profile, 2: Group, 3: Official game group

	bytes, err := get("ISteamUser/ResolveVanityURL/v1/", options)
	if err != nil {
		return resp, err
	}

	// Unmarshal JSON
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return resp, err
	}

	if resp.Response.Success != 1 {
		return resp, ErrNoUserFound
	}

	return resp, nil
}

type ResolveVanityURLBody struct {
	Response ResolveVanityURLResponse
}

type ResolveVanityURLResponse struct {
	SteamID string `json:"steamid"`
	Success int8   `json:"success"`
	Message string `json:"message"`
}

func GetPlayerSummaries(id int) (player PlayerSummary, err error) {

	idString := strconv.Itoa(id)

	options := url.Values{}
	options.Set("steamids", idString)

	bytes, err := get("ISteamUser/GetPlayerSummaries/v2/", options)
	if err != nil {
		return player, err
	}

	// Fix values that can change type, causing unmarshal errors
	var regex *regexp.Regexp
	var b = string(bytes)

	// Convert strings to ints
	regex = regexp.MustCompile(`"primaryclanid":\s?"(\d+)"`)
	b = regex.ReplaceAllString(b, `"primaryclanid": $1`)

	regex = regexp.MustCompile(`"steamid":\s?"(\d+)"`)
	b = regex.ReplaceAllString(b, `"steamid": $1`)

	bytes = []byte(b)

	// Unmarshal JSON
	var resp PlayerSummariesResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return player, err
	}

	if len(resp.Response.Players) == 0 {
		return player, errors.New("not found in steam: " + idString)
	}

	return resp.Response.Players[0], nil
}

type PlayerSummariesResponse struct {
	Response struct {
		Players []PlayerSummary `json:"players"`
	} `json:"response"`
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

func GetPlayerBans(id int) (bans GetPlayerBanResponse, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.Itoa(id))

	bytes, err := get("ISteamUser/GetPlayerBans/v1", options)
	if err != nil {
		return bans, err
	}

	// Unmarshal JSON
	var resp GetPlayerBansResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return bans, err
	}

	if len(resp.Players) == 0 {
		return bans, nil
	}

	return resp.Players[0], nil
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

func GetUserGroupList(id int) (groups []int, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))

	bytes, err := get("ISteamUser/GetUserGroupList/v1", options)
	if err != nil {
		return groups, err
	}

	// Unmarshal JSON
	var resp UserGroupListResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return groups, err
	}

	if len(resp.Response.Groups) == 0 {
		return groups, nil
	}

	for _, v := range resp.Response.Groups {
		gidx, _ := strconv.Atoi(v.GID)
		groups = append(groups, gidx)
	}

	return groups, nil
}

type UserGroupListResponse struct {
	Response struct {
		Success bool `json:"success"`
		Groups []struct {
			GID string `json:"gid"`
		} `json:"groups"`
	} `json:"response"`
}

package steamapi

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Jleagle/unmarshal-go"
)

// Gets information about a player's recently played games
func (c Client) GetRecentlyPlayedGames(playerID int64) (games []RecentlyPlayedGame, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))
	options.Set("count", "0")

	b, err := c.getFromAPI("IPlayerService/GetRecentlyPlayedGames/v1", options, true)
	if err != nil {
		return games, err
	}

	var resp RecentlyPlayedGamesResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return games, err
	}

	return resp.Response.Games, nil
}

type RecentlyPlayedGamesResponse struct {
	Response RecentlyPlayedGames `json:"response"`
}

type RecentlyPlayedGames struct {
	TotalCount int                  `json:"total_count"`
	Games      []RecentlyPlayedGame `json:"games"`
}

type RecentlyPlayedGame struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	PlayTime2Weeks  int    `json:"playtime_2weeks"`
	PlayTimeForever int    `json:"playtime_forever"`
	ImgIconURL      string `json:"img_icon_url"`
	ImgLogoURL      string `json:"img_logo_url"`
	WindowsForever  int    `json:"playtime_windows_forever"`
	MacForever      int    `json:"playtime_mac_forever"`
	LinuxForever    int    `json:"playtime_linux_forever"`
}

// Return a list of games owned by the player
func (c Client) GetOwnedGames(playerID int64) (games OwnedGames, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))
	options.Set("include_appinfo", "1")
	options.Set("include_played_free_games", "1")

	b, err := c.getFromAPI("IPlayerService/GetOwnedGames/v1", options, true)
	if err != nil {
		return games, err
	}

	var resp OwnedGamesResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return games, err
	}

	return resp.Response, nil
}

type OwnedGamesResponse struct {
	Response OwnedGames `json:"response"`
}

type OwnedGames struct {
	GameCount int `json:"game_count"`
	Games     []struct {
		AppID                    int    `json:"appid"`
		Name                     string `json:"name"`
		PlaytimeForever          int    `json:"playtime_forever"`
		PlaytimeWindows          int    `json:"playtime_windows_forever"`
		PlaytimeMac              int    `json:"playtime_mac_forever"`
		PlaytimeLinux            int    `json:"playtime_linux_forever"`
		ImgIconURL               string `json:"img_icon_url"`
		ImgLogoURL               string `json:"img_logo_url"`
		HasCommunityVisibleStats bool   `json:"has_community_visible_stats"`
	} `json:"games"`
}

// Returns the Steam Level of a user
func (c Client) GetSteamLevel(playerID int64) (level int, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))

	b, err := c.getFromAPI("IPlayerService/GetSteamLevel/v1", options, true)
	if err != nil {
		return level, err
	}

	var resp LevelResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return level, err
	}

	return resp.Response.PlayerLevel, nil
}

type LevelResponse struct {
	Response struct {
		PlayerLevel int `json:"player_level"`
	} `json:"response"`
}

// Gets badges that are owned by a specific user
func (c Client) GetBadges(playerID int64) (badges BadgesInfo, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))

	b, err := c.getFromAPI("IPlayerService/GetBadges/v1", options, true)
	if err != nil {
		return badges, err
	}

	var resp BadgesResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return badges, err
	}

	return resp.Response, nil
}

type BadgesResponse struct {
	Response BadgesInfo `json:"response"`
}

type BadgesInfo struct {
	Badges                     []BadgeResponse `json:"badges"`
	PlayerXP                   int             `json:"player_xp"`
	PlayerLevel                int             `json:"player_level"`
	PlayerXPNeededToLevelUp    int             `json:"player_xp_needed_to_level_up"`
	PlayerXPNeededCurrentLevel int             `json:"player_xp_needed_current_level"`
}

func (b BadgesInfo) GetPercentOfLevel() int {

	if b.PlayerXP == 0 {
		return 0
	}

	start := b.PlayerXPNeededCurrentLevel
	finish := b.PlayerXPNeededToLevelUp + b.PlayerXP
	levelRange := finish - start
	intoLevel := b.PlayerXP - b.PlayerXPNeededCurrentLevel

	return int((float64(intoLevel) / float64(levelRange)) * 100)
}

type BadgeResponse struct {
	AppID           int             `json:"appid"`
	BadgeID         int             `json:"badgeid"`
	BorderColor     unmarshal.Bool  `json:"border_color"`
	CommunityItemID unmarshal.Int64 `json:"communityitemid"`
	CompletionTime  int64           `json:"completion_time"`
	Level           int             `json:"level"`
	Scarcity        int             `json:"scarcity"`
	XP              int             `json:"xp"`
}

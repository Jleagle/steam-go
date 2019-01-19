package steam

import (
	"net/url"
	"strconv"

	"github.com/Jleagle/unmarshal-go/unmarshal"
)

// Gets information about a player's recently played games
func (s Steam) GetRecentlyPlayedGames(playerID int64) (games RecentlyPlayedGames, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))
	options.Set("count", "0")

	bytes, err = s.getFromAPI("IPlayerService/GetRecentlyPlayedGames/v1", options)
	if err != nil {
		return games, bytes, err
	}

	var resp RecentlyPlayedGamesResponse
	err = unmarshal.Unmarshal(bytes, &resp)
	if err != nil {
		return games, bytes, err
	}

	return resp.Response, bytes, nil
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
}

// Return a list of games owned by the player
func (s Steam) GetOwnedGames(playerID int64) (games OwnedGames, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))
	options.Set("include_appinfo", "1")
	options.Set("include_played_free_games", "1")

	bytes, err = s.getFromAPI("IPlayerService/GetOwnedGames/v1", options)
	if err != nil {
		return games, bytes, err
	}

	var resp OwnedGamesResponse
	err = unmarshal.Unmarshal(bytes, &resp)
	if err != nil {
		return games, bytes, err
	}

	return resp.Response, bytes, nil
}

type OwnedGamesResponse struct {
	Response OwnedGames `json:"response"`
}

type OwnedGames struct {
	GameCount int         `json:"game_count"`
	Games     []OwnedGame `json:"games"`
}

type OwnedGame struct {
	AppID                    int    `json:"appid"`
	Name                     string `json:"name"`
	PlaytimeForever          int    `json:"playtime_forever"`
	ImgIconURL               string `json:"img_icon_url"`
	ImgLogoURL               string `json:"img_logo_url"`
	HasCommunityVisibleStats bool   `json:"has_community_visible_stats,omitempty"`
}

// Returns the Steam Level of a user
func (s Steam) GetSteamLevel(playerID int64) (level int, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))

	bytes, err = s.getFromAPI("IPlayerService/GetSteamLevel/v1", options)
	if err != nil {
		return level, bytes, err
	}

	var resp LevelResponse
	err = unmarshal.Unmarshal(bytes, &resp)
	if err != nil {
		return level, bytes, err
	}

	return resp.Response.PlayerLevel, bytes, nil
}

type LevelResponse struct {
	Response struct {
		PlayerLevel int `json:"player_level"`
	} `json:"response"`
}

// Gets badges that are owned by a specific user
func (s Steam) GetBadges(playerID int64) (badges BadgesInfo, bytes []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatInt(playerID, 10))

	bytes, err = s.getFromAPI("IPlayerService/GetBadges/v1", options)
	if err != nil {
		return badges, bytes, err
	}

	var resp BadgesResponse
	err = unmarshal.Unmarshal(bytes, &resp)
	if err != nil {
		return badges, bytes, err
	}

	return resp.Response, bytes, nil
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
	BadgeID        int   `json:"badgeid"`
	Level          int   `json:"level"`
	CompletionTime int64 `json:"completion_time"`
	XP             int   `json:"xp"`
	Scarcity       int   `json:"scarcity"`
	AppID          int   `json:"appid"`
	BorderColor    int   `json:"border_color"`
}

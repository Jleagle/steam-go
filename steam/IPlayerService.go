package steam

import (
	"encoding/json"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/steam-authority/steam-authority/helpers"
	"github.com/steam-authority/steam-authority/logger"
)

// Gets information about a player's recently played games
func GetRecentlyPlayedGames(playerID int) (games []RecentlyPlayedGame, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(playerID))
	options.Set("count", "0")

	bytes, err := get("IPlayerService/GetRecentlyPlayedGames/v1", options)
	if err != nil {
		return games, err
	}

	// Unmarshal JSON
	var resp RecentlyPlayedGamesResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return games, err
	}

	return resp.Response.Games, nil
}

type RecentlyPlayedGamesResponse struct {
	Response struct {
		TotalCount int                  `json:"total_count"`
		Games      []RecentlyPlayedGame `json:"games"`
	} `json:"response"`
}

type RecentlyPlayedGame struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	Playtime2Weeks  int    `json:"playtime_2weeks"`
	PlaytimeForever int    `json:"playtime_forever"`
	ImgIconURL      string `json:"img_icon_url"`
	ImgLogoURL      string `json:"img_logo_url"`
}

func (r RecentlyPlayedGame) GetTwoWeeksNice() string {
	return helpers.GetTimeShort(r.Playtime2Weeks, 2)
}

func (r RecentlyPlayedGame) GetAllTimeNice() string {
	return helpers.GetTimeShort(r.PlaytimeForever, 2)
}

// Return a list of games owned by the player
func GetOwnedGames(id int) (games []OwnedGame, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))
	options.Set("include_appinfo", "1")
	options.Set("include_played_free_games", "1")

	bytes, err := get("IPlayerService/GetOwnedGames/v1", options)
	if err != nil {
		return games, err
	}

	// Unmarshal JSON
	var resp OwnedGamesResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return games, err
	}

	return resp.Response.Games, nil
}

type OwnedGamesResponse struct {
	Response struct {
		GameCount int         `json:"game_count"`
		Games     []OwnedGame `json:"games"`
	} `json:"response"`
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
func GetSteamLevel(id int) (level int, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))

	bytes, err := get("IPlayerService/GetSteamLevel/v1", options)
	if err != nil {
		return level, err
	}

	// Unmarshal JSON
	var resp LevelResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
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
func GetBadges(id int) (badges BadgesResponse, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))

	bytes, err := get("IPlayerService/GetBadges/v1", options)
	if err != nil {
		return badges, err
	}

	// Unmarshal JSON
	var resp BadgesResponseOuter
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return badges, err
	}

	return resp.Response, nil
}

type BadgesResponseOuter struct {
	Response BadgesResponse `json:"response"`
}

type BadgesResponse struct {
	Badges                     []BadgeResponse `json:"badges"`
	PlayerXP                   int             `json:"player_xp"`
	PlayerLevel                int             `json:"player_level"`
	PlayerXPNeededToLevelUp    int             `json:"player_xp_needed_to_level_up"`
	PlayerXPNeededCurrentLevel int             `json:"player_xp_needed_current_level"`
}

func (b BadgesResponse) GetPercentOfLevel() int {

	if b.PlayerXP == 0 {
		return 0
	}

	start := b.PlayerXPNeededCurrentLevel
	finish := b.PlayerXPNeededToLevelUp + b.PlayerXP
	levelRange := finish - start
	intoLevel := b.PlayerXP - b.PlayerXPNeededCurrentLevel

	return int((float64(intoLevel) / float64(levelRange)) * 100)
}

func (b BadgesResponse) GetLoadingBar() int {

	percent := b.GetPercentOfLevel()
	return int(math.Max(float64(percent), 5))
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

func (b BadgeResponse) GetTimeNice() string {
	return time.Unix(b.CompletionTime, 0).Format(helpers.DateYear)
}

func (b BadgeResponse) GetTimeUnix() int64 {
	return b.CompletionTime
}

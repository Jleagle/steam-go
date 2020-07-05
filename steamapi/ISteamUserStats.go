package steamapi

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Jleagle/unmarshal-go/ctypes"
)

// Retrieves the global achievement percentages for the specified app.
func (s Steam) GetGlobalAchievementPercentagesForApp(appID int) (percentages GlobalAchievementPercentages, b []byte, err error) {

	options := url.Values{}
	options.Set("gameid", strconv.Itoa(appID))

	b, err = s.getFromAPI("ISteamUserStats/GetGlobalAchievementPercentagesForApp/v2", options, false)
	if err != nil {
		return percentages, b, err
	}

	var resp GlobalAchievementPercentagesResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return percentages, b, err
	}

	return resp.GlobalAchievementPercentagesOuter, b, nil
}

type GlobalAchievementPercentagesResponse struct {
	GlobalAchievementPercentagesOuter GlobalAchievementPercentages `json:"achievementpercentages"`
}

type GlobalAchievementPercentages struct {
	GlobalAchievementPercentage []GlobalAchievementAchievement `json:"achievements"`
}

func (a GlobalAchievementPercentages) GetMap() map[string]float64 {
	m := map[string]float64{}
	for _, v := range a.GlobalAchievementPercentage {
		m[v.Name] = v.Percent
	}
	return m
}

type GlobalAchievementAchievement struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}

// Gets the total number of players currently active in the specified app on Steam.
func (s Steam) GetNumberOfCurrentPlayers(appID int) (players int, b []byte, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))

	b, err = s.getFromAPI("ISteamUserStats/GetNumberOfCurrentPlayers/v1", options, false)
	if err != nil {
		return players, b, err
	}

	var resp NumberOfCurrentPlayersResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return players, b, err
	}

	return resp.Response.PlayerCount, b, nil
}

type NumberOfCurrentPlayersResponse struct {
	Response struct {
		PlayerCount int `json:"player_count"`
		Result      int `json:"result"`
	} `json:"response"`
}

// Gets the complete list of stats and achievements for the specified game.
func (s Steam) GetSchemaForGame(appID int) (schema SchemaForGame, b []byte, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))
	options.Set("l", "english")

	b, err = s.getFromAPI("ISteamUserStats/GetSchemaForGame/v2", options, true)
	if err != nil {
		return schema, b, err
	}

	var resp SchemaForGameResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return schema, b, err
	}

	return resp.Game, b, nil
}

type SchemaForGameResponse struct {
	Game SchemaForGame `json:"game"`
}

type SchemaForGame struct {
	Name               string `json:"gameName"`
	Version            string `json:"gameVersion"`
	AvailableGameStats struct {
		Stats        []SchemaForGameStat        `json:"stats"`
		Achievements []SchemaForGameAchievement `json:"achievements"`
	} `json:"availableGameStats"`
}

type SchemaForGameAchievement struct {
	Name         string      `json:"name"`
	DefaultValue int         `json:"defaultvalue"`
	DisplayName  string      `json:"displayName"`
	Hidden       ctypes.Bool `json:"hidden"`
	Description  string      `json:"description"`
	Icon         string      `json:"icon"`
	IconGray     string      `json:"icongray"`
}

type SchemaForGameStat struct {
	Name         string `json:"name"`
	DefaultValue int    `json:"defaultvalue"`
	DisplayName  string `json:"displayName"`
}

func (s Steam) GetPlayerAchievements(playerID uint64, appID uint32) (schema PlayerAchievementsResponse, b []byte, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatUint(playerID, 10))
	options.Set("appid", strconv.FormatUint(uint64(appID), 10))
	options.Set("l", string(LanguageEnglish))

	b, err = s.getFromAPI("ISteamUserStats/GetPlayerAchievements/v1", options, true)
	if err != nil {
		return schema, b, err
	}

	var resp PlayerAchievementsOuterResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return schema, b, err
	}

	return resp.Playerstats, b, nil
}

type PlayerAchievementsOuterResponse struct {
	Playerstats PlayerAchievementsResponse `json:"playerstats"`
}

type PlayerAchievementsResponse struct {
	SteamID      ctypes.Int64 `json:"steamID"`
	GameName     string       `json:"gameName"`
	Achievements []struct {
		APIName     string      `json:"apiname"`
		Achieved    ctypes.Bool `json:"achieved"`
		UnlockTime  int64       `json:"unlocktime"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
	} `json:"achievements"`
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

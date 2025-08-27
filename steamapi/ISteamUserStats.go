package steamapi

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Jleagle/unmarshal-go"
)

// Retrieves the global achievement percentages for the specified app.
func (c *Client) GetGlobalAchievementPercentagesForApp(appID int) (percentages GlobalAchievementPercentages, err error) {

	options := url.Values{}
	options.Set("gameid", strconv.Itoa(appID))

	b, err := c.getFromAPI("ISteamUserStats/GetGlobalAchievementPercentagesForApp/v2", options, false)
	if err != nil {
		return percentages, err
	}

	var resp GlobalAchievementPercentagesResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return percentages, err
	}

	return resp.GlobalAchievementPercentagesOuter, nil
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
func (c *Client) GetNumberOfCurrentPlayers(appID int) (players int, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))

	b, err := c.getFromAPI("ISteamUserStats/GetNumberOfCurrentPlayers/v1", options, false)
	if err != nil {
		return players, err
	}

	var resp NumberOfCurrentPlayersResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return players, err
	}

	return resp.Response.PlayerCount, nil
}

type NumberOfCurrentPlayersResponse struct {
	Response struct {
		PlayerCount int `json:"player_count"`
		Result      int `json:"result"`
	} `json:"response"`
}

// Gets the complete list of stats and achievements for the specified game.
func (c *Client) GetSchemaForGame(appID int) (schema SchemaForGame, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))
	options.Set("l", "english")

	b, err := c.getFromAPI("ISteamUserStats/GetSchemaForGame/v2", options, true)
	if err != nil {
		return schema, err
	}

	var resp SchemaForGameResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return schema, err
	}

	return resp.Game, nil
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
	Name         string         `json:"name"`
	DefaultValue int            `json:"defaultvalue"`
	DisplayName  string         `json:"displayName"`
	Hidden       unmarshal.Bool `json:"hidden"`
	Description  string         `json:"description"`
	Icon         string         `json:"icon"`
	IconGray     string         `json:"icongray"`
}

type SchemaForGameStat struct {
	Name         string `json:"name"`
	DefaultValue int    `json:"defaultvalue"`
	DisplayName  string `json:"displayName"`
}

func (c *Client) GetPlayerAchievements(playerID uint64, appID uint32) (schema PlayerAchievementsResponse, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.FormatUint(playerID, 10))
	options.Set("appid", strconv.FormatUint(uint64(appID), 10))
	options.Set("l", string(LanguageEnglish))

	b, err := c.getFromAPI("ISteamUserStats/GetPlayerAchievements/v1", options, true)
	if err != nil {
		return schema, err
	}

	var resp PlayerAchievementsOuterResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return schema, err
	}

	return resp.Playerstats, nil
}

type PlayerAchievementsOuterResponse struct {
	Playerstats PlayerAchievementsResponse `json:"playerstats"`
}

type PlayerAchievementsResponse struct {
	SteamID      unmarshal.Int64 `json:"steamID"`
	GameName     string          `json:"gameName"`
	Achievements []struct {
		APIName     string         `json:"apiname"`
		Achieved    unmarshal.Bool `json:"achieved"`
		UnlockTime  int64          `json:"unlocktime"`
		Name        string         `json:"name"`
		Description string         `json:"description"`
	} `json:"achievements"`
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

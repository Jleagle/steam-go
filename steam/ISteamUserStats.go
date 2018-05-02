package steam

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/steam-authority/steam-authority/logger"
)

// Retrieves the global achievement percentages for the specified app.
func GetGlobalAchievementPercentagesForApp(appID int) (percentages []AchievementPercentage, err error) {

	options := url.Values{}
	options.Set("gameid", strconv.Itoa(appID))

	bytes, err := get("ISteamUserStats/GetGlobalAchievementPercentagesForApp/v2", options)
	if err != nil {
		return percentages, err
	}

	// Unmarshal JSON
	var resp GlobalAchievementPercentagesForAppResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return percentages, err
	}

	return resp.AchievementPercentagesOuter.AchievementPercentages, nil
}

type GlobalAchievementPercentagesForAppResponse struct {
	AchievementPercentagesOuter struct {
		AchievementPercentages []AchievementPercentage `json:"achievements"`
	} `json:"achievementpercentages"`
}

type AchievementPercentage struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}

// Gets the total number of players currently active in the specified app on Steam.
func GetNumberOfCurrentPlayers(appID int) (players int, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))

	bytes, err := get("ISteamUserStats/GetNumberOfCurrentPlayers/v1", options)
	if err != nil {
		return players, err
	}

	// Unmarshal JSON
	var resp NumberOfCurrentPlayersResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
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
func GetSchemaForGame(appID int) (schema GameSchema, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))
	options.Set("l", "english")

	bytes, err := get("ISteamUserStats/GetSchemaForGame/v2", options)
	if err != nil {
		return schema, err
	}

	// Unmarshal JSON
	var resp SchemaForGameResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return schema, err
	}

	return resp.Game, nil
}

type SchemaForGameResponse struct {
	Game GameSchema `json:"game"`
}

type GameSchema struct {
	GameName    string `json:"gameName"`
	GameVersion string `json:"gameVersion"`
	AvailableGameStats struct {
		Stats []struct {
			Name         string `json:"name"`
			Defaultvalue int    `json:"defaultvalue"`
			DisplayName  string `json:"displayName"`
		} `json:"stats"`
		Achievements []struct {
			Name         string `json:"name"`
			Defaultvalue int    `json:"defaultvalue"`
			DisplayName  string `json:"displayName"`
			Hidden       int    `json:"hidden"`
			Description  string `json:"description"`
			Icon         string `json:"icon"`
			Icongray     string `json:"icongray"`
		} `json:"achievements"`
	} `json:"availableGameStats"`
}

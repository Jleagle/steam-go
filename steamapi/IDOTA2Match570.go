package steamapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

func (c *Client) GetTeamInfoByTeamID(teamID string) (team TeamSummary, err error) {

	options := url.Values{}
	options.Set("start_at_team_id", teamID)
	options.Set("teams_requested", strconv.Itoa(1))

	b, err := c.getFromAPI("IDOTA2Match_570/GetTeamInfoByTeamID/v1", options, true)
	if err != nil {
		return team, err
	}

	// Unmarshal
	var resp TeamResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return team, err
	}

	if len(resp.Result.Teams) == 0 {
		return team, ErrProfileMissing
	}

	return resp.Result.Teams[0], nil
}

type TeamResponse struct {
	Result struct {
		Status int64         `json:"status"`
		Teams  []TeamSummary `json:"teams"`
	} `json:"result"`
}

type TeamSummary struct {
	Name             string `json:"name"`
	Tag              string `json:"tag"`
	Abbreviation     string `json:"abbreviation"`
	TimeCreated      int    `json:"time_created"`
	CountryCode      string `json:"country_code"`
	Url              string `json:"url"`
	GamesPlayed      int    `json:"games_played"`
	Player0AccountId int    `json:"player_0_account_id,omitempty"`
	Player1AccountId int    `json:"player_1_account_id,omitempty"`
	Player2AccountId int    `json:"player_2_account_id,omitempty"`
	Player3AccountId int    `json:"player_3_account_id,omitempty"`
	Player4AccountId int    `json:"player_4_account_id,omitempty"`
	Player5AccountId int    `json:"player_5_account_id,omitempty"`
	Player6AccountId int    `json:"player_6_account_id,omitempty"`
	AdminAccountId   int    `json:"admin_account_id"`
}

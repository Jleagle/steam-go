package steamapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

func (c *Client) GetTeamInfoByTeamID(ctx context.Context, teamID string) (team TeamSummary, err error) {

	options := url.Values{}
	options.Set("start_at_team_id", teamID)
	options.Set("teams_requested", strconv.Itoa(1))

	b, err := c.getFromAPI(ctx, "IDOTA2Match_570/GetTeamInfoByTeamID/v1", options, true)
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
	URL              string `json:"url"`
	GamesPlayed      int    `json:"games_played"`
	Player0AccountID int    `json:"player_0_account_id,omitempty"`
	Player1AccountID int    `json:"player_1_account_id,omitempty"`
	Player2AccountID int    `json:"player_2_account_id,omitempty"`
	Player3AccountID int    `json:"player_3_account_id,omitempty"`
	Player4AccountID int    `json:"player_4_account_id,omitempty"`
	Player5AccountID int    `json:"player_5_account_id,omitempty"`
	Player6AccountID int    `json:"player_6_account_id,omitempty"`
	AdminAccountID   int    `json:"admin_account_id"`
}

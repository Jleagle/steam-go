package steamapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

func (c *Client) GetAppList(limit int, offset int, afterDate int64, language LanguageCode) (apps AppList, err error) {

	q := url.Values{}
	q.Set("include_games", "1")
	q.Set("include_dlc", "1")
	q.Set("include_software", "1")
	q.Set("include_videos", "1")
	q.Set("include_hardware", "1")

	if afterDate > 0 {
		q.Set("if_modified_since", strconv.FormatInt(afterDate, 10))
	}

	if language != "" {
		q.Set("have_description_language", string(language))
	}

	if offset > 0 {
		q.Set("last_appid", strconv.Itoa(offset))
	}
	if limit > 0 {
		q.Set("max_results", strconv.Itoa(limit))
	}

	b, err := c.getFromAPI("IStoreService/GetAppList/v1", q, true)
	if err != nil {
		return apps, err
	}

	var resp AppListResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return apps, err
	}

	return resp.AppListResponseInner, nil
}

type AppListResponse struct {
	AppListResponseInner AppList `json:"response"`
}

type AppList struct {
	Apps []struct {
		AppID             int    `json:"appid"`
		Name              string `json:"name"`
		LastModified      int64  `json:"last_modified"`
		PriceChangeNumber int    `json:"price_change_number"`
	} `json:"apps"`
	HaveMoreResults bool `json:"have_more_results"`
	LastAppID       int  `json:"last_appid"`
}

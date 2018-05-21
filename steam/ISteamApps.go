package steam

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type GetAppList struct {
	ModifiedAfter int64
	AfterAppID    int
	MaxResults    int // 10,000 default
}

func (s Steam) GetAppList(options GetAppList) (apps AppList, bytes []byte, err error) {

	q := url.Values{}
	q.Set("include_games", "1")
	q.Set("include_dlc", "1")
	q.Set("include_software", "1")
	q.Set("include_videos", "1")
	q.Set("include_hardware", "1")

	if options.ModifiedAfter > 0 {
		q.Set("if_modified_since", strconv.FormatInt(options.ModifiedAfter, 10))
	}
	if options.AfterAppID > 0 {
		q.Set("last_appid", strconv.Itoa(options.AfterAppID))
	}
	if options.MaxResults > 0 {
		q.Set("max_results", strconv.Itoa(options.MaxResults))
	}

	bytes, err = s.getFromAPI("ISteamApps/GetAppList/v2", q)
	if err != nil {
		return apps, bytes, err
	}

	var resp AppListResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return apps, bytes, err
	}

	return resp.AppList, bytes, nil
}

type AppListResponse struct {
	AppList AppList `json:"applist"`
}

type AppList struct {
	Apps []App `json:"apps"`
}

type App struct {
	AppID             int    `json:"appid"`
	Name              string `json:"name"`
	LastModified      string `json:"last_modified"`
	PriceChangeNumber string `json:"price_change_number"`
}

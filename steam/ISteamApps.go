package steam

import (
	"encoding/json"
	"net/url"
)

func (s Steam) GetAppList() (apps AppList, bytes []byte, err error) {

	bytes, err = s.getFromAPI("ISteamApps/GetAppList/v2", url.Values{})
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
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}

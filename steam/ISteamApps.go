package steam

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/steam-authority/steam-authority/logger"
)

func GetAppList() (apps []GetAppListApp, err error) {

	bytes, err := get("ISteamApps/GetAppList/v2/", url.Values{})
	if err != nil {
		return apps, err
	}

	// Unmarshal JSON
	resp := GetAppListBody{}
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return apps, err
	}

	return resp.AppList.Apps, nil
}

type GetAppListBody struct {
	AppList GetAppListAppList `json:"applist"`
}

type GetAppListAppList struct {
	Apps []GetAppListApp `json:"apps"`
}

type GetAppListApp struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}

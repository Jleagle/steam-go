package steamapi

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Jleagle/unmarshal-go/ctypes"
)

func (c Client) GetNews(appID int, limit int) (articles News, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))
	options.Set("count", strconv.Itoa(limit))
	options.Set("maxlength", "0")

	b, err := c.getFromAPI("ISteamNews/GetNewsForApp/v2", options, false)
	if err != nil {
		return articles, err
	}

	// Unmarshal
	var resp NewsResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return articles, err
	}

	return resp.App, nil
}

type NewsResponse struct {
	App News `json:"appnews"`
}

type News struct {
	AppID int `json:"appid"`
	Items []struct {
		GID           ctypes.Int64 `json:"gid"`
		Title         string       `json:"title"`
		URL           string       `json:"url"`
		IsExternalURL bool         `json:"is_external_url"`
		Author        string       `json:"author"`
		Contents      string       `json:"contents"`
		Feedlabel     string       `json:"feedlabel"`
		Date          int64        `json:"date"`
		Feedname      string       `json:"feedname"`
		FeedType      int          `json:"feed_type"`
		AppID         int          `json:"appid"`
	} `json:"newsitems"`
	Count int `json:"count"`
}

package steam

import (
	"encoding/json"
	"net/url"
)

func GetNews(id string) (articles News, bytes []byte, err error) {

	options := url.Values{}
	options.Set("appid", id)
	options.Set("count", "20")

	bytes, err = get("ISteamNews/GetNewsForApp/v2/", options)
	if err != nil {
		return articles, bytes, err
	}

	var resp NewsResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return articles, bytes, err
	}

	return resp.App, bytes, nil
}

type NewsResponse struct {
	App News `json:"appnews"`
}

type News struct {
	AppID int           `json:"appid"`
	Items []NewsArticle `json:"newsitems"`
	Count int           `json:"count"`
}

type NewsArticle struct {
	GID           string `json:"gid"` // todo, make int64
	Title         string `json:"title"`
	URL           string `json:"url"`
	IsExternalURL bool   `json:"is_external_url"`
	Author        string `json:"author"`
	Contents      string `json:"contents"`
	Feedlabel     string `json:"feedlabel"`
	Date          int64  `json:"date"`
	Feedname      string `json:"feedname"`
	FeedType      int    `json:"feed_type"`
	AppID         int    `json:"appid"`
}

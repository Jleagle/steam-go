package steam

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strconv"
)

func (s Steam) GetNews(appID int, limit int) (articles News, bytes []byte, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))
	options.Set("count", strconv.Itoa(limit))
	options.Set("maxlength", "0")

	bytes, err = s.getFromAPI("ISteamNews/GetNewsForApp/v2", options)
	if err != nil {
		return articles, bytes, err
	}

	// Regex
	var regex *regexp.Regexp
	var str = string(bytes)

	regex = regexp.MustCompile(`"gid":\s?"(\d+)"`)
	str = regex.ReplaceAllString(str, `"gid": $1`)

	bytes = []byte(str)

	// Unmarshal
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
	GID           int64  `json:"gid"`
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

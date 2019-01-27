package steam

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Jleagle/unmarshal-go/ctypes"
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

	// Unmarshal
	var resp NewsResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
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
	GID           ctypes.CInt64 `json:"gid"`
	Title         string        `json:"title"`
	URL           string        `json:"url"`
	IsExternalURL bool          `json:"is_external_url"`
	Author        string        `json:"author"`
	Contents      string        `json:"contents"`
	Feedlabel     string        `json:"feedlabel"`
	Date          int64         `json:"date"`
	Feedname      string        `json:"feedname"`
	FeedType      int           `json:"feed_type"`
	AppID         int           `json:"appid"`
}

package steamapi

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/Jleagle/unmarshal-go/ctypes"
)

func (c Client) GetInventory(playerID int64, appID int) (resp CommunityInventory, b []byte, err error) {

	b, err = c.getFromStore("profiles/"+strconv.FormatInt(playerID, 10)+"/inventory/json/"+strconv.Itoa(appID)+"/2", url.Values{})
	if err != nil {
		return resp, b, err
	}

	//
	err = json.Unmarshal(b, &resp)
	return resp, b, err
}

type CommunityInventory struct {
	Success     bool `json:"success"`
	RgInventory map[string]struct {
		ID         string `json:"id"`
		Classid    string `json:"classid"`
		Instanceid string `json:"instanceid"`
		Amount     string `json:"amount"`
		Pos        int    `json:"pos"`
	} `json:"rgInventory"`
	RgCurrency     []interface{} `json:"rgCurrency"`
	RgDescriptions map[string]struct {
		Appid                       string `json:"appid"`
		Classid                     string `json:"classid"`
		Instanceid                  string `json:"instanceid"`
		IconURL                     string `json:"icon_url"`
		IconURLLarge                string `json:"icon_url_large"`
		IconDragURL                 string `json:"icon_drag_url"`
		Name                        string `json:"name"`
		MarketHashName              string `json:"market_hash_name"`
		MarketName                  string `json:"market_name"`
		NameColor                   string `json:"name_color"`
		BackgroundColor             string `json:"background_color"`
		Type                        string `json:"type"`
		Tradable                    int    `json:"tradable"`
		Marketable                  int    `json:"marketable"`
		Commodity                   int    `json:"commodity"`
		MarketTradableRestriction   string `json:"market_tradable_restriction"`
		MarketMarketableRestriction string `json:"market_marketable_restriction"`
		Descriptions                []struct {
			Value   string `json:"value"`
			Color   string `json:"color,omitempty"`
			AppData struct {
				DefIndex string `json:"def_index"`
			} `json:"app_data,omitempty"`
		} `json:"descriptions"`
		Actions []struct {
			Name string `json:"name"`
			Link string `json:"link"`
		} `json:"actions"`
		MarketActions []struct {
			Name string `json:"name"`
			Link string `json:"link"`
		} `json:"market_actions"`
		Tags []struct {
			InternalName string `json:"internal_name"`
			Name         string `json:"name"`
			Category     string `json:"category"`
			Color        string `json:"color,omitempty"`
			CategoryName string `json:"category_name"`
		} `json:"tags"`
		AppData struct {
			DefIndex string `json:"def_index"`
			Quality  string `json:"quality"`
		} `json:"app_data"`
	} `json:"rgDescriptions"`
	More      bool `json:"more"`
	MoreStart bool `json:"more_start"`
}

type MarketSearchPayload struct {
	FriendlyDescriptions bool
	SortColumn           string
	SortOrder            bool
	AppID                int
	Limit                int
	Offset               int
}

func (c Client) GetMarketSearch(payload MarketSearchPayload) (resp MarketSearch, b []byte, err error) {

	vals := url.Values{}
	if payload.FriendlyDescriptions {
		vals.Set("search_descriptions", "1")
	} else {
		vals.Set("search_descriptions", "0")
	}
	if payload.SortColumn != "" {
		vals.Set("sort_column", payload.SortColumn)
	}
	if payload.SortOrder {
		vals.Set("sort_dir", "asc")
	} else {
		vals.Set("sort_dir", "desc")
	}
	if payload.AppID > 0 {
		vals.Set("appid", strconv.Itoa(payload.AppID))
	}
	if payload.Limit > 0 {
		vals.Set("count", strconv.Itoa(payload.Limit))
	}
	vals.Set("start", strconv.Itoa(payload.Offset))
	vals.Set("norender", "1")

	b, err = c.getFromStore("market/search/render", vals)
	if err != nil {
		return resp, b, err
	}

	//
	err = json.Unmarshal(b, &resp)
	return resp, b, err
}

type MarketSearch struct {
	Success    bool `json:"success"`
	Start      int  `json:"start"`
	Pagesize   int  `json:"pagesize"`
	TotalCount int  `json:"total_count"`
	Searchdata struct {
		Query              string `json:"query"`
		SearchDescriptions bool   `json:"search_descriptions"`
		TotalCount         int    `json:"total_count"`
		Pagesize           int    `json:"pagesize"`
		Prefix             string `json:"prefix"`
		ClassPrefix        string `json:"class_prefix"`
	} `json:"searchdata"`
	Results []struct {
		Name             string `json:"name"`
		HashName         string `json:"hash_name"`
		SellListings     int    `json:"sell_listings"`
		SellPrice        int    `json:"sell_price"`
		SellPriceText    string `json:"sell_price_text"`
		AppIcon          string `json:"app_icon"`
		AppName          string `json:"app_name"`
		AssetDescription struct {
			Appid                       int    `json:"appid"`
			Classid                     string `json:"classid"`
			Instanceid                  string `json:"instanceid"`
			Currency                    int    `json:"currency"`
			BackgroundColor             string `json:"background_color"`
			IconURL                     string `json:"icon_url"`
			IconURLLarge                string `json:"icon_url_large"`
			Tradable                    int    `json:"tradable"`
			Name                        string `json:"name"`
			Type                        string `json:"type"`
			MarketName                  string `json:"market_name"`
			MarketHashName              string `json:"market_hash_name"`
			Commodity                   int    `json:"commodity"`
			MarketTradableRestriction   int    `json:"market_tradable_restriction"`
			MarketMarketableRestriction int    `json:"market_marketable_restriction"`
			Marketable                  int    `json:"marketable"`
		} `json:"asset_description"`
		SalePriceText string `json:"sale_price_text"`
	} `json:"results"`
}

func GetPriceOverview() (resp PriceOverview, b []byte, err error) {

	// http://steamcommunity.com/market/priceoverview/?appid=730&currency=3&market_hash_name=StatTrak%E2%84%A2 M4A1-S | Hyper Beast (Minimal Wear)

	return resp, b, err

}

type PriceOverview struct {
	Success     bool   `json:"success"`
	LowestPrice string `json:"lowest_price"`
	Volume      string `json:"volume"`
	MedianPrice string `json:"median_price"`
}

var ErrRateLimited = errors.New("rate limited")

// Rate limited to once per minute
func (c Client) GetGroup(id string, vanityURL string, page int) (resp GroupInfo, b []byte, err error) {

	vals := url.Values{}
	vals.Set("p", strconv.Itoa(page))
	// vals.Set("xml", "1") // Without this, it redirects to a slug, so we can get the type

	var urlx string
	if id != "" {
		b, urlx, err = c.getFromCommunity("gid/"+id+"/memberslistxml", vals)
	} else {
		b, urlx, err = c.getFromCommunity("groups/"+vanityURL+"/memberslistxml", vals)
	}

	if err != nil {
		return resp, b, err
	}

	if string(b) == "null" {
		return resp, b, ErrRateLimited
	}

	err = xml.Unmarshal(b, &resp)

	if strings.Contains(urlx, "/games/") {
		resp.Type = "game"
	} else if strings.Contains(urlx, "/groups/") {
		resp.Type = "group"
	}

	return resp, b, err
}

type GroupInfo struct {
	Type    string   `xml:"-"` // Not in Steam response
	XMLName xml.Name `xml:"memberList"`
	Text    string   `xml:",chardata"`
	ID64    string   `xml:"groupID64"` // Too big for int64
	Details struct {
		Text          string     `xml:",chardata"`
		Name          string     `xml:"groupName"`
		URL           string     `xml:"groupURL"`
		Headline      string     `xml:"headline"`
		Summary       string     `xml:"summary"`
		AvatarIcon    string     `xml:"avatarIcon"`
		AvatarMedium  string     `xml:"avatarMedium"`
		AvatarFull    string     `xml:"avatarFull"`
		MemberCount   ctypes.Int `xml:"memberCount"`
		MembersInChat ctypes.Int `xml:"membersInChat"`
		MembersInGame ctypes.Int `xml:"membersInGame"`
		MembersOnline ctypes.Int `xml:"membersOnline"`
	} `xml:"groupDetails"`
	MemberCount      ctypes.Int `xml:"memberCount"`
	TotalPages       ctypes.Int `xml:"totalPages"`
	CurrentPage      ctypes.Int `xml:"currentPage"`
	StartingMember   ctypes.Int `xml:"startingMember"`
	NextPageLink     string     `xml:"nextPageLink"`
	PreviousPageLink string     `xml:"previousPageLink"`
	Members          struct {
		Text      string         `xml:",chardata"`
		SteamID64 []ctypes.Int64 `xml:"steamID64"`
	} `xml:"members"`
}

func (c Client) GetComments(playerID int64, limit int, offset int) (resp Comments, b []byte, err error) {

	vals := url.Values{}
	vals.Set("count", strconv.Itoa(limit))
	if offset > 0 {
		vals.Set("start", strconv.Itoa(offset))
	}

	b, _, err = c.getFromCommunity("comment/Profile/render/"+strconv.FormatInt(playerID, 10), vals)
	if err != nil {
		return resp, b, err
	}

	if len(b) == 0 {
		return resp, b, nil
	}

	//
	err = json.Unmarshal(b, &resp)
	return resp, b, err
}

type Comments struct {
	Success      bool       `json:"success"`
	Name         string     `json:"name"`
	Start        int        `json:"start"`
	PageSize     ctypes.Int `json:"pagesize"`
	TotalCount   int        `json:"total_count"`
	Upvotes      int        `json:"upvotes"`
	HasUpvoted   int        `json:"has_upvoted"`
	CommentsHTML string     `json:"comments_html"`
	TimeLastPost int64      `json:"timelastpost"`
}

func (c Client) GetAliases(playerID int64) (resp []Alias, b []byte, err error) {

	b, _, err = c.getFromCommunity("profiles/"+strconv.FormatInt(playerID, 10)+"/ajaxaliases", nil)
	if err != nil {
		return resp, b, err
	}

	if strings.HasPrefix(string(b), "<") {
		return resp, b, ErrProfileMissing
	}

	err = json.Unmarshal(b, &resp)
	return resp, b, err
}

type Alias struct {
	Alias string `json:"newname"`
	Time  string `json:"timechanged"`
}

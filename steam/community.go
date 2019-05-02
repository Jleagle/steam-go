package steam

import (
	"encoding/json"
	"net/url"
	"strconv"
)

func (s Steam) GetInventory(playerID int64, appID int) (resp CommunityInventory, bytes []byte, err error) {

	bytes, err = s.getFromStore("profiles/"+strconv.FormatInt(playerID, 10)+"/inventory/json/"+strconv.Itoa(appID)+"/2", url.Values{})
	if err != nil {
		return resp, bytes, err
	}

	//
	err = json.Unmarshal(bytes, &resp)
	return resp, bytes, err
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

func (s Steam) GetMarketSearch(payload MarketSearchPayload) (resp MarketSearch, bytes []byte, err error) {

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

	bytes, err = s.getFromStore("market/search/render", vals)
	if err != nil {
		return resp, bytes, err
	}

	//
	err = json.Unmarshal(bytes, &resp)
	return resp, bytes, err
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

func GetPriceOverview() (resp PriceOverview, bytes []byte, err error) {

	// http://steamcommunity.com/market/priceoverview/?appid=730&currency=3&market_hash_name=StatTrak%E2%84%A2 M4A1-S | Hyper Beast (Minimal Wear)

	return resp, bytes, err

}

type PriceOverview struct {
	Success     bool   `json:"success"`
	LowestPrice string `json:"lowest_price"`
	Volume      string `json:"volume"`
	MedianPrice string `json:"median_price"`
}

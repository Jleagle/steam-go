package steam

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

func (s Steam) GetItemDefArchive(appID int, digest string) (archives []ItemDefArchive, bytes []byte, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))
	options.Set("digest", digest)

	bytes, err = s.getFromAPI("IGameInventory/GetItemDefArchive/v1", options)
	if err != nil {
		return archives, bytes, err
	}

	err = json.Unmarshal(bytes, &archives)
	return archives, bytes, err
}

type ItemDefArchive struct {
	AppID            string    `json:"appid"`
	ItemdefID        string    `json:"itemdefid"`
	Timestamp        time.Time `json:"Timestamp"`
	Modified         string    `json:"modified"`
	DateCreated      string    `json:"date_created"`
	Type             string    `json:"type"`
	DisplayType      string    `json:"display_type"`
	Name             string    `json:"name"`
	Quantity         int       `json:"quantity"`
	Description      string    `json:"description"`
	IconURL          string    `json:"icon_url,omitempty"`
	IconURLLarge     string    `json:"icon_url_large,omitempty"`
	Tags             string    `json:"tags,omitempty"`
	Tradable         bool      `json:"tradable"`
	Marketable       bool      `json:"marketable"`
	Commodity        bool      `json:"commodity"`
	DropInterval     int       `json:"drop_interval"`
	DropMaxPerWindow int       `json:"drop_max_per_window"`
	WorkshopID       string    `json:"workshopid"`
	Descrption       string    `json:"descrption,omitempty"`
	ItemQuality      string    `json:"item_quality,omitempty"`
	Hash             string    `json:"hash,omitempty"`
	Price            string    `json:"price,omitempty"`
	Promo            string    `json:"promo,omitempty"`
	Exchange         string    `json:"exchange,omitempty"`
	Bundle           string    `json:"bundle,omitempty"`
}

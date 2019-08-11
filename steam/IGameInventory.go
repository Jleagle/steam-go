package steam

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/Jleagle/unmarshal-go/ctypes"
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
	AppID            ctypes.CInt   `json:"appid"`
	ItemdefID        ctypes.CInt   `json:"itemdefid"`
	Timestamp        time.Time     `json:"Timestamp"`
	Modified         time.Time     `json:"modified"`
	DateCreated      time.Time     `json:"date_created"`
	Type             string        `json:"type"`
	DisplayType      string        `json:"display_type"`
	Name             string        `json:"name"`
	Quantity         int           `json:"quantity"`
	Description      string        `json:"description"`
	IconURL          string        `json:"icon_url"`
	IconURLLarge     string        `json:"icon_url_large"`
	Tags             string        `json:"tags"`
	Tradable         bool          `json:"tradable"`
	Marketable       bool          `json:"marketable"`
	Commodity        bool          `json:"commodity"`
	DropInterval     int           `json:"drop_interval"`
	DropMaxPerWindow int           `json:"drop_max_per_window"`
	WorkshopID       ctypes.CInt64 `json:"workshopid"`
	Descrption       string        `json:"descrption"`
	ItemQuality      string        `json:"item_quality"`
	Hash             string        `json:"hash"`
	Price            string        `json:"price"`
	Promo            string        `json:"promo"`
	Exchange         string        `json:"exchange"`
	Bundle           string        `json:"bundle"`
}

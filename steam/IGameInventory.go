package steam

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/Jleagle/unmarshal-go/ctypes"
)

var ErrInvalidDigest = errors.New("invalid digest")

func (s Steam) GetItemDefArchive(appID int, digest string) (archives []ItemDefArchive, b []byte, err error) {

	if digest == "" {
		return archives, b, ErrInvalidDigest
	}

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))
	options.Set("digest", digest)

	b, err = s.getFromAPI("IGameInventory/GetItemDefArchive/v1", options)
	if err != nil {
		return archives, b, err
	}

	// The response has an empty byte at the end of it causing Unmarshal to fail
	b = bytes.TrimSuffix(b, []byte{0x00})

	err = json.Unmarshal(b, &archives)
	return archives, b, err
}

type ItemDefArchive struct {
	AppID            ctypes.CInt   `json:"appid"`
	Bundle           string        `json:"bundle"`
	Commodity        bool          `json:"commodity"`
	DateCreated      string        `json:"date_created"` // Can't be time.Time, - "20161010T080316Z"
	Description      string        `json:"description"`
	DisplayType      string        `json:"display_type"`
	DropInterval     int           `json:"drop_interval"`
	DropMaxPerWindow int           `json:"drop_max_per_window"`
	Exchange         string        `json:"exchange"`
	Hash             string        `json:"hash"`
	IconURL          string        `json:"icon_url"`
	IconURLLarge     string        `json:"icon_url_large"`
	ItemdefID        ctypes.CInt   `json:"itemdefid"`
	ItemQuality      string        `json:"item_quality"`
	Marketable       bool          `json:"marketable"`
	Modified         string        `json:"modified"` // Can't be time.Time, - "20161010T080316Z"
	Name             string        `json:"name"`
	Price            string        `json:"price"`
	Promo            string        `json:"promo"`
	Quantity         int           `json:"quantity"`
	Tags             string        `json:"tags"`
	Timestamp        time.Time     `json:"Timestamp"`
	Tradable         bool          `json:"tradable"`
	Type             string        `json:"type"`
	WorkshopID       ctypes.CInt64 `json:"workshopid"`
}

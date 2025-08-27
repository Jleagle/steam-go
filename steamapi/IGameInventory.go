package steamapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/Jleagle/unmarshal-go"
)

var ErrInvalidDigest = errors.New("invalid digest")

func (c *Client) GetItemDefArchive(appID int, digest string) (archives []ItemDefArchive, err error) {

	if digest == "" {
		return archives, ErrInvalidDigest
	}

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))
	options.Set("digest", digest)

	b, err := c.getFromAPI("IGameInventory/GetItemDefArchive/v1", options, false)
	if err != nil {
		return archives, err
	}

	// The response has an empty byte at the end of it causing Unmarshal to fail
	b = bytes.TrimSuffix(b, []byte{0x00})

	if len(b) == 0 {
		return archives, nil
	}

	err = json.Unmarshal(b, &archives)
	return archives, err
}

type ItemDefArchive struct {
	AppID            unmarshal.Int    `json:"appid"`
	Bundle           string           `json:"bundle"`
	Commodity        unmarshal.Bool   `json:"commodity"`
	DateCreated      string           `json:"date_created"` // Can't be time.Time, - "20161010T080316Z"
	Description      string           `json:"description"`
	DisplayType      string           `json:"display_type"`
	DropInterval     unmarshal.Int    `json:"drop_interval"`
	DropMaxPerWindow unmarshal.Int    `json:"drop_max_per_window"`
	Exchange         string           `json:"exchange"`
	Hash             string           `json:"hash"`
	IconURL          string           `json:"icon_url"`
	IconURLLarge     string           `json:"icon_url_large"`
	ItemDefID        unmarshal.Int    `json:"itemdefid"`
	ItemQuality      unmarshal.String `json:"item_quality"` // Can be bool
	Marketable       unmarshal.Bool   `json:"marketable"`
	Modified         string           `json:"modified"` // Can't be time.Time, - "20161010T080316Z"
	Name             string           `json:"name"`
	Price            string           `json:"price"`
	Promo            string           `json:"promo"`
	Quantity         unmarshal.Int    `json:"quantity"` // Can be 0 or "0"
	Tags             string           `json:"tags"`
	Timestamp        string           `json:"Timestamp"` // Can't be time.Time, - ""
	Tradable         unmarshal.Bool   `json:"tradable"`  // Can be false or "false"
	Type             string           `json:"type"`
	WorkshopID       unmarshal.Int64  `json:"workshopid"`
}

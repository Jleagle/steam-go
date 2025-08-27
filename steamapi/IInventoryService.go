package steamapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

func (c *Client) GetItemDefMeta(appID int) (meta ItemDefMeta, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))

	b, err := c.getFromAPI("IInventoryService/GetItemDefMeta/v1", options, true)
	if err != nil {
		return meta, err
	}

	err = json.Unmarshal(b, &meta)
	return meta, err
}

type ItemDefMeta struct {
	Response struct {
		Modified int64  `json:"modified"`
		Digest   string `json:"digest"`
	} `json:"response"`
}

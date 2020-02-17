package steamapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

func (s Steam) GetItemDefMeta(appID int) (meta ItemDefMeta, bytes []byte, err error) {

	options := url.Values{}
	options.Set("appid", strconv.Itoa(appID))

	bytes, err = s.getFromAPI("IInventoryService/GetItemDefMeta/v1", options)
	if err != nil {
		return meta, bytes, err
	}

	err = json.Unmarshal(bytes, &meta)
	return meta, bytes, err
}

type ItemDefMeta struct {
	Response struct {
		Modified int64  `json:"modified"`
		Digest   string `json:"digest"`
	} `json:"response"`
}

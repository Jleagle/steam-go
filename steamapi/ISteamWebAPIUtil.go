package steamapi

import (
	"encoding/json"
	"net/url"
)

// Gets the list of supported API calls. This is used to build this documentation.
func (s Steam) GetSupportedAPIList() (percentages APIInterfaces, b []byte, err error) {

	b, err = s.getFromAPI("ISteamWebAPIUtil/GetSupportedAPIList/v1", url.Values{}, false)
	if err != nil {
		return percentages, b, err
	}

	// Unmarshal JSON
	var resp SupportedAPIListResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return percentages, b, err
	}

	return resp.APIList, b, nil
}

type SupportedAPIListResponse struct {
	APIList APIInterfaces `json:"apilist"`
}

type APIInterfaces struct {
	Interfaces []APIInterface `json:"interfaces"`
}

type APIInterface struct {
	Name    string `json:"name"`
	Methods []struct {
		Name       string `json:"name"`
		Version    int    `json:"version"`
		HTTPmethod string `json:"httpmethod"`
		Parameters []struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Optional    bool   `json:"optional"`
			Description string `json:"description"`
		} `json:"parameters"`
	} `json:"methods"`
}

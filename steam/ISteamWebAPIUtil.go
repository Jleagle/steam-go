package steam

import (
	"net/url"

	"github.com/Jleagle/unmarshal-go/unmarshal"
)

// Gets the list of supported API calls. This is used to build this documentation.
func (s Steam) GetSupportedAPIList() (percentages APIInterfaces, err error) {

	bytes, err := s.getFromAPI("ISteamWebAPIUtil/GetSupportedAPIList/v1", url.Values{})
	if err != nil {
		return percentages, err
	}

	// Unmarshal JSON
	var resp SupportedAPIListResponse
	err = unmarshal.Unmarshal(bytes, &resp)
	if err != nil {
		return percentages, err
	}

	return resp.APIList, nil
}

type SupportedAPIListResponse struct {
	APIList APIInterfaces `json:"apilist"`
}

type APIInterfaces struct {
	Interfaces []APIInterface `json:"interfaces"`
}

type APIInterface struct {
	Name    string      `json:"name"`
	Methods []APIMethod `json:"methods"`
}

type APIMethod struct {
	Name       string         `json:"name"`
	Version    int            `json:"version"`
	Httpmethod string         `json:"httpmethod"`
	Parameters []APIParameter `json:"parameters"`
}

type APIParameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Optional    bool   `json:"optional"`
	Description string `json:"description"`
}

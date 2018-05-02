package steam

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/steam-authority/steam-authority/logger"
)

// Gets the list of supported API calls. This is used to build this documentation.
func GetSupportedAPIList() (percentages []SupportedAPIListInterface, err error) {

	bytes, err := get("ISteamWebAPIUtil/GetSupportedAPIList/v1", url.Values{})
	if err != nil {
		return percentages, err
	}

	// Unmarshal JSON
	var resp SupportedAPIListResponseResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			logger.Info(err.Error() + " - " + string(bytes))
		}
		return percentages, err
	}

	return resp.SupportedAPIListResponseOuter.SupportedAPIListInterfaces, nil
}

type SupportedAPIListResponseResponse struct {
	SupportedAPIListResponseOuter struct {
		SupportedAPIListInterfaces []SupportedAPIListInterface `json:"interfaces"`
	} `json:"apilist"`
}

type SupportedAPIListInterface struct {
	Name string `json:"name"`
	Methods []struct {
		Name       string `json:"name"`
		Version    int    `json:"version"`
		Httpmethod string `json:"httpmethod"`
		Parameters []struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Optional    bool   `json:"optional"`
			Description string `json:"description"`
		} `json:"parameters"`
	} `json:"methods"`
}

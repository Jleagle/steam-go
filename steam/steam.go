package steam

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	key string

	ErrInvalidKey = errors.New("invalid api key")
)

func SetKey(k string) {
	key = k
}

// todo, handle error codes crrectly.. https://partner.steamgames.com/doc/webapi_overview/responses
func get(path string, query url.Values) (bytes []byte, err error) {

	// Build endpoint
	query.Add("format", "json")
	query.Add("key", key)

	path = "http://api.steampowered.com/" + path + "?" + query.Encode()

	// Grab the JSON from node
	response, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Convert to bytes
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Check response
	if string(contents) == "<html><head><title>Forbidden</title></head><body><h1>Forbidden</h1>Access is denied. Retrying will not help. Please verify your <pre>key=</pre> parameter.</body></html>" {
		return contents, ErrInvalidKey
	}

	return contents, nil
}

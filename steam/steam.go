package steam

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Steam struct {
	key        string      // API key
	throttle   bool        // Rate limit calls
	logChannel chan string // channel to return call URLs
	format     string      // json, vdf, xml
}

var (
	apiThrottle   chan time.Time
	storeThrottle chan time.Time
)

func init() {

	// API Throttle
	apiTick := time.NewTicker(time.Hour * 24 / 100000)
	defer apiTick.Stop()

	apiThrottle = make(chan time.Time, 1000) // Burst limit
	go func() {
		for t := range apiTick.C {
			select {
			case apiThrottle <- t:
			default:
			}
		}
	}()

	// Store throttle
	storeTick := time.NewTicker(time.Minute * 5 / 200)
	defer storeTick.Stop()

	storeThrottle = make(chan time.Time, 2) // Burst limit
	go func() {
		for t := range storeTick.C {
			select {
			case storeThrottle <- t:
			default:
			}
		}
	}()
}

var statusCodes = map[int]string{
	200: "Success!",
	400: "Please verify that all required parameters are being sent.",
	401: "Access is denied. Retrying will not help. Please verify your key= parameter.",
	403: "Access is denied. Retrying will not help. Please verify your key= parameter.",
	404: "The API requested does not exists.",
	405: "This API has been called with a the wrong HTTP method like GET or PUSH.",
	429: "You are being rate limited.",
	500: "An unrecoverable error has occurred, please try again.If this continues to persist then please post to the Steamworks developer discussion with additional details of your request.",
	503: "Server is temporarily unavailable, or too busy to respond.Please wait and try again later.",
}

func (s Steam) getFromAPI(path string, query url.Values) (bytes []byte, err error) {

	if s.throttle {
		<-apiThrottle
	}

	if s.format == "" {
		s.format = "json"
	}

	query.Add("format", s.format)
	query.Add("key", s.key)

	path = "http://api.steampowered.com/" + path + "?" + query.Encode()

	if s.logChannel != nil {
		s.logChannel <- path
	}

	response, err := http.Get(path)
	if err != nil {
		return bytes, err
	}
	defer response.Body.Close()

	// Handle errors
	if response.StatusCode != 200 {
		if val, ok := statusCodes[response.StatusCode]; ok {
			return bytes, errors.New(val)
		} else {
			return bytes, errors.New("steam: something went wrong")
		}
	}

	//
	bytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}

func (s Steam) getFromStore(path string, query url.Values) (bytes []byte, err error) {

	if s.throttle {
		<-storeThrottle
	}

	path = "http://store.steampowered.com/" + path + "?" + query.Encode()

	if s.logChannel != nil {
		s.logChannel <- path
	}

	response, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

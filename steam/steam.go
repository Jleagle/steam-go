package steam

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Steam struct {
	Key        string      // API key
	Throttle   bool        // Rate limit calls
	LogChannel chan string // channel to return call URLs
	Format     string      // json, vdf, xml
	HTTP       bool
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
	400: "steam: please verify that all required parameters are being sent.",
	401: "steam: access is denied. retrying will not help. please verify your key= parameter.",
	403: "steam: access is denied. retrying will not help. please verify your key= parameter.",
	404: "steam: the api requested does not exists.",
	405: "steam: this api has been called with a the wrong http method like get or push.",
	429: "steam: you are being rate limited.",
	500: "steam: an unrecoverable error has occurred, please try again. if this continues to persist then please post to the steamworks developer discussion with additional details of your request.",
	503: "steam: server is temporarily unavailable, or too busy to respond. please wait and try again later.",
}

type Error struct {
	err  string
	code int
}

func (e Error) Error() string {
	return e.err
}

func (e Error) Code() int {
	return e.code
}

func (e Error) IsHardFail() bool {
	for _, v := range []int{400, 401, 403, 404, 405} {
		if v == e.code {
			return true
		}
	}
	return false
}

func (s Steam) getFromAPI(path string, query url.Values) (bytes []byte, err error) {

	if s.Throttle {
		<-apiThrottle
	}

	if s.Format == "" {
		s.Format = "json"
	}

	query.Add("format", s.Format)
	query.Add("key", s.Key)

	path = s.getProtocol() + "api.steampowered.com/" + path + "?" + query.Encode()

	if s.LogChannel != nil {
		s.LogChannel <- path
	}

	response, err := http.Get(path)
	if err != nil {
		return bytes, err
	}
	defer response.Body.Close()

	// Handle errors
	if response.StatusCode != 200 {
		if val, ok := statusCodes[response.StatusCode]; ok {
			return bytes, Error{
				err:  val,
				code: response.StatusCode,
			}
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

	if s.Throttle {
		<-storeThrottle
	}

	path = "https://store.steampowered.com/" + path + "?" + query.Encode()

	if s.LogChannel != nil {
		s.LogChannel <- path
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

func (s Steam) getProtocol() string {
	if s.HTTP {
		return "http://"
	}

	return "https://"
}

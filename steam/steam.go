package steam

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const defaultUserAgent = "github.com/Jleagle/steam-go"

var statusCodes = map[int]string{
	400: "please verify that all required parameters are being sent.",
	401: "access is denied. retrying will not help. please verify your key= parameter.",
	403: "access is denied. retrying will not help. please verify your key= parameter.",
	404: "the api requested does not exists.",
	405: "this api has been called with a the wrong http method like get or push.",
	429: "you are being rate limited.",
	500: "an unrecoverable error has occurred, please try again. if this continues to persist then please post to the steamworks developer discussion with additional details of your request.",
	503: "server is temporarily unavailable, or too busy to respond. please wait and try again later.",
}

type Steam struct {
	Key        string        // API key
	LogChannel chan Log      // Channel to return call URLs
	UserAgent  string        // User agent
	APIRate    time.Duration // Min time between each call, defaults to 1 second
	StoreRate  time.Duration // Min time between each call, defaults to 1 second

	apiThrottle   *time.Ticker
	storeThrottle *time.Ticker
}

func (s *Steam) setup() {
	if s.APIRate <= 0 {
		s.APIRate = time.Second
	}
	if s.StoreRate <= 0 {
		s.StoreRate = time.Second
	}
	if s.apiThrottle == nil {
		s.apiThrottle = time.NewTicker(s.APIRate)
	}
	if s.storeThrottle == nil {
		s.storeThrottle = time.NewTicker(s.StoreRate)
	}
	if s.UserAgent == "" {
		s.UserAgent = defaultUserAgent
	}
}

func (s Steam) getFromAPI(path string, query url.Values) (bytes []byte, err error) {

	if s.Key == "" {
		return bytes, errors.New("missing api key")
	}

	// Throttle
	s.setup()
	<-s.apiThrottle.C

	query.Set("format", "json")
	query.Set("key", s.Key)

	path = "https://api.steampowered.com/" + path + "?" + query.Encode()

	// Create request
	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return bytes, err
	}

	req.Header.Set("User-Agent", s.UserAgent)

	start := time.Now()
	response, err := client.Do(req)
	elapsed := time.Since(start)

	// Send log
	if s.LogChannel != nil {
		s.LogChannel <- Log{path, response.StatusCode, elapsed}
	}

	//
	if err != nil {
		return bytes, err
	}
	defer response.Body.Close()

	// Handle errors
	if response.StatusCode != 200 {
		if val, ok := statusCodes[response.StatusCode]; ok {
			return bytes, Error{err: val, code: response.StatusCode, url: path}
		} else {
			return bytes, errors.New("steam: something went wrong")
		}
	}

	//
	return ioutil.ReadAll(response.Body)
}

func (s Steam) getFromStore(path string, query url.Values) (bytes []byte, err error) {

	// Throttle
	s.setup()
	<-s.storeThrottle.C

	path = "https://store.steampowered.com/" + path + "?" + query.Encode()
	fmt.Println(path)

	// Create request
	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return bytes, err
	}

	req.Header.Set("User-Agent", s.UserAgent)

	start := time.Now()
	response, err := client.Do(req)
	elapsed := time.Since(start)

	// Send log
	if s.LogChannel != nil {
		s.LogChannel <- Log{path, response.StatusCode, elapsed}
	}

	//
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

type Log struct {
	Path string
	Code int
	Time time.Duration
}

func (l Log) String() string {
	return strconv.Itoa(l.Code) + " " + l.Path + " " + l.Time.String()
}

type Error struct {
	err  string
	code int
	url  string
}

func (e Error) Error() string {
	return "steam-go: (" + strconv.Itoa(e.code) + ") " + e.err + " (" + e.url + ")"
}

func (e Error) Code() int {
	return e.code
}

package steam

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/juju/ratelimit"
)

const defaultUserAgent = "github.com/Jleagle/steam-go"

var (
	ErrMissingKey = errors.New("missing api key")

	apiStatusCodes = map[int]string{
		400: "please verify that all required parameters are being sent.",
		401: "access is denied. retrying will not help. please verify your key= parameter.",
		403: "access is denied. retrying will not help. please verify your key= parameter.",
		404: "the api requested does not exists.",
		405: "this api has been called with a the wrong http method like get or push.",
		429: "you are being rate limited.",
		500: "an unrecoverable error has occurred, please try again.",
		503: "server is temporarily unavailable, or too busy to respond. please wait and try again later.",
	}
)

type Steam struct {
	key             string
	logger          logger
	userAgent       string
	apiBucket       *ratelimit.Bucket
	storeBucket     *ratelimit.Bucket
	communityBucket *ratelimit.Bucket
}

func (s *Steam) SetKey(key string) {
	s.key = key
}

func (s *Steam) SetLogger(logger logger) {
	s.logger = logger
}

func (s *Steam) SetUserAgent(userAgent string) {
	s.userAgent = userAgent
}

func (s *Steam) SetAPIRateLimit(duration time.Duration, burst int64) {
	s.apiBucket = ratelimit.NewBucket(duration, burst)
}

func (s *Steam) SetStoreRateLimit(duration time.Duration, burst int64) {
	s.storeBucket = ratelimit.NewBucket(duration, burst)
}

func (s *Steam) SetCommunityRateLimit(duration time.Duration, burst int64) {
	s.communityBucket = ratelimit.NewBucket(duration, burst)
}

func (s Steam) getFromAPI(path string, query url.Values) (bytes []byte, err error) {

	if s.key == "" {
		return bytes, ErrMissingKey
	}

	if s.apiBucket != nil {
		s.apiBucket.Wait(1)
	}

	query.Set("format", "json")
	query.Set("key", s.key)

	response, err := s.get("https://api.steampowered.com/" + path + "?" + query.Encode())
	if err != nil {
		return bytes, err
	}

	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	// Handle errors
	if response.StatusCode != 200 {
		if val, ok := apiStatusCodes[response.StatusCode]; ok {
			return bytes, Error{Err: val, Code: response.StatusCode, URL: path}
		} else {
			return bytes, errors.New("steam: something went wrong")
		}
	}

	//
	return ioutil.ReadAll(response.Body)
}

func (s Steam) getFromStore(path string, query url.Values) (bytes []byte, err error) {

	if s.storeBucket != nil {
		s.storeBucket.Wait(1)
	}

	response, err := s.get("https://store.steampowered.com/" + path + "?" + query.Encode())
	if err != nil {
		return bytes, err
	}

	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

func (s Steam) getFromCommunity(path string, query url.Values) (resp *http.Response, err error) {

	if s.communityBucket != nil {
		s.communityBucket.Wait(1)
	}

	// Response body gets closed in caller function
	return s.get("https://steamcommunity.com/" + path + "?" + query.Encode())
}

func (s Steam) get(path string) (response *http.Response, err error) {

	// Create request
	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return response, err
	}

	if s.userAgent == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	} else {
		req.Header.Set("User-Agent", s.userAgent)
	}

	start := time.Now()
	response, err = client.Do(req)
	elapsed := time.Since(start)

	// Send log
	if s.logger != nil && response != nil {
		s.logger.Write(Log{path, response.StatusCode, elapsed})
	}

	//
	return response, err
}

type logger interface {
	Write(log Log)
}

type DefaultLogger struct {
}

func (l DefaultLogger) Write(log Log) {
	fmt.Println(log.String())
}

type Log struct {
	Path     string
	Code     int
	Duration time.Duration
}

func (l Log) String() string {
	return strconv.Itoa(l.Code) + " " + l.Path + " " + l.Duration.String()
}

type Error struct {
	Err  string
	Code int
	URL  string
}

func (e Error) Error() string {
	return "steam-go: (" + strconv.Itoa(e.Code) + ") " + e.Err + " (" + e.URL + ")"
}

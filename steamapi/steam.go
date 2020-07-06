package steamapi

import (
	"bytes"
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

func (s Steam) getFromAPI(path string, query url.Values, key bool) (b []byte, err error) {

	if s.key == "" {
		return b, ErrMissingKey
	}

	if s.apiBucket != nil {
		s.apiBucket.Wait(1)
	}

	query.Set("format", "json")

	if key {
		query.Set("key", s.key)
	}

	b, code, _, err := s.get("https://api.steampowered.com/" + path + "?" + query.Encode())
	if err != nil {
		return b, err
	}

	// Handle errors
	if code != 200 {
		if val, ok := apiStatusCodes[code]; ok {
			return b, Error{Err: val, Code: code, URL: path}
		} else {
			return b, Error{Err: "something went wrong", Code: 0, URL: path}
		}
	}

	return b, err
}

func (s Steam) getFromStore(path string, query url.Values) (b []byte, err error) {

	if s.storeBucket != nil {
		s.storeBucket.Wait(1)
	}

	b, _, _, err = s.get("https://store.steampowered.com/" + path + "?" + query.Encode())
	if err != nil {
		return b, err
	}

	return b, err
}

func (s Steam) getFromCommunity(path string, query url.Values) (b []byte, url string, err error) {

	if s.communityBucket != nil {
		s.communityBucket.Wait(1)
	}

	if query != nil {
		path += "?" + query.Encode()
	}

	b, _, url, err = s.get("https://steamcommunity.com/" + path)
	return b, url, err
}

func (s Steam) get(path string) (b []byte, code int, url string, err error) {

	// Create request
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return b, code, url, err
	}

	if s.userAgent == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	} else {
		req.Header.Set("User-Agent", s.userAgent)
	}

	start := time.Now()
	response, err := client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		return b, code, url, err
	}

	defer func() {
		err = req.Body.Close()
		c.logger.Err(err)
	}()

	b, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return b, code, url, err
	}

	b = bytes.TrimSpace(b)
	code = response.StatusCode
	url = response.Request.URL.Path

	c.logger.Info(path)

	//
	return b, code, url, err
}

type logger interface {
	Info(string)
	Err(error)
}

type DefaultLogger struct {
}

func (l DefaultLogger) Info(s string) {
	fmt.Println("INFO: " + s)
}

func (l DefaultLogger) Err(e error) {
	fmt.Println("ERROR: " + e.Error())
}

type Error struct {
	Err  string
	Code int
	URL  string
}

func (e Error) Error() string {
	return "steam-go: (" + strconv.Itoa(e.Code) + ") " + e.Err + " (" + e.URL + ")"
}

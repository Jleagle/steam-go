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

func NewClient() *Client {
	c := &Client{}
	c.SetLogger(DefaultLogger{})
	c.SetUserAgent("github.com/Jleagle/steam-go")
	return c
}

type Client struct {
	key             string
	logger          logger
	userAgent       string
	apiBucket       *ratelimit.Bucket
	storeBucket     *ratelimit.Bucket
	communityBucket *ratelimit.Bucket
}

func (c *Client) SetKey(key string) {
	c.key = key
}

func (c *Client) SetLogger(logger logger) {
	c.logger = logger
}

func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func (c *Client) SetAPIRateLimit(duration time.Duration, burst int64) {
	c.apiBucket = ratelimit.NewBucket(duration, burst)
}

func (c *Client) SetStoreRateLimit(duration time.Duration, burst int64) {
	c.storeBucket = ratelimit.NewBucket(duration, burst)
}

func (c *Client) SetCommunityRateLimit(duration time.Duration, burst int64) {
	c.communityBucket = ratelimit.NewBucket(duration, burst)
}

func (c Client) getFromAPI(path string, query url.Values, key bool) (b []byte, err error) {

	if c.key == "" {
		return b, ErrMissingKey
	}

	if c.apiBucket != nil {
		c.apiBucket.Wait(1)
	}

	query.Set("format", "json")

	if key {
		query.Set("key", c.key)
	}

	b, code, _, err := c.get("https://api.steampowered.com/" + path + "?" + query.Encode())
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

func (c Client) getFromStore(path string, query url.Values) (b []byte, err error) {

	if c.storeBucket != nil {
		c.storeBucket.Wait(1)
	}

	b, _, _, err = c.get("https://store.steampowered.com/" + path + "?" + query.Encode())
	if err != nil {
		return b, err
	}

	return b, err
}

func (c Client) getFromCommunity(path string, query url.Values) (b []byte, url string, err error) {

	if c.communityBucket != nil {
		c.communityBucket.Wait(1)
	}

	if query != nil {
		path += "?" + query.Encode()
	}

	b, _, url, err = c.get("https://steamcommunity.com/" + path)
	return b, url, err
}

func (c Client) get(path string) (b []byte, code int, url string, err error) {

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return b, code, url, err
	}

	req.Header.Set("User-Agent", c.userAgent)

	response, err := client.Do(req)
	if err != nil {
		return b, code, url, err
	}

	defer func() {
		err = response.Body.Close()
		c.logger.Err(err)
	}()

	b, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return b, code, url, err
	}

	b = bytes.TrimSpace(b)
	code = response.StatusCode
	url = response.Request.URL.Path

	c.logger.Info(path)

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

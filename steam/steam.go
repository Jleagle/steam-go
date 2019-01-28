package steam

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/juju/ratelimit"
)

const defaultUserAgent = "github.com/Jleagle/steam-go"

var (
	errMissingKey = errors.New("missing api key")
	statusCodes   = map[int]string{
		400: "please verify that all required parameters are being sent.",
		401: "access is denied. retrying will not help. please verify your key= parameter.",
		403: "access is denied. retrying will not help. please verify your key= parameter.",
		404: "the api requested does not exists.",
		405: "this api has been called with a the wrong http method like get or push.",
		429: "you are being rate limited.",
		500: "an unrecoverable error has occurred, please try again. if this continues to persist then please post to the steamworks developer discussion with additional details of your request.",
		503: "server is temporarily unavailable, or too busy to respond. please wait and try again later.",
	}
)

type Steam struct {
	key         string
	logChannel  chan Log
	userAgent   string
	apiBucket   *ratelimit.Bucket
	storeBucket *ratelimit.Bucket
}

func (s *Steam) SetKey(key string) {
	s.key = key
}

func (s *Steam) SetLogChannel(c chan Log) {
	s.logChannel = c
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

func (s Steam) getFromAPI(path string, query url.Values) (bytes []byte, err error) {

	if s.key == "" {
		return bytes, errMissingKey
	}

	if s.apiBucket != nil {
		s.apiBucket.Wait(1)
	}

	query.Set("format", "json")
	query.Set("key", s.key)

	path = "https://api.steampowered.com/" + path + "?" + query.Encode()

	// Create request
	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return bytes, err
	}

	if s.userAgent == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	} else {
		req.Header.Set("User-Agent", s.userAgent)
	}

	start := time.Now()
	response, err := client.Do(req)
	elapsed := time.Since(start)

	// Send log
	if s.logChannel != nil {
		s.logChannel <- Log{path, response.StatusCode, elapsed}
	}

	//
	if err != nil {
		return bytes, err
	}

	//noinspection GoUnhandledErrorResult
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

	if s.storeBucket != nil {
		s.storeBucket.Wait(1)
	}

	path = "https://store.steampowered.com/" + path + "?" + query.Encode()

	// Create request
	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return bytes, err
	}

	if s.userAgent == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	} else {
		req.Header.Set("User-Agent", s.userAgent)
	}

	start := time.Now()
	response, err := client.Do(req)
	elapsed := time.Since(start)

	// Send log
	if s.logChannel != nil {
		s.logChannel <- Log{path, response.StatusCode, elapsed}
	}

	//
	if err != nil {
		return nil, err
	}

	//noinspection GoUnhandledErrorResult
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

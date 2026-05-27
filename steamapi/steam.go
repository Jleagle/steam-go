// Package steamapi provides a client for the Steam Web API and Store API.
package steamapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/juju/ratelimit"
)

const (
	defaultAPIBaseURL       = "https://api.steampowered.com/"
	defaultStoreBaseURL     = "https://store.steampowered.com/"
	defaultCommunityBaseURL = "https://steamcommunity.com/"
)

var (
	// ErrMissingKey is returned when a Steam Web API key is missing.
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

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithKey sets the Steam Web API key.
func WithKey(key string) Option {
	return func(c *Client) {
		c.key = key
	}
}

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

// WithLogger sets a custom slog.Logger.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithUserAgent sets a custom User-Agent header for requests.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithAPIBaseURL overrides the default Steam Web API base URL.
func WithAPIBaseURL(url string) Option {
	return func(c *Client) {
		c.apiBaseURL = url
	}
}

// WithStoreBaseURL overrides the default Steam Store API base URL.
func WithStoreBaseURL(url string) Option {
	return func(c *Client) {
		c.storeBaseURL = url
	}
}

// WithCommunityBaseURL overrides the default Steam Community base URL.
func WithCommunityBaseURL(url string) Option {
	return func(c *Client) {
		c.communityBaseURL = url
	}
}

// NewClient creates a new Steam API client with the given options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		userAgent:        "github.com/Jleagle/steam-go",
		client:           http.DefaultClient,
		logger:           slog.New(slog.NewTextHandler(io.Discard, nil)),
		apiBaseURL:       defaultAPIBaseURL,
		storeBaseURL:     defaultStoreBaseURL,
		communityBaseURL: defaultCommunityBaseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Client is a Steam API client.
type Client struct {
	key              string
	userAgent        string
	logger           *slog.Logger
	client           *http.Client
	apiBucket        *ratelimit.Bucket
	storeBucket      *ratelimit.Bucket
	communityBucket  *ratelimit.Bucket
	apiBaseURL       string
	storeBaseURL     string
	communityBaseURL string
}

// SetAPIRateLimit sets the rate limit for Steam Web API requests.
func (c *Client) SetAPIRateLimit(duration time.Duration, burst int64) {
	c.apiBucket = ratelimit.NewBucket(duration, burst)
}

// SetStoreRateLimit sets the rate limit for Steam Store API requests.
func (c *Client) SetStoreRateLimit(duration time.Duration, burst int64) {
	c.storeBucket = ratelimit.NewBucket(duration, burst)
}

// SetCommunityRateLimit sets the rate limit for Steam Community requests.
func (c *Client) SetCommunityRateLimit(duration time.Duration, burst int64) {
	c.communityBucket = ratelimit.NewBucket(duration, burst)
}

func (c *Client) getFromAPI(ctx context.Context, path string, query url.Values, key bool) (b []byte, err error) {

	if c.key == "" && key {
		return b, ErrMissingKey
	}

	if c.apiBucket != nil {
		c.apiBucket.Wait(1)
	}

	query.Set("format", "json")

	if key {
		query.Set("key", c.key)
	}

	b, code, _, err := c.get(ctx, c.apiBaseURL+path+"?"+query.Encode())
	if err != nil {
		return b, fmt.Errorf("api request failed: %w", err)
	}

	// Handle errors
	if code != 200 {
		if val, ok := apiStatusCodes[code]; ok {
			return b, Error{Err: val, Code: code, URL: path}
		}
		return b, Error{Err: "something went wrong", Code: code, URL: path}
	}

	return b, err
}

func (c *Client) getFromStore(ctx context.Context, path string, query url.Values) (b []byte, err error) {

	if c.storeBucket != nil {
		c.storeBucket.Wait(1)
	}

	b, _, _, err = c.get(ctx, c.storeBaseURL+path+"?"+query.Encode())
	if err != nil {
		return b, fmt.Errorf("store request failed: %w", err)
	}

	return b, err
}

func (c *Client) getFromCommunity(ctx context.Context, path string, query url.Values) (b []byte, url string, err error) {

	if c.communityBucket != nil {
		c.communityBucket.Wait(1)
	}

	if query != nil {
		path += "?" + query.Encode()
	}

	b, _, url, err = c.get(ctx, c.communityBaseURL+path)
	if err != nil {
		return b, url, fmt.Errorf("community request failed: %w", err)
	}

	return b, url, err
}

func (c *Client) get(ctx context.Context, path string) (b []byte, code int, url string, err error) {

	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return b, code, url, err
	}

	req.Header.Set("User-Agent", c.userAgent)

	response, err := c.client.Do(req)
	if err != nil {
		return b, code, url, err
	}

	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			c.logger.Error("failed to close response body", "error", closeErr)
		}
	}()

	b, err = io.ReadAll(response.Body)
	if err != nil {
		return b, code, url, err
	}

	b = bytes.TrimSpace(b)
	b = bytes.TrimPrefix(b, []byte{239, 187, 191}) // Trim byte order mark

	code = response.StatusCode
	url = response.Request.URL.Path

	c.logger.Info("steam api request", "path", path, "code", code)

	return b, code, url, err
}

// Error represents a Steam API error.
type Error struct {
	Err  string
	Code int
	URL  string
}

func (e Error) Error() string {
	return "steam-go: (" + strconv.Itoa(e.Code) + ") " + e.Err + " (" + e.URL + ")"
}

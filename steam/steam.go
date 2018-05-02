package steam

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/steam-authority/steam-authority/logger"
)

const (
	ErrInvalidJson = "invalid character '<' looking for beginning of value"
)

const (
	countCalls = false
	logCalls   = false
)

var (
	key  = os.Getenv("STEAM_API_KEY")
	logs = steamLogs{}
)

// todo, handle error codes crrectly.. https://partner.steamgames.com/doc/webapi_overview/responses

func get(path string, query url.Values) (bytes []byte, err error) {

	query.Add("format", "json")
	query.Add("key", key)

	path = "http://api.steampowered.com/" + path + "?" + query.Encode()

	logs.AddLog(path)

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

	if string(contents) == "<html><head><title>Forbidden</title></head><body><h1>Forbidden</h1>Access is denied. Retrying will not help. Please verify your <pre>key=</pre> parameter.</body></html>" {
		errors.New("invalid api key")
	}

	return contents, nil
}

type steamLogs struct {
	logs map[int]steamLog
}

type steamLog struct {
	path string
	time time.Time
}

func (s steamLogs) AddLog(path string) {

	if s.logs == nil {
		s.logs = map[int]steamLog{}
	}

	//path = strings.Replace(path, key, "-", 1)
	//path = strings.Replace(path, "http://api.steampowered.com", "Steam API: ", 1)
	//path = strings.Replace(path, "http://store.steampowered.com", "Store API: ", 1)

	if logCalls {
		logger.Info(path)
	}

	if countCalls {
		s.logs[rand.Int()] = steamLog{
			path: path,
			time: time.Now(),
		}
	}
}

func (s steamLogs) CountStoreCalls(seconds int) (count int) {

	for k, v := range s.logs {
		if v.time.Unix() > (time.Now().Unix() - int64(seconds)) {
			if strings.Contains(v.path, "store.steampowered") {
				count++
			}
		} else {
			delete(s.logs, k)
		}
	}

	return count
}

func (s steamLogs) CountSteamCalls(seconds int) (count int) {

	for k, v := range s.logs {
		if v.time.Unix() > (time.Now().Unix() - int64(seconds)) {
			if strings.Contains(v.path, "api.steampowered.com") {
				count++
			}
		} else {
			delete(s.logs, k)
		}
	}

	return count
}

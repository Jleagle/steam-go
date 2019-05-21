package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Jleagle/steam-go/steam"
)

func TestRateLimits(t *testing.T) {

	steamClient := steam.Steam{}
	steamClient.SetStoreRateLimit(time.Second, 3)

	// Random app IDs
	var ids = []int{408070, 208370, 568670, 595588, 507580, 787980, 994300, 620520, 576250, 759300}
	var wg sync.WaitGroup

	for _, v := range ids {
		wg.Add(1)
		go func(v int) {
			_, _, err := steamClient.GetAppDetails(v, steam.CountryUS, steam.LanguageEnglish)
			fmt.Println(err)
			wg.Done()
		}(v)
	}
	wg.Wait()
}

func TestGroupTypes(t *testing.T) {

	var resp steam.GroupInfo
	// var b []byte
	// var err error

	steamClient := steam.Steam{}

	resp, _, _ = steamClient.GetGroupByID("103582791433980119")
	if resp.Type != "game" {
		t.Error("group type: " + resp.Type)
	}

	resp, _, _ = steamClient.GetGroupByID("103582791429670253")
	if resp.Type != "group" {
		t.Error("group type: " + resp.Type)
	}

}

func TestGroupName(t *testing.T) {

	steamClient := steam.Steam{}

	resp, _, _ := steamClient.GetGroupByID("103582791432805705")
	if resp.Details.Name == "" {
		t.Error("empty name")
	}
}

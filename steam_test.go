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
			_, _, err := steamClient.GetAppDetails(v, "eu", steam.LanguageEnglish, nil)
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

	resp, _, err := steamClient.GetGroupByID("103582791432805705")
	fmt.Println(err)
	if resp.Details.Name == "" {
		t.Error("empty name")
	}
}

func TestDigest(t *testing.T) {

	steamClient := steam.Steam{}
	steamClient.SetKey("x")

	resp, b, err := steamClient.GetItemDefArchive(365960, "6FD2576E9AE2C17F60D08A59D2E7E80F3265BA5B")
	fmt.Println(resp)
	fmt.Println(b)
	fmt.Println(err)
}

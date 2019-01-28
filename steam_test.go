package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Jleagle/steam-go/steam"
)

func TestSteam(t *testing.T) {

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

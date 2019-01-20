package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Jleagle/steam-go/steam"
)

func TestSteam(t *testing.T) {

	var steamClient = &steam.Steam{
		APIRate:   time.Millisecond * 1000,
		StoreRate: time.Millisecond * 1000,
	}

	var wg sync.WaitGroup

	for i := 1; i <= 1; i++ {
		wg.Add(1)
		go func() {
			_, _, err := steamClient.GetAppDetails(440, steam.CountryUS, steam.LanguageEnglish)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

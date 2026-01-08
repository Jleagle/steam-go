package steamapi

import (
	"testing"
)

//func TestWishlist(t *testing.T) {
//
//	c := NewClient()
//
//	_, err := c.GetWishlist(76561198004579722)
//	if err != nil {
//		t.Error(err)
//	}
//}

func TestApps(t *testing.T) {

	c := NewClient()

	// Paid game price overview
	app, err := c.GetAppDetails(252490, ProductCCUS, LanguageEnglish, []string{"price_overview"})
	if err != nil {
		t.Error(err)
	}
	if app.Data == nil {
		t.Error(err)
	}
	if app.Data.PriceOverview.Currency != "USD" {
		t.Error("currency")
	}

	// Free game price overview
	app, err = c.GetAppDetails(440, ProductCCUS, LanguageEnglish, []string{"price_overview"})
	if err != nil {
		t.Error(err)
	}
	if !app.Success {
		t.Error(err)
	}
	if app.Data != nil {
		t.Error(err)
	}

	// Free game
	app, err = c.GetAppDetails(440, ProductCCUS, LanguageEnglish, nil)
	if err != nil {
		t.Error(err)
	}
	if app.Data == nil {
		t.Error(err)
	}
	if app.Data.PriceOverview != nil {
		t.Error("price should be nil")
	}
}

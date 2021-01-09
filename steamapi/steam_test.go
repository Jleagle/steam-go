package steamapi

import (
	"testing"
)

func TestWishlist(t *testing.T) {

	c := NewClient()

	_, err := c.GetWishlist(76561198004579722)
	if err != nil {
		t.Error(err)
	}
}

func TestPlayers(t *testing.T) {

	c := NewClient()

	_, _, err := c.GetAliases(76561197968626192)
	if err != nil {
		t.Error(err)
	}
}

func TestGroups(t *testing.T) {

	c := NewClient()

	group, _, err := c.GetGroup("103582791434672565", "", 1)
	if err != nil {
		t.Error(err)
	}
	if group.Details.Name != "Steam Universe" {
		t.Error("name")
	}
}

func TestApps(t *testing.T) {

	c := NewClient()

	// Paid game price overview
	app, err := c.GetAppDetails(578080, ProductCCUS, LanguageEnglish, []string{"price_overview"})
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

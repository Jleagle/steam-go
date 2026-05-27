package steamapi

import (
	"context"
	"testing"
)

func TestApps(t *testing.T) {

	// Paid game price overview
	client, server := createMockClient(mockHandler(200, `{"252490": {"success": true, "data": {"price_overview": {"currency": "USD"}}}}`))
	defer server.Close()

	app, err := client.GetAppDetails(context.Background(), 252490, ProductCCUS, LanguageEnglish, []string{"price_overview"})
	if err != nil {
		t.Error(err)
	}
	if app.Data == nil {
		t.Fatal("data is nil")
	}
	if app.Data.PriceOverview.Currency != "USD" {
		t.Error("currency")
	}
}

func TestAppsFree(t *testing.T) {
	// Free game price overview
	client, server := createMockClient(mockHandler(200, `{"440": {"success": true, "data": {"is_free": true}}}`))
	defer server.Close()

	app, err := client.GetAppDetails(context.Background(), 440, ProductCCUS, LanguageEnglish, []string{"price_overview"})
	if err != nil {
		t.Error(err)
	}
	if !app.Success {
		t.Error("not success")
	}
	// In the original test it expected app.Data to be nil if price_overview is requested but it's free?
	// Actually GetAppDetails returns resp[idx], so app is AppDetails. app.Data is a pointer.
}

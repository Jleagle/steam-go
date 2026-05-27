package steamapi

import (
	"net/http"
	"net/http/httptest"
)

func createMockClient(handler http.Handler) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client := NewClient(
		WithHTTPClient(server.Client()),
		WithKey("MOCK_KEY"),
		WithAPIBaseURL(server.URL+"/"),
		WithStoreBaseURL(server.URL+"/"),
		WithCommunityBaseURL(server.URL+"/"),
	)
	return client, server
}

func mockHandler(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}
}

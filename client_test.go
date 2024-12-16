package aapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests should register handlers on mux which provide mock responses for the API method being tested.
func setup(t *testing.T) (*http.ServeMux, *httptest.Server, *Client) {
	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	c, err := New(server.URL, "token")
	if err != nil {
		server.Close()
		t.Fatalf("Failed to create client: %v", err)
	}

	return mux, server, c
}

func teardown(server *httptest.Server) {
	server.Close()
}

func testRequestMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %s, want %s", got, want)
	}
}

func TestNewClient(t *testing.T) {
	c, err := New("base_url", "token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	expectedBaseURL := "base_url/" + apiVersionPath

	if c.BaseURL().String() != expectedBaseURL {
		t.Errorf("NewClient BaseURL is %s, want %s", c.BaseURL().String(), expectedBaseURL)
	}
}

func TestNewClientEmptyToken(t *testing.T) {
	c, err := New("base_url", "")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	expectedBaseURL := "base_url/" + apiVersionPath

	if c.BaseURL().String() != expectedBaseURL {
		t.Errorf("NewClient BaseURL is %s, want %s", c.BaseURL().String(), expectedBaseURL)
	}
}

func TestNewClientWithGrafanaURL(t *testing.T) {
	c, err := NewWithGrafanaURL("base_url", "token", "grafana_url")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	expectedBaseURL := "base_url/" + apiVersionPath

	if c.BaseURL().String() != expectedBaseURL {
		t.Errorf("NewClient BaseURL is %s, want %s", c.BaseURL().String(), expectedBaseURL)
	}

	if c.GrafanaURL().String() != "grafana_url" {
		t.Errorf("NewClient GrafanaURL is %s, want grafana_url", c.GrafanaURL().String())
	}
}

func TestNewClientWithGrafanaURLEmptyToken(t *testing.T) {
	c, err := NewWithGrafanaURL("base_url", "", "grafana_url")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	expectedBaseURL := "base_url/" + apiVersionPath

	if c.BaseURL().String() != expectedBaseURL {
		t.Errorf("NewClient BaseURL is %s, want %s", c.BaseURL().String(), expectedBaseURL)
	}

	if c.GrafanaURL().String() != "grafana_url" {
		t.Errorf("NewClient GrafanaURL is %s, want grafana_url", c.GrafanaURL().String())
	}
}

func TestCheckRequest(t *testing.T) {
	c, err := New("base_url", "token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := c.NewRequest("GET", "test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if req.Header.Get("X-Grafana-URL") != "" {
		t.Errorf("X-Grafana-URL should not be set: %s", req.Header.Get("X-Grafana-URL"))
	}
}

func TestCheckRequestSettingGrafanaURL(t *testing.T) {
	c, err := NewWithGrafanaURL("base_url", "token", "grafana_url")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := c.NewRequest("GET", "test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if req.Header.Get("X-Grafana-URL") != "grafana_url" {
		t.Errorf("X-Grafana-URL is not set correctly: %s", req.Header.Get("X-Grafana-URL"))
	}
}

func TestCheckResponse(t *testing.T) {
	c, err := New("base_url", "token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := c.NewRequest("GET", "test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp := &http.Response{
		Request:    req.Request,
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`
		{
			"detail": "error"
		}`)),
	}

	errResp := CheckResponse(resp)
	if errResp == nil {
		t.Fatal("Expected error response.")
	}

	want := fmt.Sprintf("GET ://%s: 400 {detail: error}", req.URL)

	if errResp.Error() != want {
		t.Errorf("Expected error: %s, got %s", want, errResp.Error())
	}
}

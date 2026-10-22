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

type recordingTransport struct {
	next     http.RoundTripper
	requests []string
}

func (rt *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.requests = append(rt.requests, req.Method+" "+req.URL.Path)
	next := rt.next
	if next == nil {
		next = http.DefaultTransport
	}
	return next.RoundTrip(req)
}

func TestSetHTTPClient(t *testing.T) {
	mux, server, c := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/users/U1/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": "U1"}`)
	})

	rt := &recordingTransport{}
	c.SetHTTPClient(&http.Client{Transport: rt})

	if _, _, err := c.Users.GetUser("U1", &GetUserOptions{}); err != nil {
		t.Fatalf("GetUser returned error: %v", err)
	}
	if len(rt.requests) != 1 || rt.requests[0] != "GET /api/v1/users/U1/" {
		t.Errorf("custom transport saw %v, want [GET /api/v1/users/U1/]", rt.requests)
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

func TestRatelimitRetry(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	requestCount := 0
	mux.HandleFunc("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	req, err := client.NewRequest("GET", "test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	_, err = client.Do(req, nil)
	expectedErr := fmt.Sprintf("GET %s giving up after 5 attempt(s)", req.URL)
	if err.Error() != expectedErr {
		t.Fatalf("Expected error: %s, got %s", expectedErr, err.Error())
	}

	if requestCount != 5 {
		t.Errorf("Expected 5 requests (1 original + 4 retries), got %d", requestCount)
	}
}

package aapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// settingsHandler writes a grafana-irm-app plugin settings response whose
// jsonData.onCallApiUrl is onCallURL.
func settingsHandler(onCallURL string, capture *http.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if capture != nil {
			*capture = r.Header.Clone()
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonData": map[string]any{"onCallApiUrl": onCallURL},
		})
	}
}

func expectedBaseURL(rawURL string) string {
	return rawURL + "/" + apiVersionPath
}

func TestAutodiscoverySuccess(t *testing.T) {
	mux := http.NewServeMux()
	var gotHeaders http.Header
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", settingsHandler(server.URL+"/oncall", &gotHeaders))

	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	// No network at construction.
	if c.baseURL != nil {
		t.Fatal("baseURL should be unresolved before EnsureBaseURL")
	}

	if err := c.EnsureBaseURL(context.Background()); err != nil {
		t.Fatalf("EnsureBaseURL error: %v", err)
	}

	want := expectedBaseURL(server.URL + "/oncall")
	if got := c.BaseURL().String(); got != want {
		t.Errorf("BaseURL = %s, want %s", got, want)
	}
	if w := c.Warnings(); len(w) != 0 {
		t.Errorf("expected no warnings, got %v", w)
	}
	if got := gotHeaders.Get("Authorization"); got != "Bearer glsa_grafana" {
		t.Errorf("discovery Authorization = %q, want %q", got, "Bearer glsa_grafana")
	}
}

func TestAutodiscoveryFallbackToExplicitURLOnLookupFailure(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	explicit := "https://oncall.example.com/oncall"
	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", explicit)
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	if err := c.EnsureBaseURL(context.Background()); err != nil {
		t.Fatalf("EnsureBaseURL error: %v", err)
	}

	if got, want := c.BaseURL().String(), expectedBaseURL(explicit); got != want {
		t.Errorf("BaseURL = %s, want %s", got, want)
	}

	warnings := c.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
	}
	// Warnings drain: a second read returns nothing.
	if again := c.Warnings(); len(again) != 0 {
		t.Errorf("warnings should drain, got %v", again)
	}
}

func TestAutodiscoveryFallbackToDefaultOnLookupFailure(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	if err := c.EnsureBaseURL(context.Background()); err != nil {
		t.Fatalf("EnsureBaseURL error: %v", err)
	}

	if got, want := c.BaseURL().String(), expectedBaseURL(defaultOnCallURL); got != want {
		t.Errorf("BaseURL = %s, want %s", got, want)
	}
	if w := c.Warnings(); len(w) != 1 {
		t.Errorf("expected 1 warning, got %v", w)
	}
}

func TestAutodiscoveryExplicitURLWithoutGrafana(t *testing.T) {
	// No grafana URL/auth -> no discovery attempt, no warning, explicit URL used.
	explicit := "https://oncall.example.com/oncall"
	c, err := NewWithGrafanaAutodiscovery("", "", "oncall_token", explicit)
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	if err := c.EnsureBaseURL(context.Background()); err != nil {
		t.Fatalf("EnsureBaseURL error: %v", err)
	}

	if got, want := c.BaseURL().String(), expectedBaseURL(explicit); got != want {
		t.Errorf("BaseURL = %s, want %s", got, want)
	}
	if w := c.Warnings(); len(w) != 0 {
		t.Errorf("expected no warnings, got %v", w)
	}
}

func TestAutodiscoveryTwoTokens(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var discoveryAuth string
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
		discoveryAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonData": map[string]any{"onCallApiUrl": server.URL},
		})
	})

	var apiAuth string
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		apiAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":0,"results":[]}`))
	})

	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	req, err := c.NewRequest("GET", "users", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}
	if _, err := c.Do(req, nil); err != nil {
		t.Fatalf("Do error: %v", err)
	}

	if discoveryAuth != "Bearer glsa_grafana" {
		t.Errorf("discovery used Authorization %q, want %q", discoveryAuth, "Bearer glsa_grafana")
	}
	// OnCall API calls use the raw OnCall token (existing amixr convention).
	if apiAuth != "oncall_token" {
		t.Errorf("OnCall API call used Authorization %q, want %q", apiAuth, "oncall_token")
	}
}

func TestAutodiscoveryOnCallTokenFallsBackToGrafanaToken(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonData": map[string]any{"onCallApiUrl": server.URL},
		})
	})

	var apiAuth string
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		apiAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":0,"results":[]}`))
	})

	// Empty OnCall token -> OnCall API calls use the Grafana auth token.
	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "", "")
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	req, err := c.NewRequest("GET", "users", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}
	if _, err := c.Do(req, nil); err != nil {
		t.Fatalf("Do error: %v", err)
	}

	if apiAuth != "glsa_grafana" {
		t.Errorf("OnCall API call used Authorization %q, want %q", apiAuth, "glsa_grafana")
	}
}

func TestAutodiscoveryLazyOnNewRequest(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", settingsHandler(server.URL+"/oncall", nil))

	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	// Never call EnsureBaseURL; NewRequest must resolve the base URL itself.
	if _, err := c.NewRequest("GET", "users", nil); err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}
	if got, want := c.BaseURL().String(), expectedBaseURL(server.URL+"/oncall"); got != want {
		t.Errorf("BaseURL = %s, want %s", got, want)
	}
	_ = fmt.Sprint(c.GrafanaURL())
}

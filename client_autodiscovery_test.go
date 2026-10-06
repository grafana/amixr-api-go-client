package aapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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
	var lookups atomic.Int32
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
		lookups.Add(1)
		settingsHandler(server.URL+"/oncall", &gotHeaders)(w, r)
	})

	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	// No network at construction.
	if got := lookups.Load(); got != 0 {
		t.Fatalf("plugin settings lookups before EnsureBaseURL = %d, want 0", got)
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

func TestAutodiscoveryErrorsOnLookupFailureWithoutExplicitURL(t *testing.T) {
	tests := map[string]http.HandlerFunc{
		"plugin not found": func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "not found", http.StatusNotFound)
		},
		"no onCallApiUrl": settingsHandler("", nil),
	}

	for name, handler := range tests {
		t.Run(name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()
			var lookups atomic.Int32
			mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
				lookups.Add(1)
				handler(w, r)
			})

			c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
			if err != nil {
				t.Fatalf("constructor error: %v", err)
			}

			err = c.EnsureBaseURL(context.Background())
			if err == nil {
				t.Fatal("expected EnsureBaseURL error, got nil")
			}
			for _, want := range []string{"plugins.app:access", "oncall_url"} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q should mention %q", err, want)
				}
			}

			// The default OnCall URL is never used.
			if u := c.BaseURL(); u != nil {
				t.Errorf("BaseURL = %s, want nil", u)
			}
			if _, err := c.NewRequest("GET", "users", nil); err == nil {
				t.Error("expected NewRequest error, got nil")
			}

			// The failure is cached: no second lookup.
			if err := c.EnsureBaseURL(context.Background()); err == nil {
				t.Error("expected cached EnsureBaseURL error, got nil")
			}
			if got := lookups.Load(); got != 1 {
				t.Errorf("plugin settings lookups = %d, want 1", got)
			}
			if w := c.Warnings(); len(w) != 0 {
				t.Errorf("expected no warnings, got %v", w)
			}
		})
	}
}

func TestAutodiscoveryLegacyDefaultURL(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	var lookups atomic.Int32
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
		lookups.Add(1)
	})

	// Without a Grafana URL or auth token no lookup is possible, so the default is used.
	tests := map[string]struct{ grafanaURL, grafanaAuthToken string }{
		"no grafana URL":        {"", "glsa_grafana"},
		"no grafana auth token": {server.URL, ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c, err := NewWithGrafanaAutodiscovery(tt.grafanaURL, tt.grafanaAuthToken, "oncall_token", "")
			if err != nil {
				t.Fatalf("constructor error: %v", err)
			}

			if err := c.EnsureBaseURL(context.Background()); err != nil {
				t.Fatalf("EnsureBaseURL error: %v", err)
			}

			if got, want := c.BaseURL().String(), expectedBaseURL(defaultOnCallURL); got != want {
				t.Errorf("BaseURL = %s, want %s", got, want)
			}
			if w := c.Warnings(); len(w) != 0 {
				t.Errorf("expected no warnings, got %v", w)
			}
		})
	}

	if got := lookups.Load(); got != 0 {
		t.Errorf("plugin settings lookups = %d, want 0", got)
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

func TestAutodiscoveryConcurrentResolution(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var lookups atomic.Int32
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
		lookups.Add(1)
		// Keep the lookup in flight so other goroutines read the base URL meanwhile.
		time.Sleep(50 * time.Millisecond)
		settingsHandler(server.URL+"/oncall", nil)(w, r)
	})

	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	var wg sync.WaitGroup
	errs := make(chan error, 50)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				errs <- c.EnsureBaseURL(context.Background())
				return
			}
			_, err := c.NewRequest("GET", "users", nil)
			errs <- err
		}(i)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Errorf("concurrent call error: %v", err)
		}
	}
	if got := lookups.Load(); got != 1 {
		t.Errorf("plugin settings lookups = %d, want 1", got)
	}
	if got, want := c.BaseURL().String(), expectedBaseURL(server.URL+"/oncall"); got != want {
		t.Errorf("BaseURL = %s, want %s", got, want)
	}
}

func TestAutodiscoveryInvalidPluginURL(t *testing.T) {
	for _, pluginURL := range []string{"oncall.example.com", "/oncall", "http://a b"} {
		t.Run(pluginURL, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()
			mux.HandleFunc("/api/plugins/grafana-irm-app/settings", settingsHandler(pluginURL, nil))

			t.Run("without oncallURL", func(t *testing.T) {
				c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
				if err != nil {
					t.Fatalf("constructor error: %v", err)
				}

				err = c.EnsureBaseURL(context.Background())
				if err == nil || !strings.Contains(err.Error(), "invalid onCallApiUrl") {
					t.Errorf("EnsureBaseURL error = %v, want an invalid onCallApiUrl error", err)
				}
			})

			t.Run("with oncallURL", func(t *testing.T) {
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
				if w := c.Warnings(); len(w) != 1 {
					t.Errorf("expected 1 warning, got %v", w)
				}
			})
		})
	}
}

func TestAutodiscoveryInvalidExplicitURL(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	tests := map[string]struct{ grafanaURL, grafanaAuthToken string }{
		"lookup not possible": {"", ""},
		"lookup failed":       {server.URL, "glsa_grafana"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c, err := NewWithGrafanaAutodiscovery(tt.grafanaURL, tt.grafanaAuthToken, "oncall_token", "http://a b")
			if err != nil {
				t.Fatalf("constructor error: %v", err)
			}

			err = c.EnsureBaseURL(context.Background())
			if err == nil || !strings.Contains(err.Error(), "invalid oncall_url") {
				t.Errorf("EnsureBaseURL error = %v, want an invalid oncall_url error", err)
			}
			// No "Falling back to oncall_url" warning when the fallback itself is invalid.
			if w := c.Warnings(); len(w) != 0 {
				t.Errorf("expected no warnings, got %v", w)
			}
		})
	}
}

func TestAutodiscoveryConcurrentFailure(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var lookups atomic.Int32
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", func(w http.ResponseWriter, r *http.Request) {
		lookups.Add(1)
		time.Sleep(50 * time.Millisecond)
		http.Error(w, "not found", http.StatusNotFound)
	})

	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.EnsureBaseURL(context.Background()); err == nil {
				t.Error("expected EnsureBaseURL error, got nil")
			}
		}()
	}
	wg.Wait()

	if got := lookups.Load(); got != 1 {
		t.Errorf("plugin settings lookups = %d, want 1", got)
	}
}

func TestAutodiscoveryIgnoresCallerContext(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/api/plugins/grafana-irm-app/settings", settingsHandler(server.URL+"/oncall", nil))

	c, err := NewWithGrafanaAutodiscovery(server.URL, "glsa_grafana", "oncall_token", "")
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}

	// The cached result must not depend on the first caller's context.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := c.EnsureBaseURL(ctx); err != nil {
		t.Fatalf("EnsureBaseURL error: %v", err)
	}
	if got, want := c.BaseURL().String(), expectedBaseURL(server.URL+"/oncall"); got != want {
		t.Errorf("BaseURL = %s, want %s", got, want)
	}
}

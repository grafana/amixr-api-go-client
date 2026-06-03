package aapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	apiVersionPath   = "api/v1/"
	defaultUserAgent = "amixr-api-go-client"

	// grafanaIRMAppSettingsPath is the Grafana API path that exposes the
	// grafana-irm-app plugin settings, whose jsonData.onCallApiUrl holds the
	// OnCall backend URL provisioned for the stack.
	grafanaIRMAppSettingsPath = "api/plugins/grafana-irm-app/settings"

	// defaultOnCallURL is the last-resort OnCall backend URL used when the URL
	// can neither be discovered from the grafana-irm-app plugin settings nor
	// provided explicitly.
	defaultOnCallURL = "https://oncall-prod-us-central-0.grafana.net/oncall"
)

type ListOptions struct {
	Page int `url:"page,omitempty" json:"page,omitempty"`
}

type PaginatedResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

type Client struct {
	// HTTP client used to communicate with the API.
	client     *retryablehttp.Client
	token      string
	baseURL    *url.URL
	grafanaURL *url.URL
	UserAgent  string

	// Lazy OnCall backend-URL resolution. These are only used by clients built
	// with NewWithGrafanaAutodiscovery; for the other constructors baseURL is
	// set eagerly and these remain zero-valued (resolution never runs).
	grafanaAuthToken  string     // auth used for the grafana-irm-app settings lookup
	configuredBaseURL string     // explicit OnCall URL fallback (e.g. oncall_url)
	baseURLOnce       sync.Once  // guards a single resolution attempt
	warnMu            sync.Mutex // guards warnings
	warnings          []string   // resolution warnings, drained via Warnings()

	// List of Services. Keep in sync with func newClient
	Alerts                *AlertService
	AlertGroups           *AlertGroupService
	Integrations          *IntegrationService
	EscalationChains      *EscalationChainService
	Escalations           *EscalationService
	Users                 *UserService
	Schedules             *ScheduleService
	Routes                *RouteService
	SlackChannels         *SlackChannelService
	UserGroups            *UserGroupService
	OnCallShifts          *OnCallShiftService
	Teams                 *TeamService
	Webhooks              *WebhookService
	UserNotificationRules *UserNotificationRuleService
}

func NewWithGrafanaURL(base_url, token, grafana_url string) (*Client, error) {
	if base_url == "" {
		return nil, fmt.Errorf("BaseUrl required")
	}
	client, err := newClient(grafana_url)
	if err != nil {
		return nil, err
	}
	if err := client.setBaseURL(base_url); err != nil {
		return nil, err
	}
	client.token = token
	return client, nil
}

func New(base_url, token string) (*Client, error) {
	if base_url == "" {
		return nil, fmt.Errorf("BaseUrl required")
	}
	client, err := newClient("")
	if err != nil {
		return nil, err
	}
	if err := client.setBaseURL(base_url); err != nil {
		return nil, err
	}
	client.token = token
	return client, nil
}

// NewWithGrafanaAutodiscovery creates a client whose OnCall backend URL is
// resolved lazily, on first request (or an explicit EnsureBaseURL call), rather
// than at construction. No network call is made here.
//
// On first use the client resolves its base URL in this order:
//  1. the URL provisioned in the grafana-irm-app plugin settings, looked up via
//     grafanaURL + grafanaAuthToken;
//  2. oncallURL, if provided;
//  3. a built-in default.
//
// grafanaAuthToken authenticates the plugin-settings lookup (a Grafana API
// call); it follows the same convention as the Grafana provider's auth field
// ("user:password" for basic auth, otherwise a bearer token). oncallToken
// authenticates OnCall API calls and may differ from grafanaAuthToken.
func NewWithGrafanaAutodiscovery(grafanaURL, grafanaAuthToken, oncallToken, oncallURL string) (*Client, error) {
	client, err := newClient(grafanaURL)
	if err != nil {
		return nil, err
	}
	client.token = oncallToken
	client.grafanaAuthToken = grafanaAuthToken
	client.configuredBaseURL = oncallURL
	return client, nil
}

func newClient(grafana_url string) (*Client, error) {
	c := &Client{}

	// retryablehttp.Client will retry up to 4 times on recoverable errors (429, 5xx, and low-level network errors)
	c.client = retryablehttp.NewClient()

	err := c.setGrafanaURL(grafana_url)
	if err != nil {
		return nil, err
	}

	c.UserAgent = defaultUserAgent

	// Create services. Keep in sync with Client struct
	c.Alerts = NewAlertService(c)
	c.AlertGroups = NewAlertGroupService(c)
	c.Integrations = NewIntegrationService(c)
	c.EscalationChains = NewEscalationChainService(c)
	c.Escalations = NewEscalationService(c)
	c.Users = NewUserService(c)
	c.Schedules = NewScheduleService(c)
	c.Routes = NewRouteService(c)
	c.SlackChannels = NewSlackChannelService(c)
	c.UserGroups = NewUserGroupService(c)
	c.OnCallShifts = NewOnCallShiftService(c)
	c.Teams = NewTeamService(c)
	c.Webhooks = NewWebhookService(c)
	c.UserNotificationRules = NewUserNotificationRuleService(c)

	return c, nil
}

func (c *Client) setBaseURL(urlStr string) error {

	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(baseURL.Path, apiVersionPath) {
		baseURL.Path += apiVersionPath
	}
	c.baseURL = baseURL

	return nil
}

func (c *Client) setGrafanaURL(urlStr string) error {
	if urlStr != "" {
		grafanaUrl, err := url.Parse(urlStr)
		if err != nil {
			return err
		}
		c.grafanaURL = grafanaUrl
	}

	return nil
}

func (c *Client) NewRequest(method, path string, opt interface{}) (*retryablehttp.Request, error) {
	// Safety net for autodiscovery clients whose caller never invoked
	// EnsureBaseURL: resolve lazily here. For clients built with New /
	// NewWithGrafanaURL the base URL is already set and this is a no-op.
	if c.baseURL == nil {
		if err := c.EnsureBaseURL(context.Background()); err != nil {
			return nil, err
		}
	}

	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	// Set the encoded path data
	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + unescaped

	// Create a request specific headers map.
	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")
	reqHeaders.Set("Authorization", c.token)
	if c.grafanaURL != nil {
		reqHeaders.Set("X-Grafana-URL", c.grafanaURL.String())
	}
	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	var body interface{}
	switch {
	case method == "POST" || method == "PUT":
		reqHeaders.Set("Content-Type", "application/json")

		if opt != nil {
			body, err = json.Marshal(opt)
			if err != nil {
				return nil, err
			}
		}
	case opt != nil:
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := retryablehttp.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// Set the request specific headers.
	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
func (c *Client) Do(req *retryablehttp.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		// Even though there was an error, we still return the response
		// in case the caller wants to inspect it further.
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		// Very naive realization if handling errors messages
		var rawError interface{}
		if err := json.Unmarshal(data, &rawError); err != nil {
			errorResponse.Message = "failed to parse unknown error format"
		} else {
			errorResponse.Message = parseError(rawError)
		}
	}
	if err != nil {
		return err
	}

	return errorResponse
}

func parseError(raw interface{}) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []interface{}:
		var errs []string
		for _, v := range raw {
			errs = append(errs, parseError(v))
		}
		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]interface{}:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}
		sort.Strings(errs)
		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

func (c *Client) GrafanaURL() *url.URL {
	if c.grafanaURL == nil {
		return nil
	}
	u := *c.grafanaURL
	return &u
}

// EnsureBaseURL resolves the OnCall backend URL if it has not been resolved yet.
// It is safe to call multiple times and concurrently; resolution runs at most
// once. Resolution is best-effort and always yields a usable URL (falling back
// to an explicit OnCall URL or a built-in default), so it returns a nil error
// today; the error is retained for forward compatibility. Clients constructed
// with a concrete base URL (New / NewWithGrafanaURL) treat this as a no-op.
func (c *Client) EnsureBaseURL(ctx context.Context) error {
	if c.baseURL != nil {
		return nil
	}
	c.baseURLOnce.Do(func() {
		c.resolveBaseURL(ctx)
	})
	return nil
}

// resolveBaseURL determines the OnCall backend URL: it prefers the URL
// provisioned in the grafana-irm-app plugin settings, then an explicitly
// configured URL, then the built-in default. A failed plugin lookup is recorded
// as a warning (see Warnings) and resolution falls through.
func (c *Client) resolveBaseURL(ctx context.Context) {
	if c.grafanaURL != nil && c.grafanaAuthToken != "" {
		if pluginURL, err := c.fetchOnCallURLFromPlugin(ctx); err != nil {
			c.addWarning(fmt.Sprintf(
				"Could not determine the OnCall backend URL from the grafana-irm-app plugin settings: %s. "+
					"Ensure the Grafana IRM/OnCall app is installed for this stack and that the token has the plugins.app:access permission, "+
					"or set the oncall_url provider attribute explicitly. Falling back to oncall_url or the default OnCall URL.",
				err,
			))
		} else if pluginURL != "" {
			if err := c.setBaseURL(pluginURL); err == nil {
				return
			}
		}
	}

	if c.configuredBaseURL != "" {
		if err := c.setBaseURL(c.configuredBaseURL); err == nil {
			return
		}
	}

	// Last resort: the historical default. setBaseURL only fails on a parse
	// error, which cannot happen for this constant.
	_ = c.setBaseURL(defaultOnCallURL)
}

// fetchOnCallURLFromPlugin reads jsonData.onCallApiUrl from the grafana-irm-app
// plugin settings at GET {grafanaURL}/api/plugins/grafana-irm-app/settings.
func (c *Client) fetchOnCallURLFromPlugin(ctx context.Context) (string, error) {
	settingsURL, err := url.JoinPath(c.grafanaURL.String(), grafanaIRMAppSettingsPath)
	if err != nil {
		return "", fmt.Errorf("failed to build plugin settings URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, settingsURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build plugin settings request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	// Grafana auth follows the provider convention: "user:password" -> basic
	// auth, otherwise a bearer token. "anonymous"/empty -> no auth header.
	if parts := strings.SplitN(c.grafanaAuthToken, ":", 2); len(parts) == 2 {
		req.SetBasicAuth(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	} else if token := strings.TrimSpace(parts[0]); token != "" && token != "anonymous" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := pluginSettingsHTTPClient().Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request plugin settings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("plugin settings request returned status %d", resp.StatusCode)
	}

	var settings struct {
		JSONData struct {
			OnCallAPIURL string `json:"onCallApiUrl"`
		} `json:"jsonData"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return "", fmt.Errorf("failed to decode plugin settings response: %w", err)
	}
	if settings.JSONData.OnCallAPIURL == "" {
		return "", fmt.Errorf("grafana-irm-app plugin settings did not contain an onCallApiUrl")
	}
	return settings.JSONData.OnCallAPIURL, nil
}

// pluginSettingsHTTPClient returns a short-lived HTTP client for the best-effort
// plugin-settings lookup. Retries are limited so a missing plugin or unreachable
// Grafana fails fast instead of stalling the first OnCall request.
func pluginSettingsHTTPClient() *http.Client {
	rc := retryablehttp.NewClient()
	rc.RetryMax = 2
	rc.Logger = nil
	return rc.StandardClient()
}

func (c *Client) addWarning(msg string) {
	c.warnMu.Lock()
	defer c.warnMu.Unlock()
	c.warnings = append(c.warnings, msg)
}

// Warnings returns and clears any warnings accumulated during OnCall URL
// resolution (for example, a failed plugin lookup that fell back to a default).
// Because it drains, the warnings surface only once even if read by multiple
// callers.
func (c *Client) Warnings() []string {
	c.warnMu.Lock()
	defer c.warnMu.Unlock()
	if len(c.warnings) == 0 {
		return nil
	}
	out := c.warnings
	c.warnings = nil
	return out
}

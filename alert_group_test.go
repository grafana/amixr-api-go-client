package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testAlertGroup = &AlertGroup{
	ID:            "I68T24C13IFW1",
	IntegrationID: "CFRPV98RPR1U8",
	RouteID:       "RIYGUJXCPFHXY",
	AlertsCount:   3,
	State:         "resolved",
	CreatedAt:     "2020-05-19T12:37:01.430444Z",
	ResolvedAt:    "2020-05-19T13:37:01.429805Z",
	Title:         "Memory above 90% threshold",
	Permalinks: map[string]string{
		"slack":    "https://ghostbusters.slack.com/archives/C1H9RESGA/p135854651500008",
		"telegram": "https://t.me/c/5354/1234?thread=1234",
	},
}

var testAlertGroupBody = `{
	"id": "I68T24C13IFW1",
	"integration_id": "CFRPV98RPR1U8",
	"route_id": "RIYGUJXCPFHXY",
	"alerts_count": 3,
	"state": "resolved",
	"created_at": "2020-05-19T12:37:01.430444Z",
	"resolved_at": "2020-05-19T13:37:01.429805Z",
	"acknowledged_at": null,
	"title": "Memory above 90% threshold",
	"permalinks": {
	  "slack": "https://ghostbusters.slack.com/archives/C1H9RESGA/p135854651500008",
	  "telegram": "https://t.me/c/5354/1234?thread=1234"
	}
}`

func TestListAlertGroup(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/alert_groups/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testAlertGroupBody))
	})

	options := &ListAlertGroupOptions{
		AlertGroupID: "I68T24C13IFW1",
	}

	alerts, _, err := client.AlertGroups.ListAlertGroups(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedAlertGroupsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		AlertGroups: []*AlertGroup{
			testAlertGroup,
		},
	}

	if !reflect.DeepEqual(want, alerts) {
		fmt.Println(alerts.AlertGroups[0])
		fmt.Println(want.AlertGroups[0])
		t.Errorf(" returned\n %+v, \nwant\n %+v", alerts, want)
	}
}

func TestValidateTimeRange(t *testing.T) {
	tests := []struct {
		name      string
		timeRange string
		wantErr   bool
	}{
		{
			name:      "valid time range",
			timeRange: "2024-03-20T10:00:00_2024-03-21T10:00:00",
			wantErr:   false,
		},
		{
			name:      "empty time range",
			timeRange: "",
			wantErr:   false,
		},
		{
			name:      "invalid format - missing separator",
			timeRange: "2024-03-20T10:00:002024-03-21T10:00:00",
			wantErr:   true,
		},
		{
			name:      "invalid format - wrong date format",
			timeRange: "2024/03/20T10:00:00_2024/03/21T10:00:00",
			wantErr:   true,
		},
		{
			name:      "invalid time - end before start",
			timeRange: "2024-03-21T10:00:00_2024-03-20T10:00:00",
			wantErr:   true,
		},
		{
			name:      "invalid time - invalid hour",
			timeRange: "2024-03-20T25:00:00_2024-03-21T10:00:00",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTimeRange(tt.timeRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTimeRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListAlertGroupsValidation(t *testing.T) {
	tests := []struct {
		name    string
		options *ListAlertGroupOptions
		wantErr bool
	}{
		{
			name: "valid options",
			options: &ListAlertGroupOptions{
				StartedAt: "2024-03-20T10:00:00_2024-03-21T10:00:00",
			},
			wantErr: false,
		},
		{
			name: "invalid time range",
			options: &ListAlertGroupOptions{
				StartedAt: "invalid-time-range",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAlertGroupOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListAlertGroupQueryURL(t *testing.T) {
	tests := []struct {
		name        string
		options     *ListAlertGroupOptions
		expectedURL string
	}{
		{
			name: "single label",
			options: &ListAlertGroupOptions{
				Labels: []string{"env:prod"},
			},
			expectedURL: "/api/v1/alert_groups/?label=env%3Aprod",
		},
		{
			name: "multiple labels",
			options: &ListAlertGroupOptions{
				Labels: []string{"env:prod", "severity:high", "team:backend"},
			},
			expectedURL: "/api/v1/alert_groups/?label=env%3Aprod&label=severity%3Ahigh&label=team%3Abackend",
		},
		{
			name: "empty labels",
			options: &ListAlertGroupOptions{
				Labels: []string{},
			},
			expectedURL: "/api/v1/alert_groups/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux, server, client := setup(t)
			defer teardown(server)

			var capturedURL string
			mux.HandleFunc("/api/v1/alert_groups/", func(w http.ResponseWriter, r *http.Request) {
				capturedURL = r.URL.String()
				t.Logf("Request URL: %s", capturedURL)
				w.WriteHeader(http.StatusOK)
				// Add minimal response body
				fmt.Fprint(w, `{"count": 0, "next": null, "previous": null, "results": []}`)
			})

			_, _, err := client.AlertGroups.ListAlertGroups(tt.options)
			if err != nil {
				t.Fatal(err)
			}

			if capturedURL != tt.expectedURL {
				t.Errorf("Request URL = %v, want %v", capturedURL, tt.expectedURL)
			}
		})
	}
}

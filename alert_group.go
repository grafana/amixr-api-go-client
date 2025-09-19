package aapi

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// AlertGroupService handles requests to the on-call alert_groups endpoint.
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/alertgroups/
type AlertGroupService struct {
	client *Client
	url    string
}

// NewAlertGroupService creates an AlertGroupService with the defined URL.
func NewAlertGroupService(client *Client) *AlertGroupService {
	alertGroupService := AlertGroupService{}
	alertGroupService.client = client
	alertGroupService.url = "alert_groups"
	return &alertGroupService
}

// PaginatedAlertGroupsResponse represents a paginated response from the on-call alerts API.
type PaginatedAlertGroupsResponse struct {
	PaginatedResponse
	AlertGroups []*AlertGroup `json:"results"`
}

// AlertGroup represents an on-call alert group.
type AlertGroup struct {
	ID             string            `json:"id"`
	IntegrationID  string            `json:"integration_id"`
	RouteID        string            `json:"route_id"`
	AlertsCount    int               `json:"alerts_count"`
	State          string            `json:"state"`
	CreatedAt      string            `json:"created_at"`
	ResolvedAt     string            `json:"resolved_at"`
	AcknowledgedAt string            `json:"acknowledged_at"`
	Title          string            `json:"title"`
	Permalinks     map[string]string `json:"permalinks"`
}

// validateTimeRange validates if the time range string matches the expected format
// Expected format: %Y-%m-%dT%H:%M:%S_%Y-%m-%dT%H:%M:%S
func validateTimeRange(timeRange string) error {
	if timeRange == "" {
		return nil
	}

	// Check if the string matches the expected format
	pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}_\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$`
	matched, err := regexp.MatchString(pattern, timeRange)
	if err != nil {
		return fmt.Errorf("error validating time range format: %v", err)
	}
	if !matched {
		return fmt.Errorf("invalid time range format. Expected format: YYYY-MM-DDThh:mm:ss_YYYY-MM-DDThh:mm:ss")
	}

	// Split the time range into start and end times
	times := regexp.MustCompile(`_`).Split(timeRange, 2)
	if len(times) != 2 {
		return fmt.Errorf("invalid time range format: missing separator '_'")
	}

	// Parse both times to ensure they are valid
	startTime, err := time.Parse("2006-01-02T15:04:05", times[0])
	if err != nil {
		return fmt.Errorf("invalid start time format: %v", err)
	}

	endTime, err := time.Parse("2006-01-02T15:04:05", times[1])
	if err != nil {
		return fmt.Errorf("invalid end time format: %v", err)
	}

	// Validate that end time is after start time
	if endTime.Before(startTime) {
		return fmt.Errorf("end time must be after start time")
	}

	return nil
}

// ListAlertGroupOptions represent filter options supported by the on-call alert_groups API.
type ListAlertGroupOptions struct {
	ListOptions
	AlertGroupID  string `url:"id,omitempty" json:"alert_group_id,omitempty"`
	RouteID       string `url:"route_id,omitempty" json:"route_id,omitempty"`
	IntegrationID string `url:"integration_id,omitempty" json:"integration_id,omitempty" `
	State         string `url:"state,omitempty" json:"state,omitempty" `
	TeamID        string `url:"team_id,omitempty" json:"team_id,omitempty"`
	// StartedAt is a time range in ISO 8601 format with start and end timestamps separated by underscore.
	// Expected format: %Y-%m-%dT%H:%M:%S_%Y-%m-%dT%H:%M:%S
	// Example: "2024-03-20T10:00:00_2024-03-21T10:00:00"
	StartedAt string `url:"started_at,omitempty" json:"started_at,omitempty"`
	// Labels are matching labels that can be passed multiple times.
	// Expected format: key1:value1
	// Example: ["env:prod", "severity:high"]
	Labels []string `url:"label,omitempty" json:"label,omitempty"`
	Name   string   `url:"name,omitempty" json:"name,omitempty"`
}

// Validate checks if the options are valid
func (o *ListAlertGroupOptions) Validate() error {
	if err := validateTimeRange(o.StartedAt); err != nil {
		return err
	}
	return nil
}

// ListAlertGroups fetches all on-call alerts for authorized organization.
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/alertgroups/
func (service *AlertGroupService) ListAlertGroups(opt *ListAlertGroupOptions) (*PaginatedAlertGroupsResponse, *http.Response, error) {
	if opt != nil {
		if err := opt.Validate(); err != nil {
			return nil, nil, err
		}
	}

	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var alertGroups *PaginatedAlertGroupsResponse
	resp, err := service.client.Do(req, &alertGroups)
	if err != nil {
		return nil, resp, err
	}

	return alertGroups, resp, err
}

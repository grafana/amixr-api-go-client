package aapi

import (
	"fmt"
	"net/http"
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

// ListAlertGroupOptions represent filter options supported by the on-call alert_groups API.
type ListAlertGroupOptions struct {
	ListOptions
	AlertGroupID  string `url:"alert_group_id,omitempty" json:"alert_group_id,omitempty"`
	RouteID       string `url:"route_id,omitempty" json:"route_id,omitempty"`
	IntegrationID string `url:"integration_id,omitempty" json:"integration_id,omitempty" `
	State         string `url:"state,omitempty" json:"state,omitempty" `
	Name          string `url:"name,omitempty" json:"name,omitempty"`
}

// ListAlertGroups fetches all on-call alerts for authorized organization.
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/alertgroups/
func (service *AlertGroupService) ListAlertGroups(opt *ListAlertGroupOptions) (*PaginatedAlertGroupsResponse, *http.Response, error) {
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

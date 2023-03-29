package aapi

import (
	"fmt"
	"net/http"
)

// AlertService handles requests to the on-call alerts endpoint.
//
// // https://grafana.com/docs/oncall/latest/oncall-api-reference/alerts/
type AlertService struct {
	client *Client
	url    string
}

// NewAlertService creates an AlertService with the defined URL.
func NewAlertService(client *Client) *AlertService {
	alertService := AlertService{}
	alertService.client = client
	alertService.url = "alerts"
	return &alertService
}

// PaginatedAlertsResponse represents a paginated response from the on-call alerts API.
type PaginatedAlertsResponse struct {
	PaginatedResponse
	Alerts []*Alert `json:"results"`
}

// Alert represents an on-call alert.
type Alert struct {
	ID           string       `json:"id"`
	AlertGroupID string       `json:"alert_group_id"`
	CreatedAt    string       `json:"created_at"`
	Payload      AlertPayload `json:"payload"`
}

// AlertPayload represents an on-call alert payload.
type AlertPayload struct {
	State       string           `json:"state"`
	Title       string           `json:"title"`
	RuleID      int              `json:"ruleId"`
	Message     string           `json:"message"`
	RuleURL     string           `json:"ruleUrl"`
	RuleName    string           `json:"ruleName"`
	EvalMatches []AlertEvalMatch `json:"evalMatches"`
}

// AlertEvalMatch represents an on-call alert payload evalMatch.
type AlertEvalMatch struct {
	Tags   []string `json:"tags"`
	Value  int64    `json:"value"`
	Metric string   `json:"metric"`
}

// ListAlertOptions represent filter options supported by the on-call alerts API.
type ListAlertOptions struct {
	ListOptions
	AlertGroupID string `url:"alert_group_id,omitempty" json:"alert_group_id,omitempty"`
	Name         string `url:"search,omitempty" json:"search,omitempty"`
}

// ListAlerts fetches all on-call alerts for authorized organization.
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/alerts/
func (service *AlertService) ListAlerts(opt *ListAlertOptions) (*PaginatedAlertsResponse, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var alerts *PaginatedAlertsResponse
	resp, err := service.client.Do(req, &alerts)
	if err != nil {
		return nil, resp, err
	}

	return alerts, resp, err
}

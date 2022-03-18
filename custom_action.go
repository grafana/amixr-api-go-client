package aapi

import (
	"fmt"
	"net/http"
)

// CustomActionService handles requests to outgoing webhook endpoint
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
type CustomActionService struct {
	client *Client
	url    string
}

// NewCustomActionService creates CustomActionService with defined url
func NewCustomActionService(client *Client) *CustomActionService {
	customActionService := CustomActionService{}
	customActionService.client = client
	customActionService.url = "actions"
	return &customActionService
}

type PaginatedCustomActionsResponse struct {
	PaginatedResponse
	CustomActions []*CustomAction `json:"results"`
}

type CustomAction struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	TeamId string `json:"team_id"`
}

type ListCustomActionOptions struct {
	ListOptions
	Name string `url:"name,omitempty" json:"name,omitempty"`
}

// ListCustomActions fetches all customActions for authorized organization
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/#list-actions
func (service *CustomActionService) ListCustomActions(opt *ListCustomActionOptions) (*PaginatedCustomActionsResponse, *http.Response, error) {
	u := fmt.Sprintf("%s", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var customActions *PaginatedCustomActionsResponse
	resp, err := service.client.Do(req, &customActions)
	if err != nil {
		return nil, resp, err
	}

	return customActions, resp, err
}

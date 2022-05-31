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
	ID                  string `json:"id"`
	Name                string `json:"name"`
	TeamId              string `json:"team_id"`
	Webhook             string `json:"webhook"`
	Data                string `json:"data"`
	User                string `json:"user"`
	Password            string `json:"password"`
	AuthorizationHeader string `json:"authorization_header"`
	ForwardWholePayload bool   `json:"forward_whole_payload"`
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

type GetCustomActionOptions struct {
}

// GetCustomAction fetches custom action by given id.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
func (service *CustomActionService) GetCustomAction(id string, opt *GetCustomActionOptions) (*CustomAction, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	customAction := new(CustomAction)
	resp, err := service.client.Do(req, customAction)
	if err != nil {
		return nil, resp, err
	}

	return customAction, resp, err
}

type CreateCustomActionOptions struct {
	Name                string  `json:"name,omitempty"`
	TeamId              string  `json:"team_id"`
	Webhook             string  `json:"webhook,omitempty"`
	Data                *string `json:"data"`
	User                *string `json:"user"`
	Password            *string `json:"password"`
	AuthorizationHeader *string `json:"authorization_header"`
	ForwardWholePayload bool    `json:"forward_whole_payload"`
}

// CreateCustomAction creates custom action
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
func (service *CustomActionService) CreateCustomAction(opt *CreateCustomActionOptions) (*CustomAction, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	customAction := new(CustomAction)

	resp, err := service.client.Do(req, customAction)

	if err != nil {
		return nil, resp, err
	}

	return customAction, resp, err
}

type UpdateCustomActionOptions struct {
	Name                string `json:"name,omitempty"`
	Webhook             string `json:"webhook"`
	Data                *string `json:"data"`
	User                *string `json:"user"`
	Password            *string `json:"password"`
	AuthorizationHeader *string `json:"authorization_header"`
	ForwardWholePayload bool   `json:"forward_whole_payload"`
}

// UpdateCustomAction updates custom action
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
func (service *CustomActionService) UpdateCustomAction(id string, opt *UpdateCustomActionOptions) (*CustomAction, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	CustomAction := new(CustomAction)
	resp, err := service.client.Do(req, CustomAction)
	if err != nil {
		return nil, resp, err
	}

	return CustomAction, resp, err
}

type DeleteCustomActionOptions struct {
}

// DeleteCustomAction deletes custom action.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
func (service *CustomActionService) DeleteCustomAction(id string, opt *DeleteCustomActionOptions) (*http.Response, error) {

	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("DELETE", u, opt)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Do(req, nil)
	return resp, err
}

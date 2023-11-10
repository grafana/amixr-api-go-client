package aapi

import (
	"fmt"
	"net/http"
)

// WebhookService handles requests to outgoing webhook endpoint
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
type WebhookService struct {
	client *Client
	url    string
}

// NewWebhookService creates WebhookService with defined url
func NewWebhookService(client *Client) *WebhookService {
	WebhookService := WebhookService{}
	WebhookService.client = client
	WebhookService.url = "webhooks"
	return &WebhookService
}

type PaginatedWebhooksResponse struct {
	PaginatedResponse
	Webhooks []*Webhook `json:"results"`
}

type Webhook struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Team                string    `json:"team"`
	Url                 string    `json:"url"`
	TriggerType         string    `json:"trigger_type"`
	HttpMethod          string    `json:"http_method"`
	Data                *string   `json:"data"`
	Username            *string   `json:"username"`
	Password            *string   `json:"password"`
	AuthorizationHeader *string   `json:"authorization_header"`
	TriggerTemplate     *string   `json:"trigger_template"`
	Headers             *string   `json:"headers"`
	ForwardAll          bool      `json:"forward_all"`
	IntegrationFilter   *[]string `json:"integration_filter"`
	IsWebhookEnabled    bool      `json:"is_webhook_enabled"`
}

type ListWebhookOptions struct {
	ListOptions
	Name string `url:"name,omitempty" json:"name,omitempty"`
}

// ListWebhooks fetches all Webhooks for authorized organization
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/#list-actions
func (service *WebhookService) ListWebhooks(opt *ListWebhookOptions) (*PaginatedWebhooksResponse, *http.Response, error) {
	u := fmt.Sprintf("%s", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var Webhooks *PaginatedWebhooksResponse
	resp, err := service.client.Do(req, &Webhooks)
	if err != nil {
		return nil, resp, err
	}

	return Webhooks, resp, err
}

type GetWebhookOptions struct {
}

// GetWebhook fetches webhook by given id.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
func (service *WebhookService) GetWebhook(id string, opt *GetWebhookOptions) (*Webhook, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	Webhook := new(Webhook)
	resp, err := service.client.Do(req, Webhook)
	if err != nil {
		return nil, resp, err
	}

	return Webhook, resp, err
}

type CreateWebhookOptions struct {
	Name                string    `json:"name"`
	Team                string    `json:"team"`
	Url                 string    `json:"url"`
	TriggerType         string    `json:"trigger_type"`
	HttpMethod          string    `json:"http_method"`
	Data                *string   `json:"data"`
	Username            *string   `json:"username"`
	Password            *string   `json:"password"`
	AuthorizationHeader *string   `json:"authorization_header"`
	TriggerTemplate     *string   `json:"trigger_template"`
	Headers             *string   `json:"headers"`
	ForwardAll          bool      `json:"forward_all"`
	IntegrationFilter   *[]string `json:"integration_filter"`
	IsWebhookEnabled    bool      `json:"is_webhook_enabled"`
}

// CreateWebhook creates webhook
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
func (service *WebhookService) CreateWebhook(opt *CreateWebhookOptions) (*Webhook, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	Webhook := new(Webhook)

	resp, err := service.client.Do(req, Webhook)

	if err != nil {
		return nil, resp, err
	}

	return Webhook, resp, err
}

type UpdateWebhookOptions struct {
	Name                string    `json:"name"`
	Team                string    `json:"team"`
	Url                 string    `json:"url"`
	TriggerType         string    `json:"trigger_type"`
	HttpMethod          string    `json:"http_method"`
	Data                *string   `json:"data"`
	Username            *string   `json:"username"`
	Password            *string   `json:"password"`
	AuthorizationHeader *string   `json:"authorization_header"`
	TriggerTemplate     *string   `json:"trigger_template"`
	Headers             *string   `json:"headers"`
	ForwardAll          bool      `json:"forward_all"`
	IntegrationFilter   *[]string `json:"integration_filter"`
	IsWebhookEnabled    bool      `json:"is_webhook_enabled"`
}

// UpdateWebhook updates webhook
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
func (service *WebhookService) UpdateWebhook(id string, opt *UpdateWebhookOptions) (*Webhook, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	Webhook := new(Webhook)
	resp, err := service.client.Do(req, Webhook)
	if err != nil {
		return nil, resp, err
	}

	return Webhook, resp, err
}

type DeleteWebhookOptions struct {
}

// DeleteWebhook deletes webhook.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/outgoing_webhooks/
func (service *WebhookService) DeleteWebhook(id string, opt *DeleteWebhookOptions) (*http.Response, error) {

	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("DELETE", u, opt)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Do(req, nil)
	return resp, err
}

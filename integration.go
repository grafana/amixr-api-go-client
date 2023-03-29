package aapi

import (
	"fmt"
	"net/http"
)

// IntegrationService handles requests to integration endpoint
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/integrations/
type IntegrationService struct {
	client *Client
	url    string
}

// NewIntegrationService creates IntegrationService with corresponding url part
func NewIntegrationService(client *Client) *IntegrationService {
	integrationService := IntegrationService{}
	integrationService.client = client
	integrationService.url = "integrations"
	return &integrationService
}

type PaginatedIntegrationsResponse struct {
	PaginatedResponse
	Integrations []*Integration `json:"results"`
}

type Integration struct {
	ID             string        `json:"id"`
	TeamId         string        `json:"team_id"`
	Name           string        `json:"name"`
	Link           string        `json:"link"`
	IncidentsCount int           `json:"incidents_count"`
	Type           string        `json:"type"`
	DefaultRoute   *DefaultRoute `json:"default_route"`
	Templates      *Templates    `json:"templates"`
}

type DefaultRoute struct {
	ID                string         `json:"id"`
	EscalationChainId *string        `json:"escalation_chain_id"`
	SlackRoute        *SlackRoute    `json:"slack,omitempty"`
	TelegramRoute     *TelegramRoute `json:"telegram,omitempty"`
	MSTeamsRoute      *MSTeamsRoute  `json:"msteams,omitempty"`
}

type Templates struct {
	GroupingKey       *string                    `json:"grouping_key"`
	ResolveSignal     *string                    `json:"resolve_signal"`
	AcknowledgeSignal *string                    `json:"acknowledge_signal"`
	SourceLink        *string                    `json:"source_link"`
	Slack             *TitleMessageImageTemplate `json:"slack"`
	Web               *TitleMessageImageTemplate `json:"web"`
	MSTeams           *TitleMessageImageTemplate `json:"msteams"`
	Telegram          *TitleMessageImageTemplate `json:"telegram"`
	PhoneCall         *TitleTemplate             `json:"phone_call"`
	SMS               *TitleTemplate             `json:"sms"`
	Email             *TitleMessageTemplate      `json:"email"`
}

type TitleMessageImageTemplate struct {
	Title    *string `json:"title"`
	Message  *string `json:"message"`
	ImageURL *string `json:"image_url"`
}

type TitleMessageTemplate struct {
	Title   *string `json:"title"`
	Message *string `json:"message"`
}

type TitleTemplate struct {
	Title *string `json:"title"`
}

type MessageTemplate struct {
	Message *string `json:"message"`
}

type ImageURLTemplate struct {
	ImageURL *string `json:"image_url"`
}

type ListIntegrationOptions struct {
	ListOptions
}

// ListIntegrations fetches all integrations for current organization.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/integrations/#get-integration
func (service *IntegrationService) ListIntegrations(opt *ListIntegrationOptions) (*PaginatedIntegrationsResponse, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var integrations *PaginatedIntegrationsResponse
	resp, err := service.client.Do(req, &integrations)
	if err != nil {
		return nil, resp, err
	}

	return integrations, resp, err
}

type GetIntegrationOptions struct {
}

// GetIntegration fetches integration by given id.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/integrations/#get-integration
func (service *IntegrationService) GetIntegration(id string, opt *GetIntegrationOptions) (*Integration, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	integration := new(Integration)
	resp, err := service.client.Do(req, integration)
	if err != nil {
		return nil, resp, err
	}

	return integration, resp, err
}

type CreateIntegrationOptions struct {
	TeamId       string        `json:"team_id"`
	Name         string        `json:"name,omitempty"`
	Type         string        `json:"type,omitempty"`
	Templates    *Templates    `json:"templates,omitempty"`
	DefaultRoute *DefaultRoute `json:"default_route,omitempty"`
}

// CreateIntegration creates integration with type, team_id and optional given name.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/integrations/#get-integration
func (service *IntegrationService) CreateIntegration(opt *CreateIntegrationOptions) (*Integration, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	integration := new(Integration)
	resp, err := service.client.Do(req, integration)
	if err != nil {
		return nil, resp, err
	}

	return integration, resp, err
}

type UpdateIntegrationOptions struct {
	Name         string        `json:"name,omitempty"`
	TeamId       string        `json:"team_id"`
	Templates    *Templates    `json:"templates,omitempty"`
	DefaultRoute *DefaultRoute `json:"default_route,omitempty"`
}

// UpdateIntegration updates integration with new templates, name and default route.
// To update template it is enough to provide at least one field.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/integrations/#update-integration
func (service *IntegrationService) UpdateIntegration(id string, opt *UpdateIntegrationOptions) (*Integration, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	integration := new(Integration)
	resp, err := service.client.Do(req, integration)
	if err != nil {
		return nil, resp, err
	}

	return integration, resp, err
}

type DeleteIntegrationOptions struct {
}

// DeleteIntegration deletes integration.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/integrations/#delete-integration
func (service *IntegrationService) DeleteIntegration(id string, opt *DeleteIntegrationOptions) (*http.Response, error) {

	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("DELETE", u, opt)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Do(req, nil)
	return resp, err
}

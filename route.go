package aapi

import (
	"fmt"
	"log"
	"net/http"
)

// RouteService handles requests to route endpoint
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/routes/
type RouteService struct {
	client *Client
	url    string
}

// NewRouteService creates RouteService with defined url
func NewRouteService(client *Client) *RouteService {
	routeService := RouteService{}
	routeService.client = client
	routeService.url = "routes"
	return &routeService
}

type PaginatedRoutesResponse struct {
	PaginatedResponse
	Routes []*Route `json:"results"`
}

type Route struct {
	ID                string         `json:"id"`
	IntegrationId     string         `json:"integration_id"`
	EscalationChainId string         `json:"escalation_chain_id"`
	Position          int            `json:"position"`
	RoutingRegex      string         `json:"routing_regex"`
	IsTheLastRoute    bool           `json:"is_the_last_route"`
	SlackRoute        *SlackRoute    `json:"slack"`
	TelegramRoute     *TelegramRoute `json:"telegram"`
	MSTeamsRoute      *MSTeamsRoute  `json:"msteams"`
}

type SlackRoute struct {
	ChannelId *string `json:"channel_id"`
	Enabled   bool    `json:"enabled"`
}
type TelegramRoute struct {
	Id      *string `json:"id"`
    Enabled bool    `json:"enabled"`
}
type MSTeamsRoute struct {
	Id      *string `json:"id"`
    Enabled bool    `json:"enabled"`
}

type ListRouteOptions struct {
	ListOptions
	IntegrationId string `url:"integration_id,omitempty" json:"integration_id,omitempty"`
	RoutingRegex  string `url:"routing_regex,omitempty" json:"routing_regex,omitempty"`
}

// ListRoutes fetches all routes for authorized organization
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/routes/#list-routes
func (service *RouteService) ListRoutes(opt *ListRouteOptions) (*PaginatedRoutesResponse, *http.Response, error) {
	u := fmt.Sprintf("%s", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var routes *PaginatedRoutesResponse
	resp, err := service.client.Do(req, &routes)
	if err != nil {
		return nil, resp, err
	}

	return routes, resp, err
}

type GetRouteOptions struct {
}

// GetRoute fetches route by given id
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/routes/#get-a-route
func (service *RouteService) GetRoute(id string, opt *GetRouteOptions) (*Route, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	route := new(Route)
	resp, err := service.client.Do(req, route)
	if err != nil {
		return nil, resp, err
	}

	return route, resp, err
}

type CreateRouteOptions struct {
	IntegrationId     string         `json:"integration_id,omitempty"`
	EscalationChainId string         `json:"escalation_chain_id,omitempty"`
	Position          *int           `json:"position,omitempty"`
	RoutingRegex      string         `json:"routing_regex,omitempty"`
	Slack             *SlackRoute    `json:"slack,omitempty"`
	Telegram          *TelegramRoute `json:"telegram,omitempty"`
	MSTeams           *MSTeamsRoute  `json:"msteams,omitempty"`
	ManualOrder       bool           `url:"manual_order,omitempty" json:"manual_order,omitempty"`
}

// CreateRoute creates route with given name and type
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/routes/#create-a-route
func (service *RouteService) CreateRoute(opt *CreateRouteOptions) (*Route, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)
	req, err := service.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	route := new(Route)

	resp, err := service.client.Do(req, route)
	log.Printf("[DEBUG] request success")

	if err != nil {
		return nil, resp, err
	}

	return route, resp, err
}

type UpdateRouteOptions struct {
	EscalationChainId string         `json:"escalation_chain_id,omitempty"`
	Position          *int           `json:"position,omitempty"`
	Slack             *SlackRoute    `json:"slack,omitempty"`
	Telegram          *TelegramRoute `json:"telegram,omitempty"`
	MSTeams           *MSTeamsRoute  `json:"msteams,omitempty"`
	RoutingRegex      string         `json:"routing_regex,omitempty"`
	ManualOrder       bool           `url:"manual_order,omitempty" json:"manual_order,omitempty"`
}

// UpdateRoute updates route with new templates and/or name. At least one field in template is required
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/routes/#update-route
func (service *RouteService) UpdateRoute(id string, opt *UpdateRouteOptions) (*Route, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	route := new(Route)
	resp, err := service.client.Do(req, route)
	if err != nil {
		return nil, resp, err
	}

	return route, resp, err
}

type DeleteRouteOptions struct {
}

// DeleteRoute deletes route
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/routes/#delete-a-route
func (service *RouteService) DeleteRoute(id string, opt *DeleteRouteOptions) (*http.Response, error) {

	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("DELETE", u, opt)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Do(req, nil)
	return resp, err
}

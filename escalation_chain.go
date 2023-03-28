package aapi

import (
	"fmt"
	"net/http"
)

// EscalationChainService handles requests to escalation chain endpoint
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_chains/
type EscalationChainService struct {
	client *Client
	url    string
}

// NewEscalationChainService creates EscalationChainService with corresponding url part
func NewEscalationChainService(client *Client) *EscalationChainService {
	escalationChainService := EscalationChainService{}
	escalationChainService.client = client
	escalationChainService.url = "escalation_chains"
	return &escalationChainService
}

type PaginatedEscalationChainsResponse struct {
	PaginatedResponse
	EscalationChains []*EscalationChain `json:"results"`
}

type EscalationChain struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	TeamId string `json:"team_id"`
}

type ListEscalationChainOptions struct {
	ListOptions
	Name string `url:"name,omitempty" json:"name,omitempty"`
}

// ListEscalationChains fetches all escalation chains for current organization.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_chains/#list-escalation-chains
func (service *EscalationChainService) ListEscalationChains(opt *ListEscalationChainOptions) (*PaginatedEscalationChainsResponse, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var escalation_chains *PaginatedEscalationChainsResponse
	resp, err := service.client.Do(req, &escalation_chains)
	if err != nil {
		return nil, resp, err
	}

	return escalation_chains, resp, err
}

type GetEscalationChainOptions struct {
}

// GetEscalationChain fetches escalation chain by given id.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_chains/#get-an-escalation-chain
func (service *EscalationChainService) GetEscalationChain(id string, opt *GetEscalationChainOptions) (*EscalationChain, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	escalation_chain := new(EscalationChain)
	resp, err := service.client.Do(req, escalation_chain)
	if err != nil {
		return nil, resp, err
	}

	return escalation_chain, resp, err
}

type CreateEscalationChainOptions struct {
	Name   string `json:"name,omitempty"`
	TeamId string `json:"team_id"`
}

// CreateEscalationChain creates escalation chain with name and team_id.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_chains/#create-an-escalation-chain
func (service *EscalationChainService) CreateEscalationChain(opt *CreateEscalationChainOptions) (*EscalationChain, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)
	req, err := service.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	escalationChain := new(EscalationChain)

	resp, err := service.client.Do(req, escalationChain)

	if err != nil {
		return nil, resp, err
	}

	return escalationChain, resp, err
}

type UpdateEscalationChainOptions struct {
	Name	string	`json:"name,omitempty"`
	TeamId	string	`json:"team_id,omitempty"`
}

// UpdateEscalationChain updates escalation chain with name.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_chains/#update-an-escalation-chain
func (service *EscalationChainService) UpdateEscalationChain(id string, opt *UpdateEscalationChainOptions) (*EscalationChain, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	escalationChain := new(EscalationChain)
	resp, err := service.client.Do(req, escalationChain)
	if err != nil {
		return nil, resp, err
	}

	return escalationChain, resp, err
}

type DeleteEscalationChainOptions struct {
}

// DeleteEscalationChain deletes escalation chain.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_chains/#delete-an-escalation-chain
func (service *EscalationChainService) DeleteEscalationChain(id string, opt *DeleteEscalationChainOptions) (*http.Response, error) {

	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("DELETE", u, opt)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Do(req, nil)
	return resp, err
}

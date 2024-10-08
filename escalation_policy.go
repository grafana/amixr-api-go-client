package aapi

import (
	"fmt"
	"log"
	"net/http"
)

// EscalationService handles requests to escalation endpoint
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_policies/
type EscalationService struct {
	client *Client
	url    string
}

// NewEscalationService creates EscalationService with defined url
func NewEscalationService(client *Client) *EscalationService {
	escalationService := EscalationService{}
	escalationService.client = client
	escalationService.url = "escalation_policies"
	return &escalationService
}

type PaginatedEscalationsResponse struct {
	PaginatedResponse
	Escalations []*Escalation `json:"results"`
}

type Escalation struct {
	ID                       string    `json:"id"`
	EscalationChainId        string    `json:"escalation_chain_id"`
	Position                 int       `json:"position"`
	Type                     *string   `json:"type"`
	Duration                 *int      `json:"duration"`
	PersonsToNotify          *[]string `json:"persons_to_notify"`
	PersonsToNotifyEachTime  *[]string `json:"persons_to_notify_next_each_time"`
	TeamToNotify             *string   `json:"team_to_notify"`
	NotifyOnCallFromSchedule *string   `json:"notify_on_call_from_schedule"`
	ActionToTrigger          *string   `json:"action_to_trigger"`
	GroupToNotify            *string   `json:"group_to_notify"`
	Important                *bool     `json:"important"`
	NotifyIfTimeFrom         *string   `json:"notify_if_time_from"`
	NotifyIfTimeTo           *string   `json:"notify_if_time_to"`
	Severity                 *string   `json:"severity"`
}

// Empty struct is here in case we want to add request params to ListEscalations.
type ListEscalationOptions struct {
	ListOptions
}

// ListEscalations gets all escalations for authorized organization
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_policies/#list-escalation-policies
func (service *EscalationService) ListEscalations(opt *ListEscalationOptions) (*PaginatedEscalationsResponse, *http.Response, error) {
	u := fmt.Sprintf("%s", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var escalations *PaginatedEscalationsResponse
	resp, err := service.client.Do(req, &escalations)
	if err != nil {
		return nil, resp, err
	}

	return escalations, resp, err
}

type GetEscalationOptions struct {
}

// GetEscalation fetches an escalation by given id
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_policies/#get-an-escalation-policy
func (service *EscalationService) GetEscalation(id string, opt *GetEscalationOptions) (*Escalation, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	escalation := new(Escalation)
	resp, err := service.client.Do(req, escalation)
	if err != nil {
		return nil, resp, err
	}

	return escalation, resp, err
}

type CreateEscalationOptions struct {
	EscalationChainId           string    `json:"escalation_chain_id,omitempty"`
	Position                    *int      `json:"position,omitempty"`
	Type                        *string   `json:"type"`
	Duration                    int       `json:"duration,omitempty"`
	PersonsToNotify             *[]string `json:"persons_to_notify,omitempty"`
	PersonsToNotifyNextEachTime *[]string `json:"persons_to_notify_next_each_time,omitempty"`
	TeamToNotify                string    `json:"team_to_notify,omitempty"`
	NotifyOnCallFromSchedule    string    `json:"notify_on_call_from_schedule,omitempty"`
	ActionToTrigger             string    `json:"action_to_trigger,omitempty"`
	GroupToNotify               string    `json:"group_to_notify,omitempty"`
	ManualOrder                 bool      `json:"manual_order,omitempty"`
	Important                   *bool     `json:"important,omitempty"`
	NotifyIfTimeFrom            string    `json:"notify_if_time_from,omitempty"`
	NotifyIfTimeTo              string    `json:"notify_if_time_to,omitempty"`
	Severity                    string    `json:"severity,omitempty"`
}

// CreateEscalation creates an  escalation
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_policies/#create-an-escalation-policy
func (service *EscalationService) CreateEscalation(opt *CreateEscalationOptions) (*Escalation, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)
	req, err := service.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	escalation := new(Escalation)

	resp, err := service.client.Do(req, escalation)
	log.Printf("[DEBUG] request success")

	if err != nil {
		return nil, resp, err
	}

	return escalation, resp, err
}

type UpdateEscalationOptions struct {
	Position                 *int      `json:"position,omitempty"`
	Type                     *string   `json:"type"`
	Duration                 int       `json:"duration,omitempty"`
	PersonsToNotify          *[]string `json:"persons_to_notify,omitempty"`
	PersonsToNotifyEachTime  *[]string `json:"persons_to_notify_next_each_time,omitempty"`
	TeamToNotify             string    `json:"team_to_notify,omitempty"`
	NotifyOnCallFromSchedule string    `json:"notify_on_call_from_schedule,omitempty"`
	ActionToTrigger          string    `json:"action_to_trigger,omitempty"`
	GroupToNotify            string    `json:"group_to_notify,omitempty"`
	ManualOrder              bool      `json:"manual_order,omitempty"`
	Important                *bool     `json:"important,omitempty"`
	NotifyIfTimeFrom         string    `json:"notify_if_time_from,omitempty"`
	NotifyIfTimeTo           string    `json:"notify_if_time_to,omitempty"`
	Severity                 string    `json:"severity,omitempty"`
}

// UpdateEscalation updates an escalation with new templates and/or name. At least one field in template is required
func (service *EscalationService) UpdateEscalation(id string, opt *UpdateEscalationOptions) (*Escalation, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	escalation := new(Escalation)
	resp, err := service.client.Do(req, escalation)
	if err != nil {
		return nil, resp, err
	}

	return escalation, resp, err
}

type DeleteEscalationOptions struct {
}

// DeleteEscalation deletes an escalation
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/escalation_policies/#list-escalation-policies
func (service *EscalationService) DeleteEscalation(id string, opt *DeleteEscalationOptions) (*http.Response, error) {

	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("DELETE", u, opt)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Do(req, nil)
	return resp, err
}

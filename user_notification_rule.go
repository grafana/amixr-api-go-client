package aapi

import (
	"fmt"
	"log"
	"net/http"
)

// UserNotificationRuleService handles requests to user notification rule endpoints
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/personal_notification_rules/
type UserNotificationRuleService struct {
	client *Client
	url    string
}

// NewUserNotificationRuleService creates UserNotificationRuleService with defined url
func NewUserNotificationRuleService(client *Client) *UserNotificationRuleService {
	userNotificationRuleService := UserNotificationRuleService{client: client, url: "personal_notification_rules"}
	return &userNotificationRuleService
}

type PaginatedUserNotificationRulesResponse struct {
	PaginatedResponse
	UserNotificationRules []*UserNotificationRule `json:"results"`
}

type UserNotificationRule struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	Position  int    `json:"position"`
	Duration  int    `json:"duration"`
	Important bool   `json:"important"`
	Type      string `json:"type"`
}

type ListUserNotificationRuleOptions struct {
	ListOptions
	UserId    string `url:"user_id,omitempty" json:"user_id,omitempty"`
	Important string `url:"important,omitempty" json:"important,omitempty"`
}

// ListUserNotificationRules fetches all user notification rules for authorized organization
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/personal_notification_rules/#list-personal-notification-rules
func (service *UserNotificationRuleService) ListUserNotificationRules(opt *ListUserNotificationRuleOptions) (*PaginatedUserNotificationRulesResponse, *http.Response, error) {
	req, err := service.client.NewRequest("GET", service.url, opt)
	if err != nil {
		return nil, nil, err
	}

	var userNotificationRules *PaginatedUserNotificationRulesResponse
	resp, err := service.client.Do(req, &userNotificationRules)
	if err != nil {
		return nil, resp, err
	}

	return userNotificationRules, resp, err
}

type GetUserNotificationRuleOptions struct {
}

// GetUserNotificationRule fetches a user notification rule by given id
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/personal_notification_rules/#get-personal-notification-rule
func (service *UserNotificationRuleService) GetUserNotificationRule(id string, opt *GetUserNotificationRuleOptions) (*UserNotificationRule, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	userNotificationRule := new(UserNotificationRule)
	resp, err := service.client.Do(req, userNotificationRule)
	if err != nil {
		return nil, resp, err
	}

	return userNotificationRule, resp, err
}

type CreateUserNotificationRuleOptions struct {
	UserId      string `json:"user_id,omitempty"`
	Position    *int   `json:"position,omitempty"`
	Duration    *int   `json:"duration,omitempty"`
	Important   bool   `json:"important,omitempty"`
	Type        string `json:"type,omitempty"`
	ManualOrder bool   `json:"manual_order,omitempty"`
}

// CreateUserNotificationRule creates a user notification rule for the given user, type, and position
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/personal_notification_rules/#post-a-personal-notification-rule
func (service *UserNotificationRuleService) CreateUserNotificationRule(opt *CreateUserNotificationRuleOptions) (*UserNotificationRule, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)
	req, err := service.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	userNotificationRule := new(UserNotificationRule)

	resp, err := service.client.Do(req, userNotificationRule)
	log.Printf("[DEBUG] request success")

	if err != nil {
		return nil, resp, err
	}

	return userNotificationRule, resp, err
}

type UpdateUserNotificationRuleOptions struct {
	Position    *int   `json:"position,omitempty"`
	Duration    *int   `json:"duration,omitempty"`
	Type        string `json:"type,omitempty"`
	ManualOrder bool   `json:"manual_order,omitempty"`
}

// UpdateUserNotificationRule updates user notification rule with new position, duration, and type
//
// NOTE: this endpoint is not currently publicly documented, but it does exist
func (service *UserNotificationRuleService) UpdateUserNotificationRule(id string, opt *UpdateUserNotificationRuleOptions) (*UserNotificationRule, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	userNotificationRule := new(UserNotificationRule)
	resp, err := service.client.Do(req, userNotificationRule)
	if err != nil {
		return nil, resp, err
	}

	return userNotificationRule, resp, err
}

type DeleteUserNotificationRuleOptions struct {
}

// DeleteUserNotificationRule deletes user notification rule
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/personal_notification_rules/#delete-a-personal-notification-rule
func (service *UserNotificationRuleService) DeleteUserNotificationRule(id string, opt *DeleteUserNotificationRuleOptions) (*http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("DELETE", u, opt)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Do(req, nil)
	return resp, err
}

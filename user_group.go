package aapi

import (
	"fmt"
	"net/http"
)

// UserGroupService handles requests for user group endpoint
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/user_groups/
type UserGroupService struct {
	client *Client
	url    string
}

// NewUserGroupService creates UserGroupService with defined url
func NewUserGroupService(client *Client) *UserGroupService {
	userGroupService := UserGroupService{}
	userGroupService.client = client
	userGroupService.url = "user_groups"
	return &userGroupService
}

type PaginatedUserGroupsResponse struct {
	PaginatedResponse
	UserGroups []*UserGroup `json:"results"`
}

type UserGroup struct {
	ID             string          `json:"id"`
	Type           string          `json:"type"`
	SlackUserGroup *SlackUserGroup `json:"slack"`
}

type SlackUserGroup struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Handle string `json:"handle"`
}

type ListUserGroupOptions struct {
	ListOptions
	SlackHandle string `url:"slack_handle,omitempty" json:"slack_handle,omitempty"`
}

// ListUserGroups gets all UserGroups for authorized organization
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/user_groups/#list-user-groups
func (service *UserGroupService) ListUserGroups(opt *ListUserGroupOptions) (*PaginatedUserGroupsResponse, *http.Response, error) {
	u := fmt.Sprintf("%s", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var userGroups *PaginatedUserGroupsResponse
	resp, err := service.client.Do(req, &userGroups)
	if err != nil {
		return nil, resp, err
	}

	return userGroups, resp, err
}

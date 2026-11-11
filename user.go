package aapi

import (
	"fmt"
	"net/http"
)

// UserService handles requests to user endpoint
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/users/
type UserService struct {
	client *Client
	url    string
}

// NewUserService creates UserService with defined url
func NewUserService(client *Client) *UserService {
	userService := UserService{}
	userService.client = client
	userService.url = "users"
	return &userService
}

type PaginatedUsersResponse struct {
	PaginatedResponse
	Users []*User `json:"results"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	// Teams holds the public ids of the teams the user belongs to. It is nil
	// when the API omits it, for example on requests that pass short=true.
	Teams []string `json:"teams"`
}

type ListUserOptions struct {
	ListOptions
	Username string `url:"username,omitempty" json:"username,omitempty"`
	// TeamID filters by membership of the team with this public id, as returned
	// by the teams endpoint. NoTeamID is not supported here: the users endpoint
	// has no filter for users who belong to no team, and "null" matches nothing.
	TeamID string `url:"team_id,omitempty" json:"team_id,omitempty"`
}

// ListUsers fetches all users for authorized organization.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/users/
func (service *UserService) ListUsers(opt *ListUserOptions) (*PaginatedUsersResponse, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var users *PaginatedUsersResponse
	resp, err := service.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, err
}

type GetUserOptions struct {
}

// GetUser fetches a user by given id.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/users/#get-a-user
func (service *UserService) GetUser(id string, opt *GetUserOptions) (*User, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := service.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

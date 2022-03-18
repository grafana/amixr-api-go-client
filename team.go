package aapi

import (
	"fmt"
	"net/http"
)

// TeamService handles requests to team endpoint
type TeamService struct {
	client *Client
	url    string
}

// NewTeamService creates TeamService with defined url
func NewTeamService(client *Client) *TeamService {
	teamService := TeamService{}
	teamService.client = client
	teamService.url = "teams"
	return &teamService
}

type PaginatedTeamsResponse struct {
	PaginatedResponse
	Teams []*Team `json:"results"`
}

type Team struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

type ListTeamOptions struct {
	ListOptions
	Name string `url:"name,omitempty" json:"name,omitempty"`
}

// ListTeams fetchs all Teams for authorized user
func (service *TeamService) ListTeams(opt *ListTeamOptions) (*PaginatedTeamsResponse, *http.Response, error) {
	u := fmt.Sprintf("%s", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var teams *PaginatedTeamsResponse
	resp, err := service.client.Do(req, &teams)
	if err != nil {
		return nil, resp, err
	}

	return teams, resp, err
}

type GetTeamOptions struct {
}

// GetTeam fetches team by given id
func (service *TeamService) GetTeam(id string, opt *GetTeamOptions) (*Team, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	team := new(Team)
	resp, err := service.client.Do(req, team)
	if err != nil {
		return nil, resp, err
	}

	return team, resp, err
}

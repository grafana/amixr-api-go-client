package aapi

import (
	"net/http"
)

// TeamHoursService handles requests to the team_hours endpoint.
type TeamHoursService struct {
	client *Client
	url    string
}

// NewTeamHoursService creates TeamHoursService with defined url
func NewTeamHoursService(client *Client) *TeamHoursService {
	return &TeamHoursService{client: client, url: "team_hours"}
}

// TimeRange is a single start/end range within a weekday, as HH:MM or HH:MM:SS
// wall-clock times in the team's timezone.
type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// WeekHours maps each weekday to its time ranges. Omitted or empty days are days off.
type WeekHours struct {
	Monday    []TimeRange `json:"monday"`
	Tuesday   []TimeRange `json:"tuesday"`
	Wednesday []TimeRange `json:"wednesday"`
	Thursday  []TimeRange `json:"thursday"`
	Friday    []TimeRange `json:"friday"`
	Saturday  []TimeRange `json:"saturday"`
	Sunday    []TimeRange `json:"sunday"`
}

// TeamHoursInterval is one named set of weekly hours. A team is in
// hours when the current time falls in any range of any interval.
type TeamHoursInterval struct {
	Name  string    `json:"name"`
	Hours WeekHours `json:"hours"`
}

// TeamHours is a team's full set of intervals.
type TeamHours struct {
	TeamID string `json:"team_id"`
	// Timezone is read-only: it is the team's timezone, or nil (UTC) if unset.
	Timezone  *string             `json:"timezone"`
	Intervals []TeamHoursInterval `json:"team_hours"`
}

type PaginatedTeamHoursResponse struct {
	PaginatedResponse
	TeamHours []*TeamHours `json:"results"`
}

type ListTeamHoursOptions struct {
	ListOptions
	TeamID string `url:"team_id,omitempty" json:"team_id,omitempty"`
}

// ListTeamHours lists the hours of teams that have any, optionally filtered by team.
func (service *TeamHoursService) ListTeamHours(opt *ListTeamHoursOptions) (*PaginatedTeamHoursResponse, *http.Response, error) {
	req, err := service.client.NewRequest("GET", service.url, opt)
	if err != nil {
		return nil, nil, err
	}

	var out *PaginatedTeamHoursResponse
	resp, err := service.client.Do(req, &out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// GetTeamHours returns a team's hours, or nil if the team has none.
func (service *TeamHoursService) GetTeamHours(teamID string) (*TeamHours, *http.Response, error) {
	out, resp, err := service.ListTeamHours(&ListTeamHoursOptions{TeamID: teamID})
	if err != nil {
		return nil, resp, err
	}
	for _, h := range out.TeamHours {
		if h.TeamID == teamID {
			return h, resp, nil
		}
	}
	return nil, resp, nil
}

// SetTeamHoursOptions replaces all of a team's intervals. Send a nil or
// empty Intervals to clear them.
type SetTeamHoursOptions struct {
	TeamID    string              `json:"team_id"`
	Intervals []TeamHoursInterval `json:"team_hours"`
}

// SetTeamHours replaces a team's intervals (upsert).
func (service *TeamHoursService) SetTeamHours(opt *SetTeamHoursOptions) (*TeamHours, *http.Response, error) {
	req, err := service.client.NewRequest("POST", service.url+"/", opt)
	if err != nil {
		return nil, nil, err
	}

	out := new(TeamHours)
	resp, err := service.client.Do(req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

package aapi

import (
	"fmt"
	"net/http"
)

// ScheduleService handles requests to schedule endpoint
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/schedules/
type ScheduleService struct {
	client *Client
	url    string
}

// NewScheduleService creates ScheduleService with defined url
func NewScheduleService(client *Client) *ScheduleService {
	scheduleService := ScheduleService{}
	scheduleService.client = client
	scheduleService.url = "schedules"
	return &scheduleService
}

type PaginatedSchedulesResponse struct {
	PaginatedResponse
	Schedules []*Schedule `json:"results"`
}

type Schedule struct {
	ID                 string         `json:"id"`
	TeamId             string         `json:"team_id"`
	Type               string         `json:"type"`
	OnCallNow          []string       `json:"on_call_now"`
	Name               string         `json:"name"`
	ICalUrlPrimary     *string        `json:"ical_url_primary"`
	ICalUrlOverrides   *string        `json:"ical_url_overrides"`
	EnableWebOverrides bool           `json:"enable_web_overrides"`
	TimeZone           string         `json:"time_zone"`
	Slack              *SlackSchedule `json:"slack"`
	Shifts             *[]string      `json:"shifts"`
}

type SlackSchedule struct {
	ChannelId   *string `json:"channel_id"`
	UserGroupId *string `json:"user_group_id"`
}

type ListScheduleOptions struct {
	ListOptions
	Name   string `url:"name,omitempty" json:"name,omitempty"`
	TeamID string `url:"team_id,omitempty" json:"team_id,omitempty"`
}

// ListSchedules fetches all schedules for authorized organization
// Optional filter:
// - team_id: Exact match filter for team ID
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/schedules/#list-schedules
func (service *ScheduleService) ListSchedules(opt *ListScheduleOptions) (*PaginatedSchedulesResponse, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var schedules *PaginatedSchedulesResponse
	resp, err := service.client.Do(req, &schedules)
	if err != nil {
		return nil, resp, err
	}

	return schedules, resp, err
}

type GetScheduleOptions struct {
}

// GetSchedule fetches a schedule by given id
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/schedules/#get-a-schedule
func (service *ScheduleService) GetSchedule(id string, opt *GetScheduleOptions) (*Schedule, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	schedule := new(Schedule)
	resp, err := service.client.Do(req, schedule)
	if err != nil {
		return nil, resp, err
	}

	return schedule, resp, err
}

type CreateScheduleOptions struct {
	TeamId             string         `json:"team_id"`
	Name               string         `json:"name"`
	Type               string         `json:"type"`
	ICalUrlPrimary     *string        `json:"ical_url_primary"`
	ICalUrlOverrides   *string        `json:"ical_url_overrides"`
	EnableWebOverrides bool           `json:"enable_web_overrides"`
	TimeZone           string         `json:"time_zone,omitempty"`
	Slack              *SlackSchedule `json:"slack,omitempty"`
	Shifts             *[]string      `json:"shifts"`
}

// CreateSchedule creates a schedule with given name, type and other parameters depending on type/
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/schedules/#create-a-schedule
func (service *ScheduleService) CreateSchedule(opt *CreateScheduleOptions) (*Schedule, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)
	req, err := service.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	schedule := new(Schedule)

	resp, err := service.client.Do(req, schedule)

	if err != nil {
		return nil, resp, err
	}

	return schedule, resp, err
}

type UpdateScheduleOptions struct {
	Name               string         `json:"name,omitempty"`
	TeamId             string         `json:"team_id"`
	ICalUrlPrimary     *string        `json:"ical_url_primary"`
	ICalUrlOverrides   *string        `json:"ical_url_overrides"`
	TimeZone           string         `json:"time_zone,omitempty"`
	EnableWebOverrides bool           `json:"enable_web_overrides"`
	Slack              *SlackSchedule `json:"slack,omitempty"`
	Shifts             *[]string      `json:"shifts"`
}

// UpdateSchedule updates a schedule.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/schedules/#update-a-schedule
func (service *ScheduleService) UpdateSchedule(id string, opt *UpdateScheduleOptions) (*Schedule, *http.Response, error) {
	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	schedule := new(Schedule)
	resp, err := service.client.Do(req, schedule)
	if err != nil {
		return nil, resp, err
	}

	return schedule, resp, err
}

type DeleteScheduleOptions struct {
}

// DeleteSchedule deletes a schedule.
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/schedules/#delete-a-schedule
func (service *ScheduleService) DeleteSchedule(id string, opt *DeleteScheduleOptions) (*http.Response, error) {

	u := fmt.Sprintf("%s/%s/", service.url, id)

	req, err := service.client.NewRequest("DELETE", u, opt)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Do(req, nil)
	return resp, err
}

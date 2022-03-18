package aapi

import (
	"fmt"
	"net/http"
)

// SlackChannelService handles requests to slack channel endpoint
//
// // https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/slack_channels/
type SlackChannelService struct {
	client *Client
	url    string
}

// NewSlackChannelsService creates SlackChannelService with defined url
func NewSlackChannelService(client *Client) *SlackChannelService {
	slackChannelService := SlackChannelService{}
	slackChannelService.client = client
	slackChannelService.url = "slack_channels"
	return &slackChannelService
}

type PaginatedSlackChannelsResponse struct {
	PaginatedResponse
	SlackChannels []*SlackChannel `json:"results"`
}

type SlackChannel struct {
	Name    string `json:"name"`
	SlackId string `json:"slack_id"`
}

type ListSlackChannelOptions struct {
	ListOptions
	ChannelName string `url:"channel_name,omitempty" json:"channel_name,omitempty"`
}

// ListSlackChannels gets all slackChannels for authorized organization
//
// https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/slack_channels/#list-slack-channels
func (service *SlackChannelService) ListSlackChannels(opt *ListSlackChannelOptions) (*PaginatedSlackChannelsResponse, *http.Response, error) {
	u := fmt.Sprintf("%s", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var slackChannels *PaginatedSlackChannelsResponse
	resp, err := service.client.Do(req, &slackChannels)
	if err != nil {
		return nil, resp, err
	}

	return slackChannels, resp, err
}

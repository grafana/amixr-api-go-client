package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testSlackChannel = &SlackChannel{
	Name:    "general",
	SlackId: "TEST_SLACK_ID",
}

var testSlackChannelBody = `{
	"name": "general",
	"slack_id": "TEST_SLACK_ID"
}`

func TestListSlackChannels(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/slack_channels/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testSlackChannelBody))
	})

	options := &ListSlackChannelOptions{}

	slackChannels, _, err := client.SlackChannels.ListSlackChannels(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedSlackChannelsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		SlackChannels: []*SlackChannel{
			testSlackChannel,
		},
	}
	if !reflect.DeepEqual(want, slackChannels) {
		t.Errorf("returned\n %+v, \nwant\n %+v", slackChannels, want)
	}
}

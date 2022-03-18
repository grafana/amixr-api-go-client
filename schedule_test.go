package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testScheduleSlackChannelId = "TEST_SLACK_CHANNEL_ID"
var testScheduleSlackUserGroupId = "TEST_SLACK_USER_GROUP_ID"
var iCalUrl = "https://example.com/ical.ics"

var testSchedule = &Schedule{
	ID:             "SBM7DV7BKFUYU",
	TeamId:         "T3HRAP3K3IKOP",
	Name:           "Test Schedule",
	Type:           "ical",
	ICalUrlPrimary: &iCalUrl,
	Slack: &SlackSchedule{
		&testScheduleSlackChannelId,
		&testScheduleSlackUserGroupId,
	},
	OnCallNow: []string{"U4DNY931HHJS5", "U6RV9WPSL6DFW"},
}

var testScheduleBody = `{
	"id": "SBM7DV7BKFUYU",
	"team_id": "T3HRAP3K3IKOP",
	"type": "ical",
	"name": "Test Schedule",
	"ical_url_primary": "https://example.com/ical.ics",
	"slack": {
		"channel_id": "TEST_SLACK_CHANNEL_ID",
		"user_group_id": "TEST_SLACK_USER_GROUP_ID"
	},
	"on_call_now": [
		"U4DNY931HHJS5",
		"U6RV9WPSL6DFW"
	]
}`

func TestCreateSchedule(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/schedules/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		fmt.Fprint(w, testScheduleBody)
	})

	createOptions := &CreateScheduleOptions{
		Name: "Test Schedule",
		Slack: &SlackSchedule{
			&testScheduleSlackChannelId,
			&testScheduleSlackUserGroupId,
		},
	}
	schedule, _, err := client.Schedules.CreateSchedule(createOptions)

	if err != nil {
		t.Fatal(err)
	}

	want := testSchedule

	if !reflect.DeepEqual(want, schedule) {
		t.Errorf("returned\n %+v\n want\n %+v\n", schedule, want)
	}
}

func TestDeleteSchedule(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/schedules/SBM7DV7BKFUYU/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "DELETE")
	})

	options := &DeleteScheduleOptions{}

	_, err := client.Schedules.DeleteSchedule("SBM7DV7BKFUYU", options)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListSchedules(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/schedules/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testScheduleBody))
	})

	options := &ListScheduleOptions{}

	schedules, _, err := client.Schedules.ListSchedules(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedSchedulesResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		Schedules: []*Schedule{
			testSchedule,
		},
	}
	if !reflect.DeepEqual(want, schedules) {
		t.Errorf("returned\n %+v, \nwant\n %+v", schedules, want)
	}
}

func TestGetSchedule(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/schedules/SBM7DV7BKFUYU/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testScheduleBody)
	})

	options := &GetScheduleOptions{}

	schedule, _, err := client.Schedules.GetSchedule("SBM7DV7BKFUYU", options)

	if err != nil {
		t.Fatal(err)
	}

	want := testSchedule

	if !reflect.DeepEqual(want, schedule) {
		t.Errorf("returned\n %+v\n want\n %+v\n", schedule, want)
	}
}

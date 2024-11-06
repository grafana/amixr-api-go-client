package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var frequency = "weekly"
var weekStart = "SU"
var byDay = []string{"MO", "FR"}
var interval = 2
var users = []string{"U4DNY931HHJS5", "U6RV9WPSL6DFW"}

var testOnCallShift = &OnCallShift{
	ID:        "OH3V5FYQEYJ6M",
	TeamId:    "T3HRAP3K3IKOP",
	Name:      "Test On-Call Shift",
	Type:      "recurrent_event",
	Start:     "2020-09-04T13:00:00",
	Until:     "2020-09-05T13:00:00",
	Level:     0,
	Duration:  7200,
	Frequency: &frequency,
	Interval:  &interval,
	WeekStart: &weekStart,
	ByDay:     &byDay,
	Users:     &users,
}

var testOnCallShiftBody = `{
	"id": "OH3V5FYQEYJ6M",
	"team_id": "T3HRAP3K3IKOP",
	"name" : "Test On-Call Shift",
	"schedule_id" : "SBM7DV7BKFUYU",
	"type" : "recurrent_event",
	"start" : "2020-09-04T13:00:00",
	"until" : "2020-09-05T13:00:00",
	"level" : 0,
	"duration" : 7200,
	"frequency" : "weekly",
	"interval" : 2,
	"week_start" : "SU",
	"by_day" : ["MO", "FR"],
	"users" : [
		"U4DNY931HHJS5",
		"U6RV9WPSL6DFW"
	]
}`

func TestCreateOnCallShift(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/on_call_shifts/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		fmt.Fprint(w, testOnCallShiftBody)
	})

	tmp := 0
	createOptions := &CreateOnCallShiftOptions{
		Name:      "Test On-Call Shift",
		Type:      "recurrent_event",
		Start:     "2020-09-04T13:00:00",
		Level:     &tmp,
		Duration:  7200,
		Frequency: &frequency,
		Interval:  &interval,
		WeekStart: &weekStart,
		ByDay:     &byDay,
		Users:     &users,
	}
	shift, _, err := client.OnCallShifts.CreateOnCallShift(createOptions)

	if err != nil {
		t.Fatal(err)
	}

	want := testOnCallShift

	if !reflect.DeepEqual(want, shift) {
		t.Errorf("returned\n %+v\n want\n %+v\n", shift, want)
	}
}

func TestDeleteOnCallShift(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/on_call_shifts/OH3V5FYQEYJ6M/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "DELETE")
	})

	options := &DeleteOnCallShiftOptions{}

	_, err := client.OnCallShifts.DeleteOnCallShift("OH3V5FYQEYJ6M", options)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListOnCallShifts(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/on_call_shifts/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testOnCallShiftBody))
	})

	options := &ListOnCallShiftOptions{
		ScheduleId: "SBM7DV7BKFUYU",
	}

	shifts, _, err := client.OnCallShifts.ListOnCallShifts(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedOnCallShiftsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		OnCallShifts: []*OnCallShift{
			testOnCallShift,
		},
	}
	if !reflect.DeepEqual(want, shifts) {

		t.Errorf(" returned\n %+v, \nwant\n %+v", shifts, want)
	}
}

func TestGetOnCallShift(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/on_call_shifts/OH3V5FYQEYJ6M/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testOnCallShiftBody)
	})

	options := &GetOnCallShiftOptions{}

	shift, _, err := client.OnCallShifts.GetOnCallShift("OH3V5FYQEYJ6M", options)

	if err != nil {
		t.Fatal(err)
	}

	want := testOnCallShift

	if !reflect.DeepEqual(want, shift) {
		t.Errorf("returned\n %+v\n want\n %+v\n", shift, want)
	}
}

package aapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

const testTeamHoursBody = `{"team_id":"TI73TDU19W48J","timezone":"Europe/Dublin","team_hours":[{"name":"Business hours","hours":{"monday":[{"start":"09:00:00","end":"17:00:00"}],"tuesday":[],"wednesday":[],"thursday":[],"friday":[],"saturday":[],"sunday":[]}}]}`

func TestSetTeamHours(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/team_hours/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"team_id":"TI73TDU19W48J"`) || !strings.Contains(string(body), `"name":"Business hours"`) {
			t.Errorf("unexpected request body: %s", body)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, testTeamHoursBody)
	})

	got, _, err := client.TeamHours.SetTeamHours(&SetTeamHoursOptions{
		TeamID: "TI73TDU19W48J",
		Intervals: []TeamHoursInterval{{
			Name:  "Business hours",
			Hours: WeekHours{Monday: []TimeRange{{Start: "09:00", End: "17:00"}}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Timezone == nil || *got.Timezone != "Europe/Dublin" {
		t.Errorf("unexpected timezone: %v", got.Timezone)
	}
	if got.Intervals[0].Hours.Monday[0].End != "17:00:00" {
		t.Errorf("unexpected hours: %+v", got.Intervals)
	}
}

func TestGetTeamHours(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/team_hours/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		if r.URL.Query().Get("team_id") == "" {
			t.Errorf("missing team_id filter: %s", r.URL.RawQuery)
		}
		fmt.Fprintf(w, `{"count":1,"next":null,"previous":null,"results":[%s]}`, testTeamHoursBody)
	})

	got, _, err := client.TeamHours.GetTeamHours("TI73TDU19W48J")
	if err != nil || got == nil {
		t.Fatalf("got %v, err %v", got, err)
	}
	missing, _, err := client.TeamHours.GetTeamHours("other")
	if err != nil || missing != nil {
		t.Fatalf("expected nil for other team, got %v, err %v", missing, err)
	}
}

package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testTeam = &Team{
	ID:        "T3HRAP3K3IKOP",
	Name:      "test team",
	Email:     "test@test",
	AvatarUrl: "https://example.com/avatar",
}

var testTeamBody = `{
	"id": "T3HRAP3K3IKOP",
	"name": "test team",
	"email": "test@test",
	"avatar_url": "https://example.com/avatar"
}`

func TestListTeams(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/teams/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testTeamBody))
	})

	options := &ListTeamOptions{
		Name: "test team",
	}

	teams, _, err := client.Teams.ListTeams(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedTeamsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		Teams: []*Team{
			testTeam,
		},
	}
	if !reflect.DeepEqual(want, teams) {
		t.Errorf("returned\n %+v, \nwant\n %+v", teams, want)
	}
}

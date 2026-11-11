package aapi

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

var testUser = &User{
	ID:       "U4DNY931HHJS5",
	Email:    "public-api-demo-user-1@grafana.com",
	Role:     "admin",
	Username: "Alex",
	Teams:    []string{"TMKSD2R5W9JFA"},
}

var testUserBody = `{
	"id": "U4DNY931HHJS5",
	"email": "public-api-demo-user-1@grafana.com",
	"slack": [
		{
			"user_id": "UALEXSLACKDJPK",
			"team_id": "TALEXSLACKDJPK"
		}
	],
	"username": "Alex",
	"role": "admin",
	"teams": ["TMKSD2R5W9JFA"]
}`

func TestListUsers(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprintf(w, `{"count": 1, "next": null, "previous": null, "results": [%s]}`, testUserBody)
	})

	options := &ListUserOptions{}

	users, _, err := client.Users.ListUsers(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedUsersResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		Users: []*User{
			testUser,
		},
	}
	if !reflect.DeepEqual(want, users) {
		t.Errorf("returned\n %+v, \nwant\n %+v", users, want)
	}
}

func TestGetUser(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/users/U4DNY931HHJS5/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testUserBody)
	})

	options := &GetUserOptions{}

	user, _, err := client.Users.GetUser("U4DNY931HHJS5", options)

	if err != nil {
		t.Fatal(err)
	}

	want := testUser

	if !reflect.DeepEqual(want, user) {
		t.Errorf("returned\n %+v\n want\n %+v\n", user, want)
	}
}

func TestGetUser_NoTeams(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/users/U4DNY931HHJS5/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": "U4DNY931HHJS5", "username": "Alex"}`)
	})

	user, _, err := client.Users.GetUser("U4DNY931HHJS5", &GetUserOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if user.Teams != nil {
		t.Errorf("Teams = %+v, want nil", user.Teams)
	}
}

func TestListUsers_Filters(t *testing.T) {
	tests := []struct {
		name    string
		options *ListUserOptions
		want    url.Values
	}{
		{
			name:    "no filters",
			options: &ListUserOptions{},
			want:    url.Values{},
		},
		{
			name:    "team",
			options: &ListUserOptions{TeamID: "TMKSD2R5W9JFA"},
			want:    url.Values{"team_id": {"TMKSD2R5W9JFA"}},
		},
		{
			name:    "username and team",
			options: &ListUserOptions{Username: "Alex", TeamID: "TMKSD2R5W9JFA"},
			want:    url.Values{"username": {"Alex"}, "team_id": {"TMKSD2R5W9JFA"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux, server, client := setup(t)
			defer teardown(server)

			mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
				testRequestMethod(t, r, "GET")
				if got := r.URL.Query(); !reflect.DeepEqual(tt.want, got) {
					t.Errorf("query returned\n %+v, \nwant\n %+v", got, tt.want)
				}
				fmt.Fprint(w, `{"count": 0, "next": null, "previous": null, "results": []}`)
			})

			if _, _, err := client.Users.ListUsers(tt.options); err != nil {
				t.Fatal(err)
			}
		})
	}
}

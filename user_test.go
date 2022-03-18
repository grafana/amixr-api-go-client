package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testUser = &User{
	ID:       "U4DNY931HHJS5",
	Email:    "public-api-demo-user-1@grafana.com",
	Role:     "admin",
	Username: "Alex",
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
	"role": "admin"
}`

func TestListUsers(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testUserBody))
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

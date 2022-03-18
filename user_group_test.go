package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testUserGroup = &UserGroup{
	ID:   "GPFAPH7J7BKJB",
	Type: "slack_based",
	SlackUserGroup: &SlackUserGroup{
		"TEST_SLACK_ID",
		"Test",
		"test",
	},
}

var testUserGroupBody = `{
	"id": "GPFAPH7J7BKJB",
	"type": "slack_based",
	"slack": {
		"id": "TEST_SLACK_ID",
		"name": "Test",
		"handle": "test"
	}
}`

func TestListUserGroups(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/user_groups/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testUserGroupBody))
	})

	options := &ListUserGroupOptions{
		SlackHandle: "test",
	}

	userGroups, _, err := client.UserGroups.ListUserGroups(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedUserGroupsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		UserGroups: []*UserGroup{
			testUserGroup,
		},
	}
	if !reflect.DeepEqual(want, userGroups) {
		t.Errorf("returned\n %+v, \nwant\n %+v", userGroups, want)
	}
}

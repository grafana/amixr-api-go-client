package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var baseUrl = "personal_notification_rules"
var testId = "NT79GA9I7E4DJ"
var testUserId = "U4DNY931HHJS5"

var testUserNotificationRule = &UserNotificationRule{
	ID:       testId,
	UserId:   testUserId,
	Position: 0,
	Type:     "notify_by_sms",
}

var testUserNotificationRuleBody = fmt.Sprintf(`{
	"id": "%s",
	"user_id": "%s",
	"position": 0,
	"type": "notify_by_sms"
}`, testId, testUserId)

func TestCreateUserNotificationRule(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	url := fmt.Sprintf("/api/v1/%s/", baseUrl)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		fmt.Fprint(w, testUserNotificationRuleBody)
	})

	createOptions := &CreateUserNotificationRuleOptions{
		UserId:    testUserId,
		Important: true,
		Type:      "notify_by_sms",
	}
	userNotificationRule, _, err := client.UserNotificationRules.CreateUserNotificationRule(createOptions)

	if err != nil {
		t.Fatal(err)
	}

	want := testUserNotificationRule

	if !reflect.DeepEqual(want, userNotificationRule) {
		t.Errorf("returned\n %+v\n want\n %+v\n", userNotificationRule, want)
	}
}

func TestDeleteUserNotificationRule(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	url := fmt.Sprintf("/api/v1/%s/%s/", baseUrl, testId)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "DELETE")
	})

	options := &DeleteUserNotificationRuleOptions{}

	_, err := client.UserNotificationRules.DeleteUserNotificationRule(testId, options)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListUserNotificationRules(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	url := fmt.Sprintf("/api/v1/%s/", baseUrl)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testUserNotificationRuleBody))
	})

	options := &ListUserNotificationRuleOptions{
		UserId: testUserId,
	}

	userNotificationRules, _, err := client.UserNotificationRules.ListUserNotificationRules(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedUserNotificationRulesResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		UserNotificationRules: []*UserNotificationRule{
			testUserNotificationRule,
		},
	}
	if !reflect.DeepEqual(want, userNotificationRules) {
		t.Errorf(" returned\n %+v, \nwant\n %+v", userNotificationRules, want)
	}
}

func TestGetUserNotificationRule(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	url := fmt.Sprintf("/api/v1/%s/%s/", baseUrl, testId)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testUserNotificationRuleBody)
	})

	options := &GetUserNotificationRuleOptions{}

	userNotificationRule, _, err := client.UserNotificationRules.GetUserNotificationRule(testId, options)

	if err != nil {
		t.Fatal(err)
	}

	want := userNotificationRule

	if !reflect.DeepEqual(want, userNotificationRule) {
		t.Errorf("returned\n %+v\n want\n %+v\n", userNotificationRule, want)
	}
}

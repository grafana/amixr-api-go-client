package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testAlertGroup = &AlertGroup{
	ID:            "I68T24C13IFW1",
	IntegrationID: "CFRPV98RPR1U8",
	RouteID:       "RIYGUJXCPFHXY",
	AlertsCount:   3,
	State:         "resolved",
	CreatedAt:     "2020-05-19T12:37:01.430444Z",
	ResolvedAt:    "2020-05-19T13:37:01.429805Z",
	Title:         "Memory above 90% threshold",
	Permalinks: map[string]string{
		"slack":    "https://ghostbusters.slack.com/archives/C1H9RESGA/p135854651500008",
		"telegram": "https://t.me/c/5354/1234?thread=1234",
	},
}

var testAlertGroupBody = `{
	"id": "I68T24C13IFW1",
	"integration_id": "CFRPV98RPR1U8",
	"route_id": "RIYGUJXCPFHXY",
	"alerts_count": 3,
	"state": "resolved",
	"created_at": "2020-05-19T12:37:01.430444Z",
	"resolved_at": "2020-05-19T13:37:01.429805Z",
	"acknowledged_at": null,
	"title": "Memory above 90% threshold",
	"permalinks": {
	  "slack": "https://ghostbusters.slack.com/archives/C1H9RESGA/p135854651500008",
	  "telegram": "https://t.me/c/5354/1234?thread=1234"
	}
}`

func TestListAlertGroup(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/alert_groups/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testAlertGroupBody))
	})

	options := &ListAlertGroupOptions{
		AlertGroupID: "I68T24C13IFW1",
	}

	alerts, _, err := client.AlertGroups.ListAlertGroups(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedAlertGroupsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		AlertGroups: []*AlertGroup{
			testAlertGroup,
		},
	}

	if !reflect.DeepEqual(want, alerts) {
		fmt.Println(alerts.AlertGroups[0])
		fmt.Println(want.AlertGroups[0])
		t.Errorf(" returned\n %+v, \nwant\n %+v", alerts, want)
	}
}

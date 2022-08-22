package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var evalMatches = []AlertEvalMatch{
	{
		Value:  100,
		Metric: "High value",
	},
	{
		Value:  200,
		Metric: "Higher value",
	},
}

var payload = AlertPayload{
	State:       "alerting",
	Title:       "[Alerting] Test notification",
	RuleID:      0,
	Message:     "Someone is testing the alert notification within Grafana.",
	RuleURL:     "{{API_URL}}/",
	RuleName:    "Test notification",
	EvalMatches: evalMatches,
}

var testAlert = &Alert{
	ID:           "OH3V5FYQEYJ6M",
	AlertGroupID: "T3HRAP3K3IKOP",
	CreatedAt:    "2020-05-11T20:07:43Z",
	Payload:      payload,
}

var testAlertBody = `{
	"id": "OH3V5FYQEYJ6M",
	"alert_group_id": "T3HRAP3K3IKOP",
	"created_at": "2020-05-11T20:07:43Z",
	"payload": {
		"state": "alerting",
		"title": "[Alerting] Test notification",
		"ruleId": 0,
		"message": "Someone is testing the alert notification within Grafana.",
		"ruleUrl": "{{API_URL}}/",
		"ruleName": "Test notification",
		"evalMatches": [
			{
				"tags": null,
				"value": 100,
				"metric": "High value"
			},
			{
				"tags": null,
				"value": 200,
				"metric": "Higher value"
			}
		]
	}
}`

func TestListAlerts(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/alerts/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testAlertBody))
	})

	options := &ListAlertOptions{
		AlertGroupID: "SBM7DV7BKFUYU",
	}

	alerts, _, err := client.Alerts.ListAlerts(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedAlertsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		Alerts: []*Alert{
			testAlert,
		},
	}

	if !reflect.DeepEqual(want, alerts) {
		fmt.Println(alerts.Alerts[0])
		fmt.Println(want.Alerts[0])
		t.Errorf(" returned\n %+v, \nwant\n %+v", alerts, want)
	}
}

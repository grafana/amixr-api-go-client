package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var duration = 60
var typeWait = "wait"
var typeNotifyPersons = "notify_persons"

var testEscalation = &Escalation{
	ID:                "E3GA6SJETWWJS",
	EscalationChainId: "RIYGUJXCPFHXY",
	Position:          0,
	Type:              &typeWait,
	Duration:          &duration,
}

var testEscalationEmptyDuration = &Escalation{
	ID:                "E3GA6SJETWWJS",
	EscalationChainId: "RIYGUJXCPFHXY",
	Position:          0,
	Type:              &typeNotifyPersons,
}

var testEscalationBody = `{
	"id": "E3GA6SJETWWJS",
    "escalation_chain_id": "RIYGUJXCPFHXY",
    "position": 0,
    "type": "wait",
    "duration": 60
}`

var testEscalationEmptyDurationBody = `{
	"id": "E3GA6SJETWWJS",
    "escalation_chain_id": "RIYGUJXCPFHXY",
    "position": 0,
    "type": "notify_persons"
}`

var testUpdatedEscalationBody = `{
	"id": "E3GA6SJETWWJS",
    "escalation_chain_id": "RIYGUJXCPFHXY",
    "position": 1,
    "type": "wait",
    "duration": 60
}`

func TestCreateEscalation(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/escalation_policies/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		fmt.Fprint(w, testEscalationBody)
	})
	tmp := 0
	createOptions := &CreateEscalationOptions{
		Type:              &typeWait,
		Position:          &tmp,
		Duration:          60,
		ManualOrder:       true,
		EscalationChainId: "RIYGUJXCPFHXY",
	}
	escalation, _, err := client.Escalations.CreateEscalation(createOptions)

	if err != nil {
		t.Fatal(err)
	}

	want := testEscalation

	if !reflect.DeepEqual(want, escalation) {
		t.Errorf("returned\n %+v\n want\n %+v\n", escalation, want)
	}
}

func TestDeleteEscalation(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/escalation_policies/E3GA6SJETWWJS/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "DELETE")
	})

	options := &DeleteEscalationOptions{}

	_, err := client.Escalations.DeleteEscalation("E3GA6SJETWWJS", options)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListEscalations(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/escalation_policies/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testEscalationBody))
	})

	options := &ListEscalationOptions{}

	escalations, _, err := client.Escalations.ListEscalations(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedEscalationsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		Escalations: []*Escalation{
			testEscalation,
		},
	}
	if !reflect.DeepEqual(want, escalations) {

		t.Errorf(" returned\n %+v, \nwant\n %+v", escalations, want)
	}
}

func TestGetEscalation(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/escalation_policies/E3GA6SJETWWJS/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testEscalationBody)
	})

	options := &GetEscalationOptions{}

	escalation, _, err := client.Escalations.GetEscalation("E3GA6SJETWWJS", options)

	if err != nil {
		t.Fatal(err)
	}

	want := testEscalation

	if !reflect.DeepEqual(want, escalation) {
		t.Errorf("returned\n %+v\n want\n %+v\n", escalation, want)
	}
}

func TestGetEscalationWithEmptyDuration(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/escalation_policies/E3GA6SJETWWJS/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testEscalationEmptyDurationBody)
	})

	options := &GetEscalationOptions{}

	escalation, _, err := client.Escalations.GetEscalation("E3GA6SJETWWJS", options)

	if err != nil {
		t.Fatal(err)
	}

	want := testEscalationEmptyDuration

	if !reflect.DeepEqual(want, escalation) {
		t.Errorf("returned\n %+v\n want\n %+v\n", escalation, want)
	}
}

func TestUpdateEscalation(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/escalation_policies/E3GA6SJETWWJS/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "PUT")
		fmt.Fprint(w, testUpdatedEscalationBody)
	})
	tmp := 1
	options := &UpdateEscalationOptions{
		Position: &tmp,
	}

	escalation, _, err := client.Escalations.UpdateEscalation("E3GA6SJETWWJS", options)

	if err != nil {
		t.Fatal(err)
	}
	var duration = 60
	var testUpdatedEscalation = &Escalation{
		ID:                "E3GA6SJETWWJS",
		EscalationChainId: "RIYGUJXCPFHXY",
		Position:          1,
		Type:              &typeWait,
		Duration:          &duration,
	}

	want := testUpdatedEscalation

	if !reflect.DeepEqual(want, escalation) {
		t.Errorf("returned\n %+v\n want\n %+v\n", escalation, want)
	}
}

package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testCustomAction = &CustomAction{
	ID:     "KGEFG74LU1D8L",
	Name:   "Test action",
	TeamId: "T3HRAP3K3IKOP",
}

var testCustomActionBody = `{
	"id": "KGEFG74LU1D8L",
	"name": "Test action",
	"team_id": "T3HRAP3K3IKOP"
}`

func TestListCustomActions(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/actions/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testCustomActionBody))
	})

	options := &ListCustomActionOptions{
		Name: "Test action",
	}

	customActions, _, err := client.CustomActions.ListCustomActions(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedCustomActionsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		CustomActions: []*CustomAction{
			testCustomAction,
		},
	}
	if !reflect.DeepEqual(want, customActions) {
		t.Errorf("returned\n %+v, \nwant\n %+v", customActions, want)
	}
}

func TestCreateCustomAction(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/actions/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		fmt.Fprint(w, testCustomActionBody)
	})

	createOptions := &CreateCustomActionOptions{
		Name: "Test CustomAction",
		Url:  "https://example.com",
	}
	customAction, _, err := client.CustomActions.CreateCustomAction(createOptions)

	if err != nil {
		t.Fatal(err)
	}

	want := testCustomAction

	if !reflect.DeepEqual(want, customAction) {
		t.Errorf("returned\n %+v\n want\n %+v\n", customAction, want)
	}
}

func TestDeleteCustomAction(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/actions/KGEFG74LU1D8L/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "DELETE")
	})

	options := &DeleteCustomActionOptions{}

	_, err := client.CustomActions.DeleteCustomAction("KGEFG74LU1D8L", options)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCustomAction(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/actions/KGEFG74LU1D8L/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testCustomActionBody)
	})

	options := &GetCustomActionOptions{}

	customAction, _, err := client.CustomActions.GetCustomAction("KGEFG74LU1D8L", options)

	if err != nil {
		t.Fatal(err)
	}

	want := testCustomAction

	if !reflect.DeepEqual(want, customAction) {
		t.Errorf("returned\n %+v\n want\n %+v\n", customAction, want)
	}
}

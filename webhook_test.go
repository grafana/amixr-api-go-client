package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testWebhook = &Webhook{
	ID:          "KGEFG74LU1D8L",
	Name:        "Test action",
	Team:        "T3HRAP3K3IKOP",
	HttpMethod:  "POST",
	TriggerType: "escalation step",
	Url:         "http://test.com",
}

var testWebhookBody = `{
	"id": "KGEFG74LU1D8L",
	"name": "Test action",
	"team": "T3HRAP3K3IKOP",
	"http_method": "POST",
	"trigger_type": "escalation step",
	"url":"http://test.com"
}`

func TestListWebhooks(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/webhooks/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testWebhookBody))
	})

	options := &ListWebhookOptions{
		Name: "Test action",
	}

	Webhooks, _, err := client.Webhooks.ListWebhooks(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedWebhooksResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		Webhooks: []*Webhook{
			testWebhook,
		},
	}
	if !reflect.DeepEqual(want, Webhooks) {
		t.Errorf("returned\n %+v, \nwant\n %+v", Webhooks, want)
	}
}

func TestCreateWebhook(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/webhooks/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		fmt.Fprint(w, testWebhookBody)
	})

	createOptions := &CreateWebhookOptions{
		Name: "Test Webhook",
		Url:  "https://example.com",
	}
	Webhook, _, err := client.Webhooks.CreateWebhook(createOptions)

	if err != nil {
		t.Fatal(err)
	}

	want := testWebhook

	if !reflect.DeepEqual(want, Webhook) {
		t.Errorf("returned\n %+v\n want\n %+v\n", Webhook, want)
	}
}

func TestDeleteWebhook(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/webhooks/KGEFG74LU1D8L/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "DELETE")
	})

	options := &DeleteWebhookOptions{}

	_, err := client.Webhooks.DeleteWebhook("KGEFG74LU1D8L", options)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetWebhook(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/webhooks/KGEFG74LU1D8L/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testWebhookBody)
	})

	options := &GetWebhookOptions{}

	Webhook, _, err := client.Webhooks.GetWebhook("KGEFG74LU1D8L", options)

	if err != nil {
		t.Fatal(err)
	}

	want := testWebhook

	if !reflect.DeepEqual(want, Webhook) {
		t.Errorf("returned\n %+v\n want\n %+v\n", Webhook, want)
	}
}

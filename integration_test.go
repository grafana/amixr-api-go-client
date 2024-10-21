package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var key = "key"
var signal = "signal"
var source_link = "source_link"
var testIntegration = &Integration{
	ID:     "CFRPV98RPR1U8",
	TeamId: "T3HRAP3K3IKOP",
	Name:   "Test Grafana",
	Type:   "grafana",
	Link:   "https://grafana_url/integrations/v1/grafana/mReAoNwDm0eMwKo1mTeTwYo/",
	DefaultRoute: &DefaultRoute{
		ID: "RIYGUJXCPFHXY",
	},
	IncidentsCount: 0,
	Templates: &Templates{
		&key,
		&signal, // resolve signal;
		&signal, // ack signal
		&source_link,
		&TitleMessageImageTemplate{ // Slack
			nil,
			nil,
			nil,
		},
		&TitleMessageImageTemplate{ // Telegram
			nil,
			nil,
			nil,
		},
		&TitleMessageImageTemplate{ // MSTeams
			nil,
			nil,
			nil,
		},
		&TitleMessageImageTemplate{ // Web
			nil,
			nil,
			nil,
		},
		&TitleTemplate{ // PhoneCall
			nil,
		},
		&TitleTemplate{ // SMS
			nil,
		},
		&TitleMessageTemplate{ // Email
			nil,
			nil,
		},
		&TitleMessageTemplate{ // Mobile app
			nil,
			nil,
		},
	},
}

var testIntegrationBody = `{
	"id":"CFRPV98RPR1U8",
	"team_id":"T3HRAP3K3IKOP",
	"name":"Test Grafana",
	"type":"grafana",
	"link":"https://grafana_url/integrations/v1/grafana/mReAoNwDm0eMwKo1mTeTwYo/",
	"default_route":{
	   "id":"RIYGUJXCPFHXY"
	},
	"incidents_count":0,
	"templates":{
	   "grouping_key":"key",
	   "source_link":"source_link",
	   "resolve_signal":"signal",
	   "acknowledge_signal":"signal",
	   "slack":{
		  "title":null,
		  "message":null,
		  "image_url":null
	   },
	   "web":{
		  "title":null,
		  "message":null,
		  "image_url":null
	   },
	   "sms":{
		  "title":null
	   },
	   "phone_call":{
		  "title":null
	   },
	   "telegram":{
		  "title":null,
		  "message":null,
		  "image_url":null
	   },
	   "email":{
		  "title":null,
		  "message":null
	   },
	   "msteams":{
		  "title":null,
		  "message":null,
		  "image_url":null
	   },
	   "mobile_app":{
		  "title":null,
		  "message":null
	   }
	}
 }`

func TestCreateIntegration(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/integrations/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		fmt.Fprint(w, testIntegrationBody)
	})

	createOptions := &CreateIntegrationOptions{
		Name: "Test Grafana",
		Type: "grafana",
	}
	integration, _, err := client.Integrations.CreateIntegration(createOptions)

	if err != nil {
		t.Fatal(err)
	}

	want := testIntegration

	if !reflect.DeepEqual(want, integration) {
		t.Errorf("returned\n %+v\n want\n %+v\n", integration, want)
	}
}

func TestDeleteIntegration(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/integrations/CFRPV98RPR1U8/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "DELETE")
	})

	options := &DeleteIntegrationOptions{}

	_, err := client.Integrations.DeleteIntegration("CFRPV98RPR1U8", options)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListIntegrations(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/integrations/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testIntegrationBody))
	})

	options := &ListIntegrationOptions{}

	integrations, _, err := client.Integrations.ListIntegrations(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedIntegrationsResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		Integrations: []*Integration{
			testIntegration,
		},
	}
	if !reflect.DeepEqual(want, integrations) {

		t.Errorf(" returned\n %+v, \nwant\n %+v", integrations, want)
	}
}

func TestGetIntegration(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/integrations/CFRPV98RPR1U8/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testIntegrationBody)
	})

	options := &GetIntegrationOptions{}

	integration, _, err := client.Integrations.GetIntegration("CFRPV98RPR1U8", options)

	if err != nil {
		t.Fatal(err)
	}

	want := testIntegration

	if !reflect.DeepEqual(want, integration) {
		t.Errorf("returned\n %+v\n want\n %+v\n", integration, want)
	}
}

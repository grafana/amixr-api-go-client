package aapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var key = "key"
var signal = "signal"
var source_link = "source_link"
var testIntegration = &Integration{
	ID:           "CFRPV98RPR1U8",
	TeamId:       "T3HRAP3K3IKOP",
	Name:         "Test Grafana",
	Type:         "grafana",
	Link:         "https://grafana_url/integrations/v1/grafana/mReAoNwDm0eMwKo1mTeTwYo/",
	InboundEmail: "",
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
	Labels: []*Label{
		&Label{
			Key: KeyValueName{
				Name: "Flip",
			},
			Value: KeyValueName{
				Name: "Flop",
			},
		},
	},
	DynamicLabels: []*Label{
		&Label{
			Key: KeyValueName{
				Name: "",
			},
			Value: KeyValueName{
				Name: "Flop",
			},
		},
	},
}

var testIntegrationBody = `{
	"id":"CFRPV98RPR1U8",
	"team_id":"T3HRAP3K3IKOP",
	"name":"Test Grafana",
	"type":"grafana",
	"link":"https://grafana_url/integrations/v1/grafana/mReAoNwDm0eMwKo1mTeTwYo/",
	"inbound_email": "",
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
	},
	"labels":[
		{
			"key": {
				"name": "Flip"
			},
			"value": {
				"name": "Flop"
			}
		}
	]
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

func TestUpdateIntegrationOptions_MarshalLabels(t *testing.T) {
	label := &Label{
		Key:   KeyValueName{Name: "severity"},
		Value: KeyValueName{Name: "critical"},
	}
	labelsWithValue := []*Label{label}
	emptyLabels := []*Label{}

	tests := []struct {
		name              string
		opts              UpdateIntegrationOptions
		wantLabels        string
		wantDynamicLabels string
	}{
		{
			name:              "nil labels and dynamic labels are omitted",
			opts:              UpdateIntegrationOptions{Name: "integration", TeamId: "team"},
			wantLabels:        "absent",
			wantDynamicLabels: "absent",
		},
		{
			name: "empty labels slice is sent as an empty array",
			opts: UpdateIntegrationOptions{
				Name:   "integration",
				TeamId: "team",
				Labels: &emptyLabels,
			},
			wantLabels:        "empty",
			wantDynamicLabels: "absent",
		},
		{
			name: "empty dynamic labels slice is sent as an empty array",
			opts: UpdateIntegrationOptions{
				Name:          "integration",
				TeamId:        "team",
				DynamicLabels: &emptyLabels,
			},
			wantLabels:        "absent",
			wantDynamicLabels: "empty",
		},
		{
			name: "both empty label lists are sent as empty arrays",
			opts: UpdateIntegrationOptions{
				Name:          "integration",
				TeamId:        "team",
				Labels:        EmptyLabelList(),
				DynamicLabels: EmptyLabelList(),
			},
			wantLabels:        "empty",
			wantDynamicLabels: "empty",
		},
		{
			name: "populated labels are sent",
			opts: UpdateIntegrationOptions{
				Name:   "integration",
				TeamId: "team",
				Labels: &labelsWithValue,
			},
			wantLabels:        "present",
			wantDynamicLabels: "absent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.opts)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			if got := jsonFieldState(t, body, "labels"); got != tt.wantLabels {
				t.Errorf("labels in request body = %q, want %q; body = %s", got, tt.wantLabels, body)
			}
			if got := jsonFieldState(t, body, "dynamic_labels"); got != tt.wantDynamicLabels {
				t.Errorf("dynamic_labels in request body = %q, want %q; body = %s", got, tt.wantDynamicLabels, body)
			}
		})
	}
}

func TestUpdateIntegration_LabelsRequestBody(t *testing.T) {
	tests := []struct {
		name              string
		opts              *UpdateIntegrationOptions
		wantLabels        string
		wantDynamicLabels string
	}{
		{
			name: "nil labels are omitted from update request body",
			opts: &UpdateIntegrationOptions{
				Name:   "integration",
				TeamId: "team",
			},
			wantLabels:        "absent",
			wantDynamicLabels: "absent",
		},
		{
			name: "empty labels slice is sent in update request body",
			opts: &UpdateIntegrationOptions{
				Name:   "integration",
				TeamId: "team",
				Labels: EmptyLabelList(),
			},
			wantLabels:        "empty",
			wantDynamicLabels: "absent",
		},
		{
			name: "empty dynamic labels slice is sent in update request body",
			opts: &UpdateIntegrationOptions{
				Name:          "integration",
				TeamId:        "team",
				DynamicLabels: EmptyLabelList(),
			},
			wantLabels:        "absent",
			wantDynamicLabels: "empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux, server, client := setup(t)
			defer teardown(server)

			mux.HandleFunc("/api/v1/integrations/CFRPV98RPR1U8/", func(w http.ResponseWriter, r *http.Request) {
				testRequestMethod(t, r, "PUT")

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}

				if got := jsonFieldState(t, body, "labels"); got != tt.wantLabels {
					t.Errorf("labels in request body = %q, want %q; body = %s", got, tt.wantLabels, body)
				}
				if got := jsonFieldState(t, body, "dynamic_labels"); got != tt.wantDynamicLabels {
					t.Errorf("dynamic_labels in request body = %q, want %q; body = %s", got, tt.wantDynamicLabels, body)
				}

				fmt.Fprint(w, testIntegrationBody)
			})

			_, _, err := client.Integrations.UpdateIntegration("CFRPV98RPR1U8", tt.opts)
			if err != nil {
				t.Fatalf("UpdateIntegration() error = %v", err)
			}
		})
	}
}

func jsonFieldState(t *testing.T, body []byte, field string) string {
	t.Helper()

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; body = %s", err, body)
	}

	raw, ok := payload[field]
	if !ok {
		return "absent"
	}

	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "[]" {
		return "empty"
	}

	return "present"
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

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

var testSlackChannelId = "TEST_SLACK_CHANNEL_ID"
var testTelegramChannelId = "TESTTELEGRAMID"
var testMSTeamsChannelId = "TESTMSTEAMSID"

var testRoute = &Route{
	ID:             "RH2V5FYIPYJ1M",
	IntegrationId:  "CGEXJ922S7TXQ",
	Position:       0,
	RoutingRegex:   "us-west",
	RoutingType:    "regex",
	IsTheLastRoute: false,
	SlackRoute: &SlackRoute{
		ChannelId: &testSlackChannelId,
		Enabled:   true,
	},
	TelegramRoute: &TelegramRoute{
		Id:      &testTelegramChannelId,
		Enabled: true,
	},
	MSTeamsRoute: &MSTeamsRoute{
		Id:      &testMSTeamsChannelId,
		Enabled: true,
	},
}

var testRouteBody = `{
	"id": "RH2V5FYIPYJ1M",
	"integration_id": "CGEXJ922S7TXQ",
	"routing_regex": "us-west",
	"routing_type": "regex",
	"position": 0,
	"is_the_last_route": false,
	"slack": {
		"channel_id": "TEST_SLACK_CHANNEL_ID",
		"enabled": true
	},
	"telegram": {
	    "id": "TESTTELEGRAMID",
        "enabled": true
	},
	"msteams": {
	    "id": "TESTMSTEAMSID",
	    "enabled": true
	}
}`

func TestCreateRoute(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/routes/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		fmt.Fprint(w, testRouteBody)
	})

	createOptions := &CreateRouteOptions{
		IntegrationId: "CGEXJ922S7TXQ",
		RoutingRegex:  "us-west",
		RoutingType:   "regex",
		Slack: &SlackRoute{
			ChannelId: &testSlackChannelId,
			Enabled:   true,
		},
		Telegram: &TelegramRoute{
			Id:      &testTelegramChannelId,
			Enabled: true,
		},
		MSTeams: &MSTeamsRoute{
			Id:      &testMSTeamsChannelId,
			Enabled: true,
		},
	}
	route, _, err := client.Routes.CreateRoute(createOptions)

	if err != nil {
		t.Fatal(err)
	}

	want := testRoute

	if !reflect.DeepEqual(want, route) {
		t.Errorf("returned\n %+v\n want\n %+v\n", route, want)
	}
}

func TestCreateRouteOptionsJSONEscalationChainId(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		opt := &CreateRouteOptions{
			IntegrationId: "CGEXJ922S7TXQ",
			RoutingRegex:  "us-west",
		}

		data, err := json.Marshal(opt)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(data), `"escalation_chain_id":null`) {
			t.Fatalf("expected explicit null escalation_chain_id, got %s", data)
		}
	})

	t.Run("id", func(t *testing.T) {
		chainID := "RIYGUJXCPFHXY"
		opt := &CreateRouteOptions{
			IntegrationId:     "CGEXJ922S7TXQ",
			RoutingRegex:      "us-west",
			EscalationChainId: &chainID,
		}

		data, err := json.Marshal(opt)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(data), `"escalation_chain_id":"RIYGUJXCPFHXY"`) {
			t.Fatalf("expected escalation_chain_id in payload, got %s", data)
		}
	})
}

func TestCreateRouteWithoutEscalationChain(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/routes/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), `"escalation_chain_id":null`) {
			t.Fatalf("expected explicit null escalation_chain_id, got %s", body)
		}

		fmt.Fprint(w, `{
			"id": "RH2V5FYIPYJ1M",
			"integration_id": "CGEXJ922S7TXQ",
			"escalation_chain_id": null,
			"routing_regex": "us-west",
			"routing_type": "regex",
			"position": 0,
			"is_the_last_route": false,
			"slack": {"channel_id": null, "enabled": true},
			"telegram": {"id": null, "enabled": false},
			"msteams": {"id": null, "enabled": false}
		}`)
	})

	createOptions := &CreateRouteOptions{
		IntegrationId: "CGEXJ922S7TXQ",
		RoutingRegex:  "us-west",
	}
	route, _, err := client.Routes.CreateRoute(createOptions)
	if err != nil {
		t.Fatal(err)
	}

	if route.EscalationChainId != nil {
		t.Fatalf("expected nil escalation chain id, got %v", route.EscalationChainId)
	}
}

func TestGetRouteWithoutEscalationChain(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/routes/RH2V5FYIPYJ1M/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": "RH2V5FYIPYJ1M",
			"integration_id": "CGEXJ922S7TXQ",
			"escalation_chain_id": null,
			"routing_regex": "us-west",
			"routing_type": "regex",
			"position": 0,
			"is_the_last_route": false,
			"slack": {"channel_id": null, "enabled": true},
			"telegram": {"id": null, "enabled": false},
			"msteams": {"id": null, "enabled": false}
		}`)
	})

	route, _, err := client.Routes.GetRoute("RH2V5FYIPYJ1M", &GetRouteOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if route.EscalationChainId != nil {
		t.Fatalf("expected nil escalation chain id, got %v", route.EscalationChainId)
	}
}

func TestDeleteRoute(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/routes/RH2V5FYIPYJ1M/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "DELETE")
	})

	options := &DeleteRouteOptions{}

	_, err := client.Routes.DeleteRoute("RH2V5FYIPYJ1M", options)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListRoutes(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/routes/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testRouteBody))
	})

	options := &ListRouteOptions{
		IntegrationId: "CGEXJ922S7TXQ",
	}

	routes, _, err := client.Routes.ListRoutes(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedRoutesResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		Routes: []*Route{
			testRoute,
		},
	}
	if !reflect.DeepEqual(want, routes) {

		t.Errorf(" returned\n %+v, \nwant\n %+v", routes, want)
	}
}

func TestGetRoute(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/routes/RH2V5FYIPYJ1M/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, testRouteBody)
	})

	options := &GetRouteOptions{}

	route, _, err := client.Routes.GetRoute("RH2V5FYIPYJ1M", options)

	if err != nil {
		t.Fatal(err)
	}

	want := testRoute

	if !reflect.DeepEqual(want, route) {
		t.Errorf("returned\n %+v\n want\n %+v\n", route, want)
	}
}

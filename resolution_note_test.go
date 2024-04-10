package aapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testResolutionNote = &ResolutionNote{
	ID:          "M4BTQUS3PRHYQ",
	AlertGroupID: "I68T24C13IFW1",
	Author:      "U4DNY931HHJS5",
	Source:      "web",
	CreatedAt:   "2020-06-19T12:40:01.429805Z",
	Text:        "Demo resolution note",
}

var testResolutionNoteBody = `{
	"id": "M4BTQUS3PRHYQ",
	"alert_group_id": "I68T24C13IFW1",
	"author": "U4DNY931HHJS5",
	"source": "web",
	"created_at": "2020-06-19T12:40:01.429805Z",
	"text": "Demo resolution note"
}`



func TestListResolutionNote(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v1/resolution_notes/", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"count": 1, "next": null, "previous": null, "results": [%s]}`, testResolutionNoteBody))
	})

	options := &ListResolutionNoteOptions{
		AlertGroupID: "I68T24C13IFW1",
	}

	alerts, _, err := client.ResolutionNotes.ListResolutionNotes(options)
	if err != nil {
		t.Fatal(err)
	}

	want := &PaginatedResolutionNotesResponse{
		PaginatedResponse: PaginatedResponse{
			Count:    1,
			Next:     nil,
			Previous: nil,
		},
		ResolutionNotes: []*ResolutionNote{
			testResolutionNote,
		},
	}

	if !reflect.DeepEqual(want, alerts) {
		fmt.Println(alerts.ResolutionNotes[0])
		fmt.Println(want.ResolutionNotes[0])
		t.Errorf(" returned\n %+v, \nwant\n %+v", alerts, want)
	}
}
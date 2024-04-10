package aapi

import (
	"fmt"
	"net/http"
)

// ResolutionNoteService handles requests to the on-call resolution_notes endpoint.
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/resolution_notes/
type ResolutionNoteService struct {
	client *Client
	url    string
}

// NewResolutionNoteService creates an ResolutionNoteService with the defined URL.
func NewResolutionNoteService(client *Client) *ResolutionNoteService {
	resolutionNoteService := ResolutionNoteService{}
	resolutionNoteService.client = client
	resolutionNoteService.url = "resolution_notes"
	return &resolutionNoteService
}

// PaginatedResolutionNotesResponse represents a paginated response from the on-call resolution note API.
type PaginatedResolutionNotesResponse struct {
	PaginatedResponse
	ResolutionNotes []*ResolutionNote `json:"results"`
}

// ResolutionNote represents an on-call resolution note.
type ResolutionNote struct {
	ID           string `json:"id"`
	AlertGroupID string `json:"alert_group_id"`
	Author       string `json:"author"`
	Source       string `json:"source"`
	CreatedAt    string `json:"created_at"`
	Text         string `json:"text"`
}

// ListResolutionNoteOptions represent filter options supported by the on-call resolution note API.
type ListResolutionNoteOptions struct {
	ListOptions
	AlertGroupID string `url:"alert_group_id,omitempty" json:"alert_group_id,omitempty"`
}

// ListResolutionNotes fetches all on-call resolution notes associated to an alert group for authorized organization.
//
// https://grafana.com/docs/oncall/latest/oncall-api-reference/alertgroups/
func (service *ResolutionNoteService) ListResolutionNotes(opt *ListResolutionNoteOptions) (*PaginatedResolutionNotesResponse, *http.Response, error) {
	u := fmt.Sprintf("%s/", service.url)

	req, err := service.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var resolutionNotes *PaginatedResolutionNotesResponse
	resp, err := service.client.Do(req, &resolutionNotes)
	if err != nil {
		return nil, resp, err
	}

	return resolutionNotes, resp, err
}

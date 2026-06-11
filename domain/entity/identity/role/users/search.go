package users

// This file provides a JSON search handler for the user auto-complete component.
// GET /action/roles/detail/{id}/users/search?q=term
// Returns JSON: [{"value":"workspace_user_id","label":"First Last (email)"}, ...]

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
)

// searchOption is the JSON shape returned by the search handler.
type searchOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// SearchDeps holds dependencies for the user search handler.
type SearchDeps struct {
	ListWorkspaceUsers func(ctx context.Context, req *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error)
}

// NewSearchUsersAction creates an http.HandlerFunc that returns workspace users as
// JSON for the auto-complete component.
// GET /action/roles/detail/{id}/users/search?q=term
func NewSearchUsersAction(deps *SearchDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))

		// Require at least some input — don't return all users for empty query
		if query == "" {
			writeSearchJSON(w, []searchOption{})
			return
		}

		wuResp, err := deps.ListWorkspaceUsers(ctx, &workspaceuserpb.ListWorkspaceUsersRequest{})
		if err != nil {
			log.Printf("Failed to search users: %v", err)
			writeSearchJSON(w, []searchOption{})
			return
		}

		var results []searchOption
		for _, wu := range wuResp.GetData() {
			if !wu.GetActive() {
				continue
			}
			user := wu.GetUser()
			if user == nil {
				continue
			}
			name := user.GetFirstName() + " " + user.GetLastName()
			email := user.GetEmailAddress()
			label := strings.TrimSpace(name)
			if email != "" {
				label = label + " (" + email + ")"
			}

			// Filter by query
			if query != "" && !strings.Contains(strings.ToLower(label), query) {
				continue
			}

			results = append(results, searchOption{
				Value: wu.GetId(),
				Label: label,
			})
		}

		if results == nil {
			results = []searchOption{}
		}
		writeSearchJSON(w, results)
	}
}

// writeSearchJSON marshals data as JSON and writes it to the response writer.
func writeSearchJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("search users: failed to encode JSON response: %v", err)
	}
}

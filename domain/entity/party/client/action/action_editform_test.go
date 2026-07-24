package action

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	entityclient "github.com/erniealice/entydad-golang/domain/entity/party/client"
	clientform "github.com/erniealice/entydad-golang/domain/entity/party/client/form"
)

// D-CLIENT (signer input): the GET edit drawer must emit a query-less FormAction
// (the bare, signed edit path) and carry `mode` as data so the template renders
// it as a hidden field. A query-bearing FormAction ("...edit/{id}?mode=") is
// signed over the querystring, but the action_workspace_guard verifies the bare
// r.URL.Path — so the old form 409'd on every save.
func TestNewEditAction_GET_FormActionIsQueryless(t *testing.T) {
	t.Parallel()

	rec := &clientActionRecorder{}
	deps := &Deps{
		ReadClient: rec.readClient,
		Routes:     entityclient.Routes{EditURL: "/action/clients/{id}/edit"},
	}
	req := httptest.NewRequest(http.MethodGet, "/action/clients/cl-1/edit?mode=info", nil)
	res := runHandler(t, NewEditAction(deps), withPerms("client:update"), req)

	data, ok := res.Data.(*clientform.Data)
	if !ok {
		t.Fatalf("res.Data = %T, want *clientform.Data", res.Data)
	}
	if strings.Contains(data.FormAction, "?") {
		t.Fatalf("FormAction must be query-less, got %q", data.FormAction)
	}
	if data.Mode != "info" {
		t.Fatalf("Mode = %q, want %q (mode must ride as a hidden field, not the FormAction querystring)", data.Mode, "info")
	}
}

// D-CLIENT (POST side): mode now arrives in the form BODY (the hidden field),
// not the URL querystring. The handler must read it via r.FormValue so the
// detail-tab redirect still fires.
func TestNewEditAction_POST_ModeFromBody(t *testing.T) {
	t.Parallel()

	rec := &clientActionRecorder{}
	deps := &Deps{
		UpdateClient: rec.updateClient,
		ReadClient:   rec.readClient,
		Routes: entityclient.Routes{
			EditURL:   "/action/clients/{id}/edit",
			DetailURL: "/clients/{id}",
		},
	}
	// Bare path (no ?mode=); mode carried in the POST body only.
	req := makePostRequest("/action/clients/cl-1/edit", url.Values{
		"mode": {"info"},
		"name": {"Updated"},
	})
	res := runHandler(t, NewEditAction(deps), withPerms("client:update"), req)

	if got := res.Headers["HX-Redirect"]; !strings.Contains(got, "tab=info") {
		t.Fatalf("HX-Redirect = %q, want a detail redirect with tab=info (mode read from the body)", got)
	}
	if len(rec.updateCalls) != 1 {
		t.Fatalf("UpdateClient calls = %d, want 1", len(rec.updateCalls))
	}
}

// List-page edit (no mode): a bare POST with no mode still succeeds and refreshes
// the table (the default branch) — confirms the empty-mode path is intact.
func TestNewEditAction_POST_NoModeRefreshesTable(t *testing.T) {
	t.Parallel()

	rec := &clientActionRecorder{}
	deps := &Deps{
		UpdateClient: rec.updateClient,
		ReadClient:   rec.readClient,
		Routes:       entityclient.Routes{EditURL: "/action/clients/{id}/edit"},
	}
	req := makePostRequest("/action/clients/cl-1/edit", url.Values{"name": {"Updated"}})
	res := runHandler(t, NewEditAction(deps), withPerms("client:update"), req)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want 200", res.StatusCode)
	}
	assertSuccessHeader(t, res, "clients-table")
}

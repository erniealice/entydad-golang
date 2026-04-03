package action

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

type deleteCall struct {
	id string
}

type statusCall struct {
	id     string
	active bool
}

type clientActionRecorder struct {
	deleteCalls   []deleteCall
	statusCalls   []statusCall
	deleteErrByID map[string]error
	statusErrByID map[string]error
}

func (r *clientActionRecorder) deleteClient(_ context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error) {
	id := req.GetData().GetId()
	r.deleteCalls = append(r.deleteCalls, deleteCall{id: id})
	if err := r.deleteErrByID[id]; err != nil {
		return nil, err
	}
	return &clientpb.DeleteClientResponse{}, nil
}

func (r *clientActionRecorder) setClientActive(_ context.Context, id string, active bool) error {
	r.statusCalls = append(r.statusCalls, statusCall{id: id, active: active})
	return r.statusErrByID[id]
}

func testMessages() map[string]string {
	return map[string]string{
		"shared.errors.permissionDenied":    "permission denied",
		"shared.errors.idRequired":          "id required",
		"shared.errors.noIdsProvided":       "no ids provided",
		"shared.errors.invalidStatus":       "invalid status",
		"shared.errors.invalidTargetStatus": "invalid target status",
	}
}

func makePostRequest(rawURL string, form url.Values) *http.Request {
	body := ""
	if form != nil {
		body = form.Encode()
	}
	req := httptest.NewRequest(http.MethodPost, rawURL, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func withPerms(codes ...string) context.Context {
	return view.WithUserPermissions(context.Background(), pyezatypes.NewUserPermissions(codes))
}

func runHandler(t *testing.T, h view.View, ctx context.Context, req *http.Request) view.ViewResult {
	t.Helper()
	return h.Handle(ctx, &view.ViewContext{
		Request:  req,
		Messages: testMessages(),
	})
}

func assertErrorHeader(t *testing.T, res view.ViewResult, want string) {
	t.Helper()
	if got := res.Headers["HX-Error-Message"]; got != want {
		t.Fatalf("HX-Error-Message = %q, want %q", got, want)
	}
}

func assertSuccessHeader(t *testing.T, res view.ViewResult, tableID string) {
	t.Helper()
	want := `{"formSuccess":true,"refreshTable":"` + tableID + `"}`
	if got := res.Headers["HX-Trigger"]; got != want {
		t.Fatalf("HX-Trigger = %q, want %q", got, want)
	}
}

func TestNewDeleteAction_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		req             *http.Request
		deleteErrByID   map[string]error
		wantStatus      int
		wantErrorHeader string
		wantDeleteIDs   []string
	}{
		{
			name:            "permission denied",
			perms:           []string{"client:read"},
			req:             makePostRequest("/action/clients/delete?id=cl-1", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "missing id",
			perms:           []string{"client:delete"},
			req:             makePostRequest("/action/clients/delete", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:          "uses id from query and succeeds",
			perms:         []string{"client:delete"},
			req:           makePostRequest("/action/clients/delete?id=cl-q", nil),
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"cl-q"},
		},
		{
			name:          "falls back to form id and succeeds",
			perms:         []string{"client:delete"},
			req:           makePostRequest("/action/clients/delete", url.Values{"id": {"cl-f"}}),
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"cl-f"},
		},
		{
			name:            "dependency error returns htmx error",
			perms:           []string{"client:delete"},
			req:             makePostRequest("/action/clients/delete?id=cl-e", nil),
			deleteErrByID:   map[string]error{"cl-e": errors.New("delete failed")},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "delete failed",
			wantDeleteIDs:   []string{"cl-e"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &clientActionRecorder{deleteErrByID: tt.deleteErrByID}
			deps := &Deps{DeleteClient: rec.deleteClient}
			res := runHandler(t, NewDeleteAction(deps), withPerms(tt.perms...), tt.req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "clients-table")
			}

			gotIDs := make([]string, 0, len(rec.deleteCalls))
			for _, c := range rec.deleteCalls {
				gotIDs = append(gotIDs, c.id)
			}
			if strings.Join(gotIDs, ",") != strings.Join(tt.wantDeleteIDs, ",") {
				t.Fatalf("DeleteClient IDs = %v, want %v", gotIDs, tt.wantDeleteIDs)
			}
		})
	}
}

func TestNewBulkDeleteAction_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		form            url.Values
		deleteErrByID   map[string]error
		wantStatus      int
		wantErrorHeader string
		wantDeleteIDs   []string
	}{
		{
			name:            "permission denied",
			perms:           []string{"client:read"},
			form:            url.Values{"id": {"cl-1"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "no ids provided",
			perms:           []string{"client:delete"},
			form:            url.Values{},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "no ids provided",
		},
		{
			name:          "success deletes all ids",
			perms:         []string{"client:delete"},
			form:          url.Values{"id": {"cl-1", "cl-2", "cl-3"}},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"cl-1", "cl-2", "cl-3"},
		},
		{
			name:          "partial failure still returns success",
			perms:         []string{"client:delete"},
			form:          url.Values{"id": {"cl-1", "cl-2", "cl-3"}},
			deleteErrByID: map[string]error{"cl-2": errors.New("boom")},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"cl-1", "cl-2", "cl-3"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &clientActionRecorder{deleteErrByID: tt.deleteErrByID}
			deps := &Deps{DeleteClient: rec.deleteClient}
			req := makePostRequest("/action/clients/bulk-delete", tt.form)
			res := runHandler(t, NewBulkDeleteAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "clients-table")
			}

			gotIDs := make([]string, 0, len(rec.deleteCalls))
			for _, c := range rec.deleteCalls {
				gotIDs = append(gotIDs, c.id)
			}
			if strings.Join(gotIDs, ",") != strings.Join(tt.wantDeleteIDs, ",") {
				t.Fatalf("DeleteClient IDs = %v, want %v", gotIDs, tt.wantDeleteIDs)
			}
		})
	}
}

func TestNewSetStatusAction_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		req             *http.Request
		statusErrByID   map[string]error
		wantStatus      int
		wantErrorHeader string
		wantCalls       []statusCall
	}{
		{
			name:            "permission denied",
			perms:           []string{"client:read"},
			req:             makePostRequest("/action/clients/set-status?id=cl-1&status=active", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "missing id",
			perms:           []string{"client:update"},
			req:             makePostRequest("/action/clients/set-status?status=active", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:            "invalid status",
			perms:           []string{"client:update"},
			req:             makePostRequest("/action/clients/set-status?id=cl-1&status=paused", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "invalid status",
		},
		{
			name:       "query params active success",
			perms:      []string{"client:update"},
			req:        makePostRequest("/action/clients/set-status?id=cl-1&status=active", nil),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "cl-1", active: true}},
		},
		{
			name:       "form fallback inactive success",
			perms:      []string{"client:update"},
			req:        makePostRequest("/action/clients/set-status", url.Values{"id": {"cl-2"}, "status": {"inactive"}}),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "cl-2", active: false}},
		},
		{
			name:            "dependency error",
			perms:           []string{"client:update"},
			req:             makePostRequest("/action/clients/set-status?id=cl-3&status=inactive", nil),
			statusErrByID:   map[string]error{"cl-3": errors.New("set status failed")},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "set status failed",
			wantCalls:       []statusCall{{id: "cl-3", active: false}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &clientActionRecorder{statusErrByID: tt.statusErrByID}
			deps := &Deps{SetClientActive: rec.setClientActive}
			res := runHandler(t, NewSetStatusAction(deps), withPerms(tt.perms...), tt.req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "clients-table")
			}

			if len(rec.statusCalls) != len(tt.wantCalls) {
				t.Fatalf("SetClientActive call count = %d, want %d", len(rec.statusCalls), len(tt.wantCalls))
			}
			for i := range tt.wantCalls {
				if rec.statusCalls[i] != tt.wantCalls[i] {
					t.Fatalf("SetClientActive call[%d] = %+v, want %+v", i, rec.statusCalls[i], tt.wantCalls[i])
				}
			}
		})
	}
}

func TestNewBulkSetStatusAction_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		form            url.Values
		statusErrByID   map[string]error
		wantStatus      int
		wantErrorHeader string
		wantCalls       []statusCall
	}{
		{
			name:            "permission denied",
			perms:           []string{"client:read"},
			form:            url.Values{"id": {"cl-1"}, "target_status": {"active"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "no ids provided",
			perms:           []string{"client:update"},
			form:            url.Values{"target_status": {"active"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "no ids provided",
		},
		{
			name:            "invalid target status",
			perms:           []string{"client:update"},
			form:            url.Values{"id": {"cl-1"}, "target_status": {"paused"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "invalid target status",
		},
		{
			name:       "active success",
			perms:      []string{"client:update"},
			form:       url.Values{"id": {"cl-1", "cl-2"}, "target_status": {"active"}},
			wantStatus: http.StatusOK,
			wantCalls: []statusCall{
				{id: "cl-1", active: true},
				{id: "cl-2", active: true},
			},
		},
		{
			name:       "inactive success",
			perms:      []string{"client:update"},
			form:       url.Values{"id": {"cl-1", "cl-2"}, "target_status": {"inactive"}},
			wantStatus: http.StatusOK,
			wantCalls: []statusCall{
				{id: "cl-1", active: false},
				{id: "cl-2", active: false},
			},
		},
		{
			name:          "partial failure still returns success",
			perms:         []string{"client:update"},
			form:          url.Values{"id": {"cl-1", "cl-2", "cl-3"}, "target_status": {"inactive"}},
			statusErrByID: map[string]error{"cl-2": errors.New("boom")},
			wantStatus:    http.StatusOK,
			wantCalls: []statusCall{
				{id: "cl-1", active: false},
				{id: "cl-2", active: false},
				{id: "cl-3", active: false},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &clientActionRecorder{statusErrByID: tt.statusErrByID}
			deps := &Deps{SetClientActive: rec.setClientActive}
			req := makePostRequest("/action/clients/bulk-set-status", tt.form)
			res := runHandler(t, NewBulkSetStatusAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "clients-table")
			}

			if len(rec.statusCalls) != len(tt.wantCalls) {
				t.Fatalf("SetClientActive call count = %d, want %d", len(rec.statusCalls), len(tt.wantCalls))
			}
			for i := range tt.wantCalls {
				if rec.statusCalls[i] != tt.wantCalls[i] {
					t.Fatalf("SetClientActive call[%d] = %+v, want %+v", i, rec.statusCalls[i], tt.wantCalls[i])
				}
			}
		})
	}
}

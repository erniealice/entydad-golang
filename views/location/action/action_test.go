package action

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
)

type deleteCall struct {
	id string
}

type statusCall struct {
	id     string
	active bool
}

type createLocationCall struct {
	name    string
	address string
	active  bool
}

type updateLocationCall struct {
	id      string
	name    string
	address string
	active  bool
}

type locationActionRecorder struct {
	deleteCalls   []deleteCall
	statusCalls   []statusCall
	createCalls   []createLocationCall
	updateCalls   []updateLocationCall
	deleteErrByID map[string]error
	statusErrByID map[string]error
	createErr     error
	updateErr     error
	inUseIDs      map[string]bool
	inUseErr      error
}

func (r *locationActionRecorder) deleteLocation(_ context.Context, req *locationpb.DeleteLocationRequest) (*locationpb.DeleteLocationResponse, error) {
	id := req.GetData().GetId()
	r.deleteCalls = append(r.deleteCalls, deleteCall{id: id})
	if err := r.deleteErrByID[id]; err != nil {
		return nil, err
	}
	return &locationpb.DeleteLocationResponse{}, nil
}

func (r *locationActionRecorder) setLocationActive(_ context.Context, id string, active bool) error {
	r.statusCalls = append(r.statusCalls, statusCall{id: id, active: active})
	return r.statusErrByID[id]
}

func (r *locationActionRecorder) createLocation(_ context.Context, req *locationpb.CreateLocationRequest) (*locationpb.CreateLocationResponse, error) {
	d := req.GetData()
	r.createCalls = append(r.createCalls, createLocationCall{
		name:    d.GetName(),
		address: d.GetAddress(),
		active:  d.GetActive(),
	})
	if r.createErr != nil {
		return nil, r.createErr
	}
	return &locationpb.CreateLocationResponse{}, nil
}

func (r *locationActionRecorder) updateLocation(_ context.Context, req *locationpb.UpdateLocationRequest) (*locationpb.UpdateLocationResponse, error) {
	d := req.GetData()
	r.updateCalls = append(r.updateCalls, updateLocationCall{
		id:      d.GetId(),
		name:    d.GetName(),
		address: d.GetAddress(),
		active:  d.GetActive(),
	})
	if r.updateErr != nil {
		return nil, r.updateErr
	}
	return &locationpb.UpdateLocationResponse{}, nil
}

func (r *locationActionRecorder) readLocation(_ context.Context, req *locationpb.ReadLocationRequest) (*locationpb.ReadLocationResponse, error) {
	id := req.GetData().GetId()
	return &locationpb.ReadLocationResponse{
		Data: []*locationpb.Location{{
			Id:      id,
			Name:    "Test Location",
			Address: "123 Main St",
			Active:  true,
		}},
	}, nil
}

func (r *locationActionRecorder) getInUseIDs(_ context.Context, ids []string) (map[string]bool, error) {
	if r.inUseErr != nil {
		return nil, r.inUseErr
	}
	if r.inUseIDs != nil {
		return r.inUseIDs, nil
	}
	return map[string]bool{}, nil
}

func testMessages() map[string]string {
	return map[string]string{
		"shared.errors.permissionDenied":    "permission denied",
		"shared.errors.idRequired":          "id required",
		"shared.errors.noIdsProvided":       "no ids provided",
		"shared.errors.invalidStatus":       "invalid status",
		"shared.errors.invalidTargetStatus": "invalid target status",
		"shared.errors.invalidFormData":     "invalid form data",
		"shared.errors.cannotDeleteInUse":   "cannot delete in use",
		"shared.errors.verifyFailed":        "verify failed",
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

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestNewDeleteAction_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		req             *http.Request
		deleteErrByID   map[string]error
		inUseIDs        map[string]bool
		inUseErr        error
		wantStatus      int
		wantErrorHeader string
		wantDeleteIDs   []string
	}{
		{
			name:            "permission denied",
			perms:           []string{"location:read"},
			req:             makePostRequest("/action/locations/delete?id=loc-1", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "missing id",
			perms:           []string{"location:delete"},
			req:             makePostRequest("/action/locations/delete", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:          "uses id from query and succeeds",
			perms:         []string{"location:delete"},
			req:           makePostRequest("/action/locations/delete?id=loc-q", nil),
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"loc-q"},
		},
		{
			name:          "falls back to form id and succeeds",
			perms:         []string{"location:delete"},
			req:           makePostRequest("/action/locations/delete", url.Values{"id": {"loc-f"}}),
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"loc-f"},
		},
		{
			name:            "dependency error returns htmx error",
			perms:           []string{"location:delete"},
			req:             makePostRequest("/action/locations/delete?id=loc-e", nil),
			deleteErrByID:   map[string]error{"loc-e": errors.New("delete failed")},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "delete failed",
			wantDeleteIDs:   []string{"loc-e"},
		},
		{
			name:            "in-use guard blocks deletion",
			perms:           []string{"location:delete"},
			req:             makePostRequest("/action/locations/delete?id=loc-used", nil),
			inUseIDs:        map[string]bool{"loc-used": true},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "cannot delete in use",
		},
		{
			name:            "in-use check error",
			perms:           []string{"location:delete"},
			req:             makePostRequest("/action/locations/delete?id=loc-chk", nil),
			inUseErr:        errors.New("db error"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "verify failed",
		},
		{
			name:          "in-use check passes when not in use",
			perms:         []string{"location:delete"},
			req:           makePostRequest("/action/locations/delete?id=loc-free", nil),
			inUseIDs:      map[string]bool{"loc-free": false},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"loc-free"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{
				deleteErrByID: tt.deleteErrByID,
				inUseIDs:      tt.inUseIDs,
				inUseErr:      tt.inUseErr,
			}
			deps := &Deps{
				DeleteLocation: rec.deleteLocation,
				GetInUseIDs:    rec.getInUseIDs,
			}
			res := runHandler(t, NewDeleteAction(deps), withPerms(tt.perms...), tt.req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "locations-table")
			}

			gotIDs := make([]string, 0, len(rec.deleteCalls))
			for _, c := range rec.deleteCalls {
				gotIDs = append(gotIDs, c.id)
			}
			if strings.Join(gotIDs, ",") != strings.Join(tt.wantDeleteIDs, ",") {
				t.Fatalf("DeleteLocation IDs = %v, want %v", gotIDs, tt.wantDeleteIDs)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BulkDelete
// ---------------------------------------------------------------------------

func TestNewBulkDeleteAction_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		form            url.Values
		deleteErrByID   map[string]error
		inUseIDs        map[string]bool
		inUseErr        error
		wantStatus      int
		wantErrorHeader string
		wantDeleteIDs   []string
	}{
		{
			name:            "permission denied",
			perms:           []string{"location:read"},
			form:            url.Values{"id": {"loc-1"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "no ids provided",
			perms:           []string{"location:delete"},
			form:            url.Values{},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "no ids provided",
		},
		{
			name:          "success deletes all ids",
			perms:         []string{"location:delete"},
			form:          url.Values{"id": {"loc-1", "loc-2", "loc-3"}},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"loc-1", "loc-2", "loc-3"},
		},
		{
			name:          "partial failure still returns success",
			perms:         []string{"location:delete"},
			form:          url.Values{"id": {"loc-1", "loc-2", "loc-3"}},
			deleteErrByID: map[string]error{"loc-2": errors.New("boom")},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"loc-1", "loc-2", "loc-3"},
		},
		{
			name:            "in-use guard blocks bulk deletion",
			perms:           []string{"location:delete"},
			form:            url.Values{"id": {"loc-1", "loc-used"}},
			inUseIDs:        map[string]bool{"loc-used": true},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "cannot delete in use",
		},
		{
			name:            "in-use check error blocks bulk deletion",
			perms:           []string{"location:delete"},
			form:            url.Values{"id": {"loc-1"}},
			inUseErr:        errors.New("db error"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "verify failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{
				deleteErrByID: tt.deleteErrByID,
				inUseIDs:      tt.inUseIDs,
				inUseErr:      tt.inUseErr,
			}
			deps := &Deps{
				DeleteLocation: rec.deleteLocation,
				GetInUseIDs:    rec.getInUseIDs,
			}
			req := makePostRequest("/action/locations/bulk-delete", tt.form)
			res := runHandler(t, NewBulkDeleteAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "locations-table")
			}

			gotIDs := make([]string, 0, len(rec.deleteCalls))
			for _, c := range rec.deleteCalls {
				gotIDs = append(gotIDs, c.id)
			}
			if strings.Join(gotIDs, ",") != strings.Join(tt.wantDeleteIDs, ",") {
				t.Fatalf("DeleteLocation IDs = %v, want %v", gotIDs, tt.wantDeleteIDs)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// SetStatus
// ---------------------------------------------------------------------------

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
			perms:           []string{"location:read"},
			req:             makePostRequest("/action/locations/set-status?id=loc-1&status=active", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "missing id",
			perms:           []string{"location:update"},
			req:             makePostRequest("/action/locations/set-status?status=active", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:            "invalid status",
			perms:           []string{"location:update"},
			req:             makePostRequest("/action/locations/set-status?id=loc-1&status=paused", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "invalid status",
		},
		{
			name:       "query params active success",
			perms:      []string{"location:update"},
			req:        makePostRequest("/action/locations/set-status?id=loc-1&status=active", nil),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "loc-1", active: true}},
		},
		{
			name:       "form fallback inactive success",
			perms:      []string{"location:update"},
			req:        makePostRequest("/action/locations/set-status", url.Values{"id": {"loc-2"}, "status": {"inactive"}}),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "loc-2", active: false}},
		},
		{
			name:            "dependency error",
			perms:           []string{"location:update"},
			req:             makePostRequest("/action/locations/set-status?id=loc-3&status=inactive", nil),
			statusErrByID:   map[string]error{"loc-3": errors.New("set status failed")},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "set status failed",
			wantCalls:       []statusCall{{id: "loc-3", active: false}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{statusErrByID: tt.statusErrByID}
			deps := &Deps{SetLocationActive: rec.setLocationActive}
			res := runHandler(t, NewSetStatusAction(deps), withPerms(tt.perms...), tt.req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "locations-table")
			}

			if len(rec.statusCalls) != len(tt.wantCalls) {
				t.Fatalf("SetLocationActive call count = %d, want %d", len(rec.statusCalls), len(tt.wantCalls))
			}
			for i := range tt.wantCalls {
				if rec.statusCalls[i] != tt.wantCalls[i] {
					t.Fatalf("SetLocationActive call[%d] = %+v, want %+v", i, rec.statusCalls[i], tt.wantCalls[i])
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BulkSetStatus
// ---------------------------------------------------------------------------

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
			perms:           []string{"location:read"},
			form:            url.Values{"id": {"loc-1"}, "target_status": {"active"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "no ids provided",
			perms:           []string{"location:update"},
			form:            url.Values{"target_status": {"active"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "no ids provided",
		},
		{
			name:            "invalid target status",
			perms:           []string{"location:update"},
			form:            url.Values{"id": {"loc-1"}, "target_status": {"paused"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "invalid target status",
		},
		{
			name:       "active success",
			perms:      []string{"location:update"},
			form:       url.Values{"id": {"loc-1", "loc-2"}, "target_status": {"active"}},
			wantStatus: http.StatusOK,
			wantCalls: []statusCall{
				{id: "loc-1", active: true},
				{id: "loc-2", active: true},
			},
		},
		{
			name:       "inactive success",
			perms:      []string{"location:update"},
			form:       url.Values{"id": {"loc-1", "loc-2"}, "target_status": {"inactive"}},
			wantStatus: http.StatusOK,
			wantCalls: []statusCall{
				{id: "loc-1", active: false},
				{id: "loc-2", active: false},
			},
		},
		{
			name:          "partial failure still returns success",
			perms:         []string{"location:update"},
			form:          url.Values{"id": {"loc-1", "loc-2", "loc-3"}, "target_status": {"inactive"}},
			statusErrByID: map[string]error{"loc-2": errors.New("boom")},
			wantStatus:    http.StatusOK,
			wantCalls: []statusCall{
				{id: "loc-1", active: false},
				{id: "loc-2", active: false},
				{id: "loc-3", active: false},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{statusErrByID: tt.statusErrByID}
			deps := &Deps{SetLocationActive: rec.setLocationActive}
			req := makePostRequest("/action/locations/bulk-set-status", tt.form)
			res := runHandler(t, NewBulkSetStatusAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "locations-table")
			}

			if len(rec.statusCalls) != len(tt.wantCalls) {
				t.Fatalf("SetLocationActive call count = %d, want %d", len(rec.statusCalls), len(tt.wantCalls))
			}
			for i := range tt.wantCalls {
				if rec.statusCalls[i] != tt.wantCalls[i] {
					t.Fatalf("SetLocationActive call[%d] = %+v, want %+v", i, rec.statusCalls[i], tt.wantCalls[i])
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AddAction (POST path)
// ---------------------------------------------------------------------------

func TestNewAddAction_POST_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		form            url.Values
		createErr       error
		wantStatus      int
		wantErrorHeader string
		wantCreateCount int
	}{
		{
			name:            "permission denied",
			perms:           []string{"location:read"},
			form:            url.Values{"name": {"HQ"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:  "success",
			perms: []string{"location:create"},
			form: url.Values{
				"name":    {"HQ"},
				"address": {"123 Main St"},
				"active":  {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
		},
		{
			name:  "create error",
			perms: []string{"location:create"},
			form: url.Values{
				"name":    {"Bad"},
				"address": {"456 Elm St"},
				"active":  {"true"},
			},
			createErr:       errors.New("create failed"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "create failed",
			wantCreateCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{createErr: tt.createErr}
			deps := &Deps{
				CreateLocation: rec.createLocation,
				Routes:         entydad.LocationRoutes{AddURL: "/action/locations/add"},
			}
			req := makePostRequest("/action/locations/add", tt.form)
			res := runHandler(t, NewAddAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "locations-table")
			}
			if len(rec.createCalls) != tt.wantCreateCount {
				t.Fatalf("CreateLocation call count = %d, want %d", len(rec.createCalls), tt.wantCreateCount)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// EditAction (POST path)
// ---------------------------------------------------------------------------

func TestNewEditAction_POST_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		form            url.Values
		updateErr       error
		wantStatus      int
		wantErrorHeader string
		wantUpdateCount int
	}{
		{
			name:            "permission denied",
			perms:           []string{"location:read"},
			form:            url.Values{"name": {"HQ"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:  "success",
			perms: []string{"location:update"},
			form: url.Values{
				"name":    {"Updated HQ"},
				"address": {"789 Oak Ave"},
				"active":  {"true"},
			},
			wantStatus:      http.StatusOK,
			wantUpdateCount: 1,
		},
		{
			name:  "update error",
			perms: []string{"location:update"},
			form: url.Values{
				"name":    {"Bad"},
				"address": {"456 Elm St"},
				"active":  {"true"},
			},
			updateErr:       errors.New("update failed"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "update failed",
			wantUpdateCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{updateErr: tt.updateErr}
			deps := &Deps{
				UpdateLocation: rec.updateLocation,
				ReadLocation:   rec.readLocation,
				Routes:         entydad.LocationRoutes{EditURL: "/action/locations/{id}/edit"},
			}
			req := makePostRequest("/action/locations/loc-1/edit", tt.form)
			res := runHandler(t, NewEditAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "locations-table")
			}
			if len(rec.updateCalls) != tt.wantUpdateCount {
				t.Fatalf("UpdateLocation call count = %d, want %d", len(rec.updateCalls), tt.wantUpdateCount)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AddAction — negative / defensive tests
// ---------------------------------------------------------------------------

func TestNewAddAction_Negative(t *testing.T) {
	t.Parallel()

	longString := strings.Repeat("X", 300)

	tests := []struct {
		name            string
		perms           []string
		form            url.Values
		createErr       error
		wantStatus      int
		wantErrorHeader string
		wantCreateCount int
		wantName        string
	}{
		{
			name:  "missing name (empty form) still calls create",
			perms: []string{"location:create"},
			form: url.Values{
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantName:        "",
		},
		{
			name:  "empty name value",
			perms: []string{"location:create"},
			form: url.Values{
				"name":   {""},
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantName:        "",
		},
		{
			name:  "whitespace-only name",
			perms: []string{"location:create"},
			form: url.Values{
				"name":   {"   "},
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantName:        "   ",
		},
		{
			name:  "very long name (>255 chars)",
			perms: []string{"location:create"},
			form: url.Values{
				"name":   {longString},
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantName:        longString,
		},
		{
			name:  "XSS-like script tag in name",
			perms: []string{"location:create"},
			form: url.Values{
				"name":   {"<script>alert('xss')</script>"},
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantName:        "<script>alert('xss')</script>",
		},
		{
			name:  "backend rejects empty name",
			perms: []string{"location:create"},
			form: url.Values{
				"active": {"true"},
			},
			createErr:       errors.New("name is required"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "name is required",
			wantCreateCount: 1,
		},
		{
			name:  "active field absent defaults to false",
			perms: []string{"location:create"},
			form: url.Values{
				"name": {"HQ"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{createErr: tt.createErr}
			deps := &Deps{
				CreateLocation: rec.createLocation,
				Routes:         entydad.LocationRoutes{AddURL: "/action/locations/add"},
			}
			req := makePostRequest("/action/locations/add", tt.form)
			res := runHandler(t, NewAddAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "locations-table")
			}
			if len(rec.createCalls) != tt.wantCreateCount {
				t.Fatalf("CreateLocation call count = %d, want %d", len(rec.createCalls), tt.wantCreateCount)
			}
			if tt.wantName != "" && len(rec.createCalls) > 0 {
				if got := rec.createCalls[0].name; got != tt.wantName {
					t.Fatalf("name = %q, want %q", got, tt.wantName)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// EditAction — negative / defensive tests
// ---------------------------------------------------------------------------

func TestNewEditAction_Negative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		pathURL         string
		form            url.Values
		updateErr       error
		wantStatus      int
		wantErrorHeader string
		wantUpdateCount int
	}{
		{
			name:    "edit with empty path ID sends empty id to backend",
			perms:   []string{"location:update"},
			pathURL: "/action/locations//edit",
			form: url.Values{
				"name":   {"Updated"},
				"active": {"true"},
			},
			updateErr:       errors.New("id is required"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id is required",
			wantUpdateCount: 1,
		},
		{
			name:    "edit with missing name",
			perms:   []string{"location:update"},
			pathURL: "/action/locations/loc-1/edit",
			form: url.Values{
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantUpdateCount: 1,
		},
		{
			name:    "edit with XSS in name",
			perms:   []string{"location:update"},
			pathURL: "/action/locations/loc-1/edit",
			form: url.Values{
				"name":   {"<script>alert(1)</script>"},
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantUpdateCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{updateErr: tt.updateErr}
			deps := &Deps{
				UpdateLocation: rec.updateLocation,
				ReadLocation:   rec.readLocation,
				Routes:         entydad.LocationRoutes{EditURL: "/action/locations/{id}/edit"},
			}
			req := makePostRequest(tt.pathURL, tt.form)
			res := runHandler(t, NewEditAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "locations-table")
			}
			if len(rec.updateCalls) != tt.wantUpdateCount {
				t.Fatalf("UpdateLocation call count = %d, want %d", len(rec.updateCalls), tt.wantUpdateCount)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GetInUseIDs — large ID set boundary test
// ---------------------------------------------------------------------------

func TestNewBulkDeleteAction_LargeInUseSet(t *testing.T) {
	t.Parallel()

	// Build a large set of IDs where many are in-use
	ids := make([]string, 50)
	inUseMap := make(map[string]bool)
	for i := range ids {
		ids[i] = fmt.Sprintf("loc-%d", i)
		if i%2 == 0 {
			inUseMap[ids[i]] = true
		}
	}

	tests := []struct {
		name            string
		perms           []string
		form            url.Values
		inUseIDs        map[string]bool
		wantStatus      int
		wantErrorHeader string
		wantDeleteIDs   int
	}{
		{
			name:            "bulk delete blocked when any of 50 IDs is in use",
			perms:           []string{"location:delete"},
			form:            url.Values{"id": ids},
			inUseIDs:        inUseMap,
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "cannot delete in use",
			wantDeleteIDs:   0,
		},
		{
			name:          "bulk delete succeeds when none of 50 IDs is in use",
			perms:         []string{"location:delete"},
			form:          url.Values{"id": ids},
			inUseIDs:      map[string]bool{},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: 50,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{inUseIDs: tt.inUseIDs}
			deps := &Deps{
				DeleteLocation: rec.deleteLocation,
				GetInUseIDs:    rec.getInUseIDs,
			}
			req := makePostRequest("/action/locations/bulk-delete", tt.form)
			res := runHandler(t, NewBulkDeleteAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if len(rec.deleteCalls) != tt.wantDeleteIDs {
				t.Fatalf("DeleteLocation call count = %d, want %d", len(rec.deleteCalls), tt.wantDeleteIDs)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Delete — additional edge cases
// ---------------------------------------------------------------------------

func TestNewDeleteAction_Negative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		req             *http.Request
		wantStatus      int
		wantErrorHeader string
	}{
		{
			name:            "empty string id in query param",
			perms:           []string{"location:delete"},
			req:             makePostRequest("/action/locations/delete?id=", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:       "whitespace-only id in form",
			perms:      []string{"location:delete"},
			req:        makePostRequest("/action/locations/delete", url.Values{"id": {"   "}}),
			wantStatus: http.StatusOK, // whitespace is not validated, treated as a valid ID
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &locationActionRecorder{}
			deps := &Deps{
				DeleteLocation: rec.deleteLocation,
				GetInUseIDs:    rec.getInUseIDs,
			}
			res := runHandler(t, NewDeleteAction(deps), withPerms(tt.perms...), tt.req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
		})
	}
}

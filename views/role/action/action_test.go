package action

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

const rolesTableTrigger = `{"formSuccess":true,"refreshTable":"roles-table"}`

var testMessages = map[string]string{
	"shared.errors.permissionDenied":    "permission denied",
	"shared.errors.idRequired":          "id required",
	"shared.errors.noIdsProvided":       "no ids provided",
	"shared.errors.invalidStatus":       "invalid status",
	"shared.errors.invalidTargetStatus": "invalid target status",
}

type setRoleActiveCall struct {
	id     string
	active bool
}

func TestNewDeleteAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		perms         []string
		req           func(t *testing.T) *http.Request
		deleteErrByID map[string]error
		wantStatus    int
		wantError     string
		wantTrigger   string
		wantDeleted   []string
	}{
		{
			name:  "permission denied",
			perms: nil,
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/delete?id=role-1", nil)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.permissionDenied"],
		},
		{
			name:  "missing id in query and form",
			perms: []string{"role:delete"},
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/delete", nil)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.idRequired"],
		},
		{
			name:  "uses form id fallback",
			perms: []string{"role:delete"},
			req: func(t *testing.T) *http.Request {
				return newFormRequest(t, http.MethodPost, "/action/roles/delete", url.Values{"id": {"role-form"}})
			},
			wantStatus:  http.StatusOK,
			wantTrigger: rolesTableTrigger,
			wantDeleted: []string{"role-form"},
		},
		{
			name:  "dependency error returns htmx error",
			perms: []string{"role:delete"},
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/delete?id=role-err", nil)
			},
			deleteErrByID: map[string]error{
				"role-err": errors.New("delete failed"),
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantError:   "delete failed",
			wantDeleted: []string{"role-err"},
		},
		{
			name:  "success from query id",
			perms: []string{"role:delete"},
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/delete?id=role-1", nil)
			},
			wantStatus:  http.StatusOK,
			wantTrigger: rolesTableTrigger,
			wantDeleted: []string{"role-1"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			deletedIDs := make([]string, 0)
			deps := &Deps{
				DeleteRole: func(ctx context.Context, req *rolepb.DeleteRoleRequest) (*rolepb.DeleteRoleResponse, error) {
					id := req.GetData().GetId()
					deletedIDs = append(deletedIDs, id)
					if tc.deleteErrByID != nil {
						if err, ok := tc.deleteErrByID[id]; ok {
							return nil, err
						}
					}
					return &rolepb.DeleteRoleResponse{}, nil
				},
			}

			ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions(tc.perms))
			vc := &view.ViewContext{
				Request:  tc.req(t),
				Messages: testMessages,
			}

			got := NewDeleteAction(deps).Handle(ctx, vc)
			assertViewResult(t, got, tc.wantStatus, tc.wantError, tc.wantTrigger)
			if !slices.Equal(deletedIDs, tc.wantDeleted) {
				t.Fatalf("deleted IDs = %v, want %v", deletedIDs, tc.wantDeleted)
			}
		})
	}
}

func TestNewBulkDeleteAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		perms         []string
		req           func(t *testing.T) *http.Request
		deleteErrByID map[string]error
		wantStatus    int
		wantError     string
		wantTrigger   string
		wantDeleted   []string
	}{
		{
			name:  "permission denied",
			perms: nil,
			req: func(t *testing.T) *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/action/roles/bulk-delete", map[string][]string{"id": {"r1"}})
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.permissionDenied"],
		},
		{
			name:  "missing ids",
			perms: []string{"role:delete"},
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/bulk-delete", nil)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.noIdsProvided"],
		},
		{
			name:  "success with multiple ids",
			perms: []string{"role:delete"},
			req: func(t *testing.T) *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/action/roles/bulk-delete", map[string][]string{"id": {"r1", "r2", "r3"}})
			},
			wantStatus:  http.StatusOK,
			wantTrigger: rolesTableTrigger,
			wantDeleted: []string{"r1", "r2", "r3"},
		},
		{
			name:  "partial failure still succeeds and continues deleting",
			perms: []string{"role:delete"},
			req: func(t *testing.T) *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/action/roles/bulk-delete", map[string][]string{"id": {"r1", "r2", "r3"}})
			},
			deleteErrByID: map[string]error{
				"r2": errors.New("boom"),
			},
			wantStatus:  http.StatusOK,
			wantTrigger: rolesTableTrigger,
			wantDeleted: []string{"r1", "r2", "r3"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			deletedIDs := make([]string, 0)
			deps := &Deps{
				DeleteRole: func(ctx context.Context, req *rolepb.DeleteRoleRequest) (*rolepb.DeleteRoleResponse, error) {
					id := req.GetData().GetId()
					deletedIDs = append(deletedIDs, id)
					if tc.deleteErrByID != nil {
						if err, ok := tc.deleteErrByID[id]; ok {
							return nil, err
						}
					}
					return &rolepb.DeleteRoleResponse{}, nil
				},
			}

			ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions(tc.perms))
			vc := &view.ViewContext{
				Request:  tc.req(t),
				Messages: testMessages,
			}

			got := NewBulkDeleteAction(deps).Handle(ctx, vc)
			assertViewResult(t, got, tc.wantStatus, tc.wantError, tc.wantTrigger)
			if !slices.Equal(deletedIDs, tc.wantDeleted) {
				t.Fatalf("deleted IDs = %v, want %v", deletedIDs, tc.wantDeleted)
			}
		})
	}
}

func TestNewSetStatusAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		perms        []string
		req          func(t *testing.T) *http.Request
		setErrByID   map[string]error
		wantStatus   int
		wantError    string
		wantTrigger  string
		wantSetCalls []setRoleActiveCall
	}{
		{
			name:  "permission denied",
			perms: nil,
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/set-status?id=r1&status=active", nil)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.permissionDenied"],
		},
		{
			name:  "missing id",
			perms: []string{"role:update"},
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/set-status?status=active", nil)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.idRequired"],
		},
		{
			name:  "invalid status",
			perms: []string{"role:update"},
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/set-status?id=r1&status=blocked", nil)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.invalidStatus"],
		},
		{
			name:  "uses form fallback and sets active false",
			perms: []string{"role:update"},
			req: func(t *testing.T) *http.Request {
				return newFormRequest(t, http.MethodPost, "/action/roles/set-status", url.Values{
					"id":     {"r-form"},
					"status": {"inactive"},
				})
			},
			wantStatus:   http.StatusOK,
			wantTrigger:  rolesTableTrigger,
			wantSetCalls: []setRoleActiveCall{{id: "r-form", active: false}},
		},
		{
			name:  "dependency error",
			perms: []string{"role:update"},
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/set-status?id=r-err&status=active", nil)
			},
			setErrByID: map[string]error{
				"r-err": errors.New("set failed"),
			},
			wantStatus:   http.StatusUnprocessableEntity,
			wantError:    "set failed",
			wantSetCalls: []setRoleActiveCall{{id: "r-err", active: true}},
		},
		{
			name:  "success active",
			perms: []string{"role:update"},
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/action/roles/set-status?id=r1&status=active", nil)
			},
			wantStatus:   http.StatusOK,
			wantTrigger:  rolesTableTrigger,
			wantSetCalls: []setRoleActiveCall{{id: "r1", active: true}},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			setCalls := make([]setRoleActiveCall, 0)
			deps := &Deps{
				SetRoleActive: func(ctx context.Context, id string, active bool) error {
					setCalls = append(setCalls, setRoleActiveCall{id: id, active: active})
					if tc.setErrByID != nil {
						if err, ok := tc.setErrByID[id]; ok {
							return err
						}
					}
					return nil
				},
			}

			ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions(tc.perms))
			vc := &view.ViewContext{
				Request:  tc.req(t),
				Messages: testMessages,
			}

			got := NewSetStatusAction(deps).Handle(ctx, vc)
			assertViewResult(t, got, tc.wantStatus, tc.wantError, tc.wantTrigger)
			if !slices.Equal(setCalls, tc.wantSetCalls) {
				t.Fatalf("set calls = %v, want %v", setCalls, tc.wantSetCalls)
			}
		})
	}
}

func TestNewBulkSetStatusAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		perms        []string
		req          func(t *testing.T) *http.Request
		setErrByID   map[string]error
		wantStatus   int
		wantError    string
		wantTrigger  string
		wantSetCalls []setRoleActiveCall
	}{
		{
			name:  "permission denied",
			perms: nil,
			req: func(t *testing.T) *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/action/roles/bulk-set-status", map[string][]string{"id": {"r1"}, "target_status": {"active"}})
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.permissionDenied"],
		},
		{
			name:  "missing ids",
			perms: []string{"role:update"},
			req: func(t *testing.T) *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/action/roles/bulk-set-status", map[string][]string{"target_status": {"active"}})
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.noIdsProvided"],
		},
		{
			name:  "invalid target status",
			perms: []string{"role:update"},
			req: func(t *testing.T) *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/action/roles/bulk-set-status", map[string][]string{"id": {"r1"}, "target_status": {"blocked"}})
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  testMessages["shared.errors.invalidTargetStatus"],
		},
		{
			name:  "success active",
			perms: []string{"role:update"},
			req: func(t *testing.T) *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/action/roles/bulk-set-status", map[string][]string{"id": {"r1", "r2"}, "target_status": {"active"}})
			},
			wantStatus:  http.StatusOK,
			wantTrigger: rolesTableTrigger,
			wantSetCalls: []setRoleActiveCall{
				{id: "r1", active: true},
				{id: "r2", active: true},
			},
		},
		{
			name:  "partial failure still succeeds and continues updating",
			perms: []string{"role:update"},
			req: func(t *testing.T) *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/action/roles/bulk-set-status", map[string][]string{"id": {"r1", "r2", "r3"}, "target_status": {"inactive"}})
			},
			setErrByID:  map[string]error{"r2": errors.New("set failed")},
			wantStatus:  http.StatusOK,
			wantTrigger: rolesTableTrigger,
			wantSetCalls: []setRoleActiveCall{
				{id: "r1", active: false},
				{id: "r2", active: false},
				{id: "r3", active: false},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			setCalls := make([]setRoleActiveCall, 0)
			deps := &Deps{
				SetRoleActive: func(ctx context.Context, id string, active bool) error {
					setCalls = append(setCalls, setRoleActiveCall{id: id, active: active})
					if tc.setErrByID != nil {
						if err, ok := tc.setErrByID[id]; ok {
							return err
						}
					}
					return nil
				},
			}

			ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions(tc.perms))
			vc := &view.ViewContext{
				Request:  tc.req(t),
				Messages: testMessages,
			}

			got := NewBulkSetStatusAction(deps).Handle(ctx, vc)
			assertViewResult(t, got, tc.wantStatus, tc.wantError, tc.wantTrigger)
			if !slices.Equal(setCalls, tc.wantSetCalls) {
				t.Fatalf("set calls = %v, want %v", setCalls, tc.wantSetCalls)
			}
		})
	}
}

func assertViewResult(t *testing.T, got view.ViewResult, wantStatus int, wantError, wantTrigger string) {
	t.Helper()

	if got.StatusCode != wantStatus {
		t.Fatalf("status = %d, want %d", got.StatusCode, wantStatus)
	}

	gotErr := ""
	gotTrigger := ""
	if got.Headers != nil {
		gotErr = got.Headers["HX-Error-Message"]
		gotTrigger = got.Headers["HX-Trigger"]
	}

	if gotErr != wantError {
		t.Fatalf("HX-Error-Message = %q, want %q", gotErr, wantError)
	}
	if gotTrigger != wantTrigger {
		t.Fatalf("HX-Trigger = %q, want %q", gotTrigger, wantTrigger)
	}
}

func newFormRequest(t *testing.T, method, target string, form url.Values) *http.Request {
	t.Helper()

	req := httptest.NewRequest(method, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func newMultipartRequest(t *testing.T, method, target string, fields map[string][]string) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, values := range fields {
		for _, value := range values {
			if err := writer.WriteField(key, value); err != nil {
				t.Fatalf("WriteField(%q, %q): %v", key, value, err)
			}
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close: %v", err)
	}

	req := httptest.NewRequest(method, target, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

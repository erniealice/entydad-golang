package action

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

type supplierStatusCall struct {
	id     string
	status string
}

func TestNewSetStatusAction_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		req             *http.Request
		setErr          error
		wantStatus      int
		wantErrHeader   string
		wantTrigger     string
		wantStatusCalls []supplierStatusCall
	}{
		{
			name:          "permission denied when user lacks supplier:update",
			perms:         []string{},
			req:           httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-1&status=active", nil),
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrHeader: "permission denied",
		},
		{
			name:          "id required when id missing in query and form",
			perms:         []string{"supplier:update"},
			req:           newFormRequest(http.MethodPost, "/action/suppliers/status", map[string]string{"status": "active"}),
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrHeader: "id required",
		},
		{
			name:          "invalid status rejected",
			perms:         []string{"supplier:update"},
			req:           httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-1&status=paused", nil),
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrHeader: "invalid status",
		},
		{
			name:        "active persists active status",
			perms:       []string{"supplier:update"},
			req:         httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-1&status=active", nil),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-1", status: "active"},
			},
		},
		{
			name:        "blocked persists blocked status",
			perms:       []string{"supplier:update"},
			req:         httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-2&status=blocked", nil),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-2", status: "blocked"},
			},
		},
		{
			name:        "on_hold persists on_hold status",
			perms:       []string{"supplier:update"},
			req:         httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-3&status=on_hold", nil),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-3", status: "on_hold"},
			},
		},
		{
			name:        "falls back to form values when query params are absent",
			perms:       []string{"supplier:update"},
			req:         newFormRequest(http.MethodPost, "/action/suppliers/status", map[string]string{"id": "sup-4", "status": "active"}),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-4", status: "active"},
			},
		},
		{
			name:          "dependency error returns htmx error",
			perms:         []string{"supplier:update"},
			req:           httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-5&status=active", nil),
			setErr:        errors.New("set active failed"),
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrHeader: "set active failed",
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-5", status: "active"},
			},
		},
		{
			name:            "nil dependency still succeeds after validation",
			perms:           []string{"supplier:update"},
			req:             httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-6&status=active", nil),
			wantStatus:      http.StatusOK,
			wantTrigger:     `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var gotCalls []supplierStatusCall
			deps := &Deps{}
			if tt.name != "nil dependency still succeeds after validation" {
				deps.SetSupplierStatus = func(_ context.Context, id string, status string) error {
					gotCalls = append(gotCalls, supplierStatusCall{id: id, status: status})
					return tt.setErr
				}
			}

			res := NewSetStatusAction(deps).Handle(
				ctxWithPerms(tt.perms...),
				newTestViewContext(tt.req),
			)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if got := res.Headers["HX-Error-Message"]; got != tt.wantErrHeader {
				t.Fatalf("HX-Error-Message = %q, want %q", got, tt.wantErrHeader)
			}
			if got := res.Headers["HX-Trigger"]; got != tt.wantTrigger {
				t.Fatalf("HX-Trigger = %q, want %q", got, tt.wantTrigger)
			}
			if !reflect.DeepEqual(gotCalls, tt.wantStatusCalls) {
				t.Fatalf("SetSupplierStatus calls = %#v, want %#v", gotCalls, tt.wantStatusCalls)
			}
		})
	}
}

func TestNewBulkSetStatusAction_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perms           []string
		req             *http.Request
		setErrByID      map[string]error
		useNilDep       bool
		wantStatus      int
		wantErrHeader   string
		wantTrigger     string
		wantStatusCalls []supplierStatusCall
	}{
		{
			name:          "permission denied when user lacks supplier:update",
			perms:         []string{},
			req:           newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-1"}, "target_status": {"active"}}),
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrHeader: "permission denied",
		},
		{
			name:          "no ids provided",
			perms:         []string{"supplier:update"},
			req:           newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"target_status": {"active"}}),
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrHeader: "no ids provided",
		},
		{
			name:          "invalid target status",
			perms:         []string{"supplier:update"},
			req:           newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-1"}, "target_status": {"paused"}}),
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrHeader: "invalid target status",
		},
		{
			name:        "active persists active for all ids",
			perms:       []string{"supplier:update"},
			req:         newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-1", "sup-2"}, "target_status": {"active"}}),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-1", status: "active"},
				{id: "sup-2", status: "active"},
			},
		},
		{
			name:        "blocked persists blocked for all ids",
			perms:       []string{"supplier:update"},
			req:         newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-3", "sup-4"}, "target_status": {"blocked"}}),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-3", status: "blocked"},
				{id: "sup-4", status: "blocked"},
			},
		},
		{
			name:        "on_hold persists on_hold for all ids",
			perms:       []string{"supplier:update"},
			req:         newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-5", "sup-6"}, "target_status": {"on_hold"}}),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-5", status: "on_hold"},
				{id: "sup-6", status: "on_hold"},
			},
		},
		{
			name:        "partial failures still return success",
			perms:       []string{"supplier:update"},
			req:         newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-7", "sup-8"}, "target_status": {"active"}}),
			setErrByID:  map[string]error{"sup-8": errors.New("write failed")},
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-7", status: "active"},
				{id: "sup-8", status: "active"},
			},
		},
		{
			name:            "nil dependency still succeeds after validation",
			perms:           []string{"supplier:update"},
			req:             newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-9"}, "target_status": {"active"}}),
			useNilDep:       true,
			wantStatus:      http.StatusOK,
			wantTrigger:     `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var gotCalls []supplierStatusCall
			deps := &Deps{}
			if !tt.useNilDep {
				deps.SetSupplierStatus = func(_ context.Context, id string, status string) error {
					gotCalls = append(gotCalls, supplierStatusCall{id: id, status: status})
					if tt.setErrByID != nil {
						if err, ok := tt.setErrByID[id]; ok {
							return err
						}
					}
					return nil
				}
			}

			res := NewBulkSetStatusAction(deps).Handle(
				ctxWithPerms(tt.perms...),
				newTestViewContext(tt.req),
			)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if got := res.Headers["HX-Error-Message"]; got != tt.wantErrHeader {
				t.Fatalf("HX-Error-Message = %q, want %q", got, tt.wantErrHeader)
			}
			if got := res.Headers["HX-Trigger"]; got != tt.wantTrigger {
				t.Fatalf("HX-Trigger = %q, want %q", got, tt.wantTrigger)
			}
			if !reflect.DeepEqual(gotCalls, tt.wantStatusCalls) {
				t.Fatalf("SetSupplierStatus calls = %#v, want %#v", gotCalls, tt.wantStatusCalls)
			}
		})
	}
}

func newTestViewContext(req *http.Request) *view.ViewContext {
	return &view.ViewContext{
		Request: req,
		Messages: map[string]string{
			"shared.errors.permissionDenied":    "permission denied",
			"shared.errors.idRequired":          "id required",
			"shared.errors.invalidStatus":       "invalid status",
			"shared.errors.noIdsProvided":       "no ids provided",
			"shared.errors.invalidTargetStatus": "invalid target status",
		},
	}
}

func ctxWithPerms(codes ...string) context.Context {
	return view.WithUserPermissions(context.Background(), types.NewUserPermissions(codes))
}

func newFormRequest(method, target string, fields map[string]string) *http.Request {
	values := url.Values{}
	for key, value := range fields {
		values.Set(key, value)
	}

	req := httptest.NewRequest(method, target, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func newMultipartRequest(method, target string, fields map[string][]string) *http.Request {
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	for key, values := range fields {
		for _, value := range values {
			_ = writer.WriteField(key, value)
		}
	}
	_ = writer.Close()

	req := httptest.NewRequest(method, target, strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

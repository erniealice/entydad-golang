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
	active bool
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
			name:        "active maps to true",
			perms:       []string{"supplier:update"},
			req:         httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-1&status=active", nil),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-1", active: true},
			},
		},
		{
			name:        "blocked maps to false",
			perms:       []string{"supplier:update"},
			req:         httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-2&status=blocked", nil),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-2", active: false},
			},
		},
		{
			name:        "on_hold maps to false",
			perms:       []string{"supplier:update"},
			req:         httptest.NewRequest(http.MethodPost, "/action/suppliers/status?id=sup-3&status=on_hold", nil),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-3", active: false},
			},
		},
		{
			name:        "falls back to form values when query params are absent",
			perms:       []string{"supplier:update"},
			req:         newFormRequest(http.MethodPost, "/action/suppliers/status", map[string]string{"id": "sup-4", "status": "active"}),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-4", active: true},
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
				{id: "sup-5", active: true},
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
				deps.SetSupplierActive = func(_ context.Context, id string, active bool) error {
					gotCalls = append(gotCalls, supplierStatusCall{id: id, active: active})
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
				t.Fatalf("SetSupplierActive calls = %#v, want %#v", gotCalls, tt.wantStatusCalls)
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
			name:        "active maps all ids to true",
			perms:       []string{"supplier:update"},
			req:         newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-1", "sup-2"}, "target_status": {"active"}}),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-1", active: true},
				{id: "sup-2", active: true},
			},
		},
		{
			name:        "blocked maps all ids to false",
			perms:       []string{"supplier:update"},
			req:         newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-3", "sup-4"}, "target_status": {"blocked"}}),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-3", active: false},
				{id: "sup-4", active: false},
			},
		},
		{
			name:        "on_hold maps all ids to false",
			perms:       []string{"supplier:update"},
			req:         newMultipartRequest(http.MethodPost, "/action/suppliers/bulk-status", map[string][]string{"id": {"sup-5", "sup-6"}, "target_status": {"on_hold"}}),
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"suppliers-table"}`,
			wantStatusCalls: []supplierStatusCall{
				{id: "sup-5", active: false},
				{id: "sup-6", active: false},
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
				{id: "sup-7", active: true},
				{id: "sup-8", active: true},
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
				deps.SetSupplierActive = func(_ context.Context, id string, active bool) error {
					gotCalls = append(gotCalls, supplierStatusCall{id: id, active: active})
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
				t.Fatalf("SetSupplierActive calls = %#v, want %#v", gotCalls, tt.wantStatusCalls)
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

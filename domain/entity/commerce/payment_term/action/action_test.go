package action

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	paymentterm "github.com/erniealice/entydad-golang/domain/entity/commerce/payment_term"
)

type deleteCall struct {
	id string
}

type createPTCall struct {
	name               string
	code               string
	ptType             string
	netDays            int32
	discountDays       *int32
	discountPercentBps *int32
	entityScope        string
	isDefault          bool
	active             bool
}

type updatePTCall struct {
	id                 string
	name               string
	code               string
	ptType             string
	netDays            int32
	discountDays       *int32
	discountPercentBps *int32
	entityScope        string
	isDefault          bool
	active             bool
}

type ptActionRecorder struct {
	deleteCalls   []deleteCall
	createCalls   []createPTCall
	updateCalls   []updatePTCall
	deleteErrByID map[string]error
	createErr     error
	updateErr     error
}

func (r *ptActionRecorder) deletePaymentTerm(_ context.Context, req *paymenttermpb.DeletePaymentTermRequest) (*paymenttermpb.DeletePaymentTermResponse, error) {
	id := req.GetData().GetId()
	r.deleteCalls = append(r.deleteCalls, deleteCall{id: id})
	if err := r.deleteErrByID[id]; err != nil {
		return nil, err
	}
	return &paymenttermpb.DeletePaymentTermResponse{}, nil
}

func (r *ptActionRecorder) createPaymentTerm(_ context.Context, req *paymenttermpb.CreatePaymentTermRequest) (*paymenttermpb.CreatePaymentTermResponse, error) {
	d := req.GetData()
	r.createCalls = append(r.createCalls, createPTCall{
		name:               d.GetName(),
		code:               d.GetCode(),
		ptType:             d.GetType(),
		netDays:            d.GetNetDays(),
		discountDays:       d.DiscountDays,
		discountPercentBps: d.DiscountPercentBps,
		entityScope:        d.GetEntityScope(),
		isDefault:          d.GetIsDefault(),
		active:             d.GetActive(),
	})
	if r.createErr != nil {
		return nil, r.createErr
	}
	return &paymenttermpb.CreatePaymentTermResponse{}, nil
}

func (r *ptActionRecorder) updatePaymentTerm(_ context.Context, req *paymenttermpb.UpdatePaymentTermRequest) (*paymenttermpb.UpdatePaymentTermResponse, error) {
	d := req.GetData()
	r.updateCalls = append(r.updateCalls, updatePTCall{
		id:                 d.GetId(),
		name:               d.GetName(),
		code:               d.GetCode(),
		ptType:             d.GetType(),
		netDays:            d.GetNetDays(),
		discountDays:       d.DiscountDays,
		discountPercentBps: d.DiscountPercentBps,
		entityScope:        d.GetEntityScope(),
		isDefault:          d.GetIsDefault(),
		active:             d.GetActive(),
	})
	if r.updateErr != nil {
		return nil, r.updateErr
	}
	return &paymenttermpb.UpdatePaymentTermResponse{}, nil
}

func (r *ptActionRecorder) readPaymentTerm(_ context.Context, req *paymenttermpb.ReadPaymentTermRequest) (*paymenttermpb.ReadPaymentTermResponse, error) {
	id := req.GetData().GetId()
	nd := int32(30)
	dd := int32(10)
	dbps := int32(200)
	do := int32(1)
	desc := "test desc"
	return &paymenttermpb.ReadPaymentTermResponse{
		Data: []*paymenttermpb.PaymentTerm{{
			Id:                 id,
			Name:               "Net 30",
			Code:               "NET30",
			Type:               "net",
			NetDays:            nd,
			DiscountDays:       &dd,
			DiscountPercentBps: &dbps,
			EntityScope:        "both",
			IsDefault:          false,
			Description:        &desc,
			DisplayOrder:       &do,
			Active:             true,
		}},
	}, nil
}

func testMessages() map[string]string {
	return map[string]string{
		"shared.errors.permissionDenied": "permission denied",
		"shared.errors.idRequired":       "id required",
		"shared.errors.noIdsProvided":    "no ids provided",
		"shared.errors.invalidFormData":  "invalid form data",
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

func int32Ptr(v int32) *int32 {
	return &v
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
		wantStatus      int
		wantErrorHeader string
		wantDeleteIDs   []string
	}{
		{
			name:            "permission denied",
			perms:           []string{"payment_term:read"},
			req:             makePostRequest("/action/payment-terms/delete?id=pt-1", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "missing id",
			perms:           []string{"payment_term:delete"},
			req:             makePostRequest("/action/payment-terms/delete", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:          "uses id from query and succeeds",
			perms:         []string{"payment_term:delete"},
			req:           makePostRequest("/action/payment-terms/delete?id=pt-q", nil),
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"pt-q"},
		},
		{
			name:          "falls back to form id and succeeds",
			perms:         []string{"payment_term:delete"},
			req:           makePostRequest("/action/payment-terms/delete", url.Values{"id": {"pt-f"}}),
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"pt-f"},
		},
		{
			name:            "dependency error returns htmx error",
			perms:           []string{"payment_term:delete"},
			req:             makePostRequest("/action/payment-terms/delete?id=pt-e", nil),
			deleteErrByID:   map[string]error{"pt-e": errors.New("delete failed")},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "delete failed",
			wantDeleteIDs:   []string{"pt-e"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &ptActionRecorder{deleteErrByID: tt.deleteErrByID}
			deps := &Deps{DeletePaymentTerm: rec.deletePaymentTerm}
			res := runHandler(t, NewDeleteAction(deps), withPerms(tt.perms...), tt.req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "payment-terms-table")
			}

			gotIDs := make([]string, 0, len(rec.deleteCalls))
			for _, c := range rec.deleteCalls {
				gotIDs = append(gotIDs, c.id)
			}
			if strings.Join(gotIDs, ",") != strings.Join(tt.wantDeleteIDs, ",") {
				t.Fatalf("DeletePaymentTerm IDs = %v, want %v", gotIDs, tt.wantDeleteIDs)
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
		wantStatus      int
		wantErrorHeader string
		wantDeleteIDs   []string
	}{
		{
			name:            "permission denied",
			perms:           []string{"payment_term:read"},
			form:            url.Values{"id": {"pt-1"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "no ids provided",
			perms:           []string{"payment_term:delete"},
			form:            url.Values{},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "no ids provided",
		},
		{
			name:          "success deletes all ids",
			perms:         []string{"payment_term:delete"},
			form:          url.Values{"id": {"pt-1", "pt-2", "pt-3"}},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"pt-1", "pt-2", "pt-3"},
		},
		{
			name:          "partial failure still returns success",
			perms:         []string{"payment_term:delete"},
			form:          url.Values{"id": {"pt-1", "pt-2"}},
			deleteErrByID: map[string]error{"pt-2": errors.New("boom")},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"pt-1", "pt-2"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &ptActionRecorder{deleteErrByID: tt.deleteErrByID}
			deps := &Deps{DeletePaymentTerm: rec.deletePaymentTerm}
			req := makePostRequest("/action/payment-terms/bulk-delete", tt.form)
			res := runHandler(t, NewBulkDeleteAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "payment-terms-table")
			}

			gotIDs := make([]string, 0, len(rec.deleteCalls))
			for _, c := range rec.deleteCalls {
				gotIDs = append(gotIDs, c.id)
			}
			if strings.Join(gotIDs, ",") != strings.Join(tt.wantDeleteIDs, ",") {
				t.Fatalf("DeletePaymentTerm IDs = %v, want %v", gotIDs, tt.wantDeleteIDs)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AddAction (POST path) — numeric field parsing
// ---------------------------------------------------------------------------

func TestNewAddAction_POST_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		perms                  []string
		form                   url.Values
		createErr              error
		wantStatus             int
		wantErrorHeader        string
		wantCreateCount        int
		wantNetDays            int32
		wantDiscountDays       *int32
		wantDiscountPercentBps *int32
	}{
		{
			name:            "permission denied",
			perms:           []string{"payment_term:read"},
			form:            url.Values{"name": {"Net 30"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:  "success with all numeric fields",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":                 {"Net 30"},
				"code":                 {"NET30"},
				"type":                 {"net"},
				"net_days":             {"30"},
				"discount_days":        {"10"},
				"discount_percent_bps": {"200"},
				"entity_scope":         {"both"},
				"is_default":           {"true"},
				"active":               {"true"},
			},
			wantStatus:             http.StatusOK,
			wantCreateCount:        1,
			wantNetDays:            30,
			wantDiscountDays:       int32Ptr(10),
			wantDiscountPercentBps: int32Ptr(200),
		},
		{
			name:  "empty optional numeric fields yield nil",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":         {"Due on Receipt"},
				"code":         {"DOR"},
				"type":         {"due_on_receipt"},
				"net_days":     {"0"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			wantStatus:             http.StatusOK,
			wantCreateCount:        1,
			wantNetDays:            0,
			wantDiscountDays:       nil,
			wantDiscountPercentBps: nil,
		},
		{
			name:  "invalid net_days defaults to 0",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":         {"Bad"},
				"code":         {"BAD"},
				"type":         {"net"},
				"net_days":     {"abc"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantNetDays:     0,
		},
		{
			name:  "invalid discount_days yields nil",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":          {"Bad Disc"},
				"code":          {"BDISC"},
				"type":          {"net"},
				"net_days":      {"30"},
				"discount_days": {"xyz"},
				"entity_scope":  {"both"},
				"active":        {"true"},
			},
			wantStatus:       http.StatusOK,
			wantCreateCount:  1,
			wantNetDays:      30,
			wantDiscountDays: nil,
		},
		{
			name:  "invalid discount_percent_bps yields nil",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":                 {"Bad BPS"},
				"code":                 {"BBPS"},
				"type":                 {"net"},
				"net_days":             {"30"},
				"discount_percent_bps": {"not-a-number"},
				"entity_scope":         {"both"},
				"active":               {"true"},
			},
			wantStatus:             http.StatusOK,
			wantCreateCount:        1,
			wantNetDays:            30,
			wantDiscountPercentBps: nil,
		},
		{
			name:  "create error",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":         {"Fail"},
				"code":         {"FAIL"},
				"type":         {"net"},
				"net_days":     {"30"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			createErr:       errors.New("create failed"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "create failed",
			wantCreateCount: 1,
			wantNetDays:     30,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &ptActionRecorder{createErr: tt.createErr}
			deps := &Deps{
				CreatePaymentTerm: rec.createPaymentTerm,
				Routes:            paymentterm.Routes{AddURL: "/action/payment-terms/add"},
			}
			req := makePostRequest("/action/payment-terms/add", tt.form)
			res := runHandler(t, NewAddAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "payment-terms-table")
			}
			if len(rec.createCalls) != tt.wantCreateCount {
				t.Fatalf("CreatePaymentTerm call count = %d, want %d", len(rec.createCalls), tt.wantCreateCount)
			}
			if tt.wantCreateCount > 0 && len(rec.createCalls) > 0 {
				got := rec.createCalls[0]
				if got.netDays != tt.wantNetDays {
					t.Fatalf("netDays = %d, want %d", got.netDays, tt.wantNetDays)
				}
				if tt.wantDiscountDays != nil {
					if got.discountDays == nil {
						t.Fatalf("discountDays = nil, want %d", *tt.wantDiscountDays)
					} else if *got.discountDays != *tt.wantDiscountDays {
						t.Fatalf("discountDays = %d, want %d", *got.discountDays, *tt.wantDiscountDays)
					}
				} else if got.discountDays != nil {
					t.Fatalf("discountDays = %d, want nil", *got.discountDays)
				}
				if tt.wantDiscountPercentBps != nil {
					if got.discountPercentBps == nil {
						t.Fatalf("discountPercentBps = nil, want %d", *tt.wantDiscountPercentBps)
					} else if *got.discountPercentBps != *tt.wantDiscountPercentBps {
						t.Fatalf("discountPercentBps = %d, want %d", *got.discountPercentBps, *tt.wantDiscountPercentBps)
					}
				} else if got.discountPercentBps != nil {
					t.Fatalf("discountPercentBps = %d, want nil", *got.discountPercentBps)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// EditAction (POST path) — numeric field parsing
// ---------------------------------------------------------------------------

func TestNewEditAction_POST_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		perms                  []string
		form                   url.Values
		updateErr              error
		wantStatus             int
		wantErrorHeader        string
		wantUpdateCount        int
		wantNetDays            int32
		wantDiscountDays       *int32
		wantDiscountPercentBps *int32
	}{
		{
			name:            "permission denied",
			perms:           []string{"payment_term:read"},
			form:            url.Values{"name": {"Net 30"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:  "success with numeric fields",
			perms: []string{"payment_term:update"},
			form: url.Values{
				"name":                 {"Net 60"},
				"code":                 {"NET60"},
				"type":                 {"net"},
				"net_days":             {"60"},
				"discount_days":        {"15"},
				"discount_percent_bps": {"100"},
				"entity_scope":         {"client"},
				"is_default":           {"false"},
				"active":               {"true"},
			},
			wantStatus:             http.StatusOK,
			wantUpdateCount:        1,
			wantNetDays:            60,
			wantDiscountDays:       int32Ptr(15),
			wantDiscountPercentBps: int32Ptr(100),
		},
		{
			name:  "update error",
			perms: []string{"payment_term:update"},
			form: url.Values{
				"name":         {"Fail"},
				"code":         {"FAIL"},
				"type":         {"net"},
				"net_days":     {"30"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			updateErr:       errors.New("update failed"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "update failed",
			wantUpdateCount: 1,
			wantNetDays:     30,
		},
		{
			name:  "invalid discount_days in edit yields nil",
			perms: []string{"payment_term:update"},
			form: url.Values{
				"name":          {"Edge"},
				"code":          {"EDGE"},
				"type":          {"net"},
				"net_days":      {"30"},
				"discount_days": {"not-numeric"},
				"entity_scope":  {"both"},
				"active":        {"true"},
			},
			wantStatus:       http.StatusOK,
			wantUpdateCount:  1,
			wantNetDays:      30,
			wantDiscountDays: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &ptActionRecorder{updateErr: tt.updateErr}
			deps := &Deps{
				UpdatePaymentTerm: rec.updatePaymentTerm,
				ReadPaymentTerm:   rec.readPaymentTerm,
				Routes:            paymentterm.Routes{EditURL: "/action/payment-terms/{id}/edit"},
			}
			req := makePostRequest("/action/payment-terms/pt-1/edit", tt.form)
			res := runHandler(t, NewEditAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "payment-terms-table")
			}
			if len(rec.updateCalls) != tt.wantUpdateCount {
				t.Fatalf("UpdatePaymentTerm call count = %d, want %d", len(rec.updateCalls), tt.wantUpdateCount)
			}
			if tt.wantUpdateCount > 0 && len(rec.updateCalls) > 0 {
				got := rec.updateCalls[0]
				if got.netDays != tt.wantNetDays {
					t.Fatalf("netDays = %d, want %d", got.netDays, tt.wantNetDays)
				}
				if tt.wantDiscountDays != nil {
					if got.discountDays == nil {
						t.Fatalf("discountDays = nil, want %d", *tt.wantDiscountDays)
					} else if *got.discountDays != *tt.wantDiscountDays {
						t.Fatalf("discountDays = %d, want %d", *got.discountDays, *tt.wantDiscountDays)
					}
				} else if got.discountDays != nil {
					t.Fatalf("discountDays = %d, want nil", *got.discountDays)
				}
				if tt.wantDiscountPercentBps != nil {
					if got.discountPercentBps == nil {
						t.Fatalf("discountPercentBps = nil, want %d", *tt.wantDiscountPercentBps)
					} else if *got.discountPercentBps != *tt.wantDiscountPercentBps {
						t.Fatalf("discountPercentBps = %d, want %d", *got.discountPercentBps, *tt.wantDiscountPercentBps)
					}
				} else if got.discountPercentBps != nil {
					t.Fatalf("discountPercentBps = %d, want nil", *got.discountPercentBps)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AddAction — negative / defensive tests (boundary values, adversarial input)
// ---------------------------------------------------------------------------

func TestNewAddAction_Negative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		perms                  []string
		form                   url.Values
		createErr              error
		wantStatus             int
		wantErrorHeader        string
		wantCreateCount        int
		wantNetDays            int32
		wantDiscountDays       *int32
		wantDiscountPercentBps *int32
	}{
		{
			name:  "negative net_days parses as negative int32",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":         {"Negative Net"},
				"code":         {"NEGNET"},
				"type":         {"net"},
				"net_days":     {"-10"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantNetDays:     -10,
		},
		{
			name:  "negative discount_days parses as negative int32",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":          {"Neg Disc"},
				"code":          {"NEGDISC"},
				"type":          {"net"},
				"net_days":      {"30"},
				"discount_days": {"-5"},
				"entity_scope":  {"both"},
				"active":        {"true"},
			},
			wantStatus:       http.StatusOK,
			wantCreateCount:  1,
			wantNetDays:      30,
			wantDiscountDays: int32Ptr(-5),
		},
		{
			name:  "discount_percent_bps over 10000 (>100%) is accepted",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":                 {"Over 100%"},
				"code":                 {"OVER"},
				"type":                 {"net"},
				"net_days":             {"30"},
				"discount_percent_bps": {"15000"},
				"entity_scope":         {"both"},
				"active":               {"true"},
			},
			wantStatus:             http.StatusOK,
			wantCreateCount:        1,
			wantNetDays:            30,
			wantDiscountPercentBps: int32Ptr(15000),
		},
		{
			name:  "floating point net_days is treated as invalid (defaults to 0)",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":         {"Float"},
				"code":         {"FLT"},
				"type":         {"net"},
				"net_days":     {"10.5"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantNetDays:     0,
		},
		{
			name:  "very large net_days within int32 range",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":         {"MaxInt"},
				"code":         {"MAX"},
				"type":         {"net"},
				"net_days":     {"2147483647"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantNetDays:     2147483647,
		},
		{
			name:  "net_days exceeding int32 range defaults to 0",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":         {"Overflow"},
				"code":         {"OVF"},
				"type":         {"net"},
				"net_days":     {"99999999999"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantNetDays:     0,
		},
		{
			name:  "missing all fields still calls create",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
		},
		{
			name:  "backend rejects missing name",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"active": {"true"},
			},
			createErr:       errors.New("name is required"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "name is required",
			wantCreateCount: 1,
		},
		{
			name:  "XSS in name field passes through",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":         {"<script>alert('xss')</script>"},
				"code":         {"XSS"},
				"type":         {"net"},
				"net_days":     {"30"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantNetDays:     30,
		},
		{
			name:  "special characters in discount fields yield nil",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":                 {"Special"},
				"code":                 {"SPEC"},
				"type":                 {"net"},
				"net_days":             {"30"},
				"discount_days":        {"!@#"},
				"discount_percent_bps": {"$%^"},
				"entity_scope":         {"both"},
				"active":               {"true"},
			},
			wantStatus:             http.StatusOK,
			wantCreateCount:        1,
			wantNetDays:            30,
			wantDiscountDays:       nil,
			wantDiscountPercentBps: nil,
		},
		{
			name:  "zero discount_percent_bps parses as zero pointer",
			perms: []string{"payment_term:create"},
			form: url.Values{
				"name":                 {"Zero BPS"},
				"code":                 {"ZBPS"},
				"type":                 {"net"},
				"net_days":             {"30"},
				"discount_percent_bps": {"0"},
				"entity_scope":         {"both"},
				"active":               {"true"},
			},
			wantStatus:             http.StatusOK,
			wantCreateCount:        1,
			wantNetDays:            30,
			wantDiscountPercentBps: int32Ptr(0),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &ptActionRecorder{createErr: tt.createErr}
			deps := &Deps{
				CreatePaymentTerm: rec.createPaymentTerm,
				Routes:            paymentterm.Routes{AddURL: "/action/payment-terms/add"},
			}
			req := makePostRequest("/action/payment-terms/add", tt.form)
			res := runHandler(t, NewAddAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "payment-terms-table")
			}
			if len(rec.createCalls) != tt.wantCreateCount {
				t.Fatalf("CreatePaymentTerm call count = %d, want %d", len(rec.createCalls), tt.wantCreateCount)
			}
			if tt.wantCreateCount > 0 && len(rec.createCalls) > 0 {
				got := rec.createCalls[0]
				if got.netDays != tt.wantNetDays {
					t.Fatalf("netDays = %d, want %d", got.netDays, tt.wantNetDays)
				}
				if tt.wantDiscountDays != nil {
					if got.discountDays == nil {
						t.Fatalf("discountDays = nil, want %d", *tt.wantDiscountDays)
					} else if *got.discountDays != *tt.wantDiscountDays {
						t.Fatalf("discountDays = %d, want %d", *got.discountDays, *tt.wantDiscountDays)
					}
				} else if tt.wantDiscountDays == nil && got.discountDays != nil {
					t.Fatalf("discountDays = %d, want nil", *got.discountDays)
				}
				if tt.wantDiscountPercentBps != nil {
					if got.discountPercentBps == nil {
						t.Fatalf("discountPercentBps = nil, want %d", *tt.wantDiscountPercentBps)
					} else if *got.discountPercentBps != *tt.wantDiscountPercentBps {
						t.Fatalf("discountPercentBps = %d, want %d", *got.discountPercentBps, *tt.wantDiscountPercentBps)
					}
				} else if tt.wantDiscountPercentBps == nil && got.discountPercentBps != nil {
					t.Fatalf("discountPercentBps = %d, want nil", *got.discountPercentBps)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// EditAction — negative / defensive tests (missing ID)
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
			perms:   []string{"payment_term:update"},
			pathURL: "/action/payment-terms//edit",
			form: url.Values{
				"name":         {"Updated"},
				"code":         {"UPD"},
				"type":         {"net"},
				"net_days":     {"30"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			updateErr:       errors.New("id is required"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id is required",
			wantUpdateCount: 1,
		},
		{
			name:    "edit with negative net_days",
			perms:   []string{"payment_term:update"},
			pathURL: "/action/payment-terms/pt-1/edit",
			form: url.Values{
				"name":         {"Negative"},
				"code":         {"NEG"},
				"type":         {"net"},
				"net_days":     {"-1"},
				"entity_scope": {"both"},
				"active":       {"true"},
			},
			wantStatus:      http.StatusOK,
			wantUpdateCount: 1,
		},
		{
			name:    "edit with overflow discount_percent_bps",
			perms:   []string{"payment_term:update"},
			pathURL: "/action/payment-terms/pt-1/edit",
			form: url.Values{
				"name":                 {"Overflow BPS"},
				"code":                 {"OVBPS"},
				"type":                 {"net"},
				"net_days":             {"30"},
				"discount_percent_bps": {"99999"},
				"entity_scope":         {"both"},
				"active":               {"true"},
			},
			wantStatus:      http.StatusOK,
			wantUpdateCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &ptActionRecorder{updateErr: tt.updateErr}
			deps := &Deps{
				UpdatePaymentTerm: rec.updatePaymentTerm,
				ReadPaymentTerm:   rec.readPaymentTerm,
				Routes:            paymentterm.Routes{EditURL: "/action/payment-terms/{id}/edit"},
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
				assertSuccessHeader(t, res, "payment-terms-table")
			}
			if len(rec.updateCalls) != tt.wantUpdateCount {
				t.Fatalf("UpdatePaymentTerm call count = %d, want %d", len(rec.updateCalls), tt.wantUpdateCount)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Delete — negative / defensive tests (empty ID)
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
			perms:           []string{"payment_term:delete"},
			req:             makePostRequest("/action/payment-terms/delete?id=", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:            "no id in query or form",
			perms:           []string{"payment_term:delete"},
			req:             makePostRequest("/action/payment-terms/delete", url.Values{}),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:       "whitespace-only id in form",
			perms:      []string{"payment_term:delete"},
			req:        makePostRequest("/action/payment-terms/delete", url.Values{"id": {"   "}}),
			wantStatus: http.StatusOK, // whitespace is not validated, treated as a valid ID
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &ptActionRecorder{}
			deps := &Deps{DeletePaymentTerm: rec.deletePaymentTerm}
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

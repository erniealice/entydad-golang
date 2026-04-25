package action

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
)

type deleteCall struct {
	id string
}

type statusCall struct {
	id     string
	status string
}

type createClientCall struct {
	name      string
	firstName string
	lastName  string
	email     string
	active    bool
}

type updateClientCall struct {
	id        string
	name      string
	firstName string
	lastName  string
	email     string
	active    bool
}

type createClientCategoryCall struct {
	clientID   string
	categoryID string
}

type deleteClientCategoryCall struct {
	id string
}

type clientActionRecorder struct {
	deleteCalls        []deleteCall
	statusCalls        []statusCall
	createCalls        []createClientCall
	updateCalls        []updateClientCall
	createCatCalls     []createClientCategoryCall
	deleteCatCalls     []deleteClientCategoryCall
	deleteErrByID      map[string]error
	statusErrByID      map[string]error
	createErr          error
	updateErr          error
	readErr            error
	createdClientID    string
	listCategoriesResp *categorypb.ListCategoriesResponse
	listClientCatsResp *clientcategorypb.ListClientCategoriesResponse
	createCatErr       error
	deleteCatErr       error
}

func (r *clientActionRecorder) deleteClient(_ context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error) {
	id := req.GetData().GetId()
	r.deleteCalls = append(r.deleteCalls, deleteCall{id: id})
	if err := r.deleteErrByID[id]; err != nil {
		return nil, err
	}
	return &clientpb.DeleteClientResponse{}, nil
}

func (r *clientActionRecorder) setClientStatus(_ context.Context, id string, status string) error {
	r.statusCalls = append(r.statusCalls, statusCall{id: id, status: status})
	return r.statusErrByID[id]
}

func (r *clientActionRecorder) createClient(_ context.Context, req *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error) {
	d := req.GetData()
	u := d.GetUser()
	r.createCalls = append(r.createCalls, createClientCall{
		name:      d.GetName(),
		firstName: u.GetFirstName(),
		lastName:  u.GetLastName(),
		email:     u.GetEmailAddress(),
		active:    d.GetActive(),
	})
	if r.createErr != nil {
		return nil, r.createErr
	}
	cid := r.createdClientID
	if cid == "" {
		cid = "new-client-id"
	}
	return &clientpb.CreateClientResponse{
		Data: []*clientpb.Client{{Id: cid}},
	}, nil
}

func (r *clientActionRecorder) readClient(_ context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error) {
	if r.readErr != nil {
		return nil, r.readErr
	}
	id := req.GetData().GetId()
	return &clientpb.ReadClientResponse{
		Data: []*clientpb.Client{{
			Id:     id,
			Active: true,
			User: &userpb.User{
				FirstName:    "Test",
				LastName:     "Client",
				EmailAddress: "client@test.com",
				MobileNumber: "+639000000000",
			},
		}},
	}, nil
}

func (r *clientActionRecorder) updateClient(_ context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error) {
	d := req.GetData()
	u := d.GetUser()
	r.updateCalls = append(r.updateCalls, updateClientCall{
		id:        d.GetId(),
		name:      d.GetName(),
		firstName: u.GetFirstName(),
		lastName:  u.GetLastName(),
		email:     u.GetEmailAddress(),
		active:    d.GetActive(),
	})
	if r.updateErr != nil {
		return nil, r.updateErr
	}
	return &clientpb.UpdateClientResponse{}, nil
}

func (r *clientActionRecorder) listCategories(_ context.Context, _ *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error) {
	if r.listCategoriesResp != nil {
		return r.listCategoriesResp, nil
	}
	return &categorypb.ListCategoriesResponse{}, nil
}

func (r *clientActionRecorder) listClientCategories(_ context.Context, _ *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error) {
	if r.listClientCatsResp != nil {
		return r.listClientCatsResp, nil
	}
	return &clientcategorypb.ListClientCategoriesResponse{}, nil
}

func (r *clientActionRecorder) createClientCategory(_ context.Context, req *clientcategorypb.CreateClientCategoryRequest) (*clientcategorypb.CreateClientCategoryResponse, error) {
	d := req.GetData()
	r.createCatCalls = append(r.createCatCalls, createClientCategoryCall{
		clientID:   d.GetClientId(),
		categoryID: d.GetCategoryId(),
	})
	if r.createCatErr != nil {
		return nil, r.createCatErr
	}
	return &clientcategorypb.CreateClientCategoryResponse{}, nil
}

func (r *clientActionRecorder) deleteClientCategory(_ context.Context, req *clientcategorypb.DeleteClientCategoryRequest) (*clientcategorypb.DeleteClientCategoryResponse, error) {
	d := req.GetData()
	r.deleteCatCalls = append(r.deleteCatCalls, deleteClientCategoryCall{id: d.GetId()})
	if r.deleteCatErr != nil {
		return nil, r.deleteCatErr
	}
	return &clientcategorypb.DeleteClientCategoryResponse{}, nil
}

func testMessages() map[string]string {
	return map[string]string{
		"shared.errors.permissionDenied":    "permission denied",
		"shared.errors.idRequired":          "id required",
		"shared.errors.noIdsProvided":       "no ids provided",
		"shared.errors.invalidStatus":       "invalid status",
		"shared.errors.invalidTargetStatus": "invalid target status",
		"shared.errors.invalidFormData":     "invalid form data",
		"shared.errors.notFound":            "not found",
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
			wantCalls:  []statusCall{{id: "cl-1", status: "active"}},
		},
		{
			name:       "form fallback inactive success",
			perms:      []string{"client:update"},
			req:        makePostRequest("/action/clients/set-status", url.Values{"id": {"cl-2"}, "status": {"inactive"}}),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "cl-2", status: "inactive"}},
		},
		{
			name:       "block lifecycle success",
			perms:      []string{"client:update"},
			req:        makePostRequest("/action/clients/set-status?id=cl-4&status=blocked", nil),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "cl-4", status: "blocked"}},
		},
		{
			name:       "on_hold lifecycle success",
			perms:      []string{"client:update"},
			req:        makePostRequest("/action/clients/set-status?id=cl-5&status=on_hold", nil),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "cl-5", status: "on_hold"}},
		},
		{
			name:       "prospect lifecycle success",
			perms:      []string{"client:update"},
			req:        makePostRequest("/action/clients/set-status?id=cl-6&status=prospect", nil),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "cl-6", status: "prospect"}},
		},
		{
			name:            "dependency error",
			perms:           []string{"client:update"},
			req:             makePostRequest("/action/clients/set-status?id=cl-3&status=inactive", nil),
			statusErrByID:   map[string]error{"cl-3": errors.New("set status failed")},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "set status failed",
			wantCalls:       []statusCall{{id: "cl-3", status: "inactive"}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &clientActionRecorder{statusErrByID: tt.statusErrByID}
			deps := &Deps{SetClientStatus: rec.setClientStatus}
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
				t.Fatalf("SetClientStatus call count = %d, want %d", len(rec.statusCalls), len(tt.wantCalls))
			}
			for i := range tt.wantCalls {
				if rec.statusCalls[i] != tt.wantCalls[i] {
					t.Fatalf("SetClientStatus call[%d] = %+v, want %+v", i, rec.statusCalls[i], tt.wantCalls[i])
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
				{id: "cl-1", status: "active"},
				{id: "cl-2", status: "active"},
			},
		},
		{
			name:       "inactive success",
			perms:      []string{"client:update"},
			form:       url.Values{"id": {"cl-1", "cl-2"}, "target_status": {"inactive"}},
			wantStatus: http.StatusOK,
			wantCalls: []statusCall{
				{id: "cl-1", status: "inactive"},
				{id: "cl-2", status: "inactive"},
			},
		},
		{
			name:       "blocked success",
			perms:      []string{"client:update"},
			form:       url.Values{"id": {"cl-1", "cl-2"}, "target_status": {"blocked"}},
			wantStatus: http.StatusOK,
			wantCalls: []statusCall{
				{id: "cl-1", status: "blocked"},
				{id: "cl-2", status: "blocked"},
			},
		},
		{
			name:          "partial failure still returns success",
			perms:         []string{"client:update"},
			form:          url.Values{"id": {"cl-1", "cl-2", "cl-3"}, "target_status": {"inactive"}},
			statusErrByID: map[string]error{"cl-2": errors.New("boom")},
			wantStatus:    http.StatusOK,
			wantCalls: []statusCall{
				{id: "cl-1", status: "inactive"},
				{id: "cl-2", status: "inactive"},
				{id: "cl-3", status: "inactive"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &clientActionRecorder{statusErrByID: tt.statusErrByID}
			deps := &Deps{SetClientStatus: rec.setClientStatus}
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
				t.Fatalf("SetClientStatus call count = %d, want %d", len(rec.statusCalls), len(tt.wantCalls))
			}
			for i := range tt.wantCalls {
				if rec.statusCalls[i] != tt.wantCalls[i] {
					t.Fatalf("SetClientStatus call[%d] = %+v, want %+v", i, rec.statusCalls[i], tt.wantCalls[i])
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AddAction — negative / defensive tests
// ---------------------------------------------------------------------------

func TestNewAddAction_Negative(t *testing.T) {
	t.Parallel()

	longString := strings.Repeat("Z", 300)

	tests := []struct {
		name            string
		perms           []string
		form            url.Values
		createErr       error
		wantStatus      int
		wantErrorHeader string
		wantCreateCount int
		wantName        string
		wantFirstName   string
	}{
		{
			name:  "missing all fields still calls create",
			perms: []string{"client:create"},
			form: url.Values{
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
		},
		{
			name:  "missing name (company name)",
			perms: []string{"client:create"},
			form: url.Values{
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantName:        "",
		},
		{
			name:  "missing first_name and last_name",
			perms: []string{"client:create"},
			form: url.Values{
				"name":          {"ACME Corp"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantFirstName:   "",
		},
		{
			name:  "backend rejects missing required fields",
			perms: []string{"client:create"},
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
			perms: []string{"client:create"},
			form: url.Values{
				"name":          {"<script>alert('xss')</script>"},
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantName:        "<script>alert('xss')</script>",
		},
		{
			name:  "very long name (>255 chars)",
			perms: []string{"client:create"},
			form: url.Values{
				"name":          {longString},
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantName:        longString,
		},
		{
			name:  "invalid email format passes through",
			perms: []string{"client:create"},
			form: url.Values{
				"name":          {"ACME"},
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"not-an-email"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
		},
		{
			name:  "with invalid tag IDs (non-existent) — tags synced after create",
			perms: []string{"client:create"},
			form: url.Values{
				"name":          {"Tagged"},
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
				"tags":          {"invalid-tag-id-1,invalid-tag-id-2"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
		},
		{
			name:  "empty tags string results in no sync",
			perms: []string{"client:create"},
			form: url.Values{
				"name":          {"No Tags"},
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
				"tags":          {""},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &clientActionRecorder{
				createErr:       tt.createErr,
				createdClientID: "new-client-id",
			}
			deps := &Deps{
				CreateClient:         rec.createClient,
				ListClientCategories: rec.listClientCategories,
				CreateClientCategory: rec.createClientCategory,
				DeleteClientCategory: rec.deleteClientCategory,
				Routes:               entydad.ClientRoutes{AddURL: "/action/clients/add"},
			}
			req := makePostRequest("/action/clients/add", tt.form)
			res := runHandler(t, NewAddAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "clients-table")
			}
			if len(rec.createCalls) != tt.wantCreateCount {
				t.Fatalf("CreateClient call count = %d, want %d", len(rec.createCalls), tt.wantCreateCount)
			}
			if tt.wantName != "" && len(rec.createCalls) > 0 {
				if got := rec.createCalls[0].name; got != tt.wantName {
					t.Fatalf("name = %q, want %q", got, tt.wantName)
				}
			}
			if tt.wantFirstName != "" && len(rec.createCalls) > 0 {
				if got := rec.createCalls[0].firstName; got != tt.wantFirstName {
					t.Fatalf("firstName = %q, want %q", got, tt.wantFirstName)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// EditAction — negative / defensive tests (empty ID)
// ---------------------------------------------------------------------------

func TestNewEditAction_POST_Negative(t *testing.T) {
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
			perms:   []string{"client:update"},
			pathURL: "/action/clients//edit",
			form: url.Values{
				"name":          {"Updated"},
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			updateErr:       errors.New("id is required"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id is required",
			wantUpdateCount: 1,
		},
		{
			name:    "edit with XSS in name",
			perms:   []string{"client:update"},
			pathURL: "/action/clients/cl-1/edit",
			form: url.Values{
				"name":          {"<script>alert(1)</script>"},
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantUpdateCount: 1,
		},
		{
			name:    "edit with missing all optional and required fields",
			perms:   []string{"client:update"},
			pathURL: "/action/clients/cl-1/edit",
			form: url.Values{
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

			rec := &clientActionRecorder{updateErr: tt.updateErr}
			deps := &Deps{
				UpdateClient:         rec.updateClient,
				ReadClient:           rec.readClient,
				ListClientCategories: rec.listClientCategories,
				CreateClientCategory: rec.createClientCategory,
				DeleteClientCategory: rec.deleteClientCategory,
				Routes:               entydad.ClientRoutes{EditURL: "/action/clients/{id}/edit"},
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
				assertSuccessHeader(t, res, "clients-table")
			}
			if len(rec.updateCalls) != tt.wantUpdateCount {
				t.Fatalf("UpdateClient call count = %d, want %d", len(rec.updateCalls), tt.wantUpdateCount)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// syncTags — negative / defensive tests
// ---------------------------------------------------------------------------

func TestSyncTags_Negative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		clientID           string
		submittedTagIDs    []string
		existingCats       []*clientcategorypb.ClientCategory
		wantCreateCatCount int
		wantDeleteCatCount int
	}{
		{
			name:               "nil tag list results in no creates or deletes",
			clientID:           "cl-1",
			submittedTagIDs:    nil,
			wantCreateCatCount: 0,
			wantDeleteCatCount: 0,
		},
		{
			name:            "empty tag list deletes existing assignments",
			clientID:        "cl-1",
			submittedTagIDs: []string{},
			existingCats: []*clientcategorypb.ClientCategory{
				{Id: "cc-1", ClientId: "cl-1", CategoryId: "cat-1"},
			},
			wantCreateCatCount: 0,
			wantDeleteCatCount: 1,
		},
		{
			name:               "empty strings in tag list are ignored",
			clientID:           "cl-1",
			submittedTagIDs:    []string{"", "", ""},
			wantCreateCatCount: 0,
			wantDeleteCatCount: 0,
		},
		{
			name:            "new tags are created and old tags are deleted",
			clientID:        "cl-1",
			submittedTagIDs: []string{"cat-new"},
			existingCats: []*clientcategorypb.ClientCategory{
				{Id: "cc-1", ClientId: "cl-1", CategoryId: "cat-old"},
			},
			wantCreateCatCount: 1,
			wantDeleteCatCount: 1,
		},
		{
			name:               "duplicate tag IDs in submitted list result in single create",
			clientID:           "cl-1",
			submittedTagIDs:    []string{"cat-1", "cat-1", "cat-1"},
			wantCreateCatCount: 1,
			wantDeleteCatCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &clientActionRecorder{
				listClientCatsResp: &clientcategorypb.ListClientCategoriesResponse{
					Data: tt.existingCats,
				},
			}
			deps := &Deps{
				ListClientCategories: rec.listClientCategories,
				CreateClientCategory: rec.createClientCategory,
				DeleteClientCategory: rec.deleteClientCategory,
			}

			syncTags(context.Background(), deps, tt.clientID, tt.submittedTagIDs)

			if len(rec.createCatCalls) != tt.wantCreateCatCount {
				t.Fatalf("CreateClientCategory call count = %d, want %d", len(rec.createCatCalls), tt.wantCreateCatCount)
			}
			if len(rec.deleteCatCalls) != tt.wantDeleteCatCount {
				t.Fatalf("DeleteClientCategory call count = %d, want %d", len(rec.deleteCatCalls), tt.wantDeleteCatCount)
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
			perms:           []string{"client:delete"},
			req:             makePostRequest("/action/clients/delete?id=", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:       "whitespace-only id in form",
			perms:      []string{"client:delete"},
			req:        makePostRequest("/action/clients/delete", url.Values{"id": {"   "}}),
			wantStatus: http.StatusOK, // whitespace is not validated, treated as a valid ID
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &clientActionRecorder{}
			deps := &Deps{DeleteClient: rec.deleteClient}
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

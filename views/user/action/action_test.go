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

	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
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

type createUserCall struct {
	firstName    string
	lastName     string
	email        string
	mobile       string
	passwordHash string
	active       bool
}

type createWSUserCall struct {
	workspaceID string
	userID      string
}

type updateUserCall struct {
	id           string
	firstName    string
	lastName     string
	email        string
	mobile       string
	passwordHash string
	active       bool
}

type userActionRecorder struct {
	deleteCalls      []deleteCall
	statusCalls      []statusCall
	createUserCalls  []createUserCall
	createWSCalls    []createWSUserCall
	updateUserCalls  []updateUserCall
	deleteErrByID    map[string]error
	statusErrByID    map[string]error
	createUserErr    error
	createWSErr      error
	updateUserErr    error
	readUserErr      error
	createdUserID    string // returned as the new user ID
	hashPasswordFunc func(string) (string, error)
}

func (r *userActionRecorder) deleteUser(_ context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	id := req.GetData().GetId()
	r.deleteCalls = append(r.deleteCalls, deleteCall{id: id})
	if err := r.deleteErrByID[id]; err != nil {
		return nil, err
	}
	return &userpb.DeleteUserResponse{}, nil
}

func (r *userActionRecorder) setUserActive(_ context.Context, id string, active bool) error {
	r.statusCalls = append(r.statusCalls, statusCall{id: id, active: active})
	return r.statusErrByID[id]
}

func (r *userActionRecorder) createUser(_ context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	d := req.GetData()
	r.createUserCalls = append(r.createUserCalls, createUserCall{
		firstName:    d.GetFirstName(),
		lastName:     d.GetLastName(),
		email:        d.GetEmailAddress(),
		mobile:       d.GetMobileNumber(),
		passwordHash: d.GetPasswordHash(),
		active:       d.GetActive(),
	})
	if r.createUserErr != nil {
		return nil, r.createUserErr
	}
	uid := r.createdUserID
	if uid == "" {
		uid = "new-user-id"
	}
	return &userpb.CreateUserResponse{
		Data: []*userpb.User{{Id: uid}},
	}, nil
}

func (r *userActionRecorder) createWorkspaceUser(_ context.Context, req *workspaceuserpb.CreateWorkspaceUserRequest) (*workspaceuserpb.CreateWorkspaceUserResponse, error) {
	d := req.GetData()
	r.createWSCalls = append(r.createWSCalls, createWSUserCall{
		workspaceID: d.GetWorkspaceId(),
		userID:      d.GetUserId(),
	})
	if r.createWSErr != nil {
		return nil, r.createWSErr
	}
	return &workspaceuserpb.CreateWorkspaceUserResponse{}, nil
}

func (r *userActionRecorder) readUser(_ context.Context, req *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error) {
	if r.readUserErr != nil {
		return nil, r.readUserErr
	}
	id := req.GetData().GetId()
	return &userpb.ReadUserResponse{
		Data: []*userpb.User{{
			Id:           id,
			FirstName:    "Test",
			LastName:     "User",
			EmailAddress: "test@test.com",
			MobileNumber: "+639000000000",
			Active:       true,
		}},
	}, nil
}

func (r *userActionRecorder) updateUser(_ context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	d := req.GetData()
	r.updateUserCalls = append(r.updateUserCalls, updateUserCall{
		id:           d.GetId(),
		firstName:    d.GetFirstName(),
		lastName:     d.GetLastName(),
		email:        d.GetEmailAddress(),
		mobile:       d.GetMobileNumber(),
		passwordHash: d.GetPasswordHash(),
		active:       d.GetActive(),
	})
	if r.updateUserErr != nil {
		return nil, r.updateUserErr
	}
	return &userpb.UpdateUserResponse{}, nil
}

func testMessages() map[string]string {
	return map[string]string{
		"shared.errors.permissionDenied":    "permission denied",
		"shared.errors.idRequired":          "id required",
		"shared.errors.noIdsProvided":       "no ids provided",
		"shared.errors.invalidStatus":       "invalid status",
		"shared.errors.invalidTargetStatus": "invalid target status",
		"shared.errors.invalidFormData":     "invalid form data",
		"shared.errors.passwordFailed":      "password failed",
		"shared.errors.passwordRequired":    "password required",
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
			perms:           []string{"user:read"},
			req:             makePostRequest("/action/users/delete?id=u-1", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "missing id",
			perms:           []string{"user:delete"},
			req:             makePostRequest("/action/users/delete", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:          "uses id from query and succeeds",
			perms:         []string{"user:delete"},
			req:           makePostRequest("/action/users/delete?id=u-q", nil),
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"u-q"},
		},
		{
			name:          "falls back to form id and succeeds",
			perms:         []string{"user:delete"},
			req:           makePostRequest("/action/users/delete", url.Values{"id": {"u-f"}}),
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"u-f"},
		},
		{
			name:            "dependency error returns htmx error",
			perms:           []string{"user:delete"},
			req:             makePostRequest("/action/users/delete?id=u-e", nil),
			deleteErrByID:   map[string]error{"u-e": errors.New("delete failed")},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "delete failed",
			wantDeleteIDs:   []string{"u-e"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{deleteErrByID: tt.deleteErrByID}
			deps := &Deps{DeleteUser: rec.deleteUser}
			res := runHandler(t, NewDeleteAction(deps), withPerms(tt.perms...), tt.req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "users-table")
			}

			gotIDs := make([]string, 0, len(rec.deleteCalls))
			for _, c := range rec.deleteCalls {
				gotIDs = append(gotIDs, c.id)
			}
			if strings.Join(gotIDs, ",") != strings.Join(tt.wantDeleteIDs, ",") {
				t.Fatalf("DeleteUser IDs = %v, want %v", gotIDs, tt.wantDeleteIDs)
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
			perms:           []string{"user:read"},
			form:            url.Values{"id": {"u-1"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "no ids provided",
			perms:           []string{"user:delete"},
			form:            url.Values{},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "no ids provided",
		},
		{
			name:          "success deletes all ids",
			perms:         []string{"user:delete"},
			form:          url.Values{"id": {"u-1", "u-2", "u-3"}},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"u-1", "u-2", "u-3"},
		},
		{
			name:          "partial failure still returns success",
			perms:         []string{"user:delete"},
			form:          url.Values{"id": {"u-1", "u-2", "u-3"}},
			deleteErrByID: map[string]error{"u-2": errors.New("boom")},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: []string{"u-1", "u-2", "u-3"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{deleteErrByID: tt.deleteErrByID}
			deps := &Deps{DeleteUser: rec.deleteUser}
			req := makePostRequest("/action/users/bulk-delete", tt.form)
			res := runHandler(t, NewBulkDeleteAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "users-table")
			}

			gotIDs := make([]string, 0, len(rec.deleteCalls))
			for _, c := range rec.deleteCalls {
				gotIDs = append(gotIDs, c.id)
			}
			if strings.Join(gotIDs, ",") != strings.Join(tt.wantDeleteIDs, ",") {
				t.Fatalf("DeleteUser IDs = %v, want %v", gotIDs, tt.wantDeleteIDs)
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
			perms:           []string{"user:read"},
			req:             makePostRequest("/action/users/set-status?id=u-1&status=active", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "missing id",
			perms:           []string{"user:update"},
			req:             makePostRequest("/action/users/set-status?status=active", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:            "invalid status",
			perms:           []string{"user:update"},
			req:             makePostRequest("/action/users/set-status?id=u-1&status=paused", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "invalid status",
		},
		{
			name:       "query params active success",
			perms:      []string{"user:update"},
			req:        makePostRequest("/action/users/set-status?id=u-1&status=active", nil),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "u-1", active: true}},
		},
		{
			name:       "form fallback inactive success",
			perms:      []string{"user:update"},
			req:        makePostRequest("/action/users/set-status", url.Values{"id": {"u-2"}, "status": {"inactive"}}),
			wantStatus: http.StatusOK,
			wantCalls:  []statusCall{{id: "u-2", active: false}},
		},
		{
			name:            "dependency error",
			perms:           []string{"user:update"},
			req:             makePostRequest("/action/users/set-status?id=u-3&status=inactive", nil),
			statusErrByID:   map[string]error{"u-3": errors.New("set status failed")},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "set status failed",
			wantCalls:       []statusCall{{id: "u-3", active: false}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{statusErrByID: tt.statusErrByID}
			deps := &Deps{SetUserActive: rec.setUserActive}
			res := runHandler(t, NewSetStatusAction(deps), withPerms(tt.perms...), tt.req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "users-table")
			}

			if len(rec.statusCalls) != len(tt.wantCalls) {
				t.Fatalf("SetUserActive call count = %d, want %d", len(rec.statusCalls), len(tt.wantCalls))
			}
			for i := range tt.wantCalls {
				if rec.statusCalls[i] != tt.wantCalls[i] {
					t.Fatalf("SetUserActive call[%d] = %+v, want %+v", i, rec.statusCalls[i], tt.wantCalls[i])
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
			perms:           []string{"user:read"},
			form:            url.Values{"id": {"u-1"}, "target_status": {"active"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:            "no ids provided",
			perms:           []string{"user:update"},
			form:            url.Values{"target_status": {"active"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "no ids provided",
		},
		{
			name:            "invalid target status",
			perms:           []string{"user:update"},
			form:            url.Values{"id": {"u-1"}, "target_status": {"paused"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "invalid target status",
		},
		{
			name:       "active success",
			perms:      []string{"user:update"},
			form:       url.Values{"id": {"u-1", "u-2"}, "target_status": {"active"}},
			wantStatus: http.StatusOK,
			wantCalls: []statusCall{
				{id: "u-1", active: true},
				{id: "u-2", active: true},
			},
		},
		{
			name:       "inactive success",
			perms:      []string{"user:update"},
			form:       url.Values{"id": {"u-1", "u-2"}, "target_status": {"inactive"}},
			wantStatus: http.StatusOK,
			wantCalls: []statusCall{
				{id: "u-1", active: false},
				{id: "u-2", active: false},
			},
		},
		{
			name:          "partial failure still returns success",
			perms:         []string{"user:update"},
			form:          url.Values{"id": {"u-1", "u-2", "u-3"}, "target_status": {"inactive"}},
			statusErrByID: map[string]error{"u-2": errors.New("boom")},
			wantStatus:    http.StatusOK,
			wantCalls: []statusCall{
				{id: "u-1", active: false},
				{id: "u-2", active: false},
				{id: "u-3", active: false},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{statusErrByID: tt.statusErrByID}
			deps := &Deps{SetUserActive: rec.setUserActive}
			req := makePostRequest("/action/users/bulk-set-status", tt.form)
			res := runHandler(t, NewBulkSetStatusAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "users-table")
			}

			if len(rec.statusCalls) != len(tt.wantCalls) {
				t.Fatalf("SetUserActive call count = %d, want %d", len(rec.statusCalls), len(tt.wantCalls))
			}
			for i := range tt.wantCalls {
				if rec.statusCalls[i] != tt.wantCalls[i] {
					t.Fatalf("SetUserActive call[%d] = %+v, want %+v", i, rec.statusCalls[i], tt.wantCalls[i])
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
		name             string
		perms            []string
		form             url.Values
		createUserErr    error
		createWSErr      error
		hashPassword     func(string) (string, error)
		defaultWSID      string
		wantStatus       int
		wantErrorHeader  string
		wantPasswordHash string
		wantMobile       string
		wantWSCalls      int
	}{
		{
			name:            "permission denied",
			perms:           []string{"user:read"},
			form:            url.Values{"first_name": {"Alice"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "permission denied",
		},
		{
			name:  "success without password or workspace",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"alice@test.com"},
				"mobile_number": {"+639111222333"},
				"active":        {"true"},
			},
			wantStatus:       http.StatusOK,
			wantPasswordHash: "",
			wantMobile:       "+639111222333",
			wantWSCalls:      0,
		},
		{
			name:  "mobile defaults when empty",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Bob"},
				"last_name":     {"Jones"},
				"email_address": {"bob@test.com"},
				"active":        {"true"},
			},
			wantStatus:  http.StatusOK,
			wantMobile:  "+639000000000",
			wantWSCalls: 0,
		},
		{
			name:  "password hashed when HashPassword provided",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Carol"},
				"last_name":     {"Wu"},
				"email_address": {"carol@test.com"},
				"password":      {"secret123"},
				"active":        {"true"},
			},
			hashPassword:     func(pw string) (string, error) { return "hashed:" + pw, nil },
			wantStatus:       http.StatusOK,
			wantPasswordHash: "hashed:secret123",
			wantWSCalls:      0,
		},
		{
			name:  "password stored as-is when HashPassword is nil",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Dave"},
				"last_name":     {"Lee"},
				"email_address": {"dave@test.com"},
				"password":      {"plaintext"},
				"active":        {"true"},
			},
			wantStatus:       http.StatusOK,
			wantPasswordHash: "plaintext",
			wantWSCalls:      0,
		},
		{
			name:  "hash password error",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Err"},
				"last_name":     {"Hash"},
				"email_address": {"err@test.com"},
				"password":      {"pw"},
				"active":        {"true"},
			},
			hashPassword:    func(pw string) (string, error) { return "", errors.New("hash boom") },
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "password failed",
		},
		{
			name:  "create user error",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Fail"},
				"last_name":     {"Create"},
				"email_address": {"fail@test.com"},
				"active":        {"true"},
			},
			createUserErr:   errors.New("create failed"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "create failed",
		},
		{
			name:  "workspace user created on success",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Eve"},
				"last_name":     {"Fox"},
				"email_address": {"eve@test.com"},
				"active":        {"true"},
			},
			defaultWSID: "ws-default",
			wantStatus:  http.StatusOK,
			wantWSCalls: 1,
		},
		{
			name:  "workspace user error is non-fatal",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Grace"},
				"last_name":     {"Hop"},
				"email_address": {"grace@test.com"},
				"active":        {"true"},
			},
			defaultWSID: "ws-default",
			createWSErr: errors.New("ws boom"),
			wantStatus:  http.StatusOK,
			wantWSCalls: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{
				createUserErr: tt.createUserErr,
				createWSErr:   tt.createWSErr,
				createdUserID: "new-user-id",
			}
			deps := &Deps{
				CreateUser:   rec.createUser,
				HashPassword: tt.hashPassword,
			}
			if tt.defaultWSID != "" {
				deps.CreateWorkspaceUser = rec.createWorkspaceUser
				deps.DefaultWorkspaceID = tt.defaultWSID
			}

			req := makePostRequest("/action/users/add", tt.form)
			res := runHandler(t, NewAddAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "users-table")
			}

			// Verify password hash
			if tt.wantPasswordHash != "" && len(rec.createUserCalls) > 0 {
				if got := rec.createUserCalls[0].passwordHash; got != tt.wantPasswordHash {
					t.Fatalf("passwordHash = %q, want %q", got, tt.wantPasswordHash)
				}
			}

			// Verify mobile default
			if tt.wantMobile != "" && len(rec.createUserCalls) > 0 {
				if got := rec.createUserCalls[0].mobile; got != tt.wantMobile {
					t.Fatalf("mobile = %q, want %q", got, tt.wantMobile)
				}
			}

			// Verify workspace user calls
			if len(rec.createWSCalls) != tt.wantWSCalls {
				t.Fatalf("CreateWorkspaceUser call count = %d, want %d", len(rec.createWSCalls), tt.wantWSCalls)
			}
			if tt.wantWSCalls > 0 && len(rec.createWSCalls) > 0 {
				if got := rec.createWSCalls[0].workspaceID; got != tt.defaultWSID {
					t.Fatalf("workspace ID = %q, want %q", got, tt.defaultWSID)
				}
				if got := rec.createWSCalls[0].userID; got != "new-user-id" {
					t.Fatalf("user ID = %q, want %q", got, "new-user-id")
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

	longString := strings.Repeat("A", 300)

	tests := []struct {
		name            string
		perms           []string
		form            url.Values
		createUserErr   error
		wantStatus      int
		wantErrorHeader string
		wantCreateCount int
		wantFirstName   string
		wantLastName    string
		wantEmail       string
	}{
		{
			name:  "missing all required fields still calls create (no server-side form validation)",
			perms: []string{"user:create"},
			form: url.Values{
				"active": {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantFirstName:   "",
			wantLastName:    "",
			wantEmail:       "",
		},
		{
			name:  "missing first_name",
			perms: []string{"user:create"},
			form: url.Values{
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantFirstName:   "",
		},
		{
			name:  "missing last_name",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Alice"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantLastName:    "",
		},
		{
			name:  "missing email",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name": {"Alice"},
				"last_name":  {"Smith"},
				"active":     {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantEmail:       "",
		},
		{
			name:  "invalid email format passes through (no server-side email validation)",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"not-an-email"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantEmail:       "not-an-email",
		},
		{
			name:  "very long first_name (>255 chars) passes through",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {longString},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantFirstName:   longString,
		},
		{
			name:  "very long last_name (>255 chars) passes through",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Alice"},
				"last_name":     {longString},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantLastName:    longString,
		},
		{
			name:  "very long email (>255 chars) passes through",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {longString + "@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantEmail:       longString + "@test.com",
		},
		{
			name:  "XSS-like script tag in first_name passes through",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"<script>alert('xss')</script>"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantFirstName:   "<script>alert('xss')</script>",
		},
		{
			name:  "XSS-like script tag in last_name passes through",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Alice"},
				"last_name":     {"<img src=x onerror=alert(1)>"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantLastName:    "<img src=x onerror=alert(1)>",
		},
		{
			name:  "unicode and special characters in name fields",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Ñoño"},
				"last_name":     {"O'Brien-Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantFirstName:   "Ñoño",
			wantLastName:    "O'Brien-Smith",
		},
		{
			name:  "empty string email_address",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {""},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantEmail:       "",
		},
		{
			name:  "whitespace-only first_name",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"   "},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
			wantFirstName:   "   ",
		},
		{
			name:  "create user backend rejects missing required fields",
			perms: []string{"user:create"},
			form: url.Values{
				"active": {"true"},
			},
			createUserErr:   errors.New("first_name is required"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "first_name is required",
			wantCreateCount: 1,
		},
		{
			name:  "active field absent defaults to false",
			perms: []string{"user:create"},
			form: url.Values{
				"first_name":    {"Alice"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
			},
			wantStatus:      http.StatusOK,
			wantCreateCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{
				createUserErr: tt.createUserErr,
				createdUserID: "new-user-id",
			}
			deps := &Deps{CreateUser: rec.createUser}

			req := makePostRequest("/action/users/add", tt.form)
			res := runHandler(t, NewAddAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if tt.wantStatus == http.StatusOK {
				assertSuccessHeader(t, res, "users-table")
			}
			if len(rec.createUserCalls) != tt.wantCreateCount {
				t.Fatalf("CreateUser call count = %d, want %d", len(rec.createUserCalls), tt.wantCreateCount)
			}
			if tt.wantCreateCount > 0 && len(rec.createUserCalls) > 0 {
				got := rec.createUserCalls[0]
				if tt.wantFirstName != "" && got.firstName != tt.wantFirstName {
					t.Fatalf("firstName = %q, want %q", got.firstName, tt.wantFirstName)
				}
				if tt.wantLastName != "" && got.lastName != tt.wantLastName {
					t.Fatalf("lastName = %q, want %q", got.lastName, tt.wantLastName)
				}
				if tt.wantEmail != "" && got.email != tt.wantEmail {
					t.Fatalf("email = %q, want %q", got.email, tt.wantEmail)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// EditAction — negative / defensive tests (empty/missing ID)
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
			perms:   []string{"user:update"},
			pathURL: "/action/users//edit",
			form: url.Values{
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
			name:    "edit with XSS in name passes through to update",
			perms:   []string{"user:update"},
			pathURL: "/action/users/u-1/edit",
			form: url.Values{
				"first_name":    {"<script>alert(1)</script>"},
				"last_name":     {"Smith"},
				"email_address": {"test@test.com"},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantUpdateCount: 1,
		},
		{
			name:    "edit with very long strings",
			perms:   []string{"user:update"},
			pathURL: "/action/users/u-1/edit",
			form: url.Values{
				"first_name":    {strings.Repeat("B", 300)},
				"last_name":     {strings.Repeat("C", 300)},
				"email_address": {strings.Repeat("D", 300)},
				"active":        {"true"},
			},
			wantStatus:      http.StatusOK,
			wantUpdateCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{updateUserErr: tt.updateErr}
			deps := &Deps{
				UpdateUser: rec.updateUser,
				ReadUser:   rec.readUser,
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
				assertSuccessHeader(t, res, "users-table")
			}
			if len(rec.updateUserCalls) != tt.wantUpdateCount {
				t.Fatalf("UpdateUser call count = %d, want %d", len(rec.updateUserCalls), tt.wantUpdateCount)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ResetPasswordAction — negative / defensive tests
// ---------------------------------------------------------------------------

func makePostRequestWithPathValue(rawURL string, form url.Values, pathKey, pathVal string) *http.Request {
	req := makePostRequest(rawURL, form)
	req.SetPathValue(pathKey, pathVal)
	return req
}

func TestNewResetPasswordAction_Negative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		pathID          string // value set via SetPathValue("id", ...)
		form            url.Values
		readUserErr     error
		updateUserErr   error
		hashPassword    func(string) (string, error)
		wantStatus      int
		wantErrorHeader string
		wantUpdateCount int
	}{
		{
			name:            "missing user ID in path",
			pathID:          "",
			form:            url.Values{"password": {"newpass"}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:            "empty password",
			pathID:          "u-1",
			form:            url.Values{"password": {""}},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "password required",
		},
		{
			name:            "no password field at all",
			pathID:          "u-1",
			form:            url.Values{},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "password required",
		},
		{
			name:            "read user fails (user not found)",
			pathID:          "u-missing",
			form:            url.Values{"password": {"newpass"}},
			readUserErr:     errors.New("user not found"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "not found",
		},
		{
			name:            "hash password fails",
			pathID:          "u-1",
			form:            url.Values{"password": {"newpass"}},
			hashPassword:    func(string) (string, error) { return "", errors.New("hash error") },
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "password failed",
		},
		{
			name:            "update user fails",
			pathID:          "u-1",
			form:            url.Values{"password": {"newpass"}},
			updateUserErr:   errors.New("update failed"),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "update failed",
			wantUpdateCount: 1,
		},
		{
			name:            "very long password passes through",
			pathID:          "u-1",
			form:            url.Values{"password": {strings.Repeat("P", 1000)}},
			wantStatus:      http.StatusOK,
			wantUpdateCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{
				readUserErr:   tt.readUserErr,
				updateUserErr: tt.updateUserErr,
			}
			deps := &Deps{
				ReadUser:     rec.readUser,
				UpdateUser:   rec.updateUser,
				HashPassword: tt.hashPassword,
			}

			req := makePostRequestWithPathValue("/action/users/reset-password", tt.form, "id", tt.pathID)
			ctx := withPerms("user:update")
			res := runHandler(t, NewResetPasswordAction(deps), ctx, req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantErrorHeader != "" {
				assertErrorHeader(t, res, tt.wantErrorHeader)
			}
			if len(rec.updateUserCalls) != tt.wantUpdateCount {
				t.Fatalf("UpdateUser call count = %d, want %d", len(rec.updateUserCalls), tt.wantUpdateCount)
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
			perms:           []string{"user:delete"},
			req:             makePostRequest("/action/users/delete?id=", nil),
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorHeader: "id required",
		},
		{
			name:       "whitespace-only id in form",
			perms:      []string{"user:delete"},
			req:        makePostRequest("/action/users/delete", url.Values{"id": {"   "}}),
			wantStatus: http.StatusOK, // whitespace is not validated, treated as a valid ID
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{}
			deps := &Deps{DeleteUser: rec.deleteUser}
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

// ---------------------------------------------------------------------------
// BulkDelete — additional edge cases
// ---------------------------------------------------------------------------

func TestNewBulkDeleteAction_Negative(t *testing.T) {
	t.Parallel()

	largeIDSet := make([]string, 100)
	for i := range largeIDSet {
		largeIDSet[i] = fmt.Sprintf("u-%d", i)
	}

	tests := []struct {
		name          string
		perms         []string
		form          url.Values
		wantStatus    int
		wantDeleteIDs int
	}{
		{
			name:          "large bulk delete with 100 IDs",
			perms:         []string{"user:delete"},
			form:          url.Values{"id": largeIDSet},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: 100,
		},
		{
			name:          "single empty-string ID treated as valid (len=1)",
			perms:         []string{"user:delete"},
			form:          url.Values{"id": {""}},
			wantStatus:    http.StatusOK,
			wantDeleteIDs: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &userActionRecorder{}
			deps := &Deps{DeleteUser: rec.deleteUser}
			req := makePostRequest("/action/users/bulk-delete", tt.form)
			res := runHandler(t, NewBulkDeleteAction(deps), withPerms(tt.perms...), req)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if len(rec.deleteCalls) != tt.wantDeleteIDs {
				t.Fatalf("DeleteUser call count = %d, want %d", len(rec.deleteCalls), tt.wantDeleteIDs)
			}
		})
	}
}

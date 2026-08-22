// Package action provides handlers for staff mutations.
// Handles: Add, Edit, Delete, BulkDelete, and the user-search JSON autocomplete.
//
// Staff links to an EXISTING user (unlike party/delegate and party/client,
// which embed a brand-new User sub-message on create) — the drawer form
// carries only user_id (chosen via the auto-complete SearchURL, same pattern
// as identity/workspace_user). Because CreateStaffUseCase/UpdateStaffUseCase
// validate a populated User sub-message (first/last/email) on the request,
// this handler resolves the picked user_id to a full User record via ReadUser
// before calling CreateStaff/UpdateStaff — see resolveUser below.
package action

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	staffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/staff"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	entitystaff "github.com/erniealice/entydad-golang/domain/entity/party/staff"
	staffform "github.com/erniealice/entydad-golang/domain/entity/party/staff/form"
)

// Deps holds dependencies for staff action handlers.
// No payment-terms/categories/tags — Staff has only user_id, employment_type,
// status, and active (per the entity-status convention).
type Deps struct {
	Routes          entitystaff.Routes
	CreateStaff     func(ctx context.Context, req *staffpb.CreateStaffRequest) (*staffpb.CreateStaffResponse, error)
	GetItemPageData func(ctx context.Context, req *staffpb.GetStaffItemPageDataRequest) (*staffpb.GetStaffItemPageDataResponse, error)
	UpdateStaff     func(ctx context.Context, req *staffpb.UpdateStaffRequest) (*staffpb.UpdateStaffResponse, error)
	DeleteStaff     func(ctx context.Context, req *staffpb.DeleteStaffRequest) (*staffpb.DeleteStaffResponse, error)
	// ReadUser resolves a picked user_id into a full User record — required
	// because the espyna CreateStaff/UpdateStaff use cases validate a
	// populated User sub-message, not a bare user_id reference.
	ReadUser func(ctx context.Context, req *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error)
	// ListUsers is used by the user search endpoint to find users for autocomplete.
	ListUsers func(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error)
}

// searchOption is the JSON shape returned by the user search handler.
type searchOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// resolveUser resolves a user_id form value into a full *userpb.User via
// ReadUser. Returns nil, nil when userID is empty (caller decides whether
// that's an error).
func resolveUser(ctx context.Context, deps *Deps, userID string) (*userpb.User, error) {
	if userID == "" || deps.ReadUser == nil {
		return nil, nil
	}
	resp, err := deps.ReadUser(ctx, &userpb.ReadUserRequest{Data: &userpb.User{Id: userID}})
	if err != nil {
		return nil, err
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, nil
	}
	return data[0], nil
}

// userDisplayLabel formats "First Last (email)" for the auto-complete's
// pre-filled SelectedLabel on the edit drawer.
func userDisplayLabel(u *userpb.User) string {
	if u == nil {
		return ""
	}
	name := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
	email := u.GetEmailAddress()
	if email == "" {
		return name
	}
	if name == "" {
		return email
	}
	return name + " (" + email + ")"
}

// NewAddAction creates the staff add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("staff", "create") {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			labels := staffform.BuildLabels(viewCtx.T)
			return view.OK("staff-drawer-form", &staffform.Data{
				FormAction:            deps.Routes.AddURL,
				UserSearchURL:         deps.Routes.SearchURL,
				Status:                staffform.StatusAvailable,
				EmploymentTypeOptions: staffform.BuildEmploymentTypeOptions("", labels),
				StatusOptions:         staffform.BuildStatusOptions(staffform.StatusAvailable, labels),
				Labels:                labels,
				CommonLabels:          nil, // injected by ViewAdapter
			})
		}

		// POST — create staff
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalid_form_data"))
		}

		r := viewCtx.Request
		userID := r.FormValue("user_id")
		if userID == "" {
			return view.HTMXError(viewCtx.T("shared.errors.invalid_form_data"))
		}

		u, err := resolveUser(ctx, deps, userID)
		if err != nil {
			log.Printf("Failed to resolve user %s for staff create: %v", userID, err)
			return view.HTMXError(err.Error())
		}
		if u == nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalid_form_data"))
		}

		staffData := &staffpb.Staff{
			UserId: userID,
			User:   u,
			Active: true,
		}
		if et := r.FormValue("employment_type"); et != "" {
			staffData.EmploymentType = &et
		}
		if st := r.FormValue("status"); st != "" {
			staffData.Status = &st
		} else {
			defaultStatus := staffform.StatusAvailable
			staffData.Status = &defaultStatus
		}

		_, err = deps.CreateStaff(ctx, &staffpb.CreateStaffRequest{Data: staffData})
		if err != nil {
			log.Printf("Failed to create staff: %v", err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("staffs-table")
	})
}

// NewEditAction creates the staff edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("staff", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			// GetItemPageData enriches Staff.User via a SQL JOIN — needed to
			// pre-fill the user auto-complete's display label.
			resp, err := deps.GetItemPageData(ctx, &staffpb.GetStaffItemPageDataRequest{StaffId: id})
			if err != nil {
				log.Printf("Failed to read staff %s: %v", id, err)
				return view.HTMXError(viewCtx.T("shared.errors.not_found"))
			}

			s := resp.GetStaff()
			labels := staffform.BuildLabels(viewCtx.T)
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id)
			return view.OK("staff-drawer-form", &staffform.Data{
				FormAction:            formAction,
				IsEdit:                true,
				ID:                    id,
				UserID:                s.GetUserId(),
				UserSelectedLabel:     userDisplayLabel(s.GetUser()),
				UserSearchURL:         deps.Routes.SearchURL,
				EmploymentType:        s.GetEmploymentType(),
				Status:                s.GetStatus(),
				EmploymentTypeOptions: staffform.BuildEmploymentTypeOptions(s.GetEmploymentType(), labels),
				StatusOptions:         staffform.BuildStatusOptions(s.GetStatus(), labels),
				Labels:                labels,
				CommonLabels:          nil, // injected by ViewAdapter
			})
		}

		// POST — update staff
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalid_form_data"))
		}

		r := viewCtx.Request
		staffData := &staffpb.Staff{Id: id}

		// Only re-resolve the User sub-message when the operator actually
		// picked a (possibly different) user on the edit form; an empty
		// user_id leaves the existing staff.user_id column untouched (proto3
		// zero-value fields are omitted from the marshaled update).
		if userID := r.FormValue("user_id"); userID != "" {
			u, err := resolveUser(ctx, deps, userID)
			if err != nil {
				log.Printf("Failed to resolve user %s for staff update: %v", userID, err)
				return view.HTMXError(err.Error())
			}
			if u != nil {
				staffData.UserId = userID
				staffData.User = u
			}
		}
		if et := r.FormValue("employment_type"); et != "" {
			staffData.EmploymentType = &et
		}
		if st := r.FormValue("status"); st != "" {
			staffData.Status = &st
		}

		_, err := deps.UpdateStaff(ctx, &staffpb.UpdateStaffRequest{Data: staffData})
		if err != nil {
			log.Printf("Failed to update staff %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("staffs-table")
	})
}

// NewDeleteAction creates the staff delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("staff", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(viewCtx.T("shared.errors.id_required"))
		}

		_, err := deps.DeleteStaff(ctx, &staffpb.DeleteStaffRequest{Data: &staffpb.Staff{Id: id}})
		if err != nil {
			log.Printf("Failed to delete staff %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("staffs-table")
	})
}

// NewBulkDeleteAction creates the staff bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("staff", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.no_ids_provided"))
		}

		for _, id := range ids {
			_, err := deps.DeleteStaff(ctx, &staffpb.DeleteStaffRequest{Data: &staffpb.Staff{Id: id}})
			if err != nil {
				log.Printf("Failed to delete staff %s: %v", id, err)
			}
		}

		return view.HTMXSuccess("staffs-table")
	})
}

// NewUserSearchAction returns an http.HandlerFunc that searches users for the
// staff drawer's "User" autocomplete input.
// GET /action/staff/search?q=term
// Returns JSON: [{"value":"user_id","label":"First Last (email)"}, ...]
// Mirrors identity/workspace_user/action.NewUserSearchAction.
func NewUserSearchAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("staff", "create") {
			writeSearchJSON(w, http.StatusForbidden, []searchOption{})
			return
		}
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))

		if query == "" {
			writeSearchJSON(w, http.StatusOK, []searchOption{})
			return
		}

		if deps.ListUsers == nil {
			writeSearchJSON(w, http.StatusOK, []searchOption{})
			return
		}

		resp, err := deps.ListUsers(ctx, &userpb.ListUsersRequest{})
		if err != nil {
			log.Printf("staff search: failed to list users: %v", err)
			writeSearchJSON(w, http.StatusOK, []searchOption{})
			return
		}

		var results []searchOption
		for _, u := range resp.GetData() {
			if !u.GetActive() {
				continue
			}
			name := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			email := u.GetEmailAddress()
			label := name
			if email != "" {
				label = label + " (" + email + ")"
			}
			if !strings.Contains(strings.ToLower(label), query) {
				continue
			}
			results = append(results, searchOption{
				Value: u.GetId(),
				Label: label,
			})
		}

		if results == nil {
			results = []searchOption{}
		}
		writeSearchJSON(w, http.StatusOK, results)
	}
}

// writeSearchJSON marshals data as JSON and writes it to the response writer
// with the given status code. Content-Type is set BEFORE WriteHeader — see
// identity/workspace_user/action.writeSearchJSON for the nosniff rationale.
func writeSearchJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.MarshalEncode(jsontext.NewEncoder(w), data); err != nil {
		log.Printf("staff search: failed to encode JSON response: %v", err)
	}
}

package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"

	"github.com/erniealice/entydad-golang"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	FirstName            string
	FirstNamePlaceholder string
	LastName             string
	LastNamePlaceholder  string
	Email                string
	EmailPlaceholder     string
	Mobile               string
	MobilePlaceholder    string
	Password             string
	PasswordPlaceholder  string
	PasswordGenerate     string
	Active               string
}

// FormData is the template data for the user drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	FirstName    string
	LastName     string
	Email        string
	Mobile       string
	Active       bool
	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for user action handlers.
type Deps struct {
	Routes              entydad.UserRoutes
	CreateUser          func(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error)
	ReadUser            func(ctx context.Context, req *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error)
	UpdateUser          func(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error)
	DeleteUser          func(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error)
	SetUserActive       func(ctx context.Context, id string, active bool) error
	CreateWorkspaceUser func(ctx context.Context, req *workspaceuserpb.CreateWorkspaceUserRequest) (*workspaceuserpb.CreateWorkspaceUserResponse, error)
	DefaultWorkspaceID  string
	HashPassword        func(password string) (string, error) // optional; if nil, password stored as-is
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		FirstName:            t("form.firstName"),
		FirstNamePlaceholder: t("form.firstNamePlaceholder"),
		LastName:             t("form.lastName"),
		LastNamePlaceholder:  t("form.lastNamePlaceholder"),
		Email:                t("form.email"),
		EmailPlaceholder:     t("form.emailPlaceholder"),
		Mobile:               t("form.mobile"),
		MobilePlaceholder:    t("form.mobilePlaceholder"),
		Password:             t("form.password"),
		PasswordPlaceholder:  t("form.passwordPlaceholder"),
		PasswordGenerate:     t("form.passwordGenerate"),
		Active:               t("form.active"),
	}
}

// hashPassword hashes the password using the deps.HashPassword func, or returns it as-is.
func hashPassword(deps *Deps, password string) (string, error) {
	if deps.HashPassword != nil {
		return deps.HashPassword(password)
	}
	return password, nil
}

// NewAddAction creates the user add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("user-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create user
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		var pwHash string
		if pw := r.FormValue("password"); pw != "" {
			h, hashErr := hashPassword(deps, pw)
			if hashErr != nil {
				log.Printf("Failed to hash password: %v", hashErr)
				return entydad.HTMXError("Failed to process password")
			}
			pwHash = h
		}

		createResp, err := deps.CreateUser(ctx, &userpb.CreateUserRequest{
			Data: &userpb.User{
				FirstName:    r.FormValue("first_name"),
				LastName:     r.FormValue("last_name"),
				EmailAddress: r.FormValue("email_address"),
				MobileNumber: r.FormValue("mobile_number"),
				PasswordHash: pwHash,
				Active:       active,
			},
		})
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			return entydad.HTMXError(err.Error())
		}

		// Auto-create WorkspaceUser for the default workspace
		if deps.CreateWorkspaceUser != nil && deps.DefaultWorkspaceID != "" {
			newUserID := ""
			if data := createResp.GetData(); len(data) > 0 {
				newUserID = data[0].GetId()
			}
			if newUserID != "" {
				_, err := deps.CreateWorkspaceUser(ctx, &workspaceuserpb.CreateWorkspaceUserRequest{
					Data: &workspaceuserpb.WorkspaceUser{
						WorkspaceId: deps.DefaultWorkspaceID,
						UserId:      newUserID,
						Active:      true,
					},
				})
				if err != nil {
					log.Printf("Warning: Failed to create workspace user for %s: %v", newUserID, err)
				}
			}
		}

		return entydad.HTMXSuccess("users-table")
	})
}

// NewEditAction creates the user edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadUser(ctx, &userpb.ReadUserRequest{
				Data: &userpb.User{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read user %s: %v", id, err)
				return entydad.HTMXError("User not found")
			}

			u := resp.GetData()[0]

			return view.OK("user-drawer-form", &FormData{
				FormAction:   route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:       true,
				ID:           id,
				FirstName:    u.GetFirstName(),
				LastName:     u.GetLastName(),
				Email:        u.GetEmailAddress(),
				Mobile:       u.GetMobileNumber(),
				Active:       u.GetActive(),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — update user
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		userData := &userpb.User{
			Id:           id,
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			EmailAddress: r.FormValue("email_address"),
			MobileNumber: r.FormValue("mobile_number"),
			Active:       active,
		}

		// Only update password if a new one was provided
		if pw := r.FormValue("password"); pw != "" {
			pwHash, hashErr := hashPassword(deps, pw)
			if hashErr != nil {
				log.Printf("Failed to hash password: %v", hashErr)
				return entydad.HTMXError("Failed to process password")
			}
			userData.PasswordHash = pwHash
		}

		_, updateErr := deps.UpdateUser(ctx, &userpb.UpdateUserRequest{
			Data: userData,
		})
		if updateErr != nil {
			log.Printf("Failed to update user %s: %v", id, updateErr)
			return entydad.HTMXError(updateErr.Error())
		}

		return entydad.HTMXSuccess("users-table")
	})
}

// NewDeleteAction creates the user delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError("User ID is required")
		}

		_, err := deps.DeleteUser(ctx, &userpb.DeleteUserRequest{
			Data: &userpb.User{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete user %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("users-table")
	})
}

// NewBulkDeleteAction creates the user bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError("No user IDs provided")
		}

		for _, id := range ids {
			_, err := deps.DeleteUser(ctx, &userpb.DeleteUserRequest{
				Data: &userpb.User{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete user %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("users-table")
	})
}

// NewSetStatusAction creates the user activate/deactivate action (POST only).
// Expects query params: ?id={userId}&status={active|inactive}
//
// Uses SetUserActive (raw map update) instead of UpdateUser (protobuf) because
// proto3's protojson omits bool fields with value false, which means
// deactivation (active=false) would silently be skipped.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return entydad.HTMXError("User ID is required")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid status")
		}

		if err := deps.SetUserActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update user status %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("users-table")
	})
}

// NewResetPasswordAction creates the password reset action (POST only).
// Expects path param {id} and form field "password".
// Reads the existing user first so all required fields are present for the update.
func NewResetPasswordAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return entydad.HTMXError("User ID is required")
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		password := viewCtx.Request.FormValue("password")
		if password == "" {
			return entydad.HTMXError("Password is required")
		}

		// Read existing user to preserve all required fields
		resp, err := deps.ReadUser(ctx, &userpb.ReadUserRequest{
			Data: &userpb.User{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read user %s for password reset: %v", id, err)
			return entydad.HTMXError("User not found")
		}
		data := resp.GetData()
		if len(data) == 0 {
			return entydad.HTMXError("User not found")
		}
		user := data[0]

		pwHash, hashErr := hashPassword(deps, password)
		if hashErr != nil {
			log.Printf("Failed to hash password: %v", hashErr)
			return entydad.HTMXError("Failed to process password")
		}

		user.PasswordHash = pwHash

		_, updateErr := deps.UpdateUser(ctx, &userpb.UpdateUserRequest{
			Data: user,
		})
		if updateErr != nil {
			log.Printf("Failed to reset password for user %s: %v", id, updateErr)
			return entydad.HTMXError(updateErr.Error())
		}

		return entydad.HTMXSuccess("")
	})
}

// NewBulkSetStatusAction creates the user bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return entydad.HTMXError("No user IDs provided")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid target status")
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetUserActive(ctx, id, active); err != nil {
				log.Printf("Failed to update user status %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("users-table")
	})
}

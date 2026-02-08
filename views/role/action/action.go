package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"

	"github.com/erniealice/entydad-golang"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Name                   string
	NamePlaceholder        string
	Description            string
	DescriptionPlaceholder string
	Color                  string
	ColorPlaceholder       string
	Active                 string
}

// FormData is the template data for the role drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Color        string
	Active       bool
	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for role action handlers.
type Deps struct {
	CreateRole    func(ctx context.Context, req *rolepb.CreateRoleRequest) (*rolepb.CreateRoleResponse, error)
	ReadRole      func(ctx context.Context, req *rolepb.ReadRoleRequest) (*rolepb.ReadRoleResponse, error)
	UpdateRole    func(ctx context.Context, req *rolepb.UpdateRoleRequest) (*rolepb.UpdateRoleResponse, error)
	DeleteRole    func(ctx context.Context, req *rolepb.DeleteRoleRequest) (*rolepb.DeleteRoleResponse, error)
	SetRoleActive func(ctx context.Context, id string, active bool) error
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Name:                   t("form.name"),
		NamePlaceholder:        t("form.namePlaceholder"),
		Description:            t("form.description"),
		DescriptionPlaceholder: t("form.descriptionPlaceholder"),
		Color:                  t("form.color"),
		ColorPlaceholder:       t("form.colorPlaceholder"),
		Active:                 t("form.active"),
	}
}

// NewAddAction creates the role add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("role-drawer-form", &FormData{
				FormAction:   "/action/roles/add",
				Active:       true,
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- create role
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.CreateRole(ctx, &rolepb.CreateRoleRequest{
			Data: &rolepb.Role{
				Name:        r.FormValue("name"),
				Description: r.FormValue("description"),
				Color:       r.FormValue("color"),
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to create role: %v", err)
			return entydad.HTMXError("Failed to create role")
		}

		return entydad.HTMXSuccess("roles-table")
	})
}

// NewEditAction creates the role edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadRole(ctx, &rolepb.ReadRoleRequest{
				Data: &rolepb.Role{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read role %s: %v", id, err)
				return entydad.HTMXError("Role not found")
			}

			role := resp.GetData()[0]

			return view.OK("role-drawer-form", &FormData{
				FormAction:   "/action/roles/edit/" + id,
				IsEdit:       true,
				ID:           id,
				Name:         role.GetName(),
				Description:  role.GetDescription(),
				Color:        role.GetColor(),
				Active:       role.GetActive(),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- update role
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.UpdateRole(ctx, &rolepb.UpdateRoleRequest{
			Data: &rolepb.Role{
				Id:          id,
				Name:        r.FormValue("name"),
				Description: r.FormValue("description"),
				Color:       r.FormValue("color"),
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to update role %s: %v", id, err)
			return entydad.HTMXError("Failed to update role")
		}

		return entydad.HTMXSuccess("roles-table")
	})
}

// NewDeleteAction creates the role delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError("Role ID is required")
		}

		_, err := deps.DeleteRole(ctx, &rolepb.DeleteRoleRequest{
			Data: &rolepb.Role{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete role %s: %v", id, err)
			return entydad.HTMXError("Failed to delete role")
		}

		return entydad.HTMXSuccess("roles-table")
	})
}

// NewBulkDeleteAction creates the role bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError("No role IDs provided")
		}

		for _, id := range ids {
			_, err := deps.DeleteRole(ctx, &rolepb.DeleteRoleRequest{
				Data: &rolepb.Role{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete role %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("roles-table")
	})
}

// NewSetStatusAction creates the role activate/deactivate action (POST only).
// Expects query params: ?id={roleId}&status={active|inactive}
//
// Uses SetRoleActive (raw map update) instead of UpdateRole (protobuf) because
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
			return entydad.HTMXError("Role ID is required")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid status")
		}

		if err := deps.SetRoleActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update role status %s: %v", id, err)
			return entydad.HTMXError("Failed to update role status")
		}

		return entydad.HTMXSuccess("roles-table")
	})
}

// NewBulkSetStatusAction creates the role bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return entydad.HTMXError("No role IDs provided")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid target status")
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetRoleActive(ctx, id, active); err != nil {
				log.Printf("Failed to update role status %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("roles-table")
	})
}

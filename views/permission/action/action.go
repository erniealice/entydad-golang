package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"

	"github.com/erniealice/entydad-golang"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Name                      string
	NamePlaceholder           string
	PermissionCode            string
	PermissionCodePlaceholder string
	PermissionCodeHint        string
	PermissionType            string
	Description               string
	DescriptionPlaceholder    string
	Active                    string
}

// FormData is the template data for the permission drawer form.
type FormData struct {
	FormAction            string
	IsEdit                bool
	ID                    string
	Name                  string
	PermissionCode        string
	PermissionType        string
	Description           string
	Active                bool
	Labels                FormLabels
	PermissionTypeOptions []types.SelectOption
	CommonLabels          any
}

// Deps holds dependencies for permission action handlers.
type Deps struct {
	CreatePermission    func(ctx context.Context, req *permissionpb.CreatePermissionRequest) (*permissionpb.CreatePermissionResponse, error)
	ReadPermission      func(ctx context.Context, req *permissionpb.ReadPermissionRequest) (*permissionpb.ReadPermissionResponse, error)
	UpdatePermission    func(ctx context.Context, req *permissionpb.UpdatePermissionRequest) (*permissionpb.UpdatePermissionResponse, error)
	DeletePermission    func(ctx context.Context, req *permissionpb.DeletePermissionRequest) (*permissionpb.DeletePermissionResponse, error)
	SetPermissionActive func(ctx context.Context, id string, active bool) error
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Name:                      t("form.name"),
		NamePlaceholder:           t("form.namePlaceholder"),
		PermissionCode:            t("form.permissionCode"),
		PermissionCodePlaceholder: t("form.permissionCodePlaceholder"),
		PermissionCodeHint:        t("form.permissionCodeHint"),
		PermissionType:            t("form.permissionType"),
		Description:               t("form.description"),
		DescriptionPlaceholder:    t("form.descriptionPlaceholder"),
		Active:                    t("form.active"),
	}
}

func permissionTypeOptions(current string) []types.SelectOption {
	return []types.SelectOption{
		{Value: "PERMISSION_TYPE_ALLOW", Label: "Allow", Selected: current == "PERMISSION_TYPE_ALLOW"},
		{Value: "PERMISSION_TYPE_DENY", Label: "Deny", Selected: current == "PERMISSION_TYPE_DENY"},
	}
}

func parsePermissionType(s string) permissionpb.PermissionType {
	switch s {
	case "PERMISSION_TYPE_DENY":
		return permissionpb.PermissionType_PERMISSION_TYPE_DENY
	default:
		return permissionpb.PermissionType_PERMISSION_TYPE_ALLOW
	}
}

func formatPermissionType(pt permissionpb.PermissionType) string {
	switch pt {
	case permissionpb.PermissionType_PERMISSION_TYPE_DENY:
		return "PERMISSION_TYPE_DENY"
	default:
		return "PERMISSION_TYPE_ALLOW"
	}
}

// NewAddAction creates the permission add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("permission-drawer-form", &FormData{
				FormAction:            "/action/permissions/add",
				Active:                true,
				PermissionType:        "PERMISSION_TYPE_ALLOW",
				Labels:                formLabels(viewCtx.T),
				PermissionTypeOptions: permissionTypeOptions("PERMISSION_TYPE_ALLOW"),
				CommonLabels:          nil,
			})
		}

		// POST -- create permission
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.CreatePermission(ctx, &permissionpb.CreatePermissionRequest{
			Data: &permissionpb.Permission{
				Name:           r.FormValue("name"),
				PermissionCode: r.FormValue("permission_code"),
				PermissionType: parsePermissionType(r.FormValue("permission_type")),
				Description:    r.FormValue("description"),
				Active:         active,
			},
		})
		if err != nil {
			log.Printf("Failed to create permission: %v", err)
			return entydad.HTMXError("Failed to create permission")
		}

		return entydad.HTMXSuccess("permissions-table")
	})
}

// NewEditAction creates the permission edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPermission(ctx, &permissionpb.ReadPermissionRequest{
				Data: &permissionpb.Permission{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read permission %s: %v", id, err)
				return entydad.HTMXError("Permission not found")
			}

			perm := resp.GetData()[0]

			return view.OK("permission-drawer-form", &FormData{
				FormAction:            "/action/permissions/edit/" + id,
				IsEdit:                true,
				ID:                    id,
				Name:                  perm.GetName(),
				PermissionCode:        perm.GetPermissionCode(),
				PermissionType:        formatPermissionType(perm.GetPermissionType()),
				Description:           perm.GetDescription(),
				Active:                perm.GetActive(),
				Labels:                formLabels(viewCtx.T),
				PermissionTypeOptions: permissionTypeOptions(formatPermissionType(perm.GetPermissionType())),
				CommonLabels:          nil,
			})
		}

		// POST -- update permission
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.UpdatePermission(ctx, &permissionpb.UpdatePermissionRequest{
			Data: &permissionpb.Permission{
				Id:             id,
				Name:           r.FormValue("name"),
				PermissionCode: r.FormValue("permission_code"),
				PermissionType: parsePermissionType(r.FormValue("permission_type")),
				Description:    r.FormValue("description"),
				Active:         active,
			},
		})
		if err != nil {
			log.Printf("Failed to update permission %s: %v", id, err)
			return entydad.HTMXError("Failed to update permission")
		}

		return entydad.HTMXSuccess("permissions-table")
	})
}

// NewDeleteAction creates the permission delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError("Permission ID is required")
		}

		_, err := deps.DeletePermission(ctx, &permissionpb.DeletePermissionRequest{
			Data: &permissionpb.Permission{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete permission %s: %v", id, err)
			return entydad.HTMXError("Failed to delete permission")
		}

		return entydad.HTMXSuccess("permissions-table")
	})
}

// NewBulkDeleteAction creates the permission bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError("No permission IDs provided")
		}

		for _, id := range ids {
			_, err := deps.DeletePermission(ctx, &permissionpb.DeletePermissionRequest{
				Data: &permissionpb.Permission{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete permission %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("permissions-table")
	})
}

// NewSetStatusAction creates the permission activate/deactivate action (POST only).
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
			return entydad.HTMXError("Permission ID is required")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid status")
		}

		if err := deps.SetPermissionActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update permission status %s: %v", id, err)
			return entydad.HTMXError("Failed to update permission status")
		}

		return entydad.HTMXSuccess("permissions-table")
	})
}

// NewBulkSetStatusAction creates the permission bulk activate/deactivate action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return entydad.HTMXError("No permission IDs provided")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid target status")
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetPermissionActive(ctx, id, active); err != nil {
				log.Printf("Failed to update permission status %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("permissions-table")
	})
}

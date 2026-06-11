package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"

	permission "github.com/erniealice/entydad-golang/domain/entity/identity/permission"
	"github.com/erniealice/entydad-golang/domain/entity/identity/permission/form"
)

// Deps holds dependencies for permission action handlers.
type Deps struct {
	CreatePermission    func(ctx context.Context, req *permissionpb.CreatePermissionRequest) (*permissionpb.CreatePermissionResponse, error)
	ReadPermission      func(ctx context.Context, req *permissionpb.ReadPermissionRequest) (*permissionpb.ReadPermissionResponse, error)
	UpdatePermission    func(ctx context.Context, req *permissionpb.UpdatePermissionRequest) (*permissionpb.UpdatePermissionResponse, error)
	DeletePermission    func(ctx context.Context, req *permissionpb.DeletePermissionRequest) (*permissionpb.DeletePermissionResponse, error)
	SetPermissionActive func(ctx context.Context, id string, active bool) error
	Routes              permission.Routes
}

// NewAddAction creates the permission add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("permission", "create") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("permission-drawer-form", &form.Data{
				FormAction:            deps.Routes.AddURL,
				Active:                true,
				PermissionType:        "PERMISSION_TYPE_ALLOW",
				Labels:                form.BuildLabels(viewCtx.T),
				PermissionTypeOptions: form.BuildPermissionTypeOptions("PERMISSION_TYPE_ALLOW", viewCtx.T),
				CommonLabels:          nil,
			})
		}

		// POST -- create permission
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.CreatePermission(ctx, &permissionpb.CreatePermissionRequest{
			Data: &permissionpb.Permission{
				Name:           r.FormValue("name"),
				PermissionCode: r.FormValue("permission_code"),
				PermissionType: form.ParsePermissionType(r.FormValue("permission_type")),
				Description:    r.FormValue("description"),
				Active:         active,
			},
		})
		if err != nil {
			log.Printf("Failed to create permission: %v", err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("permissions-table")
	})
}

// NewEditAction creates the permission edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("permission", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPermission(ctx, &permissionpb.ReadPermissionRequest{
				Data: &permissionpb.Permission{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read permission %s: %v", id, err)
				return view.HTMXError(viewCtx.T("shared.errors.notFound"))
			}

			perm := resp.GetData()[0]

			return view.OK("permission-drawer-form", &form.Data{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				ID:                    id,
				Name:                  perm.GetName(),
				PermissionCode:        perm.GetPermissionCode(),
				PermissionType:        form.FormatPermissionType(perm.GetPermissionType()),
				Description:           perm.GetDescription(),
				Active:                perm.GetActive(),
				Labels:                form.BuildLabels(viewCtx.T),
				PermissionTypeOptions: form.BuildPermissionTypeOptions(form.FormatPermissionType(perm.GetPermissionType()), viewCtx.T),
				CommonLabels:          nil,
			})
		}

		// POST -- update permission
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.UpdatePermission(ctx, &permissionpb.UpdatePermissionRequest{
			Data: &permissionpb.Permission{
				Id:             id,
				Name:           r.FormValue("name"),
				PermissionCode: r.FormValue("permission_code"),
				PermissionType: form.ParsePermissionType(r.FormValue("permission_type")),
				Description:    r.FormValue("description"),
				Active:         active,
			},
		})
		if err != nil {
			log.Printf("Failed to update permission %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("permissions-table")
	})
}

// NewDeleteAction creates the permission delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("permission", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}

		_, err := deps.DeletePermission(ctx, &permissionpb.DeletePermissionRequest{
			Data: &permissionpb.Permission{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete permission %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("permissions-table")
	})
}

// NewBulkDeleteAction creates the permission bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("permission", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}

		for _, id := range ids {
			_, err := deps.DeletePermission(ctx, &permissionpb.DeletePermissionRequest{
				Data: &permissionpb.Permission{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete permission %s: %v", id, err)
			}
		}

		return view.HTMXSuccess("permissions-table")
	})
}

// NewSetStatusAction creates the permission activate/deactivate action (POST only).
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("permission", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return view.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return view.HTMXError(viewCtx.T("shared.errors.invalidStatus"))
		}

		if err := deps.SetPermissionActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update permission status %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("permissions-table")
	})
}

// NewBulkSetStatusAction creates the permission bulk activate/deactivate action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("permission", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return view.HTMXError(viewCtx.T("shared.errors.invalidTargetStatus"))
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetPermissionActive(ctx, id, active); err != nil {
				log.Printf("Failed to update permission status %s: %v", id, err)
			}
		}

		return view.HTMXSuccess("permissions-table")
	})
}

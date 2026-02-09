package permissions

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	rolepermissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role_permission"

	"github.com/erniealice/entydad-golang"
)

// AssignFormLabels holds i18n labels for the assign permission drawer form.
type AssignFormLabels struct {
	Permission string
}

// AssignFormData is the template data for the assign permission drawer form.
type AssignFormData struct {
	FormAction        string
	RoleID            string
	Labels            AssignFormLabels
	PermissionOptions []types.SelectOption
	CommonLabels      any
}

// ActionDeps holds dependencies for role-permission action handlers.
type ActionDeps struct {
	CreateRolePermission func(ctx context.Context, req *rolepermissionpb.CreateRolePermissionRequest) (*rolepermissionpb.CreateRolePermissionResponse, error)
	DeleteRolePermission func(ctx context.Context, req *rolepermissionpb.DeleteRolePermissionRequest) (*rolepermissionpb.DeleteRolePermissionResponse, error)
	ListPermissions      func(ctx context.Context, req *permissionpb.ListPermissionsRequest) (*permissionpb.ListPermissionsResponse, error)
	GetRoleItemPageData  func(ctx context.Context, req *rolepb.GetRoleItemPageDataRequest) (*rolepb.GetRoleItemPageDataResponse, error)
	Labels               entydad.RolePermissionLabels
}

// NewAssignAction creates the assign permission action (GET = form, POST = create).
func NewAssignAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		roleID := viewCtx.Request.PathValue("id")
		if roleID == "" {
			return entydad.HTMXError("Role ID is required")
		}

		if viewCtx.Request.Method == http.MethodGet {
			// Load all permissions for the dropdown
			permResp, err := deps.ListPermissions(ctx, &permissionpb.ListPermissionsRequest{})
			if err != nil {
				log.Printf("Failed to list permissions: %v", err)
				return entydad.HTMXError("Failed to load permissions")
			}

			// Load role with existing permissions to filter out already-assigned
			roleResp, err := deps.GetRoleItemPageData(ctx, &rolepb.GetRoleItemPageDataRequest{
				RoleId: roleID,
			})
			if err != nil {
				log.Printf("Failed to load role: %v", err)
				return entydad.HTMXError("Failed to load role")
			}

			// Build set of already-assigned permission IDs
			assignedSet := make(map[string]bool)
			for _, rp := range roleResp.GetRole().GetRolePermissions() {
				assignedSet[rp.GetPermissionId()] = true
			}

			// Build dropdown options excluding already-assigned.
			// Include the permission code in the label so the colon-notation is visible.
			options := []types.SelectOption{}
			for _, p := range permResp.GetData() {
				if !p.GetActive() {
					continue
				}
				if assignedSet[p.GetId()] {
					continue
				}
				label := p.GetName()
				if code := p.GetPermissionCode(); code != "" {
					label = label + " (" + code + ")"
				}
				options = append(options, types.SelectOption{
					Value: p.GetId(),
					Label: label,
				})
			}

			return view.OK("role-permission-assign-form", &AssignFormData{
				FormAction:        "/action/manage/roles/" + roleID + "/permissions/assign",
				RoleID:            roleID,
				Labels:            AssignFormLabels{Permission: deps.Labels.Form.Permission},
				PermissionOptions: options,
				CommonLabels:      nil,
			})
		}

		// POST -- assign permission to role
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		permissionID := viewCtx.Request.FormValue("permission_id")
		if permissionID == "" {
			return entydad.HTMXError("Permission is required")
		}

		_, err := deps.CreateRolePermission(ctx, &rolepermissionpb.CreateRolePermissionRequest{
			Data: &rolepermissionpb.RolePermission{
				RoleId:       roleID,
				PermissionId: permissionID,
				Active:       true,
			},
		})
		if err != nil {
			log.Printf("Failed to assign permission to role %s: %v", roleID, err)
			return entydad.HTMXError("Failed to assign permission")
		}

		return entydad.HTMXSuccess("role-permissions-table")
	})
}

// NewRemoveAction creates the remove permission action (POST only).
func NewRemoveAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		roleID := viewCtx.Request.PathValue("id")
		if roleID == "" {
			return entydad.HTMXError("Role ID is required")
		}

		rpID := viewCtx.Request.URL.Query().Get("id")
		if rpID == "" {
			_ = viewCtx.Request.ParseForm()
			rpID = viewCtx.Request.FormValue("id")
		}
		if rpID == "" {
			return entydad.HTMXError("Role-Permission ID is required")
		}

		_, err := deps.DeleteRolePermission(ctx, &rolepermissionpb.DeleteRolePermissionRequest{
			Data: &rolepermissionpb.RolePermission{Id: rpID},
		})
		if err != nil {
			log.Printf("Failed to remove permission %s from role %s: %v", rpID, roleID, err)
			return entydad.HTMXError("Failed to remove permission")
		}

		return entydad.HTMXSuccess("role-permissions-table")
	})
}

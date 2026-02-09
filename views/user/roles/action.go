package roles

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"

	"github.com/erniealice/entydad-golang"
)

// AssignFormLabels holds i18n labels for the assign role drawer form.
type AssignFormLabels struct {
	Role string
}

// AssignFormData is the template data for the assign role drawer form.
type AssignFormData struct {
	FormAction   string
	UserID       string
	Labels       AssignFormLabels
	RoleOptions  []types.SelectOption
	CommonLabels any
}

// ActionDeps holds dependencies for user-role action handlers.
type ActionDeps struct {
	CreateWorkspaceUserRole      func(ctx context.Context, req *workspaceuserrolepb.CreateWorkspaceUserRoleRequest) (*workspaceuserrolepb.CreateWorkspaceUserRoleResponse, error)
	DeleteWorkspaceUserRole      func(ctx context.Context, req *workspaceuserrolepb.DeleteWorkspaceUserRoleRequest) (*workspaceuserrolepb.DeleteWorkspaceUserRoleResponse, error)
	ListRoles                    func(ctx context.Context, req *rolepb.ListRolesRequest) (*rolepb.ListRolesResponse, error)
	ListWorkspaceUsers           func(ctx context.Context, req *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error)
	GetWorkspaceUserItemPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error)
	Labels                       entydad.UserRoleLabels
}

// NewAssignAction creates the assign role action (GET = form, POST = create).
func NewAssignAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		userID := viewCtx.Request.PathValue("id")
		if userID == "" {
			return entydad.HTMXError("User ID is required")
		}

		if viewCtx.Request.Method == http.MethodGet {
			// Load all roles for the dropdown
			roleResp, err := deps.ListRoles(ctx, &rolepb.ListRolesRequest{})
			if err != nil {
				log.Printf("Failed to list roles: %v", err)
				return entydad.HTMXError("Failed to load roles")
			}

			// Find workspace_user to get already-assigned roles
			wu, err := findWorkspaceUserForAction(ctx, deps, userID)
			assignedSet := make(map[string]bool)
			if err == nil && wu != nil {
				// Get full workspace_user with roles
				wuResp, err := deps.GetWorkspaceUserItemPageData(ctx, &workspaceuserpb.GetWorkspaceUserItemPageDataRequest{
					WorkspaceUserId: wu.GetId(),
				})
				if err == nil {
					for _, wur := range wuResp.GetWorkspaceUser().GetWorkspaceUserRoles() {
						assignedSet[wur.GetRoleId()] = true
					}
				}
			}

			// Build dropdown options excluding already-assigned
			options := []types.SelectOption{}
			for _, r := range roleResp.GetData() {
				if !r.GetActive() {
					continue
				}
				if assignedSet[r.GetId()] {
					continue
				}
				options = append(options, types.SelectOption{
					Value: r.GetId(),
					Label: r.GetName(),
				})
			}

			return view.OK("user-role-assign-form", &AssignFormData{
				FormAction:   "/action/manage/users/" + userID + "/roles/assign",
				UserID:       userID,
				Labels:       AssignFormLabels{Role: deps.Labels.Form.Role},
				RoleOptions:  options,
				CommonLabels: nil,
			})
		}

		// POST -- assign role to user
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		roleID := viewCtx.Request.FormValue("role_id")
		if roleID == "" {
			return entydad.HTMXError("Role is required")
		}

		// Find workspace_user for this user
		wu, err := findWorkspaceUserForAction(ctx, deps, userID)
		if err != nil {
			log.Printf("Failed to find workspace user for user %s: %v", userID, err)
			return entydad.HTMXError("Failed to find workspace membership for this user")
		}

		_, err = deps.CreateWorkspaceUserRole(ctx, &workspaceuserrolepb.CreateWorkspaceUserRoleRequest{
			Data: &workspaceuserrolepb.WorkspaceUserRole{
				WorkspaceUserId: wu.GetId(),
				RoleId:          roleID,
				Active:          true,
			},
		})
		if err != nil {
			log.Printf("Failed to assign role to user %s: %v", userID, err)
			return entydad.HTMXError("Failed to assign role")
		}

		return entydad.HTMXSuccess("user-roles-table")
	})
}

// NewRemoveAction creates the remove role action (POST only).
func NewRemoveAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		userID := viewCtx.Request.PathValue("id")
		if userID == "" {
			return entydad.HTMXError("User ID is required")
		}

		wurID := viewCtx.Request.URL.Query().Get("id")
		if wurID == "" {
			_ = viewCtx.Request.ParseForm()
			wurID = viewCtx.Request.FormValue("id")
		}
		if wurID == "" {
			return entydad.HTMXError("Workspace-User-Role ID is required")
		}

		_, err := deps.DeleteWorkspaceUserRole(ctx, &workspaceuserrolepb.DeleteWorkspaceUserRoleRequest{
			Data: &workspaceuserrolepb.WorkspaceUserRole{Id: wurID},
		})
		if err != nil {
			log.Printf("Failed to remove role %s from user %s: %v", wurID, userID, err)
			return entydad.HTMXError("Failed to remove role")
		}

		return entydad.HTMXSuccess("user-roles-table")
	})
}

// findWorkspaceUserForAction finds the workspace_user record for a given user ID.
func findWorkspaceUserForAction(ctx context.Context, deps *ActionDeps, userID string) (*workspaceuserpb.WorkspaceUser, error) {
	resp, err := deps.ListWorkspaceUsers(ctx, &workspaceuserpb.ListWorkspaceUsersRequest{})
	if err != nil {
		return nil, err
	}

	for _, wu := range resp.GetData() {
		if wu.GetUserId() == userID {
			return wu, nil
		}
	}

	return nil, fmt.Errorf("workspace user not found for user ID %s", userID)
}

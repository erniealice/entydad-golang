package users

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"

	"github.com/erniealice/entydad-golang"
)

// AssignFormLabels holds i18n labels for the assign user drawer form.
type AssignFormLabels struct {
	User string
}

// AssignFormData is the template data for the assign user drawer form.
type AssignFormData struct {
	FormAction   string
	RoleID       string
	Labels       AssignFormLabels
	UserOptions  []types.SelectOption
	CommonLabels any
}

// ActionDeps holds dependencies for role-user action handlers.
type ActionDeps struct {
	GetUsersByRoleID        func(ctx context.Context, roleID string) ([]UserByRole, error)
	ListWorkspaceUsers      func(ctx context.Context, req *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error)
	CreateWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.CreateWorkspaceUserRoleRequest) (*workspaceuserrolepb.CreateWorkspaceUserRoleResponse, error)
	DeleteWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.DeleteWorkspaceUserRoleRequest) (*workspaceuserrolepb.DeleteWorkspaceUserRoleResponse, error)
	Labels                  entydad.RoleUserLabels
}

// NewAssignAction creates the assign user action (GET = form, POST = create).
func NewAssignAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		roleID := viewCtx.Request.PathValue("id")
		if roleID == "" {
			return entydad.HTMXError("Role ID is required")
		}

		if viewCtx.Request.Method == http.MethodGet {
			// Load all workspace users for the dropdown
			wuResp, err := deps.ListWorkspaceUsers(ctx, &workspaceuserpb.ListWorkspaceUsersRequest{})
			if err != nil {
				log.Printf("Failed to list workspace users: %v", err)
				return entydad.HTMXError("Failed to load users")
			}

			// Get already-assigned users to filter them out
			assignedSet := make(map[string]bool)
			if deps.GetUsersByRoleID != nil {
				assigned, err := deps.GetUsersByRoleID(ctx, roleID)
				if err == nil {
					for _, u := range assigned {
						assignedSet[u.WorkspaceUserID] = true
					}
				}
			}

			// Build dropdown options excluding already-assigned
			options := []types.SelectOption{}
			for _, wu := range wuResp.GetData() {
				if !wu.GetActive() {
					continue
				}
				if assignedSet[wu.GetId()] {
					continue
				}
				user := wu.GetUser()
				if user == nil {
					continue
				}
				label := user.GetFirstName() + " " + user.GetLastName()
				if email := user.GetEmailAddress(); email != "" {
					label = label + " (" + email + ")"
				}
				options = append(options, types.SelectOption{
					Value: wu.GetId(),
					Label: label,
				})
			}

			return view.OK("role-user-assign-form", &AssignFormData{
				FormAction:   fmt.Sprintf("/action/roles/detail/%s/users/assign", roleID),
				RoleID:       roleID,
				Labels:       AssignFormLabels{User: deps.Labels.Form.User},
				UserOptions:  options,
				CommonLabels: nil,
			})
		}

		// POST -- assign user to role
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		workspaceUserID := viewCtx.Request.FormValue("workspace_user_id")
		if workspaceUserID == "" {
			return entydad.HTMXError("User is required")
		}

		_, err := deps.CreateWorkspaceUserRole(ctx, &workspaceuserrolepb.CreateWorkspaceUserRoleRequest{
			Data: &workspaceuserrolepb.WorkspaceUserRole{
				WorkspaceUserId: workspaceUserID,
				RoleId:          roleID,
				Active:          true,
			},
		})
		if err != nil {
			log.Printf("Failed to assign user to role %s: %v", roleID, err)
			return entydad.HTMXError("Failed to assign user")
		}

		return entydad.HTMXSuccess("role-users-table")
	})
}

// NewRemoveAction creates the remove user action (POST only).
func NewRemoveAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		roleID := viewCtx.Request.PathValue("id")
		if roleID == "" {
			return entydad.HTMXError("Role ID is required")
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
			log.Printf("Failed to remove user %s from role %s: %v", wurID, roleID, err)
			return entydad.HTMXError("Failed to remove user")
		}

		return entydad.HTMXSuccess("role-users-table")
	})
}

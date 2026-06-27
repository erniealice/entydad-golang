// identity.go — block sub-context lift (B1, block-go-anatomy).
//
// wireIdentityModule registers the identity sub-context entity modules
// (user, role, permission, workspace, workspace_user, workspace_user_role)
// into the app router. It is a PURE code-move of the corresponding
// `if cfg.enableAll || cfg.X { ... }` blocks from block.go's Block() — same
// construction order, registration order, callbacks, and nil-checks. No
// behaviour change.
//
// All deps the lifted bodies need from Block()'s scope are carried on
// identityWiring (block-go-anatomy: >6 deps → struct).
package block

import (
	"context"
	"log"
	"net/http"

	entity "github.com/erniealice/entydad-golang/domain/entity"
	identity "github.com/erniealice/entydad-golang/domain/entity/identity"
	roleusers "github.com/erniealice/entydad-golang/domain/entity/identity/role/users"
	userdashboard "github.com/erniealice/entydad-golang/domain/entity/identity/user/dashboard"
	workspaceaction "github.com/erniealice/entydad-golang/domain/entity/identity/workspace/action"
	"github.com/erniealice/espyna-golang/consumer"
	consumerapp "github.com/erniealice/espyna-golang/consumer/app"
	"github.com/erniealice/espyna-golang/ports"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
)

// identityWiring carries everything the identity cluster needs from Block()'s
// scope. Implementation detail of the wiring; never re-exported.
type identityWiring struct {
	cfg        *blockConfig
	uc         *UseCases
	labels     blockLabels
	routes     blockRoutes
	refChecker ports.Checker

	getUserWorkspacesMap func(ctx context.Context) (map[string][]pyezatypes.ChipData, error)
	getDashboardData     func(ctx context.Context) (*userdashboard.DashboardData, error)
	hashPassword         func(password string) (string, error)
	getUsersByRoleID     func(ctx context.Context, roleID string) ([]roleusers.UserByRole, error)

	uploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	listAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID  func() string
}

func wireIdentityModule(ctx *consumerapp.AppContext, w identityWiring) {
	cfg := w.cfg
	uc := w.uc
	labels := w.labels
	routes := w.routes
	refChecker := w.refChecker
	getUserWorkspacesMap := w.getUserWorkspacesMap
	getDashboardData := w.getDashboardData
	hashPassword := w.hashPassword
	getUsersByRoleID := w.getUsersByRoleID
	uploadFile := w.uploadFile
	listAttachments := w.listAttachments
	createAttachment := w.createAttachment
	deleteAttachment := w.deleteAttachment
	newAttachmentID := w.newAttachmentID

	// WS-4: capability closure sourced from the auth adapter (NOT a proto use case).
	var getUserAuthCapability func(ctx context.Context, userID string) (bool, []string, error)
	if aa, ok := ctx.AuthAdapter.(*consumer.AuthAdapter); ok && aa != nil {
		getUserAuthCapability = func(ctx context.Context, userID string) (bool, []string, error) {
			c, err := aa.GetUserAuthCapability(ctx, userID)
			return c.HasPassword, c.Providers, err
		}
	}

	if cfg.enableAll || cfg.user {
		// WS-4 fail-closed boot guard (M-5). getUserAuthCapability is the
		// AUTHORITATIVE server-side control that rejects a local password reset
		// for an IdP-federated (SSO) user: NewResetPasswordAction runs the guard
		// ONLY when this closure is non-nil, and the detail page defaults
		// CanResetPasswordHere=true. A nil closure therefore fails OPEN (an SSO
		// user's reset would be allowed). The closure is nil only when
		// ctx.AuthAdapter is not a live *consumer.AuthAdapter (unwired / typed-nil
		// / wrong type). For the real identity providers (firebase / password)
		// that is a wiring fault — refuse to boot rather than silently run with the
		// SSO password-reset guard disabled. mock / noop / unset providers
		// legitimately have no capability source and stay fail-OPEN for dev.
		if getUserAuthCapability == nil {
			switch getEnv("CONFIG_AUTH_PROVIDER", "") {
			case "firebase", "password":
				log.Fatalf("FATAL entydad.Block: WS-4 boot guard: CONFIG_AUTH_PROVIDER=%q but the "+
					"GetUserAuthCapability closure is nil (ctx.AuthAdapter is %T, want a live "+
					"*consumer.AuthAdapter). Refusing to boot with the SSO password-reset guard "+
					"disabled — it would fail OPEN and allow resetting an IdP-managed user's password.",
					getEnv("CONFIG_AUTH_PROVIDER", ""), ctx.AuthAdapter)
			}
		}
		identity.NewUserModule(&identity.UserModuleDeps{
			Routes:                       routes.User,
			CommonLabels:                 ctx.Common,
			SharedLabels:                 labels.Shared,
			Labels:                       labels.User,
			DashboardLabels:              labels.UserDashboard,
			DashboardTitleLabels:         labels.Dashboard,
			UserRoleLabels:               labels.UserRole,
			TableLabels:                  ctx.Table,
			GetListPageData:              uc.User.GetListPageData,
			GetUserWorkspacesMap:         getUserWorkspacesMap,
			CreateUser:                   uc.User.Create,
			ReadUser:                     uc.User.Read,
			UpdateUser:                   uc.User.Update,
			DeleteUser:                   uc.User.Delete,
			SetActive:                    setActiveClosure(uc, "user"),
			DisableUser:                  uc.User.Disable,
			EnableUser:                   uc.User.Enable,
			AdminResetPassword:           uc.User.ResetPassword,
			GetUserAuthCapability:        getUserAuthCapability,
			CreateWorkspaceUser:          uc.WorkspaceUser.Create,
			ListWorkspaceUsers:           uc.WorkspaceUser.List,
			GetWorkspaceUserItemPageData: uc.WorkspaceUser.GetItemPageData,
			DefaultWorkspaceID:           getDefaultWorkspaceID(),
			CreateWorkspaceUserRole:      uc.WorkspaceUserRole.Create,
			DeleteWorkspaceUserRole:      uc.WorkspaceUserRole.Delete,
			ListRoles:                    uc.Role.List,
			GetDashboardData:             getDashboardData,
			HashPassword:                 hashPassword,
			UploadFile:                   uploadFile,
			ListAttachments:              listAttachments,
			CreateAttachment:             createAttachment,
			DeleteAttachment:             deleteAttachment,
			NewID:                        newAttachmentID,
		}).RegisterRoutes(ctx.Routes)
	}

	if cfg.enableAll || cfg.role {
		identity.NewRoleModule(&identity.RoleModuleDeps{
			Routes:                  routes.Role,
			CommonLabels:            ctx.Common,
			SharedLabels:            labels.Shared,
			Labels:                  labels.Role,
			RolePermissionLabels:    labels.RolePermission,
			RoleUserLabels:          labels.RoleUser,
			TableLabels:             ctx.Table,
			GetListPageData:         uc.Role.GetListPageData,
			GetInUseIDs:             refChecker.GetRoleInUseIDs,
			CreateRole:              uc.Role.Create,
			ReadRole:                uc.Role.Read,
			UpdateRole:              uc.Role.Update,
			DeleteRole:              uc.Role.Delete,
			SetActive:               setActiveClosure(uc, "role"),
			GetItemPageData:         uc.Role.GetItemPageData,
			CreateRolePermission:    uc.RolePermission.Create,
			DeleteRolePermission:    uc.RolePermission.Delete,
			ListPermissions:         uc.Permission.List,
			GetUsersByRoleID:        getUsersByRoleID,
			ListWorkspaceUsers:      uc.WorkspaceUser.List,
			CreateWorkspaceUserRole: uc.WorkspaceUserRole.Create,
			DeleteWorkspaceUserRole: uc.WorkspaceUserRole.Delete,
			UploadFile:              uploadFile,
			ListAttachments:         listAttachments,
			CreateAttachment:        createAttachment,
			DeleteAttachment:        deleteAttachment,
			NewID:                   newAttachmentID,
		}).RegisterRoutes(ctx.Routes)

		// Role-User search (http.HandlerFunc — uses HandleFunc, not GET)
		handleFunc(ctx.Routes, "GET", routes.Role.UsersSearchURL, roleusers.NewSearchUsersAction(&roleusers.SearchDeps{
			ListWorkspaceUsers: uc.WorkspaceUser.List,
		}))
	}

	if cfg.enableAll || cfg.permission {
		identity.NewPermissionModule(&identity.PermissionModuleDeps{
			Routes:           routes.Permission,
			CommonLabels:     ctx.Common,
			SharedLabels:     labels.Shared,
			Labels:           labels.Permission,
			TableLabels:      ctx.Table,
			GetListPageData:  uc.Permission.GetListPageData,
			CreatePermission: uc.Permission.Create,
			ReadPermission:   uc.Permission.Read,
			UpdatePermission: uc.Permission.Update,
			DeletePermission: uc.Permission.Delete,
			SetActive:        setActiveClosure(uc, "permission"),
		}).RegisterRoutes(ctx.Routes)
	}

	if cfg.enableAll || cfg.workspace {
		wsMod := &identity.WorkspaceModuleDeps{
			Routes:          routes.Workspace,
			CommonLabels:    ctx.Common,
			SharedLabels:    labels.Shared,
			Labels:          labels.Workspace,
			TableLabels:     ctx.Table,
			GetListPageData: uc.Workspace.GetListPageData,
			CreateWorkspace: uc.Workspace.Create,
			ReadWorkspace:   uc.Workspace.Read,
			UpdateWorkspace: uc.Workspace.Update,
			DeleteWorkspace: uc.Workspace.Delete,
			SetActive:       setActiveClosure(uc, "workspace"),
			// Phase 2 TODO closeout: wire the workspace_user detail + add URLs
			// now that Phase 2 has registered those route constants.
			WorkspaceUserDetailURL: entity.WorkspaceUserDetailURL,
			WorkspaceUserAddURL:    entity.WorkspaceUserAddURL,
			UploadFile:             uploadFile,
			ListAttachments:        listAttachments,
			CreateAttachment:       createAttachment,
			DeleteAttachment:       deleteAttachment,
			NewID:                  newAttachmentID,
		}
		if uc.WorkspaceUser.GetListPageData != nil {
			wsMod.GetWorkspaceUserListPageData = uc.WorkspaceUser.GetListPageData
		}
		identity.NewWorkspaceModule(wsMod).RegisterRoutes(ctx.Routes)

		// Switch workspace (raw POST — uses session cookie, issues HX-Redirect)
		// Registers when EITHER the legacy SwitchWorkspace use case is wired
		// OR the host app has provided a SecureSwitch override (A1 fix
		// WKR-P0-1: service-admin wires SecureSwitch so the sidebar
		// workspace-switcher rotates + audits via executePrincipalSwitch).
		//
		// Two wire-up paths supported:
		//   1. BlockOption: WithSecureSwitch(fn, resolveUserID, setCookie)
		//      — explicit, used when host can construct closures before
		//      Block() is applied.
		//   2. AppContext fields: ctx.SecureWorkspaceSwitch + the two
		//      sibling fields — used when the host has the appBuilder
		//      ready inside buildAppContext() but constructs the entydad
		//      AppOption from a different call site. This lets
		//      service-admin populate them inside buildAppContext()
		//      without restructuring entydadBlock().
		secureSwitch := cfg.secureSwitch
		secureSwitchResolveUser := cfg.secureSwitchResolveUser
		secureSwitchSetCookie := cfg.secureSwitchSetCookie
		if secureSwitch == nil {
			if v, ok := ctx.SecureWorkspaceSwitch.(workspaceaction.SecureSwitchFn); ok {
				secureSwitch = v
			}
			if v, ok := ctx.SecureWorkspaceSwitchResolveUserID.(func(r *http.Request) string); ok {
				secureSwitchResolveUser = v
			}
			if v, ok := ctx.SecureWorkspaceSwitchSetSessionCookie.(func(w http.ResponseWriter, token string)); ok {
				secureSwitchSetCookie = v
			}
		}
		if uc.Workspace.Switch != nil || secureSwitch != nil {
			handleFunc(ctx.Routes, "POST", routes.Workspace.SwitchURL, workspaceaction.NewSwitchWorkspaceHandler(&workspaceaction.SwitchWorkspaceDeps{
				SecureSwitch:          secureSwitch,
				ResolveUserID:         secureSwitchResolveUser,
				SetSessionCookie:      secureSwitchSetCookie,
				SwitchWorkspace:       uc.Workspace.Switch,
				HomeURLForWorkspaceID: cfg.homeURLForWorkspaceID,
				HomeURL:               cfg.homeURL,
			}))
		}
	}

	if cfg.enableAll || cfg.workspaceUser {
		if uc.WorkspaceUser.GetListPageData == nil {
			log.Println("entydad.Block: warning: workspace_user use cases not initialized — workspace_user detail routes will be unavailable")
		} else {
			wuRoutes := routes.WorkspaceUser
			wuMod := &identity.WorkspaceUserModuleDeps{
				Routes:                       wuRoutes,
				WorkspaceDetailURL:           entity.WorkspaceDetailURL,
				CommonLabels:                 ctx.Common,
				Labels:                       labels.WorkspaceUser,
				TableLabels:                  ctx.Table,
				GetListPageData:              uc.WorkspaceUser.GetListPageData,
				GetWorkspaceUserItemPageData: uc.WorkspaceUser.GetItemPageData,
				CreateWorkspaceUser:          uc.WorkspaceUser.Create,
				DeleteWorkspaceUser:          uc.WorkspaceUser.Delete,
				SetWorkspaceUserActive:       setActiveClosure(uc, "workspace_user"),
				// Phase 3 closeout: wire WorkspaceUserRole routes now that Phase 3 has registered them.
				WorkspaceUserRoleAddURL:    entity.WorkspaceUserRoleAddURL,
				WorkspaceUserRoleDeleteURL: entity.WorkspaceUserRoleDeleteURL,
				UploadFile:                 uploadFile,
				ListAttachments:            listAttachments,
				CreateAttachment:           createAttachment,
				DeleteAttachment:           deleteAttachment,
				NewID:                      newAttachmentID,
			}
			// ListUsers — needed for the user-search autocomplete on the add form.
			if uc.User.List != nil {
				wuMod.ListUsers = uc.User.List
			}
			// Phase 3 closeout: wire workspace_user_role list page data.
			if uc.WorkspaceUserRole.GetListPageData != nil {
				wuMod.GetWorkspaceUserRoleListPageData = uc.WorkspaceUserRole.GetListPageData
			}
			identity.NewWorkspaceUserModule(wuMod).RegisterRoutes(ctx.Routes)
			log.Println("  ✓ WorkspaceUser module initialized (entydad.Block)")
		}
	}

	if cfg.enableAll || cfg.workspaceUserRole {
		if uc.WorkspaceUserRole.Create == nil {
			log.Println("entydad.Block: warning: workspace_user_role use cases not initialized — workspace_user_role drawer routes will be unavailable")
		} else {
			wurRoutes := routes.WorkspaceUserRole
			wurMod := &identity.WorkspaceUserRoleModuleDeps{
				Routes:                  wurRoutes,
				Labels:                  labels.WorkspaceUserRole,
				CommonLabels:            ctx.Common,
				CreateWorkspaceUserRole: uc.WorkspaceUserRole.Create,
				DeleteWorkspaceUserRole: uc.WorkspaceUserRole.Delete,
			}
			if uc.WorkspaceUser.GetItemPageData != nil {
				wurMod.GetWorkspaceUserItemPageData = uc.WorkspaceUser.GetItemPageData
			}
			if uc.Role.List != nil {
				wurMod.ListRoles = uc.Role.List
			}
			identity.NewWorkspaceUserRoleModule(wurMod).RegisterRoutes(ctx.Routes)
			log.Println("  ✓ WorkspaceUserRole module initialized (entydad.Block)")
		}
	}
}

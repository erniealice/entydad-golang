package identity

import (
	"context"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	user "github.com/erniealice/entydad-golang/domain/entity/identity/user"
	useraction "github.com/erniealice/entydad-golang/domain/entity/identity/user/action"
	userdashboard "github.com/erniealice/entydad-golang/domain/entity/identity/user/dashboard"
	userdetail "github.com/erniealice/entydad-golang/domain/entity/identity/user/detail"
	userlist "github.com/erniealice/entydad-golang/domain/entity/identity/user/list"
	userroles "github.com/erniealice/entydad-golang/domain/entity/identity/user/roles"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
)

// UserModuleDeps holds all dependencies for the user module.
type UserModuleDeps struct {
	Routes               user.Routes
	CommonLabels         pyeza.CommonLabels
	SharedLabels         entydad.SharedLabels
	Labels               user.Labels
	DashboardLabels      user.DashboardLabels
	DashboardTitleLabels entydad.DashboardLabels
	UserRoleLabels       user.RoleLabels
	TableLabels          types.TableLabels
	// User list page data
	GetListPageData      func(ctx context.Context, req *userpb.GetUserListPageDataRequest) (*userpb.GetUserListPageDataResponse, error)
	GetUserWorkspacesMap func(ctx context.Context) (map[string][]types.ChipData, error)
	// User CRUD
	CreateUser func(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error)
	ReadUser   func(ctx context.Context, req *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error)
	UpdateUser func(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error)
	DeleteUser func(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error)
	SetActive  func(ctx context.Context, id string, active bool) error
	// Workspace user (for user creation + detail)
	CreateWorkspaceUser          func(ctx context.Context, req *workspaceuserpb.CreateWorkspaceUserRequest) (*workspaceuserpb.CreateWorkspaceUserResponse, error)
	ListWorkspaceUsers           func(ctx context.Context, req *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error)
	GetWorkspaceUserItemPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error)
	DefaultWorkspaceID           string
	// User-Role assignment
	CreateWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.CreateWorkspaceUserRoleRequest) (*workspaceuserrolepb.CreateWorkspaceUserRoleResponse, error)
	DeleteWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.DeleteWorkspaceUserRoleRequest) (*workspaceuserrolepb.DeleteWorkspaceUserRoleResponse, error)
	ListRoles               func(ctx context.Context, req *rolepb.ListRolesRequest) (*rolepb.ListRolesResponse, error)
	// Dashboard
	GetDashboardData func(ctx context.Context) (*userdashboard.DashboardData, error)
	// Password hashing (optional)
	HashPassword func(password string) (string, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Audit history
	ListAuditHistory func(ctx context.Context, req *auditlog.ListAuditRequest) (*auditlog.ListAuditResponse, error)
}

// UserModule holds all constructed user views.
type UserModule struct {
	routes        user.Routes
	Dashboard     view.View
	List          view.View
	Table         view.View
	Detail        view.View
	TabAction     view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	BulkDelete    view.View
	SetStatus     view.View
	BulkSetStatus view.View
	ResetPassword view.View
	// User-Role assignment views (detail + legacy paths)
	RoleList         view.View
	RoleTable        view.View
	RoleAssign       view.View
	RoleRemove       view.View
	AttachmentUpload view.View
	AttachmentDelete view.View
	// SearchTimezones is a JSON endpoint backing the timezone autocomplete in the user drawer form.
	SearchTimezones http.HandlerFunc
}

func NewUserModule(deps *UserModuleDeps) *UserModule {
	actionDeps := &useraction.Deps{
		Routes:              deps.Routes,
		CreateUser:          deps.CreateUser,
		ReadUser:            deps.ReadUser,
		UpdateUser:          deps.UpdateUser,
		DeleteUser:          deps.DeleteUser,
		SetUserActive:       deps.SetActive,
		CreateWorkspaceUser: deps.CreateWorkspaceUser,
		DefaultWorkspaceID:  deps.DefaultWorkspaceID,
		HashPassword:        deps.HashPassword,
	}
	listDeps := &userlist.ListViewDeps{
		Routes:               deps.Routes,
		GetListPageData:      deps.GetListPageData,
		GetUserWorkspacesMap: deps.GetUserWorkspacesMap,
		RefreshURL:           deps.Routes.TableURL,
		Labels:               deps.Labels,
		SharedLabels:         deps.SharedLabels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
	}
	detailDeps := &userdetail.DetailViewDeps{
		Routes:                       deps.Routes,
		ReadUser:                     deps.ReadUser,
		GetWorkspaceUserItemPageData: deps.GetWorkspaceUserItemPageData,
		ListWorkspaceUsers:           deps.ListWorkspaceUsers,
		Labels:                       deps.Labels,
		SharedLabels:                 deps.SharedLabels,
		UserRoleLabels:               deps.UserRoleLabels,
		CommonLabels:                 deps.CommonLabels,
		TableLabels:                  deps.TableLabels,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
		AuditOps: auditlog.AuditOps{
			ListAuditHistory: deps.ListAuditHistory,
		},
	}
	roleListDeps := &userroles.Deps{
		Routes:                       deps.Routes,
		ListWorkspaceUsers:           deps.ListWorkspaceUsers,
		GetWorkspaceUserItemPageData: deps.GetWorkspaceUserItemPageData,
		ReadUser:                     deps.ReadUser,
		Labels:                       deps.UserRoleLabels,
		SharedLabels:                 deps.SharedLabels,
		CommonLabels:                 deps.CommonLabels,
		TableLabels:                  deps.TableLabels,
	}
	roleActionDeps := &userroles.ActionDeps{
		Routes:                       deps.Routes,
		CreateWorkspaceUserRole:      deps.CreateWorkspaceUserRole,
		DeleteWorkspaceUserRole:      deps.DeleteWorkspaceUserRole,
		ListRoles:                    deps.ListRoles,
		ListWorkspaceUsers:           deps.ListWorkspaceUsers,
		GetWorkspaceUserItemPageData: deps.GetWorkspaceUserItemPageData,
		CreateWorkspaceUser:          deps.CreateWorkspaceUser, // NEW
		DefaultWorkspaceID:           deps.DefaultWorkspaceID,  // NEW
		Labels:                       deps.UserRoleLabels,
	}

	return &UserModule{
		routes: deps.Routes,
		Dashboard: userdashboard.NewView(&userdashboard.Deps{
			DashboardLabels:  deps.DashboardTitleLabels,
			Dashboard:        deps.DashboardLabels,
			Routes:           deps.Routes,
			CommonLabels:     deps.CommonLabels,
			GetDashboardData: deps.GetDashboardData,
		}),
		List:             userlist.NewView(listDeps),
		Table:            userlist.NewTableView(listDeps),
		Detail:           userdetail.NewView(detailDeps),
		TabAction:        userdetail.NewTabAction(detailDeps),
		Add:              useraction.NewAddAction(actionDeps),
		Edit:             useraction.NewEditAction(actionDeps),
		Delete:           useraction.NewDeleteAction(actionDeps),
		BulkDelete:       useraction.NewBulkDeleteAction(actionDeps),
		SetStatus:        useraction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    useraction.NewBulkSetStatusAction(actionDeps),
		ResetPassword:    useraction.NewResetPasswordAction(actionDeps),
		RoleList:         userroles.NewView(roleListDeps),
		RoleTable:        userroles.NewTableView(roleListDeps),
		RoleAssign:       userroles.NewAssignAction(roleActionDeps),
		RoleRemove:       userroles.NewRemoveAction(roleActionDeps),
		AttachmentUpload: userdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: userdetail.NewAttachmentDeleteAction(detailDeps),
		SearchTimezones:  useraction.NewSearchTimezonesAction(),
	}
}

func (m *UserModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	r.POST(m.routes.ResetPasswordURL, m.ResetPassword)
	// User-Role assignment (/detail/ path)
	r.GET(m.routes.DetailRolesURL, m.RoleList)
	r.GET(m.routes.DetailRolesTableURL, m.RoleTable)
	r.GET(m.routes.DetailRolesAssignURL, m.RoleAssign)
	r.POST(m.routes.DetailRolesAssignURL, m.RoleAssign)
	r.POST(m.routes.DetailRolesRemoveURL, m.RoleRemove)
	// User-Role assignment (legacy /manage/ path)
	r.GET(m.routes.RolesURL, m.RoleList)
	r.GET(m.routes.RolesTableURL, m.RoleTable)
	r.GET(m.routes.RolesAssignURL, m.RoleAssign)
	r.POST(m.routes.RolesAssignURL, m.RoleAssign)
	r.POST(m.routes.RolesRemoveURL, m.RoleRemove)
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
	// Timezone autocomplete JSON endpoint
	identityHandleFunc(r, "GET", m.routes.SearchTimezonesURL, m.SearchTimezones)
}

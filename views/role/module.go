package role

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	roleaction "github.com/erniealice/entydad-golang/views/role/action"
	roledetail "github.com/erniealice/entydad-golang/views/role/detail"
	rolelist "github.com/erniealice/entydad-golang/views/role/list"
	rolepermissions "github.com/erniealice/entydad-golang/views/role/permissions"
	roleusers "github.com/erniealice/entydad-golang/views/role/users"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	rolepermissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role_permission"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
)

// ModuleDeps holds all dependencies for the role module.
type ModuleDeps struct {
	Routes               entydad.RoleRoutes
	CommonLabels         pyeza.CommonLabels
	SharedLabels         entydad.SharedLabels
	Labels               entydad.RoleLabels
	RolePermissionLabels entydad.RolePermissionLabels
	RoleUserLabels       entydad.RoleUserLabels
	TableLabels          types.TableLabels
	// Role list page data
	GetListPageData func(ctx context.Context, req *rolepb.GetRoleListPageDataRequest) (*rolepb.GetRoleListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	// Role CRUD
	CreateRole func(ctx context.Context, req *rolepb.CreateRoleRequest) (*rolepb.CreateRoleResponse, error)
	ReadRole   func(ctx context.Context, req *rolepb.ReadRoleRequest) (*rolepb.ReadRoleResponse, error)
	UpdateRole func(ctx context.Context, req *rolepb.UpdateRoleRequest) (*rolepb.UpdateRoleResponse, error)
	DeleteRole func(ctx context.Context, req *rolepb.DeleteRoleRequest) (*rolepb.DeleteRoleResponse, error)
	SetActive  func(ctx context.Context, id string, active bool) error
	// Role detail (item page data for permission counts)
	GetItemPageData func(ctx context.Context, req *rolepb.GetRoleItemPageDataRequest) (*rolepb.GetRoleItemPageDataResponse, error)
	// Role-Permission assignment
	CreateRolePermission func(ctx context.Context, req *rolepermissionpb.CreateRolePermissionRequest) (*rolepermissionpb.CreateRolePermissionResponse, error)
	DeleteRolePermission func(ctx context.Context, req *rolepermissionpb.DeleteRolePermissionRequest) (*rolepermissionpb.DeleteRolePermissionResponse, error)
	ListPermissions      func(ctx context.Context, req *permissionpb.ListPermissionsRequest) (*permissionpb.ListPermissionsResponse, error)
	// Role-User management
	GetUsersByRoleID        func(ctx context.Context, roleID string) ([]roleusers.UserByRole, error)
	ListWorkspaceUsers      func(ctx context.Context, req *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error)
	CreateWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.CreateWorkspaceUserRoleRequest) (*workspaceuserrolepb.CreateWorkspaceUserRoleResponse, error)
	DeleteWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.DeleteWorkspaceUserRoleRequest) (*workspaceuserrolepb.DeleteWorkspaceUserRoleResponse, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Audit history
	ListAuditHistory func(ctx context.Context, req *auditlog.ListAuditRequest) (*auditlog.ListAuditResponse, error)
}

// Module holds all constructed role views.
type Module struct {
	routes        entydad.RoleRoutes
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
	// Role-Permission assignment views
	PermissionList   view.View
	PermissionTable  view.View
	PermissionAssign view.View
	PermissionRemove view.View
	// Role-User assignment views
	UserList         view.View
	UserTable        view.View
	UserAssign       view.View
	UserRemove       view.View
	AttachmentUpload view.View
	AttachmentDelete view.View
}

func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &roleaction.Deps{
		CreateRole:    deps.CreateRole,
		ReadRole:      deps.ReadRole,
		UpdateRole:    deps.UpdateRole,
		DeleteRole:    deps.DeleteRole,
		SetRoleActive: deps.SetActive,
		Routes:        deps.Routes,
	}
	listDeps := &rolelist.ListViewDeps{
		GetListPageData: deps.GetListPageData,
		GetInUseIDs:     deps.GetInUseIDs,
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	}
	detailDeps := &roledetail.DetailViewDeps{
		ReadRole:             deps.ReadRole,
		RoleGetItemPageData:  deps.GetItemPageData,
		GetUsersByRoleID:     deps.GetUsersByRoleID,
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		SharedLabels:         deps.SharedLabels,
		RolePermissionLabels: deps.RolePermissionLabels,
		RoleUserLabels:       deps.RoleUserLabels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
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
	permListDeps := &rolepermissions.Deps{
		GetRoleItemPageData: deps.GetItemPageData,
		Routes:              deps.Routes,
		Labels:              deps.RolePermissionLabels,
		SharedLabels:        deps.SharedLabels,
		CommonLabels:        deps.CommonLabels,
		TableLabels:         deps.TableLabels,
	}
	permActionDeps := &rolepermissions.ActionDeps{
		CreateRolePermission: deps.CreateRolePermission,
		DeleteRolePermission: deps.DeleteRolePermission,
		ListPermissions:      deps.ListPermissions,
		GetRoleItemPageData:  deps.GetItemPageData,
		Routes:               deps.Routes,
		Labels:               deps.RolePermissionLabels,
	}
	userListDeps := &roleusers.Deps{
		GetUsersByRoleID: deps.GetUsersByRoleID,
		ReadRole:         deps.ReadRole,
		Routes:           deps.Routes,
		Labels:           deps.RoleUserLabels,
		SharedLabels:     deps.SharedLabels,
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
	}
	userActionDeps := &roleusers.ActionDeps{
		GetUsersByRoleID:        deps.GetUsersByRoleID,
		ListWorkspaceUsers:      deps.ListWorkspaceUsers,
		CreateWorkspaceUserRole: deps.CreateWorkspaceUserRole,
		DeleteWorkspaceUserRole: deps.DeleteWorkspaceUserRole,
		Routes:                  deps.Routes,
		Labels:                  deps.RoleUserLabels,
	}

	return &Module{
		routes:           deps.Routes,
		List:             rolelist.NewView(listDeps),
		Table:            rolelist.NewTableView(listDeps),
		Detail:           roledetail.NewView(detailDeps),
		TabAction:        roledetail.NewTabAction(detailDeps),
		Add:              roleaction.NewAddAction(actionDeps),
		Edit:             roleaction.NewEditAction(actionDeps),
		Delete:           roleaction.NewDeleteAction(actionDeps),
		BulkDelete:       roleaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        roleaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    roleaction.NewBulkSetStatusAction(actionDeps),
		PermissionList:   rolepermissions.NewView(permListDeps),
		PermissionTable:  rolepermissions.NewTableView(permListDeps),
		PermissionAssign: rolepermissions.NewAssignAction(permActionDeps),
		PermissionRemove: rolepermissions.NewRemoveAction(permActionDeps),
		UserList:         roleusers.NewView(userListDeps),
		UserTable:        roleusers.NewTableView(userListDeps),
		UserAssign:       roleusers.NewAssignAction(userActionDeps),
		UserRemove:       roleusers.NewRemoveAction(userActionDeps),
		AttachmentUpload: roledetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: roledetail.NewAttachmentDeleteAction(detailDeps),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
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
	// Role-Permission assignment (/detail/ path)
	r.GET(m.routes.DetailPermissionsURL, m.PermissionList)
	r.GET(m.routes.DetailPermissionsTableURL, m.PermissionTable)
	r.GET(m.routes.DetailPermissionsAssignURL, m.PermissionAssign)
	r.POST(m.routes.DetailPermissionsAssignURL, m.PermissionAssign)
	r.POST(m.routes.DetailPermissionsRemoveURL, m.PermissionRemove)
	// Role-Permission assignment (legacy /manage/ path)
	r.GET(m.routes.PermissionsURL, m.PermissionList)
	r.GET(m.routes.PermissionsTableURL, m.PermissionTable)
	r.GET(m.routes.PermissionsAssignURL, m.PermissionAssign)
	r.POST(m.routes.PermissionsAssignURL, m.PermissionAssign)
	r.POST(m.routes.PermissionsRemoveURL, m.PermissionRemove)
	// Role-User assignment
	r.GET(m.routes.UsersURL, m.UserList)
	r.GET(m.routes.UsersTableURL, m.UserTable)
	r.GET(m.routes.UsersAssignURL, m.UserAssign)
	r.POST(m.routes.UsersAssignURL, m.UserAssign)
	r.POST(m.routes.UsersRemoveURL, m.UserRemove)
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}

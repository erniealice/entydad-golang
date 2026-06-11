// workspace_user_module.go provides the view module for workspace_user nested detail.
// This is Phase 2 of the bootstrap-auth plan. The workspace_user surface is
// accessed by clicking a row in workspace detail's Users tab.
package identity

import (
	"context"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	workspaceuser "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user"
	workspaceuseraction "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user/action"
	workspaceuserdetail "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user/detail"
	workspaceuserlist "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user/list"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
	"github.com/erniealice/hybra-golang/views/attachment"
)

// WorkspaceUserModuleDeps holds all dependencies for the workspace_user module.
type WorkspaceUserModuleDeps struct {
	Routes                       workspaceuser.Routes
	WorkspaceDetailURL           string // /app/workspaces/detail/{id} — for "Back to workspace" link
	CommonLabels                 pyeza.CommonLabels
	Labels                       workspaceuser.Labels
	TableLabels                  types.TableLabels
	GetListPageData              func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserListPageDataRequest) (*workspaceuserpb.GetWorkspaceUserListPageDataResponse, error)
	GetWorkspaceUserItemPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error)
	CreateWorkspaceUser          func(ctx context.Context, req *workspaceuserpb.CreateWorkspaceUserRequest) (*workspaceuserpb.CreateWorkspaceUserResponse, error)
	DeleteWorkspaceUser          func(ctx context.Context, req *workspaceuserpb.DeleteWorkspaceUserRequest) (*workspaceuserpb.DeleteWorkspaceUserResponse, error)
	SetWorkspaceUserActive       func(ctx context.Context, id string, active bool) error
	// ListUsers is used by the user search autocomplete on the add form.
	ListUsers func(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error)

	// Phase 3 wired: GetWorkspaceUserRoleListPageData, WorkspaceUserRoleAddURL, WorkspaceUserRoleDeleteURL
	// are supplied by block.go after Phase 3 registered the workspace_user_role routes.
	GetWorkspaceUserRoleListPageData func(ctx context.Context, req *workspaceuserrolepb.GetWorkspaceUserRoleListPageDataRequest) (*workspaceuserrolepb.GetWorkspaceUserRoleListPageDataResponse, error)
	WorkspaceUserRoleAddURL          string
	WorkspaceUserRoleDeleteURL       string

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string
}

// WorkspaceUserModule holds all constructed workspace_user views.
type WorkspaceUserModule struct {
	routes    workspaceuser.Routes
	List      view.View
	Detail    view.View
	TabAction view.View
	Add       view.View
	Delete    view.View
	SetStatus view.View
	// UserSearch is an http.HandlerFunc for the user autocomplete endpoint.
	UserSearch       http.HandlerFunc
	AttachmentUpload view.View
	AttachmentDelete view.View
}

// NewWorkspaceUserModule constructs all workspace_user views from deps.
func NewWorkspaceUserModule(deps *WorkspaceUserModuleDeps) *WorkspaceUserModule {
	actionDeps := &workspaceuseraction.Deps{
		Routes:                 deps.Routes,
		CreateWorkspaceUser:    deps.CreateWorkspaceUser,
		DeleteWorkspaceUser:    deps.DeleteWorkspaceUser,
		SetWorkspaceUserActive: deps.SetWorkspaceUserActive,
		ListUsers:              deps.ListUsers,
	}
	listDeps := &workspaceuserlist.ListViewDeps{
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
		GetListPageData: deps.GetListPageData,
	}
	detailDeps := &workspaceuserdetail.DetailViewDeps{
		Routes:                           deps.Routes,
		WorkspaceDetailURL:               deps.WorkspaceDetailURL,
		GetWorkspaceUserItemPageData:     deps.GetWorkspaceUserItemPageData,
		Labels:                           deps.Labels,
		CommonLabels:                     deps.CommonLabels,
		TableLabels:                      deps.TableLabels,
		GetWorkspaceUserRoleListPageData: deps.GetWorkspaceUserRoleListPageData,
		WorkspaceUserRoleAddURL:          deps.WorkspaceUserRoleAddURL,
		WorkspaceUserRoleDeleteURL:       deps.WorkspaceUserRoleDeleteURL,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
	}

	m := &WorkspaceUserModule{
		routes:     deps.Routes,
		List:       workspaceuserlist.NewView(listDeps),
		Detail:     workspaceuserdetail.NewView(detailDeps),
		TabAction:  workspaceuserdetail.NewTabAction(detailDeps),
		Add:        workspaceuseraction.NewAddAction(actionDeps),
		Delete:     workspaceuseraction.NewDeleteAction(actionDeps),
		SetStatus:  workspaceuseraction.NewSetStatusAction(actionDeps),
		UserSearch: workspaceuseraction.NewUserSearchAction(actionDeps),
	}
	if deps.UploadFile != nil {
		m.AttachmentUpload = workspaceuserdetail.NewAttachmentUploadAction(detailDeps)
		m.AttachmentDelete = workspaceuserdetail.NewAttachmentDeleteAction(detailDeps)
	}
	return m
}

// RegisterRoutes registers all workspace_user routes into the app router.
func (m *WorkspaceUserModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.routes.ListURL != "" {
		r.GET(m.routes.ListURL, m.List)
	}
	if m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.routes.TabActionURL != "" {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
	if m.routes.AddURL != "" {
		r.GET(m.routes.AddURL, m.Add)
		r.POST(m.routes.AddURL, m.Add)
	}
	if m.routes.DeleteURL != "" {
		r.POST(m.routes.DeleteURL, m.Delete)
	}
	if m.routes.SetStatusURL != "" {
		r.POST(m.routes.SetStatusURL, m.SetStatus)
	}
	// User search is a raw HTTP handler (returns JSON)
	if m.routes.SearchURL != "" && m.UserSearch != nil {
		if full, ok := r.(identityRouteRegistrarFull); ok {
			full.HandleFunc("GET", m.routes.SearchURL, m.UserSearch)
		}
	}
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}

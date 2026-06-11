package identity

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	workspace "github.com/erniealice/entydad-golang/domain/entity/identity/workspace"
	workspaceaction "github.com/erniealice/entydad-golang/domain/entity/identity/workspace/action"
	workspacedetail "github.com/erniealice/entydad-golang/domain/entity/identity/workspace/detail"
	workspacelist "github.com/erniealice/entydad-golang/domain/entity/identity/workspace/list"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	"github.com/erniealice/hybra-golang/views/attachment"
)

// WorkspaceModuleDeps holds all dependencies for the workspace module.
type WorkspaceModuleDeps struct {
	Routes          workspace.Routes
	CommonLabels    pyeza.CommonLabels
	SharedLabels    entydad.SharedLabels
	Labels          workspace.Labels
	TableLabels     types.TableLabels
	GetListPageData func(ctx context.Context, req *workspacepb.GetWorkspaceListPageDataRequest) (*workspacepb.GetWorkspaceListPageDataResponse, error)
	CreateWorkspace func(ctx context.Context, req *workspacepb.CreateWorkspaceRequest) (*workspacepb.CreateWorkspaceResponse, error)
	ReadWorkspace   func(ctx context.Context, req *workspacepb.ReadWorkspaceRequest) (*workspacepb.ReadWorkspaceResponse, error)
	UpdateWorkspace func(ctx context.Context, req *workspacepb.UpdateWorkspaceRequest) (*workspacepb.UpdateWorkspaceResponse, error)
	DeleteWorkspace func(ctx context.Context, req *workspacepb.DeleteWorkspaceRequest) (*workspacepb.DeleteWorkspaceResponse, error)
	SetActive       func(ctx context.Context, id string, active bool) error

	// Detail page dependencies (Phase 1 additions).
	// Optional: when nil the detail page degrades gracefully (empty Users tab).
	GetWorkspaceUserListPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserListPageDataRequest) (*workspaceuserpb.GetWorkspaceUserListPageDataResponse, error)
	// WorkspaceUserDetailURL is the target for the "View" row action on each workspace_user row.
	// Phase 2 will register this route; Phase 1 just emits the URL in the table.
	WorkspaceUserDetailURL string
	// WorkspaceUserAddURL is the drawer action for "Add user to workspace".
	// Phase 2 will register this route; Phase 1 just emits the URL in the button.
	WorkspaceUserAddURL string

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string
}

// WorkspaceModule holds all constructed workspace views.
type WorkspaceModule struct {
	routes           workspace.Routes
	List             view.View
	Table            view.View
	Add              view.View
	Edit             view.View
	Delete           view.View
	BulkDelete       view.View
	SetStatus        view.View
	BulkSetStatus    view.View
	Detail           view.View
	TabAction        view.View
	AttachmentUpload view.View
	AttachmentDelete view.View
}

func NewWorkspaceModule(deps *WorkspaceModuleDeps) *WorkspaceModule {
	actionDeps := &workspaceaction.Deps{
		CreateWorkspace:    deps.CreateWorkspace,
		ReadWorkspace:      deps.ReadWorkspace,
		UpdateWorkspace:    deps.UpdateWorkspace,
		DeleteWorkspace:    deps.DeleteWorkspace,
		SetWorkspaceActive: deps.SetActive,
		Routes:             deps.Routes,
	}
	listDeps := &workspacelist.ListViewDeps{
		GetListPageData: deps.GetListPageData,
		RefreshURL:      deps.Routes.TableURL,
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	}
	detailDeps := &workspacedetail.DetailViewDeps{
		Routes:                       deps.Routes,
		ReadWorkspace:                deps.ReadWorkspace,
		GetWorkspaceUserListPageData: deps.GetWorkspaceUserListPageData,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		TableLabels:                  deps.TableLabels,
		WorkspaceUserDetailURL:       deps.WorkspaceUserDetailURL,
		WorkspaceUserAddURL:          deps.WorkspaceUserAddURL,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
	}

	m := &WorkspaceModule{
		routes:        deps.Routes,
		List:          workspacelist.NewView(listDeps),
		Table:         workspacelist.NewTableView(listDeps),
		Add:           workspaceaction.NewAddAction(actionDeps),
		Edit:          workspaceaction.NewEditAction(actionDeps),
		Delete:        workspaceaction.NewDeleteAction(actionDeps),
		BulkDelete:    workspaceaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     workspaceaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: workspaceaction.NewBulkSetStatusAction(actionDeps),
		Detail:        workspacedetail.NewView(detailDeps),
		TabAction:     workspacedetail.NewTabAction(detailDeps),
	}
	if deps.UploadFile != nil {
		m.AttachmentUpload = workspacedetail.NewAttachmentUploadAction(detailDeps)
		m.AttachmentDelete = workspacedetail.NewAttachmentDeleteAction(detailDeps)
	}
	return m
}

func (m *WorkspaceModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	if m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.routes.TabActionURL != "" {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}

package workspace

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	workspaceaction "github.com/erniealice/entydad-golang/views/workspace/action"
	workspacelist "github.com/erniealice/entydad-golang/views/workspace/list"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
)

// ModuleDeps holds all dependencies for the workspace module.
type ModuleDeps struct {
	Routes          entydad.WorkspaceRoutes
	CommonLabels    pyeza.CommonLabels
	SharedLabels    entydad.SharedLabels
	Labels          entydad.WorkspaceLabels
	TableLabels     types.TableLabels
	GetListPageData func(ctx context.Context, req *workspacepb.GetWorkspaceListPageDataRequest) (*workspacepb.GetWorkspaceListPageDataResponse, error)
	CreateWorkspace func(ctx context.Context, req *workspacepb.CreateWorkspaceRequest) (*workspacepb.CreateWorkspaceResponse, error)
	ReadWorkspace   func(ctx context.Context, req *workspacepb.ReadWorkspaceRequest) (*workspacepb.ReadWorkspaceResponse, error)
	UpdateWorkspace func(ctx context.Context, req *workspacepb.UpdateWorkspaceRequest) (*workspacepb.UpdateWorkspaceResponse, error)
	DeleteWorkspace func(ctx context.Context, req *workspacepb.DeleteWorkspaceRequest) (*workspacepb.DeleteWorkspaceResponse, error)
	SetActive       func(ctx context.Context, id string, active bool) error
}

// Module holds all constructed workspace views.
type Module struct {
	routes        entydad.WorkspaceRoutes
	List          view.View
	Table         view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	BulkDelete    view.View
	SetStatus     view.View
	BulkSetStatus view.View
}

func NewModule(deps *ModuleDeps) *Module {
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

	return &Module{
		routes:        deps.Routes,
		List:          workspacelist.NewView(listDeps),
		Table:         workspacelist.NewTableView(listDeps),
		Add:           workspaceaction.NewAddAction(actionDeps),
		Edit:          workspaceaction.NewEditAction(actionDeps),
		Delete:        workspaceaction.NewDeleteAction(actionDeps),
		BulkDelete:    workspaceaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     workspaceaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: workspaceaction.NewBulkSetStatusAction(actionDeps),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
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
}

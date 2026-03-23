package permission

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	permissionaction "github.com/erniealice/entydad-golang/views/permission/action"
	permissionlist "github.com/erniealice/entydad-golang/views/permission/list"
	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"
)

// ModuleDeps holds all dependencies for the permission module.
type ModuleDeps struct {
	Routes           entydad.PermissionRoutes
	CommonLabels     pyeza.CommonLabels
	SharedLabels     entydad.SharedLabels
	Labels           entydad.PermissionLabels
	TableLabels      types.TableLabels
	GetListPageData  func(ctx context.Context, req *permissionpb.GetPermissionListPageDataRequest) (*permissionpb.GetPermissionListPageDataResponse, error)
	CreatePermission func(ctx context.Context, req *permissionpb.CreatePermissionRequest) (*permissionpb.CreatePermissionResponse, error)
	ReadPermission   func(ctx context.Context, req *permissionpb.ReadPermissionRequest) (*permissionpb.ReadPermissionResponse, error)
	UpdatePermission func(ctx context.Context, req *permissionpb.UpdatePermissionRequest) (*permissionpb.UpdatePermissionResponse, error)
	DeletePermission func(ctx context.Context, req *permissionpb.DeletePermissionRequest) (*permissionpb.DeletePermissionResponse, error)
	SetActive        func(ctx context.Context, id string, active bool) error
}

// Module holds all constructed permission views.
type Module struct {
	routes        entydad.PermissionRoutes
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
	actionDeps := &permissionaction.Deps{
		CreatePermission:    deps.CreatePermission,
		ReadPermission:      deps.ReadPermission,
		UpdatePermission:    deps.UpdatePermission,
		DeletePermission:    deps.DeletePermission,
		SetPermissionActive: deps.SetActive,
		Routes:              deps.Routes,
	}
	listDeps := &permissionlist.ListViewDeps{
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
		List:          permissionlist.NewView(listDeps),
		Table:         permissionlist.NewTableView(listDeps),
		Add:           permissionaction.NewAddAction(actionDeps),
		Edit:          permissionaction.NewEditAction(actionDeps),
		Delete:        permissionaction.NewDeleteAction(actionDeps),
		BulkDelete:    permissionaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     permissionaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: permissionaction.NewBulkSetStatusAction(actionDeps),
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

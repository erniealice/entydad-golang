package identity

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	permission "github.com/erniealice/entydad-golang/domain/entity/identity/permission"
	permissionaction "github.com/erniealice/entydad-golang/domain/entity/identity/permission/action"
	permissionlist "github.com/erniealice/entydad-golang/domain/entity/identity/permission/list"
	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"
)

// PermissionModuleDeps holds all dependencies for the permission module.
type PermissionModuleDeps struct {
	Routes           permission.Routes
	CommonLabels     pyeza.CommonLabels
	SharedLabels     entydad.SharedLabels
	Labels           permission.Labels
	TableLabels      types.TableLabels
	GetListPageData  func(ctx context.Context, req *permissionpb.GetPermissionListPageDataRequest) (*permissionpb.GetPermissionListPageDataResponse, error)
	CreatePermission func(ctx context.Context, req *permissionpb.CreatePermissionRequest) (*permissionpb.CreatePermissionResponse, error)
	ReadPermission   func(ctx context.Context, req *permissionpb.ReadPermissionRequest) (*permissionpb.ReadPermissionResponse, error)
	UpdatePermission func(ctx context.Context, req *permissionpb.UpdatePermissionRequest) (*permissionpb.UpdatePermissionResponse, error)
	DeletePermission func(ctx context.Context, req *permissionpb.DeletePermissionRequest) (*permissionpb.DeletePermissionResponse, error)
	SetActive        func(ctx context.Context, id string, active bool) error
}

// PermissionModule holds all constructed permission views.
type PermissionModule struct {
	routes        permission.Routes
	List          view.View
	Table         view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	BulkDelete    view.View
	SetStatus     view.View
	BulkSetStatus view.View
}

func NewPermissionModule(deps *PermissionModuleDeps) *PermissionModule {
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

	return &PermissionModule{
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

func (m *PermissionModule) RegisterRoutes(r view.RouteRegistrar) {
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

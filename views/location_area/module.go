package location_area

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	locationareaaction "github.com/erniealice/entydad-golang/views/location_area/action"
	locationarealist "github.com/erniealice/entydad-golang/views/location_area/list"
)

// ModuleDeps holds all dependencies for the location area module.
type ModuleDeps struct {
	Routes       entydad.LocationAreaRoutes
	CommonLabels pyeza.CommonLabels
	SharedLabels entydad.SharedLabels
	Labels       entydad.LocationAreaLabels
	TableLabels  types.TableLabels

	// Data operations — caller provides these from their service/repo layer
	GetListPageData       func(ctx context.Context, status string, search string, page, pageSize int) (*locationarealist.LocationAreaListResult, error)
	GetInUseIDs           func(ctx context.Context, ids []string) (map[string]bool, error)
	CreateLocationArea    func(ctx context.Context, name, description string, active bool) (string, error)
	ReadLocationArea      func(ctx context.Context, id string) (*locationareaaction.LocationAreaRecord, error)
	UpdateLocationArea    func(ctx context.Context, id, name, description string, active bool) error
	DeleteLocationArea    func(ctx context.Context, id string) error
	SetLocationAreaActive func(ctx context.Context, id string, active bool) error
}

// Module holds all constructed location area views.
type Module struct {
	routes        entydad.LocationAreaRoutes
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
	actionDeps := &locationareaaction.Deps{
		CreateLocationArea:    deps.CreateLocationArea,
		ReadLocationArea:      deps.ReadLocationArea,
		UpdateLocationArea:    deps.UpdateLocationArea,
		DeleteLocationArea:    deps.DeleteLocationArea,
		SetLocationAreaActive: deps.SetLocationAreaActive,
		GetInUseIDs:           deps.GetInUseIDs,
		Routes:                deps.Routes,
		Labels:                deps.Labels,
	}
	listDeps := &locationarealist.ListViewDeps{
		GetListPageData: deps.GetListPageData,
		GetInUseIDs:     deps.GetInUseIDs,
		RefreshURL:      deps.Routes.TableURL,
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	}

	return &Module{
		routes:        deps.Routes,
		List:          locationarealist.NewView(listDeps),
		Table:         locationarealist.NewTableView(listDeps),
		Add:           locationareaaction.NewAddAction(actionDeps),
		Edit:          locationareaaction.NewEditAction(actionDeps),
		Delete:        locationareaaction.NewDeleteAction(actionDeps),
		BulkDelete:    locationareaaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     locationareaaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: locationareaaction.NewBulkSetStatusAction(actionDeps),
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

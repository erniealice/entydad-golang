package party

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	entitydelegate "github.com/erniealice/entydad-golang/domain/entity/party/delegate"
	delegateaction "github.com/erniealice/entydad-golang/domain/entity/party/delegate/action"
	delegatelist "github.com/erniealice/entydad-golang/domain/entity/party/delegate/list"
	delegatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/delegate"
)

// DelegateModuleDeps holds all dependencies for the delegate module.
// Trimmed vs ClientModuleDeps: no payment-terms, categories, subscriptions,
// attachments, audit, statement, revenue-run, dashboard, or status-transition
// deps — Delegate has only active bool and user identity.
type DelegateModuleDeps struct {
	Routes          entitydelegate.Routes
	CommonLabels    pyeza.CommonLabels
	SharedLabels    entydad.SharedLabels
	Labels          entitydelegate.Labels
	TableLabels     types.TableLabels
	GetListPageData func(ctx context.Context, req *delegatepb.GetDelegateListPageDataRequest) (*delegatepb.GetDelegateListPageDataResponse, error)
	// CRUD
	CreateDelegate func(ctx context.Context, req *delegatepb.CreateDelegateRequest) (*delegatepb.CreateDelegateResponse, error)
	ReadDelegate   func(ctx context.Context, req *delegatepb.ReadDelegateRequest) (*delegatepb.ReadDelegateResponse, error)
	UpdateDelegate func(ctx context.Context, req *delegatepb.UpdateDelegateRequest) (*delegatepb.UpdateDelegateResponse, error)
	DeleteDelegate func(ctx context.Context, req *delegatepb.DeleteDelegateRequest) (*delegatepb.DeleteDelegateResponse, error)
}

// DelegateModule holds all constructed delegate views.
type DelegateModule struct {
	routes     entitydelegate.Routes
	List       view.View
	Table      view.View
	Add        view.View
	Edit       view.View
	Delete     view.View
	BulkDelete view.View
}

// NewDelegateModule constructs a DelegateModule from the provided deps.
func NewDelegateModule(deps *DelegateModuleDeps) *DelegateModule {
	listDeps := &delegatelist.ListViewDeps{
		Routes:          deps.Routes,
		GetListPageData: deps.GetListPageData,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	}
	actionDeps := &delegateaction.Deps{
		Routes:         deps.Routes,
		CreateDelegate: deps.CreateDelegate,
		ReadDelegate:   deps.ReadDelegate,
		UpdateDelegate: deps.UpdateDelegate,
		DeleteDelegate: deps.DeleteDelegate,
	}

	return &DelegateModule{
		routes:     deps.Routes,
		List:       delegatelist.NewView(listDeps),
		Table:      delegatelist.NewTableView(listDeps),
		Add:        delegateaction.NewAddAction(actionDeps),
		Edit:       delegateaction.NewEditAction(actionDeps),
		Delete:     delegateaction.NewDeleteAction(actionDeps),
		BulkDelete: delegateaction.NewBulkDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers all delegate HTTP routes on the provided registrar.
func (m *DelegateModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
}

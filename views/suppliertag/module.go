package suppliertag

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	suppliertag "github.com/erniealice/entydad-golang/views/supplier/tag"
	suppliertagaction "github.com/erniealice/entydad-golang/views/supplier/tag/action"
	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	suppliercategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier_category"
)

// ModuleDeps holds all dependencies for the supplier tag module.
type ModuleDeps struct {
	Routes               entydad.SupplierTagRoutes
	Labels               entydad.SupplierTagLabels
	SharedLabels         entydad.SharedLabels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
	GetInUseIDs              func(ctx context.Context, ids []string) (map[string]bool, error)
	GetCategoryListPageData  func(ctx context.Context) ([]*categorypb.Category, error)
	ListCategories           func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	CreateCategory           func(ctx context.Context, req *categorypb.CreateCategoryRequest) (*categorypb.CreateCategoryResponse, error)
	ReadCategory             func(ctx context.Context, req *categorypb.ReadCategoryRequest) (*categorypb.ReadCategoryResponse, error)
	UpdateCategory           func(ctx context.Context, req *categorypb.UpdateCategoryRequest) (*categorypb.UpdateCategoryResponse, error)
	DeleteCategory           func(ctx context.Context, req *categorypb.DeleteCategoryRequest) (*categorypb.DeleteCategoryResponse, error)
	ListSupplierCategories   func(ctx context.Context, req *suppliercategorypb.ListSupplierCategoriesRequest) (*suppliercategorypb.ListSupplierCategoriesResponse, error)
	SetCategoryActive        func(ctx context.Context, id string, active bool) error
}

// Module holds all constructed supplier tag views.
type Module struct {
	routes        entydad.SupplierTagRoutes
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
	actionDeps := &suppliertagaction.Deps{
		Routes:            deps.Routes,
		CommonLabels:      deps.CommonLabels,
		ListCategories:    deps.ListCategories,
		CreateCategory:    deps.CreateCategory,
		ReadCategory:      deps.ReadCategory,
		UpdateCategory:    deps.UpdateCategory,
		DeleteCategory:    deps.DeleteCategory,
		SetCategoryActive: deps.SetCategoryActive,
	}

	listDeps := &suppliertag.Deps{
		Routes:                  deps.Routes,
		GetInUseIDs:             deps.GetInUseIDs,
		GetCategoryListPageData: deps.GetCategoryListPageData,
		ListSupplierCategories:  deps.ListSupplierCategories,
		RefreshURL:              deps.Routes.TableURL,
		Labels:                  deps.Labels,
		SharedLabels:            deps.SharedLabels,
		CommonLabels:            deps.CommonLabels,
		TableLabels:             deps.TableLabels,
	}

	return &Module{
		routes:        deps.Routes,
		List:          suppliertag.NewView(listDeps),
		Table:         suppliertag.NewTableView(listDeps),
		Add:           suppliertagaction.NewAddAction(actionDeps),
		Edit:          suppliertagaction.NewEditAction(actionDeps),
		Delete:        suppliertagaction.NewDeleteAction(actionDeps),
		BulkDelete:    suppliertagaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     suppliertagaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: suppliertagaction.NewBulkSetStatusAction(actionDeps),
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

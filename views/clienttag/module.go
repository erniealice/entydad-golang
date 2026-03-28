package clienttag

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	clienttag "github.com/erniealice/entydad-golang/views/client/tag"
	clienttagaction "github.com/erniealice/entydad-golang/views/client/tag/action"
	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
)

// ModuleDeps holds all dependencies for the client tag module.
type ModuleDeps struct {
	Routes               entydad.ClientTagRoutes
	Labels               entydad.ClientTagLabels
	SharedLabels         entydad.SharedLabels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
	GetInUseIDs          func(ctx context.Context, ids []string) (map[string]bool, error)
	ListCategories       func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	CreateCategory       func(ctx context.Context, req *categorypb.CreateCategoryRequest) (*categorypb.CreateCategoryResponse, error)
	ReadCategory         func(ctx context.Context, req *categorypb.ReadCategoryRequest) (*categorypb.ReadCategoryResponse, error)
	UpdateCategory       func(ctx context.Context, req *categorypb.UpdateCategoryRequest) (*categorypb.UpdateCategoryResponse, error)
	DeleteCategory       func(ctx context.Context, req *categorypb.DeleteCategoryRequest) (*categorypb.DeleteCategoryResponse, error)
	ListClientCategories func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
}

// Module holds all constructed client tag views.
type Module struct {
	routes     entydad.ClientTagRoutes
	List       view.View
	Add        view.View
	Edit       view.View
	Delete     view.View
	BulkDelete view.View
}

func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &clienttagaction.Deps{
		Routes:         deps.Routes,
		ListCategories: deps.ListCategories,
		CreateCategory: deps.CreateCategory,
		ReadCategory:   deps.ReadCategory,
		UpdateCategory: deps.UpdateCategory,
		DeleteCategory: deps.DeleteCategory,
	}

	return &Module{
		routes: deps.Routes,
		List: clienttag.NewView(&clienttag.Deps{
			Routes:               deps.Routes,
			GetInUseIDs:          deps.GetInUseIDs,
			ListCategories:       deps.ListCategories,
			ListClientCategories: deps.ListClientCategories,
			RefreshURL:           deps.Routes.ListURL,
			Labels:               deps.Labels,
			SharedLabels:         deps.SharedLabels,
			CommonLabels:         deps.CommonLabels,
			TableLabels:          deps.TableLabels,
		}),
		Add:        clienttagaction.NewAddAction(actionDeps),
		Edit:       clienttagaction.NewEditAction(actionDeps),
		Delete:     clienttagaction.NewDeleteAction(actionDeps),
		BulkDelete: clienttagaction.NewBulkDeleteAction(actionDeps),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
}

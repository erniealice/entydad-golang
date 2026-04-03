package payment_term

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	paymenttermaction "github.com/erniealice/entydad-golang/views/payment_term/action"
	paymenttermlist "github.com/erniealice/entydad-golang/views/payment_term/list"
	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"
)

// ModuleDeps holds all dependencies for the payment term module.
type ModuleDeps struct {
	Routes            entydad.PaymentTermRoutes
	CommonLabels      pyeza.CommonLabels
	SharedLabels      entydad.SharedLabels
	Labels            entydad.PaymentTermLabels
	TableLabels       types.TableLabels
	GetListPageData   func(ctx context.Context, req *paymenttermpb.GetPaymentTermListPageDataRequest) (*paymenttermpb.GetPaymentTermListPageDataResponse, error)
	CreatePaymentTerm func(ctx context.Context, req *paymenttermpb.CreatePaymentTermRequest) (*paymenttermpb.CreatePaymentTermResponse, error)
	ReadPaymentTerm   func(ctx context.Context, req *paymenttermpb.ReadPaymentTermRequest) (*paymenttermpb.ReadPaymentTermResponse, error)
	UpdatePaymentTerm func(ctx context.Context, req *paymenttermpb.UpdatePaymentTermRequest) (*paymenttermpb.UpdatePaymentTermResponse, error)
	DeletePaymentTerm func(ctx context.Context, req *paymenttermpb.DeletePaymentTermRequest) (*paymenttermpb.DeletePaymentTermResponse, error)
}

// Module holds all constructed payment term views.
type Module struct {
	routes     entydad.PaymentTermRoutes
	List       view.View
	Table      view.View
	Add        view.View
	Edit       view.View
	Delete     view.View
	BulkDelete view.View
}

// NewModule creates a new payment term module with all views wired up.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &paymenttermaction.Deps{
		Routes:            deps.Routes,
		CreatePaymentTerm: deps.CreatePaymentTerm,
		ReadPaymentTerm:   deps.ReadPaymentTerm,
		UpdatePaymentTerm: deps.UpdatePaymentTerm,
		DeletePaymentTerm: deps.DeletePaymentTerm,
	}

	listDeps := &paymenttermlist.Deps{
		GetListPageData: deps.GetListPageData,
		RefreshURL:      deps.Routes.TableURL,
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	}

	return &Module{
		routes:     deps.Routes,
		List:       paymenttermlist.NewView(listDeps),
		Table:      paymenttermlist.NewTableView(listDeps),
		Add:        paymenttermaction.NewAddAction(actionDeps),
		Edit:       paymenttermaction.NewEditAction(actionDeps),
		Delete:     paymenttermaction.NewDeleteAction(actionDeps),
		BulkDelete: paymenttermaction.NewBulkDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers all payment term routes with the given registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
}

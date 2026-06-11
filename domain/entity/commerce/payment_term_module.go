package commerce

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	paymentterm "github.com/erniealice/entydad-golang/domain/entity/commerce/payment_term"
	paymenttermaction "github.com/erniealice/entydad-golang/domain/entity/commerce/payment_term/action"
	paymenttermlist "github.com/erniealice/entydad-golang/domain/entity/commerce/payment_term/list"
	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"
)

// PaymentTermModuleDeps holds all dependencies for the payment term module.
type PaymentTermModuleDeps struct {
	Routes               paymentterm.Routes
	CommonLabels         pyeza.CommonLabels
	SharedLabels         entydad.SharedLabels
	Labels               paymentterm.Labels
	TableLabels          types.TableLabels
	GetListPageData      func(ctx context.Context, req *paymenttermpb.GetPaymentTermListPageDataRequest) (*paymenttermpb.GetPaymentTermListPageDataResponse, error)
	GetInUseIDs          func(ctx context.Context, ids []string) (map[string]bool, error)
	CreatePaymentTerm    func(ctx context.Context, req *paymenttermpb.CreatePaymentTermRequest) (*paymenttermpb.CreatePaymentTermResponse, error)
	ReadPaymentTerm      func(ctx context.Context, req *paymenttermpb.ReadPaymentTermRequest) (*paymenttermpb.ReadPaymentTermResponse, error)
	UpdatePaymentTerm    func(ctx context.Context, req *paymenttermpb.UpdatePaymentTermRequest) (*paymenttermpb.UpdatePaymentTermResponse, error)
	DeletePaymentTerm    func(ctx context.Context, req *paymenttermpb.DeletePaymentTermRequest) (*paymenttermpb.DeletePaymentTermResponse, error)
	SetPaymentTermActive func(ctx context.Context, id string, active bool) error
	// Scope filters which payment terms are shown in the list page.
	// Valid values: "client" (shows client + both), "supplier" (shows supplier + both).
	// Leave empty to show all terms (used when registering a standalone settings page).
	Scope string
}

// PaymentTermModule holds all constructed payment term views.
type PaymentTermModule struct {
	routes        paymentterm.Routes
	List          view.View
	Table         view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	BulkDelete    view.View
	SetStatus     view.View
	BulkSetStatus view.View
}

// NewPaymentTermModule creates a new payment term module with all views wired up.
func NewPaymentTermModule(deps *PaymentTermModuleDeps) *PaymentTermModule {
	actionDeps := &paymenttermaction.Deps{
		Routes:               deps.Routes,
		CreatePaymentTerm:    deps.CreatePaymentTerm,
		ReadPaymentTerm:      deps.ReadPaymentTerm,
		UpdatePaymentTerm:    deps.UpdatePaymentTerm,
		DeletePaymentTerm:    deps.DeletePaymentTerm,
		SetPaymentTermActive: deps.SetPaymentTermActive,
		Scope:                deps.Scope,
	}

	listDeps := &paymenttermlist.Deps{
		GetListPageData: deps.GetListPageData,
		GetInUseIDs:     deps.GetInUseIDs,
		RefreshURL:      deps.Routes.TableURL,
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
		Scope:           deps.Scope,
	}

	return &PaymentTermModule{
		routes:        deps.Routes,
		List:          paymenttermlist.NewView(listDeps),
		Table:         paymenttermlist.NewTableView(listDeps),
		Add:           paymenttermaction.NewAddAction(actionDeps),
		Edit:          paymenttermaction.NewEditAction(actionDeps),
		Delete:        paymenttermaction.NewDeleteAction(actionDeps),
		BulkDelete:    paymenttermaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     paymenttermaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: paymenttermaction.NewBulkSetStatusAction(actionDeps),
	}
}

// RegisterRoutes registers all payment term routes with the given registrar.
func (m *PaymentTermModule) RegisterRoutes(r view.RouteRegistrar) {
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

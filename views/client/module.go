package client

import (
	"context"
	"log"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	clientaction "github.com/erniealice/entydad-golang/views/client/action"
	clientdashboard "github.com/erniealice/entydad-golang/views/client/dashboard"
	clientdetail "github.com/erniealice/entydad-golang/views/client/detail"
	clientform "github.com/erniealice/entydad-golang/views/client/form"
	clientlist "github.com/erniealice/entydad-golang/views/client/list"
	categorypb   "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	clientpb     "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
	clientstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/client_statement"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
)

// PaymentTermOption is re-exported from form for use by callers wiring ModuleDeps.
type PaymentTermOption = clientform.PaymentTermOption

// ModuleDeps holds all dependencies for the client module.
type ModuleDeps struct {
	Routes               entydad.ClientRoutes
	// SearchTimezonesURL points to the timezone autocomplete JSON endpoint
	// (provided by the user module). The client representative section reuses
	// the same handler for its timezone picker.
	SearchTimezonesURL string
	CommonLabels         pyeza.CommonLabels
	SharedLabels         entydad.SharedLabels
	Labels               entydad.ClientLabels
	DashboardLabels      entydad.ClientDashboardLabels
	DashboardTitleLabels entydad.DashboardLabels
	TableLabels          types.TableLabels
	GetListPageData      func(ctx context.Context, req *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
	GetInUseIDs          func(ctx context.Context, ids []string) (map[string]bool, error)
	GetClientBalances    func(ctx context.Context) (map[string]int64, error)
	// Client CRUD
	CreateClient func(ctx context.Context, req *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error)
	ReadClient   func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	UpdateClient func(ctx context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error)
	DeleteClient func(ctx context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error)
	SetStatus    func(ctx context.Context, id string, status string) error
	// Payment terms dropdown
	ListPaymentTerms func(ctx context.Context) ([]*clientform.PaymentTermOption, error)
	// Categories (client tags)
	ListCategories       func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListClientCategories func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
	CreateClientCategory func(ctx context.Context, req *clientcategorypb.CreateClientCategoryRequest) (*clientcategorypb.CreateClientCategoryResponse, error)
	DeleteClientCategory func(ctx context.Context, req *clientcategorypb.DeleteClientCategoryRequest) (*clientcategorypb.DeleteClientCategoryResponse, error)
	// Revenue listing (for detail view)
	ListRevenues func(ctx context.Context, collection string) ([]map[string]any, error)
	// Client statement (for detail view)
	GetClientStatement func(ctx context.Context, req *clientstmtpb.ClientStatementRequest) (*clientstmtpb.ClientStatementResponse, error)
	// Subscription listing (for detail view)
	ListSubscriptions           func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)
	GetSubscriptionListPageData func(ctx context.Context, req *subscriptionpb.GetSubscriptionListPageDataRequest) (*subscriptionpb.GetSubscriptionListPageDataResponse, error)
	// Subscription URLs (cross-module, from centymo).
	SubscriptionAddURL    string
	SubscriptionDetailURL string
	// SubscriptionUnderClientDetailURL is the nested-route template; when set,
	// the engagements row link uses it so the subscription detail renders
	// with a "client → subscription" breadcrumb.
	SubscriptionUnderClientDetailURL string
	SubscriptionEditURL              string
	SubscriptionDeleteURL            string

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Audit history
	ListAuditHistory func(ctx context.Context, req *auditlog.ListAuditRequest) (*auditlog.ListAuditResponse, error)

	// GetFunctionalCurrency resolves the current workspace's functional currency
	// so new-client drawers can prefill billing_currency.
	GetFunctionalCurrency func(ctx context.Context) string

	// ListClientPlans fetches Plans scoped to a client_id for the Packages tab.
	// Wired from centymo via the centymo block; nil-safe (tab renders empty state).
	ListClientPlans func(ctx context.Context, clientID string) ([]clientdetail.ClientPlanRow, error)

	// PlanAddURL is the centymo Plan-add drawer URL; the Packages tab appends
	// ?context=client&client_id={cid} to pre-fill the client field.
	PlanAddURL string
}

// Module holds all constructed client views.
type Module struct {
	routes           entydad.ClientRoutes
	Dashboard        view.View
	List             view.View
	Table            view.View
	Detail           view.View
	TabAction        view.View
	Add              view.View
	Edit             view.View
	Delete           view.View
	BulkDelete       view.View
	SetStatus        view.View
	BulkSetStatus    view.View
	AttachmentUpload view.View
	AttachmentDelete view.View
	StatementExport  http.HandlerFunc
}

func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &clientaction.Deps{
		Routes:               deps.Routes,
		SearchTimezonesURL:   deps.SearchTimezonesURL,
		CreateClient:         deps.CreateClient,
		ReadClient:           deps.ReadClient,
		UpdateClient:         deps.UpdateClient,
		DeleteClient:         deps.DeleteClient,
		SetClientStatus:      deps.SetStatus,
		ListPaymentTerms:     deps.ListPaymentTerms,
		ListCategories:       deps.ListCategories,
		ListClientCategories: deps.ListClientCategories,
		CreateClientCategory:  deps.CreateClientCategory,
		DeleteClientCategory:  deps.DeleteClientCategory,
		GetFunctionalCurrency: deps.GetFunctionalCurrency,
	}
	listDeps := &clientlist.ListViewDeps{
		Routes:            deps.Routes,
		GetListPageData:   deps.GetListPageData,
		GetInUseIDs:       deps.GetInUseIDs,
		GetClientBalances: deps.GetClientBalances,
		Labels:            deps.Labels,
		SharedLabels:      deps.SharedLabels,
		CommonLabels:      deps.CommonLabels,
		TableLabels:       deps.TableLabels,
	}
	detailDeps := &clientdetail.DetailViewDeps{
		Routes:                      deps.Routes,
		ReadClient:                  deps.ReadClient,
		ListCategories:              deps.ListCategories,
		ListClientCategories:        deps.ListClientCategories,
		ListRevenues:                deps.ListRevenues,
		GetClientStatement:          deps.GetClientStatement,
		ListSubscriptions:           deps.ListSubscriptions,
		GetSubscriptionListPageData: deps.GetSubscriptionListPageData,
		SubscriptionAddURL:                deps.SubscriptionAddURL,
		SubscriptionDetailURL:             deps.SubscriptionDetailURL,
		SubscriptionUnderClientDetailURL:  deps.SubscriptionUnderClientDetailURL,
		SubscriptionEditURL:               deps.SubscriptionEditURL,
		SubscriptionDeleteURL:             deps.SubscriptionDeleteURL,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		TableLabels:           deps.TableLabels,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
		AuditOps: auditlog.AuditOps{
			ListAuditHistory: deps.ListAuditHistory,
		},
		ListClientPlans: deps.ListClientPlans,
		PlanAddURL:      deps.PlanAddURL,
	}

	return &Module{
		routes:           deps.Routes,
		Dashboard:        clientdashboard.NewView(&clientdashboard.Deps{DashboardLabels: deps.DashboardTitleLabels, CommonLabels: deps.CommonLabels, Dashboard: deps.DashboardLabels, Routes: deps.Routes}),
		List:             clientlist.NewView(listDeps),
		Table:            clientlist.NewTableView(listDeps),
		Detail:           clientdetail.NewView(detailDeps),
		TabAction:        clientdetail.NewTabAction(detailDeps),
		Add:              clientaction.NewAddAction(actionDeps),
		Edit:             clientaction.NewEditAction(actionDeps),
		Delete:           clientaction.NewDeleteAction(actionDeps),
		BulkDelete:       clientaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        clientaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    clientaction.NewBulkSetStatusAction(actionDeps),
		AttachmentUpload: clientdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: clientdetail.NewAttachmentDeleteAction(detailDeps),
		StatementExport:  clientdetail.NewStatementExportHandler(detailDeps),
	}
}

// routeRegistrarFull extends view.RouteRegistrar with HandleFunc support
// for raw http.HandlerFunc routes (e.g., CSV exports).
type routeRegistrarFull interface {
	view.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// handleFunc is a nil-safe helper that registers an http.HandlerFunc route if
// the RouteRegistrar supports it, otherwise logs a warning and skips.
func handleFunc(r view.RouteRegistrar, method, path string, handler http.HandlerFunc) {
	if handler == nil {
		return
	}
	if full, ok := r.(routeRegistrarFull); ok {
		full.HandleFunc(method, path, handler)
		return
	}
	log.Printf("client: RouteRegistrar does not support HandleFunc — skipping %s %s", method, path)
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
	// Statement CSV export
	handleFunc(r, "GET", m.routes.StatementExportURL, m.StatementExport)
}

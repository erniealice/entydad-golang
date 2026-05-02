package supplier

import (
	"context"
	"log"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	supplieraction "github.com/erniealice/entydad-golang/views/supplier/action"
	supplierdashboard "github.com/erniealice/entydad-golang/views/supplier/dashboard"
	supplierdetail "github.com/erniealice/entydad-golang/views/supplier/detail"
	supplierform "github.com/erniealice/entydad-golang/views/supplier/form"
	supplierlist "github.com/erniealice/entydad-golang/views/supplier/list"
	categorypb      "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb    "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	supplierpb      "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	suppliercategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier_category"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	suppstmtpb      "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
)

// PaymentTermOption is re-exported from form for use by callers wiring ModuleDeps.
type PaymentTermOption = supplierform.PaymentTermOption

// ModuleDeps holds all dependencies for the supplier module.
type ModuleDeps struct {
	Routes               entydad.SupplierRoutes
	CommonLabels         pyeza.CommonLabels
	SharedLabels         entydad.SharedLabels
	Labels               entydad.SupplierLabels
	DashboardLabels      entydad.SupplierDashboardLabels
	DashboardTitleLabels entydad.DashboardLabels
	TableLabels          types.TableLabels
	GetListPageData func(ctx context.Context, req *supplierpb.GetSupplierListPageDataRequest) (*supplierpb.GetSupplierListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	// Supplier CRUD
	CreateSupplier   func(ctx context.Context, req *supplierpb.CreateSupplierRequest) (*supplierpb.CreateSupplierResponse, error)
	ReadSupplier     func(ctx context.Context, req *supplierpb.ReadSupplierRequest) (*supplierpb.ReadSupplierResponse, error)
	UpdateSupplier   func(ctx context.Context, req *supplierpb.UpdateSupplierRequest) (*supplierpb.UpdateSupplierResponse, error)
	DeleteSupplier   func(ctx context.Context, req *supplierpb.DeleteSupplierRequest) (*supplierpb.DeleteSupplierResponse, error)
	SetStatus        func(ctx context.Context, id string, status string) error
	ListPaymentTerms func(ctx context.Context) ([]*PaymentTermOption, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Audit history
	ListAuditHistory func(ctx context.Context, req *auditlog.ListAuditRequest) (*auditlog.ListAuditResponse, error)

	// Purchase orders
	ListPurchaseOrders func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)

	// Supplier statement
	GetSupplierStatement func(ctx context.Context, req *suppstmtpb.SupplierStatementRequest) (*suppstmtpb.SupplierStatementResponse, error)

	// Outstanding balances for supplier list
	GetSupplierBalances func(ctx context.Context) (map[string]int64, error)

	// Tag-related deps for multi-select tags on the supplier form
	ListCategories         func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListSupplierCategories func(ctx context.Context, req *suppliercategorypb.ListSupplierCategoriesRequest) (*suppliercategorypb.ListSupplierCategoriesResponse, error)
	CreateSupplierCategory func(ctx context.Context, req *suppliercategorypb.CreateSupplierCategoryRequest) (*suppliercategorypb.CreateSupplierCategoryResponse, error)
	DeleteSupplierCategory func(ctx context.Context, req *suppliercategorypb.DeleteSupplierCategoryRequest) (*suppliercategorypb.DeleteSupplierCategoryResponse, error)
}

// Module holds all constructed supplier views.
type Module struct {
	routes           entydad.SupplierRoutes
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
	actionDeps := &supplieraction.Deps{
		Routes:                 deps.Routes,
		CreateSupplier:         deps.CreateSupplier,
		ReadSupplier:           deps.ReadSupplier,
		UpdateSupplier:         deps.UpdateSupplier,
		DeleteSupplier:         deps.DeleteSupplier,
		SetSupplierStatus:      deps.SetStatus,
		ListPaymentTerms:       deps.ListPaymentTerms,
		ListCategories:         deps.ListCategories,
		ListSupplierCategories: deps.ListSupplierCategories,
		CreateSupplierCategory: deps.CreateSupplierCategory,
		DeleteSupplierCategory: deps.DeleteSupplierCategory,
	}
	listDeps := &supplierlist.ListViewDeps{
		Routes:              deps.Routes,
		GetListPageData:     deps.GetListPageData,
		GetInUseIDs:         deps.GetInUseIDs,
		Labels:              deps.Labels,
		SharedLabels:        deps.SharedLabels,
		CommonLabels:        deps.CommonLabels,
		TableLabels:         deps.TableLabels,
		GetSupplierBalances: deps.GetSupplierBalances,
	}
	detailDeps := &supplierdetail.DetailViewDeps{
		Routes:       deps.Routes,
		ReadSupplier: deps.ReadSupplier,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
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
		ListPurchaseOrders:     deps.ListPurchaseOrders,
		GetSupplierStatement:   deps.GetSupplierStatement,
		ListCategories:         deps.ListCategories,
		ListSupplierCategories: deps.ListSupplierCategories,
	}

	return &Module{
		routes: deps.Routes,
		Dashboard: supplierdashboard.NewView(&supplierdashboard.Deps{
			DashboardLabels: deps.DashboardTitleLabels,
			Dashboard:       deps.DashboardLabels,
			CommonLabels:    deps.CommonLabels,
		}),
		List:             supplierlist.NewView(listDeps),
		Table:            supplierlist.NewTableView(listDeps),
		Detail:           supplierdetail.NewView(detailDeps),
		TabAction:        supplierdetail.NewTabAction(detailDeps),
		Add:              supplieraction.NewAddAction(actionDeps),
		Edit:             supplieraction.NewEditAction(actionDeps),
		Delete:           supplieraction.NewDeleteAction(actionDeps),
		BulkDelete:       supplieraction.NewBulkDeleteAction(actionDeps),
		SetStatus:        supplieraction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    supplieraction.NewBulkSetStatusAction(actionDeps),
		AttachmentUpload: supplierdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: supplierdetail.NewAttachmentDeleteAction(detailDeps),
		StatementExport:  supplierdetail.NewStatementExportHandler(detailDeps),
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
	log.Printf("supplier: RouteRegistrar does not support HandleFunc — skipping %s %s", method, path)
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

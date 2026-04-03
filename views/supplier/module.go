package supplier

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	supplieraction "github.com/erniealice/entydad-golang/views/supplier/action"
	supplierdetail "github.com/erniealice/entydad-golang/views/supplier/detail"
	supplierlist "github.com/erniealice/entydad-golang/views/supplier/list"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
)

// PaymentTermOption is re-exported from action for use by callers wiring ModuleDeps.
type PaymentTermOption = supplieraction.PaymentTermOption

// ModuleDeps holds all dependencies for the supplier module.
type ModuleDeps struct {
	Routes          entydad.SupplierRoutes
	CommonLabels    pyeza.CommonLabels
	SharedLabels    entydad.SharedLabels
	Labels          entydad.SupplierLabels
	TableLabels     types.TableLabels
	GetListPageData func(ctx context.Context, req *supplierpb.GetSupplierListPageDataRequest) (*supplierpb.GetSupplierListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	// Supplier CRUD
	CreateSupplier   func(ctx context.Context, req *supplierpb.CreateSupplierRequest) (*supplierpb.CreateSupplierResponse, error)
	ReadSupplier     func(ctx context.Context, req *supplierpb.ReadSupplierRequest) (*supplierpb.ReadSupplierResponse, error)
	UpdateSupplier   func(ctx context.Context, req *supplierpb.UpdateSupplierRequest) (*supplierpb.UpdateSupplierResponse, error)
	DeleteSupplier   func(ctx context.Context, req *supplierpb.DeleteSupplierRequest) (*supplierpb.DeleteSupplierResponse, error)
	SetActive        func(ctx context.Context, id string, active bool) error
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
}

// Module holds all constructed supplier views.
type Module struct {
	routes           entydad.SupplierRoutes
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
}

func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &supplieraction.Deps{
		Routes:            deps.Routes,
		CreateSupplier:    deps.CreateSupplier,
		ReadSupplier:      deps.ReadSupplier,
		UpdateSupplier:    deps.UpdateSupplier,
		DeleteSupplier:    deps.DeleteSupplier,
		SetSupplierActive: deps.SetActive,
		ListPaymentTerms:  deps.ListPaymentTerms,
	}
	listDeps := &supplierlist.ListViewDeps{
		Routes:          deps.Routes,
		GetListPageData: deps.GetListPageData,
		GetInUseIDs:     deps.GetInUseIDs,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
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
		ListPurchaseOrders: deps.ListPurchaseOrders,
	}

	return &Module{
		routes:           deps.Routes,
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
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
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
}

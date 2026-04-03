package client

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	clientaction "github.com/erniealice/entydad-golang/views/client/action"
	clientdashboard "github.com/erniealice/entydad-golang/views/client/dashboard"
	clientdetail "github.com/erniealice/entydad-golang/views/client/detail"
	clientlist "github.com/erniealice/entydad-golang/views/client/list"
	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
)

// PaymentTermOption is re-exported from action for use by callers wiring ModuleDeps.
type PaymentTermOption = clientaction.PaymentTermOption

// ModuleDeps holds all dependencies for the client module.
type ModuleDeps struct {
	Routes           entydad.ClientRoutes
	CommonLabels     pyeza.CommonLabels
	SharedLabels     entydad.SharedLabels
	Labels           entydad.ClientLabels
	DashboardLabels      entydad.ClientDashboardLabels
	DashboardTitleLabels entydad.DashboardLabels
	TableLabels          types.TableLabels
	GetListPageData func(ctx context.Context, req *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	// Client CRUD
	CreateClient func(ctx context.Context, req *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error)
	ReadClient   func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	UpdateClient func(ctx context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error)
	DeleteClient func(ctx context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error)
	SetActive    func(ctx context.Context, id string, active bool) error
	// Payment terms dropdown
	ListPaymentTerms func(ctx context.Context) ([]*clientaction.PaymentTermOption, error)
	// Categories (client tags)
	ListCategories       func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListClientCategories func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
	CreateClientCategory func(ctx context.Context, req *clientcategorypb.CreateClientCategoryRequest) (*clientcategorypb.CreateClientCategoryResponse, error)
	DeleteClientCategory func(ctx context.Context, req *clientcategorypb.DeleteClientCategoryRequest) (*clientcategorypb.DeleteClientCategoryResponse, error)
	// Revenue listing (for detail view)
	ListRevenues func(ctx context.Context, collection string) ([]map[string]any, error)
	// Subscription listing (for detail view)
	ListSubscriptions func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)
	// Subscription URLs (cross-module, from centymo).
	SubscriptionAddURL    string
	SubscriptionDetailURL string
	SubscriptionEditURL   string
	SubscriptionDeleteURL string

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Audit history
	ListAuditHistory func(ctx context.Context, req *auditlog.ListAuditRequest) (*auditlog.ListAuditResponse, error)
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
}

func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &clientaction.Deps{
		Routes:               deps.Routes,
		CreateClient:         deps.CreateClient,
		ReadClient:           deps.ReadClient,
		UpdateClient:         deps.UpdateClient,
		DeleteClient:         deps.DeleteClient,
		SetClientActive:      deps.SetActive,
		ListPaymentTerms:     deps.ListPaymentTerms,
		ListCategories:       deps.ListCategories,
		ListClientCategories: deps.ListClientCategories,
		CreateClientCategory: deps.CreateClientCategory,
		DeleteClientCategory: deps.DeleteClientCategory,
	}
	listDeps := &clientlist.ListViewDeps{
		Routes:          deps.Routes,
		GetListPageData: deps.GetListPageData,
		GetInUseIDs:     deps.GetInUseIDs,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	}
	detailDeps := &clientdetail.DetailViewDeps{
		Routes:               deps.Routes,
		ReadClient:           deps.ReadClient,
		ListCategories:       deps.ListCategories,
		ListClientCategories: deps.ListClientCategories,
		ListRevenues:         deps.ListRevenues,
		ListSubscriptions:     deps.ListSubscriptions,
		SubscriptionAddURL:    deps.SubscriptionAddURL,
		SubscriptionDetailURL: deps.SubscriptionDetailURL,
		SubscriptionEditURL:   deps.SubscriptionEditURL,
		SubscriptionDeleteURL: deps.SubscriptionDeleteURL,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
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
	}

	return &Module{
		routes:           deps.Routes,
		Dashboard:        clientdashboard.NewView(&clientdashboard.Deps{DashboardLabels: deps.DashboardTitleLabels, CommonLabels: deps.CommonLabels, Dashboard: deps.DashboardLabels}),
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
	}
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
}

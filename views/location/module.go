package location

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	locationaction "github.com/erniealice/entydad-golang/views/location/action"
	locationdashboard "github.com/erniealice/entydad-golang/views/location/dashboard"
	locationdetail "github.com/erniealice/entydad-golang/views/location/detail"
	locationlist "github.com/erniealice/entydad-golang/views/location/list"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
)

// ModuleDeps holds all dependencies for the location module.
type ModuleDeps struct {
	Routes               entydad.LocationRoutes
	CommonLabels         pyeza.CommonLabels
	SharedLabels         entydad.SharedLabels
	Labels               entydad.LocationLabels
	DashboardTitleLabels entydad.DashboardLabels
	TableLabels          types.TableLabels
	GetListPageData      func(ctx context.Context, req *locationpb.GetLocationListPageDataRequest) (*locationpb.GetLocationListPageDataResponse, error)
	GetInUseIDs          func(ctx context.Context, ids []string) (map[string]bool, error)
	CreateLocation       func(ctx context.Context, req *locationpb.CreateLocationRequest) (*locationpb.CreateLocationResponse, error)
	ReadLocation         func(ctx context.Context, req *locationpb.ReadLocationRequest) (*locationpb.ReadLocationResponse, error)
	UpdateLocation       func(ctx context.Context, req *locationpb.UpdateLocationRequest) (*locationpb.UpdateLocationResponse, error)
	DeleteLocation       func(ctx context.Context, req *locationpb.DeleteLocationRequest) (*locationpb.DeleteLocationResponse, error)
	SetActive            func(ctx context.Context, id string, active bool) error
	// ListLocationAreas is optional — if provided, the area dropdown appears in the form.
	ListLocationAreas func(ctx context.Context) ([]locationaction.LocationAreaOption, error)

	// LocationAreaRoutes provides deep-link URLs for the location-area domain
	// (passed through to the dashboard view for quick-action links).
	LocationAreaRoutes entydad.LocationAreaRoutes

	// GetLocationDashboardPageData is the workspace-scoped page-data fetch
	// for the dashboard view. The container builds this by calling the
	// GetLocationDashboardPageDataUseCase. nil-safe: when missing, the
	// view renders empty-state widgets.
	GetLocationDashboardPageData func(ctx context.Context) (*locationdashboard.LocationDashboardData, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Audit history
	ListAuditHistory func(ctx context.Context, req *auditlog.ListAuditRequest) (*auditlog.ListAuditResponse, error)
}

// Module holds all constructed location views.
type Module struct {
	routes           entydad.LocationRoutes
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
	actionDeps := &locationaction.Deps{
		CreateLocation:    deps.CreateLocation,
		ReadLocation:      deps.ReadLocation,
		UpdateLocation:    deps.UpdateLocation,
		DeleteLocation:    deps.DeleteLocation,
		SetLocationActive: deps.SetActive,
		GetInUseIDs:       deps.GetInUseIDs,
		ListLocationAreas: deps.ListLocationAreas,
		Routes:            deps.Routes,
		Labels:            deps.Labels,
	}
	listDeps := &locationlist.ListViewDeps{
		GetListPageData: deps.GetListPageData,
		GetInUseIDs:     deps.GetInUseIDs,
		RefreshURL:      deps.Routes.TableURL,
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	}
	detailDeps := &locationdetail.DetailViewDeps{
		Routes:       deps.Routes,
		ReadLocation: deps.ReadLocation,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
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
		routes: deps.Routes,
		Dashboard: locationdashboard.NewView(&locationdashboard.Deps{
			DashboardLabels:    deps.DashboardTitleLabels,
			Dashboard:          deps.Labels.Dashboard,
			Routes:             deps.Routes,
			LocationAreaRoutes: deps.LocationAreaRoutes,
			CommonLabels:       deps.CommonLabels,
			GetDashboardData:   deps.GetLocationDashboardPageData,
		}),
		List:             locationlist.NewView(listDeps),
		Table:            locationlist.NewTableView(listDeps),
		Detail:           locationdetail.NewView(detailDeps),
		TabAction:        locationdetail.NewTabAction(detailDeps),
		Add:              locationaction.NewAddAction(actionDeps),
		Edit:             locationaction.NewEditAction(actionDeps),
		Delete:           locationaction.NewDeleteAction(actionDeps),
		BulkDelete:       locationaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        locationaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    locationaction.NewBulkSetStatusAction(actionDeps),
		AttachmentUpload: locationdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: locationdetail.NewAttachmentDeleteAction(detailDeps),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)
	r.POST(m.routes.EditDetailURL, m.Edit)
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

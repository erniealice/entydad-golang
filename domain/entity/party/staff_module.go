package party

import (
	"context"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	entitystaff "github.com/erniealice/entydad-golang/domain/entity/party/staff"
	staffaction "github.com/erniealice/entydad-golang/domain/entity/party/staff/action"
	stafflist "github.com/erniealice/entydad-golang/domain/entity/party/staff/list"
	staffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/staff"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
)

// StaffModuleDeps holds all dependencies for the staff module.
// Trimmed vs ClientModuleDeps, mirroring DelegateModuleDeps: no payment-terms,
// categories, subscriptions, attachments, audit, statement, revenue-run,
// dashboard, or status-transition deps. Adds ReadUser/ListUsers — Staff links
// to an EXISTING user via auto-complete rather than embedding a new one.
type StaffModuleDeps struct {
	Routes          entitystaff.Routes
	CommonLabels    pyeza.CommonLabels
	SharedLabels    entydad.SharedLabels
	Labels          entitystaff.Labels
	TableLabels     types.TableLabels
	GetListPageData func(ctx context.Context, req *staffpb.GetStaffListPageDataRequest) (*staffpb.GetStaffListPageDataResponse, error)
	GetItemPageData func(ctx context.Context, req *staffpb.GetStaffItemPageDataRequest) (*staffpb.GetStaffItemPageDataResponse, error)
	// CRUD
	CreateStaff func(ctx context.Context, req *staffpb.CreateStaffRequest) (*staffpb.CreateStaffResponse, error)
	UpdateStaff func(ctx context.Context, req *staffpb.UpdateStaffRequest) (*staffpb.UpdateStaffResponse, error)
	DeleteStaff func(ctx context.Context, req *staffpb.DeleteStaffRequest) (*staffpb.DeleteStaffResponse, error)
	// User linkage (autocomplete + resolve-on-submit)
	ReadUser  func(ctx context.Context, req *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error)
	ListUsers func(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error)
}

// StaffModule holds all constructed staff views.
type StaffModule struct {
	routes     entitystaff.Routes
	List       view.View
	Table      view.View
	Add        view.View
	Edit       view.View
	Delete     view.View
	BulkDelete view.View
	// UserSearch is a raw http.HandlerFunc for the user autocomplete endpoint.
	UserSearch http.HandlerFunc
}

// NewStaffModule constructs a StaffModule from the provided deps.
func NewStaffModule(deps *StaffModuleDeps) *StaffModule {
	listDeps := &stafflist.ListViewDeps{
		Routes:          deps.Routes,
		GetListPageData: deps.GetListPageData,
		Labels:          deps.Labels,
		SharedLabels:    deps.SharedLabels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	}
	actionDeps := &staffaction.Deps{
		Routes:          deps.Routes,
		CreateStaff:     deps.CreateStaff,
		GetItemPageData: deps.GetItemPageData,
		UpdateStaff:     deps.UpdateStaff,
		DeleteStaff:     deps.DeleteStaff,
		ReadUser:        deps.ReadUser,
		ListUsers:       deps.ListUsers,
	}

	return &StaffModule{
		routes:     deps.Routes,
		List:       stafflist.NewView(listDeps),
		Table:      stafflist.NewTableView(listDeps),
		Add:        staffaction.NewAddAction(actionDeps),
		Edit:       staffaction.NewEditAction(actionDeps),
		Delete:     staffaction.NewDeleteAction(actionDeps),
		BulkDelete: staffaction.NewBulkDeleteAction(actionDeps),
		UserSearch: staffaction.NewUserSearchAction(actionDeps),
	}
}

// RegisterRoutes registers all staff HTTP routes on the provided registrar.
func (m *StaffModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	// User search is a raw HTTP handler (returns JSON) — reuses the
	// client-module's HandleFunc shim (same package, same registrar contract).
	if m.routes.SearchURL != "" && m.UserSearch != nil {
		clientHandleFunc(r, "GET", m.routes.SearchURL, m.UserSearch)
	}
}

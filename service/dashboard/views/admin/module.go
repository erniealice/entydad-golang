package admin

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	admindashboard "github.com/erniealice/entydad-golang/service/dashboard/views/admin/dashboard"
)

// ModuleDeps holds the dependencies needed by the admin app dashboard view.
//
// The admin app is composite — it does not own any CRUD entities of its
// own, so the module surface is intentionally small (one dashboard view
// only). The constituent CRUD modules (permission, role, workspace,
// workspace_user, workspace_user_role) continue to be wired separately
// by service-admin.
type ModuleDeps struct {
	Routes               entydad.AdminDashboardRoutes
	CommonLabels         pyeza.CommonLabels
	DashboardLabels      entydad.AdminDashboardLabels
	DashboardTitleLabels entydad.DashboardLabels

	// DashboardRoutes holds the cross-entity URLs the admin dashboard's
	// quick actions and "view all" buttons resolve to. The orchestrator
	// supplies these from the existing entydad.PermissionRoutes /
	// RoleRoutes / WorkspaceRoutes / WorkspaceUserRoutes.
	DashboardRoutes admindashboard.Routes

	// GetDashboardData is the workspace-scoped page-data fetch. The
	// container builds this by calling the
	// GetAdminDashboardPageDataUseCase. nil-safe: when missing, the
	// view renders empty-state widgets.
	GetDashboardData func(ctx context.Context) (*admindashboard.AdminDashboardData, error)
}

// Module holds the constructed admin views.
type Module struct {
	routes    entydad.AdminDashboardRoutes
	Dashboard view.View
}

// NewModule constructs the admin app module. Currently dashboard-only.
func NewModule(deps *ModuleDeps) *Module {
	dashRoutes := deps.DashboardRoutes
	if dashRoutes.DashboardURL == "" {
		dashRoutes.DashboardURL = deps.Routes.DashboardURL
	}
	return &Module{
		routes: deps.Routes,
		Dashboard: admindashboard.NewView(&admindashboard.Deps{
			DashboardLabels:  deps.DashboardTitleLabels,
			Dashboard:        deps.DashboardLabels,
			Routes:           dashRoutes,
			CommonLabels:     deps.CommonLabels,
			GetDashboardData: deps.GetDashboardData,
		}),
	}
}

// RegisterRoutes wires the admin dashboard route to the registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
}

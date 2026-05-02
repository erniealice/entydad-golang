package block

// wiring.go wires dashboard use cases from the espyna UseCases aggregate
// into entydad module ModuleDeps callbacks.
//
// Since the dashboard use-case request/response types live in espyna's
// internal packages (unreachable from entydad), we use reflection to:
//  1. Dereference the use-case pointer field (e.g. useCases.Entity.LocationDashboard)
//  2. Build the Execute request via reflect.New + field-name assignment
//  3. Call Execute reflectively
//  4. Copy matching fields from the response to the view-layer *XxxDashboardData type
//
// All helpers are nil-safe: if the Dashboard field is nil the callback is
// left unset and the dashboard view renders empty state (its existing behaviour).
//
// WorkspaceID is extracted from the request context via
// consumer.GetWorkspaceIDFromContext — the view layer callbacks have the
// signature func(ctx context.Context) (*XxxDashboardData, error) and rely on
// the orchestrator to thread the workspace ID through.

import (
	"context"
	"reflect"
	"time"

	consumer "github.com/erniealice/espyna-golang/consumer"

	locationdashboard "github.com/erniealice/entydad-golang/views/location/dashboard"
	admindashboard "github.com/erniealice/entydad-golang/views/admin/dashboard"

	locationmod "github.com/erniealice/entydad-golang/views/location"
	adminmod "github.com/erniealice/entydad-golang/views/admin"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
)

// callEntydadExecute calls a use-case's Execute method via reflection.
// The useCase value must be a non-nil pointer to a use-case struct.
// workspaceID and now are set on the request struct by field name.
// Returns the dereferenced response (as reflect.Value) and an error.
func callEntydadExecute(
	useCase reflect.Value,
	ctx context.Context,
	workspaceID string,
	now time.Time,
) (reflect.Value, error) {
	m := useCase.MethodByName("Execute")
	if !m.IsValid() {
		return reflect.Value{}, nil
	}
	reqType := m.Type().In(1).Elem() // *Request → Request
	reqPtr := reflect.New(reqType)
	if f := reqPtr.Elem().FieldByName("WorkspaceID"); f.IsValid() && f.CanSet() {
		f.SetString(workspaceID)
	}
	if f := reqPtr.Elem().FieldByName("Now"); f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(now))
	}
	results := m.Call([]reflect.Value{reflect.ValueOf(ctx), reqPtr})
	if len(results) < 2 {
		return reflect.Value{}, nil
	}
	if !results[1].IsNil() {
		return reflect.Value{}, results[1].Interface().(error)
	}
	resp := results[0]
	if resp.Kind() == reflect.Ptr && !resp.IsNil() {
		return resp.Elem(), nil
	}
	return resp, nil
}

// int64FieldE reads an int64 field by name from a reflect.Value (struct).
func int64FieldE(v reflect.Value, name string) int64 {
	if !v.IsValid() {
		return 0
	}
	f := v.FieldByName(name)
	if !f.IsValid() {
		return 0
	}
	return f.Int()
}

// mapStringInt64FieldE reads a map[string]int64 field by name.
func mapStringInt64FieldE(v reflect.Value, name string) map[string]int64 {
	if !v.IsValid() {
		return nil
	}
	f := v.FieldByName(name)
	if !f.IsValid() || f.IsNil() {
		return nil
	}
	if m, ok := f.Interface().(map[string]int64); ok {
		return m
	}
	return nil
}

// ---------------------------------------------------------------------------
// Location dashboard wiring
// ---------------------------------------------------------------------------

// wireLocationDashboard sets locationDeps.GetLocationDashboardPageData if
// useCases.Entity.LocationDashboard is non-nil.
func wireLocationDashboard(deps *locationmod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Entity == nil || useCases.Entity.LocationDashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Entity.LocationDashboard)
	deps.GetLocationDashboardPageData = func(ctx context.Context) (*locationdashboard.LocationDashboardData, error) {
		workspaceID := consumer.GetWorkspaceIDFromContext(ctx)
		resp, err := callEntydadExecute(uc, ctx, workspaceID, time.Now())
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		// Stats sub-struct
		var total, active, regions, areas int64
		if stats := resp.FieldByName("Stats"); stats.IsValid() {
			total = stats.FieldByName("TotalLocations").Int()
			active = stats.FieldByName("ActiveLocations").Int()
			regions = stats.FieldByName("RegionsCount").Int()
			areas = stats.FieldByName("AreasCount").Int()
		}
		// TopAreas: []LocationAreaCount — same field names, different package type
		var topAreas []locationdashboard.LocationAreaCount
		if f := resp.FieldByName("TopAreas"); f.IsValid() && !f.IsNil() {
			for i := 0; i < f.Len(); i++ {
				s := f.Index(i)
				topAreas = append(topAreas, locationdashboard.LocationAreaCount{
					LocationAreaID:   s.FieldByName("LocationAreaID").String(),
					LocationAreaName: s.FieldByName("LocationAreaName").String(),
					LocationCount:    s.FieldByName("LocationCount").Int(),
				})
			}
		}
		// RecentLocations: []*locationpb.Location — same proto type
		var recentLocations []*locationpb.Location
		if f := resp.FieldByName("RecentLocations"); f.IsValid() && !f.IsNil() {
			if v, ok := f.Interface().([]*locationpb.Location); ok {
				recentLocations = v
			}
		}
		return &locationdashboard.LocationDashboardData{
			TotalLocations:    total,
			ActiveLocations:   active,
			RegionsCount:      regions,
			AreasCount:        areas,
			LocationsByRegion: mapStringInt64FieldE(resp, "LocationsByRegion"),
			TopAreas:          topAreas,
			RecentLocations:   recentLocations,
		}, nil
	}
}

// ---------------------------------------------------------------------------
// Admin dashboard wiring
// ---------------------------------------------------------------------------

// wireAdminDashboard constructs the admin module deps and wires the
// GetDashboardData callback if useCases.Entity.AdminDashboard is non-nil.
func wireAdminDashboard(deps *adminmod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Entity == nil || useCases.Entity.AdminDashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Entity.AdminDashboard)
	deps.GetDashboardData = func(ctx context.Context) (*admindashboard.AdminDashboardData, error) {
		workspaceID := consumer.GetWorkspaceIDFromContext(ctx)
		resp, err := callEntydadExecute(uc, ctx, workspaceID, time.Now())
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		// Stats sub-struct
		var wuCount, roleCount, permCount, recentChanges int64
		if stats := resp.FieldByName("Stats"); stats.IsValid() {
			wuCount = stats.FieldByName("WorkspaceUsers").Int()
			roleCount = stats.FieldByName("Roles").Int()
			permCount = stats.FieldByName("Permissions").Int()
			recentChanges = stats.FieldByName("RecentRoleChanges7d").Int()
		}
		// TopRolesByPerms: []RolePermissionCount
		var topRoles []admindashboard.RolePermissionCount
		if f := resp.FieldByName("TopRolesByPerms"); f.IsValid() && !f.IsNil() {
			for i := 0; i < f.Len(); i++ {
				s := f.Index(i)
				topRoles = append(topRoles, admindashboard.RolePermissionCount{
					RoleID:          s.FieldByName("RoleID").String(),
					RoleName:        s.FieldByName("RoleName").String(),
					PermissionCount: s.FieldByName("PermissionCount").Int(),
				})
			}
		}
		// RecentAssignments: []*workspaceuserrolepb.WorkspaceUserRole — same proto type
		var recentAssignments []*workspaceuserrolepb.WorkspaceUserRole
		if f := resp.FieldByName("RecentAssignments"); f.IsValid() && !f.IsNil() {
			if v, ok := f.Interface().([]*workspaceuserrolepb.WorkspaceUserRole); ok {
				recentAssignments = v
			}
		}
		return &admindashboard.AdminDashboardData{
			WorkspaceUsers:      wuCount,
			Roles:               roleCount,
			Permissions:         permCount,
			RecentRoleChanges7d: recentChanges,
			UsersPerRole:        mapStringInt64FieldE(resp, "UsersPerRole"),
			TopRolesByPerms:     topRoles,
			RecentAssignments:   recentAssignments,
			// RoleNamesByID and UserLabelsByID: deferred — requires additional
			// enrichment queries (workspace_user + role lookups by ID). Currently
			// nil; admin dashboard renders IDs until enrichment is wired.
		}, nil
	}
}

package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	locationarea "github.com/erniealice/entydad-golang/domain/entity/location/location_area"
	"github.com/erniealice/entydad-golang/domain/entity/location/location_area/form"
)

// LocationAreaRecord holds the data for a single location area, used by the
// ReadLocationArea callback.
type LocationAreaRecord struct {
	ID          string
	Name        string
	Description string
	Active      bool
}

// Deps holds dependencies for location area action handlers.
type Deps struct {
	CreateLocationArea    func(ctx context.Context, name, description string, active bool) (string, error)
	ReadLocationArea      func(ctx context.Context, id string) (*LocationAreaRecord, error)
	UpdateLocationArea    func(ctx context.Context, id, name, description string, active bool) error
	DeleteLocationArea    func(ctx context.Context, id string) error
	SetLocationAreaActive func(ctx context.Context, id string, active bool) error
	GetInUseIDs           func(ctx context.Context, ids []string) (map[string]bool, error)
	Routes                locationarea.Routes
	Labels                locationarea.Labels
}

// NewAddAction creates the location area add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location_area", "create") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("location-area-drawer-form", &form.Data{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				Labels:       deps.Labels.Form,
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- create location area
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		name := r.FormValue("name")
		description := r.FormValue("description")

		_, err := deps.CreateLocationArea(ctx, name, description, active)
		if err != nil {
			log.Printf("Failed to create location area: %v", err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("location-areas-table")
	})
}

// NewEditAction creates the location area edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location_area", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			rec, err := deps.ReadLocationArea(ctx, id)
			if err != nil {
				log.Printf("Failed to read location area %s: %v", id, err)
				return view.HTMXError(viewCtx.T("shared.errors.notFound"))
			}

			return view.OK("location-area-drawer-form", &form.Data{
				FormAction:   route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:       true,
				ID:           id,
				Name:         rec.Name,
				Description:  rec.Description,
				Active:       rec.Active,
				Labels:       deps.Labels.Form,
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- update location area
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		err := deps.UpdateLocationArea(ctx, id, r.FormValue("name"), r.FormValue("description"), active)
		if err != nil {
			log.Printf("Failed to update location area %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("location-areas-table")
	})
}

// NewDeleteAction creates the location area delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location_area", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}

		// Server-side re-check: ensure location area is not in use
		if deps.GetInUseIDs != nil {
			inUse, err := deps.GetInUseIDs(ctx, []string{id})
			if err != nil {
				log.Printf("Failed to check location area in-use status: %v", err)
				return view.HTMXError(viewCtx.T("shared.errors.verifyFailed"))
			}
			if inUse[id] {
				return view.HTMXError(viewCtx.T("shared.errors.cannotDeleteInUse"))
			}
		}

		if err := deps.DeleteLocationArea(ctx, id); err != nil {
			log.Printf("Failed to delete location area %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("location-areas-table")
	})
}

// NewBulkDeleteAction creates the location area bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location_area", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}

		// Server-side re-check: ensure none of the location areas are in use
		if deps.GetInUseIDs != nil {
			inUse, err := deps.GetInUseIDs(ctx, ids)
			if err != nil {
				log.Printf("Failed to check location areas in-use status: %v", err)
				return view.HTMXError(viewCtx.T("shared.errors.verifyFailed"))
			}
			for _, id := range ids {
				if inUse[id] {
					return view.HTMXError(viewCtx.T("shared.errors.cannotDeleteInUse"))
				}
			}
		}

		for _, id := range ids {
			if err := deps.DeleteLocationArea(ctx, id); err != nil {
				log.Printf("Failed to delete location area %s: %v", id, err)
			}
		}

		return view.HTMXSuccess("location-areas-table")
	})
}

// NewSetStatusAction creates the location area activate/deactivate action (POST only).
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location_area", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return view.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return view.HTMXError(viewCtx.T("shared.errors.invalidStatus"))
		}

		if err := deps.SetLocationAreaActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update location area status %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("location-areas-table")
	})
}

// NewBulkSetStatusAction creates the location area bulk activate/deactivate action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location_area", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return view.HTMXError(viewCtx.T("shared.errors.invalidTargetStatus"))
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetLocationAreaActive(ctx, id, active); err != nil {
				log.Printf("Failed to update location area status %s: %v", id, err)
			}
		}

		return view.HTMXSuccess("location-areas-table")
	})
}

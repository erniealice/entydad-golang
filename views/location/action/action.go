package action

import (
	"context"
	"log"
	"net/http"
	"sort"

	"github.com/erniealice/pyeza-golang/route"
	pyeza "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"

	"github.com/erniealice/entydad-golang"
	locationform "github.com/erniealice/entydad-golang/views/location/form"
)

// LocationAreaOption is a single location area available for selection.
type LocationAreaOption struct {
	ID   string
	Name string
}

// Deps holds dependencies for location action handlers.
type Deps struct {
	CreateLocation    func(ctx context.Context, req *locationpb.CreateLocationRequest) (*locationpb.CreateLocationResponse, error)
	ReadLocation      func(ctx context.Context, req *locationpb.ReadLocationRequest) (*locationpb.ReadLocationResponse, error)
	UpdateLocation    func(ctx context.Context, req *locationpb.UpdateLocationRequest) (*locationpb.UpdateLocationResponse, error)
	DeleteLocation    func(ctx context.Context, req *locationpb.DeleteLocationRequest) (*locationpb.DeleteLocationResponse, error)
	SetLocationActive func(ctx context.Context, id string, active bool) error
	GetInUseIDs       func(ctx context.Context, ids []string) (map[string]bool, error)
	// ListLocationAreas loads active location areas for the area dropdown.
	// If nil, the area field is omitted from the form.
	ListLocationAreas func(ctx context.Context) ([]LocationAreaOption, error)
	Routes            entydad.LocationRoutes
	Labels            entydad.LocationLabels
}

// buildAreaSelectOptions converts location area options to the select component format.
// Options are sorted alphabetically by label.
func buildAreaSelectOptions(areas []LocationAreaOption, selectedID string) []pyeza.SelectOption {
	sort.Slice(areas, func(i, j int) bool {
		return areas[i].Name < areas[j].Name
	})
	opts := make([]pyeza.SelectOption, 0, len(areas))
	for _, a := range areas {
		opts = append(opts, pyeza.SelectOption{
			Value:    a.ID,
			Label:    a.Name,
			Selected: a.ID == selectedID,
		})
	}
	return opts
}

// loadAreaOptions fetches location areas if the dep is present; returns nil on no dep or error.
func loadAreaOptions(ctx context.Context, deps *Deps, selectedID string) []pyeza.SelectOption {
	if deps.ListLocationAreas == nil {
		return nil
	}
	areas, err := deps.ListLocationAreas(ctx)
	if err != nil {
		log.Printf("Failed to list location areas: %v", err)
		return nil
	}
	return buildAreaSelectOptions(areas, selectedID)
}

// NewAddAction creates the location add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location", "create") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			areaOpts := loadAreaOptions(ctx, deps, "")
			return view.OK("location-drawer-form", &locationform.Data{
				FormAction:                deps.Routes.AddURL,
				Active:                    true,
				Timezone:                  "Asia/Manila",
				LocationAreaSelectOptions: areaOpts,
				Labels:                    deps.Labels.Form,
				CommonLabels:              nil, // injected by ViewAdapter
			})
		}

		// POST -- create location
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")
		tz := r.FormValue("timezone")
		if tz == "" {
			tz = "Asia/Manila"
		}
		areaID := r.FormValue("location_area_id")
		locData := &locationpb.Location{
			Name:        r.FormValue("name"),
			Address:     r.FormValue("address"),
			Description: &desc,
			Timezone:    &tz,
			Active:      active,
		}
		if areaID != "" {
			locData.LocationAreaId = &areaID
		}
		_, err := deps.CreateLocation(ctx, &locationpb.CreateLocationRequest{
			Data: locData,
		})
		if err != nil {
			log.Printf("Failed to create location: %v", err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

// NewEditAction creates the location edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadLocation(ctx, &locationpb.ReadLocationRequest{
				Data: &locationpb.Location{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read location %s: %v", id, err)
				return entydad.HTMXError(viewCtx.T("shared.errors.notFound"))
			}

			loc := resp.GetData()[0]

			tz := loc.GetTimezone()
			if tz == "" {
				tz = "Asia/Manila"
			}
			selectedAreaID := loc.GetLocationAreaId()
			areaOpts := loadAreaOptions(ctx, deps, selectedAreaID)
			return view.OK("location-drawer-form", &locationform.Data{
				FormAction:                route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                    true,
				ID:                        id,
				Name:                      loc.GetName(),
				Address:                   loc.GetAddress(),
				Description:               loc.GetDescription(),
				Timezone:                  tz,
				Active:                    loc.GetActive(),
				SelectedLocationAreaID:    selectedAreaID,
				LocationAreaSelectOptions: areaOpts,
				Labels:                    deps.Labels.Form,
				CommonLabels:              nil, // injected by ViewAdapter
			})
		}

		// POST -- update location
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")
		tz := r.FormValue("timezone")
		if tz == "" {
			tz = "Asia/Manila"
		}
		areaID := r.FormValue("location_area_id")
		locData := &locationpb.Location{
			Id:          id,
			Name:        r.FormValue("name"),
			Address:     r.FormValue("address"),
			Description: &desc,
			Timezone:    &tz,
			Active:      active,
		}
		if areaID != "" {
			locData.LocationAreaId = &areaID
		}
		_, err := deps.UpdateLocation(ctx, &locationpb.UpdateLocationRequest{
			Data: locData,
		})
		if err != nil {
			log.Printf("Failed to update location %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

// NewDeleteAction creates the location delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location", "delete") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}

		// Server-side re-check: ensure location is not in use
		if deps.GetInUseIDs != nil {
			inUse, err := deps.GetInUseIDs(ctx, []string{id})
			if err != nil {
				log.Printf("Failed to check location in-use status: %v", err)
				return entydad.HTMXError(viewCtx.T("shared.errors.verifyFailed"))
			}
			if inUse[id] {
				return entydad.HTMXError(viewCtx.T("shared.errors.cannotDeleteInUse"))
			}
		}

		_, err := deps.DeleteLocation(ctx, &locationpb.DeleteLocationRequest{
			Data: &locationpb.Location{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete location %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

// NewBulkDeleteAction creates the location bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location", "delete") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}

		// Server-side re-check: ensure none of the locations are in use
		if deps.GetInUseIDs != nil {
			inUse, err := deps.GetInUseIDs(ctx, ids)
			if err != nil {
				log.Printf("Failed to check locations in-use status: %v", err)
				return entydad.HTMXError(viewCtx.T("shared.errors.verifyFailed"))
			}
			for _, id := range ids {
				if inUse[id] {
					return entydad.HTMXError(viewCtx.T("shared.errors.cannotDeleteInUse"))
				}
			}
		}

		for _, id := range ids {
			_, err := deps.DeleteLocation(ctx, &locationpb.DeleteLocationRequest{
				Data: &locationpb.Location{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete location %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

// NewSetStatusAction creates the location activate/deactivate action (POST only).
// Expects query params: ?id={locationId}&status={active|inactive}
//
// Uses SetLocationActive (raw map update) instead of UpdateLocation (protobuf) because
// proto3's protojson omits bool fields with value false, which means
// deactivation (active=false) would silently be skipped.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return entydad.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidStatus"))
		}

		if err := deps.SetLocationActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update location status %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

// NewBulkSetStatusAction creates the location bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("location", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return entydad.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidTargetStatus"))
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetLocationActive(ctx, id, active); err != nil {
				log.Printf("Failed to update location status %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

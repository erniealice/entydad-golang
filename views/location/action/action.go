package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"

	"github.com/erniealice/entydad-golang"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Name                   string
	NamePlaceholder        string
	Address                string
	AddressPlaceholder     string
	Description            string
	DescriptionPlaceholder string
	Active                 string
}

// FormData is the template data for the location drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Address      string
	Description  string
	Active       bool
	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for location action handlers.
type Deps struct {
	CreateLocation    func(ctx context.Context, req *locationpb.CreateLocationRequest) (*locationpb.CreateLocationResponse, error)
	ReadLocation      func(ctx context.Context, req *locationpb.ReadLocationRequest) (*locationpb.ReadLocationResponse, error)
	UpdateLocation    func(ctx context.Context, req *locationpb.UpdateLocationRequest) (*locationpb.UpdateLocationResponse, error)
	DeleteLocation    func(ctx context.Context, req *locationpb.DeleteLocationRequest) (*locationpb.DeleteLocationResponse, error)
	SetLocationActive func(ctx context.Context, id string, active bool) error
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Name:                   t("form.name"),
		NamePlaceholder:        t("form.namePlaceholder"),
		Address:                t("form.address"),
		AddressPlaceholder:     t("form.addressPlaceholder"),
		Description:            t("form.description"),
		DescriptionPlaceholder: t("form.descriptionPlaceholder"),
		Active:                 t("form.active"),
	}
}

// NewAddAction creates the location add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("location-drawer-form", &FormData{
				FormAction:   "/action/locations/add",
				Active:       true,
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- create location
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")

		_, err := deps.CreateLocation(ctx, &locationpb.CreateLocationRequest{
			Data: &locationpb.Location{
				Name:        r.FormValue("name"),
				Address:     r.FormValue("address"),
				Description: &desc,
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to create location: %v", err)
			return entydad.HTMXError("Failed to create location")
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

// NewEditAction creates the location edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadLocation(ctx, &locationpb.ReadLocationRequest{
				Data: &locationpb.Location{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read location %s: %v", id, err)
				return entydad.HTMXError("Location not found")
			}

			loc := resp.GetData()[0]

			return view.OK("location-drawer-form", &FormData{
				FormAction:   "/action/locations/edit/" + id,
				IsEdit:       true,
				ID:           id,
				Name:         loc.GetName(),
				Address:      loc.GetAddress(),
				Description:  loc.GetDescription(),
				Active:       loc.GetActive(),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- update location
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")

		_, err := deps.UpdateLocation(ctx, &locationpb.UpdateLocationRequest{
			Data: &locationpb.Location{
				Id:          id,
				Name:        r.FormValue("name"),
				Address:     r.FormValue("address"),
				Description: &desc,
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to update location %s: %v", id, err)
			return entydad.HTMXError("Failed to update location")
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

// NewDeleteAction creates the location delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError("Location ID is required")
		}

		_, err := deps.DeleteLocation(ctx, &locationpb.DeleteLocationRequest{
			Data: &locationpb.Location{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete location %s: %v", id, err)
			return entydad.HTMXError("Failed to delete location")
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

// NewBulkDeleteAction creates the location bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError("No location IDs provided")
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
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return entydad.HTMXError("Location ID is required")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid status")
		}

		if err := deps.SetLocationActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update location status %s: %v", id, err)
			return entydad.HTMXError("Failed to update location status")
		}

		return entydad.HTMXSuccess("locations-table")
	})
}

// NewBulkSetStatusAction creates the location bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return entydad.HTMXError("No location IDs provided")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid target status")
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

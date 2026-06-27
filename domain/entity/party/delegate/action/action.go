package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	delegatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/delegate"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	entitydelegate "github.com/erniealice/entydad-golang/domain/entity/party/delegate"
	delegateform "github.com/erniealice/entydad-golang/domain/entity/party/delegate/form"
)

// Deps holds dependencies for delegate action handlers.
// No payment-terms/categories/tags/status — delegate has only active bool.
type Deps struct {
	Routes         entitydelegate.Routes
	CreateDelegate func(ctx context.Context, req *delegatepb.CreateDelegateRequest) (*delegatepb.CreateDelegateResponse, error)
	ReadDelegate   func(ctx context.Context, req *delegatepb.ReadDelegateRequest) (*delegatepb.ReadDelegateResponse, error)
	UpdateDelegate func(ctx context.Context, req *delegatepb.UpdateDelegateRequest) (*delegatepb.UpdateDelegateResponse, error)
	DeleteDelegate func(ctx context.Context, req *delegatepb.DeleteDelegateRequest) (*delegatepb.DeleteDelegateResponse, error)
}

// NewAddAction creates the delegate add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("delegate", "create") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			labels := delegateform.BuildLabels(viewCtx.T)
			return view.OK("delegate-drawer-form", &delegateform.Data{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				Labels:       labels,
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create delegate
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		repUser := &userpb.User{
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			EmailAddress: r.FormValue("email_address"),
			MobileNumber: r.FormValue("mobile_number"),
			Active:       true,
		}

		_, err := deps.CreateDelegate(ctx, &delegatepb.CreateDelegateRequest{
			Data: &delegatepb.Delegate{
				User:   repUser,
				Active: true,
			},
		})
		if err != nil {
			log.Printf("Failed to create delegate: %v", err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("delegates-table")
	})
}

// NewEditAction creates the delegate edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("delegate", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadDelegate(ctx, &delegatepb.ReadDelegateRequest{
				Data: &delegatepb.Delegate{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read delegate %s: %v", id, err)
				return view.HTMXError(viewCtx.T("shared.errors.notFound"))
			}

			d := resp.GetData()[0]
			u := d.GetUser()
			labels := delegateform.BuildLabels(viewCtx.T)
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id)
			return view.OK("delegate-drawer-form", &delegateform.Data{
				FormAction:   formAction,
				IsEdit:       true,
				ID:           id,
				FirstName:    u.GetFirstName(),
				LastName:     u.GetLastName(),
				Email:        u.GetEmailAddress(),
				Mobile:       u.GetMobileNumber(),
				Active:       d.GetActive(),
				Labels:       labels,
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — update delegate
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		_, err := deps.UpdateDelegate(ctx, &delegatepb.UpdateDelegateRequest{
			Data: &delegatepb.Delegate{
				Id: id,
				User: &userpb.User{
					FirstName:    r.FormValue("first_name"),
					LastName:     r.FormValue("last_name"),
					EmailAddress: r.FormValue("email_address"),
					MobileNumber: r.FormValue("mobile_number"),
				},
			},
		})
		if err != nil {
			log.Printf("Failed to update delegate %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("delegates-table")
	})
}

// NewDeleteAction creates the delegate delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("delegate", "delete") {
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

		_, err := deps.DeleteDelegate(ctx, &delegatepb.DeleteDelegateRequest{
			Data: &delegatepb.Delegate{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete delegate %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("delegates-table")
	})
}

// NewBulkDeleteAction creates the delegate bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("delegate", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}

		for _, id := range ids {
			_, err := deps.DeleteDelegate(ctx, &delegatepb.DeleteDelegateRequest{
				Data: &delegatepb.Delegate{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete delegate %s: %v", id, err)
			}
		}

		return view.HTMXSuccess("delegates-table")
	})
}

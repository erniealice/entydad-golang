package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	clientpb "leapfor.xyz/esqyma/golang/v1/domain/entity/client"
	userpb "leapfor.xyz/esqyma/golang/v1/domain/entity/user"

	"leapfor.xyz/entydad"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	FirstName            string
	FirstNamePlaceholder string
	LastName             string
	LastNamePlaceholder  string
	Email                string
	EmailPlaceholder     string
	Mobile               string
	MobilePlaceholder    string
	Active               string
}

// FormData is the template data for the client drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	FirstName    string
	LastName     string
	Email        string
	Mobile       string
	Active       bool
	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for client action handlers.
type Deps struct {
	CreateClient func(ctx context.Context, req *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error)
	ReadClient   func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	UpdateClient func(ctx context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error)
	DeleteClient func(ctx context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error)
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		FirstName:            t("client.form.firstName"),
		FirstNamePlaceholder: t("client.form.firstNamePlaceholder"),
		LastName:             t("client.form.lastName"),
		LastNamePlaceholder:  t("client.form.lastNamePlaceholder"),
		Email:                t("client.form.email"),
		EmailPlaceholder:     t("client.form.emailPlaceholder"),
		Mobile:               t("client.form.phone"),
		MobilePlaceholder:    t("client.form.phonePlaceholder"),
		Active:               t("client.form.active"),
	}
}

// NewAddAction creates the client add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("client-drawer-form", &FormData{
				FormAction:   "/action/clients/add",
				Active:       true,
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create client
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.CreateClient(ctx, &clientpb.CreateClientRequest{
			Data: &clientpb.Client{
				Active: active,
				User: &userpb.User{
					FirstName:    r.FormValue("first_name"),
					LastName:     r.FormValue("last_name"),
					EmailAddress: r.FormValue("email_address"),
					MobileNumber: r.FormValue("mobile_number"),
					Active:       active,
				},
			},
		})
		if err != nil {
			log.Printf("Failed to create client: %v", err)
			return entydad.HTMXError("Failed to create client")
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewEditAction creates the client edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{
				Data: &clientpb.Client{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read client %s: %v", id, err)
				return entydad.HTMXError("Client not found")
			}

			c := resp.GetData()[0]
			u := c.GetUser()

			return view.OK("client-drawer-form", &FormData{
				FormAction:   "/action/clients/edit/" + id,
				IsEdit:       true,
				ID:           id,
				FirstName:    u.GetFirstName(),
				LastName:     u.GetLastName(),
				Email:        u.GetEmailAddress(),
				Mobile:       u.GetMobileNumber(),
				Active:       c.GetActive(),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — update client
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.UpdateClient(ctx, &clientpb.UpdateClientRequest{
			Data: &clientpb.Client{
				Id:     id,
				Active: active,
				User: &userpb.User{
					FirstName:    r.FormValue("first_name"),
					LastName:     r.FormValue("last_name"),
					EmailAddress: r.FormValue("email_address"),
					MobileNumber: r.FormValue("mobile_number"),
					Active:       active,
				},
			},
		})
		if err != nil {
			log.Printf("Failed to update client %s: %v", id, err)
			return entydad.HTMXError("Failed to update client")
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewDeleteAction creates the client delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError("Client ID is required")
		}

		_, err := deps.DeleteClient(ctx, &clientpb.DeleteClientRequest{
			Data: &clientpb.Client{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete client %s: %v", id, err)
			return entydad.HTMXError("Failed to delete client")
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewBulkDeleteAction creates the client bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError("No client IDs provided")
		}

		for _, id := range ids {
			_, err := deps.DeleteClient(ctx, &clientpb.DeleteClientRequest{
				Data: &clientpb.Client{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete client %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

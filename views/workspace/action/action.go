package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"

	"github.com/erniealice/entydad-golang"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Name                   string
	NamePlaceholder        string
	Description            string
	DescriptionPlaceholder string
	Private                string
	Active                 string
}

// FormData is the template data for the workspace drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Private      bool
	Active       bool
	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for workspace action handlers.
type Deps struct {
	CreateWorkspace    func(ctx context.Context, req *workspacepb.CreateWorkspaceRequest) (*workspacepb.CreateWorkspaceResponse, error)
	ReadWorkspace      func(ctx context.Context, req *workspacepb.ReadWorkspaceRequest) (*workspacepb.ReadWorkspaceResponse, error)
	UpdateWorkspace    func(ctx context.Context, req *workspacepb.UpdateWorkspaceRequest) (*workspacepb.UpdateWorkspaceResponse, error)
	DeleteWorkspace    func(ctx context.Context, req *workspacepb.DeleteWorkspaceRequest) (*workspacepb.DeleteWorkspaceResponse, error)
	SetWorkspaceActive func(ctx context.Context, id string, active bool) error
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Name:                   t("form.name"),
		NamePlaceholder:        t("form.namePlaceholder"),
		Description:            t("form.description"),
		DescriptionPlaceholder: t("form.descriptionPlaceholder"),
		Private:                t("form.private"),
		Active:                 t("form.active"),
	}
}

// NewAddAction creates the workspace add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("workspace-drawer-form", &FormData{
				FormAction:   "/action/workspaces/add",
				Active:       true,
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- create workspace
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		private := r.FormValue("private") == "true"

		_, err := deps.CreateWorkspace(ctx, &workspacepb.CreateWorkspaceRequest{
			Data: &workspacepb.Workspace{
				Name:        r.FormValue("name"),
				Description: r.FormValue("description"),
				Private:     private,
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to create workspace: %v", err)
			return entydad.HTMXError("Failed to create workspace")
		}

		return entydad.HTMXSuccess("workspaces-table")
	})
}

// NewEditAction creates the workspace edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadWorkspace(ctx, &workspacepb.ReadWorkspaceRequest{
				Data: &workspacepb.Workspace{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read workspace %s: %v", id, err)
				return entydad.HTMXError("Workspace not found")
			}

			ws := resp.GetData()[0]

			return view.OK("workspace-drawer-form", &FormData{
				FormAction:   "/action/workspaces/edit/" + id,
				IsEdit:       true,
				ID:           id,
				Name:         ws.GetName(),
				Description:  ws.GetDescription(),
				Private:      ws.GetPrivate(),
				Active:       ws.GetActive(),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- update workspace
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		private := r.FormValue("private") == "true"

		_, err := deps.UpdateWorkspace(ctx, &workspacepb.UpdateWorkspaceRequest{
			Data: &workspacepb.Workspace{
				Id:          id,
				Name:        r.FormValue("name"),
				Description: r.FormValue("description"),
				Private:     private,
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to update workspace %s: %v", id, err)
			return entydad.HTMXError("Failed to update workspace")
		}

		return entydad.HTMXSuccess("workspaces-table")
	})
}

// NewDeleteAction creates the workspace delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError("Workspace ID is required")
		}

		_, err := deps.DeleteWorkspace(ctx, &workspacepb.DeleteWorkspaceRequest{
			Data: &workspacepb.Workspace{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete workspace %s: %v", id, err)
			return entydad.HTMXError("Failed to delete workspace")
		}

		return entydad.HTMXSuccess("workspaces-table")
	})
}

// NewBulkDeleteAction creates the workspace bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError("No workspace IDs provided")
		}

		for _, id := range ids {
			_, err := deps.DeleteWorkspace(ctx, &workspacepb.DeleteWorkspaceRequest{
				Data: &workspacepb.Workspace{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete workspace %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("workspaces-table")
	})
}

// NewSetStatusAction creates the workspace activate/deactivate action (POST only).
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
			return entydad.HTMXError("Workspace ID is required")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid status")
		}

		if err := deps.SetWorkspaceActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update workspace status %s: %v", id, err)
			return entydad.HTMXError("Failed to update workspace status")
		}

		return entydad.HTMXSuccess("workspaces-table")
	})
}

// NewBulkSetStatusAction creates the workspace bulk activate/deactivate action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return entydad.HTMXError("No workspace IDs provided")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid target status")
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetWorkspaceActive(ctx, id, active); err != nil {
				log.Printf("Failed to update workspace status %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("workspaces-table")
	})
}

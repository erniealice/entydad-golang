package action

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
)

// slugify converts a name into a lowercase, hyphenated code suitable for the
// Category.Code field (e.g. "VIP Supplier" -> "vip-supplier").
var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonAlphaNum.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// TagFormLabels holds i18n labels for the supplier tag drawer form.
type TagFormLabels struct {
	TagName                string
	Code                   string
	CodeAutoPlaceholder    string
	Description            string
	DescriptionPlaceholder string
	Active                 string

	// Field-level info text surfaced via an info button beside each label.
	TagNameInfo     string
	CodeInfo        string
	DescriptionInfo string
	ActiveInfo      string
}

// FormData is the template data for the tag drawer form.
type FormData struct {
	FormAction   string
	WorkspaceID   string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit       bool
	ID           string
	Name         string
	Code         string
	Description  string
	Active       bool
	Labels       TagFormLabels
	CommonLabels any
}

// tagFormLabels builds i18n labels for the tag drawer form using dot-notation keys.
func tagFormLabels(t func(string) string) TagFormLabels {
	return TagFormLabels{
		TagName:                t("supplier.tagForm.tagName"),
		Code:                   t("supplier.tagForm.code"),
		CodeAutoPlaceholder:    t("supplier.tagForm.codeAutoPlaceholder"),
		Description:            t("supplier.tagForm.description"),
		DescriptionPlaceholder: t("supplier.tagForm.descriptionPlaceholder"),
		Active:                 t("supplier.tagForm.active"),
		TagNameInfo:            t("supplier.tagForm.tagNameInfo"),
		CodeInfo:               t("supplier.tagForm.codeInfo"),
		DescriptionInfo:        t("supplier.tagForm.descriptionInfo"),
		ActiveInfo:             t("supplier.tagForm.activeInfo"),
	}
}

// resolveCode returns the explicit code if provided, otherwise slugifies the name.
func resolveCode(code, name string) string {
	code = strings.TrimSpace(code)
	if code != "" {
		return slugify(code)
	}
	return slugify(name)
}

// Deps holds dependencies for supplier tag action handlers.
type Deps struct {
	Routes            entydad.SupplierTagRoutes
	CommonLabels      pyeza.CommonLabels
	ListCategories    func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	CreateCategory    func(ctx context.Context, req *categorypb.CreateCategoryRequest) (*categorypb.CreateCategoryResponse, error)
	ReadCategory      func(ctx context.Context, req *categorypb.ReadCategoryRequest) (*categorypb.ReadCategoryResponse, error)
	UpdateCategory    func(ctx context.Context, req *categorypb.UpdateCategoryRequest) (*categorypb.UpdateCategoryResponse, error)
	DeleteCategory    func(ctx context.Context, req *categorypb.DeleteCategoryRequest) (*categorypb.DeleteCategoryResponse, error)
	SetCategoryActive func(ctx context.Context, id string, active bool) error
}

// isDuplicateTagName checks if a tag name already exists among supplier-module categories,
// optionally excluding a specific category ID (for edit operations).
func isDuplicateTagName(ctx context.Context, listFn func(context.Context, *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error), name string, excludeID string) bool {
	if listFn == nil {
		return false
	}
	resp, err := listFn(ctx, &categorypb.ListCategoriesRequest{})
	if err != nil {
		log.Printf("Failed to list categories for duplicate check: %v", err)
		return false
	}
	for _, cat := range resp.GetData() {
		if cat.GetModule() != "supplier" {
			continue
		}
		if excludeID != "" && cat.GetId() == excludeID {
			continue
		}
		if strings.EqualFold(cat.GetName(), name) {
			return true
		}
	}
	return false
}

// NewAddAction creates the supplier tag add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("supplier-tag-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				Labels:       tagFormLabels(viewCtx.T),
				CommonLabels: deps.CommonLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		name := r.FormValue("name")
		active := r.FormValue("active") == "true"

		if isDuplicateTagName(ctx, deps.ListCategories, name, "") {
			return entydad.HTMXError(viewCtx.T("shared.errors.tagNameExists"))
		}

		_, err := deps.CreateCategory(ctx, &categorypb.CreateCategoryRequest{
			Data: &categorypb.Category{
				Name:        name,
				Code:        resolveCode(r.FormValue("code"), name),
				Description: r.FormValue("description"),
				Module:      "supplier",
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to create supplier tag: %v", err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("supplier-tags-table")
	})
}

// NewEditAction creates the supplier tag edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadCategory(ctx, &categorypb.ReadCategoryRequest{
				Data: &categorypb.Category{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read supplier tag %s: %v", id, err)
				return entydad.HTMXError(viewCtx.T("shared.errors.notFound"))
			}

			data := resp.GetData()
			if len(data) == 0 {
				return entydad.HTMXError(viewCtx.T("shared.errors.notFound"))
			}
			cat := data[0]

			return view.OK("supplier-tag-drawer-form", &FormData{
				FormAction:   route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:       true,
				ID:           id,
				Name:         cat.GetName(),
				Code:         cat.GetCode(),
				Description:  cat.GetDescription(),
				Active:       cat.GetActive(),
				Labels:       tagFormLabels(viewCtx.T),
				CommonLabels: deps.CommonLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		name := r.FormValue("name")
		active := r.FormValue("active") == "true"

		if isDuplicateTagName(ctx, deps.ListCategories, name, id) {
			return entydad.HTMXError(viewCtx.T("shared.errors.tagNameExists"))
		}

		_, err := deps.UpdateCategory(ctx, &categorypb.UpdateCategoryRequest{
			Data: &categorypb.Category{
				Id:          id,
				Name:        name,
				Code:        resolveCode(r.FormValue("code"), name),
				Description: r.FormValue("description"),
				Module:      "supplier",
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to update supplier tag %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("supplier-tags-table")
	})
}

// NewDeleteAction creates the supplier tag delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
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

		_, err := deps.DeleteCategory(ctx, &categorypb.DeleteCategoryRequest{
			Data: &categorypb.Category{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete supplier tag %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("supplier-tags-table")
	})
}

// NewBulkDeleteAction creates the supplier tag bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}

		for _, id := range ids {
			_, err := deps.DeleteCategory(ctx, &categorypb.DeleteCategoryRequest{
				Data: &categorypb.Category{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete supplier tag %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("supplier-tags-table")
	})
}

// NewSetStatusAction creates the supplier tag activate/deactivate action (POST only).
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
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

		if err := deps.SetCategoryActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update supplier tag status %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("supplier-tags-table")
	})
}

// NewBulkSetStatusAction creates the supplier tag bulk activate/deactivate action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
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
			if err := deps.SetCategoryActive(ctx, id, active); err != nil {
				log.Printf("Failed to update supplier tag status %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("supplier-tags-table")
	})
}

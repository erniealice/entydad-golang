package action

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
)

// slugify converts a name into a lowercase, hyphenated code suitable for the
// Category.Code field (e.g. "VIP Customer" -> "vip-customer").
var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonAlphaNum.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// FormData is the template data for the tag drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Code         string
	Description  string
	Active       bool
	CommonLabels any
}

// resolveCode returns the explicit code if provided, otherwise slugifies the name.
func resolveCode(code, name string) string {
	code = strings.TrimSpace(code)
	if code != "" {
		return slugify(code)
	}
	return slugify(name)
}

// Deps holds dependencies for client tag action handlers.
type Deps struct {
	ListCategories func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	CreateCategory func(ctx context.Context, req *categorypb.CreateCategoryRequest) (*categorypb.CreateCategoryResponse, error)
	ReadCategory   func(ctx context.Context, req *categorypb.ReadCategoryRequest) (*categorypb.ReadCategoryResponse, error)
	UpdateCategory func(ctx context.Context, req *categorypb.UpdateCategoryRequest) (*categorypb.UpdateCategoryResponse, error)
	DeleteCategory func(ctx context.Context, req *categorypb.DeleteCategoryRequest) (*categorypb.DeleteCategoryResponse, error)
}

// isDuplicateTagName checks if a tag name already exists among client-module categories,
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
		if cat.GetModule() != "client" {
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

// NewAddAction creates the tag add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("client-tag-drawer-form", &FormData{
				FormAction:   "/action/clients/tags/add",
				Active:       true,
				CommonLabels: nil,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		name := r.FormValue("name")
		active := r.FormValue("active") == "true"

		if isDuplicateTagName(ctx, deps.ListCategories, name, "") {
			return entydad.HTMXError("Tag name already exists")
		}

		_, err := deps.CreateCategory(ctx, &categorypb.CreateCategoryRequest{
			Data: &categorypb.Category{
				Name:        name,
				Code:        resolveCode(r.FormValue("code"), name),
				Description: r.FormValue("description"),
				Module:      "client",
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to create client tag: %v", err)
			return entydad.HTMXError("Failed to create tag")
		}

		return entydad.HTMXSuccess("client-tags-table")
	})
}

// NewEditAction creates the tag edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadCategory(ctx, &categorypb.ReadCategoryRequest{
				Data: &categorypb.Category{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read client tag %s: %v", id, err)
				return entydad.HTMXError("Tag not found")
			}

			data := resp.GetData()
			if len(data) == 0 {
				return entydad.HTMXError("Tag not found")
			}
			cat := data[0]

			return view.OK("client-tag-drawer-form", &FormData{
				FormAction:   "/action/clients/tags/edit/" + id,
				IsEdit:       true,
				ID:           id,
				Name:         cat.GetName(),
				Code:         cat.GetCode(),
				Description:  cat.GetDescription(),
				Active:       cat.GetActive(),
				CommonLabels: nil,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		name := r.FormValue("name")
		active := r.FormValue("active") == "true"

		if isDuplicateTagName(ctx, deps.ListCategories, name, id) {
			return entydad.HTMXError("Tag name already exists")
		}

		_, err := deps.UpdateCategory(ctx, &categorypb.UpdateCategoryRequest{
			Data: &categorypb.Category{
				Id:          id,
				Name:        name,
				Code:        resolveCode(r.FormValue("code"), name),
				Description: r.FormValue("description"),
				Module:      "client",
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to update client tag %s: %v", id, err)
			return entydad.HTMXError("Failed to update tag")
		}

		return entydad.HTMXSuccess("client-tags-table")
	})
}

// NewDeleteAction creates the tag delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError("Tag ID is required")
		}

		_, err := deps.DeleteCategory(ctx, &categorypb.DeleteCategoryRequest{
			Data: &categorypb.Category{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete client tag %s: %v", id, err)
			return entydad.HTMXError("Failed to delete tag")
		}

		return entydad.HTMXSuccess("client-tags-table")
	})
}

// NewBulkDeleteAction creates the tag bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError("No tag IDs provided")
		}

		for _, id := range ids {
			_, err := deps.DeleteCategory(ctx, &categorypb.DeleteCategoryRequest{
				Data: &categorypb.Category{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete client tag %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("client-tags-table")
	})
}

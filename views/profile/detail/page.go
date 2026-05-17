package detail

import (
	"context"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds dependencies needed to build the profile detail view.
type ModuleDeps struct {
	Messages map[string]string
}

// PageData carries the rendering context for the profile page.
type PageData struct {
	types.PageData
}

// NewView creates the profile detail view (full page — no tabs).
func NewView(deps *ModuleDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("user", "read") {
			return view.Forbidden("user:read")
		}

		titleKey := "memberPages.section.profile.title"
		iconKey := "memberPages.section.profile.icon"

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:    viewCtx.CacheVersion,
				Title:           lookup(deps.Messages, titleKey, "Profile"),
				CurrentPath:     viewCtx.CurrentPath,
				ActiveNav:       "home",
				ContentTemplate: "profile-page-content",
				HeaderTitle:     lookup(deps.Messages, titleKey, "Profile"),
				HeaderIcon:      lookup(deps.Messages, iconKey, "icon-user"),
				Messages:        deps.Messages,
			},
		}
		return view.OK("profile-page", pageData)
	})
}

func lookup(messages map[string]string, key, fallback string) string {
	if messages != nil {
		if v, ok := messages[key]; ok && v != "" {
			return v
		}
	}
	return fallback
}

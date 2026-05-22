// Package profile_overview is the /me/profile-overview stub view (Phase P9b).
//
// Cross-workspace profile overview: User basics + binding list across all
// workspaces. Real backing depends on `GetProfileOverview` (espyna use case
// not yet authored). Stub for now — the user's basic identity already
// flows in via Sidebar.CurrentUser; the workspace-binding aggregation
// requires a fresh use case.
package profile_overview

import (
	"context"

	me "github.com/erniealice/entydad-golang/views/me"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps bundles per-request configuration for the view.
type ModuleDeps struct {
	Messages map[string]string
}

// PageData carries rendering context for the /me/profile-overview stub.
type PageData struct {
	me.PageData
	Subtitle     string
	EmptyMessage string
}

// NewView constructs the /me/profile-overview view.
func NewView(deps *ModuleDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pd := &PageData{
			PageData: me.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           "Profile Overview",
					CurrentPath:     viewCtx.CurrentPath,
					ActiveNav:       "profile-overview",
					ContentTemplate: "me-stub-content",
					HeaderTitle:     "Profile Overview",
					HeaderIcon:      "icon-user",
					Messages:        deps.Messages,
				},
			},
			Subtitle:     "Your identity and workspace bindings.",
			EmptyMessage: "Workspace-binding aggregation coming soon — see Sidebar profile button for current identity.",
		}
		return view.OK("me-page", pd)
	})
}

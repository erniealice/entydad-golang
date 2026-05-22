// Package notifications is the /me/notifications stub view (Phase P9b).
//
// Stub — real backing arrives with the cross-workspace notification list
// use case (deferred). Renders empty-state copy until then.
package notifications

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

// PageData carries rendering context for the /me/notifications stub.
type PageData struct {
	me.PageData
	Subtitle     string
	EmptyMessage string
}

// NewView constructs the /me/notifications view.
func NewView(deps *ModuleDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pd := &PageData{
			PageData: me.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           "Notifications",
					CurrentPath:     viewCtx.CurrentPath,
					ActiveNav:       "notifications",
					ContentTemplate: "me-stub-content",
					HeaderTitle:     "Notifications",
					HeaderIcon:      "icon-bell",
					Messages:        deps.Messages,
				},
			},
			Subtitle:     "All notifications across your workspaces.",
			EmptyMessage: "No notifications yet.",
		}
		return view.OK("me-page", pd)
	})
}

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
		title := me.Msg(deps.Messages, "me.notifications.title", "Notifications")
		pd := &PageData{
			PageData: me.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           title,
					CurrentPath:     viewCtx.CurrentPath,
					ActiveNav:       "notifications",
					ContentTemplate: "me-stub-content",
					HeaderTitle:     title,
					HeaderIcon:      "icon-bell",
					Messages:        deps.Messages,
				},
			},
			Subtitle:     me.Msg(deps.Messages, "me.notifications.subtitle", "All notifications across your workspaces."),
			EmptyMessage: me.Msg(deps.Messages, "me.notifications.empty", "No notifications yet."),
		}
		return view.OK("me-page", pd)
	})
}

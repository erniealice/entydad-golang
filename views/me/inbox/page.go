// Package inbox is the /me/inbox stub view (Phase P9b).
//
// Cross-workspace notifications aggregated across the user's bindings.
// Real backing is a follow-up (espyna use case `ListNotificationsAcrossWorkspaces`
// is not yet authored; this view returns an empty list with a "coming soon"
// placeholder until that lands).
package inbox

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

// PageData carries rendering context for the /me/inbox stub.
type PageData struct {
	me.PageData
	Subtitle     string
	EmptyMessage string
}

// NewView constructs the /me/inbox view.
func NewView(deps *ModuleDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pd := &PageData{
			PageData: me.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           "Inbox",
					CurrentPath:     viewCtx.CurrentPath,
					ActiveNav:       "inbox",
					ContentTemplate: "me-stub-content",
					HeaderTitle:     "Inbox",
					HeaderIcon:      "icon-inbox",
					Messages:        deps.Messages,
				},
			},
			Subtitle:     "Notifications aggregated across your workspaces.",
			EmptyMessage: "Cross-workspace notifications coming soon — your inbox is currently empty.",
		}
		return view.OK("me-page", pd)
	})
}

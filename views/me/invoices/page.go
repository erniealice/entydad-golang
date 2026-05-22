// Package invoices is the /me/invoices stub view (Phase P9b).
//
// Cross-workspace invoices: Revenue rows where user is Client, Expenditure
// rows where user is Supplier, TenantInvoice rows where user is workspace
// Operator. Real backing depends on `ListInvoicesAcrossRoles` (espyna use
// case not yet authored) — stub for now.
package invoices

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

// PageData carries rendering context for the /me/invoices stub.
type PageData struct {
	me.PageData
	Subtitle     string
	EmptyMessage string
}

// NewView constructs the /me/invoices view.
func NewView(deps *ModuleDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pd := &PageData{
			PageData: me.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           "My Invoices",
					CurrentPath:     viewCtx.CurrentPath,
					ActiveNav:       "invoices",
					ContentTemplate: "me-stub-content",
					HeaderTitle:     "My Invoices",
					HeaderIcon:      "icon-file-text",
					Messages:        deps.Messages,
				},
			},
			Subtitle:     "Invoices across all workspaces where you are a client, supplier, or operator.",
			EmptyMessage: "Cross-workspace invoice aggregation coming soon.",
		}
		return view.OK("me-page", pd)
	})
}

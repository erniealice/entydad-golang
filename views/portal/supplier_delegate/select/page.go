// Package selectsupplier is the acting-as picker for SUPPLIER_DELEGATE principals.
//
// Shown when a SUPPLIER_DELEGATE principal represents 2+ suppliers. Clicking a card
// POSTs to /action/auth/switch-principal with:
//
//	principal_id=<current delegate grant id>
//	principal_kind=SUPPLIER_DELEGATE
//	acting_as_supplier_id=<chosen supplier id>
package selectsupplier

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// SupplierOption is one selectable supplier card on the picker page.
type SupplierOption struct {
	// SupplierID is the supplier entity id — sent as acting_as_supplier_id.
	SupplierID string
	// DisplayName is the supplier's human-readable name shown on the card.
	DisplayName string
}

// PageData holds data for the acting-as supplier picker page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
	// PrincipalID is the delegate grant id — sent as principal_id on POST.
	PrincipalID   string
	SwitchPostURL string
	Suppliers     []SupplierOption
}

// Deps holds view dependencies for the select view.
type Deps struct {
	Labels             entydad.PortalLabels
	ResolveSuppliers   func(ctx context.Context) []SupplierOption
	ResolvePrincipalID func(ctx context.Context) string
}

// NewView creates the acting-as supplier picker view (GET /portal/supplier-delegate/select).
func NewView(deps *Deps) view.View {
	switchPostURL := "/action/auth/switch-principal"
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.SupplierDelegate.Select.PageTitle == "" {
		labels.SupplierDelegate.Select.PageTitle = "Choose a supplier"
	}
	if labels.SupplierDelegate.Name == "" {
		labels.SupplierDelegate.Name = "Delegate Portal"
	}
	if labels.SupplierDelegate.Select.EmptyState == "" {
		labels.SupplierDelegate.Select.EmptyState = "No supplier accounts found. Contact your administrator."
	}
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		var suppliers []SupplierOption
		if deps != nil && deps.ResolveSuppliers != nil {
			suppliers = deps.ResolveSuppliers(ctx)
		}
		principalID := ""
		if deps != nil && deps.ResolvePrincipalID != nil {
			principalID = deps.ResolvePrincipalID(ctx)
		}

		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.SupplierDelegate.Select.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-supplier-delegate-select-content",
					ActiveNav:       "home",
				},
				PrincipalKind: "supplier-delegate",
				PortalName:    labels.SupplierDelegate.Name,
				ProfileURL:    "/portal/supplier-delegate/profile",
			},
			Labels:        labels,
			PrincipalID:   principalID,
			SwitchPostURL: switchPostURL,
			Suppliers:     suppliers,
		}
		return view.OK("portal-supplier-delegate-select", pageData)
	})
}

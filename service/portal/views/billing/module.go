// Package billing provides the /app/billing page — workspace billing
// (tabs: subscription | payment-method | invoices). Part of the four
// "personal scope" pages accessible from the sidebar bottom profile popover.
//
// Permission gating (Layer 3): workspace:read
package billing

import (
	billingdetail "github.com/erniealice/entydad-golang/service/portal/views/billing/detail"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds dependencies for the billing module.
type ModuleDeps struct {
	Messages map[string]string
	// PageURL is the route path for the billing page (e.g. "/app/billing").
	// Defaults to "/app/billing" when empty for backward compatibility.
	PageURL string
}

// Module wires the billing route.
type Module struct {
	deps *ModuleDeps
}

// NewModule creates a new billing module.
func NewModule(deps *ModuleDeps) *Module {
	return &Module{deps: deps}
}

// RegisterRoutes registers the GET handler for the billing page.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	pageURL := m.deps.PageURL
	if pageURL == "" {
		pageURL = "/app/billing"
	}
	r.GET(pageURL, billingdetail.NewView(&billingdetail.ModuleDeps{
		Messages: m.deps.Messages,
		PageURL:  pageURL,
	}))
}

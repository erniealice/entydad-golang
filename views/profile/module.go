// Package profile provides the /app/profile page — personal info (single
// section, no internal tabs). Part of the four "personal scope" pages
// accessible from the sidebar bottom profile popover.
//
// Permission gating (Layer 3): user:read
package profile

import (
	profiledetail "github.com/erniealice/entydad-golang/views/profile/detail"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds dependencies for the profile module.
type ModuleDeps struct {
	Messages map[string]string
	// PageURL is the route path for the profile page (e.g. "/app/profile").
	// Defaults to "/app/profile" when empty for backward compatibility.
	PageURL string
}

// Module wires the profile route.
type Module struct {
	deps *ModuleDeps
}

// NewModule creates a new profile module.
func NewModule(deps *ModuleDeps) *Module {
	return &Module{deps: deps}
}

// RegisterRoutes registers the GET handler for the profile page.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	pageURL := m.deps.PageURL
	if pageURL == "" {
		pageURL = "/app/profile"
	}
	r.GET(pageURL, profiledetail.NewView(&profiledetail.ModuleDeps{
		Messages: m.deps.Messages,
	}))
}

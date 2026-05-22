// Package account provides the /app/account page — account & security
// (tabs: email | password | sessions). Part of the four "personal scope"
// pages accessible from the sidebar bottom profile popover.
//
// Permission gating (Layer 3): user:update
package account

import (
	accountdetail "github.com/erniealice/entydad-golang/views/account/detail"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds dependencies for the account module.
type ModuleDeps struct {
	Messages          map[string]string
	ChangePasswordURL string
	// PageURL is the route path for the account page (e.g. "/app/account").
	// Defaults to "/app/account" when empty for backward compatibility.
	PageURL string
}

// Module wires the account route.
type Module struct {
	deps *ModuleDeps
}

// NewModule creates a new account module.
func NewModule(deps *ModuleDeps) *Module {
	return &Module{deps: deps}
}

// RegisterRoutes registers the GET handler for the account page.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	pageURL := m.deps.PageURL
	if pageURL == "" {
		pageURL = "/app/account"
	}
	r.GET(pageURL, accountdetail.NewView(&accountdetail.ModuleDeps{
		Messages:          m.deps.Messages,
		ChangePasswordURL: m.deps.ChangePasswordURL,
		PageURL:           pageURL,
	}))
}

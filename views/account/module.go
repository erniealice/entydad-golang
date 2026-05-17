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
}

// Module wires the account route.
type Module struct {
	deps *ModuleDeps
}

// NewModule creates a new account module.
func NewModule(deps *ModuleDeps) *Module {
	return &Module{deps: deps}
}

// RegisterRoutes registers the GET handler for /app/account.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET("/app/account", accountdetail.NewView(&accountdetail.ModuleDeps{
		Messages:          m.deps.Messages,
		ChangePasswordURL: m.deps.ChangePasswordURL,
	}))
}

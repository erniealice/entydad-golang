// Package preference provides the /app/preferences page — UI preferences
// (tabs: appearance | notifications | language-region). Part of the four
// "personal scope" pages accessible from the sidebar bottom profile popover.
//
// Note: package name is singular (preference) per entydad noun convention;
// the route URL stays plural (/app/preferences).
//
// Permission gating (Layer 3): user:update
package preference

import (
	preferencedetail "github.com/erniealice/entydad-golang/service/portal/views/preference/detail"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds dependencies for the preference module.
type ModuleDeps struct {
	Messages map[string]string
	// PageURL is the route path for the preferences page (e.g. "/app/preferences").
	// Defaults to "/app/preferences" when empty for backward compatibility.
	PageURL string
}

// Module wires the preference route.
type Module struct {
	deps *ModuleDeps
}

// NewModule creates a new preference module.
func NewModule(deps *ModuleDeps) *Module {
	return &Module{deps: deps}
}

// RegisterRoutes registers the GET handler for the preferences page.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	pageURL := m.deps.PageURL
	if pageURL == "" {
		pageURL = "/app/preferences"
	}
	r.GET(pageURL, preferencedetail.NewView(&preferencedetail.ModuleDeps{
		Messages: m.deps.Messages,
		PageURL:  pageURL,
	}))
}

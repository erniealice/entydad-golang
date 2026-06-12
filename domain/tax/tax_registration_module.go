// tax_registration_module.go — hoisted view-module assembler for the
// TaxRegistration entity (package tax, the domain root).
//
// Relocated 2026-06-12 from views/tax_registration/module.go (fork E4 / thread
// TX). Pure structural move — the view wiring is byte-identical. The assembler
// is HOISTED to package tax (the domain root, Option-B module-hoist standard:
// domain/<d>/<e>_module.go is package <d>) so the entity leaf package
// (tax_registration) holds only its contract (labels.go/routes.go) + embed +
// leaf view sub-packages, with no import edge back up to the assembler.
//
// Entity-local rename from the former root-resident module:
//
//	ModuleDeps  -> TaxRegistrationModuleDeps
//	Module      -> TaxRegistrationModule
//	NewModule   -> NewTaxRegistrationModule
//
// and entydad.TaxRegistration{Routes,Labels} -> tax_registration.{Routes,Labels}.
package tax

import (
	"context"

	taxregistration "github.com/erniealice/entydad-golang/domain/tax/tax_registration"
	"github.com/erniealice/entydad-golang/domain/tax/tax_registration/action"
	listview "github.com/erniealice/entydad-golang/domain/tax/tax_registration/list"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	taxregistrationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_registration"
	taxregistrationkindpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_registration_kind"
)

// TaxRegistrationModuleDeps holds all dependencies for the tax_registration module.
type TaxRegistrationModuleDeps struct {
	Routes       taxregistration.Routes
	Labels       taxregistration.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Tax registration use cases
	ListTaxRegistrations func(ctx context.Context, req *taxregistrationpb.ListTaxRegistrationsRequest) (*taxregistrationpb.ListTaxRegistrationsResponse, error)

	// CreateTaxRegistration is the espyna use case for creating a new registration.
	// When nil, the create action returns an HTMX error.
	CreateTaxRegistration func(ctx context.Context, req *taxregistrationpb.CreateTaxRegistrationRequest) (*taxregistrationpb.CreateTaxRegistrationResponse, error)

	// FindByPartyTypeTaxRegistrationKind filters kinds by party type.
	// When nil, the Kind dropdown returns empty (TODO fallback).
	FindByPartyTypeTaxRegistrationKind func(ctx context.Context, partyType string) ([]*taxregistrationkindpb.TaxRegistrationKind, error)

	// SupersedeTaxRegistration — stub until espyna exposes via consumer.
	// When nil, the supersede action returns 501.
	SupersedeTaxRegistration func(ctx context.Context, partyType, partyID, supersededID string, newReg *taxregistrationpb.TaxRegistration) error

	// RevokeTaxRegistration — stub until espyna exposes via consumer.
	// When nil, the revoke action returns 501.
	RevokeTaxRegistration func(ctx context.Context, id, effectiveTo, reason string) error
}

// TaxRegistrationModule holds all constructed tax_registration views.
type TaxRegistrationModule struct {
	// ClientList is the tax registrations tab view for a client detail page.
	ClientList view.View
	// WorkspaceList is the tax registrations tab view for the workspace settings page.
	WorkspaceList view.View
	// Create handles GET (draw form) and POST (submit) for adding a new registration.
	Create view.View
	// Supersede handles GET (draw form) and POST (submit) for superseding a registration.
	Supersede view.View
	// Revoke handles POST for revoking a registration.
	Revoke view.View
	routes taxregistration.Routes
}

// NewTaxRegistrationModule creates a tax_registration module with List + CRUD views wired for both
// client and workspace party contexts.
func NewTaxRegistrationModule(deps *TaxRegistrationModuleDeps) *TaxRegistrationModule {
	if deps == nil {
		deps = &TaxRegistrationModuleDeps{}
	}

	clientListDeps := &listview.Deps{
		PartyType:            "client",
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		ListTaxRegistrations: deps.ListTaxRegistrations,
	}

	workspaceListDeps := &listview.Deps{
		PartyType:            "workspace",
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		ListTaxRegistrations: deps.ListTaxRegistrations,
	}

	actionDeps := &action.Deps{
		Routes:                             deps.Routes,
		Labels:                             deps.Labels,
		CommonLabels:                       deps.CommonLabels,
		CreateTaxRegistration:              deps.CreateTaxRegistration,
		FindByPartyTypeTaxRegistrationKind: deps.FindByPartyTypeTaxRegistrationKind,
		SupersedeTaxRegistration:           deps.SupersedeTaxRegistration,
		RevokeTaxRegistration:              deps.RevokeTaxRegistration,
	}

	return &TaxRegistrationModule{
		ClientList:    listview.NewView(clientListDeps),
		WorkspaceList: listview.NewView(workspaceListDeps),
		Create:        action.NewCreateAction(actionDeps),
		Supersede:     action.NewSupersedeAction(actionDeps),
		Revoke:        action.NewRevokeAction(actionDeps),
		routes:        deps.Routes,
	}
}

// RegisterRoutes registers all tax_registration routes with the given route registrar.
func (m *TaxRegistrationModule) RegisterRoutes(r view.RouteRegistrar) {
	// Client-scoped list
	if m.ClientList != nil && m.routes.ClientListURL != "" {
		r.GET(m.routes.ClientListURL, m.ClientList)
	}
	// Workspace-scoped list
	if m.WorkspaceList != nil && m.routes.WorkspaceListURL != "" {
		r.GET(m.routes.WorkspaceListURL, m.WorkspaceList)
	}
	// Client-scoped create
	if m.Create != nil {
		if m.routes.ClientAddURL != "" {
			r.GET(m.routes.ClientAddURL, m.Create)
			r.POST(m.routes.ClientAddURL, m.Create)
		}
		if m.routes.WorkspaceAddURL != "" {
			r.GET(m.routes.WorkspaceAddURL, m.Create)
			r.POST(m.routes.WorkspaceAddURL, m.Create)
		}
	}
	// Client-scoped supersede
	if m.Supersede != nil {
		if m.routes.ClientEditURL != "" {
			r.GET(m.routes.ClientEditURL, m.Supersede)
			r.POST(m.routes.ClientEditURL, m.Supersede)
		}
		if m.routes.WorkspaceEditURL != "" {
			r.GET(m.routes.WorkspaceEditURL, m.Supersede)
			r.POST(m.routes.WorkspaceEditURL, m.Supersede)
		}
	}
	// Revoke (POST only)
	if m.Revoke != nil {
		if m.routes.ClientDeleteURL != "" {
			r.POST(m.routes.ClientDeleteURL, m.Revoke)
		}
		if m.routes.WorkspaceDeleteURL != "" {
			r.POST(m.routes.WorkspaceDeleteURL, m.Revoke)
		}
	}
}

// Package action provides HTTP action handlers for TaxRegistration CRUD views
// (create, supersede, revoke).
// Phase 2 C1/C3 — view side of the polymorphic tax registration form.
package action

import (
	"context"
	"log"
	"net/http"

	taxregistration "github.com/erniealice/entydad-golang/domain/tax/tax_registration"
	"github.com/erniealice/entydad-golang/domain/tax/tax_registration/form"
	taxregistrationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_registration"
	taxregistrationkindpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_registration_kind"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds the use-case callbacks required by the create/supersede/revoke handlers.
type Deps struct {
	Routes       taxregistration.Routes
	Labels       taxregistration.Labels
	CommonLabels any

	// CreateTaxRegistration is the espyna use case for creating a new registration.
	CreateTaxRegistration func(ctx context.Context, req *taxregistrationpb.CreateTaxRegistrationRequest) (*taxregistrationpb.CreateTaxRegistrationResponse, error)

	// FindByPartyTypeTaxRegistrationKind filters kinds whose applicable_party_types
	// includes the given party type. When nil, all kinds are returned (TODO fallback).
	// Phase 2 C3 wires this from espyna's FindByPartyType use case.
	FindByPartyTypeTaxRegistrationKind func(ctx context.Context, partyType string) ([]*taxregistrationkindpb.TaxRegistrationKind, error)

	// SupersedeTaxRegistration — TODO: stub until other agent exposes via consumer.
	// When nil, the supersede action returns 501.
	SupersedeTaxRegistration func(ctx context.Context, partyType, partyID, supersededID string, newReg *taxregistrationpb.TaxRegistration) error

	// RevokeTaxRegistration — TODO: stub until other agent exposes via consumer.
	// When nil, the revoke action returns 501.
	RevokeTaxRegistration func(ctx context.Context, id, effectiveTo, reason string) error
}

// loadKindOptions fetches TaxRegistrationKind options filtered by party type.
// Falls back to empty list + TODO log when FindByPartyTypeTaxRegistrationKind is not wired.
func loadKindOptions(ctx context.Context, deps *Deps, partyType, selectedID string) []form.KindOption {
	if deps.FindByPartyTypeTaxRegistrationKind == nil {
		// TODO: wire FindByPartyType use case from espyna consumer
		log.Printf("[TODO] FindByPartyTypeTaxRegistrationKind not wired; returning empty kind list for party_type=%s", partyType)
		return nil
	}
	kinds, err := deps.FindByPartyTypeTaxRegistrationKind(ctx, partyType)
	if err != nil {
		log.Printf("FindByPartyTypeTaxRegistrationKind(%s): %v", partyType, err)
		return nil
	}
	opts := make([]form.KindOption, 0, len(kinds))
	for _, k := range kinds {
		opts = append(opts, form.KindOption{
			Value:    k.GetId(),
			Label:    k.GetName(),
			Selected: k.GetId() == selectedID,
		})
	}
	return opts
}

// NewCreateAction handles GET (form) and POST (submit) for adding a new
// TaxRegistration attached to the specified party.
//
// GET query params:
//   - party_type — "client" or "workspace"
//   - party_id   — the party's record ID
func NewCreateAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("tax_registration", "create") {
			return view.HTMXError("You do not have permission to create tax registrations")
		}

		labels := form.BuildLabels(deps.Labels)
		q := viewCtx.Request.URL.Query()
		partyType := q.Get("party_type")
		partyID := q.Get("party_id")
		if partyID == "" {
			partyID = viewCtx.Request.PathValue("id")
		}

		if viewCtx.Request.Method == http.MethodGet {
			kindOpts := loadKindOptions(ctx, deps, partyType, "")
			return view.OK("tax-registration-drawer-form", &form.Data{
				FormAction:   deps.Routes.AddURL,
				IsEdit:       false,
				PartyType:    partyType,
				PartyID:      partyID,
				KindOptions:  kindOpts,
				Labels:       labels,
				CommonLabels: deps.CommonLabels,
			})
		}

		// POST — create
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError("Invalid form data")
		}
		r := viewCtx.Request
		kindID := r.FormValue("tax_registration_kind_id")
		if kindID == "" {
			return view.HTMXError("Tax registration kind is required")
		}
		partyType = r.FormValue("party_type")
		partyID = r.FormValue("party_id")

		var partyTypeEnum taxregistrationpb.TaxRegistrationPartyType
		switch partyType {
		case "client":
			partyTypeEnum = taxregistrationpb.TaxRegistrationPartyType_TAX_REGISTRATION_PARTY_TYPE_CLIENT
		case "workspace":
			partyTypeEnum = taxregistrationpb.TaxRegistrationPartyType_TAX_REGISTRATION_PARTY_TYPE_WORKSPACE
		}
		record := &taxregistrationpb.TaxRegistration{
			PartyType:             partyTypeEnum,
			PartyId:               partyID,
			TaxRegistrationKindId: kindID,
			RegistrationNumber:    r.FormValue("registration_number"),
			EffectiveFrom:         r.FormValue("effective_from"),
		}
		if deps.CreateTaxRegistration == nil {
			return view.HTMXError("Create use case not available")
		}
		if _, err := deps.CreateTaxRegistration(ctx, &taxregistrationpb.CreateTaxRegistrationRequest{
			Data: record,
		}); err != nil {
			log.Printf("CreateTaxRegistration error: %v", err)
			return view.HTMXError(err.Error())
		}

		// Reload the list view that embedded us
		return view.HTMXSuccess("tax-registrations-table")
	})
}

package action

import (
	"context"
	"net/http"

	"github.com/erniealice/entydad-golang/domain/tax/tax_registration/form"
	"github.com/erniealice/pyeza-golang/view"
)

// NewSupersedeAction handles GET (form) and POST (submit) for superseding
// an existing TaxRegistration.
//
// GET query params:
//   - party_type — "client" or "workspace"
//   - party_id   — the party's record ID
//
// When deps.SupersedeTaxRegistration is nil (other agent's espyna work not yet
// ready), this handler returns 501 on POST and a TODO form on GET.
func NewSupersedeAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("tax_registration", "update") {
			return view.HTMXError("You do not have permission to supersede tax registrations")
		}

		labels := form.BuildLabels(deps.Labels)
		q := viewCtx.Request.URL.Query()
		partyType := q.Get("party_type")
		partyID := q.Get("party_id")
		supersededID := viewCtx.Request.PathValue("reg_id")
		if supersededID == "" {
			supersededID = q.Get("reg_id")
		}

		if viewCtx.Request.Method == http.MethodGet {
			// TODO: load existing registration and pre-fill kind + registration number
			kindOpts := loadKindOptions(ctx, deps, partyType, "")
			return view.OK("tax-registration-drawer-form", &form.Data{
				FormAction:   deps.Routes.ClientEditURL,
				IsEdit:       true,
				ID:           supersededID,
				PartyType:    partyType,
				PartyID:      partyID,
				KindOptions:  kindOpts,
				Labels:       labels,
				CommonLabels: deps.CommonLabels,
			})
		}

		// POST — supersede
		if deps.SupersedeTaxRegistration == nil {
			// TODO: wire SupersedeTaxRegistration from espyna consumer once espyna agent delivers
			return view.ViewResult{StatusCode: http.StatusNotImplemented}
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError("Invalid form data")
		}
		r := viewCtx.Request
		partyType = r.FormValue("party_type")
		partyID = r.FormValue("party_id")
		effectiveFrom := r.FormValue("effective_from")
		kindID := r.FormValue("tax_registration_kind_id")
		regNum := r.FormValue("registration_number")

		if effectiveFrom == "" || kindID == "" {
			return view.HTMXError("Kind and effective date are required")
		}

		// TODO: convert to proto once SupersedeTaxRegistration request type is defined
		_ = effectiveFrom
		_ = regNum
		if err := deps.SupersedeTaxRegistration(ctx, partyType, partyID, supersededID, nil); err != nil {
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("tax-registrations-table")
	})
}

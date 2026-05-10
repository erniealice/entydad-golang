package action

import (
	"context"
	"net/http"

	entydad "github.com/erniealice/entydad-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// NewRevokeAction handles POST for revoking a TaxRegistration.
// The revoke confirm form is rendered by the tax-registration-revoke-confirm
// template; this handler processes the confirmed POST.
//
// When deps.RevokeTaxRegistration is nil (other agent's espyna work not yet
// ready), this handler returns 501.
func NewRevokeAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("tax_registration", "delete") {
			return entydad.HTMXError("You do not have permission to revoke tax registrations")
		}

		// POST only — the GET is the revoke confirm form (pre-rendered by the list view action)
		if viewCtx.Request.Method != http.MethodPost {
			return entydad.HTMXError("Method not allowed")
		}

		if deps.RevokeTaxRegistration == nil {
			// TODO: wire RevokeTaxRegistration from espyna consumer once espyna agent delivers
			return view.ViewResult{StatusCode: http.StatusNotImplemented}
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}
		r := viewCtx.Request
		id := r.FormValue("id")
		effectiveTo := r.FormValue("effective_to")
		reason := r.FormValue("reason")

		if id == "" {
			return entydad.HTMXError("Registration ID is required")
		}

		if err := deps.RevokeTaxRegistration(ctx, id, effectiveTo, reason); err != nil {
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("tax-registrations-table")
	})
}

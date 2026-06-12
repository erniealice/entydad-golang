package tax_registration

// routes.go — TaxRegistration route structs, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (the root tax leftovers). Pure structural move — route URL string values are
// byte-identical. Entity-local rename:
//   TaxRegistrationRoutes        -> Routes
//   DefaultTaxRegistrationRoutes -> DefaultRoutes
// The party-scoped URL constants keep their Client/Workspace prefixes (they are
// genuinely local to this polymorphic entity). The facade (domain/tax/tax.go)
// restores the original entydad.TaxRegistration* names for consumers.

// TaxRegistration — polymorphic (client + workspace party types in v1)
// URL convention: party_type + party_id come from the parent detail page context.
const (
	ClientTaxRegistrationListURL   = "/clients/detail/{id}/tax-registrations"
	ClientTaxRegistrationAddURL    = "/action/client/{id}/tax-registration/add"
	ClientTaxRegistrationEditURL   = "/action/client/{id}/tax-registration/edit/{reg_id}"
	ClientTaxRegistrationDeleteURL = "/action/client/{id}/tax-registration/delete"

	WorkspaceTaxRegistrationListURL   = "/workspace/settings/tax-registrations"
	WorkspaceTaxRegistrationAddURL    = "/action/workspace/tax-registration/add"
	WorkspaceTaxRegistrationEditURL   = "/action/workspace/tax-registration/edit/{reg_id}"
	WorkspaceTaxRegistrationDeleteURL = "/action/workspace/tax-registration/delete"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for the polymorphic TaxRegistration
// views. v1 surfaces client + workspace party types only.
// The AddURL and DeleteURL are party-scoped (include party_id in path).
type Routes struct {
	// Client-scoped routes
	ClientListURL   string `json:"client_list_url"`
	ClientAddURL    string `json:"client_add_url"`
	ClientEditURL   string `json:"client_edit_url"`
	ClientDeleteURL string `json:"client_delete_url"`

	// Workspace-scoped routes
	WorkspaceListURL   string `json:"workspace_list_url"`
	WorkspaceAddURL    string `json:"workspace_add_url"`
	WorkspaceEditURL   string `json:"workspace_edit_url"`
	WorkspaceDeleteURL string `json:"workspace_delete_url"`

	// AddURL and DeleteURL are the active-context URLs (set by the view wiring).
	// For client views: same as ClientAddURL / ClientDeleteURL.
	// For workspace views: same as WorkspaceAddURL / WorkspaceDeleteURL.
	AddURL    string `json:"add_url"`
	DeleteURL string `json:"delete_url"`
}

// DefaultRoutes returns a Routes populated from
// the package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ClientListURL:   ClientTaxRegistrationListURL,
		ClientAddURL:    ClientTaxRegistrationAddURL,
		ClientEditURL:   ClientTaxRegistrationEditURL,
		ClientDeleteURL: ClientTaxRegistrationDeleteURL,

		WorkspaceListURL:   WorkspaceTaxRegistrationListURL,
		WorkspaceAddURL:    WorkspaceTaxRegistrationAddURL,
		WorkspaceEditURL:   WorkspaceTaxRegistrationEditURL,
		WorkspaceDeleteURL: WorkspaceTaxRegistrationDeleteURL,

		// Default to client context; override at wiring time for workspace.
		AddURL:    ClientTaxRegistrationAddURL,
		DeleteURL: ClientTaxRegistrationDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"tax_registration.client.list":      r.ClientListURL,
		"tax_registration.client.add":       r.ClientAddURL,
		"tax_registration.client.edit":      r.ClientEditURL,
		"tax_registration.client.delete":    r.ClientDeleteURL,
		"tax_registration.workspace.list":   r.WorkspaceListURL,
		"tax_registration.workspace.add":    r.WorkspaceAddURL,
		"tax_registration.workspace.edit":   r.WorkspaceEditURL,
		"tax_registration.workspace.delete": r.WorkspaceDeleteURL,
	}
}

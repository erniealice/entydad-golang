package client

// routes.go — Client route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, party sub-context). Pure structural move — route URL string
// values are byte-identical. Entity-local rename: ClientRoutes -> Routes,
// DefaultClientRoutes -> DefaultRoutes, Client<Xxx>URL -> <Xxx>URL.

// Default route constants for the client view.
const (
	DashboardURL        = "/clients/dashboard"
	ListURL             = "/clients/list/{status}"
	TableURL            = "/action/client/table/{status}"
	AddURL              = "/action/client/add"
	EditURL             = "/action/client/edit/{id}"
	DeleteURL           = "/action/client/delete"
	BulkDeleteURL       = "/action/client/bulk-delete"
	DetailURL           = "/clients/detail/{id}"
	TabActionURL        = "/action/client/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/client/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/client/{id}/attachments/delete"
	SetStatusURL        = "/action/client/set-status"
	BulkSetStatusURL    = "/action/client/bulk-set-status"
	SearchURL           = "/action/client/search"

	StatementExportURL = "/action/client/{id}/statement/export"

	// RevenueRunURL is the per-client "Run Invoices" drawer endpoint.
	RevenueRunURL = "/action/client/revenue-run/{id}"
)

type Routes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	DetailURL        string `json:"detail_url"`
	TabActionURL     string `json:"tab_action_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	SearchURL        string `json:"search_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Statement export
	StatementExportURL string `json:"statement_export_url"`

	// Report routes
	ReceivablesAgingURL string `json:"receivables_aging_url"`

	// Settings routes
	PaymentTermsURL  string `json:"payment_terms_url"`
	ClientTagListURL string `json:"client_tag_list_url"` // cross-app link to client-tag list (dashboard quick-action)

	// RevenueRunURL is the per-client "Run Invoices" drawer endpoint.
	RevenueRunURL string `json:"revenue_run_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants.
func DefaultRoutes() Routes {
	return Routes{
		DashboardURL:     DashboardURL,
		ListURL:          ListURL,
		TableURL:         TableURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		DetailURL:        DetailURL,
		TabActionURL:     TabActionURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
		SearchURL:        SearchURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		StatementExportURL: StatementExportURL,

		// FINALIZE: cross-entity link (reports surface ReceivablesAgingURL);
		// value byte-identical to root entydad.ReceivablesAgingURL.
		ReceivablesAgingURL: "/reports/receivables-aging",

		// FINALIZE: cross-entity link (commerce/payment_term, client context
		// PaymentTermListURL); value byte-identical to root.
		PaymentTermsURL: "/clients/settings/payment-terms/list",
		// FINALIZE: cross-entity link (party/client_tag ListURL); value
		// byte-identical to root entydad.ClientTagListURL.
		ClientTagListURL: "/clients/settings/tags/list",
		RevenueRunURL:    RevenueRunURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"client.dashboard":       r.DashboardURL,
		"client.list":            r.ListURL,
		"client.table":           r.TableURL,
		"client.add":             r.AddURL,
		"client.edit":            r.EditURL,
		"client.delete":          r.DeleteURL,
		"client.bulk_delete":     r.BulkDeleteURL,
		"client.detail":          r.DetailURL,
		"client.tab_action":      r.TabActionURL,
		"client.set_status":      r.SetStatusURL,
		"client.bulk_set_status": r.BulkSetStatusURL,
		"client.search":          r.SearchURL,

		"client.attachment.upload": r.AttachmentUploadURL,
		"client.attachment.delete": r.AttachmentDeleteURL,

		"client.statement_export": r.StatementExportURL,

		"client.receivables_aging": r.ReceivablesAgingURL,

		"client.payment_terms": r.PaymentTermsURL,
		"client.revenue_run":   r.RevenueRunURL,
	}
}

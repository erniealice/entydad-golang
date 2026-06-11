package supplier

// routes.go — Supplier route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, party sub-context). Pure structural move — route URL string
// values are byte-identical. Entity-local rename: SupplierRoutes -> Routes,
// DefaultSupplierRoutes -> DefaultRoutes, Supplier<Xxx>URL -> <Xxx>URL.

// Default route constants for the supplier view.
const (
	DashboardURL        = "/suppliers/dashboard"
	ListURL             = "/suppliers/list/{status}"
	TableURL            = "/action/supplier/table/{status}"
	AddURL              = "/action/supplier/add"
	EditURL             = "/action/supplier/edit/{id}"
	DeleteURL           = "/action/supplier/delete"
	BulkDeleteURL       = "/action/supplier/bulk-delete"
	DetailURL           = "/suppliers/detail/{id}"
	TabActionURL        = "/action/supplier/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/supplier/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/supplier/{id}/attachments/delete"
	SetStatusURL        = "/action/supplier/set-status"
	BulkSetStatusURL    = "/action/supplier/bulk-set-status"

	StatementExportURL = "/action/supplier/{id}/statement/export"

	// Plan A 20260517-expense-run — Surface A per-supplier drawer URL.
	ExpenseRecognitionRunURL = "/action/supplier/expense-recognition-run/{id}"
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

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Statement export
	StatementExportURL string `json:"statement_export_url"`

	// Report routes
	PayablesAgingURL string `json:"payables_aging_url"`

	// Settings routes
	PaymentTermsURL string `json:"payment_terms_url"`

	// Plan A 20260517-expense-run — Surface A per-supplier drawer URL.
	// "Run Recognitions" CTA on the Statement tab opens this drawer.
	ExpenseRecognitionRunURL string `json:"expense_recognition_run_url"`
}

// DefaultRoutes returns a Routes populated from the
// package-level route constants.
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

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		StatementExportURL: StatementExportURL,

		// FINALIZE: cross-entity link (reports surface PayablesAgingURL);
		// value byte-identical to root entydad.PayablesAgingURL.
		PayablesAgingURL: "/suppliers/reports/payables-aging",

		// FINALIZE: cross-entity link (commerce/payment_term, supplier context
		// SupplierPaymentTermListURL); value byte-identical to root.
		PaymentTermsURL: "/suppliers/settings/payment-terms/list",

		ExpenseRecognitionRunURL: ExpenseRecognitionRunURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"supplier.dashboard":       r.DashboardURL,
		"supplier.list":            r.ListURL,
		"supplier.table":           r.TableURL,
		"supplier.add":             r.AddURL,
		"supplier.edit":            r.EditURL,
		"supplier.delete":          r.DeleteURL,
		"supplier.bulk_delete":     r.BulkDeleteURL,
		"supplier.detail":          r.DetailURL,
		"supplier.tab_action":      r.TabActionURL,
		"supplier.set_status":      r.SetStatusURL,
		"supplier.bulk_set_status": r.BulkSetStatusURL,

		"supplier.attachment.upload": r.AttachmentUploadURL,
		"supplier.attachment.delete": r.AttachmentDeleteURL,

		"supplier.statement_export": r.StatementExportURL,

		"supplier.payables_aging": r.PayablesAgingURL,

		"supplier.payment_terms": r.PaymentTermsURL,

		"supplier.expense_recognition_run": r.ExpenseRecognitionRunURL,
	}
}

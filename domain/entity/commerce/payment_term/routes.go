package payment_term

// routes.go — PaymentTerm route structs, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, commerce sub-context). Pure structural move — route URL
// string values are byte-identical.
//
// Entity-local rename:
//   PaymentTermRoutes          -> Routes
//   DefaultPaymentTermRoutes   -> DefaultRoutes
//   PaymentTerm<Xxx>URL        -> <Xxx>URL
//
// The supplier-context variant is genuinely local to this entity (a
// supplier-scoped view of payment-term routes — NOT shared with the supplier
// entity). Its names are not <E>-prefixed at the anchor, so the redundant
// PaymentTerm token is dropped consistently across the whole supplier family:
//   SupplierPaymentTermRoutes          -> SupplierRoutes
//   DefaultSupplierPaymentTermRoutes   -> DefaultSupplierRoutes
//   SupplierPaymentTerm<Xxx>URL        -> Supplier<Xxx>URL
//   (r SupplierPaymentTermRoutes) ToPaymentTermRoutes() -> (r SupplierRoutes) ToRoutes()
// Finalize's facade restores the original entydad.SupplierPaymentTerm* names.

// Payment Term routes — client context (shows client + both scopes)
const (
	ListURL          = "/clients/settings/payment-terms/list"
	TableURL         = "/action/client/settings/payment-terms/table"
	AddURL           = "/action/client/settings/payment-terms/add"
	EditURL          = "/action/client/settings/payment-terms/edit/{id}"
	DeleteURL        = "/action/client/settings/payment-terms/delete"
	BulkDeleteURL    = "/action/client/settings/payment-terms/bulk-delete"
	SetStatusURL     = "/action/client/settings/payment-terms/set-status"
	BulkSetStatusURL = "/action/client/settings/payment-terms/bulk-set-status"
)

// Payment Term routes — supplier context (shows supplier + both scopes)
const (
	SupplierListURL          = "/suppliers/settings/payment-terms/list"
	SupplierTableURL         = "/action/supplier/settings/payment-terms/table"
	SupplierAddURL           = "/action/supplier/settings/payment-terms/add"
	SupplierEditURL          = "/action/supplier/settings/payment-terms/edit/{id}"
	SupplierDeleteURL        = "/action/supplier/settings/payment-terms/delete"
	SupplierBulkDeleteURL    = "/action/supplier/settings/payment-terms/bulk-delete"
	SupplierSetStatusURL     = "/action/supplier/settings/payment-terms/set-status"
	SupplierBulkSetStatusURL = "/action/supplier/settings/payment-terms/bulk-set-status"
)

// Routes holds all route paths for payment term management.
type Routes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route
// constants.
func DefaultRoutes() Routes {
	return Routes{
		ListURL:          ListURL,
		TableURL:         TableURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"payment_term.list":            r.ListURL,
		"payment_term.table":           r.TableURL,
		"payment_term.add":             r.AddURL,
		"payment_term.edit":            r.EditURL,
		"payment_term.delete":          r.DeleteURL,
		"payment_term.bulk_delete":     r.BulkDeleteURL,
		"payment_term.set_status":      r.SetStatusURL,
		"payment_term.bulk_set_status": r.BulkSetStatusURL,
	}
}

// SupplierRoutes holds payment term route paths for the supplier context.
// These routes show only payment terms with entity_scope IN ('supplier', 'both').
type SupplierRoutes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultSupplierRoutes returns a SupplierRoutes from package-level constants.
func DefaultSupplierRoutes() SupplierRoutes {
	return SupplierRoutes{
		ListURL:          SupplierListURL,
		TableURL:         SupplierTableURL,
		AddURL:           SupplierAddURL,
		EditURL:          SupplierEditURL,
		DeleteURL:        SupplierDeleteURL,
		BulkDeleteURL:    SupplierBulkDeleteURL,
		SetStatusURL:     SupplierSetStatusURL,
		BulkSetStatusURL: SupplierBulkSetStatusURL,
	}
}

// ToRoutes converts SupplierRoutes to a Routes, allowing the payment term
// module to be reused with supplier-context paths.
func (r SupplierRoutes) ToRoutes() Routes {
	return Routes{
		ListURL:          r.ListURL,
		TableURL:         r.TableURL,
		AddURL:           r.AddURL,
		EditURL:          r.EditURL,
		DeleteURL:        r.DeleteURL,
		BulkDeleteURL:    r.BulkDeleteURL,
		SetStatusURL:     r.SetStatusURL,
		BulkSetStatusURL: r.BulkSetStatusURL,
	}
}

package entydad

// PortalLabels holds translations for the self-service portal surfaces.
// Loaded from general/portal.json (with optional per-vertical overrides via the
// standard lyngua cascade: common → general → businessType).
//
// JSON field names mirror the key hierarchy in portal.json so the lyngua
// LoadPath call ("portal.json", "portal", &labels) populates every field
// without manual mapping.
type PortalLabels struct {
	Client           PortalClientLabels           `json:"client"`
	Supplier         PortalSupplierLabels         `json:"supplier"`
	ClientDelegate   PortalClientDelegateLabels   `json:"client_delegate"`
	SupplierDelegate PortalSupplierDelegateLabels `json:"supplier_delegate"`
	Home             PortalHomeSharedLabels       `json:"home"`
	Page             PortalPageSectionLabels      `json:"page"`
	Header           PortalHeaderLabels           `json:"header"`
	Accessibility    PortalAccessibilityLabels    `json:"accessibility"`
	Sidebar          PortalSidebarLabels          `json:"sidebar"`
}

// PortalClientLabels holds client-specific portal labels.
type PortalClientLabels struct {
	Name        string             `json:"name"`
	Home        PortalHomeLabels   `json:"home"`
	Profile     PortalTitleOnly    `json:"profile"`
	Account     PortalTitleOnly    `json:"account"`
	Billing     PortalTitleOnly    `json:"billing"`
	Preferences PortalTitleOnly    `json:"preferences"`
}

// PortalSupplierLabels holds supplier-specific portal labels.
type PortalSupplierLabels struct {
	Name        string           `json:"name"`
	Home        PortalHomeLabels `json:"home"`
	Profile     PortalTitleOnly  `json:"profile"`
	Account     PortalTitleOnly  `json:"account"`
	Billing     PortalTitleOnly  `json:"billing"`
	Preferences PortalTitleOnly  `json:"preferences"`
}

// PortalClientDelegateLabels holds client-delegate-specific portal labels.
type PortalClientDelegateLabels struct {
	Name        string                   `json:"name"`
	Home        PortalHomeLabels         `json:"home"`
	Profile     PortalTitleOnly          `json:"profile"`
	Account     PortalTitleOnly          `json:"account"`
	Billing     PortalTitleOnly          `json:"billing"`
	Preferences PortalTitleOnly          `json:"preferences"`
	Select      PortalSelectLabels       `json:"select"`
}

// PortalSupplierDelegateLabels holds supplier-delegate-specific portal labels.
type PortalSupplierDelegateLabels struct {
	Name        string             `json:"name"`
	Home        PortalHomeLabels   `json:"home"`
	Profile     PortalTitleOnly    `json:"profile"`
	Account     PortalTitleOnly    `json:"account"`
	Billing     PortalTitleOnly    `json:"billing"`
	Preferences PortalTitleOnly    `json:"preferences"`
	Select      PortalSelectLabels `json:"select"`
}

// PortalHomeLabels holds labels for a portal home page (per principal kind).
type PortalHomeLabels struct {
	PageTitle  string `json:"pageTitle"`
	Heading    string `json:"heading"`
	Subheading string `json:"subheading"`
}

// PortalTitleOnly holds only a pageTitle — for stub pages that have no other
// per-kind label variation.
type PortalTitleOnly struct {
	PageTitle string `json:"pageTitle"`
}

// PortalSelectLabels holds labels for the acting-as delegate picker page.
type PortalSelectLabels struct {
	PageTitle  string `json:"pageTitle"`
	Heading    string `json:"heading"`
	Subheading string `json:"subheading"`
	EmptyState string `json:"emptyState"`
}

// PortalHomeSharedLabels holds labels shared across all portal home pages
// (recent-activity section).
type PortalHomeSharedLabels struct {
	RecentActivityTitle string `json:"recentActivityTitle"`
	NoRecentActivity    string `json:"noRecentActivity"`
}

// PortalPageSectionLabels holds labels for stub content pages.
type PortalPageSectionLabels struct {
	Profile     PortalPageLabels         `json:"profile"`
	Account     PortalPageLabels         `json:"account"`
	Billing     PortalBillingPageLabels  `json:"billing"`
	Preferences PortalPageLabels         `json:"preferences"`
}

// PortalPageLabels holds title + comingSoon for a stub page.
type PortalPageLabels struct {
	Title     string `json:"title"`
	ComingSoon string `json:"comingSoon"`
}

// PortalBillingPageLabels extends PortalPageLabels with delegate variants.
type PortalBillingPageLabels struct {
	Title                      string `json:"title"`
	ComingSoon                 string `json:"comingSoon"`
	ComingSoonDelegate         string `json:"comingSoonDelegate"`
	ComingSoonSupplierDelegate string `json:"comingSoonSupplierDelegate"`
}

// PortalHeaderLabels holds labels for the portal header.
type PortalHeaderLabels struct {
	NavAriaLabel string `json:"navAriaLabel"`
}

// PortalAccessibilityLabels holds accessibility-specific portal labels.
type PortalAccessibilityLabels struct {
	SkipToMainContent string `json:"skipToMainContent"`
}

// PortalSidebarLabels holds labels for the portal sidebar.
type PortalSidebarLabels struct {
	Nav     PortalSidebarNavLabels     `json:"nav"`
	Section PortalSidebarSectionLabels `json:"section"`
}

// PortalSidebarNavLabels holds labels for sidebar navigation items.
type PortalSidebarNavLabels struct {
	Home              string `json:"home"`
	Profile           string `json:"profile"`
	Account           string `json:"account"`
	Billing           string `json:"billing"`
	Preferences       string `json:"preferences"`
	Invoices          string `json:"invoices"`
	Messages          string `json:"messages"`
	Documents         string `json:"documents"`
	PurchaseOrders    string `json:"purchaseOrders"`
	SubmittedInvoices string `json:"submittedInvoices"`
	Contracts         string `json:"contracts"`
}

// PortalSidebarSectionLabels holds labels for sidebar section headings.
type PortalSidebarSectionLabels struct {
	Overview      string `json:"overview"`
	Summary       string `json:"summary"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentHistory string `json:"paymentHistory"`
	Outstanding   string `json:"outstanding"`
	Paid          string `json:"paid"`
	Draft         string `json:"draft"`
	Pending       string `json:"pending"`
	Active        string `json:"active"`
	Expiring      string `json:"expiring"`
	Terminated    string `json:"terminated"`
}

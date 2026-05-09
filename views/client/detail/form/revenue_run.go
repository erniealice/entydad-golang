// Package form contains template-facing data types for the client detail
// feature drawers. Per the drawer-form-subpackage-convention, each secondary
// feature drawer under views/client/detail/ contributes its Data + Labels
// types here. No repository imports; no Deps structs.
package form

import (
	entydad "github.com/erniealice/entydad-golang"
	pyeza "github.com/erniealice/pyeza-golang"
)

// RevenueRunDrawerData is the top-level template context for the
// client-revenue-run-drawer-form template.
type RevenueRunDrawerData struct {
	// FormAction is the POST URL (same URL as the GET — single endpoint).
	FormAction string
	// FragmentURL is the GET URL used by the HTMX inner-swap partial when the
	// as_of_date changes. Typically equals FormAction + query params.
	FragmentURL string
	// ClientID is the client being operated on.
	ClientID string
	// ClientName is the display name for the client.
	ClientName string
	// AsOfDate is the current as-of date value (YYYY-MM-DD).
	AsOfDate string
	// MaxAsOfDate caps the date picker to today (YYYY-MM-DD).
	MaxAsOfDate string
	// Subtitle is the pre-computed subtitle text built from SubtitleTemplate
	// with EligibleCount and SubscriptionCount substituted in. The Go handler
	// fills this so the template can emit it verbatim.
	Subtitle string
	// EligibleCount is the number of periods eligible for invoicing.
	EligibleCount int
	// SubscriptionCount is the number of distinct subscriptions with candidates.
	SubscriptionCount int
	// SubscriptionGroups holds per-subscription candidate groups.
	SubscriptionGroups []RevenueRunGroup
	// TotalsByCurrency holds per-currency totals across all eligible periods,
	// keyed by currency code.
	TotalsByCurrency map[string]int64
	// Labels carries all user-facing strings for this drawer.
	Labels entydad.ClientRevenueRunLabels
	// CommonLabels carries shared UI strings (Save / Cancel / etc.).
	CommonLabels pyeza.CommonLabels
}

// RevenueRunGroup is one subscription's candidate group rendered as a
// collapsible section in the drawer.
type RevenueRunGroup struct {
	// SubscriptionID is the proto ID.
	SubscriptionID string
	// SubscriptionName is the display name for the subscription.
	SubscriptionName string
	// PlanName is the price-plan / plan name.
	PlanName string
	// BillingCycleLabel is the formatted billing-cycle string.
	BillingCycleLabel string
	// Currency is the ISO currency code for this group.
	Currency string
	// CurrencyMismatch is true when the subscription currency differs from the
	// client's billing currency; triggers the mismatch chip.
	CurrencyMismatch bool
	// GroupTotal is the sum of all eligible periods' amounts in centavos.
	GroupTotal int64
	// GroupTotalDisplay is the pre-formatted display string (÷100).
	GroupTotalDisplay string
	// Periods is the list of candidate periods.
	Periods []RevenueRunPeriod
	// HasEligible is true if at least one period in this group is eligible.
	HasEligible bool
}

// RevenueRunPeriod is one candidate billing period row inside a group.
type RevenueRunPeriod struct {
	// SubscriptionID is repeated here for the checkbox value encoding.
	SubscriptionID string
	// PeriodStart is YYYY-MM-DD.
	PeriodStart string
	// PeriodEnd is YYYY-MM-DD.
	PeriodEnd string
	// PeriodMarker is the canonical idempotency anchor.
	PeriodMarker string
	// PeriodLabel is the human-readable range (e.g. "Jan 1 – Jan 31").
	PeriodLabel string
	// Amount is the period amount in centavos.
	Amount int64
	// AmountDisplay is the pre-formatted display string (÷100).
	AmountDisplay string
	// LineItemCount is the number of line items for this period.
	LineItemCount int
	// Eligible indicates the period can be invoiced.
	Eligible bool
	// BlockerReason is the human-readable explanation when Eligible=false.
	BlockerReason string
	// SelectionValue is the composite checkbox value encoding:
	// "{SubscriptionID}|{PeriodStart}|{PeriodEnd}|{PeriodMarker}"
	SelectionValue string
}

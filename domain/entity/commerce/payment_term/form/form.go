package form

import (
	"sort"

	pyeza "github.com/erniealice/pyeza-golang"
)

// Labels holds i18n labels for the payment term drawer form template.
type Labels struct {
	SectionInfo            string
	SectionTerms           string
	SectionSettings        string
	Name                   string
	NamePlaceholder        string
	Code                   string
	CodePlaceholder        string
	Type                   string
	NetDays                string
	DiscountDays           string
	DiscountPercentBps     string
	TypeHint               string
	NetDaysHint            string
	DiscountDaysHint       string
	DiscountPercentBpsHint string
	PriorityHint           string
	EntityScope            string
	IsDefault              string
	Description            string
	DescriptionPlaceholder string
	DisplayOrder           string
	Active                 string

	// Select option labels — Type
	TypeDueOnReceipt        string
	TypeNet                 string
	TypeCOD                 string
	TypeProximate           string
	ProximateDay            string
	ProximateDayPlaceholder string

	// Select option labels — EntityScope
	ScopesBoth         string
	ScopesSupplierOnly string
	ScopesClientOnly   string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo               string
	CodeInfo               string
	CodeHint               string
	DescriptionInfo        string
	TypeInfo               string
	ProximateDayHint       string
	ProximateDayInfo       string
	DiscountDaysInfo       string
	DiscountPercentBpsInfo string
	IsDefaultInfo          string

	// Error messages for server-side validation.
	ErrTypeRequired         string
	ErrTypeInvalid          string
	ErrNetDaysRequired      string
	ErrProximateDayRequired string
}

// Data is the template data for the payment term drawer form.
type Data struct {
	FormAction         string
	WorkspaceID        string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce              string // injected by C1: populated by ViewAdapter.injectPageData via reflection
	IsEdit             bool
	ID                 string
	Name               string
	Code               string
	Type               string
	NetDays            string
	DiscountDays       string
	DiscountPercentBps string
	EntityScope        string
	IsDefault          bool
	Description        string
	DisplayOrder       string
	ProximateDay       string
	Active             bool
	// TypeOptions holds the type select options with the current value pre-selected.
	TypeOptions  []pyeza.SelectOption
	Labels       Labels
	CommonLabels any
}

// BuildTypeOptions constructs the select options for the type field,
// marking the currently-selected value. Labels are drawn from the form labels
// so they flow through lyngua and are not hardcoded in the template.
// Options are sorted alphabetically by their (tier-translated) label so the
// dropdown order remains stable for an operator regardless of which underlying
// proto enum value backs each row.
func BuildTypeOptions(labels Labels, current string) []pyeza.SelectOption {
	// Canonical type values from payment_term.proto field 9:
	//   "net", "due_on_receipt", "cod", "proximate"
	// Default to "net" when current is empty so the initial Add form
	// starts with the most common type pre-selected.
	if current == "" {
		current = "net"
	}
	opts := []pyeza.SelectOption{
		{Value: "net", Label: labels.TypeNet, Selected: current == "net"},
		{Value: "due_on_receipt", Label: labels.TypeDueOnReceipt, Selected: current == "due_on_receipt"},
		{Value: "cod", Label: labels.TypeCOD, Selected: current == "cod"},
		{Value: "proximate", Label: labels.TypeProximate, Selected: current == "proximate"},
	}
	sort.SliceStable(opts, func(i, j int) bool {
		return opts[i].Label < opts[j].Label
	})
	return opts
}

// BuildLabels constructs a Labels struct from a translation function.
// t is typically viewCtx.T — a narrow func(string) string with no Deps or storage access.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		SectionInfo:            t("payment_term.form.section_info"),
		SectionTerms:           t("payment_term.form.section_terms"),
		SectionSettings:        t("payment_term.form.section_settings"),
		Name:                   t("payment_term.form.name"),
		NamePlaceholder:        t("payment_term.form.name_placeholder"),
		Code:                   t("payment_term.form.code"),
		CodePlaceholder:        t("payment_term.form.code_placeholder"),
		Type:                   t("payment_term.form.type"),
		NetDays:                t("payment_term.form.net_days"),
		DiscountDays:           t("payment_term.form.discount_days"),
		DiscountPercentBps:     t("payment_term.form.discount_percent_bps"),
		TypeHint:               t("payment_term.form.type_hint"),
		NetDaysHint:            t("payment_term.form.net_days_hint"),
		DiscountDaysHint:       t("payment_term.form.discount_days_hint"),
		DiscountPercentBpsHint: t("payment_term.form.discount_percent_bps_hint"),
		PriorityHint:           t("payment_term.form.priority_hint"),
		EntityScope:            t("payment_term.form.entity_scope"),
		IsDefault:              t("payment_term.form.is_default"),
		Description:            t("payment_term.form.description"),
		DescriptionPlaceholder: t("payment_term.form.description_placeholder"),
		DisplayOrder:           t("payment_term.form.display_order"),
		Active:                 t("payment_term.form.active"),

		TypeDueOnReceipt:        t("payment_term.form.type_due_on_receipt"),
		TypeNet:                 t("payment_term.form.type_net"),
		TypeCOD:                 t("payment_term.form.type_cod"),
		TypeProximate:           t("payment_term.form.type_proximate"),
		ProximateDay:            t("payment_term.form.proximate_day"),
		ProximateDayPlaceholder: t("payment_term.form.proximate_day_placeholder"),

		ScopesBoth:              t("payment_term.form.scopes_both"),
		ScopesSupplierOnly:      t("payment_term.form.scopes_supplier_only"),
		ScopesClientOnly:        t("payment_term.form.scopes_client_only"),
		NameInfo:                t("payment_term.form.name_info"),
		CodeInfo:                t("payment_term.form.code_info"),
		CodeHint:                t("payment_term.form.code_hint"),
		DescriptionInfo:         t("payment_term.form.description_info"),
		TypeInfo:                t("payment_term.form.type_info"),
		ProximateDayHint:        t("payment_term.form.proximate_day_hint"),
		ProximateDayInfo:        t("payment_term.form.proximate_day_info"),
		DiscountDaysInfo:        t("payment_term.form.discount_days_info"),
		DiscountPercentBpsInfo:  t("payment_term.form.discount_percent_bps_info"),
		IsDefaultInfo:           t("payment_term.form.is_default_info"),
		ErrTypeRequired:         t("payment_term.form.errors.type_required"),
		ErrTypeInvalid:          t("payment_term.form.errors.type_invalid"),
		ErrNetDaysRequired:      t("payment_term.form.errors.net_days_required"),
		ErrProximateDayRequired: t("payment_term.form.errors.proximate_day_required"),
	}
}

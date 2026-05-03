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
	ErrTypeRequired        string
	ErrTypeInvalid         string
	ErrNetDaysRequired     string
	ErrProximateDayRequired string
}

// Data is the template data for the payment term drawer form.
type Data struct {
	FormAction         string
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
		SectionInfo:            t("paymentTerm.form.sectionInfo"),
		SectionTerms:           t("paymentTerm.form.sectionTerms"),
		SectionSettings:        t("paymentTerm.form.sectionSettings"),
		Name:                   t("paymentTerm.form.name"),
		NamePlaceholder:        t("paymentTerm.form.namePlaceholder"),
		Code:                   t("paymentTerm.form.code"),
		CodePlaceholder:        t("paymentTerm.form.codePlaceholder"),
		Type:                   t("paymentTerm.form.type"),
		NetDays:                t("paymentTerm.form.netDays"),
		DiscountDays:           t("paymentTerm.form.discountDays"),
		DiscountPercentBps:     t("paymentTerm.form.discountPercentBps"),
		TypeHint:               t("paymentTerm.form.typeHint"),
		NetDaysHint:            t("paymentTerm.form.netDaysHint"),
		DiscountDaysHint:       t("paymentTerm.form.discountDaysHint"),
		DiscountPercentBpsHint: t("paymentTerm.form.discountPercentBpsHint"),
		PriorityHint:           t("paymentTerm.form.priorityHint"),
		EntityScope:            t("paymentTerm.form.entityScope"),
		IsDefault:              t("paymentTerm.form.isDefault"),
		Description:            t("paymentTerm.form.description"),
		DescriptionPlaceholder: t("paymentTerm.form.descriptionPlaceholder"),
		DisplayOrder:           t("paymentTerm.form.displayOrder"),
		Active:                 t("paymentTerm.form.active"),

		TypeDueOnReceipt:        t("paymentTerm.form.typeDueOnReceipt"),
		TypeNet:                 t("paymentTerm.form.typeNet"),
		TypeCOD:                 t("paymentTerm.form.typeCOD"),
		TypeProximate:           t("paymentTerm.form.typeProximate"),
		ProximateDay:            t("paymentTerm.form.proximateDay"),
		ProximateDayPlaceholder: t("paymentTerm.form.proximateDayPlaceholder"),

		ScopesBoth:         t("paymentTerm.form.scopesBoth"),
		ScopesSupplierOnly: t("paymentTerm.form.scopesSupplierOnly"),
		ScopesClientOnly:   t("paymentTerm.form.scopesClientOnly"),
		NameInfo:               t("paymentTerm.form.nameInfo"),
		CodeInfo:               t("paymentTerm.form.codeInfo"),
		CodeHint:               t("paymentTerm.form.codeHint"),
		DescriptionInfo:        t("paymentTerm.form.descriptionInfo"),
		TypeInfo:               t("paymentTerm.form.typeInfo"),
		ProximateDayHint:       t("paymentTerm.form.proximateDayHint"),
		ProximateDayInfo:       t("paymentTerm.form.proximateDayInfo"),
		DiscountDaysInfo:       t("paymentTerm.form.discountDaysInfo"),
		DiscountPercentBpsInfo: t("paymentTerm.form.discountPercentBpsInfo"),
		IsDefaultInfo:          t("paymentTerm.form.isDefaultInfo"),
		ErrTypeRequired:        t("paymentTerm.form.errors.typeRequired"),
		ErrTypeInvalid:         t("paymentTerm.form.errors.typeInvalid"),
		ErrNetDaysRequired:     t("paymentTerm.form.errors.netDaysRequired"),
		ErrProximateDayRequired: t("paymentTerm.form.errors.proximateDayRequired"),
	}
}

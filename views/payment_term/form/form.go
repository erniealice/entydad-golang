package form

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
	NameInfo        string
	CodeInfo        string
	DescriptionInfo string
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
	Labels             Labels
	CommonLabels       any
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
		NameInfo:           t("paymentTerm.form.nameInfo"),
		CodeInfo:           t("paymentTerm.form.codeInfo"),
		DescriptionInfo:    t("paymentTerm.form.descriptionInfo"),
	}
}

package form

// Labels holds i18n labels for the drawer form template.
type Labels struct {
	Name                   string
	NamePlaceholder        string
	Description            string
	DescriptionPlaceholder string
	Private                string
	Active                 string

	// Tax section labels (Phase 5)
	SectionTax                  string
	TaxInclusivePricing         string
	TaxInclusivePricingInfo     string
	TaxComputationEnabled       string
	TaxComputationEnabledInfo   string
	HomeJurisdiction            string
	HomeJurisdictionPlaceholder string
	HomeJurisdictionInfo        string
	TIN                         string
	TINPlaceholder              string
	TINInfo                     string
	// Phase 5 M2 — accessible aria-label for tax info buttons (replaces hardcoded "More info")
	MoreInfo string
}

// Data is the template data for the workspace drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit      bool
	ID          string
	Name        string
	Description string
	Private     bool
	Active      bool

	// Tax fields (Phase 5)
	TaxInclusivePricing   bool
	TaxComputationEnabled bool
	HomeJurisdiction      string
	TIN                   string

	Labels       Labels
	CommonLabels any
}

// BuildLabels constructs Labels using the translator function.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		Name:                   t("form.name"),
		NamePlaceholder:        t("form.namePlaceholder"),
		Description:            t("form.description"),
		DescriptionPlaceholder: t("form.descriptionPlaceholder"),
		Private:                t("form.private"),
		Active:                 t("form.active"),

		// Tax section
		SectionTax:                  t("form.sectionTax"),
		TaxInclusivePricing:         t("form.taxInclusivePricing"),
		TaxInclusivePricingInfo:     t("form.taxInclusivePricingInfo"),
		TaxComputationEnabled:       t("form.taxComputationEnabled"),
		TaxComputationEnabledInfo:   t("form.taxComputationEnabledInfo"),
		HomeJurisdiction:            t("form.homeJurisdiction"),
		HomeJurisdictionPlaceholder: t("form.homeJurisdictionPlaceholder"),
		HomeJurisdictionInfo:        t("form.homeJurisdictionInfo"),
		TIN:                         t("form.tin"),
		TINPlaceholder:              t("form.tinPlaceholder"),
		TINInfo:                     t("form.tinInfo"),
		MoreInfo:                    t("form.moreInfo"),
	}
}

package form

// Labels holds i18n labels for the drawer form template.
type Labels struct {
	Name                   string
	NamePlaceholder        string
	Description            string
	DescriptionPlaceholder string
	Private                string
	Active                 string
}

// Data is the template data for the workspace drawer form.
type Data struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Private      bool
	Active       bool
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
	}
}

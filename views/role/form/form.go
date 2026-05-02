package form

// Labels holds i18n labels for the drawer form template.
type Labels struct {
	Name                   string
	NamePlaceholder        string
	Description            string
	DescriptionPlaceholder string
	Color                  string
	ColorPlaceholder       string
	Active                 string
}

// Data is the template data for the role drawer form.
type Data struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Color        string
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
		Color:                  t("form.color"),
		ColorPlaceholder:       t("form.colorPlaceholder"),
		Active:                 t("form.active"),
	}
}

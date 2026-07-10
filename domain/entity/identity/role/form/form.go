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
	WorkspaceID  string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
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
		NamePlaceholder:        t("form.name_placeholder"),
		Description:            t("form.description"),
		DescriptionPlaceholder: t("form.description_placeholder"),
		Color:                  t("form.color"),
		ColorPlaceholder:       t("form.color_placeholder"),
		Active:                 t("form.active"),
	}
}

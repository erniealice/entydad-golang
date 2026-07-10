package form

// Labels holds i18n labels for the user drawer form template.
type Labels struct {
	FirstName                 string
	FirstNamePlaceholder      string
	LastName                  string
	LastNamePlaceholder       string
	Email                     string
	EmailPlaceholder          string
	Mobile                    string
	MobilePlaceholder         string
	Timezone                  string
	TimezonePlaceholder       string
	TimezoneSearchPlaceholder string
	TimezoneNoResults         string
	Password                  string
	PasswordPlaceholder       string
	PasswordGenerate          string
	Active                    string
	TogglePasswordVisibility  string

	// Field-level info text surfaced via an info button beside each label.
	EmailInfo    string
	MobileInfo   string
	TimezoneInfo string
	ActiveInfo   string
}

// Data is the template data for the user drawer form.
type Data struct {
	FormAction         string
	WorkspaceID        string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce              string // injected by C1: populated by ViewAdapter.injectPageData via reflection
	IsEdit             bool
	ID                 string
	FirstName          string
	LastName           string
	Email              string
	Mobile             string
	Timezone           string
	Active             bool
	SearchTimezonesURL string
	Labels             Labels
	CommonLabels       any
}

// BuildLabels constructs a Labels struct from a translation function.
// t is typically viewCtx.T — a narrow func(string) string with no Deps or storage access.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		FirstName:                 t("form.first_name"),
		FirstNamePlaceholder:      t("form.first_name_placeholder"),
		LastName:                  t("form.last_name"),
		LastNamePlaceholder:       t("form.last_name_placeholder"),
		Email:                     t("form.email"),
		EmailPlaceholder:          t("form.email_placeholder"),
		Mobile:                    t("form.mobile"),
		MobilePlaceholder:         t("form.mobile_placeholder"),
		Timezone:                  t("form.timezone"),
		TimezonePlaceholder:       t("form.timezone_placeholder"),
		TimezoneSearchPlaceholder: t("form.timezone_search_placeholder"),
		TimezoneNoResults:         t("form.timezone_no_results"),
		Password:                  t("form.password"),
		PasswordPlaceholder:       t("form.password_placeholder"),
		PasswordGenerate:          t("form.password_generate"),
		Active:                    t("form.active"),
		TogglePasswordVisibility:  t("form.toggle_password_visibility"),
		EmailInfo:                 t("user.form.email_info"),
		MobileInfo:                t("user.form.mobile_info"),
		TimezoneInfo:              t("user.form.timezone_info"),
		ActiveInfo:                t("user.form.active_info"),
	}
}

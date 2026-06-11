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
		FirstName:                 t("form.firstName"),
		FirstNamePlaceholder:      t("form.firstNamePlaceholder"),
		LastName:                  t("form.lastName"),
		LastNamePlaceholder:       t("form.lastNamePlaceholder"),
		Email:                     t("form.email"),
		EmailPlaceholder:          t("form.emailPlaceholder"),
		Mobile:                    t("form.mobile"),
		MobilePlaceholder:         t("form.mobilePlaceholder"),
		Timezone:                  t("form.timezone"),
		TimezonePlaceholder:       t("form.timezonePlaceholder"),
		TimezoneSearchPlaceholder: t("form.timezoneSearchPlaceholder"),
		TimezoneNoResults:         t("form.timezoneNoResults"),
		Password:                  t("form.password"),
		PasswordPlaceholder:       t("form.passwordPlaceholder"),
		PasswordGenerate:          t("form.passwordGenerate"),
		Active:                    t("form.active"),
		TogglePasswordVisibility:  t("form.togglePasswordVisibility"),
		EmailInfo:                 t("user.form.emailInfo"),
		MobileInfo:                t("user.form.mobileInfo"),
		TimezoneInfo:              t("user.form.timezoneInfo"),
		ActiveInfo:                t("user.form.activeInfo"),
	}
}

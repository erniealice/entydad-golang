package form

// form.go — Delegate drawer form data and labels.
// Mirrors party/client/form/form.go trimmed to guardian-only fields.
//
// Label keys deliberately reuse client.form.* translations (firstName,
// lastName, email, phone) rather than adding delegate.form.* to stay within
// the research §4 lyngua scope. This choice is noted for the owner; add
// delegate.form.* keys to general/delegate.json if tier-specific copy is needed.

// Labels holds i18n labels for the delegate drawer form template.
// Fields are populated via BuildLabels which calls t("client.form.*").
type Labels struct {
	SectionGuardian      string
	FirstName            string
	FirstNamePlaceholder string
	LastName             string
	LastNamePlaceholder  string
	Email                string
	EmailPlaceholder     string
	Mobile               string
	MobilePlaceholder    string
}

// Data is the template data for the delegate drawer form.
type Data struct {
	FormAction   string
	WorkspaceID  string // injected by ViewAdapter for actionForm workspace guard
	IsEdit       bool
	ID           string
	FirstName    string
	LastName     string
	Email        string
	Mobile       string
	Active       bool
	Labels       Labels
	CommonLabels any
}

// BuildLabels constructs a Labels struct from a translation function.
// t is typically viewCtx.T — func(string) string.
// Reuses client.form.* keys to stay within research §4 lyngua scope.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		SectionGuardian:      t("delegate.page.heading"),
		FirstName:            t("client.form.firstName"),
		FirstNamePlaceholder: t("client.form.firstNamePlaceholder"),
		LastName:             t("client.form.lastName"),
		LastNamePlaceholder:  t("client.form.lastNamePlaceholder"),
		Email:                t("client.form.email"),
		EmailPlaceholder:     t("client.form.emailPlaceholder"),
		Mobile:               t("client.form.phone"),
		MobilePlaceholder:    t("client.form.phonePlaceholder"),
	}
}

// Package form holds the staff drawer's template-facing data per
// ui-drawer-form-subpackage-convention.md: bare types (Data, Labels), no Deps,
// no context.Context, no repository imports.
package form

import (
	pyeza "github.com/erniealice/pyeza-golang/types"
)

// Labels holds i18n labels for the staff drawer form template.
// Keys are staff-specific (staff.form.*) — unlike party/delegate/form, which
// reuses client.form.* keys, Staff's fields (user picker, employment_type,
// status) have no equivalent on Client, so dedicated keys are warranted.
type Labels struct {
	SectionStaff string

	User                  string
	UserPlaceholder       string
	UserSearchPlaceholder string
	UserNoResults         string

	EmploymentType              string
	EmploymentTypePlaceholder   string
	EmploymentTypeEmployed      string
	EmploymentTypeContractor    string
	EmploymentTypeExternal      string
	EmploymentTypePartner       string
	EmploymentTypeRetained      string
	EmploymentTypeSubcontractor string

	Status            string
	StatusPlaceholder string
	StatusAvailable   string
	StatusAssigned    string
	StatusBench       string
	StatusOffboarded  string
}

// BuildLabels constructs a Labels struct from a translation function.
// t is typically viewCtx.T — func(string) string.
//
// The employment_type/status OPTION labels deliberately read from the
// pre-existing staff.availability.* / staff.employment.* lyngua keys (a
// staff.json skeleton already shipped ahead of this view, from the staff
// use-case wave — see staff.proto's Performance Evaluation §E1/E2 comments)
// rather than inventing a parallel form.status_*/form.employment_type_*
// tree; the field label + placeholder keys (form.status, form.status_placeholder,
// form.employment_type, form.employment_type_placeholder) were already
// present too and are reused as-is.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		SectionStaff: t("staff.page.heading"),

		User:                  t("staff.form.user"),
		UserPlaceholder:       t("staff.form.user_placeholder"),
		UserSearchPlaceholder: t("staff.form.user_search_placeholder"),
		UserNoResults:         t("staff.form.user_no_results"),

		EmploymentType:              t("staff.form.employment_type"),
		EmploymentTypePlaceholder:   t("staff.form.employment_type_placeholder"),
		EmploymentTypeEmployed:      t("staff.employment.employed"),
		EmploymentTypeContractor:    t("staff.employment.contractor"),
		EmploymentTypeExternal:      t("staff.employment.external"),
		EmploymentTypePartner:       t("staff.employment.partner"),
		EmploymentTypeRetained:      t("staff.employment.retained"),
		EmploymentTypeSubcontractor: t("staff.employment.subcontractor"),

		Status:            t("staff.form.status"),
		StatusPlaceholder: t("staff.form.status_placeholder"),
		StatusAvailable:   t("staff.availability.available"),
		StatusAssigned:    t("staff.availability.assigned"),
		StatusBench:       t("staff.availability.bench"),
		StatusOffboarded:  t("staff.availability.offboarded"),
	}
}

// Data is the template data for the staff drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter for actionForm workspace guard
	IsEdit      bool
	ID          string

	UserID            string
	UserSelectedLabel string // "First Last (email)" — pre-fill for the edit drawer
	UserSearchURL     string

	EmploymentType string
	Status         string

	EmploymentTypeOptions []pyeza.SelectOption
	StatusOptions         []pyeza.SelectOption

	Labels       Labels
	CommonLabels any
}

// EmploymentTypeValues are the canonical string forms of the staff.proto
// EmploymentType enum (Staff.employment_type stays a plain string column —
// existing-entity convention — the enum is the value vocabulary only).
const (
	EmploymentTypeEmployed      = "EMPLOYED"
	EmploymentTypeContractor    = "CONTRACTOR"
	EmploymentTypeExternal      = "EXTERNAL"
	EmploymentTypePartner       = "PARTNER"
	EmploymentTypeRetained      = "RETAINED"
	EmploymentTypeSubcontractor = "SUBCONTRACTOR"
)

// StatusValues are the canonical AVAILABILITY values documented on
// staff.proto's status field (available|assigned|bench|offboarded).
const (
	StatusAvailable  = "available"
	StatusAssigned   = "assigned"
	StatusBench      = "bench"
	StatusOffboarded = "offboarded"
)

// BuildEmploymentTypeOptions returns the employment_type select options.
func BuildEmploymentTypeOptions(selected string, labels Labels) []pyeza.SelectOption {
	return []pyeza.SelectOption{
		{Value: EmploymentTypeEmployed, Label: labels.EmploymentTypeEmployed, Selected: selected == EmploymentTypeEmployed},
		{Value: EmploymentTypeContractor, Label: labels.EmploymentTypeContractor, Selected: selected == EmploymentTypeContractor},
		{Value: EmploymentTypeExternal, Label: labels.EmploymentTypeExternal, Selected: selected == EmploymentTypeExternal},
		{Value: EmploymentTypePartner, Label: labels.EmploymentTypePartner, Selected: selected == EmploymentTypePartner},
		{Value: EmploymentTypeRetained, Label: labels.EmploymentTypeRetained, Selected: selected == EmploymentTypeRetained},
		{Value: EmploymentTypeSubcontractor, Label: labels.EmploymentTypeSubcontractor, Selected: selected == EmploymentTypeSubcontractor},
	}
}

// BuildStatusOptions returns the AVAILABILITY status select options shown in
// the staff drawer form. Default on create is "available" (the entity's
// documented default lifecycle state — analogous to client/supplier
// defaulting to "active", see proto-entity-status-conventions.md; Staff's own
// vocabulary is available|assigned|bench|offboarded, not active/inactive).
func BuildStatusOptions(selected string, labels Labels) []pyeza.SelectOption {
	return []pyeza.SelectOption{
		{Value: StatusAvailable, Label: labels.StatusAvailable, Selected: selected == StatusAvailable},
		{Value: StatusAssigned, Label: labels.StatusAssigned, Selected: selected == StatusAssigned},
		{Value: StatusBench, Label: labels.StatusBench, Selected: selected == StatusBench},
		{Value: StatusOffboarded, Label: labels.StatusOffboarded, Selected: selected == StatusOffboarded},
	}
}

// Package form provides the data shape and label builder for the
// TaxRegistration drawer form.
// Used for both client-scoped and workspace-scoped party contexts.
package form

import (
	entydad "github.com/erniealice/entydad-golang"
	pyeza "github.com/erniealice/pyeza-golang"
)

// Labels holds the translatable strings used by the drawer template.
type Labels struct {
	DrawerTitleAdd     string
	DrawerTitleEdit    string
	KindID             string
	RegistrationNumber string
	EffectiveFrom      string
	Notes              string
	Status             string
}

// BuildLabels converts TaxRegistrationLabels.Fields to the flat Labels struct.
func BuildLabels(l entydad.TaxRegistrationLabels) Labels {
	return Labels{
		DrawerTitleAdd:     l.Page.AddDrawerTitle,
		DrawerTitleEdit:    l.Page.EditDrawerTitle,
		KindID:             l.Fields.TaxRegistrationKindID,
		RegistrationNumber: l.Fields.RegistrationNumber,
		EffectiveFrom:      l.Fields.EffectiveFrom,
		Notes:              l.Fields.Notes,
		Status:             l.Fields.Status,
	}
}

// KindOption is a single entry in the Kind <select>.
type KindOption struct {
	Value    string
	Label    string
	Selected bool
}

// Data is the template data for the tax registration drawer form.
type Data struct {
	FormAction string
	WorkspaceID string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit     bool
	ID         string
	// PartyType is "client" or "workspace" — drives the filtered Kind dropdown.
	PartyType string
	PartyID   string
	// KindOptions is filtered by FindByPartyType when the use case is wired,
	// or all kinds when the use case is not yet available (TODO fallback).
	KindOptions        []KindOption
	SelectedKindID     string
	RegistrationNumber string
	EffectiveFrom      string
	Notes              string
	StatusOptions      []pyeza.SelectOption

	Labels       Labels
	CommonLabels any
}

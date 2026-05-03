package form

import pyeza "github.com/erniealice/pyeza-golang/types"

// BuildCurrencyOptions returns the currency select options sourced from
// lyngua's common currency list (CommonLabels.Currency.Options). Each entry
// is a pyeza.SelectOption with Value, Label, and Description already set by
// lyngua; this function only marks the Selected flag. The list order is
// preserved from lyngua (alphabetical by ISO code). selected is the current
// BillingCurrency value — empty string means no option is pre-selected.
// This is a pure function: no Deps, no context, no storage access.
func BuildCurrencyOptions(selected string, currencyOptions []pyeza.SelectOption) []pyeza.SelectOption {
	opts := make([]pyeza.SelectOption, len(currencyOptions))
	for i, o := range currencyOptions {
		opts[i] = pyeza.SelectOption{
			Value:       o.Value,
			Label:       o.Label,
			Description: o.Description,
			Selected:    o.Value == selected,
		}
	}
	return opts
}

// BuildPaymentTermSelectOptions converts a slice of PaymentTermOption into the
// SelectOption format expected by the pyeza form-group select component.
// Mirrors the client/form/options.go helper of the same name.
// This is a pure function: no Deps, no context, no storage access.
func BuildPaymentTermSelectOptions(terms []*PaymentTermOption, selectedID string) []pyeza.SelectOption {
	opts := make([]pyeza.SelectOption, 0, len(terms))
	for _, t := range terms {
		opts = append(opts, pyeza.SelectOption{
			Value:    t.Id,
			Label:    t.Name,
			Selected: t.Id == selectedID,
		})
	}
	return opts
}

// BuildStatusOptions returns the lifecycle-status options shown in the
// supplier drawer form's accounting Status select. Suppliers carry only
// three values today (active / on_hold / blocked) — see entity-status-
// conventions.md. Tier-specific wording lives in lyngua and reaches this
// function via labels — keep the value strings proto-generic.
func BuildStatusOptions(selected string, labels Labels) []pyeza.SelectOption {
	return []pyeza.SelectOption{
		{Value: "active", Label: labels.StatusActive, Selected: selected == "active"},
		{Value: "on_hold", Label: labels.StatusOnHold, Selected: selected == "on_hold"},
		{Value: "blocked", Label: labels.StatusBlocked, Selected: selected == "blocked"},
	}
}

// BuildSupplierTypeOptions returns the two-value supplier_type select options
// (company / individual). Same proto-generic value-string rule as Status.
func BuildSupplierTypeOptions(selected string, labels Labels) []pyeza.SelectOption {
	return []pyeza.SelectOption{
		{Value: "company", Label: labels.TypeCompany, Selected: selected == "company"},
		{Value: "individual", Label: labels.TypeIndividual, Selected: selected == "individual"},
	}
}

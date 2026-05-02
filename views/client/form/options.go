package form

import pyeza "github.com/erniealice/pyeza-golang/types"

// BuildPaymentTermSelectOptions converts a slice of PaymentTermOption into the
// SelectOption format expected by the pyeza form-group select component.
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

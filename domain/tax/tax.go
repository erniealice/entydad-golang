// Package tax is the facade for the entydad slice of the `tax` esqyma domain.
//
// `tax` is a CROSS-PACKAGE-SPLIT domain (manifest thread TX): entydad owns the
// `tax_registration` entity (this package); fycha owns `tax_rate`. Each package
// declares `package tax` on its own import path and re-exports only the entities
// it owns — placement-test R3 (entity ∈ esqyma domain) holds per-package.
//
// Consumers (block/, service-admin) use the TaxRegistration-prefixed names
// exported here (e.g. TaxRegistrationLabels, DefaultTaxRegistrationRoutes) which
// resolve to the entity-local types in domain/tax/tax_registration/. This keeps
// consumer call sites byte-identical to the former entydad.TaxRegistration*
// root names while the internals follow the Option-B entity-local convention
// (one `Labels`/`Routes`/`DefaultRoutes` per entity package).
//
// Import-cycle rule (contract D): the tax_registration entity package MUST NEVER
// import this facade. Only this facade and the hoisted assembler
// (tax_registration_module.go) live in package tax.
package tax

import (
	taxregistration "github.com/erniealice/entydad-golang/domain/tax/tax_registration"
)

// ---------------------------------------------------------------------------
// TaxRegistration (tax/tax_registration)
// ---------------------------------------------------------------------------

type TaxRegistrationLabels = taxregistration.Labels
type TaxRegistrationPageLabels = taxregistration.PageLabels
type TaxRegistrationColumnLabels = taxregistration.ColumnLabels
type TaxRegistrationButtonLabels = taxregistration.ButtonLabels
type TaxRegistrationActionLabels = taxregistration.ActionLabels
type TaxRegistrationEmptyLabels = taxregistration.EmptyLabels
type TaxRegistrationFieldLabels = taxregistration.FieldLabels
type TaxRegistrationRevokeLabels = taxregistration.RevokeLabels

type TaxRegistrationRoutes = taxregistration.Routes

// DefaultTaxRegistrationLabels returns the entity's English label defaults.
func DefaultTaxRegistrationLabels() TaxRegistrationLabels { return taxregistration.DefaultLabels() }

// DefaultTaxRegistrationRoutes returns the entity's route defaults.
func DefaultTaxRegistrationRoutes() TaxRegistrationRoutes { return taxregistration.DefaultRoutes() }

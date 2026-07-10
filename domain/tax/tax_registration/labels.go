package tax_registration

// labels.go — TaxRegistration label structs.
//
// Extracted verbatim from packages/entydad-golang/labels.go (the root tax
// leftovers). Pure structural move — no behaviour change; field names, json
// tags, and string literals are byte-identical. Entity-local rename:
// TaxRegistrationLabels -> Labels, TaxRegistration<Xxx>Labels -> <Xxx>Labels,
// DefaultTaxRegistrationLabels -> DefaultLabels. The facade (domain/tax/tax.go)
// restores the original entydad.TaxRegistration* names for consumers.
//
// Lyngua root key: "taxRegistration"

// Labels holds all translatable strings for the polymorphic
// Tax Registration views (client + workspace party types in v1).
type Labels struct {
	Page    PageLabels   `json:"page"`
	Columns ColumnLabels `json:"columns"`
	Buttons ButtonLabels `json:"buttons"`
	Actions ActionLabels `json:"actions"`
	Empty   EmptyLabels  `json:"empty"`
	Fields  FieldLabels  `json:"fields"`
	Revoke  RevokeLabels `json:"revoke"`
}

// PageLabels holds page heading strings.
type PageLabels struct {
	Heading          string `json:"heading"`
	HeadingClient    string `json:"heading_client"`
	HeadingWorkspace string `json:"heading_workspace"`
	Caption          string `json:"caption"`
	AddDrawerTitle   string `json:"add_drawer_title"`
	EditDrawerTitle  string `json:"edit_drawer_title"`
}

// ColumnLabels holds table column headers.
type ColumnLabels struct {
	KindName           string `json:"kind_name"`
	ComputePath        string `json:"compute_path"`
	PartyRole          string `json:"party_role"`
	Status             string `json:"status"`
	EffectiveFrom      string `json:"effective_from"`
	RegistrationNumber string `json:"registration_number"`
}

// ButtonLabels holds button text.
type ButtonLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

// ActionLabels holds action dropdown labels.
type ActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	NoPermission string `json:"no_permission"`
}

// EmptyLabels holds empty-state strings.
type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FieldLabels holds drawer form field labels.
type FieldLabels struct {
	TaxRegistrationKindID string `json:"tax_registration_kind_id"`
	RegistrationNumber    string `json:"registration_number"`
	EffectiveFrom         string `json:"effective_from"`
	Notes                 string `json:"notes"`
	Status                string `json:"status"`
}

// RevokeLabels holds strings for the revoke confirm dialog.
type RevokeLabels struct {
	WarningMessage        string `json:"warning_message"`
	EffectiveTo           string `json:"effective_to"`
	AffectedPeriodsNotice string `json:"affected_periods_notice"`
	// AffectedPeriodsCount is the row label for the pending-period count (Phase 5 M3).
	AffectedPeriodsCount string `json:"affected_periods_count"`
	// AffectedSubscriptionsCount is the row label for the subscription count (Phase 5 M3).
	AffectedSubscriptionsCount string `json:"affected_subscriptions_count"`
	ReasonLabel                string `json:"reason_label"`
	ReasonPlaceholder          string `json:"reason_placeholder"`
	ConfirmButton              string `json:"confirm_button"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:          "Tax Registrations",
			HeadingClient:    "Client Tax Registrations",
			HeadingWorkspace: "Workspace Tax Registrations",
			Caption:          "Active tax registrations determine compute path during revenue recognition",
			AddDrawerTitle:   "Add Tax Registration",
			EditDrawerTitle:  "Edit Tax Registration",
		},
		Columns: ColumnLabels{
			KindName:           "Kind",
			ComputePath:        "Compute Path",
			PartyRole:          "Party Role",
			Status:             "Status",
			EffectiveFrom:      "Effective From",
			RegistrationNumber: "Registration No.",
		},
		Buttons: ButtonLabels{
			Add:    "Add Tax Registration",
			Edit:   "Edit",
			Delete: "Delete",
		},
		Actions: ActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			NoPermission: "You do not have permission to manage tax registrations",
		},
		Empty: EmptyLabels{
			Title:   "No tax registrations",
			Message: "Add a tax registration to enable tax computation for this party.",
		},
		Fields: FieldLabels{
			TaxRegistrationKindID: "Tax Registration Kind",
			RegistrationNumber:    "Registration Number",
			EffectiveFrom:         "Effective From",
			Notes:                 "Notes",
			Status:                "Status",
		},
		Revoke: RevokeLabels{
			WarningMessage:             "Revoking this registration will affect pending billing periods. Ensure all outstanding periods are settled before proceeding.",
			EffectiveTo:                "Effective To",
			AffectedPeriodsNotice:      "Some pending subscription billing periods fall within the revocation window and may need to be reprocessed.",
			AffectedPeriodsCount:       "Affected billing periods",
			AffectedSubscriptionsCount: "Affected subscriptions",
			ReasonLabel:                "Reason for revocation",
			ReasonPlaceholder:          "Describe why this registration is being revoked",
			ConfirmButton:              "Revoke Registration",
		},
	}
}

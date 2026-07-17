package form

import (
	"log"
	"sort"
	"strconv"
	"strings"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	pyeza "github.com/erniealice/pyeza-golang/types"
)

// Labels holds i18n labels for the client drawer form template.
type Labels struct {
	Name                       string
	NamePlaceholder            string
	CompanyDetails             string
	Representative             string
	FirstName                  string
	FirstNamePlaceholder       string
	LastName                   string
	LastNamePlaceholder        string
	Email                      string
	EmailPlaceholder           string
	Mobile                     string
	MobilePlaceholder          string
	Active                     string
	StreetAddress              string
	StreetAddressPlaceholder   string
	City                       string
	CityPlaceholder            string
	Province                   string
	ProvincePlaceholder        string
	PostalCode                 string
	PostalCodePlaceholder      string
	Notes                      string
	NotesPlaceholder           string
	PaymentTerms               string
	SelectPaymentTerm          string
	Tags                       string
	TagsPlaceholder            string
	TagsSearchPlaceholder      string
	TagsNoResults              string
	Accounting                 string
	BillingCurrency            string
	BillingCurrencyPlaceholder string
	BillingCurrencyInfo        string
	Timezone                   string
	TimezonePlaceholder        string
	TimezoneSearchPlaceholder  string
	TimezoneNoResults          string
	TimezoneInfo               string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo         string
	EmailInfo        string
	MobileInfo       string
	NotesInfo        string
	PaymentTermsInfo string
	TagsInfo         string
	ActiveInfo       string

	// Section + status + address field labels
	Status                string
	StatusPlaceholder     string
	StatusActive          string
	StatusBlocked         string
	StatusOnHold          string
	StatusInactive        string
	StatusProspect        string
	Country               string
	CountryPlaceholder    string
	Website               string
	WebsitePlaceholder    string
	SectionCompany        string
	SectionAddress        string
	SectionRepresentative string
	SectionAccounting     string
	SectionOthers         string
	SectionAttributes     string

	// Accounting fields (mirrored from supplier)
	TaxID                         string
	TaxIDPlaceholder              string
	TaxIDInfo                     string
	RegistrationNumber            string
	RegistrationNumberPlaceholder string
	RegistrationNumberInfo        string
	CreditLimit                   string
	CreditLimitPlaceholder        string
	CreditLimitInfo               string
	LeadTimeDays                  string
	LeadTimeDaysPlaceholder       string
	LeadTimeDaysInfo              string

	// Tax identity fields (Phase 5)
	TIN            string
	TINPlaceholder string
	TINInfo        string

	// CountryCode — ISO 3166-1 alpha-2 (Phase 5 H2).
	// Separate from the legacy Country free-text field; drives jurisdiction lookup.
	CountryCode            string
	CountryCodePlaceholder string
	CountryCodeInfo        string
}

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id   string
	Name string
}

// TagOption represents a tag available for selection in the form.
// Fields named Value/Label to match the pyeza multi-select component template.
type TagOption struct {
	Value    string
	Label    string
	Selected bool
}

// SelectedTag represents a pre-selected tag for chip rendering in the multi-select.
type SelectedTag struct {
	Value string
	Label string
}

// Data is the template data for the client drawer form.
type Data struct {
	FormAction         string
	WorkspaceID        string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit             bool
	ID                 string
	Mode               string
	Name               string
	FirstName          string
	LastName           string
	Email              string
	Mobile             string
	Timezone           string
	Active             bool
	Status             string
	Country            string
	Website            string
	StreetAddress      string
	City               string
	Province           string
	PostalCode         string
	Notes              string
	BillingCurrency    string
	TaxID              string
	RegistrationNumber string
	CreditLimit        string // form-input string; converted to int64 centavos in POST
	LeadTimeDays       string
	SearchTimezonesURL string
	// Tax identity fields (Phase 5)
	TIN                      string
	CountryCode              string
	PaymentTerms             []*PaymentTermOption
	SelectedPaymentTermID    string
	PaymentTermSelectOptions []pyeza.SelectOption
	StatusOptions            []pyeza.SelectOption
	BillingCurrencyOptions   []pyeza.SelectOption
	TagOptions               []TagOption
	SelectedTags             []SelectedTag
	// AttributeFields is the optional generic-attribute (EAV) section rendered
	// between Others and Status. Populated only when active attribute definitions
	// exist for the entity/client module (Q-GSE-10). Empty => section not shown.
	AttributeFields []AttributeField
	Labels          Labels
	CommonLabels    any
}

// BuildLabels constructs a Labels struct from a translation function.
// t is typically viewCtx.T — a narrow func(string) string with no Deps or storage access.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		Name:                          t("client.form.name"),
		NamePlaceholder:               t("client.form.name_placeholder"),
		CompanyDetails:                t("client.form.company_details"),
		Representative:                t("client.form.representative"),
		FirstName:                     t("client.form.first_name"),
		FirstNamePlaceholder:          t("client.form.first_name_placeholder"),
		LastName:                      t("client.form.last_name"),
		LastNamePlaceholder:           t("client.form.last_name_placeholder"),
		Email:                         t("client.form.email"),
		EmailPlaceholder:              t("client.form.email_placeholder"),
		Mobile:                        t("client.form.phone"),
		MobilePlaceholder:             t("client.form.phone_placeholder"),
		Active:                        t("client.form.active"),
		StreetAddress:                 t("client.form.street_address"),
		StreetAddressPlaceholder:      t("client.form.street_address_placeholder"),
		City:                          t("client.form.city"),
		CityPlaceholder:               t("client.form.city_placeholder"),
		Province:                      t("client.form.province"),
		ProvincePlaceholder:           t("client.form.province_placeholder"),
		PostalCode:                    t("client.form.postal_code"),
		PostalCodePlaceholder:         t("client.form.postal_code_placeholder"),
		Notes:                         t("client.form.notes"),
		NotesPlaceholder:              t("client.form.notes_placeholder"),
		PaymentTerms:                  t("client.form.payment_terms"),
		SelectPaymentTerm:             t("client.form.select_payment_term"),
		Tags:                          t("client.form.tags"),
		TagsPlaceholder:               t("client.form.tags_placeholder"),
		TagsSearchPlaceholder:         t("client.form.tags_search_placeholder"),
		TagsNoResults:                 t("client.form.tags_no_results"),
		NameInfo:                      t("client.form.name_info"),
		EmailInfo:                     t("client.form.email_info"),
		MobileInfo:                    t("client.form.mobile_info"),
		NotesInfo:                     t("client.form.notes_info"),
		PaymentTermsInfo:              t("client.form.payment_terms_info"),
		TagsInfo:                      t("client.form.tags_info"),
		ActiveInfo:                    t("client.form.active_info"),
		Accounting:                    t("client.form.accounting"),
		BillingCurrency:               t("client.form.billing_currency"),
		BillingCurrencyPlaceholder:    t("client.form.billing_currency_placeholder"),
		BillingCurrencyInfo:           t("client.form.billing_currency_info"),
		Timezone:                      t("client.form.timezone"),
		TimezonePlaceholder:           t("client.form.timezone_placeholder"),
		TimezoneSearchPlaceholder:     t("client.form.timezone_search_placeholder"),
		TimezoneNoResults:             t("client.form.timezone_no_results"),
		TimezoneInfo:                  t("client.form.timezone_info"),
		Status:                        t("client.form.status"),
		StatusPlaceholder:             t("client.form.status_placeholder"),
		StatusActive:                  t("client.form.status_active"),
		StatusBlocked:                 t("client.form.status_blocked"),
		StatusOnHold:                  t("client.form.status_on_hold"),
		StatusInactive:                t("client.form.status_inactive"),
		StatusProspect:                t("client.form.status_prospect"),
		Country:                       t("client.form.country"),
		CountryPlaceholder:            t("client.form.country_placeholder"),
		Website:                       t("client.form.website"),
		WebsitePlaceholder:            t("client.form.website_placeholder"),
		SectionCompany:                t("client.form.section_company"),
		SectionAddress:                t("client.form.section_address"),
		SectionRepresentative:         t("client.form.section_representative"),
		SectionAccounting:             t("client.form.section_accounting"),
		SectionOthers:                 t("client.form.section_others"),
		SectionAttributes:             t("client.form.section_attributes"),
		TaxID:                         t("client.form.tax_id"),
		TaxIDPlaceholder:              t("client.form.tax_id_placeholder"),
		TaxIDInfo:                     t("client.form.tax_id_info"),
		RegistrationNumber:            t("client.form.registration_number"),
		RegistrationNumberPlaceholder: t("client.form.registration_number_placeholder"),
		RegistrationNumberInfo:        t("client.form.registration_number_info"),
		CreditLimit:                   t("client.form.credit_limit"),
		CreditLimitPlaceholder:        t("client.form.credit_limit_placeholder"),
		CreditLimitInfo:               t("client.form.credit_limit_info"),
		LeadTimeDays:                  t("client.form.lead_time_days"),
		LeadTimeDaysPlaceholder:       t("client.form.lead_time_days_placeholder"),
		LeadTimeDaysInfo:              t("client.form.lead_time_days_info"),
		// Tax identity (Phase 5)
		TIN:            t("client.form.tin"),
		TINPlaceholder: t("client.form.tin_placeholder"),
		TINInfo:        t("client.form.tin_info"),
		// CountryCode (Phase 5 H2)
		CountryCode:            t("client.form.country_code"),
		CountryCodePlaceholder: t("client.form.country_code_placeholder"),
		CountryCodeInfo:        t("client.form.country_code_info"),
	}
}

// AttributeOption is one selectable enum value for an attribute select control.
// Label is the display text (av.label ∥ av.value); Value is the stored value.
type AttributeOption struct {
	Value    string
	Label    string
	Selected bool
}

// AttributeField is one rendered input in the drawer's optional Attributes
// section. Its control family MIRRORS the espyna server validator's
// classifyAttributeControl mapping (attribute_sync.go) EXACTLY so the
// client-side control and native constraints match server-side validation.
type AttributeField struct {
	Code      string            // attribute.code (POST field is "attributes." + Code)
	Name      string            // POST field name ("attributes.<code>")
	Label     string            // display label (attribute.name)
	Type      string            // HTML input type: "text" | "number" | "date" | "select"
	Value     string            // current value (edit pre-fill)
	Required  bool              // attribute.required
	Min       string            // number min (attribute.min_value)
	Max       string            // number max (attribute.max_value)
	Step      string            // "1" (integer) | "any" (decimal)
	MinLength string            // text minlength (attribute.min_length)
	MaxLength string            // text maxlength (attribute.max_length)
	Options   []AttributeOption // select membership (enum rows) or boolean Yes/No
	TestID    string            // stable data-testid ("client-attr-<code>")
}

// attrControl is the derived control family for a definition's data_type.
type attrControl struct {
	kind    string // "text" | "number" | "select" | "boolean" | "date"
	integer bool   // number sub-kind: integral only (step 1)
}

// classifyControl maps a definition's free-string data_type to a control family.
// This MUST stay byte-for-byte equivalent to espyna's classifyAttributeControl
// (packages/espyna-golang/.../entity/client/attribute_sync.go) — the client form
// and the server validator classify the same data_type identically. ok=false =>
// unknown data_type: the field is OMITTED (fail-closed, never a silent accept).
func classifyControl(dataType string) (attrControl, bool) {
	switch strings.ToLower(strings.TrimSpace(dataType)) {
	case "text", "string", "free_text":
		return attrControl{kind: "text"}, true
	case "integer", "int", "number_list":
		return attrControl{kind: "number", integer: true}, true
	case "number", "decimal", "float", "free_number":
		return attrControl{kind: "number", integer: false}, true
	case "option", "enum", "select", "text_list", "color_list":
		return attrControl{kind: "select"}, true
	case "boolean", "bool":
		return attrControl{kind: "boolean"}, true
	case "date":
		return attrControl{kind: "date"}, true
	default:
		return attrControl{}, false
	}
}

// buildAttrOptions turns the active AttributeValue rows into select options,
// preserving the display label (label ∥ value) and ordering by sort_order.
func buildAttrOptions(rows []*commonpb.AttributeValue, current string) []AttributeOption {
	active := make([]*commonpb.AttributeValue, 0, len(rows))
	for _, r := range rows {
		if r != nil && r.GetActive() {
			active = append(active, r)
		}
	}
	sort.SliceStable(active, func(i, j int) bool {
		if active[i].GetSortOrder() != active[j].GetSortOrder() {
			return active[i].GetSortOrder() < active[j].GetSortOrder()
		}
		return active[i].GetValue() < active[j].GetValue()
	})
	opts := make([]AttributeOption, 0, len(active))
	for _, r := range active {
		label := r.GetLabel()
		if label == "" {
			label = r.GetValue()
		}
		opts = append(opts, AttributeOption{
			Value:    r.GetValue(),
			Label:    label,
			Selected: r.GetValue() == current,
		})
	}
	return opts
}

// formatAttrNumber renders a float64 constraint value without trailing zeros for
// use in an HTML min/max attribute.
func formatAttrNumber(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// BuildAttributeFields is a PURE function (no Deps, no context, no I/O) that
// builds the drawer's Attributes-section fields from the active definitions,
// their enum option rows (keyed by attribute_id), and the current client values
// (keyed by attribute code). Fields are sorted deterministically by code.
//
// Unknown data_types are OMITTED (and logged), matching the server validator's
// fail-closed classify: a control the server would reject is never rendered.
func BuildAttributeFields(
	defs []*commonpb.Attribute,
	optionsByAttr map[string][]*commonpb.AttributeValue,
	currentValues map[string]string,
) []AttributeField {
	fields := make([]AttributeField, 0, len(defs))
	for _, def := range defs {
		if def == nil || !def.GetActive() {
			continue
		}
		code := def.GetCode()
		if code == "" {
			continue
		}
		control, ok := classifyControl(def.GetDataType())
		if !ok {
			log.Printf("client attributes: omitting %q — unsupported data_type %q", code, def.GetDataType())
			continue
		}
		current := currentValues[code]
		f := AttributeField{
			Code:     code,
			Name:     "attributes." + code,
			Label:    def.GetName(),
			Value:    current,
			Required: def.GetRequired(),
			TestID:   "client-attr-" + code,
		}
		switch control.kind {
		case "text":
			f.Type = "text"
			if def.MinLength != nil {
				f.MinLength = strconv.Itoa(int(def.GetMinLength()))
			}
			if def.MaxLength != nil {
				f.MaxLength = strconv.Itoa(int(def.GetMaxLength()))
			}
		case "number":
			f.Type = "number"
			if control.integer {
				f.Step = "1"
			} else {
				f.Step = "any"
			}
			if def.MinValue != nil {
				f.Min = formatAttrNumber(def.GetMinValue())
			}
			if def.MaxValue != nil {
				f.Max = formatAttrNumber(def.GetMaxValue())
			}
		case "select":
			f.Type = "select"
			f.Options = buildAttrOptions(optionsByAttr[def.GetId()], current)
		case "boolean":
			f.Type = "select"
			f.Options = []AttributeOption{
				{Value: "true", Label: "Yes", Selected: current == "true"},
				{Value: "false", Label: "No", Selected: current == "false"},
			}
		case "date":
			f.Type = "date"
		}
		fields = append(fields, f)
	}
	sort.SliceStable(fields, func(i, j int) bool { return fields[i].Code < fields[j].Code })
	return fields
}

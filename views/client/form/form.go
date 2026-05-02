package form

import (
	pyeza "github.com/erniealice/pyeza-golang/types"
)

// Labels holds i18n labels for the client drawer form template.
type Labels struct {
	Name                     string
	NamePlaceholder          string
	CompanyDetails           string
	Representative           string
	FirstName                string
	FirstNamePlaceholder     string
	LastName                 string
	LastNamePlaceholder      string
	Email                    string
	EmailPlaceholder         string
	Mobile                   string
	MobilePlaceholder        string
	Active                   string
	StreetAddress            string
	StreetAddressPlaceholder string
	City                     string
	CityPlaceholder          string
	Province                 string
	ProvincePlaceholder      string
	PostalCode               string
	PostalCodePlaceholder    string
	Notes                    string
	NotesPlaceholder         string
	PaymentTerms             string
	SelectPaymentTerm        string
	Tags                     string
	TagsPlaceholder          string
	TagsSearchPlaceholder    string
	TagsNoResults            string
	Accounting                 string
	BillingCurrency            string
	BillingCurrencyPlaceholder string
	BillingCurrencyInfo        string
	Timezone                  string
	TimezonePlaceholder       string
	TimezoneSearchPlaceholder string
	TimezoneNoResults         string
	TimezoneInfo              string

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
	FormAction               string
	IsEdit                   bool
	ID                       string
	Mode                     string
	Name                     string
	FirstName                string
	LastName                 string
	Email                    string
	Mobile                   string
	Timezone                 string
	Active                   bool
	Status                   string
	Country                  string
	Website                  string
	StreetAddress            string
	City                     string
	Province                 string
	PostalCode               string
	Notes                    string
	BillingCurrency          string
	SearchTimezonesURL       string
	PaymentTerms             []*PaymentTermOption
	SelectedPaymentTermID    string
	PaymentTermSelectOptions []pyeza.SelectOption
	TagOptions               []TagOption
	SelectedTags             []SelectedTag
	Labels                   Labels
	CommonLabels             any
}

// BuildLabels constructs a Labels struct from a translation function.
// t is typically viewCtx.T — a narrow func(string) string with no Deps or storage access.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		Name:                     t("client.form.name"),
		NamePlaceholder:          t("client.form.namePlaceholder"),
		CompanyDetails:           t("client.form.companyDetails"),
		Representative:           t("client.form.representative"),
		FirstName:                t("client.form.firstName"),
		FirstNamePlaceholder:     t("client.form.firstNamePlaceholder"),
		LastName:                 t("client.form.lastName"),
		LastNamePlaceholder:      t("client.form.lastNamePlaceholder"),
		Email:                    t("client.form.email"),
		EmailPlaceholder:         t("client.form.emailPlaceholder"),
		Mobile:                   t("client.form.phone"),
		MobilePlaceholder:        t("client.form.phonePlaceholder"),
		Active:                   t("client.form.active"),
		StreetAddress:            t("client.form.streetAddress"),
		StreetAddressPlaceholder: t("client.form.streetAddressPlaceholder"),
		City:                     t("client.form.city"),
		CityPlaceholder:          t("client.form.cityPlaceholder"),
		Province:                 t("client.form.province"),
		ProvincePlaceholder:      t("client.form.provincePlaceholder"),
		PostalCode:               t("client.form.postalCode"),
		PostalCodePlaceholder:    t("client.form.postalCodePlaceholder"),
		Notes:                    t("client.form.notes"),
		NotesPlaceholder:         t("client.form.notesPlaceholder"),
		PaymentTerms:             t("client.form.paymentTerms"),
		SelectPaymentTerm:        t("client.form.selectPaymentTerm"),
		Tags:                     t("client.form.tags"),
		TagsPlaceholder:          t("client.form.tagsPlaceholder"),
		TagsSearchPlaceholder:    t("client.form.tagsSearchPlaceholder"),
		TagsNoResults:            t("client.form.tagsNoResults"),
		NameInfo:                   t("client.form.nameInfo"),
		EmailInfo:                  t("client.form.emailInfo"),
		MobileInfo:                 t("client.form.mobileInfo"),
		NotesInfo:                  t("client.form.notesInfo"),
		PaymentTermsInfo:           t("client.form.paymentTermsInfo"),
		TagsInfo:                   t("client.form.tagsInfo"),
		ActiveInfo:                 t("client.form.activeInfo"),
		Accounting:                 t("client.form.accounting"),
		BillingCurrency:            t("client.form.billingCurrency"),
		BillingCurrencyPlaceholder: t("client.form.billingCurrencyPlaceholder"),
		BillingCurrencyInfo:        t("client.form.billingCurrencyInfo"),
		Timezone:                  t("client.form.timezone"),
		TimezonePlaceholder:       t("client.form.timezonePlaceholder"),
		TimezoneSearchPlaceholder: t("client.form.timezoneSearchPlaceholder"),
		TimezoneNoResults:         t("client.form.timezoneNoResults"),
		TimezoneInfo:              t("client.form.timezoneInfo"),
		Status:                t("client.form.status"),
		StatusPlaceholder:     t("client.form.statusPlaceholder"),
		StatusActive:          t("client.form.statusActive"),
		StatusBlocked:         t("client.form.statusBlocked"),
		StatusOnHold:          t("client.form.statusOnHold"),
		StatusInactive:        t("client.form.statusInactive"),
		StatusProspect:        t("client.form.statusProspect"),
		Country:               t("client.form.country"),
		CountryPlaceholder:    t("client.form.countryPlaceholder"),
		Website:               t("client.form.website"),
		WebsitePlaceholder:    t("client.form.websitePlaceholder"),
		SectionCompany:        t("client.form.sectionCompany"),
		SectionAddress:        t("client.form.sectionAddress"),
		SectionRepresentative: t("client.form.sectionRepresentative"),
		SectionAccounting:     t("client.form.sectionAccounting"),
		SectionOthers:         t("client.form.sectionOthers"),
	}
}

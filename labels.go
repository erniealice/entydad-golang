package entydad

import (
	"strings"
)

// ---------------------------------------------------------------------------
// Client labels
// ---------------------------------------------------------------------------

// ClientLabels holds all translatable strings for the client module.
// JSON tags match the "client" wrapper key in retail/client.json.
type ClientLabels struct {
	Page        ClientPageLabels       `json:"page"`
	Buttons     ClientButtonLabels     `json:"buttons"`
	Columns     ClientColumnLabels     `json:"columns"`
	Empty       ClientEmptyLabels      `json:"empty"`
	Form        ClientFormLabels       `json:"form"`
	Detail      ClientDetailLabels     `json:"detail"`
	BulkActions ClientBulkActionLabels `json:"bulkActions"`
}

type ClientPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingProspect string `json:"headingProspect"`
	HeadingOnHold   string `json:"headingOnHold"`
	HeadingBlocked  string `json:"headingBlocked"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionProspect string `json:"captionProspect"`
	CaptionOnHold   string `json:"captionOnHold"`
	CaptionBlocked  string `json:"captionBlocked"`
	CaptionInactive string `json:"captionInactive"`
}

type ClientButtonLabels struct {
	AddNew string `json:"addNew"`
}

type ClientColumnLabels struct {
	ClientName          string `json:"clientName"`
	Representative      string `json:"representative"`
	Status              string `json:"status"`
	Category            string `json:"category"`
	ActiveSubscriptions string `json:"activeSubscriptions"`
	PaymentTerm         string `json:"paymentTerm"`
	DateCreated         string `json:"dateCreated"`
}

type ClientEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	ProspectTitle   string `json:"prospectTitle"`
	ProspectMessage string `json:"prospectMessage"`
	OnHoldTitle     string `json:"onHoldTitle"`
	OnHoldMessage   string `json:"onHoldMessage"`
	BlockedTitle    string `json:"blockedTitle"`
	BlockedMessage  string `json:"blockedMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type ClientFormLabels struct {
	Email                      string `json:"email"`
	Phone                      string `json:"phone"`
	Name                       string `json:"name"`
	NamePlaceholder            string `json:"namePlaceholder"`
	CompanyDetails             string `json:"companyDetails"`
	Representative             string `json:"representative"`
	StreetAddress              string `json:"streetAddress"`
	StreetAddressPlaceholder   string `json:"streetAddressPlaceholder"`
	City                       string `json:"city"`
	CityPlaceholder            string `json:"cityPlaceholder"`
	Province                   string `json:"province"`
	ProvincePlaceholder        string `json:"provincePlaceholder"`
	PostalCode                 string `json:"postalCode"`
	PostalCodePlaceholder      string `json:"postalCodePlaceholder"`
	Notes                      string `json:"notes"`
	NotesPlaceholder           string `json:"notesPlaceholder"`
	Tags                       string `json:"tags"`
	TagsPlaceholder            string `json:"tagsPlaceholder"`
	TagsSearchPlaceholder      string `json:"tagsSearchPlaceholder"`
	TagsNoResults              string `json:"tagsNoResults"`
	Accounting                 string `json:"accounting"`
	BillingCurrency            string `json:"billingCurrency"`
	BillingCurrencyPlaceholder string `json:"billingCurrencyPlaceholder"`
	BillingCurrencyInfo        string `json:"billingCurrencyInfo"`
	// New fields
	Status             string `json:"status"`
	StatusPlaceholder  string `json:"statusPlaceholder"`
	StatusActive       string `json:"statusActive"`
	StatusBlocked      string `json:"statusBlocked"`
	StatusOnHold       string `json:"statusOnHold"`
	StatusInactive     string `json:"statusInactive"`
	StatusProspect     string `json:"statusProspect"`
	Country            string `json:"country"`
	CountryPlaceholder string `json:"countryPlaceholder"`
	// Phase 5 H2 — ISO 3166 alpha-2 country code separate from legacy Country
	CountryCode            string `json:"countryCode"`
	CountryCodePlaceholder string `json:"countryCodePlaceholder"`
	CountryCodeInfo        string `json:"countryCodeInfo"`
	Website                string `json:"website"`
	WebsitePlaceholder     string `json:"websitePlaceholder"`
	SectionCompany         string `json:"sectionCompany"`
	SectionAddress         string `json:"sectionAddress"`
	SectionRepresentative  string `json:"sectionRepresentative"`
	SectionAccounting      string `json:"sectionAccounting"`
	SectionOthers          string `json:"sectionOthers"`
}

type ClientDetailLabels struct {
	CompanyDetails  ClientCompanyDetailLabels   `json:"companyDetails"`
	Actions         ClientDetailActionLabels    `json:"actions"`
	Profile         ClientDetailSectionLabels   `json:"profile"`
	Company         ClientDetailSectionLabels   `json:"company"`
	Address         ClientDetailSectionLabels   `json:"address"`
	Representative  string                      `json:"representative"`
	NotesSection    ClientDetailSectionLabels   `json:"notesSection"`
	Tags            ClientDetailTagLabels       `json:"tags"`
	PurchaseHistory ClientPurchaseHistoryLabels `json:"purchaseHistory"`
	Tabs            ClientDetailTabLabels       `json:"tabs"`
	// Flat inline labels
	Name                 string `json:"name"`
	RecentOrders         string `json:"recentOrders"`
	ColDate              string `json:"colDate"`
	ColReference         string `json:"colReference"`
	ColAmount            string `json:"colAmount"`
	ColStatus            string `json:"colStatus"`
	PurchaseHistoryEmpty string `json:"purchaseHistoryEmpty"`
	// Subscriptions tab
	AddSubscription         string `json:"addSubscription"`
	EmptySubscriptionsTitle string `json:"emptySubscriptionsTitle"`
	EmptySubscriptions      string `json:"emptySubscriptions"`
	// Statement tab stat card labels
	OutstandingBalance string `json:"outstandingBalance"`
	TotalBilled        string `json:"totalBilled"`
	TotalReceived      string `json:"totalReceived"`
	Invoices           string `json:"invoices"`
	// Statement empty state
	EmptyStatementTitle   string `json:"emptyStatementTitle"`
	EmptyStatementMessage string `json:"emptyStatementMessage"`
	// PriceSchedules tab
	PriceSchedules ClientPriceSchedulesLabels `json:"priceSchedules"`
	// Subscriptions tab column headers + confirm dialogs
	Subscriptions ClientSubscriptionLabels `json:"subscriptions"`
	// Statement tab column headers + totals row
	Statement ClientStatementLabels `json:"statement"`
	// OutstandingTable tab column headers + empty state
	OutstandingTable ClientOutstandingTableLabels `json:"outstandingTable"`
	// RevenueRun drawer labels for the per-client Run Invoices flow
	RevenueRun ClientRevenueRunLabels `json:"revenueRun"`
}

type ClientCompanyDetailLabels struct {
	Status string `json:"status"`
}

// ClientDetailSectionLabels holds a title for a detail page section.
type ClientDetailSectionLabels struct {
	Title string `json:"title"`
}

// ClientDetailTagLabels holds labels for the tags section on the detail page.
type ClientDetailTagLabels struct {
	Title  string `json:"title"`
	NoTags string `json:"noTags"`
}

// ClientPurchaseHistoryLabels holds labels for the purchase history section.
type ClientPurchaseHistoryLabels struct {
	Title         string `json:"title"`
	LifetimeSpend string `json:"lifetimeSpend"`
	TotalOrders   string `json:"totalOrders"`
	AvgOrderValue string `json:"avgOrderValue"`
	LastPurchase  string `json:"lastPurchase"`
	Empty         string `json:"empty"`
}

// ClientDetailTabLabels holds labels for the client detail page tabs.
type ClientDetailTabLabels struct {
	Info               string `json:"info"`
	Representative     string `json:"representative"`
	Subscriptions      string `json:"subscriptions"`
	SubscriptionsSlug  string `json:"subscriptionsSlug"`
	Accounting         string `json:"accounting"`
	History            string `json:"history"`
	Statement          string `json:"statement"`
	PriceSchedules     string `json:"priceSchedules"`
	PriceSchedulesSlug string `json:"priceSchedulesSlug"`
	Attachments        string `json:"attachments"`
	AuditHistory       string `json:"auditHistory"`
	// Phase 2 — polymorphic tax registrations tab
	TaxRegistrations string `json:"taxRegistrations"`
}

// ClientPriceSchedulesLabels holds labels for the PriceSchedules tab on the
// client detail page. All copy uses proto-generic vocabulary ("price schedule",
// "plan", "client"); tier-specific words ("rate card", "package") live only in
// lyngua JSON overrides.
type ClientPriceSchedulesLabels struct {
	Empty           string `json:"empty"`
	AddAction       string `json:"addAction"`
	ColumnName      string `json:"columnName"`
	ColumnDateStart string `json:"columnDateStart"`
	ColumnDateEnd   string `json:"columnDateEnd"`
	ColumnPlanCount string `json:"columnPlanCount"`
}

// ClientSubscriptionLabels holds column headers, actions, and confirm-dialog labels
// for the Subscriptions tab table on the client detail page.
type ClientSubscriptionLabels struct {
	ColumnName           string `json:"columnName"`
	ColumnPlan           string `json:"columnPlan"`
	ColumnStartDate      string `json:"columnStartDate"`
	ColumnEndDate        string `json:"columnEndDate"`
	ConfirmDeleteTitle   string `json:"confirmDeleteTitle"`
	ConfirmDeleteMessage string `json:"confirmDeleteMessage"`
}

// ClientStatementLabels holds column headers and totals-row label for the
// Statement tab table on the client detail page.
type ClientStatementLabels struct {
	ColumnDate        string `json:"columnDate"`
	ColumnType        string `json:"columnType"`
	ColumnReference   string `json:"columnReference"`
	ColumnDescription string `json:"columnDescription"`
	ColumnBilled      string `json:"columnBilled"`
	ColumnReceived    string `json:"columnReceived"`
	ColumnBalance     string `json:"columnBalance"`
	TotalsRowLabel    string `json:"totalsRowLabel"`
}

// ClientOutstandingTableLabels holds column headers and empty-state labels for
// the outstanding-revenue table on the Statement tab.
type ClientOutstandingTableLabels struct {
	Columns          ClientOutstandingTableColumnLabels `json:"columns"`
	Empty            ClientOutstandingTableEmptyLabels  `json:"empty"`
	RunInvoicesLabel string                             `json:"runInvoicesLabel"`
}

// ClientOutstandingTableColumnLabels holds column header labels for the
// outstanding-revenue table.
type ClientOutstandingTableColumnLabels struct {
	Date        string `json:"date"`
	Reference   string `json:"reference"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	Billed      string `json:"billed"`
	Paid        string `json:"paid"`
	Outstanding string `json:"outstanding"`
	Status      string `json:"status"`
}

// ClientOutstandingTableEmptyLabels holds empty-state labels for the
// outstanding-revenue table.
type ClientOutstandingTableEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// ClientRevenueRunLabels holds labels for the per-client Run Invoices
// drawer surfaced from the Statement-tab outstanding table.
type ClientRevenueRunLabels struct {
	Title                   string                      `json:"title"`
	SubtitleTemplate        string                      `json:"subtitleTemplate"`
	AsOfDateLabel           string                      `json:"asOfDateLabel"`
	AsOfDateHint            string                      `json:"asOfDateHint"`
	BillThroughTodayLabel   string                      `json:"billThroughTodayLabel"`
	ColumnSubscription      string                      `json:"columnSubscription"`
	ColumnPeriod            string                      `json:"columnPeriod"`
	ColumnAmount            string                      `json:"columnAmount"`
	ColumnLines             string                      `json:"columnLines"`
	GroupTotalLabel         string                      `json:"groupTotalLabel"`
	GroupNoPending          string                      `json:"groupNoPending"`
	GroupCurrencyMismatch   string                      `json:"groupCurrencyMismatch"`
	ColumnSelectAriaLabel   string                      `json:"columnSelectAriaLabel"`
	EmptyTitle              string                      `json:"emptyTitle"`
	EmptyMessage            string                      `json:"emptyMessage"`
	IntroMessage            string                      `json:"introMessage"`
	GenerateButton          string                      `json:"generateButton"`
	GenerateButtonCountOne  string                      `json:"generateButtonCountOne"`
	GenerateButtonCountMany string                      `json:"generateButtonCountMany"`
	CancelButton            string                      `json:"cancelButton"`
	ToastSuccess            string                      `json:"toastSuccess"`
	ViewRunLink             string                      `json:"viewRunLink"`
	Errors                  ClientRevenueRunErrorLabels `json:"errors"`
}

// ClientRevenueRunErrorLabels — error copy surfaced in the drawer.
type ClientRevenueRunErrorLabels struct {
	PermissionDenied   string `json:"permissionDenied"`
	IDRequired         string `json:"idRequired"`
	InvalidFormData    string `json:"invalidFormData"`
	UseCaseUnavailable string `json:"useCaseUnavailable"`
	SelectOne          string `json:"selectOne"`
}

// ResolveTabSlug returns the URL slug for a canonical tab key. Tier-specific
// slugs flow through here so URLs match the operator's vocabulary (e.g.
// professional ships "engagements" + "rate-cards"). Tabs without overrides
// round-trip through unchanged.
func (t ClientDetailTabLabels) ResolveTabSlug(canonical string) string {
	switch canonical {
	case "subscriptions":
		if s := strings.TrimSpace(t.SubscriptionsSlug); s != "" {
			return s
		}
	case "priceSchedules":
		if s := strings.TrimSpace(t.PriceSchedulesSlug); s != "" {
			return s
		}
	}
	return canonical
}

// CanonicalizeTab maps an incoming URL tab slug back to its canonical key so
// internal template lookups and equality checks stay tier-agnostic.
func (t ClientDetailTabLabels) CanonicalizeTab(slug string) string {
	if slug == "" {
		return ""
	}
	if s := strings.TrimSpace(t.SubscriptionsSlug); s != "" && slug == s {
		return "subscriptions"
	}
	if s := strings.TrimSpace(t.PriceSchedulesSlug); s != "" && slug == s {
		return "priceSchedules"
	}
	return slug
}

type ClientDetailActionLabels struct {
	ViewClient       string `json:"viewClient"`
	EditClient       string `json:"editClient"`
	DeleteClient     string `json:"deleteClient"`
	DeactivateClient string `json:"deactivateClient"`
	ActivateClient   string `json:"activateClient"`
	BlockClient      string `json:"blockClient"`
	HoldClient       string `json:"holdClient"`
	SetProspect      string `json:"setProspect"`
}

type ClientBulkActionLabels struct {
	SetAsInactive string `json:"setAsInactive"`
	SetAsActive   string `json:"setAsActive"`
	SetAsBlocked  string `json:"setAsBlocked"`
	SetAsOnHold   string `json:"setAsOnHold"`
	SetAsProspect string `json:"setAsProspect"`
}

// ---------------------------------------------------------------------------
// User labels
// ---------------------------------------------------------------------------

// UserLabels holds all translatable strings for the user module.
// JSON tags match retail/user.json (no wrapper key).
type UserLabels struct {
	Page    UserPageLabels   `json:"page"`
	Buttons UserButtonLabels `json:"buttons"`
	Columns UserColumnLabels `json:"columns"`
	Empty   UserEmptyLabels  `json:"empty"`
	Form    UserFormLabels   `json:"form"`
	Actions UserActionLabels `json:"actions"`
	Detail  UserDetailLabels `json:"detail"`
}

type UserPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type UserButtonLabels struct {
	AddUser string `json:"addUser"`
}

type UserColumnLabels struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Roles       string `json:"roles"`
	Workspaces  string `json:"workspaces"`
	DateCreated string `json:"dateCreated"`
	Status      string `json:"status"`
}

type UserEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type UserFormLabels struct {
	Mobile              string `json:"mobile"`
	Timezone            string `json:"timezone"`
	TimezonePlaceholder string `json:"timezonePlaceholder"`
	TimezoneInfo        string `json:"timezoneInfo"`
}

type UserActionLabels struct {
	View        string `json:"view"`
	Edit        string `json:"edit"`
	Delete      string `json:"delete"`
	Activate    string `json:"activate"`
	Deactivate  string `json:"deactivate"`
	ManageRoles string `json:"manageRoles"`
}

// UserDetailLabels holds labels for the user detail page.
type UserDetailLabels struct {
	BasicInfo   UserDetailBasicInfoLabels  `json:"basicInfo"`
	Tabs        UserDetailTabLabels        `json:"tabs"`
	Security    UserDetailSecurityLabels   `json:"security"`
	EmptyStates UserDetailEmptyStateLabels `json:"emptyStates"`
	// Inline feedback and empty-state messages
	UpdateSuccess            string `json:"updateSuccess"`
	UpdateError              string `json:"updateError"`
	NoRolesAssigned          string `json:"noRolesAssigned"`
	NoRolesDesc              string `json:"noRolesDesc"`
	NewPasswordPlaceholder   string `json:"newPasswordPlaceholder"`
	TogglePasswordVisibility string `json:"togglePasswordVisibility"`
	GeneratePassword         string `json:"generatePassword"`
	PasswordUpdated          string `json:"passwordUpdated"`
	PasswordFailed           string `json:"passwordFailed"`
	// Tab label for attachments (shared across all detail pages)
	AttachmentsTab string `json:"attachmentsTab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"auditHistoryTab"`
}

// UserDetailSecurityLabels holds labels for the security tab.
type UserDetailSecurityLabels struct {
	Title           string `json:"title"`
	LastLogin       string `json:"lastLogin"`
	MfaStatus       string `json:"mfaStatus"`
	MfaEnabled      string `json:"mfaEnabled"`
	MfaDisabled     string `json:"mfaDisabled"`
	PasswordSection string `json:"passwordSection"`
	ResetPassword   string `json:"resetPassword"`
}

// UserDetailEmptyStateLabels holds empty-state labels for user detail tabs.
type UserDetailEmptyStateLabels struct {
	AuditTitle string `json:"auditTitle"`
	AuditDesc  string `json:"auditDesc"`
}

type UserDetailBasicInfoLabels struct {
	Title                string `json:"title"`
	FirstName            string `json:"firstName"`
	FirstNamePlaceholder string `json:"firstNamePlaceholder"`
	LastName             string `json:"lastName"`
	LastNamePlaceholder  string `json:"lastNamePlaceholder"`
	Email                string `json:"email"`
	EmailPlaceholder     string `json:"emailPlaceholder"`
	Username             string `json:"username"`
	Division             string `json:"division"`
	Status               string `json:"status"`
	UserType             string `json:"userType"`
	Mobile               string `json:"mobile"`
	MobilePlaceholder    string `json:"mobilePlaceholder"`
	Active               string `json:"active"`
	Save                 string `json:"save"`
}

type UserDetailTabLabels struct {
	Info       string `json:"info"`
	Roles      string `json:"roles"`
	Security   string `json:"security"`
	AuditTrail string `json:"auditTrail"`
}

// ---------------------------------------------------------------------------
// Location labels
// ---------------------------------------------------------------------------

// LocationLabels holds all translatable strings for the location module.
// JSON tags match the "location" wrapper key in retail/location.json.
type LocationLabels struct {
	Page      LocationPageLabels      `json:"page"`
	Buttons   LocationButtonLabels    `json:"buttons"`
	Columns   LocationColumnLabels    `json:"columns"`
	Empty     LocationEmptyLabels     `json:"empty"`
	Form      LocationFormLabels      `json:"form"`
	Actions   LocationActionLabels    `json:"actions"`
	Detail    LocationDetailLabels    `json:"detail"`
	Dashboard LocationDashboardLabels `json:"dashboard"`
}

type LocationPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type LocationButtonLabels struct {
	AddLocation string `json:"addLocation"`
}

type LocationColumnLabels struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Timezone    string `json:"timezone"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
}

type LocationEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type LocationFormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Address                string `json:"address"`
	AddressPlaceholder     string `json:"addressPlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Timezone               string `json:"timezone"`
	Area                   string `json:"area"`
	AreaPlaceholder        string `json:"areaPlaceholder"`
	Active                 string `json:"active"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	AddressInfo     string `json:"addressInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	TimezoneInfo    string `json:"timezoneInfo"`
	AreaInfo        string `json:"areaInfo"`
	ActiveInfo      string `json:"activeInfo"`
}

type LocationActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type LocationDetailLabels struct {
	BasicInfo   LocationDetailBasicInfoLabels `json:"basicInfo"`
	Tabs        LocationDetailTabLabels       `json:"tabs"`
	EmptyStates LocationDetailEmptyLabels     `json:"emptyStates"`
	// Inline feedback messages
	UpdateSuccess string `json:"updateSuccess"`
	UpdateError   string `json:"updateError"`
	// Tab label for attachments
	AttachmentsTab string `json:"attachmentsTab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"auditHistoryTab"`
}

type LocationDetailBasicInfoLabels struct {
	Title                  string `json:"title"`
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Address                string `json:"address"`
	AddressPlaceholder     string `json:"addressPlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Active                 string `json:"active"`
	Save                   string `json:"save"`
}

type LocationDetailTabLabels struct {
	Info       string `json:"info"`
	Users      string `json:"users"`
	PriceLists string `json:"priceLists"`
	AuditTrail string `json:"auditTrail"`
}

type LocationDetailEmptyLabels struct {
	UsersTitle      string `json:"usersTitle"`
	UsersDesc       string `json:"usersDesc"`
	PriceListsTitle string `json:"priceListsTitle"`
	PriceListsDesc  string `json:"priceListsDesc"`
	AuditTitle      string `json:"auditTitle"`
	AuditDesc       string `json:"auditDesc"`
}

// ---------------------------------------------------------------------------
// LocationArea labels
// ---------------------------------------------------------------------------

// LocationAreaLabels holds all translatable strings for the location area module.
type LocationAreaLabels struct {
	Page    LocationAreaPageLabels   `json:"page"`
	Buttons LocationAreaButtonLabels `json:"buttons"`
	Columns LocationAreaColumnLabels `json:"columns"`
	Empty   LocationAreaEmptyLabels  `json:"empty"`
	Form    LocationAreaFormLabels   `json:"form"`
	Actions LocationAreaActionLabels `json:"actions"`
	Errors  LocationAreaErrorLabels  `json:"errors"`
}

type LocationAreaPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type LocationAreaButtonLabels struct {
	AddLocationArea string `json:"addLocationArea"`
}

type LocationAreaColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
}

type LocationAreaEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type LocationAreaFormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Active                 string `json:"active"`
}

type LocationAreaActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type LocationAreaErrorLabels struct {
	CannotDeleteInUse string `json:"cannotDeleteInUse"`
}

// DefaultLocationAreaLabels returns sensible English defaults for LocationAreaLabels.
func DefaultLocationAreaLabels() LocationAreaLabels {
	return LocationAreaLabels{
		Page: LocationAreaPageLabels{
			Heading:         "Location Areas",
			HeadingActive:   "Active Location Areas",
			HeadingInactive: "Inactive Location Areas",
			Caption:         "Manage location areas",
			CaptionActive:   "Active location areas",
			CaptionInactive: "Inactive location areas",
		},
		Buttons: LocationAreaButtonLabels{
			AddLocationArea: "Add Location Area",
		},
		Columns: LocationAreaColumnLabels{
			Name:        "Name",
			Description: "Description",
			Status:      "Status",
			DateCreated: "Date Created",
		},
		Empty: LocationAreaEmptyLabels{
			ActiveTitle:     "No active location areas",
			ActiveMessage:   "Add your first location area to get started.",
			InactiveTitle:   "No inactive location areas",
			InactiveMessage: "Inactive location areas will appear here.",
		},
		Form: LocationAreaFormLabels{
			Name:                   "Name",
			NamePlaceholder:        "Enter name...",
			Description:            "Description",
			DescriptionPlaceholder: "Enter description...",
			Active:                 "Active",
		},
		Actions: LocationAreaActionLabels{
			View:       "View",
			Edit:       "Edit",
			Delete:     "Delete",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Errors: LocationAreaErrorLabels{
			CannotDeleteInUse: "Cannot delete — this location area is in use.",
		},
	}
}

// ---------------------------------------------------------------------------
// Role labels
// ---------------------------------------------------------------------------

// RoleLabels holds all translatable strings for the role module.
// JSON tags match the "role" wrapper key in retail/role.json.
type RoleLabels struct {
	Page    RolePageLabels   `json:"page"`
	Buttons RoleButtonLabels `json:"buttons"`
	Columns RoleColumnLabels `json:"columns"`
	Empty   RoleEmptyLabels  `json:"empty"`
	Form    RoleFormLabels   `json:"form"`
	Actions RoleActionLabels `json:"actions"`
	Detail  RoleDetailLabels `json:"detail"`
}

type RolePageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type RoleButtonLabels struct {
	AddRole string `json:"addRole"`
}

type RoleColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Permissions string `json:"permissions"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
}

type RoleEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type RoleFormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Color                  string `json:"color"`
	ColorPlaceholder       string `json:"colorPlaceholder"`
	Active                 string `json:"active"`
}

type RoleActionLabels struct {
	View              string `json:"view"`
	Edit              string `json:"edit"`
	Delete            string `json:"delete"`
	Activate          string `json:"activate"`
	Deactivate        string `json:"deactivate"`
	ManagePermissions string `json:"managePermissions"`
}

// RoleDetailLabels holds labels for the role detail page.
type RoleDetailLabels struct {
	Tabs RoleDetailTabLabels  `json:"tabs"`
	Info RoleDetailInfoLabels `json:"info"`
	// Empty-state labels for role detail tabs
	NoPermissionsAssigned string `json:"noPermissionsAssigned"`
	NoPermissionsDesc     string `json:"noPermissionsDesc"`
	NoUsersAssigned       string `json:"noUsersAssigned"`
	NoUsersDesc           string `json:"noUsersDesc"`
	// Tab label for attachments
	AttachmentsTab string `json:"attachmentsTab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"auditHistoryTab"`
}

type RoleDetailTabLabels struct {
	Info        string `json:"info"`
	Permissions string `json:"permissions"`
	Users       string `json:"users"`
}

type RoleDetailInfoLabels struct {
	Title       string `json:"title"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Status      string `json:"status"`
}

// ---------------------------------------------------------------------------
// Permission labels
// ---------------------------------------------------------------------------

// PermissionLabels holds all translatable strings for the permission module.
type PermissionLabels struct {
	Page    PermissionPageLabels   `json:"page"`
	Buttons PermissionButtonLabels `json:"buttons"`
	Columns PermissionColumnLabels `json:"columns"`
	Empty   PermissionEmptyLabels  `json:"empty"`
	Form    PermissionFormLabels   `json:"form"`
	Actions PermissionActionLabels `json:"actions"`
}

type PermissionPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type PermissionButtonLabels struct {
	AddPermission string `json:"addPermission"`
}

type PermissionColumnLabels struct {
	Name           string `json:"name"`
	Entity         string `json:"entity"`
	PermissionCode string `json:"permissionCode"`
	Type           string `json:"type"`
	Status         string `json:"status"`
}

type PermissionEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type PermissionFormLabels struct {
	Name                      string `json:"name"`
	NamePlaceholder           string `json:"namePlaceholder"`
	PermissionCode            string `json:"permissionCode"`
	PermissionCodePlaceholder string `json:"permissionCodePlaceholder"`
	PermissionCodeHint        string `json:"permissionCodeHint"`
	PermissionType            string `json:"permissionType"`
	Description               string `json:"description"`
	DescriptionPlaceholder    string `json:"descriptionPlaceholder"`
	Active                    string `json:"active"`
}

type PermissionActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

// ---------------------------------------------------------------------------
// Role-Permission labels
// ---------------------------------------------------------------------------

// RolePermissionLabels holds all translatable strings for the role-permission assignment view.
type RolePermissionLabels struct {
	Page    RolePermissionPageLabels   `json:"page"`
	Buttons RolePermissionButtonLabels `json:"buttons"`
	Columns RolePermissionColumnLabels `json:"columns"`
	Empty   RolePermissionEmptyLabels  `json:"empty"`
	Form    RolePermissionFormLabels   `json:"form"`
	Actions RolePermissionActionLabels `json:"actions"`
}

type RolePermissionPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type RolePermissionButtonLabels struct {
	AssignPermission string `json:"assignPermission"`
}

type RolePermissionColumnLabels struct {
	PermissionName string `json:"permissionName"`
	Code           string `json:"code"`
	Type           string `json:"type"`
	DateAssigned   string `json:"dateAssigned"`
}

type RolePermissionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type RolePermissionFormLabels struct {
	Permission string `json:"permission"`
}

type RolePermissionActionLabels struct {
	Assign            string `json:"assign"`
	Remove            string `json:"remove"`
	ManagePermissions string `json:"managePermissions"`
}

// ---------------------------------------------------------------------------
// User-Role labels
// ---------------------------------------------------------------------------

// UserRoleLabels holds all translatable strings for the user-role assignment view.
type UserRoleLabels struct {
	Page    UserRolePageLabels   `json:"page"`
	Buttons UserRoleButtonLabels `json:"buttons"`
	Columns UserRoleColumnLabels `json:"columns"`
	Empty   UserRoleEmptyLabels  `json:"empty"`
	Form    UserRoleFormLabels   `json:"form"`
	Actions UserRoleActionLabels `json:"actions"`
}

type UserRolePageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type UserRoleButtonLabels struct {
	AssignRole string `json:"assignRole"`
}

type UserRoleColumnLabels struct {
	RoleName     string `json:"roleName"`
	Description  string `json:"description"`
	Color        string `json:"color"`
	DateAssigned string `json:"dateAssigned"`
}

type UserRoleEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type UserRoleFormLabels struct {
	Role string `json:"role"`
}

type UserRoleActionLabels struct {
	Assign      string `json:"assign"`
	Remove      string `json:"remove"`
	ManageRoles string `json:"manageRoles"`
}

// ---------------------------------------------------------------------------
// Role-User labels (reverse of User-Role: managing users on a role)
// ---------------------------------------------------------------------------

// RoleUserLabels holds all translatable strings for the role-user assignment view.
type RoleUserLabels struct {
	Page    RoleUserPageLabels   `json:"page"`
	Buttons RoleUserButtonLabels `json:"buttons"`
	Columns RoleUserColumnLabels `json:"columns"`
	Empty   RoleUserEmptyLabels  `json:"empty"`
	Form    RoleUserFormLabels   `json:"form"`
	Actions RoleUserActionLabels `json:"actions"`
}

type RoleUserPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type RoleUserButtonLabels struct {
	AssignUser string `json:"assignUser"`
}

type RoleUserColumnLabels struct {
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	DateAssigned string `json:"dateAssigned"`
}

type RoleUserEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type RoleUserFormLabels struct {
	User   string `json:"user"`
	Assign string `json:"assign"`
}

type RoleUserActionLabels struct {
	Assign string `json:"assign"`
	Remove string `json:"remove"`
}

// ---------------------------------------------------------------------------
// Workspace labels
// ---------------------------------------------------------------------------

// WorkspaceLabels holds all translatable strings for the workspace module.
type WorkspaceLabels struct {
	Page    WorkspacePageLabels   `json:"page"`
	Buttons WorkspaceButtonLabels `json:"buttons"`
	Columns WorkspaceColumnLabels `json:"columns"`
	Empty   WorkspaceEmptyLabels  `json:"empty"`
	Form    WorkspaceFormLabels   `json:"form"`
	Actions WorkspaceActionLabels `json:"actions"`
	Detail  WorkspaceDetailLabels `json:"detail"`
}

// WorkspaceDetailLabels holds i18n strings for the workspace detail page (Phase 1).
type WorkspaceDetailLabels struct {
	Tabs  WorkspaceDetailTabLabels  `json:"tabs"`
	Users WorkspaceDetailUserLabels `json:"users"`
	// Info holds the Info-tab section title + non-Form/Column field labels.
	// Name/Description/Status reuse Form.* and Columns.Status; Info supplies
	// the section heading + the Currency/Region fields not present elsewhere.
	Info WorkspaceDetailInfoLabels `json:"info"`
	// TaxReg holds the Tax Registrations tab panel copy (W4.5 label
	// remediation — previously hardcoded in detail.html).
	TaxReg WorkspaceDetailTaxRegLabels `json:"taxReg"`
}

// WorkspaceDetailTaxRegLabels holds the Tax Registrations tab panel copy
// (W4.5 label remediation).
type WorkspaceDetailTaxRegLabels struct {
	Loading       string `json:"loading"`
	Title         string `json:"title"`
	NotConfigured string `json:"notConfigured"`
}

// WorkspaceDetailInfoLabels holds the Info-tab labels on the workspace detail
// page (W4.5 label remediation). SectionTitle is the "Details" heading;
// Currency/Region are the optional workspace fields shown on the Info tab.
type WorkspaceDetailInfoLabels struct {
	SectionTitle string `json:"sectionTitle"`
	Currency     string `json:"currency"`
	Region       string `json:"region"`
}

// WorkspaceDetailTabLabels holds the tab display names for the workspace detail page.
type WorkspaceDetailTabLabels struct {
	Info        string `json:"info"`
	Users       string `json:"users"`
	Attachments string `json:"attachments"`
	// Phase 2 — polymorphic tax registrations tab
	TaxRegistrations string `json:"taxRegistrations"`
}

// WorkspaceDetailUserLabels holds i18n strings for the Users tab on the workspace detail page.
type WorkspaceDetailUserLabels struct {
	AddButton string `json:"addButton"`
	// Empty-state copy shown when the workspace has no user assignments yet
	// (W4.5 label remediation — previously hardcoded in users-tab.html).
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

type WorkspacePageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type WorkspaceButtonLabels struct {
	AddWorkspace string `json:"addWorkspace"`
}

type WorkspaceColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     string `json:"private"`
	Status      string `json:"status"`
}

type WorkspaceEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type WorkspaceFormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Private                string `json:"private"`
	Active                 string `json:"active"`
}

type WorkspaceActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

// ---------------------------------------------------------------------------
// WorkspaceUser labels
// ---------------------------------------------------------------------------

// WorkspaceUserLabels holds all translatable strings for the workspace_user module.
type WorkspaceUserLabels struct {
	Page    WorkspaceUserPageLabels   `json:"page"`
	Columns WorkspaceUserColumnLabels `json:"columns"`
	Detail  WorkspaceUserDetailLabels `json:"detail"`
	Form    WorkspaceUserFormLabels   `json:"form"`
	Actions WorkspaceUserActionLabels `json:"actions"`
}

type WorkspaceUserPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type WorkspaceUserColumnLabels struct {
	UserName   string `json:"userName"`
	Email      string `json:"email"`
	Roles      string `json:"roles"`
	Status     string `json:"status"`
	RoleName   string `json:"roleName"`
	PermCount  string `json:"permCount"`
	DateJoined string `json:"dateJoined"`
}

// WorkspaceUserDetailLabels holds i18n strings for the workspace_user detail page (Phase 2).
type WorkspaceUserDetailLabels struct {
	BackToWorkspace string                         `json:"backToWorkspace"`
	Tabs            WorkspaceUserDetailTabLabels   `json:"tabs"`
	Roles           WorkspaceUserDetailRolesLabels `json:"roles"`
	// Info holds the Info-tab section title + field labels (W4.5 label
	// remediation — previously hardcoded in info-tab.html).
	Info WorkspaceUserDetailInfoLabels `json:"info"`
}

// WorkspaceUserDetailInfoLabels holds the Info-tab labels on the
// workspace_user detail page (W4.5 label remediation).
type WorkspaceUserDetailInfoLabels struct {
	SectionTitle string `json:"sectionTitle"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Workspace    string `json:"workspace"`
	DateJoined   string `json:"dateJoined"`
	Status       string `json:"status"`
}

type WorkspaceUserDetailTabLabels struct {
	Info        string `json:"info"`
	Roles       string `json:"roles"`
	Attachments string `json:"attachments"`
}

type WorkspaceUserDetailRolesLabels struct {
	AssignButton string `json:"assignButton"`
	// Empty-state copy shown when no roles are assigned yet (W4.5 label
	// remediation — previously hardcoded in roles-tab.html).
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

type WorkspaceUserFormLabels struct {
	User                  string `json:"user"`
	UserPlaceholder       string `json:"userPlaceholder"`
	UserSearchPlaceholder string `json:"userSearchPlaceholder"`
	WorkspaceID           string `json:"workspaceId"`
	Active                string `json:"active"`
}

type WorkspaceUserActionLabels struct {
	View       string `json:"view"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

// ---------------------------------------------------------------------------
// WorkspaceUserRoleLabels
// ---------------------------------------------------------------------------

// WorkspaceUserRoleLabels holds all translatable strings for the
// workspace_user_role assignment drawer (Phase 3).
type WorkspaceUserRoleLabels struct {
	Form    WorkspaceUserRoleFormLabels   `json:"form"`
	Buttons WorkspaceUserRoleButtonLabels `json:"buttons"`
}

// WorkspaceUserRoleFormLabels holds field labels for the assign-form drawer.
type WorkspaceUserRoleFormLabels struct {
	WorkspaceUser         string `json:"workspaceUser"`
	Role                  string `json:"role"`
	RolePlaceholder       string `json:"rolePlaceholder"`
	RoleSearchPlaceholder string `json:"roleSearchPlaceholder"`
	RoleNoResults         string `json:"roleNoResults"`
	Permissions           string `json:"permissions"`
	PermissionsHint       string `json:"permissionsHint"`
}

// WorkspaceUserRoleButtonLabels holds button text for the assign-form drawer.
type WorkspaceUserRoleButtonLabels struct {
	Submit string `json:"submit"`
	Cancel string `json:"cancel"`
}

// ---------------------------------------------------------------------------
// Login labels
// ---------------------------------------------------------------------------

// LoginLabels holds i18n strings for the login page.
type LoginLabels struct {
	Title              string `json:"title"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Submit             string `json:"submit"`
	ForgotLink         string `json:"forgotLink"`
	Error              string `json:"error"`
	AdminTitle         string `json:"adminTitle"`
	AdminDescription   string `json:"adminDescription"`
	EmailPlaceholder   string `json:"emailPlaceholder"`
	StaffTitle         string `json:"staffTitle"`
	StaffDescription   string `json:"staffDescription"`
	StaffPinComingSoon string `json:"staffPinComingSoon"`
}

// Login02Labels holds i18n strings for the login02 split-screen page.
type Login02Labels struct {
	Title               string `json:"title"`
	Heading             string `json:"heading"`
	Subheading          string `json:"subheading"`
	EmailLabel          string `json:"emailLabel"`
	EmailPlaceholder    string `json:"emailPlaceholder"`
	PasswordLabel       string `json:"passwordLabel"`
	PasswordPlaceholder string `json:"passwordPlaceholder"`
	RememberMe          string `json:"rememberMe"`
	ForgotPassword      string `json:"forgotPassword"`
	SignInButton        string `json:"signInButton"`
	NoAccount           string `json:"noAccount"`
	SignUpLink          string `json:"signUpLink"`
	SocialDivider       string `json:"socialDivider"`
	Error               string `json:"error"`
	// Carousel navigation
	PreviousSlide string `json:"previousSlide"`
	NextSlide     string `json:"nextSlide"`
	ContinueWith  string `json:"continueWith"`
}

// ---------------------------------------------------------------------------
// Signup labels
// ---------------------------------------------------------------------------

// SignupLabels holds i18n strings for the signup01 page (dual-card style).
type SignupLabels struct {
	Title            string `json:"title"`
	Heading          string `json:"heading"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Email            string `json:"email"`
	EmailPlaceholder string `json:"emailPlaceholder"`
	Password         string `json:"password"`
	ConfirmPassword  string `json:"confirmPassword"`
	Submit           string `json:"submit"`
	HasAccount       string `json:"hasAccount"`
	SignInLink       string `json:"signInLink"`
	TermsPrefix      string `json:"termsPrefix"`
	TermsLink        string `json:"termsLink"`
	PrivacyLink      string `json:"privacyLink"`
	AdminTitle       string `json:"adminTitle"`
	AdminDescription string `json:"adminDescription"`
	StaffTitle       string `json:"staffTitle"`
	StaffDescription string `json:"staffDescription"`
	PasswordStrength string `json:"passwordStrength"`
}

// Signup02Labels holds i18n strings for the signup02 page (split-screen style).
type Signup02Labels struct {
	Title                      string `json:"title"`
	Heading                    string `json:"heading"`
	Subheading                 string `json:"subheading"`
	FirstNameLabel             string `json:"firstNameLabel"`
	FirstNamePlaceholder       string `json:"firstNamePlaceholder"`
	LastNameLabel              string `json:"lastNameLabel"`
	LastNamePlaceholder        string `json:"lastNamePlaceholder"`
	EmailLabel                 string `json:"emailLabel"`
	EmailPlaceholder           string `json:"emailPlaceholder"`
	PasswordLabel              string `json:"passwordLabel"`
	PasswordPlaceholder        string `json:"passwordPlaceholder"`
	ConfirmPasswordLabel       string `json:"confirmPasswordLabel"`
	ConfirmPasswordPlaceholder string `json:"confirmPasswordPlaceholder"`
	SignUpButton               string `json:"signUpButton"`
	HasAccount                 string `json:"hasAccount"`
	SignInLink                 string `json:"signInLink"`
	SocialDivider              string `json:"socialDivider"`
	TermsText                  string `json:"termsText"`
	Error                      string `json:"error"`
	// Carousel navigation + accessibility
	PreviousSlide    string `json:"previousSlide"`
	NextSlide        string `json:"nextSlide"`
	ContinueWith     string `json:"continueWith"`
	PasswordStrength string `json:"passwordStrength"`
	TermsLink        string `json:"termsLink"`
}

// ---------------------------------------------------------------------------
// Reset password labels
// ---------------------------------------------------------------------------

// ResetPasswordLabels holds i18n strings for the reset-password01 page (dual-card style).
type ResetPasswordLabels struct {
	Title              string `json:"title"`
	Heading            string `json:"heading"`
	Description        string `json:"description"`
	Email              string `json:"email"`
	EmailPlaceholder   string `json:"emailPlaceholder"`
	Submit             string `json:"submit"`
	BackToLogin        string `json:"backToLogin"`
	ConfirmHeading     string `json:"confirmHeading"`
	ConfirmDescription string `json:"confirmDescription"`
	NewPassword        string `json:"newPassword"`
	ConfirmPassword    string `json:"confirmPassword"`
	ResetButton        string `json:"resetButton"`
	SuccessHeading     string `json:"successHeading"`
	SuccessMessage     string `json:"successMessage"`
}

// ResetPassword02Labels holds i18n strings for the reset-password02 page (split-screen style).
type ResetPassword02Labels struct {
	Title                      string `json:"title"`
	Heading                    string `json:"heading"`
	Subheading                 string `json:"subheading"`
	EmailLabel                 string `json:"emailLabel"`
	EmailPlaceholder           string `json:"emailPlaceholder"`
	SendResetButton            string `json:"sendResetButton"`
	BackToLogin                string `json:"backToLogin"`
	ConfirmHeading             string `json:"confirmHeading"`
	ConfirmSubheading          string `json:"confirmSubheading"`
	NewPasswordLabel           string `json:"newPasswordLabel"`
	NewPasswordPlaceholder     string `json:"newPasswordPlaceholder"`
	ConfirmPasswordLabel       string `json:"confirmPasswordLabel"`
	ConfirmPasswordPlaceholder string `json:"confirmPasswordPlaceholder"`
	ResetButton                string `json:"resetButton"`
	SuccessHeading             string `json:"successHeading"`
	SuccessMessage             string `json:"successMessage"`
	// Generic + code-specific error messages.
	// Action handlers emit short codes via the `?error=` query param; the
	// page handler maps each code to one of these fields. Never display
	// raw err.Error() — it's not localisable and may leak internals.
	//   ?error=mismatch       → ErrorMismatch
	//   ?error=invalid_token  → ErrorInvalidToken
	//   ?error=expired_token  → ErrorExpiredToken
	//   ?error=weak_password  → ErrorWeakPassword
	//   ?error=generic (and anything unrecognized) → Error
	Error             string `json:"error"`
	ErrorMismatch     string `json:"errorMismatch"`
	ErrorInvalidToken string `json:"errorInvalidToken"`
	ErrorExpiredToken string `json:"errorExpiredToken"`
	ErrorWeakPassword string `json:"errorWeakPassword"`
	// Carousel navigation
	PreviousSlide string `json:"previousSlide"`
	NextSlide     string `json:"nextSlide"`
}

// ChangePasswordLabels holds i18n strings for the change-password page.
//
// Error fields are addressed by code: the action handler emits a short
// code via `?error=...`, and the page handler maps each code to one of
// these fields. Raw err.Error() must never be rendered.
//
//	?error=mismatch  → ErrorMismatch
//	?error=incorrect → ErrorCurrentIncorrect
//	?error=too_short → ErrorTooShort
//	?error=generic (and anything unrecognized) → Error
type ChangePasswordLabels struct {
	Title                      string `json:"title"`
	Heading                    string `json:"heading"`
	Subheading                 string `json:"subheading"`
	OldPasswordLabel           string `json:"oldPasswordLabel"`
	OldPasswordPlaceholder     string `json:"oldPasswordPlaceholder"`
	NewPasswordLabel           string `json:"newPasswordLabel"`
	NewPasswordPlaceholder     string `json:"newPasswordPlaceholder"`
	ConfirmPasswordLabel       string `json:"confirmPasswordLabel"`
	ConfirmPasswordPlaceholder string `json:"confirmPasswordPlaceholder"`
	SubmitButton               string `json:"submitButton"`
	SuccessMessage             string `json:"successMessage"`
	// Generic fallback + code-specific error messages.
	Error                 string `json:"error"`
	ErrorMismatch         string `json:"errorMismatch"`
	ErrorCurrentIncorrect string `json:"errorCurrentIncorrect"`
	ErrorTooShort         string `json:"errorTooShort"`
	BackToApp             string `json:"backToApp"`
}

// ---------------------------------------------------------------------------
// Auth email labels
// ---------------------------------------------------------------------------

// AuthEmailLabels holds i18n strings for authentication-related email templates.
type AuthEmailLabels struct {
	ResetSubject           string `json:"resetSubject"`
	ResetHeading           string `json:"resetHeading"`
	ResetBody              string `json:"resetBody"`
	ResetButtonText        string `json:"resetButtonText"`
	ResetExpiry            string `json:"resetExpiry"`
	WelcomeSubject         string `json:"welcomeSubject"`
	WelcomeHeading         string `json:"welcomeHeading"`
	WelcomeBody            string `json:"welcomeBody"`
	WelcomeButtonText      string `json:"welcomeButtonText"`
	PasswordChangedSubject string `json:"passwordChangedSubject"`
	PasswordChangedHeading string `json:"passwordChangedHeading"`
	PasswordChangedBody    string `json:"passwordChangedBody"`
	SecurityNotice         string `json:"securityNotice"`
}

// ---------------------------------------------------------------------------
// Supplier labels
// ---------------------------------------------------------------------------

// SupplierLabels holds all translatable strings for the supplier module.
type SupplierLabels struct {
	Page    SupplierPageLabels   `json:"page"`
	Buttons SupplierButtonLabels `json:"buttons"`
	Columns SupplierColumnLabels `json:"columns"`
	Empty   SupplierEmptyLabels  `json:"empty"`
	Form    SupplierFormLabels   `json:"form"`
	Detail  SupplierDetailLabels `json:"detail"`
	Actions SupplierActionLabels `json:"actions"`
}

type SupplierPageLabels struct {
	Heading        string `json:"heading"`
	HeadingActive  string `json:"headingActive"`
	HeadingBlocked string `json:"headingBlocked"`
	HeadingOnHold  string `json:"headingOnHold"`
	Caption        string `json:"caption"`
	CaptionActive  string `json:"captionActive"`
	CaptionBlocked string `json:"captionBlocked"`
	CaptionOnHold  string `json:"captionOnHold"`
}

type SupplierButtonLabels struct {
	AddNew string `json:"addNew"`
}

type SupplierColumnLabels struct {
	Name         string `json:"name"`
	SupplierType string `json:"supplierType"`
	InternalID   string `json:"internalId"`
	Status       string `json:"status"`
	Category     string `json:"category"`
	PaymentTerms string `json:"paymentTerms"`
	ContactName  string `json:"contactName"`
	DateCreated  string `json:"dateCreated"`
}

type SupplierEmptyLabels struct {
	ActiveTitle    string `json:"activeTitle"`
	ActiveMessage  string `json:"activeMessage"`
	BlockedTitle   string `json:"blockedTitle"`
	BlockedMessage string `json:"blockedMessage"`
	OnHoldTitle    string `json:"onHoldTitle"`
	OnHoldMessage  string `json:"onHoldMessage"`
}

type SupplierFormLabels struct {
	Name               string `json:"name"`
	SupplierType       string `json:"supplierType"`
	TaxID              string `json:"taxId"`
	RegistrationNumber string `json:"registrationNumber"`
	StreetAddress      string `json:"streetAddress"`
	City               string `json:"city"`
	Province           string `json:"province"`
	PostalCode         string `json:"postalCode"`
	Country            string `json:"country"`
	BillingCurrency    string `json:"billingCurrency"`
	PaymentTerms       string `json:"paymentTerms"`
	LeadTimeDays       string `json:"leadTimeDays"`
	CreditLimit        string `json:"creditLimit"`
	Status             string `json:"status"`
	Website            string `json:"website"`
	Notes              string `json:"notes"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Active             string `json:"active"`

	// Section titles
	SectionCompany        string `json:"sectionCompany"`
	SectionRepresentative string `json:"sectionRepresentative"`
	SectionAccounting     string `json:"sectionAccounting"`
	SectionAddress        string `json:"sectionAddress"`
	SectionOthers         string `json:"sectionOthers"`

	// Timezone autocomplete
	Timezone                  string `json:"timezone"`
	TimezonePlaceholder       string `json:"timezonePlaceholder"`
	TimezoneSearchPlaceholder string `json:"timezoneSearchPlaceholder"`
	TimezoneNoResults         string `json:"timezoneNoResults"`
	TimezoneInfo              string `json:"timezoneInfo"`

	// Placeholders
	NamePlaceholder               string `json:"namePlaceholder"`
	SupplierTypePlaceholder       string `json:"supplierTypePlaceholder"`
	StatusPlaceholder             string `json:"statusPlaceholder"`
	FirstNamePlaceholder          string `json:"firstNamePlaceholder"`
	LastNamePlaceholder           string `json:"lastNamePlaceholder"`
	EmailPlaceholder              string `json:"emailPlaceholder"`
	PhonePlaceholder              string `json:"phonePlaceholder"`
	PaymentTermsPlaceholder       string `json:"paymentTermsPlaceholder"`
	CreditLimitPlaceholder        string `json:"creditLimitPlaceholder"`
	BillingCurrencyPlaceholder    string `json:"billingCurrencyPlaceholder"`
	LeadTimeDaysPlaceholder       string `json:"leadTimeDaysPlaceholder"`
	TaxIDPlaceholder              string `json:"taxIdPlaceholder"`
	RegistrationNumberPlaceholder string `json:"registrationNumberPlaceholder"`
	StreetAddressPlaceholder      string `json:"streetAddressPlaceholder"`
	CityPlaceholder               string `json:"cityPlaceholder"`
	ProvincePlaceholder           string `json:"provincePlaceholder"`
	PostalCodePlaceholder         string `json:"postalCodePlaceholder"`
	CountryPlaceholder            string `json:"countryPlaceholder"`
	WebsitePlaceholder            string `json:"websitePlaceholder"`
	NotesPlaceholder              string `json:"notesPlaceholder"`

	// Select option labels
	TypeCompany    string `json:"typeCompany"`
	TypeIndividual string `json:"typeIndividual"`

	StatusActive  string `json:"statusActive"`
	StatusBlocked string `json:"statusBlocked"`
	StatusOnHold  string `json:"statusOnHold"`

	TermsImmediate string `json:"termsImmediate"`
	TermsNet30     string `json:"termsNet30"`
	TermsNet60     string `json:"termsNet60"`
	Terms2_10Net30 string `json:"terms2_10Net30"`
}

type SupplierDetailLabels struct {
	InfoTab       string                      `json:"infoTab"`
	CompanyInfo   SupplierDetailSectionLabels `json:"companyInfo"`
	ContactInfo   SupplierDetailSectionLabels `json:"contactInfo"`
	FinancialInfo SupplierDetailSectionLabels `json:"financialInfo"`
	AddressInfo   SupplierDetailSectionLabels `json:"addressInfo"`
	// Tab label for attachments
	AttachmentsTab string `json:"attachmentsTab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"auditHistoryTab"`
	// Tab label for statement
	StatementTab string `json:"statementTab"`
	// Inline labels
	DaysSuffix string `json:"daysSuffix"`
	Website    string `json:"website"`
	// Purchase Orders tab labels
	PurchaseOrders SupplierPurchaseOrdersLabels `json:"purchaseOrders"`
	// Statement tab stat card labels
	OutstandingBalance string `json:"outstandingBalance"`
	TotalBilled        string `json:"totalBilled"`
	TotalPaid          string `json:"totalPaid"`
	Bills              string `json:"bills"`
	// Statement empty state
	EmptyStatementTitle   string `json:"emptyStatementTitle"`
	EmptyStatementMessage string `json:"emptyStatementMessage"`
}

// SupplierDetailSectionLabels holds a title for a detail page section.
type SupplierDetailSectionLabels struct {
	Title string `json:"title"`
}

// SupplierPurchaseOrdersLabels holds labels for the purchase orders tab on the supplier detail page.
type SupplierPurchaseOrdersLabels struct {
	Title        string `json:"title"`
	ColPONumber  string `json:"colPONumber"`
	ColOrderDate string `json:"colOrderDate"`
	ColAmount    string `json:"colAmount"`
	ColCurrency  string `json:"colCurrency"`
	ColStatus    string `json:"colStatus"`
	EmptyPO      string `json:"emptyPO"`
}

type SupplierActionLabels struct {
	View      string `json:"view"`
	Edit      string `json:"edit"`
	Delete    string `json:"delete"`
	Activate  string `json:"activate"`
	Block     string `json:"block"`
	SetOnHold string `json:"setOnHold"`
}

// ---------------------------------------------------------------------------
// Client Tag labels
// ---------------------------------------------------------------------------

// ClientTagLabels holds all translatable strings for the client tag module.
type ClientTagLabels struct {
	Page    ClientTagPageLabels    `json:"page"`
	Buttons ClientTagButtonLabels  `json:"buttons"`
	Columns ClientTagColumnLabels  `json:"columns"`
	Empty   ClientTagEmptyLabels   `json:"empty"`
	Actions ClientTagActionLabels  `json:"actions"`
	Confirm ClientTagConfirmLabels `json:"confirm"`
}

type ClientTagPageLabels struct {
	Heading  string `json:"heading"`
	Subtitle string `json:"subtitle"`
}

type ClientTagButtonLabels struct {
	AddTag string `json:"addTag"`
}

type ClientTagColumnLabels struct {
	TagName     string `json:"tagName"`
	Customers   string `json:"customers"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ClientTagEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ClientTagActionLabels struct {
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ClientTagConfirmLabels struct {
	DeleteTitle   string `json:"deleteTitle"`
	DeleteMessage string `json:"deleteMessage"`
	CannotDelete  string `json:"cannotDelete"`
}

// ---------------------------------------------------------------------------
// Supplier Tag labels
// ---------------------------------------------------------------------------

// SupplierTagLabels holds all translatable strings for the supplier tag module.
type SupplierTagLabels struct {
	Page    SupplierTagPageLabels    `json:"page"`
	Buttons SupplierTagButtonLabels  `json:"buttons"`
	Columns SupplierTagColumnLabels  `json:"columns"`
	Empty   SupplierTagEmptyLabels   `json:"empty"`
	Actions SupplierTagActionLabels  `json:"actions"`
	Confirm SupplierTagConfirmLabels `json:"confirm"`
}

type SupplierTagPageLabels struct {
	Heading  string `json:"heading"`
	Subtitle string `json:"subtitle"`
}

type SupplierTagButtonLabels struct {
	AddTag string `json:"addTag"`
}

type SupplierTagColumnLabels struct {
	TagName     string `json:"tagName"`
	Suppliers   string `json:"suppliers"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type SupplierTagEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SupplierTagActionLabels struct {
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SupplierTagConfirmLabels struct {
	DeleteTitle   string `json:"deleteTitle"`
	DeleteMessage string `json:"deleteMessage"`
	CannotDelete  string `json:"cannotDelete"`
}

// ---------------------------------------------------------------------------
// PaymentTerm labels
// ---------------------------------------------------------------------------

// PaymentTermLabels holds all translatable strings for the payment term module.
type PaymentTermLabels struct {
	Page    PaymentTermPageLabels   `json:"page"`
	Buttons PaymentTermButtonLabels `json:"buttons"`
	Columns PaymentTermColumnLabels `json:"columns"`
	Empty   PaymentTermEmptyLabels  `json:"empty"`
	Form    PaymentTermFormLabels   `json:"form"`
	Actions PaymentTermActionLabels `json:"actions"`
}

type PaymentTermPageLabels struct {
	Heading  string `json:"heading"`
	Subtitle string `json:"subtitle"`
}

type PaymentTermButtonLabels struct {
	AddPaymentTerm string `json:"addPaymentTerm"`
}

type PaymentTermColumnLabels struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	Type      string `json:"type"`
	NetDays   string `json:"netDays"`
	IsDefault string `json:"isDefault"`
	Status    string `json:"status"`
}

type PaymentTermEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type PaymentTermFormLabels struct {
	SectionInfo            string `json:"sectionInfo"`
	SectionTerms           string `json:"sectionTerms"`
	SectionSettings        string `json:"sectionSettings"`
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Code                   string `json:"code"`
	CodePlaceholder        string `json:"codePlaceholder"`
	Type                   string `json:"type"`
	NetDays                string `json:"netDays"`
	DiscountDays           string `json:"discountDays"`
	DiscountPercentBps     string `json:"discountPercentBps"`
	TypeHint               string `json:"typeHint"`
	NetDaysHint            string `json:"netDaysHint"`
	DiscountDaysHint       string `json:"discountDaysHint"`
	DiscountPercentBpsHint string `json:"discountPercentBpsHint"`
	PriorityHint           string `json:"priorityHint"`
	EntityScope            string `json:"entityScope"`
	IsDefault              string `json:"isDefault"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	DisplayOrder           string `json:"displayOrder"`
	Active                 string `json:"active"`

	// Type select options
	TypeDueOnReceipt        string `json:"typeDueOnReceipt"`
	TypeNet                 string `json:"typeNet"`
	TypeCOD                 string `json:"typeCOD"`
	TypeProximate           string `json:"typeProximate"`
	ProximateDay            string `json:"proximateDay"`
	ProximateDayPlaceholder string `json:"proximateDayPlaceholder"`

	// Entity scope select options
	ScopesBoth         string `json:"scopesBoth"`
	ScopesSupplierOnly string `json:"scopesSupplierOnly"`
	ScopesClientOnly   string `json:"scopesClientOnly"`

	// Field-level info text for the drawer form.
	NameInfo        string `json:"nameInfo"`
	CodeInfo        string `json:"codeInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
}

type PaymentTermActionLabels struct {
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

// ---------------------------------------------------------------------------
// Shared labels (used across all modules)
// ---------------------------------------------------------------------------

// SharedLabels holds translatable strings shared across all entydad modules.
type SharedLabels struct {
	Errors  SharedErrorLabels   `json:"errors"`
	Confirm SharedConfirmLabels `json:"confirm"`
	Badges  SharedBadgeLabels   `json:"badges"`
}

// SharedErrorLabels holds HTMXError messages used across all action handlers.
type SharedErrorLabels struct {
	PermissionDenied    string `json:"permissionDenied"`
	InvalidFormData     string `json:"invalidFormData"`
	InvalidStatus       string `json:"invalidStatus"`
	InvalidTargetStatus string `json:"invalidTargetStatus"`
	NotFound            string `json:"notFound"`
	IDRequired          string `json:"idRequired"`
	NoIDsProvided       string `json:"noIdsProvided"`
	PasswordRequired    string `json:"passwordRequired"`
	PasswordFailed      string `json:"passwordFailed"`
	RoleRequired        string `json:"roleRequired"`
	PermissionRequired  string `json:"permissionRequired"`
	UserRequired        string `json:"userRequired"`
	TagNotFound         string `json:"tagNotFound"`
	TagNameExists       string `json:"tagNameExists"`
	VerifyFailed        string `json:"verifyFailed"`
	CannotDeleteInUse   string `json:"cannotDeleteInUse"`
}

// SharedConfirmLabels holds confirm dialog message templates used across modules.
type SharedConfirmLabels struct {
	Activate       string `json:"activate"`
	Deactivate     string `json:"deactivate"`
	Delete         string `json:"delete"`
	Block          string `json:"block"`
	Hold           string `json:"hold"`
	Prospect       string `json:"prospect"`
	Remove         string `json:"remove"`
	BulkActivate   string `json:"bulkActivate"`
	BulkDeactivate string `json:"bulkDeactivate"`
	BulkDelete     string `json:"bulkDelete"`
	BulkBlock      string `json:"bulkBlock"`
	BulkHold       string `json:"bulkHold"`
	BulkProspect   string `json:"bulkProspect"`
}

// SharedBadgeLabels holds translatable badge values.
type SharedBadgeLabels struct {
	Allow        string `json:"allow"`
	Deny         string `json:"deny"`
	Yes          string `json:"yes"`
	No           string `json:"no"`
	NoPermission string `json:"noPermission"`
}

// DashboardLabels holds translatable strings for dashboard pages.
type DashboardLabels struct {
	ClientTitle   string `json:"clientTitle"`
	UserTitle     string `json:"userTitle"`
	SupplierTitle string `json:"supplierTitle"`
	LocationTitle string `json:"locationTitle"`
	AdminTitle    string `json:"adminTitle"`
}

// ClientDashboardLabels holds translatable strings for the client dashboard.
type ClientDashboardLabels struct {
	TotalClients   string `json:"totalClients"`
	Active         string `json:"active"`
	Inactive       string `json:"inactive"`
	NewThisMonth   string `json:"newThisMonth"`
	ClientGrowth   string `json:"clientGrowth"`
	FilterWeek     string `json:"filterWeek"`
	FilterMonth    string `json:"filterMonth"`
	FilterYear     string `json:"filterYear"`
	RecentActivity string `json:"recentActivity"`
	ViewAll        string `json:"viewAll"`

	// Quick action labels (Phase 1b — pyeza dashboard block refactor)
	QuickNew          string `json:"quickNew"`
	QuickViewAll      string `json:"quickViewAll"`
	QuickTags         string `json:"quickTags"`
	QuickPaymentTerms string `json:"quickPaymentTerms"`

	// Activity feed titles
	ClientAdded     string `json:"clientAdded"`
	ClientActivated string `json:"clientActivated"`
	ProfileUpdated  string `json:"profileUpdated"`
	TagAssigned     string `json:"tagAssigned"`
}

// SupplierDashboardLabels holds translatable strings for the supplier dashboard.
type SupplierDashboardLabels struct {
	TotalSuppliers   string `json:"totalSuppliers"`
	Active           string `json:"active"`
	Blocked          string `json:"blocked"`
	OnHold           string `json:"onHold"`
	SupplierActivity string `json:"supplierActivity"`
	TopSuppliers     string `json:"topSuppliers"`
	FilterWeek       string `json:"filterWeek"`
	FilterMonth      string `json:"filterMonth"`
	FilterYear       string `json:"filterYear"`
	RecentActivity   string `json:"recentActivity"`
	ViewAll          string `json:"viewAll"`

	// Quick action labels (Phase 1b — pyeza dashboard block refactor)
	QuickNew        string `json:"quickNew"`
	QuickViewAll    string `json:"quickViewAll"`
	QuickTags       string `json:"quickTags"`
	QuickCategories string `json:"quickCategories"`

	// Activity feed titles
	SupplierAdded     string `json:"supplierAdded"`
	SupplierActivated string `json:"supplierActivated"`
	DetailsUpdated    string `json:"detailsUpdated"`
	TagAssigned       string `json:"tagAssigned"`
}

// LocationDashboardLabels holds translatable strings for the location dashboard.
type LocationDashboardLabels struct {
	// Stats (4): Total / Active / Regions / Areas Count
	TotalLocations string `json:"totalLocations"`
	Active         string `json:"active"`
	Regions        string `json:"regions"`
	AreasCount     string `json:"areasCount"`

	// Widget titles
	LocationsByRegion  string `json:"locationsByRegion"`
	TopLocationsByArea string `json:"topLocationsByArea"`
	RecentAdditions    string `json:"recentAdditions"`
	ViewAll            string `json:"viewAll"`

	// Chart filter labels
	FilterWeek  string `json:"filterWeek"`
	FilterMonth string `json:"filterMonth"`
	FilterYear  string `json:"filterYear"`

	// Quick action labels
	QuickNewLocation string `json:"quickNewLocation"`
	QuickNewArea     string `json:"quickNewArea"`

	// Activity / table column labels
	ColumnLocation string `json:"columnLocation"`
	ColumnAreas    string `json:"columnAreas"`
	LocationAdded  string `json:"locationAdded"`
}

// AdminDashboardLabels holds translatable strings for the admin app dashboard.
//
// The admin app is composite: its dashboard surfaces aggregates across the
// permission, role, workspace, workspace_user, and workspace_user_role
// entities — see plan.md § Phase 4b.
type AdminDashboardLabels struct {
	// Page header / subtitle
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`

	// Stats (4): Workspace Users / Roles / Permissions / Recent Role Changes (7d)
	WorkspaceUsers    string `json:"workspaceUsers"`
	Roles             string `json:"roles"`
	Permissions       string `json:"permissions"`
	RecentRoleChanges string `json:"recentRoleChanges"`

	// Widget titles
	UsersPerRole           string `json:"usersPerRole"`
	RolesByPermissionCount string `json:"rolesByPermissionCount"`
	RecentRoleChangesList  string `json:"recentRoleChangesList"`
	ViewAll                string `json:"viewAll"`

	// Quick action labels
	QuickNewUser      string `json:"quickNewUser"`
	QuickNewWorkspace string `json:"quickNewWorkspace"`
	QuickAssignRole   string `json:"quickAssignRole"`
	QuickAuditLog     string `json:"quickAuditLog"`

	// Activity / table column labels
	ColumnRole            string `json:"columnRole"`
	ColumnPermissionCount string `json:"columnPermissionCount"`
	RoleAssigned          string `json:"roleAssigned"`
}

// UserDashboardLabels holds translatable strings for the user dashboard.
type UserDashboardLabels struct {
	TotalUsers       string `json:"totalUsers"`
	Active           string `json:"active"`
	Inactive         string `json:"inactive"`
	Roles            string `json:"roles"`
	UserActivity     string `json:"userActivity"`
	FilterWeek       string `json:"filterWeek"`
	FilterMonth      string `json:"filterMonth"`
	FilterYear       string `json:"filterYear"`
	RecentActivity   string `json:"recentActivity"`
	ViewAll          string `json:"viewAll"`
	NoRecentActivity string `json:"noRecentActivity"`

	// Quick action labels (Phase 1b — pyeza dashboard block refactor)
	QuickNew         string `json:"quickNew"`
	QuickViewAll     string `json:"quickViewAll"`
	QuickRoles       string `json:"quickRoles"`
	QuickPermissions string `json:"quickPermissions"`

	// Activity feed titles
	UserAdded      string `json:"userAdded"`
	UserActivated  string `json:"userActivated"`
	RoleAssigned   string `json:"roleAssigned"`
	ProfileUpdated string `json:"profileUpdated"`
}

// ---------------------------------------------------------------------------
// Shared types
// ---------------------------------------------------------------------------

// RoleBadge holds minimal role info for display as a chip/badge in lists.
type RoleBadge struct {
	Name  string
	Color string
}

// ---------------------------------------------------------------------------
// TaxRegistrationLabels
// Lyngua root key: "taxRegistration"
// ---------------------------------------------------------------------------

// TaxRegistrationLabels holds all translatable strings for the polymorphic
// Tax Registration views (client + workspace party types in v1).
type TaxRegistrationLabels struct {
	Page    TaxRegistrationPageLabels   `json:"page"`
	Columns TaxRegistrationColumnLabels `json:"columns"`
	Buttons TaxRegistrationButtonLabels `json:"buttons"`
	Actions TaxRegistrationActionLabels `json:"actions"`
	Empty   TaxRegistrationEmptyLabels  `json:"empty"`
	Fields  TaxRegistrationFieldLabels  `json:"fields"`
	Revoke  TaxRegistrationRevokeLabels `json:"revoke"`
}

// TaxRegistrationPageLabels holds page heading strings.
type TaxRegistrationPageLabels struct {
	Heading          string `json:"heading"`
	HeadingClient    string `json:"headingClient"`
	HeadingWorkspace string `json:"headingWorkspace"`
	Caption          string `json:"caption"`
	AddDrawerTitle   string `json:"addDrawerTitle"`
	EditDrawerTitle  string `json:"editDrawerTitle"`
}

// TaxRegistrationColumnLabels holds table column headers.
type TaxRegistrationColumnLabels struct {
	KindName           string `json:"kindName"`
	ComputePath        string `json:"computePath"`
	PartyRole          string `json:"partyRole"`
	Status             string `json:"status"`
	EffectiveFrom      string `json:"effectiveFrom"`
	RegistrationNumber string `json:"registrationNumber"`
}

// TaxRegistrationButtonLabels holds button text.
type TaxRegistrationButtonLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

// TaxRegistrationActionLabels holds action dropdown labels.
type TaxRegistrationActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	NoPermission string `json:"noPermission"`
}

// TaxRegistrationEmptyLabels holds empty-state strings.
type TaxRegistrationEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// TaxRegistrationFieldLabels holds drawer form field labels.
type TaxRegistrationFieldLabels struct {
	TaxRegistrationKindID string `json:"taxRegistrationKindId"`
	RegistrationNumber    string `json:"registrationNumber"`
	EffectiveFrom         string `json:"effectiveFrom"`
	Notes                 string `json:"notes"`
	Status                string `json:"status"`
}

// TaxRegistrationRevokeLabels holds strings for the revoke confirm dialog.
type TaxRegistrationRevokeLabels struct {
	WarningMessage        string `json:"warningMessage"`
	EffectiveTo           string `json:"effectiveTo"`
	AffectedPeriodsNotice string `json:"affectedPeriodsNotice"`
	// AffectedPeriodsCount is the row label for the pending-period count (Phase 5 M3).
	AffectedPeriodsCount string `json:"affectedPeriodsCount"`
	// AffectedSubscriptionsCount is the row label for the subscription count (Phase 5 M3).
	AffectedSubscriptionsCount string `json:"affectedSubscriptionsCount"`
	ReasonLabel                string `json:"reasonLabel"`
	ReasonPlaceholder          string `json:"reasonPlaceholder"`
	ConfirmButton              string `json:"confirmButton"`
}

// DefaultTaxRegistrationLabels returns TaxRegistrationLabels with sensible
// English defaults.
func DefaultTaxRegistrationLabels() TaxRegistrationLabels {
	return TaxRegistrationLabels{
		Page: TaxRegistrationPageLabels{
			Heading:          "Tax Registrations",
			HeadingClient:    "Client Tax Registrations",
			HeadingWorkspace: "Workspace Tax Registrations",
			Caption:          "Active tax registrations determine compute path during revenue recognition",
			AddDrawerTitle:   "Add Tax Registration",
			EditDrawerTitle:  "Edit Tax Registration",
		},
		Columns: TaxRegistrationColumnLabels{
			KindName:           "Kind",
			ComputePath:        "Compute Path",
			PartyRole:          "Party Role",
			Status:             "Status",
			EffectiveFrom:      "Effective From",
			RegistrationNumber: "Registration No.",
		},
		Buttons: TaxRegistrationButtonLabels{
			Add:    "Add Tax Registration",
			Edit:   "Edit",
			Delete: "Delete",
		},
		Actions: TaxRegistrationActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			NoPermission: "You do not have permission to manage tax registrations",
		},
		Empty: TaxRegistrationEmptyLabels{
			Title:   "No tax registrations",
			Message: "Add a tax registration to enable tax computation for this party.",
		},
		Fields: TaxRegistrationFieldLabels{
			TaxRegistrationKindID: "Tax Registration Kind",
			RegistrationNumber:    "Registration Number",
			EffectiveFrom:         "Effective From",
			Notes:                 "Notes",
			Status:                "Status",
		},
		Revoke: TaxRegistrationRevokeLabels{
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

// ===========================================================================
// Conversation labels — secure messaging / ticketing (Plan-4, 2026-06-03)
//
// Loaded from translations/en/{tier}/conversation.json (root key "conversation")
// and conversation_post.json (root key "conversationPost") via LoadPathIfExists.
// All fields are nil-safe: DefaultConversationLabels() pre-populates English so a
// missing JSON file does not produce empty strings in the UI.
// ===========================================================================

// ConversationLabels is the top-level label struct for the conversation surface.
type ConversationLabels struct {
	List    ConversationListLabels    `json:"list"`
	Inbox   ConversationInboxLabels   `json:"inbox"`
	Thread  ConversationThreadLabels  `json:"thread"`
	Status  ConversationStatusLabels  `json:"status"`
	Actions ConversationActionLabels  `json:"actions"`
	Form    ConversationFormLabels    `json:"form"`
	Columns ConversationColumnLabels  `json:"columns"`
	Confirm ConversationConfirmLabels `json:"confirm"`
	Errors  ConversationErrorLabels   `json:"errors"`
}

// ConversationListLabels — staff inbox + portal thread-list headings.
type ConversationListLabels struct {
	Heading      string `json:"heading"`
	Subtitle     string `json:"subtitle"`
	Title        string `json:"title"`
	NewButton    string `json:"newButton"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// ConversationInboxLabels — staff filter chips.
type ConversationInboxLabels struct {
	FilterAll        string `json:"filterAll"`
	FilterUnassigned string `json:"filterUnassigned"`
	FilterMyQueue    string `json:"filterMyQueue"`
	FilterOpen       string `json:"filterOpen"`
	FilterInProgress string `json:"filterInProgress"`
	FilterResolved   string `json:"filterResolved"`
	FilterClosed     string `json:"filterClosed"`
}

// ConversationThreadLabels — thread-detail header + meta.
type ConversationThreadLabels struct {
	BackToInbox   string `json:"backToInbox"`
	Assignee      string `json:"assignee"`
	Unassigned    string `json:"unassigned"`
	Client        string `json:"client"`
	Created       string `json:"created"`
	LastActivity  string `json:"lastActivity"`
	ViewRequest   string `json:"viewRequest"`
	Subtitle      string `json:"subtitle"`
	EmptyTitle    string `json:"emptyTitle"`
	EmptySubtitle string `json:"emptySubtitle"`
}

// ConversationStatusLabels — human-readable status badge labels keyed by enum.
type ConversationStatusLabels struct {
	Open       string `json:"open"`
	InProgress string `json:"inProgress"`
	Resolved   string `json:"resolved"`
	Closed     string `json:"closed"`
	Unknown    string `json:"unknown"`
}

// ConversationActionLabels — action button labels.
type ConversationActionLabels struct {
	NewConversation string `json:"newConversation"`
	Open            string `json:"open"`
	Assign          string `json:"assign"`
	MarkResolved    string `json:"markResolved"`
	Close           string `json:"close"`
	Reopen          string `json:"reopen"`
	SetStatus       string `json:"setStatus"`
	Send            string `json:"send"`
}

// ConversationFormLabels — new-conversation / assign / status drawer fields.
type ConversationFormLabels struct {
	SectionTitle         string `json:"sectionTitle"`
	SubjectLabel         string `json:"subjectLabel"`
	SubjectPlaceholder   string `json:"subjectPlaceholder"`
	ClientLabel          string `json:"clientLabel"`
	ClientPlaceholder    string `json:"clientPlaceholder"`
	AssigneeLabel        string `json:"assigneeLabel"`
	AssigneePlaceholder  string `json:"assigneePlaceholder"`
	LinkLabel            string `json:"linkLabel"`
	LinkPlaceholder      string `json:"linkPlaceholder"`
	MessageLabel         string `json:"messageLabel"`
	MessagePlaceholder   string `json:"messagePlaceholder"`
	CurrentStatusLabel   string `json:"currentStatusLabel"`
	NewStatusLabel       string `json:"newStatusLabel"`
	CurrentAssigneeLabel string `json:"currentAssigneeLabel"`
}

// ConversationColumnLabels — staff inbox table headers.
type ConversationColumnLabels struct {
	Subject      string `json:"subject"`
	Client       string `json:"client"`
	LastActivity string `json:"lastActivity"`
	Assignee     string `json:"assignee"`
	Status       string `json:"status"`
}

// ConversationConfirmLabels — confirm-dialog copy for status transitions.
type ConversationConfirmLabels struct {
	ResolveTitle   string `json:"resolveTitle"`
	ResolveMessage string `json:"resolveMessage"`
	CloseTitle     string `json:"closeTitle"`
	CloseMessage   string `json:"closeMessage"`
	ReopenTitle    string `json:"reopenTitle"`
	ReopenMessage  string `json:"reopenMessage"`
}

// ConversationErrorLabels — error strings surfaced via HTMX error toast.
type ConversationErrorLabels struct {
	PermissionDenied  string `json:"permissionDenied"`
	NotFound          string `json:"notFound"`
	InvalidForm       string `json:"invalidForm"`
	SubjectRequired   string `json:"subjectRequired"`
	ClientRequired    string `json:"clientRequired"`
	MessageRequired   string `json:"messageRequired"`
	InvalidTransition string `json:"invalidTransition"`
	IDRequired        string `json:"idRequired"`
	SaveFailed        string `json:"saveFailed"`
}

// ConversationPostLabels is the label struct for the post composer + bubbles.
// Loaded from conversation_post.json (root key "conversationPost").
type ConversationPostLabels struct {
	Composer ConversationComposerLabels  `json:"composer"`
	Bubble   ConversationBubbleLabels    `json:"bubble"`
	Subtitle string                      `json:"subtitle"`
	Empty    string                      `json:"empty"`
	Errors   ConversationPostErrorLabels `json:"errors"`
}

// ConversationComposerLabels — reply composer.
type ConversationComposerLabels struct {
	Placeholder string `json:"placeholder"`
	Send        string `json:"send"`
	Attach      string `json:"attach"`
}

// ConversationBubbleLabels — sender role labels.
type ConversationBubbleLabels struct {
	You    string `json:"you"`
	Staff  string `json:"staff"`
	Client string `json:"client"`
}

// ConversationPostErrorLabels — post-specific errors.
type ConversationPostErrorLabels struct {
	EmptyBody    string `json:"emptyBody"`
	MissingToken string `json:"missingToken"`
	SendFailed   string `json:"sendFailed"`
}

// DefaultConversationLabels returns English defaults for the conversation
// surface. Override per business type via conversation.json.
func DefaultConversationLabels() ConversationLabels {
	return ConversationLabels{
		List: ConversationListLabels{
			Heading:      "Conversations",
			Subtitle:     "Secure messaging with your clients",
			Title:        "Messages",
			NewButton:    "New",
			EmptyTitle:   "No conversations yet",
			EmptyMessage: "Start a new conversation to message a client.",
		},
		Inbox: ConversationInboxLabels{
			FilterAll:        "All open",
			FilterUnassigned: "Unassigned",
			FilterMyQueue:    "My queue",
			FilterOpen:       "Open",
			FilterInProgress: "In progress",
			FilterResolved:   "Resolved",
			FilterClosed:     "Closed",
		},
		Thread: ConversationThreadLabels{
			BackToInbox:   "Back to inbox",
			Assignee:      "Assigned to",
			Unassigned:    "Unassigned",
			Client:        "Client",
			Created:       "Created",
			LastActivity:  "Last activity",
			ViewRequest:   "View request",
			Subtitle:      "Secure messaging. Every conversation is logged.",
			EmptyTitle:    "Select a conversation",
			EmptySubtitle: "Choose a thread from the list to view messages.",
		},
		Status: ConversationStatusLabels{
			Open:       "Open",
			InProgress: "In progress",
			Resolved:   "Resolved",
			Closed:     "Closed",
			Unknown:    "Unknown",
		},
		Actions: ConversationActionLabels{
			NewConversation: "New conversation",
			Open:            "Open",
			Assign:          "Assign",
			MarkResolved:    "Mark resolved",
			Close:           "Close",
			Reopen:          "Reopen",
			SetStatus:       "Change status",
			Send:            "Send",
		},
		Form: ConversationFormLabels{
			SectionTitle:         "Conversation details",
			SubjectLabel:         "Subject",
			SubjectPlaceholder:   "What is this about?",
			ClientLabel:          "Client",
			ClientPlaceholder:    "Search clients…",
			AssigneeLabel:        "Assign to",
			AssigneePlaceholder:  "Search staff…",
			LinkLabel:            "Linked request (optional)",
			LinkPlaceholder:      "e.g. REQ-0091",
			MessageLabel:         "Message",
			MessagePlaceholder:   "Type your first message…",
			CurrentStatusLabel:   "Current status",
			NewStatusLabel:       "New status",
			CurrentAssigneeLabel: "Currently",
		},
		Columns: ConversationColumnLabels{
			Subject:      "Conversation",
			Client:       "Client",
			LastActivity: "Last activity",
			Assignee:     "Assigned",
			Status:       "Status",
		},
		Confirm: ConversationConfirmLabels{
			ResolveTitle:   "Mark resolved",
			ResolveMessage: "Mark this conversation as resolved?",
			CloseTitle:     "Close conversation",
			CloseMessage:   "Close this conversation? It can be reopened later.",
			ReopenTitle:    "Reopen conversation",
			ReopenMessage:  "Reopen this conversation?",
		},
		Errors: ConversationErrorLabels{
			PermissionDenied:  "You don't have permission to perform this action.",
			NotFound:          "Conversation not found.",
			InvalidForm:       "Invalid form data.",
			SubjectRequired:   "A subject is required.",
			ClientRequired:    "Please select a client.",
			MessageRequired:   "A message is required.",
			InvalidTransition: "That status change is not allowed.",
			IDRequired:        "A conversation id is required.",
			SaveFailed:        "Could not save. Please try again.",
		},
	}
}

// DefaultConversationPostLabels returns English defaults for the composer /
// bubble surface. Override per business type via conversation_post.json.
func DefaultConversationPostLabels() ConversationPostLabels {
	return ConversationPostLabels{
		Composer: ConversationComposerLabels{
			Placeholder: "Reply…",
			Send:        "Send",
			Attach:      "Attach",
		},
		Bubble: ConversationBubbleLabels{
			You:    "You",
			Staff:  "Staff",
			Client: "Client",
		},
		Subtitle: "Secure messaging. Every conversation is logged.",
		Empty:    "No messages yet.",
		Errors: ConversationPostErrorLabels{
			EmptyBody:    "Message cannot be empty.",
			MissingToken: "Missing idempotency token. Please refresh and try again.",
			SendFailed:   "Could not send your message. Please try again.",
		},
	}
}

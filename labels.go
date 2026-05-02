package entydad

import (
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
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
	ClientName     string `json:"clientName"`
	Representative string `json:"representative"`
	Status         string `json:"status"`
	Category       string `json:"category"`
	PaymentTerm    string `json:"paymentTerm"`
	DateCreated    string `json:"dateCreated"`
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
	Status                string `json:"status"`
	StatusPlaceholder     string `json:"statusPlaceholder"`
	StatusActive          string `json:"statusActive"`
	StatusBlocked         string `json:"statusBlocked"`
	StatusOnHold          string `json:"statusOnHold"`
	StatusInactive        string `json:"statusInactive"`
	StatusProspect        string `json:"statusProspect"`
	Country               string `json:"country"`
	CountryPlaceholder    string `json:"countryPlaceholder"`
	Website               string `json:"website"`
	WebsitePlaceholder    string `json:"websitePlaceholder"`
	SectionCompany        string `json:"sectionCompany"`
	SectionAddress        string `json:"sectionAddress"`
	SectionRepresentative string `json:"sectionRepresentative"`
	SectionAccounting     string `json:"sectionAccounting"`
	SectionOthers         string `json:"sectionOthers"`
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
	// Packages tab
	Packages ClientPackagesLabels `json:"packages"`
	// Subscriptions tab column headers + confirm dialogs
	Subscriptions ClientSubscriptionLabels `json:"subscriptions"`
	// Statement tab column headers + totals row
	Statement ClientStatementLabels `json:"statement"`
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
	Info              string `json:"info"`
	Representative    string `json:"representative"`
	Subscriptions     string `json:"subscriptions"`
	SubscriptionsSlug string `json:"subscriptionsSlug"`
	Accounting        string `json:"accounting"`
	History           string `json:"history"`
	Statement         string `json:"statement"`
	Packages          string `json:"packages"`
	Attachments       string `json:"attachments"`
	AuditHistory      string `json:"auditHistory"`
}

// ClientPackagesLabels holds labels for the Packages tab on the client detail page.
type ClientPackagesLabels struct {
	Empty             string `json:"empty"`
	AddAction         string `json:"addAction"`
	ColumnName        string `json:"columnName"`
	ColumnRateCard    string `json:"columnRateCard"`
	ColumnEngagements string `json:"columnEngagements"`
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

// ResolveTabSlug returns the URL slug for a canonical tab key. The
// "subscriptions" tab can be re-slugged per tier (e.g. professional ships
// "engagements"); other tabs round-trip through as-is.
func (t ClientDetailTabLabels) ResolveTabSlug(canonical string) string {
	if canonical == "subscriptions" {
		if s := strings.TrimSpace(t.SubscriptionsSlug); s != "" {
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
	Page    LocationPageLabels   `json:"page"`
	Buttons LocationButtonLabels `json:"buttons"`
	Columns LocationColumnLabels `json:"columns"`
	Empty   LocationEmptyLabels  `json:"empty"`
	Form    LocationFormLabels   `json:"form"`
	Actions LocationActionLabels `json:"actions"`
	Detail  LocationDetailLabels `json:"detail"`
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
}

// WorkspaceDetailTabLabels holds the tab display names for the workspace detail page.
type WorkspaceDetailTabLabels struct {
	Info  string `json:"info"`
	Users string `json:"users"`
}

// WorkspaceDetailUserLabels holds i18n strings for the Users tab on the workspace detail page.
type WorkspaceDetailUserLabels struct {
	AddButton string `json:"addButton"`
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
}

type WorkspaceUserDetailTabLabels struct {
	Info  string `json:"info"`
	Roles string `json:"roles"`
}

type WorkspaceUserDetailRolesLabels struct {
	AssignButton string `json:"assignButton"`
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
	Error                      string `json:"error"`
	// Carousel navigation
	PreviousSlide string `json:"previousSlide"`
	NextSlide     string `json:"nextSlide"`
}

// ChangePasswordLabels holds i18n strings for the change-password page.
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
	ErrorCurrentIncorrect      string `json:"errorCurrentIncorrect"`
	ErrorTooShort              string `json:"errorTooShort"`
	BackToApp                  string `json:"backToApp"`
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
// Mapping helpers
// ---------------------------------------------------------------------------

// MapTableLabels maps common labels into the flat types.TableLabels structure.
func MapTableLabels(common pyeza.CommonLabels) types.TableLabels {
	return types.TableLabels{
		Search:                   common.Table.Search,
		SearchPlaceholder:        common.Table.SearchPlaceholder,
		Filters:                  common.Table.Filters,
		FilterConditions:         common.Table.FilterConditions,
		ClearAll:                 common.Table.ClearAll,
		AddCondition:             common.Table.AddCondition,
		Clear:                    common.Table.Clear,
		ApplyFilters:             common.Table.ApplyFilters,
		Sort:                     common.Table.Sort,
		Columns:                  common.Table.Columns,
		Export:                   common.Table.Export,
		DensityLabel:             common.Table.Density.Title,
		DensityDense:             common.Table.Density.Dense,
		DensityDefault:           common.Table.Density.Default,
		DensityComfortable:       common.Table.Density.Comfortable,
		DensityCompact:           common.Table.Density.Compact,
		EntriesPerPage:           common.Table.EntriesLabel,
		Show:                     common.Table.Show,
		Entries:                  common.Table.Entries,
		Showing:                  common.Table.Showing,
		To:                       common.Table.To,
		Of:                       common.Table.Of,
		EntriesLabel:             common.Table.EntriesLabel,
		SelectAll:                common.Table.SelectAll,
		BulkSelectAllPage:        common.Table.BulkSelectAllPage,
		BulkSelectAllAcrossPages: common.Table.BulkSelectAllAcrossPages,
		BulkClearSelection:       common.Table.BulkClearSelection,
		ColumnSortLockedHint:     common.Table.ColumnSortLockedHint,
		SortAscText:              common.Table.SortAscText,
		SortDescText:             common.Table.SortDescText,
		SortAscNumber:            common.Table.SortAscNumber,
		SortDescNumber:           common.Table.SortDescNumber,
		SortAscDate:              common.Table.SortAscDate,
		SortDescDate:             common.Table.SortDescDate,
		SortAscEnum:              common.Table.SortAscEnum,
		SortDescEnum:             common.Table.SortDescEnum,
		FilterOpContains:         common.Table.FilterOpContains,
		FilterOpEquals:           common.Table.FilterOpEquals,
		FilterOpStartsWith:       common.Table.FilterOpStartsWith,
		FilterOpEndsWith:         common.Table.FilterOpEndsWith,
		FilterOpNotEquals:        common.Table.FilterOpNotEquals,
		FilterOpBetween:          common.Table.FilterOpBetween,
		FilterOpEq:               common.Table.FilterOpEq,
		FilterOpNeq:              common.Table.FilterOpNeq,
		FilterOpGt:               common.Table.FilterOpGt,
		FilterOpGte:              common.Table.FilterOpGte,
		FilterOpLt:               common.Table.FilterOpLt,
		FilterOpLte:              common.Table.FilterOpLte,
		FilterOpOn:               common.Table.FilterOpOn,
		FilterOpBefore:           common.Table.FilterOpBefore,
		FilterOpAfter:            common.Table.FilterOpAfter,
		FilterOpIn:               common.Table.FilterOpIn,
		FilterOpNotIn:            common.Table.FilterOpNotIn,
		FilterPresetToday:        common.Table.FilterPresetToday,
		FilterPreset7d:           common.Table.FilterPreset7d,
		FilterPreset30d:          common.Table.FilterPreset30d,
		FilterPresetMonth:        common.Table.FilterPresetMonth,
		FilterPresetCustom:       common.Table.FilterPresetCustom,
		FilterAny:                common.Table.FilterAny,
		FilterYes:                common.Table.FilterYes,
		FilterNo:                 common.Table.FilterNo,
		FilterSearchPlaceholder:  common.Table.FilterSearchPlaceholder,
		FilterMinPlaceholder:     common.Table.FilterMinPlaceholder,
		FilterMaxPlaceholder:     common.Table.FilterMaxPlaceholder,
		Actions:                  common.Table.Actions,
		Prev:                     common.Pagination.Prev,
		Next:                     common.Pagination.Next,
	}
}

// MapBulkConfig returns a BulkActionsConfig with labels from common bulk labels.
func MapBulkConfig(common pyeza.CommonLabels) types.BulkActionsConfig {
	return types.BulkActionsConfig{
		Enabled:        true,
		SelectAllLabel: common.Bulk.SelectAll,
		SelectedLabel:  common.Bulk.Selected,
		CancelLabel:    common.Bulk.ClearSelection,
	}
}

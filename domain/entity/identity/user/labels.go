package user

// labels.go — User label structs, the user dashboard labels, and the
// user-attached user↔role junction label set that the user view module owns.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// identity sub-context). Pure structural move — no behaviour change; field
// names, json tags, and string literals are byte-identical. Entity-local
// rename: UserLabels -> Labels, User<Xxx>Labels -> <Xxx>Labels,
// UserDashboardLabels -> DashboardLabels, UserRoleLabels -> RoleLabels.

// Labels holds all translatable strings for the user module.
// JSON tags match retail/user.json (no wrapper key).
type Labels struct {
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Form    FormLabels   `json:"form"`
	Actions ActionLabels `json:"actions"`
	Detail  DetailLabels `json:"detail"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type ButtonLabels struct {
	AddUser string `json:"addUser"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Roles       string `json:"roles"`
	Workspaces  string `json:"workspaces"`
	DateCreated string `json:"dateCreated"`
	Status      string `json:"status"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type FormLabels struct {
	Mobile              string `json:"mobile"`
	Timezone            string `json:"timezone"`
	TimezonePlaceholder string `json:"timezonePlaceholder"`
	TimezoneInfo        string `json:"timezoneInfo"`
}

type ActionLabels struct {
	View        string `json:"view"`
	Edit        string `json:"edit"`
	Delete      string `json:"delete"`
	Activate    string `json:"activate"`
	Deactivate  string `json:"deactivate"`
	ManageRoles string `json:"manageRoles"`
}

// DetailLabels holds labels for the user detail page.
type DetailLabels struct {
	BasicInfo   DetailBasicInfoLabels  `json:"basicInfo"`
	Tabs        DetailTabLabels        `json:"tabs"`
	Security    DetailSecurityLabels   `json:"security"`
	EmptyStates DetailEmptyStateLabels `json:"emptyStates"`
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

// DetailSecurityLabels holds labels for the security tab.
type DetailSecurityLabels struct {
	Title           string `json:"title"`
	LastLogin       string `json:"lastLogin"`
	MfaStatus       string `json:"mfaStatus"`
	MfaEnabled      string `json:"mfaEnabled"`
	MfaDisabled     string `json:"mfaDisabled"`
	PasswordSection string `json:"passwordSection"`
	ResetPassword   string `json:"resetPassword"`
}

// DetailEmptyStateLabels holds empty-state labels for user detail tabs.
type DetailEmptyStateLabels struct {
	AuditTitle string `json:"auditTitle"`
	AuditDesc  string `json:"auditDesc"`
}

type DetailBasicInfoLabels struct {
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

type DetailTabLabels struct {
	Info       string `json:"info"`
	Roles      string `json:"roles"`
	Security   string `json:"security"`
	AuditTrail string `json:"auditTrail"`
}

// ---------------------------------------------------------------------------
// User dashboard labels
// ---------------------------------------------------------------------------

// DashboardLabels holds translatable strings for the user dashboard.
type DashboardLabels struct {
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
// User-Role labels
// ---------------------------------------------------------------------------

// RoleLabels holds all translatable strings for the user-role assignment view.
type RoleLabels struct {
	Page    RolePageLabels   `json:"page"`
	Buttons RoleButtonLabels `json:"buttons"`
	Columns RoleColumnLabels `json:"columns"`
	Empty   RoleEmptyLabels  `json:"empty"`
	Form    RoleFormLabels   `json:"form"`
	Actions RoleActionLabels `json:"actions"`
}

type RolePageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type RoleButtonLabels struct {
	AssignRole string `json:"assignRole"`
}

type RoleColumnLabels struct {
	RoleName     string `json:"roleName"`
	Description  string `json:"description"`
	Color        string `json:"color"`
	DateAssigned string `json:"dateAssigned"`
}

type RoleEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type RoleFormLabels struct {
	Role string `json:"role"`
}

type RoleActionLabels struct {
	Assign      string `json:"assign"`
	Remove      string `json:"remove"`
	ManageRoles string `json:"manageRoles"`
}

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
	HeadingActive   string `json:"heading_active"`
	HeadingInactive string `json:"heading_inactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"caption_active"`
	CaptionInactive string `json:"caption_inactive"`
}

type ButtonLabels struct {
	AddUser string `json:"add_user"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Roles       string `json:"roles"`
	Workspaces  string `json:"workspaces"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"active_title"`
	ActiveMessage   string `json:"active_message"`
	InactiveTitle   string `json:"inactive_title"`
	InactiveMessage string `json:"inactive_message"`
}

type FormLabels struct {
	Mobile              string `json:"mobile"`
	Timezone            string `json:"timezone"`
	TimezonePlaceholder string `json:"timezone_placeholder"`
	TimezoneInfo        string `json:"timezone_info"`
}

type ActionLabels struct {
	View        string `json:"view"`
	Edit        string `json:"edit"`
	Delete      string `json:"delete"`
	Activate    string `json:"activate"`
	Deactivate  string `json:"deactivate"`
	ManageRoles string `json:"manage_roles"`
}

// DetailLabels holds labels for the user detail page.
type DetailLabels struct {
	BasicInfo   DetailBasicInfoLabels  `json:"basic_info"`
	Tabs        DetailTabLabels        `json:"tabs"`
	Security    DetailSecurityLabels   `json:"security"`
	EmptyStates DetailEmptyStateLabels `json:"empty_states"`
	// Inline feedback and empty-state messages
	UpdateSuccess            string `json:"update_success"`
	UpdateError              string `json:"update_error"`
	NoRolesAssigned          string `json:"no_roles_assigned"`
	NoRolesDesc              string `json:"no_roles_desc"`
	NewPasswordPlaceholder   string `json:"new_password_placeholder"`
	TogglePasswordVisibility string `json:"toggle_password_visibility"`
	GeneratePassword         string `json:"generate_password"`
	PasswordUpdated          string `json:"password_updated"`
	PasswordFailed           string `json:"password_failed"`
	// Tab label for attachments (shared across all detail pages)
	AttachmentsTab string `json:"attachments_tab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"audit_history_tab"`
}

// DetailSecurityLabels holds labels for the security tab.
type DetailSecurityLabels struct {
	Title             string `json:"title"`
	LastLogin         string `json:"last_login"`
	MfaStatus         string `json:"mfa_status"`
	MfaEnabled        string `json:"mfa_enabled"`
	MfaDisabled       string `json:"mfa_disabled"`
	PasswordSection   string `json:"password_section"`
	ResetPassword     string `json:"reset_password"`
	AuthMethod        string `json:"auth_method"`
	ManagedByProvider string `json:"managed_by_provider"`
	ManageAccountLink string `json:"manage_account_link"`
}

// DetailEmptyStateLabels holds empty-state labels for user detail tabs.
type DetailEmptyStateLabels struct {
	AuditTitle string `json:"audit_title"`
	AuditDesc  string `json:"audit_desc"`
}

type DetailBasicInfoLabels struct {
	Title                string `json:"title"`
	FirstName            string `json:"first_name"`
	FirstNamePlaceholder string `json:"first_name_placeholder"`
	LastName             string `json:"last_name"`
	LastNamePlaceholder  string `json:"last_name_placeholder"`
	Email                string `json:"email"`
	EmailPlaceholder     string `json:"email_placeholder"`
	Username             string `json:"username"`
	Division             string `json:"division"`
	Status               string `json:"status"`
	UserType             string `json:"user_type"`
	Mobile               string `json:"mobile"`
	MobilePlaceholder    string `json:"mobile_placeholder"`
	Active               string `json:"active"`
	Save                 string `json:"save"`
}

type DetailTabLabels struct {
	Info       string `json:"info"`
	Roles      string `json:"roles"`
	Security   string `json:"security"`
	AuditTrail string `json:"audit_trail"`
}

// ---------------------------------------------------------------------------
// User dashboard labels
// ---------------------------------------------------------------------------

// DashboardLabels holds translatable strings for the user dashboard.
type DashboardLabels struct {
	TotalUsers       string `json:"total_users"`
	Active           string `json:"active"`
	Inactive         string `json:"inactive"`
	Roles            string `json:"roles"`
	UserActivity     string `json:"user_activity"`
	FilterWeek       string `json:"filter_week"`
	FilterMonth      string `json:"filter_month"`
	FilterYear       string `json:"filter_year"`
	RecentActivity   string `json:"recent_activity"`
	ViewAll          string `json:"view_all"`
	NoRecentActivity string `json:"no_recent_activity"`

	// Quick action labels (Phase 1b — pyeza dashboard block refactor)
	QuickNew         string `json:"quick_new"`
	QuickViewAll     string `json:"quick_view_all"`
	QuickRoles       string `json:"quick_roles"`
	QuickPermissions string `json:"quick_permissions"`

	// Activity feed titles
	UserAdded      string `json:"user_added"`
	UserActivated  string `json:"user_activated"`
	RoleAssigned   string `json:"role_assigned"`
	ProfileUpdated string `json:"profile_updated"`
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
	AssignRole string `json:"assign_role"`
}

type RoleColumnLabels struct {
	RoleName     string `json:"role_name"`
	Description  string `json:"description"`
	Color        string `json:"color"`
	DateAssigned string `json:"date_assigned"`
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
	ManageRoles string `json:"manage_roles"`
}

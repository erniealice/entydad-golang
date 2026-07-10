package role

// labels.go — Role label structs, plus the role-attached junction label sets
// (role↔permission and role↔user) that the role view module owns.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// identity sub-context). Pure structural move — no behaviour change; field
// names, json tags, and string literals are byte-identical. Entity-local
// rename: RoleLabels -> Labels, Role<Xxx>Labels -> <Xxx>Labels,
// RolePermissionLabels -> PermissionLabels, RoleUserLabels -> UserLabels.

// Labels holds all translatable strings for the role module.
// JSON tags match the "role" wrapper key in retail/role.json.
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
	AddRole string `json:"add_role"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Permissions string `json:"permissions"`
	Status      string `json:"status"`
	DateCreated string `json:"date_created"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"active_title"`
	ActiveMessage   string `json:"active_message"`
	InactiveTitle   string `json:"inactive_title"`
	InactiveMessage string `json:"inactive_message"`
}

type FormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"name_placeholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	Color                  string `json:"color"`
	ColorPlaceholder       string `json:"color_placeholder"`
	Active                 string `json:"active"`
}

type ActionLabels struct {
	View              string `json:"view"`
	Edit              string `json:"edit"`
	Delete            string `json:"delete"`
	Activate          string `json:"activate"`
	Deactivate        string `json:"deactivate"`
	ManagePermissions string `json:"manage_permissions"`
}

// DetailLabels holds labels for the role detail page.
type DetailLabels struct {
	Tabs DetailTabLabels  `json:"tabs"`
	Info DetailInfoLabels `json:"info"`
	// Empty-state labels for role detail tabs
	NoPermissionsAssigned string `json:"no_permissions_assigned"`
	NoPermissionsDesc     string `json:"no_permissions_desc"`
	NoUsersAssigned       string `json:"no_users_assigned"`
	NoUsersDesc           string `json:"no_users_desc"`
	// Tab label for attachments
	AttachmentsTab string `json:"attachments_tab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"audit_history_tab"`
}

type DetailTabLabels struct {
	Info        string `json:"info"`
	Permissions string `json:"permissions"`
	Users       string `json:"users"`
}

type DetailInfoLabels struct {
	Title       string `json:"title"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Status      string `json:"status"`
}

// ---------------------------------------------------------------------------
// Role-Permission labels
// ---------------------------------------------------------------------------

// PermissionLabels holds all translatable strings for the role-permission assignment view.
type PermissionLabels struct {
	Page    PermissionPageLabels   `json:"page"`
	Buttons PermissionButtonLabels `json:"buttons"`
	Columns PermissionColumnLabels `json:"columns"`
	Empty   PermissionEmptyLabels  `json:"empty"`
	Form    PermissionFormLabels   `json:"form"`
	Actions PermissionActionLabels `json:"actions"`
}

type PermissionPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type PermissionButtonLabels struct {
	AssignPermission string `json:"assign_permission"`
}

type PermissionColumnLabels struct {
	PermissionName string `json:"permission_name"`
	Code           string `json:"code"`
	Type           string `json:"type"`
	DateAssigned   string `json:"date_assigned"`
}

type PermissionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type PermissionFormLabels struct {
	Permission string `json:"permission"`
}

type PermissionActionLabels struct {
	Assign            string `json:"assign"`
	Remove            string `json:"remove"`
	ManagePermissions string `json:"manage_permissions"`
}

// ---------------------------------------------------------------------------
// Role-User labels (reverse of User-Role: managing users on a role)
// ---------------------------------------------------------------------------

// UserLabels holds all translatable strings for the role-user assignment view.
type UserLabels struct {
	Page    UserPageLabels   `json:"page"`
	Buttons UserButtonLabels `json:"buttons"`
	Columns UserColumnLabels `json:"columns"`
	Empty   UserEmptyLabels  `json:"empty"`
	Form    UserFormLabels   `json:"form"`
	Actions UserActionLabels `json:"actions"`
}

type UserPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type UserButtonLabels struct {
	AssignUser string `json:"assign_user"`
}

type UserColumnLabels struct {
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	DateAssigned string `json:"date_assigned"`
}

type UserEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type UserFormLabels struct {
	User   string `json:"user"`
	Assign string `json:"assign"`
}

type UserActionLabels struct {
	Assign string `json:"assign"`
	Remove string `json:"remove"`
}

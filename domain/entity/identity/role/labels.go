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
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type ButtonLabels struct {
	AddRole string `json:"addRole"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Permissions string `json:"permissions"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type FormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Color                  string `json:"color"`
	ColorPlaceholder       string `json:"colorPlaceholder"`
	Active                 string `json:"active"`
}

type ActionLabels struct {
	View              string `json:"view"`
	Edit              string `json:"edit"`
	Delete            string `json:"delete"`
	Activate          string `json:"activate"`
	Deactivate        string `json:"deactivate"`
	ManagePermissions string `json:"managePermissions"`
}

// DetailLabels holds labels for the role detail page.
type DetailLabels struct {
	Tabs DetailTabLabels  `json:"tabs"`
	Info DetailInfoLabels `json:"info"`
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
	AssignPermission string `json:"assignPermission"`
}

type PermissionColumnLabels struct {
	PermissionName string `json:"permissionName"`
	Code           string `json:"code"`
	Type           string `json:"type"`
	DateAssigned   string `json:"dateAssigned"`
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
	ManagePermissions string `json:"managePermissions"`
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
	AssignUser string `json:"assignUser"`
}

type UserColumnLabels struct {
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	DateAssigned string `json:"dateAssigned"`
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

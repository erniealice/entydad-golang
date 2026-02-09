package entydad

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// ---------------------------------------------------------------------------
// Client labels
// ---------------------------------------------------------------------------

// ClientLabels holds all translatable strings for the client module.
// JSON tags match the "client" wrapper key in retail/client.json.
type ClientLabels struct {
	Page        ClientPageLabels        `json:"page"`
	Buttons     ClientButtonLabels      `json:"buttons"`
	Columns     ClientColumnLabels      `json:"columns"`
	Empty       ClientEmptyLabels       `json:"empty"`
	Form        ClientFormLabels        `json:"form"`
	Detail      ClientDetailLabels      `json:"detail"`
	BulkActions ClientBulkActionLabels  `json:"bulkActions"`
}

type ClientPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingProspect string `json:"headingProspect"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionProspect string `json:"captionProspect"`
	CaptionInactive string `json:"captionInactive"`
}

type ClientButtonLabels struct {
	AddNew string `json:"addNew"`
}

type ClientColumnLabels struct {
	ClientName string `json:"clientName"`
}

type ClientEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	ProspectTitle   string `json:"prospectTitle"`
	ProspectMessage string `json:"prospectMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type ClientFormLabels struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ClientDetailLabels struct {
	CompanyDetails ClientCompanyDetailLabels `json:"companyDetails"`
	Actions        ClientDetailActionLabels  `json:"actions"`
}

type ClientCompanyDetailLabels struct {
	Status string `json:"status"`
}

type ClientDetailActionLabels struct {
	ViewClient   string `json:"viewClient"`
	EditClient   string `json:"editClient"`
	DeleteClient string `json:"deleteClient"`
}

type ClientBulkActionLabels struct {
	SetAsInactive string `json:"setAsInactive"`
}

// ---------------------------------------------------------------------------
// User labels
// ---------------------------------------------------------------------------

// UserLabels holds all translatable strings for the user module.
// JSON tags match retail/user.json (no wrapper key).
type UserLabels struct {
	Page    UserPageLabels    `json:"page"`
	Buttons UserButtonLabels  `json:"buttons"`
	Columns UserColumnLabels  `json:"columns"`
	Empty   UserEmptyLabels   `json:"empty"`
	Form    UserFormLabels    `json:"form"`
	Actions UserActionLabels  `json:"actions"`
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
	Name   string `json:"name"`
	Email  string `json:"email"`
	Roles  string `json:"roles"`
	Status string `json:"status"`
}

type UserEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type UserFormLabels struct {
	Mobile string `json:"mobile"`
}

type UserActionLabels struct {
	View        string `json:"view"`
	Edit        string `json:"edit"`
	Delete      string `json:"delete"`
	Activate    string `json:"activate"`
	Deactivate  string `json:"deactivate"`
	ManageRoles string `json:"manageRoles"`
}

// ---------------------------------------------------------------------------
// Location labels
// ---------------------------------------------------------------------------

// LocationLabels holds all translatable strings for the location module.
// JSON tags match the "location" wrapper key in retail/location.json.
type LocationLabels struct {
	Page    LocationPageLabels    `json:"page"`
	Buttons LocationButtonLabels  `json:"buttons"`
	Columns LocationColumnLabels  `json:"columns"`
	Empty   LocationEmptyLabels   `json:"empty"`
	Form    LocationFormLabels    `json:"form"`
	Actions LocationActionLabels  `json:"actions"`
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
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  string `json:"status"`
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
	Active                 string `json:"active"`
}

type LocationActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

// ---------------------------------------------------------------------------
// Role labels
// ---------------------------------------------------------------------------

// RoleLabels holds all translatable strings for the role module.
// JSON tags match the "role" wrapper key in retail/role.json.
type RoleLabels struct {
	Page    RolePageLabels    `json:"page"`
	Buttons RoleButtonLabels  `json:"buttons"`
	Columns RoleColumnLabels  `json:"columns"`
	Empty   RoleEmptyLabels   `json:"empty"`
	Form    RoleFormLabels    `json:"form"`
	Actions RoleActionLabels  `json:"actions"`
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

// ---------------------------------------------------------------------------
// Permission labels
// ---------------------------------------------------------------------------

// PermissionLabels holds all translatable strings for the permission module.
type PermissionLabels struct {
	Page    PermissionPageLabels    `json:"page"`
	Buttons PermissionButtonLabels  `json:"buttons"`
	Columns PermissionColumnLabels  `json:"columns"`
	Empty   PermissionEmptyLabels   `json:"empty"`
	Form    PermissionFormLabels    `json:"form"`
	Actions PermissionActionLabels  `json:"actions"`
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
	Name                       string `json:"name"`
	NamePlaceholder            string `json:"namePlaceholder"`
	PermissionCode             string `json:"permissionCode"`
	PermissionCodePlaceholder  string `json:"permissionCodePlaceholder"`
	PermissionType             string `json:"permissionType"`
	Description                string `json:"description"`
	DescriptionPlaceholder     string `json:"descriptionPlaceholder"`
	Active                     string `json:"active"`
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
	Page    RolePermissionPageLabels    `json:"page"`
	Buttons RolePermissionButtonLabels  `json:"buttons"`
	Columns RolePermissionColumnLabels  `json:"columns"`
	Empty   RolePermissionEmptyLabels   `json:"empty"`
	Form    RolePermissionFormLabels    `json:"form"`
	Actions RolePermissionActionLabels  `json:"actions"`
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
	Assign          string `json:"assign"`
	Remove          string `json:"remove"`
	ManagePermissions string `json:"managePermissions"`
}

// ---------------------------------------------------------------------------
// User-Role labels
// ---------------------------------------------------------------------------

// UserRoleLabels holds all translatable strings for the user-role assignment view.
type UserRoleLabels struct {
	Page    UserRolePageLabels    `json:"page"`
	Buttons UserRoleButtonLabels  `json:"buttons"`
	Columns UserRoleColumnLabels  `json:"columns"`
	Empty   UserRoleEmptyLabels   `json:"empty"`
	Form    UserRoleFormLabels    `json:"form"`
	Actions UserRoleActionLabels  `json:"actions"`
}

type UserRolePageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type UserRoleButtonLabels struct {
	AssignRole string `json:"assignRole"`
}

type UserRoleColumnLabels struct {
	RoleName    string `json:"roleName"`
	Description string `json:"description"`
	Color       string `json:"color"`
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
// Workspace labels
// ---------------------------------------------------------------------------

// WorkspaceLabels holds all translatable strings for the workspace module.
type WorkspaceLabels struct {
	Page    WorkspacePageLabels    `json:"page"`
	Buttons WorkspaceButtonLabels  `json:"buttons"`
	Columns WorkspaceColumnLabels  `json:"columns"`
	Empty   WorkspaceEmptyLabels   `json:"empty"`
	Form    WorkspaceFormLabels    `json:"form"`
	Actions WorkspaceActionLabels  `json:"actions"`
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
// Login labels
// ---------------------------------------------------------------------------

// LoginLabels holds i18n strings for the login page.
type LoginLabels struct {
	Title      string `json:"title"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Submit     string `json:"submit"`
	ForgotLink string `json:"forgotLink"`
	Error      string `json:"error"`
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
		Search:             common.Table.Search,
		SearchPlaceholder:  common.Table.SearchPlaceholder,
		Filters:            common.Table.Filters,
		FilterConditions:   common.Table.FilterConditions,
		ClearAll:           common.Table.ClearAll,
		AddCondition:       common.Table.AddCondition,
		Clear:              common.Table.Clear,
		ApplyFilters:       common.Table.ApplyFilters,
		Sort:               common.Table.Sort,
		Columns:            common.Table.Columns,
		Export:              common.Table.Export,
		DensityDefault:     common.Table.Density.Default,
		DensityComfortable: common.Table.Density.Comfortable,
		DensityCompact:     common.Table.Density.Compact,
		Show:               common.Table.Show,
		Entries:             common.Table.Entries,
		Showing:            common.Table.Showing,
		To:                 common.Table.To,
		Of:                 common.Table.Of,
		EntriesLabel:       common.Table.EntriesLabel,
		SelectAll:          common.Table.SelectAll,
		Actions:            common.Table.Actions,
		Prev:               common.Pagination.Prev,
		Next:               common.Pagination.Next,
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

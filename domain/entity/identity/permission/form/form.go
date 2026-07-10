package form

import (
	"github.com/erniealice/pyeza-golang/types"

	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"
)

// Labels holds i18n labels for the drawer form template.
type Labels struct {
	Name                      string
	NamePlaceholder           string
	PermissionCode            string
	PermissionCodePlaceholder string
	PermissionCodeHint        string
	PermissionType            string
	Description               string
	DescriptionPlaceholder    string
	Active                    string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo           string
	PermissionCodeInfo string
	PermissionTypeInfo string
	DescriptionInfo    string
	ActiveInfo         string
}

// Data is the template data for the permission drawer form.
type Data struct {
	FormAction            string
	WorkspaceID           string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit                bool
	ID                    string
	Name                  string
	PermissionCode        string
	PermissionType        string
	Description           string
	Active                bool
	Labels                Labels
	PermissionTypeOptions []types.SelectOption
	CommonLabels          any
}

// BuildLabels constructs Labels using the translator function.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		Name:                      t("form.name"),
		NamePlaceholder:           t("form.name_placeholder"),
		PermissionCode:            t("form.permission_code"),
		PermissionCodePlaceholder: t("form.permission_code_placeholder"),
		PermissionCodeHint:        t("form.permission_code_hint"),
		PermissionType:            t("form.permission_type"),
		Description:               t("form.description"),
		DescriptionPlaceholder:    t("form.description_placeholder"),
		Active:                    t("form.active"),
		NameInfo:                  t("permission.form.name_info"),
		PermissionCodeInfo:        t("permission.form.permission_code_info"),
		PermissionTypeInfo:        t("permission.form.permission_type_info"),
		DescriptionInfo:           t("permission.form.description_info"),
		ActiveInfo:                t("permission.form.active_info"),
	}
}

// BuildPermissionTypeOptions builds the permission type select options.
func BuildPermissionTypeOptions(current string, t func(string) string) []types.SelectOption {
	return []types.SelectOption{
		{Value: "PERMISSION_TYPE_ALLOW", Label: t("shared.badges.allow"), Selected: current == "PERMISSION_TYPE_ALLOW"},
		{Value: "PERMISSION_TYPE_DENY", Label: t("shared.badges.deny"), Selected: current == "PERMISSION_TYPE_DENY"},
	}
}

// ParsePermissionType maps a string token to the PermissionType proto enum.
func ParsePermissionType(s string) permissionpb.PermissionType {
	switch s {
	case "PERMISSION_TYPE_DENY":
		return permissionpb.PermissionType_PERMISSION_TYPE_DENY
	default:
		return permissionpb.PermissionType_PERMISSION_TYPE_ALLOW
	}
}

// FormatPermissionType maps a PermissionType proto enum to a string token.
func FormatPermissionType(pt permissionpb.PermissionType) string {
	switch pt {
	case permissionpb.PermissionType_PERMISSION_TYPE_DENY:
		return "PERMISSION_TYPE_DENY"
	default:
		return "PERMISSION_TYPE_ALLOW"
	}
}

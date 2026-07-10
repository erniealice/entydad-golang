package workspace_user_role

// labels.go — WorkspaceUserRole label structs.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// identity sub-context). Pure structural move — no behaviour change; field
// names, json tags, and string literals are byte-identical. Entity-local
// rename: WorkspaceUserRoleLabels -> Labels, WorkspaceUserRole<Xxx>Labels ->
// <Xxx>Labels.

// Labels holds all translatable strings for the
// workspace_user_role assignment drawer (Phase 3).
type Labels struct {
	Form    FormLabels   `json:"form"`
	Buttons ButtonLabels `json:"buttons"`
}

// FormLabels holds field labels for the assign-form drawer.
type FormLabels struct {
	WorkspaceUser         string `json:"workspace_user"`
	Role                  string `json:"role"`
	RolePlaceholder       string `json:"role_placeholder"`
	RoleSearchPlaceholder string `json:"role_search_placeholder"`
	RoleNoResults         string `json:"role_no_results"`
	Permissions           string `json:"permissions"`
	PermissionsHint       string `json:"permissions_hint"`
}

// ButtonLabels holds button text for the assign-form drawer.
type ButtonLabels struct {
	Submit string `json:"submit"`
	Cancel string `json:"cancel"`
}

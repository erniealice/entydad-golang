package workspace_user

// labels.go — WorkspaceUser label structs.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// identity sub-context). Pure structural move — no behaviour change; field
// names, json tags, and string literals are byte-identical. Entity-local
// rename: WorkspaceUserLabels -> Labels, WorkspaceUser<Xxx>Labels -> <Xxx>Labels.

// Labels holds all translatable strings for the workspace_user module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Columns ColumnLabels `json:"columns"`
	Detail  DetailLabels `json:"detail"`
	Form    FormLabels   `json:"form"`
	Actions ActionLabels `json:"actions"`
}

type PageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type ColumnLabels struct {
	UserName   string `json:"user_name"`
	Email      string `json:"email"`
	Roles      string `json:"roles"`
	Status     string `json:"status"`
	RoleName   string `json:"role_name"`
	PermCount  string `json:"perm_count"`
	DateJoined string `json:"date_joined"`
}

// DetailLabels holds i18n strings for the workspace_user detail page (Phase 2).
type DetailLabels struct {
	BackToWorkspace string            `json:"back_to_workspace"`
	Tabs            DetailTabLabels   `json:"tabs"`
	Roles           DetailRolesLabels `json:"roles"`
	// Info holds the Info-tab section title + field labels (W4.5 label
	// remediation — previously hardcoded in info-tab.html).
	Info DetailInfoLabels `json:"info"`
}

// DetailInfoLabels holds the Info-tab labels on the
// workspace_user detail page (W4.5 label remediation).
type DetailInfoLabels struct {
	SectionTitle string `json:"section_title"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Workspace    string `json:"workspace"`
	DateJoined   string `json:"date_joined"`
	Status       string `json:"status"`
}

type DetailTabLabels struct {
	Info        string `json:"info"`
	Roles       string `json:"roles"`
	Attachments string `json:"attachments"`
}

type DetailRolesLabels struct {
	AssignButton string `json:"assign_button"`
	// Empty-state copy shown when no roles are assigned yet (W4.5 label
	// remediation — previously hardcoded in roles-tab.html).
	EmptyTitle   string `json:"empty_title"`
	EmptyMessage string `json:"empty_message"`
}

type FormLabels struct {
	User                  string `json:"user"`
	UserPlaceholder       string `json:"user_placeholder"`
	UserSearchPlaceholder string `json:"user_search_placeholder"`
	WorkspaceID           string `json:"workspace_id"`
	Active                string `json:"active"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

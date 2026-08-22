package entydad

// labels.go — entydad root LEFTOVERS after the domain-first restructuring.
//
// Migrated entity label types (Client/User/Location/Role/Permission/Workspace/
// Supplier/Tag/PaymentTerm + their dashboards) now live under domain/entity/**
// and are re-exported through the entity facade (domain/entity/entity.go).
//
// What remains here is genuinely NOT an esqyma `entity` symbol:
//   - Shared* + DashboardLabels (domain-wide, imported by the entity packages —
//     cannot move to the facade without a cycle; future home: domain/entity/shared/).
//   - Auth service-surface labels (Login/Signup/ResetPassword/ChangePassword/
//     AuthEmail) — charter-exempt service surface, not an entity facade type.
//   - AdminDashboardLabels — admin service surface (dashboard proto lives under
//     service/dashboard/admin, not domain/entity/).
//   - RoleBadge — small shared helper type.
//
// TaxRegistration* relocated to domain/tax/tax_registration (entity-local
// Labels) + re-exported through the domain/tax facade (fork E4 / thread TX,
// 2026-06-12). Conversation* relocated to hybra views/conversation/model
// (cross-cutting communication surface, view-package-placement.md OCID /
// thread TC, 2026-06-12).

// ---------------------------------------------------------------------------
// User labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Location labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// LocationArea labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Role labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Permission labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Role-Permission labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// User-Role labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Role-User labels (reverse of User-Role: managing users on a role)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Workspace labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// WorkspaceUser labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// WorkspaceUserRoleLabels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Login labels
// ---------------------------------------------------------------------------

// LoginLabels holds i18n strings for the login page.
type LoginLabels struct {
	Title              string `json:"title"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Submit             string `json:"submit"`
	ForgotLink         string `json:"forgot_link"`
	Error              string `json:"error"`
	AdminTitle         string `json:"admin_title"`
	AdminDescription   string `json:"admin_description"`
	EmailPlaceholder   string `json:"email_placeholder"`
	StaffTitle         string `json:"staff_title"`
	StaffDescription   string `json:"staff_description"`
	StaffPinComingSoon string `json:"staff_pin_coming_soon"`
}

// Login02Labels holds i18n strings for the login02 split-screen page.
type Login02Labels struct {
	Title               string `json:"title"`
	Heading             string `json:"heading"`
	Subheading          string `json:"subheading"`
	EmailLabel          string `json:"email_label"`
	EmailPlaceholder    string `json:"email_placeholder"`
	PasswordLabel       string `json:"password_label"`
	PasswordPlaceholder string `json:"password_placeholder"`
	RememberMe          string `json:"remember_me"`
	ForgotPassword      string `json:"forgot_password"`
	SignInButton        string `json:"sign_in_button"`
	ImpersonateButton   string `json:"impersonate_button"`
	NoAccount           string `json:"no_account"`
	SignUpLink          string `json:"sign_up_link"`
	SocialDivider       string `json:"social_divider"`
	Error               string `json:"error"`
	// Carousel navigation
	PreviousSlide string `json:"previous_slide"`
	NextSlide     string `json:"next_slide"`
	ContinueWith  string `json:"continue_with"`
}

// ---------------------------------------------------------------------------
// Signup labels
// ---------------------------------------------------------------------------

// SignupLabels holds i18n strings for the signup01 page (dual-card style).
type SignupLabels struct {
	Title            string `json:"title"`
	Heading          string `json:"heading"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	EmailPlaceholder string `json:"email_placeholder"`
	Password         string `json:"password"`
	ConfirmPassword  string `json:"confirm_password"`
	Submit           string `json:"submit"`
	HasAccount       string `json:"has_account"`
	SignInLink       string `json:"sign_in_link"`
	TermsPrefix      string `json:"terms_prefix"`
	TermsLink        string `json:"terms_link"`
	PrivacyLink      string `json:"privacy_link"`
	AdminTitle       string `json:"admin_title"`
	AdminDescription string `json:"admin_description"`
	StaffTitle       string `json:"staff_title"`
	StaffDescription string `json:"staff_description"`
	PasswordStrength string `json:"password_strength"`
}

// Signup02Labels holds i18n strings for the signup02 page (split-screen style).
type Signup02Labels struct {
	Title                      string `json:"title"`
	Heading                    string `json:"heading"`
	Subheading                 string `json:"subheading"`
	FirstNameLabel             string `json:"first_name_label"`
	FirstNamePlaceholder       string `json:"first_name_placeholder"`
	LastNameLabel              string `json:"last_name_label"`
	LastNamePlaceholder        string `json:"last_name_placeholder"`
	EmailLabel                 string `json:"email_label"`
	EmailPlaceholder           string `json:"email_placeholder"`
	PasswordLabel              string `json:"password_label"`
	PasswordPlaceholder        string `json:"password_placeholder"`
	ConfirmPasswordLabel       string `json:"confirm_password_label"`
	ConfirmPasswordPlaceholder string `json:"confirm_password_placeholder"`
	SignUpButton               string `json:"sign_up_button"`
	HasAccount                 string `json:"has_account"`
	SignInLink                 string `json:"sign_in_link"`
	SocialDivider              string `json:"social_divider"`
	TermsText                  string `json:"terms_text"`
	Error                      string `json:"error"`
	// Carousel navigation + accessibility
	PreviousSlide    string `json:"previous_slide"`
	NextSlide        string `json:"next_slide"`
	ContinueWith     string `json:"continue_with"`
	PasswordStrength string `json:"password_strength"`
	TermsLink        string `json:"terms_link"`
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
	EmailPlaceholder   string `json:"email_placeholder"`
	Submit             string `json:"submit"`
	BackToLogin        string `json:"back_to_login"`
	ConfirmHeading     string `json:"confirm_heading"`
	ConfirmDescription string `json:"confirm_description"`
	NewPassword        string `json:"new_password"`
	ConfirmPassword    string `json:"confirm_password"`
	ResetButton        string `json:"reset_button"`
	SuccessHeading     string `json:"success_heading"`
	SuccessMessage     string `json:"success_message"`
}

// ResetPassword02Labels holds i18n strings for the reset-password02 page (split-screen style).
type ResetPassword02Labels struct {
	Title                      string `json:"title"`
	Heading                    string `json:"heading"`
	Subheading                 string `json:"subheading"`
	EmailLabel                 string `json:"email_label"`
	EmailPlaceholder           string `json:"email_placeholder"`
	SendResetButton            string `json:"send_reset_button"`
	BackToLogin                string `json:"back_to_login"`
	ConfirmHeading             string `json:"confirm_heading"`
	ConfirmSubheading          string `json:"confirm_subheading"`
	NewPasswordLabel           string `json:"new_password_label"`
	NewPasswordPlaceholder     string `json:"new_password_placeholder"`
	ConfirmPasswordLabel       string `json:"confirm_password_label"`
	ConfirmPasswordPlaceholder string `json:"confirm_password_placeholder"`
	ResetButton                string `json:"reset_button"`
	SuccessHeading             string `json:"success_heading"`
	SuccessMessage             string `json:"success_message"`
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
	ErrorMismatch     string `json:"error_mismatch"`
	ErrorInvalidToken string `json:"error_invalid_token"`
	ErrorExpiredToken string `json:"error_expired_token"`
	ErrorWeakPassword string `json:"error_weak_password"`
	// Carousel navigation
	PreviousSlide string `json:"previous_slide"`
	NextSlide     string `json:"next_slide"`
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
	OldPasswordLabel           string `json:"old_password_label"`
	OldPasswordPlaceholder     string `json:"old_password_placeholder"`
	NewPasswordLabel           string `json:"new_password_label"`
	NewPasswordPlaceholder     string `json:"new_password_placeholder"`
	ConfirmPasswordLabel       string `json:"confirm_password_label"`
	ConfirmPasswordPlaceholder string `json:"confirm_password_placeholder"`
	SubmitButton               string `json:"submit_button"`
	SuccessMessage             string `json:"success_message"`
	// Generic fallback + code-specific error messages.
	Error                 string `json:"error"`
	ErrorMismatch         string `json:"error_mismatch"`
	ErrorCurrentIncorrect string `json:"error_current_incorrect"`
	ErrorTooShort         string `json:"error_too_short"`
	BackToApp             string `json:"back_to_app"`
}

// ---------------------------------------------------------------------------
// Auth email labels
// ---------------------------------------------------------------------------

// AuthEmailLabels holds i18n strings for authentication-related email templates.
type AuthEmailLabels struct {
	ResetSubject           string `json:"reset_subject"`
	ResetHeading           string `json:"reset_heading"`
	ResetBody              string `json:"reset_body"`
	ResetButtonText        string `json:"reset_button_text"`
	ResetExpiry            string `json:"reset_expiry"`
	WelcomeSubject         string `json:"welcome_subject"`
	WelcomeHeading         string `json:"welcome_heading"`
	WelcomeBody            string `json:"welcome_body"`
	WelcomeButtonText      string `json:"welcome_button_text"`
	PasswordChangedSubject string `json:"password_changed_subject"`
	PasswordChangedHeading string `json:"password_changed_heading"`
	PasswordChangedBody    string `json:"password_changed_body"`
	SecurityNotice         string `json:"security_notice"`
}

// ---------------------------------------------------------------------------
// Supplier labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Client Tag labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Supplier Tag labels
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// PaymentTerm labels
// ---------------------------------------------------------------------------

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
	PermissionDenied          string `json:"permission_denied"`
	InvalidFormData           string `json:"invalid_form_data"`
	InvalidStatus             string `json:"invalid_status"`
	InvalidTargetStatus       string `json:"invalid_target_status"`
	NotFound                  string `json:"not_found"`
	IDRequired                string `json:"id_required"`
	NoIDsProvided             string `json:"no_ids_provided"`
	PasswordRequired          string `json:"password_required"`
	PasswordFailed            string `json:"password_failed"`
	PasswordManagedByProvider string `json:"password_managed_by_provider"`
	RoleRequired              string `json:"role_required"`
	PermissionRequired        string `json:"permission_required"`
	UserRequired              string `json:"user_required"`
	TagNotFound               string `json:"tag_not_found"`
	TagNameExists             string `json:"tag_name_exists"`
	VerifyFailed              string `json:"verify_failed"`
	CannotDeleteInUse         string `json:"cannot_delete_in_use"`
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
	BulkActivate   string `json:"bulk_activate"`
	BulkDeactivate string `json:"bulk_deactivate"`
	BulkDelete     string `json:"bulk_delete"`
	BulkBlock      string `json:"bulk_block"`
	BulkHold       string `json:"bulk_hold"`
	BulkProspect   string `json:"bulk_prospect"`
}

// SharedBadgeLabels holds translatable badge values.
type SharedBadgeLabels struct {
	Allow        string `json:"allow"`
	Deny         string `json:"deny"`
	Yes          string `json:"yes"`
	No           string `json:"no"`
	NoPermission string `json:"no_permission"`
}

// DashboardLabels holds translatable strings for dashboard pages.
type DashboardLabels struct {
	ClientTitle   string `json:"client_title"`
	UserTitle     string `json:"user_title"`
	SupplierTitle string `json:"supplier_title"`
	LocationTitle string `json:"location_title"`
	AdminTitle    string `json:"admin_title"`
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
	WorkspaceUsers    string `json:"workspace_users"`
	Roles             string `json:"roles"`
	Permissions       string `json:"permissions"`
	RecentRoleChanges string `json:"recent_role_changes"`

	// Widget titles
	UsersPerRole           string `json:"users_per_role"`
	RolesByPermissionCount string `json:"roles_by_permission_count"`
	RecentRoleChangesList  string `json:"recent_role_changes_list"`
	ViewAll                string `json:"view_all"`

	// Quick action labels
	QuickNewUser      string `json:"quick_new_user"`
	QuickNewWorkspace string `json:"quick_new_workspace"`
	QuickAssignRole   string `json:"quick_assign_role"`
	QuickAuditLog     string `json:"quick_audit_log"`

	// Activity / table column labels
	ColumnRole            string `json:"column_role"`
	ColumnPermissionCount string `json:"column_permission_count"`
	RoleAssigned          string `json:"role_assigned"`
}

// ---------------------------------------------------------------------------
// Shared types
// ---------------------------------------------------------------------------

// RoleBadge holds minimal role info for display as a chip/badge in lists.
type RoleBadge struct {
	Name  string
	Color string
}

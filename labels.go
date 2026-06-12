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
	ErrorMismatch     string `json:"errorMismatch"`
	ErrorInvalidToken string `json:"errorInvalidToken"`
	ErrorExpiredToken string `json:"errorExpiredToken"`
	ErrorWeakPassword string `json:"errorWeakPassword"`
	// Carousel navigation
	PreviousSlide string `json:"previousSlide"`
	NextSlide     string `json:"nextSlide"`
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
	OldPasswordLabel           string `json:"oldPasswordLabel"`
	OldPasswordPlaceholder     string `json:"oldPasswordPlaceholder"`
	NewPasswordLabel           string `json:"newPasswordLabel"`
	NewPasswordPlaceholder     string `json:"newPasswordPlaceholder"`
	ConfirmPasswordLabel       string `json:"confirmPasswordLabel"`
	ConfirmPasswordPlaceholder string `json:"confirmPasswordPlaceholder"`
	SubmitButton               string `json:"submitButton"`
	SuccessMessage             string `json:"successMessage"`
	// Generic fallback + code-specific error messages.
	Error                 string `json:"error"`
	ErrorMismatch         string `json:"errorMismatch"`
	ErrorCurrentIncorrect string `json:"errorCurrentIncorrect"`
	ErrorTooShort         string `json:"errorTooShort"`
	BackToApp             string `json:"backToApp"`
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
	LocationTitle string `json:"locationTitle"`
	AdminTitle    string `json:"adminTitle"`
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
	WorkspaceUsers    string `json:"workspaceUsers"`
	Roles             string `json:"roles"`
	Permissions       string `json:"permissions"`
	RecentRoleChanges string `json:"recentRoleChanges"`

	// Widget titles
	UsersPerRole           string `json:"usersPerRole"`
	RolesByPermissionCount string `json:"rolesByPermissionCount"`
	RecentRoleChangesList  string `json:"recentRoleChangesList"`
	ViewAll                string `json:"viewAll"`

	// Quick action labels
	QuickNewUser      string `json:"quickNewUser"`
	QuickNewWorkspace string `json:"quickNewWorkspace"`
	QuickAssignRole   string `json:"quickAssignRole"`
	QuickAuditLog     string `json:"quickAuditLog"`

	// Activity / table column labels
	ColumnRole            string `json:"columnRole"`
	ColumnPermissionCount string `json:"columnPermissionCount"`
	RoleAssigned          string `json:"roleAssigned"`
}

// ---------------------------------------------------------------------------
// Shared types
// ---------------------------------------------------------------------------

// RoleBadge holds minimal role info for display as a chip/badge in lists.
type RoleBadge struct {
	Name  string
	Color string
}

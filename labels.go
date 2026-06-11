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
//   - TaxRegistration* — domain/tax (fork E4 / W7), still root-resident.
//   - Conversation* — domain/communication (→ hybra, W5), still root-resident.
//   - RoleBadge — small shared helper type.

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

// ---------------------------------------------------------------------------
// TaxRegistrationLabels
// Lyngua root key: "taxRegistration"
// ---------------------------------------------------------------------------

// TaxRegistrationLabels holds all translatable strings for the polymorphic
// Tax Registration views (client + workspace party types in v1).
type TaxRegistrationLabels struct {
	Page    TaxRegistrationPageLabels   `json:"page"`
	Columns TaxRegistrationColumnLabels `json:"columns"`
	Buttons TaxRegistrationButtonLabels `json:"buttons"`
	Actions TaxRegistrationActionLabels `json:"actions"`
	Empty   TaxRegistrationEmptyLabels  `json:"empty"`
	Fields  TaxRegistrationFieldLabels  `json:"fields"`
	Revoke  TaxRegistrationRevokeLabels `json:"revoke"`
}

// TaxRegistrationPageLabels holds page heading strings.
type TaxRegistrationPageLabels struct {
	Heading          string `json:"heading"`
	HeadingClient    string `json:"headingClient"`
	HeadingWorkspace string `json:"headingWorkspace"`
	Caption          string `json:"caption"`
	AddDrawerTitle   string `json:"addDrawerTitle"`
	EditDrawerTitle  string `json:"editDrawerTitle"`
}

// TaxRegistrationColumnLabels holds table column headers.
type TaxRegistrationColumnLabels struct {
	KindName           string `json:"kindName"`
	ComputePath        string `json:"computePath"`
	PartyRole          string `json:"partyRole"`
	Status             string `json:"status"`
	EffectiveFrom      string `json:"effectiveFrom"`
	RegistrationNumber string `json:"registrationNumber"`
}

// TaxRegistrationButtonLabels holds button text.
type TaxRegistrationButtonLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

// TaxRegistrationActionLabels holds action dropdown labels.
type TaxRegistrationActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	NoPermission string `json:"noPermission"`
}

// TaxRegistrationEmptyLabels holds empty-state strings.
type TaxRegistrationEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// TaxRegistrationFieldLabels holds drawer form field labels.
type TaxRegistrationFieldLabels struct {
	TaxRegistrationKindID string `json:"taxRegistrationKindId"`
	RegistrationNumber    string `json:"registrationNumber"`
	EffectiveFrom         string `json:"effectiveFrom"`
	Notes                 string `json:"notes"`
	Status                string `json:"status"`
}

// TaxRegistrationRevokeLabels holds strings for the revoke confirm dialog.
type TaxRegistrationRevokeLabels struct {
	WarningMessage        string `json:"warningMessage"`
	EffectiveTo           string `json:"effectiveTo"`
	AffectedPeriodsNotice string `json:"affectedPeriodsNotice"`
	// AffectedPeriodsCount is the row label for the pending-period count (Phase 5 M3).
	AffectedPeriodsCount string `json:"affectedPeriodsCount"`
	// AffectedSubscriptionsCount is the row label for the subscription count (Phase 5 M3).
	AffectedSubscriptionsCount string `json:"affectedSubscriptionsCount"`
	ReasonLabel                string `json:"reasonLabel"`
	ReasonPlaceholder          string `json:"reasonPlaceholder"`
	ConfirmButton              string `json:"confirmButton"`
}

// DefaultTaxRegistrationLabels returns TaxRegistrationLabels with sensible
// English defaults.
func DefaultTaxRegistrationLabels() TaxRegistrationLabels {
	return TaxRegistrationLabels{
		Page: TaxRegistrationPageLabels{
			Heading:          "Tax Registrations",
			HeadingClient:    "Client Tax Registrations",
			HeadingWorkspace: "Workspace Tax Registrations",
			Caption:          "Active tax registrations determine compute path during revenue recognition",
			AddDrawerTitle:   "Add Tax Registration",
			EditDrawerTitle:  "Edit Tax Registration",
		},
		Columns: TaxRegistrationColumnLabels{
			KindName:           "Kind",
			ComputePath:        "Compute Path",
			PartyRole:          "Party Role",
			Status:             "Status",
			EffectiveFrom:      "Effective From",
			RegistrationNumber: "Registration No.",
		},
		Buttons: TaxRegistrationButtonLabels{
			Add:    "Add Tax Registration",
			Edit:   "Edit",
			Delete: "Delete",
		},
		Actions: TaxRegistrationActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			NoPermission: "You do not have permission to manage tax registrations",
		},
		Empty: TaxRegistrationEmptyLabels{
			Title:   "No tax registrations",
			Message: "Add a tax registration to enable tax computation for this party.",
		},
		Fields: TaxRegistrationFieldLabels{
			TaxRegistrationKindID: "Tax Registration Kind",
			RegistrationNumber:    "Registration Number",
			EffectiveFrom:         "Effective From",
			Notes:                 "Notes",
			Status:                "Status",
		},
		Revoke: TaxRegistrationRevokeLabels{
			WarningMessage:             "Revoking this registration will affect pending billing periods. Ensure all outstanding periods are settled before proceeding.",
			EffectiveTo:                "Effective To",
			AffectedPeriodsNotice:      "Some pending subscription billing periods fall within the revocation window and may need to be reprocessed.",
			AffectedPeriodsCount:       "Affected billing periods",
			AffectedSubscriptionsCount: "Affected subscriptions",
			ReasonLabel:                "Reason for revocation",
			ReasonPlaceholder:          "Describe why this registration is being revoked",
			ConfirmButton:              "Revoke Registration",
		},
	}
}

// ===========================================================================
// Conversation labels — secure messaging / ticketing (Plan-4, 2026-06-03)
//
// Loaded from translations/en/{tier}/conversation.json (root key "conversation")
// and conversation_post.json (root key "conversationPost") via LoadPathIfExists.
// All fields are nil-safe: DefaultConversationLabels() pre-populates English so a
// missing JSON file does not produce empty strings in the UI.
// ===========================================================================

// ConversationLabels is the top-level label struct for the conversation surface.
type ConversationLabels struct {
	List    ConversationListLabels    `json:"list"`
	Inbox   ConversationInboxLabels   `json:"inbox"`
	Thread  ConversationThreadLabels  `json:"thread"`
	Status  ConversationStatusLabels  `json:"status"`
	Actions ConversationActionLabels  `json:"actions"`
	Form    ConversationFormLabels    `json:"form"`
	Columns ConversationColumnLabels  `json:"columns"`
	Confirm ConversationConfirmLabels `json:"confirm"`
	Errors  ConversationErrorLabels   `json:"errors"`
}

// ConversationListLabels — staff inbox + portal thread-list headings.
type ConversationListLabels struct {
	Heading      string `json:"heading"`
	Subtitle     string `json:"subtitle"`
	Title        string `json:"title"`
	NewButton    string `json:"newButton"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// ConversationInboxLabels — staff filter chips.
type ConversationInboxLabels struct {
	FilterAll        string `json:"filterAll"`
	FilterUnassigned string `json:"filterUnassigned"`
	FilterMyQueue    string `json:"filterMyQueue"`
	FilterOpen       string `json:"filterOpen"`
	FilterInProgress string `json:"filterInProgress"`
	FilterResolved   string `json:"filterResolved"`
	FilterClosed     string `json:"filterClosed"`
}

// ConversationThreadLabels — thread-detail header + meta.
type ConversationThreadLabels struct {
	BackToInbox   string `json:"backToInbox"`
	Assignee      string `json:"assignee"`
	Unassigned    string `json:"unassigned"`
	Client        string `json:"client"`
	Created       string `json:"created"`
	LastActivity  string `json:"lastActivity"`
	ViewRequest   string `json:"viewRequest"`
	Subtitle      string `json:"subtitle"`
	EmptyTitle    string `json:"emptyTitle"`
	EmptySubtitle string `json:"emptySubtitle"`
}

// ConversationStatusLabels — human-readable status badge labels keyed by enum.
type ConversationStatusLabels struct {
	Open       string `json:"open"`
	InProgress string `json:"inProgress"`
	Resolved   string `json:"resolved"`
	Closed     string `json:"closed"`
	Unknown    string `json:"unknown"`
}

// ConversationActionLabels — action button labels.
type ConversationActionLabels struct {
	NewConversation string `json:"newConversation"`
	Open            string `json:"open"`
	Assign          string `json:"assign"`
	MarkResolved    string `json:"markResolved"`
	Close           string `json:"close"`
	Reopen          string `json:"reopen"`
	SetStatus       string `json:"setStatus"`
	Send            string `json:"send"`
}

// ConversationFormLabels — new-conversation / assign / status drawer fields.
type ConversationFormLabels struct {
	SectionTitle         string `json:"sectionTitle"`
	SubjectLabel         string `json:"subjectLabel"`
	SubjectPlaceholder   string `json:"subjectPlaceholder"`
	ClientLabel          string `json:"clientLabel"`
	ClientPlaceholder    string `json:"clientPlaceholder"`
	AssigneeLabel        string `json:"assigneeLabel"`
	AssigneePlaceholder  string `json:"assigneePlaceholder"`
	LinkLabel            string `json:"linkLabel"`
	LinkPlaceholder      string `json:"linkPlaceholder"`
	MessageLabel         string `json:"messageLabel"`
	MessagePlaceholder   string `json:"messagePlaceholder"`
	CurrentStatusLabel   string `json:"currentStatusLabel"`
	NewStatusLabel       string `json:"newStatusLabel"`
	CurrentAssigneeLabel string `json:"currentAssigneeLabel"`
}

// ConversationColumnLabels — staff inbox table headers.
type ConversationColumnLabels struct {
	Subject      string `json:"subject"`
	Client       string `json:"client"`
	LastActivity string `json:"lastActivity"`
	Assignee     string `json:"assignee"`
	Status       string `json:"status"`
}

// ConversationConfirmLabels — confirm-dialog copy for status transitions.
type ConversationConfirmLabels struct {
	ResolveTitle   string `json:"resolveTitle"`
	ResolveMessage string `json:"resolveMessage"`
	CloseTitle     string `json:"closeTitle"`
	CloseMessage   string `json:"closeMessage"`
	ReopenTitle    string `json:"reopenTitle"`
	ReopenMessage  string `json:"reopenMessage"`
}

// ConversationErrorLabels — error strings surfaced via HTMX error toast.
type ConversationErrorLabels struct {
	PermissionDenied  string `json:"permissionDenied"`
	NotFound          string `json:"notFound"`
	InvalidForm       string `json:"invalidForm"`
	SubjectRequired   string `json:"subjectRequired"`
	ClientRequired    string `json:"clientRequired"`
	MessageRequired   string `json:"messageRequired"`
	InvalidTransition string `json:"invalidTransition"`
	IDRequired        string `json:"idRequired"`
	SaveFailed        string `json:"saveFailed"`
}

// ConversationPostLabels is the label struct for the post composer + bubbles.
// Loaded from conversation_post.json (root key "conversationPost").
type ConversationPostLabels struct {
	Composer ConversationComposerLabels  `json:"composer"`
	Bubble   ConversationBubbleLabels    `json:"bubble"`
	Subtitle string                      `json:"subtitle"`
	Empty    string                      `json:"empty"`
	Errors   ConversationPostErrorLabels `json:"errors"`
}

// ConversationComposerLabels — reply composer.
type ConversationComposerLabels struct {
	Placeholder string `json:"placeholder"`
	Send        string `json:"send"`
	Attach      string `json:"attach"`
}

// ConversationBubbleLabels — sender role labels.
type ConversationBubbleLabels struct {
	You    string `json:"you"`
	Staff  string `json:"staff"`
	Client string `json:"client"`
}

// ConversationPostErrorLabels — post-specific errors.
type ConversationPostErrorLabels struct {
	EmptyBody    string `json:"emptyBody"`
	MissingToken string `json:"missingToken"`
	SendFailed   string `json:"sendFailed"`
}

// DefaultConversationLabels returns English defaults for the conversation
// surface. Override per business type via conversation.json.
func DefaultConversationLabels() ConversationLabels {
	return ConversationLabels{
		List: ConversationListLabels{
			Heading:      "Conversations",
			Subtitle:     "Secure messaging with your clients",
			Title:        "Messages",
			NewButton:    "New",
			EmptyTitle:   "No conversations yet",
			EmptyMessage: "Start a new conversation to message a client.",
		},
		Inbox: ConversationInboxLabels{
			FilterAll:        "All open",
			FilterUnassigned: "Unassigned",
			FilterMyQueue:    "My queue",
			FilterOpen:       "Open",
			FilterInProgress: "In progress",
			FilterResolved:   "Resolved",
			FilterClosed:     "Closed",
		},
		Thread: ConversationThreadLabels{
			BackToInbox:   "Back to inbox",
			Assignee:      "Assigned to",
			Unassigned:    "Unassigned",
			Client:        "Client",
			Created:       "Created",
			LastActivity:  "Last activity",
			ViewRequest:   "View request",
			Subtitle:      "Secure messaging. Every conversation is logged.",
			EmptyTitle:    "Select a conversation",
			EmptySubtitle: "Choose a thread from the list to view messages.",
		},
		Status: ConversationStatusLabels{
			Open:       "Open",
			InProgress: "In progress",
			Resolved:   "Resolved",
			Closed:     "Closed",
			Unknown:    "Unknown",
		},
		Actions: ConversationActionLabels{
			NewConversation: "New conversation",
			Open:            "Open",
			Assign:          "Assign",
			MarkResolved:    "Mark resolved",
			Close:           "Close",
			Reopen:          "Reopen",
			SetStatus:       "Change status",
			Send:            "Send",
		},
		Form: ConversationFormLabels{
			SectionTitle:         "Conversation details",
			SubjectLabel:         "Subject",
			SubjectPlaceholder:   "What is this about?",
			ClientLabel:          "Client",
			ClientPlaceholder:    "Search clients…",
			AssigneeLabel:        "Assign to",
			AssigneePlaceholder:  "Search staff…",
			LinkLabel:            "Linked request (optional)",
			LinkPlaceholder:      "e.g. REQ-0091",
			MessageLabel:         "Message",
			MessagePlaceholder:   "Type your first message…",
			CurrentStatusLabel:   "Current status",
			NewStatusLabel:       "New status",
			CurrentAssigneeLabel: "Currently",
		},
		Columns: ConversationColumnLabels{
			Subject:      "Conversation",
			Client:       "Client",
			LastActivity: "Last activity",
			Assignee:     "Assigned",
			Status:       "Status",
		},
		Confirm: ConversationConfirmLabels{
			ResolveTitle:   "Mark resolved",
			ResolveMessage: "Mark this conversation as resolved?",
			CloseTitle:     "Close conversation",
			CloseMessage:   "Close this conversation? It can be reopened later.",
			ReopenTitle:    "Reopen conversation",
			ReopenMessage:  "Reopen this conversation?",
		},
		Errors: ConversationErrorLabels{
			PermissionDenied:  "You don't have permission to perform this action.",
			NotFound:          "Conversation not found.",
			InvalidForm:       "Invalid form data.",
			SubjectRequired:   "A subject is required.",
			ClientRequired:    "Please select a client.",
			MessageRequired:   "A message is required.",
			InvalidTransition: "That status change is not allowed.",
			IDRequired:        "A conversation id is required.",
			SaveFailed:        "Could not save. Please try again.",
		},
	}
}

// DefaultConversationPostLabels returns English defaults for the composer /
// bubble surface. Override per business type via conversation_post.json.
func DefaultConversationPostLabels() ConversationPostLabels {
	return ConversationPostLabels{
		Composer: ConversationComposerLabels{
			Placeholder: "Reply…",
			Send:        "Send",
			Attach:      "Attach",
		},
		Bubble: ConversationBubbleLabels{
			You:    "You",
			Staff:  "Staff",
			Client: "Client",
		},
		Subtitle: "Secure messaging. Every conversation is logged.",
		Empty:    "No messages yet.",
		Errors: ConversationPostErrorLabels{
			EmptyBody:    "Message cannot be empty.",
			MissingToken: "Missing idempotency token. Please refresh and try again.",
			SendFailed:   "Could not send your message. Please try again.",
		},
	}
}

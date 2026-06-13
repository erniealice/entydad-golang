package auth

import (
	"context"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"
)

// AuthAdapter is the narrow interface entydad's auth module needs for
// credential operations. Satisfied by espyna's consumer.AuthAdapter.
// Deliberately does not include CreateSession or GetSessionWorkspaceContext
// — those are consumed by other layers (session middleware, principal switch).
type AuthAdapter interface {
	Login(ctx context.Context, email, password string) (token string, identity AuthIdentity, err error)
	Register(ctx context.Context, email, password, firstName, lastName, mobileNumber string) (userID string, err error)
	RequestPasswordReset(ctx context.Context, email string) (resetToken string, err error)
	ExecutePasswordReset(ctx context.Context, token, newPassword string) error
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	ValidateSession(ctx context.Context, token string) (userID string, err error)
	InvalidateSession(ctx context.Context, token string) error
}

// AuthIdentity is the minimal identity returned by Login.
// Satisfied by esqyma's authpb.Identity via its GetId() method.
type AuthIdentity interface {
	GetId() string
}

// SessionManager handles session cookie write/clear on HTTP responses.
// Satisfied by espyna's consumer.SessionMiddleware.
type SessionManager interface {
	SetSessionCookie(w http.ResponseWriter, token string)
	ClearSessionCookie(w http.ResponseWriter)
}

// PrincipalResolver loads and resolves a user's active principal bindings.
// Satisfied by adapthttp.DBPrincipalLoader.
type PrincipalResolver interface {
	Resolve(ctx context.Context, userID string) ([]Principal, error)
	IsEnabled() bool
}

// PrincipalSwitcher performs the security-critical session rotation or
// in-place mutation when switching between principals.
//
// This is a function type rather than an interface because the auth
// module has exactly one callsite pattern. The app's composition layer
// wraps appBuilder.executePrincipalSwitch into this closure.
type PrincipalSwitcher func(ctx context.Context, input PrincipalSwitchInput) (*PrincipalSwitchResult, error)

// CSRFIssuer writes workspace-claim CSRF cookies on session rotation.
// Injected as a function closure from the app's middleware layer.
// The return value (the signed token) is ignored by the auth module
// but preserved in the signature so middleware.IssueWorkspaceCSRFCookie
// can be passed directly.
type CSRFIssuer func(w http.ResponseWriter, secret []byte, sessionToken, workspaceID string) string

// Renderer renders a named template with data to an HTTP response.
// Used by the auth-shell bypass path (login, select-workspace-role).
// Satisfied by *pyeza.HTMLRenderer.
type Renderer interface {
	Render(w http.ResponseWriter, templateName string, data interface{}) error
}

// UserIDByEmail resolves a user's ID from their email address.
// Used as a fallback when the auth adapter's Login response doesn't
// include a user ID (mock providers). Injected as a closure.
type UserIDByEmail func(ctx context.Context, email string) (userID string)

// WorkspaceSlugResolver resolves a workspace's slug from its ID for
// constructing post-login redirect URLs (/w/{slug}/home).
type WorkspaceSlugResolver func(ctx context.Context, workspaceID string) (slug string)

// RouteRegistrar extends pyeza's view.RouteRegistrar with HandleFunc for
// raw http.HandlerFunc registration. The auth module needs both: GET
// (for view-based routes like login/signup pages) and HandleFunc (for
// POST handlers and non-view GETs like /auth/no-access, /auth/logout).
type RouteRegistrar interface {
	view.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

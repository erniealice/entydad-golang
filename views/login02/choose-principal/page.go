package choosePrincipal

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PrincipalCard is one tile rendered on the chooser page. There is one
// card per active principal binding the User holds. The host application
// builds the slice and passes it via the request context (see
// service-admin/domain_auth.go for the wiring).
//
// TestIDKind is the lowercase principal-type token (e.g. "operator_owner",
// "client_delegate"); TestIDID is the row id (e.g. ws-user-abc-123). The
// template derives `data-testid="choose-principal-{kind}"` and, when more
// than one card shares the same kind, `choose-principal-{kind}-{id}`.
type PrincipalCard struct {
	Kind        string // lowercase token: operator_owner / operator_staff / client / client_delegate / supplier / supplier_delegate
	PrincipalID string // row id of the underlying grant (workspace_user.id / client_portal_grant.id / etc.)
	DisplayName string // human label (e.g. "Acme Corp · Owner", "Jane Smith · Client Delegate")
	Subtitle    string // optional secondary line — e.g. workspace name for context
	IconName    string // pyeza icon template name (e.g. "icon-shield-check" for owner)
	// HasMultipleOfKind drives the testid suffix. When true, the template
	// emits `data-testid="choose-principal-{kind}-{id}"` so Playwright can
	// pick a specific card when several share the same kind (multi-target
	// delegate, multi-client-grant user).
	HasMultipleOfKind bool
}

// CardsResolver returns the cards to render for the current request.
// Service-admin's wiring injects a function that wraps principalLoader.Resolve
// and converts each Principal to a PrincipalCard. Returning an empty slice
// is legal — the template renders an empty-state line.
type CardsResolver func(ctx context.Context) []PrincipalCard

// Deps holds view dependencies for the choose-principal page.
//
// ResolveCards is called once per request — the view reads the User from
// context (via consumer.GetUserIDFromContext on the host side) and returns
// the live list of principal cards. This indirection keeps the view
// package free of any DB / loader concerns.
type Deps struct {
	Labels        Labels
	Login02       entydad.Login02Labels // re-uses the auth shell's slide labels for visual parity
	CommonLabels  pyeza.CommonLabels
	LogoText      string
	LogoIcon      string
	SwitchPostURL string        // default: /action/auth/switch-principal
	LogoutURL     string        // default: /auth/logout
	ResolveCards  CardsResolver // required for production; nil falls back to ctx cards (test-only)
}

// Labels holds i18n strings for the choose-principal page.
// JSON tags match the "choose_principal" subtree in common/auth.json so
// lyngua.LoadPath("en", bt, "auth.json", "choose_principal", &labels)
// populates every field directly.
type Labels struct {
	Page                struct {
		Title      string `json:"title"`
		Heading    string `json:"heading"`
		Subheading string `json:"subheading"`
	} `json:"page"`
	SignOutLink          string `json:"signOutLink"`
	SubmitLabel         string `json:"submitLabel"`
	EmptyState          string `json:"emptyState"`
	ErrorSwitchPrincipal string `json:"errorSwitchPrincipal"`

	// Flat convenience accessors (populated by DefaultLabels / lyngua shim).
	Title      string
	Heading    string
	Subheading string
}

// DefaultLabels returns English defaults so the page can render before
// lyngua wiring lands. Service-admin can override by passing its own Labels.
func DefaultLabels() Labels {
	l := Labels{
		SignOutLink:          "Sign out",
		SubmitLabel:         "Continue as",
		EmptyState:          "You don't have any active profiles for this account.",
		ErrorSwitchPrincipal: "Could not switch principal. Please try again.",
	}
	l.Page.Title = "Choose a profile"
	l.Page.Heading = "Choose a profile"
	l.Page.Subheading = "You have access to more than one profile. Pick the one you want to use right now — you can switch later."
	// Populate flat accessors from nested Page struct.
	l.Title = l.Page.Title
	l.Heading = l.Page.Heading
	l.Subheading = l.Page.Subheading
	return l
}

// PageData is the template-facing data shape. ContentTemplate is fixed —
// always "choose-principal-content".
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          Labels
	Login02         entydad.Login02Labels
	LogoText        string
	LogoIcon        string
	SwitchPostURL   string
	LogoutURL       string
	Cards           []PrincipalCard
	Error           string // surfaces ?error= query param messages
}

// ctxKeyCards is a typed context key used by host code to inject the
// per-request Cards slice. The key is unexported so callers must use
// WithCards / getCards — preventing accidental collisions with other
// context values.
type ctxKey int

const ctxKeyCards ctxKey = 0

// WithCards returns a derived context carrying the given cards. Host code
// (domain_auth.go) calls this from its GET /auth/choose-principal handler
// after running principalLoader.Resolve.
func WithCards(ctx context.Context, cards []PrincipalCard) context.Context {
	return context.WithValue(ctx, ctxKeyCards, cards)
}

// getCards returns cards previously installed via WithCards, or nil.
func getCards(ctx context.Context) []PrincipalCard {
	if v, ok := ctx.Value(ctxKeyCards).([]PrincipalCard); ok {
		return v
	}
	return nil
}

// NewView creates the choose-principal page view (GET /auth/choose-principal).
//
// The view reads cards from the request context via getCards — see the
// WithCards helper above. This indirection lets the host wire principal
// resolution at handler-mount time without changing the view signature.
//
// When the context has zero cards, the page renders an explanatory empty
// state. The host SHOULD redirect to /auth/no-access in that case rather
// than relying on this fallback, but the fallback exists so the page never
// crashes.
func NewView(deps *Deps) view.View {
	switchPostURL := deps.SwitchPostURL
	if switchPostURL == "" {
		switchPostURL = "/action/auth/switch-principal"
	}
	logoutURL := deps.LogoutURL
	if logoutURL == "" {
		logoutURL = entydad.AuthLogoutURL
	}
	labels := deps.Labels
	if labels.Title == "" && labels.Page.Title == "" {
		labels = DefaultLabels()
	}
	// Ensure flat Title/Heading/Subheading are populated from the nested Page
	// struct (covers the case where Labels was loaded via lyngua JSON tags).
	if labels.Title == "" {
		labels.Title = labels.Page.Title
	}
	if labels.Heading == "" {
		labels.Heading = labels.Page.Heading
	}
	if labels.Subheading == "" {
		labels.Subheading = labels.Page.Subheading
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		errorMsg := ""
		if viewCtx.Request != nil {
			if viewCtx.Request.URL.Query().Get("error") != "" {
				errorMsg = labels.ErrorSwitchPrincipal
				if errorMsg == "" {
					errorMsg = "Could not switch principal. Please try again."
				}
			}
		}

		var cards []PrincipalCard
		if deps.ResolveCards != nil {
			cards = deps.ResolveCards(ctx)
		} else {
			cards = getCards(ctx)
		}

		// Annotate cards with HasMultipleOfKind so the template can
		// disambiguate test ids when several cards of the same kind exist
		// (e.g. a delegate-of-two-clients case shouldn't be reachable here
		// because delegate cards collapse via ActingAsTargets, but we keep
		// the flag generic for future principal models).
		kindCounts := make(map[string]int, len(cards))
		for _, c := range cards {
			kindCounts[c.Kind]++
		}
		annotated := make([]PrincipalCard, len(cards))
		for i, c := range cards {
			if kindCounts[c.Kind] > 1 {
				c.HasMultipleOfKind = true
			}
			annotated[i] = c
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        labels.Title,
				CurrentPath:  viewCtx.CurrentPath,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "choose-principal-content",
			Labels:          labels,
			Login02:         deps.Login02,
			LogoText:        deps.LogoText,
			LogoIcon:        deps.LogoIcon,
			SwitchPostURL:   switchPostURL,
			LogoutURL:       logoutURL,
			Cards:           annotated,
			Error:           errorMsg,
		}

		return view.OK("choose-principal", pageData)
	})
}

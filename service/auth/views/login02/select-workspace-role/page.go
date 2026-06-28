package selectWorkspaceRole

import (
	"context"
	"strings"

	entydad "github.com/erniealice/entydad-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// KindLabels holds per-principal-kind role labels for the chooser page.
// Each field corresponds to the canonical lowercase Kind token produced by
// PrincipalTypeString (see entydad/service/auth/types.go). Fields are
// populated by lyngua.LoadPathIfExists("en", bt, "auth.json",
// "select_workspace_role", &labels): common/auth.json provides generic
// English values ("Staff", "Client") and business-type overlays supply
// vertical-specific wording (education/auth.json → staff: "Teacher").
type KindLabels struct {
	OperatorOwner    string `json:"operator_owner"`
	OperatorStaff    string `json:"operator_staff"`
	Client           string `json:"client"`
	ClientDelegate   string `json:"client_delegate"`
	Supplier         string `json:"supplier"`
	SupplierDelegate string `json:"supplier_delegate"`
	// Staff is PRINCIPAL_TYPE_STAFF=7 (anchored on staff.id). Generic
	// value: "Staff"; education override: "Teacher".
	Staff string `json:"staff"`
}

// LabelFor returns the configured label for the given Kind token, falling
// back to a humanized version of the token (e.g. "client_delegate" →
// "Client Delegate") when no label is set. This ensures cards never render
// a blank role badge even when lyngua loading is not wired.
func (kl KindLabels) LabelFor(kind string) string {
	var s string
	switch kind {
	case "operator_owner":
		s = kl.OperatorOwner
	case "operator_staff":
		s = kl.OperatorStaff
	case "client":
		s = kl.Client
	case "client_delegate":
		s = kl.ClientDelegate
	case "supplier":
		s = kl.Supplier
	case "supplier_delegate":
		s = kl.SupplierDelegate
	case "staff":
		s = kl.Staff
	}
	if s != "" {
		return s
	}
	return humanizeKind(kind)
}

// humanizeKind converts a snake_case kind token to a title-cased label.
// Example: "client_delegate" → "Client Delegate".
// Used as a last-resort fallback when no lyngua label is configured.
func humanizeKind(kind string) string {
	if kind == "" {
		return ""
	}
	parts := strings.Split(kind, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

// PrincipalCard is one tile rendered on the chooser page. There is one
// card per active principal binding the User holds. The host application
// builds the slice and passes it via the request context (see
// service-admin/domain_auth.go for the wiring).
//
// TestIDKind is the lowercase principal-type token (e.g. "operator_owner",
// "client_delegate"); TestIDID is the row id (e.g. ws-user-abc-123). The
// template derives `data-testid="select-workspace-role-{kind}"` and, when more
// than one card shares the same kind, `select-workspace-role-{kind}-{id}`.
type PrincipalCard struct {
	Kind        string // lowercase token: operator_owner / operator_staff / client / client_delegate / supplier / supplier_delegate / staff
	PrincipalID string // row id of the underlying grant (workspace_user.id / client_portal_grant.id / etc.)
	DisplayName string // human label (e.g. "Acme Corp · Owner", "Jane Smith · Client Delegate")
	Subtitle    string // optional secondary line — e.g. workspace name for context
	IconName    string // pyeza icon template name (e.g. "icon-shield-check" for owner)
	// KindLabel is the vertical-localized role label for this principal
	// kind, e.g. "Staff" (generic) or "Teacher" (education). Populated by
	// NewView from Labels.KindLabels.LabelFor(Kind); never empty (humanized
	// fallback ensures a value is always present).
	KindLabel string
	// HasMultipleOfKind drives the testid suffix. When true, the template
	// emits `data-testid="select-workspace-role-{kind}-{id}"` so Playwright can
	// pick a specific card when several share the same kind (multi-target
	// delegate, multi-client-grant user).
	HasMultipleOfKind bool
}

// CardsResolver returns the cards to render for the current request.
// Service-admin's wiring injects a function that wraps principalLoader.Resolve
// and converts each Principal to a PrincipalCard. Returning an empty slice
// is legal — the template renders an empty-state line.
type CardsResolver func(ctx context.Context) []PrincipalCard

// Deps holds view dependencies for the select-workspace-role page.
//
// ResolveCards is called once per request — the view reads the User from
// context (via identity.FromContext on the host side) and returns
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

// Labels holds i18n strings for the select-workspace-role page.
// JSON tags match the "select_workspace_role" subtree in common/auth.json so
// lyngua.LoadPathIfExists("en", bt, "auth.json", "select_workspace_role", &labels)
// populates every field directly. Business-type overlays (e.g.
// education/auth.json) can override individual kindLabels keys via lyngua's
// recursive deep-merge.
type Labels struct {
	Page struct {
		Title      string `json:"title"`
		Heading    string `json:"heading"`
		Subheading string `json:"subheading"`
	} `json:"page"`
	SignOutLink          string `json:"signOutLink"`
	SubmitLabel          string `json:"submitLabel"`
	EmptyState           string `json:"emptyState"`
	ErrorSwitchPrincipal string `json:"errorSwitchPrincipal"`

	// KindLabels maps principal Kind tokens to localized role labels.
	// Loaded from the "kindLabels" sub-object within the
	// "select_workspace_role" JSON subtree. Education overlay sets
	// kindLabels.staff = "Teacher"; all other kinds default to generic
	// English values from common/auth.json.
	KindLabels KindLabels `json:"kindLabels"`

	// Flat convenience accessors (populated by DefaultLabels / lyngua shim).
	Title      string
	Heading    string
	Subheading string
}

// DefaultLabels returns English defaults so the page can render before
// lyngua wiring lands. Service-admin can override by passing its own Labels
// loaded from auth.json via LoadPathIfExists.
func DefaultLabels() Labels {
	l := Labels{
		SignOutLink:          "Sign out",
		SubmitLabel:         "Continue as",
		EmptyState:          "You don't have any active workspace roles for this account.",
		ErrorSwitchPrincipal: "Could not switch principal. Please try again.",
		// Generic English role labels for all canonical principal kinds.
		// Business-type overlays (e.g. education/auth.json) can override
		// individual keys; the lyngua deep-merge preserves the rest.
		KindLabels: KindLabels{
			OperatorOwner:    "Owner",
			OperatorStaff:    "Staff",
			Client:           "Client",
			ClientDelegate:   "Delegate",
			Supplier:         "Supplier",
			SupplierDelegate: "Delegate",
			Staff:            "Staff", // education overlay → "Teacher"
		},
	}
	l.Page.Title = "Select workspace role"
	l.Page.Heading = "Select workspace role"
	l.Page.Subheading = "You have access to more than one workspace role. Pick the one you want to use right now — you can switch later."
	// Populate flat accessors from nested Page struct.
	l.Title = l.Page.Title
	l.Heading = l.Page.Heading
	l.Subheading = l.Page.Subheading
	return l
}

// PageData is the template-facing data shape. ContentTemplate is fixed —
// always "select-workspace-role-content".
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
// (domain_auth.go) calls this from its GET /auth/select-workspace-role handler
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

// NewView creates the select-workspace-role page view (GET /auth/select-workspace-role).
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
			// Populate the localized role badge from the lyngua-backed
			// KindLabels map. LabelFor never returns empty — it falls back
			// to a humanized token so the badge always has a value.
			c.KindLabel = labels.KindLabels.LabelFor(c.Kind)
			annotated[i] = c
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        labels.Title,
				CurrentPath:  viewCtx.CurrentPath,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "select-workspace-role-content",
			Labels:          labels,
			Login02:         deps.Login02,
			LogoText:        deps.LogoText,
			LogoIcon:        deps.LogoIcon,
			SwitchPostURL:   switchPostURL,
			LogoutURL:       logoutURL,
			Cards:           annotated,
			Error:           errorMsg,
		}

		return view.OK("select-workspace-role", pageData)
	})
}

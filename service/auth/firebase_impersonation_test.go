package auth

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	entydad "github.com/erniealice/entydad-golang"
	login02mod "github.com/erniealice/entydad-golang/service/auth/views/login02"
	"github.com/erniealice/pyeza-golang/view"
)

type authTestRoutes struct {
	views    map[string]view.View
	handlers map[string]http.HandlerFunc
}

type authTestSessionManager struct{}

func (authTestSessionManager) SetSessionCookie(http.ResponseWriter, string) {}
func (authTestSessionManager) ClearSessionCookie(http.ResponseWriter)       {}

func newAuthTestRoutes() *authTestRoutes {
	return &authTestRoutes{views: make(map[string]view.View), handlers: make(map[string]http.HandlerFunc)}
}

func (r *authTestRoutes) GET(path string, v view.View, _ ...string) {
	r.views[http.MethodGet+" "+path] = v
}

func (r *authTestRoutes) POST(path string, v view.View, _ ...string) {
	r.views[http.MethodPost+" "+path] = v
}

func (r *authTestRoutes) HandleFunc(method, path string, handler http.HandlerFunc, _ ...string) {
	r.handlers[method+" "+path] = handler
}

func completeFirebaseImpersonationDeps() *Deps {
	return &Deps{
		SessionManager: authTestSessionManager{},
		UserIDByEmail: func(_ context.Context, email string) string {
			if strings.EqualFold(email, "dev@example.com") {
				return "user-123"
			}
			return ""
		},
		CSRFIssuer: func(http.ResponseWriter, []byte, string, string) string { return "csrf" },
		FirebaseVerifier: func(context.Context, string) (string, string, error) {
			return "dev@example.com", "custom", nil
		},
		SessionMinter: func(context.Context, string) (string, error) { return "session", nil },
		FirebaseCustomTokenMinter: func(context.Context, string) (string, error) {
			return "custom-token", nil
		},
		AllowFirebaseImpersonation: true,
		FirebaseWebConfig: &FirebaseWebConfig{
			APIKey:     "public-api-key",
			AuthDomain: "example.firebaseapp.com",
			ProjectID:  "example",
		},
	}
}

func TestRegisterRoutesFirebaseImpersonationGate(t *testing.T) {
	t.Run("complete local Firebase dependencies mount route and render URL", func(t *testing.T) {
		routes := newAuthTestRoutes()
		NewAuthModule(completeFirebaseImpersonationDeps()).RegisterRoutes(routes)

		key := http.MethodPost + " " + entydad.AuthFirebaseImpersonationURL
		if routes.handlers[key] == nil {
			t.Fatalf("%s was not registered", key)
		}

		loginView := routes.views[http.MethodGet+" "+entydad.AuthLoginURL]
		if loginView == nil {
			t.Fatal("login view was not registered")
		}
		result := loginView.Handle(context.Background(), &view.ViewContext{
			Request: httptest.NewRequest(http.MethodGet, entydad.AuthLoginURL, nil),
		})
		data, ok := result.Data.(*login02mod.PageData)
		if !ok || data.FirebaseConfig == nil {
			t.Fatalf("login PageData FirebaseConfig = %#v", result.Data)
		}
		if got := data.FirebaseConfig.ImpersonationPostURL; got != entydad.AuthFirebaseImpersonationURL {
			t.Fatalf("ImpersonationPostURL = %q, want %q", got, entydad.AuthFirebaseImpersonationURL)
		}
	})

	tests := []struct {
		name   string
		mutate func(*Deps)
	}{
		{name: "allow gate false", mutate: func(d *Deps) { d.AllowFirebaseImpersonation = false }},
		{name: "token minter missing", mutate: func(d *Deps) { d.FirebaseCustomTokenMinter = nil }},
		{name: "web config missing", mutate: func(d *Deps) { d.FirebaseWebConfig = nil }},
		{name: "ID token verifier missing", mutate: func(d *Deps) { d.FirebaseVerifier = nil }},
		{name: "session minter missing", mutate: func(d *Deps) { d.SessionMinter = nil }},
		{name: "session manager missing", mutate: func(d *Deps) { d.SessionManager = nil }},
		{name: "user lookup missing", mutate: func(d *Deps) { d.UserIDByEmail = nil }},
		{name: "CSRF issuer missing", mutate: func(d *Deps) { d.CSRFIssuer = nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := completeFirebaseImpersonationDeps()
			tt.mutate(deps)
			routes := newAuthTestRoutes()
			NewAuthModule(deps).RegisterRoutes(routes)
			if routes.handlers[http.MethodPost+" "+entydad.AuthFirebaseImpersonationURL] != nil {
				t.Fatal("impersonation route registered with incomplete/disabled dependencies")
			}
		})
	}

	t.Run("incomplete local dependencies also keep custom ID tokens denied", func(t *testing.T) {
		deps := completeFirebaseImpersonationDeps()
		deps.FirebaseCustomTokenMinter = nil
		routes := newAuthTestRoutes()
		NewAuthModule(deps).RegisterRoutes(routes)

		handler := routes.handlers[http.MethodPost+" "+entydad.AuthFirebaseLoginURL]
		if handler == nil {
			t.Fatal("Firebase login route was not registered")
		}
		form := url.Values{"id_token": {"verified-custom-id-token"}}
		request := httptest.NewRequest(http.MethodPost, entydad.AuthFirebaseLoginURL, strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		handler(response, request)
		if response.Code != http.StatusForbidden || !strings.Contains(response.Body.String(), "method_not_allowed") {
			t.Fatalf("status/body = %d %s, want 403 method_not_allowed", response.Code, response.Body.String())
		}
	})
}

func TestHandleFirebaseImpersonation(t *testing.T) {
	newRequest := func(email string) *http.Request {
		form := url.Values{"email": {email}}
		request := httptest.NewRequest(http.MethodPost, entydad.AuthFirebaseImpersonationURL, strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return request
	}

	t.Run("returns no-store token for existing DB and Firebase account", func(t *testing.T) {
		const customToken = "sensitive-custom-token"
		var gotIdentifier string
		deps := completeFirebaseImpersonationDeps()
		deps.UserIDByEmail = func(_ context.Context, email string) string {
			if strings.EqualFold(email, "dev@example.com") {
				return "user-123"
			}
			return ""
		}
		deps.FirebaseCustomTokenMinter = func(_ context.Context, identifier string) (string, error) {
			gotIdentifier = identifier
			return customToken, nil
		}

		var logs bytes.Buffer
		originalWriter := log.Writer()
		log.SetOutput(&logs)
		defer log.SetOutput(originalWriter)

		response := httptest.NewRecorder()
		NewAuthModule(deps).handleFirebaseImpersonation()(response, newRequest("dev@example.com"))

		if response.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", response.Code, response.Body.String())
		}
		if got := response.Header().Get("Cache-Control"); got != "no-store" {
			t.Fatalf("Cache-Control = %q, want no-store", got)
		}
		if gotIdentifier != "dev@example.com" {
			t.Fatalf("minter identifier = %q, want dev@example.com", gotIdentifier)
		}
		if !strings.Contains(response.Body.String(), `"custom_token":"`+customToken+`"`) {
			t.Fatalf("response body does not contain custom token: %s", response.Body.String())
		}
		if strings.Contains(logs.String(), customToken) {
			t.Fatalf("custom token leaked to logs: %s", logs.String())
		}
	})

	t.Run("fails closed when disabled", func(t *testing.T) {
		deps := completeFirebaseImpersonationDeps()
		deps.AllowFirebaseImpersonation = false
		response := httptest.NewRecorder()
		NewAuthModule(deps).handleFirebaseImpersonation()(response, newRequest("dev@example.com"))
		if response.Code != http.StatusNotFound || !strings.Contains(response.Body.String(), "impersonation_not_enabled") {
			t.Fatalf("status/body = %d %s", response.Code, response.Body.String())
		}
	})

	t.Run("rejects missing DB account before minting", func(t *testing.T) {
		calls := 0
		deps := completeFirebaseImpersonationDeps()
		deps.UserIDByEmail = func(context.Context, string) string { return "" }
		deps.FirebaseCustomTokenMinter = func(context.Context, string) (string, error) {
			calls++
			return "token", nil
		}
		response := httptest.NewRecorder()
		NewAuthModule(deps).handleFirebaseImpersonation()(response, newRequest("missing@example.com"))
		if response.Code != http.StatusForbidden || calls != 0 {
			t.Fatalf("status=%d minter calls=%d, want 403/0", response.Code, calls)
		}
	})

	t.Run("hides provider error", func(t *testing.T) {
		deps := completeFirebaseImpersonationDeps()
		deps.UserIDByEmail = func(context.Context, string) string { return "user-123" }
		deps.FirebaseCustomTokenMinter = func(context.Context, string) (string, error) {
			return "", errors.New("private signer detail")
		}
		response := httptest.NewRecorder()
		NewAuthModule(deps).handleFirebaseImpersonation()(response, newRequest("dev@example.com"))
		if response.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", response.Code)
		}
		if strings.Contains(response.Body.String(), "private signer detail") || !strings.Contains(response.Body.String(), "impersonation_failed") {
			t.Fatalf("provider error was not hidden: %s", response.Body.String())
		}
	})
}

func TestFirebaseSignInMethodAllowed(t *testing.T) {
	tests := []struct {
		name               string
		allowed            []string
		method             string
		allowImpersonation bool
		want               bool
	}{
		{name: "custom denied by default", method: "custom", want: false},
		{name: "custom allowed by local gate", method: "custom", allowImpersonation: true, want: true},
		{name: "configured custom still denied without gate", allowed: []string{"custom"}, method: "custom", want: false},
		{name: "configured federated method allowed", allowed: []string{"microsoft.com"}, method: "microsoft.com", want: true},
		{name: "unconfigured method denied", allowed: []string{"password"}, method: "google.com", want: false},
		{name: "empty list preserves existing any non-custom behavior", method: "google.com", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firebaseSignInMethodAllowed(tt.allowed, tt.method, tt.allowImpersonation); got != tt.want {
				t.Fatalf("firebaseSignInMethodAllowed(%v, %q, %t) = %t, want %t", tt.allowed, tt.method, tt.allowImpersonation, got, tt.want)
			}
		})
	}
}

package entydad

import (
	"os"
	"strings"
)

// PasswordChangeEnabled reports whether the self-service change-password
// capability (the /auth/change-password page + endpoint, and the account-page
// link) should be available.
//
// Rationale: a FEDERATED sign-in (e.g. Microsoft via Firebase) has no local
// password to change — the credential lives at the external IdP — so exposing
// "change password" is misleading and the endpoint can't do anything. A Firebase
// deployment that offers the `password` method, and the legacy password
// provider, DO have a local password, so it is exposed.
//
// Resolution order:
//  1. explicit env override AUTH_FIREBASE_PASSWORD_CHANGE_ENABLED (true/false),
//  2. else derived: enabled UNLESS this is a firebase deployment whose
//     AUTH_FIREBASE_ALLOWED_SIGN_IN_METHODS is federated-only (no `password`).
func PasswordChangeEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_FIREBASE_PASSWORD_CHANGE_ENABLED"))) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	}
	// Derived default. Non-firebase providers (password / mock) always have a
	// local password.
	if strings.ToLower(strings.TrimSpace(os.Getenv("CONFIG_AUTH_PROVIDER"))) != "firebase" {
		return true
	}
	methods := strings.TrimSpace(os.Getenv("AUTH_FIREBASE_ALLOWED_SIGN_IN_METHODS"))
	if methods == "" {
		return true // firebase email/password-only deployment
	}
	for _, m := range strings.Split(methods, ",") {
		if strings.ToLower(strings.TrimSpace(m)) == "password" {
			return true // the firebase `password` method is offered
		}
	}
	return false // federated-only firebase → no local password to change
}

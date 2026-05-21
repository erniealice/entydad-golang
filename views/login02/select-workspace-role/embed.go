// Package selectWorkspaceRole renders the post-login principal chooser page —
// shown when a User holds 2+ active principal bindings (e.g., an operator
// who is also a client of the same workspace).
//
// The page deliberately reuses the login02 auth shell (no app shell, no
// sidebar) so the visual flow from /auth/login → /auth/select-workspace-role →
// /portal/{kind}/ stays seamless.
//
// Route convention:
//
//	GET  /auth/select-workspace-role     — render this page
//	POST /action/auth/switch-principal   — handle card click; rotates session
//
// Security contract: this page is reached AFTER credential verification but
// BEFORE a principal-scoped session is established. The intermediate
// session is principal-less; route handlers must reject any
// permission-requiring access until the user picks one. See
// docs/plan/20260516-self-domain/codex-design-supplement.md §1 for the full
// reasoning around the pre-principal session shape.
package selectWorkspaceRole

import "embed"

// TemplatesFS embeds the select-workspace-role templates. Service-admin
// registers this FS in its renderer alongside login02.TemplatesFS.
//
//go:embed templates/*.html
var TemplatesFS embed.FS

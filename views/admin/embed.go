// Package admin is the entydad view module for the admin app.
//
// Admin is a *composite* sidebar app that spans 5 entities (permission,
// role, workspace, workspace_user, workspace_user_role). Its dashboard sits
// at the app level — not per entity — and lives under views/admin/dashboard.
// The constituent entities keep their own view modules (views/permission,
// views/role, views/workspace, views/workspace_user, views/workspace_user_role)
// for CRUD; this package is dashboard-only by design.
package admin

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS

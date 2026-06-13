package block

import (
	"context"
	"database/sql"
	"net/http"

	entydad "github.com/erniealice/entydad-golang"
	entityclient "github.com/erniealice/entydad-golang/domain/entity/party/client"
	entitysupplier "github.com/erniealice/entydad-golang/domain/entity/party/supplier"
	entityrole "github.com/erniealice/entydad-golang/domain/entity/identity/role"
	entityuser "github.com/erniealice/entydad-golang/domain/entity/identity/user"
	roleusers "github.com/erniealice/entydad-golang/domain/entity/identity/role/users"
	userdashboard "github.com/erniealice/entydad-golang/domain/entity/identity/user/dashboard"
	workspaceaction "github.com/erniealice/entydad-golang/domain/entity/identity/workspace/action"
	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/espyna-golang/reference"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
)

// Infra carries the subset of AppContext that view modules need beyond the
// typed UseCases: attachment ops, reference checker, raw DB for registry-
// backed repos (payment_term, category), identity helpers, cross-centymo
// route structs, and the optional secure-workspace-switch closures. Built
// once by service-admin and passed into each catalog binder.
//
// Companion: block/catalog.go — the binder functions that consume Infra.
//
// What lives here vs UseCases:
//   - UseCases: proto-shaped business use-case closures (all nil-safe).
//   - Infra: infrastructure dependencies that are NOT proto use cases —
//     DB access, attachment ops, ref-checker, host-provided helper closures,
//     cross-domain label sets not owned by a single entity descriptor.
type Infra struct {
	// Attachment operations — passed verbatim to entity module Deps.
	UploadFile       func(context.Context, string, string, []byte, string) error
	ListAttachments  func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewAttachmentID  func() string

	// RefChecker provides in-use-ID gating for deletable entities.
	RefChecker reference.Checker

	// SqlDB is the raw SQL connection used by registry.CreateRepository for
	// the payment_term and category (client_tag / supplier_tag) repos. Nil
	// means those modules are skipped at mount time (non-fatal, with a log
	// warning).
	SqlDB *sql.DB

	// Identity helpers — passed into the user / role module Deps.
	// All are nil-safe: the relevant UI sections degrade gracefully when unset.
	GetUsersByRoleID     func(ctx context.Context, roleID string) ([]roleusers.UserByRole, error)
	GetDashboardData     func(ctx context.Context) (*userdashboard.DashboardData, error)
	HashPassword         func(password string) (string, error)
	GetUserWorkspacesMap func(ctx context.Context) (map[string][]pyezatypes.ChipData, error)

	// SecureSwitch — optional secure workspace-rotation primitive for the
	// switch-workspace handler. When non-nil (together with its two siblings),
	// WorkspaceUnit's Mount registers the rotation-aware path instead of the
	// legacy in-place SwitchWorkspace use case. All three fields are required
	// together; a partial set (e.g. only SecureSwitch) falls back to the
	// legacy path.
	SecureSwitch            workspaceaction.SecureSwitchFn
	SecureSwitchResolveUser func(r *http.Request) string
	SecureSwitchSetCookie   func(w http.ResponseWriter, token string)

	// Cross-domain label sets — loaded from lyngua by service-admin and
	// passed in because MountContext does not carry them. All are zero-safe:
	// the view modules render with empty strings when unset.

	// SharedLabels holds error/confirm/badge strings shared across all modules.
	SharedLabels entydad.SharedLabels

	// DashboardTitleLabels holds the cross-entity dashboard title strings
	// (client, user, supplier, location, admin section headers).
	DashboardTitleLabels entydad.DashboardLabels

	// ClientDashboardLabels holds the client-specific dashboard tab strings.
	ClientDashboardLabels entityclient.DashboardLabels

	// SupplierDashboardLabels holds the supplier-specific dashboard tab strings.
	SupplierDashboardLabels entitysupplier.DashboardLabels

	// UserDashboardLabels holds the user-specific dashboard widget strings.
	UserDashboardLabels entityuser.DashboardLabels

	// RolePermissionLabels holds the role→permission junction tab strings.
	RolePermissionLabels entityrole.PermissionLabels

	// RoleUserLabels holds the role→user junction tab strings.
	RoleUserLabels entityrole.UserLabels

	// UserRoleLabels holds the user→role junction tab strings.
	UserRoleLabels entityuser.RoleLabels

	// Cross-centymo routes — used by client / supplier detail tabs for
	// subscription and price-schedule deep links. Zero-valued routes
	// produce empty-string URLs, which the view modules treat as "not
	// available" (CTAs hidden, tabs render empty-state).
	SubscriptionRoutes  centymo.SubscriptionRoutes
	PriceScheduleRoutes centymo.PriceScheduleRoutes

	// HomeURLForWorkspaceID resolves the post-switch redirect URL given the
	// newly-active workspace_id (used by WorkspaceUnit's switch handler).
	// Nil-safe: falls back to "/home".
	HomeURLForWorkspaceID func(ctx context.Context, workspaceID string) string
}

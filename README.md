# entydad-golang

**Module:** `github.com/erniealice/entydad-golang`

Identity and entity management domain package for the Ryta OS platform. Provides reusable views, templates, labels, routes, and HTMX helpers for managing users, clients, roles, permissions, locations, workspaces, suppliers, and client tags. Consumer apps (retail-admin, service-admin) wire these views into their composition roots.

## Package Structure

```
packages/entydad-golang-ryta/
├── assets/
│   └── css/
│       ├── entydad-client-dashboard.css
│       ├── entydad-client-detail.css
│       ├── entydad-location-detail.css
│       ├── entydad-role-detail.css
│       ├── entydad-supplier-detail.css
│       ├── entydad-user-dashboard.css
│       ├── entydad-user-detail.css
│       ├── login01.css
│       └── login02.css
├── views/
│   ├── client/
│   │   ├── action/action.go          # CRUD + status + bulk actions
│   │   ├── dashboard/page.go         # Client dashboard view
│   │   ├── detail/page.go            # Client detail + tab action views
│   │   ├── list/page.go              # Client list + table-only views
│   │   ├── tag/
│   │   │   ├── action/action.go      # Client tag CRUD actions
│   │   │   └── page.go               # Client tag list view
│   │   ├── templates/
│   │   │   ├── client-drawer-form.html
│   │   │   ├── dashboard.html
│   │   │   ├── detail.html
│   │   │   ├── list.html
│   │   │   ├── tag-drawer-form.html
│   │   │   └── tag-list.html
│   │   └── embed.go                  # //go:embed templates/*.html
│   ├── user/
│   │   ├── action/action.go          # CRUD + status + password reset
│   │   ├── dashboard/page.go         # User dashboard (stats, chart, activity)
│   │   ├── detail/page.go            # User detail + tab action + roles table
│   │   ├── list/page.go              # User list + table-only views
│   │   ├── roles/
│   │   │   ├── page.go               # User-role assignment list
│   │   │   └── action.go             # User-role assign/remove actions
│   │   ├── templates/
│   │   │   ├── dashboard.html
│   │   │   ├── detail.html
│   │   │   ├── list.html
│   │   │   ├── user-drawer-form.html
│   │   │   ├── user-role-assign-form.html
│   │   │   └── user-roles.html
│   │   └── embed.go
│   ├── role/
│   │   ├── action/action.go          # Role CRUD + status actions
│   │   ├── detail/page.go            # Role detail (info/permissions/users tabs)
│   │   ├── list/page.go              # Role list + table-only views
│   │   ├── permissions/
│   │   │   ├── page.go               # Role-permission assignment list
│   │   │   └── action.go             # Role-permission assign/remove actions
│   │   ├── users/
│   │   │   ├── page.go               # Role-user assignment list
│   │   │   └── action.go             # Role-user assign/remove actions
│   │   ├── templates/
│   │   │   ├── detail.html
│   │   │   ├── list.html
│   │   │   ├── role-drawer-form.html
│   │   │   ├── role-permission-assign-form.html
│   │   │   ├── role-permissions.html
│   │   │   ├── role-user-assign-form.html
│   │   │   └── role-users.html
│   │   └── embed.go
│   ├── location/
│   │   ├── action/action.go          # Location CRUD + status actions
│   │   ├── detail/page.go            # Location detail view
│   │   ├── list/page.go              # Location list + table-only views
│   │   ├── templates/
│   │   │   ├── detail.html
│   │   │   ├── list.html
│   │   │   └── location-drawer-form.html
│   │   └── embed.go
│   ├── permission/
│   │   ├── action/action.go          # Permission CRUD + status actions
│   │   ├── list/page.go              # Permission list + table-only views
│   │   ├── templates/
│   │   │   ├── list.html
│   │   │   └── permission-drawer-form.html
│   │   └── embed.go
│   ├── workspace/
│   │   ├── action/action.go          # Workspace CRUD + status actions
│   │   ├── list/page.go              # Workspace list + table-only views
│   │   ├── templates/
│   │   │   ├── list.html
│   │   │   └── workspace-drawer-form.html
│   │   └── embed.go
│   ├── supplier/
│   │   ├── action/action.go          # Supplier CRUD + status actions
│   │   ├── detail/page.go            # Supplier detail view
│   │   ├── list/page.go              # Supplier list + table-only views
│   │   ├── templates/
│   │   │   ├── detail.html
│   │   │   ├── list.html
│   │   │   └── supplier-drawer-form.html
│   │   └── embed.go
│   ├── login01/
│   │   ├── page.go                   # Simple login page view
│   │   ├── action.go                 # Login POST handler (redirect)
│   │   ├── templates/login01.html
│   │   └── embed.go
│   └── login02/
│       ├── page.go                   # Split-screen login with carousel
│       ├── action.go                 # Login POST handler (redirect)
│       ├── templates/login02.html
│       └── embed.go
├── assets.go                         # CopyStyles(), CopyStaticAssets()
├── datasource.go                     # DataSource interface
├── htmx.go                           # HTMXSuccess(), HTMXError()
├── labels.go                         # All label struct types
├── pkgdir.go                         # packageDir() via runtime.Caller(0)
├── routes.go                         # Route constants
├── routes_config.go                  # Route config structs + defaults + RouteMap()
├── go.mod
└── go.sum
```

## Entity Modules

| Entity | List | Detail | Dashboard | Action | Sub-views |
|--------|------|--------|-----------|--------|-----------|
| Client | list, table | detail, tab-action | dashboard | add, edit, delete, bulk-delete, set-status, bulk-set-status | tag (list + CRUD) |
| User | list, table | detail, tab-action | dashboard (stats/chart) | add, edit, delete, bulk-delete, set-status, bulk-set-status, reset-password | roles (assign/remove) |
| Role | list, table | detail (info/permissions/users tabs) | -- | add, edit, delete, bulk-delete, set-status, bulk-set-status | permissions (assign/remove), users (assign/remove) |
| Location | list, table | detail | -- | add, edit, delete, bulk-delete, set-status, bulk-set-status, edit-detail | -- |
| Permission | list, table | -- | -- | add, edit, delete, bulk-delete, set-status, bulk-set-status | -- |
| Workspace | list, table | -- | -- | add, edit, delete, bulk-delete, set-status, bulk-set-status | -- |
| Supplier | list, table | detail | -- | add, edit, delete, bulk-delete, set-status, bulk-set-status | -- |
| Client Tag | list | -- | -- | add, edit, delete, bulk-delete | -- |
| Login01 | -- | -- | -- | login POST | -- |
| Login02 | -- | -- | -- | login POST | -- |

## View Sub-Packages

Every view sub-package follows the same pattern:

1. **`Deps` struct** -- dependency injection container holding routes, labels, and use-case functions
2. **`PageData` struct** -- embeds `types.PageData`, adds content-specific fields
3. **`NewView(deps)` constructor** -- returns a `view.View` (full page render)
4. **`NewTableView(deps)` constructor** (list views only) -- returns table-only HTML for HTMX refresh
5. **Action constructors** (`NewAddAction`, `NewEditAction`, etc.) -- handle GET (form) + POST (submit)

Views accept protobuf use-case functions directly (not interfaces), so the consumer app can wire in any implementation -- espyna-generated use cases, mocks, or custom adapters.

### Client List (`views/client/list`)

```go
type Deps struct {
    Routes          entydad.ClientRoutes
    GetListPageData func(ctx, *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
    GetInUseIDs     func(ctx, ids []string) (map[string]bool, error)
    Labels          entydad.ClientLabels
    SharedLabels    entydad.SharedLabels
    CommonLabels    pyeza.CommonLabels
    TableLabels     types.TableLabels
}
```

- `NewView(deps)` -- full page with table, renders template `"client-list"`
- `NewTableView(deps)` -- table card only, renders template `"table-card"`
- Table ID: `"clients-table"`
- Filters by `{status}` path param (active, prospect, inactive)
- Supports bulk actions (deactivate, activate, delete)

### Client Detail (`views/client/detail`)

```go
type Deps struct {
    Routes               entydad.ClientRoutes
    ReadClient           func(ctx, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
    ListCategories       func(ctx, *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
    ListClientCategories func(ctx, *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
    ListRevenues         func(ctx, collection string) ([]map[string]any, error)
    Labels               entydad.ClientLabels
    CommonLabels         pyeza.CommonLabels
}
```

- `NewView(deps)` -- full detail page, renders template `"client-detail"`
- `NewTabAction(deps)` -- HTMX partial, renders `"client-tab-{tab}"` (basic, history)
- Tabs: basic info (company, personal, address, notes, tags), purchase history
- Loads tags via client_category junction records
- Loads purchase history from revenue records, calculates lifetime spend / avg order / last purchase

### Client Action (`views/client/action`)

```go
type Deps struct {
    Routes          entydad.ClientRoutes
    CreateClient    func(ctx, *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error)
    ReadClient      func(ctx, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
    UpdateClient    func(ctx, *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error)
    DeleteClient    func(ctx, *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error)
    SetClientActive func(ctx, id string, active bool) error
    ListCategories       func(...)   // for tag multi-select options
    ListClientCategories func(...)   // for current tag assignments
    CreateClientCategory func(...)   // for tag sync on save
    DeleteClientCategory func(...)   // for tag sync on save
}
```

Constructors:
- `NewAddAction(deps)` -- GET returns drawer form, POST creates client + syncs tags
- `NewEditAction(deps)` -- GET returns prefilled form, POST updates client + syncs tags
- `NewDeleteAction(deps)` -- POST deletes single client
- `NewBulkDeleteAction(deps)` -- POST deletes multiple clients
- `NewSetStatusAction(deps)` -- POST activate/deactivate single client
- `NewBulkSetStatusAction(deps)` -- POST activate/deactivate multiple clients

Form template: `"client-drawer-form"`. All actions return `HTMXSuccess("clients-table")` or `HTMXError(message)`.

Uses `SetClientActive` (raw map update) instead of `UpdateClient` (protobuf) for status changes because proto3's protojson omits bool fields with value `false`.

### Client Dashboard (`views/client/dashboard`)

```go
type Deps struct {
    DashboardLabels entydad.DashboardLabels
    CommonLabels    pyeza.CommonLabels
}
```

- `NewView(deps)` -- renders template `"client-dashboard"`
- Lightweight placeholder page (no data provider)

### Client Tag (`views/client/tag`)

```go
type Deps struct {
    Routes               entydad.ClientTagRoutes
    ListCategories       func(...)
    ListClientCategories func(...)
    GetInUseIDs          func(ctx, ids []string) (map[string]bool, error)
    RefreshURL           string
    Labels               entydad.ClientTagLabels
    SharedLabels         entydad.SharedLabels
    CommonLabels         pyeza.CommonLabels
    TableLabels          types.TableLabels
}
```

- `NewView(deps)` -- full page, renders template `"client-tag-list"`
- Table ID: `"client-tags-table"`
- Shows customer count per tag via client_category junction records
- Tag actions in `views/client/tag/action/`: add, edit, delete, bulk-delete
- Duplicate name check on add/edit (case-insensitive)
- Auto-generates slug `Code` from tag name if not provided

### User List (`views/user/list`)

```go
type Deps struct {
    Routes          entydad.UserRoutes
    GetListPageData func(ctx, *userpb.GetUserListPageDataRequest) (*userpb.GetUserListPageDataResponse, error)
    GetUserRolesMap func(ctx) (map[string][]entydad.RoleBadge, error)
    RefreshURL      string
    Labels          entydad.UserLabels
    SharedLabels    entydad.SharedLabels
    CommonLabels    pyeza.CommonLabels
    TableLabels     types.TableLabels
}
```

- `NewView(deps)` / `NewTableView(deps)`
- Table ID: `"users-table"`
- Displays role chips per user via `GetUserRolesMap` (max 2 visible + "+N more")

### User Detail (`views/user/detail`)

```go
type Deps struct {
    Routes                       entydad.UserRoutes
    ReadUser                     func(...)
    GetWorkspaceUserItemPageData func(...)
    ListWorkspaceUsers           func(...)
    Labels                       entydad.UserLabels
    SharedLabels                 entydad.SharedLabels
    UserRoleLabels               entydad.UserRoleLabels
    CommonLabels                 pyeza.CommonLabels
    TableLabels                  types.TableLabels
}
```

- `NewView(deps)` / `NewTabAction(deps)` (returns `"user-tab-{tab}"`)
- Tabs: info, roles, security, audit trail
- Roles tab renders an embedded table (`"user-roles-table"`) with assign/remove actions
- Role count shown as tab badge

### User Dashboard (`views/user/dashboard`)

```go
type Deps struct {
    DashboardLabels  entydad.DashboardLabels
    CommonLabels     pyeza.CommonLabels
    GetDashboardData func(ctx) (*DashboardData, error)
}

type DashboardData struct {
    Stats          DashboardStats    // TotalUsers, ActiveUsers, InactiveUsers, TotalRoles
    RecentActivity []ActivityItem    // IconHTML, Title, Description, TimeAgo
    Chart          ChartData         // Labels, Values, Period
}
```

- Generates SVG chart paths from data values for inline rendering
- Supports period switching (year, month, etc.)

### User Action (`views/user/action`)

```go
type Deps struct {
    Routes              entydad.UserRoutes
    CreateUser          func(...)
    ReadUser            func(...)
    UpdateUser          func(...)
    DeleteUser          func(...)
    SetUserActive       func(ctx, id string, active bool) error
    CreateWorkspaceUser func(...)                              // auto-create on user creation
    DefaultWorkspaceID  string                                 // target workspace for auto-create
    HashPassword        func(password string) (string, error)  // optional hasher
}
```

Additional constructors:
- `NewResetPasswordAction(deps)` -- reads existing user, hashes new password, updates
- `NewBulkSetStatusAction(deps)` -- bulk activate/deactivate

### Role List (`views/role/list`)

```go
type Deps struct {
    GetListPageData func(ctx, *rolepb.GetRoleListPageDataRequest) (*rolepb.GetRoleListPageDataResponse, error)
    GetInUseIDs     func(ctx, ids []string) (map[string]bool, error)
    RefreshURL      string
    Routes          entydad.RoleRoutes
    Labels          entydad.RoleLabels
    SharedLabels    entydad.SharedLabels
    CommonLabels    pyeza.CommonLabels
    TableLabels     types.TableLabels
}
```

- Table ID: `"roles-table"`
- Shows permission count per role as badge

### Role Detail (`views/role/detail`)

```go
type Deps struct {
    ReadRole            func(...)
    RoleGetItemPageData func(...)
    GetUsersByRoleID    func(ctx, roleID string) ([]UserByRole, error)
    Routes              entydad.RoleRoutes
    Labels              entydad.RoleLabels
    SharedLabels        entydad.SharedLabels
    RolePermissionLabels entydad.RolePermissionLabels
    RoleUserLabels      entydad.RoleUserLabels
    CommonLabels        pyeza.CommonLabels
    TableLabels         types.TableLabels
}
```

- Tabs: info, permissions, users
- Permissions tab: shows assigned permissions with assign/remove actions
- Users tab: shows assigned users with assign/remove actions (via `GetUsersByRoleID`)

### Role Permissions (`views/role/permissions`)

```go
type Deps struct {
    GetRoleItemPageData func(...)
    Routes              entydad.RoleRoutes
    Labels              entydad.RolePermissionLabels
    SharedLabels        entydad.SharedLabels
    CommonLabels        pyeza.CommonLabels
    TableLabels         types.TableLabels
}
```

### Role Users (`views/role/users`)

```go
type UserByRole struct {
    WorkspaceUserRoleID string
    WorkspaceUserID     string
    UserID              string
    UserName            string
    Email               string
    DateAssigned        string
}

type Deps struct {
    GetUsersByRoleID func(ctx, roleID string) ([]UserByRole, error)
    Routes           entydad.RoleRoutes
    Labels           entydad.RoleUserLabels
    SharedLabels     entydad.SharedLabels
    CommonLabels     pyeza.CommonLabels
    TableLabels      types.TableLabels
}
```

### User Roles (`views/user/roles`)

```go
type Deps struct {
    Routes                       entydad.UserRoutes
    ListWorkspaceUsers           func(...)
    GetWorkspaceUserItemPageData func(...)
    ReadUser                     func(...)
    Labels                       entydad.UserRoleLabels
    SharedLabels                 entydad.SharedLabels
    CommonLabels                 pyeza.CommonLabels
    TableLabels                  types.TableLabels
}
```

### Location List (`views/location/list`)

```go
type Deps struct {
    GetListPageData func(ctx, *locationpb.GetLocationListPageDataRequest) (*locationpb.GetLocationListPageDataResponse, error)
    GetInUseIDs     func(ctx, ids []string) (map[string]bool, error)
    RefreshURL      string
    Routes          entydad.LocationRoutes
    Labels          entydad.LocationLabels
    SharedLabels    entydad.SharedLabels
    CommonLabels    pyeza.CommonLabels
    TableLabels     types.TableLabels
}
```

### Location Detail (`views/location/detail`)

```go
type Deps struct {
    Routes       entydad.LocationRoutes
    ReadLocation func(ctx, *locationpb.ReadLocationRequest) (*locationpb.ReadLocationResponse, error)
    Labels       entydad.LocationLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```

- Tab labels: info, users, price lists, audit trail

### Location Action (`views/location/action`)

```go
type Deps struct {
    CreateLocation    func(...)
    ReadLocation      func(...)
    UpdateLocation    func(...)
    DeleteLocation    func(...)
    SetLocationActive func(ctx, id string, active bool) error
    GetInUseIDs       func(ctx, ids []string) (map[string]bool, error)
    Routes            entydad.LocationRoutes
}
```

### Permission List (`views/permission/list`)

```go
type Deps struct {
    GetListPageData func(ctx, *permissionpb.GetPermissionListPageDataRequest) (*permissionpb.GetPermissionListPageDataResponse, error)
    RefreshURL      string
    Routes          entydad.PermissionRoutes
    Labels          entydad.PermissionLabels
    SharedLabels    entydad.SharedLabels
    CommonLabels    pyeza.CommonLabels
    TableLabels     types.TableLabels
}
```

### Permission Action (`views/permission/action`)

```go
type Deps struct {
    CreatePermission    func(...)
    ReadPermission      func(...)
    UpdatePermission    func(...)
    DeletePermission    func(...)
    SetPermissionActive func(ctx, id string, active bool) error
    Routes              entydad.PermissionRoutes
}
```

- Form includes `PermissionTypeOptions []types.SelectOption` for permission type dropdown

### Workspace List (`views/workspace/list`)

```go
type Deps struct {
    GetListPageData func(ctx, *workspacepb.GetWorkspaceListPageDataRequest) (*workspacepb.GetWorkspaceListPageDataResponse, error)
    RefreshURL      string
    Routes          entydad.WorkspaceRoutes
    Labels          entydad.WorkspaceLabels
    SharedLabels    entydad.SharedLabels
    CommonLabels    pyeza.CommonLabels
    TableLabels     types.TableLabels
}
```

### Workspace Action (`views/workspace/action`)

```go
type Deps struct {
    CreateWorkspace    func(...)
    ReadWorkspace      func(...)
    UpdateWorkspace    func(...)
    DeleteWorkspace    func(...)
    SetWorkspaceActive func(ctx, id string, active bool) error
    Routes             entydad.WorkspaceRoutes
}
```

### Supplier List (`views/supplier/list`)

```go
type Deps struct {
    Routes          entydad.SupplierRoutes
    GetListPageData func(ctx, *supplierpb.GetSupplierListPageDataRequest) (*supplierpb.GetSupplierListPageDataResponse, error)
    GetInUseIDs     func(ctx, ids []string) (map[string]bool, error)
    Labels          entydad.SupplierLabels
    SharedLabels    entydad.SharedLabels
    CommonLabels    pyeza.CommonLabels
    TableLabels     types.TableLabels
}
```

- Status tabs: active, blocked, on-hold (not active/inactive like other entities)

### Supplier Detail (`views/supplier/detail`)

```go
type Deps struct {
    Routes       entydad.SupplierRoutes
    ReadSupplier func(ctx, *supplierpb.ReadSupplierRequest) (*supplierpb.ReadSupplierResponse, error)
    Labels       entydad.SupplierLabels
    CommonLabels pyeza.CommonLabels
}
```

- Sections: company info, contact info, financial info, address info

### Supplier Action (`views/supplier/action`)

```go
type Deps struct {
    CreateSupplier    func(...)
    ReadSupplier      func(...)
    UpdateSupplier    func(...)
    DeleteSupplier    func(...)
    SetSupplierStatus func(ctx, id string, status string) error
    Routes            entydad.SupplierRoutes
}
```

- Includes contact person fields (first/last name, email, phone)
- Financial fields: default currency, payment terms, lead time, credit limit

### Login Views

#### login01 (Simple)

```go
type Deps struct {
    Labels       entydad.LoginLabels
    CommonLabels pyeza.CommonLabels
    RedirectURL  string  // default: /app/
}
```

- `NewView(deps)` -- renders `"login01"` template
- `NewAction(&ActionDeps{RedirectURL})` -- POST handler, redirects on success
- Includes admin login and staff PIN login sections

#### login02 (Split-Screen)

```go
type Deps struct {
    Labels          entydad.Login02Labels
    CommonLabels    pyeza.CommonLabels
    RedirectURL     string            // default: /app/
    LogoText        string            // brand name on form side
    LogoIcon        string            // icon template name
    LoginPostURL    string            // default: /login
    RegisterURL     string            // default: /register
    ForgotURL       string            // default: /auth/reset-password
    Slides          []CarouselSlide   // left panel carousel
    SocialProviders []SocialProvider  // social login buttons
}

type CarouselSlide struct {
    Title       string
    Description string
}

type SocialProvider struct {
    Name    string  // e.g. "Google", "Apple"
    IconSVG string  // raw SVG markup
}
```

- `NewView(deps)` -- renders `"login02"` template with split-screen layout
- `NewAction(&ActionDeps{RedirectURL})` -- POST handler, redirects on success
- Features carousel slides on left panel and optional social login buttons

## Template System

Each view sub-package embeds its own templates via `//go:embed templates/*.html` in an `embed.go` file, exposing a `TemplatesFS embed.FS` variable. The pyeza renderer discovers these embedded filesystems when the consumer app registers them.

### Template Names Referenced in Go Code

| Template Name | View Package | Used By |
|---------------|-------------|---------|
| `client-list` | client/list | `NewView` |
| `table-card` | (pyeza shared) | `NewTableView` (all list views) |
| `client-detail` | client/detail | `NewView` |
| `client-tab-basic` | client/detail | `NewTabAction` |
| `client-tab-history` | client/detail | `NewTabAction` |
| `client-dashboard` | client/dashboard | `NewView` |
| `client-drawer-form` | client/action | `NewAddAction`, `NewEditAction` |
| `client-tag-list` | client/tag | `NewView` |
| `client-tag-drawer-form` | client/tag/action | `NewAddAction`, `NewEditAction` |
| `user-list` | user/list | `NewView` |
| `user-detail` | user/detail | `NewView` |
| `user-tab-{tab}` | user/detail | `NewTabAction` (info, roles, security, audit) |
| `user-dashboard` | user/dashboard | `NewView` |
| `user-drawer-form` | user/action | `NewAddAction`, `NewEditAction` |
| `user-role-assign-form` | user/roles | assign action |
| `user-roles` | user/roles | `NewView` |
| `role-list` | role/list | `NewView` |
| `role-detail` | role/detail | `NewView` |
| `role-drawer-form` | role/action | `NewAddAction`, `NewEditAction` |
| `role-permissions` | role/permissions | `NewView` |
| `role-permission-assign-form` | role/permissions | assign action |
| `role-users` | role/users | `NewView` |
| `role-user-assign-form` | role/users | assign action |
| `location-list` | location/list | `NewView` |
| `location-detail` | location/detail | `NewView` |
| `location-drawer-form` | location/action | `NewAddAction`, `NewEditAction` |
| `permission-list` | permission/list | `NewView` |
| `permission-drawer-form` | permission/action | `NewAddAction`, `NewEditAction` |
| `workspace-list` | workspace/list | `NewView` |
| `workspace-drawer-form` | workspace/action | `NewAddAction`, `NewEditAction` |
| `supplier-list` | supplier/list | `NewView` |
| `supplier-detail` | supplier/detail | `NewView` |
| `supplier-drawer-form` | supplier/action | `NewAddAction`, `NewEditAction` |
| `login01` | login01 | `NewView` |
| `login02` | login02 | `NewView` |

## Route Constants

All route constants are defined in `routes.go`. They follow the pattern:
- **Page routes:** `GET /app/{entity}/...` (full page renders)
- **Action routes:** `GET|POST /action/{entity}/...` (HTMX partials and form handlers)

### Client Routes
```
/app/clients/dashboard
/app/clients/list/{status}
/app/clients/detail/{id}
/action/clients/table/{status}
/action/clients/add
/action/clients/edit/{id}
/action/clients/delete
/action/clients/bulk-delete
/action/clients/set-status
/action/clients/bulk-set-status
/action/clients/{id}/tab/{tab}
/app/clients/reports/receivables-aging
```

### User Routes
```
/app/users/dashboard
/app/users/list/{status}
/app/users/detail/{id}
/action/users/table/{status}
/action/users/add
/action/users/edit/{id}
/action/users/delete
/action/users/bulk-delete
/action/users/set-status
/action/users/bulk-set-status
/action/users/{id}/tab/{tab}
/action/users/reset-password/{id}
/app/users/detail/{id}/roles          (migrated from /manage/)
/action/users/detail/{id}/roles/table
/action/users/detail/{id}/roles/assign
/action/users/detail/{id}/roles/remove
```

### Role Routes
```
/app/roles/list
/app/roles/detail/{id}
/action/roles/table
/action/roles/add
/action/roles/edit/{id}
/action/roles/delete
/action/roles/bulk-delete
/action/roles/set-status
/action/roles/bulk-set-status
/action/roles/{id}/tab/{tab}
/app/roles/detail/{id}/permissions
/action/roles/detail/{id}/permissions/table
/action/roles/detail/{id}/permissions/assign
/action/roles/detail/{id}/permissions/remove
/app/roles/detail/{id}/users
/action/roles/detail/{id}/users/table
/action/roles/detail/{id}/users/assign
/action/roles/detail/{id}/users/remove
```

### Location Routes
```
/app/locations/list/{status}
/app/locations/{id}
/action/locations/table/{status}
/action/locations/add
/action/locations/edit/{id}
/action/locations/delete
/action/locations/bulk-delete
/action/locations/set-status
/action/locations/bulk-set-status
/action/locations/{id}/tab/{tab}
/action/locations/edit-detail/{id}
```

### Permission Routes
```
/app/permissions/list/{status}
/action/permissions/table/{status}
/action/permissions/add
/action/permissions/edit/{id}
/action/permissions/delete
/action/permissions/bulk-delete
/action/permissions/set-status
/action/permissions/bulk-set-status
```

### Workspace Routes
```
/app/workspaces/list/{status}
/action/workspaces/table/{status}
/action/workspaces/add
/action/workspaces/edit/{id}
/action/workspaces/delete
/action/workspaces/bulk-delete
/action/workspaces/set-status
/action/workspaces/bulk-set-status
```

### Supplier Routes
```
/app/suppliers/list/{status}
/app/suppliers/detail/{id}
/action/suppliers/table/{status}
/action/suppliers/add
/action/suppliers/edit/{id}
/action/suppliers/delete
/action/suppliers/bulk-delete
/action/suppliers/set-status
/action/suppliers/bulk-set-status
/action/suppliers/{id}/tab/{tab}
/app/suppliers/reports/payables-aging
```

### Client Tag Routes
```
/app/clients/settings/tags/list
/action/clients/tags/add
/action/clients/tags/edit/{id}
/action/clients/tags/delete
/action/clients/tags/bulk-delete
```

### Login Routes
```
GET  /login
POST /login
```

### Default Redirect
```
/app/  -->  (consumer app defines target, default: /app/)
```

## Route Config Structs

Each entity has a route config struct with:
- JSON tags for industry-specific overrides via lyngua
- `Default{Entity}Routes()` constructor populating from constants
- `RouteMap() map[string]string` method for template/reverse-routing lookups

Available structs: `ClientRoutes`, `UserRoutes`, `RoleRoutes`, `LocationRoutes`, `PermissionRoutes`, `WorkspaceRoutes`, `SupplierRoutes`, `ClientTagRoutes`, `LoginRoutes`.

### Three-Level Routing System

1. **Level 1 -- Generic defaults** from Go constants (`DefaultXxxRoutes()`)
2. **Level 2 -- Industry-specific overrides** via JSON (consumer apps unmarshal into route structs)
3. **Level 3 -- App-specific overrides** via direct Go field assignment

### RouteMap Example

```go
routes := entydad.DefaultClientRoutes()
rm := routes.RouteMap()
// rm["client.list"]    = "/app/clients/list/{status}"
// rm["client.detail"]  = "/app/clients/detail/{id}"
// rm["client.add"]     = "/action/clients/add"
```

## Labels

All label types are defined in `labels.go`. Each entity has a top-level label struct containing nested sub-structs for page, buttons, columns, empty states, form fields, actions, and detail sections. All structs are JSON-tagged to match lyngua translation file structure.

### Label Struct Hierarchy

| Top-Level Struct | Sub-Structs |
|-----------------|-------------|
| `ClientLabels` | Page (8 fields), Buttons (1), Columns (1), Empty (6), Form (20), Detail (CompanyDetails, Actions, Profile, Company, Personal, Address, NotesSection, Tags, PurchaseHistory, Tabs), BulkActions (1) |
| `UserLabels` | Page (6), Buttons (1), Columns (4), Empty (4), Form (1), Actions (6), Detail (BasicInfo 14 fields, Tabs 4, Security 7, EmptyStates 2) |
| `LocationLabels` | Page (6), Buttons (1), Columns (3), Empty (4), Form (7), Actions (5), Detail (BasicInfo 8, Tabs 4, EmptyStates 6) |
| `RoleLabels` | Page (6), Buttons (1), Columns (5), Empty (4), Form (7), Actions (6), Detail (Tabs 3, Info 5) |
| `PermissionLabels` | Page (6), Buttons (1), Columns (5), Empty (4), Form (7), Actions (5) |
| `RolePermissionLabels` | Page (2), Buttons (1), Columns (4), Empty (2), Form (1), Actions (3) |
| `UserRoleLabels` | Page (2), Buttons (1), Columns (4), Empty (2), Form (1), Actions (3) |
| `RoleUserLabels` | Page (2), Buttons (1), Columns (3), Empty (2), Form (2), Actions (2) |
| `WorkspaceLabels` | Page (6), Buttons (1), Columns (4), Empty (4), Form (6), Actions (5) |
| `SupplierLabels` | Page (8), Buttons (1), Columns (7), Empty (6), Form (20), Detail (InfoTab, CompanyInfo, ContactInfo, FinancialInfo, AddressInfo), Actions (6) |
| `ClientTagLabels` | Page (2), Buttons (1), Columns (4), Empty (2), Actions (2), Confirm (3) |
| `LoginLabels` | 11 flat fields |
| `Login02Labels` | 14 flat fields |
| `DashboardLabels` | 2 flat fields (ClientTitle, UserTitle) |

### SharedLabels

Used across all entydad modules. Loaded separately from common translation files:

```go
type SharedLabels struct {
    Errors  SharedErrorLabels   // 14 fields: PermissionDenied, InvalidFormData, NotFound, IDRequired, etc.
    Confirm SharedConfirmLabels // 8 fields: Activate, Deactivate, Delete, Block, Remove, BulkActivate, etc.
    Badges  SharedBadgeLabels   // 5 fields: Allow, Deny, Yes, No, NoPermission
}
```

### Mapping Helpers

```go
// Map pyeza.CommonLabels into types.TableLabels for table views
tableLabels := entydad.MapTableLabels(commonLabels)

// Map pyeza.CommonLabels into types.BulkActionsConfig for bulk operations
bulkCfg := entydad.MapBulkConfig(commonLabels)
```

### RoleBadge Type

```go
type RoleBadge struct {
    Name  string
    Color string
}
```

Used by `user/list` to display role chips per user. Consumer apps provide a `GetUserRolesMap func(ctx) (map[string][]RoleBadge, error)`.

## HTMX Helpers

Defined in `htmx.go`:

```go
// HTMXSuccess returns a header-only 200 response that closes the sheet
// and triggers a table refresh via HX-Trigger header:
// {"formSuccess":true,"refreshTable":"<tableID>"}
func HTMXSuccess(tableID string) view.ViewResult

// HTMXError returns a header-only 422 response with an error message
// in the HX-Error-Message header. The sheet.js handleResponse reads this.
func HTMXError(message string) view.ViewResult
```

All action handlers use these helpers for consistent HTMX response patterns:
- Success: closes the drawer/sheet and refreshes the relevant table
- Error: keeps the sheet open and displays the error message

## DataSource Interface

```go
type DataSource interface {
    ListSimple(ctx context.Context, collection string) ([]map[string]any, error)
}
```

Technology-agnostic data access for views that need raw record lists (e.g., client detail's purchase history). Espyna's `DatabaseAdapter` satisfies this directly.

## Asset Pipeline

The package ships CSS assets in `assets/css/`. Consumer apps copy them at startup:

```go
cssDir := filepath.Join("assets", "css")
if err := entydad.CopyStyles(cssDir); err != nil {
    log.Printf("Warning: Failed to copy entydad styles: %v", err)
}

jsDir := filepath.Join("assets", "js")
if err := entydad.CopyStaticAssets(jsDir); err != nil {
    log.Printf("Warning: Failed to copy entydad assets: %v", err)
}
```

Both functions use `runtime.Caller(0)` via `packageDir()` to discover the package directory at runtime. Files are copied to `{targetDir}/entydad/` to keep them namespaced.

### Available CSS Assets
- `entydad-client-dashboard.css` -- client dashboard layout
- `entydad-client-detail.css` -- client detail page
- `entydad-user-dashboard.css` -- user dashboard with chart
- `entydad-user-detail.css` -- user detail page
- `entydad-role-detail.css` -- role detail page
- `entydad-location-detail.css` -- location detail page
- `entydad-supplier-detail.css` -- supplier detail page
- `login01.css` -- simple login page
- `login02.css` -- split-screen login page

## Consumer App Wiring Guide

### Step 1: Create an App-Level Module

Each entity needs a thin app-level module in `apps/{app}/internal/presentation/{entity}/module.go`. The module translates app-level dependencies into entydad view `Deps`:

```go
package client

import (
    "context"

    pyeza "github.com/erniealice/pyeza-golang"
    "github.com/erniealice/pyeza-golang/types"
    "github.com/erniealice/pyeza-golang/view"

    "github.com/erniealice/entydad-golang"
    clientaction "github.com/erniealice/entydad-golang/views/client/action"
    clientdashboard "github.com/erniealice/entydad-golang/views/client/dashboard"
    clientdetail "github.com/erniealice/entydad-golang/views/client/detail"
    clientlist "github.com/erniealice/entydad-golang/views/client/list"

    clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
)

type ModuleDeps struct {
    Routes          entydad.ClientRoutes
    CommonLabels    pyeza.CommonLabels
    Labels          entydad.ClientLabels
    TableLabels     types.TableLabels
    GetListPageData func(ctx context.Context, req *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
    CreateClient    func(ctx context.Context, req *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error)
    ReadClient      func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
    UpdateClient    func(ctx context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error)
    DeleteClient    func(ctx context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error)
    SetActive       func(ctx context.Context, id string, active bool) error
}

type Module struct {
    routes    entydad.ClientRoutes
    Dashboard view.View
    List      view.View
    Table     view.View
    Detail    view.View
    TabAction view.View
    Add       view.View
    Edit      view.View
    Delete    view.View
    BulkDelete    view.View
    SetStatus     view.View
    BulkSetStatus view.View
}

func NewModule(deps *ModuleDeps) *Module {
    listDeps := &clientlist.Deps{
        Routes: deps.Routes, GetListPageData: deps.GetListPageData,
        Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels,
    }
    actionDeps := &clientaction.Deps{
        Routes: deps.Routes, CreateClient: deps.CreateClient,
        ReadClient: deps.ReadClient, UpdateClient: deps.UpdateClient,
        DeleteClient: deps.DeleteClient, SetClientActive: deps.SetActive,
    }
    detailDeps := &clientdetail.Deps{
        Routes: deps.Routes, ReadClient: deps.ReadClient,
        Labels: deps.Labels, CommonLabels: deps.CommonLabels,
    }
    return &Module{
        routes:        deps.Routes,
        Dashboard:     clientdashboard.NewView(&clientdashboard.Deps{CommonLabels: deps.CommonLabels}),
        List:          clientlist.NewView(listDeps),
        Table:         clientlist.NewTableView(listDeps),
        Detail:        clientdetail.NewView(detailDeps),
        TabAction:     clientdetail.NewTabAction(detailDeps),
        Add:           clientaction.NewAddAction(actionDeps),
        Edit:          clientaction.NewEditAction(actionDeps),
        Delete:        clientaction.NewDeleteAction(actionDeps),
        BulkDelete:    clientaction.NewBulkDeleteAction(actionDeps),
        SetStatus:     clientaction.NewSetStatusAction(actionDeps),
        BulkSetStatus: clientaction.NewBulkSetStatusAction(actionDeps),
    }
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
    r.GET(m.routes.DashboardURL, m.Dashboard)
    r.GET(m.routes.ListURL, m.List)
    r.GET(m.routes.TableURL, m.Table)
    r.GET(m.routes.DetailURL, m.Detail)
    r.GET(m.routes.TabActionURL, m.TabAction)
    r.GET(m.routes.AddURL, m.Add)
    r.POST(m.routes.AddURL, m.Add)
    r.GET(m.routes.EditURL, m.Edit)
    r.POST(m.routes.EditURL, m.Edit)
    r.POST(m.routes.DeleteURL, m.Delete)
    r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
    r.POST(m.routes.SetStatusURL, m.SetStatus)
    r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
}
```

### Step 2: Wire in views.go

In `apps/{app}/internal/composition/views.go`, instantiate modules and register routes:

```go
func RegisterAllRoutes(deps *ViewDeps, routes *RouteRegistry) {
    clientmod.NewModule(&clientmod.ModuleDeps{
        Routes:          entydad.DefaultClientRoutes(),
        CommonLabels:    deps.CommonLabels,
        Labels:          deps.ClientLabels,
        TableLabels:     deps.TableLabels,
        GetListPageData: deps.ClientGetListPageData,
        CreateClient:    deps.CreateClient,
        ReadClient:      deps.ReadClient,
        UpdateClient:    deps.UpdateClient,
        DeleteClient:    deps.DeleteClient,
        SetActive:       deps.SetClientActive,
    }).RegisterRoutes(routes)

    usermod.NewModule(&usermod.ModuleDeps{
        Routes:          entydad.DefaultUserRoutes(),
        CommonLabels:    deps.CommonLabels,
        Labels:          deps.UserLabels,
        // ... other deps
    }).RegisterRoutes(routes)

    // ... repeat for role, location, permission, workspace, supplier
}
```

### Step 3: Register Template Filesystems in container.go

Import the entity view packages to register their embedded template filesystems:

```go
import (
    entydad_client "github.com/erniealice/entydad-golang/views/client"
    entydad_user "github.com/erniealice/entydad-golang/views/user"
    entydad_role "github.com/erniealice/entydad-golang/views/role"
    entydad_location "github.com/erniealice/entydad-golang/views/location"
    entydad_permission "github.com/erniealice/entydad-golang/views/permission"
    entydad_workspace "github.com/erniealice/entydad-golang/views/workspace"
    entydad_supplier "github.com/erniealice/entydad-golang/views/supplier"
    entydad_login01 "github.com/erniealice/entydad-golang/views/login01"
    entydad_login02 "github.com/erniealice/entydad-golang/views/login02"
)
```

Each package exposes `TemplatesFS embed.FS` which pyeza discovers at startup.

### Step 4: Copy Assets

In `container.go`, copy entydad CSS/JS assets at startup alongside pyeza and centymo assets:

```go
cssDir := filepath.Join("assets", "css")
if err := entydad.CopyStyles(cssDir); err != nil {
    log.Printf("Warning: Failed to copy entydad styles: %v", err)
}
```

### Step 5: Load Labels

Load labels from lyngua JSON files and pass them through the composition:

```go
var clientLabels entydad.ClientLabels
translations.LoadPath("en", cfg.BusinessType, "client.json", "client", &clientLabels)

var sharedLabels entydad.SharedLabels
translations.LoadPath("en", cfg.BusinessType, "shared.json", "", &sharedLabels)

tableLabels := entydad.MapTableLabels(commonLabels)
```

## Consumer Apps

### retail-admin

Uses all entydad modules: client (with dashboard, detail, tags), user (with dashboard, detail, roles), role (with detail, permissions, users), location (with detail), permission, workspace, supplier, login01, login02.

Wiring: `apps/retail-admin/internal/composition/views.go` and `container.go`.

### service-admin

Uses the same set of entydad modules. Lyngua route overrides remap labels for service businesses: clients to customers, products to services, sales to bookings.

Wiring: `apps/service-admin/internal/composition/views.go` and `container.go`.

## RBAC Integration

All action handlers check permissions via `view.GetUserPermissions(ctx)`:

```go
perms := view.GetUserPermissions(ctx)
if !perms.Can("client", "create") {
    return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
}
```

Table primary actions and row actions set `Disabled: !perms.Can(entity, operation)` with tooltip from `SharedLabels.Badges.NoPermission`.

## In-Use Protection

List views accept an optional `GetInUseIDs func(ctx, ids []string) (map[string]bool, error)` dependency. When provided, rows for in-use entities show disabled delete buttons with `SharedLabels.Errors.CannotDeleteInUse` tooltip, and bulk delete skips items with `RequiresDataAttr: "deletable"`.

## Key Dependencies

- `github.com/erniealice/pyeza-golang` -- UI framework (view, types, route helpers)
- `github.com/erniealice/esqyma` -- proto schemas (client, user, role, etc.)
- `github.com/erniealice/lyngua` -- i18n translation loading (used by consumer apps)

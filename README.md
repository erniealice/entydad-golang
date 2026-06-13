# entydad-golang

Identity and entity domain package for Ichizen OS. Owns the esqyma **`entity`** domain -- the largest in the monorepo -- organized into four sub-contexts (identity, party, location, commerce) plus a cross-package slice of the **`tax`** domain (tax_registration). Also provides three service surfaces: auth, dashboard, and portal.

**Module path:** `github.com/erniealice/entydad-golang`

## Domain ownership

entydad maps to two esqyma proto domains:

- **`entity`** (`proto/v1/domain/entity/`) -- the primary domain. Entities are grouped into four navigation-only sub-contexts:
  - **identity**: workspace, user, role, permission, workspace_user, workspace_user_role
  - **party**: client, client_tag, supplier, supplier_tag
  - **location**: location, location_area
  - **commerce**: payment_term
- **`tax`** (`proto/v1/domain/tax/`) -- a cross-package-split domain (thread TX). entydad owns `tax_registration`; fycha owns `tax_rate`. Each package declares its own `package tax` facade on its own import path, re-exporting only the entities it owns.

Two entities (`client_tag`, `supplier_tag`) are entydad-local tag entities with no 1:1 esqyma proto entity (esqyma's only `*_tag` entity is `event/event_tag`). They carry `legacyAllow` entries pending either esqyma promotion or folding into their parent entity as sub-views.

## Package structure (Option B)

Under Option B the ENTITY is the contract package. entydad realizes this with a **navigation-only sub-context layer** unique to the large `entity` domain: `domain/entity/<subcontext>/<entity>/`. Each `<subcontext>` is a folder grouping only (not a Go contract boundary). There is exactly ONE Go contract package per entity, and exactly ONE import namespace for consumers: `entity.*`.

```
entydad-golang/
  placement_test.go            # B-STRICT placement gate -- the ONLY test file at root
  go.mod / go.sum
  domain/
    entity/                    # package entity -- facade for the entity domain
      entity.go                # facade: type ClientLabels = client.Labels, etc. (aliases only)
      identity/                # sub-context: identity management
        helpers.go             # sub-context-shared wiring helpers
        permission_module.go   # NewPermissionModule() assembler
        role_module.go         # NewRoleModule() assembler
        user_module.go         # NewUserModule() assembler
        workspace_module.go    # NewWorkspaceModule() assembler
        workspace_user_module.go
        workspace_user_role_module.go
        permission/            # entity: esqyma entity/permission
          labels.go            # Labels struct (PermissionLabels in the facade)
          routes.go            # Routes struct + DefaultRoutes()
          embed.go             # template embed.FS
          action/action.go     # CRUD action handlers
          form/form.go         # form.Data struct
          list/page.go         # list page handler
          templates/           # permission-drawer-form.html, list.html
        role/                  # entity: esqyma entity/role
          labels.go
          routes.go
          embed.go
          action/ detail/ form/ list/ templates/
          permissions/         # role-permission tab handlers
          users/               # role-user tab handlers
        user/                  # entity: esqyma entity/user
          labels.go
          routes.go
          embed.go
          action/ detail/ form/ list/ templates/
          dashboard/page.go    # user dashboard handler
          roles/               # user-role tab handlers
        workspace/             # entity: esqyma entity/workspace
          labels.go
          routes.go
          embed.go
          action/ detail/ form/ list/ templates/
        workspace_user/        # entity: esqyma entity/workspace_user
          labels.go
          routes.go
          embed.go
          action/ detail/ form/ list/ templates/
        workspace_user_role/   # entity: esqyma entity/workspace_user_role
          labels.go
          routes.go
          embed.go (no templates/ -- inline-rendered)
          action/ form/ templates/
      party/                   # sub-context: party management
        client_module.go       # NewClientModule() assembler
        client_tag_module.go
        supplier_module.go
        supplier_tag_module.go
        client/                # entity: esqyma entity/client
          labels.go
          routes.go
          embed.go
          action/ dashboard/ detail/ form/ list/ tag/ templates/
        client_tag/            # entydad-local tag entity (no esqyma proto)
          labels.go
          routes.go
        supplier/              # entity: esqyma entity/supplier
          labels.go
          routes.go
          embed.go
          action/ dashboard/ detail/ form/ list/ tag/ templates/
        supplier_tag/          # entydad-local tag entity (no esqyma proto)
          labels.go
          routes.go
      location/                # sub-context: location management
        location_module.go
        location_area_module.go
        location/              # entity: esqyma entity/location
          labels.go
          routes.go
          embed.go
          action/ dashboard/ detail/ form/ list/ templates/
        location_area/         # entity: esqyma entity/location_area
          labels.go
          routes.go
          embed.go
          action/ form/ list/ templates/
      commerce/                # sub-context: commerce
        payment_term_module.go
        payment_term/          # entity: esqyma entity/payment_term
          labels.go
          routes.go
          embed.go
          action/ form/ list/ templates/
    tax/                       # package tax -- facade for the entydad slice of tax domain
      tax.go                   # facade: type TaxRegistrationLabels = taxregistration.Labels, etc.
      tax_registration_module.go
      tax_registration/        # entity: esqyma tax/tax_registration
        labels.go
        routes.go
        embed.go
        action/ form/ list/ templates/
  block/
    block.go                   # Block() constructor -- pyeza.AppOption entry point
    usecases.go                # *UseCases typed wiring contract + RequireFor + MustValidate
    catalog.go                 # compose-v2 Unit binders (AllUnits aggregator)
    identity.go                # per-subcontext module registration (identity entities)
    party.go                   # per-subcontext module registration (party entities)
    commerce.go                # per-subcontext module registration (commerce entities)
    infra.go                   # Infra struct -- AppContext subset for view modules
    helpers.go                 # categoryListPageDataGetter interface, statement helpers
    route_loading.go           # blockLabels / blockRoutes types + lyngua loaders
    wiring.go                  # placeholder for future cross-cutting wiring helpers
    block_test.go              # MustValidate fail-closed wiring tests
    routes_test.go             # route override tests
    usecases_test.go           # UseCases contract tests
  service/
    auth/views/                # auth service surface
      login01/                 # login variant 1 (page.go, action.go, embed.go, templates/)
      login02/                 # login variant 2
        select-workspace-role/ # workspace/role selection step
      signup01/                # signup variant 1
      signup02/                # signup variant 2
      reset-password01/        # password reset variant 1
      reset-password02/        # password reset variant 2
      change-password/         # change password (logged-in)
    dashboard/views/admin/     # admin dashboard surface
      module.go                # admin dashboard module assembler
      embed.go
      dashboard/page.go        # dashboard page handler
      templates/
    portal/views/              # member portal surface
      labels_portal.go         # PortalLabels struct (from member.json)
      me/                      # member home page
        page.go                # "me" landing
        profile_overview/      # profile overview tab
        notifications/         # notifications tab
        inbox/ invoices/ recent_activity/
        templates/
      profile/                 # profile management
      account/                 # account management
      billing/                 # billing management
      preference/              # preferences management
```

## Placement gate (`placement_test.go`)

entydad carries a **B-STRICT** placement gate (v2, Option B). `legacyAllow` holds six dated residuals; the target state is empty (STRICT).

| Rule | What it checks |
|------|----------------|
| **R1** Empty root | No package `.go` files at module root -- only `_test.go` permitted |
| **R2** Canonical dirs | Every first-level dir is an allowed infra surface; every `domain/<d>` is an esqyma proto domain |
| **R2'** Entity dirs | Every `domain/<d>/<child>/` DIR is an esqyma entity of domain `<d>`, `shared`, or a domain-view (name starts with `<d>`) |
| **R3'** Entity contract | No real `*Labels`/`*Routes` type declaration at the domain root -- only alias re-exports (`type X = pkg.Y`) are allowed |
| **R4** No god-files | No `.go` file (excl. `_test.go`) may exceed 1,200 lines |
| **R5** Facade exists | A facade `domain/<d>/<d>.go` must exist for every domain dir with >=1 entity subdir |
| **R6** No cycles | Enforced by `lint-no-domain-cycles.sh` (external, go-list based) |

`crossCutting = false` -- the domain variant applies. `subContexts = ["commerce", "identity", "location", "party"]` enables the one-level recursion unique to entydad: R2' recurses into each sub-context folder and validates grandchildren as entities of the `entity` domain. This is the ONLY package in the fleet with a non-empty `subContexts`.

Current `legacyAllow` residuals (all EXPIRES 2026-07-15):
- `assets.go` -- embed FS for the auth/dashboard service surface (capstone: move to pyeza)
- `labels.go` -- auth/login/signup/reset/change-password label types (capstone: relocate under service/auth)
- `routes.go` -- root route constants for auth/dashboard (capstone: relocate under service/)
- `routes_config.go` -- Login/Auth/AdminDashboard route structs (capstone: relocate under service/)
- `domain/entity/party/client_tag` -- entydad-local tag entity, no 1:1 esqyma proto
- `domain/entity/party/supplier_tag` -- entydad-local tag entity, no 1:1 esqyma proto

## Fail-closed wiring (`block/usecases.go`)

`*UseCases` is the typed wiring contract between service-admin's composition layer and entydad's view modules. `RequireFor(cfg)` lists every missing REQUIRED closure for the enabled modules. `MustValidate(cfg)` adds fail-closed posture:

- **dev/test** (`testing.Testing()` true or `ENTYDAD_BLOCK_STRICT` truthy): PANIC with the full field list -- uncatchable-by-accident, stack-traced, fails CI loudly.
- **prod**: `log.Printf("FATAL: ...")` at the seam AND returns the error -- `Block()` propagates -- `NewServiceAdmin` halts boot.

OPTIONAL closures (tag entities, dashboards, reports, conversation assign/set-status drawers) are never flagged -- they degrade gracefully to empty-state.

The former `DataSource` interface duck was DELETED 2026-06-12. All data writes now flow through narrow typed `UseCases.SetActive`/`SetStatus` primitives bound by service-admin.

## Compose-v2 catalog (`block/catalog.go`)

Each binder function constructs a `compose.Unit` with a typed `Mount` closure that wires the entity module's `Deps` from `UseCases` + `Infra`, then calls `NewXxxModule(deps).RegisterRoutes(mc.Routes)`. The `AllUnits(uc, infra)` aggregator returns the full ordered slice for the compose engine. The structure mirrors the three sub-context files (`identity.go`, `party.go`, `commerce.go`) plus the tax domain.

## Service surfaces

entydad hosts three service surfaces under `service/` -- these are NOT entity views but cross-cutting app-level pages:

- **`service/auth/`** -- login (two variants), signup (two variants), reset-password (two variants), change-password. Each variant is a self-contained view package with `page.go`, `action.go`, `embed.go`, and `templates/`. Login02 includes a `select-workspace-role/` sub-step for multi-workspace users.
- **`service/dashboard/`** -- admin dashboard (`views/admin/`). Module assembler + dashboard page handler.
- **`service/portal/`** -- member self-service portal. Five view sections: me (home landing with profile overview, notifications, inbox, invoices, recent activity tabs), profile, account, billing, preference. `PortalLabels` loaded from `member.json` via lyngua.

## Cross-package split: tax domain

The `tax` esqyma domain is split across two packages:
- **entydad** owns `tax_registration` (this package, `domain/tax/`)
- **fycha** owns `tax_rate`

Each package has its own `package tax` facade on its own import path. The placement gate validates each package's entity set independently against esqyma.

## Dependencies

- `github.com/erniealice/pyeza-golang` -- UI framework (view system, template engine, types)
- `github.com/erniealice/esqyma` -- proto schemas (entity + tax domains, plus cross-domain refs)
- `github.com/erniealice/lyngua` -- translation/i18n
- `github.com/erniealice/espyna-golang` -- typed use cases (appcontext, reference checker)
- `github.com/erniealice/hybra-golang` -- cross-cutting views (conversation module, imported by block)
- `github.com/erniealice/centymo-golang` -- route structs for cross-centymo navigation

## Role in the monorepo

entydad sits in the domain layer above pyeza and espyna. It is the foundational identity package -- workspaces, users, roles, and permissions are prerequisites for every other domain. Consumer apps (e.g., `apps/service-admin`) call `block.Block()` to mount the entity modules, supplying a `*UseCases` via `block.WithUseCases(...)`. The typed contract ensures any drift between espyna and entydad is a compile error, not a silent nil.

The auth and portal service surfaces are also mounted via `Block()`, making entydad the single entry point for identity, entity CRUD, authentication, and member self-service.

See `docs/wiki/articles/vertical-slices.md` for the full entity trace and `docs/wiki/articles/package-map.md` for the monorepo dependency graph.

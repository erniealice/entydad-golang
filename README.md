# entydad-golang

**Module:** `github.com/erniealice/entydad-golang` · **Go:** 1.25.1

Domain package for the esqyma **`entity`** domain (identity + party + location +
commerce) plus the auth / dashboard / portal **service surfaces**. Provides
reusable views, templates, labels, routes, and HTMX helpers for users, clients,
roles, permissions, locations, workspaces, suppliers, payment terms, and tags.
Consumer apps (service-admin) register entydad through its `block` assembler.

## Architecture — Option B (entity-as-package)

Under Option B the **entity** is the contract package and `domain/<d>/` is a
domain facade over its entities. entydad realizes this with a **navigation-only
sub-context layer** unique to the large `entity` domain:

```
domain/entity/<subcontext>/<entity>/
```

`<subcontext>` ∈ `{ identity, party, location, commerce }` is a **folder grouping
only** (not a Go contract boundary). There is exactly ONE Go contract package per
entity, and exactly ONE import namespace for consumers: `entity.*`.

```
packages/entydad-golang/
├── domain/
│   └── entity/                         # the esqyma `entity` domain
│       ├── entity.go                   # ← FACADE: one hand-written facade for the
│       │                               #   whole domain. Re-exports E-prefixed
│       │                               #   names (ClientLabels, DefaultClientRoutes)
│       │                               #   as Go ALIASES to the entity-local types.
│       ├── identity/
│       │   ├── <entity>_module.go      # hoisted assembler per entity (permission,
│       │   │                           #   role, user, workspace, workspace_user,
│       │   │                           #   workspace_user_role)
│       │   ├── helpers.go              # subcontext-shared wiring helpers
│       │   └── <entity>/               # CONTRACT package `package <entity>`
│       │       ├── labels.go           #   entity-local Labels / PageLabels / …
│       │       ├── routes.go           #   entity-local Routes + DefaultRoutes()
│       │       ├── embed.go            #   per-entity template FS
│       │       ├── action|form|list|detail|…/   # render leaves (page.go etc.)
│       │       └── templates/*.html
│       ├── party/                      # client, client_tag, supplier, supplier_tag
│       │   ├── <entity>_module.go
│       │   └── <entity>/ …
│       ├── location/                   # location, location_area
│       │   ├── <entity>_module.go
│       │   └── <entity>/ …
│       └── commerce/                   # payment_term
│           ├── payment_term_module.go
│           └── payment_term/ …
├── block/                              # the ASSEMBLER (Lego pattern)
│   ├── block.go                        #   Block() → pyeza.AppOption; per-module BlockOption
│   ├── usecases.go                     #   typed UseCases contract + RequireFor + MustValidate
│   ├── identity.go / party.go / commerce.go  # per-subcontext module registration
│   ├── route_loading.go / helpers.go / wiring.go
├── service/                            # cross-cutting service SURFACES (not entities)
│   ├── auth/views/…                    #   login01/02, signup, reset/change-password
│   ├── dashboard/views/admin/…         #   admin dashboard
│   └── portal/views/…                  #   account, billing, me, preference, profile
├── assets/css/                         # per-view CSS (lf-/entydad- prefixed)
└── tests/                              # Playwright e2e
```

### Facade contract (`domain/entity/entity.go`)

- Consumers use the **E-prefixed** names (`ClientLabels`, `DefaultClientRoutes`,
  `WorkspaceUserAddURL`) which resolve via Go aliases to the entity-local
  `Labels` / `Routes` / `DefaultRoutes` in each entity package.
- **Import-cycle rule (contract D):** entity packages MUST NEVER import the
  facade. Cross-entity references go DIRECT to the sibling package. Only the
  facade and the `<entity>_module.go` assemblers import the entity packages —
  keeping each `domain/entity/<sub>/` subtree a child→parent DAG (enforced by
  `lint-no-domain-cycles.sh`).

## The `block` assembler

`Block()` returns a `pyeza.AppOption` that registers entydad's modules using the
shared `AppContext` infrastructure carrier:

```go
app.Apply(entydadblock.Block(entydadblock.WithUseCases(uc)))          // all modules
app.Apply(entydadblock.Block(entydadblock.WithUseCases(uc), entydadblock.WithClient()))  // one
```

### Fail-closed wiring (`MustValidate`)

`block/usecases.go` declares the typed `UseCases` contract entydad needs from the
consumer. Two guards protect against the nil-closure trap (a nil func-field
compiles and silently renders an empty feature):

- **`RequireFor(cfg)`** — the *policy*: for each ENABLED module it asserts every
  REQUIRED closure is non-nil, accumulating a named-field list. Optional ports
  (e.g. `clientTag`, dashboards, reports, the Conversation assign/set-status
  drawers) are intentionally never asserted.
- **`MustValidate(cfg)`** — the *fail-closed posture*, called at `Block()` entry:
  - **dev / test / CI** (`testing.Testing()`, or `ENTYDAD_BLOCK_STRICT` truthy):
    a missing REQUIRED closure **PANICS** with the full field list — loud,
    stack-traced, uncatchable-by-accident.
  - **prod**: logs a screaming `FATAL:` line at the seam **and** returns the
    error so `Block()` propagates it and `NewServiceAdmin` halts boot.

  This mirrors the AUTHZ_ENFORCE boot-guard: dev can't merge a gap, prod can't
  boot one. See `docs/audit/architecture/fail-closed-wiring.md`.

## Structural invariants (placement gate)

`placement_test.go` (run `go test -run Placement .`) is the canonical Option-B
placement gate (`docs/orchestrate/20260610-package-cleanup/placement_test.template.go`).
It derives the esqyma domain + entity sets live from `packages/esqyma/proto/`, so
the rules cannot drift from proto. It enforces: empty root (R1), canonical
domains (R2), entity-dir placement (R2′), facade alias-only at the domain root
(R3′), no god-files >1200 lines (R4), and facade presence (R5). Intra-domain
import cycles (R6) are caught by `lint-no-domain-cycles.sh entydad`.

`legacyAllow` is the shrinking migration ledger; each entry carries a dated
`EXPIRES` stamp. The capstone target is **STRICT** (empty `legacyAllow`).
</content>

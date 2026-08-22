# entydad-golang/block

This package is the **composition entrypoint** for mounting entydad features into a pyeza app.

`block.Block()` returns an `espyna-golang` `consumerapp.AppOption` and is called by app composition layers (for example `apps/service-admin`) to register domain modules, services, and routes.

The package enforces a **typed composition seam**:

- `block.WithUseCases(...)` injects the adapter-shared `*UseCases` contract from the app.
- `Block()` validates that contract with `UseCases.MustValidate` before registration.
- `block.With...()` feature flags (for example `WithAdmin`, `WithClient`, `WithWorkspace`) control which modules are enabled; no options means "register all".
- Integration objects (`blockLabels`, `blockRoutes`) are loaded from Lyngua in `route_loading.go` and passed to module constructors.

Each domain surface then assembles through `catalog.go` and sub-files:

- `identity.go`, `party.go`, `commerce.go` map entydad sub-context modules into compose units.
- `infra.go` and `helpers.go` provide the narrow app context needed by view modules.
- `auth_bridge.go`, `principal_*`, and `workspace_*` wire auth/identity and principal-switching actions.
- `usecases.go` defines the boundary contract that prevents direct DB/repo access from this package.

`auth_bridge.go` also owns environment-to-auth composition. For local Firebase impersonation it
computes the exact `APP_ENVIRONMENT=local` plus
`AUTH_FIREBASE_ALLOW_IMPERSONATION=true` gate and injects only a narrow custom-token function. The
auth service and view do not read environment variables or import the Firebase SDK.

Route and label overrides are loaded at composition time (not embedded at runtime mutation), and dashboard wiring is intentionally kept thin in `wiring.go`.

There are no side effects until `Block()` is invoked by app startup, so this package remains a pure adapter layer between:

1) app composition (consumer modules/services)
2) domain module registration (entydad views and actions)

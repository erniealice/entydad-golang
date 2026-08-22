# Authentication Service Module

This package assembles Entydad's reusable HTTP authentication surface for host
applications. It connects login, signup, password lifecycle, logout, principal
selection, session rotation, and CSRF issuance without owning the provider or
database implementation.

## Key files

- `composition.go` defines module dependencies and registers view and action
  routes. It also maps configured Firebase methods to login providers.
- `handlers.go` implements password and Firebase login tails, principal resolution,
  session rotation, redirects, and other auth actions.
- `interfaces.go` defines the narrow provider, session, principal, rendering, and
  route-registration seams supplied by an application.
- `types.go`, `helpers.go`, and `carousel.go` hold shared request and presentation
  data used across the auth screens.

Espyna supplies authentication and session adapters; Pyeza supplies rendering and
route contracts. A host app supplies provider-specific configuration. Federated
buttons such as Microsoft appear only when Firebase web configuration is present
and the corresponding sign-in method is allowed.

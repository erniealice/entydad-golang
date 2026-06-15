# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0-alpha] - 2026-06-15

Identity domain — first published alpha.

### Added
- Clients, suppliers, users, roles, permissions, locations, location areas, workspaces, client tags, payment terms, and auth screens (login, signup, reset-password). Consumes centymo route-config types (SubscriptionRoutes, PriceScheduleRoutes).

### Changed
- `go.mod` now references published tags (`v0.1.0-alpha`) instead of local `replace` directives; local development continues via `go.work`.

[Unreleased]: https://github.com/erniealice/entydad-golang/compare/v0.1.0-alpha...HEAD
[0.1.0-alpha]: https://github.com/erniealice/entydad-golang/releases/tag/v0.1.0-alpha

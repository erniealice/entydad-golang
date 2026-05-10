package block

// wiring.go previously held reflective dashboard wiring helpers that called
// espyna's internal use cases via reflect to avoid import cycles.
//
// Those helpers are now obsolete: block.go wires dashboard data directly from
// the UseCases.GetLocationDashboardPageData and UseCases.GetAdminDashboardPageData
// closure fields, which are supplied by service-admin's adapter function.
//
// This file is kept as a placeholder so that any future cross-cutting wiring
// helpers have a natural home. Add helpers here when 2+ callers in block.go
// need the same wiring logic.

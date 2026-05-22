// Package me holds the shared types used by all /me/* personal-scope view
// sub-packages.
//
// Per docs/plan/20260521-workspace-keyed-routing/phases.md Phase P9b the /me/*
// tree is cross-workspace personal context (notifications, invoices, profile
// overview, recent workspace switches). It bypasses the workspace_path
// middleware (`/me/` is in the early-exit list) and never carries an
// implicit-workspace prefix.
//
// v1 lock per phases.md 9b: /me/* is read-only. Cross-workspace actions
// return the user to the originating /w/{ws}/* page for execution.
//
// Created 2026-05-22 in Phase P9b.
package me

import (
	"github.com/erniealice/pyeza-golang/types"
)

// PageData is the shared base embedded by every /me/* view's concrete PageData.
// Mirrors the shape of staff /app/* pages so app-shell.html renders without
// surprises (header title/icon driven by HeaderTitle/HeaderIcon).
type PageData struct {
	types.PageData
}

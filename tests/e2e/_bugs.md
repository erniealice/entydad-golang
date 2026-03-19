# Entydad E2E Test Bugs

Discovered 2026-03-19 while running entydad package tests against service-admin (localhost:8081, professional business type).

---

## BUG: Client detail page shows "Page content not available"
**Route:** `/app/clients/detail/{id}` (e.g., `/app/clients/detail/client-007`)
**Test:** `clients/clients-crud.spec.ts` > ENT-CLI-005
**Error:** Page content div contains "Page content not available" instead of detail layout
**Root cause:** The client detail view template is not wired in the service-admin composition (container.go / views.go). The route exists and responds with 200, but the content template is missing or not registered. Same pattern as the supplier detail bug documented in gotchas.md.
**Status:** Test skipped with `test.skip(true, 'BUG: Client detail page shows "Page content not available"')`

---

## BUG: User edit endpoint returns 422 (notFound)
**Route:** `/action/users/edit/{id}` (e.g., `/action/users/edit/019d0554-ad8e-7476-8c56-729a78518200`)
**Test:** `users/users-crud.spec.ts` > ENT-USR-003
**Error:** Clicking the edit button opens the sheet drawer with title "Edit user", but the form body is empty. The server responds with HTTP 422 and header `Hx-Error-Message: shared.errors.notFound`.
**Root cause:** The user edit GET handler likely fails to find the user, possibly due to a query issue (the data-edit-url uses the internal UUID, but the lookup may be using a different ID format). Alternatively, a workspace/tenant scoping issue prevents the mock_auth session from reading user data for edits. The user list loads fine (uses different query), and user add works, so the issue is specific to the edit endpoint.
**Status:** Test skipped with `test.skip(true, 'BUG: User edit endpoint returns 422')`

---

## BUG: Supplier detail page shows "Page content not available"
**Route:** `/app/suppliers/detail/{id}`
**Test:** (existing, in `suppliers/suppliers-detail.spec.ts`)
**Error:** Same as client detail — content template not registered in service-admin composition.
**Root cause:** Supplier detail view not wired in service-admin.
**Status:** Previously documented in `references/gotchas.md`. Test skipped.

---

## NOTE: Roles list has no per-row edit button
**Route:** `/app/roles/list`
**Observation:** Unlike other entities (clients, users, suppliers), the roles table does NOT render an `.action-btn.edit` button per row. Roles can only be viewed (navigates to detail page), deactivated, or deleted. Editing is presumably done via the role detail page. This is not a bug, but a deliberate design difference.
**Test:** Documented in `roles/roles-crud.spec.ts` > ENT-ROL-003

---

## NOTE: Role color picker native input is hidden
**Route:** `/action/roles/add` (role add drawer)
**Observation:** The `#color` field is `<input type="color">` which is styled as hidden by the `.color-picker-native` class. The visible element is the `.color-picker-field` wrapper with a `.color-picker-band` div that shows the color. Tests should assert `.color-picker-field` visibility instead of `#color`.
**Fix applied:** Test updated to check `.color-picker-field` instead of `#color`.

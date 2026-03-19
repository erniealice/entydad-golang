# Entydad Route Test Coverage

**Source:** `packages/entydad-golang-ryta/routes.go`
**Generated:** 2026-03-19
**Last updated:** 2026-03-19

Legend: `[x]` = covered, `[ ]` = not started, `[~]` = partial (some assertions, not full CRUD)

---

## Clients (13 routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `ClientDashboardURL` | `/app/clients/dashboard` | page | — | `[ ]` |
| `ClientListURL` | `/app/clients/list/{status}` | page | `customers/customers-crud.spec.ts` | `[x]` retail, `[x]` service |
| `ClientTableURL` | `/action/clients/table/{status}` | action | (implicit via list) | `[~]` |
| `ClientAddURL` | `/action/clients/add` | action | `customers/customers-crud.spec.ts` | `[x]` retail, `[x]` service |
| `ClientEditURL` | `/action/clients/edit/{id}` | action | `customers/customers-crud.spec.ts` | `[x]` retail, `[x]` service |
| `ClientDeleteURL` | `/action/clients/delete` | action | — | `[ ]` |
| `ClientBulkDeleteURL` | `/action/clients/bulk-delete` | action | — | `[ ]` |
| `ClientDetailURL` | `/app/clients/detail/{id}` | page | `client-detail-crm.spec.ts` | `[x]` retail |
| `ClientTabActionURL` | `/action/clients/{id}/tab/{tab}` | action | `client-detail-crm.spec.ts` | `[~]` retail |
| `ClientAttachmentUploadURL` | `/action/clients/{id}/attachments/upload` | action | — | `[ ]` |
| `ClientAttachmentDeleteURL` | `/action/clients/{id}/attachments/delete` | action | — | `[ ]` |
| `ClientSetStatusURL` | `/action/clients/set-status` | action | — | `[ ]` |
| `ClientBulkSetStatusURL` | `/action/clients/bulk-set-status` | action | — | `[ ]` |
| `ClientSearchURL` | `/action/clients/search` | action | — | `[ ]` |

**Coverage: 4/13 routes (31%)**

---

## Users (9 routes + 5 detail routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `UserDashboardURL` | `/app/users/dashboard` | page | `user-dashboard-dynamic.spec.ts` | `[x]` retail |
| `UserListURL` | `/app/users/list/{status}` | page | `users/users-list.spec.ts` | `[x]` retail |
| `UserTableURL` | `/action/users/table/{status}` | action | (implicit via list) | `[~]` |
| `UserAddURL` | `/action/users/add` | action | `users/users-add.spec.ts` | `[x]` retail |
| `UserEditURL` | `/action/users/edit/{id}` | action | `users/users-edit.spec.ts` | `[x]` retail |
| `UserDeleteURL` | `/action/users/delete` | action | `users/users-delete.spec.ts` | `[x]` retail |
| `UserBulkDeleteURL` | `/action/users/bulk-delete` | action | — | `[ ]` |
| `UserSetStatusURL` | `/action/users/set-status` | action | `users/users-status.spec.ts` | `[x]` retail |
| `UserBulkSetStatusURL` | `/action/users/bulk-set-status` | action | — | `[ ]` |
| `UserDetailURL` | `/app/users/detail/{id}` | page | `users/users-detail.spec.ts` | `[x]` retail |
| `UserTabActionURL` | `/action/users/{id}/tab/{tab}` | action | `users/users-detail.spec.ts` | `[~]` retail |
| `UserAttachmentUploadURL` | `/action/users/{id}/attachments/upload` | action | — | `[ ]` |
| `UserAttachmentDeleteURL` | `/action/users/{id}/attachments/delete` | action | — | `[ ]` |
| `UserResetPasswordURL` | `/action/users/reset-password/{id}` | action | — | `[ ]` |

**Coverage: 7/14 routes (50%)**

---

## Roles (12 routes + 4 detail routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `RoleListURL` | `/app/roles/list` | page | `roles/roles-list.spec.ts` | `[x]` retail |
| `RoleTableURL` | `/action/roles/table` | action | (implicit via list) | `[~]` |
| `RoleAddURL` | `/action/roles/add` | action | `roles/roles-add.spec.ts` | `[x]` retail |
| `RoleEditURL` | `/action/roles/edit/{id}` | action | `roles/roles-edit.spec.ts` | `[x]` retail |
| `RoleDeleteURL` | `/action/roles/delete` | action | `roles/roles-delete.spec.ts` | `[x]` retail |
| `RoleBulkDeleteURL` | `/action/roles/bulk-delete` | action | — | `[ ]` |
| `RoleSetStatusURL` | `/action/roles/set-status` | action | `roles/roles-status.spec.ts` | `[x]` retail |
| `RoleBulkSetStatusURL` | `/action/roles/bulk-set-status` | action | — | `[ ]` |
| `RoleDetailURL` | `/app/roles/detail/{id}` | page | `roles/roles-detail.spec.ts` | `[x]` retail |
| `RoleTabActionURL` | `/action/roles/{id}/tab/{tab}` | action | `roles/roles-detail.spec.ts` | `[~]` retail |
| `RoleAttachmentUploadURL` | `/action/roles/{id}/attachments/upload` | action | — | `[ ]` |
| `RoleAttachmentDeleteURL` | `/action/roles/{id}/attachments/delete` | action | — | `[ ]` |
| `RoleUsersURL` | `/app/roles/detail/{id}/users` | page | `roles/roles-users.spec.ts` | `[x]` retail |
| `RoleUsersTableURL` | `/action/roles/detail/{id}/users/table` | action | (implicit) | `[~]` |
| `RoleUsersAssignURL` | `/action/roles/detail/{id}/users/assign` | action | `roles/roles-users.spec.ts` | `[~]` retail |
| `RoleUsersRemoveURL` | `/action/roles/detail/{id}/users/remove` | action | `roles/roles-users.spec.ts` | `[~]` retail |

**Coverage: 8/16 routes (50%)**

---

## Role Permissions (4 routes, migrated)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `RoleDetailPermissionsURL` | `/app/roles/detail/{id}/permissions` | page | `roles/roles-permissions.spec.ts` | `[x]` retail |
| `RoleDetailPermissionsTableURL` | `/action/roles/detail/{id}/permissions/table` | action | (implicit) | `[~]` |
| `RoleDetailPermissionsAssignURL` | `/action/roles/detail/{id}/permissions/assign` | action | `roles/roles-permissions.spec.ts` | `[~]` retail |
| `RoleDetailPermissionsRemoveURL` | `/action/roles/detail/{id}/permissions/remove` | action | `roles/roles-permissions.spec.ts` | `[~]` retail |

**Coverage: 3/4 routes (75%)**

---

## User Roles (4 routes, migrated)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `UserDetailRolesURL` | `/app/users/detail/{id}/roles` | page | `users/users-roles.spec.ts` | `[x]` retail |
| `UserDetailRolesTableURL` | `/action/users/detail/{id}/roles/table` | action | (implicit) | `[~]` |
| `UserDetailRolesAssignURL` | `/action/users/detail/{id}/roles/assign` | action | `users/users-roles.spec.ts` | `[~]` retail |
| `UserDetailRolesRemoveURL` | `/action/users/detail/{id}/roles/remove` | action | `users/users-roles.spec.ts` | `[~]` retail |

**Coverage: 3/4 routes (75%)**

---

## Locations (12 routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `LocationListURL` | `/app/locations/list/{status}` | page | `locations/locations-crud.spec.ts` | `[x]` service |
| `LocationTableURL` | `/action/locations/table/{status}` | action | (implicit) | `[~]` |
| `LocationAddURL` | `/action/locations/add` | action | `locations/locations-crud.spec.ts` | `[x]` service |
| `LocationEditURL` | `/action/locations/edit/{id}` | action | `locations/locations-crud.spec.ts` | `[~]` service |
| `LocationDeleteURL` | `/action/locations/delete` | action | — | `[ ]` |
| `LocationBulkDeleteURL` | `/action/locations/bulk-delete` | action | — | `[ ]` |
| `LocationSetStatusURL` | `/action/locations/set-status` | action | — | `[ ]` |
| `LocationBulkSetStatusURL` | `/action/locations/bulk-set-status` | action | — | `[ ]` |
| `LocationDetailURL` | `/app/locations/{id}` | page | — | `[ ]` |
| `LocationTabActionURL` | `/action/locations/{id}/tab/{tab}` | action | — | `[ ]` |
| `LocationAttachmentUploadURL` | `/action/locations/{id}/attachments/upload` | action | — | `[ ]` |
| `LocationAttachmentDeleteURL` | `/action/locations/{id}/attachments/delete` | action | — | `[ ]` |
| `LocationEditDetailURL` | `/action/locations/edit-detail/{id}` | action | — | `[ ]` |

**Coverage: 3/13 routes (23%)**

---

## Permissions (8 routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `PermissionListURL` | `/app/permissions/list/{status}` | page | — | `[ ]` |
| `PermissionTableURL` | `/action/permissions/table/{status}` | action | — | `[ ]` |
| `PermissionAddURL` | `/action/permissions/add` | action | — | `[ ]` |
| `PermissionEditURL` | `/action/permissions/edit/{id}` | action | — | `[ ]` |
| `PermissionDeleteURL` | `/action/permissions/delete` | action | — | `[ ]` |
| `PermissionBulkDeleteURL` | `/action/permissions/bulk-delete` | action | — | `[ ]` |
| `PermissionSetStatusURL` | `/action/permissions/set-status` | action | — | `[ ]` |
| `PermissionBulkSetStatusURL` | `/action/permissions/bulk-set-status` | action | — | `[ ]` |

**Coverage: 0/8 routes (0%)**

---

## Workspaces (8 routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `WorkspaceListURL` | `/app/workspaces/list/{status}` | page | — | `[ ]` |
| `WorkspaceTableURL` | `/action/workspaces/table/{status}` | action | — | `[ ]` |
| `WorkspaceAddURL` | `/action/workspaces/add` | action | — | `[ ]` |
| `WorkspaceEditURL` | `/action/workspaces/edit/{id}` | action | — | `[ ]` |
| `WorkspaceDeleteURL` | `/action/workspaces/delete` | action | — | `[ ]` |
| `WorkspaceBulkDeleteURL` | `/action/workspaces/bulk-delete` | action | — | `[ ]` |
| `WorkspaceSetStatusURL` | `/action/workspaces/set-status` | action | — | `[ ]` |
| `WorkspaceBulkSetStatusURL` | `/action/workspaces/bulk-set-status` | action | — | `[ ]` |

**Coverage: 0/8 routes (0%)**

---

## Suppliers (12 routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `SupplierListURL` | `/app/suppliers/list/{status}` | page | — | `[ ]` |
| `SupplierTableURL` | `/action/suppliers/table/{status}` | action | — | `[ ]` |
| `SupplierAddURL` | `/action/suppliers/add` | action | — | `[ ]` |
| `SupplierEditURL` | `/action/suppliers/edit/{id}` | action | — | `[ ]` |
| `SupplierDeleteURL` | `/action/suppliers/delete` | action | — | `[ ]` |
| `SupplierBulkDeleteURL` | `/action/suppliers/bulk-delete` | action | — | `[ ]` |
| `SupplierDetailURL` | `/app/suppliers/detail/{id}` | page | — | `[ ]` |
| `SupplierTabActionURL` | `/action/suppliers/{id}/tab/{tab}` | action | — | `[ ]` |
| `SupplierAttachmentUploadURL` | `/action/suppliers/{id}/attachments/upload` | action | — | `[ ]` |
| `SupplierAttachmentDeleteURL` | `/action/suppliers/{id}/attachments/delete` | action | — | `[ ]` |
| `SupplierSetStatusURL` | `/action/suppliers/set-status` | action | — | `[ ]` |
| `SupplierBulkSetStatusURL` | `/action/suppliers/bulk-set-status` | action | — | `[ ]` |

**Coverage: 0/12 routes (0%)**

---

## Reports (2 routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `ReceivablesAgingURL` | `/app/clients/reports/receivables-aging` | page | — | `[ ]` |
| `PayablesAgingURL` | `/app/suppliers/reports/payables-aging` | page | — | `[ ]` |

**Coverage: 0/2 routes (0%)**

---

## Client Tags (5 routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `ClientTagListURL` | `/app/clients/settings/tags/list` | page | `client-tags-crud.spec.ts` | `[x]` retail |
| `ClientTagAddURL` | `/action/clients/tags/add` | action | `client-tags-crud.spec.ts` | `[x]` retail |
| `ClientTagEditURL` | `/action/clients/tags/edit/{id}` | action | `client-tags-crud.spec.ts` | `[x]` retail |
| `ClientTagDeleteURL` | `/action/clients/tags/delete` | action | `client-tags-crud.spec.ts` | `[x]` retail |
| `ClientTagBulkDeleteURL` | `/action/clients/tags/bulk-delete` | action | — | `[ ]` |

**Coverage: 4/5 routes (80%)**

---

## Auth (9 routes)

| Route Constant | URL | Type | Test File | Status |
|---------------|-----|------|-----------|--------|
| `AuthLoginURL` | `/auth/login` | page | — | `[ ]` |
| `AuthLoginPostURL` | `/auth/login` | action | — | `[ ]` |
| `AuthSignupURL` | `/auth/signup` | page | — | `[ ]` |
| `AuthSignupPostURL` | `/auth/signup` | action | — | `[ ]` |
| `AuthResetPasswordURL` | `/auth/reset-password` | page | — | `[ ]` |
| `AuthResetPasswordPostURL` | `/auth/reset-password` | action | — | `[ ]` |
| `AuthResetConfirmURL` | `/auth/reset-password/confirm` | page | — | `[ ]` |
| `AuthResetConfirmPostURL` | `/auth/reset-password/confirm` | action | — | `[ ]` |
| `AuthLogoutURL` | `/auth/logout` | action | — | `[ ]` |

**Coverage: 0/9 routes (0%)** — Auth tests are app-level (not package-level)

---

## Summary

| Entity | Total Routes | Covered | Partial | Not Started | Coverage |
|--------|-------------|---------|---------|-------------|----------|
| Clients | 13 | 3 | 1 | 9 | 31% |
| Users | 14 | 7 | 2 | 5 | 50% |
| Roles | 16 | 7 | 4 | 5 | 50% |
| Role Permissions | 4 | 1 | 2 | 1 | 75% |
| User Roles | 4 | 1 | 2 | 1 | 75% |
| Locations | 13 | 2 | 1 | 10 | 23% |
| Permissions | 8 | 0 | 0 | 8 | 0% |
| Workspaces | 8 | 0 | 0 | 8 | 0% |
| Suppliers | 12 | 0 | 0 | 12 | 0% |
| Reports | 2 | 0 | 0 | 2 | 0% |
| Client Tags | 5 | 4 | 0 | 1 | 80% |
| Auth | 9 | 0 | 0 | 9 | 0% (app-level) |
| **Total** | **108** | **25** | **12** | **71** | **23%** |

### Existing Tests by App

| App | Test Files Covering Entydad | Routes Covered |
|-----|----------------------------|----------------|
| retail-admin | 18 spec files (users/, roles/, client-*) | Users, Roles, Clients, Client Tags |
| service-admin | 2 spec files (customers/, locations/) | Clients, Locations |

### Priority Gaps (highest value untested routes)

1. **Suppliers** — 0% coverage, 12 routes, full CRUD + detail
2. **Permissions** — 0% coverage, 8 routes, CRUD
3. **Workspaces** — 0% coverage, 8 routes, CRUD
4. **Locations detail** — detail page, tabs, attachments not tested
5. **Client delete/status** — CRUD partially covered (no delete, no status change)
6. **Auth** — 0% coverage but these are app-level tests

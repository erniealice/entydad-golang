import { test, expect } from '@playwright/test';

/**
 * ENT-PRM-001: Permission List
 * ENT-PRM-002: Permission Table Structure
 *
 * Routes: PermissionListURL
 * Verifies: list page loads, table structure, column headers, data rows
 *
 * NOTE: Permissions list has edit/deactivate/delete per row but NO view action.
 */

test.describe('ENT-PRM-001: Permission List', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/permissions/list/active');
    await expect(page.locator('#permissions-table')).toBeVisible();
  });

  test('displays permission table with correct column headers', async ({ page }) => {
    const headers = page.locator('thead th .column-label');
    const count = await headers.count();
    // Name, Entity, Permission Code, Type, Status = 5 data columns
    expect(count).toBeGreaterThanOrEqual(5);
  });

  test('shows data rows with permission data', async ({ page }) => {
    const rows = page.locator('#permissions-table tbody tr[data-id]');
    const count = await rows.count();
    expect(count).toBeGreaterThanOrEqual(1);
  });

  test('has primary action button in toolbar', async ({ page }) => {
    const primaryAction = page.locator('.toolbar-primary-action');
    await expect(primaryAction).toBeVisible();
    await expect(primaryAction).toBeEnabled();
  });

  test('shows pagination with entry count', async ({ page }) => {
    const pagination = page.locator('.table-footer, .pagination-info');
    await expect(pagination).toBeVisible();
  });

  test('row has action buttons (edit, deactivate, delete)', async ({ page }) => {
    const firstRow = page.locator('#permissions-table tbody tr[data-id]').first();
    const editBtn = firstRow.locator('.action-btn.edit');
    const deactivateBtn = firstRow.locator('.action-btn.deactivate');
    const deleteBtn = firstRow.locator('.action-btn.delete');

    await expect(editBtn).toBeVisible();
    await expect(deactivateBtn).toBeVisible();
    await expect(deleteBtn).toBeVisible();
  });

  test('permissions list has no view action per row', async ({ page }) => {
    // Permissions do not have a view/detail page
    const firstRow = page.locator('#permissions-table tbody tr[data-id]').first();
    const viewLink = firstRow.locator('.action-btn.view');
    await expect(viewLink).toHaveCount(0);
  });
});

test.describe('ENT-PRM-002: Permission Table Data', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/permissions/list/active');
    await expect(page.locator('#permissions-table')).toBeVisible();
  });

  test('permission rows have expected data attributes', async ({ page }) => {
    const firstRow = page.locator('#permissions-table tbody tr[data-id]').first();
    const dataId = await firstRow.getAttribute('data-id');
    expect(dataId).toBeTruthy();
    expect(dataId!.startsWith('perm-')).toBeTruthy();
  });

  test('has multiple permission rows (seed data)', async ({ page }) => {
    const rows = page.locator('#permissions-table tbody tr[data-id]');
    const count = await rows.count();
    // Service-admin seeds many permissions
    expect(count).toBeGreaterThanOrEqual(5);
  });
});

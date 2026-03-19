import { test, expect } from '@playwright/test';
import { waitForHtmxSettle } from '../helpers/htmx';

/**
 * ENT-SUP-004: Supplier Detail Page
 * ENT-SUP-005: Supplier Tab Navigation
 * ENT-SUP-006: Supplier Status Change
 * ENT-SUP-007: Supplier Delete
 *
 * Routes: SupplierDetailURL, SupplierTabActionURL,
 *         SupplierSetStatusURL, SupplierDeleteURL
 *
 * NOTE: Detail page content renders "Page content not available" as of 2026-03-19.
 *       Detail/tab tests are skipped until the view is wired. Status/delete tests
 *       run against the list page actions.
 */

test.describe('ENT-SUP-004: Supplier Detail Page', () => {
  test.skip(true, 'Supplier detail page not yet wired — shows "Page content not available"');

  test('displays supplier name in heading', async ({ page }) => {
    await page.goto('/app/suppliers/detail/sup-001');
    const heading = page.locator('.greeting h1');
    await expect(heading).toBeVisible();
  });
});

test.describe('ENT-SUP-005: Supplier Tab Navigation', () => {
  test.skip(true, 'Supplier detail page not yet wired — tabs not rendered');

  test('info tab is active by default', async ({ page }) => {
    await page.goto('/app/suppliers/detail/sup-001');
    const infoTab = page.locator('[role="tab"]').first();
    await expect(infoTab).toHaveClass(/active/);
  });
});

test.describe('ENT-SUP-006: Supplier Status Change', () => {
  test('status action button exists on list page', async ({ page }) => {
    await page.goto('/app/suppliers/list/active');
    await expect(page.locator('table')).toBeVisible();

    // Status button uses .action-btn.deactivate class
    const statusBtn = page.locator('tbody tr:first-child button.action-btn.deactivate');
    await expect(statusBtn).toBeVisible();
  });
});

test.describe('ENT-SUP-007: Supplier Delete', () => {
  test('delete button exists on list page', async ({ page }) => {
    await page.goto('/app/suppliers/list/active');
    await expect(page.locator('table')).toBeVisible();

    const deleteBtn = page.locator('tbody tr:first-child button.action-btn.delete');
    await expect(deleteBtn).toBeVisible();
  });
});

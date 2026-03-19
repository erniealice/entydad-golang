import { test, expect } from '@playwright/test';
import { waitForHtmxSettle } from '../helpers/htmx';

/**
 * ENT-LOC-001: Location List
 * ENT-LOC-002: Location Add via Drawer
 * ENT-LOC-003: Location Table Structure (empty state / data)
 *
 * Routes: LocationListURL, LocationAddURL
 * Verifies: list page loads, table structure, add via drawer
 *
 * NOTE: Locations list shows per-row action buttons (view/edit/deactivate/delete) when data rows exist.
 * Empty table shows empty state with no row actions.
 */

test.describe('ENT-LOC-001: Location List', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/locations/list/active');
    await expect(page.locator('#locations-table')).toBeVisible();
  });

  test('displays location table with correct column headers', async ({ page }) => {
    const headers = page.locator('thead th .column-label');
    const count = await headers.count();
    // Name, Address, Status = 3 data columns minimum
    expect(count).toBeGreaterThanOrEqual(3);
  });

  test('has primary action button in toolbar', async ({ page }) => {
    const primaryAction = page.locator('.toolbar-primary-action');
    await expect(primaryAction).toBeVisible();
    await expect(primaryAction).toBeEnabled();
  });

  test('shows table footer with entry count', async ({ page }) => {
    const footer = page.locator('#locations-table-footer');
    await expect(footer).toBeVisible();
  });

  test('shows empty state or data rows', async ({ page }) => {
    // Locations table may be empty (no seed data) or have data
    const dataRows = page.locator('#locations-table tbody tr[data-id]');
    const emptyState = page.locator('#locations-table tbody .empty-state');
    const dataCount = await dataRows.count();
    const emptyCount = await emptyState.count();

    // One of them should be present
    expect(dataCount + emptyCount).toBeGreaterThanOrEqual(1);
  });
});

test.describe('ENT-LOC-002: Location Add via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/locations/list/active');
    await expect(page.locator('#locations-table')).toBeVisible();
  });

  test('opens drawer when primary action clicked', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Verify form fields exist by ID
    await expect(page.locator('#name')).toBeVisible();
    await expect(page.locator('#address')).toBeVisible();
    await expect(page.locator('#description')).toBeVisible();
  });

  test('creates location via drawer form', async ({ page }) => {
    const ts = Date.now();

    // Open drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Fill required fields
    await page.locator('#name').fill(`E2E Location ${ts}`);
    await page.locator('#address').fill(`123 Test Street ${ts}`);
    await page.locator('#description').fill(`Test location created by E2E at ${ts}`);

    // Submit
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();

    // Wait for HTMX response
    await waitForHtmxSettle(page);

    // Verify drawer closes
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 10000 });
  });

  test('cancel closes drawer without creating', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    await page.locator('#name').fill('ShouldNotSave');

    // Cancel
    await page.locator('#sheet .sheet-footer .btn-secondary').click();

    await expect(page.locator('#sheet').first()).not.toHaveClass(/open/, { timeout: 5000 });
  });
});

test.describe('ENT-LOC-003: Location Row Actions', () => {
  test('location rows have action buttons when data exists', async ({ page }) => {
    const ts = Date.now();

    // Ensure at least one location exists by creating one
    await page.goto('/app/locations/list/active');
    await expect(page.locator('#locations-table')).toBeVisible();

    const dataRows = page.locator('#locations-table tbody tr[data-id]');
    const initialCount = await dataRows.count();

    if (initialCount === 0) {
      // No rows — create one via the add drawer so we can test action buttons
      await page.locator('.toolbar-primary-action').click();
      await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

      await page.locator('#name').fill(`E2E LOC Actions ${ts}`);
      await page.locator('#address').fill(`123 Actions St ${ts}`);

      await page.locator('#sheet .sheet-footer button[type="submit"]').click();
      await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 10000 });

      // Navigate back to reload the table
      await page.goto('/app/locations/list/active');
      await expect(page.locator('#locations-table')).toBeVisible();
    }

    const finalRows = page.locator('#locations-table tbody tr[data-id]');
    const finalCount = await finalRows.count();
    expect(finalCount).toBeGreaterThan(0);

    // Locations have view, edit, deactivate, delete per row
    const firstRow = finalRows.first();
    const viewBtn = firstRow.locator('.action-btn.view');
    const editBtn = firstRow.locator('.action-btn.edit');
    const deleteBtn = firstRow.locator('.action-btn.delete');

    await expect(viewBtn).toBeVisible();
    await expect(editBtn).toBeVisible();
    await expect(deleteBtn).toBeVisible();
  });
});

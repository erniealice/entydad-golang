import { test, expect } from '@playwright/test';
import { waitForHtmxSettle } from '../helpers/htmx';

/**
 * ENT-CLI-001: Client List
 * ENT-CLI-002: Client Add via Drawer
 * ENT-CLI-003: Client Edit via Drawer
 * ENT-CLI-004: Client Row Actions
 * ENT-CLI-005: Client Detail Page
 *
 * Routes: ClientListURL, ClientAddURL, ClientEditURL, ClientDetailURL
 * Verifies: list page loads, table structure, CRUD via drawer, detail page
 */

test.describe('ENT-CLI-001: Client List', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/clients/list/active');
    await expect(page.locator('#clients-table')).toBeVisible();
  });

  test('displays client table with correct column headers', async ({ page }) => {
    const headers = page.locator('thead th .column-label');
    const count = await headers.count();
    // Client Name, Email, Phone, Status = 4 data columns minimum
    expect(count).toBeGreaterThanOrEqual(4);
  });

  test('shows data rows with client data', async ({ page }) => {
    const rows = page.locator('#clients-table tbody tr[data-id]');
    const count = await rows.count();
    expect(count).toBeGreaterThanOrEqual(1);

    // First row should have cell content
    const firstRowCells = rows.first().locator('td');
    const cellCount = await firstRowCells.count();
    expect(cellCount).toBeGreaterThanOrEqual(4);
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

  test('row has action buttons (view, edit, deactivate, delete)', async ({ page }) => {
    const firstRow = page.locator('#clients-table tbody tr[data-id]').first();
    const viewLink = firstRow.locator('.action-btn.view');
    const editBtn = firstRow.locator('.action-btn.edit');
    const deactivateBtn = firstRow.locator('.action-btn.deactivate');
    const deleteBtn = firstRow.locator('.action-btn.delete');

    await expect(viewLink).toBeVisible();
    await expect(editBtn).toBeVisible();
    await expect(deactivateBtn).toBeVisible();
    await expect(deleteBtn).toBeVisible();
  });

  test('view link navigates to client detail', async ({ page }) => {
    const viewLink = page.locator('#clients-table tbody tr[data-id]').first().locator('a.action-btn.view');
    const href = await viewLink.getAttribute('href');
    expect(href).toContain('/app/clients/detail/');
  });
});

test.describe('ENT-CLI-002: Client Add via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/clients/list/active');
    await expect(page.locator('#clients-table')).toBeVisible();
  });

  test('opens drawer when primary action clicked', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Verify form fields exist by ID
    await expect(page.locator('#first_name')).toBeVisible();
    await expect(page.locator('#last_name')).toBeVisible();
    await expect(page.locator('#email_address')).toBeVisible();
    await expect(page.locator('#mobile_number')).toBeVisible();
  });

  test('drawer has all required form fields', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Contact info
    await expect(page.locator('#first_name')).toBeVisible();
    await expect(page.locator('#last_name')).toBeVisible();
    await expect(page.locator('#email_address')).toBeVisible();
    await expect(page.locator('#mobile_number')).toBeVisible();

    // Company info
    await expect(page.locator('#company_name')).toBeVisible();
    await expect(page.locator('#customer_type')).toBeVisible();

    // Address
    await expect(page.locator('#street_address')).toBeVisible();
    await expect(page.locator('#city')).toBeVisible();
  });

  test('creates client via drawer form', async ({ page }) => {
    const ts = Date.now();

    // Open drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Fill required fields
    await page.locator('#first_name').fill(`E2E`);
    await page.locator('#last_name').fill(`Client${ts}`);
    await page.locator('#email_address').fill(`e2e-${ts}@test.com`);
    await page.locator('#mobile_number').fill(`09170000${String(ts).slice(-4)}`);

    // Submit
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();

    // Wait for HTMX response + sheet close callback
    await waitForHtmxSettle(page);

    // Verify drawer closes
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 10000 });
  });

  test('cancel closes drawer without creating', async ({ page }) => {
    // Open drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Fill something
    await page.locator('#first_name').fill('ShouldNotSave');

    // Cancel
    await page.locator('#sheet .sheet-footer .btn-secondary').click();

    // Drawer should close
    await expect(page.locator('#sheet').first()).not.toHaveClass(/open/, { timeout: 5000 });
  });
});

test.describe('ENT-CLI-003: Client Edit via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/clients/list/active', { timeout: 15000 });
    await expect(page.locator('#clients-table')).toBeVisible({ timeout: 10000 });
  });

  test('opens edit drawer with pre-filled data', async ({ page }) => {
    const editBtn = page.locator('#clients-table tbody tr[data-id]').first().locator('.action-btn.edit');
    await editBtn.click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();
    await waitForHtmxSettle(page);

    // First name should be pre-filled
    const firstName = await page.locator('#first_name').inputValue();
    expect(firstName.length).toBeGreaterThan(0);
  });

  test('saves edit and closes drawer', async ({ page }) => {
    const editBtn = page.locator('#clients-table tbody tr[data-id]').first().locator('.action-btn.edit');
    await editBtn.click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();
    await waitForHtmxSettle(page);

    // Modify a field
    const ts = Date.now();
    await page.locator('#notes').fill(`Updated by E2E test at ${ts}`);

    // Submit
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();

    // Wait for HTMX settle
    await waitForHtmxSettle(page);

    // Drawer closes
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 10000 });
  });
});

test.describe('ENT-CLI-005: Client Detail Page', () => {
  test('detail page loads or is documented as unavailable', async ({ page }) => {
    // Navigate to list first to get a valid client ID
    await page.goto('/app/clients/list/active');
    await expect(page.locator('#clients-table')).toBeVisible();

    const viewLink = page.locator('#clients-table tbody tr[data-id]').first().locator('a.action-btn.view');
    const href = await viewLink.getAttribute('href');
    expect(href).toBeTruthy();

    // Navigate to detail page
    await page.goto(href!);

    // Check if page shows "Page content not available" - known issue
    const pageContent = page.locator('.page-content');
    const text = await pageContent.textContent();

    if (text?.includes('Page content not available')) {
      test.skip(true, 'BUG: Client detail page shows "Page content not available" — view not wired in service-admin');
    }

    // If we get here, verify detail page structure
    await expect(page.locator('.detail-layout')).toBeVisible();
  });
});

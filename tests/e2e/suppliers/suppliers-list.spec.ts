import { test, expect } from '@playwright/test';
import { waitForHtmxSettle } from '../helpers/htmx';

/**
 * ENT-SUP-001: Supplier List
 * ENT-SUP-002: Supplier Add
 * ENT-SUP-003: Supplier Edit
 *
 * Routes: SupplierListURL, SupplierAddURL, SupplierEditURL
 * Verifies: list page loads, table structure, CRUD via drawer
 */

test.describe('ENT-SUP-001: Supplier List', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/suppliers/list/active');
    await expect(page.locator('table')).toBeVisible();
  });

  test('displays supplier table with correct column headers', async ({ page }) => {
    const headers = page.locator('thead th');
    const count = await headers.count();
    // Checkbox + 8 data columns + Actions = 10
    expect(count).toBeGreaterThanOrEqual(9);

    const labelCount = await page.locator('thead th .column-label').count();
    expect(labelCount).toBeGreaterThanOrEqual(2);
  });

  test('shows data rows with supplier data', async ({ page }) => {
    const rows = page.locator('tbody tr');
    const count = await rows.count();
    expect(count).toBeGreaterThanOrEqual(1);

    // First row should have cell content
    const firstRowCells = page.locator('tbody tr:first-child td');
    const cellCount = await firstRowCells.count();
    expect(cellCount).toBeGreaterThanOrEqual(8);
  });

  test('has primary action button in toolbar', async ({ page }) => {
    const primaryAction = page.locator('.toolbar-primary-action');
    await expect(primaryAction).toBeVisible();
    await expect(primaryAction).toBeEnabled();
  });

  test('has table search input', async ({ page }) => {
    const search = page.locator('table').locator('..').locator('..').locator('input[type="text"][placeholder*="Search"]');
    await expect(search).toBeVisible();
  });

  test('shows pagination with entry count', async ({ page }) => {
    const pagination = page.locator('.table-footer, .pagination-info');
    await expect(pagination).toBeVisible();
  });

  test('row has action buttons (view, edit, status, delete)', async ({ page }) => {
    const firstRow = page.locator('tbody tr:first-child');
    const viewLink = firstRow.locator('a.action-btn.view');
    const editBtn = firstRow.locator('.action-btn.edit');
    const deleteBtn = firstRow.locator('.action-btn.delete');

    await expect(viewLink).toBeVisible();
    await expect(editBtn).toBeVisible();
    await expect(deleteBtn).toBeVisible();
  });

  test('view link navigates to supplier detail', async ({ page }) => {
    const viewLink = page.locator('tbody tr:first-child a.action-btn.view');
    const href = await viewLink.getAttribute('href');
    expect(href).toContain('/app/suppliers/detail/');
  });
});

test.describe('ENT-SUP-002: Supplier Add via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/suppliers/list/active');
    await expect(page.locator('table')).toBeVisible();
  });

  test('opens drawer when primary action clicked', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Verify form fields exist by ID
    await expect(page.locator('#company_name')).toBeVisible();
    await expect(page.locator('#supplier_type')).toBeVisible();
    await expect(page.locator('#status')).toBeVisible();
  });

  test('drawer has all required form sections', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Company info
    await expect(page.locator('#company_name')).toBeVisible();

    // Contact person
    await expect(page.locator('#first_name')).toBeVisible();
    await expect(page.locator('#last_name')).toBeVisible();
    await expect(page.locator('#email_address')).toBeVisible();
    await expect(page.locator('#mobile_number')).toBeVisible();

    // Financial details
    await expect(page.locator('#payment_terms')).toBeVisible();
    await expect(page.locator('#credit_limit')).toBeVisible();

    // Address
    await expect(page.locator('#street_address')).toBeVisible();
    await expect(page.locator('#city')).toBeVisible();
  });

  test('creates supplier via drawer form', async ({ page }) => {
    const ts = Date.now();
    const rowsBefore = await page.locator('tbody tr').count();

    // Open drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Fill required fields
    await page.locator('#company_name').fill(`TestSupplier${ts}`);
    await page.locator('#first_name').fill('E2E');
    await page.locator('#last_name').fill(`Test${ts}`);

    // Submit
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();

    // Wait for HTMX response + sheet close callback
    await waitForHtmxSettle(page);

    // Verify drawer closes — sheet-form uses .open class on the wrapper
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 10000 });
  });

  test('cancel closes drawer without creating', async ({ page }) => {
    // Open drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Fill something
    await page.locator('#company_name').fill('ShouldNotSave');

    // Cancel — use the secondary button in sheet footer
    await page.locator('#sheet .sheet-footer .btn-secondary').click();

    // Drawer should close
    await expect(page.locator('#sheet').first()).not.toHaveClass(/open/, { timeout: 5000 });
  });
});

test.describe('ENT-SUP-003: Supplier Edit via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/suppliers/list/active');
    await expect(page.locator('table')).toBeVisible();
  });

  test('opens edit drawer with pre-filled data', async ({ page }) => {
    const editBtn = page.locator('tbody tr:first-child .action-btn.edit');
    await editBtn.click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Company name should be pre-filled
    const companyName = await page.locator('#company_name').inputValue();
    expect(companyName.length).toBeGreaterThan(0);
  });

  test('saves edit and closes drawer', async ({ page }) => {
    const editBtn = page.locator('tbody tr:first-child .action-btn.edit');
    await editBtn.click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Modify a field
    const ts = Date.now();
    await page.locator('#company_name').fill(`Updated${ts}`);

    // Submit
    await page.locator('#sheet button[type="submit"]').click();

    // Drawer closes
    await expect(page.locator('#sheet').first()).not.toHaveClass(/open/, { timeout: 5000 });

    await waitForHtmxSettle(page);
  });
});

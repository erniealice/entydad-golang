// Package block — workspace helpers and statement translators shared across domain wirings.
//
// This file holds:
//   - getDefaultWorkspaceID: small, stateless workspace helper.
//   - Statement request/response translation helpers — bridge the new
//     service.reporting.v1 proto package (used by the typed
//     `uc.Reports.Statements.*` closures) and the legacy
//     domain.ledger/treasury.v1 proto package consumed by entydad view
//     deps. Both shapes are field-for-field identical so the translation
//     is mechanical.
//
// (The former DataSource/UpdateableSource/CRUDSource ducks were deleted
// 2026-06-12; active/status writes + payment_term/client reads now flow through
// the narrow typed UseCases primitives/closures bound by service-admin.)
//
// (The former categoryListPageDataGetter local interface was deleted
// 2026-06-14 as part of the E6 provider hardcode fix — category list page
// data is now wired through the typed UseCases.Category.GetListPageData
// closure bound by service-admin.)
//
// Rule of thumb: 1 caller → live with the caller; 2+ callers → live here.
package block

import (
	"os"

	clientstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/client_statement"
	suppstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
	stmtspb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/statements"
)

// translateClientStatementReq translates the legacy
// `domain.ledger/client_statement.v1` request shape consumed by entydad
// views into the service-layer `service.reporting.v1` shape that the
// `uc.Reports.Statements.GetClientStatement` closure expects.
func translateClientStatementReq(req *clientstmtpb.ClientStatementRequest) *stmtspb.GetClientStatementRequest {
	if req == nil {
		return nil
	}
	return &stmtspb.GetClientStatementRequest{
		ClientId:   req.GetClientId(),
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Currency:   req.Currency,
		Pagination: req.GetPagination(),
	}
}

// translateClientStatementResp translates the service-layer response back
// into the legacy `domain.ledger/client_statement.v1` response shape that
// entydad view deps consume.
func translateClientStatementResp(resp *stmtspb.GetClientStatementResponse) *clientstmtpb.ClientStatementResponse {
	if resp == nil {
		return nil
	}
	out := &clientstmtpb.ClientStatementResponse{
		Success: resp.GetSuccess(),
	}
	for _, e := range resp.GetEntries() {
		if e == nil {
			continue
		}
		out.Entries = append(out.Entries, &clientstmtpb.StatementEntry{
			Date:            e.GetDate(),
			Type:            e.GetType(),
			ReferenceNumber: e.GetReferenceNumber(),
			Description:     e.GetDescription(),
			Billed:          e.GetBilled(),
			Received:        e.GetReceived(),
			Balance:         e.GetBalance(),
			EntityId:        e.GetEntityId(),
			Status:          e.GetStatus(),
		})
	}
	if s := resp.GetSummary(); s != nil {
		out.Summary = &clientstmtpb.ClientStatementSummary{
			TotalBilled:        s.GetTotalBilled(),
			TotalReceived:      s.GetTotalReceived(),
			OutstandingBalance: s.GetOutstandingBalance(),
			InvoiceCount:       s.GetInvoiceCount(),
			CollectionCount:    s.GetCollectionCount(),
			Currency:           s.GetCurrency(),
			StartDate:          s.StartDate,
			EndDate:            s.EndDate,
			ClientName:         s.GetClientName(),
		}
	}
	out.Pagination = resp.GetPagination()
	out.Error = resp.GetError()
	return out
}

// translateSupplierStatementReq translates the legacy
// `domain.treasury/supplier_statement.v1` request shape into the
// service-layer shape.
func translateSupplierStatementReq(req *suppstmtpb.SupplierStatementRequest) *stmtspb.GetSupplierStatementRequest {
	if req == nil {
		return nil
	}
	return &stmtspb.GetSupplierStatementRequest{
		SupplierId: req.GetSupplierId(),
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Currency:   req.Currency,
		Pagination: req.GetPagination(),
	}
}

// translateSupplierStatementResp translates the service-layer response
// back into the legacy `domain.treasury/supplier_statement.v1` shape.
func translateSupplierStatementResp(resp *stmtspb.GetSupplierStatementResponse) *suppstmtpb.SupplierStatementResponse {
	if resp == nil {
		return nil
	}
	out := &suppstmtpb.SupplierStatementResponse{
		Success: resp.GetSuccess(),
	}
	for _, e := range resp.GetEntries() {
		if e == nil {
			continue
		}
		out.Entries = append(out.Entries, &suppstmtpb.SupplierStatementEntry{
			Date:            e.GetDate(),
			Type:            e.GetType(),
			ReferenceNumber: e.GetReferenceNumber(),
			Description:     e.GetDescription(),
			Billed:          e.GetBilled(),
			Paid:            e.GetPaid(),
			Balance:         e.GetBalance(),
			EntityId:        e.GetEntityId(),
			Status:          e.GetStatus(),
		})
	}
	if s := resp.GetSummary(); s != nil {
		out.Summary = &suppstmtpb.SupplierStatementSummary{
			TotalBilled:        s.GetTotalBilled(),
			TotalPaid:          s.GetTotalPaid(),
			OutstandingBalance: s.GetOutstandingBalance(),
			BillCount:          s.GetBillCount(),
			PaymentCount:       s.GetPaymentCount(),
			Currency:           s.GetCurrency(),
			StartDate:          s.StartDate,
			EndDate:            s.EndDate,
			SupplierName:       s.GetSupplierName(),
		}
	}
	out.Pagination = resp.GetPagination()
	out.Error = resp.GetError()
	return out
}

// getDefaultWorkspaceID returns the default workspace ID from the environment,
// falling back to "default-workspace" if the env var is not set.
func getDefaultWorkspaceID() string {
	if v := os.Getenv("DEFAULT_WORKSPACE_ID"); v != "" {
		return v
	}
	return "default-workspace"
}

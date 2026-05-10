// Package block — interface contracts and workspace helpers shared across domain wirings.
//
// This file holds:
//   - DB interface types (categoryListPageDataGetter, UpdateableSource, CRUDSource)
//     used by Block() when type-asserting ctx.DB.
//   - getDefaultWorkspaceID: small, stateless workspace helper.
//
// Rule of thumb: 1 caller → live with the caller; 2+ callers → live here.
// If this file grows past ~150 lines, consider a more specific name.
package block

import (
	"context"
	"os"

	"github.com/erniealice/entydad-golang"
	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/client_statement"
	suppstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
)

// LedgerReportingService is the subset of the espyna consumer.LedgerReportingService
// interface that entydad's block actually calls. The concrete implementation
// (espyna's postgres adapter) satisfies this interface implicitly.
// Defined here so the block package does not need to import consumer.
type LedgerReportingService interface {
	GetClientStatement(ctx context.Context, req *clientstmtpb.ClientStatementRequest) (*clientstmtpb.ClientStatementResponse, error)
	GetClientBalances(ctx context.Context) (map[string]int64, error)
	GetSupplierStatement(ctx context.Context, req *suppstmtpb.SupplierStatementRequest) (*suppstmtpb.SupplierStatementResponse, error)
	GetSupplierBalances(ctx context.Context) (map[string]int64, error)
}

// categoryListPageDataGetter is a local interface satisfied by the PostgresCategoryRepository
// concrete type, allowing GetCategoryListPageData to be called via type assertion without
// importing the espyna postgres adapter package.
type categoryListPageDataGetter interface {
	GetCategoryListPageData(ctx context.Context) ([]*categorypb.Category, error)
}

// UpdateableSource extends entydad.DataSource with the Update method that
// SetActive closures need. espyna's DatabaseAdapter satisfies this interface.
type UpdateableSource interface {
	entydad.DataSource
	Update(ctx context.Context, collection, id string, data map[string]any) (map[string]any, error)
}

// CRUDSource extends UpdateableSource with Create, Read, and Delete operations.
// espyna's DatabaseAdapter satisfies this interface. Used by simpler entities
// (e.g. LocationArea) that do not yet have dedicated proto service use-cases.
type CRUDSource interface {
	UpdateableSource
	Create(ctx context.Context, collection string, data map[string]any) (map[string]any, error)
	Read(ctx context.Context, collection, id string) (map[string]any, error)
	Delete(ctx context.Context, collection, id string) error
}

// getDefaultWorkspaceID returns the default workspace ID from the environment,
// falling back to "default-workspace" if the env var is not set.
func getDefaultWorkspaceID() string {
	if v := os.Getenv("DEFAULT_WORKSPACE_ID"); v != "" {
		return v
	}
	return "default-workspace"
}

// Package list provides the polymorphic Tax Registration list view.
// Used on both Client and Workspace detail pages (party_type scoped).
package list

import (
	"context"
	"log"

	taxregistration "github.com/erniealice/entydad-golang/domain/tax/tax_registration"
	taxregistrationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_registration"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the tax registration list view.
type Deps struct {
	// PartyType and PartyID identify which party's registrations to show.
	// Populated from URL path parameters at route registration time.
	PartyType string
	PartyID   string

	Routes       taxregistration.Routes
	Labels       taxregistration.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Tax registration use cases
	ListTaxRegistrations func(ctx context.Context, req *taxregistrationpb.ListTaxRegistrationsRequest) (*taxregistrationpb.ListTaxRegistrationsResponse, error)
}

// PageData holds the data for the tax registration list.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	Labels          taxregistration.Labels
}

// TaxRegistrationRow is the view-model for a single tax registration row.
type TaxRegistrationRow struct {
	ID                 string
	KindName           string
	ComputePath        string
	PartyRole          string
	Status             string
	EffectiveFrom      string
	RegistrationNumber string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the tax registration list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("tax_registration", "list") {
			return view.Forbidden("tax_registration:list")
		}

		// Party ID may come from URL path parameter if not pre-populated in deps.
		partyID := deps.PartyID
		if partyID == "" {
			partyID = viewCtx.Request.PathValue("id")
		}

		rows := fetchRegistrations(ctx, deps, partyID)
		tableConfig := buildTableConfig(deps, rows, partyID, perms)

		heading := headingForPartyType(deps.Labels, deps.PartyType)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        heading,
				CurrentPath:  viewCtx.CurrentPath,
				HeaderTitle:  heading,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "tax-registration-list-content",
			Table:           tableConfig,
			Labels:          deps.Labels,
		}

		return view.OK("tax-registration-list", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchRegistrations(ctx context.Context, deps *Deps, partyID string) []TaxRegistrationRow {
	if deps.ListTaxRegistrations == nil {
		return []TaxRegistrationRow{}
	}

	resp, err := deps.ListTaxRegistrations(ctx, &taxregistrationpb.ListTaxRegistrationsRequest{})
	if err != nil {
		log.Printf("ListTaxRegistrations error: %v", err)
		return []TaxRegistrationRow{}
	}
	if resp == nil {
		return []TaxRegistrationRow{}
	}

	rows := make([]TaxRegistrationRow, 0)
	for _, tr := range resp.GetData() {
		// Filter by party_id if provided
		if partyID != "" && tr.GetPartyId() != partyID {
			continue
		}
		rows = append(rows, protoToRow(tr))
	}
	return rows
}

func protoToRow(tr *taxregistrationpb.TaxRegistration) TaxRegistrationRow {
	return TaxRegistrationRow{
		ID:                 tr.GetId(),
		KindName:           tr.GetTaxRegistrationKindId(), // Kind name resolved via ID lookup in Phase 5
		ComputePath:        computePathLabel(tr.GetComputePathSnapshot()),
		PartyRole:          partyRoleLabel(tr.GetPartyRoleSnapshot()),
		Status:             statusString(tr.GetStatus()),
		EffectiveFrom:      tr.GetEffectiveFrom(),
		RegistrationNumber: tr.GetRegistrationNumber(),
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, rows []TaxRegistrationRow, partyID string, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "kind", Label: l.Columns.KindName},
		{Key: "compute_path", Label: l.Columns.ComputePath, WidthClass: "col-2xl"},
		{Key: "party_role", Label: l.Columns.PartyRole, WidthClass: "col-2xl"},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
		{Key: "effective_from", Label: l.Columns.EffectiveFrom, WidthClass: "col-3xl"},
		{Key: "reg_number", Label: l.Columns.RegistrationNumber, WidthClass: "col-3xl"},
	}

	canCreate := perms.Can("tax_registration", "create")
	canDelete := perms.Can("tax_registration", "delete")

	tableRows := []types.TableRow{}
	for _, r := range rows {
		actions := []types.TableAction{
			{
				Type:            "delete",
				Label:           l.Actions.Delete,
				Action:          "delete",
				Href:            deps.Routes.DeleteURL,
				Disabled:        !canDelete,
				DisabledTooltip: l.Actions.NoPermission,
			},
		}

		computeVariant := computePathVariant(r.ComputePath)

		tableRows = append(tableRows, types.TableRow{
			ID: r.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: r.KindName},
				{Type: "badge", Value: r.ComputePath, Variant: computeVariant},
				{Type: "text", Value: r.PartyRole},
				{Type: "text", Value: r.Status},
				{Type: "text", Value: r.EffectiveFrom},
				{Type: "text", Value: r.RegistrationNumber},
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, tableRows)

	tableConfig := &types.TableConfig{
		ID:          "tax-registrations-table",
		Columns:     columns,
		Rows:        tableRows,
		ShowSearch:  false,
		ShowActions: true,
		ShowEntries: false,
		Labels:      deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !canCreate,
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func headingForPartyType(l taxregistration.Labels, partyType string) string {
	switch partyType {
	case "client":
		return l.Page.HeadingClient
	case "workspace":
		return l.Page.HeadingWorkspace
	default:
		return l.Page.Heading
	}
}

func statusString(s taxregistrationpb.TaxRegistrationStatus) string {
	switch s {
	case taxregistrationpb.TaxRegistrationStatus_TAX_REGISTRATION_STATUS_ACTIVE:
		return "active"
	case taxregistrationpb.TaxRegistrationStatus_TAX_REGISTRATION_STATUS_SUPERSEDED:
		return "superseded"
	case taxregistrationpb.TaxRegistrationStatus_TAX_REGISTRATION_STATUS_CANCELLED:
		return "cancelled"
	default:
		return "active"
	}
}

func computePathLabel(cp taxregistrationpb.TaxRegistrationComputePathSnapshot) string {
	switch cp {
	case taxregistrationpb.TaxRegistrationComputePathSnapshot_TAX_REGISTRATION_COMPUTE_PATH_SNAPSHOT_SURCHARGE:
		return "Surcharge"
	case taxregistrationpb.TaxRegistrationComputePathSnapshot_TAX_REGISTRATION_COMPUTE_PATH_SNAPSHOT_WITHHOLDING:
		return "Withholding"
	case taxregistrationpb.TaxRegistrationComputePathSnapshot_TAX_REGISTRATION_COMPUTE_PATH_SNAPSHOT_PERIODIC_ONLY:
		return "Periodic Only"
	case taxregistrationpb.TaxRegistrationComputePathSnapshot_TAX_REGISTRATION_COMPUTE_PATH_SNAPSHOT_NONE:
		return "None"
	default:
		return "Unknown"
	}
}

func partyRoleLabel(pr taxregistrationpb.TaxRegistrationPartyRoleSnapshot) string {
	switch pr {
	case taxregistrationpb.TaxRegistrationPartyRoleSnapshot_TAX_REGISTRATION_PARTY_ROLE_SNAPSHOT_SELLER:
		return "Seller"
	case taxregistrationpb.TaxRegistrationPartyRoleSnapshot_TAX_REGISTRATION_PARTY_ROLE_SNAPSHOT_BUYER:
		return "Buyer"
	default:
		return "Unknown"
	}
}

func computePathVariant(label string) string {
	switch label {
	case "Surcharge":
		return "navy"
	case "Withholding":
		return "amber"
	case "Periodic Only":
		return "sage"
	default:
		return "default"
	}
}

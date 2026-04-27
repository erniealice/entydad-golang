package detail

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/pyeza-golang/types"
)

// ClientPlanRow represents a single client-scoped Plan row for the Packages tab.
// It is populated by the ListClientPlans callback wired from centymo's ListPlans
// helper filtered to client_id = current client. The PlanDetailURL is the
// standard Plan detail page URL (/app/services/packages/detail/{id}).
type ClientPlanRow struct {
	PlanID          string
	PlanName        string
	RateCardName    string // resolved from the joined PriceSchedule name
	EngagementCount int    // count of active subscriptions on any PricePlan of this Plan
	PlanDetailURL   string
}

// buildPackagesTable constructs a TableConfig for the client Packages tab.
// It calls deps.ListClientPlans and builds a table with three columns:
// Name, Rate card (resolved from joined PriceSchedule), Engagements (count).
// The primary action opens the centymo Plan-add drawer with context=client.
// On error or empty result the table's own empty-state copy is rendered.
func buildPackagesTable(ctx context.Context, deps *DetailViewDeps, clientID, clientName string) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "name", Label: deps.Labels.Detail.Packages.ColumnName, Sortable: true},
		{Key: "rate_card", Label: deps.Labels.Detail.Packages.ColumnRateCard, Sortable: true},
		{Key: "engagements", Label: deps.Labels.Detail.Packages.ColumnEngagements, Sortable: true, WidthClass: "col-2xl"},
	}

	emptyLabel := deps.Labels.Detail.Packages.Empty
	addActionLabel := deps.Labels.Detail.Packages.AddAction

	var rows []types.TableRow

	planRows, err := deps.ListClientPlans(ctx, clientID)
	if err != nil {
		log.Printf("Failed to load client plans for client %s: %v", clientID, err)
		// Render an empty table rather than a hard error — the tab should still load.
		planRows = nil
	}

	for _, p := range planRows {
		engLabel := fmt.Sprintf("%d", p.EngagementCount)

		href := p.PlanDetailURL

		rows = append(rows, types.TableRow{
			ID:   p.PlanID,
			Href: href,
			Cells: []types.TableCell{
				{Type: "text", Value: p.PlanName},
				{Type: "text", Value: p.RateCardName},
				{Type: "text", Value: engLabel},
			},
			DataAttrs: map[string]string{
				"name":      p.PlanName,
				"rate_card": p.RateCardName,
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	tc := &types.TableConfig{
		ID:                   "client-packages-table",
		Columns:              columns,
		Rows:                 rows,
		Labels:               deps.TableLabels,
		ShowSearch:           true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		EmptyState: types.TableEmptyState{
			Title:   addActionLabel,
			Message: emptyLabel,
		},
	}

	// Primary action — opens centymo Plan-add drawer with client context pre-filled.
	// The ?context=client&client_id={cid} query params instruct the drawer to render
	// the client_id field as a read-only badge instead of an editable picker (§6.6).
	if deps.PlanAddURL != "" {
		addURL := deps.PlanAddURL + "?context=client&client_id=" + clientID
		tc.PrimaryAction = &types.PrimaryAction{
			Label:     addActionLabel,
			ActionURL: addURL,
			Icon:      "icon-plus",
		}
	}

	types.ApplyTableSettings(tc)
	return tc
}

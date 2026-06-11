package detail

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/pyeza-golang/types"
)

// ClientPriceScheduleRow represents a single client-scoped PriceSchedule row
// for the PriceSchedules tab. It is populated by the ListClientPriceSchedules
// callback wired from the PriceSchedule use case filtered to
// price_schedule.client_id = current client. DetailURL points to the
// schedule detail page (resolved by the wiring closure in centymo's block
// using the active route config — tier-overridden via lyngua).
//
// Date+time fields are pre-split into date and time strings via
// types.FormatTimestampSplitInTZ so the table can render them via
// DateTimeCellSplit (matching the price_schedule list table format).
type ClientPriceScheduleRow struct {
	ID            string
	Name          string
	DateStartDate string
	DateStartTime string
	DateEndDate   string
	DateEndTime   string
	PlanCount     int
	DetailURL     string
}

// buildPriceSchedulesTable constructs a TableConfig for the client PriceSchedules tab.
// Columns: Name, Start Date, End Date, Plan Count. The primary action opens the
// PriceSchedule-add drawer with ?context=client&client_id={cid} so the new
// schedule cascades the current client onto its client_id field.
func buildPriceSchedulesTable(ctx context.Context, deps *DetailViewDeps, clientID, clientName string) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "name", Label: deps.Labels.Detail.PriceSchedules.ColumnName},
		{Key: "date_start", Label: deps.Labels.Detail.PriceSchedules.ColumnDateStart, WidthClass: "col-2xl"},
		{Key: "date_end", Label: deps.Labels.Detail.PriceSchedules.ColumnDateEnd, WidthClass: "col-2xl"},
		{Key: "plan_count", Label: deps.Labels.Detail.PriceSchedules.ColumnPlanCount, WidthClass: "col-2xl"},
	}

	emptyLabel := deps.Labels.Detail.PriceSchedules.Empty
	addActionLabel := deps.Labels.Detail.PriceSchedules.AddAction

	var rows []types.TableRow

	var scheduleRows []ClientPriceScheduleRow
	if deps.ListClientPriceSchedules != nil {
		loaded, err := deps.ListClientPriceSchedules(ctx, clientID)
		if err != nil {
			log.Printf("Failed to load client price schedules for client %s: %v", clientID, err)
		} else {
			scheduleRows = loaded
		}
	}

	for _, s := range scheduleRows {
		planCountLabel := fmt.Sprintf("%d", s.PlanCount)

		rows = append(rows, types.TableRow{
			ID:   s.ID,
			Href: s.DetailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: s.Name},
				types.DateTimeCellSplit(s.DateStartDate, s.DateStartTime),
				types.DateTimeCellSplit(s.DateEndDate, s.DateEndTime),
				{Type: "text", Value: planCountLabel},
			},
			DataAttrs: map[string]string{
				"name": s.Name,
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	tc := &types.TableConfig{
		ID:                   "client-price-schedules-table",
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

	// Primary action — opens the PriceSchedule-add drawer with client context
	// pre-filled. The ?context=client&client_id={cid} query params instruct the
	// drawer to render the client_id field as a read-only badge instead of an
	// editable picker (§6.6).
	if deps.PriceScheduleAddURL != "" {
		addURL := deps.PriceScheduleAddURL + "?context=client&client_id=" + clientID
		tc.PrimaryAction = &types.PrimaryAction{
			Label:     addActionLabel,
			ActionURL: addURL,
			Icon:      "icon-plus",
		}
	}

	types.ApplyTableSettings(tc)
	return tc
}

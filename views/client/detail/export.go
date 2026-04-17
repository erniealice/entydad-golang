package detail

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"time"

	clientstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/client_statement"
)

// NewStatementExportHandler creates an http.HandlerFunc that exports the
// client statement as CSV. Columns: Date, Type, Reference, Description,
// Billed, Received, Balance.
func NewStatementExportHandler(deps *DetailViewDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing client id", http.StatusBadRequest)
			return
		}

		if deps.GetClientStatement == nil {
			http.Error(w, "statement service unavailable", http.StatusServiceUnavailable)
			return
		}

		resp, err := deps.GetClientStatement(ctx, &clientstmtpb.ClientStatementRequest{
			ClientId: id,
		})
		if err != nil {
			log.Printf("client statement export: failed to get statement for %s: %v", id, err)
			http.Error(w, "failed to generate statement", http.StatusInternalServerError)
			return
		}

		// Set CSV response headers
		filename := fmt.Sprintf("client-statement-%s-%s.csv", id, time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Header row — client uses "Received" instead of "Paid"
		if err := writer.Write([]string{"Date", "Type", "Reference", "Description", "Billed", "Received", "Balance"}); err != nil {
			log.Printf("client statement export: failed to write CSV header: %v", err)
			return
		}

		// Data rows
		for _, entry := range resp.GetEntries() {
			billedStr := ""
			if entry.GetBilled() > 0 {
				billedStr = csvCurrency(entry.GetBilled())
			}
			receivedStr := ""
			if entry.GetReceived() > 0 {
				receivedStr = csvCurrency(entry.GetReceived())
			}
			record := []string{
				entry.GetDate(),
				entry.GetType(),
				entry.GetReferenceNumber(),
				entry.GetDescription(),
				billedStr,
				receivedStr,
				csvCurrency(entry.GetBalance()),
			}
			if err := writer.Write(record); err != nil {
				log.Printf("client statement export: failed to write CSV row: %v", err)
				return
			}
		}

		// Summary/totals row
		if s := resp.GetSummary(); s != nil {
			if err := writer.Write([]string{
				"", "", "", "TOTAL",
				csvCurrency(s.GetTotalBilled()),
				csvCurrency(s.GetTotalReceived()),
				csvCurrency(s.GetOutstandingBalance()),
			}); err != nil {
				log.Printf("client statement export: failed to write CSV totals: %v", err)
				return
			}
		}
	}
}

// csvCurrency formats a centavo integer as a plain decimal string (e.g. "15000.50").
func csvCurrency(centavos int64) string {
	return fmt.Sprintf("%.2f", float64(centavos)/100)
}

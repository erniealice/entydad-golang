package entydad

import (
	"net/http"
	"testing"
)

func TestHTMXSuccess_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tableID      string
		wantStatus   int
		wantTrigger  string
		wantTemplate string
	}{
		{
			name:        "users table",
			tableID:     "users-table",
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"users-table"}`,
		},
		{
			name:        "clients table",
			tableID:     "clients-table",
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"clients-table"}`,
		},
		{
			name:        "empty table ID",
			tableID:     "",
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":""}`,
		},
		{
			name:        "locations table",
			tableID:     "locations-table",
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"locations-table"}`,
		},
		{
			name:        "payment terms table",
			tableID:     "payment-terms-table",
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"payment-terms-table"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := HTMXSuccess(tt.tableID)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if got := res.Headers["HX-Trigger"]; got != tt.wantTrigger {
				t.Fatalf("HX-Trigger = %q, want %q", got, tt.wantTrigger)
			}
			if res.Template != "" {
				t.Fatalf("Template = %q, want empty (header-only response)", res.Template)
			}
		})
	}
}

func TestHTMXError_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		message    string
		wantStatus int
		wantHeader string
	}{
		{
			name:       "simple error message",
			message:    "permission denied",
			wantStatus: http.StatusUnprocessableEntity,
			wantHeader: "permission denied",
		},
		{
			name:       "empty message",
			message:    "",
			wantStatus: http.StatusUnprocessableEntity,
			wantHeader: "",
		},
		{
			name:       "long error message",
			message:    "The requested resource could not be found in the database",
			wantStatus: http.StatusUnprocessableEntity,
			wantHeader: "The requested resource could not be found in the database",
		},
		{
			name:       "message with special characters",
			message:    `error: "field" is required`,
			wantStatus: http.StatusUnprocessableEntity,
			wantHeader: `error: "field" is required`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := HTMXError(tt.message)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("StatusCode = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if got := res.Headers["HX-Error-Message"]; got != tt.wantHeader {
				t.Fatalf("HX-Error-Message = %q, want %q", got, tt.wantHeader)
			}
			if res.Template != "" {
				t.Fatalf("Template = %q, want empty (header-only response)", res.Template)
			}
		})
	}
}

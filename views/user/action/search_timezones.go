package action

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	pyezatypes "github.com/erniealice/pyeza-golang/types"
)

// timezoneOption is the JSON shape returned to the auto-complete component.
type timezoneOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// NewSearchTimezonesAction returns an http.HandlerFunc that filters
// pyezatypes.CommonTimezones by ?q= and returns JSON [{value,label}, ...].
// Empty query returns the full curated list.
func NewSearchTimezonesAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		results := make([]timezoneOption, 0, len(pyezatypes.CommonTimezones))
		for _, tz := range pyezatypes.CommonTimezones {
			if q != "" && !strings.Contains(strings.ToLower(tz), q) {
				continue
			}
			results = append(results, timezoneOption{Value: tz, Label: tz})
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			log.Printf("search timezones: encode failed: %v", err)
		}
	}
}

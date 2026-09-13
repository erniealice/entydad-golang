package location

import (
	"github.com/erniealice/pyeza-golang/route"
	"testing"
)

func TestLocationAppEntryStatus(t *testing.T) {
	unit := Describe()
	entry := unit.Nav.AppEntry
	if got := route.ResolveURL(DefaultRoutes().ListURL, "status", entry.Params["status"]); got != "/locations/list/active" {
		t.Fatalf("app entry must resolve a concrete active list: %q", got)
	}
}

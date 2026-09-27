package location

import (
	"testing"
	"testing/fstest"

	v1 "github.com/erniealice/lyngua/golang/v1"
)

func TestLocationLabelsAcceptProvidedValues(t *testing.T) {
	provider := v1.NewTranslationProviderFromFS(fstest.MapFS{
		"translations/en/general/location.json": &fstest.MapFile{Data: []byte(`{"page":{"heading_active":"Open","heading_inactive":"Closed"},"buttons":{"add_location":"Add place"},"columns":{"name":"Name","address":"Address"},"empty":{"active_title":"No places"},"form":{"name":"Place name","active":"Enabled"}}`)},
	})
	var labels Labels
	if err := provider.LoadFile("en", "sample", "location.json", &labels); err != nil {
		t.Fatal(err)
	}
	for field, tc := range map[string]struct{ got, want string }{
		"page.heading_active":   {labels.Page.HeadingActive, "Open"},
		"page.heading_inactive": {labels.Page.HeadingInactive, "Closed"},
		"buttons.add_location": {labels.Buttons.AddLocation, "Add place"},
		"columns.name":         {labels.Columns.Name, "Name"},
		"columns.address":      {labels.Columns.Address, "Address"},
		"empty.active_title":   {labels.Empty.ActiveTitle, "No places"},
		"form.name":            {labels.Form.Name, "Place name"},
		"form.active":          {labels.Form.Active, "Enabled"},
	} {
		if tc.got != tc.want {
			t.Errorf("provided %s label = %q, want %q", field, tc.got, tc.want)
		}
	}
}

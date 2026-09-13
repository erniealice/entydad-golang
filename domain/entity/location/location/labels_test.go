package location

import (
	lyngua "github.com/erniealice/lyngua"
	v1 "github.com/erniealice/lyngua/golang/v1"
	"testing"
)

func TestLocationLabelDefaults(t *testing.T) {
	p := v1.NewTranslationProviderFromFS(lyngua.TranslationsFS)
	for _, vertical := range []string{"general", "education", "leasing"} {
		t.Run(vertical, func(t *testing.T) {
			var labels Labels
			if err := p.LoadFile("en", vertical, "location.json", &labels); err != nil {
				t.Fatal(err)
			}
			if labels.Page.HeadingActive == "" || labels.Page.HeadingInactive == "" || labels.Buttons.AddLocation == "" || labels.Columns.Name == "" || labels.Columns.Address == "" || labels.Empty.ActiveTitle == "" || labels.Form.Name == "" {
				t.Fatal("incomplete location defaults")
			}
			if labels.Form.Active != "Active" {
				t.Fatalf("legacy active label changed: %q", labels.Form.Active)
			}
		})
	}
}

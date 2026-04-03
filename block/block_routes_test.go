package block

import (
	"testing"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/lyngua"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

func TestLoadBlockRoutes_ServiceBusinessTypeLoadsSubscriptionOverrides(t *testing.T) {
	provider := lynguaV1.NewTranslationProviderFromFS(lyngua.TranslationsFS)
	routes := loadBlockRoutes(provider, "service")

	if got, want := routes.Subscription.ListURL, "/app/memberships/list/{status}"; got != want {
		t.Fatalf("subscription list_url mismatch: got=%q want=%q", got, want)
	}
	if got, want := routes.Subscription.DetailURL, "/app/memberships/{id}"; got != want {
		t.Fatalf("subscription detail_url mismatch: got=%q want=%q", got, want)
	}
	if got, want := routes.Subscription.AddURL, "/action/memberships/add"; got != want {
		t.Fatalf("subscription add_url mismatch: got=%q want=%q", got, want)
	}
	if got, want := routes.Subscription.EditURL, "/action/memberships/edit/{id}"; got != want {
		t.Fatalf("subscription edit_url mismatch: got=%q want=%q", got, want)
	}
	if got, want := routes.Subscription.DeleteURL, "/action/memberships/delete"; got != want {
		t.Fatalf("subscription delete_url mismatch: got=%q want=%q", got, want)
	}
}

func TestLoadBlockRoutes_GeneralBusinessTypeFallsBackToDefaultSubscriptionRoutes(t *testing.T) {
	provider := lynguaV1.NewTranslationProviderFromFS(lyngua.TranslationsFS)
	routes := loadBlockRoutes(provider, "general")
	defaultRoutes := centymo.DefaultSubscriptionRoutes()

	if routes.Subscription != defaultRoutes {
		t.Fatalf("subscription routes mismatch for general business type: got=%+v want=%+v", routes.Subscription, defaultRoutes)
	}
}

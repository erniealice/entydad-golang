package block

import (
	"reflect"
	"testing"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/lyngua"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

func TestClientSubscriptionRoutesUseVerticalActions(t *testing.T) {
	provider := lynguaV1.NewTranslationProviderFromFS(lyngua.TranslationsFS)
	for _, tc := range []struct{ tier, add, detail string }{
		{"general", "/action/subscription/add", "/subscriptions/detail/{id}"},
		{"leasing", "/action/subscription/add", "/agreements/detail/{id}"},
		{"service", "/action/membership/add", "/memberships/detail/{id}"},
	} {
		r := loadClientSubscriptionRoutes(provider, tc.tier)
		if r.AddURL != tc.add || r.DetailURL != tc.detail {
			t.Errorf("%s routes: add=%q detail=%q", tc.tier, r.AddURL, r.DetailURL)
		}
	}
}

func TestLoadBlockRoutes_ServiceBusinessTypeLoadsSubscriptionOverrides(t *testing.T) {
	provider := lynguaV1.NewTranslationProviderFromFS(lyngua.TranslationsFS)
	routes := loadBlockRoutes(provider, "service")

	if got, want := routes.Subscription.ListURL, "/memberships/list/{status}"; got != want {
		t.Fatalf("subscription list_url mismatch: got=%q want=%q", got, want)
	}
	if got, want := routes.Subscription.DetailURL, "/memberships/detail/{id}"; got != want {
		t.Fatalf("subscription detail_url mismatch: got=%q want=%q", got, want)
	}
	if got, want := routes.Subscription.AddURL, "/action/membership/add"; got != want {
		t.Fatalf("subscription add_url mismatch: got=%q want=%q", got, want)
	}
	if got, want := routes.Subscription.EditURL, "/action/membership/edit/{id}"; got != want {
		t.Fatalf("subscription edit_url mismatch: got=%q want=%q", got, want)
	}
	if got, want := routes.Subscription.DeleteURL, "/action/membership/delete"; got != want {
		t.Fatalf("subscription delete_url mismatch: got=%q want=%q", got, want)
	}
}

func TestLoadBlockRoutes_GeneralBusinessTypeFallsBackToDefaultSubscriptionRoutes(t *testing.T) {
	provider := lynguaV1.NewTranslationProviderFromFS(lyngua.TranslationsFS)
	routes := loadBlockRoutes(provider, "general")
	defaultRoutes := centymo.DefaultSubscriptionRoutes()

	if !reflect.DeepEqual(routes.Subscription, defaultRoutes) {
		t.Fatalf("subscription routes mismatch for general business type: got=%+v want=%+v", routes.Subscription, defaultRoutes)
	}
}

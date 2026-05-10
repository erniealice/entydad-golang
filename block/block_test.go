package block

import (
	"testing"

	lyngua "github.com/erniealice/lyngua"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

func TestLoadBlockRoutes_ServiceOverridesCrossPackageSubscriptionRoutes(t *testing.T) {
	t.Parallel()

	translations := lynguaV1.NewTranslationProviderFromFS(lyngua.TranslationsFS)
	routes := loadBlockRoutes(translations, "service")

	if got, want := routes.Subscription.ListURL, "/app/memberships/list/{status}"; got != want {
		t.Fatalf("Subscription.ListURL = %q, want %q", got, want)
	}
	if got, want := routes.Subscription.DetailURL, "/app/memberships/detail/{id}"; got != want {
		t.Fatalf("Subscription.DetailURL = %q, want %q", got, want)
	}
	if got, want := routes.Subscription.AddURL, "/action/membership/add"; got != want {
		t.Fatalf("Subscription.AddURL = %q, want %q", got, want)
	}
	if got, want := routes.Subscription.EditURL, "/action/membership/edit/{id}"; got != want {
		t.Fatalf("Subscription.EditURL = %q, want %q", got, want)
	}
	if got, want := routes.Subscription.DeleteURL, "/action/membership/delete"; got != want {
		t.Fatalf("Subscription.DeleteURL = %q, want %q", got, want)
	}
	if got, want := routes.Client.ListURL, "/app/customers/list/{status}"; got != want {
		t.Fatalf("Client.ListURL = %q, want %q", got, want)
	}
}

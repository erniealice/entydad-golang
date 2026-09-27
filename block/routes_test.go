package block

import (
	"reflect"
	"testing"
	"testing/fstest"

	centymo "github.com/erniealice/centymo-golang"
	lyngua "github.com/erniealice/lyngua"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

func testRouteProvider() *lynguaV1.TranslationProvider {
	return lynguaV1.NewTranslationProviderFromFS(fstest.MapFS{
		"translations/en/general/route.json": &fstest.MapFile{Data: []byte(`{"subscription":{"add_url":"/generic/action/add","detail_url":"/generic/detail/{id}"}}`)},
		"translations/en/sample/route.json":  &fstest.MapFile{Data: []byte(`{"subscription":{"list_url":"/sample/list/{status}","detail_url":"/sample/detail/{id}"},"client":{"list_url":"/sample/clients"}}`)},
	})
}

func TestClientSubscriptionRoutesApplyProvidedOverride(t *testing.T) {
	routes := loadClientSubscriptionRoutes(testRouteProvider(), "sample")
	if got, want := routes.AddURL, "/generic/action/add"; got != want {
		t.Errorf("inherited add route = %q, want %q", got, want)
	}
	if got, want := routes.DetailURL, "/sample/detail/{id}"; got != want {
		t.Errorf("overridden detail route = %q, want %q", got, want)
	}
}

func TestLoadBlockRoutesAppliesProvidedOverrides(t *testing.T) {
	routes := loadBlockRoutes(testRouteProvider(), "sample")
	if got, want := routes.Subscription.ListURL, "/sample/list/{status}"; got != want {
		t.Errorf("subscription list route = %q, want %q", got, want)
	}
	if got, want := routes.Client.ListURL, "/sample/clients"; got != want {
		t.Errorf("client list route = %q, want %q", got, want)
	}
}

func TestLoadBlockRoutesFallsBackToDefaultSubscriptionRoutes(t *testing.T) {
	provider := lynguaV1.NewTranslationProviderFromFS(fstest.MapFS{})
	routes := loadBlockRoutes(provider, "sample")
	defaultRoutes := centymo.DefaultSubscriptionRoutes()

	if !reflect.DeepEqual(routes.Subscription, defaultRoutes) {
		t.Fatalf("subscription routes mismatch without an override: got=%+v want=%+v", routes.Subscription, defaultRoutes)
	}
}

func TestLoadBlockRoutesGeneralTierFallsBackToDefaults(t *testing.T) {
	provider := lynguaV1.NewTranslationProviderFromFS(lyngua.TranslationsFS)
	routes := loadBlockRoutes(provider, "general")
	if want := centymo.DefaultSubscriptionRoutes(); !reflect.DeepEqual(routes.Subscription, want) {
		t.Fatalf("general tier changed default subscription routes: got=%+v want=%+v", routes.Subscription, want)
	}
}

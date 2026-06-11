package block

import (
	"context"
	"strings"
	"testing"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
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

// ---------------------------------------------------------------------------
// MustValidate — FAIL-CLOSED wiring guard (architecture-roast burn #1).
//
// RequireFor returns an error; MustValidate adds the posture: in dev/test
// (testing.Testing() is true here) a missing REQUIRED closure PANICS — loud,
// stack-traced, uncatchable-by-accident — so a nil-closure wiring gap can never
// be silently dropped into an empty-state render. OPTIONAL nils never trip it.
// ---------------------------------------------------------------------------

// wireClientRequired sets every closure RequireFor checks for the Client
// module: the five Client CRUD/page closures (Category + cross-domain deps are
// optional, nil-safe).
func wireClientRequired(uc *UseCases) {
	c := &uc.Client
	c.GetListPageData = func(context.Context, *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error) {
		return nil, nil
	}
	c.Create = func(context.Context, *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error) { return nil, nil }
	c.Read = func(context.Context, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error) { return nil, nil }
	c.Update = func(context.Context, *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error) { return nil, nil }
	c.Delete = func(context.Context, *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error) { return nil, nil }
}

// TestMustValidate_NilRequiredClosure_Panics is the core burn-#1 proof: with
// the Client module enabled but one REQUIRED closure (GetListPageData) left nil,
// MustValidate must PANIC under test — not return an empty render, not silently
// degrade. This is the loud failure the bare-return path lacked.
func TestMustValidate_NilRequiredClosure_Panics(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	wireClientRequired(uc)
	uc.Client.GetListPageData = nil // drop exactly one REQUIRED closure

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("MustValidate(Client enabled, GetListPageData nil) should PANIC in dev/test, but did not")
		}
		msg, _ := r.(string)
		if !strings.Contains(msg, "GetListPageData") {
			t.Fatalf("panic message should name the missing field; got %q", msg)
		}
	}()

	// Should not reach the next line — MustValidate panics first.
	_ = uc.MustValidate(&blockConfig{client: true})
	t.Fatal("MustValidate returned instead of panicking on a nil REQUIRED closure")
}

// TestMustValidate_EmptyUseCases_EnableAll_Panics: a fully empty UseCases with
// every module enabled (the "permanently nil dashboard" trap) must panic loudly
// in dev/test rather than register a wall of empty views.
func TestMustValidate_EmptyUseCases_EnableAll_Panics(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	defer func() {
		if recover() == nil {
			t.Fatal("MustValidate(empty UseCases, enableAll) should PANIC in dev/test")
		}
	}()
	_ = uc.MustValidate(&blockConfig{enableAll: true})
	t.Fatal("MustValidate returned instead of panicking on an empty enableAll wiring")
}

// TestMustValidate_NilOptionalClosure_OK proves the required-vs-optional
// discrimination survives the fail-closed wrapper: the OPTIONAL clientTag module
// (not in RequireFor) enabled with nil closures must pass MustValidate with NO
// panic and NO error — optional features stay legitimately nil.
func TestMustValidate_NilOptionalClosure_OK(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	// An optional module enabled, its closures left nil. clientTag has no
	// RequireFor assertions, so nil closures must not trip the guard.
	cfg := &blockConfig{clientTag: true}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustValidate(optional nil closures) must NOT panic; panicked with %v", r)
		}
	}()
	if err := uc.MustValidate(cfg); err != nil {
		t.Fatalf("MustValidate(optional nil closures) should be nil, got %v", err)
	}
}

// TestMustValidate_FullyWired_OK: a completely wired REQUIRED set passes with no
// panic and no error (happy path — guard is silent when wiring is complete).
func TestMustValidate_FullyWired_OK(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	wireClientRequired(uc)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustValidate(fully wired Client) must NOT panic; panicked with %v", r)
		}
	}()
	if err := uc.MustValidate(&blockConfig{client: true}); err != nil {
		t.Fatalf("MustValidate(fully wired Client) should be nil, got %v", err)
	}
}

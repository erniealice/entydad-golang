package detail

import (
	"net/url"
	"testing"
)

func TestSubscriptionDrawerURLPreservesContext(t *testing.T) {
	if got := buildSubscriptionAddURL("", "client-1", "Tenant", "PHP"); got != "" {
		t.Fatalf("missing action route must not become a current-page query: %q", got)
	}
	u, err := url.Parse(buildSubscriptionAddURL("/action/subscription/add", "client-1", "A & B", "PHP"))
	if err != nil {
		t.Fatal(err)
	}
	if u.Path != "/action/subscription/add" || u.Query().Get("client_id") != "client-1" || u.Query().Get("client_name") != "A & B" || u.Query().Get("billing_currency") != "PHP" {
		t.Fatalf("unexpected drawer URL: %s", u)
	}
}

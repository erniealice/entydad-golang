package block

import "testing"

// Wave 0 of plan 20260927-tenant-boundary-hardening: the generic SetActive /
// SetStatus closures do a raw ops.Update with no use-case gate, so they must
// refuse the platform-global collections (user, workspace) that carry no
// workspace_id, and keep working for tenant-scoped collections.
func TestDenyPlatformGlobalToggle(t *testing.T) {
	for _, c := range []string{"user", "workspace"} {
		if err := denyPlatformGlobalToggle(c); err == nil {
			t.Fatalf("collection %q must be denied", c)
		}
	}
	for _, c := range []string{"role", "permission", "location", "workspace_user", "client"} {
		if err := denyPlatformGlobalToggle(c); err != nil {
			t.Fatalf("collection %q must stay allowed: %v", c, err)
		}
	}
}

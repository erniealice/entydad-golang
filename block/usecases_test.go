package block

import (
	"strings"
	"testing"
)

// TestRequireFor_Empty asserts that a completely empty UseCases fails
// RequireFor when enableAll is true — all required groups are missing.
func TestRequireFor_Empty(t *testing.T) {
	cfg := &blockConfig{enableAll: true}
	uc := &UseCases{}
	err := uc.RequireFor(cfg)
	if err == nil {
		t.Fatal("expected error for empty UseCases with enableAll, got nil")
	}
}

// TestRequireFor_PartialClientOnly asserts that an empty UseCases still
// fails RequireFor when only cfg.client is true, because none of the
// Client CRUD closures are set.
func TestRequireFor_PartialClientOnly(t *testing.T) {
	cfg := &blockConfig{client: true}
	uc := &UseCases{}
	err := uc.RequireFor(cfg)
	if err == nil {
		t.Fatal("expected error when Client CRUD fields are nil, got nil")
	}
	if !strings.Contains(err.Error(), "UseCases.Client") {
		t.Errorf("expected error to mention UseCases.Client, got: %v", err)
	}
}

// TestRequireFor_NoModulesEnabled asserts that an empty UseCases passes
// RequireFor when no module flags are set (nothing to require).
func TestRequireFor_NoModulesEnabled(t *testing.T) {
	cfg := &blockConfig{} // enableAll=false, all flags false
	uc := &UseCases{}
	if err := uc.RequireFor(cfg); err != nil {
		t.Fatalf("expected no error when no modules enabled, got: %v", err)
	}
}

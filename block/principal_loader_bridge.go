// Package block — principal_loader_bridge.go (relocated from app composition, Wave B D2a)
//
// Thin wrapper around consumer principal loader utilities. The proto-to-local
// type mapping lives in consumer/principal_loader_bridge.go; this file bridges
// the consumer PrincipalData type to espyna's adapthttp.Principal type and
// wires the auth-chain deps' use cases into a DBPrincipalLoader.

package block

import (
	"context"

	consumer "github.com/erniealice/espyna-golang/consumer"

	adapthttp "github.com/erniealice/espyna-golang/consumer/http"
)

// newPrincipalLoader creates a DBPrincipalLoader backed by the espyna
// ResolvePrincipals and ResolveBinding use cases. This bridges the use case
// layer (using proto types) to the adapthttp layer (using local Principal
// types). Proto mapping is delegated to consumer.BuildPrincipalResolveFn
// and consumer.BuildBindingResolveFn.
func (d *authChainDeps) newPrincipalLoader() *adapthttp.DBPrincipalLoader {
	if d.uc == nil || d.uc.Service == nil || d.uc.Service.Auth == nil {
		return nil
	}
	resolvePrincipalsUC := d.uc.Service.Auth.ResolvePrincipals
	resolveBindingUC := d.uc.Service.Auth.ResolveBinding
	if resolvePrincipalsUC == nil {
		return nil
	}

	// Build the consumer-level resolve function, then bridge to adapthttp types.
	consumerResolveFn := consumer.BuildPrincipalResolveFn(resolvePrincipalsUC)
	resolveFn := func(ctx context.Context, userID string) ([]adapthttp.Principal, error) {
		data, err := consumerResolveFn(ctx, userID)
		if err != nil {
			return nil, err
		}
		return consumerDataToAdaptPrincipals(data), nil
	}

	// Build the consumer-level binding function, then bridge to adapthttp types.
	var bindingFn adapthttp.BindingResolveFunc
	consumerBindingFn := consumer.BuildBindingResolveFn(resolveBindingUC)
	if consumerBindingFn != nil {
		bindingFn = func(
			ctx context.Context,
			userID, workspaceID string,
			sessionPrincipalKind adapthttp.PrincipalType,
			sessionPrincipalID string,
		) (*adapthttp.Principal, error) {
			data, err := consumerBindingFn(ctx, userID, workspaceID, int32(sessionPrincipalKind), sessionPrincipalID)
			if err != nil {
				// Map sentinel error messages to adapthttp sentinel errors.
				errMsg := err.Error()
				if errMsg == "resolve_binding: no active binding in workspace" {
					return nil, adapthttp.ErrNoBinding
				}
				if errMsg == "resolve_binding: ambiguous binding (multiple bindings, no session principal match)" {
					return nil, adapthttp.ErrAmbiguousBinding
				}
				return nil, err
			}
			if data == nil {
				return nil, adapthttp.ErrNoBinding
			}
			p := consumerDataToAdaptPrincipal(*data)
			return &p, nil
		}
	}

	return adapthttp.NewDBPrincipalLoader(resolveFn, bindingFn)
}

// consumerDataToAdaptPrincipals converts consumer PrincipalData to adapthttp Principals.
func consumerDataToAdaptPrincipals(data []consumer.PrincipalData) []adapthttp.Principal {
	out := make([]adapthttp.Principal, 0, len(data))
	for _, d := range data {
		out = append(out, consumerDataToAdaptPrincipal(d))
	}
	return out
}

// consumerDataToAdaptPrincipal converts a single consumer PrincipalData to adapthttp.Principal.
func consumerDataToAdaptPrincipal(d consumer.PrincipalData) adapthttp.Principal {
	targets := make([]adapthttp.ActingAsTarget, 0, len(d.ActingAsTargets))
	for _, t := range d.ActingAsTargets {
		targets = append(targets, adapthttp.ActingAsTarget{
			ID:          t.ID,
			WorkspaceID: t.WorkspaceID,
			DisplayName: t.DisplayName,
		})
	}
	return adapthttp.Principal{
		Type:            adapthttp.PrincipalType(d.Type),
		PrincipalID:     d.PrincipalID,
		WorkspaceID:     d.WorkspaceID,
		DisplayName:     d.DisplayName,
		ActingAsTargets: targets,
	}
}

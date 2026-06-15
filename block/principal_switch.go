// Package block — principal_switch.go (relocated from app composition, Wave B D2a)
//
// Thin wrapper around consumer.ExecutePrincipalSwitch, bridging the auth-chain
// deps' use-case lookup to the framework-level proto request builder. The
// portable logic (types, proto mapping, request execution) lives in
// packages/espyna-golang/consumer/principal_switch.go.

package block

import (
	"context"
	"errors"

	consumer "github.com/erniealice/espyna-golang/consumer"

	adapthttp "github.com/erniealice/espyna-golang/consumer/http"
)

// authPrincipalSwitchInput is the auth-chain switch operation input. It uses the
// espyna adapthttp.Principal type. executeAuthPrincipalSwitch bridges this to
// consumer.PrincipalSwitchInput.
type authPrincipalSwitchInput struct {
	UserID             string
	Token              string
	TargetPrincipal    adapthttp.Principal
	ActingAsClientID   string
	ActingAsSupplierID string
	UseCase            string
	RequestURL         string
	Referer            string
	SecFetchSite       string
	UserAgent          string
	URLDriven          bool
	RequireAudit       bool
}

// authPrincipalSwitchResult is the auth-chain switch result.
type authPrincipalSwitchResult = consumer.PrincipalSwitchResult

// executeAuthPrincipalSwitch bridges the auth-chain deps' use-case lookup to the
// framework-level consumer.ExecutePrincipalSwitch function.
func (d *authChainDeps) executePrincipalSwitch(
	ctx context.Context,
	in authPrincipalSwitchInput,
) (*authPrincipalSwitchResult, error) {
	if d.uc == nil || d.uc.Service == nil ||
		d.uc.Service.Auth == nil || d.uc.Service.Auth.SwitchPrincipal == nil {
		return nil, errors.New("principal switch: SwitchPrincipal use case not wired")
	}

	return consumer.ExecutePrincipalSwitch(ctx, d.uc.Service.Auth.SwitchPrincipal, consumer.PrincipalSwitchInput{
		UserID:             in.UserID,
		Token:              in.Token,
		TargetPrincipal:    adaptPrincipalToConsumer(in.TargetPrincipal),
		ActingAsClientID:   in.ActingAsClientID,
		ActingAsSupplierID: in.ActingAsSupplierID,
		UseCase:            in.UseCase,
		RequestURL:         in.RequestURL,
		Referer:            in.Referer,
		SecFetchSite:       in.SecFetchSite,
		UserAgent:          in.UserAgent,
		URLDriven:          in.URLDriven,
		RequireAudit:       in.RequireAudit,
	})
}

// adaptPrincipalToConsumer maps adapthttp.Principal to consumer.PrincipalData.
func adaptPrincipalToConsumer(p adapthttp.Principal) consumer.PrincipalData {
	targets := make([]consumer.ActingAsTargetData, 0, len(p.ActingAsTargets))
	for _, t := range p.ActingAsTargets {
		targets = append(targets, consumer.ActingAsTargetData{
			ID:          t.ID,
			WorkspaceID: t.WorkspaceID,
			DisplayName: t.DisplayName,
		})
	}
	return consumer.PrincipalData{
		Type:            int32(p.Type),
		PrincipalID:     p.PrincipalID,
		WorkspaceID:     p.WorkspaceID,
		DisplayName:     p.DisplayName,
		ActingAsTargets: targets,
	}
}

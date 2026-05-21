package profile

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	espynaconsumer "github.com/erniealice/espyna-golang/consumer"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ReadUserFunc matches *ReadUserUseCase.Execute in
// packages/espyna-golang/internal/application/usecases/domain/entity/user/.
type ReadUserFunc func(ctx context.Context, req *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error)

// PageData holds data for the client profile page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels   entydad.PortalLabels
	ReadUser ReadUserFunc // Loads the logged-in user's basics (Pre-B real backing).
}

// NewView creates the client profile view (GET /portal/client/profile).
//
// First Pre-B wave per docs/plan/20260521-workspace-keyed-routing/phases.md
// Pre-B and docs/plan/20260516-self-domain/ §P6 — the page now loads the
// logged-in User via ReadUser and renders First/Last name, email, and mobile.
// Read-only for v1; edit/upsert action handler is a follow-up.
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	var readUser ReadUserFunc
	if deps != nil {
		labels = deps.Labels
		readUser = deps.ReadUser
	}
	if labels.Client.Profile.PageTitle == "" {
		labels.Client.Profile.PageTitle = "Profile"
	}
	if labels.Client.Name == "" {
		labels.Client.Name = "My Account"
	}
	if labels.Page.Profile.Title == "" {
		labels.Page.Profile.Title = "Your profile"
	}
	if labels.Page.Profile.ComingSoon == "" {
		labels.Page.Profile.ComingSoon = "Profile information unavailable."
	}
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// Load the logged-in user. If session has no userID or readUser is
		// unwired, fall back to empty ProfileUser (template guards on the
		// FirstName field). The page still renders so the navigation stays
		// usable; the ComingSoon hint signals when data isn't loaded.
		var profileUser portal.ProfileUser
		if readUser != nil {
			if userID, err := espynaconsumer.RequireUserIDFromContext(ctx); err == nil {
				resp, err := readUser(ctx, &userpb.ReadUserRequest{Data: &userpb.User{Id: userID}})
				if err == nil && resp != nil && len(resp.Data) > 0 {
					u := resp.Data[0]
					profileUser = portal.ProfileUser{
						ID:           u.Id,
						FirstName:    u.FirstName,
						LastName:     u.LastName,
						EmailAddress: u.EmailAddress,
						MobileNumber: u.MobileNumber,
					}
				}
			}
		}

		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.Client.Profile.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-client-profile-content",
					ActiveNav:       "profile",
				},
				PrincipalKind: "client",
				PortalName:    labels.Client.Name,
				ProfileURL:    "/portal/client/profile",
				User:          profileUser,
			},
			Labels: labels,
		}
		return view.OK("portal-client-profile", pageData)
	})
}

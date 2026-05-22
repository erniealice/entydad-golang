// Package recent_activity is the /me/recent-activity view (Phase P9b).
//
// Lists the user's recent workspace switches — forensic surface required by
// the security red-team (A-2). Each switch is one row of audit_entry where
// `actor_user_id = current_user_id AND use_case LIKE 'switch_%'`, ordered by
// most recent. Per phases.md 9b lock (a), this view is READ-ONLY; there are
// no actions on /me/* in v1.
//
// The view consumes a ListRecentSwitchesFunc closure injected by the
// composition layer. The closure encapsulates the audit-entry query so
// entydad does not take a direct dependency on the postgres adapter.
package recent_activity

import (
	"context"

	me "github.com/erniealice/entydad-golang/views/me"
	espynaconsumer "github.com/erniealice/espyna-golang/consumer"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// SwitchEntry is one row in the recent-activity table.
type SwitchEntry struct {
	OccurredAt string // ISO 8601 timestamp
	UseCase    string // e.g. switch_url_rotate, switch_explicit
	RequestURL string
	Referer    string
}

// ListRecentSwitchesFunc looks up the most-recent N workspace switch rows
// for the given user. Limit is the cap on rows returned.
type ListRecentSwitchesFunc func(ctx context.Context, userID string, limit int) ([]SwitchEntry, error)

// ModuleDeps bundles per-request configuration for the view.
type ModuleDeps struct {
	Messages           map[string]string
	ListRecentSwitches ListRecentSwitchesFunc
}

// PageData carries rendering context for the /me/recent-activity page.
type PageData struct {
	me.PageData
	Subtitle     string
	EmptyMessage string
	Switches     []SwitchEntry
}

// NewView constructs the /me/recent-activity view.
func NewView(deps *ModuleDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		var switches []SwitchEntry
		if deps != nil && deps.ListRecentSwitches != nil {
			if userID, err := espynaconsumer.RequireUserIDFromContext(ctx); err == nil {
				if rows, err := deps.ListRecentSwitches(ctx, userID, 50); err == nil {
					switches = rows
				}
			}
		}

		pd := &PageData{
			PageData: me.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           "Recent Activity",
					CurrentPath:     viewCtx.CurrentPath,
					ActiveNav:       "recent-activity",
					ContentTemplate: "me-recent-activity-content",
					HeaderTitle:     "Recent Activity",
					HeaderIcon:      "icon-clock",
					Messages:        deps.Messages,
				},
			},
			Subtitle:     "Your recent workspace switches across all sessions.",
			EmptyMessage: "No recent workspace switches.",
			Switches:     switches,
		}
		return view.OK("me-page", pd)
	})
}

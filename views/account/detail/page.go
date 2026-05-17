package detail

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds dependencies needed to build the account detail view.
type ModuleDeps struct {
	Messages          map[string]string
	ChangePasswordURL string
}

// PageData carries the rendering context for the account page.
type PageData struct {
	types.PageData
	TabItems          []pyeza.TabItem
	ActiveTab         string
	ChangePasswordURL string
}

// NewView creates the account detail view (full page — tabs: email | password | sessions).
func NewView(deps *ModuleDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("user", "update") {
			return view.Forbidden("user:update")
		}

		activeTab := ""
		if viewCtx.Request != nil {
			activeTab = viewCtx.Request.URL.Query().Get("tab")
		}
		tabs := buildTabs(deps.Messages)
		if activeTab == "" || !validTab(tabs, activeTab) {
			activeTab = "email"
		}

		titleKey := "memberPages.section.account.title"
		iconKey := "memberPages.section.account.icon"

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:    viewCtx.CacheVersion,
				Title:           lookup(deps.Messages, titleKey, "Account"),
				CurrentPath:     viewCtx.CurrentPath,
				ActiveNav:       "home",
				ContentTemplate: "account-page-content",
				HeaderTitle:     lookup(deps.Messages, titleKey, "Account"),
				HeaderIcon:      lookup(deps.Messages, iconKey, "icon-user"),
				Messages:        deps.Messages,
			},
			TabItems:          tabs,
			ActiveTab:         activeTab,
			ChangePasswordURL: deps.ChangePasswordURL,
		}
		return view.OK("account-page", pageData)
	})
}

func buildTabs(messages map[string]string) []pyeza.TabItem {
	return []pyeza.TabItem{
		{Key: "email", Label: lookup(messages, "memberPages.account.tab.email", "Sign-in email"), Href: "/app/account?tab=email"},
		{Key: "password", Label: lookup(messages, "memberPages.account.tab.password", "Password"), Href: "/app/account?tab=password"},
		{Key: "sessions", Label: lookup(messages, "memberPages.account.tab.sessions", "Sessions"), Href: "/app/account?tab=sessions"},
	}
}

func validTab(tabs []pyeza.TabItem, key string) bool {
	for _, t := range tabs {
		if t.Key == key {
			return true
		}
	}
	return false
}

func lookup(messages map[string]string, key, fallback string) string {
	if messages != nil {
		if v, ok := messages[key]; ok && v != "" {
			return v
		}
	}
	return fallback
}

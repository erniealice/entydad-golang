package detail

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds dependencies needed to build the preference detail view.
type ModuleDeps struct {
	Messages map[string]string
}

// PageData carries the rendering context for the preference page.
type PageData struct {
	types.PageData
	TabItems  []pyeza.TabItem
	ActiveTab string
}

// NewView creates the preference detail view (full page — tabs: appearance | notifications | language-region).
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
			activeTab = "appearance"
		}

		titleKey := "memberPages.section.preferences.title"
		iconKey := "memberPages.section.preferences.icon"

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:    viewCtx.CacheVersion,
				Title:           lookup(deps.Messages, titleKey, "Preferences"),
				CurrentPath:     viewCtx.CurrentPath,
				ActiveNav:       "home",
				ContentTemplate: "preferences-page-content",
				HeaderTitle:     lookup(deps.Messages, titleKey, "Preferences"),
				HeaderIcon:      lookup(deps.Messages, iconKey, "icon-settings"),
				Messages:        deps.Messages,
			},
			TabItems:  tabs,
			ActiveTab: activeTab,
		}
		return view.OK("preferences-page", pageData)
	})
}

func buildTabs(messages map[string]string) []pyeza.TabItem {
	return []pyeza.TabItem{
		{Key: "appearance", Label: lookup(messages, "memberPages.preferences.tab.appearance", "Appearance"), Href: "/app/preferences?tab=appearance"},
		{Key: "notifications", Label: lookup(messages, "memberPages.preferences.tab.notifications", "Notifications"), Href: "/app/preferences?tab=notifications"},
		{Key: "language-region", Label: lookup(messages, "memberPages.preferences.tab.languageRegion", "Language & region"), Href: "/app/preferences?tab=language-region"},
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

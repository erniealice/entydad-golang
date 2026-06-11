package detail

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds dependencies needed to build the billing detail view.
type ModuleDeps struct {
	Messages map[string]string
	// PageURL is the base URL for the billing page used to build tab hrefs.
	// Defaults to "/app/billing" when empty for backward compatibility.
	PageURL string
}

// PageData carries the rendering context for the billing page.
type PageData struct {
	types.PageData
	TabItems  []pyeza.TabItem
	ActiveTab string
}

// NewView creates the billing detail view (full page — tabs: subscription | payment-method | invoices).
func NewView(deps *ModuleDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("workspace", "read") {
			return view.Forbidden("workspace:read")
		}

		activeTab := ""
		if viewCtx.Request != nil {
			activeTab = viewCtx.Request.URL.Query().Get("tab")
		}
		tabs := buildTabs(deps.Messages, deps.PageURL)
		if activeTab == "" || !validTab(tabs, activeTab) {
			activeTab = "subscription"
		}

		titleKey := "memberPages.section.billing.title"
		iconKey := "memberPages.section.billing.icon"

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:    viewCtx.CacheVersion,
				Title:           lookup(deps.Messages, titleKey, "Billing"),
				CurrentPath:     viewCtx.CurrentPath,
				ActiveNav:       "home",
				ContentTemplate: "billing-page-content",
				HeaderTitle:     lookup(deps.Messages, titleKey, "Billing"),
				HeaderIcon:      lookup(deps.Messages, iconKey, "icon-credit-card"),
				Messages:        deps.Messages,
			},
			TabItems:  tabs,
			ActiveTab: activeTab,
		}
		return view.OK("billing-page", pageData)
	})
}

func buildTabs(messages map[string]string, pageURL string) []pyeza.TabItem {
	if pageURL == "" {
		pageURL = "/app/billing"
	}
	return []pyeza.TabItem{
		{Key: "subscription", Label: lookup(messages, "memberPages.billing.tab.subscription", "Subscription"), Href: pageURL + "?tab=subscription"},
		{Key: "payment-method", Label: lookup(messages, "memberPages.billing.tab.paymentMethod", "Payment method"), Href: pageURL + "?tab=payment-method"},
		{Key: "invoices", Label: lookup(messages, "memberPages.billing.tab.invoices", "Invoices"), Href: pageURL + "?tab=invoices"},
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

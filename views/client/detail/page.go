package detail

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes                entydad.ClientRoutes
	ReadClient            func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	ListCategories        func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListClientCategories  func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
	ListRevenues                func(ctx context.Context, collection string) ([]map[string]any, error)
	ListSubscriptions           func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)
	GetSubscriptionListPageData func(ctx context.Context, req *subscriptionpb.GetSubscriptionListPageDataRequest) (*subscriptionpb.GetSubscriptionListPageDataResponse, error)
	SubscriptionAddURL    string
	SubscriptionDetailURL string
	SubscriptionEditURL   string
	SubscriptionDeleteURL string
	Labels                entydad.ClientLabels
	CommonLabels          pyeza.CommonLabels
	TableLabels           types.TableLabels

	// Attachment operations (embedded from hybra)
	attachment.AttachmentOps

	// Audit log operations (embedded from hybra)
	auditlog.AuditOps
}

// TagChip represents a tag displayed as a chip on the detail page.
type TagChip struct {
	Name string
}

// PurchaseStats holds aggregated purchase statistics for a client.
type PurchaseStats struct {
	LifetimeSpend string
	TotalOrders   int
	AvgOrderValue string
	LastPurchase  string
}

// OrderRow represents a single order in the purchase history table.
type OrderRow struct {
	ID        string
	Reference string
	Date      string
	Amount    string
	Status    string
	Variant   string
}

// SubscriptionRow represents a single subscription/engagement in the subscriptions tab.
type SubscriptionRow struct {
	ID        string
	Name      string
	Plan      string
	DateStart string
	DateEnd   string
}

// PageData holds the data for the client detail page.
type PageData struct {
	types.PageData
	ContentTemplate    string
	Client             *clientpb.Client
	Labels             entydad.ClientLabels
	ActiveTab          string
	TabItems           []pyeza.TabItem
	EditURL            string
	SubscriptionAddURL string
	ClientName         string
	RepresentativeName string
	ClientEmail        string
	ClientPhone        string
	ClientStatus       string
	StatusVariant      string
	// CRM fields
	Name          string
	StreetAddress string
	City          string
	Province      string
	PostalCode    string
	Notes         string
	FullAddress   string
	Tags          []TagChip
	// Has* booleans for conditional rendering in templates
	HasName    bool
	HasAddress bool
	HasNotes   bool
	HasTags    bool
	// Purchase history
	PurchaseStats PurchaseStats
	Orders        []OrderRow
	HasOrders     bool
	// Engagements tab
	Subscriptions    []SubscriptionRow
	EngagementsTable *types.TableConfig
}

// NewView creates the client detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		resp, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{
			Data: &clientpb.Client{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read client %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load client: %w", err))
		}

		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("client not found"))
		}
		client := data[0]
		u := client.GetUser()

		clientName := clientDisplayName(client)
		representativeName := ""
		clientEmail := ""
		clientPhone := ""
		if u != nil {
			representativeName = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			clientEmail = u.GetEmailAddress()
			clientPhone = u.GetMobileNumber()
		}

		clientStatus := "active"
		if !client.GetActive() {
			clientStatus = "inactive"
		}
		statusVariant := "success"
		if clientStatus == "inactive" {
			statusVariant = "warning"
		}

		tabItems := buildTabItems(id, deps)

		// CRM fields
		name := client.GetName()
		streetAddress := client.GetStreetAddress()
		city := client.GetCity()
		province := client.GetProvince()
		postalCode := client.GetPostalCode()
		notes := client.GetNotes()
		fullAddress := buildFullAddress(streetAddress, city, province, postalCode)

		hasName := name != ""
		hasAddress := streetAddress != "" || city != "" || province != "" || postalCode != ""
		hasNotes := notes != ""

		// Load tags for this client
		tags := loadClientTags(ctx, deps, id)
		hasTags := len(tags) > 0

		// Load purchase history
		stats, orders := loadPurchaseHistory(ctx, deps, id)
		hasOrders := len(orders) > 0

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          clientName,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "client",
				HeaderTitle:    clientName,
				HeaderSubtitle: clientEmail,
				HeaderIcon:     "icon-user",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:    "client-detail-content",
			Client:             client,
			Labels:             deps.Labels,
			ActiveTab:          activeTab,
			TabItems:           tabItems,
			EditURL:            route.ResolveURL(deps.Routes.EditURL, "id", id),
			SubscriptionAddURL: deps.SubscriptionAddURL + "?client_id=" + id + "&client_name=" + url.QueryEscape(clientName),
			ClientName:         clientName,
			RepresentativeName: representativeName,
			ClientEmail:        clientEmail,
			ClientPhone:        clientPhone,
			ClientStatus:       clientStatus,
			StatusVariant:      statusVariant,
			Name:               name,
			StreetAddress:      streetAddress,
			City:               city,
			Province:           province,
			PostalCode:         postalCode,
			Notes:              notes,
			FullAddress:        fullAddress,
			Tags:               tags,
			HasName:            hasName,
			HasAddress:         hasAddress,
			HasNotes:           hasNotes,
			HasTags:            hasTags,
			PurchaseStats:      stats,
			Orders:             orders,
			HasOrders:          hasOrders,
		}

		// Load tab-specific data for the active tab on full page load
		switch activeTab {
		case "engagements":
			subs := loadClientSubscriptions(ctx, deps, id)
			pageData.Subscriptions = subs
			pageData.EngagementsTable = buildEngagementsTable(subs, pageData.SubscriptionAddURL, id, clientName, deps)
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "client-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("client-detail", pageData)
	})
}

func buildTabItems(id string, deps *DetailViewDeps) []pyeza.TabItem {
	routes := deps.Routes
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: deps.Labels.Detail.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "representative", Label: deps.Labels.Detail.Tabs.Representative, Href: base + "?tab=representative", HxGet: action + "representative", Icon: "icon-user"},
		{Key: "engagements", Label: deps.Labels.Detail.Tabs.Engagements, Href: base + "?tab=engagements", HxGet: action + "engagements", Icon: "icon-file-text"},
		{Key: "history", Label: deps.Labels.Detail.Tabs.History, Href: base + "?tab=history", HxGet: action + "history", Icon: "icon-shopping-bag"},
	}
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{
			Data: &clientpb.Client{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read client %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load client: %w", err))
		}

		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("client not found"))
		}
		client := data[0]
		u := client.GetUser()

		clientName := clientDisplayName(client)
		representativeName := ""
		clientEmail := ""
		clientPhone := ""
		if u != nil {
			representativeName = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			clientEmail = u.GetEmailAddress()
			clientPhone = u.GetMobileNumber()
		}

		clientStatus := "active"
		if !client.GetActive() {
			clientStatus = "inactive"
		}
		statusVariant := "success"
		if clientStatus == "inactive" {
			statusVariant = "warning"
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Client:             client,
			Labels:             deps.Labels,
			ActiveTab:          tab,
			TabItems:           buildTabItems(id, deps),
			ClientName:         clientName,
			RepresentativeName: representativeName,
			ClientEmail:        clientEmail,
			ClientPhone:        clientPhone,
			ClientStatus:       clientStatus,
			StatusVariant:      statusVariant,
			EditURL:            route.ResolveURL(deps.Routes.EditURL, "id", id),
			SubscriptionAddURL: deps.SubscriptionAddURL + "?client_id=" + id + "&client_name=" + url.QueryEscape(clientName),
		}

		switch tab {
		case "info":
			pageData.Name = client.GetName()
			pageData.StreetAddress = client.GetStreetAddress()
			pageData.City = client.GetCity()
			pageData.Province = client.GetProvince()
			pageData.PostalCode = client.GetPostalCode()
			pageData.Notes = client.GetNotes()
			pageData.FullAddress = buildFullAddress(pageData.StreetAddress, pageData.City, pageData.Province, pageData.PostalCode)
			pageData.HasName = pageData.Name != ""
			pageData.HasAddress = pageData.StreetAddress != "" || pageData.City != "" || pageData.Province != "" || pageData.PostalCode != ""
			pageData.HasNotes = pageData.Notes != ""
			pageData.Tags = loadClientTags(ctx, deps, id)
			pageData.HasTags = len(pageData.Tags) > 0
		case "representative":
			// user fields already on client via GetUser()
		case "engagements":
			pageData.Subscriptions = loadClientSubscriptions(ctx, deps, id)
			pageData.EngagementsTable = buildEngagementsTable(pageData.Subscriptions, pageData.SubscriptionAddURL, id, clientName, deps)
		case "history":
			pageData.PurchaseStats, pageData.Orders = loadPurchaseHistory(ctx, deps, id)
			pageData.HasOrders = len(pageData.Orders) > 0
		}

		return view.OK("client-tab-"+tab, pageData)
	})
}

// loadClientSubscriptions fetches active subscriptions for a client
// using GetSubscriptionListPageData with a client_id filter so that
// PricePlan (and its embedded Plan) are populated via JOIN, enabling
// the Plan name column to display correctly.
func loadClientSubscriptions(ctx context.Context, deps *DetailViewDeps, clientID string) []SubscriptionRow {
	if deps.GetSubscriptionListPageData == nil {
		return nil
	}
	resp, err := deps.GetSubscriptionListPageData(ctx, &subscriptionpb.GetSubscriptionListPageDataRequest{
		Filters: &categorypb.FilterRequest{
			Filters: []*categorypb.TypedFilter{
				{
					Field: "client_id",
					FilterType: &categorypb.TypedFilter_StringFilter{
						StringFilter: &categorypb.StringFilter{
							Value:         clientID,
							Operator:      categorypb.StringOperator_STRING_EQUALS,
							CaseSensitive: true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to load subscriptions for client %s: %v", clientID, err)
		return nil
	}
	var rows []SubscriptionRow
	for _, s := range resp.GetSubscriptionList() {
		if !s.GetActive() {
			continue
		}
		planName := ""
		if pp := s.GetPricePlan(); pp != nil {
			if p := pp.GetPlan(); p != nil {
				planName = p.GetName()
			}
			if planName == "" {
				planName = pp.GetName()
			}
		}
		dateStart := s.GetDateStart()
		dateEnd := s.GetDateEnd()
		rows = append(rows, SubscriptionRow{
			ID:        s.GetId(),
			Name:      s.GetName(),
			Plan:      planName,
			DateStart: dateStart,
			DateEnd:   dateEnd,
		})
	}
	return rows
}

// clientDisplayName returns the client's display name.
// buildEngagementsTable builds a TableConfig for the engagements tab.
func buildEngagementsTable(rows []SubscriptionRow, addURL string, clientID string, clientName string, deps *DetailViewDeps) *types.TableConfig {
	if len(rows) == 0 {
		return nil
	}

	columns := []types.TableColumn{
		{Key: "name", Label: "Name", Sortable: true},
		{Key: "plan", Label: "Package", Sortable: true},
		{Key: "start_date", Label: "Start Date", Sortable: true, Width: "140px"},
		{Key: "end_date", Label: "End Date", Sortable: true, Width: "140px"},
	}

	// Build locked client query params for edit URLs
	clientParams := "?client_id=" + clientID + "&client_name=" + url.QueryEscape(clientName)

	var tableRows []types.TableRow
	for _, r := range rows {
		editURL := route.ResolveURL(deps.SubscriptionEditURL, "id", r.ID) + clientParams
		detailURL := route.ResolveURL(deps.SubscriptionDetailURL, "id", r.ID)

		tableRows = append(tableRows, types.TableRow{
			ID: r.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: r.Name},
				{Type: "text", Value: r.Plan},
				{Type: "text", Value: r.DateStart},
				{Type: "text", Value: r.DateEnd},
			},
			DataAttrs: map[string]string{
				"name": r.Name,
				"plan": r.Plan,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: "View", Action: "view", Href: detailURL},
				{Type: "edit", Label: "Edit", Action: "edit", URL: editURL, DrawerTitle: r.Name},
				{Type: "delete", Label: "Delete", Action: "delete", URL: deps.SubscriptionDeleteURL, ItemName: r.Name, ConfirmTitle: "Delete Engagement", ConfirmMessage: "Are you sure you want to delete " + r.Name + "?"},
			},
		})
	}

	types.ApplyColumnStyles(columns, tableRows)

	tc := &types.TableConfig{
		ID:                   "engagements-table",
		Columns:              columns,
		Rows:                 tableRows,
		Labels:               deps.TableLabels,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		EmptyState: types.TableEmptyState{
			Title:   "No engagements",
			Message: "No engagements found for this client.",
		},
	}

	if addURL != "" {
		tc.PrimaryAction = &types.PrimaryAction{
			Label:     deps.Labels.Detail.Tabs.Engagements,
			ActionURL: addURL,
			Icon:      "icon-plus",
		}
	}

	types.ApplyTableSettings(tc)
	return tc
}

// clientDisplayName returns the client's display name.
// Prefers client.name, falls back to user first+last name, then email.
func clientDisplayName(c *clientpb.Client) string {
	if name := c.GetName(); name != "" {
		return name
	}
	if u := c.GetUser(); u != nil {
		first := u.GetFirstName()
		last := u.GetLastName()
		if first != "" || last != "" {
			return first + " " + last
		}
		if u.GetEmailAddress() != "" {
			return u.GetEmailAddress()
		}
	}
	return c.GetId()
}

// buildFullAddress joins non-empty address parts into a single line.
func buildFullAddress(street, city, province, postalCode string) string {
	var parts []string
	if street != "" {
		parts = append(parts, street)
	}
	if city != "" {
		parts = append(parts, city)
	}
	if province != "" {
		parts = append(parts, province)
	}
	if postalCode != "" {
		parts = append(parts, postalCode)
	}
	return strings.Join(parts, ", ")
}

// loadClientTags fetches the tags assigned to a client by looking up
// client_category junction records and resolving category names.
func loadClientTags(ctx context.Context, deps *DetailViewDeps, clientID string) []TagChip {
	if deps.ListClientCategories == nil || deps.ListCategories == nil {
		return nil
	}

	// Load all categories to build ID -> Name map
	catResp, err := deps.ListCategories(ctx, &categorypb.ListCategoriesRequest{})
	if err != nil {
		log.Printf("Failed to load categories for client detail: %v", err)
		return nil
	}
	catNames := make(map[string]string)
	for _, cat := range catResp.GetData() {
		if cat.GetModule() == "client" {
			catNames[cat.GetId()] = cat.GetName()
		}
	}

	// Load junction records for this client
	ccResp, err := deps.ListClientCategories(ctx, &clientcategorypb.ListClientCategoriesRequest{})
	if err != nil {
		log.Printf("Failed to load client categories for detail: %v", err)
		return nil
	}

	var chips []TagChip
	for _, cc := range ccResp.GetData() {
		if cc.GetClientId() == clientID {
			if name, ok := catNames[cc.GetCategoryId()]; ok {
				chips = append(chips, TagChip{Name: name})
			}
		}
	}
	return chips
}

// loadPurchaseHistory fetches revenue records for a client, calculates stats,
// and returns sorted order rows (most recent first).
func loadPurchaseHistory(ctx context.Context, deps *DetailViewDeps, clientID string) (PurchaseStats, []OrderRow) {
	empty := PurchaseStats{
		LifetimeSpend: "PHP 0.00",
		AvgOrderValue: "PHP 0.00",
		LastPurchase:  "N/A",
	}

	if deps.ListRevenues == nil {
		return empty, nil
	}

	records, err := deps.ListRevenues(ctx, "revenue")
	if err != nil {
		log.Printf("Failed to load revenues for client %s: %v", clientID, err)
		return empty, nil
	}

	// Filter revenue records for this client
	var orders []OrderRow
	var totalSpend float64
	var lastPurchase string

	for _, r := range records {
		cid, _ := r["client_id"].(string)
		if cid != clientID {
			continue
		}

		id, _ := r["id"].(string)
		ref, _ := r["reference_number"].(string)
		date, _ := r["revenue_date_string"].(string)
		status, _ := r["status"].(string)
		currency, _ := r["currency"].(string)
		if currency == "" {
			currency = "PHP"
		}

		// Parse amount — can be string or float64 from DB
		var amount float64
		switch v := r["total_amount"].(type) {
		case float64:
			amount = v
		case string:
			amount, _ = strconv.ParseFloat(v, 64)
		}

		totalSpend += amount
		amountStr := fmt.Sprintf("%s %.2f", currency, amount)

		variant := "default"
		switch status {
		case "active":
			variant = "info"
		case "completed":
			variant = "success"
		case "cancelled":
			variant = "warning"
		}

		orders = append(orders, OrderRow{
			ID:        id,
			Reference: ref,
			Date:      date,
			Amount:    amountStr,
			Status:    status,
			Variant:   variant,
		})

		// Track most recent purchase date
		if date > lastPurchase {
			lastPurchase = date
		}
	}

	// Sort orders by date descending (most recent first)
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Date > orders[j].Date
	})

	totalOrders := len(orders)
	avgOrder := 0.0
	if totalOrders > 0 {
		avgOrder = totalSpend / float64(totalOrders)
	}

	if lastPurchase == "" {
		lastPurchase = "N/A"
	}

	stats := PurchaseStats{
		LifetimeSpend: fmt.Sprintf("PHP %.2f", totalSpend),
		TotalOrders:   totalOrders,
		AvgOrderValue: fmt.Sprintf("PHP %.2f", avgOrder),
		LastPurchase:  lastPurchase,
	}

	return stats, orders
}

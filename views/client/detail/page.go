package detail

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
)

// Deps holds view dependencies.
type Deps struct {
	ReadClient           func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	ListCategories       func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListClientCategories func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
	ListRevenues         func(ctx context.Context, collection string) ([]map[string]any, error)
	Labels               entydad.ClientLabels
	CommonLabels         pyeza.CommonLabels
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

// PageData holds the data for the client detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Client          *clientpb.Client
	Labels          entydad.ClientLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	ClientName      string
	ClientEmail     string
	ClientPhone     string
	ClientStatus    string
	StatusVariant   string
	// CRM fields
	CompanyName   string
	CustomerType  string
	DateOfBirth   string
	StreetAddress string
	City          string
	Province      string
	PostalCode    string
	Notes         string
	FullAddress   string
	Tags          []TagChip
	// Has* booleans for conditional rendering in templates
	HasCompany  bool
	HasPersonal bool
	HasAddress  bool
	HasNotes    bool
	HasTags     bool
	// Purchase history
	PurchaseStats PurchaseStats
	Orders        []OrderRow
	HasOrders     bool
}

// NewView creates the client detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "basic"
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
		clientEmail := ""
		clientPhone := ""
		if u != nil {
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

		tabItems := buildTabItems(id)

		// CRM fields
		companyName := client.GetCompanyName()
		customerType := client.GetCustomerType()
		dateOfBirth := client.GetDateOfBirth()
		streetAddress := client.GetStreetAddress()
		city := client.GetCity()
		province := client.GetProvince()
		postalCode := client.GetPostalCode()
		notes := client.GetNotes()
		fullAddress := buildFullAddress(streetAddress, city, province, postalCode)

		hasCompany := companyName != "" || customerType != ""
		hasPersonal := dateOfBirth != ""
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
				ActiveNav:      "clients",
				HeaderTitle:    clientName,
				HeaderSubtitle: clientEmail,
				HeaderIcon:     "icon-user",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "client-detail-content",
			Client:          client,
			Labels:          deps.Labels,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
			ClientName:      clientName,
			ClientEmail:     clientEmail,
			ClientPhone:     clientPhone,
			ClientStatus:    clientStatus,
			StatusVariant:   statusVariant,
			CompanyName:     companyName,
			CustomerType:    customerType,
			DateOfBirth:     dateOfBirth,
			StreetAddress:   streetAddress,
			City:            city,
			Province:        province,
			PostalCode:      postalCode,
			Notes:           notes,
			FullAddress:     fullAddress,
			Tags:            tags,
			HasCompany:      hasCompany,
			HasPersonal:     hasPersonal,
			HasAddress:      hasAddress,
			HasNotes:        hasNotes,
			HasTags:         hasTags,
			PurchaseStats:   stats,
			Orders:          orders,
			HasOrders:       hasOrders,
		}

		return view.OK("client-detail", pageData)
	})
}

func buildTabItems(id string) []pyeza.TabItem {
	base := "/app/clients/detail/" + id
	action := "/action/clients/" + id + "/tab/"
	return []pyeza.TabItem{
		{Key: "basic", Label: "Basic Information", Href: base + "?tab=basic", HxGet: action + "basic", Icon: "icon-info"},
		{Key: "history", Label: "Purchase History", Href: base + "?tab=history", HxGet: action + "history", Icon: "icon-shopping-bag"},
	}
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "basic"
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
		clientEmail := ""
		clientPhone := ""
		if u != nil {
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
			Client:        client,
			Labels:        deps.Labels,
			ActiveTab:     tab,
			TabItems:      buildTabItems(id),
			ClientName:    clientName,
			ClientEmail:   clientEmail,
			ClientPhone:   clientPhone,
			ClientStatus:  clientStatus,
			StatusVariant: statusVariant,
		}

		switch tab {
		case "basic":
			pageData.CompanyName = client.GetCompanyName()
			pageData.CustomerType = client.GetCustomerType()
			pageData.DateOfBirth = client.GetDateOfBirth()
			pageData.StreetAddress = client.GetStreetAddress()
			pageData.City = client.GetCity()
			pageData.Province = client.GetProvince()
			pageData.PostalCode = client.GetPostalCode()
			pageData.Notes = client.GetNotes()
			pageData.FullAddress = buildFullAddress(pageData.StreetAddress, pageData.City, pageData.Province, pageData.PostalCode)
			pageData.HasCompany = pageData.CompanyName != "" || pageData.CustomerType != ""
			pageData.HasPersonal = pageData.DateOfBirth != ""
			pageData.HasAddress = pageData.StreetAddress != "" || pageData.City != "" || pageData.Province != "" || pageData.PostalCode != ""
			pageData.HasNotes = pageData.Notes != ""
			pageData.Tags = loadClientTags(ctx, deps, id)
			pageData.HasTags = len(pageData.Tags) > 0
		case "history":
			pageData.PurchaseStats, pageData.Orders = loadPurchaseHistory(ctx, deps, id)
			pageData.HasOrders = len(pageData.Orders) > 0
		}

		templateName := "client-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// clientDisplayName returns the client's display name from the embedded user.
func clientDisplayName(c *clientpb.Client) string {
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
func loadClientTags(ctx context.Context, deps *Deps, clientID string) []TagChip {
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
func loadPurchaseHistory(ctx context.Context, deps *Deps, clientID string) (PurchaseStats, []OrderRow) {
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

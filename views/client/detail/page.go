package detail

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"

	categorypb       "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb         "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
	clientstmtpb     "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/client_statement"
	revenuepb        "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	collectionpb     "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	subscriptionpb   "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes                entydad.ClientRoutes
	ReadClient            func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	ListCategories        func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListClientCategories  func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
	ListRevenues                func(ctx context.Context, collection string) ([]map[string]any, error)
	GetClientStatement          func(ctx context.Context, req *clientstmtpb.ClientStatementRequest) (*clientstmtpb.ClientStatementResponse, error)
	ListSubscriptions           func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)
	GetSubscriptionListPageData func(ctx context.Context, req *subscriptionpb.GetSubscriptionListPageDataRequest) (*subscriptionpb.GetSubscriptionListPageDataResponse, error)
	SubscriptionAddURL    string
	SubscriptionDetailURL string
	// SubscriptionUnderClientDetailURL is the nested-route template (e.g.
	// "/app/clients/detail/{client_id}/subscriptions/{id}"). When set, the
	// engagements row link uses this so the subscription detail page renders
	// with a "client → subscription" breadcrumb. Falls back to the flat
	// SubscriptionDetailURL when empty.
	SubscriptionUnderClientDetailURL string
	SubscriptionEditURL   string
	SubscriptionDeleteURL string
	Labels                entydad.ClientLabels
	CommonLabels          pyeza.CommonLabels
	TableLabels           types.TableLabels

	// Attachment operations (embedded from hybra)
	attachment.AttachmentOps

	// Audit log operations (embedded from hybra)
	auditlog.AuditOps

	// ListClientPriceSchedules is an optional callback that fetches PriceSchedules
	// scoped to a specific client_id (price_schedule.client_id == clientID).
	// When nil the PriceSchedules tab silently renders empty state.
	// It is wired from centymo's ListPriceSchedules use case filtered by client_id
	// so entydad never imports the centymo PriceSchedule repo directly.
	ListClientPriceSchedules func(ctx context.Context, clientID string) ([]ClientPriceScheduleRow, error)

	// PriceScheduleAddURL is the centymo PriceSchedule-add drawer URL.
	// The PriceSchedules tab appends ?context=client&client_id={cid} to
	// pre-fill and lock the client field.
	PriceScheduleAddURL string

	// ListRevenuesByClient returns all Revenue rows for the given client_id,
	// ordered by revenue_date desc. Used by the outstanding-revenue table on
	// the Statement tab.
	ListRevenuesByClient func(ctx context.Context, clientID string) ([]*revenuepb.Revenue, error)

	// ListCollectionsByClient returns all Collection rows whose linked revenue
	// has the given client_id. Used to compute paid-amount per revenue on the
	// outstanding-revenue table.
	ListCollectionsByClient func(ctx context.Context, clientID string) ([]*collectionpb.Collection, error)

	// ListRevenueRunCandidates enumerates un-invoiced billing periods for the
	// given scope. Wired from espyna's consumer.ListRevenueRunCandidates via
	// block.go. Nil-safe: if nil, the Revenue Run drawer returns an error.
	ListRevenueRunCandidates func(ctx context.Context, scope RevenueRunScope) ([]RevenueRunCandidate, string, error)

	// GenerateRevenueRun executes a batch revenue generation run. Wired from
	// espyna's consumer.GenerateRevenueRun via block.go. Nil-safe: if nil, the
	// Revenue Run drawer returns an error.
	GenerateRevenueRun func(ctx context.Context, scope RevenueRunScope, selections RevenueRunSelections) (*RevenueRunResult, error)
}

// TagChip represents a tag displayed as a chip on the detail page.
type TagChip struct {
	Name string
}

// SubscriptionRow represents a single subscription in the subscriptions tab.
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
	Country       string
	Website       string
	Notes         string
	FullAddress   string
	Tags          []TagChip
	// Has* booleans for conditional rendering in templates
	HasName    bool
	HasAddress bool
	HasNotes   bool
	HasTags    bool
	// Subscriptions tab
	Subscriptions      []SubscriptionRow
	SubscriptionsTable *types.TableConfig
	// PriceSchedules tab
	PriceSchedulesTable *types.TableConfig
	// Accounting tab
	BillingCurrency string
	// Statement tab
	StatementEntries        []*clientstmtpb.StatementEntry
	StatementSummary        *clientstmtpb.ClientStatementSummary
	StatementSummaryDisplay *StatementSummaryDisplay
	StatementTable          *types.TableConfig
	OutstandingTable        *types.TableConfig
	// Attachments tab
	AttachmentTable *types.TableConfig
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// StatementSummaryDisplay holds pre-formatted money cells for the statement summary bar.
type StatementSummaryDisplay struct {
	OutstandingBalance types.TableCell
	TotalBilled        types.TableCell
	TotalReceived      types.TableCell
}

// NewView creates the client detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := deps.Labels.Detail.Tabs.CanonicalizeTab(viewCtx.Request.URL.Query().Get("tab"))
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

		tabItems := buildTabItems(id, deps, countClientSubscriptions(ctx, deps, id), countClientPriceSchedules(ctx, deps, id))

		// CRM fields
		name := client.GetName()
		streetAddress := client.GetStreetAddress()
		city := client.GetCity()
		province := client.GetProvince()
		postalCode := client.GetPostalCode()
		country := client.GetCountry()
		website := client.GetWebsite()
		notes := client.GetNotes()
		fullAddress := buildFullAddress(streetAddress, city, province, postalCode)

		hasName := name != ""
		hasAddress := streetAddress != "" || city != "" || province != "" || postalCode != ""
		hasNotes := notes != ""

		// Load tags for this client
		tags := loadClientTags(ctx, deps, id)
		hasTags := len(tags) > 0

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
			SubscriptionAddURL: buildSubscriptionAddURL(deps.SubscriptionAddURL, id, clientName, client.GetBillingCurrency()),
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
			Country:            country,
			Website:            website,
			Notes:              notes,
			FullAddress:        fullAddress,
			Tags:               tags,
			HasName:            hasName,
			HasAddress:         hasAddress,
			HasNotes:           hasNotes,
			HasTags:            hasTags,
			BillingCurrency:    client.GetBillingCurrency(),
		}

		// Load tab-specific data for the active tab on full page load
		switch activeTab {
		case "subscriptions":
			subs := loadClientSubscriptions(ctx, deps, id)
			pageData.Subscriptions = subs
			pageData.SubscriptionsTable = buildSubscriptionsTable(subs, pageData.SubscriptionAddURL, id, clientName, deps)
		case "priceSchedules":
			pageData.PriceSchedulesTable = buildPriceSchedulesTable(ctx, deps, id, clientName)
		case "statement":
			if deps.GetClientStatement != nil {
				req := &clientstmtpb.ClientStatementRequest{
					ClientId: id,
				}
				resp, err := deps.GetClientStatement(ctx, req)
				if err == nil && resp.Success {
					pageData.StatementEntries = resp.Entries
					pageData.StatementSummary = resp.Summary
					if resp.Summary != nil {
						pageData.StatementSummaryDisplay = &StatementSummaryDisplay{
							OutstandingBalance: types.MoneyCell(float64(resp.Summary.OutstandingBalance), "", true),
							TotalBilled:        types.MoneyCell(float64(resp.Summary.TotalBilled), "", true),
							TotalReceived:      types.MoneyCell(float64(resp.Summary.TotalReceived), "", true),
						}
					}
				}
			}
			pageData.OutstandingTable = buildOutstandingRevenueTable(ctx, deps, id)
		case "attachments":
			loadAttachments(ctx, deps, id, pageData)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.Request.URL.Query().Get("cursor")
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "client",
					EntityID:    id,
					Limit:       20,
					CursorToken: cursor,
				})
				if err != nil {
					log.Printf("Failed to load audit history: %v", err)
				}
				if auditResp != nil {
					pageData.AuditEntries = auditResp.Entries
					pageData.AuditHasNext = auditResp.HasNext
					pageData.AuditNextCursor = auditResp.NextCursor
				}
			}
			pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
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

func buildTabItems(id string, deps *DetailViewDeps, subscriptionCount, priceScheduleCount int) []pyeza.TabItem {
	routes := deps.Routes
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	subscriptionsSlug := deps.Labels.Detail.Tabs.ResolveTabSlug("subscriptions")
	priceSchedulesSlug := deps.Labels.Detail.Tabs.ResolveTabSlug("priceSchedules")
	return []pyeza.TabItem{
		{Key: "info", Label: deps.Labels.Detail.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "representative", Label: deps.Labels.Detail.Tabs.Representative, Href: base + "?tab=representative", HxGet: action + "representative", Icon: "icon-user"},
		{Key: "priceSchedules", Label: deps.Labels.Detail.Tabs.PriceSchedules, Href: base + "?tab=" + priceSchedulesSlug, HxGet: action + priceSchedulesSlug, Icon: "icon-calendar", Count: priceScheduleCount},
		{Key: "subscriptions", Label: deps.Labels.Detail.Tabs.Subscriptions, Href: base + "?tab=" + subscriptionsSlug, HxGet: action + subscriptionsSlug, Icon: "icon-file-text", Count: subscriptionCount},
		{Key: "statement", Label: deps.Labels.Detail.Tabs.Statement, Href: base + "?tab=statement", HxGet: action + "statement", Icon: "icon-file-text"},
		{Key: "attachments", Label: deps.Labels.Detail.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		{Key: "audit-history", Label: deps.Labels.Detail.Tabs.AuditHistory, Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := deps.Labels.Detail.Tabs.CanonicalizeTab(viewCtx.Request.PathValue("tab"))
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
			TabItems:           buildTabItems(id, deps, countClientSubscriptions(ctx, deps, id), countClientPriceSchedules(ctx, deps, id)),
			ClientName:         clientName,
			RepresentativeName: representativeName,
			ClientEmail:        clientEmail,
			ClientPhone:        clientPhone,
			ClientStatus:       clientStatus,
			StatusVariant:      statusVariant,
			EditURL:            route.ResolveURL(deps.Routes.EditURL, "id", id),
			SubscriptionAddURL: buildSubscriptionAddURL(deps.SubscriptionAddURL, id, clientName, client.GetBillingCurrency()),
		}

		switch tab {
		case "info":
			pageData.Name = client.GetName()
			pageData.StreetAddress = client.GetStreetAddress()
			pageData.City = client.GetCity()
			pageData.Province = client.GetProvince()
			pageData.PostalCode = client.GetPostalCode()
			pageData.Country = client.GetCountry()
			pageData.Website = client.GetWebsite()
			pageData.Notes = client.GetNotes()
			pageData.FullAddress = buildFullAddress(pageData.StreetAddress, pageData.City, pageData.Province, pageData.PostalCode)
			pageData.HasName = pageData.Name != ""
			pageData.HasAddress = pageData.StreetAddress != "" || pageData.City != "" || pageData.Province != "" || pageData.PostalCode != ""
			pageData.HasNotes = pageData.Notes != ""
			pageData.Tags = loadClientTags(ctx, deps, id)
			pageData.HasTags = len(pageData.Tags) > 0
		case "representative":
			// user fields already on client via GetUser()
		case "subscriptions":
			pageData.Subscriptions = loadClientSubscriptions(ctx, deps, id)
			pageData.SubscriptionsTable = buildSubscriptionsTable(pageData.Subscriptions, pageData.SubscriptionAddURL, id, clientName, deps)
		case "priceSchedules":
			pageData.PriceSchedulesTable = buildPriceSchedulesTable(ctx, deps, id, clientName)
		case "statement":
			if deps.GetClientStatement != nil {
				req := &clientstmtpb.ClientStatementRequest{
					ClientId: id,
				}
				resp, err := deps.GetClientStatement(ctx, req)
				if err == nil && resp.Success {
					pageData.StatementEntries = resp.Entries
					pageData.StatementSummary = resp.Summary
					if resp.Summary != nil {
						pageData.StatementSummaryDisplay = &StatementSummaryDisplay{
							OutstandingBalance: types.MoneyCell(float64(resp.Summary.OutstandingBalance), "", true),
							TotalBilled:        types.MoneyCell(float64(resp.Summary.TotalBilled), "", true),
							TotalReceived:      types.MoneyCell(float64(resp.Summary.TotalReceived), "", true),
						}
					}
				}
			}
			pageData.OutstandingTable = buildOutstandingRevenueTable(ctx, deps, id)
		case "attachments":
			loadAttachments(ctx, deps, id, pageData)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.Request.URL.Query().Get("cursor")
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "client",
					EntityID:    id,
					Limit:       20,
					CursorToken: cursor,
				})
				if err != nil {
					log.Printf("Failed to load audit history: %v", err)
				}
				if auditResp != nil {
					pageData.AuditEntries = auditResp.Entries
					pageData.AuditHasNext = auditResp.HasNext
					pageData.AuditNextCursor = auditResp.NextCursor
				}
			}
			pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
		}

		templateName := "client-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		if tab == "priceSchedules" {
			templateName = "client-tab-priceSchedules"
		}
		return view.OK(templateName, pageData)
	})
}

// countClientPriceSchedules returns the count of price_schedules scoped to
// this client. Reuses the same callback that powers the PriceSchedules tab —
// returns 0 if the dep is unwired or the call errors. The count surfaces in
// the tab header badge alongside the Subscriptions count.
func countClientPriceSchedules(ctx context.Context, deps *DetailViewDeps, clientID string) int {
	if deps.ListClientPriceSchedules == nil {
		return 0
	}
	rows, err := deps.ListClientPriceSchedules(ctx, clientID)
	if err != nil {
		log.Printf("Failed to count price schedules for client %s: %v", clientID, err)
		return 0
	}
	return len(rows)
}

// countClientSubscriptions returns the active-subscription count for a client.
// Uses GetSubscriptionListPageData with a single-row pagination so the response
// payload stays small; the count is taken from pagination.TotalItems. Returns 0
// on any error.
func countClientSubscriptions(ctx context.Context, deps *DetailViewDeps, clientID string) int {
	if deps.GetSubscriptionListPageData == nil {
		return 0
	}
	limit := int32(1)
	resp, err := deps.GetSubscriptionListPageData(ctx, &subscriptionpb.GetSubscriptionListPageDataRequest{
		Pagination: &categorypb.PaginationRequest{Limit: limit},
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
		log.Printf("Failed to count subscriptions for client %s: %v", clientID, err)
		return 0
	}
	return int(resp.GetPagination().GetTotalItems())
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
	tz := types.LocationFromContext(ctx)
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
		dateStart := types.FormatTimestampInTZ(s.GetDateTimeStart(), tz, types.DateTimeReadable)
		dateEnd := types.FormatTimestampInTZ(s.GetDateTimeEnd(), tz, types.DateTimeReadable)
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

// buildSubscriptionAddURL appends client_id, client_name, and (when set)
// billing_currency so the subscription drawer can scope the plan search.
func buildSubscriptionAddURL(base, clientID, clientName, billingCurrency string) string {
	u := base + "?client_id=" + clientID + "&client_name=" + url.QueryEscape(clientName)
	if billingCurrency != "" {
		u += "&billing_currency=" + url.QueryEscape(billingCurrency)
	}
	return u
}

// buildSubscriptionsTable builds a TableConfig for the subscriptions tab.
// The table is always returned (even when empty) so the primary action
// stays visible and the table's own empty state renders.
func buildSubscriptionsTable(rows []SubscriptionRow, addURL string, clientID string, clientName string, deps *DetailViewDeps) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "name", Label: deps.Labels.Detail.Subscriptions.ColumnName},
		{Key: "plan", Label: deps.Labels.Detail.Subscriptions.ColumnPlan},
		{Key: "start_date", Label: deps.Labels.Detail.Subscriptions.ColumnStartDate, WidthClass: "col-3xl"},
		{Key: "end_date", Label: deps.Labels.Detail.Subscriptions.ColumnEndDate, WidthClass: "col-3xl"},
	}

	// Build locked client query params for edit URLs
	clientParams := "?client_id=" + clientID + "&client_name=" + url.QueryEscape(clientName)

	var tableRows []types.TableRow
	for _, r := range rows {
		editURL := route.ResolveURL(deps.SubscriptionEditURL, "id", r.ID) + clientParams
		detailURL := route.ResolveURL(deps.SubscriptionDetailURL, "id", r.ID)
		if deps.SubscriptionUnderClientDetailURL != "" {
			detailURL = route.ResolveURL(deps.SubscriptionUnderClientDetailURL, "client_id", clientID, "id", r.ID)
		}

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
				{Type: "view", Label: deps.CommonLabels.Actions.View, Action: "view", Href: detailURL},
				{Type: "edit", Label: deps.CommonLabels.Actions.Edit, Action: "edit", URL: editURL, DrawerTitle: r.Name},
				{Type: "delete", Label: deps.CommonLabels.Actions.Delete, Action: "delete", URL: deps.SubscriptionDeleteURL, ItemName: r.Name, ConfirmTitle: deps.Labels.Detail.Subscriptions.ConfirmDeleteTitle, ConfirmMessage: fmt.Sprintf(deps.Labels.Detail.Subscriptions.ConfirmDeleteMessage, r.Name)},
			},
		})
	}

	types.ApplyColumnStyles(columns, tableRows)

	tc := &types.TableConfig{
		ID:                   "subscriptions-table",
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
			Title:   deps.Labels.Detail.EmptySubscriptionsTitle,
			Message: deps.Labels.Detail.EmptySubscriptions,
		},
	}

	if addURL != "" {
		tc.PrimaryAction = &types.PrimaryAction{
			Label:     deps.Labels.Detail.AddSubscription,
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

// capitalizeType capitalizes the first letter of a type string.
func capitalizeType(t string) string {
	if len(t) == 0 {
		return t
	}
	return strings.ToUpper(t[:1]) + t[1:]
}

// buildStatementTable builds a TableConfig for the statement tab.
func buildStatementTable(resp *clientstmtpb.ClientStatementResponse, deps *DetailViewDeps) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: deps.Labels.Detail.Statement.ColumnDate, WidthClass: "col-2xl"},
		{Key: "type", Label: deps.Labels.Detail.Statement.ColumnType, WidthClass: "col-lg"},
		{Key: "reference", Label: deps.Labels.Detail.Statement.ColumnReference, WidthClass: "col-3xl"},
		{Key: "description", Label: deps.Labels.Detail.Statement.ColumnDescription},
		{Key: "billed", Label: deps.Labels.Detail.Statement.ColumnBilled, WidthClass: "col-3xl", Align: "right"},
		{Key: "received", Label: deps.Labels.Detail.Statement.ColumnReceived, WidthClass: "col-3xl", Align: "right"},
		{Key: "balance", Label: deps.Labels.Detail.Statement.ColumnBalance, WidthClass: "col-3xl", Align: "right"},
	}

	var rows []types.TableRow
	for _, entry := range resp.Entries {
		var billedCell types.TableCell
		if entry.Billed > 0 {
			billedCell = types.MoneyCell(float64(entry.Billed), "", true)
		}
		var receivedCell types.TableCell
		if entry.Received > 0 {
			receivedCell = types.MoneyCell(float64(entry.Received), "", true)
		}
		rows = append(rows, types.TableRow{
			ID: entry.EntityId,
			Cells: []types.TableCell{
				{Type: "text", Value: entry.Date},
				{Type: "text", Value: capitalizeType(entry.Type)},
				{Type: "text", Value: entry.ReferenceNumber},
				{Type: "text", Value: entry.Description},
				billedCell,
				receivedCell,
				types.MoneyCell(float64(entry.Balance), "", true),
			},
		})
	}

	// Add summary/totals row if summary exists
	if resp.Summary != nil {
		s := resp.Summary
		rows = append(rows, types.TableRow{
			ID: "__totals__",
			Cells: []types.TableCell{
				{Type: "text", Value: ""},
				{Type: "text", Value: ""},
				{Type: "text", Value: ""},
				{Type: "text", Value: strings.ToUpper(deps.Labels.Detail.Statement.TotalsRowLabel)},
				types.MoneyCell(float64(s.TotalBilled), "", true),
				types.MoneyCell(float64(s.TotalReceived), "", true),
				types.MoneyCell(float64(s.OutstandingBalance), "", true),
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	tc := &types.TableConfig{
		ID:                   "clientStatementTable",
		Columns:              columns,
		Rows:                 rows,
		Labels:               deps.TableLabels,
		ShowSearch:           false,
		ShowSort:             false,
		ShowExport:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "date",
		DefaultSortDirection: "asc",
		EmptyState: types.TableEmptyState{
			Title:   deps.Labels.Detail.EmptyStatementTitle,
			Message: deps.Labels.Detail.EmptyStatementMessage,
		},
	}

	types.ApplyTableSettings(tc)
	return tc
}

// buildOutstandingRevenueTable builds a TableConfig showing only revenue rows
// that are not fully paid (outstanding amount > 0) for the given client.
// It calls ListRevenuesByClient and ListCollectionsByClient; if either dep is
// nil it returns nil so the template falls back to its empty-state panel.
func buildOutstandingRevenueTable(ctx context.Context, deps *DetailViewDeps, clientID string) *types.TableConfig {
	if deps.ListRevenuesByClient == nil || deps.ListCollectionsByClient == nil {
		return nil
	}

	revenues, err := deps.ListRevenuesByClient(ctx, clientID)
	if err != nil {
		log.Printf("buildOutstandingRevenueTable: failed to list revenues for client %s: %v", clientID, err)
		return nil
	}

	collections, err := deps.ListCollectionsByClient(ctx, clientID)
	if err != nil {
		log.Printf("buildOutstandingRevenueTable: failed to list collections for client %s: %v", clientID, err)
		return nil
	}

	// Aggregate paid amount by revenue_id (skip reversed/voided/inactive collections)
	paidByRevenue := make(map[string]int64, len(collections))
	for _, c := range collections {
		if !c.GetActive() {
			continue
		}
		s := c.GetStatus()
		if s == "reversed" || s == "voided" {
			continue
		}
		paidByRevenue[c.GetRevenueId()] += c.GetAmount()
	}

	labels := deps.Labels.Detail.OutstandingTable
	columns := []types.TableColumn{
		{Key: "revenue_date", Label: labels.Columns.Date, WidthClass: "col-2xl"},
		{Key: "reference_number", Label: labels.Columns.Reference, WidthClass: "col-3xl"},
		{Key: "name", Label: labels.Columns.Description},
		{Key: "due_date", Label: labels.Columns.DueDate, WidthClass: "col-2xl"},
		{Key: "billed", Label: labels.Columns.Billed, WidthClass: "col-3xl", Align: "right"},
		{Key: "paid", Label: labels.Columns.Paid, WidthClass: "col-3xl", Align: "right"},
		{Key: "outstanding", Label: labels.Columns.Outstanding, WidthClass: "col-3xl", Align: "right"},
		{Key: "status", Label: labels.Columns.Status, WidthClass: "col-lg"},
	}

	var rows []types.TableRow
	for _, r := range revenues {
		outstanding := r.GetTotalAmount() - paidByRevenue[r.GetId()]
		if outstanding <= 0 {
			continue
		}

		dueDate := r.GetDueDate()
		if dueDate == "" {
			dueDate = "—"
		}

		currency := r.GetCurrency()
		paid := paidByRevenue[r.GetId()]

		rows = append(rows, types.TableRow{
			ID: r.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: r.GetRevenueDate()},
				{Type: "text", Value: r.GetReferenceNumber()},
				{Type: "text", Value: r.GetName()},
				{Type: "text", Value: dueDate},
				types.MoneyCell(float64(r.GetTotalAmount()), currency, true),
				types.MoneyCell(float64(paid), currency, true),
				types.MoneyCell(float64(outstanding), currency, true),
				{Type: "badge", Value: r.GetStatus(), Variant: revenueStatusVariant(r.GetStatus())},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	tc := &types.TableConfig{
		ID:                   "client-statement-outstanding-table",
		Columns:              columns,
		Rows:                 rows,
		Labels:               deps.TableLabels,
		ShowSearch:           true,
		ShowSort:             true,
		ShowFilters:          false,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "revenue_date",
		DefaultSortDirection: "desc",
		BulkActions:          nil,
		EmptyState: types.TableEmptyState{
			Title:   labels.Empty.Title,
			Message: labels.Empty.Message,
		},
	}

	// Wire the "Run Invoices" CTA when the route and callbacks are available
	// and the operator has the required permissions.
	if deps.Routes.RevenueRunURL != "" &&
		deps.ListRevenueRunCandidates != nil &&
		deps.GenerateRevenueRun != nil &&
		labels.RunInvoicesLabel != "" {
		perms := view.GetUserPermissions(ctx)
		canRun := perms != nil && perms.Can("revenue", "create") && perms.Can("subscription", "read")
		tc.PrimaryAction = &types.PrimaryAction{
			Label:     labels.RunInvoicesLabel,
			ActionURL: route.ResolveURL(deps.Routes.RevenueRunURL, "id", clientID),
			Icon:      "icon-zap",
			Disabled:  !canRun,
		}
	}

	types.ApplyTableSettings(tc)
	return tc
}

// revenueStatusVariant maps a revenue status string to a badge variant.
func revenueStatusVariant(status string) string {
	switch status {
	case "complete":
		return "success"
	case "draft":
		return "warning"
	case "cancelled":
		return "default"
	default:
		return "default"
	}
}

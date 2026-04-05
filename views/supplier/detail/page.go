package detail

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"

	supplierpb      "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	suppstmtpb      "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes       entydad.SupplierRoutes
	ReadSupplier func(ctx context.Context, req *supplierpb.ReadSupplierRequest) (*supplierpb.ReadSupplierResponse, error)
	Labels       entydad.SupplierLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Attachment operations (embedded from hybra)
	attachment.AttachmentOps

	// Audit log operations (embedded from hybra)
	auditlog.AuditOps

	ListPurchaseOrders   func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)
	GetSupplierStatement func(ctx context.Context, req *suppstmtpb.SupplierStatementRequest) (*suppstmtpb.SupplierStatementResponse, error)
}

// PurchaseOrderRow holds display data for a single purchase order row.
type PurchaseOrderRow struct {
	ID          string
	PONumber    string
	Status      string
	TotalAmount string
	Currency    string
	OrderDate   string
}

// PageData holds the data for the supplier detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Supplier        *supplierpb.Supplier
	Labels          entydad.SupplierLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	// Company info
	CompanyName   string
	SupplierType  string
	InternalID    string
	Status        string
	StatusVariant string
	// Contact info (from user)
	ContactName  string
	ContactEmail string
	ContactPhone string
	// Financial info
	PaymentTerms       string
	CreditLimit        string
	DefaultCurrency    string
	LeadTimeDays       string
	TaxID              string
	RegistrationNumber string
	// Address info
	StreetAddress string
	City          string
	Province      string
	PostalCode    string
	Country       string
	FullAddress   string
	// Other
	Website string
	Notes   string
	// Has* booleans for conditional rendering in templates
	HasContact   bool
	HasFinancial bool
	HasAddress   bool
	HasNotes     bool
	// Attachments
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
	// Purchase Orders tab
	PurchaseOrders []PurchaseOrderRow
	// Statement tab
	StatementSummary        *suppstmtpb.SupplierStatementSummary
	StatementSummaryDisplay *StatementSummaryDisplay
	StatementTable          *types.TableConfig
}

// StatementSummaryDisplay holds pre-formatted string values for the statement summary bar.
type StatementSummaryDisplay struct {
	OutstandingBalance string
	TotalBilled        string
	TotalPaid          string
}

// NewView creates the supplier detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		resp, err := deps.ReadSupplier(ctx, &supplierpb.ReadSupplierRequest{
			Data: &supplierpb.Supplier{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read supplier %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load supplier: %w", err))
		}

		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("supplier not found"))
		}
		supplier := data[0]

		pageData := buildPageData(supplier, id, activeTab, viewCtx, deps)

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "client-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		switch activeTab {
		case "purchase-orders":
			if deps.ListPurchaseOrders != nil {
				supplierId := id
				poResp, poErr := deps.ListPurchaseOrders(ctx, &purchaseorderpb.ListPurchaseOrdersRequest{
					SupplierId: &supplierId,
				})
				if poErr == nil && poResp != nil {
					for _, po := range poResp.GetData() {
						pageData.PurchaseOrders = append(pageData.PurchaseOrders, PurchaseOrderRow{
							ID:          po.GetId(),
							PONumber:    po.GetPoNumber(),
							Status:      po.GetStatus(),
							TotalAmount: fmt.Sprintf("%.2f", float64(po.GetTotalAmount())/100.0),
							Currency:    po.GetCurrency(),
							OrderDate:   po.GetOrderDateString(),
						})
					}
				}
			}
		case "statement":
			if deps.GetSupplierStatement != nil {
				req := &suppstmtpb.SupplierStatementRequest{
					SupplierId: id,
				}
				stmtResp, stmtErr := deps.GetSupplierStatement(ctx, req)
				if stmtErr == nil && stmtResp.GetSuccess() {
					pageData.StatementSummary = stmtResp.GetSummary()
					if stmtResp.GetSummary() != nil {
						pageData.StatementSummaryDisplay = &StatementSummaryDisplay{
							OutstandingBalance: fmt.Sprintf("%.2f", float64(stmtResp.GetSummary().GetOutstandingBalance())/100),
							TotalBilled:        fmt.Sprintf("%.2f", float64(stmtResp.GetSummary().GetTotalBilled())/100),
							TotalPaid:          fmt.Sprintf("%.2f", float64(stmtResp.GetSummary().GetTotalPaid())/100),
						}
					}
					pageData.StatementTable = buildSupplierStatementTable(stmtResp, deps.TableLabels)
				}
			}
		case "attachments":
			loadAttachments(ctx, deps, id, pageData)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.Request.URL.Query().Get("cursor")
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "supplier",
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

		return view.OK("supplier-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial -- returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadSupplier(ctx, &supplierpb.ReadSupplierRequest{
			Data: &supplierpb.Supplier{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read supplier %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load supplier: %w", err))
		}

		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("supplier not found"))
		}
		supplier := data[0]

		pageData := buildPageData(supplier, id, tab, viewCtx, deps)

		switch tab {
		case "purchase-orders":
			if deps.ListPurchaseOrders != nil {
				supplierId := id
				poResp, poErr := deps.ListPurchaseOrders(ctx, &purchaseorderpb.ListPurchaseOrdersRequest{
					SupplierId: &supplierId,
				})
				if poErr == nil && poResp != nil {
					for _, po := range poResp.GetData() {
						pageData.PurchaseOrders = append(pageData.PurchaseOrders, PurchaseOrderRow{
							ID:          po.GetId(),
							PONumber:    po.GetPoNumber(),
							Status:      po.GetStatus(),
							TotalAmount: fmt.Sprintf("%.2f", float64(po.GetTotalAmount())/100.0),
							Currency:    po.GetCurrency(),
							OrderDate:   po.GetOrderDateString(),
						})
					}
				}
			}
		case "statement":
			if deps.GetSupplierStatement != nil {
				req := &suppstmtpb.SupplierStatementRequest{
					SupplierId: id,
				}
				stmtResp, stmtErr := deps.GetSupplierStatement(ctx, req)
				if stmtErr == nil && stmtResp.GetSuccess() {
					pageData.StatementSummary = stmtResp.GetSummary()
					if stmtResp.GetSummary() != nil {
						pageData.StatementSummaryDisplay = &StatementSummaryDisplay{
							OutstandingBalance: fmt.Sprintf("%.2f", float64(stmtResp.GetSummary().GetOutstandingBalance())/100),
							TotalBilled:        fmt.Sprintf("%.2f", float64(stmtResp.GetSummary().GetTotalBilled())/100),
							TotalPaid:          fmt.Sprintf("%.2f", float64(stmtResp.GetSummary().GetTotalPaid())/100),
						}
					}
					pageData.StatementTable = buildSupplierStatementTable(stmtResp, deps.TableLabels)
				}
			}
			return view.OK("supplier-tab-statement", pageData)
		case "attachments":
			loadAttachments(ctx, deps, id, pageData)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.Request.URL.Query().Get("cursor")
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "supplier",
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

		templateName := "supplier-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

func buildPageData(supplier *supplierpb.Supplier, id, activeTab string, viewCtx *view.ViewContext, deps *DetailViewDeps) *PageData {
	u := supplier.GetUser()

	companyName := supplier.GetCompanyName()
	supplierType := supplier.GetSupplierType()
	internalID := supplier.GetInternalId()

	status := supplier.GetStatus()
	if status == "" {
		if supplier.GetActive() {
			status = "active"
		} else {
			status = "blocked"
		}
	}
	statusVariant := "success"
	switch status {
	case "blocked":
		statusVariant = "danger"
	case "on_hold":
		statusVariant = "warning"
	}

	contactName := ""
	contactEmail := ""
	contactPhone := ""
	if u != nil {
		contactName = u.GetFirstName() + " " + u.GetLastName()
		contactEmail = u.GetEmailAddress()
		contactPhone = u.GetMobileNumber()
	}
	hasContact := contactName != "" || contactEmail != "" || contactPhone != ""

	paymentTerms := supplier.GetPaymentTerms()
	creditLimit := ""
	if cl := supplier.GetCreditLimit(); cl > 0 {
		creditLimit = fmt.Sprintf("%.2f", float64(cl)/100.0)
	}
	defaultCurrency := supplier.GetDefaultCurrency()
	leadTimeDays := ""
	if ltd := supplier.GetLeadTimeDays(); ltd > 0 {
		leadTimeDays = fmt.Sprintf("%d", ltd)
	}
	taxID := supplier.GetTaxId()
	registrationNumber := supplier.GetRegistrationNumber()
	hasFinancial := paymentTerms != "" || creditLimit != "" || defaultCurrency != "" || leadTimeDays != "" || taxID != "" || registrationNumber != ""

	streetAddress := supplier.GetStreetAddress()
	city := supplier.GetCity()
	province := supplier.GetProvince()
	postalCode := supplier.GetPostalCode()
	country := supplier.GetCountry()
	fullAddress := buildFullAddress(streetAddress, city, province, postalCode, country)
	hasAddress := streetAddress != "" || city != "" || province != "" || postalCode != "" || country != ""

	website := supplier.GetWebsite()
	notes := supplier.GetNotes()
	hasNotes := notes != ""

	displayName := companyName
	if displayName == "" {
		displayName = contactName
	}
	if displayName == "" {
		displayName = id
	}

	tabItems := buildTabItems(id, deps)

	return &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          displayName,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "supplier",
			HeaderTitle:    displayName,
			HeaderSubtitle: supplierType,
			HeaderIcon:     "icon-truck",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:    "supplier-detail-content",
		Supplier:           supplier,
		Labels:             deps.Labels,
		ActiveTab:          activeTab,
		TabItems:           tabItems,
		CompanyName:        companyName,
		SupplierType:       supplierType,
		InternalID:         internalID,
		Status:             status,
		StatusVariant:      statusVariant,
		ContactName:        contactName,
		ContactEmail:       contactEmail,
		ContactPhone:       contactPhone,
		PaymentTerms:       paymentTerms,
		CreditLimit:        creditLimit,
		DefaultCurrency:    defaultCurrency,
		LeadTimeDays:       leadTimeDays,
		TaxID:              taxID,
		RegistrationNumber: registrationNumber,
		StreetAddress:      streetAddress,
		City:               city,
		Province:           province,
		PostalCode:         postalCode,
		Country:            country,
		FullAddress:        fullAddress,
		Website:            website,
		Notes:              notes,
		HasContact:         hasContact,
		HasFinancial:       hasFinancial,
		HasAddress:         hasAddress,
		HasNotes:           hasNotes,
	}
}

func buildTabItems(id string, deps *DetailViewDeps) []pyeza.TabItem {
	routes := deps.Routes
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: deps.Labels.Detail.InfoTab, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "purchase-orders", Label: "Purchase Orders", Href: base + "?tab=purchase-orders", HxGet: action + "purchase-orders", Icon: "icon-shopping-cart"},
		{Key: "statement", Label: deps.Labels.Detail.StatementTab, Href: base + "?tab=statement", HxGet: action + "statement", Icon: "icon-file-text"},
		{Key: "attachments", Label: deps.Labels.Detail.AttachmentsTab, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		{Key: "audit-history", Label: "History", Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

func buildSupplierStatementTable(resp *suppstmtpb.SupplierStatementResponse, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: "Date"},
		{Key: "type", Label: "Type"},
		{Key: "reference", Label: "Reference"},
		{Key: "description", Label: "Description"},
		{Key: "billed", Label: "Billed", Align: "right"},
		{Key: "paid", Label: "Paid", Align: "right"},
		{Key: "balance", Label: "Balance", Align: "right"},
	}

	var rows []types.TableRow
	for _, entry := range resp.Entries {
		billedStr := ""
		if entry.Billed > 0 {
			billedStr = fmt.Sprintf("%.2f", float64(entry.Billed)/100)
		}
		paidStr := ""
		if entry.Paid > 0 {
			paidStr = fmt.Sprintf("%.2f", float64(entry.Paid)/100)
		}
		entryType := entry.Type
		if len(entryType) > 0 {
			entryType = strings.ToUpper(entryType[:1]) + entryType[1:]
		}
		rows = append(rows, types.TableRow{
			ID: entry.EntityId,
			Cells: []types.TableCell{
				{Value: entry.Date},
				{Value: entryType},
				{Value: entry.ReferenceNumber},
				{Value: entry.Description},
				{Value: billedStr},
				{Value: paidStr},
				{Value: fmt.Sprintf("%.2f", float64(entry.Balance)/100)},
			},
		})
	}

	if resp.Summary != nil {
		s := resp.Summary
		rows = append(rows, types.TableRow{
			ID: "__totals__",
			Cells: []types.TableCell{
				{}, {}, {},
				{Value: "TOTAL"},
				{Value: fmt.Sprintf("%.2f", float64(s.TotalBilled)/100)},
				{Value: fmt.Sprintf("%.2f", float64(s.TotalPaid)/100)},
				{Value: fmt.Sprintf("%.2f", float64(s.OutstandingBalance)/100)},
			},
		})
	}

	return &types.TableConfig{
		ID:              "supplierStatementTable",
		NameColumnLabel: "Date",
		Columns:         columns,
		Rows:            rows,
		ShowSearch:      false,
		ShowFilters:     false,
		ShowSort:        false,
		ShowExport:      true,
		ShowEntries:     true,
		Labels:          tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   "No Statement Entries",
			Message: "There are no transactions for this supplier.",
		},
	}
}

// buildFullAddress joins non-empty address parts into a single line.
func buildFullAddress(street, city, province, postalCode, country string) string {
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
	if country != "" {
		parts = append(parts, country)
	}
	return strings.Join(parts, ", ")
}

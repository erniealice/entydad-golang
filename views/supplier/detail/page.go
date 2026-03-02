package detail

import (
	"context"
	"fmt"
	"log"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"

	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
)

// Deps holds view dependencies.
type Deps struct {
	Routes       entydad.SupplierRoutes
	ReadSupplier func(ctx context.Context, req *supplierpb.ReadSupplierRequest) (*supplierpb.ReadSupplierResponse, error)
	Labels       entydad.SupplierLabels
	CommonLabels pyeza.CommonLabels
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
	CompanyName  string
	SupplierType string
	InternalID   string
	Status       string
	StatusVariant string
	// Contact info (from user)
	ContactName string
	ContactEmail string
	ContactPhone string
	// Financial info
	PaymentTerms    string
	CreditLimit     string
	DefaultCurrency string
	LeadTimeDays    string
	TaxID           string
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
}

// NewView creates the supplier detail view.
func NewView(deps *Deps) view.View {
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

		return view.OK("supplier-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial -- returns only the tab content).
func NewTabAction(deps *Deps) view.View {
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

		templateName := "supplier-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

func buildPageData(supplier *supplierpb.Supplier, id, activeTab string, viewCtx *view.ViewContext, deps *Deps) *PageData {
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
		creditLimit = fmt.Sprintf("%.2f", cl)
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

	tabItems := buildTabItems(id, deps.Routes)

	return &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          displayName,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "suppliers",
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

func buildTabItems(id string, routes entydad.SupplierRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: "Supplier Information", Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
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

package action

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	suppliercategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier_category"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	"github.com/erniealice/entydad-golang"
)

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id   string
	Name string
}

// TagOption represents a tag available for selection in the form.
type TagOption struct {
	Value    string
	Label    string
	Selected bool
}

// SelectedTag represents a pre-selected tag for chip rendering in the multi-select.
type SelectedTag struct {
	Value string
	Label string
}

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	CompanyName        string
	SupplierType       string
	TaxID              string
	RegistrationNumber string
	StreetAddress      string
	City               string
	Province           string
	PostalCode         string
	Country            string
	DefaultCurrency    string
	PaymentTerms       string
	LeadTimeDays       string
	CreditLimit        string
	Status             string
	Website            string
	Notes              string
	FirstName          string
	LastName           string
	Email              string
	Phone              string
	Active             string

	// Section titles
	SectionCompany   string
	SectionContact   string
	SectionFinancial string
	SectionAddress   string

	// Placeholders
	CompanyNamePlaceholder        string
	SupplierTypePlaceholder       string
	StatusPlaceholder             string
	FirstNamePlaceholder          string
	LastNamePlaceholder           string
	EmailPlaceholder              string
	PhonePlaceholder              string
	PaymentTermsPlaceholder       string
	CreditLimitPlaceholder        string
	DefaultCurrencyPlaceholder    string
	LeadTimeDaysPlaceholder       string
	TaxIDPlaceholder              string
	RegistrationNumberPlaceholder string
	StreetAddressPlaceholder      string
	CityPlaceholder               string
	ProvincePlaceholder           string
	PostalCodePlaceholder         string
	CountryPlaceholder            string
	WebsitePlaceholder            string
	NotesPlaceholder              string

	// Select option labels
	TypeCompany    string
	TypeIndividual string

	StatusActive  string
	StatusBlocked string
	StatusOnHold  string

	SelectPaymentTerm string

	Tags                  string
	TagsPlaceholder       string
	TagsSearchPlaceholder string
	TagsNoResults         string
}

// FormData is the template data for the supplier drawer form.
type FormData struct {
	FormAction            string
	IsEdit                bool
	ID                    string
	CompanyName           string
	SupplierType          string
	TaxID                 string
	RegistrationNumber    string
	StreetAddress         string
	City                  string
	Province              string
	PostalCode            string
	Country               string
	DefaultCurrency       string
	PaymentTerms          []*PaymentTermOption
	SelectedPaymentTermID string
	LeadTimeDays          string
	CreditLimit           string
	Status                string
	Website               string
	Notes                 string
	FirstName             string
	LastName              string
	Email                 string
	Phone                 string
	Active                bool
	Labels                FormLabels
	CommonLabels          any
	TagOptions            []TagOption
	SelectedTags          []SelectedTag
}

// Deps holds dependencies for supplier action handlers.
type Deps struct {
	Routes            entydad.SupplierRoutes
	CreateSupplier    func(ctx context.Context, req *supplierpb.CreateSupplierRequest) (*supplierpb.CreateSupplierResponse, error)
	ReadSupplier      func(ctx context.Context, req *supplierpb.ReadSupplierRequest) (*supplierpb.ReadSupplierResponse, error)
	UpdateSupplier    func(ctx context.Context, req *supplierpb.UpdateSupplierRequest) (*supplierpb.UpdateSupplierResponse, error)
	DeleteSupplier    func(ctx context.Context, req *supplierpb.DeleteSupplierRequest) (*supplierpb.DeleteSupplierResponse, error)
	SetSupplierActive func(ctx context.Context, id string, active bool) error
	ListPaymentTerms  func(ctx context.Context) ([]*PaymentTermOption, error)

	// Tag-related deps for multi-select tags on the supplier form
	ListCategories         func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListSupplierCategories func(ctx context.Context, req *suppliercategorypb.ListSupplierCategoriesRequest) (*suppliercategorypb.ListSupplierCategoriesResponse, error)
	CreateSupplierCategory func(ctx context.Context, req *suppliercategorypb.CreateSupplierCategoryRequest) (*suppliercategorypb.CreateSupplierCategoryResponse, error)
	DeleteSupplierCategory func(ctx context.Context, req *suppliercategorypb.DeleteSupplierCategoryRequest) (*suppliercategorypb.DeleteSupplierCategoryResponse, error)
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		CompanyName:        t("supplier.form.companyName"),
		SupplierType:       t("supplier.form.supplierType"),
		TaxID:              t("supplier.form.taxId"),
		RegistrationNumber: t("supplier.form.registrationNumber"),
		StreetAddress:      t("supplier.form.streetAddress"),
		City:               t("supplier.form.city"),
		Province:           t("supplier.form.province"),
		PostalCode:         t("supplier.form.postalCode"),
		Country:            t("supplier.form.country"),
		DefaultCurrency:    t("supplier.form.defaultCurrency"),
		PaymentTerms:       t("supplier.form.paymentTerms"),
		LeadTimeDays:       t("supplier.form.leadTimeDays"),
		CreditLimit:        t("supplier.form.creditLimit"),
		Status:             t("supplier.form.status"),
		Website:            t("supplier.form.website"),
		Notes:              t("supplier.form.notes"),
		FirstName:          t("supplier.form.firstName"),
		LastName:           t("supplier.form.lastName"),
		Email:              t("supplier.form.email"),
		Phone:              t("supplier.form.phone"),
		Active:             t("supplier.form.active"),

		// Section titles
		SectionCompany:   t("supplier.form.sectionCompany"),
		SectionContact:   t("supplier.form.sectionContact"),
		SectionFinancial: t("supplier.form.sectionFinancial"),
		SectionAddress:   t("supplier.form.sectionAddress"),

		// Placeholders
		CompanyNamePlaceholder:        t("supplier.form.companyNamePlaceholder"),
		SupplierTypePlaceholder:       t("supplier.form.supplierTypePlaceholder"),
		StatusPlaceholder:             t("supplier.form.statusPlaceholder"),
		FirstNamePlaceholder:          t("supplier.form.firstNamePlaceholder"),
		LastNamePlaceholder:           t("supplier.form.lastNamePlaceholder"),
		EmailPlaceholder:              t("supplier.form.emailPlaceholder"),
		PhonePlaceholder:              t("supplier.form.phonePlaceholder"),
		PaymentTermsPlaceholder:       t("supplier.form.paymentTermsPlaceholder"),
		CreditLimitPlaceholder:        t("supplier.form.creditLimitPlaceholder"),
		DefaultCurrencyPlaceholder:    t("supplier.form.defaultCurrencyPlaceholder"),
		LeadTimeDaysPlaceholder:       t("supplier.form.leadTimeDaysPlaceholder"),
		TaxIDPlaceholder:              t("supplier.form.taxIdPlaceholder"),
		RegistrationNumberPlaceholder: t("supplier.form.registrationNumberPlaceholder"),
		StreetAddressPlaceholder:      t("supplier.form.streetAddressPlaceholder"),
		CityPlaceholder:               t("supplier.form.cityPlaceholder"),
		ProvincePlaceholder:           t("supplier.form.provincePlaceholder"),
		PostalCodePlaceholder:         t("supplier.form.postalCodePlaceholder"),
		CountryPlaceholder:            t("supplier.form.countryPlaceholder"),
		WebsitePlaceholder:            t("supplier.form.websitePlaceholder"),
		NotesPlaceholder:              t("supplier.form.notesPlaceholder"),

		// Select option labels
		TypeCompany:    t("supplier.form.typeCompany"),
		TypeIndividual: t("supplier.form.typeIndividual"),

		StatusActive:  t("supplier.form.statusActive"),
		StatusBlocked: t("supplier.form.statusBlocked"),
		StatusOnHold:  t("supplier.form.statusOnHold"),

		SelectPaymentTerm: t("supplier.form.selectPaymentTerm"),

		Tags:                  t("supplier.form.tags"),
		TagsPlaceholder:       t("supplier.form.tagsPlaceholder"),
		TagsSearchPlaceholder: t("supplier.form.tagsSearchPlaceholder"),
		TagsNoResults:         t("supplier.form.tagsNoResults"),
	}
}

// optionalString returns a pointer to the string if non-empty, nil otherwise.
func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// optionalInt32 parses a string as int32, returning nil if empty or invalid.
func optionalInt32(s string) *int32 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil
	}
	i := int32(v)
	return &i
}

// optionalInt64Money parses a money string (e.g. "123.45") as int64 centavos, returning nil if empty or invalid.
func optionalInt64Money(s string) *int64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	i := int64(math.Round(v * 100))
	return &i
}

// loadPaymentTerms fetches the payment term options. Returns nil slice on error (graceful degradation).
func loadPaymentTerms(ctx context.Context, deps *Deps) []*PaymentTermOption {
	if deps.ListPaymentTerms == nil {
		return nil
	}
	terms, err := deps.ListPaymentTerms(ctx)
	if err != nil {
		log.Printf("Failed to load payment terms: %v", err)
		return nil
	}
	return terms
}

// loadTagData returns available tag options and the pre-selected tags for the form.
func loadTagData(ctx context.Context, deps *Deps, supplierID string) ([]TagOption, []SelectedTag) {
	if deps.ListCategories == nil {
		return nil, nil
	}

	catResp, err := deps.ListCategories(ctx, &categorypb.ListCategoriesRequest{})
	if err != nil {
		log.Printf("Failed to load supplier tag options: %v", err)
		return nil, nil
	}

	assigned := make(map[string]bool)
	if supplierID != "" && deps.ListSupplierCategories != nil {
		scResp, err := deps.ListSupplierCategories(ctx, &suppliercategorypb.ListSupplierCategoriesRequest{})
		if err != nil {
			log.Printf("Failed to load supplier categories: %v", err)
		} else {
			for _, sc := range scResp.GetData() {
				if sc.GetSupplierId() == supplierID {
					assigned[sc.GetCategoryId()] = true
				}
			}
		}
	}

	var options []TagOption
	var selected []SelectedTag
	for _, cat := range catResp.GetData() {
		if cat.GetModule() != "supplier" || !cat.GetActive() {
			continue
		}
		isAssigned := assigned[cat.GetId()]
		options = append(options, TagOption{
			Value:    cat.GetId(),
			Label:    cat.GetName(),
			Selected: isAssigned,
		})
		if isAssigned {
			selected = append(selected, SelectedTag{
				Value: cat.GetId(),
				Label: cat.GetName(),
			})
		}
	}
	return options, selected
}

// syncTags reconciles the submitted tag IDs with existing supplier_category junction records.
func syncTags(ctx context.Context, deps *Deps, supplierID string, submittedTagIDs []string) {
	if deps.ListSupplierCategories == nil || deps.CreateSupplierCategory == nil || deps.DeleteSupplierCategory == nil {
		return
	}

	scResp, err := deps.ListSupplierCategories(ctx, &suppliercategorypb.ListSupplierCategoriesRequest{})
	if err != nil {
		log.Printf("Failed to list supplier categories for sync: %v", err)
		return
	}

	current := make(map[string]string)
	for _, sc := range scResp.GetData() {
		if sc.GetSupplierId() == supplierID {
			current[sc.GetCategoryId()] = sc.GetId()
		}
	}

	desired := make(map[string]bool)
	for _, tagID := range submittedTagIDs {
		if tagID != "" {
			desired[tagID] = true
		}
	}

	for tagID := range desired {
		if _, exists := current[tagID]; !exists {
			_, err := deps.CreateSupplierCategory(ctx, &suppliercategorypb.CreateSupplierCategoryRequest{
				Data: &suppliercategorypb.SupplierCategory{
					SupplierId: supplierID,
					CategoryId: tagID,
					Active:     true,
				},
			})
			if err != nil {
				log.Printf("Failed to assign tag %s to supplier %s: %v", tagID, supplierID, err)
			}
		}
	}

	for tagID, junctionID := range current {
		if !desired[tagID] {
			_, err := deps.DeleteSupplierCategory(ctx, &suppliercategorypb.DeleteSupplierCategoryRequest{
				Data: &suppliercategorypb.SupplierCategory{Id: junctionID},
			})
			if err != nil {
				log.Printf("Failed to remove tag %s from supplier %s: %v", tagID, supplierID, err)
			}
		}
	}
}

// parseTagIDs splits a comma-separated string of tag IDs from the multi-select
// hidden input into a slice of individual IDs.
func parseTagIDs(csv string) []string {
	if csv == "" {
		return nil
	}
	parts := strings.Split(csv, ",")
	var ids []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			ids = append(ids, p)
		}
	}
	return ids
}

// NewAddAction creates the supplier add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "create") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			tagOptions, _ := loadTagData(ctx, deps, "")
			return view.OK("supplier-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				Status:       "active",
				PaymentTerms: loadPaymentTerms(ctx, deps),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
				TagOptions:   tagOptions,
			})
		}

		// POST -- create supplier
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		resp, err := deps.CreateSupplier(ctx, &supplierpb.CreateSupplierRequest{
			Data: &supplierpb.Supplier{
				Active:             active,
				CompanyName:        r.FormValue("company_name"),
				SupplierType:       r.FormValue("supplier_type"),
				TaxId:              optionalString(r.FormValue("tax_id")),
				RegistrationNumber: optionalString(r.FormValue("registration_number")),
				StreetAddress:      optionalString(r.FormValue("street_address")),
				City:               optionalString(r.FormValue("city")),
				Province:           optionalString(r.FormValue("province")),
				PostalCode:         optionalString(r.FormValue("postal_code")),
				Country:            optionalString(r.FormValue("country")),
				DefaultCurrency:    optionalString(r.FormValue("default_currency")),
				PaymentTermId:      optionalString(r.FormValue("payment_term_id")),
				LeadTimeDays:       optionalInt32(r.FormValue("lead_time_days")),
				CreditLimit:        optionalInt64Money(r.FormValue("credit_limit")),
				Status:             optionalString(r.FormValue("status")),
				Website:            optionalString(r.FormValue("website")),
				Notes:              optionalString(r.FormValue("notes")),
				User: &userpb.User{
					FirstName:    r.FormValue("first_name"),
					LastName:     r.FormValue("last_name"),
					EmailAddress: r.FormValue("email_address"),
					MobileNumber: r.FormValue("mobile_number"),
					Active:       active,
				},
			},
		})
		if err != nil {
			log.Printf("Failed to create supplier: %v", err)
			return entydad.HTMXError(err.Error())
		}

		if resp != nil && len(resp.GetData()) > 0 {
			syncTags(ctx, deps, resp.GetData()[0].GetId(), parseTagIDs(r.FormValue("tags")))
		}

		return entydad.HTMXSuccess("suppliers-table")
	})
}

// NewEditAction creates the supplier edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSupplier(ctx, &supplierpb.ReadSupplierRequest{
				Data: &supplierpb.Supplier{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read supplier %s: %v", id, err)
				return entydad.HTMXError(viewCtx.T("shared.errors.notFound"))
			}

			s := resp.GetData()[0]
			u := s.GetUser()

			firstName := ""
			lastName := ""
			email := ""
			phone := ""
			if u != nil {
				firstName = u.GetFirstName()
				lastName = u.GetLastName()
				email = u.GetEmailAddress()
				phone = u.GetMobileNumber()
			}

			status := s.GetStatus()
			if status == "" {
				if s.GetActive() {
					status = "active"
				} else {
					status = "blocked"
				}
			}

			leadTimeDays := ""
			if ltd := s.GetLeadTimeDays(); ltd > 0 {
				leadTimeDays = strconv.FormatInt(int64(ltd), 10)
			}
			creditLimit := ""
			if cl := s.GetCreditLimit(); cl > 0 {
				creditLimit = strconv.FormatFloat(float64(cl)/100.0, 'f', 2, 64)
			}

			tagOptions, selectedTags := loadTagData(ctx, deps, id)
			return view.OK("supplier-drawer-form", &FormData{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				ID:                    id,
				CompanyName:           s.GetCompanyName(),
				SupplierType:          s.GetSupplierType(),
				TaxID:                 s.GetTaxId(),
				RegistrationNumber:    s.GetRegistrationNumber(),
				StreetAddress:         s.GetStreetAddress(),
				City:                  s.GetCity(),
				Province:              s.GetProvince(),
				PostalCode:            s.GetPostalCode(),
				Country:               s.GetCountry(),
				DefaultCurrency:       s.GetDefaultCurrency(),
				PaymentTerms:          loadPaymentTerms(ctx, deps),
				SelectedPaymentTermID: s.GetPaymentTermId(),
				LeadTimeDays:          leadTimeDays,
				CreditLimit:           creditLimit,
				Status:                status,
				Website:               s.GetWebsite(),
				Notes:                 s.GetNotes(),
				FirstName:             firstName,
				LastName:              lastName,
				Email:                 email,
				Phone:                 phone,
				Active:                s.GetActive(),
				Labels:                formLabels(viewCtx.T),
				CommonLabels:          nil, // injected by ViewAdapter
				TagOptions:            tagOptions,
				SelectedTags:          selectedTags,
			})
		}

		// POST -- update supplier
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.UpdateSupplier(ctx, &supplierpb.UpdateSupplierRequest{
			Data: &supplierpb.Supplier{
				Id:                 id,
				Active:             active,
				CompanyName:        r.FormValue("company_name"),
				SupplierType:       r.FormValue("supplier_type"),
				TaxId:              optionalString(r.FormValue("tax_id")),
				RegistrationNumber: optionalString(r.FormValue("registration_number")),
				StreetAddress:      optionalString(r.FormValue("street_address")),
				City:               optionalString(r.FormValue("city")),
				Province:           optionalString(r.FormValue("province")),
				PostalCode:         optionalString(r.FormValue("postal_code")),
				Country:            optionalString(r.FormValue("country")),
				DefaultCurrency:    optionalString(r.FormValue("default_currency")),
				PaymentTermId:      optionalString(r.FormValue("payment_term_id")),
				LeadTimeDays:       optionalInt32(r.FormValue("lead_time_days")),
				CreditLimit:        optionalInt64Money(r.FormValue("credit_limit")),
				Status:             optionalString(r.FormValue("status")),
				Website:            optionalString(r.FormValue("website")),
				Notes:              optionalString(r.FormValue("notes")),
				User: &userpb.User{
					FirstName:    r.FormValue("first_name"),
					LastName:     r.FormValue("last_name"),
					EmailAddress: r.FormValue("email_address"),
					MobileNumber: r.FormValue("mobile_number"),
					Active:       active,
				},
			},
		})
		if err != nil {
			log.Printf("Failed to update supplier %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		syncTags(ctx, deps, id, parseTagIDs(r.FormValue("tags")))

		return entydad.HTMXSuccess("suppliers-table")
	})
}

// NewDeleteAction creates the supplier delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "delete") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}

		_, err := deps.DeleteSupplier(ctx, &supplierpb.DeleteSupplierRequest{
			Data: &supplierpb.Supplier{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete supplier %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("suppliers-table")
	})
}

// NewBulkDeleteAction creates the supplier bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "delete") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}

		for _, id := range ids {
			_, err := deps.DeleteSupplier(ctx, &supplierpb.DeleteSupplierRequest{
				Data: &supplierpb.Supplier{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete supplier %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("suppliers-table")
	})
}

// NewSetStatusAction creates the supplier set-status action (POST only).
// Expects query params: ?id={supplierId}&status={active|blocked|on_hold}
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return entydad.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}
		if targetStatus != "active" && targetStatus != "blocked" && targetStatus != "on_hold" {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidStatus"))
		}

		if deps.SetSupplierActive != nil {
			if err := deps.SetSupplierActive(ctx, id, targetStatus == "active"); err != nil {
				log.Printf("Failed to update supplier status %s: %v", id, err)
				return entydad.HTMXError(err.Error())
			}
		}

		return entydad.HTMXSuccess("suppliers-table")
	})
}

// NewBulkSetStatusAction creates the supplier bulk set-status action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return entydad.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}
		if targetStatus != "active" && targetStatus != "blocked" && targetStatus != "on_hold" {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidTargetStatus"))
		}

		active := targetStatus == "active"

		if deps.SetSupplierActive != nil {
			for _, id := range ids {
				if err := deps.SetSupplierActive(ctx, id, active); err != nil {
					log.Printf("Failed to update supplier status %s: %v", id, err)
				}
			}
		}

		return entydad.HTMXSuccess("suppliers-table")
	})
}

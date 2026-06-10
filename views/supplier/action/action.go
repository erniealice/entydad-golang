package action

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	suppliercategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier_category"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	"github.com/erniealice/entydad-golang"
	supplierform "github.com/erniealice/entydad-golang/views/supplier/form"
)

// PaymentTermOption is a type alias so callers wired through module.go
// retain the same surface: form.PaymentTermOption is the source of truth.
type PaymentTermOption = supplierform.PaymentTermOption

// Deps holds dependencies for supplier action handlers.
type Deps struct {
	Routes entydad.SupplierRoutes
	// SearchTimezonesURL is the URL of the timezone autocomplete JSON endpoint.
	SearchTimezonesURL string
	CreateSupplier     func(ctx context.Context, req *supplierpb.CreateSupplierRequest) (*supplierpb.CreateSupplierResponse, error)
	ReadSupplier       func(ctx context.Context, req *supplierpb.ReadSupplierRequest) (*supplierpb.ReadSupplierResponse, error)
	UpdateSupplier     func(ctx context.Context, req *supplierpb.UpdateSupplierRequest) (*supplierpb.UpdateSupplierResponse, error)
	DeleteSupplier     func(ctx context.Context, req *supplierpb.DeleteSupplierRequest) (*supplierpb.DeleteSupplierResponse, error)
	SetSupplierStatus  func(ctx context.Context, id string, status string) error
	ListPaymentTerms   func(ctx context.Context) ([]*PaymentTermOption, error)

	// Tag-related deps for multi-select tags on the supplier form
	ListCategories         func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListSupplierCategories func(ctx context.Context, req *suppliercategorypb.ListSupplierCategoriesRequest) (*suppliercategorypb.ListSupplierCategoriesResponse, error)
	CreateSupplierCategory func(ctx context.Context, req *suppliercategorypb.CreateSupplierCategoryRequest) (*suppliercategorypb.CreateSupplierCategoryResponse, error)
	DeleteSupplierCategory func(ctx context.Context, req *suppliercategorypb.DeleteSupplierCategoryRequest) (*suppliercategorypb.DeleteSupplierCategoryResponse, error)
	// CurrencyOptions is the pre-built list of currency select options sourced
	// from lyngua's CommonLabels.Currency.Options. Populated by module.go from
	// ModuleDeps.CommonLabels so the action handler can call
	// supplierform.BuildCurrencyOptions without importing pyeza.CommonLabels.
	CurrencyOptions []pyezatypes.SelectOption
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
func loadPaymentTerms(ctx context.Context, deps *Deps) []*supplierform.PaymentTermOption {
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
func loadTagData(ctx context.Context, deps *Deps, supplierID string) ([]supplierform.TagOption, []supplierform.SelectedTag) {
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

	var options []supplierform.TagOption
	var selected []supplierform.SelectedTag
	for _, cat := range catResp.GetData() {
		if cat.GetModule() != "supplier" || !cat.GetActive() {
			continue
		}
		isAssigned := assigned[cat.GetId()]
		options = append(options, supplierform.TagOption{
			Value:    cat.GetId(),
			Label:    cat.GetName(),
			Selected: isAssigned,
		})
		if isAssigned {
			selected = append(selected, supplierform.SelectedTag{
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
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			tagOptions, _ := loadTagData(ctx, deps, "")
			paymentTerms := loadPaymentTerms(ctx, deps)
			labels := supplierform.BuildLabels(viewCtx.T)
			return view.OK("supplier-drawer-form", &supplierform.Data{
				FormAction:               deps.Routes.AddURL,
				Active:                   true,
				Status:                   "active",
				PaymentTerms:             paymentTerms,
				PaymentTermSelectOptions: supplierform.BuildPaymentTermSelectOptions(paymentTerms, ""),
				StatusOptions:            supplierform.BuildStatusOptions("active", labels),
				SupplierTypeOptions:      supplierform.BuildSupplierTypeOptions("", labels),
				BillingCurrencyOptions:   supplierform.BuildCurrencyOptions("", deps.CurrencyOptions),
				Labels:                   labels,
				CommonLabels:             nil, // injected by ViewAdapter
				TagOptions:               tagOptions,
				SearchTimezonesURL:       deps.SearchTimezonesURL,
			})
		}

		// POST -- create supplier
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		// New suppliers default to active=true; the form no longer exposes the
		// active flag (it's derived from status). The use case keeps active
		// in sync with status on subsequent updates.
		resp, err := deps.CreateSupplier(ctx, &supplierpb.CreateSupplierRequest{
			Data: &supplierpb.Supplier{
				Active:             true,
				Name:               r.FormValue("name"),
				SupplierType:       r.FormValue("supplier_type"),
				TaxId:              optionalString(r.FormValue("tax_id")),
				RegistrationNumber: optionalString(r.FormValue("registration_number")),
				StreetAddress:      optionalString(r.FormValue("street_address")),
				City:               optionalString(r.FormValue("city")),
				Province:           optionalString(r.FormValue("province")),
				PostalCode:         optionalString(r.FormValue("postal_code")),
				Country:            optionalString(r.FormValue("country")),
				BillingCurrency:    optionalString(r.FormValue("billing_currency")),
				PaymentTermId:      optionalString(r.FormValue("payment_term_id")),
				LeadTimeDays:       optionalInt32(r.FormValue("lead_time_days")),
				CreditLimit:        optionalInt64Money(r.FormValue("credit_limit")),
				Status:             optionalString(r.FormValue("status")),
				Website:            optionalString(r.FormValue("website")),
				Notes:              optionalString(r.FormValue("notes")),
				Timezone:           optionalString(r.FormValue("timezone")),
				User: &userpb.User{
					FirstName:    r.FormValue("first_name"),
					LastName:     r.FormValue("last_name"),
					EmailAddress: r.FormValue("email_address"),
					MobileNumber: r.FormValue("mobile_number"),
					Timezone:     optionalString(r.FormValue("timezone")),
					Active:       true,
				},
			},
		})
		if err != nil {
			log.Printf("Failed to create supplier: %v", err)
			return view.HTMXError(err.Error())
		}

		if resp != nil && len(resp.GetData()) > 0 {
			syncTags(ctx, deps, resp.GetData()[0].GetId(), parseTagIDs(r.FormValue("tags")))
		}

		return view.HTMXSuccess("suppliers-table")
	})
}

// NewEditAction creates the supplier edit action (GET = form, POST = update).
// When the GET request includes ?clone=1, the handler returns the drawer form
// pre-populated from the source record but wired to AddURL (submission creates
// a new supplier) with " (Copy)" appended to the name.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("supplier", requiredAction) {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSupplier(ctx, &supplierpb.ReadSupplierRequest{
				Data: &supplierpb.Supplier{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read supplier %s: %v", id, err)
				return view.HTMXError(viewCtx.T("shared.errors.notFound"))
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

			name := s.GetName()
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id)
			formID := id
			if isClone {
				name = strings.TrimSpace(name) + viewCtx.T("actions.copySuffix")
				formAction = deps.Routes.AddURL
				formID = ""
			}

			timezone := ""
			if u != nil {
				timezone = u.GetTimezone()
			}

			tagOptions, selectedTags := loadTagData(ctx, deps, id)
			paymentTerms := loadPaymentTerms(ctx, deps)
			selectedPaymentTermID := s.GetPaymentTermId()
			labels := supplierform.BuildLabels(viewCtx.T)
			return view.OK("supplier-drawer-form", &supplierform.Data{
				FormAction:               formAction,
				IsEdit:                   !isClone,
				ID:                       formID,
				Name:                     name,
				Timezone:                 timezone,
				SearchTimezonesURL:       deps.SearchTimezonesURL,
				SupplierType:             s.GetSupplierType(),
				TaxID:                    s.GetTaxId(),
				RegistrationNumber:       s.GetRegistrationNumber(),
				StreetAddress:            s.GetStreetAddress(),
				City:                     s.GetCity(),
				Province:                 s.GetProvince(),
				PostalCode:               s.GetPostalCode(),
				Country:                  s.GetCountry(),
				BillingCurrency:          s.GetBillingCurrency(),
				PaymentTerms:             paymentTerms,
				SelectedPaymentTermID:    selectedPaymentTermID,
				PaymentTermSelectOptions: supplierform.BuildPaymentTermSelectOptions(paymentTerms, selectedPaymentTermID),
				StatusOptions:            supplierform.BuildStatusOptions(status, labels),
				SupplierTypeOptions:      supplierform.BuildSupplierTypeOptions(s.GetSupplierType(), labels),
				BillingCurrencyOptions:   supplierform.BuildCurrencyOptions(s.GetBillingCurrency(), deps.CurrencyOptions),
				LeadTimeDays:             leadTimeDays,
				CreditLimit:              creditLimit,
				Status:                   status,
				Website:                  s.GetWebsite(),
				Notes:                    s.GetNotes(),
				FirstName:                firstName,
				LastName:                 lastName,
				Email:                    email,
				Phone:                    phone,
				Active:                   s.GetActive(),
				Labels:                   labels,
				CommonLabels:             nil, // injected by ViewAdapter
				TagOptions:               tagOptions,
				SelectedTags:             selectedTags,
			})
		}

		// POST -- update supplier
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		// The drawer form no longer exposes an "active" toggle. Active is
		// derived from status server-side; the Update Supplier use case copies
		// the value from the existing record so a missing field never deactivates
		// the supplier.
		_, err := deps.UpdateSupplier(ctx, &supplierpb.UpdateSupplierRequest{
			Data: &supplierpb.Supplier{
				Id:                 id,
				Name:               r.FormValue("name"),
				SupplierType:       r.FormValue("supplier_type"),
				TaxId:              optionalString(r.FormValue("tax_id")),
				RegistrationNumber: optionalString(r.FormValue("registration_number")),
				StreetAddress:      optionalString(r.FormValue("street_address")),
				City:               optionalString(r.FormValue("city")),
				Province:           optionalString(r.FormValue("province")),
				PostalCode:         optionalString(r.FormValue("postal_code")),
				Country:            optionalString(r.FormValue("country")),
				BillingCurrency:    optionalString(r.FormValue("billing_currency")),
				PaymentTermId:      optionalString(r.FormValue("payment_term_id")),
				LeadTimeDays:       optionalInt32(r.FormValue("lead_time_days")),
				CreditLimit:        optionalInt64Money(r.FormValue("credit_limit")),
				Status:             optionalString(r.FormValue("status")),
				Website:            optionalString(r.FormValue("website")),
				Notes:              optionalString(r.FormValue("notes")),
				Timezone:           optionalString(r.FormValue("timezone")),
				User: &userpb.User{
					FirstName:    r.FormValue("first_name"),
					LastName:     r.FormValue("last_name"),
					EmailAddress: r.FormValue("email_address"),
					MobileNumber: r.FormValue("mobile_number"),
					Timezone:     optionalString(r.FormValue("timezone")),
				},
			},
		})
		if err != nil {
			log.Printf("Failed to update supplier %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		syncTags(ctx, deps, id, parseTagIDs(r.FormValue("tags")))

		return view.HTMXSuccess("suppliers-table")
	})
}

// NewDeleteAction creates the supplier delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}

		_, err := deps.DeleteSupplier(ctx, &supplierpb.DeleteSupplierRequest{
			Data: &supplierpb.Supplier{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete supplier %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("suppliers-table")
	})
}

// NewBulkDeleteAction creates the supplier bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}

		for _, id := range ids {
			_, err := deps.DeleteSupplier(ctx, &supplierpb.DeleteSupplierRequest{
				Data: &supplierpb.Supplier{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete supplier %s: %v", id, err)
			}
		}

		return view.HTMXSuccess("suppliers-table")
	})
}

// NewSetStatusAction creates the supplier set-status action (POST only).
// Expects query params: ?id={supplierId}&status={active|blocked|on_hold}
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return view.HTMXError(viewCtx.T("shared.errors.idRequired"))
		}
		if targetStatus != "active" && targetStatus != "blocked" && targetStatus != "on_hold" {
			return view.HTMXError(viewCtx.T("shared.errors.invalidStatus"))
		}

		if deps.SetSupplierStatus != nil {
			if err := deps.SetSupplierStatus(ctx, id, targetStatus); err != nil {
				log.Printf("Failed to update supplier status %s: %v", id, err)
				return view.HTMXError(err.Error())
			}
		}

		return view.HTMXSuccess("suppliers-table")
	})
}

// NewBulkSetStatusAction creates the supplier bulk set-status action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}
		if targetStatus != "active" && targetStatus != "blocked" && targetStatus != "on_hold" {
			return view.HTMXError(viewCtx.T("shared.errors.invalidTargetStatus"))
		}

		if deps.SetSupplierStatus != nil {
			for _, id := range ids {
				if err := deps.SetSupplierStatus(ctx, id, targetStatus); err != nil {
					log.Printf("Failed to update supplier status %s: %v", id, err)
				}
			}
		}

		return view.HTMXSuccess("suppliers-table")
	})
}

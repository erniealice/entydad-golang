package action

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	"github.com/erniealice/entydad-golang"
)

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
}

// FormData is the template data for the supplier drawer form.
type FormData struct {
	FormAction         string
	IsEdit             bool
	ID                 string
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
	Active             bool
	Labels             FormLabels
	CommonLabels       any
}

// Deps holds dependencies for supplier action handlers.
type Deps struct {
	Routes            entydad.SupplierRoutes
	CreateSupplier    func(ctx context.Context, req *supplierpb.CreateSupplierRequest) (*supplierpb.CreateSupplierResponse, error)
	ReadSupplier      func(ctx context.Context, req *supplierpb.ReadSupplierRequest) (*supplierpb.ReadSupplierResponse, error)
	UpdateSupplier    func(ctx context.Context, req *supplierpb.UpdateSupplierRequest) (*supplierpb.UpdateSupplierResponse, error)
	DeleteSupplier    func(ctx context.Context, req *supplierpb.DeleteSupplierRequest) (*supplierpb.DeleteSupplierResponse, error)
	SetSupplierActive func(ctx context.Context, id string, active bool) error
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

// optionalFloat64 parses a string as float64, returning nil if empty or invalid.
func optionalFloat64(s string) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

// NewAddAction creates the supplier add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier", "create") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("supplier-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				Status:       "active",
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- create supplier
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.CreateSupplier(ctx, &supplierpb.CreateSupplierRequest{
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
				PaymentTerms:       optionalString(r.FormValue("payment_terms")),
				LeadTimeDays:       optionalInt32(r.FormValue("lead_time_days")),
				CreditLimit:        optionalFloat64(r.FormValue("credit_limit")),
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
				creditLimit = strconv.FormatFloat(cl, 'f', 2, 64)
			}

			return view.OK("supplier-drawer-form", &FormData{
				FormAction:         route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:             true,
				ID:                 id,
				CompanyName:        s.GetCompanyName(),
				SupplierType:       s.GetSupplierType(),
				TaxID:              s.GetTaxId(),
				RegistrationNumber: s.GetRegistrationNumber(),
				StreetAddress:      s.GetStreetAddress(),
				City:               s.GetCity(),
				Province:           s.GetProvince(),
				PostalCode:         s.GetPostalCode(),
				Country:            s.GetCountry(),
				DefaultCurrency:    s.GetDefaultCurrency(),
				PaymentTerms:       s.GetPaymentTerms(),
				LeadTimeDays:       leadTimeDays,
				CreditLimit:        creditLimit,
				Status:             status,
				Website:            s.GetWebsite(),
				Notes:              s.GetNotes(),
				FirstName:          firstName,
				LastName:           lastName,
				Email:              email,
				Phone:              phone,
				Active:             s.GetActive(),
				Labels:             formLabels(viewCtx.T),
				CommonLabels:       nil, // injected by ViewAdapter
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
				PaymentTerms:       optionalString(r.FormValue("payment_terms")),
				LeadTimeDays:       optionalInt32(r.FormValue("lead_time_days")),
				CreditLimit:        optionalFloat64(r.FormValue("credit_limit")),
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

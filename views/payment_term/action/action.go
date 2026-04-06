package action

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"

	"github.com/erniealice/entydad-golang"
)

// FormLabels holds i18n labels for the payment term drawer form.
type FormLabels struct {
	SectionInfo            string
	SectionTerms           string
	SectionSettings        string
	Name                   string
	NamePlaceholder        string
	Code                   string
	CodePlaceholder        string
	Type                   string
	NetDays                string
	DiscountDays           string
	DiscountPercentBps     string
	TypeHint               string
	NetDaysHint            string
	DiscountDaysHint       string
	DiscountPercentBpsHint string
	PriorityHint           string
	EntityScope            string
	IsDefault              string
	Description            string
	DescriptionPlaceholder string
	DisplayOrder           string
	Active                 string

	// Select option labels — Type
	TypeDueOnReceipt        string
	TypeNet                 string
	TypeCOD                 string
	TypeProximate           string
	ProximateDay            string
	ProximateDayPlaceholder string

	// Select option labels — EntityScope
	ScopesBoth         string
	ScopesSupplierOnly string
	ScopesClientOnly   string
}

// FormData is the template data for the payment term drawer form.
type FormData struct {
	FormAction         string
	IsEdit             bool
	ID                 string
	Name               string
	Code               string
	Type               string
	NetDays            string
	DiscountDays       string
	DiscountPercentBps string
	EntityScope        string
	IsDefault          bool
	Description        string
	DisplayOrder       string
	ProximateDay       string
	Active             bool
	Labels             FormLabels
	CommonLabels       any
}

// Deps holds dependencies for payment term action handlers.
type Deps struct {
	Routes               entydad.PaymentTermRoutes
	CreatePaymentTerm    func(ctx context.Context, req *paymenttermpb.CreatePaymentTermRequest) (*paymenttermpb.CreatePaymentTermResponse, error)
	ReadPaymentTerm      func(ctx context.Context, req *paymenttermpb.ReadPaymentTermRequest) (*paymenttermpb.ReadPaymentTermResponse, error)
	UpdatePaymentTerm    func(ctx context.Context, req *paymenttermpb.UpdatePaymentTermRequest) (*paymenttermpb.UpdatePaymentTermResponse, error)
	DeletePaymentTerm    func(ctx context.Context, req *paymenttermpb.DeletePaymentTermRequest) (*paymenttermpb.DeletePaymentTermResponse, error)
	SetPaymentTermActive func(ctx context.Context, id string, active bool) error
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		SectionInfo:            t("paymentTerm.form.sectionInfo"),
		SectionTerms:           t("paymentTerm.form.sectionTerms"),
		SectionSettings:        t("paymentTerm.form.sectionSettings"),
		Name:                   t("paymentTerm.form.name"),
		NamePlaceholder:        t("paymentTerm.form.namePlaceholder"),
		Code:                   t("paymentTerm.form.code"),
		CodePlaceholder:        t("paymentTerm.form.codePlaceholder"),
		Type:                   t("paymentTerm.form.type"),
		NetDays:                t("paymentTerm.form.netDays"),
		DiscountDays:           t("paymentTerm.form.discountDays"),
		DiscountPercentBps:     t("paymentTerm.form.discountPercentBps"),
		TypeHint:               t("paymentTerm.form.typeHint"),
		NetDaysHint:            t("paymentTerm.form.netDaysHint"),
		DiscountDaysHint:       t("paymentTerm.form.discountDaysHint"),
		DiscountPercentBpsHint: t("paymentTerm.form.discountPercentBpsHint"),
		PriorityHint:           t("paymentTerm.form.priorityHint"),
		EntityScope:            t("paymentTerm.form.entityScope"),
		IsDefault:              t("paymentTerm.form.isDefault"),
		Description:            t("paymentTerm.form.description"),
		DescriptionPlaceholder: t("paymentTerm.form.descriptionPlaceholder"),
		DisplayOrder:           t("paymentTerm.form.displayOrder"),
		Active:                 t("paymentTerm.form.active"),

		TypeDueOnReceipt:        t("paymentTerm.form.typeDueOnReceipt"),
		TypeNet:                 t("paymentTerm.form.typeNet"),
		TypeCOD:                 t("paymentTerm.form.typeCOD"),
		TypeProximate:           t("paymentTerm.form.typeProximate"),
		ProximateDay:            t("paymentTerm.form.proximateDay"),
		ProximateDayPlaceholder: t("paymentTerm.form.proximateDayPlaceholder"),

		ScopesBoth:         t("paymentTerm.form.scopesBoth"),
		ScopesSupplierOnly: t("paymentTerm.form.scopesSupplierOnly"),
		ScopesClientOnly:   t("paymentTerm.form.scopesClientOnly"),
	}
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

// optionalString returns a pointer to the string if non-empty, nil otherwise.
func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// requiredInt32 parses a string as int32, returning 0 if empty or invalid.
func requiredInt32(s string) int32 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(v)
}

// NewAddAction creates the payment term add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("payment_term", "create") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("payment-term-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				EntityScope:  "both",
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil,
			})
		}

		// POST — create payment term
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		isDefault := r.FormValue("is_default") == "true"

		_, err := deps.CreatePaymentTerm(ctx, &paymenttermpb.CreatePaymentTermRequest{
			Data: &paymenttermpb.PaymentTerm{
				Active:             active,
				Name:               r.FormValue("name"),
				Code:               r.FormValue("code"),
				Type:               r.FormValue("type"),
				NetDays:            requiredInt32(r.FormValue("net_days")),
				DiscountDays:       optionalInt32(r.FormValue("discount_days")),
				DiscountPercentBps: optionalInt32(r.FormValue("discount_percent_bps")),
				EntityScope:        r.FormValue("entity_scope"),
				IsDefault:          isDefault,
				Description:        optionalString(r.FormValue("description")),
				DisplayOrder:       optionalInt32(r.FormValue("display_order")),
				ProximateDay:       optionalInt32(r.FormValue("proximate_day")),
			},
		})
		if err != nil {
			log.Printf("Failed to create payment term: %v", err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("payment-terms-table")
	})
}

// NewEditAction creates the payment term edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("payment_term", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPaymentTerm(ctx, &paymenttermpb.ReadPaymentTermRequest{
				Data: &paymenttermpb.PaymentTerm{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read payment term %s: %v", id, err)
				return entydad.HTMXError(viewCtx.T("shared.errors.notFound"))
			}

			data := resp.GetData()
			if len(data) == 0 {
				return entydad.HTMXError(viewCtx.T("shared.errors.notFound"))
			}
			pt := data[0]

			netDays := ""
			if v := pt.GetNetDays(); v > 0 {
				netDays = strconv.FormatInt(int64(v), 10)
			}
			discountDays := ""
			if v := pt.GetDiscountDays(); v > 0 {
				discountDays = strconv.FormatInt(int64(v), 10)
			}
			discountPercentBps := ""
			if v := pt.GetDiscountPercentBps(); v > 0 {
				discountPercentBps = strconv.FormatInt(int64(v), 10)
			}
			displayOrder := ""
			if v := pt.GetDisplayOrder(); v > 0 {
				displayOrder = strconv.FormatInt(int64(v), 10)
			}
			proximateDay := ""
			if v := pt.GetProximateDay(); v > 0 {
				proximateDay = strconv.FormatInt(int64(v), 10)
			}

			return view.OK("payment-term-drawer-form", &FormData{
				FormAction:         route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:             true,
				ID:                 id,
				Name:               pt.GetName(),
				Code:               pt.GetCode(),
				Type:               pt.GetType(),
				NetDays:            netDays,
				DiscountDays:       discountDays,
				DiscountPercentBps: discountPercentBps,
				EntityScope:        pt.GetEntityScope(),
				IsDefault:          pt.GetIsDefault(),
				Description:        pt.GetDescription(),
				DisplayOrder:       displayOrder,
				ProximateDay:       proximateDay,
				Active:             pt.GetActive(),
				Labels:             formLabels(viewCtx.T),
				CommonLabels:       nil,
			})
		}

		// POST — update payment term
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		isDefault := r.FormValue("is_default") == "true"

		_, err := deps.UpdatePaymentTerm(ctx, &paymenttermpb.UpdatePaymentTermRequest{
			Data: &paymenttermpb.PaymentTerm{
				Id:                 id,
				Active:             active,
				Name:               r.FormValue("name"),
				Code:               r.FormValue("code"),
				Type:               r.FormValue("type"),
				NetDays:            requiredInt32(r.FormValue("net_days")),
				DiscountDays:       optionalInt32(r.FormValue("discount_days")),
				DiscountPercentBps: optionalInt32(r.FormValue("discount_percent_bps")),
				EntityScope:        r.FormValue("entity_scope"),
				IsDefault:          isDefault,
				Description:        optionalString(r.FormValue("description")),
				DisplayOrder:       optionalInt32(r.FormValue("display_order")),
				ProximateDay:       optionalInt32(r.FormValue("proximate_day")),
			},
		})
		if err != nil {
			log.Printf("Failed to update payment term %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("payment-terms-table")
	})
}

// NewDeleteAction creates the payment term delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("payment_term", "delete") {
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

		_, err := deps.DeletePaymentTerm(ctx, &paymenttermpb.DeletePaymentTermRequest{
			Data: &paymenttermpb.PaymentTerm{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete payment term %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("payment-terms-table")
	})
}

// NewBulkDeleteAction creates the payment term bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("payment_term", "delete") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}

		for _, id := range ids {
			_, err := deps.DeletePaymentTerm(ctx, &paymenttermpb.DeletePaymentTermRequest{
				Data: &paymenttermpb.PaymentTerm{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete payment term %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("payment-terms-table")
	})
}

// NewSetStatusAction creates the payment term activate/deactivate action (POST only).
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("payment_term", "update") {
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
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidStatus"))
		}

		if err := deps.SetPaymentTermActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update payment term status %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("payment-terms-table")
	})
}

// NewBulkSetStatusAction creates the payment term bulk activate/deactivate action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("payment_term", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return entydad.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidTargetStatus"))
		}

		active := targetStatus == "active"
		for _, id := range ids {
			if err := deps.SetPaymentTermActive(ctx, id, active); err != nil {
				log.Printf("Failed to update payment term status %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("payment-terms-table")
	})
}

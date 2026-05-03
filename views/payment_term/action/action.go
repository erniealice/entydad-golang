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
	paymenttermform "github.com/erniealice/entydad-golang/views/payment_term/form"
)

// Deps holds dependencies for payment term action handlers.
type Deps struct {
	Routes               entydad.PaymentTermRoutes
	CreatePaymentTerm    func(ctx context.Context, req *paymenttermpb.CreatePaymentTermRequest) (*paymenttermpb.CreatePaymentTermResponse, error)
	ReadPaymentTerm      func(ctx context.Context, req *paymenttermpb.ReadPaymentTermRequest) (*paymenttermpb.ReadPaymentTermResponse, error)
	UpdatePaymentTerm    func(ctx context.Context, req *paymenttermpb.UpdatePaymentTermRequest) (*paymenttermpb.UpdatePaymentTermResponse, error)
	DeletePaymentTerm    func(ctx context.Context, req *paymenttermpb.DeletePaymentTermRequest) (*paymenttermpb.DeletePaymentTermResponse, error)
	SetPaymentTermActive func(ctx context.Context, id string, active bool) error
	// Scope is the route-derived entity scope context: "client", "supplier", or ""
	// (empty = standalone, shows all scopes). Used to pre-fill entity_scope on Add.
	Scope string
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

// validateAndNilOutTypeFields enforces the type-driven requiredness + nil-out matrix:
//
//	type            net_days   proximate_day
//	net             A (req)    N (nil-out)
//	due_on_receipt  N (nil-out) N (nil-out)
//	cod             N (nil-out) N (nil-out)
//	proximate       N (nil-out) A (req)
//
// A = required; N = nil-out (caller must not pass value to proto, regardless of client).
// Returns (netDays, proximateDay, errorMsg). errorMsg is non-empty on validation failure.
// Mirrors the JS `rules` object in the template tail script.
func validateAndNilOutTypeFields(termType string, r *http.Request, t func(string) string) (netDays int32, proximateDay *int32, errMsg string) {
	switch termType {
	case "net":
		// net_days is required for Net type.
		rawNetDays := r.FormValue("net_days")
		if rawNetDays == "" {
			return 0, nil, t("paymentTerm.form.errors.netDaysRequired")
		}
		return requiredInt32(rawNetDays), nil, ""
	case "due_on_receipt", "cod":
		// Both net_days and proximate_day are N — nil-out regardless of client value.
		return 0, nil, ""
	case "proximate":
		// proximate_day is required for Proximate type.
		rawProximateDay := r.FormValue("proximate_day")
		if rawProximateDay == "" {
			return 0, nil, t("paymentTerm.form.errors.proximateDayRequired")
		}
		return 0, optionalInt32(rawProximateDay), ""
	default:
		// Unknown or empty type.
		if termType == "" {
			return 0, nil, t("paymentTerm.form.errors.typeRequired")
		}
		return 0, nil, t("paymentTerm.form.errors.typeInvalid")
	}
}

// NewAddAction creates the payment term add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("payment_term", "create") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			// Default entity_scope from the route context so client/supplier
			// settings pages pre-fill the right scope without operator input.
			entityScope := deps.Scope
			if entityScope == "" {
				entityScope = "both"
			}
			labels := paymenttermform.BuildLabels(viewCtx.T)
			// Default type to "net" for Add so the form starts in a coherent state
			// (net_days visible+required, proximate_day hidden).
			defaultType := "net"
			return view.OK("payment-term-drawer-form", &paymenttermform.Data{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				EntityScope:  entityScope,
				Type:         defaultType,
				TypeOptions:  paymenttermform.BuildTypeOptions(labels, defaultType),
				Labels:       labels,
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

		// Validate and nil-out type-driven fields.
		// This mirrors the JS rules object in the template tail script.
		termType := r.FormValue("type")
		netDays, proximateDay, validErr := validateAndNilOutTypeFields(termType, r, viewCtx.T)
		if validErr != "" {
			return entydad.HTMXError(validErr)
		}

		_, err := deps.CreatePaymentTerm(ctx, &paymenttermpb.CreatePaymentTermRequest{
			Data: &paymenttermpb.PaymentTerm{
				Active:             active,
				Name:               r.FormValue("name"),
				Code:               r.FormValue("code"),
				Type:               termType,
				NetDays:            netDays,
				DiscountDays:       optionalInt32(r.FormValue("discount_days")),
				DiscountPercentBps: optionalInt32(r.FormValue("discount_percent_bps")),
				EntityScope:        r.FormValue("entity_scope"),
				IsDefault:          isDefault,
				Description:        optionalString(r.FormValue("description")),
				DisplayOrder:       optionalInt32(r.FormValue("display_order")),
				ProximateDay:       proximateDay,
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

			currentType := pt.GetType()
			labels := paymenttermform.BuildLabels(viewCtx.T)
			return view.OK("payment-term-drawer-form", &paymenttermform.Data{
				FormAction:         route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:             true,
				ID:                 id,
				Name:               pt.GetName(),
				Code:               pt.GetCode(),
				Type:               currentType,
				NetDays:            netDays,
				DiscountDays:       discountDays,
				DiscountPercentBps: discountPercentBps,
				EntityScope:        pt.GetEntityScope(),
				IsDefault:          pt.GetIsDefault(),
				Description:        pt.GetDescription(),
				DisplayOrder:       displayOrder,
				ProximateDay:       proximateDay,
				Active:             pt.GetActive(),
				TypeOptions:        paymenttermform.BuildTypeOptions(labels, currentType),
				Labels:             labels,
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

		// Validate and nil-out type-driven fields.
		// This mirrors the JS rules object in the template tail script.
		termType := r.FormValue("type")
		netDays, proximateDay, validErr := validateAndNilOutTypeFields(termType, r, viewCtx.T)
		if validErr != "" {
			return entydad.HTMXError(validErr)
		}

		_, err := deps.UpdatePaymentTerm(ctx, &paymenttermpb.UpdatePaymentTermRequest{
			Data: &paymenttermpb.PaymentTerm{
				Id:                 id,
				Active:             active,
				Name:               r.FormValue("name"),
				Code:               r.FormValue("code"),
				Type:               termType,
				NetDays:            netDays,
				DiscountDays:       optionalInt32(r.FormValue("discount_days")),
				DiscountPercentBps: optionalInt32(r.FormValue("discount_percent_bps")),
				EntityScope:        r.FormValue("entity_scope"),
				IsDefault:          isDefault,
				Description:        optionalString(r.FormValue("description")),
				DisplayOrder:       optionalInt32(r.FormValue("display_order")),
				ProximateDay:       proximateDay,
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

package action

import (
	"context"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientattributepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_attribute"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	entityclient "github.com/erniealice/entydad-golang/domain/entity/party/client"
	clientform "github.com/erniealice/entydad-golang/domain/entity/party/client/form"
)

// PaymentTermOption is a type alias so callers wired through module.go
// retain the same surface: form.PaymentTermOption is the source of truth.
type PaymentTermOption = clientform.PaymentTermOption

// Deps holds dependencies for client action handlers.
type Deps struct {
	Routes entityclient.Routes
	// SearchTimezonesURL is the URL of the timezone autocomplete JSON endpoint
	// (owned by the user module; passed through so the client representative
	// section can wire its own auto-complete to the same handler).
	SearchTimezonesURL string
	CreateClient       func(ctx context.Context, req *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error)
	ReadClient         func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	UpdateClient       func(ctx context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error)
	DeleteClient       func(ctx context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error)
	SetClientStatus    func(ctx context.Context, id string, status string) error
	// Payment terms dropdown
	ListPaymentTerms func(ctx context.Context) ([]*PaymentTermOption, error)
	// Tag-related deps for multi-select tags on the client form
	ListCategories       func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListClientCategories func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
	CreateClientCategory func(ctx context.Context, req *clientcategorypb.CreateClientCategoryRequest) (*clientcategorypb.CreateClientCategoryResponse, error)
	DeleteClientCategory func(ctx context.Context, req *clientcategorypb.DeleteClientCategoryRequest) (*clientcategorypb.DeleteClientCategoryResponse, error)
	// GetFunctionalCurrency resolves the current workspace's functional currency
	// so new-client drawers can prefill billing_currency. Optional; returns
	// empty string (or a nil func) means no prefill.
	GetFunctionalCurrency func(ctx context.Context) string
	// CurrencyOptions is the pre-built list of currency select options sourced
	// from lyngua's CommonLabels.Currency.Options. Populated by module.go from
	// ModuleDeps.CommonLabels so the action handler can call
	// clientform.BuildCurrencyOptions without importing pyeza.CommonLabels.
	CurrencyOptions []pyezatypes.SelectOption

	// --- Client attributes (generic EAV overlay, Q-GSE-10) ---
	// EnableAttributes gates the optional Attributes drawer section. When false
	// (or the list closures below are nil) the section is never rendered and the
	// POST never sets req.Attributes/AttributesPresent — the client_attribute
	// rows are left untouched. The espyna use case remains the real validator;
	// these closures only feed the drawer's read/display path.
	EnableAttributes bool
	// ListAttributes returns the active attribute definitions (read-only). The
	// loader filters to the entity/client module (B3: Attribute.module is loose).
	ListAttributes func(ctx context.Context, req *categorypb.ListAttributesRequest) (*categorypb.ListAttributesResponse, error)
	// ListAttributeValues returns the enum option rows for a definition. Uses the
	// generic list path that preserves av.label (W3-part1 MED#6).
	ListAttributeValues func(ctx context.Context, req *categorypb.ListAttributeValuesRequest) (*categorypb.ListAttributeValuesResponse, error)
	// ListClientAttributes returns the current client_attribute rows for a client
	// (edit pre-fill).
	ListClientAttributes func(ctx context.Context, req *clientattributepb.ListClientAttributesRequest) (*clientattributepb.ListClientAttributesResponse, error)
}

// clientAttrModule is the set of Attribute.module values whose definitions are
// relevant to the client drawer. B3: Attribute.module is a loose free-text label
// ("entity"/"general"/"client"/…), not a strict per-junction filter — the live
// education1 attr-gender row carries module="entity". Accept the entity-scope
// buckets a client can carry.
func isClientAttributeModule(module string) bool {
	switch module {
	case "entity", "client", "general":
		return true
	default:
		return false
	}
}

// stringEqualsFilter builds a single-field STRING_EQUALS FilterRequest, matching
// the shape the espyna list repositories expect.
func stringEqualsFilter(field, value string) *categorypb.FilterRequest {
	return &categorypb.FilterRequest{
		Filters: []*categorypb.TypedFilter{{
			Field: field,
			FilterType: &categorypb.TypedFilter_StringFilter{
				StringFilter: &categorypb.StringFilter{
					Value:    value,
					Operator: categorypb.StringOperator_STRING_EQUALS,
				},
			},
		}},
	}
}

// loadClientAttributeData builds the drawer's Attributes-section fields by
// calling the three read use cases: active definitions (entity/client module),
// their enum options (preserving label), and — for edit — the client's current
// values. Every dependency is nil-safe and every failure degrades to an empty
// (unrendered) section: the drawer never blocks on the attribute overlay, and
// the espyna use case remains the authoritative validator on submit.
func loadClientAttributeData(ctx context.Context, deps *Deps, clientID string) []clientform.AttributeField {
	if !deps.EnableAttributes || deps.ListAttributes == nil {
		return nil
	}
	defResp, err := deps.ListAttributes(ctx, &categorypb.ListAttributesRequest{})
	if err != nil {
		log.Printf("Failed to load attribute definitions: %v", err)
		return nil
	}
	var defs []*categorypb.Attribute
	for _, def := range defResp.GetData() {
		if def == nil || !def.GetActive() || def.GetCode() == "" {
			continue
		}
		if !isClientAttributeModule(def.GetModule()) {
			continue
		}
		defs = append(defs, def)
	}
	if len(defs) == 0 {
		return nil
	}

	// Enum options per attribute_id (only the select-family definitions need them,
	// but fetching for all is harmless — non-select fields ignore the map).
	//
	// W3-HIGH-1 FAIL-SAFE: on ANY prefill-data-load error we return nil so the
	// drawer does NOT render the Attributes section. The template gates both the
	// fields AND the hidden `attributes_present=1` marker on a non-empty field list,
	// so an unrendered section means the marker is never posted — and the espyna
	// update-sync path (which is the ONLY thing that deletes client_attribute rows)
	// is skipped entirely. Rendering the section with partially-loaded (blank)
	// options/values would instead post the marker and let a Save silently CLEAR the
	// client's stored attributes. Fail-closed (no section) is the safe degrade.
	optionsByAttr := make(map[string][]*categorypb.AttributeValue, len(defs))
	if deps.ListAttributeValues != nil {
		for _, def := range defs {
			avResp, err := deps.ListAttributeValues(ctx, &categorypb.ListAttributeValuesRequest{
				Filters: stringEqualsFilter("attribute_id", def.GetId()),
			})
			if err != nil {
				log.Printf("Failed to load attribute options for %q: %v — suppressing Attributes section (fail-safe)", def.GetCode(), err)
				return nil
			}
			optionsByAttr[def.GetId()] = avResp.GetData()
		}
	}

	// Current values (edit only), keyed by attribute code.
	currentValues := map[string]string{}
	if clientID != "" && deps.ListClientAttributes != nil {
		caResp, err := deps.ListClientAttributes(ctx, &clientattributepb.ListClientAttributesRequest{
			Filters: stringEqualsFilter("client_id", clientID),
		})
		if err != nil {
			// A transient read failure here previously rendered the section with
			// EMPTY current values; a subsequent Save then posted `attributes_present`
			// with blank fields, which the sync layer interprets as DELETE — silently
			// erasing stored attributes. Suppress the whole section instead.
			log.Printf("Failed to load client attributes for %s: %v — suppressing Attributes section (fail-safe)", clientID, err)
			return nil
		}
		codeByID := make(map[string]string, len(defs))
		for _, def := range defs {
			codeByID[def.GetId()] = def.GetCode()
		}
		for _, ca := range caResp.GetData() {
			if ca == nil || ca.GetClientId() != clientID {
				continue
			}
			if code := codeByID[ca.GetAttributeId()]; code != "" {
				currentValues[code] = ca.GetValue()
			}
		}
	}

	return clientform.BuildAttributeFields(defs, optionsByAttr, currentValues)
}

// parseSubmittedAttributes reads the "attributes.<code>" inputs from a parsed
// form into []AttributeCodeValue (deterministic order). Unknown/crafted codes
// are harmless — the espyna use case validates every code fail-closed. Call only
// after r.ParseForm().
func parseSubmittedAttributes(r *http.Request) []*categorypb.AttributeCodeValue {
	const prefix = "attributes."
	out := make([]*categorypb.AttributeCodeValue, 0, len(r.Form))
	for key, vals := range r.Form {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		code := strings.TrimPrefix(key, prefix)
		if code == "" || len(vals) == 0 {
			continue
		}
		out = append(out, &categorypb.AttributeCodeValue{Code: code, Value: vals[0]})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetCode() < out[j].GetCode() })
	return out
}

// loadPaymentTerms fetches the payment term options. Returns nil slice on error (graceful degradation).
func loadPaymentTerms(ctx context.Context, deps *Deps) []*clientform.PaymentTermOption {
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
// If clientID is provided, marks tags that are currently assigned and populates selected.
func loadTagData(ctx context.Context, deps *Deps, clientID string) ([]clientform.TagOption, []clientform.SelectedTag) {
	if deps.ListCategories == nil {
		return nil, nil
	}

	catResp, err := deps.ListCategories(ctx, &categorypb.ListCategoriesRequest{})
	if err != nil {
		log.Printf("Failed to load tag options: %v", err)
		return nil, nil
	}

	// Build set of assigned category IDs for this client
	assigned := make(map[string]bool)
	if clientID != "" && deps.ListClientCategories != nil {
		ccResp, err := deps.ListClientCategories(ctx, &clientcategorypb.ListClientCategoriesRequest{})
		if err != nil {
			log.Printf("Failed to load client categories: %v", err)
		} else {
			for _, cc := range ccResp.GetData() {
				if cc.GetClientId() == clientID {
					assigned[cc.GetCategoryId()] = true
				}
			}
		}
	}

	var options []clientform.TagOption
	var selected []clientform.SelectedTag
	for _, cat := range catResp.GetData() {
		if cat.GetModule() != "client" || !cat.GetActive() {
			continue
		}
		isAssigned := assigned[cat.GetId()]
		options = append(options, clientform.TagOption{
			Value:    cat.GetId(),
			Label:    cat.GetName(),
			Selected: isAssigned,
		})
		if isAssigned {
			selected = append(selected, clientform.SelectedTag{
				Value: cat.GetId(),
				Label: cat.GetName(),
			})
		}
	}
	return options, selected
}

// syncTags reconciles the submitted tag IDs with existing junction records.
// Creates missing assignments, deletes removed ones.
func syncTags(ctx context.Context, deps *Deps, clientID string, submittedTagIDs []string) {
	if deps.ListClientCategories == nil || deps.CreateClientCategory == nil || deps.DeleteClientCategory == nil {
		return
	}

	// Get current assignments for this client
	ccResp, err := deps.ListClientCategories(ctx, &clientcategorypb.ListClientCategoriesRequest{})
	if err != nil {
		log.Printf("Failed to list client categories for sync: %v", err)
		return
	}

	// Build map of current: categoryID -> junction record ID
	current := make(map[string]string)
	for _, cc := range ccResp.GetData() {
		if cc.GetClientId() == clientID {
			current[cc.GetCategoryId()] = cc.GetId()
		}
	}

	// Build set of desired tag IDs
	desired := make(map[string]bool)
	for _, tagID := range submittedTagIDs {
		if tagID != "" {
			desired[tagID] = true
		}
	}

	// Create new assignments
	for tagID := range desired {
		if _, exists := current[tagID]; !exists {
			_, err := deps.CreateClientCategory(ctx, &clientcategorypb.CreateClientCategoryRequest{
				Data: &clientcategorypb.ClientCategory{
					ClientId:   clientID,
					CategoryId: tagID,
					Active:     true,
				},
			})
			if err != nil {
				log.Printf("Failed to assign tag %s to client %s: %v", tagID, clientID, err)
			}
		}
	}

	// Delete removed assignments
	for tagID, junctionID := range current {
		if !desired[tagID] {
			_, err := deps.DeleteClientCategory(ctx, &clientcategorypb.DeleteClientCategoryRequest{
				Data: &clientcategorypb.ClientCategory{Id: junctionID},
			})
			if err != nil {
				log.Printf("Failed to remove tag %s from client %s: %v", tagID, clientID, err)
			}
		}
	}
}

// optionalString returns a pointer to the string if non-empty, nil otherwise.
// Needed for proto3 optional fields that should not be set when empty.
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

// centavosToDisplay converts int64 centavos to a display string (e.g. 12345 → "123.45").
// Returns empty string when the value is 0 (treat as unset in form pre-fill).
func centavosToDisplay(v int64) string {
	if v == 0 {
		return ""
	}
	return strconv.FormatFloat(float64(v)/100.0, 'f', 2, 64)
}

// parseTagIDs splits a comma-separated string of tag IDs from the multi-select
// hidden input into a slice of individual IDs. Empty strings are filtered out.
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

// NewAddAction creates the client add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("client", "create") {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			mode := viewCtx.Request.URL.Query().Get("mode")
			tagOptions, _ := loadTagData(ctx, deps, "")
			paymentTerms := loadPaymentTerms(ctx, deps)
			attributeFields := loadClientAttributeData(ctx, deps, "")
			labels := clientform.BuildLabels(viewCtx.T)
			functionalCurrency := ""
			if deps.GetFunctionalCurrency != nil {
				functionalCurrency = deps.GetFunctionalCurrency(ctx)
			}
			return view.OK("client-drawer-form", &clientform.Data{
				FormAction:               deps.Routes.AddURL,
				Active:                   true,
				Status:                   "active",
				Mode:                     mode,
				BillingCurrency:          functionalCurrency,
				TaxID:                    "",
				RegistrationNumber:       "",
				CreditLimit:              "",
				LeadTimeDays:             "",
				SearchTimezonesURL:       deps.SearchTimezonesURL,
				PaymentTerms:             paymentTerms,
				PaymentTermSelectOptions: clientform.BuildPaymentTermSelectOptions(paymentTerms, ""),
				StatusOptions:            clientform.BuildStatusOptions("active", labels),
				BillingCurrencyOptions:   clientform.BuildCurrencyOptions(functionalCurrency, deps.CurrencyOptions),
				TagOptions:               tagOptions,
				AttributeFields:          attributeFields,
				Labels:                   labels,
				CommonLabels:             nil, // injected by ViewAdapter
			})
		}

		// POST — create client
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalid_form_data"))
		}

		r := viewCtx.Request
		// New clients default to active=true; the form no longer exposes the
		// active flag (it's derived from status). The use case keeps active
		// in sync with status on subsequent updates.
		repUser := &userpb.User{
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			EmailAddress: r.FormValue("email_address"),
			MobileNumber: r.FormValue("mobile_number"),
			Active:       true,
		}
		if tz := r.FormValue("timezone"); tz != "" {
			repUser.Timezone = &tz
		}

		createReq := &clientpb.CreateClientRequest{
			Data: &clientpb.Client{
				Active:             true,
				Name:               optionalString(r.FormValue("name")),
				Status:             optionalString(r.FormValue("status")),
				Country:            optionalString(r.FormValue("country")),
				Website:            optionalString(r.FormValue("website")),
				StreetAddress:      optionalString(r.FormValue("street_address")),
				City:               optionalString(r.FormValue("city")),
				Province:           optionalString(r.FormValue("province")),
				PostalCode:         optionalString(r.FormValue("postal_code")),
				Notes:              optionalString(r.FormValue("notes")),
				BillingCurrency:    optionalString(r.FormValue("billing_currency")),
				PaymentTermId:      optionalString(r.FormValue("payment_term_id")),
				TaxId:              optionalString(r.FormValue("tax_id")),
				RegistrationNumber: optionalString(r.FormValue("registration_number")),
				CreditLimit:        optionalInt64Money(r.FormValue("credit_limit")),
				LeadTimeDays:       optionalInt32(r.FormValue("lead_time_days")),
				Tin:                optionalString(r.FormValue("tin")),
				CountryCode:        optionalString(r.FormValue("country_code")),
				User:               repUser,
			},
		}
		// Attach generic attributes when the drawer's Attributes section was
		// rendered (marker present). The espyna use case validates + persists them
		// in the same transaction as the client (Q-GSE-10).
		if r.FormValue("attributes_present") == "1" {
			createReq.Attributes = parseSubmittedAttributes(r)
			createReq.AttributesPresent = true
		}
		resp, err := deps.CreateClient(ctx, createReq)
		if err != nil {
			log.Printf("Failed to create client: %v", err)
			return view.HTMXError(err.Error())
		}

		// Sync tags for the newly created client
		if data := resp.GetData(); len(data) > 0 {
			newClientID := data[0].GetId()
			tagIDs := parseTagIDs(r.FormValue("tags"))
			if len(tagIDs) > 0 {
				syncTags(ctx, deps, newClientID, tagIDs)
			}
		}

		return view.HTMXSuccess("clients-table")
	})
}

// NewEditAction creates the client edit action (GET = form, POST = update).
// When the GET request includes ?clone=1, the handler returns the drawer form
// pre-populated from the source record but wired to AddURL (so submission
// creates a new client) with " (Copy)" appended to the name. Tag assignments
// are loaded for UI display but — because the form POSTs to the add handler —
// get attached to the newly created client, not the source.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("client", requiredAction) {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			mode := viewCtx.Request.URL.Query().Get("mode")
			resp, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{
				Data: &clientpb.Client{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read client %s: %v", id, err)
				return view.HTMXError(viewCtx.T("shared.errors.not_found"))
			}

			c := resp.GetData()[0]
			u := c.GetUser()
			tagOptions, selectedTags := loadTagData(ctx, deps, id)
			paymentTerms := loadPaymentTerms(ctx, deps)
			selectedPaymentTermID := c.GetPaymentTermId()
			// Attributes render only on the full drawer (list-page edit / clone),
			// not on the detail-page partial tabs (info/accounting/representative),
			// so a partial-tab save never touches the client_attribute rows.
			var attributeFields []clientform.AttributeField
			if mode == "" {
				attributeFields = loadClientAttributeData(ctx, deps, id)
			}

			name := c.GetName()
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id) + "?mode=" + mode
			formID := id
			if isClone {
				name = strings.TrimSpace(name) + viewCtx.T("actions.copy_suffix")
				formAction = deps.Routes.AddURL
				formID = ""
			}

			labels := clientform.BuildLabels(viewCtx.T)
			creditLimitDisplay := ""
			if c.CreditLimit != nil {
				creditLimitDisplay = centavosToDisplay(c.GetCreditLimit())
			}
			leadTimeDaysDisplay := ""
			if c.LeadTimeDays != nil {
				leadTimeDaysDisplay = strconv.Itoa(int(c.GetLeadTimeDays()))
			}
			return view.OK("client-drawer-form", &clientform.Data{
				FormAction:               formAction,
				IsEdit:                   !isClone,
				ID:                       formID,
				Mode:                     mode,
				Name:                     name,
				FirstName:                u.GetFirstName(),
				LastName:                 u.GetLastName(),
				Email:                    u.GetEmailAddress(),
				Mobile:                   u.GetMobileNumber(),
				Timezone:                 u.GetTimezone(),
				Active:                   c.GetActive(),
				Status:                   c.GetStatus(),
				Country:                  c.GetCountry(),
				Website:                  c.GetWebsite(),
				StreetAddress:            c.GetStreetAddress(),
				City:                     c.GetCity(),
				Province:                 c.GetProvince(),
				PostalCode:               c.GetPostalCode(),
				Notes:                    c.GetNotes(),
				BillingCurrency:          c.GetBillingCurrency(),
				TaxID:                    c.GetTaxId(),
				RegistrationNumber:       c.GetRegistrationNumber(),
				TIN:                      c.GetTin(),
				CountryCode:              c.GetCountryCode(),
				CreditLimit:              creditLimitDisplay,
				LeadTimeDays:             leadTimeDaysDisplay,
				SearchTimezonesURL:       deps.SearchTimezonesURL,
				PaymentTerms:             paymentTerms,
				SelectedPaymentTermID:    selectedPaymentTermID,
				PaymentTermSelectOptions: clientform.BuildPaymentTermSelectOptions(paymentTerms, selectedPaymentTermID),
				StatusOptions:            clientform.BuildStatusOptions(c.GetStatus(), labels),
				BillingCurrencyOptions:   clientform.BuildCurrencyOptions(c.GetBillingCurrency(), deps.CurrencyOptions),
				TagOptions:               tagOptions,
				SelectedTags:             selectedTags,
				AttributeFields:          attributeFields,
				Labels:                   labels,
				CommonLabels:             nil, // injected by ViewAdapter
			})
		}

		// POST — update client
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(viewCtx.T("shared.errors.invalid_form_data"))
		}

		r := viewCtx.Request
		mode := r.URL.Query().Get("mode")

		// The drawer form no longer exposes an "active" toggle. Active is
		// derived from status server-side; when the request payload omits
		// it (the proto3 zero, which we can't distinguish from an explicit
		// false), the Update Client use case copies the value from the
		// existing record so a missing field never deactivates the client.
		clientData := &clientpb.Client{Id: id}
		userData := &userpb.User{}

		switch mode {
		case "info":
			// Only update company-related fields; leave representative fields untouched
			clientData.Name = optionalString(r.FormValue("name"))
			clientData.Website = optionalString(r.FormValue("website"))
			clientData.StreetAddress = optionalString(r.FormValue("street_address"))
			clientData.City = optionalString(r.FormValue("city"))
			clientData.Province = optionalString(r.FormValue("province"))
			clientData.PostalCode = optionalString(r.FormValue("postal_code"))
			clientData.Country = optionalString(r.FormValue("country"))
			clientData.Notes = optionalString(r.FormValue("notes"))
			clientData.PaymentTermId = optionalString(r.FormValue("payment_term_id"))
		case "accounting":
			clientData.Status = optionalString(r.FormValue("status"))
			clientData.BillingCurrency = optionalString(r.FormValue("billing_currency"))
			clientData.PaymentTermId = optionalString(r.FormValue("payment_term_id"))
			clientData.TaxId = optionalString(r.FormValue("tax_id"))
			clientData.RegistrationNumber = optionalString(r.FormValue("registration_number"))
			clientData.CreditLimit = optionalInt64Money(r.FormValue("credit_limit"))
			clientData.LeadTimeDays = optionalInt32(r.FormValue("lead_time_days"))
			clientData.Tin = optionalString(r.FormValue("tin"))
			clientData.CountryCode = optionalString(r.FormValue("country_code"))
		case "representative":
			// Only update representative (user) fields; leave company fields untouched
			userData.FirstName = r.FormValue("first_name")
			userData.LastName = r.FormValue("last_name")
			userData.EmailAddress = r.FormValue("email_address")
			userData.MobileNumber = r.FormValue("mobile_number")
			if tz := r.FormValue("timezone"); tz != "" {
				userData.Timezone = &tz
			}
			clientData.User = userData
		default:
			// List page edit — update all fields
			clientData.Name = optionalString(r.FormValue("name"))
			clientData.Status = optionalString(r.FormValue("status"))
			clientData.Country = optionalString(r.FormValue("country"))
			clientData.Website = optionalString(r.FormValue("website"))
			clientData.StreetAddress = optionalString(r.FormValue("street_address"))
			clientData.City = optionalString(r.FormValue("city"))
			clientData.Province = optionalString(r.FormValue("province"))
			clientData.PostalCode = optionalString(r.FormValue("postal_code"))
			clientData.Notes = optionalString(r.FormValue("notes"))
			clientData.BillingCurrency = optionalString(r.FormValue("billing_currency"))
			clientData.PaymentTermId = optionalString(r.FormValue("payment_term_id"))
			clientData.TaxId = optionalString(r.FormValue("tax_id"))
			clientData.RegistrationNumber = optionalString(r.FormValue("registration_number"))
			clientData.CreditLimit = optionalInt64Money(r.FormValue("credit_limit"))
			clientData.LeadTimeDays = optionalInt32(r.FormValue("lead_time_days"))
			clientData.Tin = optionalString(r.FormValue("tin"))
			clientData.CountryCode = optionalString(r.FormValue("country_code"))
			userData.FirstName = r.FormValue("first_name")
			userData.LastName = r.FormValue("last_name")
			userData.EmailAddress = r.FormValue("email_address")
			userData.MobileNumber = r.FormValue("mobile_number")
			if tz := r.FormValue("timezone"); tz != "" {
				userData.Timezone = &tz
			}
			clientData.User = userData
		}

		updateReq := &clientpb.UpdateClientRequest{
			Data: clientData,
		}
		// Sync attributes ONLY on the full drawer (mode==""), and only when the
		// section was actually rendered (marker present). The detail-page partial
		// tabs never carry the section, so their saves leave attributes untouched
		// (AttributesPresent stays false => the use case skips the attribute sync).
		if mode == "" && r.FormValue("attributes_present") == "1" {
			updateReq.Attributes = parseSubmittedAttributes(r)
			updateReq.AttributesPresent = true
		}
		_, err := deps.UpdateClient(ctx, updateReq)
		if err != nil {
			log.Printf("Failed to update client %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		// Sync tags only when the current mode renders the tags field.
		// representative + accounting forms don't, so FormValue("tags") would
		// be empty and wipe all existing tags.
		if mode != "representative" && mode != "accounting" {
			syncTags(ctx, deps, id, parseTagIDs(r.FormValue("tags")))
		}

		// If mode is set, we're in the detail page context — redirect to the correct tab
		if mode == "info" || mode == "representative" || mode == "accounting" {
			detailURL := route.ResolveURL(deps.Routes.DetailURL, "id", id) + "?tab=" + mode
			return view.ViewResult{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"HX-Trigger":  `{"formSuccess":true}`,
					"HX-Redirect": detailURL,
				},
			}
		}

		return view.HTMXSuccess("clients-table")
	})
}

// NewDeleteAction creates the client delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("client", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(viewCtx.T("shared.errors.id_required"))
		}

		_, err := deps.DeleteClient(ctx, &clientpb.DeleteClientRequest{
			Data: &clientpb.Client{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete client %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("clients-table")
	})
}

// NewBulkDeleteAction creates the client bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("client", "delete") {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.no_ids_provided"))
		}

		for _, id := range ids {
			_, err := deps.DeleteClient(ctx, &clientpb.DeleteClientRequest{
				Data: &clientpb.Client{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete client %s: %v", id, err)
			}
		}

		return view.HTMXSuccess("clients-table")
	})
}

// validClientStatus returns true for the five canonical client lifecycle states.
func validClientStatus(s string) bool {
	switch s {
	case "prospect", "active", "on_hold", "blocked", "inactive":
		return true
	}
	return false
}

// NewSetStatusAction creates the client set-status action (POST only).
// Expects query params: ?id={clientId}&status={prospect|active|on_hold|blocked|inactive}
//
// Uses SetClientStatus (raw map update) instead of UpdateClient (protobuf) because
// proto3's protojson omits bool fields with value false, which means deactivation
// (active=false) would silently be skipped. The closure also keeps the active
// boolean in sync with the status string (active for prospect/active/on_hold/blocked,
// inactive for "inactive") so legacy consumers reading c.active still see consistent values.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("client", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return view.HTMXError(viewCtx.T("shared.errors.id_required"))
		}
		if !validClientStatus(targetStatus) {
			return view.HTMXError(viewCtx.T("shared.errors.invalid_status"))
		}

		if err := deps.SetClientStatus(ctx, id, targetStatus); err != nil {
			log.Printf("Failed to update client status %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("clients-table")
	})
}

// NewBulkSetStatusAction creates the client bulk set-status action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("client", "update") {
			return view.HTMXError(viewCtx.T("shared.errors.permission_denied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return view.HTMXError(viewCtx.T("shared.errors.no_ids_provided"))
		}
		if !validClientStatus(targetStatus) {
			return view.HTMXError(viewCtx.T("shared.errors.invalid_target_status"))
		}

		for _, id := range ids {
			if err := deps.SetClientStatus(ctx, id, targetStatus); err != nil {
				log.Printf("Failed to update client status %s: %v", id, err)
			}
		}

		return view.HTMXSuccess("clients-table")
	})
}

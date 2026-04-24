package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	pyeza "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	"github.com/erniealice/entydad-golang"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Name                     string
	NamePlaceholder          string
	CompanyDetails           string
	Representative           string
	FirstName                string
	FirstNamePlaceholder     string
	LastName                 string
	LastNamePlaceholder      string
	Email                    string
	EmailPlaceholder         string
	Mobile                   string
	MobilePlaceholder        string
	Active                   string
	StreetAddress            string
	StreetAddressPlaceholder string
	City                     string
	CityPlaceholder          string
	Province                 string
	ProvincePlaceholder      string
	PostalCode               string
	PostalCodePlaceholder    string
	Notes                    string
	NotesPlaceholder         string
	PaymentTerms             string
	SelectPaymentTerm        string
	Tags                     string
	TagsPlaceholder          string
	TagsSearchPlaceholder    string
	TagsNoResults            string
	Accounting                 string
	BillingCurrency            string
	BillingCurrencyPlaceholder string
	BillingCurrencyInfo        string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo          string
	EmailInfo         string
	MobileInfo        string
	NotesInfo         string
	PaymentTermsInfo  string
	TagsInfo          string
	ActiveInfo        string
}

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id   string
	Name string
}

// TagOption represents a tag available for selection in the form.
// Fields named Value/Label to match the pyeza multi-select component template.
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

// FormData is the template data for the client drawer form.
type FormData struct {
	FormAction               string
	IsEdit                   bool
	ID                       string
	Mode                     string
	Name                     string
	FirstName                string
	LastName                 string
	Email                    string
	Mobile                   string
	Active                   bool
	StreetAddress            string
	City                     string
	Province                 string
	PostalCode               string
	Notes                    string
	BillingCurrency          string
	PaymentTerms             []*PaymentTermOption
	SelectedPaymentTermID    string
	PaymentTermSelectOptions []pyeza.SelectOption
	TagOptions               []TagOption
	SelectedTags             []SelectedTag
	Labels                   FormLabels
	CommonLabels             any
}

// Deps holds dependencies for client action handlers.
type Deps struct {
	Routes          entydad.ClientRoutes
	CreateClient    func(ctx context.Context, req *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error)
	ReadClient      func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	UpdateClient    func(ctx context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error)
	DeleteClient    func(ctx context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error)
	SetClientActive func(ctx context.Context, id string, active bool) error
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
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Name:                     t("client.form.name"),
		NamePlaceholder:          t("client.form.namePlaceholder"),
		CompanyDetails:           t("client.form.companyDetails"),
		Representative:           t("client.form.representative"),
		FirstName:                t("client.form.firstName"),
		FirstNamePlaceholder:     t("client.form.firstNamePlaceholder"),
		LastName:                 t("client.form.lastName"),
		LastNamePlaceholder:      t("client.form.lastNamePlaceholder"),
		Email:                    t("client.form.email"),
		EmailPlaceholder:         t("client.form.emailPlaceholder"),
		Mobile:                   t("client.form.phone"),
		MobilePlaceholder:        t("client.form.phonePlaceholder"),
		Active:                   t("client.form.active"),
		StreetAddress:            t("client.form.streetAddress"),
		StreetAddressPlaceholder: t("client.form.streetAddressPlaceholder"),
		City:                     t("client.form.city"),
		CityPlaceholder:          t("client.form.cityPlaceholder"),
		Province:                 t("client.form.province"),
		ProvincePlaceholder:      t("client.form.provincePlaceholder"),
		PostalCode:               t("client.form.postalCode"),
		PostalCodePlaceholder:    t("client.form.postalCodePlaceholder"),
		Notes:                    t("client.form.notes"),
		NotesPlaceholder:         t("client.form.notesPlaceholder"),
		PaymentTerms:             t("client.form.paymentTerms"),
		SelectPaymentTerm:        t("client.form.selectPaymentTerm"),
		Tags:                     t("client.form.tags"),
		TagsPlaceholder:          t("client.form.tagsPlaceholder"),
		TagsSearchPlaceholder:    t("client.form.tagsSearchPlaceholder"),
		TagsNoResults:            t("client.form.tagsNoResults"),
		NameInfo:                   t("client.form.nameInfo"),
		EmailInfo:                  t("client.form.emailInfo"),
		MobileInfo:                 t("client.form.mobileInfo"),
		NotesInfo:                  t("client.form.notesInfo"),
		PaymentTermsInfo:           t("client.form.paymentTermsInfo"),
		TagsInfo:                   t("client.form.tagsInfo"),
		ActiveInfo:                 t("client.form.activeInfo"),
		Accounting:                 t("client.form.accounting"),
		BillingCurrency:            t("client.form.billingCurrency"),
		BillingCurrencyPlaceholder: t("client.form.billingCurrencyPlaceholder"),
		BillingCurrencyInfo:        t("client.form.billingCurrencyInfo"),
	}
}

// buildPaymentTermSelectOptions converts a slice of PaymentTermOption into the SelectOption
// format expected by the pyeza form-group select component.
func buildPaymentTermSelectOptions(terms []*PaymentTermOption, selectedID string) []pyeza.SelectOption {
	opts := make([]pyeza.SelectOption, 0, len(terms))
	for _, t := range terms {
		opts = append(opts, pyeza.SelectOption{
			Value:    t.Id,
			Label:    t.Name,
			Selected: t.Id == selectedID,
		})
	}
	return opts
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
// If clientID is provided, marks tags that are currently assigned and populates selected.
func loadTagData(ctx context.Context, deps *Deps, clientID string) ([]TagOption, []SelectedTag) {
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

	var options []TagOption
	var selected []SelectedTag
	for _, cat := range catResp.GetData() {
		if cat.GetModule() != "client" || !cat.GetActive() {
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
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			mode := viewCtx.Request.URL.Query().Get("mode")
			tagOptions, _ := loadTagData(ctx, deps, "")
			paymentTerms := loadPaymentTerms(ctx, deps)
			return view.OK("client-drawer-form", &FormData{
				FormAction:               deps.Routes.AddURL,
				Active:                   true,
				Mode:                     mode,
				BillingCurrency: func() string {
				if deps.GetFunctionalCurrency == nil {
					return ""
				}
				return deps.GetFunctionalCurrency(ctx)
			}(),
				PaymentTerms:             paymentTerms,
				PaymentTermSelectOptions: buildPaymentTermSelectOptions(paymentTerms, ""),
				TagOptions:               tagOptions,
				Labels:                   formLabels(viewCtx.T),
				CommonLabels:             nil, // injected by ViewAdapter
			})
		}

		// POST — create client
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		resp, err := deps.CreateClient(ctx, &clientpb.CreateClientRequest{
			Data: &clientpb.Client{
				Active:          active,
				Name:            optionalString(r.FormValue("name")),
				StreetAddress:   optionalString(r.FormValue("street_address")),
				City:            optionalString(r.FormValue("city")),
				Province:        optionalString(r.FormValue("province")),
				PostalCode:      optionalString(r.FormValue("postal_code")),
				Notes:           optionalString(r.FormValue("notes")),
				BillingCurrency: optionalString(r.FormValue("billing_currency")),
				PaymentTermId:   optionalString(r.FormValue("payment_term_id")),
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
			log.Printf("Failed to create client: %v", err)
			return entydad.HTMXError(err.Error())
		}

		// Sync tags for the newly created client
		if data := resp.GetData(); len(data) > 0 {
			newClientID := data[0].GetId()
			tagIDs := parseTagIDs(r.FormValue("tags"))
			if len(tagIDs) > 0 {
				syncTags(ctx, deps, newClientID, tagIDs)
			}
		}

		return entydad.HTMXSuccess("clients-table")
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
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			mode := viewCtx.Request.URL.Query().Get("mode")
			resp, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{
				Data: &clientpb.Client{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read client %s: %v", id, err)
				return entydad.HTMXError(viewCtx.T("shared.errors.notFound"))
			}

			c := resp.GetData()[0]
			u := c.GetUser()
			tagOptions, selectedTags := loadTagData(ctx, deps, id)
			paymentTerms := loadPaymentTerms(ctx, deps)
			selectedPaymentTermID := c.GetPaymentTermId()

			name := c.GetName()
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id) + "?mode=" + mode
			formID := id
			if isClone {
				name = strings.TrimSpace(name) + viewCtx.T("actions.copySuffix")
				formAction = deps.Routes.AddURL
				formID = ""
			}

			return view.OK("client-drawer-form", &FormData{
				FormAction:               formAction,
				IsEdit:                   !isClone,
				ID:                       formID,
				Mode:                     mode,
				Name:                     name,
				FirstName:                u.GetFirstName(),
				LastName:                 u.GetLastName(),
				Email:                    u.GetEmailAddress(),
				Mobile:                   u.GetMobileNumber(),
				Active:                   c.GetActive(),
				StreetAddress:            c.GetStreetAddress(),
				City:                     c.GetCity(),
				Province:                 c.GetProvince(),
				PostalCode:               c.GetPostalCode(),
				Notes:                    c.GetNotes(),
				BillingCurrency:          c.GetBillingCurrency(),
				PaymentTerms:             paymentTerms,
				SelectedPaymentTermID:    selectedPaymentTermID,
				PaymentTermSelectOptions: buildPaymentTermSelectOptions(paymentTerms, selectedPaymentTermID),
				TagOptions:               tagOptions,
				SelectedTags:             selectedTags,
				Labels:                   formLabels(viewCtx.T),
				CommonLabels:             nil, // injected by ViewAdapter
			})
		}

		// POST — update client
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		mode := r.URL.Query().Get("mode")
		active := r.FormValue("active") == "true"

		clientData := &clientpb.Client{Id: id}
		userData := &userpb.User{}

		switch mode {
		case "info":
			// Only update company-related fields; leave representative fields untouched
			clientData.Active = active
			clientData.Name = optionalString(r.FormValue("name"))
			clientData.StreetAddress = optionalString(r.FormValue("street_address"))
			clientData.City = optionalString(r.FormValue("city"))
			clientData.Province = optionalString(r.FormValue("province"))
			clientData.PostalCode = optionalString(r.FormValue("postal_code"))
			clientData.Notes = optionalString(r.FormValue("notes"))
			clientData.PaymentTermId = optionalString(r.FormValue("payment_term_id"))
		case "accounting":
			// Only update accounting fields. Active is always present on every
			// mode's form so we keep it authoritative.
			clientData.Active = active
			clientData.BillingCurrency = optionalString(r.FormValue("billing_currency"))
		case "representative":
			// Only update representative (user) fields; leave company fields untouched
			userData.FirstName = r.FormValue("first_name")
			userData.LastName = r.FormValue("last_name")
			userData.EmailAddress = r.FormValue("email_address")
			userData.MobileNumber = r.FormValue("mobile_number")
			userData.Active = active
			clientData.User = userData
		default:
			// List page edit — update all fields
			clientData.Active = active
			clientData.Name = optionalString(r.FormValue("name"))
			clientData.StreetAddress = optionalString(r.FormValue("street_address"))
			clientData.City = optionalString(r.FormValue("city"))
			clientData.Province = optionalString(r.FormValue("province"))
			clientData.PostalCode = optionalString(r.FormValue("postal_code"))
			clientData.Notes = optionalString(r.FormValue("notes"))
			clientData.BillingCurrency = optionalString(r.FormValue("billing_currency"))
			clientData.PaymentTermId = optionalString(r.FormValue("payment_term_id"))
			userData.FirstName = r.FormValue("first_name")
			userData.LastName = r.FormValue("last_name")
			userData.EmailAddress = r.FormValue("email_address")
			userData.MobileNumber = r.FormValue("mobile_number")
			userData.Active = active
			clientData.User = userData
		}

		_, err := deps.UpdateClient(ctx, &clientpb.UpdateClientRequest{
			Data: clientData,
		})
		if err != nil {
			log.Printf("Failed to update client %s: %v", id, err)
			return entydad.HTMXError(err.Error())
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

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewDeleteAction creates the client delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("client", "delete") {
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

		_, err := deps.DeleteClient(ctx, &clientpb.DeleteClientRequest{
			Data: &clientpb.Client{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete client %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewBulkDeleteAction creates the client bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("client", "delete") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError(viewCtx.T("shared.errors.noIdsProvided"))
		}

		for _, id := range ids {
			_, err := deps.DeleteClient(ctx, &clientpb.DeleteClientRequest{
				Data: &clientpb.Client{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete client %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewSetStatusAction creates the client activate/deactivate action (POST only).
// Expects query params: ?id={clientId}&status={active|inactive}
//
// Uses SetClientActive (raw map update) instead of UpdateClient (protobuf) because
// proto3's protojson omits bool fields with value false, which means
// deactivation (active=false) would silently be skipped.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("client", "update") {
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

		if err := deps.SetClientActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update client status %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewBulkSetStatusAction creates the client bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("client", "update") {
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
			if err := deps.SetClientActive(ctx, id, active); err != nil {
				log.Printf("Failed to update client status %s: %v", id, err)
			}
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

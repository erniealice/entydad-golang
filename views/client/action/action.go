package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	"github.com/erniealice/entydad-golang"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	FirstName               string
	FirstNamePlaceholder    string
	LastName                string
	LastNamePlaceholder     string
	Email                   string
	EmailPlaceholder        string
	Mobile                  string
	MobilePlaceholder       string
	Active                  string
	CompanyName             string
	CompanyNamePlaceholder  string
	CustomerType            string
	DateOfBirth             string
	StreetAddress           string
	StreetAddressPlaceholder string
	City                    string
	CityPlaceholder         string
	Province                string
	ProvincePlaceholder     string
	PostalCode              string
	PostalCodePlaceholder   string
	Notes                   string
	NotesPlaceholder        string
	Tags                    string
	TagsPlaceholder         string
	TagsSearchPlaceholder   string
	TagsNoResults           string
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
	FormAction    string
	IsEdit        bool
	ID            string
	FirstName     string
	LastName      string
	Email         string
	Mobile        string
	Active        bool
	CompanyName   string
	CustomerType  string
	DateOfBirth   string
	StreetAddress string
	City          string
	Province      string
	PostalCode    string
	Notes         string
	TagOptions    []TagOption
	SelectedTags  []SelectedTag
	Labels        FormLabels
	CommonLabels  any
}

// Deps holds dependencies for client action handlers.
type Deps struct {
	CreateClient    func(ctx context.Context, req *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error)
	ReadClient      func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	UpdateClient    func(ctx context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error)
	DeleteClient    func(ctx context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error)
	SetClientActive func(ctx context.Context, id string, active bool) error
	// Tag-related deps for multi-select tags on the client form
	ListCategories       func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListClientCategories func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
	CreateClientCategory func(ctx context.Context, req *clientcategorypb.CreateClientCategoryRequest) (*clientcategorypb.CreateClientCategoryResponse, error)
	DeleteClientCategory func(ctx context.Context, req *clientcategorypb.DeleteClientCategoryRequest) (*clientcategorypb.DeleteClientCategoryResponse, error)
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		FirstName:                t("client.form.firstName"),
		FirstNamePlaceholder:     t("client.form.firstNamePlaceholder"),
		LastName:                 t("client.form.lastName"),
		LastNamePlaceholder:      t("client.form.lastNamePlaceholder"),
		Email:                    t("client.form.email"),
		EmailPlaceholder:         t("client.form.emailPlaceholder"),
		Mobile:                   t("client.form.phone"),
		MobilePlaceholder:        t("client.form.phonePlaceholder"),
		Active:                   t("client.form.active"),
		CompanyName:              t("client.form.companyName"),
		CompanyNamePlaceholder:   t("client.form.companyNamePlaceholder"),
		CustomerType:             t("client.form.customerType"),
		DateOfBirth:              t("client.form.dateOfBirth"),
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
		Tags:                     t("client.form.tags"),
		TagsPlaceholder:          t("client.form.tagsPlaceholder"),
		TagsSearchPlaceholder:    t("client.form.tagsSearchPlaceholder"),
		TagsNoResults:            t("client.form.tagsNoResults"),
	}
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
		if viewCtx.Request.Method == http.MethodGet {
			tagOptions, _ := loadTagData(ctx, deps, "")
			return view.OK("client-drawer-form", &FormData{
				FormAction:   "/action/clients/add",
				Active:       true,
				TagOptions:   tagOptions,
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create client
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		resp, err := deps.CreateClient(ctx, &clientpb.CreateClientRequest{
			Data: &clientpb.Client{
				Active:        active,
				CompanyName:   optionalString(r.FormValue("company_name")),
				CustomerType:  optionalString(r.FormValue("customer_type")),
				DateOfBirth:   optionalString(r.FormValue("date_of_birth")),
				StreetAddress: optionalString(r.FormValue("street_address")),
				City:          optionalString(r.FormValue("city")),
				Province:      optionalString(r.FormValue("province")),
				PostalCode:    optionalString(r.FormValue("postal_code")),
				Notes:         optionalString(r.FormValue("notes")),
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
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{
				Data: &clientpb.Client{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read client %s: %v", id, err)
				return entydad.HTMXError("Client not found")
			}

			c := resp.GetData()[0]
			u := c.GetUser()
			tagOptions, selectedTags := loadTagData(ctx, deps, id)

			return view.OK("client-drawer-form", &FormData{
				FormAction:    "/action/clients/edit/" + id,
				IsEdit:        true,
				ID:            id,
				FirstName:     u.GetFirstName(),
				LastName:      u.GetLastName(),
				Email:         u.GetEmailAddress(),
				Mobile:        u.GetMobileNumber(),
				Active:        c.GetActive(),
				CompanyName:   c.GetCompanyName(),
				CustomerType:  c.GetCustomerType(),
				DateOfBirth:   c.GetDateOfBirth(),
				StreetAddress: c.GetStreetAddress(),
				City:          c.GetCity(),
				Province:      c.GetProvince(),
				PostalCode:    c.GetPostalCode(),
				Notes:         c.GetNotes(),
				TagOptions:    tagOptions,
				SelectedTags:  selectedTags,
				Labels:        formLabels(viewCtx.T),
				CommonLabels:  nil, // injected by ViewAdapter
			})
		}

		// POST — update client
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		_, err := deps.UpdateClient(ctx, &clientpb.UpdateClientRequest{
			Data: &clientpb.Client{
				Id:            id,
				Active:        active,
				CompanyName:   optionalString(r.FormValue("company_name")),
				CustomerType:  optionalString(r.FormValue("customer_type")),
				DateOfBirth:   optionalString(r.FormValue("date_of_birth")),
				StreetAddress: optionalString(r.FormValue("street_address")),
				City:          optionalString(r.FormValue("city")),
				Province:      optionalString(r.FormValue("province")),
				PostalCode:    optionalString(r.FormValue("postal_code")),
				Notes:         optionalString(r.FormValue("notes")),
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
			log.Printf("Failed to update client %s: %v", id, err)
			return entydad.HTMXError("Failed to update client")
		}

		// Sync tags — multi-select sends comma-separated IDs in a single hidden input
		syncTags(ctx, deps, id, parseTagIDs(r.FormValue("tags")))

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewDeleteAction creates the client delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return entydad.HTMXError("Client ID is required")
		}

		_, err := deps.DeleteClient(ctx, &clientpb.DeleteClientRequest{
			Data: &clientpb.Client{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete client %s: %v", id, err)
			return entydad.HTMXError("Failed to delete client")
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewBulkDeleteAction creates the client bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return entydad.HTMXError("No client IDs provided")
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
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return entydad.HTMXError("Client ID is required")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid status")
		}

		if err := deps.SetClientActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update client status %s: %v", id, err)
			return entydad.HTMXError("Failed to update client status")
		}

		return entydad.HTMXSuccess("clients-table")
	})
}

// NewBulkSetStatusAction creates the client bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return entydad.HTMXError("No client IDs provided")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return entydad.HTMXError("Invalid target status")
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

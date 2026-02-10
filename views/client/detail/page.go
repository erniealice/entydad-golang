package detail

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
)

// Deps holds view dependencies.
type Deps struct {
	ReadClient   func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	Labels       entydad.ClientLabels
	CommonLabels pyeza.CommonLabels
}

// PageData holds the data for the client detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Client          *clientpb.Client
	Labels          entydad.ClientLabels
	ActiveTab       string
	TabItems        []TabItem
	ClientName      string
	ClientEmail     string
	ClientPhone     string
	ClientStatus    string
	StatusVariant   string
}

// NewView creates the client detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "basic"
		}

		resp, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{
			Data: &clientpb.Client{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read client %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load client: %w", err))
		}

		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("client not found"))
		}
		client := data[0]
		u := client.GetUser()

		clientName := clientDisplayName(client)
		clientEmail := ""
		clientPhone := ""
		if u != nil {
			clientEmail = u.GetEmailAddress()
			clientPhone = u.GetMobileNumber()
		}

		clientStatus := "active"
		if !client.GetActive() {
			clientStatus = "inactive"
		}
		statusVariant := "success"
		if clientStatus == "inactive" {
			statusVariant = "warning"
		}

		tabItems := buildTabItems(id)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          clientName,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "clients",
				HeaderTitle:    clientName,
				HeaderSubtitle: clientEmail,
				HeaderIcon:     "icon-user",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "client-detail-content",
			Client:          client,
			Labels:          deps.Labels,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
			ClientName:      clientName,
			ClientEmail:     clientEmail,
			ClientPhone:     clientPhone,
			ClientStatus:    clientStatus,
			StatusVariant:   statusVariant,
		}

		return view.OK("client-detail", pageData)
	})
}

// TabItem represents a tab in the detail view.
// Fields match the pyeza tabs.html template expectations.
type TabItem struct {
	Key      string
	Label    string
	Href     string
	Icon     string
	Count    int
	Disabled bool
}

func buildTabItems(id string) []TabItem {
	base := "/app/clients/" + id
	return []TabItem{
		{Key: "basic", Label: "Basic Information", Href: base + "?tab=basic", Icon: "icon-info"},
		{Key: "history", Label: "Purchase History", Href: base + "?tab=history", Icon: "icon-shopping-bag"},
	}
}

// clientDisplayName returns the client's display name from the embedded user.
func clientDisplayName(c *clientpb.Client) string {
	if u := c.GetUser(); u != nil {
		first := u.GetFirstName()
		last := u.GetLastName()
		if first != "" || last != "" {
			return first + " " + last
		}
		if u.GetEmailAddress() != "" {
			return u.GetEmailAddress()
		}
	}
	return c.GetId()
}

package detail

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	detailform "github.com/erniealice/entydad-golang/views/client/detail/form"
)

// ---------------------------------------------------------------------------
// View-local types — NOT imported from espyna. Views never import espyna
// internals; block.go provides typed callbacks that translate between
// consumer.* shapes and these view-local shapes.
// ---------------------------------------------------------------------------

// RevenueRunScope is the view-layer scope passed to the list / generate callbacks.
type RevenueRunScope struct {
	WorkspaceID    string
	ClientID       string
	SubscriptionID string
	AsOfDate       string // YYYY-MM-DD; empty → use today
	Cursor         string
	Limit          int32
}

// RevenueRunCandidate is the view-layer representation of one pending period.
type RevenueRunCandidate struct {
	SubscriptionID    string
	SubscriptionName  string
	ClientID          string
	ClientName        string
	PlanName          string
	BillingCycleLabel string
	Currency          string
	PeriodStart       string // YYYY-MM-DD
	PeriodEnd         string // YYYY-MM-DD
	PeriodLabel       string
	PeriodMarker      string
	Amount            int64
	AmountDisplay     string
	LineItemCount     int
	Eligible          bool
	BlockerReason     string
}

// SelectedRevenueRunCandidate is one confirmed selection.
type SelectedRevenueRunCandidate struct {
	SubscriptionID string
	PeriodStart    string
	PeriodEnd      string
	PeriodMarker   string
}

// RevenueRunSelections carries either an explicit list or a filter token.
type RevenueRunSelections struct {
	ExplicitList []SelectedRevenueRunCandidate
	FilterToken  string
}

// RevenueRunResult is the output of a successful GenerateRevenueRun call.
type RevenueRunResult struct {
	RunID    string
	Status   string
	Created  int32
	Skipped  int32
	Errored  int32
}

// ---------------------------------------------------------------------------
// NewRevenueRunAction
// ---------------------------------------------------------------------------

// NewRevenueRunAction returns a view.View that serves the per-client
// "Run Invoices" drawer.
//
// GET  → renders the drawer form populated with ListRevenueRunCandidates.
// POST → submits the selected periods via GenerateRevenueRun; on success
//
//	returns headers HX-Trigger(pyeza:toast + refreshTable) so the
//	outstanding-revenue table refreshes in place and a toast fires.
func NewRevenueRunAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("revenue", "create") || !perms.Can("subscription", "read") {
			return entydad.HTMXError(deps.Labels.Detail.RevenueRun.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return entydad.HTMXError(deps.Labels.Detail.RevenueRun.Errors.IDRequired)
		}

		if deps.ListRevenueRunCandidates == nil || deps.GenerateRevenueRun == nil {
			return entydad.HTMXError(deps.Labels.Detail.RevenueRun.Errors.UseCaseUnavailable)
		}

		switch viewCtx.Request.Method {
		case http.MethodGet:
			return renderRevenueRunDrawer(ctx, viewCtx, deps, id)
		case http.MethodPost:
			return submitRevenueRun(ctx, viewCtx, deps, id)
		default:
			return entydad.HTMXError(deps.Labels.Detail.RevenueRun.Errors.InvalidFormData)
		}
	})
}

// renderRevenueRunDrawer handles GET — loads candidates and renders the drawer form.
func renderRevenueRunDrawer(
	ctx context.Context,
	viewCtx *view.ViewContext,
	deps *DetailViewDeps,
	clientID string,
) view.ViewResult {
	l := deps.Labels.Detail.RevenueRun

	// Resolve as-of date: prefer query param, fall back to today.
	asOfDate := viewCtx.Request.URL.Query().Get("as_of_date")
	if asOfDate == "" {
		asOfDate = time.Now().Format("2006-01-02")
	}
	today := time.Now().Format("2006-01-02")

	scope := RevenueRunScope{
		ClientID: clientID,
		AsOfDate: asOfDate,
	}

	candidates, _, err := deps.ListRevenueRunCandidates(ctx, scope)
	if err != nil {
		log.Printf("NewRevenueRunAction GET: failed to list candidates for client %s: %v", clientID, err)
		return entydad.HTMXError(l.Errors.UseCaseUnavailable)
	}

	// Client name is passed by the opener via query param when available.
	clientName := viewCtx.Request.URL.Query().Get("client_name")

	// Client billing currency is passed by the opener for mismatch detection.
	clientCurrency := viewCtx.Request.URL.Query().Get("billing_currency")

	formAction := route.ResolveURL(deps.Routes.RevenueRunURL, "id", clientID)
	fragmentURL := formAction + "?partial=candidates&as_of_date=" + asOfDate

	data := buildRevenueRunDrawerData(
		candidates, clientID, clientName, clientCurrency,
		asOfDate, today, formAction, fragmentURL, l, deps.CommonLabels,
	)

	// Determine which template to render: the outer form or the inner partial.
	// The HTMX inner-swap on date change targets the candidates partial.
	templateName := "client-revenue-run-drawer-form"
	if viewCtx.Request.URL.Query().Get("partial") == "candidates" {
		templateName = "client-revenue-run-candidates"
	}

	return view.OK(templateName, data)
}

// submitRevenueRun handles POST — parses selections, calls GenerateRevenueRun,
// and returns the appropriate HX-Trigger header.
func submitRevenueRun(
	ctx context.Context,
	viewCtx *view.ViewContext,
	deps *DetailViewDeps,
	clientID string,
) view.ViewResult {
	l := deps.Labels.Detail.RevenueRun

	if err := viewCtx.Request.ParseForm(); err != nil {
		return entydad.HTMXError(l.Errors.InvalidFormData)
	}

	asOfDate := viewCtx.Request.FormValue("as_of_date")
	if asOfDate == "" {
		asOfDate = time.Now().Format("2006-01-02")
	}

	// Parse "selection" form values: each is "{sub_id}|{start}|{end}|{marker}"
	rawSelections := viewCtx.Request.Form["selection"]
	if len(rawSelections) == 0 {
		return entydad.HTMXError(l.Errors.SelectOne)
	}

	var selections RevenueRunSelections
	for _, raw := range rawSelections {
		parts := strings.Split(raw, "|")
		if len(parts) != 4 {
			continue
		}
		selections.ExplicitList = append(selections.ExplicitList, SelectedRevenueRunCandidate{
			SubscriptionID: parts[0],
			PeriodStart:    parts[1],
			PeriodEnd:      parts[2],
			PeriodMarker:   parts[3],
		})
	}
	if len(selections.ExplicitList) == 0 {
		return entydad.HTMXError(l.Errors.SelectOne)
	}

	scope := RevenueRunScope{
		ClientID: clientID,
		AsOfDate: asOfDate,
	}

	result, err := deps.GenerateRevenueRun(ctx, scope, selections)
	if err != nil {
		log.Printf("NewRevenueRunAction POST: GenerateRevenueRun failed for client %s: %v", clientID, err)
		return entydad.HTMXError(l.Errors.UseCaseUnavailable)
	}
	if result == nil {
		return entydad.HTMXError(l.Errors.UseCaseUnavailable)
	}

	// Resolve the lyngua-translated toast text. The template uses Go-template
	// placeholders ({{.Created}}, {{.Skipped}}, {{.Errored}}) — substitute
	// inline so the JS-side lf.Toast receives a fully-rendered string.
	toastMessage := strings.NewReplacer(
		"{{.Created}}", fmt.Sprintf("%d", result.Created),
		"{{.Skipped}}", fmt.Sprintf("%d", result.Skipped),
		"{{.Errored}}", fmt.Sprintf("%d", result.Errored),
	).Replace(l.ToastSuccess)

	// Surface A omits the View-run link — entydad does not import centymo's
	// RevenueRunDetailURL constant. Users navigate to run history via the
	// sidebar Revenue Run app.
	triggerPayload, _ := json.Marshal(map[string]any{
		"pyeza:toast": map[string]any{
			"message": toastMessage,
			"state":   toastStateFromResult(result),
		},
		"refreshTable": "client-statement-outstanding-table",
	})

	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": string(triggerPayload),
		},
	}
}

// toastStateFromResult maps a generate-run outcome to a toast state.
// All-errored = error, any-errored = warning, otherwise success.
func toastStateFromResult(r *RevenueRunResult) string {
	if r == nil {
		return "info"
	}
	if r.Errored > 0 && r.Created == 0 {
		return "error"
	}
	if r.Errored > 0 {
		return "warning"
	}
	return "success"
}

// ---------------------------------------------------------------------------
// Builder helpers
// ---------------------------------------------------------------------------

// buildRevenueRunDrawerData constructs the template-facing data from the raw
// candidate slice.
func buildRevenueRunDrawerData(
	candidates []RevenueRunCandidate,
	clientID, clientName, clientCurrency string,
	asOfDate, maxAsOfDate string,
	formAction, fragmentURL string,
	l entydad.ClientRevenueRunLabels,
	commonLabels pyeza.CommonLabels,
) *detailform.RevenueRunDrawerData {
	// Group candidates by SubscriptionID, maintaining insertion order.
	groupOrder := make([]string, 0)
	groupMap := make(map[string]*detailform.RevenueRunGroup)

	for _, c := range candidates {
		if _, exists := groupMap[c.SubscriptionID]; !exists {
			groupOrder = append(groupOrder, c.SubscriptionID)
			g := &detailform.RevenueRunGroup{
				SubscriptionID:    c.SubscriptionID,
				SubscriptionName:  c.SubscriptionName,
				PlanName:          c.PlanName,
				BillingCycleLabel: c.BillingCycleLabel,
				Currency:          c.Currency,
				CurrencyMismatch:  clientCurrency != "" && c.Currency != clientCurrency,
			}
			groupMap[c.SubscriptionID] = g
		}
		g := groupMap[c.SubscriptionID]

		period := detailform.RevenueRunPeriod{
			SubscriptionID: c.SubscriptionID,
			PeriodStart:    c.PeriodStart,
			PeriodEnd:      c.PeriodEnd,
			PeriodMarker:   c.PeriodMarker,
			PeriodLabel:    c.PeriodLabel,
			Amount:         c.Amount,
			AmountDisplay:  c.AmountDisplay,
			LineItemCount:  c.LineItemCount,
			Eligible:       c.Eligible,
			BlockerReason:  c.BlockerReason,
			SelectionValue: fmt.Sprintf("%s|%s|%s|%s",
				c.SubscriptionID, c.PeriodStart, c.PeriodEnd, c.PeriodMarker),
		}
		g.Periods = append(g.Periods, period)

		if c.Eligible {
			g.GroupTotal += c.Amount
			g.HasEligible = true
		}
	}

	// Compute group total display strings and collect subscription groups in order.
	groups := make([]detailform.RevenueRunGroup, 0, len(groupOrder))
	totalsByCurrency := make(map[string]int64)
	eligibleCount := 0
	for _, subID := range groupOrder {
		g := groupMap[subID]
		g.GroupTotalDisplay = formatCentavos(g.GroupTotal)
		if g.HasEligible {
			totalsByCurrency[g.Currency] += g.GroupTotal
		}
		for _, p := range g.Periods {
			if p.Eligible {
				eligibleCount++
			}
		}
		groups = append(groups, *g)
	}

	subtitle := buildSubtitle(l.SubtitleTemplate, eligibleCount, len(groups))

	return &detailform.RevenueRunDrawerData{
		FormAction:         formAction,
		FragmentURL:        fragmentURL,
		ClientID:           clientID,
		ClientName:         clientName,
		AsOfDate:           asOfDate,
		MaxAsOfDate:        maxAsOfDate,
		Subtitle:           subtitle,
		EligibleCount:      eligibleCount,
		SubscriptionCount:  len(groups),
		SubscriptionGroups: groups,
		TotalsByCurrency:   totalsByCurrency,
		Labels:             l,
		CommonLabels:       commonLabels,
	}
}

// buildSubtitle substitutes eligible count and subscription count into the
// SubtitleTemplate. Template uses {eligible} and {subscriptions} tokens.
func buildSubtitle(tmpl string, eligibleCount, subscriptionCount int) string {
	s := strings.ReplaceAll(tmpl, "{eligible}", fmt.Sprintf("%d", eligibleCount))
	s = strings.ReplaceAll(s, "{subscriptions}", fmt.Sprintf("%d", subscriptionCount))
	return s
}

// formatCentavos formats an integer centavo value as a decimal display string
// (e.g. 100050 → "1,000.50"). Plain formatting without currency symbol.
func formatCentavos(centavos int64) string {
	if centavos == 0 {
		return "0.00"
	}
	negative := centavos < 0
	if negative {
		centavos = -centavos
	}
	whole := centavos / 100
	frac := centavos % 100
	// Group with commas every three digits.
	s := fmt.Sprintf("%d", whole)
	if len(s) > 3 {
		n := len(s)
		var b strings.Builder
		for i, ch := range s {
			if i > 0 && (n-i)%3 == 0 {
				b.WriteRune(',')
			}
			b.WriteRune(ch)
		}
		s = b.String()
	}
	result := fmt.Sprintf("%s.%02d", s, frac)
	if negative {
		return "-" + result
	}
	return result
}


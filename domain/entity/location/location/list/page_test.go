package list

import (
	"context"
	"errors"
	"fmt"
	location "github.com/erniealice/entydad-golang/domain/entity/location/location"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	"testing"
)

func TestLocationListBooleanFilter(t *testing.T) {
	stop := errors.New("captured request")
	for _, status := range []string{"active", "inactive"} {
		t.Run(status, func(t *testing.T) {
			called := false
			deps := &ListViewDeps{GetListPageData: func(_ context.Context, req *locationpb.GetLocationListPageDataRequest) (*locationpb.GetLocationListPageDataResponse, error) {
				called = true
				filters := req.GetFilters().GetFilters()
				if len(filters) != 1 {
					t.Fatalf("expected one status filter, got %d", len(filters))
				}
				f := filters[0]
				if f.GetField() != "active" {
					t.Fatalf("adapter requires canonical active field, got %q", f.GetField())
				}
				if f.GetBooleanFilter() == nil || f.GetBooleanFilter().GetValue() != (status == "active") {
					t.Fatalf("wrong boolean for %s: %v", status, f)
				}
				return nil, stop
			}}
			_, err := buildTableConfig(context.Background(), deps, nil, status, tableparams.TableQueryParams{})
			if !called || !errors.Is(err, stop) {
				t.Fatalf("expected captured request, called=%v err=%v", called, err)
			}
		})
	}
}

func TestLocationAddRequiresActiveStatus(t *testing.T) {
	for _, status := range []string{"active", "inactive", "planned", "operational", "under_maintenance", "temporarily_closed", "closed", "unknown", "", "ACTIVE"} {
		for _, canCreate := range []bool{true, false} {
			t.Run(status+fmt.Sprint(canCreate), func(t *testing.T) {
				codes := []string{"location:list"}
				if canCreate {
					codes = append(codes, "location:create")
				}
				ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions(codes))
				deps := &ListViewDeps{
					Labels: location.Labels{Buttons: location.ButtonLabels{AddLocationActiveOnly: "active only"}},
					GetListPageData: func(context.Context, *locationpb.GetLocationListPageDataRequest) (*locationpb.GetLocationListPageDataResponse, error) {
						return &locationpb.GetLocationListPageDataResponse{}, nil
					},
				}
				table, err := buildTableConfig(ctx, deps, nil, status, tableparams.TableQueryParams{Page: 1, PageSize: 25})
				if err != nil {
					t.Fatal(err)
				}
				if got, want := table.PrimaryAction.Disabled, status != "active" || !canCreate; got != want {
					t.Fatalf("disabled=%v want %v", got, want)
				}
				if status != "active" && table.PrimaryAction.DisabledTooltip != "active only" {
					t.Fatal("missing active-only explanation")
				}
			})
		}
	}
}

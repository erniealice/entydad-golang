package action

import (
	"context"
	"errors"
	"testing"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientattributepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_attribute"
)

// genderDefsResp is one ACTIVE entity-module select definition — enough for the
// drawer to want to render the Attributes section.
func genderDefsResp() *categorypb.ListAttributesResponse {
	return &categorypb.ListAttributesResponse{Data: []*categorypb.Attribute{
		{Id: "attr-gender", Code: "gender", DataType: "option", Module: "entity", Active: true},
	}}
}

func genderValsResp() *categorypb.ListAttributeValuesResponse {
	return &categorypb.ListAttributeValuesResponse{Data: []*categorypb.AttributeValue{
		{Id: "v1", AttributeId: "attr-gender", Value: "male", Active: true},
		{Id: "v2", AttributeId: "attr-gender", Value: "female", Active: true},
	}}
}

// TestLoadClientAttributeData_PrefillErrorSuppressesSection proves the W3-HIGH-1
// fail-safe: a transient client-attribute (or option) prefill read error yields
// NO AttributeFields. Because the template gates both the fields and the hidden
// `attributes_present=1` marker on a non-empty field list, an empty result means
// the marker is never posted — so a subsequent Save cannot be misread as a request
// to CLEAR the client's stored attributes. Before the fix the section rendered with
// blank values and this test would have returned a non-empty slice.
func TestLoadClientAttributeData_PrefillErrorSuppressesSection(t *testing.T) {
	ctx := context.Background()

	t.Run("client-attribute read error => no section (no marker)", func(t *testing.T) {
		deps := &Deps{
			EnableAttributes: true,
			ListAttributes: func(context.Context, *categorypb.ListAttributesRequest) (*categorypb.ListAttributesResponse, error) {
				return genderDefsResp(), nil
			},
			ListAttributeValues: func(context.Context, *categorypb.ListAttributeValuesRequest) (*categorypb.ListAttributeValuesResponse, error) {
				return genderValsResp(), nil
			},
			ListClientAttributes: func(context.Context, *clientattributepb.ListClientAttributesRequest) (*clientattributepb.ListClientAttributesResponse, error) {
				return nil, errors.New("transient db error")
			},
		}
		got := loadClientAttributeData(ctx, deps, "client-1")
		if len(got) != 0 {
			t.Fatalf("prefill error must suppress the Attributes section (no attributes_present marker); got %d fields", len(got))
		}
	})

	t.Run("option read error => no section (no marker)", func(t *testing.T) {
		deps := &Deps{
			EnableAttributes: true,
			ListAttributes: func(context.Context, *categorypb.ListAttributesRequest) (*categorypb.ListAttributesResponse, error) {
				return genderDefsResp(), nil
			},
			ListAttributeValues: func(context.Context, *categorypb.ListAttributeValuesRequest) (*categorypb.ListAttributeValuesResponse, error) {
				return nil, errors.New("transient option read error")
			},
			ListClientAttributes: func(context.Context, *clientattributepb.ListClientAttributesRequest) (*clientattributepb.ListClientAttributesResponse, error) {
				return &clientattributepb.ListClientAttributesResponse{}, nil
			},
		}
		got := loadClientAttributeData(ctx, deps, "client-1")
		if len(got) != 0 {
			t.Fatalf("option-load error must suppress the section; got %d fields", len(got))
		}
	})

	t.Run("all reads succeed => section renders", func(t *testing.T) {
		deps := &Deps{
			EnableAttributes: true,
			ListAttributes: func(context.Context, *categorypb.ListAttributesRequest) (*categorypb.ListAttributesResponse, error) {
				return genderDefsResp(), nil
			},
			ListAttributeValues: func(context.Context, *categorypb.ListAttributeValuesRequest) (*categorypb.ListAttributeValuesResponse, error) {
				return genderValsResp(), nil
			},
			ListClientAttributes: func(context.Context, *clientattributepb.ListClientAttributesRequest) (*clientattributepb.ListClientAttributesResponse, error) {
				return &clientattributepb.ListClientAttributesResponse{Data: []*clientattributepb.ClientAttribute{
					{Id: "row1", ClientId: "client-1", AttributeId: "attr-gender", Value: "male", Active: true},
				}}, nil
			},
		}
		got := loadClientAttributeData(ctx, deps, "client-1")
		if len(got) == 0 {
			t.Fatal("expected the Attributes section to render when every prefill read succeeds")
		}
	})
}

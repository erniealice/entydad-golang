// commerce.go — block sub-context lift (B, block-go-anatomy).
//
// wireCommerceModule registers the commerce/location sub-context entity
// modules (location, location_area, payment_term) into the app router. It is
// a PURE code-move of the corresponding `if cfg.enableAll || cfg.X { ... }`
// blocks from block.go's Block() — same construction order, registration
// order, callbacks, and nil-checks. No behaviour change.
//
// All deps the lifted bodies need from Block()'s scope are carried on
// commerceWiring (block-go-anatomy: >6 deps → struct).
package block

import (
	"context"
	"fmt"
	"log"

	commerce "github.com/erniealice/entydad-golang/domain/entity/commerce"
	location "github.com/erniealice/entydad-golang/domain/entity/location"
	locationaction "github.com/erniealice/entydad-golang/domain/entity/location/location/action"
	locationareaaction "github.com/erniealice/entydad-golang/domain/entity/location/location_area/action"
	locationarealist "github.com/erniealice/entydad-golang/domain/entity/location/location_area/list"
	"github.com/erniealice/espyna-golang/reference"
	"github.com/erniealice/espyna-golang/registry"
	entityid "github.com/erniealice/espyna-golang/registry/entityid"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	locationareapb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location_area"
	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"
	pyeza "github.com/erniealice/pyeza-golang"
)

// commerceWiring carries everything the commerce/location cluster needs from
// Block()'s scope. Implementation detail of the wiring; never re-exported.
type commerceWiring struct {
	cfg        *blockConfig
	uc         *UseCases
	labels     blockLabels
	routes     blockRoutes
	refChecker reference.Checker

	uploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	listAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID  func() string
}

func wireCommerceModule(ctx *pyeza.AppContext, w commerceWiring) error {
	cfg := w.cfg
	uc := w.uc
	labels := w.labels
	routes := w.routes
	refChecker := w.refChecker
	uploadFile := w.uploadFile
	listAttachments := w.listAttachments
	createAttachment := w.createAttachment
	deleteAttachment := w.deleteAttachment
	newAttachmentID := w.newAttachmentID

	if cfg.enableAll || cfg.location {
		locationDeps := &location.LocationModuleDeps{
			Routes:             routes.Location,
			LocationAreaRoutes: routes.LocationArea,
			CommonLabels:       ctx.Common,
			SharedLabels:       labels.Shared,
			Labels:             labels.Location,
			TableLabels:        ctx.Table,
			GetListPageData:    uc.Location.GetListPageData,
			GetInUseIDs:        refChecker.GetLocationInUseIDs,
			CreateLocation:     uc.Location.Create,
			ReadLocation:       uc.Location.Read,
			UpdateLocation:     uc.Location.Update,
			DeleteLocation:     uc.Location.Delete,
			SetActive:          setActiveClosure(uc, "location"),
			UploadFile:         uploadFile,
			ListAttachments:    listAttachments,
			CreateAttachment:   createAttachment,
			DeleteAttachment:   deleteAttachment,
			NewID:              newAttachmentID,
		}
		if uc.LocationArea.List != nil {
			listLocationAreas := uc.LocationArea.List
			locationDeps.ListLocationAreas = func(fctx context.Context) ([]locationaction.LocationAreaOption, error) {
				resp, err := listLocationAreas(fctx, &locationareapb.ListLocationAreasRequest{})
				if err != nil {
					return nil, err
				}
				rows := resp.GetData()
				opts := make([]locationaction.LocationAreaOption, 0, len(rows))
				for _, row := range rows {
					if !row.GetActive() {
						continue
					}
					id := row.GetId()
					if id == "" {
						continue
					}
					opts = append(opts, locationaction.LocationAreaOption{ID: id, Name: row.GetName()})
				}
				return opts, nil
			}
		}
		if uc.GetLocationDashboardPageData != nil {
			locationDeps.GetLocationDashboardPageData = uc.GetLocationDashboardPageData
		}
		location.NewLocationModule(locationDeps).RegisterRoutes(ctx.Routes)
	}

	if cfg.enableAll || cfg.locationArea {
		la := uc.LocationArea
		if la.List == nil || la.Create == nil || la.Read == nil || la.Update == nil || la.Delete == nil {
			log.Println("entydad.Block: warning: LocationArea use cases not wired — skipping location_area module")
		} else {
			location.NewLocationAreaModule(&location.LocationAreaModuleDeps{
				Routes:       routes.LocationArea,
				CommonLabels: ctx.Common,
				SharedLabels: labels.Shared,
				Labels:       labels.LocationArea,
				TableLabels:  ctx.Table,
				GetListPageData: func(fctx context.Context, status string, search string, page, pageSize int) (*locationarealist.LocationAreaListResult, error) {
					resp, err := la.List(fctx, &locationareapb.ListLocationAreasRequest{})
					if err != nil {
						return nil, err
					}
					rows := resp.GetData()
					items := make([]*locationarealist.LocationAreaItem, 0, len(rows))
					for _, row := range rows {
						active := row.GetActive()
						recordStatus := "active"
						if !active {
							recordStatus = "inactive"
						}
						if recordStatus != status {
							continue
						}
						items = append(items, &locationarealist.LocationAreaItem{
							ID:          row.GetId(),
							Name:        row.GetName(),
							Description: row.GetDescription(),
							Active:      active,
							DateCreated: row.GetDateCreatedString(),
						})
					}
					return &locationarealist.LocationAreaListResult{Items: items, TotalItems: len(items)}, nil
				},
				GetInUseIDs: refChecker.GetLocationAreaInUseIDs,
				CreateLocationArea: func(fctx context.Context, name, description string, active bool) (string, error) {
					resp, err := la.Create(fctx, &locationareapb.CreateLocationAreaRequest{
						Data: &locationareapb.LocationArea{
							Name:        name,
							Description: description,
							Active:      active,
						},
					})
					if err != nil {
						return "", err
					}
					if data := resp.GetData(); len(data) > 0 {
						return data[0].GetId(), nil
					}
					return "", nil
				},
				ReadLocationArea: func(fctx context.Context, id string) (*locationareaaction.LocationAreaRecord, error) {
					resp, err := la.Read(fctx, &locationareapb.ReadLocationAreaRequest{
						Data: &locationareapb.LocationArea{Id: id},
					})
					if err != nil {
						return nil, err
					}
					data := resp.GetData()
					if len(data) == 0 {
						return nil, nil
					}
					row := data[0]
					return &locationareaaction.LocationAreaRecord{
						ID:          row.GetId(),
						Name:        row.GetName(),
						Description: row.GetDescription(),
						Active:      row.GetActive(),
					}, nil
				},
				UpdateLocationArea: func(fctx context.Context, id, name, description string, active bool) error {
					_, err := la.Update(fctx, &locationareapb.UpdateLocationAreaRequest{
						Data: &locationareapb.LocationArea{
							Id:          id,
							Name:        name,
							Description: description,
							Active:      active,
						},
					})
					return err
				},
				DeleteLocationArea: func(fctx context.Context, id string) error {
					_, err := la.Delete(fctx, &locationareapb.DeleteLocationAreaRequest{
						Data: &locationareapb.LocationArea{Id: id},
					})
					return err
				},
				SetLocationAreaActive: func(fctx context.Context, id string, active bool) error {
					// Read-modify-write: the typed UpdateLocationArea use case
					// validates Name is required, so flipping `active` must
					// preserve the existing name/description (the former duck
					// path issued a partial column update on {active} alone).
					readResp, err := la.Read(fctx, &locationareapb.ReadLocationAreaRequest{
						Data: &locationareapb.LocationArea{Id: id},
					})
					if err != nil {
						return err
					}
					name, description := "", ""
					if data := readResp.GetData(); len(data) > 0 {
						name = data[0].GetName()
						description = data[0].GetDescription()
					}
					_, err = la.Update(fctx, &locationareapb.UpdateLocationAreaRequest{
						Data: &locationareapb.LocationArea{
							Id:          id,
							Name:        name,
							Description: description,
							Active:      active,
						},
					})
					return err
				},
			}).RegisterRoutes(ctx.Routes)
		}
	}

	if cfg.enableAll || cfg.paymentTerm {
		if ctx.SqlDB == nil {
			log.Println("entydad.Block: warning: SqlDB is nil — skipping payment_term module")
		} else {
			repoAny, err := registry.CreateRepository("postgresql", entityid.PaymentTerm, ctx.SqlDB, entityid.PaymentTerm)
			if err != nil {
				return fmt.Errorf("entydad.Block: failed to create payment_term repository: %w", err)
			}
			ptRepo, ok := repoAny.(paymenttermpb.PaymentTermDomainServiceServer)
			if !ok {
				return fmt.Errorf("entydad.Block: payment_term repository does not implement PaymentTermDomainServiceServer")
			}
			setPaymentTermActive := setActiveClosure(uc, "payment_term")
			// Client-context payment term list: shows terms with entity_scope IN ('client', 'both')
			commerce.NewPaymentTermModule(&commerce.PaymentTermModuleDeps{
				Routes:               routes.PaymentTerm,
				CommonLabels:         ctx.Common,
				SharedLabels:         labels.Shared,
				Labels:               labels.PaymentTerm,
				TableLabels:          ctx.Table,
				GetListPageData:      ptRepo.GetPaymentTermListPageData,
				GetInUseIDs:          refChecker.GetPaymentTermInUseIDs,
				CreatePaymentTerm:    ptRepo.CreatePaymentTerm,
				ReadPaymentTerm:      ptRepo.ReadPaymentTerm,
				UpdatePaymentTerm:    ptRepo.UpdatePaymentTerm,
				DeletePaymentTerm:    ptRepo.DeletePaymentTerm,
				SetPaymentTermActive: setPaymentTermActive,
				Scope:                "client",
			}).RegisterRoutes(ctx.Routes)
			// Supplier-context payment term list: shows terms with entity_scope IN ('supplier', 'both')
			commerce.NewPaymentTermModule(&commerce.PaymentTermModuleDeps{
				Routes:               routes.SupplierPaymentTerm.ToRoutes(),
				CommonLabels:         ctx.Common,
				SharedLabels:         labels.Shared,
				Labels:               labels.PaymentTerm,
				TableLabels:          ctx.Table,
				GetListPageData:      ptRepo.GetPaymentTermListPageData,
				GetInUseIDs:          refChecker.GetPaymentTermInUseIDs,
				CreatePaymentTerm:    ptRepo.CreatePaymentTerm,
				ReadPaymentTerm:      ptRepo.ReadPaymentTerm,
				UpdatePaymentTerm:    ptRepo.UpdatePaymentTerm,
				DeletePaymentTerm:    ptRepo.DeletePaymentTerm,
				SetPaymentTermActive: setPaymentTermActive,
				Scope:                "supplier",
			}).RegisterRoutes(ctx.Routes)
		}
	}
	return nil
}

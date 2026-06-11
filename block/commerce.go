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
	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"
	pyeza "github.com/erniealice/pyeza-golang"
)

// commerceWiring carries everything the commerce/location cluster needs from
// Block()'s scope. Implementation detail of the wiring; never re-exported.
type commerceWiring struct {
	cfg        *blockConfig
	uc         *UseCases
	db         UpdateableSource
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
	db := w.db
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
			SetActive: func(fctx context.Context, id string, active bool) error {
				_, err := db.Update(fctx, "location", id, map[string]any{"active": active})
				return err
			},
			UploadFile:       uploadFile,
			ListAttachments:  listAttachments,
			CreateAttachment: createAttachment,
			DeleteAttachment: deleteAttachment,
			NewID:            newAttachmentID,
		}
		if crudDB, hasCRUD := db.(CRUDSource); hasCRUD {
			locationDeps.ListLocationAreas = func(fctx context.Context) ([]locationaction.LocationAreaOption, error) {
				rows, err := crudDB.ListSimple(fctx, "location_area")
				if err != nil {
					return nil, err
				}
				opts := make([]locationaction.LocationAreaOption, 0, len(rows))
				for _, row := range rows {
					active, _ := row["active"].(bool)
					if !active {
						continue
					}
					id, _ := row["id"].(string)
					name, _ := row["name"].(string)
					if id == "" {
						continue
					}
					opts = append(opts, locationaction.LocationAreaOption{ID: id, Name: name})
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
		crudDB, hasCRUD := db.(CRUDSource)
		if !hasCRUD {
			log.Println("entydad.Block: warning: DB does not implement CRUDSource — skipping location_area module")
		} else {
			location.NewLocationAreaModule(&location.LocationAreaModuleDeps{
				Routes:       routes.LocationArea,
				CommonLabels: ctx.Common,
				SharedLabels: labels.Shared,
				Labels:       labels.LocationArea,
				TableLabels:  ctx.Table,
				GetListPageData: func(fctx context.Context, status string, search string, page, pageSize int) (*locationarealist.LocationAreaListResult, error) {
					rows, err := crudDB.ListSimple(fctx, "location_area")
					if err != nil {
						return nil, err
					}
					items := make([]*locationarealist.LocationAreaItem, 0, len(rows))
					for _, row := range rows {
						active, _ := row["active"].(bool)
						recordStatus := "active"
						if !active {
							recordStatus = "inactive"
						}
						if recordStatus != status {
							continue
						}
						id, _ := row["id"].(string)
						name, _ := row["name"].(string)
						description, _ := row["description"].(string)
						dateCreated, _ := row["date_created"].(string)
						items = append(items, &locationarealist.LocationAreaItem{
							ID:          id,
							Name:        name,
							Description: description,
							Active:      active,
							DateCreated: dateCreated,
						})
					}
					return &locationarealist.LocationAreaListResult{Items: items, TotalItems: len(items)}, nil
				},
				GetInUseIDs: refChecker.GetLocationAreaInUseIDs,
				CreateLocationArea: func(fctx context.Context, name, description string, active bool) (string, error) {
					row, err := crudDB.Create(fctx, "location_area", map[string]any{
						"name":        name,
						"description": description,
						"active":      active,
					})
					if err != nil {
						return "", err
					}
					id, _ := row["id"].(string)
					return id, nil
				},
				ReadLocationArea: func(fctx context.Context, id string) (*locationareaaction.LocationAreaRecord, error) {
					row, err := crudDB.Read(fctx, "location_area", id)
					if err != nil {
						return nil, err
					}
					name, _ := row["name"].(string)
					description, _ := row["description"].(string)
					active, _ := row["active"].(bool)
					return &locationareaaction.LocationAreaRecord{
						ID:          id,
						Name:        name,
						Description: description,
						Active:      active,
					}, nil
				},
				UpdateLocationArea: func(fctx context.Context, id, name, description string, active bool) error {
					_, err := crudDB.Update(fctx, "location_area", id, map[string]any{
						"name":        name,
						"description": description,
						"active":      active,
					})
					return err
				},
				DeleteLocationArea: func(fctx context.Context, id string) error {
					return crudDB.Delete(fctx, "location_area", id)
				},
				SetLocationAreaActive: func(fctx context.Context, id string, active bool) error {
					_, err := crudDB.Update(fctx, "location_area", id, map[string]any{"active": active})
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
			setPaymentTermActive := func(fctx context.Context, id string, active bool) error {
				_, err := db.Update(fctx, "payment_term", id, map[string]any{"active": active})
				return err
			}
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

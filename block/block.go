// Package block provides the Block() function for registering all entydad
// entity modules into a pyeza app via AppOption composition.
//
// It lives in a sub-package to avoid the import cycle that would arise if it
// were placed in the root entydad package (the view sub-packages already import
// the root package for route/label types).
//
// Usage:
//
//	import entydadblock "github.com/erniealice/entydad-golang/block"
//
//	app, err := pyeza.NewApp(
//	    entydadblock.Block(),                           // all modules
//	    entydadblock.Block(entydadblock.WithClient()),  // client only
//	)
package block

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/entydad-golang"
	clientmod "github.com/erniealice/entydad-golang/views/client"
	clienttagmod "github.com/erniealice/entydad-golang/views/clienttag"
	locationmod "github.com/erniealice/entydad-golang/views/location"
	locationareamod "github.com/erniealice/entydad-golang/views/location_area"
	locationareaaction "github.com/erniealice/entydad-golang/views/location_area/action"
	locationarealist "github.com/erniealice/entydad-golang/views/location_area/list"
	paymenttermmod "github.com/erniealice/entydad-golang/views/payment_term"
	permissionmod "github.com/erniealice/entydad-golang/views/permission"
	rolemod "github.com/erniealice/entydad-golang/views/role"
	roleusers "github.com/erniealice/entydad-golang/views/role/users"
	suppliermod "github.com/erniealice/entydad-golang/views/supplier"
	usermod "github.com/erniealice/entydad-golang/views/user"
	userdashboard "github.com/erniealice/entydad-golang/views/user/dashboard"
	workspacemod "github.com/erniealice/entydad-golang/views/workspace"
	"github.com/erniealice/espyna-golang/consumer"
	"github.com/erniealice/espyna-golang/contrib/postgres/reference"
	"github.com/erniealice/espyna-golang/registry"
	entityid "github.com/erniealice/espyna-golang/registry/entityid"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
)

// routeRegistrarFull — optional extension for raw http.HandlerFunc routes.
// Apps that implement HandleFunc (e.g., service-admin chi router wrapper) can
// register JSON endpoints and other raw HTTP handlers. Apps that don't will skip.
type routeRegistrarFull interface {
	pyeza.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// handleFunc registers an http.HandlerFunc route if the registrar supports it.
// Silently skips if the registrar does not implement HandleFunc.
func handleFunc(r pyeza.RouteRegistrar, method, path string, handler http.HandlerFunc) {
	if full, ok := r.(routeRegistrarFull); ok {
		full.HandleFunc(method, path, handler)
		return
	}
	log.Printf("entydad.Block: RouteRegistrar does not support HandleFunc — skipping %s %s", method, path)
}

// BlockOption configures which entydad modules are registered by Block().
type BlockOption func(*blockConfig)

type blockConfig struct {
	enableAll    bool
	client       bool
	clientTag    bool
	paymentTerm  bool
	user         bool
	role         bool
	location     bool
	locationArea bool
	permission   bool
	workspace    bool
	supplier     bool
}

// WithClient enables the Client module in Block().
func WithClient() BlockOption { return func(c *blockConfig) { c.client = true } }

// WithClientTag enables the ClientTag module in Block().
func WithClientTag() BlockOption { return func(c *blockConfig) { c.clientTag = true } }

// WithPaymentTerm enables the PaymentTerm module in Block().
func WithPaymentTerm() BlockOption { return func(c *blockConfig) { c.paymentTerm = true } }

// WithUser enables the User module in Block().
func WithUser() BlockOption { return func(c *blockConfig) { c.user = true } }

// WithRole enables the Role module in Block().
func WithRole() BlockOption { return func(c *blockConfig) { c.role = true } }

// WithLocation enables the Location module in Block().
func WithLocation() BlockOption { return func(c *blockConfig) { c.location = true } }

// WithLocationArea enables the LocationArea module in Block().
func WithLocationArea() BlockOption { return func(c *blockConfig) { c.locationArea = true } }

// WithPermission enables the Permission module in Block().
func WithPermission() BlockOption { return func(c *blockConfig) { c.permission = true } }

// WithWorkspace enables the Workspace module in Block().
func WithWorkspace() BlockOption { return func(c *blockConfig) { c.workspace = true } }

// WithSupplier enables the Supplier module in Block().
func WithSupplier() BlockOption { return func(c *blockConfig) { c.supplier = true } }

// Block returns a pyeza.AppOption that registers entydad entity modules into the app.
// When called with no options, all modules are registered (enableAll mode).
// When called with specific WithXxx() options, only those modules are registered.
//
// Expected ctx fields (type-asserted from any):
//   - ctx.UseCases     → *consumer.UseCases
//   - ctx.DB           → UpdateableSource (entydad.DataSource + Update method)
//   - ctx.RefChecker   → *reference.Checker
//   - ctx.Translations → *lynguaV1.TranslationProvider
//   - ctx.UploadFile, ctx.ListAttachments, ctx.CreateAttachment,
//     ctx.DeleteAttachment, ctx.NewAttachmentID — attachment funcs
//   - ctx.GetUsersByRoleID, ctx.GetDashboardData, ctx.HashPassword,
//     ctx.GetUserRolesMap — user/role helpers
//   - ctx.Routes, ctx.Common, ctx.Table, ctx.BusinessType — from pyeza.AppContext
func Block(opts ...BlockOption) pyeza.AppOption {
	cfg := &blockConfig{enableAll: len(opts) == 0}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(ctx *pyeza.AppContext) error {
		// --- type-assert opaque fields ---

		uc, ok := ctx.UseCases.(*consumer.UseCases)
		if !ok || uc == nil {
			return fmt.Errorf("entydad.Block: UseCases must be *consumer.UseCases")
		}
		if uc.Entity == nil {
			return fmt.Errorf("entydad.Block: entity use cases not initialized")
		}

		db, ok := ctx.DB.(UpdateableSource)
		if !ok {
			return fmt.Errorf("entydad.Block: DB must implement block.UpdateableSource (DataSource + Update)")
		}

		refChecker, ok := ctx.RefChecker.(*reference.Checker)
		if !ok {
			return fmt.Errorf("entydad.Block: RefChecker must be *reference.Checker")
		}

		translations, ok := ctx.Translations.(*lynguaV1.TranslationProvider)
		if !ok {
			return fmt.Errorf("entydad.Block: Translations must be *lynguaV1.TranslationProvider")
		}

		// type-assert attachment operations (nil-safe — attachment funcs may be absent)
		uploadFile, _ := ctx.UploadFile.(func(ctx context.Context, bucket, key string, content []byte, contentType string) error)
		listAttachments, _ := ctx.ListAttachments.(func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error))
		createAttachment, _ := ctx.CreateAttachment.(func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error))
		deleteAttachment, _ := ctx.DeleteAttachment.(func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error))
		newAttachmentID, _ := ctx.NewAttachmentID.(func() string)

		// type-assert user/role helpers (nil-safe)
		getUsersByRoleID, _ := ctx.GetUsersByRoleID.(func(ctx context.Context, roleID string) ([]roleusers.UserByRole, error))
		getDashboardData, _ := ctx.GetDashboardData.(func(ctx context.Context) (*userdashboard.DashboardData, error))
		hashPassword, _ := ctx.HashPassword.(func(password string) (string, error))
		getUserRolesMap, _ := ctx.GetUserRolesMap.(func(ctx context.Context) (map[string][]entydad.RoleBadge, error))

		// --- load labels from lyngua ---
		labels := loadBlockLabels(translations, ctx.BusinessType)

		// --- load routes (defaults + lyngua JSON overrides) ---
		routes := loadBlockRoutes(translations, ctx.BusinessType)

		// --- register modules ---

		if cfg.enableAll || cfg.client {
			if uc.Entity.Client == nil {
				return fmt.Errorf("entydad.Block: client use cases not initialized")
			}
			clientDeps := &clientmod.ModuleDeps{
				Routes:               routes.Client,
				CommonLabels:         ctx.Common,
				SharedLabels:         labels.Shared,
				Labels:               labels.Client,
				DashboardLabels:      labels.ClientDashboard,
				DashboardTitleLabels: labels.Dashboard,
				TableLabels:          ctx.Table,
				GetListPageData:      uc.Entity.Client.GetClientListPageData.Execute,
				GetInUseIDs:          refChecker.GetClientInUseIDs,
				CreateClient:         uc.Entity.Client.CreateClient.Execute,
				ReadClient:           uc.Entity.Client.ReadClient.Execute,
				UpdateClient:         uc.Entity.Client.UpdateClient.Execute,
				DeleteClient:         uc.Entity.Client.DeleteClient.Execute,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "client", id, map[string]any{"active": active})
					return err
				},
				ListPaymentTerms: func(fctx context.Context) ([]*clientmod.PaymentTermOption, error) {
					rows, err := db.ListSimple(fctx, "payment_term")
					if err != nil {
						return nil, err
					}
					opts := make([]*clientmod.PaymentTermOption, 0, len(rows))
					for _, row := range rows {
						id, _ := row["id"].(string)
						name, _ := row["name"].(string)
						if id == "" {
							continue
						}
						opts = append(opts, &clientmod.PaymentTermOption{Id: id, Name: name})
					}
					return opts, nil
				},
				ListRevenues:          db.ListSimple,
				SubscriptionAddURL:    routes.Subscription.AddURL,
				SubscriptionDetailURL: routes.Subscription.DetailURL,
				SubscriptionEditURL:   routes.Subscription.EditURL,
				SubscriptionDeleteURL: routes.Subscription.DeleteURL,
				UploadFile:            uploadFile,
				ListAttachments:       listAttachments,
				CreateAttachment:      createAttachment,
				DeleteAttachment:      deleteAttachment,
				NewID:                 newAttachmentID,
			}
			if uc.Common != nil && uc.Common.Category != nil {
				clientDeps.ListCategories = uc.Common.Category.ListCategories.Execute
			}
			if uc.Entity.ClientCategory != nil {
				clientDeps.ListClientCategories = uc.Entity.ClientCategory.ListClientCategories.Execute
				clientDeps.CreateClientCategory = uc.Entity.ClientCategory.CreateClientCategory.Execute
				clientDeps.DeleteClientCategory = uc.Entity.ClientCategory.DeleteClientCategory.Execute
			}
			if uc.Subscription != nil && uc.Subscription.Subscription != nil {
				clientDeps.ListSubscriptions = uc.Subscription.Subscription.ListSubscriptions.Execute
				clientDeps.GetSubscriptionListPageData = uc.Subscription.Subscription.GetSubscriptionListPageData.Execute
			}
			clientmod.NewModule(clientDeps).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.user {
			if uc.Entity.User == nil {
				return fmt.Errorf("entydad.Block: user use cases not initialized")
			}
			usermod.NewModule(&usermod.ModuleDeps{
				Routes:          routes.User,
				CommonLabels:    ctx.Common,
				SharedLabels:    labels.Shared,
				Labels:          labels.User,
				DashboardLabels: labels.UserDashboard,
				UserRoleLabels:  labels.UserRole,
				TableLabels:     ctx.Table,
				GetListPageData: uc.Entity.User.GetUserListPageData.Execute,
				GetUserRolesMap: getUserRolesMap,
				CreateUser:      uc.Entity.User.CreateUser.Execute,
				ReadUser:        uc.Entity.User.ReadUser.Execute,
				UpdateUser:      uc.Entity.User.UpdateUser.Execute,
				DeleteUser:      uc.Entity.User.DeleteUser.Execute,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "user", id, map[string]any{"active": active})
					return err
				},
				CreateWorkspaceUser:          uc.Entity.WorkspaceUser.CreateWorkspaceUser.Execute,
				ListWorkspaceUsers:           uc.Entity.WorkspaceUser.ListWorkspaceUsers.Execute,
				GetWorkspaceUserItemPageData: uc.Entity.WorkspaceUser.GetWorkspaceUserItemPageData.Execute,
				DefaultWorkspaceID:           getDefaultWorkspaceID(),
				CreateWorkspaceUserRole:      uc.Entity.WorkspaceUserRole.CreateWorkspaceUserRole.Execute,
				DeleteWorkspaceUserRole:      uc.Entity.WorkspaceUserRole.DeleteWorkspaceUserRole.Execute,
				ListRoles:                    uc.Entity.Role.ListRoles.Execute,
				GetDashboardData:             getDashboardData,
				HashPassword:                 hashPassword,
				UploadFile:                   uploadFile,
				ListAttachments:              listAttachments,
				CreateAttachment:             createAttachment,
				DeleteAttachment:             deleteAttachment,
				NewID:                        newAttachmentID,
			}).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.role {
			if uc.Entity.Role == nil {
				return fmt.Errorf("entydad.Block: role use cases not initialized")
			}
			rolemod.NewModule(&rolemod.ModuleDeps{
				Routes:               routes.Role,
				CommonLabels:         ctx.Common,
				SharedLabels:         labels.Shared,
				Labels:               labels.Role,
				RolePermissionLabels: labels.RolePermission,
				RoleUserLabels:       labels.RoleUser,
				TableLabels:          ctx.Table,
				GetListPageData:      uc.Entity.Role.GetRoleListPageData.Execute,
				GetInUseIDs:          refChecker.GetRoleInUseIDs,
				CreateRole:           uc.Entity.Role.CreateRole.Execute,
				ReadRole:             uc.Entity.Role.ReadRole.Execute,
				UpdateRole:           uc.Entity.Role.UpdateRole.Execute,
				DeleteRole:           uc.Entity.Role.DeleteRole.Execute,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "role", id, map[string]any{"active": active})
					return err
				},
				GetItemPageData:         uc.Entity.Role.GetRoleItemPageData.Execute,
				CreateRolePermission:    uc.Entity.RolePermission.CreateRolePermission.Execute,
				DeleteRolePermission:    uc.Entity.RolePermission.DeleteRolePermission.Execute,
				ListPermissions:         uc.Entity.Permission.ListPermissions.Execute,
				GetUsersByRoleID:        getUsersByRoleID,
				ListWorkspaceUsers:      uc.Entity.WorkspaceUser.ListWorkspaceUsers.Execute,
				CreateWorkspaceUserRole: uc.Entity.WorkspaceUserRole.CreateWorkspaceUserRole.Execute,
				DeleteWorkspaceUserRole: uc.Entity.WorkspaceUserRole.DeleteWorkspaceUserRole.Execute,
				UploadFile:              uploadFile,
				ListAttachments:         listAttachments,
				CreateAttachment:        createAttachment,
				DeleteAttachment:        deleteAttachment,
				NewID:                   newAttachmentID,
			}).RegisterRoutes(ctx.Routes)

			// Role-User search (http.HandlerFunc — uses HandleFunc, not GET)
			handleFunc(ctx.Routes, "GET", routes.Role.UsersSearchURL, roleusers.NewSearchUsersAction(&roleusers.SearchDeps{
				ListWorkspaceUsers: uc.Entity.WorkspaceUser.ListWorkspaceUsers.Execute,
			}))
		}

		if cfg.enableAll || cfg.location {
			if uc.Entity.Location == nil {
				return fmt.Errorf("entydad.Block: location use cases not initialized")
			}
			locationmod.NewModule(&locationmod.ModuleDeps{
				Routes:          routes.Location,
				CommonLabels:    ctx.Common,
				SharedLabels:    labels.Shared,
				Labels:          labels.Location,
				TableLabels:     ctx.Table,
				GetListPageData: uc.Entity.Location.GetLocationListPageData.Execute,
				GetInUseIDs:     refChecker.GetLocationInUseIDs,
				CreateLocation:  uc.Entity.Location.CreateLocation.Execute,
				ReadLocation:    uc.Entity.Location.ReadLocation.Execute,
				UpdateLocation:  uc.Entity.Location.UpdateLocation.Execute,
				DeleteLocation:  uc.Entity.Location.DeleteLocation.Execute,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "location", id, map[string]any{"active": active})
					return err
				},
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
				NewID:            newAttachmentID,
			}).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.locationArea {
			crudDB, hasCRUD := db.(CRUDSource)
			if !hasCRUD {
				log.Println("entydad.Block: warning: DB does not implement CRUDSource — skipping location_area module")
			} else {
				locationareamod.NewModule(&locationareamod.ModuleDeps{
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
					GetInUseIDs: func(fctx context.Context, ids []string) (map[string]bool, error) {
						return nil, nil
					},
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

		if cfg.enableAll || cfg.permission {
			if uc.Entity.Permission == nil {
				return fmt.Errorf("entydad.Block: permission use cases not initialized")
			}
			permissionmod.NewModule(&permissionmod.ModuleDeps{
				Routes:           routes.Permission,
				CommonLabels:     ctx.Common,
				SharedLabels:     labels.Shared,
				Labels:           labels.Permission,
				TableLabels:      ctx.Table,
				GetListPageData:  uc.Entity.Permission.GetPermissionListPageData.Execute,
				CreatePermission: uc.Entity.Permission.CreatePermission.Execute,
				ReadPermission:   uc.Entity.Permission.ReadPermission.Execute,
				UpdatePermission: uc.Entity.Permission.UpdatePermission.Execute,
				DeletePermission: uc.Entity.Permission.DeletePermission.Execute,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "permission", id, map[string]any{"active": active})
					return err
				},
			}).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.workspace {
			if uc.Entity.Workspace == nil {
				return fmt.Errorf("entydad.Block: workspace use cases not initialized")
			}
			workspacemod.NewModule(&workspacemod.ModuleDeps{
				Routes:          routes.Workspace,
				CommonLabels:    ctx.Common,
				SharedLabels:    labels.Shared,
				Labels:          labels.Workspace,
				TableLabels:     ctx.Table,
				GetListPageData: uc.Entity.Workspace.GetWorkspaceListPageData.Execute,
				CreateWorkspace: uc.Entity.Workspace.CreateWorkspace.Execute,
				ReadWorkspace:   uc.Entity.Workspace.ReadWorkspace.Execute,
				UpdateWorkspace: uc.Entity.Workspace.UpdateWorkspace.Execute,
				DeleteWorkspace: uc.Entity.Workspace.DeleteWorkspace.Execute,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "workspace", id, map[string]any{"active": active})
					return err
				},
			}).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.supplier {
			supplierDeps := &suppliermod.ModuleDeps{
				Routes:       routes.Supplier,
				CommonLabels: ctx.Common,
				SharedLabels: labels.Shared,
				Labels:       labels.Supplier,
				TableLabels:  ctx.Table,
				GetInUseIDs: func(fctx context.Context, ids []string) (map[string]bool, error) {
					return nil, nil
				},
				SetActive: func(fctx context.Context, id string, active bool) error {
					status := "blocked"
					if active {
						status = "active"
					}
					_, err := db.Update(fctx, "supplier", id, map[string]any{"active": active, "status": status})
					return err
				},
				ListPaymentTerms: func(fctx context.Context) ([]*suppliermod.PaymentTermOption, error) {
					rows, err := db.ListSimple(fctx, "payment_term")
					if err != nil {
						return nil, err
					}
					opts := make([]*suppliermod.PaymentTermOption, 0, len(rows))
					for _, row := range rows {
						id, _ := row["id"].(string)
						name, _ := row["name"].(string)
						if id == "" {
							continue
						}
						opts = append(opts, &suppliermod.PaymentTermOption{Id: id, Name: name})
					}
					return opts, nil
				},
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
				NewID:            newAttachmentID,
			}
			if uc.Entity.Supplier != nil {
				supplierDeps.GetListPageData = uc.Entity.Supplier.GetSupplierListPageData.Execute
				supplierDeps.CreateSupplier = uc.Entity.Supplier.CreateSupplier.Execute
				supplierDeps.ReadSupplier = uc.Entity.Supplier.ReadSupplier.Execute
				supplierDeps.UpdateSupplier = uc.Entity.Supplier.UpdateSupplier.Execute
				supplierDeps.DeleteSupplier = uc.Entity.Supplier.DeleteSupplier.Execute
			}
			if uc.Expenditure != nil && uc.Expenditure.PurchaseOrder != nil && uc.Expenditure.PurchaseOrder.ListPurchaseOrders != nil {
				supplierDeps.ListPurchaseOrders = uc.Expenditure.PurchaseOrder.ListPurchaseOrders.Execute
			}
			suppliermod.NewModule(supplierDeps).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.clientTag {
			clienttagDeps := &clienttagmod.ModuleDeps{
				Routes:       routes.ClientTag,
				Labels:       labels.ClientTag,
				SharedLabels: labels.Shared,
				CommonLabels: ctx.Common,
				TableLabels:  ctx.Table,
				GetInUseIDs:  refChecker.GetCategoryInUseIDs,
			}
			if uc.Common != nil && uc.Common.Category != nil {
				clienttagDeps.ListCategories = uc.Common.Category.ListCategories.Execute
				clienttagDeps.CreateCategory = uc.Common.Category.CreateCategory.Execute
				clienttagDeps.ReadCategory = uc.Common.Category.ReadCategory.Execute
				clienttagDeps.UpdateCategory = uc.Common.Category.UpdateCategory.Execute
				clienttagDeps.DeleteCategory = uc.Common.Category.DeleteCategory.Execute
			}
			if uc.Entity.ClientCategory != nil {
				clienttagDeps.ListClientCategories = uc.Entity.ClientCategory.ListClientCategories.Execute
			}
			clienttagmod.NewModule(clienttagDeps).RegisterRoutes(ctx.Routes)
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
				paymenttermmod.NewModule(&paymenttermmod.ModuleDeps{
					Routes:            routes.PaymentTerm,
					CommonLabels:      ctx.Common,
					SharedLabels:      labels.Shared,
					Labels:            labels.PaymentTerm,
					TableLabels:       ctx.Table,
					GetListPageData:   ptRepo.GetPaymentTermListPageData,
					CreatePaymentTerm: ptRepo.CreatePaymentTerm,
					ReadPaymentTerm:   ptRepo.ReadPaymentTerm,
					UpdatePaymentTerm: ptRepo.UpdatePaymentTerm,
					DeletePaymentTerm: ptRepo.DeletePaymentTerm,
				}).RegisterRoutes(ctx.Routes)
			}
		}

		log.Println("  ✓ Entity domain initialized (entydad.Block)")
		return nil
	}
}

// UpdateableSource extends entydad.DataSource with the Update method that
// SetActive closures need. espyna's DatabaseAdapter satisfies this interface.
type UpdateableSource interface {
	entydad.DataSource
	Update(ctx context.Context, collection, id string, data map[string]any) (map[string]any, error)
}

// CRUDSource extends UpdateableSource with Create, Read, and Delete operations.
// espyna's DatabaseAdapter satisfies this interface. Used by simpler entities
// (e.g. LocationArea) that do not yet have dedicated proto service use-cases.
type CRUDSource interface {
	UpdateableSource
	Create(ctx context.Context, collection string, data map[string]any) (map[string]any, error)
	Read(ctx context.Context, collection, id string) (map[string]any, error)
	Delete(ctx context.Context, collection, id string) error
}

// getDefaultWorkspaceID returns the default workspace ID from the environment,
// falling back to "default-workspace" if the env var is not set.
func getDefaultWorkspaceID() string {
	if v := os.Getenv("DEFAULT_WORKSPACE_ID"); v != "" {
		return v
	}
	return "default-workspace"
}

// ---------------------------------------------------------------------------
// Internal helpers: typed label/route loaders
// ---------------------------------------------------------------------------

// blockLabels holds the subset of entydad label structs needed by Block().
type blockLabels struct {
	Shared          entydad.SharedLabels
	Dashboard       entydad.DashboardLabels
	Client          entydad.ClientLabels
	ClientDashboard entydad.ClientDashboardLabels
	ClientTag       entydad.ClientTagLabels
	PaymentTerm     entydad.PaymentTermLabels
	User            entydad.UserLabels
	UserDashboard   entydad.UserDashboardLabels
	UserRole        entydad.UserRoleLabels
	RoleUser        entydad.RoleUserLabels
	Role            entydad.RoleLabels
	RolePermission  entydad.RolePermissionLabels
	Location        entydad.LocationLabels
	LocationArea    entydad.LocationAreaLabels
	Permission      entydad.PermissionLabels
	Workspace       entydad.WorkspaceLabels
	Supplier        entydad.SupplierLabels
}

// blockRoutes holds the subset of entydad route structs needed by Block().
type blockRoutes struct {
	Client       entydad.ClientRoutes
	ClientTag    entydad.ClientTagRoutes
	PaymentTerm  entydad.PaymentTermRoutes
	Subscription centymo.SubscriptionRoutes
	User         entydad.UserRoutes
	Role         entydad.RoleRoutes
	Location     entydad.LocationRoutes
	LocationArea entydad.LocationAreaRoutes
	Permission   entydad.PermissionRoutes
	Workspace    entydad.WorkspaceRoutes
	Supplier     entydad.SupplierRoutes
}

// loadBlockLabels loads all entydad typed label structs from lyngua.
// Mirrors the entydad section of translations.go in service-admin/retail-admin.
func loadBlockLabels(t *lynguaV1.TranslationProvider, businessType string) blockLabels {
	l := blockLabels{}

	_ = t.LoadPathIfExists("en", businessType, "dashboard.json", "", &l.Dashboard)

	if err := t.LoadPath("en", businessType, "client.json", "client", &l.Client); err != nil {
		log.Printf("entydad.Block: warning: failed to load client labels: %v", err)
	}
	_ = t.LoadPathIfExists("en", businessType, "client.json", "client.dashboard", &l.ClientDashboard)
	_ = t.LoadPathIfExists("en", businessType, "client_tag.json", "", &l.ClientTag)
	_ = t.LoadPathIfExists("en", businessType, "payment_term.json", "paymentTerm", &l.PaymentTerm)

	if err := t.LoadPath("en", businessType, "user.json", "", &l.User); err != nil {
		log.Printf("entydad.Block: warning: failed to load user labels: %v", err)
	}
	_ = t.LoadPathIfExists("en", businessType, "user.json", "user.dashboard", &l.UserDashboard)

	if err := t.LoadPath("en", businessType, "role.json", "", &l.Role); err != nil {
		log.Printf("entydad.Block: warning: failed to load role labels: %v", err)
	}
	if err := t.LoadPath("en", businessType, "location.json", "", &l.Location); err != nil {
		log.Printf("entydad.Block: warning: failed to load location labels: %v", err)
	}
	l.LocationArea = entydad.DefaultLocationAreaLabels()
	_ = t.LoadPathIfExists("en", businessType, "location_area.json", "", &l.LocationArea)
	if err := t.LoadPath("en", businessType, "permission.json", "", &l.Permission); err != nil {
		log.Printf("entydad.Block: warning: failed to load permission labels: %v", err)
	}
	if err := t.LoadPath("en", businessType, "role_permission.json", "", &l.RolePermission); err != nil {
		log.Printf("entydad.Block: warning: failed to load role_permission labels: %v", err)
	}
	if err := t.LoadPath("en", businessType, "user_role.json", "", &l.UserRole); err != nil {
		log.Printf("entydad.Block: warning: failed to load user_role labels: %v", err)
	}
	if err := t.LoadPath("en", businessType, "role_user.json", "", &l.RoleUser); err != nil {
		log.Printf("entydad.Block: warning: failed to load role_user labels: %v", err)
	}
	if err := t.LoadPath("en", businessType, "workspace.json", "", &l.Workspace); err != nil {
		log.Printf("entydad.Block: warning: failed to load workspace labels: %v", err)
	}
	if err := t.LoadPath("en", businessType, "supplier.json", "supplier", &l.Supplier); err != nil {
		log.Printf("entydad.Block: warning: failed to load supplier labels: %v", err)
	}
	if err := t.LoadPath("en", businessType, "shared.json", "", &l.Shared); err != nil {
		log.Printf("entydad.Block: warning: failed to load shared labels: %v", err)
	}

	return l
}

// loadBlockRoutes loads all entydad route configs with lyngua JSON overrides.
// Mirrors the entydad section of route_config.go in service-admin/retail-admin.
func loadBlockRoutes(t *lynguaV1.TranslationProvider, businessType string) blockRoutes {
	r := blockRoutes{}

	r.Client = entydad.DefaultClientRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "client", &r.Client)

	r.ClientTag = entydad.DefaultClientTagRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "client_tag", &r.ClientTag)

	r.PaymentTerm = entydad.DefaultPaymentTermRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "payment_term", &r.PaymentTerm)

	r.Subscription = centymo.DefaultSubscriptionRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "subscription", &r.Subscription)

	r.User = entydad.DefaultUserRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "user", &r.User)

	r.Role = entydad.DefaultRoleRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "role", &r.Role)

	r.Location = entydad.DefaultLocationRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "location", &r.Location)

	r.LocationArea = entydad.DefaultLocationAreaRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "location_area", &r.LocationArea)

	r.Permission = entydad.DefaultPermissionRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "permission", &r.Permission)

	r.Workspace = entydad.DefaultWorkspaceRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "workspace", &r.Workspace)

	r.Supplier = entydad.DefaultSupplierRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "supplier", &r.Supplier)

	return r
}

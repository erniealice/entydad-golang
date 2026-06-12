// Package block — typed label and route loaders for entydad.Block.
//
// This file owns:
//   - blockLabels struct: all entydad label types needed by Block().
//   - blockRoutes struct: all entydad route types needed by Block().
//   - loadBlockLabels: populates blockLabels from lyngua per business type.
//   - loadBlockRoutes: populates blockRoutes from lyngua per business type.
//
// Adding a new module means:
//  1. Add a field to blockLabels and/or blockRoutes.
//  2. Wire the load call in loadBlockLabels and/or loadBlockRoutes.
//  3. Wire the dependency in Block() (block.go).
//
// Nothing else in this file is load-bearing — it is a flat list by design
// so a reader can scan every label/route binding in one scroll.
package block

import (
	"log"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/entydad-golang"
	entity "github.com/erniealice/entydad-golang/domain/entity"
	tax "github.com/erniealice/entydad-golang/domain/tax"
	convmodel "github.com/erniealice/hybra-golang/views/conversation/model"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// blockLabels holds the subset of entydad label structs needed by Block().
// Migrated entity labels resolve through the entity facade (entity.*);
// shared/auth/admin/tax/conversation leftovers stay on the root package (entydad.*).
type blockLabels struct {
	Shared            entydad.SharedLabels
	Dashboard         entydad.DashboardLabels
	Admin             entydad.AdminDashboardLabels
	Client            entity.ClientLabels
	ClientDashboard   entity.ClientDashboardLabels
	ClientTag         entity.ClientTagLabels
	SupplierTag       entity.SupplierTagLabels
	PaymentTerm       entity.PaymentTermLabels
	User              entity.UserLabels
	UserDashboard     entity.UserDashboardLabels
	UserRole          entity.UserRoleLabels
	RoleUser          entity.RoleUserLabels
	Role              entity.RoleLabels
	RolePermission    entity.RolePermissionLabels
	Location          entity.LocationLabels
	LocationArea      entity.LocationAreaLabels
	Permission        entity.PermissionLabels
	Workspace         entity.WorkspaceLabels
	WorkspaceUser     entity.WorkspaceUserLabels
	WorkspaceUserRole entity.WorkspaceUserRoleLabels
	Supplier          entity.SupplierLabels
	SupplierDashboard entity.SupplierDashboardLabels
	TaxRegistration   tax.TaxRegistrationLabels
	Conversation      convmodel.ConversationLabels
	ConversationPost  convmodel.ConversationPostLabels
}

// blockRoutes holds the subset of entydad route structs needed by Block().
type blockRoutes struct {
	Admin               entydad.AdminDashboardRoutes
	Client              entity.ClientRoutes
	ClientTag           entity.ClientTagRoutes
	SupplierTag         entity.SupplierTagRoutes
	PaymentTerm         entity.PaymentTermRoutes
	SupplierPaymentTerm entity.SupplierPaymentTermRoutes
	Subscription        centymo.SubscriptionRoutes
	PriceSchedule       centymo.PriceScheduleRoutes
	User                entity.UserRoutes
	Role                entity.RoleRoutes
	Location            entity.LocationRoutes
	LocationArea        entity.LocationAreaRoutes
	Permission          entity.PermissionRoutes
	Workspace           entity.WorkspaceRoutes
	WorkspaceUser       entity.WorkspaceUserRoutes
	WorkspaceUserRole   entity.WorkspaceUserRoleRoutes
	Supplier            entity.SupplierRoutes
	TaxRegistration     tax.TaxRegistrationRoutes
	Conversation        convmodel.ConversationRoutes
}

// loadBlockLabels loads all entydad typed label structs from lyngua.
// Mirrors the entydad section of translations.go in service-admin/retail-admin.
func loadBlockLabels(t *lynguaV1.TranslationProvider, businessType string) blockLabels {
	l := blockLabels{}

	_ = t.LoadPathIfExists("en", businessType, "dashboard.json", "", &l.Dashboard)
	_ = t.LoadPathIfExists("en", businessType, "admin.json", "admin.dashboard", &l.Admin)

	if err := t.LoadPath("en", businessType, "client.json", "client", &l.Client); err != nil {
		log.Printf("entydad.Block: warning: failed to load client labels: %v", err)
	}
	_ = t.LoadPathIfExists("en", businessType, "client.json", "client.dashboard", &l.ClientDashboard)
	_ = t.LoadPathIfExists("en", businessType, "client_tag.json", "", &l.ClientTag)
	_ = t.LoadPathIfExists("en", businessType, "supplier_tag.json", "", &l.SupplierTag)
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
	l.LocationArea = entity.DefaultLocationAreaLabels()
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
	_ = t.LoadPathIfExists("en", businessType, "workspace_user.json", "", &l.WorkspaceUser)
	_ = t.LoadPathIfExists("en", businessType, "workspace_user_role.json", "workspace_user_role", &l.WorkspaceUserRole)
	if err := t.LoadPath("en", businessType, "supplier.json", "supplier", &l.Supplier); err != nil {
		log.Printf("entydad.Block: warning: failed to load supplier labels: %v", err)
	}
	_ = t.LoadPathIfExists("en", businessType, "supplier.json", "supplier.dashboard", &l.SupplierDashboard)
	if err := t.LoadPath("en", businessType, "shared.json", "", &l.Shared); err != nil {
		log.Printf("entydad.Block: warning: failed to load shared labels: %v", err)
	}

	l.TaxRegistration = tax.DefaultTaxRegistrationLabels()
	_ = t.LoadPathIfExists("en", businessType, "tax_registration.json", "", &l.TaxRegistration)

	// Conversation (secure messaging — Plan-4). Optional on non-messaging
	// business types — LoadPathIfExists (no boot warning when absent).
	l.Conversation = convmodel.DefaultConversationLabels()
	_ = t.LoadPathIfExists("en", businessType, "conversation.json", "conversation", &l.Conversation)
	l.ConversationPost = convmodel.DefaultConversationPostLabels()
	_ = t.LoadPathIfExists("en", businessType, "conversation_post.json", "conversationPost", &l.ConversationPost)

	return l
}

// loadBlockRoutes loads all entydad route configs with lyngua JSON overrides.
// Mirrors the entydad section of route_config.go in service-admin/retail-admin.
func loadBlockRoutes(t *lynguaV1.TranslationProvider, businessType string) blockRoutes {
	r := blockRoutes{}

	r.Admin = entydad.DefaultAdminDashboardRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "admin", &r.Admin)

	r.Client = entity.DefaultClientRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "client", &r.Client)

	r.ClientTag = entity.DefaultClientTagRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "client_tag", &r.ClientTag)

	r.SupplierTag = entity.DefaultSupplierTagRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "supplier_tag", &r.SupplierTag)

	r.PaymentTerm = entity.DefaultPaymentTermRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "payment_term", &r.PaymentTerm)

	r.SupplierPaymentTerm = entity.DefaultSupplierPaymentTermRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "supplier_payment_term", &r.SupplierPaymentTerm)

	r.Subscription = centymo.DefaultSubscriptionRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "subscription", &r.Subscription)

	r.PriceSchedule = centymo.DefaultPriceScheduleRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "price_schedule", &r.PriceSchedule)

	r.User = entity.DefaultUserRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "user", &r.User)

	r.Role = entity.DefaultRoleRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "role", &r.Role)

	r.Location = entity.DefaultLocationRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "location", &r.Location)

	r.LocationArea = entity.DefaultLocationAreaRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "location_area", &r.LocationArea)

	r.Permission = entity.DefaultPermissionRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "permission", &r.Permission)

	r.Workspace = entity.DefaultWorkspaceRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "workspace", &r.Workspace)

	r.WorkspaceUser = entity.DefaultWorkspaceUserRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "workspace_user", &r.WorkspaceUser)

	r.WorkspaceUserRole = entity.DefaultWorkspaceUserRoleRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "workspace_user_role", &r.WorkspaceUserRole)

	r.Supplier = entity.DefaultSupplierRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "supplier", &r.Supplier)

	r.TaxRegistration = tax.DefaultTaxRegistrationRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "tax_registration", &r.TaxRegistration)

	r.Conversation = convmodel.DefaultConversationRoutes()
	_ = t.LoadPathIfExists("en", businessType, "route.json", "conversation", &r.Conversation)

	return r
}

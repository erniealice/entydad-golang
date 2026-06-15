package block

// workspace_loader.go — relocated from app
// internal/infrastructure/input/http/workspace_loader.go (Wave B D2a). The
// proto-backed sidebar workspace switcher loader. Imports workspacepb (LEGAL in
// entydad; ILLEGAL in espyna) and satisfies consumerhttp.WorkspaceLoader
// structurally. The entydad EngineBlock constructs it and sets
// ctx.WorkspaceLoader so the espyna Server.finalizeHTTPAdapter (Build() path)
// can pick it up; on the live Handler() path the app still constructs its own.

import (
	"context"
	"log"

	"github.com/erniealice/espyna-golang/consumer"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	"github.com/erniealice/pyeza-golang/types"
)

// workspaceListUseCase is the minimal interface for listing user workspaces.
type workspaceListUseCase interface {
	Execute(ctx context.Context, req *workspacepb.ListUserWorkspacesRequest) (*workspacepb.ListUserWorkspacesResponse, error)
}

// DBWorkspaceLoader loads workspace data for the current user from the database.
// It uses the ListUserWorkspaces use case via the espyna workspace domain service.
type DBWorkspaceLoader struct {
	useCase workspaceListUseCase
}

// NewDBWorkspaceLoader creates a WorkspaceLoader backed by the given use case.
func NewDBWorkspaceLoader(uc workspaceListUseCase) *DBWorkspaceLoader {
	return &DBWorkspaceLoader{useCase: uc}
}

// LoadWorkspaces queries the database for all workspaces accessible to the current user
// and identifies the currently active workspace from the session context.
func (l *DBWorkspaceLoader) LoadWorkspaces(ctx context.Context) ([]types.SidebarWorkspace, types.SidebarWorkspace) {
	userID := consumer.GetUserIDFromContext(ctx)
	if userID == "" {
		userID = consumer.ExtractUserIDFromContext(ctx)
	}
	if userID == "" {
		return nil, types.SidebarWorkspace{}
	}

	wsID := consumer.GetWorkspaceIDFromContext(ctx)

	resp, err := l.useCase.Execute(ctx, &workspacepb.ListUserWorkspacesRequest{UserId: userID})
	if err != nil {
		log.Printf("WorkspaceLoader: failed to list workspaces for user %s: %v", userID, err)
		return nil, types.SidebarWorkspace{}
	}
	if !resp.Success {
		return nil, types.SidebarWorkspace{}
	}

	var current types.SidebarWorkspace
	all := make([]types.SidebarWorkspace, 0, len(resp.Workspaces))
	for _, ws := range resp.Workspaces {
		sw := types.SidebarWorkspace{ID: ws.WorkspaceId, Name: ws.WorkspaceName}
		all = append(all, sw)
		if ws.WorkspaceId == wsID {
			current = sw
		}
	}

	return all, current
}

// IsEnabled returns true — workspace loading is always enabled when this loader exists.
func (l *DBWorkspaceLoader) IsEnabled() bool {
	return l.useCase != nil
}

package detail

import (
	"context"
	"log"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
)

// loadAttachments populates the AttachmentTable and AttachmentUploadURL on PageData.
func loadAttachments(ctx context.Context, deps *DetailViewDeps, id string, pageData *PageData) {
	cfg := attachmentConfig(deps, id)

	if deps.ListAttachments != nil {
		resp, err := deps.ListAttachments(ctx, "client", id)
		if err != nil {
			log.Printf("Failed to list attachments for client %s: %v", id, err)
		}
		var items []*attachmentpb.Attachment
		if resp != nil {
			items = resp.GetData()
		}
		pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
	}

	pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)
}

// attachmentConfig builds the shared attachment.Config for the client entity.
func attachmentConfig(deps *DetailViewDeps, id string) *attachment.Config {
	return &attachment.Config{
		EntityType:       "client",
		BucketName:       "attachments",
		UploadURL:        deps.Routes.AttachmentUploadURL,
		DeleteURL:        deps.Routes.AttachmentDeleteURL,
		RedirectURL:      route.ResolveURL(deps.Routes.DetailURL, "id", id) + "?tab=attachments",
		Labels:           attachment.DefaultLabels(),
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		NewID:            deps.NewAttachmentID,
		UploadFile:       deps.UploadFile,
		ListAttachments:  deps.ListAttachments,
		CreateAttachment: deps.CreateAttachment,
		DeleteAttachment: deps.DeleteAttachment,
	}
}

func NewAttachmentUploadAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		cfg := attachmentConfig(deps, id)
		return attachment.NewUploadAction(cfg).Handle(ctx, viewCtx)
	})
}

func NewAttachmentDeleteAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		cfg := attachmentConfig(deps, id)
		return attachment.NewDeleteAction(cfg).Handle(ctx, viewCtx)
	})
}

package detail

import (
	"context"
	"log"

	"github.com/erniealice/fycha-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
)

// loadAttachments populates the AttachmentTable and AttachmentUploadURL on PageData.
func loadAttachments(ctx context.Context, deps *Deps, id string, pageData *PageData) {
	cfg := attachmentConfig(deps, id)

	if deps.ListAttachments != nil {
		resp, err := deps.ListAttachments(ctx, "location", id)
		if err != nil {
			log.Printf("Failed to list attachments for location %s: %v", id, err)
		}
		var items []*attachmentpb.Attachment
		if resp != nil {
			items = resp.GetData()
		}
		pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
	}

	pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)
}

// attachmentConfig builds the shared attachment.Config for the location entity.
func attachmentConfig(deps *Deps, id string) *attachment.Config {
	return &attachment.Config{
		EntityType:       "location",
		BucketName:       "attachments",
		UploadURL:        deps.Routes.AttachmentUploadURL,
		DeleteURL:        deps.Routes.AttachmentDeleteURL,
		RedirectURL:      route.ResolveURL(deps.Routes.DetailURL, "id", id) + "?tab=attachments",
		Labels:           attachment.DefaultLabels(),
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		NewID:            deps.NewID,
		UploadFile:       deps.UploadFile,
		ListAttachments:  deps.ListAttachments,
		CreateAttachment: deps.CreateAttachment,
		DeleteAttachment: deps.DeleteAttachment,
	}
}

// NewAttachmentUploadAction creates the upload handler for location attachments.
func NewAttachmentUploadAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		cfg := attachmentConfig(deps, id)
		return attachment.NewUploadAction(cfg).Handle(ctx, viewCtx)
	})
}

// NewAttachmentDeleteAction creates the delete handler for location attachments.
func NewAttachmentDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		cfg := attachmentConfig(deps, id)
		return attachment.NewDeleteAction(cfg).Handle(ctx, viewCtx)
	})
}

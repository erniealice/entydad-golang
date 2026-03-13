package detail

import (
	"github.com/erniealice/fycha-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/view"
)

func attachmentConfig(deps *Deps) *attachment.Config {
	return &attachment.Config{
		EntityType:       "role",
		BucketName:       "attachments",
		UploadURL:        deps.Routes.AttachmentUploadURL,
		DeleteURL:        deps.Routes.AttachmentDeleteURL,
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

func NewAttachmentUploadAction(deps *Deps) view.View {
	return attachment.NewUploadAction(attachmentConfig(deps))
}

func NewAttachmentDeleteAction(deps *Deps) view.View {
	return attachment.NewDeleteAction(attachmentConfig(deps))
}

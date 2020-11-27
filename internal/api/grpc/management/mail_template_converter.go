package management

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func mailTemplateRequestToModel(MtemplateailTemplate *management.MailTemplateUpdate) *iam_model.MailTemplate {
	return &iam_model.MailTemplate{
		Template: MtemplateailTemplate.Template,
	}
}

func mailTemplateFromModel(MtemplateailTemplate *iam_model.MailTemplate) *management.MailTemplate {
	creationDate, err := ptypes.TimestampProto(MtemplateailTemplate.CreationDate)
	logging.Log("MANAG-ULKZ6").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(MtemplateailTemplate.ChangeDate)
	logging.Log("MANAG-451rI").OnError(err).Debug("date parse failed")

	return &management.MailTemplate{
		Template:     MtemplateailTemplate.Template,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}

func mailTemplateViewFromModel(MtemplateailTemplate *iam_model.MailTemplateView) *management.MailTemplateView {
	creationDate, err := ptypes.TimestampProto(MtemplateailTemplate.CreationDate)
	logging.Log("MANAG-koQnB").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(MtemplateailTemplate.ChangeDate)
	logging.Log("MANAG-ToDhD").OnError(err).Debug("date parse failed")

	return &management.MailTemplateView{
		Default:      MtemplateailTemplate.Default,
		Template:     MtemplateailTemplate.Template,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}

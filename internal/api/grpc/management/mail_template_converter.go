package management

import (
	"fmt"

	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func mailTemplateRequestToModel(mailTemplate *management.MailTemplateUpdate) *iam_model.MailTemplate {
	return &iam_model.MailTemplate{
		Template: mailTemplate.Template,
	}
}

func mailTemplateFromModel(mailTemplate *iam_model.MailTemplate) *management.MailTemplate {
	creationDate, err := ptypes.TimestampProto(mailTemplate.CreationDate)
	logging.Log("MANAG-ULKZ6").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(mailTemplate.ChangeDate)
	logging.Log("MANAG-451rI").OnError(err).Debug("date parse failed")
	fmt.Println(string(mailTemplate.Template))
	return &management.MailTemplate{
		Template:     mailTemplate.Template,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}

func mailTemplateViewFromModel(mailTemplate *iam_model.MailTemplateView) *management.MailTemplateView {
	creationDate, err := ptypes.TimestampProto(mailTemplate.CreationDate)
	logging.Log("MANAG-koQnB").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(mailTemplate.ChangeDate)
	logging.Log("MANAG-ToDhD").OnError(err).Debug("date parse failed")

	return &management.MailTemplateView{
		Default:      mailTemplate.Default,
		Template:     mailTemplate.Template,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}

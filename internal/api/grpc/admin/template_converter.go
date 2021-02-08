package admin

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func templateToModel(policy *admin.DefaultMailTemplateUpdate) *iam_model.MailTemplate {
	return &iam_model.MailTemplate{
		Template: policy.Template,
	}
}

func templateFromModel(policy *iam_model.MailTemplate) *admin.DefaultMailTemplate {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("ADMIN-CAA7T").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("ADMIN-H52Zx").OnError(err).Debug("date parse failed")

	return &admin.DefaultMailTemplate{
		Template:     policy.Template,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}

func templateViewFromModel(policy *iam_model.MailTemplateView) *admin.DefaultMailTemplateView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("ADMIN-yWFs5").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("ADMIN-JRpIO").OnError(err).Debug("date parse failed")

	return &admin.DefaultMailTemplateView{
		Template:     policy.Template,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}

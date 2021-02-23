package management

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mailTemplateRequestToDomain(mailTemplate *management.MailTemplateUpdate) *domain.MailTemplate {
	return &domain.MailTemplate{
		Template: mailTemplate.Template,
	}
}

func mailTemplateFromDomain(mailTemplate *domain.MailTemplate) *management.MailTemplate {
	return &management.MailTemplate{
		Template:     mailTemplate.Template,
		CreationDate: timestamppb.New(mailTemplate.CreationDate),
		ChangeDate:   timestamppb.New(mailTemplate.ChangeDate),
	}
}

func mailTemplateViewFromModel(mailTemplate *iam_model.MailTemplateView) *management.MailTemplateView {
	return &management.MailTemplateView{
		Default:      mailTemplate.Default,
		Template:     mailTemplate.Template,
		CreationDate: timestamppb.New(mailTemplate.CreationDate),
		ChangeDate:   timestamppb.New(mailTemplate.ChangeDate),
	}
}

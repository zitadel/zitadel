package admin

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func templateToDomain(policy *admin.DefaultMailTemplateUpdate) *domain.MailTemplate {
	return &domain.MailTemplate{
		Template: policy.Template,
	}
}

func templateFromDomain(policy *domain.MailTemplate) *admin.DefaultMailTemplate {
	return &admin.DefaultMailTemplate{
		Template:     policy.Template,
		CreationDate: timestamppb.New(policy.CreationDate),
		ChangeDate:   timestamppb.New(policy.ChangeDate),
	}
}

func templateViewFromModel(policy *iam_model.MailTemplateView) *admin.DefaultMailTemplateView {
	return &admin.DefaultMailTemplateView{
		Template:     policy.Template,
		CreationDate: timestamppb.New(policy.CreationDate),
		ChangeDate:   timestamppb.New(policy.ChangeDate),
	}
}

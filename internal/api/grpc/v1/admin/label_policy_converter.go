package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/interna/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func labelPolicyToDomain(policy *admin.DefaultLabelPolicyUpdate) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:   policy.PrimaryColor,
		SecondaryColor: policy.SecondaryColor,
	}
}

func labelPolicyFromDomain(policy *domain.LabelPolicy) *admin.DefaultLabelPolicy {
	return &admin.DefaultLabelPolicy{
		PrimaryColor:   policy.PrimaryColor,
		SecondaryColor: policy.SecondaryColor,
		ChangeDate:     timestamppb.New(policy.ChangeDate),
	}
}

func labelPolicyViewFromModel(policy *iam_model.LabelPolicyView) *admin.DefaultLabelPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("ADMIN-zMnlF").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("ADMIN-Vhvfp").OnError(err).Debug("date parse failed")

	return &admin.DefaultLabelPolicyView{
		PrimaryColor:   policy.PrimaryColor,
		SecondaryColor: policy.SecondaryColor,
		CreationDate:   creationDate,
		ChangeDate:     changeDate,
	}
}

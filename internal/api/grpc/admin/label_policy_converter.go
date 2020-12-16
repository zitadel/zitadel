package admin

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func labelPolicyToModel(policy *admin.DefaultLabelPolicyUpdate) *iam_model.LabelPolicy {
	return &iam_model.LabelPolicy{
		PrimaryColor:   policy.PrimaryColor,
		SecondaryColor: policy.SecondaryColor,
	}
}

func labelPolicyFromModel(policy *iam_model.LabelPolicy) *admin.DefaultLabelPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("ADMIN-QwQG9").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("ADMIN-mAgcI").OnError(err).Debug("date parse failed")

	return &admin.DefaultLabelPolicy{
		PrimaryColor:   policy.PrimaryColor,
		SecondaryColor: policy.SecondaryColor,
		CreationDate:   creationDate,
		ChangeDate:     changeDate,
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

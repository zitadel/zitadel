package management

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func labelPolicyRequestToModel(policy *management.LabelPolicyRequest) *iam_model.LabelPolicy {
	return &iam_model.LabelPolicy{
		SecondaryColor:      policy.SecondaryColor,
		PrimaryColor:        policy.PrimaryColor,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
	}
}

func labelPolicyFromModel(policy *iam_model.LabelPolicy) *management.LabelPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-2Fsm8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-3Flo0").OnError(err).Debug("date parse failed")

	return &management.LabelPolicy{
		SecondaryColor:      policy.SecondaryColor,
		PrimaryColor:        policy.PrimaryColor,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
		CreationDate:        creationDate,
		ChangeDate:          changeDate,
	}
}

func labelPolicyViewFromModel(policy *iam_model.LabelPolicyView) *management.LabelPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-5Tsm8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-8dJgs").OnError(err).Debug("date parse failed")

	return &management.LabelPolicyView{
		Default:             policy.Default,
		SecondaryColor:      policy.SecondaryColor,
		PrimaryColor:        policy.PrimaryColor,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
		CreationDate:        creationDate,
		ChangeDate:          changeDate,
	}
}

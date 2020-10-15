package admin

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func passwordAgePolicyToModel(policy *admin.DefaultPasswordAgePolicyRequest) *iam_model.PasswordAgePolicy {
	return &iam_model.PasswordAgePolicy{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
	}
}

func passwordAgePolicyFromModel(policy *iam_model.PasswordAgePolicy) *admin.DefaultPasswordAgePolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-mH9os").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-3tGs9").OnError(err).Debug("date parse failed")

	return &admin.DefaultPasswordAgePolicy{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		CreationDate:   creationDate,
		ChangeDate:     changeDate,
	}
}

func passwordAgePolicyViewFromModel(policy *iam_model.PasswordAgePolicyView) *admin.DefaultPasswordAgePolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-2Gs9o").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-8Hjss").OnError(err).Debug("date parse failed")

	return &admin.DefaultPasswordAgePolicyView{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		CreationDate:   creationDate,
		ChangeDate:     changeDate,
	}
}

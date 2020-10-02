package management

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func passwordAgePolicyAddToModel(policy *management.PasswordAgePolicyAdd) *iam_model.PasswordAgePolicy {
	return &iam_model.PasswordAgePolicy{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
	}
}
func passwordAgePolicyToModel(policy *management.PasswordAgePolicy) *iam_model.PasswordAgePolicy {
	return &iam_model.PasswordAgePolicy{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
	}
}

func passwordAgePolicyFromModel(policy *iam_model.PasswordAgePolicy) *management.PasswordAgePolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-bMd9o").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-jFs89").OnError(err).Debug("date parse failed")

	return &management.PasswordAgePolicy{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		CreationDate:   changeDate,
		ChangeDate:     creationDate,
	}
}

func passwordAgePolicyViewFromModel(policy *iam_model.PasswordAgePolicyView) *management.PasswordAgePolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-4Bms9").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-6Hmlo").OnError(err).Debug("date parse failed")

	return &management.PasswordAgePolicyView{
		Default:        policy.Default,
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		ChangeDate:     changeDate,
		CreationDate:   creationDate,
	}
}

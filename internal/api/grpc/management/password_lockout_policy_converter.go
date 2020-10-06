package management

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func passwordLockoutPolicyRequestToModel(policy *management.PasswordLockoutPolicyRequest) *iam_model.PasswordLockoutPolicy {
	return &iam_model.PasswordLockoutPolicy{
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockoutFailure,
	}
}

func passwordLockoutPolicyFromModel(policy *iam_model.PasswordLockoutPolicy) *management.PasswordLockoutPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-bMd9o").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-jFs89").OnError(err).Debug("date parse failed")

	return &management.PasswordLockoutPolicy{
		MaxAttempts:        policy.MaxAttempts,
		ShowLockoutFailure: policy.ShowLockOutFailures,
		CreationDate:       changeDate,
		ChangeDate:         creationDate,
	}
}

func passwordLockoutPolicyViewFromModel(policy *iam_model.PasswordLockoutPolicyView) *management.PasswordLockoutPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-4Bms9").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-6Hmlo").OnError(err).Debug("date parse failed")

	return &management.PasswordLockoutPolicyView{
		Default:            policy.Default,
		MaxAttempts:        policy.MaxAttempts,
		ShowLockoutFailure: policy.ShowLockOutFailures,
		ChangeDate:         changeDate,
		CreationDate:       creationDate,
	}
}

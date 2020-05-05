package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
	"github.com/golang/protobuf/ptypes"
)

func passwordLockoutPolicyFromModel(policy *model.PasswordLockoutPolicy) *PasswordLockoutPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-JRSbT").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-1sizr").OnError(err).Debug("unable to parse timestamp")

	return &PasswordLockoutPolicy{
		Id:                  policy.AggregateID,
		CreationDate:        creationDate,
		ChangeDate:          changeDate,
		Sequence:            policy.Sequence,
		Description:         policy.Description,
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
	}
}

func passwordLockoutPolicyToModel(policy *PasswordLockoutPolicy) *model.PasswordLockoutPolicy {
	creationDate, err := ptypes.Timestamp(policy.CreationDate)
	logging.Log("GRPC-8a511").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.Timestamp(policy.ChangeDate)
	logging.Log("GRPC-2rdGv").OnError(err).Debug("unable to parse timestamp")

	return &model.PasswordLockoutPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  policy.Id,
			CreationDate: creationDate,
			ChangeDate:   changeDate,
			Sequence:     policy.Sequence,
		},
		Description: policy.Description,

		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
	}
}

func passwordLockoutPolicyCreateToModel(policy *PasswordLockoutPolicyCreate) *model.PasswordLockoutPolicy {
	return &model.PasswordLockoutPolicy{
		Description:         policy.Description,
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
	}
}

func passwordLockoutPolicyUpdateToModel(policy *PasswordLockoutPolicyUpdate) *model.PasswordLockoutPolicy {
	return &model.PasswordLockoutPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: policy.Id,
		},
		Description:         policy.Description,
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
	}
}

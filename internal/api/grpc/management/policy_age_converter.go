package management

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func passwordAgePolicyFromModel(policy *model.PasswordAgePolicy) *management.PasswordAgePolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-6ILdB").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-ngUzJ").OnError(err).Debug("unable to parse timestamp")

	return &management.PasswordAgePolicy{
		Id:             policy.AggregateID,
		CreationDate:   creationDate,
		ChangeDate:     changeDate,
		Sequence:       policy.Sequence,
		Description:    policy.Description,
		ExpireWarnDays: policy.ExpireWarnDays,
		MaxAgeDays:     policy.MaxAgeDays,
		IsDefault:      policy.AggregateID == "",
	}
}

func passwordAgePolicyToModel(policy *management.PasswordAgePolicy) *model.PasswordAgePolicy {
	creationDate, err := ptypes.Timestamp(policy.CreationDate)
	logging.Log("GRPC-2QSfU").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.Timestamp(policy.ChangeDate)
	logging.Log("GRPC-LdU91").OnError(err).Debug("unable to parse timestamp")

	return &model.PasswordAgePolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  policy.Id,
			CreationDate: creationDate,
			ChangeDate:   changeDate,
			Sequence:     policy.Sequence,
		},
		Description:    policy.Description,
		ExpireWarnDays: policy.ExpireWarnDays,
		MaxAgeDays:     policy.MaxAgeDays,
	}
}

func passwordAgePolicyCreateToModel(policy *management.PasswordAgePolicyCreate) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		Description:    policy.Description,
		ExpireWarnDays: policy.ExpireWarnDays,
		MaxAgeDays:     policy.MaxAgeDays,
	}
}

func passwordAgePolicyUpdateToModel(policy *management.PasswordAgePolicyUpdate) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: policy.Id,
		},
		Description:    policy.Description,
		ExpireWarnDays: policy.ExpireWarnDays,
		MaxAgeDays:     policy.MaxAgeDays,
	}
}

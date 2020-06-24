package grpc

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
	"github.com/caos/zitadel/pkg/management/grpc"

	"github.com/golang/protobuf/ptypes"
)

func passwordAgePolicyFromModel(policy *model.PasswordAgePolicy) *grpc.PasswordAgePolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-6ILdB").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-ngUzJ").OnError(err).Debug("unable to parse timestamp")

	return &grpc.PasswordAgePolicy{
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

func passwordAgePolicyToModel(policy *grpc.PasswordAgePolicy) *model.PasswordAgePolicy {
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

func passwordAgePolicyCreateToModel(policy *grpc.PasswordAgePolicyCreate) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		Description:    policy.Description,
		ExpireWarnDays: policy.ExpireWarnDays,
		MaxAgeDays:     policy.MaxAgeDays,
	}
}

func passwordAgePolicyUpdateToModel(policy *grpc.PasswordAgePolicyUpdate) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: policy.Id,
		},
		Description:    policy.Description,
		ExpireWarnDays: policy.ExpireWarnDays,
		MaxAgeDays:     policy.MaxAgeDays,
	}
}

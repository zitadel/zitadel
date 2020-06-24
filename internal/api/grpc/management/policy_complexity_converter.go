package grpc

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
	"github.com/caos/zitadel/pkg/management/grpc"

	"github.com/golang/protobuf/ptypes"
)

func passwordComplexityPolicyFromModel(policy *model.PasswordComplexityPolicy) *grpc.PasswordComplexityPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-cQRHE").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-PVA1c").OnError(err).Debug("unable to parse timestamp")

	return &grpc.PasswordComplexityPolicy{
		Id:           policy.AggregateID,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Description:  policy.Description,
		Sequence:     policy.Sequence,
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
		IsDefault:    policy.AggregateID == "",
	}
}

func passwordComplexityPolicyToModel(policy *grpc.PasswordComplexityPolicy) *model.PasswordComplexityPolicy {
	creationDate, err := ptypes.Timestamp(policy.CreationDate)
	logging.Log("GRPC-asmEZ").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.Timestamp(policy.ChangeDate)
	logging.Log("GRPC-MCE4o").OnError(err).Debug("unable to parse timestamp")

	return &model.PasswordComplexityPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  policy.Id,
			CreationDate: creationDate,
			ChangeDate:   changeDate,
			Sequence:     policy.Sequence,
		},
		Description:  policy.Description,
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}

func passwordComplexityPolicyCreateToModel(policy *grpc.PasswordComplexityPolicyCreate) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		Description:  policy.Description,
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}

func passwordComplexityPolicyUpdateToModel(policy *grpc.PasswordComplexityPolicyUpdate) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: policy.Id,
		},
		Description:  policy.Description,
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}

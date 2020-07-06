package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/policy/model"
	"github.com/golang/protobuf/ptypes"
)

func passwordComplexityPolicyFromModel(policy *model.PasswordComplexityPolicy) *PasswordComplexityPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-Lsi3d").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-P0wr4").OnError(err).Debug("unable to parse timestamp")

	return &PasswordComplexityPolicy{
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

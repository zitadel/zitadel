package auth

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/policy/model"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func passwordComplexityPolicyFromModel(policy *model.PasswordComplexityPolicy) *auth.PasswordComplexityPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-Lsi3d").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-P0wr4").OnError(err).Debug("unable to parse timestamp")

	return &auth.PasswordComplexityPolicy{
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

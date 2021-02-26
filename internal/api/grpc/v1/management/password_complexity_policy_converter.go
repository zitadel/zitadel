package management

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func passwordComplexityPolicyRequestToDomain(ctx context.Context, policy *management.PasswordComplexityPolicyRequest) *domain.PasswordComplexityPolicy {
	return &domain.PasswordComplexityPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasSymbol:    policy.HasSymbol,
		HasNumber:    policy.HasNumber,
	}
}

func passwordComplexityPolicyFromDomain(policy *domain.PasswordComplexityPolicy) *management.PasswordComplexityPolicy {
	return &management.PasswordComplexityPolicy{
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasSymbol:    policy.HasSymbol,
		HasNumber:    policy.HasNumber,
		ChangeDate:   timestamppb.New(policy.ChangeDate),
	}
}

func passwordComplexityPolicyViewFromModel(policy *iam_model.PasswordComplexityPolicyView) *management.PasswordComplexityPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-wmi8f").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-dmOp0").OnError(err).Debug("date parse failed")

	return &management.PasswordComplexityPolicyView{
		Default:      policy.Default,
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasSymbol:    policy.HasSymbol,
		HasNumber:    policy.HasNumber,
		CreationDate: changeDate,
		ChangeDate:   creationDate,
	}
}

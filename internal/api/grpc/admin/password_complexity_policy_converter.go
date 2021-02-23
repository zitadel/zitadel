package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func passwordComplexityPolicyToDomain(policy *admin.DefaultPasswordComplexityPolicyRequest) *domain.PasswordComplexityPolicy {
	return &domain.PasswordComplexityPolicy{
		MinLength:    policy.MinLength,
		HasUppercase: policy.HasUppercase,
		HasLowercase: policy.HasLowercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}

func passwordComplexityPolicyFromDomain(policy *domain.PasswordComplexityPolicy) *admin.DefaultPasswordComplexityPolicy {
	return &admin.DefaultPasswordComplexityPolicy{
		MinLength:    policy.MinLength,
		HasUppercase: policy.HasUppercase,
		HasLowercase: policy.HasLowercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
		ChangeDate:   timestamppb.New(policy.ChangeDate),
	}
}

func passwordComplexityPolicyViewFromModel(policy *iam_model.PasswordComplexityPolicyView) *admin.DefaultPasswordComplexityPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-rTs9f").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-ks9Zt").OnError(err).Debug("date parse failed")

	return &admin.DefaultPasswordComplexityPolicyView{
		MinLength:    policy.MinLength,
		HasUppercase: policy.HasUppercase,
		HasLowercase: policy.HasLowercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}

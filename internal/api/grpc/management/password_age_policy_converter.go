package management

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func passwordAgePolicyRequestToDomain(ctx context.Context, policy *management.PasswordAgePolicyRequest) *domain.PasswordAgePolicy {
	return &domain.PasswordAgePolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
	}
}

func passwordAgePolicyFromDomain(policy *domain.PasswordAgePolicy) *management.PasswordAgePolicy {
	return &management.PasswordAgePolicy{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		CreationDate:   timestamppb.New(policy.CreationDate),
		ChangeDate:     timestamppb.New(policy.ChangeDate),
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

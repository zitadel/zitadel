package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func updateDefaultLabelPolicyToDomain(policy *admin_pb.UpdateDefaultLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:   policy.PrimaryColor,
		SecondaryColor: policy.SecondaryColor,
	}
}

package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func AddIAMMemberToDomain(req *admin.AddIAMMemberRequest) *domain.Member {
	return &domain.Member{
		UserID: req.UserId,
		Roles:  req.Roles,
	}
}

func UpdateIAMMemberToDomain(req *admin.UpdateIAMMemberRequest) *domain.Member {
	return &domain.Member{
		UserID: req.UserId,
		Roles:  req.Roles,
	}
}

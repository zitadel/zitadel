package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/user/model"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func ListMyLinkedIDPsRequestToModel(req *auth_pb.ListMyLinkedIDPsRequest) *model.ExternalIDPSearchRequest {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &model.ExternalIDPSearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
	}
}

func RemoveMyLinkedIDPRequestToDomain(ctx context.Context, req *auth_pb.RemoveMyLinkedIDPRequest) *domain.ExternalIDP {
	return &domain.ExternalIDP{
		ObjectRoot:     ctxToObjectRoot(ctx),
		IDPConfigID:    req.IdpId,
		ExternalUserID: req.LinkedUserId,
	}
}

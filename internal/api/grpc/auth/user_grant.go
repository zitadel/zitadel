package auth

import (
	"github.com/caos/zitadel/internal/usergrant/model"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func ListMyUserGrantsRequestToModel(req *auth_pb.ListMyUserGrantsRequest) *model.UserGrantSearchRequest {
	return &model.UserGrantSearchRequest{
		Offset: req.Query.Offset,
		Limit:  uint64(req.Query.Limit),
		Asc:    req.Query.Asc,
	}
}

func UserGrantsToPb(grants []*model.UserGrantView) []*auth_pb.UserGrant {
	userGrants := make([]*auth_pb.UserGrant, len(grants))
	for i, grant := range grants {
		userGrants[i] = UserGrantToPb(grant)
	}
	return userGrants
}

func UserGrantToPb(grant *model.UserGrantView) *auth_pb.UserGrant {
	return &auth_pb.UserGrant{
		GrantId:   grant.ID,
		OrgId:     grant.ResourceOwner,
		OrgName:   grant.OrgName,
		ProjectId: grant.ProjectID,
		UserId:    grant.UserID,
		Roles:     grant.RoleKeys,
	}
}

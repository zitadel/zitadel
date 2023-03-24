package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func ListMyUserGrantsRequestToQuery(ctx context.Context, req *auth_pb.ListMyUserGrantsRequest) (*query.UserGrantsQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	userGrantUserID, err := query.NewUserGrantUserIDSearchQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &query.UserGrantsQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: []query.SearchQuery{
			userGrantUserID,
		},
	}, nil
}

func UserGrantsToPb(grants []*query.UserGrant) []*auth_pb.UserGrant {
	userGrants := make([]*auth_pb.UserGrant, len(grants))
	for i, grant := range grants {
		userGrants[i] = UserGrantToPb(grant)
	}
	return userGrants
}

func UserGrantToPb(grant *query.UserGrant) *auth_pb.UserGrant {
	return &auth_pb.UserGrant{
		GrantId:   grant.ID,
		OrgId:     grant.ResourceOwner,
		OrgName:   grant.OrgName,
		ProjectId: grant.ProjectID,
		UserId:    grant.UserID,
		Roles:     grant.Roles,
		Details: object.ToViewDetailsPb(
			grant.Sequence,
			grant.CreationDate,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
		OrgDomain:      grant.OrgPrimaryDomain,
		ProjectName:    grant.ProjectName,
		ProjectGrantId: grant.GrantID,
		RoleKeys:       grant.Roles,
	}
}

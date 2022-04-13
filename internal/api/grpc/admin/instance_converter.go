package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	instance_grpc "github.com/caos/zitadel/internal/api/grpc/instance"
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/query"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func ListInstanceDomainsRequestToModel(req *admin_pb.ListInstanceDomainsRequest) (*query.InstanceDomainSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := instance_grpc.DomainQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.InstanceDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

func AddOrgDomainRequestToDomain(ctx context.Context, req *mgmt_pb.AddOrgDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain: req.Domain,
	}
}

func RemoveOrgDomainRequestToDomain(ctx context.Context, req *mgmt_pb.RemoveOrgDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain: req.Domain,
	}
}

func ValidateOrgDomainRequestToDomain(ctx context.Context, req *mgmt_pb.ValidateOrgDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain: req.Domain,
	}
}

func SetPrimaryOrgDomainRequestToDomain(ctx context.Context, req *mgmt_pb.SetPrimaryOrgDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain: req.Domain,
	}
}

func AddOrgMemberRequestToDomain(ctx context.Context, req *mgmt_pb.AddOrgMemberRequest) *domain.Member {
	return domain.NewMember(authz.GetCtxData(ctx).OrgID, req.UserId, req.Roles...)
}

func UpdateOrgMemberRequestToDomain(ctx context.Context, req *mgmt_pb.UpdateOrgMemberRequest) *domain.Member {
	return domain.NewMember(authz.GetCtxData(ctx).OrgID, req.UserId, req.Roles...)
}

func ListOrgMembersRequestToModel(ctx context.Context, req *mgmt_pb.ListOrgMembersRequest) (*query.OrgMembersQuery, error) {
	ctxData := authz.GetCtxData(ctx)
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := member_grpc.MemberQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	ownerQuery, err := query.NewMemberResourceOwnerSearchQuery(ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	queries = append(queries, ownerQuery)
	return &query.OrgMembersQuery{
		MembersQuery: query.MembersQuery{
			SearchRequest: query.SearchRequest{
				Offset: offset,
				Limit:  limit,
				Asc:    asc,
				//SortingColumn: //TODO: sorting
			},
			Queries: queries,
		},
		OrgID: ctxData.OrgID,
	}, nil
}

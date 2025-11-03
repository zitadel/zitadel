package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	member_grpc "github.com/zitadel/zitadel/internal/api/grpc/member"
	"github.com/zitadel/zitadel/internal/api/grpc/metadata"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	org_grpc "github.com/zitadel/zitadel/internal/api/grpc/org"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func ListOrgDomainsRequestToModel(req *mgmt_pb.ListOrgDomainsRequest) (*query.OrgDomainSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := org_grpc.DomainQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.OrgDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		//  SortingColumn: //TODO: sorting
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

func AddOrgMemberRequestToCommand(req *mgmt_pb.AddOrgMemberRequest, orgID string) *command.AddOrgMember {
	return &command.AddOrgMember{
		OrgID:  orgID,
		UserID: req.UserId,
		Roles:  req.Roles,
	}
}

func UpdateOrgMemberRequestToCommand(req *mgmt_pb.UpdateOrgMemberRequest, orgID string) *command.ChangeOrgMember {
	return &command.ChangeOrgMember{
		OrgID:  orgID,
		UserID: req.UserId,
		Roles:  req.Roles,
	}
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
				// SortingColumn: //TODO: sorting
			},
			Queries: queries,
		},
		OrgID: ctxData.OrgID,
	}, nil
}

func BulkSetOrgMetadataToDomain(req *mgmt_pb.BulkSetOrgMetadataRequest) []*domain.Metadata {
	metadata := make([]*domain.Metadata, len(req.Metadata))
	for i, data := range req.Metadata {
		metadata[i] = &domain.Metadata{
			Key:   data.Key,
			Value: data.Value,
		}
	}
	return metadata
}

func ListOrgMetadataToDomain(req *mgmt_pb.ListOrgMetadataRequest) (*query.OrgMetadataSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := metadata.OrgMetadataQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.OrgMetadataSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

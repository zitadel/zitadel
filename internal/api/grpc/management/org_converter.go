package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	org_grpc "github.com/caos/zitadel/internal/api/grpc/org"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func ListOrgDomainsRequestToModel(req *mgmt_pb.ListOrgDomainsRequest) (*org_model.OrgDomainSearchRequest, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := org_grpc.DomainQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &org_model.OrgDomainSearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
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

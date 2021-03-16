package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	change_grpc "github.com/caos/zitadel/internal/api/grpc/change"
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	"github.com/caos/zitadel/internal/api/grpc/object"
	org_grpc "github.com/caos/zitadel/internal/api/grpc/org"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetMyOrg(ctx context.Context, req *mgmt_pb.GetMyOrgRequest) (*mgmt_pb.GetMyOrgResponse, error) {
	org, err := s.org.OrgByID(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetMyOrgResponse{Org: org_grpc.OrgViewToPb(org)}, nil
}

func (s *Server) GetOrgByDomainGlobal(ctx context.Context, req *mgmt_pb.GetOrgByDomainGlobalRequest) (*mgmt_pb.GetOrgByDomainGlobalResponse, error) {
	org, err := s.org.OrgByDomainGlobal(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetOrgByDomainGlobalResponse{Org: org_grpc.OrgViewToPb(org)}, nil
}

func (s *Server) ListOrgChanges(ctx context.Context, req *mgmt_pb.ListOrgChangesRequest) (*mgmt_pb.ListOrgChangesResponse, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	features, err := s.features.GetOrgFeatures(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	response, err := s.org.OrgChanges(ctx, authz.GetCtxData(ctx).OrgID, offset, limit, asc, features.AuditLogRetention)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgChangesResponse{
		Result: change_grpc.OrgChangesToPb(response.Changes),
	}, nil
}

func (s *Server) AddOrg(ctx context.Context, req *mgmt_pb.AddOrgRequest) (*mgmt_pb.AddOrgResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	org, err := s.command.AddOrg(ctx, req.Name, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgResponse{
		Id: org.AggregateID,
		Details: object.ToDetailsPb(
			org.Sequence,
			org.ChangeDate,
			org.ResourceOwner,
		),
	}, err
}

func (s *Server) DeactivateOrg(ctx context.Context, req *mgmt_pb.DeactivateOrgRequest) (*mgmt_pb.DeactivateOrgResponse, error) {
	objectDetails, err := s.command.DeactivateOrg(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateOrgResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ReactivateOrg(ctx context.Context, req *mgmt_pb.ReactivateOrgRequest) (*mgmt_pb.ReactivateOrgResponse, error) {
	objectDetails, err := s.command.ReactivateOrg(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateOrgResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, err
}

func (s *Server) GetOrgIAMPolicy(ctx context.Context, req *mgmt_pb.GetOrgIAMPolicyRequest) (*mgmt_pb.GetOrgIAMPolicyResponse, error) {
	policy, err := s.org.GetMyOrgIamPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetOrgIAMPolicyResponse{
		Policy: policy_grpc.OrgIAMPolicyToPb(policy),
	}, nil
}

func (s *Server) ListOrgDomains(ctx context.Context, req *mgmt_pb.ListOrgDomainsRequest) (*mgmt_pb.ListOrgDomainsResponse, error) {
	queries, err := ListOrgDomainsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	domains, err := s.org.SearchMyOrgDomains(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgDomainsResponse{
		Result: org_grpc.DomainsToPb(domains.Result),
		Details: object.ToListDetails(
			domains.TotalResult,
			domains.Sequence,
			domains.Timestamp,
		),
	}, nil
}

func (s *Server) AddOrgDomain(ctx context.Context, req *mgmt_pb.AddOrgDomainRequest) (*mgmt_pb.AddOrgDomainResponse, error) {
	domain, err := s.command.AddOrgDomain(ctx, AddOrgDomainRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgDomainResponse{
		Details: object.ToDetailsPb(
			domain.Sequence,
			domain.ChangeDate,
			domain.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveOrgDomain(ctx context.Context, req *mgmt_pb.RemoveOrgDomainRequest) (*mgmt_pb.RemoveOrgDomainResponse, error) {
	details, err := s.command.RemoveOrgDomain(ctx, RemoveOrgDomainRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveOrgDomainResponse{
		Details: object.DomainToDetailsPb(details),
	}, err
}

func (s *Server) GenerateOrgDomainValidation(ctx context.Context, req *mgmt_pb.GenerateOrgDomainValidationRequest) (*mgmt_pb.GenerateOrgDomainValidationResponse, error) {
	token, url, err := s.command.GenerateOrgDomainValidation(ctx, GenerateOrgDomainValidationRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GenerateOrgDomainValidationResponse{
		Token: token,
		Url:   url,
		//TODO: remove details from proto
	}, nil
}

func GenerateOrgDomainValidationRequestToDomain(ctx context.Context, req *mgmt_pb.GenerateOrgDomainValidationRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain:         req.Domain,
		ValidationType: org_grpc.DomainValidationTypeToDomain(req.Type),
	}
}

func (s *Server) ValidateOrgDomain(ctx context.Context, req *mgmt_pb.ValidateOrgDomainRequest) (*mgmt_pb.ValidateOrgDomainResponse, error) {
	details, err := s.command.ValidateOrgDomain(ctx, ValidateOrgDomainRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ValidateOrgDomainResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) SetPrimaryOrgDomain(ctx context.Context, req *mgmt_pb.SetPrimaryOrgDomainRequest) (*mgmt_pb.SetPrimaryOrgDomainResponse, error) {
	details, err := s.command.SetPrimaryOrgDomain(ctx, SetPrimaryOrgDomainRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetPrimaryOrgDomainResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ListOrgMemberRoles(ctx context.Context, req *mgmt_pb.ListOrgMemberRolesRequest) (*mgmt_pb.ListOrgMemberRolesResponse, error) {
	roles := s.org.GetOrgMemberRoles()
	return &mgmt_pb.ListOrgMemberRolesResponse{
		Result: roles,
	}, nil
}

func (s *Server) ListOrgMembers(ctx context.Context, req *mgmt_pb.ListOrgMembersRequest) (*mgmt_pb.ListOrgMembersResponse, error) {
	queries, err := ListOrgMembersRequestToModel(req)
	if err != nil {
		return nil, err
	}
	members, err := s.org.SearchMyOrgMembers(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgMembersResponse{
		Result: member_grpc.OrgMembersToPb(members.Result),
		Details: object.ToListDetails(
			members.TotalResult,
			members.Sequence,
			members.Timestamp,
		),
	}, nil
}

func ListOrgMembersRequestToModel(req *mgmt_pb.ListOrgMembersRequest) (*org_model.OrgMemberSearchRequest, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries := member_grpc.MemberQueriesToOrgMember(req.Queries)
	return &org_model.OrgMemberSearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

func (s *Server) AddOrgMember(ctx context.Context, req *mgmt_pb.AddOrgMemberRequest) (*mgmt_pb.AddOrgMemberResponse, error) {
	addedMember, err := s.command.AddOrgMember(ctx, AddOrgMemberRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgMemberResponse{
		Details: object.ToDetailsPb(
			addedMember.Sequence,
			addedMember.ChangeDate,
			addedMember.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateOrgMember(ctx context.Context, req *mgmt_pb.UpdateOrgMemberRequest) (*mgmt_pb.UpdateOrgMemberResponse, error) {
	changedMember, err := s.command.ChangeOrgMember(ctx, UpdateOrgMemberRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgMemberResponse{
		Details: object.ToDetailsPb(
			changedMember.Sequence,
			changedMember.ChangeDate,
			changedMember.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveOrgMember(ctx context.Context, req *mgmt_pb.RemoveOrgMemberRequest) (*mgmt_pb.RemoveOrgMemberResponse, error) {
	details, err := s.command.RemoveOrgMember(ctx, authz.GetCtxData(ctx).OrgID, req.UserId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveOrgMemberResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

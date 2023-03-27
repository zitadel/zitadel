package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	change_grpc "github.com/zitadel/zitadel/internal/api/grpc/change"
	member_grpc "github.com/zitadel/zitadel/internal/api/grpc/member"
	"github.com/zitadel/zitadel/internal/api/grpc/metadata"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	obj_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	org_grpc "github.com/zitadel/zitadel/internal/api/grpc/org"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/org"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetMyOrg(ctx context.Context, req *mgmt_pb.GetMyOrgRequest) (*mgmt_pb.GetMyOrgResponse, error) {
	org, err := s.query.OrgByID(ctx, true, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetMyOrgResponse{Org: org_grpc.OrgViewToPb(org)}, nil
}

func (s *Server) GetOrgByDomainGlobal(ctx context.Context, req *mgmt_pb.GetOrgByDomainGlobalRequest) (*mgmt_pb.GetOrgByDomainGlobalResponse, error) {
	org, err := s.query.OrgByPrimaryDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}

	return &mgmt_pb.GetOrgByDomainGlobalResponse{Org: org_grpc.OrgViewToPb(org)}, nil
}

func (s *Server) ListOrgChanges(ctx context.Context, req *mgmt_pb.ListOrgChangesRequest) (*mgmt_pb.ListOrgChangesResponse, error) {
	var (
		limit    uint64
		sequence uint64
		asc      bool
	)
	if req.Query != nil {
		limit = uint64(req.Query.Limit)
		sequence = req.Query.Sequence
		asc = req.Query.Asc
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AllowTimeTravel().
		Limit(limit).
		OrderDesc().
		ResourceOwner(authz.GetCtxData(ctx).OrgID).
		AddQuery().
		SequenceGreater(sequence).
		AggregateTypes(org.AggregateType).
		AggregateIDs(authz.GetCtxData(ctx).OrgID).
		Builder()
	if asc {
		query.OrderAsc()
	}

	response, err := s.query.SearchEvents(ctx, query, s.auditLogRetention)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgChangesResponse{
		Result: change_grpc.EventsToChangesPb(response, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) AddOrg(ctx context.Context, req *mgmt_pb.AddOrgRequest) (*mgmt_pb.AddOrgResponse, error) {
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, domain.NewIAMDomainName(req.Name, authz.GetInstance(ctx).RequestedDomain()), "")
	if err != nil {
		return nil, err
	}
	ctxData := authz.GetCtxData(ctx)
	org, err := s.command.AddOrg(ctx, req.Name, ctxData.UserID, ctxData.ResourceOwner, userIDs)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgResponse{
		Id: org.AggregateID,
		Details: object.AddToDetailsPb(
			org.Sequence,
			org.ChangeDate,
			org.ResourceOwner,
		),
	}, err
}

func (s *Server) UpdateOrg(ctx context.Context, req *mgmt_pb.UpdateOrgRequest) (*mgmt_pb.UpdateOrgResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	org, err := s.command.ChangeOrg(ctx, ctxData.OrgID, req.Name)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgResponse{
		Details: object.AddToDetailsPb(
			org.Sequence,
			org.EventDate,
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
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ReactivateOrg(ctx context.Context, req *mgmt_pb.ReactivateOrgRequest) (*mgmt_pb.ReactivateOrgResponse, error) {
	objectDetails, err := s.command.ReactivateOrg(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateOrgResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, err
}

func (s *Server) RemoveOrg(ctx context.Context, req *mgmt_pb.RemoveOrgRequest) (*mgmt_pb.RemoveOrgResponse, error) {
	details, err := s.command.RemoveOrg(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}

	return &mgmt_pb.RemoveOrgResponse{Details: object.DomainToChangeDetailsPb(details)}, nil
}

func (s *Server) GetDomainPolicy(ctx context.Context, req *mgmt_pb.GetDomainPolicyRequest) (*mgmt_pb.GetDomainPolicyResponse, error) {
	policy, err := s.query.DomainPolicyByOrg(ctx, true, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDomainPolicyResponse{
		Policy: policy_grpc.DomainPolicyToPb(policy),
	}, nil
}

func (s *Server) GetOrgIAMPolicy(ctx context.Context, _ *mgmt_pb.GetOrgIAMPolicyRequest) (*mgmt_pb.GetOrgIAMPolicyResponse, error) {
	policy, err := s.query.DomainPolicyByOrg(ctx, true, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetOrgIAMPolicyResponse{
		Policy: policy_grpc.DomainPolicyToOrgIAMPb(policy),
	}, nil
}

func (s *Server) ListOrgDomains(ctx context.Context, req *mgmt_pb.ListOrgDomainsRequest) (*mgmt_pb.ListOrgDomainsResponse, error) {
	queries, err := ListOrgDomainsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	orgIDQuery, err := query.NewOrgDomainOrgIDSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	queries.Queries = append(queries.Queries, orgIDQuery)

	domains, err := s.query.SearchOrgDomains(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgDomainsResponse{
		Result:  org_grpc.DomainsToPb(domains.Domains),
		Details: object.ToListDetails(domains.Count, domains.Sequence, domains.Timestamp),
	}, nil
}

func (s *Server) AddOrgDomain(ctx context.Context, req *mgmt_pb.AddOrgDomainRequest) (*mgmt_pb.AddOrgDomainResponse, error) {
	orgID := authz.GetCtxData(ctx).OrgID
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, req.Domain, orgID)
	if err != nil {
		return nil, err
	}
	details, err := s.command.AddOrgDomain(ctx, orgID, req.Domain, userIDs)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgDomainResponse{
		Details: object.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) RemoveOrgDomain(ctx context.Context, req *mgmt_pb.RemoveOrgDomainRequest) (*mgmt_pb.RemoveOrgDomainResponse, error) {
	details, err := s.command.RemoveOrgDomain(ctx, RemoveOrgDomainRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveOrgDomainResponse{
		Details: object.DomainToChangeDetailsPb(details),
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
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, req.Domain, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	details, err := s.command.ValidateOrgDomain(ctx, ValidateOrgDomainRequestToDomain(ctx, req), userIDs)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ValidateOrgDomainResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) SetPrimaryOrgDomain(ctx context.Context, req *mgmt_pb.SetPrimaryOrgDomainRequest) (*mgmt_pb.SetPrimaryOrgDomainResponse, error) {
	details, err := s.command.SetPrimaryOrgDomain(ctx, SetPrimaryOrgDomainRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetPrimaryOrgDomainResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ListOrgMemberRoles(ctx context.Context, _ *mgmt_pb.ListOrgMemberRolesRequest) (*mgmt_pb.ListOrgMemberRolesResponse, error) {
	instance, err := s.query.Instance(ctx, false)
	if err != nil {
		return nil, err
	}
	roles := s.query.GetOrgMemberRoles(authz.GetCtxData(ctx).OrgID == instance.DefaultOrgID)
	return &mgmt_pb.ListOrgMemberRolesResponse{
		Result: roles,
	}, nil
}

func (s *Server) ListOrgMembers(ctx context.Context, req *mgmt_pb.ListOrgMembersRequest) (*mgmt_pb.ListOrgMembersResponse, error) {
	queries, err := ListOrgMembersRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	members, err := s.query.OrgMembers(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgMembersResponse{
		Result:  member_grpc.MembersToPb(s.assetAPIPrefix(ctx), members.Members),
		Details: object.ToListDetails(members.Count, members.Sequence, members.Timestamp),
	}, nil
}

func (s *Server) AddOrgMember(ctx context.Context, req *mgmt_pb.AddOrgMemberRequest) (*mgmt_pb.AddOrgMemberResponse, error) {
	addedMember, err := s.command.AddOrgMember(ctx, authz.GetCtxData(ctx).OrgID, req.UserId, req.Roles...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgMemberResponse{
		Details: object.AddToDetailsPb(
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
		Details: object.ChangeToDetailsPb(
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
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) getClaimedUserIDsOfOrgDomain(ctx context.Context, orgDomain, orgID string) ([]string, error) {
	queries := make([]query.SearchQuery, 0, 2)
	loginName, err := query.NewUserPreferredLoginNameSearchQuery("@"+orgDomain, query.TextEndsWithIgnoreCase)
	if err != nil {
		return nil, err
	}
	queries = append(queries, loginName)
	if orgID != "" {
		owner, err := query.NewUserResourceOwnerSearchQuery(orgID, query.TextNotEquals)
		if err != nil {
			return nil, err
		}
		queries = append(queries, owner)
	}
	users, err := s.query.SearchUsers(ctx, &query.UserSearchQueries{Queries: queries}, false)
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, len(users.Users))
	for i, user := range users.Users {
		userIDs[i] = user.ID
	}
	return userIDs, nil
}

func (s *Server) ListOrgMetadata(ctx context.Context, req *mgmt_pb.ListOrgMetadataRequest) (*mgmt_pb.ListOrgMetadataResponse, error) {
	metadataQueries, err := ListOrgMetadataToDomain(req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchOrgMetadata(ctx, true, authz.GetCtxData(ctx).OrgID, metadataQueries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgMetadataResponse{
		Result:  metadata.OrgMetadataListToPb(res.Metadata),
		Details: obj_grpc.ToListDetails(res.Count, res.Sequence, res.Timestamp),
	}, nil
}

func (s *Server) GetOrgMetadata(ctx context.Context, req *mgmt_pb.GetOrgMetadataRequest) (*mgmt_pb.GetOrgMetadataResponse, error) {
	data, err := s.query.GetOrgMetadataByKey(ctx, true, authz.GetCtxData(ctx).OrgID, req.Key, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetOrgMetadataResponse{
		Metadata: metadata.OrgMetadataToPb(data),
	}, nil
}

func (s *Server) SetOrgMetadata(ctx context.Context, req *mgmt_pb.SetOrgMetadataRequest) (*mgmt_pb.SetOrgMetadataResponse, error) {
	result, err := s.command.SetOrgMetadata(ctx, authz.GetCtxData(ctx).OrgID, &domain.Metadata{Key: req.Key, Value: req.Value})
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetOrgMetadataResponse{
		Details: obj_grpc.AddToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) BulkSetOrgMetadata(ctx context.Context, req *mgmt_pb.BulkSetOrgMetadataRequest) (*mgmt_pb.BulkSetOrgMetadataResponse, error) {
	result, err := s.command.BulkSetOrgMetadata(ctx, authz.GetCtxData(ctx).OrgID, BulkSetOrgMetadataToDomain(req)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.BulkSetOrgMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) RemoveOrgMetadata(ctx context.Context, req *mgmt_pb.RemoveOrgMetadataRequest) (*mgmt_pb.RemoveOrgMetadataResponse, error) {
	result, err := s.command.RemoveOrgMetadata(ctx, authz.GetCtxData(ctx).OrgID, req.Key)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveOrgMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) BulkRemoveOrgMetadata(ctx context.Context, req *mgmt_pb.BulkRemoveOrgMetadataRequest) (*mgmt_pb.BulkRemoveOrgMetadataResponse, error) {
	result, err := s.command.BulkRemoveOrgMetadata(ctx, authz.GetCtxData(ctx).OrgID, req.Keys...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.BulkRemoveOrgMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

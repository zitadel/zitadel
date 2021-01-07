package admin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetOrgByID(ctx context.Context, orgID *admin.OrgID) (_ *admin.Org, err error) {
	org, err := s.org.OrgByID(ctx, orgID.Id)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) SearchOrgs(ctx context.Context, request *admin.OrgSearchRequest) (_ *admin.OrgSearchResponse, err error) {
	result, err := s.org.SearchOrgs(ctx, orgSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return orgSearchResponseFromModel(result), nil
}

func (s *Server) IsOrgUnique(ctx context.Context, request *admin.UniqueOrgRequest) (org *admin.UniqueOrgResponse, err error) {
	isUnique, err := s.org.IsOrgUnique(ctx, request.Name, request.Domain)

	return &admin.UniqueOrgResponse{IsUnique: isUnique}, err
}

func (s *Server) SetUpOrg(ctx context.Context, orgSetUp *admin.OrgSetUpRequest) (_ *admin.OrgSetUpResponse, err error) {
	_, err = s.command.SetUpOrg(ctx, orgCreateRequestToDomain(orgSetUp.Org), userCreateRequestToDomain(orgSetUp.User))
	if err != nil {
		return nil, err
	}
	//return setUpOrgResponseFromModel(setUp), err
	return nil, nil //TODO: !
}

func (s *Server) GetDefaultOrgIamPolicy(ctx context.Context, _ *empty.Empty) (_ *admin.OrgIamPolicyView, err error) {
	policy, err := s.iam.GetDefaultOrgIAMPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return orgIAMPolicyViewFromModel(policy), err
}

func (s *Server) UpdateDefaultOrgIamPolicy(ctx context.Context, in *admin.OrgIamPolicyRequest) (_ *admin.OrgIamPolicy, err error) {
	policy, err := s.command.ChangeDefaultOrgIAMPolicy(ctx, orgIAMPolicyRequestToDomain(in))
	if err != nil {
		return nil, err
	}
	return orgIAMPolicyFromDomain(policy), err
}

func (s *Server) GetOrgIamPolicy(ctx context.Context, in *admin.OrgIamPolicyID) (_ *admin.OrgIamPolicyView, err error) {
	policy, err := s.org.GetOrgIAMPolicyByID(ctx, in.OrgId)
	if err != nil {
		return nil, err
	}
	return orgIAMPolicyViewFromModel(policy), err
}

func (s *Server) CreateOrgIamPolicy(ctx context.Context, in *admin.OrgIamPolicyRequest) (_ *admin.OrgIamPolicy, err error) {
	policy, err := s.command.AddOrgIAMPolicy(ctx, orgIAMPolicyRequestToDomain(in))
	if err != nil {
		return nil, err
	}
	return orgIAMPolicyFromDomain(policy), err
}

func (s *Server) UpdateOrgIamPolicy(ctx context.Context, in *admin.OrgIamPolicyRequest) (_ *admin.OrgIamPolicy, err error) {
	policy, err := s.command.ChangeOrgIAMPolicy(ctx, orgIAMPolicyRequestToDomain(in))
	if err != nil {
		return nil, err
	}
	return orgIAMPolicyFromDomain(policy), err
}

func (s *Server) RemoveOrgIamPolicy(ctx context.Context, in *admin.OrgIamPolicyID) (_ *empty.Empty, err error) {
	err = s.org.RemoveOrgIAMPolicy(ctx, in.OrgId)
	return &empty.Empty{}, err
}

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
	setUp, err := s.org.SetUpOrg(ctx, setUpRequestToModel(orgSetUp))
	if err != nil {
		return nil, err
	}
	return setUpOrgResponseFromModel(setUp), err
}

func (s *Server) GetOrgIamPolicy(ctx context.Context, in *admin.OrgIamPolicyID) (_ *admin.OrgIamPolicy, err error) {
	policy, err := s.org.GetOrgIamPolicyByID(ctx, in.OrgId)
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}

func (s *Server) CreateOrgIamPolicy(ctx context.Context, in *admin.OrgIamPolicyRequest) (_ *admin.OrgIamPolicy, err error) {
	policy, err := s.org.CreateOrgIamPolicy(ctx, orgIamPolicyRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}

func (s *Server) UpdateOrgIamPolicy(ctx context.Context, in *admin.OrgIamPolicyRequest) (_ *admin.OrgIamPolicy, err error) {
	policy, err := s.org.ChangeOrgIamPolicy(ctx, orgIamPolicyRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}

func (s *Server) DeleteOrgIamPolicy(ctx context.Context, in *admin.OrgIamPolicyID) (_ *empty.Empty, err error) {
	err = s.org.RemoveOrgIamPolicy(ctx, in.OrgId)
	return &empty.Empty{}, err
}

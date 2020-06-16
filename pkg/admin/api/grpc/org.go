package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetOrgByID(ctx context.Context, orgID *OrgID) (_ *Org, err error) {
	org, err := s.org.OrgByID(ctx, orgID.Id)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) SearchOrgs(ctx context.Context, request *OrgSearchRequest) (_ *OrgSearchResponse, err error) {
	result, err := s.org.SearchOrgs(ctx, orgSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return &OrgSearchResponse{
		Result:      orgViewsFromModel(result.Result),
		Limit:       request.Limit,
		Offset:      request.Offset,
		TotalResult: result.TotalResult,
	}, nil
}

func (s *Server) IsOrgUnique(ctx context.Context, request *UniqueOrgRequest) (org *UniqueOrgResponse, err error) {
	isUnique, err := s.org.IsOrgUnique(ctx, request.Name, request.Domain)

	return &UniqueOrgResponse{IsUnique: isUnique}, err
}

func (s *Server) SetUpOrg(ctx context.Context, orgSetUp *OrgSetUpRequest) (_ *OrgSetUpResponse, err error) {
	setUp, err := s.org.SetUpOrg(ctx, setUpRequestToModel(orgSetUp))
	if err != nil {
		return nil, err
	}
	return setUpOrgResponseFromModel(setUp), err
}

func (s *Server) GetOrgIamPolicy(ctx context.Context, in *OrgIamPolicyID) (_ *OrgIamPolicy, err error) {
	policy, err := s.org.GetOrgIamPolicyByID(ctx, in.OrgId)
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}

func (s *Server) CreateOrgIamPolicy(ctx context.Context, in *OrgIamPolicyRequest) (_ *OrgIamPolicy, err error) {
	policy, err := s.org.CreateOrgIamPolicy(ctx, orgIamPolicyRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}

func (s *Server) UpdateOrgIamPolicy(ctx context.Context, in *OrgIamPolicyRequest) (_ *OrgIamPolicy, err error) {
	policy, err := s.org.ChangeOrgIamPolicy(ctx, orgIamPolicyRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}

func (s *Server) DeleteOrgIamPolicy(ctx context.Context, in *OrgIamPolicyID) (_ *empty.Empty, err error) {
	err = s.org.RemoveOrgIamPolicy(ctx, in.OrgId)
	return &empty.Empty{}, err
}

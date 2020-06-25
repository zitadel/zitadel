package admin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/admin/grpc"
)

func (s *Server) GetOrgByID(ctx context.Context, orgID *grpc.OrgID) (_ *grpc.Org, err error) {
	org, err := s.org.OrgByID(ctx, orgID.Id)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) SearchOrgs(ctx context.Context, request *grpc.OrgSearchRequest) (_ *grpc.OrgSearchResponse, err error) {
	result, err := s.org.SearchOrgs(ctx, orgSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return &grpc.OrgSearchResponse{
		Result:      orgViewsFromModel(result.Result),
		Limit:       request.Limit,
		Offset:      request.Offset,
		TotalResult: result.TotalResult,
	}, nil
}

func (s *Server) IsOrgUnique(ctx context.Context, request *grpc.UniqueOrgRequest) (org *grpc.UniqueOrgResponse, err error) {
	isUnique, err := s.org.IsOrgUnique(ctx, request.Name, request.Domain)

	return &grpc.UniqueOrgResponse{IsUnique: isUnique}, err
}

func (s *Server) SetUpOrg(ctx context.Context, orgSetUp *grpc.OrgSetUpRequest) (_ *grpc.OrgSetUpResponse, err error) {
	setUp, err := s.org.SetUpOrg(ctx, setUpRequestToModel(orgSetUp))
	if err != nil {
		return nil, err
	}
	return setUpOrgResponseFromModel(setUp), err
}

func (s *Server) GetOrgIamPolicy(ctx context.Context, in *grpc.OrgIamPolicyID) (_ *grpc.OrgIamPolicy, err error) {
	policy, err := s.org.GetOrgIamPolicyByID(ctx, in.OrgId)
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}

func (s *Server) CreateOrgIamPolicy(ctx context.Context, in *grpc.OrgIamPolicyRequest) (_ *grpc.OrgIamPolicy, err error) {
	policy, err := s.org.CreateOrgIamPolicy(ctx, orgIamPolicyRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}

func (s *Server) UpdateOrgIamPolicy(ctx context.Context, in *grpc.OrgIamPolicyRequest) (_ *grpc.OrgIamPolicy, err error) {
	policy, err := s.org.ChangeOrgIamPolicy(ctx, orgIamPolicyRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}

func (s *Server) DeleteOrgIamPolicy(ctx context.Context, in *grpc.OrgIamPolicyID) (_ *empty.Empty, err error) {
	err = s.org.RemoveOrgIamPolicy(ctx, in.OrgId)
	return &empty.Empty{}, err
}

package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/caos/zitadel/internal/api/authz"
	change_grpc "github.com/caos/zitadel/internal/api/grpc/change"
	"github.com/caos/zitadel/internal/api/grpc/object"
	org_grpc "github.com/caos/zitadel/internal/api/grpc/org"
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
	response, err := s.org.OrgChanges(ctx, authz.GetCtxData(ctx).OrgID, req.Query.Offset, uint64(req.Query.Limit), req.Query.Asc)
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
			org.CreationDate,
			org.ChangeDate,
			org.ResourceOwner,
		),
	}, err
}

func (s *Server) DeactivateOrg(ctx context.Context, req *mgmt_pb.DeactivateOrgRequest) (*mgmt_pb.DeactivateOrgResponse, error) {
	err := s.command.DeactivateOrg(ctx, authz.GetCtxData(ctx).OrgID)
	return &mgmt_pb.DeactivateOrgResponse{
		//TODO: details
	}, err
}

func (s *Server) ReactivateOrg(ctx context.Context, req *mgmt_pb.ReactivateOrgRequest) (*mgmt_pb.ReactivateOrgResponse, error) {
	err := s.command.ReactivateOrg(ctx, authz.GetCtxData(ctx).OrgID)
	return &mgmt_pb.ReactivateOrgResponse{
		//TODO: details
	}, err
}

func (s *Server) GetOrgIAMPolicy(ctx context.Context, req *mgmt_pb.GetOrgIAMPolicyRequest) (*mgmt_pb.GetOrgIAMPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrgIAMPolicy not implemented")
}

//func (s *Server) ListOrgDomains(ctx context.Context, req *mgmt_pb.ListOrgDomainsRequest) (*mgmt_pb.ListOrgDomainsResponse, error) {
//	domains, err := s.org.SearchMyOrgDomains(ctx, ListOrgDomainsRequestToModel(req))
//	if err != nil {
//		return nil, err
//	}
//	return orgDomainSearchResponseFromModel(domains), nil
//}
func (s *Server) AddOrgDomain(ctx context.Context, req *mgmt_pb.AddOrgDomainRequest) (*mgmt_pb.AddOrgDomainResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddOrgDomain not implemented")
}
func (s *Server) RemoveOrgDomain(ctx context.Context, req *mgmt_pb.RemoveOrgDomainRequest) (*mgmt_pb.RemoveOrgDomainResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveOrgDomain not implemented")
}
func (s *Server) GenerateOrgDomainValidation(ctx context.Context, req *mgmt_pb.GenerateOrgDomainValidationRequest) (*mgmt_pb.GenerateOrgDomainValidationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateOrgDomainValidation not implemented")
}
func (s *Server) ValidateOrgDomain(ctx context.Context, req *mgmt_pb.ValidateOrgDomainRequest) (*mgmt_pb.ValidateOrgDomainResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateOrgDomain not implemented")
}
func (s *Server) SetPrimaryOrgDomain(ctx context.Context, req *mgmt_pb.SetPrimaryOrgDomainRequest) (*mgmt_pb.SetPrimaryOrgDomainResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPrimaryOrgDomain not implemented")
}
func (s *Server) ListOrgMemberRoles(ctx context.Context, req *mgmt_pb.ListOrgMemberRolesRequest) (*mgmt_pb.ListOrgMemberRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrgMemberRoles not implemented")
}
func (s *Server) ListOrgMembers(ctx context.Context, req *mgmt_pb.ListOrgMembersRequest) (*mgmt_pb.ListOrgMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrgMembers not implemented")
}
func (s *Server) AddOrgMember(ctx context.Context, req *mgmt_pb.AddOrgMemberRequest) (*mgmt_pb.AddOrgMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddOrgMember not implemented")
}
func (s *Server) UpdateOrgMember(ctx context.Context, req *mgmt_pb.UpdateOrgMemberRequest) (*mgmt_pb.UpdateOrgMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrgMember not implemented")
}
func (s *Server) RemoveOrgMember(ctx context.Context, req *mgmt_pb.RemoveOrgMemberRequest) (*mgmt_pb.RemoveOrgMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveOrgMember not implemented")
}

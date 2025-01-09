package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/group"
	obj_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetGroupGrantByID(ctx context.Context, req *mgmt_pb.GetGroupGrantByIDRequest) (*mgmt_pb.GetGroupGrantByIDResponse, error) {
	idQuery, err := query.NewGroupGrantIDSearchQuery(req.GrantId)
	if err != nil {
		return nil, err
	}
	ownerQuery, err := query.NewGroupGrantResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	grant, err := s.query.GroupGrant(ctx, true, idQuery, ownerQuery)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetGroupGrantByIDResponse{
		GroupGrant: group.GroupGrantToPb(s.assetAPIPrefix(ctx), grant),
	}, nil
}

func (s *Server) ListGroupGrants(ctx context.Context, req *mgmt_pb.ListGroupGrantRequest) (*mgmt_pb.ListGroupGrantResponse, error) {
	queries, err := ListGroupGrantsRequestToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.GroupGrants(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListGroupGrantResponse{
		Result:  group.GroupGrantsToPb(s.assetAPIPrefix(ctx), res.GroupGrants),
		Details: obj_grpc.ToListDetails(res.Count, res.Sequence, res.LastRun),
	}, nil
}

func (s *Server) AddGroupGrant(ctx context.Context, req *mgmt_pb.AddGroupGrantRequest) (*mgmt_pb.AddGroupGrantResponse, error) {
	grant := AddGroupGrantRequestToDomain(req)
	if err := checkExplicitProjectPermission(ctx, grant.ProjectGrantID, grant.ProjectID); err != nil {
		return nil, err
	}
	grant, err := s.command.AddGroupGrant(ctx, grant, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGroupGrantResponse{
		GroupGrantId: grant.AggregateID,
		Details: obj_grpc.AddToDetailsPb(
			grant.Sequence,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateGroupGrant(ctx context.Context, req *mgmt_pb.UpdateGroupGrantRequest) (*mgmt_pb.UpdateGroupGrantResponse, error) {
	grant, err := s.command.ChangeGroupGrant(ctx, UpdateGroupGrantRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGroupGrantResponse{
		Details: obj_grpc.ChangeToDetailsPb(
			grant.Sequence,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateGroupGrant(ctx context.Context, req *mgmt_pb.DeactivateGroupGrantRequest) (*mgmt_pb.DeactivateGroupGrantResponse, error) {
	objectDetails, err := s.command.DeactivateGroupGrant(ctx, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateGroupGrantResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ReactivateGroupGrant(ctx context.Context, req *mgmt_pb.ReactivateGroupGrantRequest) (*mgmt_pb.ReactivateGroupGrantResponse, error) {
	objectDetails, err := s.command.ReactivateGroupGrant(ctx, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateGroupGrantResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveGroupGrant(ctx context.Context, req *mgmt_pb.RemoveGroupGrantRequest) (*mgmt_pb.RemoveGroupGrantResponse, error) {
	objectDetails, err := s.command.RemoveGroupGrant(ctx, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveGroupGrantResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) BulkRemoveGroupGrant(ctx context.Context, req *mgmt_pb.BulkRemoveGroupGrantRequest) (*mgmt_pb.BulkRemoveGroupGrantResponse, error) {
	err := s.command.BulkRemoveGroupGrant(ctx, req.GrantId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.BulkRemoveGroupGrantResponse{}, nil
}

package auth

import (
	"context"

	idp_grpc "github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/api/grpc/object"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) ListMyLinkedIDPs(ctx context.Context, req *auth_pb.ListMyLinkedIDPsRequest) (*auth_pb.ListMyLinkedIDPsResponse, error) {
	q, err := ListMyLinkedIDPsRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	idps, err := s.query.LinkedIDPsByUser(ctx, q)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyLinkedIDPsResponse{
		Result: idp_grpc.IDPsToUserLinkPb(idps.IDPs),
		Details: object.ToListDetails(
			idps.Count,
			idps.Sequence,
			idps.Timestamp,
		),
	}, nil
}

func (s *Server) RemoveMyLinkedIDP(ctx context.Context, req *auth_pb.RemoveMyLinkedIDPRequest) (*auth_pb.RemoveMyLinkedIDPResponse, error) {
	objectDetails, err := s.command.RemoveUserIDPLink(ctx, RemoveMyLinkedIDPRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyLinkedIDPResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

package auth

import (
	"context"

	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
	idp_grpc "github.com/caos/zitadel/v2/internal/api/grpc/idp"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
)

func (s *Server) ListMyLinkedIDPs(ctx context.Context, req *auth_pb.ListMyLinkedIDPsRequest) (*auth_pb.ListMyLinkedIDPsResponse, error) {
	q, err := ListMyLinkedIDPsRequestToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	links, err := s.query.IDPUserLinks(ctx, q)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyLinkedIDPsResponse{
		Result: idp_grpc.IDPUserLinksToPb(links.Links),
		Details: object.ToListDetails(
			links.Count,
			links.Sequence,
			links.Timestamp,
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

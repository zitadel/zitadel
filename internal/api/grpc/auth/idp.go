package auth

import (
	"context"

	idp_grpc "github.com/zitadel/zitadel/internal/api/grpc/idp"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) ListMyLinkedIDPs(ctx context.Context, req *auth_pb.ListMyLinkedIDPsRequest) (*auth_pb.ListMyLinkedIDPsResponse, error) {
	q, err := ListMyLinkedIDPsRequestToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	links, err := s.query.IDPUserLinks(ctx, q, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyLinkedIDPsResponse{
		Result:  idp_grpc.IDPUserLinksToPb(links.Links),
		Details: object.ToListDetails(links.Count, links.Sequence, links.LastRun),
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

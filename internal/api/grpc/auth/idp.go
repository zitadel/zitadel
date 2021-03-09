package auth

import (
	"context"

	idp_grpc "github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/api/grpc/object"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) ListMyLinkedIDPs(ctx context.Context, req *auth_pb.ListMyLinkedIDPsRequest) (*auth_pb.ListMyLinkedIDPsResponse, error) {
	idps, err := s.repo.SearchMyExternalIDPs(ctx, ListMyLinkedIDPsRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyLinkedIDPsResponse{
		Result: idp_grpc.IDPsToUserLinkPb(idps.Result),
		Details: object.ToListDetails(
			idps.TotalResult,
			idps.Sequence,
			idps.Timestamp,
		),
	}, nil
}

func (s *Server) RemoveMyLinkedIDP(ctx context.Context, req *auth_pb.RemoveMyLinkedIDPRequest) (*auth_pb.RemoveMyLinkedIDPResponse, error) {
	objectDetails, err := s.command.RemoveHumanExternalIDP(ctx, RemoveMyLinkedIDPRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyLinkedIDPResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

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
		Result: idp_grpc.IDPsToUserLinkPb(idps),
		MetaData: object.ToListDetails(
			idps.TotalResult,
			idps.Sequence,
			idps.Timestamp,
		),
	}, nil
}

func (s *Server) RemoveMyLinkedIDP(ctx context.Context, req *auth_pb.RemoveMyLinkedIDPRequest) (*auth_pb.RemoveMyLinkedIDPResponse, error) {
	err := s.command.RemoveHumanExternalIDP(ctx, RemoveMyLinkedIDPRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	//TODO: response from business
	return &auth_pb.RemoveMyLinkedIDPResponse{}, nil
}

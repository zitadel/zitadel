package user

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/user/v2/convert"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) GetUserByID(ctx context.Context, req *connect.Request[user.GetUserByIDRequest]) (_ *connect.Response[user.GetUserByIDResponse], err error) {
	resp, err := s.query.GetUserByIDWithPermission(ctx, true, req.Msg.GetUserId(), s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.GetUserByIDResponse{
		Details: object.DomainToDetailsPb(&domain.ObjectDetails{
			Sequence:      resp.Sequence,
			CreationDate:  resp.CreationDate,
			EventDate:     resp.ChangeDate,
			ResourceOwner: resp.ResourceOwner,
		}),
		User: convert.UserToPb(resp, s.assetAPIPrefix(ctx)),
	}), nil
}

func (s *Server) ListUsers(ctx context.Context, req *connect.Request[user.ListUsersRequest]) (*connect.Response[user.ListUsersResponse], error) {
	queries, err := convert.ListUsersRequestToModel(req.Msg)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchUsers(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.ListUsersResponse{
		Result:  convert.UsersToPb(res.Users, s.assetAPIPrefix(ctx)),
		Details: object.ToListDetails(res.SearchResponse),
	}), nil
}

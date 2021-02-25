package admin

import (
	"context"

	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) ListViews(context.Context, *admin_pb.ListViewsRequest) (*admin_pb.ListViewsResponse, error) {
	views, err := s.administrator.GetViews()
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListViewsResponse{Result: ViewsToPb(views)}, nil
}

func (s *Server) ClearView(ctx context.Context, req *admin_pb.ClearViewRequest) (*admin_pb.ClearViewResponse, error) {
	err := s.administrator.ClearView(ctx, req.Database, req.ViewName)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ClearViewResponse{}, nil
}

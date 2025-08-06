package project

import (
	"context"

	"connectrpc.com/connect"
	metadata "github.com/zitadel/zitadel/internal/api/grpc/metadata/v2beta"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) SetProjectMetadata(
	ctx context.Context, request *connect.Request[project_pb.SetProjectMetadataRequest],
) (*connect.Response[project_pb.SetProjectMetadataResponse], error) {
	result, err := s.command.BulkSetProjectMetadata(ctx, request.Msg.GetProjectId(), "", BulkSetProjectMetadataToDomain(request.Msg)...)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&project_pb.SetProjectMetadataResponse{
		SetDate: timestamppb.New(result.EventDate),
	}), nil
}

func (s *Server) ListProjectMetadata(
	ctx context.Context, request *connect.Request[project_pb.ListProjectMetadataRequest],
) (*connect.Response[project_pb.ListProjectMetadataResponse], error) {
	metadataQueries, err := ListProjectMetadataToDomain(s.systemDefaults, request.Msg)
	if err != nil {
		return nil, err
	}

	res, err := s.query.SearchProjectMetadata(ctx, true, request.Msg.GetProjectId(), metadataQueries)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&project_pb.ListProjectMetadataResponse{
		Metadata: metadata.ProjectMetadataListToPb(res.Metadata),
		Pagination: &filter.PaginationResponse{
			TotalResult:  res.Count,
			AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
		},
	}), nil
}

func (s *Server) DeleteProjectMetadata(
	ctx context.Context, request *connect.Request[project_pb.DeleteProjectMetadataRequest],
) (*connect.Response[project_pb.DeleteProjectMetadataResponse], error) {
	result, err := s.command.BulkRemoveProjectMetadata(ctx, request.Msg.GetProjectId(), "", request.Msg.GetKeys()...)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&project_pb.DeleteProjectMetadataResponse{
		DeletionDate: timestamppb.New(result.EventDate),
	}), nil
}

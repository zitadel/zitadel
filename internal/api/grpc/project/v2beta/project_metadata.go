package project

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	metadata "github.com/zitadel/zitadel/internal/api/grpc/metadata/v2beta"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) SetProjectMetadata(ctx context.Context, request *project_pb.SetProjectMetadataRequest) (*project_pb.SetProjectMetadataResponse, error) {
	result, err := s.command.BulkSetProjectMetadata(ctx, request.GetProjectId(), authz.GetCtxData(ctx).OrgID, BulkSetProjectMetadataToDomain(request)...)
	if err != nil {
		return nil, err
	}

	return &project_pb.SetProjectMetadataResponse{
		SetDate: timestamppb.New(result.EventDate),
	}, nil
}

func (s *Server) ListProjectMetadata(ctx context.Context, request *project_pb.ListProjectMetadataRequest) (*project_pb.ListProjectMetadataResponse, error) {
	metadataQueries, err := ListProjectMetadataToDomain(s.systemDefaults, request)
	if err != nil {
		return nil, err
	}

	res, err := s.query.SearchProjectMetadata(ctx, true, request.GetProjectId(), metadataQueries, false)
	if err != nil {
		return nil, err
	}

	return &project_pb.ListProjectMetadataResponse{
		Metadata: metadata.ProjectMetadataListToPb(res.Metadata),
		Pagination: &filter.PaginationResponse{
			TotalResult:  res.Count,
			AppliedLimit: uint64(request.GetPagination().GetLimit()),
		},
	}, nil
}

func (s *Server) DeleteProjectMetadata(ctx context.Context, request *project_pb.DeleteProjectMetadataRequest) (*project_pb.DeleteProjectMetadataResponse, error) {
	result, err := s.command.BulkRemoveProjectMetadata(ctx, request.GetProjectId(), authz.GetCtxData(ctx).OrgID, request.GetKeys()...)
	if err != nil {
		return nil, err
	}

	return &project_pb.DeleteProjectMetadataResponse{
		DeletionDate: timestamppb.New(result.EventDate),
	}, nil
}

package instance

import (
	"context"

	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetInstance(ctx context.Context, _ *emptypb.Empty) (*instance.GetInstanceResponse, error) {
	inst, err := s.query.Instance(ctx, true)
	if err != nil {
		return nil, err
	}

	return &instance.GetInstanceResponse{
		Instance: ToProtoObject(inst),
	}, nil
}

func (s *Server) ListInstances(ctx context.Context, req *instance.ListInstancesRequest) (*instance.ListInstancesResponse, error) {
	queries, err := ListInstancesRequestToModel(req, s.systemDefaults)
	if err != nil {
		return nil, err
	}

	instances, err := s.query.SearchInstances(ctx, queries)
	if err != nil {
		return nil, err
	}

	return &instance.ListInstancesResponse{
		Instances: InstancesToPb(instances.Instances),
		Pagination: &filter.PaginationResponse{
			TotalResult:  instances.Count,
			AppliedLimit: uint64(req.GetPagination().GetLimit()),
		},
	}, nil
}

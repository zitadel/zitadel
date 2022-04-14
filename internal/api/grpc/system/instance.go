package system

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	instance_grpc "github.com/caos/zitadel/internal/api/grpc/instance"
	"github.com/caos/zitadel/internal/api/grpc/object"
	system_pb "github.com/caos/zitadel/pkg/grpc/system"
)

func (s *Server) ListInstances(ctx context.Context, req *system_pb.ListInstancesRequest) (*system_pb.ListInstancesResponse, error) {
	//queries, err := ListInstancesRequestToModel(req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//domains, err := s.query.Instance(ctx, queries)
	//if err != nil {
	//	return nil, err
	//}
	//return &system_pb.ListInstancesResponse{
	//	Result: instance_grpc.DomainsToPb(domains.Domains),
	//	Details: object.ToListDetails(
	//		domains.Count,
	//		domains.Sequence,
	//		domains.Timestamp,
	//	),
	//}, nil
	return nil, nil
}

func (s *Server) GetInstance(ctx context.Context, req *system_pb.GetInstanceRequest) (*system_pb.GetInstanceResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.Id)
	instance, err := s.query.Instance(ctx)
	if err != nil {
		return nil, err
	}
	return &system_pb.GetInstanceResponse{
		Instance: instance_grpc.InstanceToPb(instance),
	}, nil
}

func (s *Server) AddInstance(ctx context.Context, req *system_pb.AddInstanceRequest) (*system_pb.AddInstanceResponse, error) {
	//TODO: Add instance command
	return nil, nil
}

func (s *Server) ListDomains(ctx context.Context, req *system_pb.ListDomainsRequest) (*system_pb.ListDomainsResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.Id)
	queries, err := ListInstanceDomainsRequestToModel(req)
	if err != nil {
		return nil, err
	}

	domains, err := s.query.SearchInstanceDomains(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &system_pb.ListDomainsResponse{
		Result: instance_grpc.DomainsToPb(domains.Domains),
		Details: object.ToListDetails(
			domains.Count,
			domains.Sequence,
			domains.Timestamp,
		),
	}, nil
}

func (s *Server) AddDomain(ctx context.Context, req *system_pb.AddDomainRequest) (*system_pb.AddDomainResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.Id)
	details, err := s.command.AddInstanceDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &system_pb.AddDomainResponse{
		Details: object.AddToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveDomain(ctx context.Context, req *system_pb.RemoveDomainRequest) (*system_pb.RemoveDomainResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.Id)
	details, err := s.command.RemoveInstanceDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &system_pb.RemoveDomainResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, nil
}

func (s *Server) SetPrimaryDomain(ctx context.Context, req *system_pb.SetPrimaryDomainRequest) (*system_pb.SetPrimaryDomainResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.Id)
	details, err := s.command.SetPrimaryInstanceDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &system_pb.SetPrimaryDomainResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, nil
}

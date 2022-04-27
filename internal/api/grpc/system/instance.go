package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	instance_grpc "github.com/zitadel/zitadel/internal/api/grpc/instance"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) ListInstances(ctx context.Context, req *system_pb.ListInstancesRequest) (*system_pb.ListInstancesResponse, error) {
	queries, err := ListInstancesRequestToModel(req)
	if err != nil {
		return nil, err
	}

	result, err := s.query.SearchInstances(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &system_pb.ListInstancesResponse{
		Result: instance_grpc.InstancesToPb(result.Instances),
		Details: &object_pb.ListDetails{
			TotalResult: result.Count,
		},
	}, nil
}

func (s *Server) GetInstance(ctx context.Context, req *system_pb.GetInstanceRequest) (*system_pb.GetInstanceResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.InstanceId)
	instance, err := s.query.Instance(ctx)
	if err != nil {
		return nil, err
	}
	return &system_pb.GetInstanceResponse{
		Instance: instance_grpc.InstanceDetailToPb(instance),
	}, nil
}

func (s *Server) AddInstance(ctx context.Context, req *system_pb.AddInstanceRequest) (*system_pb.AddInstanceResponse, error) {
	id, details, err := s.command.SetUpInstance(ctx, AddInstancePbToSetupInstance(req, s.DefaultInstance), s.ExternalSecure)
	if err != nil {
		return nil, err
	}
	return &system_pb.AddInstanceResponse{
		InstanceId: id,
		Details: object.AddToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, nil
	return nil, nil
}

func (s *Server) ExistsDomain(ctx context.Context, req *system_pb.ExistsDomainRequest) (*system_pb.ExistsDomainResponse, error) {
	domainQuery, err := query.NewInstanceDomainDomainSearchQuery(query.TextEquals, req.Domain)

	query := &query.InstanceDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: 0,
			Limit:  1,
			Asc:    true,
		},
		Queries: []query.SearchQuery{
			domainQuery,
		},
	}
	domains, err := s.query.SearchInstanceDomains(ctx, query)
	if err != nil {
		return nil, err
	}
	return &system_pb.ExistsDomainResponse{
		Exists: domains.Count > 0,
	}, nil
}

func (s *Server) ListDomains(ctx context.Context, req *system_pb.ListDomainsRequest) (*system_pb.ListDomainsResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.InstanceId)
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
	ctx = authz.WithInstanceID(ctx, req.InstanceId)
	instance, err := s.query.Instance(ctx)
	if err != nil {
		return nil, err
	}
	ctx = authz.WithInstance(ctx, instance)
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
	ctx = authz.WithInstanceID(ctx, req.InstanceId)
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
	ctx = authz.WithInstanceID(ctx, req.InstanceId)
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

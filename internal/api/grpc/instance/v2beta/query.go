package instance

import (
	"context"

	"connectrpc.com/connect"

	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func (s *Server) GetInstance(ctx context.Context, _ *connect.Request[instance.GetInstanceRequest]) (*connect.Response[instance.GetInstanceResponse], error) {
	inst, err := s.query.Instance(ctx, true)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.GetInstanceResponse{
		Instance: ToProtoObject(inst),
	}), nil
}

func (s *Server) ListInstances(ctx context.Context, req *connect.Request[instance.ListInstancesRequest]) (*connect.Response[instance.ListInstancesResponse], error) {
	queries, err := ListInstancesRequestToModel(req.Msg, s.systemDefaults)
	if err != nil {
		return nil, err
	}

	instances, err := s.query.SearchInstances(ctx, queries)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.ListInstancesResponse{
		Instances:  InstancesToPb(instances.Instances),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, instances.SearchResponse),
	}), nil
}

func (s *Server) ListCustomDomains(ctx context.Context, req *connect.Request[instance.ListCustomDomainsRequest]) (*connect.Response[instance.ListCustomDomainsResponse], error) {
	queries, err := ListCustomDomainsRequestToModel(req.Msg, s.systemDefaults)
	if err != nil {
		return nil, err
	}

	domains, err := s.query.SearchInstanceDomains(ctx, queries)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.ListCustomDomainsResponse{
		Domains:    DomainsToPb(domains.Domains),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, domains.SearchResponse),
	}), nil
}

func (s *Server) ListTrustedDomains(ctx context.Context, req *connect.Request[instance.ListTrustedDomainsRequest]) (*connect.Response[instance.ListTrustedDomainsResponse], error) {
	queries, err := ListTrustedDomainsRequestToModel(req.Msg, s.systemDefaults)
	if err != nil {
		return nil, err
	}

	domains, err := s.query.SearchInstanceTrustedDomains(ctx, queries)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.ListTrustedDomainsResponse{
		TrustedDomain: trustedDomainsToPb(domains.Domains),
		Pagination:    filter.QueryToPaginationPb(queries.SearchRequest, domains.SearchResponse),
	}), nil
}

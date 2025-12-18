package instance

import (
	"context"

	"connectrpc.com/connect"

	instancev2 "github.com/zitadel/zitadel/backend/v3/api/instance/v2"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func (s *Server) GetInstance(ctx context.Context, req *connect.Request[instance.GetInstanceRequest]) (*connect.Response[instance.GetInstanceResponse], error) {
	if authz.GetFeatures(ctx).EnableRelationalTables {
		return instancev2.GetInstance(ctx, req)
	}
	if err := s.checkPermission(ctx, domain.PermissionSystemInstanceRead, domain.PermissionInstanceRead); err != nil {
		return nil, err
	}
	inst, err := s.query.Instance(ctx, true)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.GetInstanceResponse{
		Instance: ToProtoObject(inst),
	}), nil
}

func (s *Server) ListInstances(ctx context.Context, req *connect.Request[instance.ListInstancesRequest]) (*connect.Response[instance.ListInstancesResponse], error) {
	if authz.GetFeatures(ctx).EnableRelationalTables {
		return instancev2.ListInstances(ctx, req)
	}

	// List instances is currently only allowed with system permissions,
	// so we directly check for them in the auth interceptor and do not check here again.
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
	if authz.GetFeatures(ctx).EnableRelationalTables {
		return instancev2.ListCustomDomains(ctx, req)
	}

	if err := s.checkPermission(ctx, domain.PermissionSystemInstanceRead, domain.PermissionInstanceRead); err != nil {
		return nil, err
	}
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
	if authz.GetFeatures(ctx).EnableRelationalTables {
		return instancev2.ListTrustedDomains(ctx, req)
	}

	if err := s.checkPermission(ctx, domain.PermissionSystemInstanceRead, domain.PermissionInstanceRead); err != nil {
		return nil, err
	}
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

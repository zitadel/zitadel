package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	instance_grpc "github.com/zitadel/zitadel/internal/api/grpc/instance"
	"github.com/zitadel/zitadel/internal/api/grpc/member"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
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
	instance, err := s.query.Instance(ctx, true)
	if err != nil {
		return nil, err
	}
	return &system_pb.GetInstanceResponse{
		Instance: instance_grpc.InstanceDetailToPb(instance),
	}, nil
}

func (s *Server) AddInstance(ctx context.Context, req *system_pb.AddInstanceRequest) (*system_pb.AddInstanceResponse, error) {
	id, _, _, details, err := s.command.SetUpInstance(ctx, AddInstancePbToSetupInstance(req, s.defaultInstance, s.externalDomain))
	if err != nil {
		return nil, err
	}
	return &system_pb.AddInstanceResponse{
		InstanceId: id,
		Details:    object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) UpdateInstance(ctx context.Context, req *system_pb.UpdateInstanceRequest) (*system_pb.UpdateInstanceResponse, error) {
	details, err := s.command.UpdateInstance(ctx, req.InstanceName)
	if err != nil {
		return nil, err
	}
	return &system_pb.UpdateInstanceResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) CreateInstance(ctx context.Context, req *system_pb.CreateInstanceRequest) (*system_pb.CreateInstanceResponse, error) {
	id, pat, key, details, err := s.command.SetUpInstance(ctx, CreateInstancePbToSetupInstance(req, s.defaultInstance, s.externalDomain))
	if err != nil {
		return nil, err
	}

	var machineKey []byte
	if key != nil {
		machineKey, err = key.Detail()
		if err != nil {
			return nil, err
		}
	}

	return &system_pb.CreateInstanceResponse{
		Pat:        pat,
		MachineKey: machineKey,
		InstanceId: id,
		Details:    object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) RemoveInstance(ctx context.Context, req *system_pb.RemoveInstanceRequest) (*system_pb.RemoveInstanceResponse, error) {
	details, err := s.command.RemoveInstance(ctx, req.InstanceId)
	if err != nil {
		return nil, err
	}
	return &system_pb.RemoveInstanceResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) ListIAMMembers(ctx context.Context, req *system_pb.ListIAMMembersRequest) (*system_pb.ListIAMMembersResponse, error) {
	queries, err := ListIAMMembersRequestToQuery(req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.IAMMembers(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &system_pb.ListIAMMembersResponse{
		Details: object.ToListDetails(res.Count, res.Sequence, res.Timestamp),
		//TODO: resource owner of user of the member instead of the membership resource owner
		Result: member.MembersToPb("", res.Members),
	}, nil
}

func (s *Server) ExistsDomain(ctx context.Context, req *system_pb.ExistsDomainRequest) (*system_pb.ExistsDomainResponse, error) {
	domainQuery, err := query.NewInstanceDomainDomainSearchQuery(query.TextEqualsIgnoreCase, req.Domain)
	if err != nil {
		return nil, err
	}

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
	domains, err := s.query.SearchInstanceDomainsGlobal(ctx, query)
	if err != nil {
		return nil, err
	}
	return &system_pb.ExistsDomainResponse{
		Exists: domains.Count > 0,
	}, nil
}

func (s *Server) ListDomains(ctx context.Context, req *system_pb.ListDomainsRequest) (*system_pb.ListDomainsResponse, error) {
	queries, err := ListInstanceDomainsRequestToModel(req)
	if err != nil {
		return nil, err
	}

	domains, err := s.query.SearchInstanceDomains(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &system_pb.ListDomainsResponse{
		Result:  instance_grpc.DomainsToPb(domains.Domains),
		Details: object.ToListDetails(domains.Count, domains.Sequence, domains.Timestamp),
	}, nil
}

func (s *Server) AddDomain(ctx context.Context, req *system_pb.AddDomainRequest) (*system_pb.AddDomainResponse, error) {
	instance, err := s.query.Instance(ctx, true)
	if err != nil {
		return nil, err
	}
	ctx = authz.WithInstance(ctx, instance)

	details, err := s.command.AddInstanceDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &system_pb.AddDomainResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) RemoveDomain(ctx context.Context, req *system_pb.RemoveDomainRequest) (*system_pb.RemoveDomainResponse, error) {
	details, err := s.command.RemoveInstanceDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &system_pb.RemoveDomainResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) SetPrimaryDomain(ctx context.Context, req *system_pb.SetPrimaryDomainRequest) (*system_pb.SetPrimaryDomainResponse, error) {
	details, err := s.command.SetPrimaryInstanceDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &system_pb.SetPrimaryDomainResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

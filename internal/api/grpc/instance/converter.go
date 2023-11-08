package org

import (
	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	instance_pb "github.com/zitadel/zitadel/pkg/grpc/instance"
)

func InstancesToPb(instances []*query.Instance) []*instance_pb.Instance {
	list := make([]*instance_pb.Instance, len(instances))
	for i, instance := range instances {
		list[i] = InstanceToPb(instance)
	}
	return list
}

func InstanceToPb(instance *query.Instance) *instance_pb.Instance {
	return &instance_pb.Instance{
		Details: object.ToViewDetailsPb(
			instance.Sequence,
			instance.CreationDate,
			instance.ChangeDate,
			instance.InstanceID(),
		),
		Id:      instance.InstanceID(),
		Name:    instance.Name,
		Domains: DomainsToPb(instance.Domains),
		Version: build.Version(),
		State:   instance_pb.State_STATE_RUNNING, //TODO: change when delete is implemented
	}
}

func InstanceDetailToPb(instance *query.Instance) *instance_pb.InstanceDetail {
	return &instance_pb.InstanceDetail{
		Details: object.ToViewDetailsPb(
			instance.Sequence,
			instance.CreationDate,
			instance.ChangeDate,
			instance.InstanceID(),
		),
		Id:      instance.InstanceID(),
		Name:    instance.Name,
		Domains: DomainsToPb(instance.Domains),
		Version: build.Version(),
		State:   instance_pb.State_STATE_RUNNING, //TODO: change when delete is implemented
	}
}

func InstanceQueriesToModel(queries []*instance_pb.Query) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = InstanceQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func InstanceQueryToModel(searchQuery *instance_pb.Query) (query.SearchQuery, error) {
	switch q := searchQuery.Query.(type) {
	case *instance_pb.Query_IdQuery:
		return query.NewInstanceIDsListSearchQuery(q.IdQuery.Ids...)
	case *instance_pb.Query_DomainQuery:
		return query.NewInstanceDomainsListSearchQuery(q.DomainQuery.Domains...)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "INST-3m0se", "List.Query.Invalid")
	}
}

func DomainQueriesToModel(queries []*instance_pb.DomainSearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = DomainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func DomainQueryToModel(searchQuery *instance_pb.DomainSearchQuery) (query.SearchQuery, error) {
	switch q := searchQuery.Query.(type) {
	case *instance_pb.DomainSearchQuery_DomainQuery:
		return query.NewInstanceDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainQuery.Method), q.DomainQuery.Domain)
	case *instance_pb.DomainSearchQuery_GeneratedQuery:
		return query.NewInstanceDomainGeneratedSearchQuery(q.GeneratedQuery.Generated)
	case *instance_pb.DomainSearchQuery_PrimaryQuery:
		return query.NewInstanceDomainPrimarySearchQuery(q.PrimaryQuery.Primary)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "INST-Ags42", "List.Query.Invalid")
	}
}

func DomainsToPb(domains []*query.InstanceDomain) []*instance_pb.Domain {
	d := make([]*instance_pb.Domain, len(domains))
	for i, domain := range domains {
		d[i] = DomainToPb(domain)
	}
	return d
}

func DomainToPb(d *query.InstanceDomain) *instance_pb.Domain {
	return &instance_pb.Domain{
		Domain:    d.Domain,
		Primary:   d.IsPrimary,
		Generated: d.IsGenerated,
		Details: object.ToViewDetailsPb(
			d.Sequence,
			d.CreationDate,
			d.ChangeDate,
			d.InstanceID,
		),
	}
}

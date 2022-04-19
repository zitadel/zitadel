package org

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query"
	instance_pb "github.com/caos/zitadel/pkg/grpc/instance"
)

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
		return nil, errors.ThrowInvalidArgument(nil, "ORG-Ags42", "List.Query.Invalid")
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

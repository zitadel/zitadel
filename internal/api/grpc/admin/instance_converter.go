package admin

import (
	instance_grpc "github.com/zitadel/zitadel/internal/api/grpc/instance"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func ListInstanceDomainsRequestToModel(req *admin_pb.ListInstanceDomainsRequest) (*query.InstanceDomainSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := instance_grpc.DomainQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.InstanceDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

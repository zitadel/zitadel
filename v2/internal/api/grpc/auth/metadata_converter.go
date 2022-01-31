package auth

import (
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/pkg/grpc/auth"
	"github.com/caos/zitadel/v2/internal/api/grpc/metadata"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
)

func ListUserMetadataToQuery(req *auth.ListMyMetadataRequest) (*query.UserMetadataSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := metadata.MetadataQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.UserMetadataSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

package project

import (
	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	metadata "github.com/zitadel/zitadel/internal/api/grpc/metadata/v2beta"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	v2beta_project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func BulkSetProjectMetadataToDomain(req *v2beta_project.SetProjectMetadataRequest) []*domain.Metadata {
	metadata := make([]*domain.Metadata, len(req.Metadata))
	for i, data := range req.Metadata {
		metadata[i] = &domain.Metadata{
			Key:   data.Key,
			Value: data.Value,
		}
	}
	return metadata
}

func ListProjectMetadataToDomain(systemDefaults systemdefaults.SystemDefaults, request *v2beta_project.ListProjectMetadataRequest) (*query.ProjectMetadataSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(systemDefaults, request.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := metadata.ProjectMetadataQueriesToQuery(request.Filter)
	if err != nil {
		return nil, err
	}
	return &query.ProjectMetadataSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

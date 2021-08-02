package auth

import (
	"github.com/caos/zitadel/internal/api/grpc/metadata"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func BulkSetMetadataToDomain(req *auth.BulkSetMyMetadataRequest) []*domain.Metadata {
	metaData := make([]*domain.Metadata, len(req.Metadata))
	for i, data := range req.Metadata {
		metaData[i] = &domain.Metadata{
			Key:   data.Key,
			Value: data.Value,
		}
	}
	return metaData
}

func ListUserMetadataToDomain(req *auth.ListMyMetadataRequest) *domain.MetadataSearchRequest {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &domain.MetadataSearchRequest{
		Offset:  offset,
		Limit:   limit,
		Asc:     asc,
		Queries: metadata.MetadataQueriesToModel(req.Queries),
	}
}

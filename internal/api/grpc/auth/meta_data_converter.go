package auth

import (
	"github.com/caos/zitadel/internal/api/grpc/metadata"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func BulkSetMetaDataToDomain(req *auth.BulkSetMyMetaDataRequest) []*domain.MetaData {
	metaData := make([]*domain.MetaData, len(req.MetaData))
	for i, data := range req.MetaData {
		metaData[i] = &domain.MetaData{
			Key:   data.Key,
			Value: data.Value,
		}
	}
	return metaData
}

func ListUserMetaDataToDomain(req *auth.ListMyMetaDataRequest) *domain.MetaDataSearchRequest {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &domain.MetaDataSearchRequest{
		Offset:  offset,
		Limit:   limit,
		Asc:     asc,
		Queries: metadata.MetaDataQueriesToModel(req.Queries),
	}
}

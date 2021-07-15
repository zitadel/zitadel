package auth

import (
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

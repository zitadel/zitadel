package metadata

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	meta_pb "github.com/caos/zitadel/pkg/grpc/metadata"
)

func MetadataListToPb(dataList []*domain.Metadata) []*meta_pb.Metadata {
	mds := make([]*meta_pb.Metadata, len(dataList))
	for i, data := range dataList {
		mds[i] = DomainMetadataToPb(data)
	}
	return mds
}

func DomainMetadataToPb(data *domain.Metadata) *meta_pb.Metadata {
	return &meta_pb.Metadata{
		Key:   data.Key,
		Value: data.Value,
		Details: object.ToViewDetailsPb(
			data.Sequence,
			data.CreationDate,
			data.ChangeDate,
			data.ResourceOwner,
		),
	}
}

func MetadataQueriesToModel(queries []*meta_pb.MetadataQuery) []*domain.MetadataSearchQuery {
	q := make([]*domain.MetadataSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = MetadataQueryToModel(query)
	}
	return q
}

func MetadataQueryToModel(query *meta_pb.MetadataQuery) *domain.MetadataSearchQuery {
	switch q := query.Query.(type) {
	case *meta_pb.MetadataQuery_KeyQuery:
		return MetadataKeyQueryToModel(q.KeyQuery)
	case *meta_pb.MetadataQuery_ValueQuery:
		return MetadataValueQueryToModel(q.ValueQuery)
	default:
		return nil
	}
}

func MetadataKeyQueryToModel(q *meta_pb.MetadataKeyQuery) *domain.MetadataSearchQuery {
	return &domain.MetadataSearchQuery{
		Key:    domain.MetadataSearchKeyKey,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.Key,
	}
}

func MetadataValueQueryToModel(q *meta_pb.MetadataValueQuery) *domain.MetadataSearchQuery {
	return &domain.MetadataSearchQuery{
		Key:    domain.MetadataSearchKeyValue,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.Value,
	}
}

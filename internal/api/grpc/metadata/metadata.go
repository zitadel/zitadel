package metadata

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	meta_pb "github.com/caos/zitadel/pkg/grpc/metadata"
)

func MetaDataListToPb(dataList []*domain.MetaData) []*meta_pb.MetaData {
	u := make([]*meta_pb.MetaData, len(dataList))
	for i, data := range dataList {
		u[i] = DomainMetaDataToPb(data)
	}
	return u
}

func DomainMetaDataToPb(data *domain.MetaData) *meta_pb.MetaData {
	return &meta_pb.MetaData{
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

func MetaDataQueriesToModel(queries []*meta_pb.MetaDataQuery) []*domain.MetaDataSearchQuery {
	q := make([]*domain.MetaDataSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = MetaDataQueryToModel(query)
	}
	return q
}

func MetaDataQueryToModel(query *meta_pb.MetaDataQuery) *domain.MetaDataSearchQuery {
	switch q := query.Query.(type) {
	case *meta_pb.MetaDataQuery_KeyQuery:
		return MetaDataKeyQueryToModel(q.KeyQuery)
	case *meta_pb.MetaDataQuery_ValueQuery:
		return MetaDataValueQueryToModel(q.ValueQuery)
	default:
		return nil
	}
}

func MetaDataKeyQueryToModel(q *meta_pb.MetaDataKeyQuery) *domain.MetaDataSearchQuery {
	return &domain.MetaDataSearchQuery{
		Key:    domain.MetaDataSearchKeyKey,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.Key,
	}
}

func MetaDataValueQueryToModel(q *meta_pb.MetaDataValueQuery) *domain.MetaDataSearchQuery {
	return &domain.MetaDataSearchQuery{
		Key:    domain.MetaDataSearchKeyValue,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.Value,
	}
}

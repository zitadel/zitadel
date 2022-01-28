package metadata

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query"
	meta_pb "github.com/caos/zitadel/pkg/grpc/metadata"
)

func MetadataListToPb(dataList []*query.UserMetadata) []*meta_pb.Metadata {
	mds := make([]*meta_pb.Metadata, len(dataList))
	for i, data := range dataList {
		mds[i] = DomainMetadataToPb(data)
	}
	return mds
}

func DomainMetadataToPb(data *query.UserMetadata) *meta_pb.Metadata {
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

func MetadataQueriesToQuery(queries []*meta_pb.MetadataQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = MetadataQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func MetadataQueryToQuery(query *meta_pb.MetadataQuery) (query.SearchQuery, error) {
	switch q := query.Query.(type) {
	case *meta_pb.MetadataQuery_KeyQuery:
		return MetadataKeyQueryToQuery(q.KeyQuery)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "METAD-fdg23", "List.Query.Invalid")
	}
}

func MetadataKeyQueryToQuery(q *meta_pb.MetadataKeyQuery) (query.SearchQuery, error) {
	return query.NewUserMetadataKeySearchQuery(q.Key, object.TextMethodToQuery(q.Method))
}

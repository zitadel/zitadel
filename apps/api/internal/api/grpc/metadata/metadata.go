package metadata

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	meta_pb "github.com/zitadel/zitadel/pkg/grpc/metadata"
)

func UserMetadataListToPb(dataList []*query.UserMetadata) []*meta_pb.Metadata {
	mds := make([]*meta_pb.Metadata, len(dataList))
	for i, data := range dataList {
		mds[i] = UserMetadataToPb(data)
	}
	return mds
}

func UserMetadataToPb(data *query.UserMetadata) *meta_pb.Metadata {
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

func OrgMetadataListToPb(dataList []*query.OrgMetadata) []*meta_pb.Metadata {
	mds := make([]*meta_pb.Metadata, len(dataList))
	for i, data := range dataList {
		mds[i] = OrgMetadataToPb(data)
	}
	return mds
}

func OrgMetadataToPb(data *query.OrgMetadata) *meta_pb.Metadata {
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

func OrgMetadataQueriesToQuery(queries []*meta_pb.MetadataQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = OrgMetadataQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func OrgMetadataQueryToQuery(metadataQuery *meta_pb.MetadataQuery) (query.SearchQuery, error) {
	switch q := metadataQuery.Query.(type) {
	case *meta_pb.MetadataQuery_KeyQuery:
		return query.NewOrgMetadataKeySearchQuery(q.KeyQuery.Key, object.TextMethodToQuery(q.KeyQuery.Method))
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "METAD-fdg23", "List.Query.Invalid")
	}
}

func UserMetadataQueriesToQuery(queries []*meta_pb.MetadataQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = UserMetadataQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func UserMetadataQueryToQuery(metadataQuery *meta_pb.MetadataQuery) (query.SearchQuery, error) {
	switch q := metadataQuery.Query.(type) {
	case *meta_pb.MetadataQuery_KeyQuery:
		return query.NewUserMetadataKeySearchQuery(q.KeyQuery.Key, object.TextMethodToQuery(q.KeyQuery.Method))
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "METAD-Vn7qy", "List.Query.Invalid")
	}
}

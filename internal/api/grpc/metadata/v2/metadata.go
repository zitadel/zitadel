package metadata

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	filter_v2 "github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	meta_pb "github.com/zitadel/zitadel/pkg/grpc/metadata/v2"
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
		Key:          data.Key,
		Value:        data.Value,
		CreationDate: timestamppb.New(data.CreationDate),
		ChangeDate:   timestamppb.New(data.ChangeDate),
	}
}

func UserMetadataFiltersToQuery(queries []*meta_pb.MetadataSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = UserMetadataFilterToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func UserMetadataFilterToQuery(filter *meta_pb.MetadataSearchFilter) (query.SearchQuery, error) {
	switch q := filter.Filter.(type) {
	case *meta_pb.MetadataSearchFilter_KeyFilter:
		return query.NewUserMetadataKeySearchQuery(q.KeyFilter.Key, filter_v2.TextMethodPbToQuery(q.KeyFilter.Method))
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "METAD-fdg23", "List.Query.Invalid")
	}
}

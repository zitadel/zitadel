package user

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter_pb "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) ListKeys(ctx context.Context, req *connect.Request[user.ListKeysRequest]) (*connect.Response[user.ListKeysResponse], error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Msg.GetPagination())
	if err != nil {
		return nil, err
	}

	filters, err := keyFiltersToQueries(req.Msg.GetFilters())
	if err != nil {
		return nil, err
	}
	search := &query.AuthNKeySearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: authnKeyFieldNameToSortingColumn(req.Msg.SortingColumn),
		},
		Queries: filters,
	}
	result, err := s.query.SearchAuthNKeys(ctx, search, query.JoinFilterUserMachine, s.checkPermission)
	if err != nil {
		return nil, err
	}
	resp := &user.ListKeysResponse{
		Result:     make([]*user.Key, len(result.AuthNKeys)),
		Pagination: filter.QueryToPaginationPb(search.SearchRequest, result.SearchResponse),
	}
	for i, key := range result.AuthNKeys {
		resp.Result[i] = &user.Key{
			CreationDate:   timestamppb.New(key.CreationDate),
			ChangeDate:     timestamppb.New(key.ChangeDate),
			Id:             key.ID,
			UserId:         key.AggregateID,
			OrganizationId: key.ResourceOwner,
			ExpirationDate: timestamppb.New(key.Expiration),
		}
	}
	return connect.NewResponse(resp), nil
}

func keyFiltersToQueries(filters []*user.KeysSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(filters))
	for i, filter := range filters {
		q[i], err = keyFilterToQuery(filter)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func keyFilterToQuery(filter *user.KeysSearchFilter) (query.SearchQuery, error) {
	switch q := filter.Filter.(type) {
	case *user.KeysSearchFilter_CreatedDateFilter:
		return authnKeyCreatedFilterToQuery(q.CreatedDateFilter)
	case *user.KeysSearchFilter_ExpirationDateFilter:
		return authnKeyExpirationFilterToQuery(q.ExpirationDateFilter)
	case *user.KeysSearchFilter_KeyIdFilter:
		return authnKeyIdFilterToQuery(q.KeyIdFilter)
	case *user.KeysSearchFilter_UserIdFilter:
		return authnKeyUserIdFilterToQuery(q.UserIdFilter)
	case *user.KeysSearchFilter_OrganizationIdFilter:
		return authnKeyOrgIdFilterToQuery(q.OrganizationIdFilter)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func authnKeyIdFilterToQuery(f *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyIDQuery(f.Id)
}

func authnKeyUserIdFilterToQuery(f *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyIdentifyerQuery(f.Id)
}

func authnKeyOrgIdFilterToQuery(f *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyResourceOwnerQuery(f.Id)
}

func authnKeyCreatedFilterToQuery(f *filter_pb.TimestampFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyCreationDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
}

func authnKeyExpirationFilterToQuery(f *filter_pb.TimestampFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyExpirationDateDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
}

// authnKeyFieldNameToSortingColumn defaults to the creation date because this ensures deterministic pagination
func authnKeyFieldNameToSortingColumn(field *user.KeyFieldName) query.Column {
	if field == nil {
		return query.AuthNKeyColumnCreationDate
	}
	switch *field {
	case user.KeyFieldName_KEY_FIELD_NAME_UNSPECIFIED:
		return query.AuthNKeyColumnCreationDate
	case user.KeyFieldName_KEY_FIELD_NAME_ID:
		return query.AuthNKeyColumnID
	case user.KeyFieldName_KEY_FIELD_NAME_USER_ID:
		return query.AuthNKeyColumnIdentifier
	case user.KeyFieldName_KEY_FIELD_NAME_ORGANIZATION_ID:
		return query.AuthNKeyColumnResourceOwner
	case user.KeyFieldName_KEY_FIELD_NAME_CREATED_DATE:
		return query.AuthNKeyColumnCreationDate
	case user.KeyFieldName_KEY_FIELD_NAME_KEY_EXPIRATION_DATE:
		return query.AuthNKeyColumnExpiration
	default:
		return query.AuthNKeyColumnCreationDate
	}
}

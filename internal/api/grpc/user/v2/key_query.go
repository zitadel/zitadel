package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) ListKeys(ctx context.Context, req *user.ListKeysRequest) (*user.ListKeysResponse, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	if err != nil {
		return nil, err
	}

	filters, err := keyFiltersToQueries(req.Filters, 0)
	if err != nil {
		return nil, err
	}
	search := &query.AuthNKeySearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: authnKeyFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: filters,
	}
	result, err := s.query.SearchAuthNKeys(ctx, search, query.JoinFilterUserMachine)
	if err != nil {
		return nil, err
	}
	resp := &user.ListKeysResponse{
		Result:     make([]*user.Key, len(result.AuthNKeys)),
		Pagination: filter.QueryToPaginationPb(search.SearchRequest, result.SearchResponse),
	}
	for i := range result.AuthNKeys {
		key := result.AuthNKeys[i]
		resp.Result[i] = &user.Key{
			CreationDate:   timestamppb.New(key.CreationDate),
			ChangeDate:     timestamppb.New(key.ChangeDate),
			Id:             key.ID,
			UserId:         key.AggregateID,
			OrganizationId: key.ResourceOwner,
			ExpirationDate: timestamppb.New(key.Expiration),
		}
	}
	return resp, nil
}

func keyFiltersToQueries(filters []*user.KeysSearchFilter, level uint8) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(filters))
	for i, filter := range filters {
		q[i], err = keyFilterToQuery(filter, level)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func keyFilterToQuery(filter *user.KeysSearchFilter, level uint8) (query.SearchQuery, error) {
	if level > 20 {
		// can't go deeper than 20 levels of nesting.
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-zsQ97", "Errors.Query.TooManyNestingLevels")
	}
	switch q := filter.Filter.(type) {
	case *user.KeysSearchFilter_CreatedDateFilter:
		return authnKeyCreatedFilterToQuery(q.CreatedDateFilter)
	case *user.KeysSearchFilter_ChangedDateFilter:
		return authnKeyChangedFilterToQuery(q.ChangedDateFilter)
	case *user.KeysSearchFilter_ExpirationDateFilter:
		return authnKeyExpirationFilterToQuery(q.ExpirationDateFilter)
	case *user.KeysSearchFilter_KeyIdFilter:
		return authnKeyIdFilterToQuery(q.KeyIdFilter)
	case *user.KeysSearchFilter_UserIdFilter:
		return authnKeyUserIdFilterToQuery(q.UserIdFilter)
	case *user.KeysSearchFilter_OrganizationIdFilter:
		return authnKeyOrgIdFilterToQuery(q.OrganizationIdFilter)
	case *user.KeysSearchFilter_OrFilter:
		return authnKeyOrFilterToQuery(q.OrFilter, level)
	case *user.KeysSearchFilter_AndFilter:
		return authnKeyAndFilterToQuery(q.AndFilter, level)
	case *user.KeysSearchFilter_NotFilter:
		return authnKeyNotFilterToQuery(q.NotFilter, level)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func authnKeyIdFilterToQuery(f *user.IDFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyObjectIDQuery(f.Id)
}

func authnKeyUserIdFilterToQuery(f *user.IDFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyAggregateIDQuery(f.Id)
}

func authnKeyOrgIdFilterToQuery(f *user.IDFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyResourceOwnerQuery(f.Id)
}

func authnKeyCreatedFilterToQuery(f *user.TimestampFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyCreationDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
}

func authnKeyChangedFilterToQuery(f *user.TimestampFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyChangedDateDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
}

func authnKeyExpirationFilterToQuery(f *user.TimestampFilter) (query.SearchQuery, error) {
	return query.NewAuthNKeyExpirationDateDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
}

func authnKeyOrFilterToQuery(q *user.KeysOrFilter, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := keyFiltersToQueries(q.Filters, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewOrQuery(mappedQueries...)
}
func authnKeyAndFilterToQuery(q *user.KeysAndFilter, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := keyFiltersToQueries(q.Filters, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewAndQuery(mappedQueries...)
}
func authnKeyNotFilterToQuery(q *user.KeysNotFilter, level uint8) (query.SearchQuery, error) {
	mappedQuery, err := keyFilterToQuery(q.Filter, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewNotQuery(mappedQuery)
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
		return query.AuthNKeyColumnObjectID
	case user.KeyFieldName_KEY_FIELD_NAME_USER_ID:
		return query.AuthNKeyColumnAggregateID
	case user.KeyFieldName_KEY_FIELD_NAME_ORGANIZATION_ID:
		return query.AuthNKeyColumnResourceOwner
	case user.KeyFieldName_KEY_FIELD_NAME_CREATED_DATE:
		return query.AuthNKeyColumnCreationDate
	case user.KeyFieldName_KEY_FIELD_NAME_CHANGED_DATE:
		return query.AuthNKeyColumnChangeDate
	case user.KeyFieldName_KEY_FIELD_NAME_KEY_EXPIRATION_DATE:
		return query.AuthNKeyColumnExpiration
	default:
		return query.AuthNKeyColumnCreationDate
	}
}

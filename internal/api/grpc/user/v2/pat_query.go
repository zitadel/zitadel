package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) ListPersonalAccessTokens(ctx context.Context, req *user.ListPersonalAccessTokensRequest) (*user.ListPersonalAccessTokensResponse, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	if err != nil {
		return nil, err
	}
	filters, err := patFiltersToQueries(req.Filters, 0)
	if err != nil {
		return nil, err
	}
	search := &query.PersonalAccessTokenSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: authnPersonalAccessTokenFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: filters,
	}
	result, err := s.query.SearchPersonalAccessTokens(ctx, search, false)
	if err != nil {
		return nil, err
	}
	resp := &user.ListPersonalAccessTokensResponse{
		Result:     make([]*user.PersonalAccessToken, len(result.PersonalAccessTokens)),
		Pagination: filter.QueryToPaginationPb(search.SearchRequest, result.SearchResponse),
	}
	for i := range result.PersonalAccessTokens {
		pat := result.PersonalAccessTokens[i]
		resp.Result[i] = &user.PersonalAccessToken{
			CreationDate:   timestamppb.New(pat.CreationDate),
			ChangeDate:     timestamppb.New(pat.ChangeDate),
			Id:             pat.ID,
			UserId:         pat.UserID,
			OrganizationId: pat.ResourceOwner,
			ExpirationDate: timestamppb.New(pat.Expiration),
		}
	}
	return resp, nil
}

func patFiltersToQueries(filters []*user.PersonalAccessTokensSearchFilter, level uint8) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(filters))
	for i, filter := range filters {
		q[i], err = patFilterToQuery(filter, level)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func patFilterToQuery(filter *user.PersonalAccessTokensSearchFilter, level uint8) (query.SearchQuery, error) {
	if level > 20 {
		// can't go deeper than 20 levels of nesting.
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-zsQ97", "Errors.Query.TooManyNestingLevels")
	}
	switch q := filter.Filter.(type) {
	case *user.PersonalAccessTokensSearchFilter_CreatedDateFilter:
		return authnPersonalAccessTokenCreatedFilterToQuery(q.CreatedDateFilter)
	case *user.PersonalAccessTokensSearchFilter_ChangedDateFilter:
		return authnPersonalAccessTokenChangedFilterToQuery(q.ChangedDateFilter)
	case *user.PersonalAccessTokensSearchFilter_ExpirationDateFilter:
		return authnPersonalAccessTokenExpirationFilterToQuery(q.ExpirationDateFilter)
	case *user.PersonalAccessTokensSearchFilter_TokenIdFilter:
		return authnPersonalAccessTokenIdFilterToQuery(q.TokenIdFilter)
	case *user.PersonalAccessTokensSearchFilter_UserIdFilter:
		return authnPersonalAccessTokenUserIdFilterToQuery(q.UserIdFilter)
	case *user.PersonalAccessTokensSearchFilter_OrganizationIdFilter:
		return authnPersonalAccessTokenOrgIdFilterToQuery(q.OrganizationIdFilter)
	case *user.PersonalAccessTokensSearchFilter_OrFilter:
		return authnPersonalAccessTokenOrFilterToQuery(q.OrFilter, level)
	case *user.PersonalAccessTokensSearchFilter_AndFilter:
		return authnPersonalAccessTokenAndFilterToQuery(q.AndFilter, level)
	case *user.PersonalAccessTokensSearchFilter_NotFilter:
		return authnPersonalAccessTokenNotFilterToQuery(q.NotFilter, level)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func authnPersonalAccessTokenIdFilterToQuery(f *user.IDFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenIDQuery(f.Id)
}

func authnPersonalAccessTokenUserIdFilterToQuery(f *user.IDFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenUserIDSearchQuery(f.Id)
}

func authnPersonalAccessTokenOrgIdFilterToQuery(f *user.IDFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenResourceOwnerSearchQuery(f.Id)
}

func authnPersonalAccessTokenCreatedFilterToQuery(f *user.TimestampFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenCreationDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
}

func authnPersonalAccessTokenChangedFilterToQuery(f *user.TimestampFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenChangedDateDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
}

func authnPersonalAccessTokenExpirationFilterToQuery(f *user.TimestampFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenExpirationDateDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
}

func authnPersonalAccessTokenOrFilterToQuery(q *user.PersonalAccessTokensOrFilter, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := patFiltersToQueries(q.Filters, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewOrQuery(mappedQueries...)
}
func authnPersonalAccessTokenAndFilterToQuery(q *user.PersonalAccessTokensAndFilter, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := patFiltersToQueries(q.Filters, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewAndQuery(mappedQueries...)
}
func authnPersonalAccessTokenNotFilterToQuery(q *user.PersonalAccessTokensNotFilter, level uint8) (query.SearchQuery, error) {
	mappedQuery, err := patFilterToQuery(q.Filter, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewNotQuery(mappedQuery)
}

// authnPersonalAccessTokenFieldNameToSortingColumn defaults to the creation date because this ensures deterministic pagination
func authnPersonalAccessTokenFieldNameToSortingColumn(field *user.PersonalAccessTokenFieldName) query.Column {
	if field == nil {
		return query.PersonalAccessTokenColumnCreationDate
	}
	switch *field {
	case user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_UNSPECIFIED:
		return query.PersonalAccessTokenColumnCreationDate
	case user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_ID:
		return query.PersonalAccessTokenColumnID
	case user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_USER_ID:
		return query.PersonalAccessTokenColumnUserID
	case user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_ORGANIZATION_ID:
		return query.PersonalAccessTokenColumnResourceOwner
	case user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_CREATED_DATE:
		return query.PersonalAccessTokenColumnCreationDate
	case user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_CHANGED_DATE:
		return query.PersonalAccessTokenColumnChangeDate
	case user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_EXPIRATION_DATE:
		return query.PersonalAccessTokenColumnExpiration
	default:
		return query.PersonalAccessTokenColumnCreationDate
	}
}

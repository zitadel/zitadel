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

func (s *Server) ListPersonalAccessTokens(ctx context.Context, req *connect.Request[user.ListPersonalAccessTokensRequest]) (*connect.Response[user.ListPersonalAccessTokensResponse], error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Msg.GetPagination())
	if err != nil {
		return nil, err
	}
	filters, err := patFiltersToQueries(req.Msg.GetFilters())
	if err != nil {
		return nil, err
	}
	search := &query.PersonalAccessTokenSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: authnPersonalAccessTokenFieldNameToSortingColumn(req.Msg.SortingColumn),
		},
		Queries: filters,
	}
	result, err := s.query.SearchPersonalAccessTokens(ctx, search, s.checkPermission)
	if err != nil {
		return nil, err
	}
	resp := &user.ListPersonalAccessTokensResponse{
		Result:     make([]*user.PersonalAccessToken, len(result.PersonalAccessTokens)),
		Pagination: filter.QueryToPaginationPb(search.SearchRequest, result.SearchResponse),
	}
	for i, pat := range result.PersonalAccessTokens {
		resp.Result[i] = &user.PersonalAccessToken{
			CreationDate:   timestamppb.New(pat.CreationDate),
			ChangeDate:     timestamppb.New(pat.ChangeDate),
			Id:             pat.ID,
			UserId:         pat.UserID,
			OrganizationId: pat.ResourceOwner,
			ExpirationDate: timestamppb.New(pat.Expiration),
		}
	}
	return connect.NewResponse(resp), nil
}

func patFiltersToQueries(filters []*user.PersonalAccessTokensSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(filters))
	for i, filter := range filters {
		q[i], err = patFilterToQuery(filter)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func patFilterToQuery(filter *user.PersonalAccessTokensSearchFilter) (query.SearchQuery, error) {
	switch q := filter.Filter.(type) {
	case *user.PersonalAccessTokensSearchFilter_CreatedDateFilter:
		return authnPersonalAccessTokenCreatedFilterToQuery(q.CreatedDateFilter)
	case *user.PersonalAccessTokensSearchFilter_ExpirationDateFilter:
		return authnPersonalAccessTokenExpirationFilterToQuery(q.ExpirationDateFilter)
	case *user.PersonalAccessTokensSearchFilter_TokenIdFilter:
		return authnPersonalAccessTokenIdFilterToQuery(q.TokenIdFilter)
	case *user.PersonalAccessTokensSearchFilter_UserIdFilter:
		return authnPersonalAccessTokenUserIdFilterToQuery(q.UserIdFilter)
	case *user.PersonalAccessTokensSearchFilter_OrganizationIdFilter:
		return authnPersonalAccessTokenOrgIdFilterToQuery(q.OrganizationIdFilter)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func authnPersonalAccessTokenIdFilterToQuery(f *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenIDQuery(f.Id)
}

func authnPersonalAccessTokenUserIdFilterToQuery(f *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenUserIDSearchQuery(f.Id)
}

func authnPersonalAccessTokenOrgIdFilterToQuery(f *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenResourceOwnerSearchQuery(f.Id)
}

func authnPersonalAccessTokenCreatedFilterToQuery(f *filter_pb.TimestampFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenCreationDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
}

func authnPersonalAccessTokenExpirationFilterToQuery(f *filter_pb.TimestampFilter) (query.SearchQuery, error) {
	return query.NewPersonalAccessTokenExpirationDateDateQuery(f.Timestamp.AsTime(), filter.TimestampMethodPbToQuery(f.Method))
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
	case user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_EXPIRATION_DATE:
		return query.PersonalAccessTokenColumnExpiration
	default:
		return query.PersonalAccessTokenColumnCreationDate
	}
}

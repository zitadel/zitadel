package convert

import (
	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	metadata "github.com/zitadel/zitadel/pkg/grpc/metadata/v2"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func ListUsersByMetadataRequestToModel(req *user.ListUsersByMetadataRequest, sysDefaults systemdefaults.SystemDefaults) (*query.UserSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(sysDefaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := usersByMetadataQueries(req.GetFilters(), 0)
	if err != nil {
		return nil, err
	}

	return &query.UserSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: usersByMetadataSorting(req.GetSortingColumn()),
		},
		Queries: queries,
	}, nil

}

func usersByMetadataSorting(sortingColumn user.UsersByMetadataSorting) query.Column {
	switch sortingColumn {
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_DISPLAY_NAME:
		return query.HumanDisplayNameCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_EMAIL:
		return query.HumanEmailCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_FIRST_NAME:
		return query.HumanFirstNameCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_LAST_NAME:
		return query.HumanLastNameCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_METADATA_KEY:
		return query.UserMetadataKeyCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_METADATA_VALUE:
		return query.UserMetadataValueCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_NICK_NAME:
		return query.HumanNickNameCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_STATE:
		return query.UserStateCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_TYPE:
		return query.UserTypeCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_USER_NAME:
		return query.UserUsernameCol
	case user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_USER_ID, user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_UNSPECIFIED:
		fallthrough
	default:
		return query.UserIDCol
	}
}

func usersByMetadataQueries(queries []*metadata.UserByMetadataSearchFilter, nesting uint) ([]query.SearchQuery, error) {
	toReturn := make([]query.SearchQuery, len(queries))

	for i, query := range queries {
		res, err := userByMetadataQuery(query, nesting)
		if err != nil {
			return nil, err
		}
		toReturn[i] = res
	}

	return toReturn, nil
}

func userByMetadataQuery(q *metadata.UserByMetadataSearchFilter, nesting uint) (query.SearchQuery, error) {
	if nesting > 20 {
		return nil, zerrors.ThrowInvalidArgument(nil, "CONV-Jhaltm", "Errors.Query.TooManyNestingLevels")
	}

	switch t := q.GetFilter().(type) {

	case *metadata.UserByMetadataSearchFilter_KeyFilter:
		return query.NewUserMetadataKeySearchQuery(t.KeyFilter.GetKey(), filter.TextMethodPbToQuery(t.KeyFilter.GetMethod()))

	case *metadata.UserByMetadataSearchFilter_ValueFilter:
		return query.NewUserMetadataValueSearchQuery(t.ValueFilter.GetValue(), filter.ByteMethodPbToQuery(t.ValueFilter.GetMethod()))

	case *metadata.UserByMetadataSearchFilter_AndFilter:
		mappedQueries, err := usersByMetadataQueries(t.AndFilter.GetQueries(), nesting+1)
		if err != nil {
			return nil, err
		}
		return query.NewUserMetadataAndSearchQuery(mappedQueries)

	case *metadata.UserByMetadataSearchFilter_OrFilter:
		mappedQueries, err := usersByMetadataQueries(t.OrFilter.GetQueries(), nesting+1)
		if err != nil {
			return nil, err
		}
		return query.NewUserMetadataOrSearchQuery(mappedQueries)

	case *metadata.UserByMetadataSearchFilter_NotFilter:
		mappedQuery, err := userByMetadataQuery(t.NotFilter.GetQuery(), nesting+1)
		if err != nil {
			return nil, err
		}
		return query.NewUserMetadataNotSearchQuery(mappedQuery)

	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "CONV-GG1Jnh", "List.Query.Invalid")
	}
}

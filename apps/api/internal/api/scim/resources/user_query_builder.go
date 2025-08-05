package resources

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// fieldPathColumnMapping maps lowercase json field names of the scim user to the matching column in the projection
// only a limited set of fields is supported
// to ensure database performance.
var fieldPathColumnMapping = filter.FieldPathMapping{
	"meta.created": {
		Column:    query.UserCreationDateCol,
		FieldType: filter.FieldTypeTimestamp,
	},
	"meta.lastmodified": {
		Column:    query.UserChangeDateCol,
		FieldType: filter.FieldTypeTimestamp,
	},
	"id": {
		Column:    query.UserIDCol,
		FieldType: filter.FieldTypeString,
	},
	"username": {
		Column:          query.UserUsernameCol,
		FieldType:       filter.FieldTypeString,
		CaseInsensitive: true,
	},
	"name.familyname": {
		Column:    query.HumanLastNameCol,
		FieldType: filter.FieldTypeString,
	},
	"name.givenname": {
		Column:    query.HumanFirstNameCol,
		FieldType: filter.FieldTypeString,
	},
	"emails": {
		Column:    query.HumanEmailCol,
		FieldType: filter.FieldTypeString,
	},
	"emails.value": {
		Column:    query.HumanEmailCol,
		FieldType: filter.FieldTypeString,
	},
	"active": {
		FieldType:        filter.FieldTypeCustom,
		BuildMappedQuery: buildActiveUserStateQuery,
	},
	"externalid": {
		FieldType:        filter.FieldTypeCustom,
		BuildMappedQuery: newMetadataQueryBuilder(metadata.KeyExternalId),
	},
}

func (h *UsersHandler) buildListQuery(ctx context.Context, request *ListRequest) (*query.UserSearchQueries, error) {
	searchRequest, err := request.toSearchRequest(query.UserIDCol, fieldPathColumnMapping)
	if err != nil {
		return nil, err
	}

	q := &query.UserSearchQueries{
		SearchRequest: searchRequest,
	}

	// the zitadel scim implementation only supports humans for now
	userTypeQuery, err := query.NewUserTypeSearchQuery(domain.UserTypeHuman)
	if err != nil {
		return nil, err
	}

	// the scim service is always limited to one organization
	// the organization is the resource owner
	orgIDQuery, err := query.NewUserResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID, query.TextEquals)
	if err != nil {
		return nil, err
	}

	q.Queries = append(q.Queries, orgIDQuery, userTypeQuery)

	if request.Filter == nil {
		return q, nil
	}

	filterQuery, err := request.Filter.BuildQuery(ctx, h.schema.ID, fieldPathColumnMapping)
	if err != nil {
		return nil, err
	}

	q.Queries = append(q.Queries, filterQuery)
	return q, nil
}

func newMetadataQueryBuilder(key metadata.Key) filter.MappedQueryBuilderFunc {
	return func(ctx context.Context, compareValue *filter.CompValue, op *filter.CompareOp) (query.SearchQuery, error) {
		return buildMetadataQuery(ctx, key, compareValue, op)
	}
}

func buildMetadataQuery(ctx context.Context, key metadata.Key, value *filter.CompValue, op *filter.CompareOp) (query.SearchQuery, error) {
	if value.StringValue == nil {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-EXid1", "invalid filter expression: unsupported comparison value"))
	}

	var comparisonOperator query.BytesComparison

	switch {
	case op.Equal:
		comparisonOperator = query.BytesEquals
	case op.NotEqual:
		comparisonOperator = query.BytesNotEquals
	default:
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-EXid1", "invalid filter expression: unsupported comparison operator"))
	}

	scopedKey := string(metadata.ScopeKey(ctx, key))
	return query.NewUserMetadataExistsQuery(scopedKey, []byte(*value.StringValue), query.TextEquals, comparisonOperator)
}

func buildActiveUserStateQuery(_ context.Context, compareValue *filter.CompValue, op *filter.CompareOp) (query.SearchQuery, error) {
	if !op.Equal && !op.NotEqual {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-MGdg", "invalid filter expression: active unsupported comparison operator"))
	}

	if !compareValue.BooleanTrue && !compareValue.BooleanFalse {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-MGdr", "invalid filter expression: active unsupported comparison value"))
	}

	active := compareValue.BooleanTrue && op.Equal || compareValue.BooleanFalse && op.NotEqual
	if active {
		activeQuery, err := query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateActive), query.NumberEquals)
		if err != nil {
			return nil, err
		}

		initialQuery, err := query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateInitial), query.NumberEquals)
		if err != nil {
			return nil, err
		}

		return query.NewOrQuery(initialQuery, activeQuery)
	}

	activeQuery, err := query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateActive), query.NumberNotEquals)
	if err != nil {
		return nil, err
	}

	initialQuery, err := query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateInitial), query.NumberNotEquals)
	if err != nil {
		return nil, err
	}

	return query.NewAndQuery(initialQuery, activeQuery)
}

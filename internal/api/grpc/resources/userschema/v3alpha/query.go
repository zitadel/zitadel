package userschema

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	schema "github.com/zitadel/zitadel/pkg/grpc/resources/userschema/v3alpha"
)

func (s *Server) SearchUserSchemas(ctx context.Context, req *schema.SearchUserSchemasRequest) (*schema.SearchUserSchemasResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	queries, err := s.searchUserSchemaToModel(req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchUserSchema(ctx, queries)
	if err != nil {
		return nil, err
	}
	userSchemas, err := userSchemasToPb(res.UserSchemas)
	if err != nil {
		return nil, err
	}
	return &schema.SearchUserSchemasResponse{
		Details: resource_object.ToSearchDetailsPb(queries.SearchRequest, res.SearchResponse),
		Result:  userSchemas,
	}, nil
}

func (s *Server) GetUserSchema(ctx context.Context, req *schema.GetUserSchemaRequest) (*schema.GetUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	res, err := s.query.GetUserSchemaByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	userSchema, err := userSchemaToPb(res)
	if err != nil {
		return nil, err
	}
	return &schema.GetUserSchemaResponse{
		UserSchema: userSchema,
	}, nil
}

func (s *Server) searchUserSchemaToModel(req *schema.SearchUserSchemasRequest) (*query.UserSchemaSearchQueries, error) {
	offset, limit, asc, err := resource_object.SearchQueryPbToQuery(s.systemDefaults, req.Query)
	if err != nil {
		return nil, err
	}
	queries, err := userSchemaFiltersToQuery(req.Filters, 0) // start at level 0
	if err != nil {
		return nil, err
	}
	return &query.UserSchemaSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: userSchemaFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func userSchemaFieldNameToSortingColumn(field *schema.FieldName) query.Column {
	if field == nil {
		return query.UserSchemaCreationDateCol
	}
	switch *field {
	case schema.FieldName_FIELD_NAME_TYPE:
		return query.UserSchemaTypeCol
	case schema.FieldName_FIELD_NAME_STATE:
		return query.UserSchemaStateCol
	case schema.FieldName_FIELD_NAME_REVISION:
		return query.UserSchemaRevisionCol
	case schema.FieldName_FIELD_NAME_CHANGE_DATE:
		return query.UserSchemaChangeDateCol
	case schema.FieldName_FIELD_NAME_CREATION_DATE:
		return query.UserSchemaCreationDateCol
	case schema.FieldName_FIELD_NAME_UNSPECIFIED:
		return query.UserSchemaIDCol
	default:
		return query.UserSchemaIDCol
	}
}

func userSchemasToPb(schemas []*query.UserSchema) (_ []*schema.GetUserSchema, err error) {
	userSchemas := make([]*schema.GetUserSchema, len(schemas))
	for i, userSchema := range schemas {
		userSchemas[i], err = userSchemaToPb(userSchema)
		if err != nil {
			return nil, err
		}
	}
	return userSchemas, nil
}

func userSchemaToPb(userSchema *query.UserSchema) (*schema.GetUserSchema, error) {
	s := new(structpb.Struct)
	if err := s.UnmarshalJSON(userSchema.Schema); err != nil {
		return nil, err
	}
	return &schema.GetUserSchema{
		Details: resource_object.DomainToDetailsPb(&userSchema.ObjectDetails, object.OwnerType_OWNER_TYPE_INSTANCE, userSchema.ResourceOwner),
		Config: &schema.UserSchema{
			Type: userSchema.Type,
			DataType: &schema.UserSchema_Schema{
				Schema: s,
			},
			PossibleAuthenticators: authenticatorTypesToPb(userSchema.PossibleAuthenticators),
		},
		State:    userSchemaStateToPb(userSchema.State),
		Revision: userSchema.Revision,
	}, nil
}

func authenticatorTypesToPb(authenticators []domain.AuthenticatorType) []schema.AuthenticatorType {
	authTypes := make([]schema.AuthenticatorType, len(authenticators))
	for i, authenticator := range authenticators {
		authTypes[i] = authenticatorTypeToPb(authenticator)
	}
	return authTypes
}

func authenticatorTypeToPb(authenticator domain.AuthenticatorType) schema.AuthenticatorType {
	switch authenticator {
	case domain.AuthenticatorTypeUsername:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_USERNAME
	case domain.AuthenticatorTypePassword:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_PASSWORD
	case domain.AuthenticatorTypeWebAuthN:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_WEBAUTHN
	case domain.AuthenticatorTypeTOTP:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_TOTP
	case domain.AuthenticatorTypeOTPEmail:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_OTP_EMAIL
	case domain.AuthenticatorTypeOTPSMS:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_OTP_SMS
	case domain.AuthenticatorTypeAuthenticationKey:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_AUTHENTICATION_KEY
	case domain.AuthenticatorTypeIdentityProvider:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_IDENTITY_PROVIDER
	case domain.AuthenticatorTypeUnspecified:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_UNSPECIFIED
	default:
		return schema.AuthenticatorType_AUTHENTICATOR_TYPE_UNSPECIFIED
	}
}

func userSchemaStateToPb(state domain.UserSchemaState) schema.State {
	switch state {
	case domain.UserSchemaStateActive:
		return schema.State_STATE_ACTIVE
	case domain.UserSchemaStateInactive:
		return schema.State_STATE_INACTIVE
	case domain.UserSchemaStateUnspecified,
		domain.UserSchemaStateDeleted:
		return schema.State_STATE_UNSPECIFIED
	default:
		return schema.State_STATE_UNSPECIFIED
	}
}

func userSchemaStateToDomain(state schema.State) domain.UserSchemaState {
	switch state {
	case schema.State_STATE_ACTIVE:
		return domain.UserSchemaStateActive
	case schema.State_STATE_INACTIVE:
		return domain.UserSchemaStateInactive
	case schema.State_STATE_UNSPECIFIED:
		return domain.UserSchemaStateUnspecified
	default:
		return domain.UserSchemaStateUnspecified
	}
}

func userSchemaFiltersToQuery(queries []*schema.SearchFilter, level uint8) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = userSchemaFilterToQuery(query, level)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func userSchemaFilterToQuery(query *schema.SearchFilter, level uint8) (query.SearchQuery, error) {
	if level > 20 {
		// can't go deeper than 20 levels of nesting.
		return nil, zerrors.ThrowInvalidArgument(nil, "SCHEMA-zsQ97", "Errors.Query.TooManyNestingLevels")
	}
	switch q := query.Filter.(type) {
	case *schema.SearchFilter_StateFilter:
		return stateQueryToQuery(q.StateFilter)
	case *schema.SearchFilter_TypeFilter:
		return typeQueryToQuery(q.TypeFilter)
	case *schema.SearchFilter_IdFilter:
		return idQueryToQuery(q.IdFilter)
	case *schema.SearchFilter_OrFilter:
		return orQueryToQuery(q.OrFilter, level)
	case *schema.SearchFilter_AndFilter:
		return andQueryToQuery(q.AndFilter, level)
	case *schema.SearchFilter_NotFilter:
		return notQueryToQuery(q.NotFilter, level)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "SCHEMA-vR9nC", "List.Query.Invalid")
	}
}

func stateQueryToQuery(q *schema.StateFilter) (query.SearchQuery, error) {
	return query.NewUserSchemaStateSearchQuery(userSchemaStateToDomain(q.GetState()))
}

func typeQueryToQuery(q *schema.TypeFilter) (query.SearchQuery, error) {
	return query.NewUserSchemaTypeSearchQuery(q.GetType(), resource_object.TextMethodPbToQuery(q.GetMethod()))
}

func idQueryToQuery(q *schema.IDFilter) (query.SearchQuery, error) {
	return query.NewUserSchemaIDSearchQuery(q.GetId(), resource_object.TextMethodPbToQuery(q.GetMethod()))
}

func orQueryToQuery(q *schema.OrFilter, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := userSchemaFiltersToQuery(q.GetQueries(), level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserOrSearchQuery(mappedQueries)
}

func andQueryToQuery(q *schema.AndFilter, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := userSchemaFiltersToQuery(q.GetQueries(), level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserAndSearchQuery(mappedQueries)
}

func notQueryToQuery(q *schema.NotFilter, level uint8) (query.SearchQuery, error) {
	mappedQuery, err := userSchemaFilterToQuery(q.GetFilter(), level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserNotSearchQuery(mappedQuery)
}

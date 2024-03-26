package schema

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	schema "github.com/zitadel/zitadel/pkg/grpc/user/schema/v3alpha"
)

func (s *Server) CreateUserSchema(ctx context.Context, req *schema.CreateUserSchemaRequest) (*schema.CreateUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	userSchema, err := createUserSchemaToCommand(req, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	id, details, err := s.command.CreateUserSchema(ctx, userSchema)
	if err != nil {
		return nil, err
	}
	return &schema.CreateUserSchemaResponse{
		Id:      id,
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) UpdateUserSchema(ctx context.Context, req *schema.UpdateUserSchemaRequest) (*schema.UpdateUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	userSchema, err := updateUserSchemaToCommand(req, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	details, err := s.command.UpdateUserSchema(ctx, userSchema)
	if err != nil {
		return nil, err
	}
	return &schema.UpdateUserSchemaResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeactivateUserSchema(ctx context.Context, req *schema.DeactivateUserSchemaRequest) (*schema.DeactivateUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeactivateUserSchema(ctx, req.GetId(), authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &schema.DeactivateUserSchemaResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ReactivateUserSchema(ctx context.Context, req *schema.ReactivateUserSchemaRequest) (*schema.ReactivateUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.ReactivateUserSchema(ctx, req.GetId(), authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &schema.ReactivateUserSchemaResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeleteUserSchema(ctx context.Context, req *schema.DeleteUserSchemaRequest) (*schema.DeleteUserSchemaResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteUserSchema(ctx, req.GetId(), authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &schema.DeleteUserSchemaResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ListUserSchemas(ctx context.Context, req *schema.ListUserSchemasRequest) (*schema.ListUserSchemasResponse, error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	queries, err := listUserSchemaToQuery(req)
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
	return &schema.ListUserSchemasResponse{
		Details: object.ToListDetails(res.SearchResponse),
		Result:  userSchemas,
	}, nil
}

func (s *Server) GetUserSchemaByID(ctx context.Context, req *schema.GetUserSchemaByIDRequest) (*schema.GetUserSchemaByIDResponse, error) {
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
	return &schema.GetUserSchemaByIDResponse{
		Schema: userSchema,
	}, nil
}

func userSchemasToPb(schemas []*query.UserSchema) (_ []*schema.UserSchema, err error) {
	userSchemas := make([]*schema.UserSchema, len(schemas))
	for i, userSchema := range schemas {
		userSchemas[i], err = userSchemaToPb(userSchema)
		if err != nil {
			return nil, err
		}
	}
	return userSchemas, nil
}

func userSchemaToPb(userSchema *query.UserSchema) (*schema.UserSchema, error) {
	s := new(structpb.Struct)
	if err := s.UnmarshalJSON(userSchema.Schema); err != nil {
		return nil, err
	}
	return &schema.UserSchema{
		Id:                     userSchema.ID,
		Details:                object.DomainToDetailsPb(&userSchema.ObjectDetails),
		Type:                   userSchema.Type,
		State:                  userSchemaStateToPb(userSchema.State),
		Revision:               userSchema.Revision,
		Schema:                 s,
		PossibleAuthenticators: authenticatorTypesToPb(userSchema.PossibleAuthenticators),
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

func listUserSchemaToQuery(req *schema.ListUserSchemasRequest) (*query.UserSchemaSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	queries, err := userSchemaQueriesToQuery(req.Queries, 0) // start at level 0
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

func userSchemaFieldNameToSortingColumn(column schema.FieldName) query.Column {
	switch column {
	case schema.FieldName_FIELD_NAME_TYPE:
		return query.UserSchemaTypeCol
	case schema.FieldName_FIELD_NAME_STATE:
		return query.UserSchemaStateCol
	case schema.FieldName_FIELD_NAME_REVISION:
		return query.UserSchemaRevisionCol
	case schema.FieldName_FIELD_NAME_CHANGE_DATE:
		return query.UserSchemaChangeDateCol
	case schema.FieldName_FIELD_NAME_UNSPECIFIED:
		return query.UserSchemaIDCol
	default:
		return query.UserSchemaIDCol
	}
}

func checkUserSchemaEnabled(ctx context.Context) error {
	if authz.GetInstance(ctx).Features().UserSchema {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "SCHEMA-SFjk3", "Errors.UserSchema.NotEnabled")
}

func createUserSchemaToCommand(req *schema.CreateUserSchemaRequest, resourceOwner string) (*command.CreateUserSchema, error) {
	schema, err := req.GetSchema().MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &command.CreateUserSchema{
		ResourceOwner:          resourceOwner,
		Type:                   req.GetType(),
		Schema:                 schema,
		PossibleAuthenticators: authenticatorsToDomain(req.GetPossibleAuthenticators()),
	}, nil
}

func updateUserSchemaToCommand(req *schema.UpdateUserSchemaRequest, resourceOwner string) (*command.UpdateUserSchema, error) {
	schema, err := req.GetSchema().MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &command.UpdateUserSchema{
		ID:                     req.GetId(),
		ResourceOwner:          resourceOwner,
		Type:                   req.Type,
		Schema:                 schema,
		PossibleAuthenticators: authenticatorsToDomain(req.GetPossibleAuthenticators()),
	}, nil
}

func authenticatorsToDomain(authenticators []schema.AuthenticatorType) []domain.AuthenticatorType {
	types := make([]domain.AuthenticatorType, len(authenticators))
	for i, authenticator := range authenticators {
		types[i] = authenticatorTypeToDomain(authenticator)
	}
	return types
}

func authenticatorTypeToDomain(authenticator schema.AuthenticatorType) domain.AuthenticatorType {
	switch authenticator {
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_UNSPECIFIED:
		return domain.AuthenticatorTypeUnspecified
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_USERNAME:
		return domain.AuthenticatorTypeUsername
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_PASSWORD:
		return domain.AuthenticatorTypePassword
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_WEBAUTHN:
		return domain.AuthenticatorTypeWebAuthN
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_TOTP:
		return domain.AuthenticatorTypeTOTP
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_OTP_EMAIL:
		return domain.AuthenticatorTypeOTPEmail
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_OTP_SMS:
		return domain.AuthenticatorTypeOTPSMS
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_AUTHENTICATION_KEY:
		return domain.AuthenticatorTypeAuthenticationKey
	case schema.AuthenticatorType_AUTHENTICATOR_TYPE_IDENTITY_PROVIDER:
		return domain.AuthenticatorTypeIdentityProvider
	default:
		return domain.AuthenticatorTypeUnspecified
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

func userSchemaQueriesToQuery(queries []*schema.SearchQuery, level uint8) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = userSchemaQueryToQuery(query, level)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func userSchemaQueryToQuery(query *schema.SearchQuery, level uint8) (query.SearchQuery, error) {
	if level > 20 {
		// can't go deeper than 20 levels of nesting.
		return nil, zerrors.ThrowInvalidArgument(nil, "SCHEMA-zsQ97", "Errors.Query.TooManyNestingLevels")
	}
	switch q := query.Query.(type) {
	case *schema.SearchQuery_StateQuery:
		return stateQueryToQuery(q.StateQuery)
	case *schema.SearchQuery_TypeQuery:
		return typeQueryToQuery(q.TypeQuery)
	case *schema.SearchQuery_IdQuery:
		return idQueryToQuery(q.IdQuery)
	case *schema.SearchQuery_OrQuery:
		return orQueryToQuery(q.OrQuery, level)
	case *schema.SearchQuery_AndQuery:
		return andQueryToQuery(q.AndQuery, level)
	case *schema.SearchQuery_NotQuery:
		return notQueryToQuery(q.NotQuery, level)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "SCHEMA-vR9nC", "List.Query.Invalid")
	}
}

func stateQueryToQuery(q *schema.StateQuery) (query.SearchQuery, error) {
	return query.NewUserSchemaStateSearchQuery(userSchemaStateToDomain(q.GetState()))
}

func typeQueryToQuery(q *schema.TypeQuery) (query.SearchQuery, error) {
	return query.NewUserSchemaTypeSearchQuery(q.GetType(), object.TextMethodToQuery(q.GetMethod()))
}

func idQueryToQuery(q *schema.IDQuery) (query.SearchQuery, error) {
	return query.NewUserSchemaIDSearchQuery(q.GetId(), object.TextMethodToQuery(q.GetMethod()))
}

func orQueryToQuery(q *schema.OrQuery, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := userSchemaQueriesToQuery(q.GetQueries(), level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserOrSearchQuery(mappedQueries)
}

func andQueryToQuery(q *schema.AndQuery, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := userSchemaQueriesToQuery(q.GetQueries(), level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserAndSearchQuery(mappedQueries)
}

func notQueryToQuery(q *schema.NotQuery, level uint8) (query.SearchQuery, error) {
	mappedQuery, err := userSchemaQueryToQuery(q.GetQuery(), level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserNotSearchQuery(mappedQuery)
}

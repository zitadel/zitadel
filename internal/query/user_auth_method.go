package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	userAuthMethodTable = table{
		name:          projection.UserAuthMethodTable,
		instanceIDCol: projection.UserAuthMethodInstanceIDCol,
	}
	UserAuthMethodColumnTokenID = Column{
		name:  projection.UserAuthMethodTokenIDCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnCreationDate = Column{
		name:  projection.UserAuthMethodCreationDateCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnChangeDate = Column{
		name:  projection.UserAuthMethodChangeDateCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnResourceOwner = Column{
		name:  projection.UserAuthMethodResourceOwnerCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnInstanceID = Column{
		name:  projection.UserAuthMethodInstanceIDCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnUserID = Column{
		name:  projection.UserAuthMethodUserIDCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnSequence = Column{
		name:  projection.UserAuthMethodSequenceCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnName = Column{
		name:  projection.UserAuthMethodNameCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnState = Column{
		name:  projection.UserAuthMethodStateCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnMethodType = Column{
		name:  projection.UserAuthMethodTypeCol,
		table: userAuthMethodTable,
	}
	UserAuthMethodColumnDomain = Column{
		name:  projection.UserAuthMethodDomainCol,
		table: userAuthMethodTable,
	}

	authMethodTypeTable      = userAuthMethodTable.setAlias("auth_method_types")
	authMethodTypeUserID     = UserAuthMethodColumnUserID.setTable(authMethodTypeTable)
	authMethodTypeInstanceID = UserAuthMethodColumnInstanceID.setTable(authMethodTypeTable)
	authMethodTypeType       = UserAuthMethodColumnMethodType.setTable(authMethodTypeTable)
	authMethodTypeState      = UserAuthMethodColumnState.setTable(authMethodTypeTable)
	authMethodTypeDomain     = UserAuthMethodColumnDomain.setTable(authMethodTypeTable)

	userIDPsCountTable      = idpUserLinkTable.setAlias("user_idps_count")
	userIDPsCountUserID     = IDPUserLinkUserIDCol.setTable(userIDPsCountTable)
	userIDPsCountInstanceID = IDPUserLinkInstanceIDCol.setTable(userIDPsCountTable)
	userIDPsCountCount      = Column{
		name:  "count",
		table: userIDPsCountTable,
	}

	forceMFATable          = loginPolicyTable.setAlias("auth_methods_force_mfa")
	forceMFAInstanceID     = LoginPolicyColumnInstanceID.setTable(forceMFATable)
	forceMFAOrgID          = LoginPolicyColumnOrgID.setTable(forceMFATable)
	forceMFAIsDefault      = LoginPolicyColumnIsDefault.setTable(forceMFATable)
	forceMFAForce          = LoginPolicyColumnForceMFA.setTable(forceMFATable)
	forceMFAForceLocalOnly = LoginPolicyColumnForceMFALocalOnly.setTable(forceMFATable)
)

type AuthMethods struct {
	SearchResponse
	AuthMethods []*AuthMethod
}

func authMethodsCheckPermission(ctx context.Context, methods *AuthMethods, permissionCheck domain.PermissionCheck) {
	methods.AuthMethods = slices.DeleteFunc(methods.AuthMethods,
		func(method *AuthMethod) bool {
			return userCheckPermission(ctx, method.ResourceOwner, method.UserID, permissionCheck) != nil
		},
	)
}

func userAuthMethodPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		UserAuthMethodColumnResourceOwner,
		domain.PermissionUserRead,
		OwnedRowsPermissionOption(UserAuthMethodColumnUserID),
	)
	return query.JoinClause(join, args...)
}

type AuthMethod struct {
	UserID        string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.MFAState
	Sequence      uint64

	TokenID string
	Name    string
	Type    domain.UserAuthMethodType
}

type AuthMethodTypes struct {
	SearchResponse
	AuthMethodTypes []domain.UserAuthMethodType
}

type UserAuthMethodSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *UserAuthMethodSearchQueries) hasUserID() bool {
	for _, query := range q.Queries {
		if query.Col() == UserAuthMethodColumnUserID {
			return true
		}
	}
	return false
}

func (q *Queries) SearchUserAuthMethods(ctx context.Context, queries *UserAuthMethodSearchQueries, permissionCheck domain.PermissionCheck) (userAuthMethods *AuthMethods, err error) {
	permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	methods, err := q.searchUserAuthMethods(ctx, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && len(methods.AuthMethods) > 0 && !permissionCheckV2 {
		// when userID for query is provided, only one check has to be done
		if queries.hasUserID() {
			if err := userCheckPermission(ctx, methods.AuthMethods[0].ResourceOwner, methods.AuthMethods[0].UserID, permissionCheck); err != nil {
				return nil, err
			}
		} else {
			authMethodsCheckPermission(ctx, methods, permissionCheck)
		}
	}
	return methods, nil
}

func (q *Queries) searchUserAuthMethods(ctx context.Context, queries *UserAuthMethodSearchQueries, permissionCheckV2 bool) (userAuthMethods *AuthMethods, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareUserAuthMethodsQuery()
	query = userAuthMethodPermissionCheckV2(ctx, query, permissionCheckV2)
	stmt, args, err := queries.toQuery(query).Where(sq.Eq{UserAuthMethodColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-j9NJd", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		userAuthMethods, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	userAuthMethods.State, err = q.latestState(ctx, userAuthMethodTable)
	return userAuthMethods, err
}

func (q *Queries) ListUserAuthMethodTypes(ctx context.Context, userID string, activeOnly bool, includeWithoutDomain bool, queryDomain string) (userAuthMethodTypes *AuthMethodTypes, err error) {
	ctxData := authz.GetCtxData(ctx)
	if ctxData.UserID != userID {
		if err := q.checkPermission(ctx, domain.PermissionUserRead, ctxData.OrgID, userID); err != nil {
			return nil, err
		}
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareUserAuthMethodTypesQuery(activeOnly, includeWithoutDomain, queryDomain)
	eq := sq.Eq{
		UserIDCol.identifier():         userID,
		UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-Sfdrg", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		userAuthMethodTypes, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	userAuthMethodTypes.State, err = q.latestState(ctx, userTable, notifyTable, userAuthMethodTable, idpUserLinkTable)
	return userAuthMethodTypes, err
}

type UserAuthMethodRequirements struct {
	UserType          domain.UserType
	ForceMFA          bool
	ForceMFALocalOnly bool
}

//go:embed user_auth_method_types_required.sql
var listUserAuthMethodTypesStmt string

func (q *Queries) ListUserAuthMethodTypesRequired(ctx context.Context, userID string) (requirements *UserAuthMethodRequirements, err error) {
	ctxData := authz.GetCtxData(ctx)
	if ctxData.UserID != userID {
		if err := q.checkPermission(ctx, domain.PermissionUserRead, ctxData.OrgID, userID); err != nil {
			return nil, err
		}
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = q.client.QueryRowContext(ctx,
		func(row *sql.Row) error {
			var userType sql.NullInt32
			var forceMFA sql.NullBool
			var forceMFALocalOnly sql.NullBool
			err := row.Scan(
				&userType,
				&forceMFA,
				&forceMFALocalOnly,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return zerrors.ThrowNotFound(err, "QUERY-SF3h2", "Errors.Internal")
				}
				return zerrors.ThrowInternal(err, "QUERY-Sf3rt", "Errors.Internal")
			}
			requirements = &UserAuthMethodRequirements{
				UserType:          domain.UserType(userType.Int32),
				ForceMFA:          forceMFA.Bool,
				ForceMFALocalOnly: forceMFALocalOnly.Bool,
			}
			return nil
		},
		listUserAuthMethodTypesStmt,
		userID,
		authz.GetInstance(ctx).InstanceID(),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dun75", "Errors.Internal")
	}
	return requirements, nil
}

func NewUserAuthMethodUserIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(UserAuthMethodColumnUserID, value, TextEquals)
}

func NewUserAuthMethodTokenIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(UserAuthMethodColumnTokenID, value, TextEquals)
}

func NewUserAuthMethodResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(UserAuthMethodColumnResourceOwner, value, TextEquals)
}

func NewUserAuthMethodTypeSearchQuery(value domain.UserAuthMethodType) (SearchQuery, error) {
	return NewNumberQuery(UserAuthMethodColumnMethodType, value, NumberEquals)
}

func NewUserAuthMethodStateSearchQuery(value domain.MFAState) (SearchQuery, error) {
	return NewNumberQuery(UserAuthMethodColumnState, value, NumberEquals)
}

func NewUserAuthMethodTypesSearchQuery(values ...domain.UserAuthMethodType) (SearchQuery, error) {
	list := make([]interface{}, len(values))
	for i, value := range values {
		list[i] = value
	}
	return NewListQuery(UserAuthMethodColumnMethodType, list, ListIn)
}

func NewUserAuthMethodStatesSearchQuery(values ...domain.MFAState) (SearchQuery, error) {
	list := make([]interface{}, len(values))
	for i, value := range values {
		list[i] = value
	}
	return NewListQuery(UserAuthMethodColumnState, list, ListIn)
}

func (r *UserAuthMethodSearchQueries) AppendResourceOwnerQuery(orgID string) error {
	query, err := NewUserAuthMethodResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *UserAuthMethodSearchQueries) AppendUserIDQuery(userID string) error {
	query, err := NewUserAuthMethodUserIDSearchQuery(userID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *UserAuthMethodSearchQueries) AppendTokenIDQuery(tokenID string) error {
	query, err := NewUserAuthMethodTokenIDSearchQuery(tokenID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *UserAuthMethodSearchQueries) AppendStateQuery(state domain.MFAState) error {
	query, err := NewUserAuthMethodStateSearchQuery(state)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *UserAuthMethodSearchQueries) AppendStatesQuery(state ...domain.MFAState) error {
	query, err := NewUserAuthMethodStatesSearchQuery(state...)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *UserAuthMethodSearchQueries) AppendAuthMethodQuery(authMethod domain.UserAuthMethodType) error {
	query, err := NewUserAuthMethodTypeSearchQuery(authMethod)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *UserAuthMethodSearchQueries) AppendAuthMethodsQuery(authMethod ...domain.UserAuthMethodType) error {
	query, err := NewUserAuthMethodTypesSearchQuery(authMethod...)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (q *UserAuthMethodSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareUserAuthMethodsQuery() (sq.SelectBuilder, func(*sql.Rows) (*AuthMethods, error)) {
	return sq.Select(
			UserAuthMethodColumnTokenID.identifier(),
			UserAuthMethodColumnCreationDate.identifier(),
			UserAuthMethodColumnChangeDate.identifier(),
			UserAuthMethodColumnResourceOwner.identifier(),
			UserAuthMethodColumnUserID.identifier(),
			UserAuthMethodColumnSequence.identifier(),
			UserAuthMethodColumnName.identifier(),
			UserAuthMethodColumnState.identifier(),
			UserAuthMethodColumnMethodType.identifier(),
			countColumn.identifier()).
			From(userAuthMethodTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*AuthMethods, error) {
			userAuthMethods := make([]*AuthMethod, 0)
			var count uint64
			for rows.Next() {
				authMethod := new(AuthMethod)
				err := rows.Scan(
					&authMethod.TokenID,
					&authMethod.CreationDate,
					&authMethod.ChangeDate,
					&authMethod.ResourceOwner,
					&authMethod.UserID,
					&authMethod.Sequence,
					&authMethod.Name,
					&authMethod.State,
					&authMethod.Type,
					&count,
				)
				if err != nil {
					return nil, err
				}
				userAuthMethods = append(userAuthMethods, authMethod)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-3n9fl", "Errors.Query.CloseRows")
			}

			return &AuthMethods{
				AuthMethods: userAuthMethods,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareUserAuthMethodTypesQuery(activeOnly bool, includeWithoutDomain bool, queryDomain string) (sq.SelectBuilder, func(*sql.Rows) (*AuthMethodTypes, error)) {
	authMethodsQuery, authMethodsArgs, err := prepareAuthMethodQuery(activeOnly, includeWithoutDomain, queryDomain)
	if err != nil {
		return sq.SelectBuilder{}, nil
	}
	idpsQuery, err := prepareAuthMethodsIDPsQuery()
	if err != nil {
		return sq.SelectBuilder{}, nil
	}
	return sq.Select(
			NotifyPasswordSetCol.identifier(),
			authMethodTypeType.identifier(),
			userIDPsCountCount.identifier()).
			From(userTable.identifier()).
			LeftJoin(join(NotifyUserIDCol, UserIDCol)).
			LeftJoin("("+authMethodsQuery+") AS "+authMethodTypeTable.alias+" ON "+
				authMethodTypeUserID.identifier()+" = "+UserIDCol.identifier()+" AND "+
				authMethodTypeInstanceID.identifier()+" = "+UserInstanceIDCol.identifier(),
				authMethodsArgs...).
			LeftJoin("(" + idpsQuery + ") AS " + userIDPsCountTable.alias + " ON " +
				userIDPsCountUserID.identifier() + " = " + UserIDCol.identifier() + " AND " +
				userIDPsCountInstanceID.identifier() + " = " + UserInstanceIDCol.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*AuthMethodTypes, error) {
			userAuthMethodTypes := make([]domain.UserAuthMethodType, 0)
			var passwordSet sql.NullBool
			var idp sql.NullInt64
			for rows.Next() {
				var authMethodType sql.NullInt16
				err := rows.Scan(
					&passwordSet,
					&authMethodType,
					&idp,
				)
				if err != nil {
					return nil, err
				}
				if authMethodType.Valid {
					userAuthMethodTypes = append(userAuthMethodTypes, domain.UserAuthMethodType(authMethodType.Int16))
				}
			}
			if passwordSet.Valid && passwordSet.Bool {
				userAuthMethodTypes = append(userAuthMethodTypes, domain.UserAuthMethodTypePassword)
			}
			if idp.Valid && idp.Int64 > 0 {
				logging.Error("IDP", idp.Int64)
				userAuthMethodTypes = append(userAuthMethodTypes, domain.UserAuthMethodTypeIDP)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-3n9fl", "Errors.Query.CloseRows")
			}

			return &AuthMethodTypes{
				AuthMethodTypes: userAuthMethodTypes,
				SearchResponse: SearchResponse{
					Count: uint64(len(userAuthMethodTypes)),
				},
			}, nil
		}
}

func prepareAuthMethodsIDPsQuery() (string, error) {
	idpsQuery, _, err := sq.Select(
		userIDPsCountUserID.identifier(),
		userIDPsCountInstanceID.identifier(),
		"COUNT("+userIDPsCountUserID.identifier()+") AS "+userIDPsCountCount.name).
		From(userIDPsCountTable.identifier()).
		GroupBy(
			userIDPsCountUserID.identifier(),
			userIDPsCountInstanceID.identifier(),
		).
		ToSql()
	return idpsQuery, err
}

func prepareAuthMethodQuery(activeOnly bool, includeWithoutDomain bool, queryDomain string) (string, []interface{}, error) {
	q := sq.Select(
		"DISTINCT("+authMethodTypeType.identifier()+")",
		authMethodTypeUserID.identifier(),
		authMethodTypeInstanceID.identifier()).
		From(authMethodTypeTable.identifier())
	if activeOnly {
		q = q.Where(sq.Eq{authMethodTypeState.identifier(): domain.MFAStateReady})
	}
	if queryDomain != "" {
		conditions := sq.Or{
			sq.Eq{authMethodTypeDomain.identifier(): nil},
			sq.Eq{authMethodTypeDomain.identifier(): queryDomain},
		}
		if includeWithoutDomain {
			conditions = append(conditions, sq.Eq{authMethodTypeDomain.identifier(): ""})
		}
		q = q.Where(conditions)
	}

	return q.ToSql()
}

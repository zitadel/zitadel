package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
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
	UserAuthMethodColumnOwnerRemoved = Column{
		name:  projection.UserAuthMethodOwnerRemovedCol,
		table: userAuthMethodTable,
	}

	authMethodTypeTable      = userAuthMethodTable.setAlias("auth_method_types")
	authMethodTypeUserID     = UserAuthMethodColumnUserID.setTable(authMethodTypeTable)
	authMethodTypeInstanceID = UserAuthMethodColumnInstanceID.setTable(authMethodTypeTable)
	authMethodTypeType       = UserAuthMethodColumnMethodType.setTable(authMethodTypeTable)
	authMethodTypeTypes      = Column{
		name:  "method_types",
		table: authMethodTypeTable,
	}
	authMethodTypeState = UserAuthMethodColumnState.setTable(authMethodTypeTable)

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

func (l *AuthMethods) RemoveNoPermission(ctx context.Context, permissionCheck domain.PermissionCheck) {
	removableIndexes := make([]int, 0)
	for i := range l.AuthMethods {
		ctxData := authz.GetCtxData(ctx)
		if ctxData.UserID != l.AuthMethods[i].UserID {
			if err := permissionCheck(ctx, domain.PermissionUserRead, l.AuthMethods[i].ResourceOwner, l.AuthMethods[i].UserID); err != nil {
				removableIndexes = append(removableIndexes, i)
			}
		}
	}
	removed := 0
	for _, removeIndex := range removableIndexes {
		l.AuthMethods = removeAuthMethod(l.AuthMethods, removeIndex-removed)
		removed++
	}
	// reset count as some users could be removed
	l.SearchResponse.Count = uint64(len(l.AuthMethods))
}

func removeAuthMethod(slice []*AuthMethod, s int) []*AuthMethod {
	return append(slice[:s], slice[s+1:]...)
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

func (q *Queries) SearchUserAuthMethods(ctx context.Context, queries *UserAuthMethodSearchQueries, withOwnerRemoved bool) (userAuthMethods *AuthMethods, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareUserAuthMethodsQuery(ctx, q.client)
	eq := sq.Eq{UserAuthMethodColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[UserAuthMethodColumnOwnerRemoved.identifier()] = false
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
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

func (q *Queries) ListActiveUserAuthMethodTypes(ctx context.Context, userID string) (userAuthMethodTypes *AuthMethodTypes, err error) {
	ctxData := authz.GetCtxData(ctx)
	if ctxData.UserID != userID {
		if err := q.checkPermission(ctx, domain.PermissionUserRead, ctxData.OrgID, userID); err != nil {
			return nil, err
		}
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareActiveUserAuthMethodTypesQuery(ctx, q.client)
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

func (q *Queries) ListUserAuthMethodTypesRequired(ctx context.Context, userID string) (requirements *UserAuthMethodRequirements, err error) {
	ctxData := authz.GetCtxData(ctx)
	if ctxData.UserID != userID {
		if err := q.checkPermission(ctx, domain.PermissionUserRead, ctxData.OrgID, userID); err != nil {
			return nil, err
		}
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareUserAuthMethodTypesRequiredQuery(ctx, q.client)
	eq := sq.Eq{
		UserIDCol.identifier():         userID,
		UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-E5ut4", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		requirements, err = scan(row)
		return err
	}, stmt, args...)
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

func prepareUserAuthMethodsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*AuthMethods, error)) {
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
			From(userAuthMethodTable.identifier() + db.Timetravel(call.Took(ctx))).
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

func prepareActiveUserAuthMethodTypesQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*AuthMethodTypes, error)) {
	authMethodsQuery, authMethodsArgs, err := prepareAuthMethodQuery()
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
				userIDPsCountInstanceID.identifier() + " = " + UserInstanceIDCol.identifier() + db.Timetravel(call.Took(ctx))).
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

func prepareUserAuthMethodTypesRequiredQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*UserAuthMethodRequirements, error)) {
	loginPolicyQuery, err := prepareAuthMethodsForceMFAQuery()
	if err != nil {
		return sq.SelectBuilder{}, nil
	}
	return sq.Select(
			UserTypeCol.identifier(),
			forceMFAForce.identifier(),
			forceMFAForceLocalOnly.identifier()).
			From(userTable.identifier()).
			LeftJoin("(" + loginPolicyQuery + ") AS " + forceMFATable.alias + " ON " +
				"(" + forceMFAOrgID.identifier() + " = " + UserInstanceIDCol.identifier() + " OR " + forceMFAOrgID.identifier() + " = " + UserResourceOwnerCol.identifier() + ") AND " +
				forceMFAInstanceID.identifier() + " = " + UserInstanceIDCol.identifier()).
			OrderBy(forceMFAIsDefault.identifier()).
			Limit(1).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*UserAuthMethodRequirements, error) {
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
					return nil, zerrors.ThrowNotFound(err, "QUERY-SF3h2", "Errors.Internal")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-Sf3rt", "Errors.Internal")
			}
			return &UserAuthMethodRequirements{
				UserType:          domain.UserType(userType.Int32),
				ForceMFA:          forceMFA.Bool,
				ForceMFALocalOnly: forceMFALocalOnly.Bool,
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

func prepareAuthMethodQuery() (string, []interface{}, error) {
	return sq.Select(
		"DISTINCT("+authMethodTypeType.identifier()+")",
		authMethodTypeUserID.identifier(),
		authMethodTypeInstanceID.identifier()).
		From(authMethodTypeTable.identifier()).
		Where(sq.Eq{authMethodTypeState.identifier(): domain.MFAStateReady}).
		ToSql()
}

func prepareAuthMethodsForceMFAQuery() (string, error) {
	loginPolicyQuery, _, err := sq.Select(
		forceMFAForce.identifier(),
		forceMFAForceLocalOnly.identifier(),
		forceMFAInstanceID.identifier(),
		forceMFAOrgID.identifier(),
		forceMFAIsDefault.identifier(),
	).
		From(forceMFATable.identifier()).
		ToSql()
	return loginPolicyQuery, err
}

package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Apps struct {
	SearchResponse
	Apps []*App
}

type App struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.AppState
	Sequence      uint64

	ProjectID string
	Name      string

	OIDCConfig *OIDCApp
	SAMLConfig *SAMLApp
	APIConfig  *APIApp
}

type OIDCApp struct {
	RedirectURIs           database.StringArray
	ResponseTypes          database.EnumArray[domain.OIDCResponseType]
	GrantTypes             database.EnumArray[domain.OIDCGrantType]
	AppType                domain.OIDCApplicationType
	ClientID               string
	AuthMethodType         domain.OIDCAuthMethodType
	PostLogoutRedirectURIs database.StringArray
	Version                domain.OIDCVersion
	ComplianceProblems     database.StringArray
	IsDevMode              bool
	AccessTokenType        domain.OIDCTokenType
	AssertAccessTokenRole  bool
	AssertIDTokenRole      bool
	AssertIDTokenUserinfo  bool
	ClockSkew              time.Duration
	AdditionalOrigins      database.StringArray
	AllowedOrigins         database.StringArray
}

type SAMLApp struct {
	Metadata    []byte
	MetadataURL string
	EntityID    string
}

type APIApp struct {
	ClientID       string
	AuthMethodType domain.APIAuthMethodType
}

type AppSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *AppSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

var (
	appsTable = table{
		name:          projection.AppProjectionTable,
		instanceIDCol: projection.AppColumnInstanceID,
	}
	AppColumnID = Column{
		name:  projection.AppColumnID,
		table: appsTable,
	}
	AppColumnName = Column{
		name:  projection.AppColumnName,
		table: appsTable,
	}
	AppColumnProjectID = Column{
		name:  projection.AppColumnProjectID,
		table: appsTable,
	}
	AppColumnCreationDate = Column{
		name:  projection.AppColumnCreationDate,
		table: appsTable,
	}
	AppColumnChangeDate = Column{
		name:  projection.AppColumnChangeDate,
		table: appsTable,
	}
	AppColumnResourceOwner = Column{
		name:  projection.AppColumnResourceOwner,
		table: appsTable,
	}
	AppColumnInstanceID = Column{
		name:  projection.AppColumnInstanceID,
		table: appsTable,
	}
	AppColumnState = Column{
		name:  projection.AppColumnState,
		table: appsTable,
	}
	AppColumnSequence = Column{
		name:  projection.AppColumnSequence,
		table: appsTable,
	}
	AppColumnOwnerRemoved = Column{
		name:  projection.AppColumnOwnerRemoved,
		table: appsTable,
	}
)

var (
	appSAMLConfigsTable = table{
		name:          projection.AppSAMLTable,
		instanceIDCol: projection.AppSAMLConfigColumnInstanceID,
	}
	AppSAMLConfigColumnAppID = Column{
		name:  projection.AppSAMLConfigColumnAppID,
		table: appSAMLConfigsTable,
	}
	AppSAMLConfigColumnEntityID = Column{
		name:  projection.AppSAMLConfigColumnEntityID,
		table: appSAMLConfigsTable,
	}
	AppSAMLConfigColumnMetadata = Column{
		name:  projection.AppSAMLConfigColumnMetadata,
		table: appSAMLConfigsTable,
	}
	AppSAMLConfigColumnMetadataURL = Column{
		name:  projection.AppSAMLConfigColumnMetadataURL,
		table: appSAMLConfigsTable,
	}
)

var (
	appAPIConfigsTable = table{
		name:          projection.AppAPITable,
		instanceIDCol: projection.AppAPIConfigColumnInstanceID,
	}
	AppAPIConfigColumnAppID = Column{
		name:  projection.AppAPIConfigColumnAppID,
		table: appAPIConfigsTable,
	}
	AppAPIConfigColumnClientID = Column{
		name:  projection.AppAPIConfigColumnClientID,
		table: appAPIConfigsTable,
	}
	AppAPIConfigColumnAuthMethod = Column{
		name:  projection.AppAPIConfigColumnAuthMethod,
		table: appAPIConfigsTable,
	}
)

var (
	appOIDCConfigsTable = table{
		name:          projection.AppOIDCTable,
		instanceIDCol: projection.AppOIDCConfigColumnInstanceID,
	}
	AppOIDCConfigColumnAppID = Column{
		name:  projection.AppOIDCConfigColumnAppID,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnVersion = Column{
		name:  projection.AppOIDCConfigColumnVersion,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnClientID = Column{
		name:  projection.AppOIDCConfigColumnClientID,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnRedirectUris = Column{
		name:  projection.AppOIDCConfigColumnRedirectUris,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnResponseTypes = Column{
		name:  projection.AppOIDCConfigColumnResponseTypes,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnGrantTypes = Column{
		name:  projection.AppOIDCConfigColumnGrantTypes,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnApplicationType = Column{
		name:  projection.AppOIDCConfigColumnApplicationType,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnAuthMethodType = Column{
		name:  projection.AppOIDCConfigColumnAuthMethodType,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnPostLogoutRedirectUris = Column{
		name:  projection.AppOIDCConfigColumnPostLogoutRedirectUris,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnDevMode = Column{
		name:  projection.AppOIDCConfigColumnDevMode,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnAccessTokenType = Column{
		name:  projection.AppOIDCConfigColumnAccessTokenType,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnAccessTokenRoleAssertion = Column{
		name:  projection.AppOIDCConfigColumnAccessTokenRoleAssertion,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnIDTokenRoleAssertion = Column{
		name:  projection.AppOIDCConfigColumnIDTokenRoleAssertion,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnIDTokenUserinfoAssertion = Column{
		name:  projection.AppOIDCConfigColumnIDTokenUserinfoAssertion,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnClockSkew = Column{
		name:  projection.AppOIDCConfigColumnClockSkew,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnAdditionalOrigins = Column{
		name:  projection.AppOIDCConfigColumnAdditionalOrigins,
		table: appOIDCConfigsTable,
	}
)

func (q *Queries) AppByProjectAndAppID(ctx context.Context, shouldTriggerBulk bool, projectID, appID string, withOwnerRemoved bool) (_ *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.AppProjection.Trigger(ctx)
	}

	stmt, scan := prepareAppQuery(ctx, q.client)
	eq := sq.Eq{
		AppColumnID.identifier():         appID,
		AppColumnProjectID.identifier():  projectID,
		AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-AFDgg", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) AppByID(ctx context.Context, appID string, withOwnerRemoved bool) (_ *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAppQuery(ctx, q.client)
	eq := sq.Eq{
		AppColumnID.identifier():         appID,
		AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-immt9", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) AppBySAMLEntityID(ctx context.Context, entityID string, withOwnerRemoved bool) (_ *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAppQuery(ctx, q.client)
	eq := sq.Eq{
		AppSAMLConfigColumnEntityID.identifier(): entityID,
		AppColumnInstanceID.identifier():         authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-JgUop", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ProjectByClientID(ctx context.Context, appID string, withOwnerRemoved bool) (_ *Project, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareProjectByAppQuery(ctx, q.client)
	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[ProjectColumnOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(sq.And{
		eq,
		sq.Or{
			sq.Eq{AppOIDCConfigColumnClientID.identifier(): appID},
			sq.Eq{AppAPIConfigColumnClientID.identifier(): appID},
			sq.Eq{AppSAMLConfigColumnAppID.identifier(): appID},
		},
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-XhJi3", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ProjectIDFromOIDCClientID(ctx context.Context, appID string, withOwnerRemoved bool) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareProjectIDByAppQuery(ctx, q.client)
	eq := sq.Eq{
		AppOIDCConfigColumnClientID.identifier(): appID,
		AppColumnInstanceID.identifier():         authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return "", errors.ThrowInternal(err, "QUERY-7d92U", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ProjectIDFromClientID(ctx context.Context, appID string, withOwnerRemoved bool) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareProjectIDByAppQuery(ctx, q.client)
	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	where := sq.And{
		eq,
		sq.Or{
			sq.Eq{AppOIDCConfigColumnClientID.identifier(): appID},
			sq.Eq{AppAPIConfigColumnClientID.identifier(): appID},
			sq.Eq{AppSAMLConfigColumnAppID.identifier(): appID},
		},
	}
	query, args, err := stmt.Where(where).ToSql()
	if err != nil {
		return "", errors.ThrowInternal(err, "QUERY-SDfg3", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ProjectByOIDCClientID(ctx context.Context, id string, withOwnerRemoved bool) (_ *Project, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareProjectByAppQuery(ctx, q.client)
	eq := sq.Eq{
		AppOIDCConfigColumnClientID.identifier(): id,
		AppColumnInstanceID.identifier():         authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-XhJi4", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) AppByOIDCClientID(ctx context.Context, clientID string, withOwnerRemoved bool) (_ *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAppQuery(ctx, q.client)
	eq := sq.Eq{
		AppOIDCConfigColumnClientID.identifier(): clientID,
		AppColumnInstanceID.identifier():         authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-JgVop", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) AppByClientID(ctx context.Context, clientID string, withOwnerRemoved bool) (_ *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAppQuery(ctx, q.client)
	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(sq.And{
		eq,
		sq.Or{
			sq.Eq{AppOIDCConfigColumnClientID.identifier(): clientID},
			sq.Eq{AppAPIConfigColumnClientID.identifier(): clientID},
		},
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dfge2", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) SearchApps(ctx context.Context, queries *AppSearchQueries, withOwnerRemoved bool) (_ *Apps, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareAppsQuery(ctx, q.client)
	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-fajp8", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-aJnZL", "Errors.Internal")
	}
	apps, err := scan(rows)
	if err != nil {
		return nil, err
	}
	apps.LatestSequence, err = q.latestSequence(ctx, appsTable)
	return apps, err
}

func (q *Queries) SearchClientIDs(ctx context.Context, queries *AppSearchQueries, withOwnerRemoved bool) (_ []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareClientIDsQuery(ctx, q.client)
	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[AppColumnOwnerRemoved.identifier()] = false
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-fajp8", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-aJnZL", "Errors.Internal")
	}
	return scan(rows)
}

func NewAppNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(AppColumnName, value, method)
}

func NewAppProjectIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AppColumnProjectID, id, TextEquals)
}

func prepareAppQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*App, error)) {
	return sq.Select(
			AppColumnID.identifier(),
			AppColumnName.identifier(),
			AppColumnProjectID.identifier(),
			AppColumnCreationDate.identifier(),
			AppColumnChangeDate.identifier(),
			AppColumnResourceOwner.identifier(),
			AppColumnState.identifier(),
			AppColumnSequence.identifier(),

			AppAPIConfigColumnAppID.identifier(),
			AppAPIConfigColumnClientID.identifier(),
			AppAPIConfigColumnAuthMethod.identifier(),

			AppOIDCConfigColumnAppID.identifier(),
			AppOIDCConfigColumnVersion.identifier(),
			AppOIDCConfigColumnClientID.identifier(),
			AppOIDCConfigColumnRedirectUris.identifier(),
			AppOIDCConfigColumnResponseTypes.identifier(),
			AppOIDCConfigColumnGrantTypes.identifier(),
			AppOIDCConfigColumnApplicationType.identifier(),
			AppOIDCConfigColumnAuthMethodType.identifier(),
			AppOIDCConfigColumnPostLogoutRedirectUris.identifier(),
			AppOIDCConfigColumnDevMode.identifier(),
			AppOIDCConfigColumnAccessTokenType.identifier(),
			AppOIDCConfigColumnAccessTokenRoleAssertion.identifier(),
			AppOIDCConfigColumnIDTokenRoleAssertion.identifier(),
			AppOIDCConfigColumnIDTokenUserinfoAssertion.identifier(),
			AppOIDCConfigColumnClockSkew.identifier(),
			AppOIDCConfigColumnAdditionalOrigins.identifier(),

			AppSAMLConfigColumnAppID.identifier(),
			AppSAMLConfigColumnEntityID.identifier(),
			AppSAMLConfigColumnMetadata.identifier(),
			AppSAMLConfigColumnMetadataURL.identifier(),
		).From(appsTable.identifier()).
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*App, error) {
			app := new(App)

			var (
				apiConfig  = sqlAPIConfig{}
				oidcConfig = sqlOIDCConfig{}
				samlConfig = sqlSAMLConfig{}
			)

			err := row.Scan(
				&app.ID,
				&app.Name,
				&app.ProjectID,
				&app.CreationDate,
				&app.ChangeDate,
				&app.ResourceOwner,
				&app.State,
				&app.Sequence,

				&apiConfig.appID,
				&apiConfig.clientID,
				&apiConfig.authMethod,

				&oidcConfig.appID,
				&oidcConfig.version,
				&oidcConfig.clientID,
				&oidcConfig.redirectUris,
				&oidcConfig.responseTypes,
				&oidcConfig.grantTypes,
				&oidcConfig.applicationType,
				&oidcConfig.authMethodType,
				&oidcConfig.postLogoutRedirectUris,
				&oidcConfig.devMode,
				&oidcConfig.accessTokenType,
				&oidcConfig.accessTokenRoleAssertion,
				&oidcConfig.iDTokenRoleAssertion,
				&oidcConfig.iDTokenUserinfoAssertion,
				&oidcConfig.clockSkew,
				&oidcConfig.additionalOrigins,

				&samlConfig.appID,
				&samlConfig.entityID,
				&samlConfig.metadata,
				&samlConfig.metadataURL,
			)

			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-pCP8P", "Errors.App.NotExisting")
				}
				return nil, errors.ThrowInternal(err, "QUERY-0R2Nw", "Errors.Internal")
			}

			apiConfig.set(app)
			oidcConfig.set(app)
			samlConfig.set(app)

			return app, nil
		}
}

func prepareProjectIDByAppQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (projectID string, err error)) {
	return sq.Select(
			AppColumnProjectID.identifier(),
		).From(appsTable.identifier()).
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (projectID string, err error) {
			err = row.Scan(
				&projectID,
			)

			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return "", errors.ThrowNotFound(err, "QUERY-aKcc2", "Errors.Project.NotExisting")
				}
				return "", errors.ThrowInternal(err, "QUERY-3A5TG", "Errors.Internal")
			}

			return projectID, nil
		}
}

func prepareProjectByAppQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Project, error)) {
	return sq.Select(
			ProjectColumnID.identifier(),
			ProjectColumnCreationDate.identifier(),
			ProjectColumnChangeDate.identifier(),
			ProjectColumnResourceOwner.identifier(),
			ProjectColumnState.identifier(),
			ProjectColumnSequence.identifier(),
			ProjectColumnName.identifier(),
			ProjectColumnProjectRoleAssertion.identifier(),
			ProjectColumnProjectRoleCheck.identifier(),
			ProjectColumnHasProjectCheck.identifier(),
			ProjectColumnPrivateLabelingSetting.identifier(),
		).From(projectsTable.identifier()).
			Join(join(AppColumnProjectID, ProjectColumnID)).
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Project, error) {
			p := new(Project)
			err := row.Scan(
				&p.ID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.State,
				&p.Sequence,
				&p.Name,
				&p.ProjectRoleAssertion,
				&p.ProjectRoleCheck,
				&p.HasProjectCheck,
				&p.PrivateLabelingSetting,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-fk2fs", "Errors.Project.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-dj2FF", "Errors.Internal")
			}
			return p, nil
		}
}

func prepareAppsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Apps, error)) {
	return sq.Select(
			AppColumnID.identifier(),
			AppColumnName.identifier(),
			AppColumnProjectID.identifier(),
			AppColumnCreationDate.identifier(),
			AppColumnChangeDate.identifier(),
			AppColumnResourceOwner.identifier(),
			AppColumnState.identifier(),
			AppColumnSequence.identifier(),

			AppAPIConfigColumnAppID.identifier(),
			AppAPIConfigColumnClientID.identifier(),
			AppAPIConfigColumnAuthMethod.identifier(),

			AppOIDCConfigColumnAppID.identifier(),
			AppOIDCConfigColumnVersion.identifier(),
			AppOIDCConfigColumnClientID.identifier(),
			AppOIDCConfigColumnRedirectUris.identifier(),
			AppOIDCConfigColumnResponseTypes.identifier(),
			AppOIDCConfigColumnGrantTypes.identifier(),
			AppOIDCConfigColumnApplicationType.identifier(),
			AppOIDCConfigColumnAuthMethodType.identifier(),
			AppOIDCConfigColumnPostLogoutRedirectUris.identifier(),
			AppOIDCConfigColumnDevMode.identifier(),
			AppOIDCConfigColumnAccessTokenType.identifier(),
			AppOIDCConfigColumnAccessTokenRoleAssertion.identifier(),
			AppOIDCConfigColumnIDTokenRoleAssertion.identifier(),
			AppOIDCConfigColumnIDTokenUserinfoAssertion.identifier(),
			AppOIDCConfigColumnClockSkew.identifier(),
			AppOIDCConfigColumnAdditionalOrigins.identifier(),

			AppSAMLConfigColumnAppID.identifier(),
			AppSAMLConfigColumnEntityID.identifier(),
			AppSAMLConfigColumnMetadata.identifier(),
			AppSAMLConfigColumnMetadataURL.identifier(),
			countColumn.identifier(),
		).From(appsTable.identifier()).
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(row *sql.Rows) (*Apps, error) {
			apps := &Apps{Apps: []*App{}}

			for row.Next() {
				app := new(App)
				var (
					apiConfig  = sqlAPIConfig{}
					oidcConfig = sqlOIDCConfig{}
					samlConfig = sqlSAMLConfig{}
				)

				err := row.Scan(
					&app.ID,
					&app.Name,
					&app.ProjectID,
					&app.CreationDate,
					&app.ChangeDate,
					&app.ResourceOwner,
					&app.State,
					&app.Sequence,

					&apiConfig.appID,
					&apiConfig.clientID,
					&apiConfig.authMethod,

					&oidcConfig.appID,
					&oidcConfig.version,
					&oidcConfig.clientID,
					&oidcConfig.redirectUris,
					&oidcConfig.responseTypes,
					&oidcConfig.grantTypes,
					&oidcConfig.applicationType,
					&oidcConfig.authMethodType,
					&oidcConfig.postLogoutRedirectUris,
					&oidcConfig.devMode,
					&oidcConfig.accessTokenType,
					&oidcConfig.accessTokenRoleAssertion,
					&oidcConfig.iDTokenRoleAssertion,
					&oidcConfig.iDTokenUserinfoAssertion,
					&oidcConfig.clockSkew,
					&oidcConfig.additionalOrigins,

					&samlConfig.appID,
					&samlConfig.entityID,
					&samlConfig.metadata,
					&samlConfig.metadataURL,

					&apps.Count,
				)

				if err != nil {
					return nil, errors.ThrowInternal(err, "QUERY-0R2Nw", "Errors.Internal")
				}

				apiConfig.set(app)
				oidcConfig.set(app)
				samlConfig.set(app)

				apps.Apps = append(apps.Apps, app)
			}

			return apps, nil
		}
}

func prepareClientIDsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) ([]string, error)) {
	return sq.Select(
			AppAPIConfigColumnClientID.identifier(),
			AppOIDCConfigColumnClientID.identifier(),
		).From(appsTable.identifier()).
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(rows *sql.Rows) ([]string, error) {
			ids := database.StringArray{}

			for rows.Next() {
				var apiID sql.NullString
				var oidcID sql.NullString
				if err := rows.Scan(
					&apiID,
					&oidcID,
				); err != nil {
					return nil, errors.ThrowInternal(err, "QUERY-0R2Nw", "Errors.Internal")
				}
				if apiID.Valid {
					ids = append(ids, apiID.String)
				} else if oidcID.Valid {
					ids = append(ids, oidcID.String)
				}
			}

			return ids, nil
		}
}

type sqlOIDCConfig struct {
	appID                    sql.NullString
	version                  sql.NullInt32
	clientID                 sql.NullString
	redirectUris             database.StringArray
	applicationType          sql.NullInt16
	authMethodType           sql.NullInt16
	postLogoutRedirectUris   database.StringArray
	devMode                  sql.NullBool
	accessTokenType          sql.NullInt16
	accessTokenRoleAssertion sql.NullBool
	iDTokenRoleAssertion     sql.NullBool
	iDTokenUserinfoAssertion sql.NullBool
	clockSkew                sql.NullInt64
	additionalOrigins        database.StringArray
	responseTypes            database.EnumArray[domain.OIDCResponseType]
	grantTypes               database.EnumArray[domain.OIDCGrantType]
}

func (c sqlOIDCConfig) set(app *App) {
	if !c.appID.Valid {
		return
	}
	app.OIDCConfig = &OIDCApp{
		Version:                domain.OIDCVersion(c.version.Int32),
		ClientID:               c.clientID.String,
		RedirectURIs:           c.redirectUris,
		AppType:                domain.OIDCApplicationType(c.applicationType.Int16),
		AuthMethodType:         domain.OIDCAuthMethodType(c.authMethodType.Int16),
		PostLogoutRedirectURIs: c.postLogoutRedirectUris,
		IsDevMode:              c.devMode.Bool,
		AccessTokenType:        domain.OIDCTokenType(c.accessTokenType.Int16),
		AssertAccessTokenRole:  c.accessTokenRoleAssertion.Bool,
		AssertIDTokenRole:      c.iDTokenRoleAssertion.Bool,
		AssertIDTokenUserinfo:  c.iDTokenUserinfoAssertion.Bool,
		ClockSkew:              time.Duration(c.clockSkew.Int64),
		AdditionalOrigins:      c.additionalOrigins,
		ResponseTypes:          c.responseTypes,
		GrantTypes:             c.grantTypes,
	}
	compliance := domain.GetOIDCCompliance(app.OIDCConfig.Version, app.OIDCConfig.AppType, app.OIDCConfig.GrantTypes, app.OIDCConfig.ResponseTypes, app.OIDCConfig.AuthMethodType, app.OIDCConfig.RedirectURIs)
	app.OIDCConfig.ComplianceProblems = compliance.Problems

	var err error
	app.OIDCConfig.AllowedOrigins, err = domain.OIDCOriginAllowList(app.OIDCConfig.RedirectURIs, app.OIDCConfig.AdditionalOrigins)
	logging.LogWithFields("app", app.ID).OnError(err).Warn("unable to set allowed origins")
}

type sqlSAMLConfig struct {
	appID       sql.NullString
	entityID    sql.NullString
	metadataURL sql.NullString
	metadata    []byte
}

func (c sqlSAMLConfig) set(app *App) {
	if !c.appID.Valid {
		return
	}
	app.SAMLConfig = &SAMLApp{
		MetadataURL: c.metadataURL.String,
		Metadata:    c.metadata,
		EntityID:    c.entityID.String,
	}
}

type sqlAPIConfig struct {
	appID      sql.NullString
	clientID   sql.NullString
	authMethod sql.NullInt16
}

func (c sqlAPIConfig) set(app *App) {
	if !c.appID.Valid {
		return
	}
	app.APIConfig = &APIApp{
		ClientID:       c.clientID.String,
		AuthMethodType: domain.APIAuthMethodType(c.authMethod.Int16),
	}
}

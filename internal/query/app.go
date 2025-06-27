package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/muhlemmer/gu"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	RedirectURIs             database.TextArray[string]
	ResponseTypes            database.NumberArray[domain.OIDCResponseType]
	GrantTypes               database.NumberArray[domain.OIDCGrantType]
	AppType                  domain.OIDCApplicationType
	ClientID                 string
	AuthMethodType           domain.OIDCAuthMethodType
	PostLogoutRedirectURIs   database.TextArray[string]
	Version                  domain.OIDCVersion
	ComplianceProblems       database.TextArray[string]
	IsDevMode                bool
	AccessTokenType          domain.OIDCTokenType
	AssertAccessTokenRole    bool
	AssertIDTokenRole        bool
	AssertIDTokenUserinfo    bool
	ClockSkew                time.Duration
	AdditionalOrigins        database.TextArray[string]
	AllowedOrigins           database.TextArray[string]
	SkipNativeAppSuccessPage bool
	BackChannelLogoutURI     string
	LoginVersion             domain.LoginVersion
	LoginBaseURI             *string
}

type SAMLApp struct {
	Metadata     []byte
	MetadataURL  string
	EntityID     string
	LoginVersion domain.LoginVersion
	LoginBaseURI *string
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
)

var (
	appSAMLConfigsTable = table{
		name:          projection.AppSAMLTable,
		instanceIDCol: projection.AppSAMLConfigColumnInstanceID,
	}
	AppSAMLConfigColumnInstanceID = Column{
		name:  projection.AppSAMLConfigColumnInstanceID,
		table: appSAMLConfigsTable,
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
	AppSAMLConfigColumnLoginVersion = Column{
		name:  projection.AppSAMLConfigColumnLoginVersion,
		table: appSAMLConfigsTable,
	}
	AppSAMLConfigColumnLoginBaseURI = Column{
		name:  projection.AppSAMLConfigColumnLoginBaseURI,
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
	AppOIDCConfigColumnInstanceID = Column{
		name:  projection.AppOIDCConfigColumnInstanceID,
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
	AppOIDCConfigColumnSkipNativeAppSuccessPage = Column{
		name:  projection.AppOIDCConfigColumnSkipNativeAppSuccessPage,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnBackChannelLogoutURI = Column{
		name:  projection.AppOIDCConfigColumnBackChannelLogoutURI,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnLoginVersion = Column{
		name:  projection.AppOIDCConfigColumnLoginVersion,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnLoginBaseURI = Column{
		name:  projection.AppOIDCConfigColumnLoginBaseURI,
		table: appOIDCConfigsTable,
	}
)

func (q *Queries) AppByProjectAndAppID(ctx context.Context, shouldTriggerBulk bool, projectID, appID string) (app *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerAppProjection")
		ctx, err = projection.AppProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	stmt, scan := prepareAppQuery(false)
	eq := sq.Eq{
		AppColumnID.identifier():         appID,
		AppColumnProjectID.identifier():  projectID,
		AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-AFDgg", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		app, err = scan(row)
		return err
	}, query, args...)
	return app, err
}

func (q *Queries) AppByIDWithPermission(ctx context.Context, appID string, activeOnly bool, permissionCheck domain.PermissionCheck) (*App, error) {
	app, err := q.AppByID(ctx, appID, activeOnly)
	if err != nil {
		return nil, err
	}

	if err := appCheckPermission(ctx, app.ResourceOwner, app.ProjectID, permissionCheck); err != nil {
		return nil, err
	}

	return app, nil
}

func (q *Queries) AppByID(ctx context.Context, appID string, activeOnly bool) (app *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAppQuery(activeOnly)
	eq := sq.Eq{
		AppColumnID.identifier():         appID,
		AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if activeOnly {
		eq[AppColumnState.identifier()] = domain.AppStateActive
		eq[ProjectColumnState.identifier()] = domain.ProjectStateActive
		eq[OrgColumnState.identifier()] = domain.OrgStateActive
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-immt9", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		app, err = scan(row)
		return err
	}, query, args...)
	return app, err
}

func (q *Queries) ProjectByClientID(ctx context.Context, appID string) (project *Project, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareProjectByAppQuery()
	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	query, args, err := stmt.Where(sq.And{
		eq,
		sq.Or{
			sq.Eq{AppOIDCConfigColumnClientID.identifier(): appID},
			sq.Eq{AppAPIConfigColumnClientID.identifier(): appID},
			sq.Eq{AppSAMLConfigColumnAppID.identifier(): appID},
		},
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-XhJi3", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		project, err = scan(row)
		return err
	}, query, args...)
	return project, err
}

//go:embed app_oidc_project_permission.sql
var appOIDCProjectPermissionQuery string

func (q *Queries) CheckProjectPermissionByClientID(ctx context.Context, clientID, userID string) (_ *projectPermission, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var p *projectPermission
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		p, err = scanProjectPermissionByClientID(row)
		return err
	}, appOIDCProjectPermissionQuery,
		authz.GetInstance(ctx).InstanceID(),
		clientID,
		domain.AppStateActive,
		domain.ProjectStateActive,
		userID,
		domain.UserStateActive,
		domain.ProjectGrantStateActive,
		domain.UserGrantStateActive,
	)
	return p, err
}

//go:embed app_saml_project_permission.sql
var appSAMLProjectPermissionQuery string

func (q *Queries) CheckProjectPermissionByEntityID(ctx context.Context, entityID, userID string) (_ *projectPermission, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var p *projectPermission
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		p, err = scanProjectPermissionByClientID(row)
		return err
	}, appSAMLProjectPermissionQuery,
		authz.GetInstance(ctx).InstanceID(),
		entityID,
		domain.AppStateActive,
		domain.ProjectStateActive,
		userID,
		domain.UserStateActive,
		domain.ProjectGrantStateActive,
		domain.UserGrantStateActive,
	)
	return p, err
}

type projectPermission struct {
	HasProjectChecked  bool
	ProjectRoleChecked bool
}

func scanProjectPermissionByClientID(row *sql.Row) (*projectPermission, error) {
	var hasProjectChecked, projectRoleChecked sql.NullBool
	err := row.Scan(
		&hasProjectChecked,
		&projectRoleChecked,
	)
	if err != nil || !hasProjectChecked.Valid || !projectRoleChecked.Valid {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowNotFound(err, "QUERY-4tq8wCTCgf", "Errors.App.NotFound")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-NwH4lAqlZC", "Errors.Internal")
	}
	return &projectPermission{
		HasProjectChecked:  hasProjectChecked.Bool,
		ProjectRoleChecked: projectRoleChecked.Bool,
	}, nil
}

func (q *Queries) ProjectIDFromClientID(ctx context.Context, appID string) (id string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareProjectIDByAppQuery()
	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
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
		return "", zerrors.ThrowInternal(err, "QUERY-SDfg3", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		id, err = scan(row)
		return err
	}, query, args...)
	return id, err
}

func (q *Queries) AppByOIDCClientID(ctx context.Context, clientID string) (app *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareOIDCAppQuery()
	eq := sq.Eq{
		AppOIDCConfigColumnClientID.identifier(): clientID,
		AppColumnInstanceID.identifier():         authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-JgVop", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		app, err = scan(row)
		return err
	}, query, args...)
	return app, err
}

func (q *Queries) AppByClientID(ctx context.Context, clientID string) (app *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAppQuery(true)
	eq := sq.Eq{
		AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		AppColumnState.identifier():      domain.AppStateActive,
		ProjectColumnState.identifier():  domain.ProjectStateActive,
		OrgColumnState.identifier():      domain.OrgStateActive,
	}
	query, args, err := stmt.Where(sq.And{
		eq,
		sq.Or{
			sq.Eq{AppOIDCConfigColumnClientID.identifier(): clientID},
			sq.Eq{AppAPIConfigColumnClientID.identifier(): clientID},
		},
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dfge2", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		app, err = scan(row)
		return err
	}, query, args...)
	return app, err
}

func (q *Queries) SearchApps(ctx context.Context, queries *AppSearchQueries, permissionCheck domain.PermissionCheck) (*Apps, error) {
	apps, err := q.searchApps(ctx, queries, PermissionV2(ctx, permissionCheck))
	if err != nil {
		return nil, err
	}

	if permissionCheck != nil && !authz.GetFeatures(ctx).PermissionCheckV2 {
		apps.Apps = appsCheckPermission(ctx, apps.Apps, permissionCheck)
	}
	return apps, nil
}

func (q *Queries) searchApps(ctx context.Context, queries *AppSearchQueries, isPermissionV2Enabled bool) (apps *Apps, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareAppsQuery()
	query = appPermissionCheckV2(ctx, query, isPermissionV2Enabled, queries)

	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-fajp8", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		apps, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-h9TeF", "Errors.Internal")
	}
	apps.State, err = q.latestState(ctx, appsTable)
	return apps, err
}

func appPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool, queries *AppSearchQueries) sq.SelectBuilder {
	if !enabled {
		return query
	}

	join, args := PermissionClause(
		ctx,
		AppColumnResourceOwner,
		domain.PermissionProjectAppRead,
		SingleOrgPermissionOption(queries.Queries),
		WithProjectsPermissionOption(AppColumnProjectID),
	)
	return query.JoinClause(join, args...)
}

func (q *Queries) SearchClientIDs(ctx context.Context, queries *AppSearchQueries, shouldTriggerBulk bool) (ids []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerAppProjection")
		ctx, err = projection.AppProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareClientIDsQuery()
	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-fajp8", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		ids, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-aJnZL", "Errors.Internal")
	}
	return ids, nil
}

func (q *Queries) OIDCClientLoginVersion(ctx context.Context, clientID string) (loginVersion domain.LoginVersion, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareLoginVersionByOIDCClientID()
	eq := sq.Eq{
		AppOIDCConfigColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		AppOIDCConfigColumnClientID.identifier():   clientID,
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return domain.LoginVersionUnspecified, zerrors.ThrowInvalidArgument(err, "QUERY-WEh31", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		loginVersion, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return domain.LoginVersionUnspecified, zerrors.ThrowInternal(err, "QUERY-W2gsa", "Errors.Internal")
	}
	return loginVersion, nil
}

func (q *Queries) SAMLAppLoginVersion(ctx context.Context, appID string) (loginVersion domain.LoginVersion, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareLoginVersionBySAMLAppID()
	eq := sq.Eq{
		AppSAMLConfigColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		AppSAMLConfigColumnAppID.identifier():      appID,
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return domain.LoginVersionUnspecified, zerrors.ThrowInvalidArgument(err, "QUERY-TnaciwZfp3", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		loginVersion, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return domain.LoginVersionUnspecified, zerrors.ThrowInternal(err, "QUERY-lvDDwRzIoP", "Errors.Internal")
	}
	return loginVersion, nil
}

func appCheckPermission(ctx context.Context, resourceOwner string, projectID string, permissionCheck domain.PermissionCheck) error {
	return permissionCheck(ctx, domain.PermissionProjectAppRead, resourceOwner, projectID)
}

// appsCheckPermission returns only the apps that the user in context has permission to read
func appsCheckPermission(ctx context.Context, apps []*App, permissionCheck domain.PermissionCheck) []*App {
	return slices.DeleteFunc(apps, func(app *App) bool {
		return permissionCheck(ctx, domain.PermissionProjectAppRead, app.ResourceOwner, app.ProjectID) != nil
	})
}

func NewAppNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(AppColumnName, value, method)
}

func NewAppStateSearchQuery(value domain.AppState) (SearchQuery, error) {
	return NewNumberQuery(AppColumnState, int(value), NumberEquals)
}

func NewAppProjectIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AppColumnProjectID, id, TextEquals)
}

func prepareAppQuery(activeOnly bool) (sq.SelectBuilder, func(*sql.Row) (*App, error)) {
	query := sq.Select(
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
		AppOIDCConfigColumnSkipNativeAppSuccessPage.identifier(),
		AppOIDCConfigColumnBackChannelLogoutURI.identifier(),
		AppOIDCConfigColumnLoginVersion.identifier(),
		AppOIDCConfigColumnLoginBaseURI.identifier(),

		AppSAMLConfigColumnAppID.identifier(),
		AppSAMLConfigColumnEntityID.identifier(),
		AppSAMLConfigColumnMetadata.identifier(),
		AppSAMLConfigColumnMetadataURL.identifier(),
		AppSAMLConfigColumnLoginVersion.identifier(),
		AppSAMLConfigColumnLoginBaseURI.identifier(),
	).From(appsTable.identifier()).
		PlaceholderFormat(sq.Dollar)

	if activeOnly {
		return query.
				LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
				LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
				LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID)).
				LeftJoin(join(ProjectColumnID, AppColumnProjectID)).
				LeftJoin(join(OrgColumnID, AppColumnResourceOwner)),
			scanApp
	}
	return query.
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID)),
		scanApp
}

func scanApp(row *sql.Row) (*App, error) {
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
		&oidcConfig.skipNativeAppSuccessPage,
		&oidcConfig.backChannelLogoutURI,
		&oidcConfig.loginVersion,
		&oidcConfig.loginBaseURI,

		&samlConfig.appID,
		&samlConfig.entityID,
		&samlConfig.metadata,
		&samlConfig.metadataURL,
		&samlConfig.loginVersion,
		&samlConfig.loginBaseURI,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowNotFound(err, "QUERY-pCP8P", "Errors.App.NotExisting")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-4SJlx", "Errors.Internal")
	}

	apiConfig.set(app)
	oidcConfig.set(app)
	samlConfig.set(app)

	return app, nil
}

func prepareOIDCAppQuery() (sq.SelectBuilder, func(*sql.Row) (*App, error)) {
	return sq.Select(
			AppColumnID.identifier(),
			AppColumnName.identifier(),
			AppColumnProjectID.identifier(),
			AppColumnCreationDate.identifier(),
			AppColumnChangeDate.identifier(),
			AppColumnResourceOwner.identifier(),
			AppColumnState.identifier(),
			AppColumnSequence.identifier(),

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
			AppOIDCConfigColumnSkipNativeAppSuccessPage.identifier(),
			AppOIDCConfigColumnBackChannelLogoutURI.identifier(),
			AppOIDCConfigColumnLoginVersion.identifier(),
			AppOIDCConfigColumnLoginBaseURI.identifier(),
		).From(appsTable.identifier()).
			Join(join(AppOIDCConfigColumnAppID, AppColumnID)).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*App, error) {
			app := new(App)

			var (
				oidcConfig = sqlOIDCConfig{}
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
				&oidcConfig.skipNativeAppSuccessPage,
				&oidcConfig.backChannelLogoutURI,
				&oidcConfig.loginVersion,
				&oidcConfig.loginBaseURI,
			)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-Fdfax", "Errors.App.NotExisting")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-aE7iE", "Errors.Internal")
			}

			oidcConfig.set(app)

			return app, nil
		}
}

func prepareProjectIDByAppQuery() (sq.SelectBuilder, func(*sql.Row) (projectID string, err error)) {
	return sq.Select(
			AppColumnProjectID.identifier(),
		).From(appsTable.identifier()).
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID)).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (projectID string, err error) {
			err = row.Scan(
				&projectID,
			)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return "", zerrors.ThrowNotFound(err, "QUERY-aKcc2", "Errors.Project.NotExisting")
				}
				return "", zerrors.ThrowInternal(err, "QUERY-3A5TG", "Errors.Internal")
			}

			return projectID, nil
		}
}

func prepareProjectByAppQuery() (sq.SelectBuilder, func(*sql.Row) (*Project, error)) {
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
			LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID)).
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
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-yxTMh", "Errors.Project.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-dj2FF", "Errors.Internal")
			}
			return p, nil
		}
}

func prepareAppsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Apps, error)) {
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
			AppOIDCConfigColumnSkipNativeAppSuccessPage.identifier(),
			AppOIDCConfigColumnBackChannelLogoutURI.identifier(),
			AppOIDCConfigColumnLoginVersion.identifier(),
			AppOIDCConfigColumnLoginBaseURI.identifier(),

			AppSAMLConfigColumnAppID.identifier(),
			AppSAMLConfigColumnEntityID.identifier(),
			AppSAMLConfigColumnMetadata.identifier(),
			AppSAMLConfigColumnMetadataURL.identifier(),
			AppSAMLConfigColumnLoginVersion.identifier(),
			AppSAMLConfigColumnLoginBaseURI.identifier(),
			countColumn.identifier(),
		).From(appsTable.identifier()).
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID)).
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
					&oidcConfig.skipNativeAppSuccessPage,
					&oidcConfig.backChannelLogoutURI,
					&oidcConfig.loginVersion,
					&oidcConfig.loginBaseURI,

					&samlConfig.appID,
					&samlConfig.entityID,
					&samlConfig.metadata,
					&samlConfig.metadataURL,
					&samlConfig.loginVersion,
					&samlConfig.loginBaseURI,

					&apps.Count,
				)

				if err != nil {
					return nil, zerrors.ThrowInternal(err, "QUERY-XGWAX", "Errors.Internal")
				}

				apiConfig.set(app)
				oidcConfig.set(app)
				samlConfig.set(app)

				apps.Apps = append(apps.Apps, app)
			}

			return apps, nil
		}
}

func prepareClientIDsQuery() (sq.SelectBuilder, func(*sql.Rows) ([]string, error)) {
	return sq.Select(
			AppAPIConfigColumnClientID.identifier(),
			AppOIDCConfigColumnClientID.identifier(),
		).From(appsTable.identifier()).
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
			PlaceholderFormat(sq.Dollar), func(rows *sql.Rows) ([]string, error) {
			ids := database.TextArray[string]{}

			for rows.Next() {
				var apiID sql.NullString
				var oidcID sql.NullString
				if err := rows.Scan(
					&apiID,
					&oidcID,
				); err != nil {
					return nil, zerrors.ThrowInternal(err, "QUERY-0R2Nw", "Errors.Internal")
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

func prepareLoginVersionByOIDCClientID() (sq.SelectBuilder, func(*sql.Row) (domain.LoginVersion, error)) {
	return sq.Select(
			AppOIDCConfigColumnLoginVersion.identifier(),
		).From(appOIDCConfigsTable.identifier()).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (domain.LoginVersion, error) {
			var loginVersion sql.NullInt16
			if err := row.Scan(
				&loginVersion,
			); err != nil {
				return domain.LoginVersionUnspecified, zerrors.ThrowInternal(err, "QUERY-KL2io", "Errors.Internal")
			}
			return domain.LoginVersion(loginVersion.Int16), nil
		}
}

func prepareLoginVersionBySAMLAppID() (sq.SelectBuilder, func(*sql.Row) (domain.LoginVersion, error)) {
	return sq.Select(
			AppSAMLConfigColumnLoginVersion.identifier(),
		).From(appSAMLConfigsTable.identifier()).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (domain.LoginVersion, error) {
			var loginVersion sql.NullInt16
			if err := row.Scan(
				&loginVersion,
			); err != nil {
				return domain.LoginVersionUnspecified, zerrors.ThrowInternal(err, "QUERY-KbzaCnaziI", "Errors.Internal")
			}
			return domain.LoginVersion(loginVersion.Int16), nil
		}
}

type sqlOIDCConfig struct {
	appID                    sql.NullString
	version                  sql.NullInt32
	clientID                 sql.NullString
	redirectUris             database.TextArray[string]
	applicationType          sql.NullInt16
	authMethodType           sql.NullInt16
	postLogoutRedirectUris   database.TextArray[string]
	devMode                  sql.NullBool
	accessTokenType          sql.NullInt16
	accessTokenRoleAssertion sql.NullBool
	iDTokenRoleAssertion     sql.NullBool
	iDTokenUserinfoAssertion sql.NullBool
	clockSkew                sql.NullInt64
	additionalOrigins        database.TextArray[string]
	responseTypes            database.NumberArray[domain.OIDCResponseType]
	grantTypes               database.NumberArray[domain.OIDCGrantType]
	skipNativeAppSuccessPage sql.NullBool
	backChannelLogoutURI     sql.NullString
	loginVersion             sql.NullInt16
	loginBaseURI             sql.NullString
}

func (c sqlOIDCConfig) set(app *App) {
	if !c.appID.Valid {
		return
	}
	app.OIDCConfig = &OIDCApp{
		Version:                  domain.OIDCVersion(c.version.Int32),
		ClientID:                 c.clientID.String,
		RedirectURIs:             c.redirectUris,
		AppType:                  domain.OIDCApplicationType(c.applicationType.Int16),
		AuthMethodType:           domain.OIDCAuthMethodType(c.authMethodType.Int16),
		PostLogoutRedirectURIs:   c.postLogoutRedirectUris,
		IsDevMode:                c.devMode.Bool,
		AccessTokenType:          domain.OIDCTokenType(c.accessTokenType.Int16),
		AssertAccessTokenRole:    c.accessTokenRoleAssertion.Bool,
		AssertIDTokenRole:        c.iDTokenRoleAssertion.Bool,
		AssertIDTokenUserinfo:    c.iDTokenUserinfoAssertion.Bool,
		ClockSkew:                time.Duration(c.clockSkew.Int64),
		AdditionalOrigins:        c.additionalOrigins,
		ResponseTypes:            c.responseTypes,
		GrantTypes:               c.grantTypes,
		SkipNativeAppSuccessPage: c.skipNativeAppSuccessPage.Bool,
		BackChannelLogoutURI:     c.backChannelLogoutURI.String,
		LoginVersion:             domain.LoginVersion(c.loginVersion.Int16),
	}
	if c.loginBaseURI.Valid {
		app.OIDCConfig.LoginBaseURI = &c.loginBaseURI.String
	}
	compliance := domain.GetOIDCCompliance(gu.Ptr(app.OIDCConfig.Version), gu.Ptr(app.OIDCConfig.AppType), app.OIDCConfig.GrantTypes, app.OIDCConfig.ResponseTypes, gu.Ptr(app.OIDCConfig.AuthMethodType), app.OIDCConfig.RedirectURIs)
	app.OIDCConfig.ComplianceProblems = compliance.Problems

	var err error
	app.OIDCConfig.AllowedOrigins, err = domain.OIDCOriginAllowList(app.OIDCConfig.RedirectURIs, app.OIDCConfig.AdditionalOrigins)
	logging.LogWithFields("app", app.ID).OnError(err).Warn("unable to set allowed origins")
}

type sqlSAMLConfig struct {
	appID        sql.NullString
	entityID     sql.NullString
	metadataURL  sql.NullString
	metadata     []byte
	loginVersion sql.NullInt16
	loginBaseURI sql.NullString
}

func (c sqlSAMLConfig) set(app *App) {
	if !c.appID.Valid {
		return
	}
	app.SAMLConfig = &SAMLApp{
		EntityID:     c.entityID.String,
		MetadataURL:  c.metadataURL.String,
		Metadata:     c.metadata,
		LoginVersion: domain.LoginVersion(c.loginVersion.Int16),
	}
	if c.loginBaseURI.Valid {
		app.SAMLConfig.LoginBaseURI = &c.loginBaseURI.String
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

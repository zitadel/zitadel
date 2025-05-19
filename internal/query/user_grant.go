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
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UserGrant struct {
	// ID represents the aggregate id (id of the user grant)
	ID           string                     `json:"id,omitempty"`
	CreationDate time.Time                  `json:"creation_date,omitempty"`
	ChangeDate   time.Time                  `json:"change_date,omitempty"`
	Sequence     uint64                     `json:"sequence,omitempty"`
	Roles        database.TextArray[string] `json:"roles,omitempty"`
	// GrantID represents the project grant id
	GrantID string                `json:"grant_id,omitempty"`
	State   domain.UserGrantState `json:"state,omitempty"`

	UserID             string          `json:"user_id,omitempty"`
	Username           string          `json:"username,omitempty"`
	UserType           domain.UserType `json:"user_type,omitempty"`
	UserResourceOwner  string          `json:"user_resource_owner,omitempty"`
	FirstName          string          `json:"first_name,omitempty"`
	LastName           string          `json:"last_name,omitempty"`
	Email              string          `json:"email,omitempty"`
	DisplayName        string          `json:"display_name,omitempty"`
	AvatarURL          string          `json:"avatar_url,omitempty"`
	PreferredLoginName string          `json:"preferred_login_name,omitempty"`

	ResourceOwner    string `json:"resource_owner,omitempty"`
	OrgName          string `json:"org_name,omitempty"`
	OrgPrimaryDomain string `json:"org_primary_domain,omitempty"`

	ProjectID   string `json:"project_id,omitempty"`
	ProjectName string `json:"project_name,omitempty"`

	GrantedOrgID     string `json:"granted_org_id,omitempty"`
	GrantedOrgName   string `json:"granted_org_name,omitempty"`
	GrantedOrgDomain string `json:"granted_org_domain,omitempty"`
}

type UserGrants struct {
	SearchResponse
	UserGrants []*UserGrant
}

type UserGrantsQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *UserGrantsQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func NewUserGrantUserIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantUserID, id, TextEquals)
}

func NewUserGrantProjectIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantProjectID, id, TextEquals)
}

func NewUserGrantProjectIDsSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(UserGrantProjectID, list, ListIn)
}

func NewUserGrantProjectOwnerSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnResourceOwner, id, TextEquals)
}

func NewUserGrantResourceOwnerSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantResourceOwner, id, TextEquals)
}

func NewUserGrantGrantIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantGrantID, id, TextEquals)
}

func NewUserGrantIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantID, id, TextEquals)
}

func NewUserGrantUserTypeQuery(typ domain.UserType) (SearchQuery, error) {
	return NewNumberQuery(UserTypeCol, typ, NumberEquals)
}

func NewUserGrantDisplayNameQuery(displayName string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanDisplayNameCol, displayName, method)
}

func NewUserGrantEmailQuery(email string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanEmailCol, email, method)
}

func NewUserGrantFirstNameQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanFirstNameCol, value, method)
}

func NewUserGrantLastNameQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanLastNameCol, value, method)
}

func NewUserGrantUsernameQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(UserUsernameCol, value, method)
}

func NewUserGrantDomainQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(OrgColumnDomain, value, method)
}

func NewUserGrantOrgNameQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(OrgColumnName, value, method)
}

func NewUserGrantProjectNameQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnName, value, method)
}

func NewUserGrantRoleQuery(value string) (SearchQuery, error) {
	return NewTextQuery(UserGrantRoles, value, TextListContains)
}

func NewUserGrantStateQuery(value domain.UserGrantState) (SearchQuery, error) {
	return NewNumberQuery(UserGrantState, value, NumberEquals)
}

func NewUserGrantWithGrantedQuery(owner string) (SearchQuery, error) {
	orgQuery, err := NewUserGrantResourceOwnerSearchQuery(owner)
	if err != nil {
		return nil, err
	}
	projectQuery, err := NewUserGrantProjectOwnerSearchQuery(owner)
	if err != nil {
		return nil, err
	}
	return NewOrQuery(orgQuery, projectQuery)
}

func NewUserGrantContainsRolesSearchQuery(roles ...string) (SearchQuery, error) {
	r := make([]interface{}, len(roles))
	for i, role := range roles {
		r[i] = role
	}
	return NewListQuery(UserGrantRoles, r, ListIn)
}

var (
	userGrantTable = table{
		name:          projection.UserGrantProjectionTable,
		instanceIDCol: projection.UserGrantInstanceID,
	}
	UserGrantID = Column{
		name:  projection.UserGrantID,
		table: userGrantTable,
	}
	UserGrantResourceOwner = Column{
		name:  projection.UserGrantResourceOwner,
		table: userGrantTable,
	}
	UserGrantInstanceID = Column{
		name:  projection.UserGrantInstanceID,
		table: userGrantTable,
	}
	UserGrantCreationDate = Column{
		name:  projection.UserGrantCreationDate,
		table: userGrantTable,
	}
	UserGrantChangeDate = Column{
		name:  projection.UserGrantChangeDate,
		table: userGrantTable,
	}
	UserGrantSequence = Column{
		name:  projection.UserGrantSequence,
		table: userGrantTable,
	}
	UserGrantUserID = Column{
		name:  projection.UserGrantUserID,
		table: userGrantTable,
	}
	UserGrantProjectID = Column{
		name:  projection.UserGrantProjectID,
		table: userGrantTable,
	}
	UserGrantGrantID = Column{
		name:  projection.UserGrantGrantID,
		table: userGrantTable,
	}
	UserGrantRoles = Column{
		name:  projection.UserGrantRoles,
		table: userGrantTable,
	}
	UserGrantState = Column{
		name:  projection.UserGrantState,
		table: userGrantTable,
	}
	GrantedOrgsTable = table{
		name:          projection.OrgProjectionTable,
		alias:         "granted_orgs",
		instanceIDCol: projection.OrgColumnInstanceID,
	}
	GrantedOrgColumnId = Column{
		name:  projection.OrgColumnID,
		table: GrantedOrgsTable,
	}
	GrantedOrgColumnName = Column{
		name:  projection.OrgColumnName,
		table: GrantedOrgsTable,
	}
	GrantedOrgColumnDomain = Column{
		name:  projection.OrgColumnDomain,
		table: GrantedOrgsTable,
	}
)

func (q *Queries) UserGrant(ctx context.Context, shouldTriggerBulk bool, queries ...SearchQuery) (grant *UserGrant, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerUserGrantProjection")
		ctx, err = projection.UserGrantProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareUserGrantQuery(ctx, q.client)
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{UserGrantInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Fa1KW", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		grant, err = scan(row)
		return err
	}, stmt, args...)
	return grant, err
}

func (q *Queries) UserGrants(ctx context.Context, queries *UserGrantsQueries, shouldTriggerBulk bool) (grants *UserGrants, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerUserGrantProjection")
		ctx, err = projection.UserGrantProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("unable to trigger")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareUserGrantsQuery(ctx, q.client)
	eq := sq.Eq{UserGrantInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-wXnQR", "Errors.Query.SQLStatement")
	}

	latestState, err := q.latestState(ctx, userGrantTable)
	if err != nil {
		return nil, err
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		grants, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}

	grants.State = latestState
	return grants, nil
}

func prepareUserGrantQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*UserGrant, error)) {
	return sq.Select(
			UserGrantID.identifier(),
			UserGrantCreationDate.identifier(),
			UserGrantChangeDate.identifier(),
			UserGrantSequence.identifier(),
			UserGrantGrantID.identifier(),
			UserGrantRoles.identifier(),
			UserGrantState.identifier(),

			UserGrantUserID.identifier(),
			UserUsernameCol.identifier(),
			UserTypeCol.identifier(),
			UserResourceOwnerCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanDisplayNameCol.identifier(),
			HumanAvatarURLCol.identifier(),
			LoginNameNameCol.identifier(),

			UserGrantResourceOwner.identifier(),
			OrgColumnName.identifier(),
			OrgColumnDomain.identifier(),

			UserGrantProjectID.identifier(),
			ProjectColumnName.identifier(),

			GrantedOrgColumnId.identifier(),
			GrantedOrgColumnName.identifier(),
			GrantedOrgColumnDomain.identifier(),
		).
			From(userGrantTable.identifier()).
			LeftJoin(join(UserIDCol, UserGrantUserID)).
			LeftJoin(join(HumanUserIDCol, UserGrantUserID)).
			LeftJoin(join(OrgColumnID, UserGrantResourceOwner)).
			LeftJoin(join(ProjectColumnID, UserGrantProjectID)).
			LeftJoin(join(GrantedOrgColumnId, UserResourceOwnerCol)).
			LeftJoin(join(LoginNameUserIDCol, UserGrantUserID) + db.Timetravel(call.Took(ctx))).
			Where(
				sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
			).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*UserGrant, error) {
			g := new(UserGrant)

			var (
				username           sql.NullString
				firstName          sql.NullString
				userType           sql.NullInt32
				userOwner          sql.NullString
				lastName           sql.NullString
				email              sql.NullString
				displayName        sql.NullString
				avatarURL          sql.NullString
				preferredLoginName sql.NullString

				orgName   sql.NullString
				orgDomain sql.NullString

				projectName sql.NullString

				grantedOrgID     sql.NullString
				grantedOrgName   sql.NullString
				grantedOrgDomain sql.NullString
			)

			err := row.Scan(
				&g.ID,
				&g.CreationDate,
				&g.ChangeDate,
				&g.Sequence,
				&g.GrantID,
				&g.Roles,
				&g.State,

				&g.UserID,
				&username,
				&userType,
				&userOwner,
				&firstName,
				&lastName,
				&email,
				&displayName,
				&avatarURL,
				&preferredLoginName,

				&g.ResourceOwner,
				&orgName,
				&orgDomain,

				&g.ProjectID,
				&projectName,

				&grantedOrgID,
				&grantedOrgName,
				&grantedOrgDomain,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-wIPkA", "Errors.UserGrant.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-oQPcP", "Errors.Internal")
			}

			g.Username = username.String
			g.UserType = domain.UserType(userType.Int32)
			g.UserResourceOwner = userOwner.String
			g.FirstName = firstName.String
			g.LastName = lastName.String
			g.Email = email.String
			g.DisplayName = displayName.String
			g.AvatarURL = avatarURL.String
			g.PreferredLoginName = preferredLoginName.String
			g.OrgName = orgName.String
			g.OrgPrimaryDomain = orgDomain.String
			g.ProjectName = projectName.String
			g.GrantedOrgID = grantedOrgID.String
			g.GrantedOrgName = grantedOrgName.String
			g.GrantedOrgDomain = grantedOrgDomain.String
			return g, nil
		}
}

func prepareUserGrantsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*UserGrants, error)) {
	return sq.Select(
			UserGrantID.identifier(),
			UserGrantCreationDate.identifier(),
			UserGrantChangeDate.identifier(),
			UserGrantSequence.identifier(),
			UserGrantGrantID.identifier(),
			UserGrantRoles.identifier(),
			UserGrantState.identifier(),

			UserGrantUserID.identifier(),
			UserUsernameCol.identifier(),
			UserTypeCol.identifier(),
			UserResourceOwnerCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanDisplayNameCol.identifier(),
			HumanAvatarURLCol.identifier(),
			LoginNameNameCol.identifier(),

			UserGrantResourceOwner.identifier(),
			OrgColumnName.identifier(),
			OrgColumnDomain.identifier(),

			UserGrantProjectID.identifier(),
			ProjectColumnName.identifier(),

			GrantedOrgColumnId.identifier(),
			GrantedOrgColumnName.identifier(),
			GrantedOrgColumnDomain.identifier(),

			countColumn.identifier(),
		).
			From(userGrantTable.identifier()).
			LeftJoin(join(UserIDCol, UserGrantUserID)).
			LeftJoin(join(HumanUserIDCol, UserGrantUserID)).
			LeftJoin(join(OrgColumnID, UserGrantResourceOwner)).
			LeftJoin(join(ProjectColumnID, UserGrantProjectID)).
			LeftJoin(join(GrantedOrgColumnId, UserResourceOwnerCol)).
			LeftJoin(join(LoginNameUserIDCol, UserGrantUserID) + db.Timetravel(call.Took(ctx))).
			Where(
				sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
			).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*UserGrants, error) {
			userGrants := make([]*UserGrant, 0)
			var count uint64
			for rows.Next() {
				g := new(UserGrant)

				var (
					username           sql.NullString
					userType           sql.NullInt32
					userOwner          sql.NullString
					firstName          sql.NullString
					lastName           sql.NullString
					email              sql.NullString
					displayName        sql.NullString
					avatarURL          sql.NullString
					preferredLoginName sql.NullString

					orgName   sql.NullString
					orgDomain sql.NullString

					grantedOrgID     sql.NullString
					grantedOrgName   sql.NullString
					grantedOrgDomain sql.NullString

					projectName sql.NullString
				)

				err := rows.Scan(
					&g.ID,
					&g.CreationDate,
					&g.ChangeDate,
					&g.Sequence,
					&g.GrantID,
					&g.Roles,
					&g.State,

					&g.UserID,
					&username,
					&userType,
					&userOwner,
					&firstName,
					&lastName,
					&email,
					&displayName,
					&avatarURL,
					&preferredLoginName,

					&g.ResourceOwner,
					&orgName,
					&orgDomain,

					&g.ProjectID,
					&projectName,

					&grantedOrgID,
					&grantedOrgName,
					&grantedOrgDomain,

					&count,
				)
				if err != nil {
					return nil, err
				}

				g.Username = username.String
				g.UserType = domain.UserType(userType.Int32)
				g.UserResourceOwner = userOwner.String
				g.FirstName = firstName.String
				g.LastName = lastName.String
				g.Email = email.String
				g.DisplayName = displayName.String
				g.AvatarURL = avatarURL.String
				g.PreferredLoginName = preferredLoginName.String
				g.OrgName = orgName.String
				g.OrgPrimaryDomain = orgDomain.String
				g.ProjectName = projectName.String
				g.GrantedOrgID = grantedOrgID.String
				g.GrantedOrgName = grantedOrgName.String
				g.GrantedOrgDomain = grantedOrgDomain.String

				userGrants = append(userGrants, g)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-iGvmP", "Errors.Query.CloseRows")
			}

			return &UserGrants{
				UserGrants: userGrants,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

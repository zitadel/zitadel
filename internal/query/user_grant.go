package query

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UserGrant struct {
	// ID represents the aggregate id (id of the user grant)
	ID              string                     `json:"id,omitempty"`
	CreationDate    time.Time                  `json:"creation_date,omitempty"`
	ChangeDate      time.Time                  `json:"change_date,omitempty"`
	Sequence        uint64                     `json:"sequence,omitempty"`
	Roles           database.TextArray[string] `json:"roles,omitempty"`
	RoleInformation []Role                     `json:"-"`
	// GrantID represents the project grant id
	GrantID string                `json:"grant_id,omitempty"`
	State   domain.UserGrantState `json:"state,omitempty"`

	UserID                  string          `json:"user_id,omitempty"`
	Username                string          `json:"username,omitempty"`
	UserType                domain.UserType `json:"user_type,omitempty"`
	UserResourceOwner       string          `json:"user_resource_owner,omitempty"`
	UserResourceOwnerName   string          `json:"user_resource_owner_name,omitempty"`
	UserResourceOwnerDomain string          `json:"user_resource_owner_domain,omitempty"`
	FirstName               string          `json:"first_name,omitempty"`
	LastName                string          `json:"last_name,omitempty"`
	Email                   string          `json:"email,omitempty"`
	DisplayName             string          `json:"display_name,omitempty"`
	AvatarURL               string          `json:"avatar_url,omitempty"`
	PreferredLoginName      string          `json:"preferred_login_name,omitempty"`

	ResourceOwner    string `json:"resource_owner,omitempty"`
	OrgName          string `json:"org_name,omitempty"`
	OrgPrimaryDomain string `json:"org_primary_domain,omitempty"`

	ProjectResourceOwner string `json:"project_resource_owner,omitempty"`
	ProjectID            string `json:"project_id,omitempty"`
	ProjectName          string `json:"project_name,omitempty"`

	GrantedOrgID     string `json:"granted_org_id,omitempty"`
	GrantedOrgName   string `json:"granted_org_name,omitempty"`
	GrantedOrgDomain string `json:"granted_org_domain,omitempty"`
}

type Role struct {
	Key         string `json:"role_key,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	GroupName   string `json:"group_name,omitempty"`
}

type UserGrants struct {
	SearchResponse
	UserGrants []*UserGrant
}

func userGrantsCheckPermission(ctx context.Context, grants *UserGrants, permissionCheck domain.PermissionCheck) {
	grants.UserGrants = slices.DeleteFunc(grants.UserGrants,
		func(grant *UserGrant) bool {
			return userGrantCheckPermission(ctx, grant.ResourceOwner, grant.ProjectID, grant.GrantID, grant.UserID, permissionCheck) != nil
		},
	)
}

func userGrantCheckPermission(ctx context.Context, resourceOwner, projectID, grantID, userID string, permissionCheck domain.PermissionCheck) error {
	// you should always be able to read your own permissions
	if authz.GetCtxData(ctx).UserID == userID {
		return nil
	}
	// check permission on the project grant
	if grantID != "" {
		return permissionCheck(ctx, domain.PermissionUserGrantRead, resourceOwner, grantID)
	}
	// check on project
	return permissionCheck(ctx, domain.PermissionUserGrantRead, resourceOwner, projectID)
}

type UserGrantsQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func userGrantPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool, queries *UserGrantsQueries) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		UserGrantResourceOwner,
		domain.PermissionUserGrantRead,
		SingleOrgPermissionOption(queries.Queries),
		WithProjectsPermissionOption(UserGrantProjectID),
		OwnedRowsPermissionOption(UserGrantUserID),
	)
	return query.JoinClause(join, args...)
}

func NewUserGrantUserIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantUserID, id, TextEquals)
}

func NewUserGrantInUserIDsSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(UserGrantUserID, list, ListIn)
}

func NewUserGrantProjectIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantProjectID, id, TextEquals)
}

func NewUserGrantProjectOwnerSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnResourceOwner, id, TextEquals)
}

func NewUserGrantResourceOwnerSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantResourceOwner, id, TextEquals)
}

func NewUserGrantUserResourceOwnerSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserResourceOwnerCol, id, TextEquals)
}

func NewUserGrantGrantIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantGrantID, id, TextEquals)
}

func NewUserGrantIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(UserGrantID, id, TextEquals)
}

func NewUserGrantInIDsSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(UserGrantID, list, ListIn)
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

	UserOrgsTable = table{
		name:          projection.OrgProjectionTable,
		alias:         "user_orgs",
		instanceIDCol: projection.OrgColumnInstanceID,
	}
	UserOrgColumnId = Column{
		name:  projection.OrgColumnID,
		table: UserOrgsTable,
	}
	UserOrgColumnName = Column{
		name:  projection.OrgColumnName,
		table: UserOrgsTable,
	}
	UserOrgColumnDomain = Column{
		name:  projection.OrgColumnDomain,
		table: UserOrgsTable,
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

	query, scan := prepareUserGrantQuery()
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

func (q *Queries) UserGrants(ctx context.Context, queries *UserGrantsQueries, shouldTriggerBulk bool, permissionCheck domain.PermissionCheck) (*UserGrants, error) {
	// removed as permission v2 is not implemented yet for project grant level permissions
	// permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	grants, err := q.userGrants(ctx, queries, shouldTriggerBulk, false)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil { // && !authz.GetFeatures(ctx).PermissionCheckV2 {
		userGrantsCheckPermission(ctx, grants, permissionCheck)
	}
	return grants, nil
}

func (q *Queries) userGrants(ctx context.Context, queries *UserGrantsQueries, shouldTriggerBulk bool, permissionCheckV2 bool) (grants *UserGrants, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerUserGrantProjection")
		ctx, err = projection.UserGrantProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("unable to trigger")
		traceSpan.EndWithError(err)
	}

	query, scan := queries.prepareUserGrantsQuery(ctx, permissionCheckV2)
	stmt, args, err := query.ToSql()
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

func prepareUserGrantQuery() (sq.SelectBuilder, func(*sql.Row) (*UserGrant, error)) {
	return sq.Select(
			UserGrantID.identifier(),
			UserGrantCreationDate.identifier(),
			UserGrantChangeDate.identifier(),
			UserGrantSequence.identifier(),
			UserGrantGrantID.identifier(),
			UserGrantRoles.identifier(),
			"roles.role_information",
			UserGrantState.identifier(),

			UserGrantUserID.identifier(),
			UserUsernameCol.identifier(),
			UserTypeCol.identifier(),
			UserOrgColumnId.identifier(),
			UserOrgColumnName.identifier(),
			UserOrgColumnDomain.identifier(),
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
			ProjectColumnResourceOwner.identifier(),

			GrantedOrgColumnId.identifier(),
			GrantedOrgColumnName.identifier(),
			GrantedOrgColumnDomain.identifier(),
		).
			From(userGrantTable.identifier()).
			LeftJoin(join(UserIDCol, UserGrantUserID)).
			LeftJoin(join(HumanUserIDCol, UserGrantUserID)).
			LeftJoin(join(OrgColumnID, UserGrantResourceOwner)).
			LeftJoin(join(ProjectColumnID, UserGrantProjectID)).
			LeftJoin(join(UserOrgColumnId, UserResourceOwnerCol)).
			LeftJoin(join(ProjectGrantColumnGrantID, UserGrantGrantID) + " AND " + ProjectGrantColumnProjectID.identifier() + " = " + UserGrantProjectID.identifier()).
			LeftJoin(join(GrantedOrgColumnId, ProjectGrantColumnGrantedOrgID)).
			LeftJoin(join(LoginNameUserIDCol, UserGrantUserID)).
			LeftJoin("LATERAL (" + rolesInfoQuery + ") as roles ON true").
			Where(
				sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
			).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*UserGrant, error) {
			g := new(UserGrant)

			var (
				roles []byte

				username           sql.NullString
				firstName          sql.NullString
				userType           sql.NullInt32
				userOwner          sql.NullString
				userOwnerName      sql.NullString
				userOwnerDomain    sql.NullString
				lastName           sql.NullString
				email              sql.NullString
				displayName        sql.NullString
				avatarURL          sql.NullString
				preferredLoginName sql.NullString

				orgName   sql.NullString
				orgDomain sql.NullString

				projectName          sql.NullString
				projectResourceOwner sql.NullString

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
				&roles,
				&g.State,

				&g.UserID,
				&username,
				&userType,
				&userOwner,
				&userOwnerName,
				&userOwnerDomain,
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
				&projectResourceOwner,

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
			if len(roles) > 0 {
				err = json.Unmarshal(roles, &g.RoleInformation)
				if err != nil {
					return nil, err
				}
			}

			g.Username = username.String
			g.UserType = domain.UserType(userType.Int32)
			g.UserResourceOwner = userOwner.String
			g.UserResourceOwnerName = userOwnerName.String
			g.UserResourceOwnerDomain = userOwnerDomain.String
			g.FirstName = firstName.String
			g.LastName = lastName.String
			g.Email = email.String
			g.DisplayName = displayName.String
			g.AvatarURL = avatarURL.String
			g.PreferredLoginName = preferredLoginName.String
			g.OrgName = orgName.String
			g.OrgPrimaryDomain = orgDomain.String
			g.ProjectName = projectName.String
			g.ProjectResourceOwner = projectResourceOwner.String
			g.GrantedOrgID = grantedOrgID.String
			g.GrantedOrgName = grantedOrgName.String
			g.GrantedOrgDomain = grantedOrgDomain.String
			return g, nil
		}
}

// userGrantJoins tracks which joins of the user grants list query a statement needs.
type userGrantJoins struct {
	users         bool
	humans        bool
	orgs          bool
	projects      bool
	userOrgs      bool
	projectGrants bool
	grantedOrgs   bool
	loginNames    bool
}

func allUserGrantJoins() userGrantJoins {
	return userGrantJoins{
		users:         true,
		humans:        true,
		orgs:          true,
		projects:      true,
		userOrgs:      true,
		projectGrants: true,
		grantedOrgs:   true,
		loginNames:    true,
	}
}

// add marks the joins needed to reference col. It returns false if col
// belongs to a table the list query does not join.
func (j *userGrantJoins) add(col Column) bool {
	switch col.table {
	case userGrantTable:
	case userTable:
		j.users = true
	case humanTable:
		j.humans = true
	case orgsTable:
		j.orgs = true
	case projectsTable:
		j.projects = true
	case UserOrgsTable:
		j.users, j.userOrgs = true, true
	case projectGrantsTable:
		j.projectGrants = true
	case GrantedOrgsTable:
		j.projectGrants, j.grantedOrgs = true, true
	case loginNameTable:
		j.loginNames = true
	default:
		return false
	}
	return true
}

// apply adds the marked joins in the order of the full list query.
// Without the login_names3 join, the primary login name filter of the full
// query is kept as a check that the user has a login name row (see userGrantHasLoginName).
func (j userGrantJoins) apply(query sq.SelectBuilder) sq.SelectBuilder {
	if j.users {
		query = query.LeftJoin(join(UserIDCol, UserGrantUserID))
	}
	if j.humans {
		query = query.LeftJoin(join(HumanUserIDCol, UserGrantUserID))
	}
	if j.orgs {
		query = query.LeftJoin(join(OrgColumnID, UserGrantResourceOwner))
	}
	if j.projects {
		query = query.LeftJoin(join(ProjectColumnID, UserGrantProjectID))
	}
	if j.userOrgs {
		query = query.LeftJoin(join(UserOrgColumnId, UserResourceOwnerCol))
	}
	if j.projectGrants {
		query = query.LeftJoin(userGrantProjectGrantJoin)
	}
	if j.grantedOrgs {
		query = query.LeftJoin(join(GrantedOrgColumnId, ProjectGrantColumnGrantedOrgID))
	}
	if j.loginNames {
		return query.LeftJoin(join(LoginNameUserIDCol, UserGrantUserID)).
			Where(sq.Eq{LoginNameIsPrimaryCol.identifier(): true})
	}
	return query.Where(userGrantHasLoginName)
}

var (
	userGrantProjectGrantJoin = join(ProjectGrantColumnGrantID, UserGrantGrantID) + " AND " + ProjectGrantColumnProjectID.identifier() + " = " + UserGrantProjectID.identifier()

	// userGrantHasLoginName replaces the login_names3 join and its is_primary filter
	// when only matching grants are counted or paged. login_names3 returns exactly one
	// primary row for every user in login_names3_users: the user name, or the user name
	// at the primary domain of the org if the domain policy requires it (every org with
	// domains has exactly one primary domain). So both keep the same grants, but this
	// check is a primary key lookup instead of expanding the view for every grant.
	userGrantHasLoginName = sq.Expr("EXISTS (SELECT 1 FROM " + projection.LoginNameUserProjectionTable +
		" WHERE " + projection.LoginNameUserProjectionTable + "." + projection.LoginNameUserInstanceIDCol + " = " + UserGrantInstanceID.identifier() +
		" AND " + projection.LoginNameUserProjectionTable + "." + projection.LoginNameUserIDCol + " = " + UserGrantUserID.identifier() + ")")
)

// filterJoins returns the joins needed by the search queries, or all joins if a
// query references a column the list query does not join itself.
func (q *UserGrantsQueries) filterJoins() userGrantJoins {
	var joins userGrantJoins
	unknown := false
	// the match func never reports a match, so every column of every query is visited
	searchQueriesUse(q.Queries, func(col Column) bool {
		if !joins.add(col) {
			unknown = true
		}
		return false
	})
	if unknown {
		return allUserGrantJoins()
	}
	return joins
}

// matchingGrants selects columns of the grants matching the instance, permission
// check and search queries, joining only joins.
func (q *UserGrantsQueries) matchingGrants(ctx context.Context, permissionCheckV2 bool, joins userGrantJoins, columns ...string) sq.SelectBuilder {
	query := joins.apply(sq.Select(columns...).From(userGrantTable.identifier()))
	query = userGrantPermissionCheckV2(ctx, query, permissionCheckV2, q)
	for _, filter := range q.Queries {
		query = filter.toQuery(query)
	}
	return query.Where(sq.Eq{UserGrantInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()})
}

func (q *UserGrantsQueries) pageOrderClause(sortExpr, idExpr string) string {
	dir := ""
	if !q.Asc {
		dir = " DESC"
	}
	clause := sortExpr + dir
	if q.SortingColumn.identifier() != UserGrantID.identifier() {
		clause += ", " + idExpr + dir
	}
	return clause
}

// userGrantsColumns are the columns scanned by prepareUserGrantsQuery, without the total.
var userGrantsColumns = []string{
	UserGrantID.identifier(),
	UserGrantCreationDate.identifier(),
	UserGrantChangeDate.identifier(),
	UserGrantSequence.identifier(),
	UserGrantGrantID.identifier(),
	UserGrantRoles.identifier(),
	"roles.role_information",
	UserGrantState.identifier(),

	UserGrantUserID.identifier(),
	UserUsernameCol.identifier(),
	UserTypeCol.identifier(),
	UserOrgColumnId.identifier(),
	UserOrgColumnName.identifier(),
	UserOrgColumnDomain.identifier(),
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
	ProjectColumnResourceOwner.identifier(),

	GrantedOrgColumnId.identifier(),
	GrantedOrgColumnName.identifier(),
	GrantedOrgColumnDomain.identifier(),
}

// prepareUserGrantsQuery creates the statement listing user grants and returns a matching scan function.
// The total is a scalar COUNT(*) over the filters, joining only what they reference.
// For a page (limit or offset), the IDs of the page are selected first, joining only what the
// filters and the sorting need, so OFFSET skips narrow rows. Users, orgs, projects, login names
// and role information are joined for the returned page only.
// Without limit and offset nothing is skipped, so the grants are joined in a single pass.
// With a sorting column the grant id breaks ties, so pages are stable.
func (q *UserGrantsQueries) prepareUserGrantsQuery(ctx context.Context, permissionCheckV2 bool) (sq.SelectBuilder, func(*sql.Rows) (*UserGrants, error)) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	filterJoins := q.filterJoins()
	total := sq.ConcatExpr("(", q.matchingGrants(ctx, permissionCheckV2, filterJoins, "COUNT(*)"), ") AS total")
	rolesJoin := "LATERAL (" + rolesInfoQuery + ") as roles ON true"
	sorting := !q.SortingColumn.isZero()

	if q.Limit == 0 && q.Offset == 0 {
		query := q.matchingGrants(ctx, permissionCheckV2, allUserGrantJoins(), userGrantsColumns...).
			Column(total).
			LeftJoin(rolesJoin)
		if sorting {
			query = query.OrderByClause(q.pageOrderClause(q.SortingColumn.orderBy(), UserGrantID.identifier()))
		}
		return query.PlaceholderFormat(sq.Dollar), scanUserGrants
	}

	pageJoins := filterJoins
	pageColumns := []string{UserGrantID.identifier()}
	if sorting {
		if !pageJoins.add(q.SortingColumn) {
			pageJoins = allUserGrantJoins()
		}
		pageColumns = append(pageColumns, q.SortingColumn.orderBy()+" AS sort_col")
	}
	page := q.matchingGrants(ctx, permissionCheckV2, pageJoins, pageColumns...)
	if sorting {
		page = page.OrderByClause(q.pageOrderClause("sort_col", UserGrantID.identifier()))
	}
	if q.Offset > 0 {
		page = page.Offset(q.Offset)
	}
	if q.Limit > 0 {
		page = page.Limit(q.Limit)
	}

	query := sq.Select(userGrantsColumns...).
		Column(total).
		FromSelect(page, "page").
		Join(userGrantTable.identifier()+" ON "+UserGrantID.identifier()+" = page.id AND "+UserGrantInstanceID.identifier()+" = ?", instanceID).
		LeftJoin(join(UserIDCol, UserGrantUserID)).
		LeftJoin(join(HumanUserIDCol, UserGrantUserID)).
		LeftJoin(join(OrgColumnID, UserGrantResourceOwner)).
		LeftJoin(join(ProjectColumnID, UserGrantProjectID)).
		LeftJoin(join(UserOrgColumnId, UserResourceOwnerCol)).
		LeftJoin(userGrantProjectGrantJoin).
		LeftJoin(join(GrantedOrgColumnId, ProjectGrantColumnGrantedOrgID)).
		LeftJoin(join(LoginNameUserIDCol, UserGrantUserID)).
		LeftJoin(rolesJoin).
		Where(sq.Eq{LoginNameIsPrimaryCol.identifier(): true})
	if sorting {
		query = query.OrderByClause(q.pageOrderClause("page.sort_col", "page.id"))
	}
	return query.PlaceholderFormat(sq.Dollar), scanUserGrants
}

func scanUserGrants(rows *sql.Rows) (*UserGrants, error) {
	userGrants := make([]*UserGrant, 0)
	var count uint64
	for rows.Next() {
		g := new(UserGrant)

		var (
			roles []byte

			username           sql.NullString
			userType           sql.NullInt32
			userOwner          sql.NullString
			userOwnerName      sql.NullString
			userOwnerDomain    sql.NullString
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

			projectName          sql.NullString
			projectResourceOwner sql.NullString
		)

		err := rows.Scan(
			&g.ID,
			&g.CreationDate,
			&g.ChangeDate,
			&g.Sequence,
			&g.GrantID,
			&g.Roles,
			&roles,
			&g.State,

			&g.UserID,
			&username,
			&userType,
			&userOwner,
			&userOwnerName,
			&userOwnerDomain,
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
			&projectResourceOwner,

			&grantedOrgID,
			&grantedOrgName,
			&grantedOrgDomain,

			&count,
		)
		if err != nil {
			return nil, err
		}
		if len(roles) > 0 {
			err = json.Unmarshal(roles, &g.RoleInformation)
			if err != nil {
				return nil, err
			}
		}

		g.Username = username.String
		g.UserType = domain.UserType(userType.Int32)
		g.UserResourceOwner = userOwner.String
		g.UserResourceOwnerName = userOwnerName.String
		g.UserResourceOwnerDomain = userOwnerDomain.String
		g.FirstName = firstName.String
		g.LastName = lastName.String
		g.Email = email.String
		g.DisplayName = displayName.String
		g.AvatarURL = avatarURL.String
		g.PreferredLoginName = preferredLoginName.String
		g.OrgName = orgName.String
		g.OrgPrimaryDomain = orgDomain.String
		g.ProjectName = projectName.String
		g.ProjectResourceOwner = projectResourceOwner.String
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

var rolesInfoQuery = fmt.Sprintf(
	`SELECT JSON_AGG(
                 JSON_BUILD_OBJECT(
                     'role_key', pr.role_key,
                     'display_name', pr.display_name,
                     'group_name', pr.group_name
                 )
             ) as role_information
			 FROM %[1]s pr
			 WHERE pr.instance_id = %[2]s
			   AND pr.project_id = %[3]s
			   AND pr.role_key = ANY(%[4]s)`,
	projectRolesTable.identifier(),
	UserGrantInstanceID.identifier(),
	UserGrantProjectID.identifier(),
	UserGrantRoles.identifier(),
)

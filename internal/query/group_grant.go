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

var (
	groupGrantTable = table{
		name:          projection.GroupGrantProjectionTable,
		instanceIDCol: projection.GroupGrantInstanceID,
	}
	GroupGrantID = Column{
		name:  projection.GroupGrantID,
		table: groupGrantTable,
	}
	GroupGrantResourceOwner = Column{
		name:  projection.GroupGrantResourceOwner,
		table: groupGrantTable,
	}
	GroupGrantInstanceID = Column{
		name:  projection.GroupGrantInstanceID,
		table: groupGrantTable,
	}
	GroupGrantCreationDate = Column{
		name:  projection.GroupGrantCreationDate,
		table: groupGrantTable,
	}
	GroupGrantChangeDate = Column{
		name:  projection.GroupGrantChangeDate,
		table: groupGrantTable,
	}
	GroupGrantSequence = Column{
		name:  projection.GroupGrantSequence,
		table: groupGrantTable,
	}
	GroupGrantGroupID = Column{
		name:  projection.GroupGrantGroupID,
		table: groupGrantTable,
	}
	GroupGrantProjectID = Column{
		name:  projection.GroupGrantProjectID,
		table: groupGrantTable,
	}
	GroupGrantGrantID = Column{
		name:  projection.GroupGrantGrantID,
		table: groupGrantTable,
	}
	GroupGrantRoles = Column{
		name:  projection.GroupGrantRoles,
		table: groupGrantTable,
	}
	GroupGrantState = Column{
		name:  projection.GroupGrantState,
		table: groupGrantTable,
	}
	// GrantedOrgsTable = table{
	// 	name:          projection.OrgProjectionTable,
	// 	alias:         "granted_orgs",
	// 	instanceIDCol: projection.OrgColumnInstanceID,
	// }
	// GrantedOrgColumnId = Column{
	// 	name:  projection.OrgColumnID,
	// 	table: GrantedOrgsTable,
	// }
	// GrantedOrgColumnName = Column{
	// 	name:  projection.OrgColumnName,
	// 	table: GrantedOrgsTable,
	// }
	// GrantedOrgColumnDomain = Column{
	// 	name:  projection.OrgColumnDomain,
	// 	table: GrantedOrgsTable,
	// }
)

type GroupGrant struct {
	// ID represents the aggregate id (id of the group grant)
	ID           string                     `json:"id,omitempty"`
	CreationDate time.Time                  `json:"creation_date,omitempty"`
	ChangeDate   time.Time                  `json:"change_date,omitempty"`
	Sequence     uint64                     `json:"sequence,omitempty"`
	Roles        database.TextArray[string] `json:"roles,omitempty"`
	// GrantID represents the project grant id
	GrantID            string                 `json:"grant_id,omitempty"`
	State              domain.GroupGrantState `json:"state,omitempty"`
	GroupResourceOwner string                 `json:"group_resource_owner,omitempty"`

	// UserID             string          `json:"user_id,omitempty"`
	// Username           string          `json:"username,omitempty"`
	// UserType           domain.UserType `json:"user_type,omitempty"`
	// FirstName          string          `json:"first_name,omitempty"`
	// UserResourceOwner  string          `json:"user_resource_owner,omitempty"`
	// LastName           string          `json:"last_name,omitempty"`
	// Email              string          `json:"email,omitempty"`
	// DisplayName        string          `json:"display_name,omitempty"`
	// AvatarURL          string          `json:"avatar_url,omitempty"`
	// PreferredLoginName string          `json:"preferred_login_name,omitempty"`

	ResourceOwner    string `json:"resource_owner,omitempty"`
	OrgName          string `json:"org_name,omitempty"`
	OrgPrimaryDomain string `json:"org_primary_domain,omitempty"`

	GroupID          string `json:"group_id,omitempty"`
	GroupName        string `json:"group_name,omitempty"`
	GroupDescription string `json:"group_description,omitempty"`

	ProjectID   string `json:"project_id,omitempty"`
	ProjectName string `json:"project_name,omitempty"`

	GrantedOrgID     string `json:"granted_org_id,omitempty"`
	GrantedOrgName   string `json:"granted_org_name,omitempty"`
	GrantedOrgDomain string `json:"granted_org_domain,omitempty"`
}

type GroupGrants struct {
	SearchResponse
	GroupGrants []*GroupGrant
}

type GroupGrantsQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *GroupGrantsQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func NewGroupGrantGroupIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(GroupGrantGroupID, id, TextEquals)
}

func NewGroupGrantGroupNameQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(GroupColumnName, value, method)
}

func NewGroupGrantProjectIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(GroupGrantProjectID, id, TextEquals)
}

func NewGroupGrantProjectIDsSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(GroupGrantProjectID, list, ListIn)
}

func NewGroupGrantProjectOwnerSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(GroupColumnResourceOwner, id, TextEquals)
}

func NewGroupGrantResourceOwnerSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(GroupGrantResourceOwner, id, TextEquals)
}

func NewGroupGrantGrantIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(GroupGrantGrantID, id, TextEquals)
}

func NewGroupGrantIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(GroupGrantID, id, TextEquals)
}

// func NewGroupGrantUserTypeQuery(typ domain.UserType) (SearchQuery, error) {
// 	return NewNumberQuery(UserTypeCol, typ, NumberEquals)
// }

// func NewGroupGrantDisplayNameQuery(displayName string, method TextComparison) (SearchQuery, error) {
// 	return NewTextQuery(HumanDisplayNameCol, displayName, method)
// }

// func NewGroupGrantEmailQuery(email string, method TextComparison) (SearchQuery, error) {
// 	return NewTextQuery(HumanEmailCol, email, method)
// }

// func NewGroupGrantFirstNameQuery(value string, method TextComparison) (SearchQuery, error) {
// 	return NewTextQuery(HumanFirstNameCol, value, method)
// }

// func NewGroupGrantLastNameQuery(value string, method TextComparison) (SearchQuery, error) {
// 	return NewTextQuery(HumanLastNameCol, value, method)
// }

// func NewGroupGrantUsernameQuery(value string, method TextComparison) (SearchQuery, error) {
// 	return NewTextQuery(UserUsernameCol, value, method)
// }

func NewGroupGrantDomainQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(OrgColumnDomain, value, method)
}

func NewGroupGrantOrgNameQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(OrgColumnName, value, method)
}

func NewGroupGrantProjectNameQuery(value string, method TextComparison) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnName, value, method)
}

func NewGroupGrantRoleQuery(value string) (SearchQuery, error) {
	return NewTextQuery(GroupGrantRoles, value, TextListContains)
}

func NewGroupGrantStateQuery(value domain.GroupGrantState) (SearchQuery, error) {
	return NewNumberQuery(GroupGrantState, value, NumberEquals)
}

func NewGroupGrantWithGrantedQuery(owner string) (SearchQuery, error) {
	orgQuery, err := NewGroupGrantResourceOwnerSearchQuery(owner)
	if err != nil {
		return nil, err
	}
	projectQuery, err := NewGroupGrantProjectOwnerSearchQuery(owner)
	if err != nil {
		return nil, err
	}
	return NewOrQuery(orgQuery, projectQuery)
}

func NewGroupGrantContainsRolesSearchQuery(roles ...string) (SearchQuery, error) {
	r := make([]interface{}, len(roles))
	for i, role := range roles {
		r[i] = role
	}
	return NewListQuery(GroupGrantRoles, r, ListIn)
}

func (q *Queries) GroupGrant(ctx context.Context, shouldTriggerBulk bool, queries ...SearchQuery) (grant *GroupGrant, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerGroupGrantProjection")
		ctx, err = projection.GroupGrantProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareGroupGrantQuery(ctx, q.client)
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{GroupGrantInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Eb2KW", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		grant, err = scan(row)
		return err
	}, stmt, args...)
	return grant, err
}

func (q *Queries) GroupGrants(ctx context.Context, queries *GroupGrantsQueries, shouldTriggerBulk bool) (grants *GroupGrants, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerGroupGrantProjection")
		ctx, err = projection.GroupGrantProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("unable to trigger")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareGroupGrantsQuery(ctx, q.client)
	eq := sq.Eq{GroupGrantInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-xYmQR", "Errors.Query.SQLStatement")
	}

	latestSequence, err := q.latestState(ctx, groupGrantTable)
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

	grants.State = latestSequence
	return grants, nil
}

func prepareGroupGrantQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*GroupGrant, error)) {
	return sq.Select(
			GroupGrantID.identifier(),
			GroupGrantCreationDate.identifier(),
			GroupGrantChangeDate.identifier(),
			GroupGrantSequence.identifier(),
			GroupGrantGrantID.identifier(),
			GroupGrantRoles.identifier(),
			GroupGrantState.identifier(),
			GroupGrantGroupID.identifier(),
			GroupColumnName.identifier(),
			GroupColumnDescription.identifier(),
			GroupColumnResourceOwner.identifier(),

			GroupGrantResourceOwner.identifier(),
			OrgColumnName.identifier(),
			OrgColumnDomain.identifier(),

			GroupGrantProjectID.identifier(),
			ProjectColumnName.identifier(),

			GrantedOrgColumnId.identifier(),
			GrantedOrgColumnName.identifier(),
			GrantedOrgColumnDomain.identifier(),
		).
			From(groupGrantTable.identifier()).
			LeftJoin(join(GroupColumnID, GroupGrantGroupID)).
			LeftJoin(join(OrgColumnID, GroupGrantResourceOwner)).
			LeftJoin(join(ProjectColumnID, GroupGrantProjectID)).
			LeftJoin(join(GrantedOrgColumnId, GroupColumnResourceOwner) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*GroupGrant, error) {
			g := new(GroupGrant)

			var (
				groupName        sql.NullString
				groupDescription sql.NullString

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

				&g.GroupID,
				&groupName,
				&groupDescription,
				&g.GroupResourceOwner,

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
					return nil, zerrors.ThrowNotFound(err, "QUERY-wIPkA", "Errors.GroupGrant.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-oQPcP", "Errors.Internal")
			}

			g.GroupName = groupName.String
			g.GroupDescription = groupDescription.String
			g.OrgName = orgName.String
			g.OrgPrimaryDomain = orgDomain.String
			g.ProjectName = projectName.String
			g.GrantedOrgID = grantedOrgID.String
			g.GrantedOrgName = grantedOrgName.String
			g.GrantedOrgDomain = grantedOrgDomain.String
			return g, nil
		}
}

func prepareGroupGrantsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*GroupGrants, error)) {
	return sq.Select(
			GroupGrantID.identifier(),
			GroupGrantCreationDate.identifier(),
			GroupGrantChangeDate.identifier(),
			GroupGrantSequence.identifier(),
			GroupGrantGrantID.identifier(),
			GroupGrantRoles.identifier(),
			GroupGrantState.identifier(),
			GroupGrantGroupID.identifier(),
			GroupColumnName.identifier(),
			GroupColumnDescription.identifier(),
			GroupColumnResourceOwner.identifier(),

			GroupGrantResourceOwner.identifier(),
			OrgColumnName.identifier(),
			OrgColumnDomain.identifier(),

			GroupGrantProjectID.identifier(),
			ProjectColumnName.identifier(),

			GrantedOrgColumnId.identifier(),
			GrantedOrgColumnName.identifier(),
			GrantedOrgColumnDomain.identifier(),

			countColumn.identifier(),
		).
			From(groupGrantTable.identifier()).
			LeftJoin(join(GroupColumnID, GroupGrantGroupID)).
			LeftJoin(join(OrgColumnID, GroupGrantResourceOwner)).
			LeftJoin(join(ProjectColumnID, GroupGrantProjectID)).
			LeftJoin(join(GrantedOrgColumnId, GroupColumnResourceOwner) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*GroupGrants, error) {
			groupGrants := make([]*GroupGrant, 0)
			var count uint64
			for rows.Next() {
				g := new(GroupGrant)

				var (
					groupName        sql.NullString
					groupDescription sql.NullString

					orgName   sql.NullString
					orgDomain sql.NullString

					projectName sql.NullString

					grantedOrgID     sql.NullString
					grantedOrgName   sql.NullString
					grantedOrgDomain sql.NullString
				)

				err := rows.Scan(
					&g.ID,
					&g.CreationDate,
					&g.ChangeDate,
					&g.Sequence,
					&g.GrantID,
					&g.Roles,
					&g.State,

					&g.GroupID,
					&groupName,
					&groupDescription,
					&g.GroupResourceOwner,

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
				g.GroupName = groupName.String
				g.GroupDescription = groupDescription.String
				g.OrgName = orgName.String
				g.OrgPrimaryDomain = orgDomain.String
				g.ProjectName = projectName.String
				g.GrantedOrgID = grantedOrgID.String
				g.GrantedOrgName = grantedOrgName.String
				g.GrantedOrgDomain = grantedOrgDomain.String

				groupGrants = append(groupGrants, g)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-JHumP", "Errors.Query.CloseRows")
			}

			return &GroupGrants{
				GroupGrants: groupGrants,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

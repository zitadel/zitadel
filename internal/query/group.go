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
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	groupsTable = table{
		name:          projection.GroupProjectionTable,
		instanceIDCol: projection.GroupColumnInstanceID,
	}
	GroupColumnID = Column{
		name:  projection.GroupColumnID,
		table: groupsTable,
	}
	GroupColumnName = Column{
		name:           projection.GroupColumnName,
		table:          groupsTable,
		isOrderByLower: true,
	}
	GroupColumnDescription = Column{
		name:           projection.GroupColumnDescription,
		table:          groupsTable,
		isOrderByLower: true,
	}
	GroupColumnCreationDate = Column{
		name:  projection.GroupColumnCreationDate,
		table: groupsTable,
	}
	GroupColumnChangeDate = Column{
		name:  projection.GroupColumnChangeDate,
		table: groupsTable,
	}
	GroupColumnResourceOwner = Column{
		name:  projection.GroupColumnResourceOwner,
		table: groupsTable,
	}
	GroupColumnInstanceID = Column{
		name:  projection.GroupColumnInstanceID,
		table: groupsTable,
	}
	GroupColumnSequence = Column{
		name:  projection.GroupColumnSequence,
		table: groupsTable,
	}
	GroupColumnState = Column{
		name:  projection.GroupColumnState,
		table: groupsTable,
	}
)

type Groups struct {
	SearchResponse
	Groups []*Group
}

type Group struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.GroupState
	Sequence      uint64

	Name        string
	Description string

	UserID string
}

type GroupSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) GroupByID(ctx context.Context, shouldTriggerBulk bool, id string) (group *Group, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerGroupProjection")
		ctx, err = projection.GroupProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	stmt, scan := prepareGroupQuery(ctx, q.client)
	eq := sq.Eq{
		GroupColumnID.identifier():         id,
		GroupColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-2m00Q", "Errors.Query.SQLStatment")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		group, err = scan(row)
		return err
	}, query, args...)
	return group, err
}

func (q *Queries) SearchGroups(ctx context.Context, queries *GroupSearchQueries) (groups *Groups, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareGroupsQuery(ctx, q.client)
	eq := sq.Eq{GroupColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-em9ew", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		groups, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-3k11f", "Errors.Internal")
	}
	groups.State, err = q.latestState(ctx, groupsTable)
	return groups, err
}

func NewGroupNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(GroupColumnName, value, method)
}

func NewGroupIDSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(GroupColumnID, list, ListIn)
}

func NewGroupResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(GroupColumnResourceOwner, value, TextEquals)
}

func (r *GroupSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewGroupResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *GroupSearchQueries) AppendPermissionQueries(permissions []string) error {
	if !authz.HasGlobalPermission(permissions) {
		ids := authz.GetAllPermissionCtxIDs(permissions)
		query, err := NewGroupIDSearchQuery(ids)
		if err != nil {
			return err
		}
		r.Queries = append(r.Queries, query)
	}
	return nil
}

func (q *GroupSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareGroupQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Group, error)) {
	return sq.Select(
			GroupColumnID.identifier(),
			GroupColumnCreationDate.identifier(),
			GroupColumnChangeDate.identifier(),
			GroupColumnResourceOwner.identifier(),
			GroupColumnState.identifier(),
			GroupColumnSequence.identifier(),
			GroupColumnName.identifier(),
			GroupColumnDescription.identifier()).
			From(groupsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Group, error) {
			g := new(Group)
			err := row.Scan(
				&g.ID,
				&g.CreationDate,
				&g.ChangeDate,
				&g.ResourceOwner,
				&g.State,
				&g.Sequence,
				&g.Name,
				&g.Description,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-el3fs", "Errors.Group.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-ek3EE", "Errors.Internal")
			}
			return g, nil
		}
}

func prepareGroupsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Groups, error)) {
	return sq.Select(
			GroupColumnID.identifier(),
			GroupColumnCreationDate.identifier(),
			GroupColumnChangeDate.identifier(),
			GroupColumnResourceOwner.identifier(),
			GroupColumnState.identifier(),
			GroupColumnSequence.identifier(),
			GroupColumnName.identifier(),
			GroupColumnDescription.identifier(),
			countColumn.identifier()).
			From(groupsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Groups, error) {
			groups := make([]*Group, 0)
			var count uint64
			for rows.Next() {
				group := new(Group)
				err := rows.Scan(
					&group.ID,
					&group.CreationDate,
					&group.ChangeDate,
					&group.ResourceOwner,
					&group.State,
					&group.Sequence,
					&group.Name,
					&group.Description,
					&count,
				)
				if err != nil {
					return nil, err
				}
				groups = append(groups, group)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-PNXJv", "Errors.Query.CloseRows")
			}

			return &Groups{
				Groups: groups,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func (q *Queries) GroupByUserID(ctx context.Context, shouldTriggerBulk bool, id string) (groups *Groups, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerGroupProjection")
		ctx, err = projection.GroupProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	stmt, scan := prepareUserGroupsQuery(ctx, q.client)
	eq := sq.Eq{
		GroupMemberUserID.identifier():     id,
		GroupColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-2n10Q", "Errors.Query.SQLStatment")
	}
	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		groups, err = scan(rows)
		return err
	}, query, args...)
	if err != nil {
		return nil, err
	}
	groups.State, err = q.latestState(ctx, groupsTable)
	return groups, err
}

func prepareUserGroupsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Groups, error)) {
	return sq.Select(
			GroupColumnID.identifier(),
			GroupColumnCreationDate.identifier(),
			GroupColumnChangeDate.identifier(),
			GroupColumnResourceOwner.identifier(),
			GroupColumnState.identifier(),
			GroupColumnSequence.identifier(),
			GroupColumnName.identifier(),
			GroupColumnDescription.identifier(),

			GroupMemberUserID.identifier(),
			countColumn.identifier()).
			From(groupsTable.identifier()).
			LeftJoin(join(GroupMemberGroupID, GroupColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Groups, error) {
			groups := make([]*Group, 0)
			var count uint64

			for rows.Next() {
				group := new(Group)
				var (
					userId = sql.NullString{}
				)
				err := rows.Scan(
					&group.ID,
					&group.CreationDate,
					&group.ChangeDate,
					&group.ResourceOwner,
					&group.State,
					&group.Sequence,
					&group.Name,
					&group.Description,
					&userId,
					&count,
				)

				if err != nil {
					return nil, err
				}

				group.UserID = userId.String
				groups = append(groups, group)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-QMcJv", "Errors.Query.CloseRows")
			}

			return &Groups{
				Groups: groups,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
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
		name:  projection.GroupColumnName,
		table: groupsTable,
	}
	GroupColumnDescription = Column{
		name:  projection.GroupColumnDescription,
		table: groupsTable,
	}
	GroupColumnResourceOwner = Column{
		name:  projection.GroupColumnResourceOwner,
		table: groupsTable,
	}
	GroupColumnCreationDate = Column{
		name:  projection.GroupColumnCreationDate,
		table: groupsTable,
	}
	GroupColumnChangeDate = Column{
		name:  projection.GroupColumnChangeDate,
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
	Name          string
	Description   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	InstanceID    string
	State         domain.GroupState
	Sequence      uint64
}

type GroupSearchQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) GetGroupByID(ctx context.Context, id string) (group *Group, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// todo: add permission check

	stmt, scan := prepareGroupQuery()
	eq := sq.Eq{
		GroupColumnID.identifier():         id,
		GroupColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-8bde1", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		group, err = scan(row)
		return err
	}, query, args...)
	return group, err

}

// SearchGroups returns the list of groups that match the search criteria
func (q *Queries) SearchGroups(ctx context.Context, queries *GroupSearchQuery) (_ *Groups, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// todo: add permission check

	groups, err := q.searchGroups(ctx, queries)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func NewGroupNameSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(GroupColumnName, value, comparison)
}

func NewGroupIDsSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(GroupColumnID, list, ListIn)
}

func NewGroupOrganizationIdSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(GroupColumnResourceOwner, id, TextEquals)
}

func prepareGroupQuery() (sq.SelectBuilder, func(*sql.Row) (*Group, error)) {
	return sq.Select(
			GroupColumnID.identifier(),
			GroupColumnName.identifier(),
			GroupColumnDescription.identifier(),
			GroupColumnCreationDate.identifier(),
			GroupColumnChangeDate.identifier(),
			GroupColumnResourceOwner.identifier(),
			GroupColumnInstanceID.identifier(),
			GroupColumnSequence.identifier(),
			GroupColumnState.identifier()).
			From(groupsTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Group, error) {
			group := new(Group)
			err := row.Scan(
				&group.ID,
				&group.Name,
				&group.Description,
				&group.CreationDate,
				&group.ChangeDate,
				&group.ResourceOwner,
				&group.InstanceID,
				&group.Sequence,
				&group.State,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-SG4WbR", "Errors.Group.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-6yHJEz", "Errors.Internal")
			}
			return group, nil
		}
}

func (q *Queries) searchGroups(ctx context.Context, queries *GroupSearchQuery) (groups *Groups, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareGroupsQuery()
	eq := sq.And{
		sq.Eq{
			GroupColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		},
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-FpBnrv", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		groups, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-vnQf5N", "Errors.Internal")
	}
	groups.State, err = q.latestState(ctx, groupsTable)
	return groups, err
}

func prepareGroupsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Groups, error)) {
	return sq.Select(
			GroupColumnID.identifier(),
			GroupColumnName.identifier(),
			GroupColumnDescription.identifier(),
			GroupColumnCreationDate.identifier(),
			GroupColumnChangeDate.identifier(),
			GroupColumnResourceOwner.identifier(),
			GroupColumnInstanceID.identifier(),
			GroupColumnSequence.identifier(),
			GroupColumnState.identifier(),
			countColumn.identifier()).
			From(groupsTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Groups, error) {
			groups := make([]*Group, 0)
			var count uint64
			for rows.Next() {
				group := new(Group)
				err := rows.Scan(
					&group.ID,
					&group.Name,
					&group.Description,
					&group.CreationDate,
					&group.ChangeDate,
					&group.ResourceOwner,
					&group.InstanceID,
					&group.Sequence,
					&group.State,
					&count,
				)
				if err != nil {
					return nil, err
				}
				groups = append(groups, group)
			}
			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-ndNVod", "Errors.Query.CloseRows")
			}

			return &Groups{
				Groups: groups,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func (q *GroupSearchQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

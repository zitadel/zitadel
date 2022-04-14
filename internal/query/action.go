package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/api/authz"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

var (
	actionTable = table{
		name: projection.ActionTable,
	}
	ActionColumnID = Column{
		name:  projection.ActionIDCol,
		table: actionTable,
	}
	ActionColumnCreationDate = Column{
		name:  projection.ActionCreationDateCol,
		table: actionTable,
	}
	ActionColumnChangeDate = Column{
		name:  projection.ActionChangeDateCol,
		table: actionTable,
	}
	ActionColumnResourceOwner = Column{
		name:  projection.ActionResourceOwnerCol,
		table: actionTable,
	}
	ActionColumnInstanceID = Column{
		name:  projection.ActionInstanceIDCol,
		table: actionTable,
	}
	ActionColumnSequence = Column{
		name:  projection.ActionSequenceCol,
		table: actionTable,
	}
	ActionColumnState = Column{
		name:  projection.ActionStateCol,
		table: actionTable,
	}
	ActionColumnName = Column{
		name:  projection.ActionNameCol,
		table: actionTable,
	}
	ActionColumnScript = Column{
		name:  projection.ActionScriptCol,
		table: actionTable,
	}
	ActionColumnTimeout = Column{
		name:  projection.ActionTimeoutCol,
		table: actionTable,
	}
	ActionColumnAllowedToFail = Column{
		name:  projection.ActionAllowedToFailCol,
		table: actionTable,
	}
)

type Actions struct {
	SearchResponse
	Actions []*Action
}

type Action struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.ActionState
	Sequence      uint64

	Name          string
	Script        string
	Timeout       time.Duration
	AllowedToFail bool
}

type ActionSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *ActionSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) SearchActions(ctx context.Context, queries *ActionSearchQueries) (actions *Actions, err error) {
	query, scan := prepareActionsQuery()
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			ActionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).
		ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-SDgwg", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-SDfr52", "Errors.Internal")
	}
	actions, err = scan(rows)
	if err != nil {
		return nil, err
	}
	actions.LatestSequence, err = q.latestSequence(ctx, actionTable)
	return actions, err
}

func (q *Queries) GetActionByID(ctx context.Context, id string, orgID string) (*Action, error) {
	stmt, scan := prepareActionQuery()
	query, args, err := stmt.Where(
		sq.Eq{
			ActionColumnID.identifier():            id,
			ActionColumnResourceOwner.identifier(): orgID,
			ActionColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dgff3", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func NewActionResourceOwnerQuery(id string) (SearchQuery, error) {
	return NewTextQuery(ActionColumnResourceOwner, id, TextEquals)
}

func NewActionNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ActionColumnName, value, method)
}

func NewActionStateSearchQuery(value domain.ActionState) (SearchQuery, error) {
	return NewNumberQuery(ActionColumnState, int(value), NumberEquals)
}

func prepareActionsQuery() (sq.SelectBuilder, func(rows *sql.Rows) (*Actions, error)) {
	return sq.Select(
			ActionColumnID.identifier(),
			ActionColumnCreationDate.identifier(),
			ActionColumnChangeDate.identifier(),
			ActionColumnResourceOwner.identifier(),
			ActionColumnSequence.identifier(),
			ActionColumnState.identifier(),
			ActionColumnName.identifier(),
			ActionColumnScript.identifier(),
			ActionColumnTimeout.identifier(),
			ActionColumnAllowedToFail.identifier(),
			countColumn.identifier(),
		).From(actionTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Actions, error) {
			actions := make([]*Action, 0)
			var count uint64
			for rows.Next() {
				action := new(Action)
				err := rows.Scan(
					&action.ID,
					&action.CreationDate,
					&action.ChangeDate,
					&action.ResourceOwner,
					&action.Sequence,
					&action.State,
					&action.Name,
					&action.Script,
					&action.Timeout,
					&action.AllowedToFail,
					&count,
				)
				if err != nil {
					return nil, err
				}
				actions = append(actions, action)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-EGdff", "Errors.Query.CloseRows")
			}

			return &Actions{
				Actions: actions,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareActionQuery() (sq.SelectBuilder, func(row *sql.Row) (*Action, error)) {
	return sq.Select(
			ActionColumnID.identifier(),
			ActionColumnCreationDate.identifier(),
			ActionColumnChangeDate.identifier(),
			ActionColumnResourceOwner.identifier(),
			ActionColumnSequence.identifier(),
			ActionColumnState.identifier(),
			ActionColumnName.identifier(),
			ActionColumnScript.identifier(),
			ActionColumnTimeout.identifier(),
			ActionColumnAllowedToFail.identifier(),
		).From(actionTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Action, error) {
			action := new(Action)
			err := row.Scan(
				&action.ID,
				&action.CreationDate,
				&action.ChangeDate,
				&action.ResourceOwner,
				&action.Sequence,
				&action.State,
				&action.Name,
				&action.Script,
				&action.Timeout,
				&action.AllowedToFail,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-GEfnb", "Errors.Action.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-Dbnt4", "Errors.Internal")
			}
			return action, nil
		}
}

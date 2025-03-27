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

const (
	maxTimeout = 20 * time.Second
)

var (
	actionTable = table{
		name:          projection.ActionTable,
		instanceIDCol: projection.ActionInstanceIDCol,
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
	ActionColumnOwnerRemoved = Column{
		name:  projection.ActionOwnerRemovedCol,
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
	timeout       time.Duration
	AllowedToFail bool
}

func (a *Action) Timeout() time.Duration {
	if a.timeout > 0 && a.timeout < maxTimeout {
		return a.timeout
	}
	return maxTimeout
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

func (q *Queries) SearchActions(ctx context.Context, queries *ActionSearchQueries, withOwnerRemoved bool) (actions *Actions, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareActionsQuery()
	eq := sq.Eq{
		ActionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[ActionColumnOwnerRemoved.identifier()] = false
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-SDgwg", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		actions, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-SDfr52", "Errors.Internal")
	}

	actions.State, err = q.latestState(ctx, actionTable)
	return actions, err
}

func (q *Queries) GetActionByID(ctx context.Context, id string, orgID string, withOwnerRemoved bool) (action *Action, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareActionQuery()
	eq := sq.Eq{
		ActionColumnID.identifier():            id,
		ActionColumnResourceOwner.identifier(): orgID,
		ActionColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[ActionColumnOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dgff3", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		action, err = scan(row)
		return err
	}, query, args...)
	return action, err
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

func NewActionIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(ActionColumnID, id, TextEquals)
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
		).From(actionTable.identifier()).
			PlaceholderFormat(sq.Dollar),
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
					&action.timeout,
					&action.AllowedToFail,
					&count,
				)
				if err != nil {
					return nil, err
				}
				actions = append(actions, action)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-EGdff", "Errors.Query.CloseRows")
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
		).From(actionTable.identifier()).
			PlaceholderFormat(sq.Dollar),
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
				&action.timeout,
				&action.AllowedToFail,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-GEfnb", "Errors.Action.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-Dbnt4", "Errors.Internal")
			}
			return action, nil
		}
}

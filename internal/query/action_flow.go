package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

var (
	flowsTriggersTable = table{
		name: projection.FlowTriggerTable,
	}
	FlowsTriggersColumnFlowType = Column{
		name:  projection.FlowTypeCol,
		table: flowsTriggersTable,
	}
	FlowsTriggersColumnTriggerType = Column{
		name:  projection.FlowTriggerTypeCol,
		table: flowsTriggersTable,
	}
	FlowsTriggersColumnResourceOwner = Column{
		name:  projection.FlowResourceOwnerCol,
		table: flowsTriggersTable,
	}
	FlowsTriggersColumnTriggerSequence = Column{
		name:  projection.FlowActionTriggerSequenceCol,
		table: flowsTriggersTable,
	}
	FlowsTriggersColumnActionID = Column{
		name:  projection.FlowActionIDCol,
		table: flowsTriggersTable,
	}
)

type Flow struct {
	CreationDate  time.Time //TODO: add in projection
	ChangeDate    time.Time //TODO: add in projection
	ResourceOwner string    //TODO: add in projection
	Sequence      uint64    //TODO: add in projection
	Type          domain.FlowType

	TriggerActions map[domain.TriggerType][]*Action
}

func (q *Queries) GetFlow(ctx context.Context, flowType domain.FlowType, orgID string) (*Flow, error) {
	query, scan := prepareFlowQuery()
	stmt, args, err := query.Where(
		sq.Eq{
			FlowsTriggersColumnFlowType.identifier():      flowType,
			FlowsTriggersColumnResourceOwner.identifier(): orgID,
		}).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-HBRh3", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Gg42f", "Errors.Internal")
	}
	return scan(rows)
}

func (q *Queries) GetActionsByFlowAndTriggerType(ctx context.Context, flowType domain.FlowType, triggerType domain.TriggerType, orgID string) ([]*Action, error) {
	stmt, scan := prepareTriggerActionsQuery()
	query, args, err := stmt.Where(
		sq.Eq{
			FlowsTriggersColumnFlowType.identifier():      flowType,
			FlowsTriggersColumnTriggerType.identifier():   triggerType,
			FlowsTriggersColumnResourceOwner.identifier(): orgID,
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dgff3", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-SDf52", "Errors.Internal")
	}
	return scan(rows)
}

func (q *Queries) GetFlowTypesOfActionID(ctx context.Context, actionID string) ([]domain.FlowType, error) {
	stmt, scan := prepareFlowTypesQuery()
	query, args, err := stmt.Where(
		sq.Eq{
			FlowsTriggersColumnActionID.identifier(): actionID,
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-Dh311", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Bhj4w", "Errors.Internal")
	}

	return scan(rows)
}

func prepareFlowTypesQuery() (sq.SelectBuilder, func(*sql.Rows) ([]domain.FlowType, error)) {
	return sq.Select(
			FlowsTriggersColumnFlowType.identifier(),
		).
			From(flowsTriggersTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) ([]domain.FlowType, error) {
			types := []domain.FlowType{}
			for rows.Next() {
				var flowType domain.FlowType
				err := rows.Scan(
					&flowType,
				)
				if err != nil {
					return nil, err
				}
				types = append(types, flowType)
			}
			return types, nil
		}

}

func prepareTriggerActionsQuery() (sq.SelectBuilder, func(*sql.Rows) ([]*Action, error)) {
	return sq.Select(
			ActionColumnID.identifier(),
			ActionColumnCreationDate.identifier(),
			ActionColumnChangeDate.identifier(),
			ActionColumnResourceOwner.identifier(),
			//ActionColumnState.identifier(),
			ActionColumnSequence.identifier(),
			ActionColumnName.identifier(),
			ActionColumnScript.identifier(),
		).
			From(flowsTriggersTable.name).
			LeftJoin(join(ActionColumnID, FlowsTriggersColumnActionID)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) ([]*Action, error) {
			actions := make([]*Action, 0)
			for rows.Next() {
				action := new(Action)
				err := rows.Scan(
					&action.ID,
					&action.CreationDate,
					&action.ChangeDate,
					&action.ResourceOwner,
					//&action.State, //TODO: state in next release
					&action.Sequence,
					&action.Name,
					&action.Script,
				)
				if err != nil {
					return nil, err
				}
				actions = append(actions, action)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-Df42d", "Errors.Query.CloseRows")
			}

			return actions, nil
		}
}

func prepareFlowQuery() (sq.SelectBuilder, func(*sql.Rows) (*Flow, error)) {
	return sq.Select(
			ActionColumnID.identifier(),
			ActionColumnCreationDate.identifier(),
			ActionColumnChangeDate.identifier(),
			ActionColumnResourceOwner.identifier(),
			//ActionColumnState.identifier(),
			ActionColumnSequence.identifier(),
			ActionColumnName.identifier(),
			ActionColumnScript.identifier(),
			FlowsTriggersColumnTriggerType.identifier(),
			FlowsTriggersColumnTriggerSequence.identifier(),
			FlowsTriggersColumnFlowType.identifier(),
		).
			From(flowsTriggersTable.name).
			LeftJoin(join(ActionColumnID, FlowsTriggersColumnActionID)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Flow, error) {
			flow := &Flow{
				TriggerActions: make(map[domain.TriggerType][]*Action),
			}
			for rows.Next() {
				// action := new(Action)
				var (
					actionID            sql.NullString
					actionCreationDate  pq.NullTime
					actionChangeDate    pq.NullTime
					actionResourceOwner sql.NullString
					actionSequence      sql.NullInt64
					actionName          sql.NullString
					actionScript        sql.NullString

					triggerType     domain.TriggerType
					triggerSequence int
					flowType        domain.FlowType
				)
				err := rows.Scan(
					&actionID,
					&actionCreationDate,
					&actionChangeDate,
					&actionResourceOwner,
					//&action.State, //TODO: state in next release
					&actionSequence,
					&actionName,
					&actionScript,
					&triggerType,
					&triggerSequence,
					&flowType,
				)
				if err != nil {
					return nil, err
				}
				flow.Type = flowType
				if !actionID.Valid {
					continue
				}
				flow.TriggerActions[triggerType] = append(flow.TriggerActions[triggerType], &Action{
					ID:            actionID.String,
					CreationDate:  actionCreationDate.Time,
					ChangeDate:    actionChangeDate.Time,
					ResourceOwner: actionResourceOwner.String,
					Sequence:      uint64(actionSequence.Int64),
					Name:          actionName.String,
					Script:        actionScript.String,
				})
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-Dfbe2", "Errors.Query.CloseRows")
			}

			return flow, nil
		}
}

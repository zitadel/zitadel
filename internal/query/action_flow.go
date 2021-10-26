package query

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"

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
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	Type          domain.FlowType

	TriggerActions map[domain.TriggerType][]*Action
}

func (q *Queries) GetFlow(ctx context.Context, flowType domain.FlowType, orgID string) (*Flow, error) {
	query, scan := q.prepareFlowQuery(flowType)
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
	stmt, scan := q.prepareTriggerActionsQuery()
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
	stmt, args, err := squirrel.StatementBuilder.
		Select(FlowsTriggersColumnFlowType.identifier()).
		From(flowsTriggersTable.identifier()).
		Where(sq.Eq{
			FlowsTriggersColumnActionID.identifier(): actionID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-Dh311", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Bhj4w", "Errors.Internal")
	}
	flowTypes := make([]domain.FlowType, 0)
	for rows.Next() {
		var flowType domain.FlowType
		err := rows.Scan(
			&flowType,
		)
		if err != nil {
			return nil, err
		}
		flowTypes = append(flowTypes, flowType)
	}

	if err := rows.Close(); err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Fbgnh", "Errors.Query.CloseRows")
	}

	return flowTypes, nil
}

func (q *Queries) prepareTriggerActionsQuery() (sq.SelectBuilder, func(*sql.Rows) ([]*Action, error)) {
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
		).
			From(flowsTriggersTable.name).
			LeftJoin(join(ActionColumnID, FlowsTriggersColumnActionID)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) ([]*Action, error) {
			actions := make([]*Action, 0)
			for rows.Next() {
				action := new(Action)
				var triggerType domain.TriggerType
				var triggerSequence int
				err := rows.Scan(
					&action.ID,
					&action.CreationDate,
					&action.ChangeDate,
					&action.ResourceOwner,
					//&action.State, //TODO: state in next release
					&action.Sequence,
					&action.Name,
					&action.Script,
					&triggerType,
					&triggerSequence,
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

func (q *Queries) prepareFlowQuery(flowType domain.FlowType) (sq.SelectBuilder, func(*sql.Rows) (*Flow, error)) {
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
		).
			From(flowsTriggersTable.name).
			LeftJoin(join(ActionColumnID, FlowsTriggersColumnActionID)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Flow, error) {
			flow := &Flow{
				Type:           flowType,
				TriggerActions: make(map[domain.TriggerType][]*Action),
			}
			for rows.Next() {
				action := new(Action)
				var triggerType domain.TriggerType
				var triggerSequence int
				err := rows.Scan(
					&action.ID,
					&action.CreationDate,
					&action.ChangeDate,
					&action.ResourceOwner,
					//&action.State, //TODO: state in next release
					&action.Sequence,
					&action.Name,
					&action.Script,
					&triggerType,
					&triggerSequence,
				)
				if err != nil {
					return nil, err
				}
				flow.TriggerActions[triggerType] = append(flow.TriggerActions[triggerType], action)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-Dfbe2", "Errors.Query.CloseRows")
			}

			return flow, nil
		}
}

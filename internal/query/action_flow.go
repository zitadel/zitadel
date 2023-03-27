package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	flowsTriggersTable = table{
		name:          projection.FlowTriggerTable,
		instanceIDCol: projection.FlowInstanceIDCol,
	}
	FlowsTriggersColumnFlowType = Column{
		name:  projection.FlowTypeCol,
		table: flowsTriggersTable,
	}
	FlowsTriggersColumnChangeDate = Column{
		name:  projection.FlowChangeDateCol,
		table: flowsTriggersTable,
	}
	FlowsTriggersColumnSequence = Column{
		name:  projection.FlowSequenceCol,
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
	FlowsTriggersColumnInstanceID = Column{
		name:  projection.FlowInstanceIDCol,
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
	FlowsTriggersOwnerRemovedCol = Column{
		name:  projection.FlowOwnerRemovedCol,
		table: flowsTriggersTable,
	}
)

type Flow struct {
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	Type          domain.FlowType

	TriggerActions map[domain.TriggerType][]*Action
}

func (q *Queries) GetFlow(ctx context.Context, flowType domain.FlowType, orgID string, withOwnerRemoved bool) (_ *Flow, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareFlowQuery(ctx, q.client, flowType)
	eq := sq.Eq{
		FlowsTriggersColumnFlowType.identifier():      flowType,
		FlowsTriggersColumnResourceOwner.identifier(): orgID,
		FlowsTriggersColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[FlowsTriggersOwnerRemovedCol.identifier()] = false
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-HBRh3", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Gg42f", "Errors.Internal")
	}
	return scan(rows)
}

func (q *Queries) GetActiveActionsByFlowAndTriggerType(ctx context.Context, flowType domain.FlowType, triggerType domain.TriggerType, orgID string, withOwnerRemoved bool) (_ []*Action, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareTriggerActionsQuery(ctx, q.client)
	eq := sq.Eq{
		FlowsTriggersColumnFlowType.identifier():      flowType,
		FlowsTriggersColumnTriggerType.identifier():   triggerType,
		FlowsTriggersColumnResourceOwner.identifier(): orgID,
		FlowsTriggersColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
		ActionColumnState.identifier():                domain.ActionStateActive,
	}
	if !withOwnerRemoved {
		eq[FlowsTriggersOwnerRemovedCol.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dgff3", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-SDf52", "Errors.Internal")
	}
	return scan(rows)
}

func (q *Queries) GetFlowTypesOfActionID(ctx context.Context, actionID string, withOwnerRemoved bool) (_ []domain.FlowType, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareFlowTypesQuery(ctx, q.client)
	eq := sq.Eq{
		FlowsTriggersColumnActionID.identifier():   actionID,
		FlowsTriggersColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[FlowsTriggersOwnerRemovedCol.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-Dh311", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Bhj4w", "Errors.Internal")
	}

	return scan(rows)
}

func prepareFlowTypesQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) ([]domain.FlowType, error)) {
	return sq.Select(
			FlowsTriggersColumnFlowType.identifier(),
		).
			From(flowsTriggersTable.identifier() + db.Timetravel(call.Took(ctx))).
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

func prepareTriggerActionsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) ([]*Action, error)) {
	return sq.Select(
			ActionColumnID.identifier(),
			ActionColumnCreationDate.identifier(),
			ActionColumnChangeDate.identifier(),
			ActionColumnResourceOwner.identifier(),
			ActionColumnState.identifier(),
			ActionColumnSequence.identifier(),
			ActionColumnName.identifier(),
			ActionColumnScript.identifier(),
			ActionColumnAllowedToFail.identifier(),
			ActionColumnTimeout.identifier(),
		).
			From(flowsTriggersTable.name).
			LeftJoin(join(ActionColumnID, FlowsTriggersColumnActionID) + db.Timetravel(call.Took(ctx))).
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
					&action.State,
					&action.Sequence,
					&action.Name,
					&action.Script,
					&action.AllowedToFail,
					&action.timeout,
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

func prepareFlowQuery(ctx context.Context, db prepareDatabase, flowType domain.FlowType) (sq.SelectBuilder, func(*sql.Rows) (*Flow, error)) {
	return sq.Select(
			ActionColumnID.identifier(),
			ActionColumnCreationDate.identifier(),
			ActionColumnChangeDate.identifier(),
			ActionColumnResourceOwner.identifier(),
			ActionColumnState.identifier(),
			ActionColumnSequence.identifier(),
			ActionColumnName.identifier(),
			ActionColumnScript.identifier(),
			ActionColumnAllowedToFail.identifier(),
			ActionColumnTimeout.identifier(),
			FlowsTriggersColumnTriggerType.identifier(),
			FlowsTriggersColumnTriggerSequence.identifier(),
			FlowsTriggersColumnFlowType.identifier(),
			FlowsTriggersColumnChangeDate.identifier(),
			FlowsTriggersColumnSequence.identifier(),
			FlowsTriggersColumnResourceOwner.identifier(),
		).
			From(flowsTriggersTable.name).
			LeftJoin(join(ActionColumnID, FlowsTriggersColumnActionID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Flow, error) {
			flow := &Flow{
				TriggerActions: make(map[domain.TriggerType][]*Action),
				Type:           flowType,
			}
			for rows.Next() {
				var (
					actionID            sql.NullString
					actionCreationDate  sql.NullTime
					actionChangeDate    sql.NullTime
					actionResourceOwner sql.NullString
					actionState         sql.NullInt32
					actionSequence      sql.NullInt64
					actionName          sql.NullString
					actionScript        sql.NullString
					actionAllowedToFail sql.NullBool
					actionTimeout       sql.NullInt64

					triggerType     domain.TriggerType
					triggerSequence int
				)
				err := rows.Scan(
					&actionID,
					&actionCreationDate,
					&actionChangeDate,
					&actionResourceOwner,
					&actionState,
					&actionSequence,
					&actionName,
					&actionScript,
					&actionAllowedToFail,
					&actionTimeout,
					&triggerType,
					&triggerSequence,
					&flow.Type,
					&flow.ChangeDate,
					&flow.Sequence,
					&flow.ResourceOwner,
				)
				if err != nil {
					return nil, err
				}
				if !actionID.Valid {
					continue
				}
				flow.TriggerActions[triggerType] = append(flow.TriggerActions[triggerType], &Action{
					ID:            actionID.String,
					CreationDate:  actionCreationDate.Time,
					ChangeDate:    actionChangeDate.Time,
					ResourceOwner: actionResourceOwner.String,
					State:         domain.ActionState(actionState.Int32),
					Sequence:      uint64(actionSequence.Int64),
					Name:          actionName.String,
					Script:        actionScript.String,
					AllowedToFail: actionAllowedToFail.Bool,
					timeout:       time.Duration(actionTimeout.Int64),
				})
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-Dfbe2", "Errors.Query.CloseRows")
			}

			return flow, nil
		}
}

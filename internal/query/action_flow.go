package query

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

func (q *Queries) GetActionsByFlowAndTriggerType(ctx context.Context, flowType domain.FlowType, triggerType domain.TriggerType) ([]*Action, error) {
	flowTypeQuery, _ := NewTriggerActionFlowTypeSearchQuery(flowType)
	triggerTypeQuery, _ := NewTriggerActionTriggerTypeSearchQuery(triggerType)
	return q.SearchActionsFromFlow(ctx, &TriggerActionSearchQueries{Queries: []SearchQuery{flowTypeQuery, triggerTypeQuery}})
}

var triggerActionsQuery = squirrel.StatementBuilder.Select("creation_date", "change_date", "resource_owner", "sequence", "action_id", "name", "script", "trigger_type", "trigger_sequence").
	From("zitadel.projections.flows_actions_triggers").PlaceholderFormat(squirrel.Dollar)

func (q *Queries) SearchActionsFromFlow(ctx context.Context, query *TriggerActionSearchQueries) ([]*Action, error) {
	stmt, args, err := query.ToQuery(triggerActionsQuery).OrderBy("flow_type", "trigger_type", "trigger_sequence").ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-wQ3by", "Errors.orgs.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M6mYN", "Errors.orgs.internal")
	}

	actions := []*Action{}
	for rows.Next() {
		action := new(Action)
		var triggerType domain.TriggerType
		var triggerSequence int
		rows.Scan(
			&action.CreationDate,
			&action.ChangeDate,
			&action.ResourceOwner,
			&action.Sequence,
			//&action.State, //TODO: state in next release
			&action.ID,
			&action.Name,
			&action.Script,
			&triggerType,
			&triggerSequence,
		)
		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-pA0Wj", "Errors.actions.internal")
	}

	return actions, nil
}

func (q *Queries) GetFlow(ctx context.Context, flowType domain.FlowType) (*Flow, error) {
	flowTypeQuery, _ := NewTriggerActionFlowTypeSearchQuery(flowType)
	return q.SearchFlow(ctx, &TriggerActionSearchQueries{Queries: []SearchQuery{flowTypeQuery}})
}

func (q *Queries) SearchFlow(ctx context.Context, query *TriggerActionSearchQueries) (*Flow, error) {
	stmt, args, err := query.ToQuery(triggerActionsQuery.OrderBy("flow_type", "trigger_type", "trigger_sequence")).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-wQ3by", "Errors.orgs.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M6mYN", "Errors.orgs.internal")
	}

	flow := &Flow{
		TriggerActions: make(map[domain.TriggerType][]*Action),
	}
	for rows.Next() {
		action := new(Action)
		var triggerType domain.TriggerType
		var triggerSequence int
		rows.Scan(
			&action.CreationDate,
			&action.ChangeDate,
			&action.ResourceOwner,
			&action.Sequence,
			//&action.State, //TODO: state in next release
			&action.ID,
			&action.Name,
			&action.Script,
			&triggerType,
			&triggerSequence,
		)

		flow.TriggerActions[triggerType] = append(flow.TriggerActions[triggerType], action)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-pA0Wj", "Errors.actions.internal")
	}

	return flow, nil
}

func (q *Queries) GetFlowTypesOfActionID(ctx context.Context, actionID string) ([]domain.FlowType, error) {
	actionIDQuery, _ := NewTriggerActionActionIDSearchQuery(actionID)
	query := &TriggerActionSearchQueries{Queries: []SearchQuery{actionIDQuery}}
	stmt, args, err := query.ToQuery(
		squirrel.StatementBuilder.
			Select("flow_type").
			From("zitadel.projections.flows_actions_triggers").
			PlaceholderFormat(squirrel.Dollar)).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-wQ3by", "Errors.orgs.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M6mYN", "Errors.orgs.internal")
	}
	flowTypes := make([]domain.FlowType, 0)
	for rows.Next() {
		var flow_type domain.FlowType
		rows.Scan(
			&flow_type,
		)

		flowTypes = append(flowTypes, flow_type)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-pA0Wj", "Errors.actions.internal")
	}

	return flowTypes, nil
}

type Flow struct {
	ID            string          `col:"id"`
	CreationDate  time.Time       `col:"creation_date"`
	ChangeDate    time.Time       `col:"change_date"`
	ResourceOwner string          `col:"resource_owner"`
	Sequence      uint64          `col:"sequence"`
	Type          domain.FlowType `col:"flow_type"`

	TriggerActions map[domain.TriggerType][]*Action
}

type TriggerActionSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *TriggerActionSearchQueries) ToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	query = q.SearchRequest.ToQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

func NewTriggerActionTriggerTypeSearchQuery(value domain.TriggerType) (SearchQuery, error) {
	return NewIntQuery("trigger_type", int(value), IntEquals)
}

func NewTriggerActionFlowTypeSearchQuery(value domain.FlowType) (SearchQuery, error) {
	return NewIntQuery("flow_type", int(value), IntEquals)
}

func NewTriggerActionActionIDSearchQuery(actionID string) (SearchQuery, error) {
	return NewTextQuery("action_id", actionID, TextEquals)
}

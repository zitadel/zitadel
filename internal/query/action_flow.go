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

var triggerActionsQuery = squirrel.StatementBuilder.Select("creation_date", "change_date", "resource_owner", "sequence", "name", "script").
	From("zitadel.projections.flows_actions_triggers").PlaceholderFormat(squirrel.Dollar)

func (q *Queries) SearchActionsFromFlow(ctx context.Context, query *TriggerActionSearchQueries) ([]*Action, error) {
	stmt, args, err := query.ToQuery(triggerActionsQuery).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-wQ3by", "Errors.orgs.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M6mYN", "Errors.orgs.internal")
	}

	actions := []*Action{}
	for rows.Next() {
		org := new(Action)
		rows.Scan(
			&org.CreationDate,
			&org.ChangeDate,
			&org.ResourceOwner,
			//&org.State,
			&org.Sequence,
			&org.Name,
			&org.Script,
		)
		actions = append(actions, org)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-pA0Wj", "Errors.actions.internal")
	}

	return actions, nil
}

type TriggerAction struct {
	ID            string    `col:"id"`
	CreationDate  time.Time `col:"creation_date"`
	ChangeDate    time.Time `col:"change_date"`
	ResourceOwner string    `col:"resource_owner"`
	Sequence      uint64    `col:"sequence"`

	Name          string        `col:"name"`
	Script        string        `col:"script"`
	Timeout       time.Duration `col:"-"`
	AllowedToFail bool          `col:"-"`
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

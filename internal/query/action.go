package query

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/errors"
)

var actionsQuery = squirrel.StatementBuilder.Select("creation_date", "change_date", "resource_owner", "sequence", "name", "script").
	From("zitadel.projections.flows_actions").PlaceholderFormat(squirrel.Dollar)

func (q *Queries) SearchActions(ctx context.Context, query *ActionSearchQueries) ([]*Action, error) {
	stmt, args, err := query.ToQuery(actionsQuery).ToSql()
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

type Action struct {
	ID            string    `col:"id"`
	CreationDate  time.Time `col:"creation_date"`
	ChangeDate    time.Time `col:"change_date"`
	ResourceOwner string    `col:"resource_owner"`
	//State         domain.ActionState `col:"action_state"`
	Sequence uint64 `col:"sequence"`

	Name          string        `col:"name"`
	Script        string        `col:"script"`
	Timeout       time.Duration `col:"-"`
	AllowedToFail bool          `col:"-"`
}

type ActionSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *ActionSearchQueries) ToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	query = q.SearchRequest.ToQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

func NewActionNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery("name", value, method)
}

func newActionIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery("id", id, TextEquals)
}

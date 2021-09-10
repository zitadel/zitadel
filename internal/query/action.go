package query

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

var actionsQuery = squirrel.StatementBuilder.Select("creation_date", "change_date", "resource_owner", "sequence", "id", "name", "script", "timeout", "allowed_to_fail").
	From("zitadel.projections.actions").PlaceholderFormat(squirrel.Dollar)

func (q *Queries) GetAction(ctx context.Context, id string, orgID string) (*Action, error) {
	idQuery, _ := newActionIDSearchQuery(id)
	actions, err := q.SearchActions(ctx, &ActionSearchQueries{Queries: []SearchQuery{idQuery}})
	if err != nil {
		return nil, err
	}
	if len(actions) != 1 {

	}
	return actions[0], err
}

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
			&org.ID,
			&org.Name,
			&org.Script,
			&org.Timeout,
			&org.AllowedToFail,
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

func NewActionResourceOwnerQuery(id string) (SearchQuery, error) {
	return NewTextQuery("resource_owner", id, TextEquals)
}

func NewActionNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery("name", value, method)
}

func NewActionStateSearchQuery(value domain.ActionState) (SearchQuery, error) {
	return NewIntQuery("state", int(value), IntEquals)
}

func newActionIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery("id", id, TextEquals)
}

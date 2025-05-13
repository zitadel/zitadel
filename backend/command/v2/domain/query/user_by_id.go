package query

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/v2/domain"
	"github.com/zitadel/zitadel/backend/command/v2/pattern"
	"github.com/zitadel/zitadel/backend/command/v2/storage/database"
)

type UserByIDQuery struct {
	querier database.Querier
	UserID  string `json:"userId"`
	res     *domain.User
}

var _ pattern.Query[*domain.User] = (*UserByIDQuery)(nil)

// Name implements [pattern.Command].
func (q *UserByIDQuery) Name() string {
	return "user.v2.by_id"
}

// Execute implements [pattern.Command].
func (q *UserByIDQuery) Execute(ctx context.Context) error {
	var res *domain.User
	err := q.querier.QueryRow(ctx, "SELECT id, username, email FROM users WHERE id = $1", q.UserID).Scan(&res.ID, &res.Username, &res.Email.Address)
	if err != nil {
		return err
	}
	q.res = res
	return nil
}

// Result implements [pattern.Query].
func (q *UserByIDQuery) Result() *domain.User {
	return q.res
}

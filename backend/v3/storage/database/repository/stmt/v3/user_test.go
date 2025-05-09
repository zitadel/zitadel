package v3_test

import (
	"context"
	"testing"

	v3 "github.com/zitadel/zitadel/backend/v3/storage/database/repository/stmt/v3"
)

type user struct{}

func TestUser(t *testing.T) {
	query := v3.NewUserQuery()
	query.Where(
		v3.Or(
			v3.UserByID("123"),
			v3.UserByUsername("test", v3.TextOperatorStartsWithIgnoreCase),
		),
	)
	query.Limit(10)
	query.Offset(5)
	// query.OrderBy(

	query.Result(context.TODO(), nil)
}

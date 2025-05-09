package stmt_test

import (
	"context"
	"testing"

	"github.com/zitadel/zitadel/backend/v3/storage/database/repository/stmt"
)

func Test_Bla(t *testing.T) {
	stmt.User(nil).Where(
		stmt.Or(
			stmt.UserIDCondition("123"),
			stmt.UserIDCondition("123"),
			stmt.UserUsernameCondition(stmt.TextOperationEqualIgnoreCase, "test"),
		),
	).Limit(1).Offset(1).Get(context.Background())
}

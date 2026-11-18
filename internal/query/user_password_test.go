package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestHumanPasswordReadModel_Query_withoutPasswordChecks(t *testing.T) {
	query := NewHumanPasswordReadModel("user1", "org1").Query()

	require.Len(t, query.GetQueries(), 1)
	// password checks don't change the password or the user state, but grow with every login of the user
	for _, eventType := range []string{
		string(user.HumanPasswordCheckFailedType),
		string(user.HumanPasswordCheckSucceededType),
		string(user.UserV1PasswordCheckFailedType),
		string(user.UserV1PasswordCheckSucceededType),
	} {
		for _, queried := range query.GetQueries()[0].GetEventTypes() {
			assert.NotEqual(t, eventType, string(queried))
		}
	}
	assert.Equal(t, "org1", query.GetResourceOwner())
}

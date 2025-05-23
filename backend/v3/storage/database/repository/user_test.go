package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"go.uber.org/mock/gomock"
)

func TestQueryUser(t *testing.T) {
	t.Skip("tests are meant as examples and are not real tests")
	t.Run("User filters", func(t *testing.T) {
		client := dbmock.NewMockClient(gomock.NewController(t))

		user := repository.UserRepository(client)
		u, err := user.Get(context.Background(),
			database.WithCondition(
				database.And(
					database.Or(
						user.IDCondition("test"),
						user.IDCondition("2"),
					),
					user.UsernameCondition(database.TextOperationStartsWithIgnoreCase, "test"),
				),
			),
			database.WithOrderBy(user.CreatedAtColumn()),
		)

		assert.NoError(t, err)
		assert.NotNil(t, u)
	})

	t.Run("machine and human filters", func(t *testing.T) {
		client := dbmock.NewMockClient(gomock.NewController(t))

		user := repository.UserRepository(client)
		machine := user.Machine()
		human := user.Human()
		email, err := human.GetEmail(context.Background(), database.And(
			user.UsernameCondition(database.TextOperationStartsWithIgnoreCase, "test"),
			database.Or(
				machine.DescriptionCondition(database.TextOperationStartsWithIgnoreCase, "test"),
				human.EmailVerifiedCondition(true),
				database.IsNotNull(machine.DescriptionColumn()),
			),
		))

		assert.NoError(t, err)
		assert.NotNil(t, email)
	})
}

type dbInstruction string

func TestArg(t *testing.T) {
	var bla any = "asdf"
	instr, ok := bla.(dbInstruction)
	assert.False(t, ok)
	assert.Empty(t, instr)
	bla = dbInstruction("asdf")
	instr, ok = bla.(dbInstruction)
	assert.True(t, ok)
	assert.Equal(t, instr, dbInstruction("asdf"))
}

func TestWriteUser(t *testing.T) {
	t.Skip("tests are meant as examples and are not real tests")
	t.Run("update user", func(t *testing.T) {
		user := repository.UserRepository(nil)
		user.Human().Update(context.Background(), user.IDCondition("test"), user.SetUsername("test"))
	})
}

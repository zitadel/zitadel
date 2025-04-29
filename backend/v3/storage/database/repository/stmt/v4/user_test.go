package v4_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v4 "github.com/zitadel/zitadel/backend/v3/storage/database/repository/stmt/v4"
)

func TestQueryUser(t *testing.T) {
	t.Run("User filters", func(t *testing.T) {
		user := v4.UserRepository(nil)
		user.WithCondition(
			v4.And(
				v4.Or(
					user.IDCondition("test"),
					user.IDCondition("2"),
				),
				user.UsernameCondition(v4.TextOperatorStartsWithIgnoreCase, "test"),
			),
		).Get(context.Background())
	})

	t.Run("machine and human filters", func(t *testing.T) {
		user := v4.UserRepository(nil)
		machine := user.Machine()
		human := user.Human()
		user.WithCondition(
			v4.And(
				user.UsernameCondition(v4.TextOperatorStartsWithIgnoreCase, "test"),
				v4.Or(
					machine.DescriptionCondition(v4.TextOperatorStartsWithIgnoreCase, "test"),
					human.EmailAddressVerifiedCondition(true),
					v4.IsNotNull(machine.DescriptionColumn()),
				),
			),
		)
		human.GetEmail(context.Background())
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
	t.Run("update user", func(t *testing.T) {
		user := v4.UserRepository(nil)
		user.WithCondition(user.IDCondition("test")).Human().Update(
			context.Background(),
			user.SetUsername("test"),
		)

	})
}

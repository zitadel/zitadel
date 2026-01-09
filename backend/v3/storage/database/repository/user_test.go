package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func Test_user_Get(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	userRepo := repository.UserRepository()
	// humanRepo := repository.HumanUserRepository()
	// machineRepo := repository.MachineUserRepository()

	instanceID := createInstance(t, tx)

	orgID1 := createOrganization(t, tx, instanceID)
	orgID2 := createOrganization(t, tx, instanceID)

	human := &domain.User{
		ID:             gofakeit.UUID(),
		InstanceID:     instanceID,
		OrganizationID: orgID1,
		Username:       "human-user",
		State:          domain.UserStateActive,
		Human: &domain.HumanUser{
			FirstName:         "John",
			LastName:          "Doe",
			DisplayName:       "JohnD",
			Nickname:          "johnny",
			PreferredLanguage: language.English,
			Gender:            domain.HumanGenderMale,
			AvatarKey:         "https://my.avatar/key",
			Email: domain.HumanEmail{
				Address:    "john@doe.com",
				VerifiedAt: time.Now(),
			},
		},
	}

	machine := &domain.User{
		ID:             gofakeit.UUID(),
		InstanceID:     instanceID,
		OrganizationID: orgID2,
		Username:       "machine-user",
		State:          domain.UserStateActive,
		Machine: &domain.MachineUser{
			Name:        "My Machine",
			Description: "This is my machine user",
		},
	}

	err := userRepo.Create(t.Context(), tx, human)
	require.NoError(t, err)
	err = userRepo.Create(t.Context(), tx, machine)
	require.NoError(t, err)

	tests := []struct {
		name      string
		setup     func(t *testing.T, tx database.Transaction) error
		condition database.Condition
		wantErr   error
		expected  *domain.User
	}{
		{
			name:      "incomplete condition",
			condition: userRepo.IDCondition(human.ID),
			wantErr:   database.NewMissingConditionError(userRepo.InstanceIDColumn()),
		},
		// {
		// 	name: "not found",
		// },
		// {
		// 	name: "too many",
		// },
		// {
		// 	name: "get human",
		// },
		// {
		// 	name: "get machine",
		// },
		// {
		// 	name: "wrong type results in not found",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			t.Cleanup(func() {
				err := savepoint.Rollback(context.Background())
				if err != nil {
					t.Log("rollback savepoint failed", err)
				}
			})
			if tt.setup != nil {
				err := tt.setup(t, savepoint)
				require.NoError(t, err)
			}

			user, err := userRepo.Get(t.Context(), savepoint, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.expected, user)
		})
	}
}

package projection

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestUserRelationalProjection_Reducers(t *testing.T) {
	handler := &userRelationalProjection{}
	rawTx, tx := getTransactions(t)
	t.Cleanup(func() {
		require.NoError(t, rawTx.Rollback())
	})

	ctx := t.Context()

	instanceReop := repository.InstanceRepository()
	instanceID := gofakeit.UUID()
	err := instanceReop.Create(ctx, tx, &domain.Instance{
		ID:   instanceID,
		Name: "test-instance",
	})
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository()
	orgID := gofakeit.UUID()
	err = orgRepo.Create(ctx, tx, &domain.Organization{
		InstanceID: instanceID,
		ID:         orgID,
		Name:       "test-org",
		State:      domain.OrgStateActive,
	})
	require.NoError(t, err)

	userRepo := repository.UserRepository()
	userID := gofakeit.UUID()
	err = userRepo.Create(ctx, tx, &domain.User{
		InstanceID:     instanceID,
		OrganizationID: orgID,
		ID:             userID,
		Username:       "testuser@example.com",
		State:          domain.UserStateActive,
		Human: &domain.HumanUser{
			FirstName: "Test",
			LastName:  "User",
			Email: domain.HumanEmail{
				Address:    "testuser@example.com",
				VerifiedAt: time.Now(),
			},
		},
		Metadata: []*domain.Metadata{
			{
				Key:   "key1",
				Value: []byte("value1"),
			},
		},
	})
	require.NoError(t, err)

	// TODO: add tests for other reducers in user_relational.go

	// recovery code reducers
	t.Run("reduce user.human.mfa.recoverycode.added event", func(t *testing.T) {
		// add recovery codes
		recoveryCodes := []string{"code1", "code2", "code3"}
		recoveryCodesAddedEvent := user.NewHumanRecoveryCodesAddedEvent(
			ctx,
			&user.NewAggregate(userID, orgID).Aggregate,
			recoveryCodes,
			nil)

		// reduce the event
		res := callReduce(t, ctx, rawTx, handler, recoveryCodesAddedEvent)
		require.True(t, res)

		// assert that the recovery codes are stored
		gotUser, err := userRepo.Get(ctx, tx, database.WithCondition(
			database.And(
				userRepo.IDCondition(userID),
				userRepo.InstanceIDCondition(instanceID),
			),
		))
		require.NoError(t, err)
		assert.ElementsMatch(t, recoveryCodes, gotUser.Human.RecoveryCodes.Codes)
		assert.Zero(t, gotUser.Human.RecoveryCodes.LastSuccessfullyCheckedAt)
		assert.Zero(t, gotUser.Human.RecoveryCodes.FailedAttempts)
	})

	t.Run("reduce user.human.mfa.recoverycode.removed event", func(t *testing.T) {
		// add recovery codes
		_, err = userRepo.Update(ctx, tx,
			userRepo.PrimaryKeyCondition(instanceID, userID),
			userRepo.Human().AddRecoveryCodes([]string{"code1", "code2", "code3"}))
		require.NoError(t, err)

		// set remove recovery codes event
		recoveryCodesRemovedEvent := user.NewHumanRecoveryCodeRemovedEvent(
			ctx,
			&user.NewAggregate(userID, orgID).Aggregate,
			nil)

		// reduce the event
		res := callReduce(t, ctx, rawTx, handler, recoveryCodesRemovedEvent)
		require.True(t, res)

		// assert that the recovery codes are removed
		gotUser, err := userRepo.Get(ctx, tx, database.WithCondition(
			database.And(
				userRepo.IDCondition(userID),
				userRepo.InstanceIDCondition(instanceID),
			),
		))
		require.NoError(t, err)
		assert.Empty(t, gotUser.Human.RecoveryCodes.Codes)
		assert.Zero(t, gotUser.Human.RecoveryCodes.LastSuccessfullyCheckedAt)
		assert.Zero(t, gotUser.Human.RecoveryCodes.FailedAttempts)
	})

	t.Run("reduce user.human.mfa.recoverycode.check.succeeded event", func(t *testing.T) {
		// add recovery codes
		_, err = userRepo.Update(ctx, tx,
			userRepo.PrimaryKeyCondition(instanceID, userID),
			userRepo.Human().AddRecoveryCodes([]string{"code1", "code2", "code3"}))
		require.NoError(t, err)

		// set recovery code check succeeded event
		recoveryCodeCheckSucceededEvent := user.NewHumanRecoveryCodeCheckSucceededEvent(
			ctx,
			&user.NewAggregate(userID, orgID).Aggregate,
			"code1",
			nil)

		// reduce the event
		res := callReduce(t, ctx, rawTx, handler, recoveryCodeCheckSucceededEvent)
		require.True(t, res)

		// assert that the recovery code is removed and the last_successfully_checked_at timestamp is set
		gotUser, err := userRepo.Get(ctx, tx, database.WithCondition(
			database.And(
				userRepo.IDCondition(userID),
				userRepo.InstanceIDCondition(instanceID),
			),
		))
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"code2", "code3"}, gotUser.Human.RecoveryCodes.Codes)
		assert.NotZero(t, gotUser.Human.RecoveryCodes.LastSuccessfullyCheckedAt)
		assert.Zero(t, gotUser.Human.RecoveryCodes.FailedAttempts)
	})

	t.Run("reduce user.human.mfa.recoverycode.check.failed event", func(t *testing.T) {
		// add recovery codes
		_, err = userRepo.Update(ctx, tx,
			userRepo.PrimaryKeyCondition(instanceID, userID),
			userRepo.Human().AddRecoveryCodes([]string{"code1", "code2", "code3"}))
		require.NoError(t, err)

		// set recovery code failed event
		recoveryCodeCheckFailedEvent := user.NewHumanRecoveryCodeCheckFailedEvent(
			ctx,
			&user.NewAggregate(userID, orgID).Aggregate,
			nil,
		)

		// reduce the event
		res := callReduce(t, ctx, rawTx, handler, recoveryCodeCheckFailedEvent)
		require.True(t, res)

		// assert that the recovery code failed_attempts is incremented
		gotUser, err := userRepo.Get(ctx, tx, database.WithCondition(
			database.And(
				userRepo.IDCondition(userID),
				userRepo.InstanceIDCondition(instanceID),
			),
		))
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"code1", "code2", "code3"}, gotUser.Human.RecoveryCodes.Codes)
		assert.Zero(t, gotUser.Human.RecoveryCodes.LastSuccessfullyCheckedAt)
		assert.Equal(t, 1, gotUser.Human.RecoveryCodes.FailedAttempts)
	})
}

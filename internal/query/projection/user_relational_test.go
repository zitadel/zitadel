package projection

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestUserRelationalProjection_Reducers(t *testing.T) {
	handler := &relationalTablesProjection{}
	rawTx, tx := getTransactions(t)
	t.Cleanup(func() {
		require.NoError(t, rawTx.Rollback())
	})

	// create instance
	instanceRepo := repository.InstanceRepository()
	instanceID := gofakeit.UUID()
	orgID := gofakeit.UUID()
	err := instanceRepo.Create(t.Context(), tx, &domain.Instance{
		ID:           instanceID,
		Name:         "test-instance",
		DefaultOrgID: orgID,
	})
	require.NoError(t, err)

	// create org
	orgRepo := repository.OrganizationRepository()
	err = orgRepo.Create(t.Context(), tx, &domain.Organization{
		InstanceID: instanceID,
		ID:         orgID,
		Name:       "test-org",
		State:      domain.OrgStateActive,
	})
	require.NoError(t, err)

	userRepo := repository.UserRepository()

	// TODO: add tests for other reducers in user_relational.go

	// recovery code reducers
	t.Run("reduce user.human.mfa.recoverycode.added event", func(t *testing.T) {
		// create user
		existingUserID := gofakeit.UUID()
		existingUserAgg := createUser(t, tx, userRepo, instanceID, orgID, existingUserID)

		// set `user.human.mfa.recoverycode.added` event
		recoveryCodes := []string{"code1", "code2", "code3"}
		recoveryCodesAddedEvent := user.NewHumanRecoveryCodesAddedEvent(
			t.Context(),
			&existingUserAgg.Aggregate,
			recoveryCodes,
			nil)

		// reduce the event
		eventReduced := callReduce(t, rawTx, handler, recoveryCodesAddedEvent)
		require.True(t, eventReduced)

		// assert that the recovery codes are stored
		gotUser, err := userRepo.Get(t.Context(), tx, database.WithCondition(
			database.And(
				userRepo.IDCondition(existingUserID),
				userRepo.InstanceIDCondition(instanceID),
			),
		))
		require.NoError(t, err)
		assert.ElementsMatch(t, recoveryCodes, gotUser.Human.RecoveryCodes.Codes)
		assert.Zero(t, gotUser.Human.RecoveryCodes.LastSuccessfullyCheckedAt)
		assert.Zero(t, gotUser.Human.RecoveryCodes.FailedAttempts)
	})

	t.Run("reduce user.human.mfa.recoverycode.removed event", func(t *testing.T) {
		// create user
		existingUserID := gofakeit.UUID()
		existingUserAgg := createUser(t, tx, userRepo, instanceID, orgID, existingUserID)

		// add recovery codes to the user
		_, err = userRepo.Update(t.Context(), tx,
			userRepo.PrimaryKeyCondition(instanceID, existingUserID),
			userRepo.Human().AddRecoveryCodes([]string{"code1", "code2", "code3"}))
		require.NoError(t, err)

		// set `user.human.mfa.recoverycode.removed` event
		recoveryCodesRemovedEvent := user.NewHumanRecoveryCodeRemovedEvent(
			t.Context(),
			&existingUserAgg.Aggregate,
			nil)

		// reduce the event
		eventReduced := callReduce(t, rawTx, handler, recoveryCodesRemovedEvent)
		require.True(t, eventReduced)

		// assert that the recovery codes are removed
		gotUser, err := userRepo.Get(t.Context(), tx, database.WithCondition(
			database.And(
				userRepo.IDCondition(existingUserID),
				userRepo.InstanceIDCondition(instanceID),
			),
		))
		require.NoError(t, err)
		assert.Empty(t, gotUser.Human.RecoveryCodes.Codes)
		assert.Zero(t, gotUser.Human.RecoveryCodes.LastSuccessfullyCheckedAt)
		assert.Zero(t, gotUser.Human.RecoveryCodes.FailedAttempts)
	})

	t.Run("reduce user.human.mfa.recoverycode.check.succeeded event", func(t *testing.T) {
		// create user
		existingUserID := gofakeit.UUID()
		existingUserAgg := createUser(t, tx, userRepo, instanceID, orgID, existingUserID)

		// add recovery codes to the user
		_, err = userRepo.Update(t.Context(), tx,
			userRepo.PrimaryKeyCondition(instanceID, existingUserID),
			userRepo.Human().AddRecoveryCodes([]string{"code1", "code2", "code3"}))
		require.NoError(t, err)

		// set `user.human.mfa.recoverycode.check.succeeded` event
		recoveryCodeCheckSucceededEvent := user.NewHumanRecoveryCodeCheckSucceededEvent(
			t.Context(),
			&existingUserAgg.Aggregate,
			"code1",
			nil)

		// reduce the event
		eventReduced := callReduce(t, rawTx, handler, recoveryCodeCheckSucceededEvent)
		require.True(t, eventReduced)

		// assert that the recovery code is removed and the last_successfully_checked_at timestamp is set
		gotUser, err := userRepo.Get(t.Context(), tx, database.WithCondition(
			database.And(
				userRepo.IDCondition(existingUserID),
				userRepo.InstanceIDCondition(instanceID),
			),
		))
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"code2", "code3"}, gotUser.Human.RecoveryCodes.Codes)
		assert.NotZero(t, gotUser.Human.RecoveryCodes.LastSuccessfullyCheckedAt)
		assert.Zero(t, gotUser.Human.RecoveryCodes.FailedAttempts)
	})

	t.Run("reduce user.human.mfa.recoverycode.check.failed event", func(t *testing.T) {
		// create user
		existingUserID := gofakeit.UUID()
		existingUserAgg := createUser(t, tx, userRepo, instanceID, orgID, existingUserID)

		// add recovery codes to the user
		_, err = userRepo.Update(t.Context(), tx,
			userRepo.PrimaryKeyCondition(instanceID, existingUserID),
			userRepo.Human().AddRecoveryCodes([]string{"code1", "code2", "code3"}))
		require.NoError(t, err)

		// set `user.human.mfa.recoverycode.check.failed` event
		recoveryCodeCheckFailedEvent := user.NewHumanRecoveryCodeCheckFailedEvent(
			t.Context(),
			&existingUserAgg.Aggregate,
			nil,
		)

		// reduce the event
		eventReduced := callReduce(t, rawTx, handler, recoveryCodeCheckFailedEvent)
		require.True(t, eventReduced)

		// assert that the recovery code failed_attempts is incremented
		gotUser, err := userRepo.Get(t.Context(), tx, database.WithCondition(
			database.And(
				userRepo.IDCondition(existingUserID),
				userRepo.InstanceIDCondition(instanceID),
			),
		))
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"code1", "code2", "code3"}, gotUser.Human.RecoveryCodes.Codes)
		assert.Zero(t, gotUser.Human.RecoveryCodes.LastSuccessfullyCheckedAt)
		assert.Equal(t, uint8(1), gotUser.Human.RecoveryCodes.FailedAttempts)
	})

	t.Run("reduce user.HumanWebAuthNVerifiedEvent event", func(t *testing.T) {
		// Given

		// create user
		existingUserID := gofakeit.UUID()
		existingUserAgg := createUser(t, tx, userRepo, instanceID, orgID, existingUserID)
		now := time.Now()

		// add passkey to the user
		pkey := &domain.Passkey{
			ID:             "pkey-id",
			Challenge:      []byte("some challenge"),
			RelyingPartyID: "rpID",
			CreatedAt:      now,
			UpdatedAt:      now,
			Type:           domain.PasskeyTypePasswordless,
		}

		_, err = userRepo.Update(t.Context(), tx,
			userRepo.PrimaryKeyCondition(instanceID, existingUserID),
			userRepo.Human().AddPasskey(pkey),
		)
		require.NoError(t, err)

		pkeyVerifiedEvt := user.NewHumanPasswordlessVerifiedEvent(
			t.Context(),
			&existingUserAgg.Aggregate,
			pkey.ID,
			"some-token-name",
			"attestation-type",
			[]byte("key-id"),
			[]byte("public-key-id"),
			[]byte("authenticator-attestation-guid"),
			33,
			"",
		)
		pkeyVerifiedEvt.Creation = time.Now()

		// Test
		eventReduced := callReduce(t, rawTx, handler, pkeyVerifiedEvt)
		require.True(t, eventReduced)

		// Verify
		gotUser, err := userRepo.Get(t.Context(), tx, database.WithCondition(
			database.And(
				userRepo.IDCondition(existingUserID),
				userRepo.InstanceIDCondition(instanceID),
			),
		))
		require.NoError(t, err)
		require.Len(t, gotUser.Human.Passkeys, 1)
		gotPKey := gotUser.Human.Passkeys[0]
		assert.Equal(t, pkeyVerifiedEvt.KeyID, gotPKey.KeyID)
		assert.Equal(t, pkeyVerifiedEvt.PublicKey, gotPKey.PublicKey)
		assert.Equal(t, pkeyVerifiedEvt.AttestationType, gotPKey.AttestationType)
		assert.Equal(t, pkeyVerifiedEvt.AAGUID, gotPKey.AuthenticatorAttestationGUID)
		assert.Equal(t, pkeyVerifiedEvt.SignCount, gotPKey.SignCount)
		assert.Equal(t, pkeyVerifiedEvt.WebAuthNTokenName, gotPKey.Name)
		assert.Equal(t, pkeyVerifiedEvt.Creation.Round(time.Millisecond), gotPKey.UpdatedAt.Round(time.Millisecond))
		assert.Equal(t, pkey.CreatedAt.Round(time.Millisecond), gotPKey.CreatedAt.Round(time.Millisecond))
		assert.NotZero(t, gotPKey.VerifiedAt)
	})
}

func createUser(t *testing.T,
	tx *sql.Transaction,
	userRepo domain.UserRepository,
	instanceID,
	orgID,
	userID string,
) *user.Aggregate {
	userAgg := user.NewAggregate(userID, orgID)
	userAgg.InstanceID = instanceID

	err := userRepo.Create(t.Context(), tx, &domain.User{
		InstanceID:     instanceID,
		OrganizationID: orgID,
		ID:             userID,
		Username:       userID + gofakeit.Username(),
		State:          domain.UserStateActive,
		Human: &domain.HumanUser{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Email: domain.HumanEmail{
				Address:    gofakeit.Email(),
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

	return userAgg
}

package repository_test

import (
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestCreateIDPIntent(t *testing.T) {
	beforeCreate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		require.NoError(t, err)
	}()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	idpID := createIdentityProvider(t, tx, instanceID, orgID)

	successURL, err := url.Parse("https://example.com/success")
	require.NoError(t, err)
	failURL, err := url.Parse("https://example.com/fail")
	require.NoError(t, err)

	emptyURL, err := url.Parse("")
	require.NoError(t, err)

	intentID := gofakeit.UUID()
	tt := []struct {
		testName string

		inputID          string
		inputInstanceID  string
		inputSuccessURL  *url.URL
		inputFailureURL  *url.URL
		inputIDPID       string
		inputIDPArgs     map[string]any
		inputCreatedAt   time.Time
		inputUpdatedAt   time.Time
		inputMaxLifetime time.Duration

		expectedError error
	}{
		{
			testName:        "create and update timestamps not set / empty fail url / should generate timestamps and retrieve successfully",
			inputID:         intentID,
			inputInstanceID: instanceID,
			inputSuccessURL: successURL,
			inputFailureURL: emptyURL,
			inputIDPID:      idpID,
			inputIDPArgs: map[string]any{
				"arg1": map[string]any{"k1": 1, "k2": "v2"},
				"arg2": 2,
				"arg3": "3",
				"arg4": true,
			},
			inputMaxLifetime: time.Hour * 2,
		},
		{
			testName:         "all set / should retrieve successfully",
			inputID:          intentID,
			inputInstanceID:  instanceID,
			inputSuccessURL:  successURL,
			inputFailureURL:  failURL,
			inputIDPID:       idpID,
			inputIDPArgs:     map[string]any{"arg1": map[string]any{"k1": 1, "k2": "v2"}},
			inputMaxLifetime: time.Hour * 2,
			inputCreatedAt:   time.Now(),
			inputUpdatedAt:   time.Now(),
		},
		{
			testName:        "id not set / should return check error",
			inputCreatedAt:  time.Now(),
			inputUpdatedAt:  time.Now(),
			inputInstanceID: instanceID,
			inputSuccessURL: successURL,
			inputFailureURL: emptyURL,
			inputIDPID:      idpID,
			inputIDPArgs:    map[string]any{"arg1": map[string]any{"k1": 1, "k2": "v2"}},
			expectedError:   database.NewCheckError("", "", nil),
		},
		{
			testName:        "instance id not found / should return foreign key error",
			inputID:         intentID,
			inputCreatedAt:  time.Now(),
			inputUpdatedAt:  time.Now(),
			inputSuccessURL: successURL,
			inputFailureURL: emptyURL,
			inputIDPID:      idpID,
			inputIDPArgs:    map[string]any{"arg1": map[string]any{"k1": 1, "k2": "v2"}},
			expectedError:   database.NewForeignKeyError("", "", nil),
		},
		{
			testName:        "idp id not found / should return foreign key error",
			inputID:         intentID,
			inputInstanceID: instanceID,
			inputCreatedAt:  time.Now(),
			inputUpdatedAt:  time.Now(),
			inputSuccessURL: successURL,
			inputFailureURL: failURL,
			inputIDPArgs:    map[string]any{"arg1": map[string]any{"k1": 1, "k2": "v2"}},
			expectedError:   database.NewForeignKeyError("", "", nil),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			savePoint, savePointErr := tx.Begin(t.Context())
			require.NoError(t, savePointErr)
			defer func() {
				err := savePoint.Rollback(t.Context())
				require.NoError(t, err)
			}()

			idpIntentRepo := repository.IDPIntentRepository()
			intent := domain.IDPIntent{
				ID:           tc.inputID,
				InstanceID:   tc.inputInstanceID,
				SuccessURL:   tc.inputSuccessURL,
				FailureURL:   tc.inputFailureURL,
				IDPID:        tc.inputIDPID,
				IDPArguments: tc.inputIDPArgs,

				State:     domain.IDPIntentStateFailed, // Should disregard
				IDPUserID: "should disregard",
			}

			// Test
			createErr := idpIntentRepo.Create(t.Context(), tx, &intent)
			afterCreate := time.Now()

			// Verify
			assert.ErrorIs(t, createErr, tc.expectedError)
			if tc.expectedError == nil {
				retrieved, getErr := idpIntentRepo.Get(t.Context(), tx, database.WithCondition(idpIntentRepo.PrimaryKeyCondition(instanceID, intentID)))
				require.NoError(t, getErr)
				require.NotNil(t, retrieved)

				assert.WithinRange(t, retrieved.CreatedAt, beforeCreate, afterCreate)
				assert.WithinRange(t, retrieved.UpdatedAt, beforeCreate, afterCreate)
				assert.NotZero(t, retrieved.ID)
				if tc.inputID != "" {
					assert.Equal(t, tc.inputID, retrieved.ID)
				}
				assert.Equal(t, tc.inputInstanceID, retrieved.InstanceID)
				assert.Equal(t, tc.inputSuccessURL, retrieved.SuccessURL)
				assert.Equal(t, tc.inputFailureURL, retrieved.FailureURL)
				assert.Equal(t, tc.inputIDPID, retrieved.IDPID)
				assert.NotZero(t, retrieved.IDPArguments)
				assert.Equal(t, domain.IDPIntentStateStarted, retrieved.State)
				assert.Zero(t, retrieved.IDPUserID)
			}
		})
	}
}

func TestGetIDPIntent(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		require.NoError(t, err)
	}()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	idpID := createIdentityProvider(t, tx, instanceID, orgID)

	idpIntentRepo := repository.IDPIntentRepository()

	intentID1 := createIDPIntent(t, tx, instanceID, idpID)
	successURL, err := url.Parse("https://example.com/success")
	require.NoError(t, err)
	failURL, err := url.Parse("https://example.com/fail")
	require.NoError(t, err)

	intent2 := domain.IDPIntent{
		ID:           gofakeit.UUID(),
		InstanceID:   instanceID,
		SuccessURL:   successURL,
		FailureURL:   failURL,
		IDPID:        idpID,
		IDPArguments: map[string]any{"arg1": map[string]any{"k1": 1, "k2": "v2"}},
		CreatedAt:    time.Now().AddDate(0, 0, 1),
	}
	err = idpIntentRepo.Create(t.Context(), tx, &intent2)
	require.NoError(t, err)

	tt := []struct {
		testName         string
		inputQueryOpts   []database.QueryOption
		expectedError    error
		expectedIntentID string
	}{
		{
			testName:      "when no condition set should return missing condition error",
			expectedError: database.NewMissingConditionError(nil),
		},
		{
			testName: "when no instance condition set should return missing condition error",
			inputQueryOpts: []database.QueryOption{
				database.WithCondition(idpIntentRepo.IDCondition(intentID1)),
			},
			expectedError: database.NewMissingConditionError(nil),
		},
		{
			testName: "when primary key condition set should return matching intent",
			inputQueryOpts: []database.QueryOption{
				database.WithCondition(idpIntentRepo.PrimaryKeyCondition(instanceID, intentID1)),
			},
			expectedIntentID: intentID1,
		},
		{
			testName: "when filtering by creation date descending should return matching intent with highest PK",
			inputQueryOpts: []database.QueryOption{
				database.WithCondition(database.And(
					idpIntentRepo.InstanceIDCondition(instanceID),
				)),
				database.WithOrderByDescending(idpIntentRepo.CreatedAtColumn()),
				database.WithLimit(1),
			},
			expectedIntentID: intent2.ID,
		},
		{
			testName: "when filtering by non-existent PK should return no row found error",
			inputQueryOpts: []database.QueryOption{
				database.WithCondition(database.And(
					idpIntentRepo.InstanceIDCondition(instanceID),
					idpIntentRepo.IDCondition("not existing"),
				)),
			},
			expectedError: database.NewNoRowFoundError(nil),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			retrieved, err := idpIntentRepo.Get(t.Context(), tx, tc.inputQueryOpts...)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expectedIntentID, retrieved.ID)
				assert.Equal(t, "https://example.com/success", retrieved.SuccessURL.String())
				assert.Equal(t, "https://example.com/fail", retrieved.FailureURL.String())
			}
		})
	}
}

func TestUpdateIDPIntent(t *testing.T) {
	beforeUpdate := time.Now()

	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		require.NoError(t, err)
	}()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	idpID := createIdentityProvider(t, tx, instanceID, orgID)

	idpIntentRepo := repository.IDPIntentRepository()

	successURL, err := url.Parse("https://example.com/success")
	require.NoError(t, err)
	failURL, err := url.Parse("https://example.com/fail")
	require.NoError(t, err)

	intentID1 := gofakeit.UUID()
	intent1 := domain.IDPIntent{
		ID:           intentID1,
		InstanceID:   instanceID,
		SuccessURL:   successURL,
		FailureURL:   failURL,
		IDPID:        idpID,
		IDPArguments: map[string]any{"arg1": map[string]any{"k1": 1, "k2": "v2"}, "arg2": []int32{1}},
		CreatedAt:    time.Now(),
	}
	err = idpIntentRepo.Create(t.Context(), tx, &intent1)
	require.NoError(t, err)

	entryAttrs := map[string][]string{
		"attr1": {"1", "2"},
		"attr2": {"3", "4"},
	}
	marshalledAttrs, marshallErr := json.Marshal(entryAttrs)
	require.NoError(t, marshallErr)

	intentID2 := gofakeit.UUID()
	intent2 := domain.IDPIntent{
		ID:           intentID2,
		InstanceID:   instanceID,
		SuccessURL:   successURL,
		FailureURL:   failURL,
		IDPID:        idpID,
		IDPArguments: map[string]any{"arg1": map[string]any{"k1": 1, "k2": "v2"}},
		CreatedAt:    time.Now(),
	}
	err = idpIntentRepo.Create(t.Context(), tx, &intent2)
	require.NoError(t, err)

	userID := createHumanUser(t, tx, instanceID, orgID)

	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	tt := []struct {
		testName               string
		inputConditions        database.Condition
		inputChanges           []database.Change
		expectedError          error
		expectedUpdatedRecords func() []*domain.IDPIntent
		expectedUpdatedRows    int64
	}{
		{
			testName:      "when no changes should return no changes error",
			expectedError: database.ErrNoChanges,
		},
		{
			testName: "when no condition set should return missing condition error",
			inputChanges: database.Changes{
				idpIntentRepo.SetIDPUser([]byte("not a user")),
			},
			expectedError: database.NewMissingConditionError(idpIntentRepo.InstanceIDColumn()),
		},
		{
			testName:        "when simulating SucceededEvent should update fields accordingly",
			inputConditions: idpIntentRepo.PrimaryKeyCondition(instanceID, intentID1),
			inputChanges: database.Changes{
				idpIntentRepo.SetState(domain.IDPIntentStateSucceeded),
				idpIntentRepo.SetIDPUser([]byte("update user")),
				idpIntentRepo.SetIDPUserID("idp user id updated"),
				idpIntentRepo.SetIDPUsername("idp username updated"),
				idpIntentRepo.SetUserID(userID),
				idpIntentRepo.SetIDPAccessToken([]byte(`{"secret": "g1g4 s3cr3t 4cc3ss t0k3n"}`)),
				idpIntentRepo.SetIDPIDToken("idp id token"),
				idpIntentRepo.SetSucceededAt(now),
				idpIntentRepo.SetExpiresAt(tomorrow),
			},
			expectedUpdatedRecords: func() []*domain.IDPIntent {
				toReturn := intent1
				toReturn.IDPUser = []byte("update user")
				toReturn.IDPUserID = "idp user id updated"
				toReturn.IDPUsername = "idp username updated"
				toReturn.UserID = userID
				toReturn.IDPAccessToken = []byte(`{"secret": "g1g4 s3cr3t 4cc3ss t0k3n"}`)
				toReturn.IDPIDToken = "idp id token"
				toReturn.SucceededAt = &now
				toReturn.ExpiresAt = &tomorrow
				toReturn.State = domain.IDPIntentStateSucceeded
				return []*domain.IDPIntent{&toReturn}
			},
			expectedUpdatedRows: 1,
		},
		{
			testName:        "when simulating SAMLSucceededEvent should update fields accordingly",
			inputConditions: idpIntentRepo.PrimaryKeyCondition(instanceID, intentID1),
			inputChanges: database.Changes{
				idpIntentRepo.SetState(domain.IDPIntentStateSucceeded),
				idpIntentRepo.SetIDPUser([]byte("update user")),
				idpIntentRepo.SetIDPUserID("idp user id updated"),
				idpIntentRepo.SetIDPUsername("idp username updated"),
				idpIntentRepo.SetUserID(userID),
				idpIntentRepo.SetAssertion([]byte(`{"assertion": "val1"}`)),
				idpIntentRepo.SetSucceededAt(now),
				idpIntentRepo.SetExpiresAt(tomorrow),
			},
			expectedUpdatedRecords: func() []*domain.IDPIntent {
				toReturn := intent1
				toReturn.IDPUser = []byte("update user")
				toReturn.IDPUserID = "idp user id updated"
				toReturn.IDPUsername = "idp username updated"
				toReturn.UserID = userID
				toReturn.Assertion = []byte(`{"assertion": "val1"}`)
				toReturn.SucceededAt = &now
				toReturn.ExpiresAt = &tomorrow
				toReturn.State = domain.IDPIntentStateSucceeded
				return []*domain.IDPIntent{&toReturn}
			},
			expectedUpdatedRows: 1,
		},
		{
			testName:        "when simulating LDAPSucceededEvent should update fields accordingly",
			inputConditions: idpIntentRepo.PrimaryKeyCondition(instanceID, intentID1),
			inputChanges: database.Changes{
				idpIntentRepo.SetState(domain.IDPIntentStateSucceeded),
				idpIntentRepo.SetIDPUser([]byte("update user")),
				idpIntentRepo.SetIDPUserID("idp user id updated"),
				idpIntentRepo.SetIDPUsername("idp username updated"),
				idpIntentRepo.SetUserID(userID),
				idpIntentRepo.SetIDPEntryAttributes(marshalledAttrs),
				idpIntentRepo.SetSucceededAt(now),
				idpIntentRepo.SetExpiresAt(tomorrow),
			},
			expectedUpdatedRecords: func() []*domain.IDPIntent {
				toReturn := intent1
				toReturn.IDPUser = []byte("update user")
				toReturn.IDPUserID = "idp user id updated"
				toReturn.IDPUsername = "idp username updated"
				toReturn.UserID = userID
				toReturn.EntryAttributes = entryAttrs
				toReturn.SucceededAt = &now
				toReturn.ExpiresAt = &tomorrow
				toReturn.State = domain.IDPIntentStateSucceeded
				return []*domain.IDPIntent{&toReturn}
			},
			expectedUpdatedRows: 1,
		},
		{
			testName:        "when simulating FailedEvent should update reason, state and failedAt",
			inputConditions: idpIntentRepo.PrimaryKeyCondition(instanceID, intentID1),
			inputChanges: database.Changes{
				idpIntentRepo.SetState(domain.IDPIntentStateFailed),
				idpIntentRepo.SetFailReason("mock failure"),
				idpIntentRepo.SetFailedAt(now),
			},
			expectedUpdatedRecords: func() []*domain.IDPIntent {
				toReturn := intent1
				toReturn.FailReason = "mock failure"
				toReturn.State = domain.IDPIntentStateFailed
				toReturn.FailedAt = &now
				return []*domain.IDPIntent{&toReturn}
			},
			expectedUpdatedRows: 1,
		},
		{
			testName:        "when simulating SAMLRequestEvent should update request ID",
			inputConditions: idpIntentRepo.PrimaryKeyCondition(instanceID, intentID1),
			inputChanges: database.Changes{
				idpIntentRepo.SetRequestID("req-123"),
			},
			expectedUpdatedRecords: func() []*domain.IDPIntent {
				toReturn := intent1
				toReturn.RequestID = "req-123"
				return []*domain.IDPIntent{&toReturn}
			},
			expectedUpdatedRows: 1,
		},
		{
			testName: "when updating multiple records should return updated records",
			inputConditions: database.Or(
				idpIntentRepo.PrimaryKeyCondition(instanceID, intentID1),
				idpIntentRepo.PrimaryKeyCondition(instanceID, intentID2),
			),
			inputChanges: database.Changes{
				idpIntentRepo.SetState(domain.IDPIntentStateFailed),
				idpIntentRepo.SetRequestID("req-123"),
				idpIntentRepo.SetFailedAt(now),
			},
			expectedUpdatedRecords: func() []*domain.IDPIntent {
				toReturn1, toReturn2 := intent1, intent2
				toReturn1.State = domain.IDPIntentStateFailed
				toReturn1.RequestID = "req-123"
				toReturn1.FailedAt = &now
				toReturn2.State = domain.IDPIntentStateFailed
				toReturn2.RequestID = "req-123"
				toReturn2.FailedAt = &now
				return []*domain.IDPIntent{&toReturn1, &toReturn2}
			},
			expectedUpdatedRows: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			savePoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err := savePoint.Rollback(t.Context())
				require.NoError(t, err)
			}()

			// Test
			updatedCount, err := idpIntentRepo.Update(t.Context(), savePoint, tc.inputConditions, tc.inputChanges...)
			afterUpdate := time.Now()

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expectedUpdatedRows, updatedCount)
			if tc.expectedError == nil {
				for _, expectedIntent := range tc.expectedUpdatedRecords() {
					retrievedIntent, err := idpIntentRepo.Get(t.Context(), savePoint, database.WithCondition(idpIntentRepo.PrimaryKeyCondition(instanceID, expectedIntent.ID)))
					require.NoError(t, err)
					assert.Equal(t, expectedIntent.ID, retrievedIntent.ID)
					assert.Equal(t, expectedIntent.InstanceID, retrievedIntent.InstanceID)
					assert.Equal(t, expectedIntent.State, retrievedIntent.State)
					assert.Equal(t, expectedIntent.SuccessURL, retrievedIntent.SuccessURL)
					assert.Equal(t, expectedIntent.FailureURL, retrievedIntent.FailureURL)
					assert.Equal(t, expectedIntent.CreatedAt, retrievedIntent.CreatedAt)
					assert.WithinRange(t, retrievedIntent.UpdatedAt, beforeUpdate, afterUpdate)
					assert.Equal(t, expectedIntent.IDPID, retrievedIntent.IDPID)
					assert.Equal(t, expectedIntent.IDPUser, retrievedIntent.IDPUser)
					assert.Equal(t, expectedIntent.IDPUserID, retrievedIntent.IDPUserID)
					assert.Equal(t, expectedIntent.IDPUsername, retrievedIntent.IDPUsername)
					assert.Equal(t, expectedIntent.UserID, retrievedIntent.UserID)
					assert.Equal(t, expectedIntent.IDPAccessToken, retrievedIntent.IDPAccessToken)
					assert.Equal(t, expectedIntent.IDPIDToken, retrievedIntent.IDPIDToken)
					assert.Equal(t, expectedIntent.EntryAttributes, retrievedIntent.EntryAttributes)
					assert.Equal(t, expectedIntent.RequestID, retrievedIntent.RequestID)
					assert.Equal(t, expectedIntent.Assertion, retrievedIntent.Assertion)
					if expectedIntent.State == domain.IDPIntentStateSucceeded {
						require.NotNil(t, retrievedIntent.SucceededAt)
						assert.WithinRange(t, *retrievedIntent.SucceededAt, beforeUpdate, afterUpdate)
						require.NotNil(t, retrievedIntent.ExpiresAt)
						assert.WithinRange(t, *retrievedIntent.ExpiresAt, beforeUpdate.AddDate(0, 0, 1), afterUpdate.AddDate(0, 0, 1))
					}
					if expectedIntent.State == domain.IDPIntentStateFailed {
						require.NotNil(t, retrievedIntent.FailedAt)
						assert.WithinRange(t, *retrievedIntent.FailedAt, beforeUpdate, afterUpdate)
					}
					assert.Equal(t, expectedIntent.FailReason, retrievedIntent.FailReason)
				}
			}
		})
	}
}

func TestDeleteIDPIntent(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		require.NoError(t, err)
	}()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	idpID := createIdentityProvider(t, tx, instanceID, orgID)
	intentID := createIDPIntent(t, tx, instanceID, idpID)
	otherIntentID := createIDPIntent(t, tx, instanceID, idpID)
	createIDPIntent(t, tx, instanceID, idpID)
	idpIntentRepo := repository.IDPIntentRepository()

	tt := []struct {
		testName             string
		inputConditions      database.Condition
		expectedError        error
		expectedDeletedCount int64
	}{
		{
			testName:        "when PK condition not set should return missing condition error",
			inputConditions: idpIntentRepo.StateCondition(domain.IDPIntentStateFailed),
			expectedError:   database.NewMissingConditionError(nil),
		},
		{
			testName:             "when condition matches record should delete matching record",
			inputConditions:      idpIntentRepo.PrimaryKeyCondition(instanceID, intentID),
			expectedDeletedCount: 1,
		},
		{
			testName: "when conditions matches record should delete matching records",
			inputConditions: database.Or(
				idpIntentRepo.PrimaryKeyCondition(instanceID, intentID),
				idpIntentRepo.PrimaryKeyCondition(instanceID, otherIntentID),
			),
			expectedDeletedCount: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			savePoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err := savePoint.Rollback(t.Context())
				require.NoError(t, err)
			}()

			// Test
			deleteCount, err := idpIntentRepo.Delete(t.Context(), savePoint, tc.inputConditions)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedDeletedCount, deleteCount)

		})
	}
}

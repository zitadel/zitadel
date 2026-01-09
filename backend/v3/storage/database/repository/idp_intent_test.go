package repository_test

import (
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
		inputUpdateddAt  time.Time
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
			inputUpdateddAt:  time.Now(),
		},
		{
			testName:        "id not set / should return check error",
			inputCreatedAt:  time.Now(),
			inputUpdateddAt: time.Now(),
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
			inputUpdateddAt: time.Now(),
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
			inputUpdateddAt: time.Now(),
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
				ID:                   tc.inputID,
				InstanceID:           tc.inputInstanceID,
				SuccessURL:           tc.inputSuccessURL,
				FailureURL:           tc.inputFailureURL,
				IDPID:                tc.inputIDPID,
				IDPArguments:         tc.inputIDPArgs,
				MaxIDPIntentLifetime: tc.inputMaxLifetime,

				State:     domain.IDPIntentStateConsumed, // Should disregard
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
				assert.Equal(t, tc.inputMaxLifetime, retrieved.MaxIDPIntentLifetime)
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
	intentID2 := createIDPIntent(t, tx, instanceID, idpID)

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
			testName: "when filtering by state should return matching intent with lowest PK",
			inputQueryOpts: []database.QueryOption{
				database.WithCondition(database.And(
					idpIntentRepo.InstanceIDCondition(instanceID),
					idpIntentRepo.StateCondition(domain.IDPIntentStateStarted),
				)),
				database.WithLimit(1),
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
			expectedIntentID: intentID2,
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

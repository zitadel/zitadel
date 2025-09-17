package repository_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestGetOrganizationMetadata(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	// create instance
	instanceID := gofakeit.UUID()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = repository.InstanceRepository(tx).Create(t.Context(), &instance)
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository(tx)
	metadataRepo := orgRepo.Metadata(false)

	// create organization
	orgA := domain.Organization{
		ID:         "1",
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), &orgA)
	require.NoError(t, err)

	err = metadataRepo.Set(
		t.Context(),
		&domain.OrganizationMetadata{
			OrgID: orgA.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key",
				Value:      []byte("1234"),
			},
		},
		&domain.OrganizationMetadata{
			OrgID: orgA.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key2",
				Value:      []byte("asdf"),
			},
		},
	)
	require.NoError(t, err)

	orgB := domain.Organization{
		ID:         "2",
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), &orgB)
	require.NoError(t, err)

	err = metadataRepo.Set(
		t.Context(),
		&domain.OrganizationMetadata{
			OrgID: orgB.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key",
				Value:      []byte("5678"),
			},
		},
		&domain.OrganizationMetadata{
			OrgID: orgB.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key2",
				Value:      []byte(`"asdf"`),
			},
		},
	)
	require.NoError(t, err)

	t.Run("metadata without instance condition", func(t *testing.T) {
		metadata, err := metadataRepo.Get(
			t.Context(),
			database.WithCondition(metadataRepo.OrgIDCondition(orgA.ID)),
		)
		assert.ErrorIs(t, err, new(database.MissingConditionError))
		assert.Nil(t, metadata)
	})

	t.Run("no metadata found", func(t *testing.T) {
		metadata, err := metadataRepo.Get(
			t.Context(),
			database.WithCondition(metadataRepo.InstanceIDCondition("non-existing")),
		)
		assert.ErrorIs(t, err, new(database.NoRowFoundError))
		assert.Empty(t, metadata)
	})

	t.Run("multiple metadata found", func(t *testing.T) {
		metadata, err := metadataRepo.Get(
			t.Context(),
			database.WithCondition(metadataRepo.InstanceIDCondition(instanceID)),
		)
		require.ErrorIs(t, err, new(database.MultipleRowsFoundError))
		assert.Empty(t, metadata)
	})

	t.Run("metadata by key", func(t *testing.T) {
		metadata, err := metadataRepo.Get(
			t.Context(),
			database.WithCondition(
				database.And(
					metadataRepo.InstanceIDCondition(instanceID),
					metadataRepo.OrgIDCondition(orgA.ID),
					metadataRepo.KeyCondition(database.TextOperationEqual, "urn:zitadel:key2"),
				),
			),
			database.WithOrderByAscending(metadataRepo.OrgIDColumn()),
		)
		require.NoError(t, err)
		assert.Equal(t, "urn:zitadel:key2", metadata.Key)
		assert.Equal(t, []byte(`asdf`), metadata.Value)
	})

	t.Run("metadata by value", func(t *testing.T) {
		metadata, err := metadataRepo.Get(
			t.Context(),
			database.WithCondition(
				database.And(
					metadataRepo.InstanceIDCondition(instanceID),
					metadataRepo.OrgIDCondition(orgA.ID),
					metadataRepo.ValueCondition(database.BytesOperationEqual, []byte("asdf")),
				),
			),
			database.WithOrderByAscending(metadataRepo.OrgIDColumn()),
		)
		require.NoError(t, err)
		assert.Equal(t, "urn:zitadel:key2", metadata.Key)
	})
}

func TestListOrganizationMetadata(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	// create instance
	instanceID := gofakeit.UUID()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = repository.InstanceRepository(tx).Create(t.Context(), &instance)
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository(tx)
	metadataRepo := orgRepo.Metadata(false)

	// create organization
	orgA := domain.Organization{
		ID:         "1",
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), &orgA)
	require.NoError(t, err)

	err = metadataRepo.Set(
		t.Context(),
		&domain.OrganizationMetadata{
			OrgID: orgA.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key",
				Value:      []byte("1234"),
			},
		},
		&domain.OrganizationMetadata{
			OrgID: orgA.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key2",
				Value:      []byte("asdf"),
			},
		},
	)
	require.NoError(t, err)

	metadataOrgA, err := metadataRepo.List(
		t.Context(),
		database.WithCondition(
			database.And(
				metadataRepo.OrgIDCondition(orgA.ID),
				metadataRepo.InstanceIDCondition(instanceID),
			),
		),
	)
	require.NoError(t, err)

	orgB := domain.Organization{
		ID:         "2",
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), &orgB)
	require.NoError(t, err)

	err = metadataRepo.Set(
		t.Context(),
		&domain.OrganizationMetadata{
			OrgID: orgB.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key",
				Value:      []byte("5678"),
			},
		},
		&domain.OrganizationMetadata{
			OrgID: orgB.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key2",
				Value:      []byte(`"asdf"`),
			},
		},
	)
	require.NoError(t, err)

	metadataOrgB, err := metadataRepo.List(
		t.Context(),
		database.WithCondition(
			database.And(
				metadataRepo.OrgIDCondition(orgB.ID),
				metadataRepo.InstanceIDCondition(instanceID),
			),
		),
	)
	require.NoError(t, err)

	t.Run("metadata without instance condition", func(t *testing.T) {
		metadata, err := metadataRepo.List(
			t.Context(),
			database.WithCondition(metadataRepo.OrgIDCondition(orgA.ID)),
		)
		assert.ErrorIs(t, err, new(database.MissingConditionError))
		assert.Nil(t, metadata)
	})

	t.Run("no metadata found", func(t *testing.T) {
		metadata, err := metadataRepo.List(
			t.Context(),
			database.WithCondition(metadataRepo.InstanceIDCondition("non-existing")),
		)
		assert.NoError(t, err)
		assert.Len(t, metadata, 0)
	})

	t.Run("all metadata of instance", func(t *testing.T) {
		metadata, err := metadataRepo.List(
			t.Context(),
			database.WithCondition(metadataRepo.InstanceIDCondition(instanceID)),
		)
		require.NoError(t, err)
		assert.ElementsMatch(t, metadata, append(metadataOrgA, metadataOrgB...))
	})

	t.Run("metadata by org id", func(t *testing.T) {
		metadata, err := metadataRepo.List(
			t.Context(),
			database.WithCondition(
				database.And(
					metadataRepo.InstanceIDCondition(instanceID),
					metadataRepo.OrgIDCondition(orgA.ID),
				),
			),
		)
		require.NoError(t, err)
		assert.ElementsMatch(t, metadata, metadataOrgA)
	})

	t.Run("metadata by key", func(t *testing.T) {
		metadata, err := metadataRepo.List(
			t.Context(),
			database.WithCondition(
				database.And(
					metadataRepo.InstanceIDCondition(instanceID),
					metadataRepo.KeyCondition(database.TextOperationEqual, "urn:zitadel:key2"),
				),
			),
			database.WithOrderByAscending(metadataRepo.OrgIDColumn()),
		)
		require.NoError(t, err)
		require.Len(t, metadata, 2)
		assert.Equal(t, []byte("asdf"), metadata[0].Value)
		assert.Equal(t, []byte(`"asdf"`), metadata[1].Value)
	})

	t.Run("metadata by value", func(t *testing.T) {
		metadata, err := metadataRepo.List(
			t.Context(),
			database.WithCondition(
				database.And(
					metadataRepo.InstanceIDCondition(instanceID),
					metadataRepo.ValueCondition(database.BytesOperationEqual, []byte("asdf")),
				),
			),
			database.WithOrderByAscending(metadataRepo.OrgIDColumn()),
		)
		require.NoError(t, err)
		require.Len(t, metadata, 1)
		assert.Equal(t, "urn:zitadel:key2", metadata[0].Key)
		assert.Equal(t, "1", metadata[0].OrgID)
	})
}

func TestSetOrganizationMetadata(t *testing.T) {
	beforeSet := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	// create instance
	instanceID := gofakeit.UUID()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = repository.InstanceRepository(tx).Create(t.Context(), &instance)
	require.NoError(t, err)

	// create organization
	orgID := gofakeit.UUID()
	organization := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = repository.OrganizationRepository(tx).Create(t.Context(), &organization)
	require.NoError(t, err)

	t.Run("set one organization metadata", func(t *testing.T) {
		savepoint, err := tx.Begin(t.Context())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, savepoint.Rollback(t.Context()))
		}()
		metadataRepo := repository.OrganizationRepository(savepoint).Metadata(false)

		err = metadataRepo.Set(t.Context(), &domain.OrganizationMetadata{
			OrgID: orgID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key",
				Value:      []byte("some-value"),
			},
		})
		require.NoError(t, err)
		afterSet := time.Now()

		savedMetadata, err := metadataRepo.Get(
			t.Context(), database.WithCondition(
				database.And(
					metadataRepo.InstanceIDCondition(instanceID),
					metadataRepo.OrgIDCondition(orgID),
					metadataRepo.KeyCondition(database.TextOperationEqual, "urn:zitadel:key"),
				),
			),
		)
		require.NoError(t, err)
		assert.Equal(t, orgID, savedMetadata.OrgID)
		assert.Equal(t, instanceID, savedMetadata.InstanceID)
		assert.Equal(t, "urn:zitadel:key", savedMetadata.Key)
		assert.Equal(t, []byte("some-value"), savedMetadata.Value)
		assert.WithinRange(t, savedMetadata.CreatedAt, beforeSet, afterSet)
		assert.WithinRange(t, savedMetadata.UpdatedAt, beforeSet, afterSet)
		assert.Equal(t, savedMetadata.CreatedAt, savedMetadata.UpdatedAt)
	})

	t.Run("set multiple organization metadata", func(t *testing.T) {
		savepoint, err := tx.Begin(t.Context())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, savepoint.Rollback(t.Context()))
		}()
		metadataRepo := repository.OrganizationRepository(savepoint).Metadata(false)

		err = metadataRepo.Set(t.Context(),
			&domain.OrganizationMetadata{
				OrgID: orgID,
				Metadata: domain.Metadata{
					InstanceID: instanceID,
					Key:        "urn:zitadel:key",
					Value:      []byte("some-value"),
				},
			},
			&domain.OrganizationMetadata{
				OrgID: orgID,
				Metadata: domain.Metadata{
					InstanceID: instanceID,
					Key:        "urn:zitadel:key2",
					Value:      []byte("1234"),
				},
			},
		)
		require.NoError(t, err)
		afterSet := time.Now()

		savedMetadata, err := metadataRepo.List(
			t.Context(),
			database.WithCondition(
				database.And(
					metadataRepo.InstanceIDCondition(instanceID),
					metadataRepo.OrgIDCondition(orgID),
					metadataRepo.KeyCondition(database.TextOperationStartsWith, "urn:zitadel:key"),
				),
			),
			database.WithOrderByAscending(metadataRepo.KeyColumn()),
		)
		require.NoError(t, err)
		require.Len(t, savedMetadata, 2)

		for _, saved := range savedMetadata {
			assert.Equal(t, orgID, saved.OrgID)
			assert.Equal(t, instanceID, saved.InstanceID)
			assert.WithinRange(t, saved.CreatedAt, beforeSet, afterSet)
			assert.WithinRange(t, saved.UpdatedAt, beforeSet, afterSet)
			assert.Equal(t, saved.CreatedAt, saved.UpdatedAt)
		}

		assert.Equal(t, "urn:zitadel:key", savedMetadata[0].Key)
		assert.Equal(t, []byte("some-value"), savedMetadata[0].Value)

		assert.Equal(t, "urn:zitadel:key2", savedMetadata[1].Key)
		assert.Equal(t, []byte("1234"), savedMetadata[1].Value)
	})

	t.Run("set no organization metadata", func(t *testing.T) {
		savepoint, err := tx.Begin(t.Context())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, savepoint.Rollback(t.Context()))
		}()
		metadataRepo := repository.OrganizationRepository(savepoint).Metadata(false)

		err = metadataRepo.Set(t.Context())
		require.ErrorIs(t, err, database.ErrNoChanges)
	})

	t.Run("overwrite organization metadata", func(t *testing.T) {
		savepoint, err := tx.Begin(t.Context())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, savepoint.Rollback(t.Context()))
		}()
		metadataRepo := repository.OrganizationRepository(savepoint).Metadata(false)

		err = metadataRepo.Set(t.Context(), &domain.OrganizationMetadata{
			OrgID: orgID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key",
				Value:      []byte("some-value"),
			},
		})
		require.NoError(t, err)
		afterSet := time.Now()

		savedMetadata, err := metadataRepo.Get(
			t.Context(), database.WithCondition(
				database.And(
					metadataRepo.InstanceIDCondition(instanceID),
					metadataRepo.OrgIDCondition(orgID),
					metadataRepo.KeyCondition(database.TextOperationEqual, "urn:zitadel:key"),
				),
			),
		)
		require.NoError(t, err)
		assert.Equal(t, orgID, savedMetadata.OrgID)
		assert.Equal(t, instanceID, savedMetadata.InstanceID)
		assert.Equal(t, "urn:zitadel:key", savedMetadata.Key)
		assert.Equal(t, []byte("some-value"), savedMetadata.Value)
		assert.WithinRange(t, savedMetadata.CreatedAt, beforeSet, afterSet)
		assert.WithinRange(t, savedMetadata.UpdatedAt, beforeSet, afterSet)
		assert.Equal(t, savedMetadata.CreatedAt, savedMetadata.UpdatedAt)

		err = metadataRepo.Set(t.Context(), &domain.OrganizationMetadata{
			OrgID: orgID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key",
				Value:      []byte("1234"),
			},
		})
		require.NoError(t, err)
		afterSet = time.Now()

		savedMetadata, err = metadataRepo.Get(
			t.Context(), database.WithCondition(
				database.And(
					metadataRepo.InstanceIDCondition(instanceID),
					metadataRepo.OrgIDCondition(orgID),
					metadataRepo.KeyCondition(database.TextOperationEqual, "urn:zitadel:key"),
				),
			),
		)
		require.NoError(t, err)

		assert.Equal(t, orgID, savedMetadata.OrgID)
		assert.Equal(t, instanceID, savedMetadata.InstanceID)
		assert.Equal(t, "urn:zitadel:key", savedMetadata.Key)
		assert.Equal(t, []byte("1234"), savedMetadata.Value)
		assert.WithinRange(t, savedMetadata.CreatedAt, beforeSet, afterSet)
		assert.WithinRange(t, savedMetadata.UpdatedAt, beforeSet, afterSet)
		// we cannot check if the updated at did change because we are in the same transaction, so we skip this check
	})

	t.Run("from events", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			savepoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				require.NoError(t, savepoint.Rollback(t.Context()))
			}()
			metadataRepo := repository.OrganizationRepository(savepoint).Metadata(false)

			err = metadataRepo.Set(t.Context(), &domain.OrganizationMetadata{
				OrgID: orgID,
				Metadata: domain.Metadata{
					InstanceID: instanceID,
					Key:        "urn:zitadel:key",
					Value:      []byte("some-value"),
					CreatedAt:  beforeSet.Add(-time.Hour),
					UpdatedAt:  beforeSet.Add(-time.Hour),
				},
			})
			require.NoError(t, err)

			savedMetadata, err := metadataRepo.Get(
				t.Context(), database.WithCondition(
					database.And(
						metadataRepo.InstanceIDCondition(instanceID),
						metadataRepo.OrgIDCondition(orgID),
						metadataRepo.KeyCondition(database.TextOperationEqual, "urn:zitadel:key"),
					),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, orgID, savedMetadata.OrgID)
			assert.Equal(t, instanceID, savedMetadata.InstanceID)
			assert.Equal(t, "urn:zitadel:key", savedMetadata.Key)
			assert.Equal(t, []byte("some-value"), savedMetadata.Value)
			assert.Less(t, savedMetadata.CreatedAt, beforeSet)
			assert.Less(t, savedMetadata.UpdatedAt, beforeSet)
			assert.Equal(t, savedMetadata.CreatedAt, savedMetadata.UpdatedAt)
		})

		t.Run("update", func(t *testing.T) {
			savepoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				require.NoError(t, savepoint.Rollback(t.Context()))
			}()
			metadataRepo := repository.OrganizationRepository(savepoint).Metadata(false)

			// first event
			firstEventCreatedAt := beforeSet.Add(-2 * time.Hour)
			err = metadataRepo.Set(t.Context(), &domain.OrganizationMetadata{
				OrgID: orgID,
				Metadata: domain.Metadata{
					InstanceID: instanceID,
					Key:        "urn:zitadel:key",
					Value:      []byte("some-value"),
					CreatedAt:  firstEventCreatedAt,
					UpdatedAt:  firstEventCreatedAt,
				},
			})
			require.NoError(t, err)

			// second event
			secondEventCreatedAt := beforeSet.Add(-time.Hour)
			err = metadataRepo.Set(t.Context(), &domain.OrganizationMetadata{
				OrgID: orgID,
				Metadata: domain.Metadata{
					InstanceID: instanceID,
					Key:        "urn:zitadel:key",
					Value:      []byte("some-other-value"),
					CreatedAt:  secondEventCreatedAt,
					UpdatedAt:  secondEventCreatedAt,
				},
			})
			require.NoError(t, err)

			savedMetadata, err := metadataRepo.Get(
				t.Context(), database.WithCondition(
					database.And(
						metadataRepo.InstanceIDCondition(instanceID),
						metadataRepo.OrgIDCondition(orgID),
						metadataRepo.KeyCondition(database.TextOperationEqual, "urn:zitadel:key"),
					),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, orgID, savedMetadata.OrgID)
			assert.Equal(t, instanceID, savedMetadata.InstanceID)
			assert.Equal(t, "urn:zitadel:key", savedMetadata.Key)
			assert.Equal(t, []byte("some-other-value"), savedMetadata.Value)
			assert.True(t, savedMetadata.CreatedAt.Equal(firstEventCreatedAt))
			assert.True(t, savedMetadata.UpdatedAt.Equal(secondEventCreatedAt), "updated at should have been %v, but was %v", secondEventCreatedAt, savedMetadata.UpdatedAt)
		})
	})

	t.Run("invalid input", func(t *testing.T) {
		for _, testCase := range []struct {
			name        string
			metadata    *domain.OrganizationMetadata
			expectedErr error
		}{
			{
				name: "missing org_id",
				metadata: &domain.OrganizationMetadata{
					Metadata: domain.Metadata{
						InstanceID: instanceID,
						Key:        "urn:zitadel:key",
						Value:      []byte("1234"),
					},
				},
				expectedErr: new(database.ForeignKeyError),
			},
			{
				name: "missing instance_id",
				metadata: &domain.OrganizationMetadata{
					OrgID: orgID,
					Metadata: domain.Metadata{
						Key:   "urn:zitadel:key",
						Value: []byte("1234"),
					},
				},
				expectedErr: new(database.ForeignKeyError),
			},
			{
				name: "missing key",
				metadata: &domain.OrganizationMetadata{
					OrgID: orgID,
					Metadata: domain.Metadata{
						InstanceID: instanceID,
						Value:      []byte("1234"),
					},
				},
				expectedErr: new(database.CheckError),
			},
			{
				name: "missing value",
				metadata: &domain.OrganizationMetadata{
					OrgID: orgID,
					Metadata: domain.Metadata{
						InstanceID: instanceID,
						Key:        "urn:zitadel:key",
					},
				},
				expectedErr: new(database.NotNullError),
			},
		} {
			t.Run(testCase.name, func(t *testing.T) {
				savepoint, err := tx.Begin(t.Context())
				require.NoError(t, err)
				defer func() {
					require.NoError(t, savepoint.Rollback(t.Context()))
				}()

				metadataRepo := repository.OrganizationRepository(savepoint).Metadata(false)
				err = metadataRepo.Set(t.Context(), testCase.metadata)
				assert.ErrorIs(t, err, testCase.expectedErr)
			})
		}
	})
}

func TestRemoveOrganizationMetadata(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	// create instance
	instanceID := gofakeit.UUID()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = repository.InstanceRepository(tx).Create(t.Context(), &instance)
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository(tx)
	metadataRepo := orgRepo.Metadata(false)

	// create organization
	orgA := domain.Organization{
		ID:         "1",
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), &orgA)
	require.NoError(t, err)

	err = metadataRepo.Set(
		t.Context(),
		&domain.OrganizationMetadata{
			OrgID: orgA.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key",
				Value:      []byte("1234"),
			},
		},
		&domain.OrganizationMetadata{
			OrgID: orgA.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key2",
				Value:      []byte("asdf"),
			},
		},
	)
	require.NoError(t, err)

	orgB := domain.Organization{
		ID:         "2",
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), &orgB)
	require.NoError(t, err)

	err = metadataRepo.Set(
		t.Context(),
		&domain.OrganizationMetadata{
			OrgID: orgB.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key",
				Value:      []byte("5678"),
			},
		},
		&domain.OrganizationMetadata{
			OrgID: orgB.ID,
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "urn:zitadel:key2",
				Value:      []byte(`"asdf"`),
			},
		},
	)
	require.NoError(t, err)

	t.Run("without instance condition", func(t *testing.T) {
		affected, err := metadataRepo.Remove(
			t.Context(),
			metadataRepo.OrgIDCondition(orgA.ID),
		)
		assert.ErrorIs(t, err, new(database.MissingConditionError))
		assert.Equal(t, int64(0), affected)
	})

	t.Run("without org condition", func(t *testing.T) {
		affected, err := metadataRepo.Remove(
			t.Context(),
			metadataRepo.InstanceIDCondition("non-existing"),
		)
		assert.ErrorIs(t, err, new(database.MissingConditionError))
		assert.Equal(t, int64(0), affected)
	})

	t.Run("successful", func(t *testing.T) {
		savepoint, err := tx.Begin(t.Context())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, savepoint.Rollback(t.Context()))
		}()
		metadataRepo := repository.OrganizationRepository(savepoint).Metadata(false)
		affected, err := metadataRepo.Remove(
			t.Context(),
			database.And(
				metadataRepo.InstanceIDCondition(instanceID),
				metadataRepo.OrgIDCondition(orgA.ID),
			),
		)
		require.NoError(t, err)
		assert.Equal(t, int64(2), affected)
	})
}

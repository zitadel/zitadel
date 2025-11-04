package repository_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

// NOTE there are 2 kind of create functions for settings;
// 1. Create() which is used in projections
// 2. Create.*() which are used to create specific types of settings
func TestCreateGenericSetting(t *testing.T) {
	now := time.Now().Round(time.Microsecond)
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.Setting
		setting  domain.Setting
		err      error
	}

	// TESTS
	tests := []test{
		{
			name: "happy path",
			setting: domain.Setting{
				InstanceID:     instanceId,
				OrganizationID: &orgId,
				ID:             gofakeit.Name(),
				Type:           domain.SettingTypeLogin,
				OwnerType:      domain.OwnerTypeInstance,
				Settings:       []byte("{}"),
				CreatedAt:      now,
				UpdatedAt:      &now,
			},
		},
		{
			name: "happy path, no org",
			setting: domain.Setting{
				InstanceID: instanceId,
				// OrgID:      &orgId,
				ID:        gofakeit.Name(),
				Type:      domain.SettingTypeLogin,
				OwnerType: domain.OwnerTypeInstance,
				Settings:  []byte("{}"),
				CreatedAt: now,
				UpdatedAt: &now,
			},
		},
		{
			name: "adding setting with same instance id, org id twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
				settingRepo := repository.SettingsRepository()

				setting := domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: &orgId,
					Type:           domain.SettingTypePasswordExpiry,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
					CreatedAt:      now,
					UpdatedAt:      &now,
				}

				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			err: new(database.UniqueError),
		},
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same org (org on different instance), type, different instance",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
					// create instance
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository()
					setting := domain.Setting{
						InstanceID:     newInstId,
						OrganizationID: &newOrgId,
						Type:           domain.SettingTypeLockout,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
						CreatedAt:      now,
						UpdatedAt:      &now,
					}

					err = settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.InstanceID = instanceId
					return &setting
				},
				err: new(database.IntegrityViolationError),
			}
		}(),
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same instance, type, different org (org on different instance)",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
					// create instance
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository()
					setting := domain.Setting{
						InstanceID:     newInstId,
						OrganizationID: &newOrgId,
						Type:           domain.SettingTypeLockout,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
						CreatedAt:      now,
						UpdatedAt:      &now,
					}

					err = settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.OrganizationID = &orgId
					return &setting
				},
				setting: domain.Setting{
					InstanceID:     newInstId,
					OrganizationID: &orgId,
					Type:           domain.SettingTypePasswordExpiry,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
				},
				err: new(database.IntegrityViolationError),
			}
		}(),
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same instance, type different org (org on same instance)",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
					// create instance
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository()
					setting := domain.Setting{
						InstanceID:     newInstId,
						OrganizationID: &org.ID,
						Type:           domain.SettingTypePasswordExpiry,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
						CreatedAt:      now,
						UpdatedAt:      &now,
					}

					err = settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					// create another org
					org = domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// change the org id
					setting.OrganizationID = &newOrgId
					return &setting
				},
				setting: domain.Setting{
					InstanceID:     newInstId,
					OrganizationID: &newOrgId,
					Type:           domain.SettingTypePasswordExpiry,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
					CreatedAt:      now,
					UpdatedAt:      &now,
				},
			}
		}(),
		{
			name: "adding setting with no type",
			setting: domain.Setting{
				InstanceID:     instanceId,
				OrganizationID: &orgId,
				// Type:     domain.SettingTypeBranding,
				OwnerType: domain.OwnerTypeInstance,
				Settings:  []byte("{}"),
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding setting with no instance id",
			setting: domain.Setting{
				// InstanceID:        instanceId,
				OrganizationID: &orgId,
				ID:             gofakeit.Name(),
				Type:           domain.SettingTypeLockout,
				OwnerType:      domain.OwnerTypeInstance,
				Settings:       []byte("{}"),
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding idp with non existent instance id",
			setting: domain.Setting{
				InstanceID:     gofakeit.Name(),
				OrganizationID: &orgId,
				ID:             gofakeit.Name(),
				Type:           domain.SettingTypeLockout,
				OwnerType:      domain.OwnerTypeInstance,
				Settings:       []byte("{}"),
				CreatedAt:      now,
				UpdatedAt:      &now,
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.Setting{
				InstanceID:     instanceId,
				OrganizationID: gu.Ptr(gofakeit.Name()),
				ID:             gofakeit.Name(),
				Type:           domain.SettingTypeLockout,
				OwnerType:      domain.OwnerTypeInstance,
				Settings:       []byte("{}"),
				CreatedAt:      now,
				UpdatedAt:      &now,
			},
			err: new(database.ForeignKeyError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.Setting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			} else {
				setting = &tt.setting
			}
			settingRepo := repository.SettingsRepository()

			// create setting
			err = settingRepo.Create(t.Context(), tx, setting)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}

			// check organization values
			setting, err = settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
					),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, tt.setting.InstanceID, setting.InstanceID)
			assert.Equal(t, tt.setting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, tt.setting.Type, setting.Type)
			assert.Equal(t, tt.setting.OwnerType, setting.OwnerType)

			assert.Equal(t, setting.CreatedAt.UTC(), now.UTC())
			assert.Equal(t, setting.UpdatedAt.UTC(), now.UTC())
		})
	}
}

// NOTE all Set.*() functions are just a wrapper around repository.createSetting()
// for the sake of testing, CreatePasswordExpiry was used, but the underlying code is the same
// across all Create.*() functions
func TestSetSpecificSetting(t *testing.T) {
	beforeCreate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting
		setting  domain.PasswordExpirySetting
		err      error
	}

	// TESTS
	tests := []test{
		{
			name: "happy path",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: &orgId,
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypePasswordExpiry,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
				},
			},
		},
		{
			name: "adding setting with same instance id, org id twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting {
				settingRepo := repository.PasswordExpiryRepository()

				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						ID:             gofakeit.Name(),
						Type:           domain.SettingTypePasswordExpiry,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
					},
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: &orgId,
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypePasswordExpiry,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
				},
			},
		},
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same org (org on different instance), type, different instance",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting {
					// create instance
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.PasswordExpiryRepository()
					setting := domain.PasswordExpirySetting{
						Setting: &domain.Setting{
							InstanceID:     newInstId,
							OrganizationID: &newOrgId,
							Type:           domain.SettingTypeLockout,
							OwnerType:      domain.OwnerTypeInstance,
							Settings:       []byte("{}"),
						},
					}

					err = settingRepo.Set(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.InstanceID = instanceId
					return &setting
				},
				setting: domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &newOrgId,
						Type:           domain.SettingTypePasswordExpiry,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
					},
				},
				err: new(database.IntegrityViolationError),
			}
		}(),
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same instance, type, different org (org on different instance)",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting {
					// create instance
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.PasswordExpiryRepository()
					setting := domain.PasswordExpirySetting{
						Setting: &domain.Setting{
							InstanceID:     newInstId,
							OrganizationID: &newOrgId,
							Type:           domain.SettingTypeLockout,
							OwnerType:      domain.OwnerTypeInstance,
							Settings:       []byte("{}"),
						},
					}

					err = settingRepo.Set(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.OrganizationID = &orgId
					return &setting
				},
				setting: domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID:     newInstId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypePasswordExpiry,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
					},
				},
				err: new(database.IntegrityViolationError),
			}
		}(),
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same instance, type different org (org on same instance)",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting {
					// create instance
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.PasswordExpiryRepository()
					setting := domain.PasswordExpirySetting{
						Setting: &domain.Setting{
							InstanceID:     newInstId,
							OrganizationID: &org.ID,
							Type:           domain.SettingTypePasswordExpiry,
							OwnerType:      domain.OwnerTypeInstance,
							Settings:       []byte("{}"),
						},
					}

					err = settingRepo.Set(t.Context(), tx, &setting)
					require.NoError(t, err)

					// create another org
					org = domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// change the org id
					setting.OrganizationID = &newOrgId
					return &setting
				},
				setting: domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID:     newInstId,
						OrganizationID: &newOrgId,
						Type:           domain.SettingTypePasswordExpiry,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
					},
				},
			}
		}(),
		{
			name: "adding idp with non existent instance id",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID:     gofakeit.Name(),
					OrganizationID: &orgId,
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: gu.Ptr(gofakeit.Name()),
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: gu.Ptr(gofakeit.Name()),
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "instance id not set",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					// InstanceID: instanceId,
					OrganizationID: gu.Ptr(gofakeit.Name()),
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
				},
			},
			err: repository.ErrMissingInstanceID,
		},
		{
			name: "instance id not set",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					// OrgID:     gu.Ptr(gofakeit.Name()),
					ID:        gofakeit.Name(),
					Type:      domain.SettingTypeLockout,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
				},
			},
			err: repository.ErrMissingOrgID,
		},
		{
			name: "owner type not set",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: gu.Ptr(gofakeit.Name()),
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					// OwnerType: domain.OwnerTypeInstance,
					Settings: []byte("{}"),
				},
			},
			err: repository.ErrMissingOwnerType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.PasswordExpirySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			} else {
				setting = &tt.setting
			}
			settingRepo := repository.PasswordExpiryRepository()

			// create setting
			err = settingRepo.Set(t.Context(), tx, setting)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check organization values
			setting, err = settingRepo.Get(
				t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
					),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.setting.InstanceID, setting.InstanceID)
			assert.Equal(t, tt.setting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, tt.setting.Type, setting.Type)
			assert.Equal(t, tt.setting.OwnerType, setting.OwnerType)
			assert.WithinRange(t, setting.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, *setting.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

// Label has different indexes, so needs to have its own test case for Set
func TestSetLabelSetting(t *testing.T) {
	beforeCreate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting
		setting  domain.LabelSetting
		err      error
	}

	// TESTS
	tests := []test{
		{
			name: "happy path",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: &orgId,
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLabel,
					OwnerType:      domain.OwnerTypeInstance,
					LabelState:     gu.Ptr(domain.LabelStatePreview),
					Settings:       []byte("{}"),
				},
			},
		},
		{
			name: "adding setting with same instance id, org id twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				settingRepo := repository.LabelRepository()

				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						ID:             gofakeit.Name(),
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: &orgId,
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLabel,
					OwnerType:      domain.OwnerTypeInstance,
					LabelState:     gu.Ptr(domain.LabelStatePreview),
					Settings:       []byte("{}"),
				},
			},
		},
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same org (org on different instance), type, different instance",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
					// create instance
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.LabelRepository()
					setting := domain.LabelSetting{
						Setting: &domain.Setting{
							InstanceID:     newInstId,
							OrganizationID: &newOrgId,
							Type:           domain.SettingTypeLockout,
							OwnerType:      domain.OwnerTypeInstance,
							Settings:       []byte("{}"),
						},
					}

					err = settingRepo.Set(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.InstanceID = instanceId
					return &setting
				},
				setting: domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &newOrgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
				},
				err: new(database.IntegrityViolationError),
			}
		}(),
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same instance, type, different org (org on different instance)",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
					// create instance
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.LabelRepository()
					setting := domain.LabelSetting{
						Setting: &domain.Setting{
							InstanceID:     newInstId,
							OrganizationID: &newOrgId,
							Type:           domain.SettingTypeLockout,
							OwnerType:      domain.OwnerTypeInstance,
							LabelState:     gu.Ptr(domain.LabelStatePreview),
							Settings:       []byte("{}"),
						},
					}

					err = settingRepo.Set(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.OrganizationID = &orgId
					return &setting
				},
				setting: domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID:     newInstId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
				},
				err: new(database.IntegrityViolationError),
			}
		}(),
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same instance, type different org (org on same instance)",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
					// create instance
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.LabelRepository()
					setting := domain.LabelSetting{
						Setting: &domain.Setting{
							InstanceID:     newInstId,
							OrganizationID: &org.ID,
							Type:           domain.SettingTypeLabel,
							OwnerType:      domain.OwnerTypeInstance,
							LabelState:     gu.Ptr(domain.LabelStatePreview),
							Settings:       []byte("{}"),
						},
					}

					err = settingRepo.Set(t.Context(), tx, &setting)
					require.NoError(t, err)

					// create another org
					org = domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// change the org id
					setting.OrganizationID = &newOrgId
					return &setting
				},
				setting: domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID:     newInstId,
						OrganizationID: &newOrgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
				},
			}
		}(),
		{
			name: "adding idp with non existent instance id",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID:     gofakeit.Name(),
					OrganizationID: &orgId,
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					OwnerType:      domain.OwnerTypeInstance,
					LabelState:     gu.Ptr(domain.LabelStatePreview),
					Settings:       []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: gu.Ptr(gofakeit.Name()),
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					OwnerType:      domain.OwnerTypeInstance,
					LabelState:     gu.Ptr(domain.LabelStatePreview),
					Settings:       []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: gu.Ptr(gofakeit.Name()),
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					OwnerType:      domain.OwnerTypeInstance,
					LabelState:     gu.Ptr(domain.LabelStatePreview),
					Settings:       []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "instance id not set",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					// InstanceID: instanceId,
					OrganizationID: gu.Ptr(gofakeit.Name()),
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
				},
			},
			err: repository.ErrMissingInstanceID,
		},
		{
			name: "org id not set",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					// OrgID:     gu.Ptr(gofakeit.Name()),
					ID:        gofakeit.Name(),
					Type:      domain.SettingTypeLockout,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
				},
			},
			err: repository.ErrMissingOrgID,
		},
		{
			name: "owner type not set",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: gu.Ptr(gofakeit.Name()),
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeLockout,
					// OwnerType: domain.OwnerTypeInstance,
					Settings: []byte("{}"),
				},
			},
			err: repository.ErrMissingOwnerType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.LabelSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			} else {
				setting = &tt.setting
			}
			settingRepo := repository.LabelRepository()

			// create setting
			err = settingRepo.Set(t.Context(), tx, setting)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			setting, err = settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.setting.InstanceID, setting.InstanceID)
			assert.Equal(t, tt.setting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, tt.setting.Type, setting.Type)
			assert.Equal(t, tt.setting.OwnerType, setting.OwnerType)
			assert.WithinRange(t, setting.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, *setting.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

func TestSetSettingsError(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	tests := []struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) error
		err      error
	}{
		{
			name: "create login no settings error",
			testFunc: func(t *testing.T, tx database.QueryExecutor) error {
				settingRepo := repository.LoginRepository()
				err := settingRepo.Set(t.Context(), tx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create label no settings error",
			testFunc: func(t *testing.T, tx database.QueryExecutor) error {
				settingRepo := repository.LabelRepository()
				err := settingRepo.Set(t.Context(), tx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create password complexity no settings error",
			testFunc: func(t *testing.T, tx database.QueryExecutor) error {
				settingRepo := repository.PasswordComplexityRepository()
				err := settingRepo.Set(t.Context(), tx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create password expiry no settings error",
			testFunc: func(t *testing.T, tx database.QueryExecutor) error {
				settingRepo := repository.PasswordComplexityRepository()
				err := settingRepo.Set(t.Context(), tx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create lockout no settings error",
			testFunc: func(t *testing.T, tx database.QueryExecutor) error {
				settingRepo := repository.LabelRepository()
				err := settingRepo.Set(t.Context(), tx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create security no settings error",
			testFunc: func(t *testing.T, tx database.QueryExecutor) error {
				settingRepo := repository.SecurityRepository()
				err := settingRepo.Set(t.Context(), tx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create domain no settings error",
			testFunc: func(t *testing.T, tx database.QueryExecutor) error {
				settingRepo := repository.DomainRepository()
				err := settingRepo.Set(t.Context(), tx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create organization no settings error",
			testFunc: func(t *testing.T, tx database.QueryExecutor) error {
				settingRepo := repository.OrganizationSettingRepository()
				err := settingRepo.Set(t.Context(), tx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()
			err = tt.testFunc(t, tx)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

// NOTE all updated functions are just a wrapper around repository.Update()
// for the sake of testing, UpdateLabel was used, but the underlying code is the same
// across all Update.*() functions
func TestUpdateSetting(t *testing.T) {
	beforeUpdate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.OrganizationSettingRepository()

	tests := []struct {
		name         string
		testFunc     func(t *testing.T, tx database.QueryExecutor) *domain.OrganizationSetting
		conditions   database.Condition
		changes      database.Changes
		rowsAffected int64
		err          error
	}{
		{
			name: "missing instance id condition",
			conditions: database.And(
				// settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.InstanceIDColumn()),
		},
		{
			name: "missing org id condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				// settingRepo.OrgIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.OrganizationIDColumn()),
		},
		{
			name: "missing type condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				// settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.TypeColumn()),
		},
		{
			name: "missing owner type condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeOrganization),
				// settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.OwnerTypeColumn()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.OrganizationSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			// update setting
			_, err = settingRepo.Update(
				t.Context(), tx,
				tt.conditions,
				tt.changes,
			)
			afterUpdate := time.Now()
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			require.NoError(t, err)

			assert.ErrorIs(t, err, tt.err)

			updatedSetting, err := settingRepo.Get(
				t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, setting.InstanceID, updatedSetting.InstanceID)
			assert.Equal(t, setting.OrganizationID, updatedSetting.OrganizationID)
			assert.Equal(t, setting.ID, updatedSetting.ID)
			assert.Equal(t, setting.Type, updatedSetting.Type)

			assert.Equal(t, setting.Settings, updatedSetting.Settings)

			assert.WithinRange(t, *updatedSetting.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

// TODO
// NOTE all the Get.*() functions are a wrapper around Get()
func TestGetSetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	now := time.Now()
	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	// create setting
	// this setting is created as an additional org which should NOT
	// be returned in the results of the tests
	preExistingSetting := domain.Setting{
		InstanceID:     instanceId,
		OrganizationID: &orgId,
		ID:             gofakeit.Name(),
		Type:           domain.SettingTypePasswordExpiry,
		OwnerType:      domain.OwnerTypeInstance,
		Settings:       []byte("{}"),
		CreatedAt:      now,
		UpdatedAt:      &now,
	}
	settingRepo := repository.SettingsRepository()
	err = settingRepo.Create(t.Context(), tx, &preExistingSetting)
	require.NoError(t, err)

	type test struct {
		name       string
		testFunc   func(t *testing.T, tx database.QueryExecutor) *domain.Setting
		conditions database.Condition
		err        error
	}

	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
				setting := domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: &orgId,
					Type:           domain.SettingTypeLogin,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
					CreatedAt:      now,
					UpdatedAt:      &now,
				}

				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLogin),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance)),
		},
		{
			name: "no org instance",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
				setting := domain.Setting{
					InstanceID: instanceId,
					// OrgID:      &orgId,
					Type:      domain.SettingTypeLogin,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
				}

				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			conditions: database.And(
				// settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLogin),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance)),
			err: database.NewMissingConditionError(settingRepo.InstanceIDColumn()),
		},
		{
			name: "no org condition",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
				setting := domain.Setting{
					InstanceID: instanceId,
					// OrgID:      &orgId,
					Type:      domain.SettingTypeLogin,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
				}

				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				// settingRepo.OrgIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLogin),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance)),
			err: database.NewMissingConditionError(settingRepo.OrganizationIDColumn()),
		},
		{
			name: "no type condition",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
				setting := domain.Setting{
					InstanceID: instanceId,
					// OrgID:      &orgId,
					Type:      domain.SettingTypeLogin,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
				}

				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				// settingRepo.TypeCondition(domain.SettingTypeLogin),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance)),
			err: database.NewMissingConditionError(settingRepo.TypeColumn()),
		},
		{
			name: "no type condition",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
				setting := domain.Setting{
					InstanceID: instanceId,
					// OrgID:      &orgId,
					Type:      domain.SettingTypeLogin,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
				}

				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLogin)),
			// settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance)),
			err: database.NewMissingConditionError(settingRepo.OwnerTypeColumn()),
		},
		{
			name: "get using non existent id",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
				setting := domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: &orgId,
					ID:             gofakeit.Name(),

					Type:      domain.SettingTypeDomain,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
				}

				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLockout),
				settingRepo.IDCondition("non-existent-id"),

				settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance)),
			err: new(database.NoRowFoundError),
		},
		func() test {
			return test{
				name: "non existent orgID",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
					setting := domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypePasswordComplexity,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
						CreatedAt:      now,
						UpdatedAt:      &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)
					setting.OrganizationID = gu.Ptr("non-existent-orgID")
					return &setting
				},
				conditions: database.And(
					settingRepo.InstanceIDCondition(instanceId),
					settingRepo.OrganizationIDCondition(&orgId),
					settingRepo.TypeCondition(domain.SettingTypeLockout),
					settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance)),
				err: new(database.NoRowFoundError),
			}
		}(),
		func() test {
			return test{
				name: "non existent instanceID",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
					setting := domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLockout,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
						CreatedAt:      now,
						UpdatedAt:      &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)
					setting.InstanceID = "non-existent-instanceID"
					return &setting
				},
				conditions: database.And(
					settingRepo.InstanceIDCondition("non-existent-instanceID"),
					settingRepo.OrganizationIDCondition(&orgId),
					settingRepo.TypeCondition(domain.SettingTypeLockout),
					settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance)),
				err: new(database.NoRowFoundError),
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.Setting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						tt.conditions,
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

// -------------------------------------------------------------
// Get Specific Settings
// -------------------------------------------------------------

// testing that values written to the db are successfully retrieved
func TestCreateGetLoginPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.LoginRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.LoginSetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LoginSetting {
				setting := domain.LoginSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					AllowUserNamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "wwww.example.com",
					PasswordCheckLifetime:      time.Second * 50,
					ExternalLoginCheckLifetime: time.Second * 50,
					MFAInitSkipLifetime:        time.Second * 50,
					SecondFactorCheckLifetime:  time.Second * 50,
					MultiFactorCheckLifetime:   time.Second * 50,
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.LoginSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

// testing that values written to the db are successfully retrieved
// label uses a different code path than all the other Get.*() functions,
// so it has additional tests around missing conditions
func TestCreateGetLabelPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.LabelRepository()

	type test struct {
		name       string
		testFunc   func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting
		conditions database.Condition
		err        error
	}

	tests := []test{
		{
			name: "happy path, state preview",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeOrganization,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
					},
					PrimaryColor:        "value",
					BackgroundColor:     "value",
					WarnColor:           "value",
					FontColor:           "value",
					PrimaryColorDark:    "value",
					BackgroundColorDark: "value",
					WarnColorDark:       "value",
					FontColorDark:       "value",
					HideLoginNameSuffix: true,
					ErrorMsgPopup:       true,
					DisableWatermark:    true,
					ThemeMode:           domain.LabelPolicyThemeAuto,
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
		},
		{
			name: "missing instance id condition",
			conditions: database.And(
				// settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.InstanceIDColumn()),
		},
		{
			name: "missing org id condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				// settingRepo.OrgIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.OrganizationIDColumn()),
		},
		{
			name: "missing type condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				// settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.TypeColumn()),
		},
		{
			name: "missing owner type condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				// settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.OwnerTypeColumn()),
		},
		{
			name: "missing owner label state condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization)),
			// settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.LabelStateColumn()),
		},
		{
			name: "happy path, state activated",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeOrganization,
						LabelState:     gu.Ptr(domain.LabelStateActivated),
					},
					PrimaryColor:        "value",
					BackgroundColor:     "value",
					WarnColor:           "value",
					FontColor:           "value",
					PrimaryColorDark:    "value",
					BackgroundColorDark: "value",
					WarnColorDark:       "value",
					FontColorDark:       "value",
					HideLoginNameSuffix: true,
					ErrorMsgPopup:       true,
					DisableWatermark:    true,
					ThemeMode:           domain.LabelPolicyThemeAuto,
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStateActivated)),
		},
		{
			name: "get label policy using wrong state",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStateActivated),
					},
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.LabelState = gu.Ptr(domain.LabelStatePreview)
				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: new(database.NoRowFoundError),
		},
		{
			name: "get non existent label policy",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStateActivated),
						Settings:       []byte("{}"),
					},
				}

				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: new(database.NoRowFoundError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.LabelSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						tt.conditions,
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.PrimaryColor, setting.PrimaryColor)
			assert.Equal(t, returnedSetting.BackgroundColor, setting.BackgroundColor)
			assert.Equal(t, returnedSetting.WarnColor, setting.WarnColor)
			assert.Equal(t, returnedSetting.FontColor, setting.FontColor)
			assert.Equal(t, returnedSetting.PrimaryColorDark, setting.PrimaryColorDark)
			assert.Equal(t, returnedSetting.BackgroundColorDark, setting.BackgroundColorDark)
			assert.Equal(t, returnedSetting.WarnColorDark, setting.WarnColorDark)
			assert.Equal(t, returnedSetting.FontColorDark, setting.FontColorDark)
			assert.Equal(t, returnedSetting.HideLoginNameSuffix, setting.HideLoginNameSuffix)
			assert.Equal(t, returnedSetting.ErrorMsgPopup, setting.ErrorMsgPopup)
			assert.Equal(t, returnedSetting.DisableWatermark, setting.DisableWatermark)
			assert.Equal(t, returnedSetting.ThemeMode, setting.ThemeMode)
		})
	}
}

// testing that values written to the db are successfully retrieved
func TestCreateGetPasswordComplexityPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.PasswordComplexityRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.PasswordComplexitySetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.PasswordComplexitySetting {
				setting := domain.PasswordComplexitySetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					MinLength:    89,
					HasLowercase: true,
					HasUppercase: true,
					HasNumber:    true,
					HasSymbol:    true,
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.PasswordComplexitySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

func TestCreateGetPasswordExpiryPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.PasswordExpiryRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting {
				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					ExpireWarnDays: 30,
					MaxAgeDays:     30,
				}

				err = settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.PasswordExpirySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

func TestCreateGetLockoutPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.LockoutRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.LockoutSetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LockoutSetting {
				setting := domain.LockoutSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					MaxPasswordAttempts: 50,
					MaxOTPAttempts:      50,
					ShowLockOutFailures: true,
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.LockoutSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

func TestCreateGetSecurityPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.SecurityRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.SecuritySetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.SecuritySetting {
				setting := domain.SecuritySetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					Enabled:               true,
					EnableIframeEmbedding: true,
					AllowedOrigins:        []string{"value"},
					EnableImpersonation:   true,
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var setting *domain.SecuritySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

func TestCreateGetDomainPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.DomainRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.DomainSetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.DomainSetting {
				setting := domain.DomainSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					UserLoginMustBeDomain:                  true,
					ValidateOrgDomains:                     true,
					SMTPSenderAddressMatchesInstanceDomain: true,
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var setting *domain.DomainSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

func TestCreateGetOrgPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.OrganizationSettingRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.OrganizationSetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.OrganizationSetting {
				setting := domain.OrganizationSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					OrganizationScopedUsernames: true,
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var setting *domain.OrganizationSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

// -------------------------------------------------------------
// Update Specific Settings
// -------------------------------------------------------------

func TestCreateUpdateLoginPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.LoginRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.LoginSetting
		changes  database.Changes
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LoginSetting {
				setting := domain.LoginSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					AllowUserNamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "wwww.example.com",
					PasswordCheckLifetime:      time.Second * 50,
					ExternalLoginCheckLifetime: time.Second * 50,
					MFAInitSkipLifetime:        time.Second * 50,
					SecondFactorCheckLifetime:  time.Second * 50,
					MultiFactorCheckLifetime:   time.Second * 50,
					MFAType:                    []domain.MultiFactorType{domain.MultiFactorTypeUnspecified},
					SecondFactorTypes:          []domain.SecondFactorType{domain.SecondFactorTypeOTPEmail},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.AllowUserNamePassword = false
				setting.AllowRegister = false
				setting.AllowExternalIDP = false
				setting.ForceMFA = false
				setting.ForceMFALocalOnly = false
				setting.HidePasswordReset = false
				setting.IgnoreUnknownUsernames = false
				setting.AllowDomainDiscovery = false
				setting.DisableLoginWithEmail = false
				setting.DisableLoginWithPhone = false
				setting.PasswordlessType = domain.PasswordlessTypeNotAllowed
				setting.DefaultRedirectURI = "www.not-example.com"
				setting.PasswordCheckLifetime = time.Minute
				setting.ExternalLoginCheckLifetime = time.Minute
				setting.MFAInitSkipLifetime = time.Minute
				setting.SecondFactorCheckLifetime = time.Minute
				setting.MultiFactorCheckLifetime = time.Minute
				setting.MFAType = []domain.MultiFactorType{domain.MultiFactorTypeUnspecified, domain.MultiFactorTypeU2FWithPIN}
				setting.SecondFactorTypes = []domain.SecondFactorType{domain.SecondFactorTypeOTPEmail, domain.SecondFactorTypeOTPSMS}

				return &setting
			},
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetAllowUserNamePasswordField(false),
					settingRepo.SetAllowDomainDiscoveryField(false),
					settingRepo.SetAllowRegisterField(false),
					settingRepo.SetAllowExternalIDPField(false),
					settingRepo.SetForceMFAField(false),
					settingRepo.SetForceMFALocalOnlyField(false),
					settingRepo.SetHidePasswordResetField(false),
					settingRepo.SetIgnoreUnknownUsernamesField(false),
					settingRepo.SetAllowDomainDiscoveryField(false),
					settingRepo.SetDisableLoginWithEmailField(false),
					settingRepo.SetDisableLoginWithPhoneField(false),
					settingRepo.SetPasswordlessTypeField(domain.PasswordlessTypeNotAllowed),
					settingRepo.SetDefaultRedirectURIField("www.not-example.com"),
					settingRepo.SetPasswordCheckLifetimeField(time.Minute),
					settingRepo.SetExternalLoginCheckLifetimeField(time.Minute),
					settingRepo.SetMFAInitSkipLifetimeField(time.Minute),
					settingRepo.SetSecondFactorCheckLifetimeField(time.Minute),
					settingRepo.SetMultiFactorCheckLifetimeField(time.Minute),
					settingRepo.AddMFAType(domain.MultiFactorTypeU2FWithPIN),
					settingRepo.AddSecondFactorTypesField(domain.SecondFactorTypeOTPSMS),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.LoginSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			_, err = settingRepo.Update(t.Context(), tx,
				database.And(
					settingRepo.InstanceIDCondition(setting.InstanceID),
					settingRepo.OrganizationIDCondition(setting.OrganizationID),
					settingRepo.TypeCondition(setting.Type),
					settingRepo.OwnerTypeCondition(setting.OwnerType),
					settingRepo.LabelStateCondition(*setting.LabelState),
				),
				tt.changes,
			)
			require.NoError(t, err)

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, setting.InstanceID, returnedSetting.InstanceID)
			assert.Equal(t, setting.OrganizationID, returnedSetting.OrganizationID)
			assert.Equal(t, setting.ID, returnedSetting.ID)
			assert.Equal(t, setting.Type, returnedSetting.Type)

			assert.Equal(t, setting.AllowUserNamePassword, returnedSetting.AllowUserNamePassword)
			assert.Equal(t, setting.AllowRegister, returnedSetting.AllowUserNamePassword)
			assert.Equal(t, setting.AllowExternalIDP, returnedSetting.AllowExternalIDP)
			assert.Equal(t, setting.ForceMFA, returnedSetting.ForceMFA)
			assert.Equal(t, setting.ForceMFALocalOnly, returnedSetting.ForceMFALocalOnly)
			assert.Equal(t, setting.HidePasswordReset, returnedSetting.HidePasswordReset)
			assert.Equal(t, setting.IgnoreUnknownUsernames, returnedSetting.IgnoreUnknownUsernames)
			assert.Equal(t, setting.AllowDomainDiscovery, returnedSetting.AllowDomainDiscovery)
			assert.Equal(t, setting.DisableLoginWithEmail, returnedSetting.DisableLoginWithEmail)
			assert.Equal(t, setting.DisableLoginWithPhone, returnedSetting.DisableLoginWithPhone)
			assert.Equal(t, setting.PasswordlessType, returnedSetting.PasswordlessType)
			assert.Equal(t, setting.DefaultRedirectURI, returnedSetting.DefaultRedirectURI)
			assert.Equal(t, setting.PasswordCheckLifetime, returnedSetting.PasswordCheckLifetime)
			assert.Equal(t, setting.ExternalLoginCheckLifetime, returnedSetting.ExternalLoginCheckLifetime)
			assert.Equal(t, setting.MFAInitSkipLifetime, returnedSetting.MFAInitSkipLifetime)
			assert.Equal(t, setting.SecondFactorCheckLifetime, returnedSetting.SecondFactorCheckLifetime)
			assert.Equal(t, setting.MultiFactorCheckLifetime, returnedSetting.MultiFactorCheckLifetime)
		})
	}
}

func TestCreateUpdateLabelPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.LabelRepository()

	type test struct {
		name       string
		testFunc   func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting
		conditions database.Condition
		changes    database.Changes
		err        error
	}

	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeOrganization,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
					},
					PrimaryColor:        "value",
					BackgroundColor:     "value",
					WarnColor:           "value",
					FontColor:           "value",
					PrimaryColorDark:    "value",
					BackgroundColorDark: "value",
					WarnColorDark:       "value",
					FontColorDark:       "value",
					HideLoginNameSuffix: true,
					ErrorMsgPopup:       true,
					DisableWatermark:    true,
					ThemeMode:           domain.LabelPolicyThemeDark,
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.PrimaryColor = "new_value"
				setting.BackgroundColor = "new_value"
				setting.WarnColor = "new_value"
				setting.FontColor = "new_value"
				setting.PrimaryColorDark = "new_value"
				setting.BackgroundColorDark = "new_value"
				setting.WarnColorDark = "new_value"
				setting.FontColorDark = "new_value"
				setting.HideLoginNameSuffix = false
				setting.ErrorMsgPopup = false
				setting.DisableWatermark = false
				setting.ThemeMode = domain.LabelPolicyThemeAuto

				return &setting
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetPrimaryColorField("new_value"),
					settingRepo.SetBackgroundColorField("new_value"),
					settingRepo.SetWarnColorField("new_value"),
					settingRepo.SetFontColorField("new_value"),
					settingRepo.SetPrimaryColorDarkField("new_value"),
					settingRepo.SetBackgroundColorDarkField("new_value"),
					settingRepo.SetWarnColorDarkField("new_value"),
					settingRepo.SetFontColorDarkField("new_value"),
					settingRepo.SetHideLoginNameSuffixField(false),
					settingRepo.SetErrorMsgPopupField(false),
					settingRepo.SetDisableWatermarkField(false),
					settingRepo.SetThemeModeField(domain.LabelPolicyThemeAuto),
				),
			},
		},
		{
			name: "missing instance id condition",
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetBackgroundColorDarkField("new_value"),
				),
			},
			conditions: database.And(
				// settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.InstanceIDColumn()),
		},
		{
			name: "missing org id condition",
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetBackgroundColorDarkField("new_value"),
				),
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				// settingRepo.OrgIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.OrganizationIDColumn()),
		},
		{
			name: "missing type condition",
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetBackgroundColorDarkField("new_value"),
				),
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				// settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.TypeColumn()),
		},
		{
			name: "missing owner type condition",
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetBackgroundColorDarkField("new_value"),
				),
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				// settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.OwnerTypeColumn()),
		},
		{
			name: "missing owner label state condition",
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetBackgroundColorDarkField("new_value"),
				),
			},
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization)),
			// settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.LabelStateColumn()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.LabelSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			_, err = settingRepo.Update(t.Context(), tx,
				tt.conditions,
				tt.changes,
			)
			require.NoError(t, err)

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					tt.conditions),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.PrimaryColor, setting.PrimaryColor)
			assert.Equal(t, returnedSetting.BackgroundColor, setting.BackgroundColor)
			assert.Equal(t, returnedSetting.WarnColor, setting.WarnColor)
			assert.Equal(t, returnedSetting.FontColor, setting.FontColor)
			assert.Equal(t, returnedSetting.PrimaryColorDark, setting.PrimaryColorDark)
			assert.Equal(t, returnedSetting.BackgroundColorDark, setting.BackgroundColorDark)
			assert.Equal(t, returnedSetting.WarnColorDark, setting.WarnColorDark)
			assert.Equal(t, returnedSetting.FontColorDark, setting.FontColorDark)
			assert.Equal(t, returnedSetting.HideLoginNameSuffix, setting.HideLoginNameSuffix)
			assert.Equal(t, returnedSetting.ErrorMsgPopup, setting.ErrorMsgPopup)
			assert.Equal(t, returnedSetting.DisableWatermark, setting.DisableWatermark)
			assert.Equal(t, returnedSetting.ThemeMode, setting.ThemeMode)
		})
	}
}

func TestCreateUpdatePasswordComplexityPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.PasswordComplexityRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.PasswordComplexitySetting
		changes  database.Changes
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.PasswordComplexitySetting {
				setting := domain.PasswordComplexitySetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					MinLength:    89,
					HasLowercase: true,
					HasUppercase: true,
					HasNumber:    true,
					HasSymbol:    true,
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.MinLength = 200
				setting.HasLowercase = false
				setting.HasUppercase = false
				setting.HasNumber = false
				setting.HasSymbol = false

				return &setting
			},
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetMinLengthField(200),
					settingRepo.SetHasLowercaseField(false),
					settingRepo.SetHasUppercaseField(false),
					settingRepo.SetHasNumberField(false),
					settingRepo.SetHasSymbolField(false),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.PasswordComplexitySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			_, err = settingRepo.Update(t.Context(), tx,
				database.And(
					settingRepo.InstanceIDCondition(setting.InstanceID),
					settingRepo.OrganizationIDCondition(setting.OrganizationID),
					settingRepo.TypeCondition(setting.Type),
					settingRepo.OwnerTypeCondition(setting.OwnerType),
					settingRepo.LabelStateCondition(*setting.LabelState),
				),
				tt.changes,
			)
			require.NoError(t, err)

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.MinLength, setting.MinLength)
			assert.Equal(t, returnedSetting.HasLowercase, setting.HasLowercase)
			assert.Equal(t, returnedSetting.HasUppercase, setting.HasUppercase)
			assert.Equal(t, returnedSetting.HasNumber, setting.HasNumber)
			assert.Equal(t, returnedSetting.HasSymbol, setting.HasSymbol)
		})
	}
}

func TestCreateUpdatePasswordExpiryPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.PasswordExpiryRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting
		changes  database.Changes
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting {
				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					ExpireWarnDays: 30,
					MaxAgeDays:     30,
				}
				err = settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.ExpireWarnDays = 200
				setting.MaxAgeDays = 200

				return &setting
			},
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetExpireWarnDays(200),
					settingRepo.SetMaxAgeDays(200),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.PasswordExpirySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			_, err = settingRepo.Update(t.Context(), tx,
				database.And(
					settingRepo.InstanceIDCondition(setting.InstanceID),
					settingRepo.OrganizationIDCondition(setting.OrganizationID),
					settingRepo.TypeCondition(setting.Type),
					settingRepo.OwnerTypeCondition(setting.OwnerType),
					settingRepo.LabelStateCondition(*setting.LabelState),
				),
				tt.changes,
			)
			require.NoError(t, err)

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.ExpireWarnDays, setting.ExpireWarnDays)
			assert.Equal(t, returnedSetting.MaxAgeDays, setting.MaxAgeDays)
		})
	}
}

func TestCreateUpdateLockoutPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.LockoutRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.LockoutSetting
		changes  database.Changes
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LockoutSetting {
				setting := domain.LockoutSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					MaxPasswordAttempts: 50,
					MaxOTPAttempts:      50,
					ShowLockOutFailures: true,
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.MaxPasswordAttempts = 100
				setting.MaxOTPAttempts = 100
				setting.ShowLockOutFailures = false

				return &setting
			},
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetMaxPasswordAttempts(100),
					settingRepo.SetMaxOTPAttempts(100),
					settingRepo.SetShowLockOutFailures(false),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			var setting *domain.LockoutSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			_, err = settingRepo.Update(t.Context(), tx,
				database.And(
					settingRepo.InstanceIDCondition(setting.InstanceID),
					settingRepo.OrganizationIDCondition(setting.OrganizationID),
					settingRepo.TypeCondition(setting.Type),
					settingRepo.OwnerTypeCondition(setting.OwnerType),
					settingRepo.LabelStateCondition(*setting.LabelState),
				),
				tt.changes,
			)
			require.NoError(t, err)

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.MaxPasswordAttempts, setting.MaxPasswordAttempts)
			assert.Equal(t, returnedSetting.MaxOTPAttempts, setting.MaxOTPAttempts)
			assert.Equal(t, returnedSetting.ShowLockOutFailures, setting.ShowLockOutFailures)
		})
	}
}

func TestCreateUpdateSecurityPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.SecurityRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.SecuritySetting
		changes  database.Changes
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.SecuritySetting {
				setting := domain.SecuritySetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeSecurity,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
					},
					Enabled:               true,
					EnableIframeEmbedding: true,
					AllowedOrigins:        []string{"value"},
					EnableImpersonation:   true,
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.Enabled = false
				setting.EnableIframeEmbedding = false
				setting.AllowedOrigins = []string{"value", "new_value"}
				setting.EnableImpersonation = false

				return &setting
			},
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetEnabled(false),
					settingRepo.SetEnableIframeEmbedding(false),
					settingRepo.SetEnableImpersonation(false),
					settingRepo.AddAllowedOrigins("new_value"),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var setting *domain.SecuritySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			_, err = settingRepo.Update(t.Context(), tx,
				database.And(
					settingRepo.InstanceIDCondition(setting.InstanceID),
					settingRepo.OrganizationIDCondition(setting.OrganizationID),
					settingRepo.TypeCondition(setting.Type),
					settingRepo.OwnerTypeCondition(setting.OwnerType),
				),
				tt.changes,
			)
			require.NoError(t, err)

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Enabled, setting.Enabled)
			assert.Equal(t, returnedSetting.EnableIframeEmbedding, setting.EnableIframeEmbedding)
			assert.Equal(t, returnedSetting.AllowedOrigins, setting.AllowedOrigins)
			assert.Equal(t, returnedSetting.EnableImpersonation, setting.EnableImpersonation)
		})
	}
}

func TestCreateUpdateDomainPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.DomainRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.DomainSetting
		changes  database.Changes
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.DomainSetting {
				setting := domain.DomainSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					UserLoginMustBeDomain:                  true,
					ValidateOrgDomains:                     true,
					SMTPSenderAddressMatchesInstanceDomain: true,
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.UserLoginMustBeDomain = false
				setting.ValidateOrgDomains = false
				setting.SMTPSenderAddressMatchesInstanceDomain = false

				return &setting
			},
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetUserLoginMustBeDomain(false),
					settingRepo.SetValidateOrgDomains(false),
					settingRepo.SetSMTPSenderAddressMatchesInstanceDomain(false),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var setting *domain.DomainSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			_, err = settingRepo.Update(t.Context(), tx,
				database.And(
					settingRepo.InstanceIDCondition(setting.InstanceID),
					settingRepo.OrganizationIDCondition(setting.OrganizationID),
					settingRepo.TypeCondition(setting.Type),
					settingRepo.OwnerTypeCondition(setting.OwnerType),
					settingRepo.LabelStateCondition(*setting.LabelState),
				),
				tt.changes,
			)
			require.NoError(t, err)

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.UserLoginMustBeDomain, setting.UserLoginMustBeDomain)
			assert.Equal(t, returnedSetting.ValidateOrgDomains, setting.ValidateOrgDomains)
			assert.Equal(t, returnedSetting.SMTPSenderAddressMatchesInstanceDomain, setting.SMTPSenderAddressMatchesInstanceDomain)
		})
	}
}

func TestCreateUpdateOrgPolicySetting(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.OrganizationSettingRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.OrganizationSetting
		changes  database.Changes
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.OrganizationSetting {
				setting := domain.OrganizationSetting{
					Setting: &domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						Type:           domain.SettingTypeLabel,
						OwnerType:      domain.OwnerTypeInstance,
						LabelState:     gu.Ptr(domain.LabelStatePreview),
						Settings:       []byte("{}"),
					},
					OrganizationScopedUsernames: true,
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.OrganizationScopedUsernames = false

				return &setting
			},
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetOrganizationScopedUsernames(false),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var setting *domain.OrganizationSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(t, tx)
			}

			_, err = settingRepo.Update(t.Context(), tx,
				database.And(
					settingRepo.InstanceIDCondition(setting.InstanceID),
					settingRepo.OrganizationIDCondition(setting.OrganizationID),
					settingRepo.TypeCondition(setting.Type),
					settingRepo.OwnerTypeCondition(setting.OwnerType),
					settingRepo.LabelStateCondition(*setting.LabelState),
				),
				tt.changes,
			)
			require.NoError(t, err)

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrganizationIDCondition(setting.OrganizationID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrganizationID, setting.OrganizationID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.OrganizationScopedUsernames, setting.OrganizationScopedUsernames)
		})
	}
}

func TestListSetting(t *testing.T) {
	now := time.Now()

	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.SettingsRepository()

	type test struct {
		name               string
		testFunc           func(t *testing.T, tx database.QueryExecutor) []*domain.Setting
		conditionClauses   database.QueryOption
		noSettingsReturned bool
	}
	tests := []test{
		{
			name: "multiple settings filter on instance",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Setting {
				// create instance
				newInstanceId := gofakeit.Name()
				instance := domain.Instance{
					ID:              newInstanceId,
					Name:            gofakeit.Name(),
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}
				err = instanceRepo.Create(t.Context(), tx, &instance)
				require.NoError(t, err)

				// create org
				newOrgId := gofakeit.Name()
				org := domain.Organization{
					ID:         newOrgId,
					Name:       gofakeit.Name(),
					InstanceID: newInstanceId,
					State:      domain.OrgStateActive,
				}
				organizationRepo := repository.OrganizationRepository()
				err = organizationRepo.Create(t.Context(), tx, &org)
				require.NoError(t, err)

				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID:     newInstanceId,
					OrganizationID: &newOrgId,
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeSecurity,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
					CreatedAt:      now,
					UpdatedAt:      &now,
				}
				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)

				noOfSettings := 5
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						ID:             gofakeit.Name(),
						// Type:              domain.SettingTypeLogin,
						Type:      domain.SettingType(i + 1),
						OwnerType: domain.OwnerTypeInstance,
						Settings:  []byte("{}"),
						CreatedAt: now,
						UpdatedAt: &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
			conditionClauses: database.WithCondition(settingRepo.InstanceIDCondition(instanceId)),
		},
		{
			name: "multiple settings filter on instance, no org",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Setting {
				// create instance
				newInstanceId := gofakeit.Name()
				instance := domain.Instance{
					ID:              newInstanceId,
					Name:            gofakeit.Name(),
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}
				err = instanceRepo.Create(t.Context(), tx, &instance)
				require.NoError(t, err)

				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID: newInstanceId,
					// OrgID:      &newOrgId,
					ID:        gofakeit.Name(),
					Type:      domain.SettingTypeSecurity,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
				}
				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)

				noOfSettings := 5
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID: instanceId,
						// OrgID:      &orgId,
						ID: gofakeit.Name(),
						// Type:              domain.SettingTypeLogin,
						Type:      domain.SettingType(i + 1),
						OwnerType: domain.OwnerTypeInstance,
						Settings:  []byte("{}"),
						CreatedAt: now,
						UpdatedAt: &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
			conditionClauses: database.WithCondition(settingRepo.InstanceIDCondition(instanceId)),
		},
		{
			name: "multiple settings filter on org",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Setting {
				// create org
				newOrgId := gofakeit.Name()
				org := domain.Organization{
					ID:         newOrgId,
					Name:       gofakeit.Name(),
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				}
				organizationRepo := repository.OrganizationRepository()
				err = organizationRepo.Create(t.Context(), tx, &org)
				require.NoError(t, err)

				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: &newOrgId,
					ID:             gofakeit.Name(),
					Type:           domain.SettingTypeSecurity,
					OwnerType:      domain.OwnerTypeInstance,
					Settings:       []byte("{}"),
					CreatedAt:      now,
					UpdatedAt:      &now,
				}
				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)

				noOfSettings := 5
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						ID:             gofakeit.Name(),
						Type:           domain.SettingType(i + 1),
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
						CreatedAt:      now,
						UpdatedAt:      &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
			conditionClauses: database.WithCondition(settingRepo.OrganizationIDCondition(&orgId)),
		},
		{
			name: "happy path single setting no filter",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Setting {
				noOfSettings := 1
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						ID:             gofakeit.Name(),
						Type:           domain.SettingType(i + 1),
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
						CreatedAt:      now,
						UpdatedAt:      &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
		},
		{
			name: "happy path multiple settings no filter",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Setting {
				noOfSettings := 5
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						ID:             gofakeit.Name(),
						Type:           domain.SettingType(i + 1),
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
						CreatedAt:      now,
						UpdatedAt:      &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
		},
		{
			name: "multiple settings filter on type",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Setting {
				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID:     instanceId,
					OrganizationID: &orgId,
					ID:             gofakeit.Name(),

					Type:      domain.SettingTypeLogin,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
				}
				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)

				noOfSettings := 1
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID:     instanceId,
						OrganizationID: &orgId,
						ID:             gofakeit.Name(),
						Type:           domain.SettingTypePasswordExpiry,
						OwnerType:      domain.OwnerTypeInstance,
						Settings:       []byte("{}"),
						CreatedAt:      now,
						UpdatedAt:      &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
			conditionClauses: database.WithCondition(settingRepo.TypeCondition(domain.SettingTypePasswordExpiry)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			settings := tt.testFunc(t, tx)

			// check setting values
			var returnedSettings []*domain.Setting
			if tt.conditionClauses != nil {
				returnedSettings, err = settingRepo.List(t.Context(), tx,
					tt.conditionClauses,
				)
			} else {
				returnedSettings, err = settingRepo.List(t.Context(), tx)
			}
			require.NoError(t, err)
			if tt.noSettingsReturned {
				assert.Nil(t, returnedSettings)
				return
			}

			assert.Equal(t, len(settings), len(returnedSettings))
			for i, setting := range settings {

				assert.Equal(t, returnedSettings[i].InstanceID, setting.InstanceID)
				assert.Equal(t, returnedSettings[i].OrganizationID, setting.OrganizationID)
				assert.Equal(t, returnedSettings[i].Type, setting.Type)
				assert.Equal(t, returnedSettings[i].OwnerType, setting.OwnerType)
				assert.Equal(t, returnedSettings[i].Settings, setting.Settings)
			}
		})
	}
}

func TestDeleteSetting(t *testing.T) {
	now := time.Now()

	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	settingRepo := repository.SettingsRepository()

	type test struct {
		name            string
		testFunc        func(t *testing.T, tx database.QueryExecutor)
		conditions      database.Condition
		noOfDeletedRows int64
		err             error
	}
	tests := []test{
		func() test {
			var noOfSettings int64 = 1
			return test{
				name: "happy path delete",
				testFunc: func(t *testing.T, tx database.QueryExecutor) {
					for range noOfSettings {
						setting := domain.Setting{
							InstanceID:     instanceId,
							OrganizationID: &orgId,
							Type:           domain.SettingTypeLogin,
							OwnerType:      domain.OwnerTypeInstance,
							Settings:       []byte("{}"),
							CreatedAt:      now,
							UpdatedAt:      &now,
						}

						err := settingRepo.Create(t.Context(), tx, &setting)
						require.NoError(t, err)
					}
				},
				noOfDeletedRows: noOfSettings,
				conditions: database.And(
					settingRepo.InstanceIDCondition(instanceId),
					settingRepo.OrganizationIDCondition(&orgId),
					settingRepo.TypeCondition(domain.SettingTypeLogin),
					settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance),
				),
			}
		}(),
		func() test {
			return test{
				name: "deleted setting with wrong type",
				testFunc: func(t *testing.T, tx database.QueryExecutor) {
					noOfSettings := 1
					for range noOfSettings {
						setting := domain.Setting{
							InstanceID:     instanceId,
							OrganizationID: &orgId,
							ID:             gofakeit.Name(),
							Type:           domain.SettingTypeLogin,
							OwnerType:      domain.OwnerTypeInstance,
							Settings:       []byte("{}"),
							CreatedAt:      now,
							UpdatedAt:      &now,
						}

						err := settingRepo.Create(t.Context(), tx, &setting)
						require.NoError(t, err)
					}
				},
				noOfDeletedRows: 0,
				conditions: database.And(
					settingRepo.InstanceIDCondition(instanceId),
					settingRepo.OrganizationIDCondition(&orgId),
					// wrong type
					settingRepo.TypeCondition(domain.SettingTypePasswordComplexity),
					settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance),
				),
			}
		}(),
		func() test {
			return test{
				name: "deleted already deleted setting",
				testFunc: func(t *testing.T, tx database.QueryExecutor) {
					noOfSettings := 1
					for range noOfSettings {
						setting := domain.Setting{
							InstanceID:     instanceId,
							OrganizationID: &orgId,
							ID:             gofakeit.Name(),
							Type:           domain.SettingTypeLogin,
							OwnerType:      domain.OwnerTypeInstance,
							Settings:       []byte("{}"),
							CreatedAt:      now,
							UpdatedAt:      &now,
						}

						err := settingRepo.Create(t.Context(), tx, &setting)
						require.NoError(t, err)
					}

					affectedRows, err := settingRepo.Delete(t.Context(), tx,
						database.And(
							settingRepo.InstanceIDCondition(instanceId),
							settingRepo.OrganizationIDCondition(&orgId),
							settingRepo.TypeCondition(domain.SettingTypeLogin),
							settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance),
						),
					)
					assert.Equal(t, int64(1), affectedRows)
					require.NoError(t, err)
				},
				// this test should return 0 affected rows as the setting was already deleted
				noOfDeletedRows: 0,
				conditions: database.And(
					settingRepo.InstanceIDCondition(instanceId),
					settingRepo.OrganizationIDCondition(&orgId),
					settingRepo.TypeCondition(domain.SettingTypeLogin),
					settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance),
				),
			}
		}(),
		{
			name: "missing instance id condition",
			conditions: database.And(
				// settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization)),
			err: database.NewMissingConditionError(settingRepo.InstanceIDColumn()),
		},
		{
			name: "missing org id condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				// settingRepo.OrgIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.OrganizationIDColumn()),
		},
		{
			name: "missing type condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				// settingRepo.TypeCondition(domain.SettingTypeLabel),
				settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.TypeColumn()),
		},
		{
			name: "missing owner type condition",
			conditions: database.And(
				settingRepo.InstanceIDCondition(instanceId),
				settingRepo.OrganizationIDCondition(&orgId),
				settingRepo.TypeCondition(domain.SettingTypeLabel),
				// settingRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingRepo.LabelStateCondition(domain.LabelStatePreview)),
			err: database.NewMissingConditionError(settingRepo.OwnerTypeColumn()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			if tt.testFunc != nil {
				tt.testFunc(t, tx)
			}

			// delete setting
			noOfDeletedRows, err := settingRepo.Delete(t.Context(), tx,
				tt.conditions,
			)
			require.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check setting was deleted
			organization, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					tt.conditions,
				),
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, organization)
		})
	}
}

package repository_test

import (
	"fmt"
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
				InstanceID: instanceId,
				OrgID:      &orgId,
				ID:         gofakeit.Name(),
				Type:       domain.SettingTypeLogin,
				OwnerType:  domain.OwnerTypeInstance,
				Settings:   []byte("{}"),
				CreatedAt:  now,
				UpdatedAt:  &now,
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
					InstanceID: instanceId,
					OrgID:      &orgId,
					Type:       domain.SettingTypePasswordExpiry,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
					CreatedAt:  now,
					UpdatedAt:  &now,
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
						InstanceID: newInstId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypeLockout,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
					}

					err = settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.InstanceID = instanceId
					return &setting
				},
				// setting: domain.Setting{
				// 	InstanceID: instanceId,
				// 	OrgID:      &newOrgId,
				// 	Type:       domain.SettingTypePasswordExpiry,
				// 	OwnerType:  domain.OwnerTypeInstance,
				// 	Settings:   []byte("{}"),
				// },
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
						InstanceID: newInstId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypeLockout,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
					}

					err = settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.OrgID = &orgId
					return &setting
				},
				setting: domain.Setting{
					InstanceID: newInstId,
					OrgID:      &orgId,
					Type:       domain.SettingTypePasswordExpiry,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
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
						InstanceID: newInstId,
						OrgID:      &org.ID,
						Type:       domain.SettingTypePasswordExpiry,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
					setting.OrgID = &newOrgId
					return &setting
				},
				setting: domain.Setting{
					InstanceID: newInstId,
					OrgID:      &newOrgId,
					Type:       domain.SettingTypePasswordExpiry,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
					CreatedAt:  now,
					UpdatedAt:  &now,
				},
			}
		}(),
		{
			name: "adding setting with no type",
			setting: domain.Setting{
				InstanceID: instanceId,
				OrgID:      &orgId,
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
				OrgID:     &orgId,
				ID:        gofakeit.Name(),
				Type:      domain.SettingTypeLockout,
				OwnerType: domain.OwnerTypeInstance,
				Settings:  []byte("{}"),
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding idp with non existent instance id",
			setting: domain.Setting{
				InstanceID: gofakeit.Name(),
				OrgID:      &orgId,
				ID:         gofakeit.Name(),
				Type:       domain.SettingTypeLockout,
				OwnerType:  domain.OwnerTypeInstance,
				Settings:   []byte("{}"),
				CreatedAt:  now,
				UpdatedAt:  &now,
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.Setting{
				InstanceID: instanceId,
				OrgID:      gu.Ptr(gofakeit.Name()),
				ID:         gofakeit.Name(),
				Type:       domain.SettingTypeLockout,
				OwnerType:  domain.OwnerTypeInstance,
				Settings:   []byte("{}"),
				CreatedAt:  now,
				UpdatedAt:  &now,
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
			// beforeCreate := time.Now()
			err = settingRepo.Create(t.Context(), tx, setting)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			// afterCreate := time.Now()

			// check organization values
			setting, err = settingRepo.Get(t.Context(), tx,
				// setting.InstanceID,
				// setting.OrgID,
				// setting.Type,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrgIDCondition(setting.OrgID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
					),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, tt.setting.InstanceID, setting.InstanceID)
			assert.Equal(t, tt.setting.OrgID, setting.OrgID)
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
func TestCreateSpecificSetting(t *testing.T) {
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
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypePasswordExpiry,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
				},
			},
		},
		{
			name: "happy path, no org",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					// OrgID:      &orgId,
					ID:        gofakeit.Name(),
					Type:      domain.SettingTypePasswordExpiry,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
				},
			},
		},
		{
			name: "adding setting with same instance id, org id twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.PasswordExpirySetting {
				settingRepo := repository.PasswordExpiryRepository()

				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypePasswordExpiry,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
					},
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypePasswordExpiry,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
				},
			},
			// err: new(database.UniqueError),
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
							InstanceID: newInstId,
							OrgID:      &newOrgId,
							Type:       domain.SettingTypeLockout,
							OwnerType:  domain.OwnerTypeInstance,
							Settings:   []byte("{}"),
						},
					}

					// err = settingRepo.CreatePasswordExpiry(t.Context(), tx, &setting)
					err = settingRepo.Set(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.InstanceID = instanceId
					return &setting
				},
				setting: domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypePasswordExpiry,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
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
							InstanceID: newInstId,
							OrgID:      &newOrgId,
							Type:       domain.SettingTypeLockout,
							OwnerType:  domain.OwnerTypeInstance,
							Settings:   []byte("{}"),
						},
					}

					err = settingRepo.Set(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.OrgID = &orgId
					return &setting
				},
				setting: domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: newInstId,
						OrgID:      &orgId,
						Type:       domain.SettingTypePasswordExpiry,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
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
							InstanceID: newInstId,
							OrgID:      &org.ID,
							Type:       domain.SettingTypePasswordExpiry,
							OwnerType:  domain.OwnerTypeInstance,
							Settings:   []byte("{}"),
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
					setting.OrgID = &newOrgId
					return &setting
				},
				setting: domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: newInstId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypePasswordExpiry,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
					},
				},
			}
		}(),
		{
			name: "adding setting with no instance id",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					// InstanceID:        instanceId,
					OrgID:     &orgId,
					ID:        gofakeit.Name(),
					Type:      domain.SettingTypeLockout,
					OwnerType: domain.OwnerTypeInstance,
					Settings:  []byte("{}"),
				},
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding idp with non existent instance id",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID: gofakeit.Name(),
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLockout,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					OrgID:      gu.Ptr(gofakeit.Name()),
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLockout,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.PasswordExpirySetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					OrgID:      gu.Ptr(gofakeit.Name()),
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLockout,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
				},
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
				// setting.InstanceID,
				// setting.OrgID,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrgIDCondition(setting.OrgID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
					),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.setting.InstanceID, setting.InstanceID)
			assert.Equal(t, tt.setting.OrgID, setting.OrgID)
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
		changes  []database.Change
		err      error
	}

	// TESTS
	tests := []test{
		{
			name: "happy path",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLabel,
					OwnerType:  domain.OwnerTypeInstance,
					LabelState: gu.Ptr(domain.LabelStatePreview),
					Settings:   []byte("{}"),
				},
			},
		},
		{
			name: "happy path, no org",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					// OrgID:      &orgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLabel,
					OwnerType:  domain.OwnerTypeInstance,
					LabelState: gu.Ptr(domain.LabelStatePreview),
					Settings:   []byte("{}"),
				},
			},
		},
		{
			name: "adding setting with same instance id, org id twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				settingRepo := repository.LabelRepository()

				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLabel,
					OwnerType:  domain.OwnerTypeInstance,
					LabelState: gu.Ptr(domain.LabelStatePreview),
					Settings:   []byte("{}"),
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
							InstanceID: newInstId,
							OrgID:      &newOrgId,
							Type:       domain.SettingTypeLockout,
							OwnerType:  domain.OwnerTypeInstance,
							Settings:   []byte("{}"),
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
						InstanceID: instanceId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
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
							InstanceID: newInstId,
							OrgID:      &newOrgId,
							Type:       domain.SettingTypeLockout,
							OwnerType:  domain.OwnerTypeInstance,
							LabelState: gu.Ptr(domain.LabelStatePreview),
							Settings:   []byte("{}"),
						},
					}

					err = settingRepo.Set(t.Context(), tx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.OrgID = &orgId
					return &setting
				},
				setting: domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: newInstId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
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
							InstanceID: newInstId,
							OrgID:      &org.ID,
							Type:       domain.SettingTypeLabel,
							OwnerType:  domain.OwnerTypeInstance,
							LabelState: gu.Ptr(domain.LabelStatePreview),
							Settings:   []byte("{}"),
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
					setting.OrgID = &newOrgId
					return &setting
				},
				setting: domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: newInstId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
				},
			}
		}(),
		{
			name: "adding setting with no instance id",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					// InstanceID:        instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLockout,
					OwnerType:  domain.OwnerTypeInstance,
					LabelState: gu.Ptr(domain.LabelStatePreview),
					Settings:   []byte("{}"),
				},
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding idp with non existent instance id",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID: gofakeit.Name(),
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLockout,
					OwnerType:  domain.OwnerTypeInstance,
					LabelState: gu.Ptr(domain.LabelStatePreview),
					Settings:   []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					OrgID:      gu.Ptr(gofakeit.Name()),
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLockout,
					OwnerType:  domain.OwnerTypeInstance,
					LabelState: gu.Ptr(domain.LabelStatePreview),
					Settings:   []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			setting: domain.LabelSetting{
				Setting: &domain.Setting{
					InstanceID: instanceId,
					OrgID:      gu.Ptr(gofakeit.Name()),
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeLockout,
					OwnerType:  domain.OwnerTypeInstance,
					LabelState: gu.Ptr(domain.LabelStatePreview),
					Settings:   []byte("{}"),
				},
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
						settingRepo.OrgIDCondition(setting.OrgID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.setting.InstanceID, setting.InstanceID)
			assert.Equal(t, tt.setting.OrgID, setting.OrgID)
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
		// {
		// 	name: "create label state not set",
		// 	testFunc: func(t *testing.T, tx database.QueryExecutor) error {
		// 		settingRepo := repository.LabelRepository()
		// 		err := settingRepo.Set(t.Context(), tx, &domain.LabelSetting{
		// 			Setting: &domain.Setting{},
		// 		})
		// 		return err
		// 	},
		// 	err: repository.ErrLabelStateMustBeDefined,
		// },
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

// NOTE all updated functions are just a wrapper around repository.updateSetting()
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

	settingRepo := repository.LabelRepository()

	tests := []struct {
		name         string
		testFunc     func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting
		rowsAffected int64
		err          error
	}{
		{
			name: "happy path update settings",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
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
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				// update settings values
				setting.OwnerType = domain.OwnerTypeInstance
				setting.Settings.PrimaryColor = "new_value"
				setting.Settings.BackgroundColor = "new_value"
				setting.Settings.WarnColor = "new_value"
				setting.Settings.FontColor = "new_value"
				setting.Settings.PrimaryColorDark = "new_value"
				setting.Settings.BackgroundColorDark = "new_value"
				setting.Settings.WarnColorDark = "new_value"
				setting.Settings.FontColorDark = "new_value"
				setting.Settings.HideLoginNameSuffix = false
				setting.Settings.ErrorMsgPopup = false
				setting.Settings.DisableWatermark = false
				setting.Settings.ThemeMode = domain.LabelPolicyThemeLight
				return &setting
			},
			rowsAffected: 1,
		},
		{
			name: "happy path update settings, no org",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						// OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
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
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				// update settings values
				setting.OwnerType = domain.OwnerTypeInstance
				setting.Settings.PrimaryColor = "new_value"
				setting.Settings.BackgroundColor = "new_value"
				setting.Settings.WarnColor = "new_value"
				setting.Settings.FontColor = "new_value"
				setting.Settings.PrimaryColorDark = "new_value"
				setting.Settings.BackgroundColorDark = "new_value"
				setting.Settings.WarnColorDark = "new_value"
				setting.Settings.FontColorDark = "new_value"
				setting.Settings.HideLoginNameSuffix = false
				setting.Settings.ErrorMsgPopup = false
				setting.Settings.DisableWatermark = false
				setting.Settings.ThemeMode = domain.LabelPolicyThemeLight
				return &setting
			},
			rowsAffected: 1,
		},
		{
			name: "update setting with non existent instance id",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: "non-existent-instanceID",
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
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
					},
				}

				return &setting
			},
			// expect nothing to get updated
			rowsAffected: 0,
			err:          new(database.ForeignKeyError),
		},
		{
			name: "update setting with non existent org id",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      gu.Ptr("non-existent-orgID"),
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
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
					},
				}
				return &setting
			},
			// expect nothing to get updated
			rowsAffected: 0,
			err:          new(database.ForeignKeyError),
		},
		{
			name: "update setting with wrong type",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
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
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				// change type
				setting.Type = domain.SettingTypeOrganization
				return &setting
			},
			// expect nothing to get updated
			rowsAffected: 0,
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

			setting := tt.testFunc(t, tx)

			// update setting
			// rowsAffected, err := settingRepo.UpdateLabel(
			err = settingRepo.Set(
				t.Context(), tx,
				setting,
			)
			afterUpdate := time.Now()
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			require.NoError(t, err)

			// assert.Equal(t, tt.rowsAffected, rowsAffected)

			assert.ErrorIs(t, err, tt.err)

			// if rowsAffected == 0 {
			// 	return
			// }

			updatedSetting, err := settingRepo.Get(
				t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrgIDCondition(setting.OrgID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, setting.InstanceID, updatedSetting.InstanceID)
			assert.Equal(t, setting.OrgID, updatedSetting.OrgID)
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
		InstanceID: instanceId,
		OrgID:      &orgId,
		ID:         gofakeit.Name(),
		Type:       domain.SettingTypePasswordExpiry,
		OwnerType:  domain.OwnerTypeInstance,
		Settings:   []byte("{}"),
		CreatedAt:  now,
		UpdatedAt:  &now,
	}
	settingRepo := repository.SettingsRepository()
	err = settingRepo.Create(t.Context(), tx, &preExistingSetting)
	require.NoError(t, err)

	type test struct {
		name                       string
		testFunc                   func(t *testing.T, tx database.QueryExecutor) *domain.Setting
		settingIdentifierCondition database.Condition
		err                        error
	}

	tests := []test{
		{
			name: "happy path get using type",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
				setting := domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					Type:       domain.SettingTypeLogin,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
					CreatedAt:  now,
					UpdatedAt:  &now,
				}

				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
			settingIdentifierCondition: settingRepo.TypeCondition(domain.SettingTypeLogin),
		},
		{
			name: "happy path get using type, no org",
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
			settingIdentifierCondition: settingRepo.TypeCondition(domain.SettingTypeLogin),
		},
		{
			name: "get using non existent id",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
				setting := domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),

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
			settingIdentifierCondition: settingRepo.IDCondition("non-existent-id"),
			err:                        new(database.NoRowFoundError),
		},
		func() test {
			return test{
				name: "non existent orgID",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypePasswordComplexity,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)
					setting.OrgID = gu.Ptr("non-existent-orgID")
					return &setting
				},
				err: new(database.NoRowFoundError),
			}
		}(),
		func() test {
			return test{
				name: "non existent instanceID",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Setting {
					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLockout,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)
					setting.InstanceID = "non-existent-instanceID"
					return &setting
				},
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

			// get setting
			returnedSetting, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrgIDCondition(setting.OrgID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						// settingRepo.LabelStateCondition(*setting.LabelState),
					),
				),
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedSetting.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LoginSettings{
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
					},
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

// testing that values written to the db are successfully retrieved
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
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting
		err      error
	}

	tests := []test{
		{
			name: "happy path, state preview",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
					},
					Settings: domain.LabelSettings{
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
					},
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
		{
			name: "happy path, state activated",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStateActivated),
					},
					Settings: domain.LabelSettings{
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
					},
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
		{
			name: "get label policy using wrong state",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStateActivated),
					},
				}

				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.LabelState = gu.Ptr(domain.LabelStatePreview)
				return &setting
			},
			err: new(database.NoRowFoundError),
		},
		{
			name: "get non existent label policy",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStateActivated),
						Settings:   []byte("{}"),
					},
				}

				return &setting
			},
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
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordComplexitySettings{
						MinLength:    89,
						HasLowercase: true,
						HasUppercase: true,
						HasNumber:    true,
						HasSymbol:    true,
					},
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
				// setting.InstanceID,
				// setting.OrgID,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordExpirySettings{
						ExpireWarnDays: 30,
						MaxAgeDays:     30,
					},
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
				// setting.InstanceID,
				// setting.OrgID,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LockoutSettings{
						MaxPasswordAttempts: 50,
						MaxOTPAttempts:      50,
						ShowLockOutFailures: true,
					},
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.SecuritySettings{
						Enabled:               true,
						EnableIframeEmbedding: true,
						AllowedOrigins:        []string{"value"},
						EnableImpersonation:   true,
					},
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.DomainSettings{
						UserLoginMustBeDomain:                  true,
						ValidateOrgDomains:                     true,
						SMTPSenderAddressMatchesInstanceDomain: true,
					},
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.OrganizationSettings{
						OrganizationScopedUsernames: true,
					},
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LoginSettings{
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
						// MFAType:                    []domain.MultiFactorType{domain.MultiFactorTypeUnspecified},
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.Settings = domain.LoginSettings{
					AllowUserNamePassword:      false,
					AllowRegister:              false,
					AllowExternalIDP:           false,
					ForceMFA:                   false,
					ForceMFALocalOnly:          false,
					HidePasswordReset:          false,
					IgnoreUnknownUsernames:     false,
					AllowDomainDiscovery:       false,
					DisableLoginWithEmail:      false,
					DisableLoginWithPhone:      false,
					PasswordlessType:           domain.PasswordlessTypeNotAllowed,
					DefaultRedirectURI:         "www.not-example.com",
					PasswordCheckLifetime:      time.Minute,
					ExternalLoginCheckLifetime: time.Minute,
					MFAInitSkipLifetime:        time.Minute,
					SecondFactorCheckLifetime:  time.Minute,
					MultiFactorCheckLifetime:   time.Minute,

					// TODO
					// MFAType: []domain.MultiFactorType{domain.MultiFactorTypeU2FWithPIN},
				}
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
					// TODO
					// settingRepo.AddMFAType(domain.MultiFactorTypeU2FWithPIN),
					// settingRepo.RemoveMFAType(domain.MultiFactorTypeU2FWithPIN),
					// settingRepo.RemoveMFAType(domain.MultiFactorTypeUnspecified),
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
					settingRepo.OrgIDCondition(setting.OrgID),
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
						settingRepo.OrgIDCondition(setting.OrgID),
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

			fmt.Printf("\033[43m[DBUGPRINT]\033[0m[:1]\033[43m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m setting.Settings.MFAType = %+v\n", setting.Settings.MFAType)
			fmt.Printf("\033[43m[DBUGPRINT]\033[0m[:1]\033[43m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m returnedSetting.Settings.MFAType = %+v\n", returnedSetting.Settings.MFAType)
			assert.Equal(t, setting.InstanceID, returnedSetting.InstanceID)
			assert.Equal(t, setting.OrgID, returnedSetting.OrgID)
			assert.Equal(t, setting.ID, returnedSetting.ID)
			assert.Equal(t, setting.Type, returnedSetting.Type)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
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
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting
		changes  database.Changes
		err      error
	}

	tests := []test{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
					},
					Settings: domain.LabelSettings{
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
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.Settings = domain.LabelSettings{
					PrimaryColor:        "new_value",
					BackgroundColor:     "new_value",
					WarnColor:           "new_value",
					FontColor:           "new_value",
					PrimaryColorDark:    "new_value",
					BackgroundColorDark: "new_value",
					WarnColorDark:       "new_value",
					FontColorDark:       "new_value",
					HideLoginNameSuffix: false,
					ErrorMsgPopup:       false,
					DisableWatermark:    false,
					ThemeMode:           domain.LabelPolicyThemeAuto,
				}

				return &setting
			},
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
				database.And(
					settingRepo.InstanceIDCondition(setting.InstanceID),
					settingRepo.OrgIDCondition(setting.OrgID),
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordComplexitySettings{
						MinLength:    89,
						HasLowercase: true,
						HasUppercase: true,
						HasNumber:    true,
						HasSymbol:    true,
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.Settings = domain.PasswordComplexitySettings{
					MinLength:    200,
					HasLowercase: false,
					HasUppercase: false,
					HasNumber:    false,
					HasSymbol:    false,
				}

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
					settingRepo.OrgIDCondition(setting.OrgID),
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordExpirySettings{
						ExpireWarnDays: 30,
						MaxAgeDays:     30,
					},
				}
				err = settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.Settings = domain.PasswordExpirySettings{
					ExpireWarnDays: 200,
					MaxAgeDays:     200,
				}

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
					settingRepo.OrgIDCondition(setting.OrgID),
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LockoutSettings{
						MaxPasswordAttempts: 50,
						MaxOTPAttempts:      50,
						ShowLockOutFailures: true,
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.Settings = domain.LockoutSettings{
					MaxPasswordAttempts: 100,
					MaxOTPAttempts:      100,
					ShowLockOutFailures: false,
				}

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
					settingRepo.OrgIDCondition(setting.OrgID),
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.SecuritySettings{
						Enabled:               true,
						EnableIframeEmbedding: true,
						AllowedOrigins:        []string{"value"},
						EnableImpersonation:   true,
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.Settings = domain.SecuritySettings{
					Enabled:               false,
					EnableIframeEmbedding: false,
					AllowedOrigins:        []string{"new_value"},
					EnableImpersonation:   false,
				}

				return &setting
			},
			changes: database.Changes{
				settingRepo.SetLabelSettings(
					settingRepo.SetEnabled(false),
					settingRepo.SetEnableIframeEmbedding(false),
					// TODO
					// settingRepo.SetAllowedOrigins([]),
					settingRepo.SetEnableImpersonation(false),
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
					settingRepo.OrgIDCondition(setting.OrgID),
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.DomainSettings{
						UserLoginMustBeDomain:                  true,
						ValidateOrgDomains:                     true,
						SMTPSenderAddressMatchesInstanceDomain: true,
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.Settings = domain.DomainSettings{
					UserLoginMustBeDomain:                  false,
					ValidateOrgDomains:                     false,
					SMTPSenderAddressMatchesInstanceDomain: false,
				}

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
					settingRepo.OrgIDCondition(setting.OrgID),
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeInstance,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.OrganizationSettings{
						OrganizationScopedUsernames: true,
					},
				}
				err := settingRepo.Set(t.Context(), tx, &setting)
				require.NoError(t, err)

				setting.Settings = domain.OrganizationSettings{
					OrganizationScopedUsernames: false,
				}

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
					settingRepo.OrgIDCondition(setting.OrgID),
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
						settingRepo.OrgIDCondition(setting.OrgID),
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
			assert.Equal(t, returnedSetting.OrgID, setting.OrgID)
			assert.Equal(t, returnedSetting.ID, setting.ID)
			assert.Equal(t, returnedSetting.Type, setting.Type)
			assert.Equal(t, returnedSetting.OwnerType, setting.OwnerType)

			assert.Equal(t, returnedSetting.Settings, setting.Settings)
		})
	}
}

// gocognit linting fails due to number of test cases
// and the fact that each test case has a testFunc()
//
//nolint:gocognit
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
					InstanceID: newInstanceId,
					OrgID:      &newOrgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeSecurity,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
					CreatedAt:  now,
					UpdatedAt:  &now,
				}
				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)

				noOfSettings := 5
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
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
					InstanceID: instanceId,
					OrgID:      &newOrgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeSecurity,
					OwnerType:  domain.OwnerTypeInstance,
					Settings:   []byte("{}"),
					CreatedAt:  now,
					UpdatedAt:  &now,
				}
				err := settingRepo.Create(t.Context(), tx, &setting)
				require.NoError(t, err)

				noOfSettings := 5
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingType(i + 1),
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
			conditionClauses: database.WithCondition(settingRepo.OrgIDCondition(&orgId)),
		},
		{
			name: "happy path single setting no filter",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Setting {
				noOfSettings := 1
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingType(i + 1),
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingType(i + 1),
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
					}

					err := settingRepo.Create(t.Context(), tx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
		},
		// func() test {
		// 	id := gofakeit.Name()
		// 	return test{
		// 		name: "setting filter on id",
		// 		testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Setting {
		// 			// create setting
		// 			// this setting is created as an additional setting which should NOT
		// 			// be returned in the results of this test case
		// 			setting := domain.Setting{
		// 				InstanceID: instanceId,
		// 				OrgID:      &orgId,
		// 				ID:         gofakeit.Name(),
		// 				Type:       domain.SettingTypeSecurity,
		// OwnerType:  domain.OwnerTypeInstance,
		// 				Settings:   []byte("{}"),
		// 				CreatedAt:  now,
		// 				UpdatedAt:  &now,
		// 			}
		// 			err := settingRepo.Create(t.Context(), tx, &setting)
		// 			require.NoError(t, err)

		// 			noOfSettings := 1
		// 			settings := make([]*domain.Setting, noOfSettings)
		// 			for i := range noOfSettings {

		// 				setting := domain.Setting{
		// 					InstanceID: instanceId,
		// 					OrgID:      &orgId,
		// 					ID:         id,
		// 					Type:       domain.SettingType(i + 1),
		// OwnerType:  domain.OwnerTypeInstance,
		// 					Settings:   []byte("{}"),
		// 					CreatedAt:  now,
		// 					UpdatedAt:  &now,
		// 				}

		// 				err := settingRepo.Create(t.Context(), tx, &setting)
		// 				require.NoError(t, err)

		// 				settings[i] = &setting
		// 				// id = setting.ID
		// 			}

		// 			return settings
		// 		},
		// 		conditionClauses: database.WithCondition(settingRepo.IDCondition(id)),
		// 	}
		// }(),
		{
			name: "multiple settings filter on type",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Setting {
				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),

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
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypePasswordExpiry,
						OwnerType:  domain.OwnerTypeInstance,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
				assert.Equal(t, returnedSettings[i].OrgID, setting.OrgID)
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
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor)
		// settingType     domain.SettingType
		condtions       database.Condition
		noOfDeletedRows int64
	}
	tests := []test{
		func() test {
			var noOfSettings int64 = 1
			return test{
				name: "happy path delete",
				testFunc: func(t *testing.T, tx database.QueryExecutor) {
					for range noOfSettings {
						setting := domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
							Type:       domain.SettingTypeLogin,
							OwnerType:  domain.OwnerTypeInstance,
							Settings:   []byte("{}"),
							CreatedAt:  now,
							UpdatedAt:  &now,
						}

						err := settingRepo.Create(t.Context(), tx, &setting)
						require.NoError(t, err)
					}
				},
				noOfDeletedRows: noOfSettings,
				// settingType:     domain.SettingTypeLogin,
				condtions: database.And(
					settingRepo.InstanceIDCondition(instanceId),
					settingRepo.OrgIDCondition(&orgId),
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
							InstanceID: instanceId,
							OrgID:      &orgId,
							ID:         gofakeit.Name(),
							Type:       domain.SettingTypeLogin,
							OwnerType:  domain.OwnerTypeInstance,
							Settings:   []byte("{}"),
							CreatedAt:  now,
							UpdatedAt:  &now,
						}

						err := settingRepo.Create(t.Context(), tx, &setting)
						require.NoError(t, err)
					}

					// // delete organization
					// affectedRows, err := settingRepo.Delete(t.Context(), tx,
					// 	database.And(
					// 		settingRepo.InstanceIDCondition(instanceId),
					// 		settingRepo.OrgIDCondition(&orgId),
					// 		settingRepo.TypeCondition(domain.SettingTypeLogin),
					// 		settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance),
					// 	),
					// )
					// assert.Equal(t, int64(1), affectedRows)
					// require.NoError(t, err)
				},
				// this test should return 0 affected rows as the setting was already deleted
				noOfDeletedRows: 0,
				condtions: database.And(
					settingRepo.InstanceIDCondition(instanceId),
					settingRepo.OrgIDCondition(&orgId),
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
							InstanceID: instanceId,
							OrgID:      &orgId,
							ID:         gofakeit.Name(),
							Type:       domain.SettingTypeLogin,
							OwnerType:  domain.OwnerTypeInstance,
							Settings:   []byte("{}"),
							CreatedAt:  now,
							UpdatedAt:  &now,
						}

						err := settingRepo.Create(t.Context(), tx, &setting)
						require.NoError(t, err)
					}

					// delete organization
					// affectedRows, err := settingRepo.Delete(t.Context(), tx,
					// 	instanceId,
					// 	&orgId,
					// 	domain.SettingTypeLogin,
					// )
					affectedRows, err := settingRepo.Delete(t.Context(), tx,
						database.And(
							settingRepo.InstanceIDCondition(instanceId),
							settingRepo.OrgIDCondition(&orgId),
							settingRepo.TypeCondition(domain.SettingTypeLogin),
							settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance),
						),
					)
					assert.Equal(t, int64(1), affectedRows)
					require.NoError(t, err)
				},
				// this test should return 0 affected rows as the setting was already deleted
				noOfDeletedRows: 0,
				condtions: database.And(
					settingRepo.InstanceIDCondition(instanceId),
					settingRepo.OrgIDCondition(&orgId),
					settingRepo.TypeCondition(domain.SettingTypeLogin),
					settingRepo.OwnerTypeCondition(domain.OwnerTypeInstance),
				),
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

			if tt.testFunc != nil {
				tt.testFunc(t, tx)
			}

			// delete setting
			noOfDeletedRows, err := settingRepo.Delete(t.Context(), tx,
				tt.condtions,
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check setting was deleted
			organization, err := settingRepo.Get(t.Context(), tx,
				database.WithCondition(
					tt.condtions,
				),
				// instanceId,
				// &orgId,
				// tt.settingType,
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, organization)
		})
	}
}

// func TestDeleteSettingsForInstance(t *testing.T) {
// 	tx, err := pool.Begin(t.Context(), nil)
// 	require.NoError(t, err)
// 	defer func() {
// 		err = tx.Rollback(t.Context())
// 		if err != nil {
// 			t.Logf("error during rollback: %v", err)
// 		}
// 	}()

// 	// create instance
// 	instanceId := gofakeit.Name()
// 	instance := domain.Instance{
// 		ID:              instanceId,
// 		Name:            gofakeit.Name(),
// 		DefaultOrgID:    "defaultOrgId",
// 		IAMProjectID:    "iamProject",
// 		ConsoleClientID: "consoleCLient",
// 		ConsoleAppID:    "consoleApp",
// 		DefaultLanguage: "defaultLanguage",
// 	}
// 	instanceRepo := repository.InstanceRepository()
// 	err = instanceRepo.Create(t.Context(), tx, &instance)
// 	require.NoError(t, err)

// 	// create org
// 	orgId := gofakeit.Name()
// 	org := domain.Organization{
// 		ID:         orgId,
// 		Name:       gofakeit.Name(),
// 		InstanceID: instanceId,
// 		State:      domain.OrgStateActive,
// 	}
// 	organizationRepo := repository.OrganizationRepository()
// 	err = organizationRepo.Create(t.Context(), tx, &org)
// 	require.NoError(t, err)

// 	settingRepo := repository.SettingsRepository()

// 	type test struct {
// 		name            string
// 		testFunc        func(t *testing.T, tx database.QueryExecutor)
// 		noOfDeletedRows int64
// 	}
// 	tests := []test{
// 		func() test {
// 			return test{
// 				name: "delete all settings for instance",
// 				testFunc: func(t *testing.T, tx database.QueryExecutor) {
// 					// create with org
// 					err = settingRepo.CreateLogin(t.Context(), tx, &domain.LoginSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreateLabel(t.Context(), tx, &domain.LabelSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 							LabelState: gu.Ptr(domain.LabelStatePreview),
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreatePasswordComplexity(t.Context(), tx, &domain.PasswordComplexitySetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreatePasswordExpiry(t.Context(), tx, &domain.PasswordExpirySetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreateLockout(t.Context(), tx, &domain.LockoutSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreateSecurity(t.Context(), tx, &domain.SecuritySetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreateDomain(t.Context(), tx, &domain.DomainSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					// create without org
// 					err = settingRepo.CreateLogin(t.Context(), tx, &domain.LoginSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							// OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreateLabel(t.Context(), tx, &domain.LabelSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							// OrgID:      &orgId,
// 							LabelState: gu.Ptr(domain.LabelStatePreview),
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreatePasswordComplexity(t.Context(), tx, &domain.PasswordComplexitySetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							// OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreatePasswordExpiry(t.Context(), tx, &domain.PasswordExpirySetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							// OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreateLockout(t.Context(), tx, &domain.LockoutSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							// OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					err = settingRepo.CreateSecurity(t.Context(), tx, &domain.SecuritySetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							// OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)
// 				},
// 				noOfDeletedRows: 13,
// 			}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tx, err := tx.Begin(t.Context())
// 			require.NoError(t, err)
// 			defer func() {
// 				err = tx.Rollback(t.Context())
// 				if err != nil {
// 					t.Logf("error during rollback: %v", err)
// 				}
// 			}()

// 			if tt.testFunc != nil {
// 				tt.testFunc(t, tx)
// 			}

// 			// delete settings
// 			noOfDeletedRows, err := settingRepo.DeleteSettingsForInstance(
// 				t.Context(), tx,
// 				instanceId,
// 			)
// 			require.NoError(t, err)
// 			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)
// 		})
// 	}
// }

func TestActivateLabelSettings(t *testing.T) {
	settingRepo := repository.LabelRepository()

	type test struct {
		name            string
		testFunc        func(t *testing.T, pool database.QueryExecutor) *domain.LabelSetting
		noOfDeletedRows int64
	}
	tests := []test{
		{
			name: "happy path activate label settings, no preexisting active label",
			testFunc: func(t *testing.T, pool database.QueryExecutor) *domain.LabelSetting {
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
				err := instanceRepo.Create(t.Context(), pool, &instance)
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
				err = organizationRepo.Create(t.Context(), pool, &org)
				require.NoError(t, err)

				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
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
					},
				}
				return &setting
			},
		},
		{
			name: "happy path activate label settings, with preexisting preview label",
			testFunc: func(t *testing.T, pool database.QueryExecutor) *domain.LabelSetting {
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
				err := instanceRepo.Create(t.Context(), pool, &instance)
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
				err = organizationRepo.Create(t.Context(), pool, &org)
				require.NoError(t, err)

				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
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
					},
				}
				err = settingRepo.Set(t.Context(), pool, &setting)
				require.NoError(t, err)

				// update settings values
				setting.OwnerType = domain.OwnerTypeInstance
				setting.Settings.PrimaryColor = "new_value"
				setting.Settings.BackgroundColor = "new_value"
				setting.Settings.WarnColor = "new_value"
				setting.Settings.FontColor = "new_value"
				setting.Settings.PrimaryColorDark = "new_value"
				setting.Settings.BackgroundColorDark = "new_value"
				setting.Settings.WarnColorDark = "new_value"
				setting.Settings.FontColorDark = "new_value"
				setting.Settings.HideLoginNameSuffix = false
				setting.Settings.ErrorMsgPopup = false
				setting.Settings.DisableWatermark = false
				setting.Settings.ThemeMode = domain.LabelPolicyThemeLight

				return &setting
			},
		},
		{
			name: "happy path activate label settings, with preexisting active label",
			testFunc: func(t *testing.T, pool database.QueryExecutor) *domain.LabelSetting {
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
				err := instanceRepo.Create(t.Context(), pool, &instance)
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
				err = organizationRepo.Create(t.Context(), pool, &org)
				require.NoError(t, err)

				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						OwnerType:  domain.OwnerTypeOrganization,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
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
					},
				}
				err = settingRepo.Set(t.Context(), pool, &setting)
				require.NoError(t, err)

				// update settings values
				setting.OwnerType = domain.OwnerTypeInstance
				setting.Settings.PrimaryColor = "new_value"
				setting.Settings.BackgroundColor = "new_value"
				setting.Settings.WarnColor = "new_value"
				setting.Settings.FontColor = "new_value"
				setting.Settings.PrimaryColorDark = "new_value"
				setting.Settings.BackgroundColorDark = "new_value"
				setting.Settings.WarnColorDark = "new_value"
				setting.Settings.FontColorDark = "new_value"
				setting.Settings.HideLoginNameSuffix = false
				setting.Settings.ErrorMsgPopup = false
				setting.Settings.DisableWatermark = false
				setting.Settings.ThemeMode = domain.LabelPolicyThemeLight

				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setting := tt.testFunc(t, pool)

			before := time.Now()
			err := settingRepo.ActivateLabelSetting(
				t.Context(),
				pool,
				setting,
			)
			require.NoError(t, err)
			after := time.Now()

			updatedSetting, err := settingRepo.Get(
				t.Context(),
				pool,
				database.WithCondition(
					database.And(
						settingRepo.InstanceIDCondition(setting.InstanceID),
						settingRepo.OrgIDCondition(setting.OrgID),
						settingRepo.TypeCondition(setting.Type),
						settingRepo.OwnerTypeCondition(setting.OwnerType),
						settingRepo.LabelStateCondition(domain.LabelStateActivated),
					),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, setting.Settings, updatedSetting.Settings)
			assert.Equal(t, domain.LabelStateActivated, *updatedSetting.LabelState)
			assert.WithinRange(t, *updatedSetting.UpdatedAt, before, after)
		})
	}
}

// func TestDeleteSettingsForOrg(t *testing.T) {
// 	tx, err := pool.Begin(t.Context(), nil)
// 	require.NoError(t, err)
// 	defer func() {
// 		err = tx.Rollback(t.Context())
// 		if err != nil {
// 			t.Logf("error during rollback: %v", err)
// 		}
// 	}()

// 	// create instance
// 	instanceId := gofakeit.Name()
// 	instance := domain.Instance{
// 		ID:              instanceId,
// 		Name:            gofakeit.Name(),
// 		DefaultOrgID:    "defaultOrgId",
// 		IAMProjectID:    "iamProject",
// 		ConsoleClientID: "consoleCLient",
// 		ConsoleAppID:    "consoleApp",
// 		DefaultLanguage: "defaultLanguage",
// 	}
// 	instanceRepo := repository.InstanceRepository()
// 	err = instanceRepo.Create(t.Context(), tx, &instance)
// 	require.NoError(t, err)

// 	// create org
// 	orgId := gofakeit.Name()
// 	org := domain.Organization{
// 		ID:         orgId,
// 		Name:       gofakeit.Name(),
// 		InstanceID: instanceId,
// 		State:      domain.OrgStateActive,
// 	}
// 	organizationRepo := repository.OrganizationRepository()
// 	err = organizationRepo.Create(t.Context(), tx, &org)
// 	require.NoError(t, err)

// 	settingRepo := repository.SettingsRepository()

// 	type test struct {
// 		name            string
// 		testFunc        func(t *testing.T, tx database.QueryExecutor)
// 		noOfDeletedRows int64
// 	}
// 	tests := []test{
// 		func() test {
// 			return test{
// 				name: "delete all settings for instance",
// 				testFunc: func(t *testing.T, tx database.QueryExecutor) {
// 					// create with org
// 					settingRepo := repository.LoginRepository()
// 					err = settingRepo.Set(t.Context(), tx, &domain.LoginSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					labelRepo := repository.LabelRepository()
// 					err = labelRepo.Set(t.Context(), tx, &domain.LabelSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 							LabelState: gu.Ptr(domain.LabelStatePreview),
// 						},
// 					})
// 					require.NoError(t, err)

// 					passwordComplexityRepo := repository.PasswordComplexityRepository()
// 					err = passwordComplexityRepo.Set(t.Context(), tx, &domain.PasswordComplexitySetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					passwordExpiryRepo := repository.PasswordExpiryRepository()
// 					err = passwordExpiryRepo.Set(t.Context(), tx, &domain.PasswordExpirySetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					lockoutRepo := repository.LockoutRepository()
// 					err = lockoutRepo.Set(t.Context(), tx, &domain.LockoutSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					securityRepo := repository.SecurityRepository()
// 					err = securityRepo.Set(t.Context(), tx, &domain.SecuritySetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)

// 					domainRepo := repository.DomainRepository()
// 					err = domainRepo.Set(t.Context(), tx, &domain.DomainSetting{
// 						Setting: &domain.Setting{
// 							InstanceID: instanceId,
// 							OrgID:      &orgId,
// 						},
// 					})
// 					require.NoError(t, err)
// 				},
// 				noOfDeletedRows: 7,
// 			}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tx, err := tx.Begin(t.Context())
// 			require.NoError(t, err)
// 			defer func() {
// 				err = tx.Rollback(t.Context())
// 				if err != nil {
// 					t.Logf("error during rollback: %v", err)
// 				}
// 			}()

// 			if tt.testFunc != nil {
// 				tt.testFunc(t, tx)
// 			}

// 			// delete settings
// 			noOfDeletedRows, err := settingRepo.DeleteSettingsForOrg(
// 				t.Context(), tx,
// 				orgId,
// 			)
// 			require.NoError(t, err)
// 			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)
// 		})
// 	}
// }

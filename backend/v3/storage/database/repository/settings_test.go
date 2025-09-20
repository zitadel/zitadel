package repository_test

import (
	"context"
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
	instanceRepo := repository.InstanceRepository(pool)
	err := instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.Setting
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
				Settings:   []byte("{}"),
			},
		},
		{
			name: "happy path, no org",
			setting: domain.Setting{
				InstanceID: instanceId,
				// OrgID:      &orgId,
				ID:       gofakeit.Name(),
				Type:     domain.SettingTypeLogin,
				Settings: []byte("{}"),
			},
		},
		{
			name: "adding setting with same instance id, org id twice",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Setting {
				settingRepo := repository.SettingsRepository(pool)

				setting := domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					Type:       domain.SettingTypePasswordExpiry,
					Settings:   []byte("{}"),
				}

				err := settingRepo.Create(ctx, &setting)
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
				testFunc: func(ctx context.Context, t *testing.T) *domain.Setting {
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
					instanceRepo := repository.InstanceRepository(pool)
					err := instanceRepo.Create(ctx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository(pool)
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository(pool)
					setting := domain.Setting{
						InstanceID: newInstId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypeLockout,
						Settings:   []byte("{}"),
					}

					err = settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.InstanceID = instanceId
					return &setting
				},
				setting: domain.Setting{
					InstanceID: instanceId,
					OrgID:      &newOrgId,
					Type:       domain.SettingTypePasswordExpiry,
					Settings:   []byte("{}"),
				},
				err: new(database.IntegrityViolationError),
			}
		}(),
		func() test {
			newInstId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return test{
				name: "adding setting with same instance, type, different org (org on different instance)",
				testFunc: func(ctx context.Context, t *testing.T) *domain.Setting {
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
					instanceRepo := repository.InstanceRepository(pool)
					err := instanceRepo.Create(ctx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository(pool)
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository(pool)
					setting := domain.Setting{
						InstanceID: newInstId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypeLockout,
						Settings:   []byte("{}"),
					}

					err = settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					// change the instance
					setting.OrgID = &orgId
					return &setting
				},
				setting: domain.Setting{
					InstanceID: newInstId,
					OrgID:      &orgId,
					Type:       domain.SettingTypePasswordExpiry,
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
				testFunc: func(ctx context.Context, t *testing.T) *domain.Setting {
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
					instanceRepo := repository.InstanceRepository(pool)
					err := instanceRepo.Create(ctx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository(pool)
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository(pool)
					setting := domain.Setting{
						InstanceID: newInstId,
						OrgID:      &org.ID,
						Type:       domain.SettingTypePasswordExpiry,
						Settings:   []byte("{}"),
					}

					err = settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					// create another org
					org = domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					// change the org id
					setting.OrgID = &newOrgId
					return &setting
				},
				setting: domain.Setting{
					InstanceID: newInstId,
					OrgID:      &newOrgId,
					Type:       domain.SettingTypePasswordExpiry,
					Settings:   []byte("{}"),
				},
			}
		}(),
		{
			name: "adding setting with no type",
			setting: domain.Setting{
				InstanceID: instanceId,
				OrgID:      &orgId,
				// Type:     domain.SettingTypeBranding,
				Settings: []byte("{}"),
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding setting with no instance id",
			setting: domain.Setting{
				// InstanceID:        instanceId,
				OrgID:    &orgId,
				ID:       gofakeit.Name(),
				Type:     domain.SettingTypeLockout,
				Settings: []byte("{}"),
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
				Settings:   []byte("{}"),
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
				Settings:   []byte("{}"),
			},
			err: new(database.ForeignKeyError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var setting *domain.Setting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			} else {
				setting = &tt.setting
			}
			settingRepo := repository.SettingsRepository(pool)

			// create setting
			beforeCreate := time.Now()
			err = settingRepo.Create(ctx, setting)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check organization values
			setting, err = settingRepo.Get(ctx,
				setting.InstanceID,
				setting.OrgID,
				setting.Type,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.setting.InstanceID, setting.InstanceID)
			assert.Equal(t, tt.setting.OrgID, setting.OrgID)
			assert.Equal(t, tt.setting.Type, setting.Type)
			assert.WithinRange(t, setting.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, setting.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

// NOTE all create functions are just a wrapper around repository.createSetting()
// for the sake of testing, CreatePasswordExpiry was used, but the underlying code is the same
// across all Create.*() functions
func TestCreateSpecificSetting(t *testing.T) {
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
	instanceRepo := repository.InstanceRepository(pool)
	err := instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting
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
					ID:       gofakeit.Name(),
					Type:     domain.SettingTypePasswordExpiry,
					Settings: []byte("{}"),
				},
			},
		},
		{
			name: "adding setting with same instance id, org id twice",
			testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
				settingRepo := repository.SettingsRepository(pool)

				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypePasswordExpiry,
						Settings:   []byte("{}"),
					},
				}

				err := settingRepo.CreatePasswordExpiry(ctx, &setting)
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
				testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
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
					instanceRepo := repository.InstanceRepository(pool)
					err := instanceRepo.Create(ctx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository(pool)
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository(pool)
					setting := domain.PasswordExpirySetting{
						Setting: &domain.Setting{
							InstanceID: newInstId,
							OrgID:      &newOrgId,
							Type:       domain.SettingTypeLockout,
							Settings:   []byte("{}"),
						},
					}

					err = settingRepo.CreatePasswordExpiry(ctx, &setting)
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
				testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
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
					instanceRepo := repository.InstanceRepository(pool)
					err := instanceRepo.Create(ctx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository(pool)
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository(pool)
					setting := domain.PasswordExpirySetting{
						Setting: &domain.Setting{
							InstanceID: newInstId,
							OrgID:      &newOrgId,
							Type:       domain.SettingTypeLockout,
							Settings:   []byte("{}"),
						},
					}

					err = settingRepo.CreatePasswordExpiry(ctx, &setting)
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
				testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
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
					instanceRepo := repository.InstanceRepository(pool)
					err := instanceRepo.Create(ctx, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository(pool)
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository(pool)
					setting := domain.PasswordExpirySetting{
						Setting: &domain.Setting{
							InstanceID: newInstId,
							OrgID:      &org.ID,
							Type:       domain.SettingTypePasswordExpiry,
							Settings:   []byte("{}"),
						},
					}

					err = settingRepo.CreatePasswordExpiry(ctx, &setting)
					require.NoError(t, err)

					// create another org
					org = domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					err = organizationRepo.Create(ctx, &org)
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
					OrgID:    &orgId,
					ID:       gofakeit.Name(),
					Type:     domain.SettingTypeLockout,
					Settings: []byte("{}"),
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
					Settings:   []byte("{}"),
				},
			},
			err: new(database.ForeignKeyError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := pool.Exec(ctx, "DELETE FROM zitadel.settings")
			require.NoError(t, err)

			var setting *domain.PasswordExpirySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			} else {
				setting = &tt.setting
			}
			settingRepo := repository.SettingsRepository(pool)

			// create setting
			beforeCreate := time.Now()
			err = settingRepo.CreatePasswordExpiry(ctx, setting)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check organization values
			setting, err = settingRepo.GetPasswordExpiry(
				ctx,
				setting.InstanceID,
				setting.OrgID,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.setting.InstanceID, setting.InstanceID)
			assert.Equal(t, tt.setting.OrgID, setting.OrgID)
			assert.Equal(t, tt.setting.Type, setting.Type)
			assert.WithinRange(t, setting.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, setting.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

// NOTE all updated functions are just a wrapper around repository.updateSetting()
// for the sake of testing, UpdatePasswordExpiry was used, but the underlying code is the same
// across all Update.*() functions
func TestUpdateSetting(t *testing.T) {
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
	instanceRepo := repository.InstanceRepository(pool)
	err := instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	settingRepo := repository.SettingsRepository(pool)

	tests := []struct {
		name         string
		testFunc     func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting
		rowsAffected int64
		err          error
	}{
		{
			name: "happy path update settings",
			testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypePasswordExpiry,
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordExpirySettings{
						IsDefault:      true,
						ExpireWarnDays: 10,
						MaxAgeDays:     40,
					},
				}
				err := settingRepo.CreatePasswordExpiry(ctx, &setting)
				require.NoError(t, err)

				// update settings values
				setting.Settings.IsDefault = false
				setting.Settings.ExpireWarnDays = 20
				setting.Settings.MaxAgeDays = 40
				return &setting
			},
			rowsAffected: 1,
		},
		{
			name: "happy path update settings, no org",
			testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						// OrgID:      &orgId,
						ID:       gofakeit.Name(),
						Type:     domain.SettingTypePasswordExpiry,
						Settings: []byte("{}"),
					},
					Settings: domain.PasswordExpirySettings{
						IsDefault:      true,
						ExpireWarnDays: 10,
						MaxAgeDays:     40,
					},
				}
				err := settingRepo.CreatePasswordExpiry(ctx, &setting)
				require.NoError(t, err)

				// update settings values
				setting.Settings.IsDefault = false
				setting.Settings.ExpireWarnDays = 20
				setting.Settings.MaxAgeDays = 40
				return &setting
			},
			rowsAffected: 1,
		},
		{
			name: "update setting with non existent instance id",
			testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: "non-existent-instanceID",
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypePasswordExpiry,
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordExpirySettings{
						IsDefault:      true,
						ExpireWarnDays: 10,
						MaxAgeDays:     40,
					},
				}

				return &setting
			},
			// expect nothing to get updated
			rowsAffected: 0,
		},
		{
			name: "update setting with non existent org id",
			testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      gu.Ptr("non-existent-orgID"),
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypePasswordExpiry,
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordExpirySettings{
						IsDefault:      true,
						ExpireWarnDays: 10,
						MaxAgeDays:     40,
					},
				}
				return &setting
			},
			// expect nothing to get updated
			rowsAffected: 0,
		},
		{
			name: "update setting with wrong type",
			testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypePasswordExpiry,
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordExpirySettings{
						IsDefault:      true,
						ExpireWarnDays: 10,
						MaxAgeDays:     40,
					},
				}
				err := settingRepo.CreatePasswordExpiry(ctx, &setting)
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
			ctx := t.Context()
			_, err := pool.Exec(ctx, "DELETE FROM zitadel.settings")
			require.NoError(t, err)

			settingRepo := repository.SettingsRepository(pool)

			setting := tt.testFunc(ctx, t)

			// update setting
			beforeUpdate := time.Now()
			rowsAffected, err := settingRepo.UpdatePasswordExpiry(
				ctx,
				setting,
			)
			afterUpdate := time.Now()
			require.NoError(t, err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			assert.ErrorIs(t, err, tt.err)

			if rowsAffected == 0 {
				return
			}

			// check updatedSetting values
			updatedSetting, err := settingRepo.GetPasswordExpiry(
				ctx,
				setting.InstanceID,
				setting.OrgID,
			)
			require.NoError(t, err)

			assert.Equal(t, setting.InstanceID, updatedSetting.InstanceID)
			assert.Equal(t, setting.OrgID, updatedSetting.OrgID)
			assert.Equal(t, setting.ID, updatedSetting.ID)
			assert.Equal(t, setting.Type, updatedSetting.Type)
			assert.Equal(t, setting.Settings.IsDefault, updatedSetting.Settings.IsDefault)
			assert.Equal(t, setting.Settings.ExpireWarnDays, updatedSetting.Settings.ExpireWarnDays)
			assert.Equal(t, setting.Settings.MaxAgeDays, updatedSetting.Settings.MaxAgeDays)
			assert.WithinRange(t, updatedSetting.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

func TestGetSetting(t *testing.T) {
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
	instanceRepo := repository.InstanceRepository(pool)
	err := instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	// create setting
	// this setting is created as an additional org which should NOT
	// be returned in the results of the tests
	prexistingSetting := domain.Setting{
		InstanceID: instanceId,
		OrgID:      &orgId,
		ID:         gofakeit.Name(),
		Type:       domain.SettingTypePasswordExpiry,
		Settings:   []byte("{}"),
	}
	settingRepo := repository.SettingsRepository(pool)
	err = settingRepo.Create(t.Context(), &prexistingSetting)
	require.NoError(t, err)

	type test struct {
		name                       string
		testFunc                   func(ctx context.Context, t *testing.T) *domain.Setting
		settingIdentifierCondition domain.OrgIdentifierCondition
		err                        error
	}

	tests := []test{
		{
			name: "happy path get using type",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Setting {
				setting := domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					Type:       domain.SettingTypeLogin,
					Settings:   []byte("{}"),
				}

				err := settingRepo.Create(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
			settingIdentifierCondition: settingRepo.TypeCondition(domain.SettingTypeLogin),
		},
		{
			name: "happy path get using type, no org",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Setting {
				setting := domain.Setting{
					InstanceID: instanceId,
					// OrgID:      &orgId,
					Type:     domain.SettingTypeLogin,
					Settings: []byte("{}"),
				}

				err := settingRepo.Create(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
			settingIdentifierCondition: settingRepo.TypeCondition(domain.SettingTypeLogin),
		},
		{
			name: "get using non existent id",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Setting {
				setting := domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),

					Type:     domain.SettingTypeDomain,
					Settings: []byte("{}"),
				}

				err := settingRepo.Create(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
			settingIdentifierCondition: settingRepo.IDCondition("non-existent-id"),
			err:                        new(database.NoRowFoundError),
		},
		func() test {
			return test{
				name: "non existent orgID",
				testFunc: func(ctx context.Context, t *testing.T) *domain.Setting {
					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypePasswordComplexity,
						Settings:   []byte("{}"),
					}

					err := settingRepo.Create(ctx, &setting)
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
				testFunc: func(ctx context.Context, t *testing.T) *domain.Setting {
					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLockout,
						Settings:   []byte("{}"),
					}

					err := settingRepo.Create(ctx, &setting)
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
			ctx := t.Context()

			var setting *domain.Setting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			}

			// get setting
			returnedIDP, err := settingRepo.Get(ctx,
				setting.InstanceID,
				setting.OrgID,
				setting.Type,
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}

			assert.Equal(t, returnedIDP.InstanceID, setting.InstanceID)
			assert.Equal(t, returnedIDP.OrgID, setting.OrgID)
			assert.Equal(t, returnedIDP.ID, setting.ID)
			assert.Equal(t, returnedIDP.Type, setting.Type)
			assert.Equal(t, returnedIDP.Settings, setting.Settings)
		})
	}
}

// gocognit linting fails due to number of test cases
// and the fact that each test case has a testFunc()
//
//nolint:gocognit
func TestListSetting(t *testing.T) {
	ctx := t.Context()
	pool, stop, err := newEmbeddedDB(ctx)
	require.NoError(t, err)
	defer stop()

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
	instanceRepo := repository.InstanceRepository(pool)
	err = instanceRepo.Create(ctx, &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name             string
		testFunc         func(ctx context.Context, t *testing.T) []*domain.Setting
		conditionClauses func() []database.Condition
		noIDPsReturned   bool
	}
	tests := []test{
		{
			name: "multiple settings filter on instance",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Setting {
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
				err = instanceRepo.Create(ctx, &instance)
				require.NoError(t, err)

				// create org
				newOrgId := gofakeit.Name()
				org := domain.Organization{
					ID:         newOrgId,
					Name:       gofakeit.Name(),
					InstanceID: newInstanceId,
					State:      domain.OrgStateActive,
				}
				organizationRepo := repository.OrganizationRepository(pool)
				err = organizationRepo.Create(ctx, &org)
				require.NoError(t, err)

				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID: newInstanceId,
					OrgID:      &newOrgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeSecurity,
					Settings:   []byte("{}"),
				}
				err := settingRepo.Create(ctx, &setting)
				require.NoError(t, err)

				noOfIDPs := 5
				settings := make([]*domain.Setting, noOfIDPs)
				for i := range noOfIDPs {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						// Type:              domain.SettingTypeLogin,
						Type:     domain.SettingType(i + 1),
						Settings: []byte("{}"),
					}

					err := settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
			conditionClauses: func() []database.Condition {
				return []database.Condition{settingRepo.InstanceIDCondition(instanceId)}
			},
		},
		{
			name: "multiple settings filter on instance, no org",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Setting {
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
				err = instanceRepo.Create(ctx, &instance)
				require.NoError(t, err)

				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID: newInstanceId,
					// OrgID:      &newOrgId,
					ID:       gofakeit.Name(),
					Type:     domain.SettingTypeSecurity,
					Settings: []byte("{}"),
				}
				err := settingRepo.Create(ctx, &setting)
				require.NoError(t, err)

				noOfIDPs := 5
				settings := make([]*domain.Setting, noOfIDPs)
				for i := range noOfIDPs {

					setting := domain.Setting{
						InstanceID: instanceId,
						// OrgID:      &orgId,
						ID: gofakeit.Name(),
						// Type:              domain.SettingTypeLogin,
						Type:     domain.SettingType(i + 1),
						Settings: []byte("{}"),
					}

					err := settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
			conditionClauses: func() []database.Condition {
				return []database.Condition{settingRepo.InstanceIDCondition(instanceId)}
			},
		},
		{
			name: "multiple settings filter on org",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Setting {
				// create org
				newOrgId := gofakeit.Name()
				org := domain.Organization{
					ID:         newOrgId,
					Name:       gofakeit.Name(),
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				}
				organizationRepo := repository.OrganizationRepository(pool)
				err = organizationRepo.Create(ctx, &org)
				require.NoError(t, err)

				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID: instanceId,
					OrgID:      &newOrgId,
					ID:         gofakeit.Name(),
					Type:       domain.SettingTypeSecurity,
					Settings:   []byte("{}"),
				}
				err := settingRepo.Create(ctx, &setting)
				require.NoError(t, err)

				noOfIDPs := 5
				settings := make([]*domain.Setting, noOfIDPs)
				for i := range noOfIDPs {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingType(i + 1),
						Settings:   []byte("{}"),
					}

					err := settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
			conditionClauses: func() []database.Condition {
				return []database.Condition{settingRepo.OrgIDCondition(&orgId)}
			},
		},
		{
			name: "happy path single setting no filter",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Setting {
				noOfIDPs := 1
				settings := make([]*domain.Setting, noOfIDPs)
				for i := range noOfIDPs {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingType(i + 1),
						Settings:   []byte("{}"),
					}

					err := settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
		},
		{
			name: "happy path multiple settings no filter",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Setting {
				noOfIDPs := 5
				settings := make([]*domain.Setting, noOfIDPs)
				for i := range noOfIDPs {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingType(i + 1),
						Settings:   []byte("{}"),
					}

					err := settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
		},
		func() test {
			id := gofakeit.Name()
			return test{
				name: "setting filter on id",
				testFunc: func(ctx context.Context, t *testing.T) []*domain.Setting {
					// create setting
					// this setting is created as an additional setting which should NOT
					// be returned in the results of this test case
					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeSecurity,
						Settings:   []byte("{}"),
					}
					err := settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					noOfIDPs := 1
					settings := make([]*domain.Setting, noOfIDPs)
					for i := range noOfIDPs {

						setting := domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
							// ID:         id,
							Type:     domain.SettingType(i + 1),
							Settings: []byte("{}"),
						}

						err := settingRepo.Create(ctx, &setting)
						require.NoError(t, err)

						settings[i] = &setting
						id = setting.ID
					}

					return settings
				},
				conditionClauses: func() []database.Condition {
					return []database.Condition{settingRepo.IDCondition(id)}
				},
			}
		}(),
		{
			name: "multiple settings filter on type",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Setting {
				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),

					Type:     domain.SettingTypeLogin,
					Settings: []byte("{}"),
				}
				err := settingRepo.Create(ctx, &setting)
				require.NoError(t, err)

				noOfIDPs := 1
				settings := make([]*domain.Setting, noOfIDPs)
				for i := range noOfIDPs {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypePasswordExpiry,
						Settings:   []byte("{}"),
					}

					err := settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					settings[i] = &setting
				}

				return settings
			},
			conditionClauses: func() []database.Condition {
				return []database.Condition{settingRepo.TypeCondition(domain.SettingTypePasswordExpiry)}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithoutCancel(t.Context())

			t.Cleanup(func() {
				_, err := pool.Exec(ctx, "DELETE FROM zitadel.settings")
				require.NoError(t, err)
			})

			settings := tt.testFunc(ctx, t)

			var conditions []database.Condition
			if tt.conditionClauses != nil {
				conditions = tt.conditionClauses()
			}

			// check setting values
			returnedIDPs, err := settingRepo.List(ctx,
				conditions...,
			)
			require.NoError(t, err)
			if tt.noIDPsReturned {
				assert.Nil(t, returnedIDPs)
				return
			}

			assert.Equal(t, len(settings), len(returnedIDPs))
			for i, setting := range settings {

				assert.Equal(t, returnedIDPs[i].InstanceID, setting.InstanceID)
				assert.Equal(t, returnedIDPs[i].OrgID, setting.OrgID)
				assert.Equal(t, returnedIDPs[i].Type, setting.Type)
				assert.Equal(t, returnedIDPs[i].Settings, setting.Settings)
			}
		})
	}
}

func TestDeleteSetting(t *testing.T) {
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
	instanceRepo := repository.InstanceRepository(pool)
	err := instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name            string
		testFunc        func(ctx context.Context, t *testing.T)
		settingType     domain.SettingType
		noOfDeletedRows int64
	}
	tests := []test{
		func() test {
			var noOfIDPs int64 = 1
			return test{
				name: "happy path delete",
				testFunc: func(ctx context.Context, t *testing.T) {
					for range noOfIDPs {
						setting := domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
							Type:       domain.SettingTypeLogin,
							Settings:   []byte("{}"),
						}

						err := settingRepo.Create(ctx, &setting)
						require.NoError(t, err)
					}
				},
				noOfDeletedRows: noOfIDPs,
				settingType:     domain.SettingTypeLogin,
			}
		}(),
		func() test {
			return test{
				name: "deleted setting with wrong type",
				testFunc: func(ctx context.Context, t *testing.T) {
					noOfIDPs := 1
					for range noOfIDPs {
						setting := domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
							ID:         gofakeit.Name(),

							Type:     domain.SettingTypeLogin,
							Settings: []byte("{}"),
						}

						err := settingRepo.Create(ctx, &setting)
						require.NoError(t, err)
					}

					// delete organization
					affectedRows, err := settingRepo.Delete(ctx,
						instanceId,
						&orgId,
						domain.SettingTypeLogin,
					)
					assert.Equal(t, int64(1), affectedRows)
					require.NoError(t, err)
				},
				// this test should return 0 affected rows as the setting was already deleted
				noOfDeletedRows: 0,
				settingType:     domain.SettingTypePasswordComplexity,
			}
		}(),
		func() test {
			return test{
				name: "deleted already deleted setting",
				testFunc: func(ctx context.Context, t *testing.T) {
					noOfIDPs := 1
					for range noOfIDPs {
						setting := domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
							ID:         gofakeit.Name(),

							Type:     domain.SettingTypeLogin,
							Settings: []byte("{}"),
						}

						err := settingRepo.Create(ctx, &setting)
						require.NoError(t, err)
					}

					// delete organization
					affectedRows, err := settingRepo.Delete(ctx,
						instanceId,
						&orgId,
						domain.SettingTypeLogin,
					)
					assert.Equal(t, int64(1), affectedRows)
					require.NoError(t, err)
				},
				// this test should return 0 affected rows as the setting was already deleted
				noOfDeletedRows: 0,
				settingType:     domain.SettingTypeLogin,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			if tt.testFunc != nil {
				tt.testFunc(ctx, t)
			}

			// delete setting
			noOfDeletedRows, err := settingRepo.Delete(ctx,
				instanceId,
				&orgId,
				tt.settingType,
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check setting was deleted
			organization, err := settingRepo.Get(ctx,
				instanceId,
				&orgId,
				tt.settingType,
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, organization)
		})
	}
}

func TestDeleteSettingsForInstance(t *testing.T) {
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
	instanceRepo := repository.InstanceRepository(pool)
	err := instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name            string
		testFunc        func(ctx context.Context, t *testing.T)
		noOfDeletedRows int64
	}
	tests := []test{
		func() test {
			return test{
				name: "delete all settings for instance",
				testFunc: func(ctx context.Context, t *testing.T) {
					// create with org
					err = settingRepo.CreateLogin(ctx, &domain.LoginSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateLabel(ctx, &domain.LabelSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreatePasswordComplexity(ctx, &domain.PasswordComplexitySetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreatePasswordExpiry(ctx, &domain.PasswordExpirySetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateLockout(ctx, &domain.LockoutSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateSecurity(ctx, &domain.SecuritySetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateDomain(ctx, &domain.DomainSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					// create without org
					err = settingRepo.CreateLogin(ctx, &domain.LoginSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							// OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateLabel(ctx, &domain.LabelSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							// OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreatePasswordComplexity(ctx, &domain.PasswordComplexitySetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							// OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreatePasswordExpiry(ctx, &domain.PasswordExpirySetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							// OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateLockout(ctx, &domain.LockoutSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							// OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateSecurity(ctx, &domain.SecuritySetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							// OrgID:      &orgId,
						},
					})
					require.NoError(t, err)
				},
				noOfDeletedRows: 13,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			if tt.testFunc != nil {
				tt.testFunc(ctx, t)
			}

			// delete settings
			noOfDeletedRows, err := settingRepo.DeleteSettingsForInstance(
				ctx,
				instanceId,
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)
		})
	}
}

func TestDeleteSettingsForOrg(t *testing.T) {
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
	instanceRepo := repository.InstanceRepository(pool)
	err := instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// create org
	orgId := gofakeit.Name()
	org := domain.Organization{
		ID:         orgId,
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name            string
		testFunc        func(ctx context.Context, t *testing.T)
		noOfDeletedRows int64
	}
	tests := []test{
		func() test {
			return test{
				name: "delete all settings for instance",
				testFunc: func(ctx context.Context, t *testing.T) {
					// create with org
					err = settingRepo.CreateLogin(ctx, &domain.LoginSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateLabel(ctx, &domain.LabelSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreatePasswordComplexity(ctx, &domain.PasswordComplexitySetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreatePasswordExpiry(ctx, &domain.PasswordExpirySetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateLockout(ctx, &domain.LockoutSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateSecurity(ctx, &domain.SecuritySetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)

					err = settingRepo.CreateDomain(ctx, &domain.DomainSetting{
						Setting: &domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
						},
					})
					require.NoError(t, err)
				},
				noOfDeletedRows: 7,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			if tt.testFunc != nil {
				tt.testFunc(ctx, t)
			}

			// delete settings
			noOfDeletedRows, err := settingRepo.DeleteSettingsForOrg(
				ctx,
				orgId,
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)
		})
	}
}

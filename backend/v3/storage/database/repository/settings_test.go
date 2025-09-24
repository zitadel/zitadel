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
				Settings:  []byte("{}"),
				CreatedAt: now,
				UpdatedAt: &now,
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
					CreatedAt:  now,
					UpdatedAt:  &now,
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
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(ctx, pool, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(ctx, pool, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository(pool)
					setting := domain.Setting{
						InstanceID: newInstId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypeLockout,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(ctx, pool, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(ctx, pool, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository(pool)
					setting := domain.Setting{
						InstanceID: newInstId,
						OrgID:      &newOrgId,
						Type:       domain.SettingTypeLockout,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(ctx, pool, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(ctx, pool, &org)
					require.NoError(t, err)

					// create setting
					settingRepo := repository.SettingsRepository(pool)
					setting := domain.Setting{
						InstanceID: newInstId,
						OrgID:      &org.ID,
						Type:       domain.SettingTypePasswordExpiry,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
					err = organizationRepo.Create(ctx, pool, &org)
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
				Settings:   []byte("{}"),
				CreatedAt:  now,
				UpdatedAt:  &now,
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
			// beforeCreate := time.Now()
			err = settingRepo.Create(ctx, setting)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			// afterCreate := time.Now()

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
			assert.Equal(t, tt.setting.Type, setting.Type)

			assert.Equal(t, setting.CreatedAt.UTC(), now.UTC())
			assert.Equal(t, setting.UpdatedAt.UTC(), now.UTC())
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
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(ctx, pool, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(ctx, pool, &org)
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
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(ctx, pool, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(ctx, pool, &org)
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
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(ctx, pool, &instance)
					assert.Nil(t, err)

					// create org
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(ctx, pool, &org)
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
					err = organizationRepo.Create(ctx, pool, &org)
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
			assert.WithinRange(t, *setting.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

func TestCreateSpecificSettingError(t *testing.T) {
	settingRepo := repository.SettingsRepository(pool)
	tests := []struct {
		name     string
		testFunc func(ctx context.Context) error
		err      error
	}{
		{
			name: "create login no settings error",
			testFunc: func(ctx context.Context) error {
				err := settingRepo.CreateLogin(ctx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create label no settings error",
			testFunc: func(ctx context.Context) error {
				err := settingRepo.CreateLabel(ctx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create label state not set",
			testFunc: func(ctx context.Context) error {
				err := settingRepo.CreateLabel(ctx, &domain.LabelSetting{
					Setting: &domain.Setting{},
				})
				return err
			},
			err: repository.ErrLabelStateMustBeDefined,
		},
		{
			name: "create password complexity no settings error",
			testFunc: func(ctx context.Context) error {
				err := settingRepo.CreatePasswordComplexity(ctx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create password expiry no settings error",
			testFunc: func(ctx context.Context) error {
				err := settingRepo.CreatePasswordExpiry(ctx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create lockout no settings error",
			testFunc: func(ctx context.Context) error {
				err := settingRepo.CreateLockout(ctx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create security no settings error",
			testFunc: func(ctx context.Context) error {
				err := settingRepo.CreateSecurity(ctx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create domain no settings error",
			testFunc: func(ctx context.Context) error {
				err := settingRepo.CreateDomain(ctx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
		{
			name: "create org no settings error",
			testFunc: func(ctx context.Context) error {
				err := settingRepo.CreateOrg(ctx, nil)
				return err
			},
			err: repository.ErrSettingObjectMustNotBeNil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			err := tt.testFunc(ctx)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

// NOTE all updated functions are just a wrapper around repository.updateSetting()
// for the sake of testing, UpdateLabel was used, but the underlying code is the same
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

	settingRepo := repository.SettingsRepository(pool)

	tests := []struct {
		name         string
		testFunc     func(ctx context.Context, t *testing.T) *domain.LabelSetting
		rowsAffected int64
		err          error
	}{
		{
			name: "happy path update settings",
			testFunc: func(ctx context.Context, t *testing.T) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
						IsDefault:           false,
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
				err := settingRepo.CreateLabel(ctx, &setting)
				require.NoError(t, err)

				// update settings values
				setting.Settings.IsDefault = true
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
			testFunc: func(ctx context.Context, t *testing.T) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						// OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
						IsDefault:           false,
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
				err := settingRepo.CreateLabel(ctx, &setting)
				require.NoError(t, err)

				// update settings values
				setting.Settings.IsDefault = true
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
			testFunc: func(ctx context.Context, t *testing.T) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: "non-existent-instanceID",
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
						IsDefault:           false,
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
		},
		{
			name: "update setting with non existent org id",
			testFunc: func(ctx context.Context, t *testing.T) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      gu.Ptr("non-existent-orgID"),
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
						IsDefault:           false,
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
		},
		{
			name: "update setting with wrong type",
			testFunc: func(ctx context.Context, t *testing.T) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LabelSettings{
						IsDefault:           false,
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
				err := settingRepo.CreateLabel(ctx, &setting)
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
			rowsAffected, err := settingRepo.UpdateLabel(
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
			updatedSetting, err := settingRepo.GetLabel(
				ctx,
				setting.InstanceID,
				setting.OrgID,
				*setting.LabelState,
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

// NOTE all the Get.*() functions are a wrapper around Get()
func TestGetSetting(t *testing.T) {
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

	// create setting
	// this setting is created as an additional org which should NOT
	// be returned in the results of the tests
	prexistingSetting := domain.Setting{
		InstanceID: instanceId,
		OrgID:      &orgId,
		ID:         gofakeit.Name(),
		Type:       domain.SettingTypePasswordExpiry,
		Settings:   []byte("{}"),
		CreatedAt:  now,
		UpdatedAt:  &now,
	}
	settingRepo := repository.SettingsRepository(pool)
	err = settingRepo.Create(t.Context(), &prexistingSetting)
	require.NoError(t, err)

	type test struct {
		name                       string
		testFunc                   func(ctx context.Context, t *testing.T) *domain.Setting
		settingIdentifierCondition database.Condition
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
					CreatedAt:  now,
					UpdatedAt:  &now,
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
					Type:      domain.SettingTypeLogin,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
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

					Type:      domain.SettingTypeDomain,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
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
						CreatedAt:  now,
						UpdatedAt:  &now,
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
						CreatedAt:  now,
						UpdatedAt:  &now,
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
			returnedSetting, err := settingRepo.Get(ctx,
				setting.InstanceID,
				setting.OrgID,
				setting.Type,
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
func TestCreateGetLoginPolicySetting(t *testing.T) {
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.LoginSetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(ctx context.Context, t *testing.T) *domain.LoginSetting {
				setting := domain.LoginSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LoginSettings{
						IsDefault:                  true,
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
						PasswordCheckLifetime:      time.Duration(time.Second * 50),
						ExternalLoginCheckLifetime: time.Duration(time.Second * 50),
						MFAInitSkipLifetime:        time.Duration(time.Second * 50),
						SecondFactorCheckLifetime:  time.Duration(time.Second * 50),
						MultiFactorCheckLifetime:   time.Duration(time.Second * 50),
					},
				}

				err := settingRepo.CreateLogin(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			// t.Cleanup(func() {
			// 	_, err := pool.Exec(context.Background(), "DELETE FROM zitadel.settings")
			// 	require.NoError(t, err)
			// })

			var setting *domain.LoginSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			}

			// get setting
			returnedSetting, err := settingRepo.GetLogin(ctx,
				setting.InstanceID,
				setting.OrgID,
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
func TestGetLabelPolicySetting(t *testing.T) {
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.LabelSetting
		err      error
	}

	tests := []test{
		{
			name: "happy path, state preview",
			testFunc: func(ctx context.Context, t *testing.T) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
					},
					Settings: domain.LabelSettings{
						IsDefault:           false,
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

				err := settingRepo.CreateLabel(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
		{
			name: "happy path, state activated",
			testFunc: func(ctx context.Context, t *testing.T) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStateActivated),
					},
					Settings: domain.LabelSettings{
						IsDefault:           false,
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

				err := settingRepo.CreateLabel(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
		{
			name: "get label policy using wrong state",
			testFunc: func(ctx context.Context, t *testing.T) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStateActivated),
					},
				}

				err := settingRepo.CreateLabel(ctx, &setting)
				require.NoError(t, err)

				setting.LabelState = gu.Ptr(domain.LabelStatePreview)
				return &setting
			},
			err: new(database.NoRowFoundError),
		},
		{
			name: "get non existent label policy",
			testFunc: func(ctx context.Context, t *testing.T) *domain.LabelSetting {
				setting := domain.LabelSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
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
			ctx := t.Context()
			t.Cleanup(func() {
				_, err := pool.Exec(context.Background(), "DELETE FROM zitadel.settings")
				require.NoError(t, err)
			})

			var setting *domain.LabelSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			}

			// get setting
			returnedSetting, err := settingRepo.GetLabel(ctx,
				setting.InstanceID,
				setting.OrgID,
				*setting.LabelState,
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
func TestCreateGetPasswordComplexityPolicySetting(t *testing.T) {
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.PasswordComplexitySetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordComplexitySetting {
				setting := domain.PasswordComplexitySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordComplexitySettings{
						IsDefault:    true,
						MinLength:    89,
						HasLowercase: true,
						HasUppercase: true,
						HasNumber:    true,
						HasSymbol:    true,
					},
				}

				err := settingRepo.CreatePasswordComplexity(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			var setting *domain.PasswordComplexitySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			}

			// get setting
			returnedSetting, err := settingRepo.GetPasswordComplexity(ctx,
				setting.InstanceID,
				setting.OrgID,
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

func TestCreateGetPasswordExpiryPolicySetting(t *testing.T) {
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(ctx context.Context, t *testing.T) *domain.PasswordExpirySetting {
				setting := domain.PasswordExpirySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.PasswordExpirySettings{
						IsDefault:      true,
						ExpireWarnDays: 30,
						MaxAgeDays:     30,
					},
				}

				err := settingRepo.CreatePasswordExpiry(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			var setting *domain.PasswordExpirySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			}

			// get setting
			returnedSetting, err := settingRepo.GetPasswordExpiry(ctx,
				setting.InstanceID,
				setting.OrgID,
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

func TestCreateGetLockoutPolicySetting(t *testing.T) {
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.LockoutSetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(ctx context.Context, t *testing.T) *domain.LockoutSetting {
				setting := domain.LockoutSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.LockoutSettings{
						IsDefault:           true,
						MaxPasswordAttempts: 50,
						MaxOTPAttempts:      50,
						ShowLockOutFailures: true,
					},
				}

				err := settingRepo.CreateLockout(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			var setting *domain.LockoutSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			}

			// get setting
			returnedSetting, err := settingRepo.GetLockout(ctx,
				setting.InstanceID,
				setting.OrgID,
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

func TestCreateGetSecurityPolicySetting(t *testing.T) {
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.SecuritySetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(ctx context.Context, t *testing.T) *domain.SecuritySetting {
				setting := domain.SecuritySetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
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

				err := settingRepo.CreateSecurity(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			var setting *domain.SecuritySetting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			}

			// get setting
			returnedSetting, err := settingRepo.GetSecurity(ctx,
				setting.InstanceID,
				setting.OrgID,
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

func TestCreateGetDomainPolicySetting(t *testing.T) {
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.DomainSetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(ctx context.Context, t *testing.T) *domain.DomainSetting {
				setting := domain.DomainSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.DomainSettings{
						IsDefault:                              true,
						UserLoginMustBeDomain:                  true,
						ValidateOrgDomains:                     true,
						SMTPSenderAddressMatchesInstanceDomain: true,
					},
				}

				err := settingRepo.CreateDomain(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			var setting *domain.DomainSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			}

			// get setting
			returnedSetting, err := settingRepo.GetDomain(ctx,
				setting.InstanceID,
				setting.OrgID,
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

func TestCreateGetOrgPolicySetting(t *testing.T) {
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.OrgSetting
		err      error
	}
	tests := []test{
		{
			name: "happy path",
			testFunc: func(ctx context.Context, t *testing.T) *domain.OrgSetting {
				setting := domain.OrgSetting{
					Setting: &domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						Type:       domain.SettingTypeLabel,
						LabelState: gu.Ptr(domain.LabelStatePreview),
						Settings:   []byte("{}"),
					},
					Settings: domain.OrgSettings{
						OrganizationScopedUsernames:    true,
						OldOrganizationScopedUsernames: true,
						UsernameChanges:                []string{"value"},
					},
				}

				err := settingRepo.CreateOrg(ctx, &setting)
				require.NoError(t, err)
				return &setting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			var setting *domain.OrgSetting
			if tt.testFunc != nil {
				setting = tt.testFunc(ctx, t)
			}

			// get setting
			returnedSetting, err := settingRepo.GetOrg(ctx,
				setting.InstanceID,
				setting.OrgID,
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

// gocognit linting fails due to number of test cases
// and the fact that each test case has a testFunc()
//
//nolint:gocognit
func TestListSetting(t *testing.T) {
	now := time.Now()
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
	instanceRepo := repository.InstanceRepository()
	err = instanceRepo.Create(ctx, pool, &instance)
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name               string
		testFunc           func(ctx context.Context, t *testing.T) []*domain.Setting
		conditionClauses   func() []database.Condition
		noSettingsReturned bool
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
				err = instanceRepo.Create(ctx, pool, &instance)
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
				err = organizationRepo.Create(ctx, pool, &org)
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
					CreatedAt:  now,
					UpdatedAt:  &now,
				}
				err := settingRepo.Create(ctx, &setting)
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
						Settings:  []byte("{}"),
						CreatedAt: now,
						UpdatedAt: &now,
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
				err = instanceRepo.Create(ctx, pool, &instance)
				require.NoError(t, err)

				// create setting
				// this setting is created as an additional setting which should NOT
				// be returned in the results of this test case
				setting := domain.Setting{
					InstanceID: newInstanceId,
					// OrgID:      &newOrgId,
					ID:        gofakeit.Name(),
					Type:      domain.SettingTypeSecurity,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
				}
				err := settingRepo.Create(ctx, &setting)
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
						Settings:  []byte("{}"),
						CreatedAt: now,
						UpdatedAt: &now,
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
				organizationRepo := repository.OrganizationRepository()
				err = organizationRepo.Create(ctx, pool, &org)
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
					CreatedAt:  now,
					UpdatedAt:  &now,
				}
				err := settingRepo.Create(ctx, &setting)
				require.NoError(t, err)

				noOfSettings := 5
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingType(i + 1),
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
				noOfSettings := 1
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingType(i + 1),
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
				noOfSettings := 5
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingType(i + 1),
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
						CreatedAt:  now,
						UpdatedAt:  &now,
					}
					err := settingRepo.Create(ctx, &setting)
					require.NoError(t, err)

					noOfSettings := 1
					settings := make([]*domain.Setting, noOfSettings)
					for i := range noOfSettings {

						setting := domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
							// ID:         id,
							Type:      domain.SettingType(i + 1),
							Settings:  []byte("{}"),
							CreatedAt: now,
							UpdatedAt: &now,
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

					Type:      domain.SettingTypeLogin,
					Settings:  []byte("{}"),
					CreatedAt: now,
					UpdatedAt: &now,
				}
				err := settingRepo.Create(ctx, &setting)
				require.NoError(t, err)

				noOfSettings := 1
				settings := make([]*domain.Setting, noOfSettings)
				for i := range noOfSettings {

					setting := domain.Setting{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						Type:       domain.SettingTypePasswordExpiry,
						Settings:   []byte("{}"),
						CreatedAt:  now,
						UpdatedAt:  &now,
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
			returnedSettings, err := settingRepo.List(ctx,
				conditions...,
			)
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
				assert.Equal(t, returnedSettings[i].Settings, setting.Settings)
			}
		})
	}
}

func TestDeleteSetting(t *testing.T) {
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

	settingRepo := repository.SettingsRepository(pool)

	type test struct {
		name            string
		testFunc        func(ctx context.Context, t *testing.T)
		settingType     domain.SettingType
		noOfDeletedRows int64
	}
	tests := []test{
		func() test {
			var noOfSettings int64 = 1
			return test{
				name: "happy path delete",
				testFunc: func(ctx context.Context, t *testing.T) {
					for range noOfSettings {
						setting := domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
							Type:       domain.SettingTypeLogin,
							Settings:   []byte("{}"),
							CreatedAt:  now,
							UpdatedAt:  &now,
						}

						err := settingRepo.Create(ctx, &setting)
						require.NoError(t, err)
					}
				},
				noOfDeletedRows: noOfSettings,
				settingType:     domain.SettingTypeLogin,
			}
		}(),
		func() test {
			return test{
				name: "deleted setting with wrong type",
				testFunc: func(ctx context.Context, t *testing.T) {
					noOfSettings := 1
					for range noOfSettings {
						setting := domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
							ID:         gofakeit.Name(),
							Type:       domain.SettingTypeLogin,
							Settings:   []byte("{}"),
							CreatedAt:  now,
							UpdatedAt:  &now,
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
					noOfSettings := 1
					for range noOfSettings {
						setting := domain.Setting{
							InstanceID: instanceId,
							OrgID:      &orgId,
							ID:         gofakeit.Name(),

							Type:      domain.SettingTypeLogin,
							Settings:  []byte("{}"),
							CreatedAt: now,
							UpdatedAt: &now,
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
							LabelState: gu.Ptr(domain.LabelStatePreview),
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
							LabelState: gu.Ptr(domain.LabelStatePreview),
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
							LabelState: gu.Ptr(domain.LabelStatePreview),
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

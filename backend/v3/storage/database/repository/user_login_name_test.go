package repository_test

import (
	"context"
	"fmt"
	"maps"
	"strconv"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestUserLoginNameManipulations(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := tx.Rollback(context.Background())
		if err != nil {
			t.Log("failed to rollback transaction:", err)
		}
	})

	orgDomainRepo := repository.OrganizationDomainRepository()
	userRepo := repository.UserRepository()
	settingsRepo := repository.DomainSettingsRepository()

	instanceID := createInstance(t, tx)
	err = settingsRepo.Set(t.Context(), tx, &domain.DomainSettings{
		Settings: domain.Settings{
			InstanceID: instanceID,
			Type:       domain.SettingTypeDomain,
		},
		DomainSettingsAttributes: domain.DomainSettingsAttributes{
			LoginNameIncludesDomain: gu.Ptr(false),
		},
	})
	require.NoError(t, err)

	org1 := createOrganization(t, tx, instanceID)
	err = orgDomainRepo.Add(t.Context(), tx, &domain.AddOrganizationDomain{
		InstanceID: instanceID,
		OrgID:      org1,
		Domain:     "verified.example.com",
		IsVerified: true,
	})
	require.NoError(t, err)
	err = orgDomainRepo.Add(t.Context(), tx, &domain.AddOrganizationDomain{
		InstanceID: instanceID,
		OrgID:      org1,
		Domain:     "primary.example.com",
		IsVerified: true,
		IsPrimary:  true,
	})
	require.NoError(t, err)
	err = orgDomainRepo.Add(t.Context(), tx, &domain.AddOrganizationDomain{
		InstanceID: instanceID,
		OrgID:      org1,
		Domain:     "unverified.example.com",
	})
	require.NoError(t, err)
	err = settingsRepo.Set(t.Context(), tx, &domain.DomainSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: &org1,
			Type:           domain.SettingTypeDomain,
		},
		DomainSettingsAttributes: domain.DomainSettingsAttributes{
			LoginNameIncludesDomain: gu.Ptr(true),
		},
	})
	require.NoError(t, err)

	org2 := createOrganization(t, tx, instanceID)
	err = orgDomainRepo.Add(t.Context(), tx, &domain.AddOrganizationDomain{
		InstanceID: instanceID,
		OrgID:      org2,
		Domain:     "primary.example2.com",
		IsVerified: true,
		IsPrimary:  true,
	})
	require.NoError(t, err)

	loginNames := setupLoginNameTestUsers(t, tx, instanceID, org1, org2)

	for _, tt := range []struct {
		name               string
		opts               []database.QueryOption
		expectedLoginNames map[string][]domain.LoginName
		before             func(t *testing.T, tx database.QueryExecutor)
	}{
		{
			name: "list all users",
			opts: []database.QueryOption{
				database.WithCondition(userRepo.InstanceIDCondition(instanceID)),
			},
			expectedLoginNames: maps.Clone(loginNames),
		},
		{
			name: "update usernames",
			before: func(t *testing.T, tx database.QueryExecutor) {
				for i := range 100 {
					_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID, strconv.Itoa(i)),
						userRepo.SetUsername("updated-user"+strconv.Itoa(i)),
					)
					require.NoError(t, err)
				}
			},
			opts: []database.QueryOption{
				database.WithCondition(userRepo.InstanceIDCondition(instanceID)),
			},
			expectedLoginNames: expectedLoginNames(loginNames, func(id string, loginNames []domain.LoginName) []domain.LoginName {
				updated := make([]domain.LoginName, len(loginNames))
				for i, ln := range loginNames {
					updated[i] = ln
					updated[i].LoginName = "updated-" + ln.LoginName
				}
				return updated
			}),
		},
		{
			name: "add org setting without impact",
			before: func(t *testing.T, tx database.QueryExecutor) {
				err := settingsRepo.Set(t.Context(), tx, &domain.DomainSettings{
					Settings: domain.Settings{
						InstanceID:     instanceID,
						OrganizationID: &org2,
						Type:           domain.SettingTypeDomain,
					},
					DomainSettingsAttributes: domain.DomainSettingsAttributes{
						RequireOrgDomainVerification: gu.Ptr(false),
					},
				})
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(userRepo.InstanceIDCondition(instanceID)),
			},
			expectedLoginNames: maps.Clone(loginNames),
		},
		{
			name: "add and update org setting without impact",
			before: func(t *testing.T, tx database.QueryExecutor) {
				err := settingsRepo.Set(t.Context(), tx, &domain.DomainSettings{
					Settings: domain.Settings{
						InstanceID:     instanceID,
						OrganizationID: &org2,
						Type:           domain.SettingTypeDomain,
					},
					DomainSettingsAttributes: domain.DomainSettingsAttributes{
						RequireOrgDomainVerification: gu.Ptr(false),
					},
				})
				require.NoError(t, err)

				err = settingsRepo.Set(t.Context(), tx, &domain.DomainSettings{
					Settings: domain.Settings{
						InstanceID:     instanceID,
						OrganizationID: &org2,
						Type:           domain.SettingTypeDomain,
					},
					DomainSettingsAttributes: domain.DomainSettingsAttributes{
						RequireOrgDomainVerification: gu.Ptr(true),
					},
				})
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(userRepo.InstanceIDCondition(instanceID)),
			},
			expectedLoginNames: maps.Clone(loginNames),
		},
		{
			name: "add without impact and update with impact org setting",
			before: func(t *testing.T, tx database.QueryExecutor) {
				err := settingsRepo.Set(t.Context(), tx, &domain.DomainSettings{
					Settings: domain.Settings{
						InstanceID:     instanceID,
						OrganizationID: &org2,
						Type:           domain.SettingTypeDomain,
					},
					DomainSettingsAttributes: domain.DomainSettingsAttributes{
						RequireOrgDomainVerification: gu.Ptr(false),
					},
				})
				require.NoError(t, err)

				err = settingsRepo.Set(t.Context(), tx, &domain.DomainSettings{
					Settings: domain.Settings{
						InstanceID:     instanceID,
						OrganizationID: &org2,
						Type:           domain.SettingTypeDomain,
					},
					DomainSettingsAttributes: domain.DomainSettingsAttributes{
						LoginNameIncludesDomain: gu.Ptr(true),
					},
				})
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(userRepo.InstanceIDCondition(instanceID)),
			},
			expectedLoginNames: expectedLoginNames(loginNames, func(id string, loginNames []domain.LoginName) []domain.LoginName {
				if len(loginNames) == 2 {
					return loginNames
				}
				return []domain.LoginName{
					{
						LoginName:   fmt.Sprintf("user%s@primary.example2.com", id),
						IsPreferred: true,
					},
				}
			}),
		},
		{
			name: "change org setting with impact",
			before: func(t *testing.T, tx database.QueryExecutor) {
				err := settingsRepo.Set(t.Context(), tx, &domain.DomainSettings{
					Settings: domain.Settings{
						InstanceID:     instanceID,
						OrganizationID: &org1,
						Type:           domain.SettingTypeDomain,
					},
					DomainSettingsAttributes: domain.DomainSettingsAttributes{
						LoginNameIncludesDomain: gu.Ptr(false),
					},
				})
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(userRepo.InstanceIDCondition(instanceID)),
			},
			expectedLoginNames: expectedLoginNames(loginNames, func(id string, loginNames []domain.LoginName) []domain.LoginName {
				if len(loginNames) == 1 {
					return loginNames
				}
				return []domain.LoginName{
					{
						LoginName:   "user" + id,
						IsPreferred: true,
					},
				}
			}),
		},
		{
			name: "reset org setting",
			before: func(t *testing.T, tx database.QueryExecutor) {
				_, err := settingsRepo.Delete(t.Context(), tx, settingsRepo.UniqueCondition(instanceID, &org1, domain.SettingTypeDomain, domain.SettingStateActive))
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(userRepo.InstanceIDCondition(instanceID)),
			},
			expectedLoginNames: expectedLoginNames(loginNames, func(id string, loginNames []domain.LoginName) []domain.LoginName {
				if len(loginNames) == 1 {
					return loginNames
				}
				return []domain.LoginName{
					{
						LoginName:   "user" + id,
						IsPreferred: true,
					},
				}
			}),
		},
		{
			name: "change instance setting",
			before: func(t *testing.T, tx database.QueryExecutor) {
				err := settingsRepo.Set(t.Context(), tx, &domain.DomainSettings{
					Settings: domain.Settings{
						InstanceID: instanceID,
						Type:       domain.SettingTypeDomain,
					},
					DomainSettingsAttributes: domain.DomainSettingsAttributes{
						LoginNameIncludesDomain: gu.Ptr(true),
					},
				})
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(userRepo.InstanceIDCondition(instanceID)),
			},
			expectedLoginNames: expectedLoginNames(loginNames, func(id string, loginNames []domain.LoginName) []domain.LoginName {
				if len(loginNames) == 2 {
					return loginNames
				}
				return []domain.LoginName{
					{
						LoginName:   fmt.Sprintf("user%s@primary.example2.com", id),
						IsPreferred: true,
					},
				}
			}),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			t.Cleanup(func() {
				err := savepoint.Rollback(context.Background())
				if err != nil {
					t.Log("failed to rollback savepoint:", err)
				}
			})

			if tt.before != nil {
				tt.before(t, savepoint)
			}

			users, err := userRepo.List(t.Context(), savepoint, tt.opts...)
			require.NoError(t, err)

			for _, user := range users {
				expected, ok := tt.expectedLoginNames[user.ID]
				require.True(t, ok, "unexpected user ID: %s", user.ID)
				assertLoginNames(t, expected, user.LoginNames)
				delete(tt.expectedLoginNames, user.ID)
			}
			assert.Empty(t, tt.expectedLoginNames, "some expected users were not found")
		})
	}
}

func setupLoginNameTestUsers(t *testing.T, tx database.QueryExecutor, instanceID, org1, org2 string) map[string][]domain.LoginName {
	userRepo := repository.UserRepository()

	expectedLoginNames := make(map[string][]domain.LoginName, 100)
	for i := range 100 {
		orgID := org1
		id := strconv.Itoa(i)

		expectedLoginNames[id] = []domain.LoginName{
			{
				LoginName:   fmt.Sprintf("user%d@verified.example.com", i),
				IsPreferred: false,
			},
			{
				LoginName:   fmt.Sprintf("user%d@primary.example.com", i),
				IsPreferred: true,
			},
		}

		if i%2 == 0 {
			orgID = org2
			expectedLoginNames[id] = []domain.LoginName{
				{
					LoginName:   "user" + id,
					IsPreferred: true,
				},
			}
		}
		err := userRepo.Create(t.Context(), tx, &domain.User{
			InstanceID:     instanceID,
			OrganizationID: orgID,
			ID:             id,
			Username:       "user" + id,
			State:          domain.UserStateActive,
			Machine: &domain.MachineUser{
				Name: "machine" + id,
			},
		})
		require.NoError(t, err)
	}
	return expectedLoginNames
}

func expectedLoginNames(loginNames map[string][]domain.LoginName, replace func(id string, loginNames []domain.LoginName) []domain.LoginName) map[string][]domain.LoginName {
	updated := maps.Clone(loginNames)
	for id, userLoginNames := range updated {
		updated[id] = replace(id, userLoginNames)
	}
	return updated
}

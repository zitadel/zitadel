package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestCreateHumanUser(t *testing.T) {
	// beforeCreate := time.Now()
	// tx, err := pool.Begin(t.Context(), nil)
	// require.NoError(t, err)
	// defer func() {
	// 	err := tx.Rollback(t.Context())
	// 	if err != nil {
	// 		t.Logf("error during rollback: %v", err)
	// 	}
	// }()
	tx := pool

	instanceRepo := repository.InstanceRepository()
	// create instance
	instanceID := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err := instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create organization
	orgID := gofakeit.UUID()
	orgRepo := repository.OrganizationRepository()
	organization := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), tx, &organization)
	require.NoError(t, err)

	userRepo := repository.UserRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, client database.QueryExecutor)
		user     *domain.Human
		opts     []database.QueryOption
		err      error
	}

	tests := []test{
		// user
		func() test {
			id := gofakeit.UUID()
			return test{
				name: "happy path",
				user: &domain.Human{
					User: domain.User{
						ID:                id,
						InstanceID:        instanceID,
						OrgID:             orgID,
						Username:          gofakeit.Username(),
						UsernameOrgUnique: true,
						State:             domain.UserStateActive,
					},
					FirstName:         gofakeit.FirstName(),
					LastName:          gofakeit.LastName(),
					NickName:          gofakeit.Username(),
					DisplayName:       gofakeit.Name(),
					PreferredLanguage: "en",
					Gender:            1,
					AvatarKey:         gofakeit.Animal(),
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(instanceID),
						userRepo.Human().OrgIDCondition(orgID),
						userRepo.Human().IDCondition(id),
					)),
				},
			}
		}(),
		{
			name: "no id",
			user: &domain.Human{
				User: domain.User{
					// ID:         id,
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			},
			err: new(database.CheckError),
		},
		{
			name: "no username",
			user: &domain.Human{
				User: domain.User{
					ID:         gofakeit.UUID(),
					InstanceID: instanceID,
					OrgID:      orgID,
					// Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			},
			err: new(database.CheckError),
		},
		func() test {
			id := gofakeit.UUID()
			user := &domain.Human{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			}
			return test{
				name: "same instance + org + username twice, with username_org_unique = true",
				user: user,
				testFunc: func(t *testing.T, client database.QueryExecutor) {
					_, err := userRepo.CreateHuman(t.Context(), pool, user)
					require.NoError(t, err)
				},
				err: new(database.UniqueError),
			}
		}(),
		func() test {
			id := gofakeit.UUID()
			user := &domain.Human{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: false,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			}
			_, err := userRepo.CreateHuman(t.Context(), pool, user)
			require.NoError(t, err)
			return test{
				name: "same instance + org + id twice, with username_org_unique = false",
				user: user,
				err:  new(database.UniqueError),
			}
		}(),
		func() test {
			id := gofakeit.UUID()
			newOrgID := gofakeit.UUID()

			user := &domain.Human{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			}
			return test{
				name: "same instance + username twice, different org with username_org_unique = true",
				user: user,
				testFunc: func(t *testing.T, client database.QueryExecutor) {
					// create organization
					orgRepo := repository.OrganizationRepository()
					organization := domain.Organization{
						ID:         newOrgID,
						Name:       gofakeit.Name(),
						InstanceID: instanceID,
						State:      domain.OrgStateActive,
					}
					err = orgRepo.Create(t.Context(), tx, &organization)
					require.NoError(t, err)

					_, err := userRepo.CreateHuman(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different org
					user.OrgID = orgID
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(instanceID),
						userRepo.Human().OrgIDCondition(orgID),
						userRepo.Human().IDCondition(id),
					)),
				},
			}
		}(),
		func() test {
			id := gofakeit.UUID()
			// create organization
			newOrgID := gofakeit.UUID()

			user := &domain.Human{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: false,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			}
			return test{
				name: "same instance + username twice, different org with username_org_unique = false",
				user: user,
				testFunc: func(t *testing.T, client database.QueryExecutor) {
					// create organization
					orgRepo := repository.OrganizationRepository()
					organization := domain.Organization{
						ID:         newOrgID,
						Name:       gofakeit.Name(),
						InstanceID: instanceID,
						State:      domain.OrgStateActive,
					}
					err = orgRepo.Create(t.Context(), tx, &organization)
					require.NoError(t, err)

					_, err := userRepo.CreateHuman(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different org
					user.OrgID = orgID
				},
				err: new(database.UniqueError),
			}
		}(),
		func() test {
			id := gofakeit.UUID()
			// create organization
			newOrgID := gofakeit.UUID()

			user := &domain.Human{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			}
			return test{
				name: "same username, instance with username_org_unique = true",
				user: user,
				testFunc: func(t *testing.T, client database.QueryExecutor) {
					// create organization
					orgRepo := repository.OrganizationRepository()
					organization := domain.Organization{
						ID:         newOrgID,
						Name:       gofakeit.Name(),
						InstanceID: instanceID,
						State:      domain.OrgStateActive,
					}
					err = orgRepo.Create(t.Context(), tx, &organization)
					require.NoError(t, err)

					_, err = userRepo.CreateHuman(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different org
					user.OrgID = orgID
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(instanceID),
						userRepo.Human().OrgIDCondition(orgID),
						userRepo.Human().IDCondition(id),
					)),
				},
			}
		}(),
		func() test {
			id := gofakeit.UUID()
			newOrgID := gofakeit.UUID()

			user := &domain.Human{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: false,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			}
			return test{
				name: "same username, instance with username_org_unique = true",
				user: user,
				testFunc: func(t *testing.T, client database.QueryExecutor) {
					// create organization
					orgRepo := repository.OrganizationRepository()
					organization := domain.Organization{
						ID:         newOrgID,
						Name:       gofakeit.Name(),
						InstanceID: instanceID,
						State:      domain.OrgStateActive,
					}
					err = orgRepo.Create(t.Context(), tx, &organization)
					require.NoError(t, err)

					_, err = userRepo.CreateHuman(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different instance
					user.InstanceID = instanceID
					user.OrgID = orgID
				},
				err: new(database.UniqueError),
			}
		}(),
		func() test {
			id := gofakeit.UUID()
			newInstanceID := gofakeit.Name()
			newOrgID := gofakeit.UUID()

			user := &domain.Human{
				User: domain.User{
					ID:                id,
					InstanceID:        newInstanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			}
			return test{
				name: "same username, different instance with username_org_unique = true",
				user: user,
				testFunc: func(t *testing.T, client database.QueryExecutor) {
					// create instance
					instance := domain.Instance{
						ID:              newInstanceID,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					err := instanceRepo.Create(t.Context(), tx, &instance)
					require.NoError(t, err)
					// create organization
					orgRepo := repository.OrganizationRepository()
					organization := domain.Organization{
						ID:         newOrgID,
						Name:       gofakeit.Name(),
						InstanceID: newInstanceID,
						State:      domain.OrgStateActive,
					}
					err = orgRepo.Create(t.Context(), tx, &organization)
					require.NoError(t, err)

					_, err = userRepo.CreateHuman(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different instance
					user.InstanceID = instanceID
					user.OrgID = orgID
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(instanceID),
						userRepo.Human().OrgIDCondition(orgID),
						userRepo.Human().IDCondition(id),
					)),
				},
			}
		}(),
		func() test {
			id := gofakeit.UUID()
			// // create instance
			newInstanceID := gofakeit.Name()
			newOrgID := gofakeit.UUID()

			user := &domain.Human{
				User: domain.User{
					ID:                id,
					InstanceID:        newInstanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: false,
					State:             domain.UserStateActive,
				},
				FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			}
			return test{
				name: "same username, different instance with username_org_unique = false",
				user: user,
				testFunc: func(t *testing.T, client database.QueryExecutor) {
					// create instance
					instance := domain.Instance{
						ID:              newInstanceID,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					err := instanceRepo.Create(t.Context(), tx, &instance)
					require.NoError(t, err)
					// create organization
					orgRepo := repository.OrganizationRepository()
					organization := domain.Organization{
						ID:         newOrgID,
						Name:       gofakeit.Name(),
						InstanceID: newInstanceID,
						State:      domain.OrgStateActive,
					}
					err = orgRepo.Create(t.Context(), tx, &organization)
					require.NoError(t, err)

					_, err = userRepo.CreateHuman(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different instance
					user.InstanceID = instanceID
					user.OrgID = orgID
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(instanceID),
						userRepo.Human().OrgIDCondition(orgID),
						userRepo.Human().IDCondition(id),
					)),
				},
				// err: new(database.UniqueError),
			}
		}(),
		// human
		{
			name: "no first name",
			user: &domain.Human{
				User: domain.User{
					ID:                gofakeit.UUID(),
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				// FirstName:         gofakeit.FirstName(),
				LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			},
			err: new(database.CheckError),
		},
		{
			name: "no last name",
			user: &domain.Human{
				User: domain.User{
					ID:                gofakeit.UUID(),
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				FirstName: gofakeit.FirstName(),
				// LastName:          gofakeit.LastName(),
				NickName:          gofakeit.Username(),
				DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			},
			err: new(database.CheckError),
		},
		{
			name: "no display name",
			user: &domain.Human{
				User: domain.User{
					ID:                gofakeit.UUID(),
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				NickName:  gofakeit.Username(),
				// DisplayName:       gofakeit.Name(),
				PreferredLanguage: "en",
				Gender:            1,
				AvatarKey:         gofakeit.Animal(),
			},
			err: new(database.CheckError),
		},
		{
			name: "no preferred language",
			user: &domain.Human{
				User: domain.User{
					ID:                gofakeit.UUID(),
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				NickName:    gofakeit.Username(),
				DisplayName: gofakeit.Name(),
				// PreferredLanguage: "en",
				Gender:    1,
				AvatarKey: gofakeit.Animal(),
			},
			err: new(database.CheckError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// savepoint, err := tx.Begin(t.Context())
			// require.NoError(t, err)
			// defer func() {
			// 	err = savepoint.Rollback(t.Context())
			// 	if err != nil {
			// 		t.Logf("error during rollback: %v", err)
			// 	}
			// }()
			savepoint := tx

			defer func() {
				_, err := tx.Exec(t.Context(), "DELETE FROM zitadel.identity_providers")
				require.NoError(t, err)
			}()

			if tt.testFunc != nil {
				tt.testFunc(t, savepoint)
			}

			// create user
			beforeCreate := time.Now()
			user, err := userRepo.CreateHuman(t.Context(), savepoint, tt.user)
			if err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			afterCreate := time.Now()

			createdUser, err := userRepo.GetHuman(t.Context(), savepoint, tt.opts...)
			require.NoError(t, err)

			// user
			assert.Equal(t, tt.user.InstanceID, createdUser.InstanceID)
			assert.Equal(t, tt.user.OrgID, createdUser.OrgID)
			assert.Equal(t, tt.user.State, createdUser.State)
			assert.Equal(t, tt.user.ID, createdUser.ID)
			assert.Equal(t, tt.user.Username, createdUser.Username)
			assert.Equal(t, tt.user.UsernameOrgUnique, createdUser.UsernameOrgUnique)
			assert.Equal(t, tt.user.State, createdUser.State)

			// // human
			assert.Equal(t, tt.user.FirstName, createdUser.FirstName)
			assert.Equal(t, tt.user.LastName, createdUser.LastName)
			assert.Equal(t, tt.user.NickName, createdUser.NickName)
			assert.Equal(t, tt.user.DisplayName, createdUser.DisplayName)
			assert.Equal(t, tt.user.PreferredLanguage, createdUser.PreferredLanguage)
			assert.Equal(t, tt.user.Gender, createdUser.Gender)
			assert.Equal(t, tt.user.AvatarKey, createdUser.AvatarKey)

			assert.WithinRange(t, user.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, user.CreatedAt, beforeCreate, afterCreate)
		})
	}
}

func TestUpdateHumanUser(t *testing.T) {
	beforeUpdate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Errorf("error rolling back transaction: %v", err)
		}
	}()

	// create instance
	instanceID := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceID,
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
	orgID := gofakeit.Name()
	org := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	organizationRepo := repository.OrganizationRepository()
	err = organizationRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	userRepo := repository.UserRepository()

	type test struct {
		name         string
		testFunc     func(t *testing.T, tx database.QueryExecutor) *domain.Human
		condition    database.Condition
		changes      []database.Change
		rowsAffected int64
		err          error
	}
	tests := []*test{
		// user
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update username",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.FirstName(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.Username = "new_username"
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetUsername("new_username")},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "update username invalid value",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.FirstName(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes: []database.Change{userRepo.Human().SetUsername("")},
				err:     new(database.CheckError),
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update username_org_unique",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.FirstName(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.UsernameOrgUnique = true
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetUsernameOrgUnique(true)},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update state",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.FirstName(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.State = domain.UserStateInactive
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetState(domain.UserStateInactive)},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "update invalid state",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.FirstName(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.State = domain.UserStateInactive
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes: []database.Change{userRepo.Human().SetState(domain.UserStateUnspecified)},
				err:     new(database.CheckError),
			}
			return test
		}(),
		// human
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update first name",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         "first_name",
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.FirstName = "new_first_name"
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetFirstName("new_first_name")},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "update first name to invalid value",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         "first_name",
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes: []database.Change{userRepo.Human().SetFirstName("")},
				err:     new(database.CheckError),
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update last name",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.Username(),
						LastName:          "last_name",
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.LastName = "new_last_name"
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetLastName("new_last_name")},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "update last name to invalid value",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.Username(),
						LastName:          "last_name",
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes: []database.Change{userRepo.Human().SetLastName("")},
				err:     new(database.CheckError),
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update nick name",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.Username(),
						LastName:          gofakeit.LastName(),
						NickName:          "nick_name",
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.NickName = "new_nick_name"
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetNickName("new_nick_name")},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update display name",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.Username(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       "display_name",
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.DisplayName = "new_display_name"
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetDisplayName("new_display_name")},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "update display name to invalid value",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.Username(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       "display_name",
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetDisplayName("")},
				err:          new(database.CheckError),
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update gender",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.Username(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.Gender = 2
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetGender(2)},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update preferred language",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.Username(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.PreferredLanguage = "de"
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetPreferredLanguage("de")},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "update preferred language invalid value",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.Username(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.PreferredLanguage = "de"
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes: []database.Change{userRepo.Human().SetPreferredLanguage("")},
				err:     new(database.CheckError),
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			avatarKey := gofakeit.Animal()
			test := &test{
				name: "happy path update avatar key",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Human {
					user := &domain.Human{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						FirstName:         gofakeit.Username(),
						LastName:          gofakeit.LastName(),
						NickName:          gofakeit.Username(),
						DisplayName:       gofakeit.Name(),
						PreferredLanguage: "en",
						Gender:            1,
						AvatarKey:         gofakeit.Animal(),
					}

					user, err := userRepo.CreateHuman(t.Context(), tx, user)
					require.NoError(t, err)
					user.AvatarKey = avatarKey
					return user
				},
				condition: database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Human().SetAvatarKey(avatarKey)},
				rowsAffected: 1,
			}
			return test
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err := tx.Rollback(t.Context())
				if err != nil {
					t.Errorf("error rolling back savepoint: %v", err)
				}
			}()

			createdUser := tt.testFunc(t, tx)

			// update user
			rowsAffected, err := userRepo.UpdateHuman(
				t.Context(),
				tx,
				tt.condition,
				tt.changes...,
			)
			afterUpdate := time.Now()

			if err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)

			if rowsAffected == 0 {
				return
			}
			assert.Equal(t, tt.rowsAffected, rowsAffected)

			// check user values
			user, err := userRepo.GetHuman(
				t.Context(),
				tx,
				database.WithCondition(database.And(
					tt.condition,
				)),
			)
			require.NoError(t, err)

			// user
			assert.Equal(t, createdUser.InstanceID, user.InstanceID)
			assert.Equal(t, createdUser.OrgID, user.OrgID)
			assert.Equal(t, createdUser.State, user.State)
			assert.Equal(t, createdUser.ID, user.ID)
			assert.Equal(t, createdUser.Username, user.Username)
			assert.Equal(t, createdUser.UsernameOrgUnique, user.UsernameOrgUnique)
			assert.Equal(t, createdUser.State, user.State)

			// // human
			assert.Equal(t, createdUser.FirstName, user.FirstName)
			assert.Equal(t, createdUser.LastName, user.LastName)
			assert.Equal(t, createdUser.NickName, user.NickName)
			assert.Equal(t, createdUser.DisplayName, user.DisplayName)
			assert.Equal(t, createdUser.PreferredLanguage, user.PreferredLanguage)
			assert.Equal(t, createdUser.Gender, user.Gender)
			assert.Equal(t, createdUser.AvatarKey, user.AvatarKey)

			assert.WithinRange(t, user.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

func TestListHumanUser(t *testing.T) {
	// tx, err := pool.Begin(t.Context(), nil)
	// require.NoError(t, err)
	// defer func() {
	// 	err := tx.Rollback(t.Context())
	// 	if err != nil {
	// 		t.Logf("error during rollback: %v", err)
	// 	}
	// }()
	tx := pool

	instanceRepo := repository.InstanceRepository()
	// create instance
	instanceID := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err := instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create organization
	orgID := gofakeit.UUID()
	orgRepo := repository.OrganizationRepository()
	organization := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), tx, &organization)
	require.NoError(t, err)

	// create user, this user is created to have at least one user in the dbdatabase
	// that should not be returned in any of the results in these tests
	user := domain.Human{
		User: domain.User{
			ID:                gofakeit.UUID(),
			InstanceID:        instanceID,
			OrgID:             orgID,
			Username:          gofakeit.Username(),
			UsernameOrgUnique: true,
			State:             domain.UserStateActive,
		},
		FirstName:         gofakeit.FirstName(),
		LastName:          gofakeit.LastName(),
		NickName:          gofakeit.Username(),
		DisplayName:       gofakeit.Name(),
		PreferredLanguage: "en",
		Gender:            1,
		AvatarKey:         gofakeit.Animal(),
	}
	userRepo := repository.UserRepository()
	_, err = userRepo.CreateHuman(t.Context(), tx, &user)
	require.NoError(t, err)

	type test struct {
		name     string
		testFunc func(t *testing.T, client database.QueryExecutor) []*domain.Human
		user     []*domain.Human
		opts     []database.QueryOption
		err      error
	}

	tests := []*test{
		func() *test {
			newInstanceId := gofakeit.Name()
			return &test{
				name: "multiple users filter on instance, across 2 orgs",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Human {
					// create instance
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

					// create first org
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

					noOfUsers := 5
					users := make([]*domain.Human, noOfUsers*2)
					for i := range noOfUsers {

						user := domain.Human{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							FirstName:         gofakeit.FirstName(),
							LastName:          gofakeit.LastName(),
							NickName:          gofakeit.Username(),
							DisplayName:       gofakeit.Name(),
							PreferredLanguage: "en",
							Gender:            1,
							AvatarKey:         gofakeit.Animal(),
						}

						_, err := userRepo.CreateHuman(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					// create second org
					newOrgId = gofakeit.Name()
					org = domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstanceId,
						State:      domain.OrgStateActive,
					}
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					for i := range noOfUsers {
						i += noOfUsers

						user := domain.Human{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							FirstName:         gofakeit.FirstName(),
							LastName:          gofakeit.LastName(),
							NickName:          gofakeit.Username(),
							DisplayName:       gofakeit.Name(),
							PreferredLanguage: "en",
							Gender:            1,
							AvatarKey:         gofakeit.Animal(),
						}

						_, err := userRepo.CreateHuman(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					return users
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(newInstanceId),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return &test{
				name: "multiple users filter on non existing instance",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Human {
					// create instance
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

					// create first org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstanceId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					noOfUsers := 5
					users := make([]*domain.Human, noOfUsers)
					for i := range noOfUsers {

						user := domain.Human{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							FirstName:         gofakeit.FirstName(),
							LastName:          gofakeit.LastName(),
							NickName:          gofakeit.Username(),
							DisplayName:       gofakeit.Name(),
							PreferredLanguage: "en",
							Gender:            1,
							AvatarKey:         gofakeit.Animal(),
						}

						_, err := userRepo.CreateHuman(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					// return users
					return nil
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition("non-existing-instance"),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return &test{
				name: "multiple users filter on orgs",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Human {
					// create instance
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

					// create first org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstanceId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					noOfUsers := 5
					users := make([]*domain.Human, noOfUsers)
					for i := range noOfUsers {

						user := domain.Human{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							FirstName:         gofakeit.FirstName(),
							LastName:          gofakeit.LastName(),
							NickName:          gofakeit.Username(),
							DisplayName:       gofakeit.Name(),
							PreferredLanguage: "en",
							Gender:            1,
							AvatarKey:         gofakeit.Animal(),
						}

						_, err := userRepo.CreateHuman(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					return users
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(newInstanceId),
						userRepo.Human().OrgIDCondition(newOrgId),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return &test{
				name: "multiple users filter on orgs, with non existing org",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Human {
					// create instance
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

					// create first org
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstanceId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					noOfUsers := 5
					users := make([]*domain.Human, noOfUsers)
					for i := range noOfUsers {

						user := domain.Human{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							FirstName:         gofakeit.FirstName(),
							LastName:          gofakeit.LastName(),
							NickName:          gofakeit.Username(),
							DisplayName:       gofakeit.Name(),
							PreferredLanguage: "en",
							Gender:            1,
							AvatarKey:         gofakeit.Animal(),
						}

						_, err := userRepo.CreateHuman(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					// return users
					return nil
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(newInstanceId),
						userRepo.Human().OrgIDCondition("non-existing-org"),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			return &test{
				name: "multiple users filter on state",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Human {
					// create instance
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

					// create first org
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

					noOfUsers := 5
					users := make([]*domain.Human, noOfUsers)
					for i := range noOfUsers {

						user := domain.Human{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							FirstName:         gofakeit.FirstName(),
							LastName:          gofakeit.LastName(),
							NickName:          gofakeit.Username(),
							DisplayName:       gofakeit.Name(),
							PreferredLanguage: "en",
							Gender:            1,
							AvatarKey:         gofakeit.Animal(),
						}

						_, err := userRepo.CreateHuman(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					return users
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(newInstanceId),
						userRepo.Human().StateCondition(domain.UserStateInactive),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			return &test{
				name: "multiple users filter on state no users found",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Human {
					// create instance
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

					// create first org
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

					noOfUsers := 5
					users := make([]*domain.Human, noOfUsers)
					for i := range noOfUsers {

						user := domain.Human{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							FirstName:         gofakeit.FirstName(),
							LastName:          gofakeit.LastName(),
							NickName:          gofakeit.Username(),
							DisplayName:       gofakeit.Name(),
							PreferredLanguage: "en",
							Gender:            1,
							AvatarKey:         gofakeit.Animal(),
						}

						_, err := userRepo.CreateHuman(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					// return nil as we expect no users to be returned
					return nil
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Human().InstanceIDCondition(newInstanceId),
						// should be no users with active state
						userRepo.Human().StateCondition(domain.UserStateActive),
					)),
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// savepoint, err := tx.Begin(t.Context())
			// require.NoError(t, err)
			// defer func() {
			// 	err = savepoint.Rollback(t.Context())
			// 	if err != nil {
			// 		t.Logf("error during rollback: %v", err)
			// 	}
			// }()

			users := tt.testFunc(t, pool)

			// create user
			returnedUsers, err := userRepo.ListHuman(t.Context(), pool, tt.opts...)
			require.NoError(t, err)
			if err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}

			assert.Equal(t, len(users), len(returnedUsers))
			for i, user := range users {
				// user
				assert.Equal(t, returnedUsers[i].InstanceID, user.InstanceID)
				assert.Equal(t, returnedUsers[i].OrgID, user.OrgID)
				assert.Equal(t, returnedUsers[i].State, user.State)
				assert.Equal(t, returnedUsers[i].ID, user.ID)
				assert.Equal(t, returnedUsers[i].Username, user.Username)
				assert.Equal(t, returnedUsers[i].UsernameOrgUnique, user.UsernameOrgUnique)
				assert.Equal(t, returnedUsers[i].State, user.State)
				assert.Equal(t, returnedUsers[i].CreatedAt, user.CreatedAt)
				assert.Equal(t, returnedUsers[i].UpdatedAt, user.UpdatedAt)

				// // human
				assert.Equal(t, returnedUsers[i].FirstName, user.FirstName)
				assert.Equal(t, returnedUsers[i].LastName, user.LastName)
				assert.Equal(t, returnedUsers[i].NickName, user.NickName)
				assert.Equal(t, returnedUsers[i].DisplayName, user.DisplayName)
				assert.Equal(t, returnedUsers[i].PreferredLanguage, user.PreferredLanguage)
				assert.Equal(t, returnedUsers[i].Gender, user.Gender)
				assert.Equal(t, returnedUsers[i].AvatarKey, user.AvatarKey)
			}
		})
	}
}

func TestDeleteHumanUser(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Errorf("error rolling back transaction: %v", err)
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

	userRepo := repository.UserRepository()

	type test struct {
		name                   string
		testFunc               func(t *testing.T)
		idpIdentifierCondition domain.IDPIdentifierCondition
		noOfDeletedRows        int64
	}
	tests := []test{
		func() test {
			id := gofakeit.Name()
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

			// create first org
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
			var noOfUsers int64 = 1
			return test{
				name: "happy path delete user filter id",
				testFunc: func(t *testing.T) {
					for range noOfUsers {
						user := domain.Human{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							FirstName:         gofakeit.FirstName(),
							LastName:          gofakeit.LastName(),
							NickName:          gofakeit.Username(),
							DisplayName:       gofakeit.Name(),
							PreferredLanguage: "en",
							Gender:            1,
							AvatarKey:         gofakeit.Animal(),
						}

						err := userRepo.CreateHuman(t.Context(), tx, &user)
						require.NoError(t, err)
					}
				},
				idpIdentifierCondition: userRepo.IDCondition(id),
				noOfDeletedRows:        noOfUsers,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.testFunc != nil {
				tt.testFunc(t)
			}

			// delete user
			noOfDeletedRows, err := userRepo.Delete(t.Context(),
				tx,
				tt.idpIdentifierCondition,
				instanceId,
				&orgId,
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check idp was deleted
			organization, err := userRepo.Get(t.Context(), tx,
				tt.idpIdentifierCondition,
				instanceId,
				&orgId,
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, organization)
		})
	}
}

func TestCreateMachineUser(t *testing.T) {
	beforeCreate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()
	// create instance
	instanceID := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create organization
	orgID := gofakeit.UUID()
	orgRepo := repository.OrganizationRepository()
	organization := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), tx, &organization)
	require.NoError(t, err)

	userRepo := repository.UserRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T, client database.QueryExecutor) *domain.Organization
		user     *domain.Machine
		opts     []database.QueryOption
		err      error
	}

	tests := []test{
		func() test {
			id := gofakeit.UUID()
			return test{
				name: "happy path",
				user: &domain.Machine{
					User: domain.User{
						ID:                id,
						InstanceID:        instanceID,
						OrgID:             orgID,
						Username:          gofakeit.Username(),
						UsernameOrgUnique: true,
						State:             domain.UserStateActive,
					},
					Name:        gofakeit.Name(),
					Description: gofakeit.HipsterSentence(10),
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(instanceID),
						userRepo.Machine().OrgIDCondition(orgID),
						userRepo.Machine().IDCondition(id),
					)),
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = savepoint.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			// create user
			beforeCreate = time.Now()
			user, err := userRepo.CreateMachine(t.Context(), savepoint, tt.user)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			afterCreate := time.Now()

			createdUser, err := userRepo.GetMachine(t.Context(), savepoint, tt.opts...)
			require.NoError(t, err)
			fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> err = %+v\n", err)
			fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> createdUser) = %+v\n", createdUser)
			// assert.Equal(t, tt.organization.ID, organization.ID)
			// assert.Equal(t, tt.organization.Name, organization.Name)
			// assert.Equal(t, tt.organization.InstanceID, organization.InstanceID)
			// assert.Equal(t, tt.organization.State, organization.State)
			assert.WithinRange(t, user.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, user.CreatedAt, beforeCreate, afterCreate)
		})
	}
}

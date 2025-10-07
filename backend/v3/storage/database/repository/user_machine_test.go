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

func TestCreateMachineUser(t *testing.T) {
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
		user     *domain.Machine
		opts     []database.QueryOption
		err      error
	}

	tests := []test{
		// user
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
					Description: gu.Ptr("description"),
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
		{
			name: "no id",
			user: &domain.Machine{
				User: domain.User{
					// ID:         id,
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
			},
			err: new(database.CheckError),
		},
		{
			name: "no username",
			user: &domain.Machine{
				User: domain.User{
					ID:         gofakeit.UUID(),
					InstanceID: instanceID,
					OrgID:      orgID,
					// Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
			},
			err: new(database.CheckError),
		},
		func() test {
			id := gofakeit.UUID()
			user := &domain.Machine{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
			}
			return test{
				name: "same instance + org + username twice, with username_org_unique = true",
				user: user,
				testFunc: func(t *testing.T, client database.QueryExecutor) {
					_, err := userRepo.CreateMachine(t.Context(), pool, user)
					require.NoError(t, err)
				},
				err: new(database.UniqueError),
			}
		}(),
		func() test {
			id := gofakeit.UUID()
			user := &domain.Machine{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: false,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
			}
			_, err := userRepo.CreateMachine(t.Context(), pool, user)
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

			user := &domain.Machine{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
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

					_, err := userRepo.CreateMachine(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different org
					user.OrgID = orgID
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
		func() test {
			id := gofakeit.UUID()
			// create organization
			newOrgID := gofakeit.UUID()

			user := &domain.Machine{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: false,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
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

					_, err := userRepo.CreateMachine(t.Context(), pool, user)
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

			user := &domain.Machine{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
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

					_, err = userRepo.CreateMachine(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different org
					user.OrgID = orgID
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
		func() test {
			id := gofakeit.UUID()
			newOrgID := gofakeit.UUID()

			user := &domain.Machine{
				User: domain.User{
					ID:                id,
					InstanceID:        instanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: false,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
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

					_, err = userRepo.CreateMachine(t.Context(), pool, user)
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

			user := &domain.Machine{
				User: domain.User{
					ID:                id,
					InstanceID:        newInstanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
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

					_, err = userRepo.CreateMachine(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different instance
					user.InstanceID = instanceID
					user.OrgID = orgID
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
		func() test {
			id := gofakeit.UUID()
			// // create instance
			newInstanceID := gofakeit.Name()
			newOrgID := gofakeit.UUID()

			user := &domain.Machine{
				User: domain.User{
					ID:                id,
					InstanceID:        newInstanceID,
					OrgID:             newOrgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: false,
					State:             domain.UserStateActive,
				},
				Name: gofakeit.Name(),
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

					_, err = userRepo.CreateMachine(t.Context(), pool, user)
					require.NoError(t, err)
					// change to different instance
					user.InstanceID = instanceID
					user.OrgID = orgID
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(instanceID),
						userRepo.Machine().OrgIDCondition(orgID),
						userRepo.Machine().IDCondition(id),
					)),
				},
				// err: new(database.UniqueError),
			}
		}(),
		// machine
		{
			name: "no name",
			user: &domain.Machine{
				User: domain.User{
					ID:                gofakeit.UUID(),
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				// Name: gofakeit.Name(),
			},
			err: new(database.CheckError),
		},
		{
			name: "no first name",
			user: &domain.Machine{
				User: domain.User{
					ID:                gofakeit.UUID(),
					InstanceID:        instanceID,
					OrgID:             orgID,
					Username:          gofakeit.Username(),
					UsernameOrgUnique: true,
					State:             domain.UserStateActive,
				},
				// Name:         gofakeit.Name(),
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
			user, err := userRepo.CreateMachine(t.Context(), savepoint, tt.user)
			if err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			afterCreate := time.Now()

			createdUser, err := userRepo.GetMachine(t.Context(), savepoint, tt.opts...)
			require.NoError(t, err)

			// user
			assert.Equal(t, tt.user.InstanceID, createdUser.InstanceID)
			assert.Equal(t, tt.user.OrgID, createdUser.OrgID)
			assert.Equal(t, tt.user.State, createdUser.State)
			assert.Equal(t, tt.user.ID, createdUser.ID)
			assert.Equal(t, tt.user.Username, createdUser.Username)
			assert.Equal(t, tt.user.UsernameOrgUnique, createdUser.UsernameOrgUnique)
			assert.Equal(t, tt.user.State, createdUser.State)

			// machine
			assert.Equal(t, tt.user.Name, createdUser.Name)
			assert.Equal(t, tt.user.Description, createdUser.Description)

			assert.WithinRange(t, user.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, user.CreatedAt, beforeCreate, afterCreate)
		})
	}
}

func TestUpdateMachineUser(t *testing.T) {
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
		testFunc     func(t *testing.T, tx database.QueryExecutor) *domain.Machine
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
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Machine {
					user := &domain.Machine{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						Name: gofakeit.Name(),
					}

					user, err := userRepo.CreateMachine(t.Context(), tx, user)
					require.NoError(t, err)
					user.Username = "new_username"
					return user
				},
				condition: database.And(
					userRepo.Machine().InstanceIDCondition(instanceID),
					userRepo.Machine().OrgIDCondition(orgID),
					userRepo.Machine().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Machine().SetUsername("new_username")},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "update username invalid value",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Machine {
					user := &domain.Machine{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						Name: gofakeit.Name(),
					}

					user, err := userRepo.CreateMachine(t.Context(), tx, user)
					require.NoError(t, err)
					return user
				},
				condition: database.And(
					userRepo.Machine().InstanceIDCondition(instanceID),
					userRepo.Machine().OrgIDCondition(orgID),
					userRepo.Machine().IDCondition(id),
				),
				changes: []database.Change{userRepo.Machine().SetUsername("")},
				err:     new(database.CheckError),
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update username_org_unique",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Machine {
					user := &domain.Machine{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						Name: gofakeit.Name(),
					}

					user, err := userRepo.CreateMachine(t.Context(), tx, user)
					require.NoError(t, err)
					user.UsernameOrgUnique = true
					return user
				},
				condition: database.And(
					userRepo.Machine().InstanceIDCondition(instanceID),
					userRepo.Machine().OrgIDCondition(orgID),
					userRepo.Machine().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Machine().SetUsernameOrgUnique(true)},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update state",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Machine {
					user := &domain.Machine{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						Name: gofakeit.Name(),
					}

					user, err := userRepo.CreateMachine(t.Context(), tx, user)
					require.NoError(t, err)
					user.State = domain.UserStateInactive
					return user
				},
				condition: database.And(
					userRepo.Machine().InstanceIDCondition(instanceID),
					userRepo.Machine().OrgIDCondition(orgID),
					userRepo.Machine().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Machine().SetState(domain.UserStateInactive)},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "update invalid state",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Machine {
					user := &domain.Machine{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						Name: gofakeit.Name(),
					}

					user, err := userRepo.CreateMachine(t.Context(), tx, user)
					require.NoError(t, err)
					user.State = domain.UserStateInactive
					return user
				},
				condition: database.And(
					userRepo.Machine().InstanceIDCondition(instanceID),
					userRepo.Machine().OrgIDCondition(orgID),
					userRepo.Machine().IDCondition(id),
				),
				changes: []database.Change{userRepo.Machine().SetState(domain.UserStateUnspecified)},
				err:     new(database.CheckError),
			}
			return test
		}(),
		// machine
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update  name",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Machine {
					user := &domain.Machine{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						Name: "first_name",
					}

					user, err := userRepo.CreateMachine(t.Context(), tx, user)
					require.NoError(t, err)
					user.Name = "new_name"
					return user
				},
				condition: database.And(
					userRepo.Machine().InstanceIDCondition(instanceID),
					userRepo.Machine().OrgIDCondition(orgID),
					userRepo.Machine().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Machine().SetName("new_name")},
				rowsAffected: 1,
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "update name to invalid value",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Machine {
					user := &domain.Machine{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          "username",
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						Name: "name",
					}

					user, err := userRepo.CreateMachine(t.Context(), tx, user)
					require.NoError(t, err)
					return user
				},
				condition: database.And(
					userRepo.Machine().InstanceIDCondition(instanceID),
					userRepo.Machine().OrgIDCondition(orgID),
					userRepo.Machine().IDCondition(id),
				),
				changes: []database.Change{userRepo.Machine().SetName("")},
				err:     new(database.CheckError),
			}
			return test
		}(),
		func() *test {
			id := gofakeit.UUID()
			test := &test{
				name: "happy path update description",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Machine {
					user := &domain.Machine{
						User: domain.User{
							ID:                id,
							InstanceID:        instanceID,
							OrgID:             orgID,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: false,
							State:             domain.UserStateActive,
						},
						Name:        gofakeit.Username(),
						Description: gu.Ptr("description"),
					}

					user, err := userRepo.CreateMachine(t.Context(), tx, user)
					require.NoError(t, err)
					user.Description = gu.Ptr("new_description")
					return user
				},
				condition: database.And(
					userRepo.Machine().InstanceIDCondition(instanceID),
					userRepo.Machine().OrgIDCondition(orgID),
					userRepo.Machine().IDCondition(id),
				),
				changes:      []database.Change{userRepo.Machine().SetDescription("new_description")},
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
			rowsAffected, err := userRepo.UpdateMachine(
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
			user, err := userRepo.GetMachine(
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

			// machine
			assert.Equal(t, createdUser.Name, user.Name)
			assert.Equal(t, createdUser.Description, user.Description)

			assert.WithinRange(t, user.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

func TestGetMachineUser(t *testing.T) {
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

	// create user, this user is created to have at least one user in the database
	// that should not be returned in any of the results in these tests
	user := domain.Machine{
		User: domain.User{
			ID:                gofakeit.UUID(),
			InstanceID:        instanceID,
			OrgID:             orgID,
			Username:          gofakeit.Username(),
			UsernameOrgUnique: true,
			State:             domain.UserStateActive,
		},
		Name: gofakeit.Name(),
	}
	userRepo := repository.UserRepository()
	_, err = userRepo.CreateMachine(t.Context(), tx, &user)
	require.NoError(t, err)

	type test struct {
		name     string
		testFunc func(t *testing.T, client database.QueryExecutor) *domain.Machine
		user     []*domain.Machine
		opts     []database.QueryOption
		err      error
	}

	tests := []*test{
		func() *test {
			newInstanceId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			id := gofakeit.UUID()
			return &test{
				name: "Get happy path",
				testFunc: func(t *testing.T, client database.QueryExecutor) *domain.Machine {
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

					user := domain.Machine{
						User: domain.User{
							ID:                id,
							InstanceID:        newInstanceId,
							OrgID:             newOrgId,
							Username:          gofakeit.Username(),
							UsernameOrgUnique: true,
							State:             domain.UserStateActive,
						},
						Name:        gofakeit.Name(),
						Description: gu.Ptr("description"),
					}
					userRepo := repository.UserRepository()
					_, err = userRepo.CreateMachine(t.Context(), tx, &user)
					require.NoError(t, err)

					return &user
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
						userRepo.Machine().OrgIDCondition(newOrgId),
						userRepo.Machine().IDCondition(id),
					)),
				},
			}
		}(),
		func() *test {
			newOrgId := gofakeit.Name()
			id := gofakeit.UUID()
			return &test{
				name: "no instance id",
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						// userRepo.Machine().InstanceIDCondition(newInstanceId),
						userRepo.Machine().OrgIDCondition(newOrgId),
						userRepo.Machine().IDCondition(id),
					)),
				},
				err: new(database.MissingConditionError),
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			// newOrgId := gofakeit.Name()
			id := gofakeit.UUID()
			return &test{
				name: "no org id",
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
						// userRepo.Machine().OrgIDCondition(newOrgId),
						userRepo.Machine().IDCondition(id),
					)),
				},
				err: new(database.MissingConditionError),
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			// id := gofakeit.UUID()
			return &test{
				name: "no user id",
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
						userRepo.Machine().OrgIDCondition(newOrgId),
						// userRepo.Machine().IDCondition(id),
					)),
				},
				err: new(database.MissingConditionError),
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			// id := gofakeit.UUID()
			return &test{
				name: "get non existing user",
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
						userRepo.Machine().OrgIDCondition(newOrgId),
						// userRepo.Machine().IDCondition(id),
					)),
				},
				err: new(database.MissingConditionError),
			}
		}(),
		func() *test {
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
			return &test{
				name: "delete already deleted user",
				testFunc: func(t *testing.T, client database.QueryExecutor) *domain.Machine {
					for range noOfUsers {
						user := domain.Machine{
							User: domain.User{
								ID:                id,
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)

						// delete user
						noOfDeletedRows, err := userRepo.DeleteMachine(
							t.Context(),
							tx,
							database.And(
								userRepo.Machine().InstanceIDCondition(newInstanceId),
								userRepo.Machine().OrgIDCondition(newOrgId),
								userRepo.Machine().IDCondition(id),
							),
						)
						require.NoError(t, err)
						assert.Equal(t, int64(1), noOfDeletedRows)
					}
					return &user
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
						userRepo.Machine().OrgIDCondition(newOrgId),
						userRepo.Machine().IDCondition(id),
					)),
				},
				err: new(database.NoRowFoundError),
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

			var user *domain.Machine
			if tt.testFunc != nil {
				user = tt.testFunc(t, pool)
			}

			// create returnedUser
			returnedUser, err := userRepo.GetMachine(t.Context(), pool, tt.opts...)
			if err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)

			// user
			assert.Equal(t, user.InstanceID, returnedUser.InstanceID)
			assert.Equal(t, user.OrgID, returnedUser.OrgID)
			assert.Equal(t, user.State, returnedUser.State)
			assert.Equal(t, user.ID, returnedUser.ID)
			assert.Equal(t, user.Username, returnedUser.Username)
			assert.Equal(t, user.UsernameOrgUnique, returnedUser.UsernameOrgUnique)
			assert.Equal(t, user.State, returnedUser.State)
			assert.Equal(t, user.CreatedAt, returnedUser.CreatedAt)
			assert.Equal(t, user.UpdatedAt, returnedUser.UpdatedAt)

			// machine
			assert.Equal(t, user.Name, returnedUser.Name)
			assert.Equal(t, user.Description, returnedUser.Description)
		})
	}
}

func TestListMachineUser(t *testing.T) {
	// tx, err := pool.Begin.Context(), nil)
	// require.NoError(t, er
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
	user := domain.Machine{
		User: domain.User{
			ID:                gofakeit.UUID(),
			InstanceID:        instanceID,
			OrgID:             orgID,
			Username:          gofakeit.Username(),
			UsernameOrgUnique: true,
			State:             domain.UserStateActive,
		},
		Name: gofakeit.Name(),
	}
	userRepo := repository.UserRepository()
	_, err = userRepo.CreateMachine(t.Context(), tx, &user)
	require.NoError(t, err)

	type test struct {
		name     string
		testFunc func(t *testing.T, client database.QueryExecutor) []*domain.Machine
		user     []*domain.Machine
		opts     []database.QueryOption
		err      error
	}

	tests := []*test{
		func() *test {
			newInstanceId := gofakeit.Name()
			return &test{
				name: "multiple users filter on instance, across 2 orgs",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Machine {
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
					users := make([]*domain.Machine, noOfUsers*2)
					for i := range noOfUsers {

						user := domain.Machine{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
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

						user := domain.Machine{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					return users
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return &test{
				name: "multiple users filter on non existing instance",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Machine {
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
					users := make([]*domain.Machine, noOfUsers)
					for i := range noOfUsers {

						user := domain.Machine{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					// return users
					return nil
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition("non-existing-instance"),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return &test{
				name: "multiple users filter on orgs",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Machine {
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
					users := make([]*domain.Machine, noOfUsers)
					for i := range noOfUsers {

						user := domain.Machine{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					return users
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
						userRepo.Machine().OrgIDCondition(newOrgId),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			newOrgId := gofakeit.Name()
			return &test{
				name: "multiple users filter on orgs, with non existing org",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Machine {
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
					users := make([]*domain.Machine, noOfUsers)
					for i := range noOfUsers {

						user := domain.Machine{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateActive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					// return users
					return nil
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
						userRepo.Machine().OrgIDCondition("non-existing-org"),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			return &test{
				name: "multiple users filter on state",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Machine {
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
					users := make([]*domain.Machine, noOfUsers)
					for i := range noOfUsers {

						user := domain.Machine{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					return users
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
						userRepo.Machine().StateCondition(domain.UserStateInactive),
					)),
				},
			}
		}(),
		func() *test {
			newInstanceId := gofakeit.Name()
			return &test{
				name: "multiple users filter on state no users found",
				testFunc: func(t *testing.T, client database.QueryExecutor) []*domain.Machine {
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
					users := make([]*domain.Machine, noOfUsers)
					for i := range noOfUsers {

						user := domain.Machine{
							User: domain.User{
								ID:                gofakeit.UUID(),
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)

						users[i] = &user
					}

					// return nil as we expect no users to be returned
					return nil
				},
				opts: []database.QueryOption{
					database.WithCondition(database.And(
						userRepo.Machine().InstanceIDCondition(newInstanceId),
						// should be no users with active state
						userRepo.Machine().StateCondition(domain.UserStateActive),
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
			returnedUsers, err := userRepo.ListMachine(t.Context(), pool, tt.opts...)
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

				// machine
				assert.Equal(t, returnedUsers[i].Name, user.Name)
				assert.Equal(t, returnedUsers[i].Description, user.Description)
			}
		})
	}
}

func TestDeleteMachineUser(t *testing.T) {
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
		name            string
		testFunc        func(t *testing.T)
		conditions      database.Condition
		noOfDeletedRows int64
		err             error
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
				name: "happy path delete user",
				testFunc: func(t *testing.T) {
					for range noOfUsers {
						user := domain.Machine{
							User: domain.User{
								ID:                id,
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)
					}
				},
				conditions: database.And(
					userRepo.Machine().InstanceIDCondition(newInstanceId),
					userRepo.Machine().OrgIDCondition(newOrgId),
					userRepo.Machine().IDCondition(id),
				),
				noOfDeletedRows: noOfUsers,
			}
		}(),
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
				name: "delete user, no instance specified",
				testFunc: func(t *testing.T) {
					for range noOfUsers {
						user := domain.Machine{
							User: domain.User{
								ID:                id,
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)
					}
				},
				conditions: database.And(
					// userRepo.Machine().InstanceIDCondition(newInstanceId),
					userRepo.Machine().OrgIDCondition(newOrgId),
					userRepo.Machine().IDCondition(id),
				),
				noOfDeletedRows: 0,
				err:             new(database.MissingConditionError),
			}
		}(),
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
				name: "delete user, no organization specified",
				testFunc: func(t *testing.T) {
					for range noOfUsers {
						user := domain.Machine{
							User: domain.User{
								ID:                id,
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)
					}
				},
				conditions: database.And(
					userRepo.Machine().InstanceIDCondition(newInstanceId),
					// userRepo.Machine().OrgIDCondition(newOrgId),
					userRepo.Machine().IDCondition(id),
				),
				noOfDeletedRows: 0,
				err:             new(database.MissingConditionError),
			}
		}(),
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
				name: "delete user, no id specified",
				testFunc: func(t *testing.T) {
					for range noOfUsers {
						user := domain.Machine{
							User: domain.User{
								ID:                id,
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)
					}
				},
				conditions: database.And(
					userRepo.Machine().InstanceIDCondition(newInstanceId),
					userRepo.Machine().OrgIDCondition(newOrgId),
					// userRepo.Machine().IDCondition(id),
				),
				noOfDeletedRows: 0,
				err:             new(database.MissingConditionError),
			}
		}(),
		func() test {
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
			return test{
				name: "delete user, no id specified",
				conditions: database.And(
					userRepo.Machine().InstanceIDCondition(newInstanceId),
					userRepo.Machine().OrgIDCondition(newOrgId),
					userRepo.Machine().IDCondition("non-existing-id"),
				),
				noOfDeletedRows: 0,
			}
		}(),
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
				name: "delete already deleted user",
				testFunc: func(t *testing.T) {
					for range noOfUsers {
						user := domain.Machine{
							User: domain.User{
								ID:                id,
								InstanceID:        newInstanceId,
								OrgID:             newOrgId,
								Username:          gofakeit.Username(),
								UsernameOrgUnique: true,
								State:             domain.UserStateInactive,
							},
							Name: gofakeit.Name(),
						}

						_, err := userRepo.CreateMachine(t.Context(), tx, &user)
						require.NoError(t, err)

						// delete user
						noOfDeletedRows, err := userRepo.DeleteMachine(
							t.Context(),
							tx,
							database.And(
								userRepo.Machine().InstanceIDCondition(newInstanceId),
								userRepo.Machine().OrgIDCondition(newOrgId),
								userRepo.Machine().IDCondition(id),
							),
						)
						require.NoError(t, err)
						assert.Equal(t, int64(1), noOfDeletedRows)
					}
				},
				conditions: database.And(
					userRepo.Machine().InstanceIDCondition(newInstanceId),
					userRepo.Machine().OrgIDCondition(newOrgId),
					userRepo.Machine().IDCondition(id),
				),
				noOfDeletedRows: 0,
				// err:             new(database.MissingConditionError),
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.testFunc != nil {
				tt.testFunc(t)
			}

			// delete user
			noOfDeletedRows, err := userRepo.DeleteMachine(
				t.Context(),
				tx,
				tt.conditions,
			)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.noOfDeletedRows, noOfDeletedRows)

			// check user was deleted
			user, err := userRepo.GetMachine(
				t.Context(),
				tx,
				database.WithCondition(
					tt.conditions),
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, user)
		})
	}
}

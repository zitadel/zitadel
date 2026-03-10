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

func TestCreateIDProvider(t *testing.T) {
	beforeCreate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Errorf("error rolling back transaction: %v", err)
		}
	}()

	// create instance
	instanceId := createInstance(t, tx)

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
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider
		idp      domain.IdentityProvider
		err      error
	}

	tests := []test{
		{
			name: "happy path",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrgID:             &orgId,
				ID:                gofakeit.Name(),
				State:             domain.IDPStateActive,
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				Payload:           []byte("{}"),
			},
		},
		{
			name: "create idp without name",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrgID:             &orgId,
				ID:                gofakeit.Name(),
				State:             domain.IDPStateActive,
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				Payload:           []byte("{}"),
			},
			err: new(database.CheckError),
		},
		{
			name: "adding idp with same id twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
				idpRepo := repository.IDProviderRepository()

				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				// change the name to make sure same only the id clashes
				org.Name = gofakeit.Name()
				return &idp
			},
			err: new(database.UniqueError),
		},
		{
			name: "adding idp with same name twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
				idpRepo := repository.IDProviderRepository()

				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				// change the id to make sure same name causes an error
				idp.ID = gofakeit.Name()
				return &idp
			},
			err: new(database.UniqueError),
		},
		func() test {
			id := gofakeit.Name()
			name := gofakeit.Name()
			return test{
				name: "adding idp with same name, id, different instance",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
					// create instance
					newInstId := gofakeit.Name()
					instance := domain.Instance{
						ID:              newInstId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "managementConsoleCLient",
						ConsoleAppID:    "managementConsoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					instanceRepo := repository.InstanceRepository()
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create org
					newOrgId := gofakeit.Name()
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive,
					}
					organizationRepo := repository.OrganizationRepository()
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					idpRepo := repository.IDProviderRepository()
					idp := domain.IdentityProvider{
						InstanceID:        newInstId,
						OrgID:             &newOrgId,
						ID:                id,
						State:             domain.IDPStateActive,
						Name:              name,
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						Payload:           []byte("{}"),
					}

					err = idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					// change the instanceID to a different instance
					idp.InstanceID = instanceId
					// change the OrgId to a different organization
					idp.OrgID = &orgId
					return &idp
				},
				idp: domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                id,
					State:             domain.IDPStateActive,
					Name:              name,
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					Payload:           []byte("{}"),
				},
			}
		}(),
		{
			name: "adding idp with no id",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrgID:             &orgId,
				State:             domain.IDPStateActive,
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				Payload:           []byte("{}"),
			},
			err: new(database.CheckError),
		},
		{
			name: "adding idp with no name",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrgID:             &orgId,
				ID:                gofakeit.Name(),
				State:             domain.IDPStateActive,
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				Payload:           []byte("{}"),
			},
			err: new(database.CheckError),
		},
		{
			name: "adding idp with no instance id",
			idp: domain.IdentityProvider{
				OrgID:             &orgId,
				State:             domain.IDPStateActive,
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				Payload:           []byte("{}"),
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding idp with non existent instance id",
			idp: domain.IdentityProvider{
				InstanceID:        gofakeit.Name(),
				OrgID:             &orgId,
				State:             domain.IDPStateActive,
				ID:                gofakeit.Name(),
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				Payload:           []byte("{}"),
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrgID:             gu.Ptr(gofakeit.Name()),
				State:             domain.IDPStateActive,
				ID:                gofakeit.Name(),
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				Payload:           []byte("{}"),
			},
			err: new(database.ForeignKeyError),
		},
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

			var idp *domain.IdentityProvider
			if tt.testFunc != nil {
				idp = tt.testFunc(t, tx)
			} else {
				idp = &tt.idp
			}
			idpRepo := repository.IDProviderRepository()

			// create idp

			err = idpRepo.Create(t.Context(), tx, idp)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check organization values
			idp, err = idpRepo.Get(t.Context(), tx, database.WithCondition(database.And(idpRepo.PrimaryKeyCondition(idp.InstanceID, idp.ID), idpRepo.OrgIDCondition(idp.OrgID))))
			require.NoError(t, err)

			assert.Equal(t, tt.idp.InstanceID, idp.InstanceID)
			assert.Equal(t, tt.idp.OrgID, idp.OrgID)
			assert.Equal(t, tt.idp.State, idp.State)
			assert.Equal(t, tt.idp.ID, idp.ID)
			assert.Equal(t, tt.idp.Name, idp.Name)
			assert.Equal(t, tt.idp.Type, idp.Type)
			assert.Equal(t, tt.idp.AllowCreation, idp.AllowCreation)
			assert.Equal(t, tt.idp.AllowAutoCreation, idp.AllowAutoCreation)
			assert.Equal(t, tt.idp.AllowAutoUpdate, idp.AllowAutoUpdate)
			assert.Equal(t, tt.idp.AllowLinking, idp.AllowLinking)
			assert.Equal(t, tt.idp.Payload, idp.Payload)
			assert.WithinRange(t, idp.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

func TestUpdateIDProvider(t *testing.T) {
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
	instanceId := createInstance(t, tx)

	// create org
	orgId := createOrganization(t, tx, instanceId)

	idpRepo := repository.IDProviderRepository()

	tests := []struct {
		name         string
		testFunc     func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider
		update       []database.Change
		rowsAffected int64
	}{
		{
			name: "happy path update name",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				idp.Name = "new_name"
				idp.State = domain.IDPStateInactive
				idp.AllowCreation = false
				idp.AllowAutoCreation = false
				idp.AllowLinking = false
				idp.Payload = []byte(`{"json": {}}`)
				return &idp
			},
			update: []database.Change{
				idpRepo.SetName("new_name"),
				idpRepo.SetState(domain.IDPStateInactive),
				idpRepo.SetAllowCreation(false),
				idpRepo.SetAllowAutoCreation(false),
				idpRepo.SetAllowLinking(false),
				idpRepo.SetPayload(`{"json": {}}`),
			},
			rowsAffected: 1,
		},
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
			idpRepo := repository.IDProviderRepository()

			createdIDP := tt.testFunc(t, tx)

			// update idp
			rowsAffected, err := idpRepo.Update(t.Context(), tx,
				database.And(
					idpRepo.PrimaryKeyCondition(
						createdIDP.InstanceID,
						createdIDP.ID,
					),
					idpRepo.OrgIDCondition(createdIDP.OrgID),
				),
				tt.update...,
			)
			afterUpdate := time.Now()
			require.NoError(t, err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			if rowsAffected == 0 {
				return
			}

			// check idp values
			idp, err := idpRepo.Get(t.Context(), tx,
				database.WithCondition(database.And(
					idpRepo.PrimaryKeyCondition(createdIDP.InstanceID, createdIDP.ID),
					idpRepo.OrgIDCondition(createdIDP.OrgID),
				)),
			)
			require.NoError(t, err)

			assert.Equal(t, createdIDP.InstanceID, idp.InstanceID)
			assert.Equal(t, createdIDP.OrgID, idp.OrgID)
			assert.Equal(t, createdIDP.State, idp.State)
			assert.Equal(t, createdIDP.ID, idp.ID)
			assert.Equal(t, createdIDP.Name, idp.Name)
			assert.Equal(t, createdIDP.Type, idp.Type)
			assert.Equal(t, createdIDP.AllowCreation, idp.AllowCreation)
			assert.Equal(t, createdIDP.AllowAutoCreation, idp.AllowAutoCreation)
			assert.Equal(t, createdIDP.AllowAutoUpdate, idp.AllowAutoUpdate)
			assert.Equal(t, createdIDP.AllowLinking, idp.AllowLinking)
			assert.Equal(t, createdIDP.Payload, idp.Payload)
			assert.WithinRange(t, idp.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

func TestGetIDProvider(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Errorf("error rolling back transaction: %v", err)
		}
	}()

	instanceId := createInstance(t, tx)
	orgId := createOrganization(t, tx, instanceId)

	// this org is created as an additional org which should NOT
	// be returned in the results of the tests
	createOrganization(t, tx, instanceId)

	idpRepo := repository.IDProviderRepository()

	idp := domain.IdentityProvider{
		InstanceID:        instanceId,
		OrgID:             &orgId,
		ID:                gofakeit.UUID(),
		State:             domain.IDPStateActive,
		Name:              gofakeit.Name(),
		Type:              gu.Ptr(domain.IDPTypeOIDC),
		AllowCreation:     true,
		AllowAutoCreation: true,
		AllowAutoUpdate:   true,
		AllowLinking:      true,
		Payload:           []byte("{}"),
	}
	err = idpRepo.Create(t.Context(), tx, &idp)
	require.NoError(t, err)

	type test struct {
		name      string
		condition database.Condition
		err       error
	}

	tests := []test{
		{
			name:      "happy path get using id",
			condition: idpRepo.PrimaryKeyCondition(idp.InstanceID, idp.ID),
		},
		{
			name: "happy path get using name",
			condition: database.And(
				idpRepo.InstanceIDCondition(idp.InstanceID),
				idpRepo.NameCondition(database.TextOperationEqual, idp.Name),
			),
		},
		{
			name:      "get using non existent id",
			condition: idpRepo.PrimaryKeyCondition(idp.InstanceID, "non-existent-id"),
			err:       new(database.NoRowFoundError),
		},
		{
			name: "get using non existent name",
			condition: database.And(
				idpRepo.InstanceIDCondition(idp.InstanceID),
				idpRepo.NameCondition(database.TextOperationEqual, "non-existent-name"),
			),
			err: new(database.NoRowFoundError),
		},
		{
			name: "non existent orgID",
			condition: database.And(
				idpRepo.PrimaryKeyCondition(idp.InstanceID, idp.ID),
				idpRepo.OrgIDCondition(gu.Ptr("non-existing-org-id")),
			),
			err: new(database.NoRowFoundError),
		},
		{
			name:      "non existent instanceID",
			condition: idpRepo.PrimaryKeyCondition("non-existent-instance-id", idp.ID),
			err:       new(database.NoRowFoundError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// get idp
			returnedIDP, err := idpRepo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}

			assert.Equal(t, returnedIDP.InstanceID, idp.InstanceID)
			assert.Equal(t, returnedIDP.OrgID, idp.OrgID)
			assert.Equal(t, returnedIDP.State, idp.State)
			assert.Equal(t, returnedIDP.ID, idp.ID)
			assert.Equal(t, returnedIDP.Name, idp.Name)
			assert.Equal(t, returnedIDP.Type, idp.Type)
			assert.Equal(t, returnedIDP.AllowCreation, idp.AllowCreation)
			assert.Equal(t, returnedIDP.AllowAutoCreation, idp.AllowAutoCreation)
			assert.Equal(t, returnedIDP.AllowAutoUpdate, idp.AllowAutoUpdate)
			assert.Equal(t, returnedIDP.AllowLinking, idp.AllowLinking)
			assert.Equal(t, returnedIDP.Payload, idp.Payload)
		})
	}
}

// gocognit linting fails due to number of test cases
// and the fact that each test case has a testFunc()
//

func TestListIDProvider(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Errorf("error rolling back transaction: %v", err)
		}
	}()

	// create instance
	instanceId := createInstance(t, tx)

	// create org
	orgId := createOrganization(t, tx, instanceId)

	idpRepo := repository.IDProviderRepository()

	type test struct {
		name             string
		testFunc         func(t *testing.T) []*domain.IdentityProvider
		conditionClauses []database.Condition
		noIDPsReturned   bool
	}
	tests := []test{
		{
			name: "multiple idps filter on instance",
			testFunc: func(t *testing.T) []*domain.IdentityProvider {
				// create instance
				newInstanceId := createInstance(t, tx)

				// create org
				newOrgId := createOrganization(t, tx, newInstanceId)

				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID:        newInstanceId,
					OrgID:             &newOrgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					Payload:           []byte("{}"),
				}
				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive,
						Name:              gofakeit.Name(),
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.InstanceIDCondition(instanceId)},
		},
		{
			name: "multiple idps filter on org",
			testFunc: func(t *testing.T) []*domain.IdentityProvider {
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

				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &newOrgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					Payload:           []byte("{}"),
				}
				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive,
						Name:              gofakeit.Name(),
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.OrgIDCondition(&orgId)},
		},
		{
			name: "happy path single idp no filter",
			testFunc: func(t *testing.T) []*domain.IdentityProvider {
				noOfIDPs := 1
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive,
						Name:              gofakeit.Name(),
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
		},
		{
			name: "happy path multiple idps no filter",
			testFunc: func(t *testing.T) []*domain.IdentityProvider {
				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive,
						Name:              gofakeit.Name(),
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
		},
		func() test {
			id := gofakeit.Name()
			return test{
				name: "idp filter on id",
				testFunc: func(t *testing.T) []*domain.IdentityProvider {
					// create idp
					// this idp is created as an additional idp which should NOT
					// be returned in the results of this test case
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive,
						Name:              gofakeit.Name(),
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						Payload:           []byte("{}"),
					}
					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					noOfIDPs := 1
					idps := make([]*domain.IdentityProvider, noOfIDPs)
					for i := range noOfIDPs {

						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                id,
							State:             domain.IDPStateActive,
							Name:              gofakeit.Name(),
							Type:              gu.Ptr(domain.IDPTypeOIDC),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							Payload:           []byte("{}"),
						}

						err := idpRepo.Create(t.Context(), tx, &idp)
						require.NoError(t, err)

						idps[i] = &idp
					}

					return idps
				},
				conditionClauses: []database.Condition{idpRepo.IDCondition(id)},
			}
		}(),
		{
			name: "multiple idps filter on state",
			testFunc: func(t *testing.T) []*domain.IdentityProvider {
				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					// state inactive
					State:             domain.IDPStateInactive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					Payload:           []byte("{}"),
				}
				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						// state active
						State:             domain.IDPStateActive,
						Name:              gofakeit.Name(),
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.StateCondition(domain.IDPStateActive)},
		},
		func() test {
			name := gofakeit.Name()
			return test{
				name: "multiple idps filter on name",
				testFunc: func(t *testing.T) []*domain.IdentityProvider {
					// create idp
					// this idp is created as an additional idp which should NOT
					// be returned in the results of this test case
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive,
						Name:              gofakeit.Name(),
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						Payload:           []byte("{}"),
					}
					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					noOfIDPs := 1
					idps := make([]*domain.IdentityProvider, noOfIDPs)
					for i := range noOfIDPs {

						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                gofakeit.Name(),
							State:             domain.IDPStateActive,
							Name:              name,
							Type:              gu.Ptr(domain.IDPTypeOIDC),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							Payload:           []byte("{}"),
						}

						err := idpRepo.Create(t.Context(), tx, &idp)
						require.NoError(t, err)

						idps[i] = &idp
					}

					return idps
				},
				conditionClauses: []database.Condition{idpRepo.NameCondition(database.TextOperationEqual, name)},
			}
		}(),
		{
			name: "multiple idps filter on type",
			testFunc: func(t *testing.T) []*domain.IdentityProvider {
				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateInactive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					Payload:           []byte("{}"),
				}
				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						State:      domain.IDPStateActive,
						Name:       gofakeit.Name(),
						// type LDAP
						Type:              gu.Ptr(domain.IDPTypeLDAP),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.TypeCondition(domain.IDPTypeLDAP)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				_, err := tx.Exec(t.Context(), "DELETE FROM zitadel.identity_providers")
				require.NoError(t, err)
			}()

			idps := tt.testFunc(t)

			// check idp values
			returnedIDPs, err := idpRepo.List(t.Context(), tx,
				database.WithCondition(database.And(append(tt.conditionClauses, idpRepo.InstanceIDCondition(instanceId))...)),
			)
			require.NoError(t, err)
			if tt.noIDPsReturned {
				assert.Nil(t, returnedIDPs)
				return
			}

			assert.Equal(t, len(idps), len(returnedIDPs))
			assert.ElementsMatch(t, idps, returnedIDPs)
		})
	}
}

func TestDeleteIDProvider(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Errorf("error rolling back transaction: %v", err)
		}
	}()

	// create instance
	instanceId := createInstance(t, tx)

	// create org
	orgId := createOrganization(t, tx, instanceId)

	idpRepo := repository.IDProviderRepository()

	type test struct {
		name                   string
		testFunc               func(t *testing.T)
		idpIdentifierCondition database.Condition
		noOfDeletedRows        int64
	}
	tests := []test{
		func() test {
			id := gofakeit.Name()
			var noOfIDPs int64 = 1
			return test{
				name: "happy path delete idp filter id",
				testFunc: func(t *testing.T) {
					for range noOfIDPs {
						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                id,
							State:             domain.IDPStateActive,
							Name:              gofakeit.Name(),
							Type:              gu.Ptr(domain.IDPTypeOIDC),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							Payload:           []byte("{}"),
						}

						err := idpRepo.Create(t.Context(), tx, &idp)
						require.NoError(t, err)
					}
				},
				idpIdentifierCondition: idpRepo.IDCondition(id),
				noOfDeletedRows:        noOfIDPs,
			}
		}(),
		{
			name:                   "delete non existent idp",
			idpIdentifierCondition: idpRepo.IDCondition(gofakeit.UUID()),
		},
		func() test {
			id := gofakeit.UUID()
			return test{
				name: "deleted already deleted idp",
				testFunc: func(t *testing.T) {
					noOfIDPs := 1
					for range noOfIDPs {
						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                id,
							State:             domain.IDPStateActive,
							Name:              gofakeit.Name(),
							Type:              gu.Ptr(domain.IDPTypeOIDC),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							Payload:           []byte("{}"),
						}

						err := idpRepo.Create(t.Context(), tx, &idp)
						require.NoError(t, err)
					}

					// delete organization
					affectedRows, err := idpRepo.Delete(t.Context(), tx,
						database.And(
							idpRepo.InstanceIDCondition(instanceId),
							idpRepo.OrgIDCondition(&orgId),
							idpRepo.IDCondition(id),
						),
					)
					assert.Equal(t, int64(1), affectedRows)
					require.NoError(t, err)
				},
				idpIdentifierCondition: idpRepo.IDCondition(id),
				// this test should return 0 affected rows as the idp was already deleted
				noOfDeletedRows: 0,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.testFunc != nil {
				tt.testFunc(t)
			}

			// delete idp
			noOfDeletedRows, err := idpRepo.Delete(t.Context(), tx,
				database.And(
					tt.idpIdentifierCondition,
					idpRepo.InstanceIDCondition(instanceId),
					idpRepo.OrgIDCondition(&orgId),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check idp was deleted
			organization, err := idpRepo.Get(t.Context(), tx,
				database.WithCondition(database.And(
					tt.idpIdentifierCondition,
					idpRepo.InstanceIDCondition(instanceId),
					idpRepo.OrgIDCondition(&orgId),
				)),
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, organization)
		})
	}
}

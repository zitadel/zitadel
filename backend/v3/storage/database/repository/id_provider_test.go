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

var stylingType int16 = 1

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
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider
		idp      domain.IdentityProvider
		err      error
	}

	// TESTS
	tests := []test{
		{
			name: "happy path",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrganizationID:    &orgId,
				ID:                gofakeit.Name(),
				State:             domain.IDPStateActive,
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &stylingType,
				Payload:           []byte("{}"),
			},
		},
		{
			name: "create idp without name",
			idp: domain.IdentityProvider{
				InstanceID:     instanceId,
				OrganizationID: &orgId,
				ID:             gofakeit.Name(),
				State:          domain.IDPStateActive,
				// Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &stylingType,
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
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
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
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
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
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
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
						OrganizationID:    &newOrgId,
						ID:                id,
						State:             domain.IDPStateActive,
						Name:              name,
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &stylingType,
						Payload:           []byte("{}"),
					}

					err = idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)

					// change the instanceID to a different instance
					idp.InstanceID = instanceId
					// change the OrgId to a different organization
					idp.OrganizationID = &orgId
					return &idp
				},
				idp: domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                id,
					State:             domain.IDPStateActive,
					Name:              name,
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				},
			}
		}(),
		{
			name: "adding idp with no id",
			idp: domain.IdentityProvider{
				InstanceID:     instanceId,
				OrganizationID: &orgId,
				// ID:                gofakeit.Name(),
				State:             domain.IDPStateActive,
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &stylingType,
				Payload:           []byte("{}"),
			},
			err: new(database.CheckError),
		},
		{
			name: "adding idp with no name",
			idp: domain.IdentityProvider{
				InstanceID:     instanceId,
				OrganizationID: &orgId,
				ID:             gofakeit.Name(),
				State:          domain.IDPStateActive,
				// Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &stylingType,
				Payload:           []byte("{}"),
			},
			err: new(database.CheckError),
		},
		{
			name: "adding idp with no instance id",
			idp: domain.IdentityProvider{
				// InstanceID:        instanceId,
				OrganizationID:    &orgId,
				State:             domain.IDPStateActive,
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &stylingType,
				Payload:           []byte("{}"),
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding idp with non existent instance id",
			idp: domain.IdentityProvider{
				InstanceID:        gofakeit.Name(),
				OrganizationID:    &orgId,
				State:             domain.IDPStateActive,
				ID:                gofakeit.Name(),
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &stylingType,
				Payload:           []byte("{}"),
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding idp with non existent org id",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrganizationID:    gu.Ptr(gofakeit.Name()),
				State:             domain.IDPStateActive,
				ID:                gofakeit.Name(),
				Name:              gofakeit.Name(),
				Type:              gu.Ptr(domain.IDPTypeOIDC),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &stylingType,
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
			idp, err = idpRepo.Get(t.Context(), tx,
				database.WithCondition(idpRepo.PrimaryKeyCondition(idp.InstanceID, idp.ID)),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.idp.InstanceID, idp.InstanceID)
			assert.Equal(t, tt.idp.OrganizationID, idp.OrganizationID)
			assert.Equal(t, tt.idp.State, idp.State)
			assert.Equal(t, tt.idp.ID, idp.ID)
			assert.Equal(t, tt.idp.Name, idp.Name)
			assert.Equal(t, tt.idp.Type, idp.Type)
			assert.Equal(t, tt.idp.AllowCreation, idp.AllowCreation)
			assert.Equal(t, tt.idp.AllowAutoCreation, idp.AllowAutoCreation)
			assert.Equal(t, tt.idp.AllowAutoUpdate, idp.AllowAutoUpdate)
			assert.Equal(t, tt.idp.AllowLinking, idp.AllowLinking)
			assert.Equal(t, tt.idp.AutoLinkingField, idp.AutoLinkingField)
			assert.Equal(t, tt.idp.StylingType, idp.StylingType)
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
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				idp.Name = "new_name"
				return &idp
			},
			update:       []database.Change{idpRepo.SetName("new_name")},
			rowsAffected: 1,
		},
		{
			name: "happy path update state",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				idp.State = domain.IDPStateInactive
				return &idp
			},
			update:       []database.Change{idpRepo.SetState(domain.IDPStateInactive)},
			rowsAffected: 1,
		},
		{
			name: "happy path update AllowCreation",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				idp.AllowCreation = false
				return &idp
			},
			update:       []database.Change{idpRepo.SetAllowCreation(false)},
			rowsAffected: 1,
		},
		{
			name: "happy path update AllowAutoCreation",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				idp.AllowAutoCreation = false
				return &idp
			},
			update:       []database.Change{idpRepo.SetAllowAutoCreation(false)},
			rowsAffected: 1,
		},
		{
			name: "happy path update AllowLinking",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				idp.AllowLinking = false
				return &idp
			},
			update:       []database.Change{idpRepo.SetAllowLinking(false)},
			rowsAffected: 1,
		},
		{
			name: "happy path update StylingType",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				newStyleType := int16(2)
				idp.StylingType = &newStyleType
				return &idp
			},
			update:       []database.Change{idpRepo.SetStylingType(2)},
			rowsAffected: 1,
		},
		{
			name: "happy path update Payload",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				idp.Payload = []byte(`{"json": {}}`)
				return &idp
			},
			// update:       []database.Change{idpRepo.SetPayload("{{}}")},
			update:       []database.Change{idpRepo.SetPayload(`{"json": {}}`)},
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
			// organizationRepo := repository.OrganizationRepository()
			idpRepo := repository.IDProviderRepository()

			createdIDP := tt.testFunc(t, tx)

			// update idp
			rowsAffected, err := idpRepo.Update(t.Context(), tx,
				idpRepo.PrimaryKeyCondition(createdIDP.InstanceID, createdIDP.ID),
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
				database.WithCondition(idpRepo.PrimaryKeyCondition(createdIDP.InstanceID, createdIDP.ID)),
			)
			require.NoError(t, err)

			assert.Equal(t, createdIDP.InstanceID, idp.InstanceID)
			assert.Equal(t, createdIDP.OrganizationID, idp.OrganizationID)
			assert.Equal(t, createdIDP.State, idp.State)
			assert.Equal(t, createdIDP.ID, idp.ID)
			assert.Equal(t, createdIDP.Name, idp.Name)
			assert.Equal(t, createdIDP.Type, idp.Type)
			assert.Equal(t, createdIDP.AllowCreation, idp.AllowCreation)
			assert.Equal(t, createdIDP.AllowAutoCreation, idp.AllowAutoCreation)
			assert.Equal(t, createdIDP.AllowAutoUpdate, idp.AllowAutoUpdate)
			assert.Equal(t, createdIDP.AllowLinking, idp.AllowLinking)
			assert.Equal(t, createdIDP.AutoLinkingField, idp.AutoLinkingField)
			assert.Equal(t, createdIDP.StylingType, idp.StylingType)
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

	// create organization
	// this org is created as an additional org which should NOT
	// be returned in the results of the tests
	preexistingOrg := domain.Organization{
		ID:         gofakeit.Name(),
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	err = organizationRepo.Create(t.Context(), tx, &preexistingOrg)
	require.NoError(t, err)

	idpRepo := repository.IDProviderRepository()
	type test struct {
		name                   string
		testFunc               func(t *testing.T) *domain.IdentityProvider
		idpIdentifierCondition database.Condition
		err                    error
	}

	tests := []test{
		func() test {
			id := gofakeit.Name()
			return test{
				name: "happy path get using id",
				testFunc: func(t *testing.T) *domain.IdentityProvider {
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrganizationID:    &orgId,
						ID:                id,
						State:             domain.IDPStateActive,
						Name:              gofakeit.Name(),
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &stylingType,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)
					return &idp
				},
				idpIdentifierCondition: idpRepo.IDCondition(id),
			}
		}(),
		func() test {
			name := gofakeit.Name()
			return test{
				name: "happy path get using name",
				testFunc: func(t *testing.T) *domain.IdentityProvider {
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrganizationID:    &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive,
						Name:              name,
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &stylingType,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)
					return &idp
				},
				idpIdentifierCondition: idpRepo.NameCondition(name),
			}
		}(),
		{
			name: "happy path using styling type",
			testFunc: func(t *testing.T) *domain.IdentityProvider {
				stylingType := 2
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       gu.Ptr(int16(stylingType)),
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				return &idp
			},
			idpIdentifierCondition: idpRepo.StylingTypeCondition(2),
		},
		{
			name: "get using non existent id",
			testFunc: func(t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				return &idp
			},
			idpIdentifierCondition: idpRepo.IDCondition("non-existent-id"),
			err:                    new(database.NoRowFoundError),
		},
		{
			name: "get using non existent name",
			testFunc: func(t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrganizationID:    &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive,
					Name:              gofakeit.Name(),
					Type:              gu.Ptr(domain.IDPTypeOIDC),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &stylingType,
					Payload:           []byte("{}"),
				}

				err := idpRepo.Create(t.Context(), tx, &idp)
				require.NoError(t, err)
				return &idp
			},
			idpIdentifierCondition: idpRepo.NameCondition("non-existent-name"),
			err:                    new(database.NoRowFoundError),
		},
		func() test {
			id := gofakeit.Name()
			return test{
				name: "non existent orgID",
				testFunc: func(t *testing.T) *domain.IdentityProvider {
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrganizationID:    &orgId,
						ID:                id,
						State:             domain.IDPStateActive,
						Name:              gofakeit.Name(),
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &stylingType,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)
					idp.OrganizationID = gu.Ptr("non-existent-orgID")
					return &idp
				},
				idpIdentifierCondition: idpRepo.IDCondition(id),
				err:                    new(database.NoRowFoundError),
			}
		}(),
		func() test {
			name := gofakeit.Name()
			return test{
				name: "non existent instanceID",
				testFunc: func(t *testing.T) *domain.IdentityProvider {
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrganizationID:    &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive,
						Name:              name,
						Type:              gu.Ptr(domain.IDPTypeOIDC),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &stylingType,
						Payload:           []byte("{}"),
					}

					err := idpRepo.Create(t.Context(), tx, &idp)
					require.NoError(t, err)
					idp.InstanceID = "non-existent-instanceID"
					return &idp
				},
				idpIdentifierCondition: idpRepo.NameCondition(name),
				err:                    new(database.NoRowFoundError),
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var idp *domain.IdentityProvider
			if tt.testFunc != nil {
				idp = tt.testFunc(t)
			}

			// get idp
			returnedIDP, err := idpRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						idpRepo.InstanceIDCondition(idp.InstanceID),
						idpRepo.OrgIDCondition(idp.OrganizationID),
						tt.idpIdentifierCondition,
					),
				),
			)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, returnedIDP.InstanceID, idp.InstanceID)
			assert.Equal(t, returnedIDP.OrganizationID, idp.OrganizationID)
			assert.Equal(t, returnedIDP.State, idp.State)
			assert.Equal(t, returnedIDP.ID, idp.ID)
			assert.Equal(t, returnedIDP.Name, idp.Name)
			assert.Equal(t, returnedIDP.Type, idp.Type)
			assert.Equal(t, returnedIDP.AllowCreation, idp.AllowCreation)
			assert.Equal(t, returnedIDP.AllowAutoCreation, idp.AllowAutoCreation)
			assert.Equal(t, returnedIDP.AllowAutoUpdate, idp.AllowAutoUpdate)
			assert.Equal(t, returnedIDP.AllowLinking, idp.AllowLinking)
			assert.Equal(t, returnedIDP.AutoLinkingField, idp.AutoLinkingField)
			assert.Equal(t, returnedIDP.StylingType, idp.StylingType)
			assert.Equal(t, returnedIDP.Payload, idp.Payload)
		})
	}
}

func TestListIDProvider(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)

	idpRepo := repository.IDProviderRepository()
	var stylingType1 int16 = 1
	var stylingType4 int16 = 4
	payload := []byte(`{"json": {}}`)

	idps := [...]*domain.IdentityProvider{
		{
			InstanceID:        instanceID,
			OrganizationID:    &orgID,
			ID:                "1",
			State:             domain.IDPStateActive,
			Name:              "1",
			Type:              gu.Ptr(domain.IDPTypeOIDC),
			AllowCreation:     false,
			AllowAutoCreation: true,
			AllowAutoUpdate:   true,
			AllowLinking:      true,
			StylingType:       &stylingType1,
			Payload:           []byte("{}"),
		},
		{
			InstanceID:        instanceID,
			OrganizationID:    &orgID,
			ID:                "2",
			State:             domain.IDPStateActive,
			Name:              "2",
			Type:              gu.Ptr(domain.IDPTypeSAML),
			AllowCreation:     true,
			AllowAutoCreation: false,
			AllowAutoUpdate:   true,
			AllowLinking:      true,
			StylingType:       &stylingType4,
			Payload:           []byte("{}"),
		},
		{

			InstanceID:        instanceID,
			OrganizationID:    &orgID,
			ID:                "3",
			State:             domain.IDPStateInactive,
			Name:              "3",
			Type:              gu.Ptr(domain.IDPTypeLDAP),
			AllowCreation:     true,
			AllowAutoCreation: true,
			AllowAutoUpdate:   false,
			AllowLinking:      true,
			AutoLinkingField:  gu.Ptr(domain.IDPAutoLinkingFieldUserName),
			StylingType:       &stylingType1,
			Payload:           []byte("{}"),
		},
		{

			InstanceID:        instanceID,
			OrganizationID:    nil,
			ID:                "4",
			State:             domain.IDPStateActive,
			Name:              "4",
			Type:              gu.Ptr(domain.IDPTypeAzure),
			AllowCreation:     true,
			AllowAutoCreation: true,
			AllowAutoUpdate:   true,
			AllowLinking:      false,
			StylingType:       &stylingType1,
			Payload:           []byte("{}"),
		},
		{

			InstanceID:        instanceID,
			OrganizationID:    nil,
			ID:                "5",
			State:             domain.IDPStateActive,
			Name:              "5",
			Type:              gu.Ptr(domain.IDPTypeApple),
			AllowCreation:     false,
			AllowAutoCreation: true,
			AllowAutoUpdate:   true,
			AllowLinking:      true,
			StylingType:       &stylingType4,
			Payload:           []byte("{}"),
		},
		{

			InstanceID:        instanceID,
			OrganizationID:    nil,
			ID:                "6",
			State:             domain.IDPStateInactive,
			Name:              "6",
			Type:              gu.Ptr(domain.IDPTypeGitHub),
			AllowCreation:     true,
			AllowAutoCreation: false,
			AllowAutoUpdate:   true,
			AllowLinking:      true,
			StylingType:       &stylingType1,
			Payload:           payload,
		},
	}
	for _, idp := range idps {
		err := idpRepo.Create(t.Context(), tx, idp)
		require.NoError(t, err)
	}

	type test struct {
		name      string
		condition database.Condition
		want      []*domain.IdentityProvider
		wantErr   error
	}
	tests := []test{
		{
			name:      "incomplete condition",
			condition: idpRepo.OrgIDCondition(&orgID),
			wantErr:   database.NewMissingConditionError(idpRepo.OrgIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: idpRepo.PrimaryKeyCondition(instanceID, "nix"),
		},
		{
			name:      "primary key condition",
			condition: idpRepo.PrimaryKeyCondition(instanceID, "1"),
			want:      []*domain.IdentityProvider{idps[0]},
		},
		{
			name:      "all from instance",
			condition: idpRepo.InstanceIDCondition(instanceID),
			want:      idps[:],
		},
		{
			name: "all only instance",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.OrgIDCondition(nil),
			),
			want: idps[3:6],
		},
		{
			name: "all from organization",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.OrgIDCondition(&orgID),
			),
			want: idps[0:3],
		},
		{
			name: "id condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.IDCondition("1"),
			),
			want: []*domain.IdentityProvider{idps[0]},
		},
		{
			name: "org and id condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.OrgIDCondition(&orgID),
				idpRepo.IDCondition("2"),
			),
			want: []*domain.IdentityProvider{idps[1]},
		},
		{
			name: "org and id condition, no match",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.OrgIDCondition(nil),
				idpRepo.IDCondition("1"),
			),
			want: []*domain.IdentityProvider{},
		},
		{
			name: "name condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.NameCondition("3"),
			),
			want: []*domain.IdentityProvider{idps[2]},
		},
		{
			name: "name condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.StateCondition(domain.IDPStateInactive),
			),
			want: []*domain.IdentityProvider{idps[2], idps[5]},
		},
		{
			name: "allow creation condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.AllowCreationCondition(false),
			),
			want: []*domain.IdentityProvider{idps[0], idps[4]},
		},
		{
			name: "allow auto creation condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.AllowAutoCreationCondition(false),
			),
			want: []*domain.IdentityProvider{idps[1], idps[5]},
		},
		{
			name: "allow update condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.AllowAutoUpdateCondition(false),
			),
			want: []*domain.IdentityProvider{idps[2]},
		},
		{
			name: "allow linking condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.AllowLinkingCondition(false),
			),
			want: []*domain.IdentityProvider{idps[3]},
		},
		{
			name: "allow auto linking condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.AllowAutoLinkingCondition(domain.IDPAutoLinkingFieldUserName),
			),
			want: []*domain.IdentityProvider{idps[2]},
		},
		{
			name: "styling condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.StylingTypeCondition(stylingType4),
			),
			want: []*domain.IdentityProvider{idps[1], idps[4]},
		},
		{
			name: "payload condition",
			condition: database.And(
				idpRepo.InstanceIDCondition(instanceID),
				idpRepo.PayloadCondition(string(payload)),
			),
			want: []*domain.IdentityProvider{idps[5]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := idpRepo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(idpRepo.PrimaryKeyColumns()...),
			)
			require.ErrorIs(t, err, tt.wantErr)

			if !assert.Len(t, got, len(tt.want)) {
				return
			}
			for i, want := range tt.want {
				assert.Equal(t, got[i].InstanceID, want.InstanceID)
				assert.Equal(t, got[i].OrganizationID, want.OrganizationID)
				assert.Equal(t, got[i].State, want.State)
				assert.Equal(t, got[i].ID, want.ID)
				assert.Equal(t, got[i].Name, want.Name)
				assert.Equal(t, got[i].Type, want.Type)
				assert.Equal(t, got[i].AllowCreation, want.AllowCreation)
				assert.Equal(t, got[i].AllowAutoCreation, want.AllowAutoCreation)
				assert.Equal(t, got[i].AllowAutoUpdate, want.AllowAutoUpdate)
				assert.Equal(t, got[i].AllowLinking, want.AllowLinking)
				assert.Equal(t, got[i].AutoLinkingField, want.AutoLinkingField)
				assert.Equal(t, got[i].StylingType, want.StylingType)
				assert.Equal(t, got[i].Payload, want.Payload)
			}
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
							OrganizationID:    &orgId,
							ID:                id,
							State:             domain.IDPStateActive,
							Name:              gofakeit.Name(),
							Type:              gu.Ptr(domain.IDPTypeOIDC),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							StylingType:       &stylingType,
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
		func() test {
			name := gofakeit.Name()
			var noOfIDPs int64 = 1
			return test{
				name: "happy path delete idp filter name",
				testFunc: func(t *testing.T) {
					for range noOfIDPs {
						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrganizationID:    &orgId,
							ID:                gofakeit.Name(),
							State:             domain.IDPStateActive,
							Name:              name,
							Type:              gu.Ptr(domain.IDPTypeOIDC),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							StylingType:       &stylingType,
							Payload:           []byte("{}"),
						}

						err := idpRepo.Create(t.Context(), tx, &idp)
						require.NoError(t, err)

					}
				},
				idpIdentifierCondition: idpRepo.NameCondition(name),
				noOfDeletedRows:        noOfIDPs,
			}
		}(),
		{
			name:                   "delete non existent idp",
			idpIdentifierCondition: idpRepo.NameCondition(gofakeit.Name()),
		},
		func() test {
			name := gofakeit.Name()
			return test{
				name: "deleted already deleted idp",
				testFunc: func(t *testing.T) {
					noOfIDPs := 1
					for range noOfIDPs {
						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrganizationID:    &orgId,
							ID:                gofakeit.Name(),
							State:             domain.IDPStateActive,
							Name:              name,
							Type:              gu.Ptr(domain.IDPTypeOIDC),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							StylingType:       &stylingType,
							Payload:           []byte("{}"),
						}

						err := idpRepo.Create(t.Context(), tx, &idp)
						require.NoError(t, err)
					}

					// delete organization
					affectedRows, err := idpRepo.Delete(t.Context(), tx,
						database.And(
							idpRepo.NameCondition(name),
							idpRepo.InstanceIDCondition(instanceId),
							idpRepo.OrgIDCondition(&orgId),
						),
					)
					assert.Equal(t, int64(1), affectedRows)
					require.NoError(t, err)
				},
				idpIdentifierCondition: idpRepo.NameCondition(name),
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
					idpRepo.InstanceIDCondition(instanceId),
					idpRepo.OrgIDCondition(&orgId),
					tt.idpIdentifierCondition,
				),
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check idp was deleted
			organization, err := idpRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						idpRepo.InstanceIDCondition(instanceId),
						idpRepo.OrgIDCondition(&orgId),
						tt.idpIdentifierCondition,
					),
				),
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, organization)
		})
	}
}

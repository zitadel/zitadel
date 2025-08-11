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

var sytlingType int16 = 1

func TestCreateIDProvider(t *testing.T) {
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
		State:      domain.OrgStateActive.String(),
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.IdentityProvider
		idp      domain.IdentityProvider
		err      error
	}

	// TESTS
	tests := []test{
		{
			name: "happy path",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrgID:             &orgId,
				ID:                gofakeit.Name(),
				State:             domain.IDPStateActive.String(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &sytlingType,
				Payload:           gu.Ptr("{}"),
			},
		},
		{
			name: "create organization without name",
			idp: domain.IdentityProvider{
				InstanceID: instanceId,
				OrgID:      &orgId,
				ID:         gofakeit.Name(),
				State:      domain.IDPStateActive.String(),
				// Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &sytlingType,
				Payload:           gu.Ptr("{}"),
			},
			err: new(database.CheckError),
		},
		{
			name: "adding org with same id twice",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idpRepo := repository.IDProviderRepository(pool)

				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)
				// change the name to make sure same only the id clashes
				org.Name = gofakeit.Name()
				return &idp
			},
			err: new(database.UniqueError),
		},
		{
			name: "adding org with same name twice",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idpRepo := repository.IDProviderRepository(pool)

				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
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
				name: "adding org with same name, id, different instance",
				testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
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
					instanceRepo := repository.InstanceRepository(pool)
					err := instanceRepo.Create(ctx, &instance)
					assert.Nil(t, err)

					// create org
					newOrgId := gofakeit.Name()
					org := domain.Organization{
						ID:         newOrgId,
						Name:       gofakeit.Name(),
						InstanceID: newInstId,
						State:      domain.OrgStateActive.String(),
					}
					organizationRepo := repository.OrganizationRepository(pool)
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					idpRepo := repository.IDProviderRepository(pool)
					idp := domain.IdentityProvider{
						InstanceID:        newInstId,
						OrgID:             &newOrgId,
						ID:                id,
						State:             domain.IDPStateActive.String(),
						Name:              name,
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err = idpRepo.Create(ctx, &idp)
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
					State:             domain.IDPStateActive.String(),
					Name:              name,
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				},
			}
		}(),
		{
			name: "adding idp with no id",
			idp: domain.IdentityProvider{
				InstanceID: instanceId,
				OrgID:      &orgId,
				// ID:                gofakeit.Name(),
				State:             domain.IDPStateActive.String(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &sytlingType,
				Payload:           gu.Ptr("{}"),
			},
			err: new(database.CheckError),
		},
		{
			name: "adding idp with no name",
			idp: domain.IdentityProvider{
				InstanceID: instanceId,
				OrgID:      &orgId,
				ID:         gofakeit.Name(),
				State:      domain.IDPStateActive.String(),
				// Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &sytlingType,
				Payload:           gu.Ptr("{}"),
			},
			err: new(database.CheckError),
		},
		{
			name: "adding idp with no instance id",
			idp: domain.IdentityProvider{
				// InstanceID:        instanceId,
				OrgID:             &orgId,
				State:             domain.IDPStateActive.String(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &sytlingType,
				Payload:           gu.Ptr("{}"),
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding organization with non existent instance id",
			idp: domain.IdentityProvider{
				InstanceID:        gofakeit.Name(),
				OrgID:             &orgId,
				State:             domain.IDPStateActive.String(),
				ID:                gofakeit.Name(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &sytlingType,
				Payload:           gu.Ptr("{}"),
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding organization with non existent org id",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrgID:             gu.Ptr(gofakeit.Name()),
				State:             domain.IDPStateActive.String(),
				ID:                gofakeit.Name(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       &sytlingType,
				Payload:           gu.Ptr("{}"),
			},
			err: new(database.ForeignKeyError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var idp *domain.IdentityProvider
			if tt.testFunc != nil {
				idp = tt.testFunc(ctx, t)
			} else {
				idp = &tt.idp
			}
			idpRepo := repository.IDProviderRepository(pool)

			// create idp
			beforeCreate := time.Now()
			err = idpRepo.Create(ctx, idp)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check organization values
			idp, err = idpRepo.Get(ctx,
				idpRepo.IDCondition(idp.ID),
				idp.InstanceID,
				idp.OrgID,
			)
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
			assert.Equal(t, tt.idp.StylingType, idp.StylingType)
			assert.Equal(t, tt.idp.Payload, idp.Payload)
			assert.WithinRange(t, idp.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

func TestUpdateIDProvider(t *testing.T) {
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
		State:      domain.OrgStateActive.String(),
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	idpRepo := repository.IDProviderRepository(pool)

	tests := []struct {
		name         string
		testFunc     func(ctx context.Context, t *testing.T) *domain.IdentityProvider
		update       []database.Change
		rowsAffected int64
	}{
		{
			name: "happy path update name",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)
				idp.Name = "new_name"
				return &idp
			},
			update:       []database.Change{idpRepo.SetName("new_name")},
			rowsAffected: 1,
		},
		{
			name: "happy path update state",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)
				idp.State = domain.IDPStateInactive.String()
				return &idp
			},
			update:       []database.Change{idpRepo.SetState(domain.IDPStateInactive)},
			rowsAffected: 1,
		},
		{
			name: "happy path update AllowCreation",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)
				idp.AllowCreation = false
				return &idp
			},
			update:       []database.Change{idpRepo.SetAllowCreation(false)},
			rowsAffected: 1,
		},
		{
			name: "happy path update AllowAutoCreation",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)
				idp.AllowAutoCreation = false
				return &idp
			},
			update:       []database.Change{idpRepo.SetAllowAutoCreation(false)},
			rowsAffected: 1,
		},
		{
			name: "happy path update AllowLinking",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)
				idp.AllowLinking = false
				return &idp
			},
			update:       []database.Change{idpRepo.SetAllowLinking(false)},
			rowsAffected: 1,
		},
		{
			name: "happy path update StylingType",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
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
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)
				idp.Payload = gu.Ptr(`{"json": {}}`)
				return &idp
			},
			// update:       []database.Change{idpRepo.SetPayload("{{}}")},
			update:       []database.Change{idpRepo.SetPayload(`{"json": {}}`)},
			rowsAffected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			organizationRepo := repository.OrganizationRepository(pool)
			idpRepo := repository.IDProviderRepository(pool)

			createdIDP := tt.testFunc(ctx, t)

			// update idp
			beforeUpdate := time.Now()
			rowsAffected, err := idpRepo.Update(ctx,
				idpRepo.IDCondition(createdIDP.ID),
				createdIDP.InstanceID,
				createdIDP.OrgID,
				tt.update...,
			)
			afterUpdate := time.Now()
			require.NoError(t, err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			if rowsAffected == 0 {
				return
			}

			// check idp values
			idp, err := idpRepo.Get(ctx,
				organizationRepo.IDCondition(createdIDP.ID),
				createdIDP.InstanceID,
				createdIDP.OrgID,
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
			assert.Equal(t, createdIDP.StylingType, idp.StylingType)
			assert.Equal(t, createdIDP.Payload, idp.Payload)
			assert.WithinRange(t, idp.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

func TestGetIDProvider(t *testing.T) {
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
		State:      domain.OrgStateActive.String(),
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	// create organization
	// this org is created as an additional org which should NOT
	// be returned in the results of the tests
	preexistingOrg := domain.Organization{
		ID:         gofakeit.Name(),
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive.String(),
	}
	err = organizationRepo.Create(t.Context(), &preexistingOrg)
	require.NoError(t, err)

	idpRepo := repository.IDProviderRepository(pool)
	type test struct {
		name                   string
		testFunc               func(ctx context.Context, t *testing.T) *domain.IdentityProvider
		idpIdentifierCondition domain.OrgIdentifierCondition
		err                    error
	}

	tests := []test{
		func() test {
			id := gofakeit.Name()
			return test{
				name: "happy path get using id",
				testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                id,
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
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
				testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              name,
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)
					return &idp
				},
				idpIdentifierCondition: idpRepo.NameCondition(name),
			}
		}(),
		{
			name: "get using non existent id",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)
				return &idp
			},
			idpIdentifierCondition: idpRepo.IDCondition("non-existent-id"),
			err:                    new(database.NoRowFoundError),
		},
		{
			name: "get using non existent name",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}

				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)
				return &idp
			},
			idpIdentifierCondition: idpRepo.NameCondition("non-existent-name"),
			err:                    new(database.NoRowFoundError),
		},
		////////
		func() test {
			id := gofakeit.Name()
			return test{
				name: "non existent orgID",
				testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                id,
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)
					idp.OrgID = gu.Ptr("non-existent-orgID")
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
				testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              name,
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)
					idp.InstanceID = "non-existent-instnaceID"
					return &idp
				},
				idpIdentifierCondition: idpRepo.NameCondition(name),
				err:                    new(database.NoRowFoundError),
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			var idp *domain.IdentityProvider
			if tt.testFunc != nil {
				idp = tt.testFunc(ctx, t)
			}

			// get idp
			returnedIDP, err := idpRepo.Get(ctx,
				tt.idpIdentifierCondition,
				idp.InstanceID,
				idp.OrgID,
			)
			if err != nil {
				require.ErrorIs(t, tt.err, err)
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
			assert.Equal(t, returnedIDP.StylingType, idp.StylingType)
			assert.Equal(t, returnedIDP.Payload, idp.Payload)
		})
	}
}

// gocognit linting fails due to number of test cases
// and the fact that each test case has a testFunc()
//
//nolint:gocognit
func TestListIDProvider(t *testing.T) {
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
		State:      domain.OrgStateActive.String(),
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	idpRepo := repository.IDProviderRepository(pool)

	type test struct {
		name             string
		testFunc         func(ctx context.Context, t *testing.T) []*domain.IdentityProvider
		conditionClauses []database.Condition
		noIDPsReturned   bool
	}
	tests := []test{
		{
			name: "multiple idps filter on instance",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
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
					State:      domain.OrgStateActive.String(),
				}
				organizationRepo := repository.OrganizationRepository(pool)
				err = organizationRepo.Create(ctx, &org)
				require.NoError(t, err)

				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID:        newInstanceId,
					OrgID:             &newOrgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}
				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.InstanceIDCondition(instanceId)},
		},
		{
			name: "multiple idps filter on org",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				// create org
				newOrgId := gofakeit.Name()
				org := domain.Organization{
					ID:         newOrgId,
					Name:       gofakeit.Name(),
					InstanceID: instanceId,
					State:      domain.OrgStateActive.String(),
				}
				organizationRepo := repository.OrganizationRepository(pool)
				err = organizationRepo.Create(ctx, &org)
				require.NoError(t, err)

				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &newOrgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}
				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.OrgIDCondition(&orgId)},
		},
		{
			name: "happy path single idp no filter",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				noOfIDPs := 1
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
		},
		{
			name: "happy path multiple idps no filter",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
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
				testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
					// create idp
					// this idp is created as an additional idp which should NOT
					// be returned in the results of this test case
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}
					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					noOfIDPs := 1
					idps := make([]*domain.IdentityProvider, noOfIDPs)
					for i := range noOfIDPs {

						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                id,
							State:             domain.IDPStateActive.String(),
							Name:              gofakeit.Name(),
							Type:              domain.IDPTypeOIDC.String(),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							StylingType:       &sytlingType,
							Payload:           gu.Ptr("{}"),
						}

						err := idpRepo.Create(ctx, &idp)
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
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					// state inactive
					State:             domain.IDPStateInactive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}
				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						// state active
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
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
				testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
					// create idp
					// this idp is created as an additional idp which should NOT
					// be returned in the results of this test case
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}
					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					noOfIDPs := 1
					idps := make([]*domain.IdentityProvider, noOfIDPs)
					for i := range noOfIDPs {

						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                gofakeit.Name(),
							State:             domain.IDPStateActive.String(),
							Name:              name,
							Type:              domain.IDPTypeOIDC.String(),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							StylingType:       &sytlingType,
							Payload:           gu.Ptr("{}"),
						}

						err := idpRepo.Create(ctx, &idp)
						require.NoError(t, err)

						idps[i] = &idp
					}

					return idps
				},
				conditionClauses: []database.Condition{idpRepo.NameCondition(name)},
			}
		}(),
		{
			name: "multiple idps filter on type",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateInactive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}
				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						State:      domain.IDPStateActive.String(),
						Name:       gofakeit.Name(),
						// type LDAP
						Type:              domain.IDPTypeLDAP.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.TypeCondition(domain.IDPTypeLDAP)},
		},
		{
			name: "multiple idps filter on AllowCreation",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateInactive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}
				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID: instanceId,
						OrgID:      &orgId,
						ID:         gofakeit.Name(),
						State:      domain.IDPStateActive.String(),
						Name:       gofakeit.Name(),
						Type:       domain.IDPTypeLDAP.String(),
						// AllowCreation set to false
						AllowCreation:     false,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.AllowCreationCondition(false)},
		},
		{
			name: "multiple idps filter on AllowAutoCreation",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateInactive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}
				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:    instanceId,
						OrgID:         &orgId,
						ID:            gofakeit.Name(),
						State:         domain.IDPStateActive.String(),
						Name:          gofakeit.Name(),
						Type:          domain.IDPTypeLDAP.String(),
						AllowCreation: true,
						// AllowAutoCreation set to false
						AllowAutoCreation: false,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.AllowAutoCreationCondition(false)},
		},
		{
			name: "multiple idps filter on AllowAutoUpdate",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					// state inactive
					State:             domain.IDPStateInactive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}
				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeLDAP.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						// AllowAutoUpdate set to false
						AllowAutoUpdate: false,
						AllowLinking:    true,
						StylingType:     &sytlingType,
						Payload:         gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.AllowAutoUpdateCondition(false)},
		},
		{
			name: "multiple idps filter on AllowLinking",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID: instanceId,
					OrgID:      &orgId,
					ID:         gofakeit.Name(),
					// state inactive
					State:             domain.IDPStateInactive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}
				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)

				noOfIDPs := 5
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeLDAP.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						// AllowLinking set to false
						AllowLinking: false,
						StylingType:  &sytlingType,
						Payload:      gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.AllowLinkingCondition(false)},
		},
		{
			name: "multiple idps filter on StylingType",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
				// create idp
				// this idp is created as an additional idp which should NOT
				// be returned in the results of this test case
				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             &orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       &sytlingType,
					Payload:           gu.Ptr("{}"),
				}
				err := idpRepo.Create(ctx, &idp)
				require.NoError(t, err)

				noOfIDPs := 1
				idps := make([]*domain.IdentityProvider, noOfIDPs)
				sytlingType := int16(4)
				for i := range noOfIDPs {

					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						// StylingType set to 4
						StylingType: &sytlingType,
						Payload:     gu.Ptr("{}"),
					}

					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					idps[i] = &idp
				}

				return idps
			},
			conditionClauses: []database.Condition{idpRepo.StylingTypeCondition(4)},
		},
		func() test {
			payload := `{"json": {}}`
			return test{
				name: "multiple idps filter on Payload",
				testFunc: func(ctx context.Context, t *testing.T) []*domain.IdentityProvider {
					// create idp
					// this idp is created as an additional idp which should NOT
					// be returned in the results of this test case
					idp := domain.IdentityProvider{
						InstanceID:        instanceId,
						OrgID:             &orgId,
						ID:                gofakeit.Name(),
						State:             domain.IDPStateActive.String(),
						Name:              gofakeit.Name(),
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       &sytlingType,
						Payload:           gu.Ptr("{}"),
					}
					err := idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					noOfIDPs := 1
					idps := make([]*domain.IdentityProvider, noOfIDPs)
					for i := range noOfIDPs {

						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                gofakeit.Name(),
							State:             domain.IDPStateActive.String(),
							Name:              gofakeit.Name(),
							Type:              domain.IDPTypeOIDC.String(),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							StylingType:       &sytlingType,
							Payload:           &payload,
						}

						err := idpRepo.Create(ctx, &idp)
						require.NoError(t, err)

						idps[i] = &idp
					}

					return idps
				},
				conditionClauses: []database.Condition{idpRepo.PayloadCondition(payload)},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, err := pool.Exec(ctx, "DELETE FROM zitadel.identity_providers")
				require.NoError(t, err)
			})

			ctx := context.WithoutCancel(t.Context())

			idps := tt.testFunc(ctx, t)

			// check idp values
			returnedIDPs, err := idpRepo.List(ctx,
				tt.conditionClauses...,
			)
			require.NoError(t, err)
			if tt.noIDPsReturned {
				assert.Nil(t, returnedIDPs)
				return
			}

			assert.Equal(t, len(idps), len(returnedIDPs))
			for i, idp := range idps {

				assert.Equal(t, returnedIDPs[i].InstanceID, idp.InstanceID)
				assert.Equal(t, returnedIDPs[i].OrgID, idp.OrgID)
				assert.Equal(t, returnedIDPs[i].State, idp.State)
				assert.Equal(t, returnedIDPs[i].ID, idp.ID)
				assert.Equal(t, returnedIDPs[i].Name, idp.Name)
				assert.Equal(t, returnedIDPs[i].Type, idp.Type)
				assert.Equal(t, returnedIDPs[i].AllowCreation, idp.AllowCreation)
				assert.Equal(t, returnedIDPs[i].AllowAutoCreation, idp.AllowAutoCreation)
				assert.Equal(t, returnedIDPs[i].AllowAutoUpdate, idp.AllowAutoUpdate)
				assert.Equal(t, returnedIDPs[i].AllowLinking, idp.AllowLinking)
				assert.Equal(t, returnedIDPs[i].StylingType, idp.StylingType)
				assert.Equal(t, returnedIDPs[i].Payload, idp.Payload)
			}
		})
	}
}

func TestDeleteIDProvider(t *testing.T) {
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
		State:      domain.OrgStateActive.String(),
	}
	organizationRepo := repository.OrganizationRepository(pool)
	err = organizationRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	idpRepo := repository.IDProviderRepository(pool)

	type test struct {
		name                   string
		testFunc               func(ctx context.Context, t *testing.T)
		idpIdentifierCondition domain.IDPIdentifierCondition
		noOfDeletedRows        int64
	}
	tests := []test{
		func() test {
			id := gofakeit.Name()
			var noOfIDPs int64 = 1
			return test{
				name: "happy path delete idp filter id",
				testFunc: func(ctx context.Context, t *testing.T) {
					for range noOfIDPs {
						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                id,
							State:             domain.IDPStateActive.String(),
							Name:              gofakeit.Name(),
							Type:              domain.IDPTypeOIDC.String(),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							StylingType:       &sytlingType,
							Payload:           gu.Ptr("{}"),
						}

						err := idpRepo.Create(ctx, &idp)
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
				testFunc: func(ctx context.Context, t *testing.T) {
					for range noOfIDPs {
						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                gofakeit.Name(),
							State:             domain.IDPStateActive.String(),
							Name:              name,
							Type:              domain.IDPTypeOIDC.String(),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							StylingType:       &sytlingType,
							Payload:           gu.Ptr("{}"),
						}

						err := idpRepo.Create(ctx, &idp)
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
				testFunc: func(ctx context.Context, t *testing.T) {
					noOfIDPs := 1
					for range noOfIDPs {
						idp := domain.IdentityProvider{
							InstanceID:        instanceId,
							OrgID:             &orgId,
							ID:                gofakeit.Name(),
							State:             domain.IDPStateActive.String(),
							Name:              name,
							Type:              domain.IDPTypeOIDC.String(),
							AllowCreation:     true,
							AllowAutoCreation: true,
							AllowAutoUpdate:   true,
							AllowLinking:      true,
							StylingType:       &sytlingType,
							Payload:           gu.Ptr("{}"),
						}

						err := idpRepo.Create(ctx, &idp)
						require.NoError(t, err)
					}

					// delete organization
					affectedRows, err := idpRepo.Delete(ctx,
						idpRepo.NameCondition(name),
						instanceId,
						&orgId,
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
			ctx := t.Context()

			if tt.testFunc != nil {
				tt.testFunc(ctx, t)
			}

			// delete idp
			noOfDeletedRows, err := idpRepo.Delete(ctx,
				tt.idpIdentifierCondition,
				instanceId,
				&orgId,
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check idp was deleted
			organization, err := idpRepo.Get(ctx,
				tt.idpIdentifierCondition,
				instanceId,
				&orgId,
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, organization)
		})
	}
}

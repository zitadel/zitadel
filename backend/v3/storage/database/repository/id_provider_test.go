package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

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
				OrgID:             orgId,
				ID:                gofakeit.Name(),
				State:             domain.IDPStateActive.String(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       1,
				Payload:           "{}",
			},
		},
		{
			name: "create organization without name",
			idp: domain.IdentityProvider{
				InstanceID: instanceId,
				OrgID:      orgId,
				ID:         gofakeit.Name(),
				State:      domain.IDPStateActive.String(),
				// Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       1,
				Payload:           "{}",
			},
			err: new(database.CheckError),
		},
		{
			name: "adding org with same id twice",
			testFunc: func(ctx context.Context, t *testing.T) *domain.IdentityProvider {
				idpRepo := repository.IDProviderRepository(pool)

				idp := domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       1,
					Payload:           "{}",
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
					OrgID:             orgId,
					ID:                gofakeit.Name(),
					State:             domain.IDPStateActive.String(),
					Name:              gofakeit.Name(),
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       1,
					Payload:           "{}",
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
					err = organizationRepo.Create(t.Context(), &org)
					require.NoError(t, err)

					idpRepo := repository.IDProviderRepository(pool)
					idp := domain.IdentityProvider{
						InstanceID:        newInstId,
						OrgID:             newOrgId,
						ID:                id,
						State:             domain.IDPStateActive.String(),
						Name:              name,
						Type:              domain.IDPTypeOIDC.String(),
						AllowCreation:     true,
						AllowAutoCreation: true,
						AllowAutoUpdate:   true,
						AllowLinking:      true,
						StylingType:       1,
						Payload:           "{}",
					}

					err = idpRepo.Create(ctx, &idp)
					require.NoError(t, err)

					// change the instanceID to a different instance
					idp.InstanceID = instanceId
					// change the OrgId to a different organization
					idp.OrgID = orgId
					return &idp
				},
				idp: domain.IdentityProvider{
					InstanceID:        instanceId,
					OrgID:             orgId,
					ID:                id,
					State:             domain.IDPStateActive.String(),
					Name:              name,
					Type:              domain.IDPTypeOIDC.String(),
					AllowCreation:     true,
					AllowAutoCreation: true,
					AllowAutoUpdate:   true,
					AllowLinking:      true,
					StylingType:       1,
					Payload:           "{}",
				},
			}
		}(),
		{
			name: "adding idp with no id",
			idp: domain.IdentityProvider{
				InstanceID: instanceId,
				OrgID:      orgId,
				// ID:                gofakeit.Name(),
				State:             domain.IDPStateActive.String(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       1,
				Payload:           "{}",
			},
			err: new(database.CheckError),
		},
		{
			name: "adding idp with no name",
			idp: domain.IdentityProvider{
				InstanceID: instanceId,
				OrgID:      orgId,
				ID:         gofakeit.Name(),
				State:      domain.IDPStateActive.String(),
				// Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       1,
				Payload:           "{}",
			},
			err: new(database.CheckError),
		},
		{
			name: "adding idp with no instance id",
			idp: domain.IdentityProvider{
				// InstanceID:        instanceId,
				OrgID:             orgId,
				State:             domain.IDPStateActive.String(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       1,
				Payload:           "{}",
			},
			err: new(database.IntegrityViolationError),
		},
		{
			name: "adding organization with non existent instance id",
			idp: domain.IdentityProvider{
				InstanceID:        gofakeit.Name(),
				OrgID:             orgId,
				State:             domain.IDPStateActive.String(),
				ID:                gofakeit.Name(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       1,
				Payload:           "{}",
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "adding organization with non existent org id",
			idp: domain.IdentityProvider{
				InstanceID:        instanceId,
				OrgID:             gofakeit.Name(),
				State:             domain.IDPStateActive.String(),
				ID:                gofakeit.Name(),
				Name:              gofakeit.Name(),
				Type:              domain.IDPTypeOIDC.String(),
				AllowCreation:     true,
				AllowAutoCreation: true,
				AllowAutoUpdate:   true,
				AllowLinking:      true,
				StylingType:       1,
				Payload:           "{}",
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

			assert.Equal(t, tt.idp.ID, idp.ID)
			assert.Equal(t, tt.idp.Name, idp.Name)
			assert.Equal(t, tt.idp.InstanceID, idp.InstanceID)
			assert.Equal(t, tt.idp.State, idp.State)
			assert.WithinRange(t, idp.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

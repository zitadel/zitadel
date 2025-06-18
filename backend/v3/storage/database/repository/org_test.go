package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestCreateOrganization(t *testing.T) {
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
	assert.Nil(t, err)

	tests := []struct {
		name         string
		testFunc     func(ctx context.Context, t *testing.T) *domain.Organization
		organization domain.Organization
		err          error
	}{
		{
			name: "happy path",
			organization: func() domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				organization := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive.String(),
				}
				return organization
			}(),
		},
		{
			name: "create organization without name",
			organization: func() domain.Organization {
				organizationId := gofakeit.Name()
				// organizationName := gofakeit.Name()
				organization := domain.Organization{
					ID:         organizationId,
					Name:       "",
					InstanceID: instanceId,
					State:      domain.OrgStateActive.String(),
				}
				return organization
			}(),
			err: errors.New("organization name not provided"),
		},
		{
			name: "adding same organization twice",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Organization {
				organizationRepo := repository.OrganizationRepository(pool)
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()

				inst := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive.String(),
				}

				err := organizationRepo.Create(ctx, &inst)
				require.NoError(t, err)
				return &inst
			},
			err: errors.New("organization id already exists"),
		},
		{
			name: "adding organization with no id",
			organization: func() domain.Organization {
				// organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				organization := domain.Organization{
					// ID:              organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive.String(),
				}
				return organization
			}(),
			err: errors.New("organization id not provided"),
		},
		{
			name: "adding organization with no instance id",
			organization: func() domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				organization := domain.Organization{
					ID:    organizationId,
					Name:  organizationName,
					State: domain.OrgStateActive.String(),
				}
				return organization
			}(),
			err: errors.New("instance id not provided"),
		},
		{
			name: "adding organization with non existent instance id",
			organization: func() domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				organization := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: gofakeit.Name(),
					State:      domain.OrgStateActive.String(),
				}
				return organization
			}(),
			err: errors.New("invalid instance id"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var organization *domain.Organization
			if tt.testFunc != nil {
				organization = tt.testFunc(ctx, t)
			} else {
				organization = &tt.organization
			}
			organizationRepo := repository.OrganizationRepository(pool)

			// create organization
			beforeCreate := time.Now()
			err = organizationRepo.Create(ctx, organization)
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check organization values
			organization, err = organizationRepo.Get(ctx,
				organizationRepo.NameCondition(database.TextOperationEqual, organization.Name),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.organization.ID, organization.ID)
			assert.Equal(t, tt.organization.Name, organization.Name)
			assert.Equal(t, tt.organization.InstanceID, organization.InstanceID)
			assert.Equal(t, tt.organization.State, organization.State)
			assert.WithinRange(t, organization.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, organization.UpdatedAt, beforeCreate, afterCreate)
			assert.Nil(t, organization.DeletedAt)
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
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
	assert.Nil(t, err)
	organizationRepo := repository.OrganizationRepository(pool)

	tests := []struct {
		name         string
		testFunc     func(ctx context.Context, t *testing.T) *domain.Organization
		update       []database.Change
		rowsAffected int64
	}{
		{
			name: "happy path update name",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()

				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive.String(),
				}

				// create organization
				err := organizationRepo.Create(ctx, &org)
				require.NoError(t, err)

				// update with updated value
				org.Name = "new_name"
				return &org
			},
			update:       []database.Change{organizationRepo.SetName("new_name")},
			rowsAffected: 1,
		},
		{
			name: "update deleted organization",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()

				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive.String(),
				}

				// create organization
				err := organizationRepo.Create(ctx, &org)
				require.NoError(t, err)

				// delete instance
				_, err = organizationRepo.Delete(ctx,
					organizationRepo.IDCondition(org.ID),
				)
				require.NoError(t, err)

				return &org
			},
			update:       []database.Change{organizationRepo.SetName("new_name")},
			rowsAffected: 0,
		},
		{
			name: "happy path change state",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()

				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive.String(),
				}

				// create organization
				err := organizationRepo.Create(ctx, &org)
				require.NoError(t, err)

				// update with updated value
				org.State = domain.OrgStateInactive.String()
				return &org
			},
			update:       []database.Change{organizationRepo.SetState(domain.OrgStateInactive)},
			rowsAffected: 1,
		},
		{
			name: "update non existent organization",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Organization {
				organizationId := gofakeit.Name()

				org := domain.Organization{
					ID: organizationId,
				}
				return &org
			},
			update:       []database.Change{organizationRepo.SetName("new_name")},
			rowsAffected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			organizationRepo := repository.OrganizationRepository(pool)

			createdOrg := tt.testFunc(ctx, t)

			// update org
			beforeUpdate := time.Now()
			rowsAffected, err := organizationRepo.Update(ctx,
				organizationRepo.IDCondition(createdOrg.ID),
				tt.update...,
			)
			afterUpdate := time.Now()
			require.NoError(t, err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			if rowsAffected == 0 {
				return
			}

			// check organization values
			organization, err := organizationRepo.Get(ctx,
				organizationRepo.IDCondition(createdOrg.ID),
			)
			require.NoError(t, err)

			assert.Equal(t, createdOrg.ID, organization.ID)
			assert.Equal(t, createdOrg.Name, organization.Name)
			assert.Equal(t, createdOrg.State, organization.State)
			assert.WithinRange(t, organization.UpdatedAt, beforeUpdate, afterUpdate)
			assert.Nil(t, organization.DeletedAt)
		})
	}
}

func TestGetOrganization(t *testing.T) {
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
	assert.Nil(t, err)

	orgRepo := repository.OrganizationRepository(pool)

	// create organization
	// this org is created as an additional org which should NOT
	// be returned in the results of the tests
	org := domain.Organization{
		ID:         gofakeit.Name(),
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive.String(),
	}
	err = orgRepo.Create(t.Context(), &org)
	require.NoError(t, err)

	type test struct {
		name             string
		testFunc         func(ctx context.Context, t *testing.T) *domain.Organization
		conditionClauses []database.Condition
	}

	tests := []test{
		func() test {
			organizationId := gofakeit.Name()
			return test{
				name: "happy path get using id",
				testFunc: func(ctx context.Context, t *testing.T) *domain.Organization {
					organizationName := gofakeit.Name()

					org := domain.Organization{
						ID:         organizationId,
						Name:       organizationName,
						InstanceID: instanceId,
						State:      domain.OrgStateActive.String(),
					}

					// create organization
					err := orgRepo.Create(ctx, &org)
					require.NoError(t, err)

					return &org
				},
				conditionClauses: []database.Condition{orgRepo.IDCondition(organizationId)},
			}
		}(),
		func() test {
			organizationName := gofakeit.Name()
			return test{
				name: "happy path get using name",
				testFunc: func(ctx context.Context, t *testing.T) *domain.Organization {
					organizationId := gofakeit.Name()

					org := domain.Organization{
						ID:         organizationId,
						Name:       organizationName,
						InstanceID: instanceId,
						State:      domain.OrgStateActive.String(),
					}

					// create organization
					err := orgRepo.Create(ctx, &org)
					require.NoError(t, err)

					return &org
				},
				conditionClauses: []database.Condition{orgRepo.NameCondition(database.TextOperationEqual, organizationName)},
			}
		}(),
		{
			name: "get non existent organization",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Organization {
				_ = domain.Instance{
					ID: instanceId,
				}
				return nil
			},
			conditionClauses: []database.Condition{orgRepo.NameCondition(database.TextOperationEqual, "non-existent-instance-name")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			orgRepo := repository.OrganizationRepository(pool)

			var org *domain.Organization
			if tt.testFunc != nil {
				org = tt.testFunc(ctx, t)
			}

			// get org values
			returnedOrg, err := orgRepo.Get(ctx,
				tt.conditionClauses...,
			)
			require.NoError(t, err)
			if org == nil {
				assert.Nil(t, org, returnedOrg)
				return
			}

			assert.Equal(t, returnedOrg.ID, org.ID)
			assert.Equal(t, returnedOrg.Name, org.Name)
			assert.Equal(t, returnedOrg.InstanceID, org.InstanceID)
			assert.Equal(t, returnedOrg.State, org.State)
		})
	}
}

func TestListOrganization(t *testing.T) {
	ctx := t.Context()
	pool, stop, err := newEmbeddedDB(ctx)
	require.NoError(t, err)
	defer stop()
	organizationRepo := repository.OrganizationRepository(pool)

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
	assert.Nil(t, err)

	type test struct {
		name                   string
		testFunc               func(ctx context.Context, t *testing.T) []*domain.Organization
		conditionClauses       []database.Condition
		noOrganizationReturned bool
	}
	tests := []test{
		{
			name: "happy path single organization no filter",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Organization {
				noOfOrganizations := 1
				organizations := make([]*domain.Organization, noOfOrganizations)
				for i := range noOfOrganizations {

					inst := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: instanceId,
						State:      domain.OrgStateActive.String(),
					}

					// create organization
					err := organizationRepo.Create(ctx, &inst)
					require.NoError(t, err)

					organizations[i] = &inst
				}

				return organizations
			},
		},
		{
			name: "happy path multiple organization no filter",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Organization {
				noOfOrganizations := 5
				organizations := make([]*domain.Organization, noOfOrganizations)
				for i := range noOfOrganizations {

					inst := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: instanceId,
						State:      domain.OrgStateActive.String(),
					}

					// create organization
					err := organizationRepo.Create(ctx, &inst)
					require.NoError(t, err)

					organizations[i] = &inst
				}

				return organizations
			},
		},
		func() test {
			organizationId := gofakeit.Name()
			return test{
				name: "organization filter on id",
				testFunc: func(ctx context.Context, t *testing.T) []*domain.Organization {
					// create organization
					// this org is created as an additional org which should NOT
					// be returned in the results of this test case
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: instanceId,
						State:      domain.OrgStateActive.String(),
					}
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					noOfOrganizations := 1
					organizations := make([]*domain.Organization, noOfOrganizations)
					for i := range noOfOrganizations {

						inst := domain.Organization{
							ID:         organizationId,
							Name:       gofakeit.Name(),
							InstanceID: instanceId,
							State:      domain.OrgStateActive.String(),
						}

						// create organization
						err := organizationRepo.Create(ctx, &inst)
						require.NoError(t, err)

						organizations[i] = &inst
					}

					return organizations
				},
				conditionClauses: []database.Condition{organizationRepo.IDCondition(organizationId)},
			}
		}(),
		{
			name: "multiple organization filter on state",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Organization {
				// create organization
				// this org is created as an additional org which should NOT
				// be returned in the results of this test case
				org := domain.Organization{
					ID:         gofakeit.Name(),
					Name:       gofakeit.Name(),
					InstanceID: instanceId,
					State:      domain.OrgStateActive.String(),
				}
				err = organizationRepo.Create(ctx, &org)
				require.NoError(t, err)

				noOfOrganizations := 5
				organizations := make([]*domain.Organization, noOfOrganizations)
				for i := range noOfOrganizations {

					inst := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: instanceId,
						State:      domain.OrgStateInactive.String(),
					}

					// create organization
					err := organizationRepo.Create(ctx, &inst)
					require.NoError(t, err)

					organizations[i] = &inst
				}

				return organizations
			},
			conditionClauses: []database.Condition{organizationRepo.StateCondition(domain.OrgStateInactive)},
		},
		func() test {
			organizationName := gofakeit.Name()
			return test{
				name: "multiple organization filter on name",
				testFunc: func(ctx context.Context, t *testing.T) []*domain.Organization {
					// create organization
					// this org is created as an additional org which should NOT
					// be returned in the results of this test case
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: instanceId,
						State:      domain.OrgStateActive.String(),
					}
					err = organizationRepo.Create(ctx, &org)
					require.NoError(t, err)

					noOfOrganizations := 5
					organizations := make([]*domain.Organization, noOfOrganizations)
					for i := range noOfOrganizations {

						inst := domain.Organization{
							ID:         gofakeit.Name(),
							Name:       organizationName,
							InstanceID: instanceId,
							State:      domain.OrgStateActive.String(),
						}

						// create organization
						err := organizationRepo.Create(ctx, &inst)
						require.NoError(t, err)

						organizations[i] = &inst
					}

					return organizations
				},
				conditionClauses: []database.Condition{organizationRepo.NameCondition(database.TextOperationEqual, organizationName)},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, err := pool.Exec(ctx, "DELETE FROM zitadel.organizations")
				require.NoError(t, err)
			})

			organizations := tt.testFunc(ctx, t)

			// check organization values
			returnedOrgs, err := organizationRepo.List(ctx,
				tt.conditionClauses...,
			)
			require.NoError(t, err)
			if tt.noOrganizationReturned {
				assert.Nil(t, returnedOrgs)
				return
			}

			assert.Equal(t, len(organizations), len(returnedOrgs))
			for i, org := range organizations {
				assert.Equal(t, returnedOrgs[i].ID, org.ID)
				assert.Equal(t, returnedOrgs[i].Name, org.Name)
				assert.Equal(t, returnedOrgs[i].InstanceID, org.InstanceID)
				assert.Equal(t, returnedOrgs[i].State, org.State)
			}
		})
	}
}

func TestDeleteOrganization(t *testing.T) {
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
	assert.Nil(t, err)

	type test struct {
		name             string
		testFunc         func(ctx context.Context, t *testing.T)
		conditionClauses database.Condition
		noOfDeletedRows  int64
	}
	tests := []test{
		func() test {
			organizationRepo := repository.OrganizationRepository(pool)
			organizationId := gofakeit.Name()
			var noOfOrganizations int64 = 1
			return test{
				name: "happy path delete single organization filter id",
				testFunc: func(ctx context.Context, t *testing.T) {
					organizations := make([]*domain.Organization, noOfOrganizations)
					for i := range noOfOrganizations {

						inst := domain.Organization{
							ID:         organizationId,
							Name:       gofakeit.Name(),
							InstanceID: instanceId,
							State:      domain.OrgStateActive.String(),
						}

						// create organization
						err := organizationRepo.Create(ctx, &inst)
						require.NoError(t, err)

						organizations[i] = &inst
					}
				},
				conditionClauses: organizationRepo.IDCondition(organizationId),
				noOfDeletedRows:  noOfOrganizations,
			}
		}(),
		func() test {
			organizationRepo := repository.OrganizationRepository(pool)
			organizationName := gofakeit.Name()
			var noOfOrganizations int64 = 1
			return test{
				name: "happy path delete single organization filter name",
				testFunc: func(ctx context.Context, t *testing.T) {
					organizations := make([]*domain.Organization, noOfOrganizations)
					for i := range noOfOrganizations {

						inst := domain.Organization{
							ID:         gofakeit.Name(),
							Name:       organizationName,
							InstanceID: instanceId,
							State:      domain.OrgStateActive.String(),
						}

						// create organization
						err := organizationRepo.Create(ctx, &inst)
						require.NoError(t, err)

						organizations[i] = &inst
					}
				},
				conditionClauses: organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
				noOfDeletedRows:  noOfOrganizations,
			}
		}(),
		func() test {
			organizationRepo := repository.OrganizationRepository(pool)
			non_existent_organization_name := gofakeit.Name()
			return test{
				name:             "delete non existent organization",
				conditionClauses: organizationRepo.NameCondition(database.TextOperationEqual, non_existent_organization_name),
			}
		}(),
		func() test {
			organizationRepo := repository.OrganizationRepository(pool)
			organizationName := gofakeit.Name()
			var noOfOrganizations int64 = 5
			return test{
				name: "multiple organization filter on name",
				testFunc: func(ctx context.Context, t *testing.T) {
					organizations := make([]*domain.Organization, noOfOrganizations)
					for i := range noOfOrganizations {

						inst := domain.Organization{
							ID:         gofakeit.Name(),
							Name:       organizationName,
							InstanceID: instanceId,
							State:      domain.OrgStateActive.String(),
						}

						// create organization
						err := organizationRepo.Create(ctx, &inst)
						require.NoError(t, err)

						organizations[i] = &inst
					}
				},
				conditionClauses: organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
				noOfDeletedRows:  noOfOrganizations,
			}
		}(),
		func() test {
			organizationRepo := repository.OrganizationRepository(pool)
			organizationName := gofakeit.Name()
			return test{
				name: "deleted already deleted organization",
				testFunc: func(ctx context.Context, t *testing.T) {
					noOfOrganizations := 1
					organizations := make([]*domain.Organization, noOfOrganizations)
					for i := range noOfOrganizations {

						inst := domain.Organization{
							ID:         gofakeit.Name(),
							Name:       organizationName,
							InstanceID: instanceId,
							State:      domain.OrgStateActive.String(),
						}

						// create organization
						err := organizationRepo.Create(ctx, &inst)
						require.NoError(t, err)

						organizations[i] = &inst
					}

					// delete organization
					affectedRows, err := organizationRepo.Delete(ctx,
						organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
					)
					assert.Equal(t, int64(1), affectedRows)
					require.NoError(t, err)
				},
				conditionClauses: organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
				// this test should return 0 affected rows as the org was already deleted
				noOfDeletedRows: 0,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			organizationRepo := repository.OrganizationRepository(pool)

			if tt.testFunc != nil {
				tt.testFunc(ctx, t)
			}

			// delete organization
			noOfDeletedRows, err := organizationRepo.Delete(ctx,
				tt.conditionClauses,
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check organization was deleted
			organization, err := organizationRepo.Get(ctx,
				tt.conditionClauses,
			)
			require.NoError(t, err)
			assert.Nil(t, organization)
		})
	}
}

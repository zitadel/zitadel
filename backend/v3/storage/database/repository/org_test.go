package repository_test

import (
	"strconv"
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
	organizationRepo := repository.OrganizationRepository()
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
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	tests := []struct {
		name         string
		testFunc     func(t *testing.T, client database.QueryExecutor) *domain.Organization
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
					State:      domain.OrgStateActive,
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
					State:      domain.OrgStateActive,
				}
				return organization
			}(),
			err: new(database.CheckError),
		},
		{
			name: "adding org with same id twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()

				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				}

				err := organizationRepo.Create(t.Context(), tx, &org)
				require.NoError(t, err)
				// change the name to make sure same only the id clashes
				org.Name = gofakeit.Name()
				return &org
			},
			err: new(database.UniqueError),
		},
		{
			name: "adding org with same name twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()

				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				}

				err := organizationRepo.Create(t.Context(), tx, &org)
				require.NoError(t, err)
				// change the id to make sure same name+instance causes an error
				org.ID = gofakeit.Name()
				return &org
			},
			err: new(database.UniqueError),
		},
		func() struct {
			name         string
			testFunc     func(t *testing.T, tx database.QueryExecutor) *domain.Organization
			organization domain.Organization
			err          error
		} {
			orgID := gofakeit.Name()
			organizationName := gofakeit.Name()

			return struct {
				name         string
				testFunc     func(t *testing.T, tx database.QueryExecutor) *domain.Organization
				organization domain.Organization
				err          error
			}{
				name: "adding org with same name, different instance",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Organization {
					// create instance
					instId := gofakeit.Name()
					instance := domain.Instance{
						ID:              instId,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					err := instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       organizationName,
						InstanceID: instId,
						State:      domain.OrgStateActive,
					}

					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// change the id to make it unique
					org.ID = orgID
					// change the instanceID to a different instance
					org.InstanceID = instanceId
					return &org
				},
				organization: domain.Organization{
					ID:         orgID,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				},
			}
		}(),
		{
			name: "adding organization with no id",
			organization: func() domain.Organization {
				// organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				organization := domain.Organization{
					// ID:              organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				}
				return organization
			}(),
			err: new(database.CheckError),
		},
		{
			name: "adding organization with no instance id",
			organization: func() domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				organization := domain.Organization{
					ID:    organizationId,
					Name:  organizationName,
					State: domain.OrgStateActive,
				}
				return organization
			}(),
			err: new(database.ForeignKeyError),
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
					State:      domain.OrgStateActive,
				}
				return organization
			}(),
			err: new(database.ForeignKeyError),
		},
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

			var organization *domain.Organization
			if tt.testFunc != nil {
				organization = tt.testFunc(t, savepoint)
			} else {
				organization = &tt.organization
			}

			// create organization

			err = organizationRepo.Create(t.Context(), savepoint, organization)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check organization values
			organization, err = organizationRepo.Get(t.Context(), savepoint,
				database.WithCondition(
					database.And(
						organizationRepo.IDCondition(organization.ID),
						organizationRepo.InstanceIDCondition(organization.InstanceID),
					),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.organization.ID, organization.ID)
			assert.Equal(t, tt.organization.Name, organization.Name)
			assert.Equal(t, tt.organization.InstanceID, organization.InstanceID)
			assert.Equal(t, tt.organization.State, organization.State)
			assert.WithinRange(t, organization.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, organization.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	beforeUpdate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()
	organizationRepo := repository.OrganizationRepository()

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
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	tests := []struct {
		name         string
		testFunc     func(t *testing.T) *domain.Organization
		update       []database.Change
		rowsAffected int64
	}{
		{
			name: "happy path update name",
			testFunc: func(t *testing.T) *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()

				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				}

				// create organization
				err := organizationRepo.Create(t.Context(), tx, &org)
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
			testFunc: func(t *testing.T) *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()

				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				}

				// create organization
				err := organizationRepo.Create(t.Context(), tx, &org)
				require.NoError(t, err)

				// delete instance
				_, err = organizationRepo.Delete(t.Context(), tx,
					database.And(
						organizationRepo.InstanceIDCondition(org.InstanceID),
						organizationRepo.IDCondition(org.ID),
					),
				)
				require.NoError(t, err)

				return &org
			},
			update:       []database.Change{organizationRepo.SetName("new_name")},
			rowsAffected: 0,
		},
		{
			name: "happy path change state",
			testFunc: func(t *testing.T) *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()

				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				}

				// create organization
				err := organizationRepo.Create(t.Context(), tx, &org)
				require.NoError(t, err)

				// update with updated value
				org.State = domain.OrgStateInactive
				return &org
			},
			update:       []database.Change{organizationRepo.SetState(domain.OrgStateInactive)},
			rowsAffected: 1,
		},
		{
			name: "update non existent organization",
			testFunc: func(t *testing.T) *domain.Organization {
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
			createdOrg := tt.testFunc(t)

			// update org
			rowsAffected, err := organizationRepo.Update(t.Context(), tx,
				database.And(
					organizationRepo.InstanceIDCondition(createdOrg.InstanceID),
					organizationRepo.IDCondition(createdOrg.ID),
				),
				tt.update...,
			)
			afterUpdate := time.Now()
			require.NoError(t, err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			if rowsAffected == 0 {
				return
			}

			// check organization values
			organization, err := organizationRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						organizationRepo.IDCondition(createdOrg.ID),
						organizationRepo.InstanceIDCondition(createdOrg.InstanceID),
					),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, createdOrg.ID, organization.ID)
			assert.Equal(t, createdOrg.Name, organization.Name)
			assert.Equal(t, createdOrg.State, organization.State)
			assert.WithinRange(t, organization.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

func TestGetOrganization(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()
	orgRepo := repository.OrganizationRepository()
	orgDomainRepo := repository.OrganizationDomainRepository()

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

	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	// create organization
	// this org is created as an additional org which should NOT
	// be returned in the results of the tests
	org := domain.Organization{
		ID:         gofakeit.Name(),
		Name:       gofakeit.Name(),
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	type test struct {
		name                   string
		testFunc               func(t *testing.T) *domain.Organization
		orgIdentifierCondition database.Condition
		err                    error
	}

	tests := []test{
		func() test {
			organizationId := gofakeit.Name()
			return test{
				name: "happy path get using id",
				testFunc: func(t *testing.T) *domain.Organization {
					organizationName := gofakeit.Name()

					org := domain.Organization{
						ID:         organizationId,
						Name:       organizationName,
						InstanceID: instanceId,
						State:      domain.OrgStateInactive,
					}

					// create organization
					err := orgRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					return &org
				},
				orgIdentifierCondition: orgRepo.IDCondition(organizationId),
			}
		}(),
		func() test {
			organizationId := gofakeit.Name()
			return test{
				name: "happy path get using id including domain",
				testFunc: func(t *testing.T) *domain.Organization {
					organizationName := gofakeit.Name()

					org := domain.Organization{
						ID:         organizationId,
						Name:       organizationName,
						InstanceID: instanceId,
						State:      domain.OrgStateActive,
					}

					// create organization
					err := orgRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					d := &domain.AddOrganizationDomain{
						InstanceID: org.InstanceID,
						OrgID:      org.ID,
						Domain:     gofakeit.DomainName(),
						IsVerified: true,
						IsPrimary:  true,
					}
					err = orgDomainRepo.Add(t.Context(), tx, d)
					require.NoError(t, err)

					org.Domains = []*domain.OrganizationDomain{
						{
							InstanceID:     d.InstanceID,
							OrgID:          d.OrgID,
							ValidationType: d.ValidationType,
							Domain:         d.Domain,
							IsPrimary:      d.IsPrimary,
							IsVerified:     d.IsVerified,
							CreatedAt:      d.CreatedAt,
							UpdatedAt:      d.UpdatedAt,
						},
					}

					return &org
				},
				orgIdentifierCondition: orgRepo.IDCondition(organizationId),
			}
		}(),
		func() test {
			organizationName := gofakeit.Name()
			return test{
				name: "happy path get using name",
				testFunc: func(t *testing.T) *domain.Organization {
					organizationId := gofakeit.Name()

					org := domain.Organization{
						ID:         organizationId,
						Name:       organizationName,
						InstanceID: instanceId,
						State:      domain.OrgStateActive,
					}

					// create organization
					err := orgRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					return &org
				},
				orgIdentifierCondition: orgRepo.NameCondition(database.TextOperationEqual, organizationName),
			}
		}(),
		{
			name: "get non existent organization",
			testFunc: func(t *testing.T) *domain.Organization {
				org := domain.Organization{
					ID:   "non existent org",
					Name: "non existent org",
				}
				return &org
			},
			orgIdentifierCondition: orgRepo.NameCondition(database.TextOperationEqual, "non-existent-instance-name"),
			err:                    new(database.NoRowFoundError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var org *domain.Organization
			if tt.testFunc != nil {
				org = tt.testFunc(t)
			}

			// get org values
			returnedOrg, err := orgRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						tt.orgIdentifierCondition,
						orgRepo.InstanceIDCondition(org.InstanceID),
					),
				),
			)
			if tt.err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, returnedOrg.ID, org.ID)
			assert.Equal(t, returnedOrg.Name, org.Name)
			assert.Equal(t, returnedOrg.InstanceID, org.InstanceID)
			assert.Equal(t, returnedOrg.State, org.State)
		})
	}
}

func TestListOrganization(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()
	organizationRepo := repository.OrganizationRepository()

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
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	type test struct {
		name                   string
		testFunc               func(t *testing.T, tx database.QueryExecutor) []*domain.Organization
		conditionClauses       []database.Condition
		noOrganizationReturned bool
	}
	tests := []test{
		{
			name: "happy path single organization no filter",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Organization {
				noOfOrganizations := 1
				organizations := make([]*domain.Organization, noOfOrganizations)
				for i := range noOfOrganizations {

					org := domain.Organization{
						ID:         strconv.Itoa(i),
						Name:       gofakeit.Name(),
						InstanceID: instanceId,
						State:      domain.OrgStateActive,
					}

					// create organization
					err := organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					organizations[i] = &org
				}

				return organizations
			},
		},
		{
			name: "happy path multiple organization no filter",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Organization {
				noOfOrganizations := 5
				organizations := make([]*domain.Organization, noOfOrganizations)
				for i := range noOfOrganizations {

					org := domain.Organization{
						ID:         strconv.Itoa(i),
						Name:       gofakeit.Name(),
						InstanceID: instanceId,
						State:      domain.OrgStateActive,
					}

					// create organization
					err := organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					organizations[i] = &org
				}

				return organizations
			},
		},
		func() test {
			organizationId := gofakeit.Name()
			return test{
				name: "organization filter on id",
				testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Organization {
					// create organization
					// this org is created as an additional org which should NOT
					// be returned in the results of this test case
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: instanceId,
						State:      domain.OrgStateActive,
					}
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					noOfOrganizations := 1
					organizations := make([]*domain.Organization, noOfOrganizations)
					for i := range noOfOrganizations {

						org := domain.Organization{
							ID:         organizationId,
							Name:       gofakeit.Name(),
							InstanceID: instanceId,
							State:      domain.OrgStateActive,
						}

						// create organization
						err := organizationRepo.Create(t.Context(), tx, &org)
						require.NoError(t, err)

						organizations[i] = &org
					}

					return organizations
				},
				conditionClauses: []database.Condition{
					organizationRepo.InstanceIDCondition(instanceId),
					organizationRepo.IDCondition(organizationId),
				},
			}
		}(),
		{
			name: "multiple organization filter on state",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Organization {
				// create organization
				// this org is created as an additional org which should NOT
				// be returned in the results of this test case
				org := domain.Organization{
					ID:         gofakeit.Name(),
					Name:       gofakeit.Name(),
					InstanceID: instanceId,
					State:      domain.OrgStateActive,
				}
				err = organizationRepo.Create(t.Context(), tx, &org)
				require.NoError(t, err)

				noOfOrganizations := 5
				organizations := make([]*domain.Organization, noOfOrganizations)
				for i := range noOfOrganizations {

					org := domain.Organization{
						ID:         strconv.Itoa(i),
						Name:       gofakeit.Name(),
						InstanceID: instanceId,
						State:      domain.OrgStateInactive,
					}

					// create organization
					err := organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					organizations[i] = &org
				}

				return organizations
			},
			conditionClauses: []database.Condition{
				organizationRepo.InstanceIDCondition(instanceId),
				organizationRepo.StateCondition(domain.OrgStateInactive),
			},
		},
		func() test {
			instanceId_2 := gofakeit.Name()
			return test{
				name: "multiple organization filter on instance",
				testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Organization {
					// create instance 1
					instanceId_1 := gofakeit.Name()
					instance := domain.Instance{
						ID:              instanceId_1,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					err = instanceRepo.Create(t.Context(), tx, &instance)
					assert.Nil(t, err)

					// create organization
					// this org is created as an additional org which should NOT
					// be returned in the results of this test case
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       gofakeit.Name(),
						InstanceID: instanceId_1,
						State:      domain.OrgStateActive,
					}
					err = organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// create instance 2
					instance_2 := domain.Instance{
						ID:              instanceId_2,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}
					err = instanceRepo.Create(t.Context(), tx, &instance_2)
					assert.Nil(t, err)

					noOfOrganizations := 5
					organizations := make([]*domain.Organization, noOfOrganizations)
					for i := range noOfOrganizations {

						org := domain.Organization{
							ID:         strconv.Itoa(i),
							Name:       gofakeit.Name(),
							InstanceID: instanceId_2,
							State:      domain.OrgStateActive,
						}

						// create organization
						err := organizationRepo.Create(t.Context(), tx, &org)
						require.NoError(t, err)

						organizations[i] = &org
					}

					return organizations
				},
				conditionClauses: []database.Condition{organizationRepo.InstanceIDCondition(instanceId_2)},
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
			organizations := tt.testFunc(t, savepoint)

			condition := organizationRepo.InstanceIDCondition(instanceId)
			if len(tt.conditionClauses) > 0 {
				condition = database.And(tt.conditionClauses...)
			}

			// check organization values
			returnedOrgs, err := organizationRepo.List(t.Context(), tx,
				database.WithCondition(condition),
				database.WithOrderByAscending(organizationRepo.CreatedAtColumn(), organizationRepo.IDColumn()),
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
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(t.Context())
	})

	instanceRepo := repository.InstanceRepository()
	organizationRepo := repository.OrganizationRepository()

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
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	type test struct {
		name                   string
		testFunc               func(t *testing.T)
		orgIdentifierCondition database.Condition
		noOfDeletedRows        int64
	}
	tests := []test{
		func() test {
			organizationId := gofakeit.Name()
			var noOfOrganizations int64 = 1
			return test{
				name: "happy path delete organization filter id",
				testFunc: func(t *testing.T) {
					for range noOfOrganizations {
						org := domain.Organization{
							ID:         organizationId,
							Name:       gofakeit.Name(),
							InstanceID: instanceId,
							State:      domain.OrgStateActive,
						}

						// create organization
						err := organizationRepo.Create(t.Context(), tx, &org)
						require.NoError(t, err)

					}
				},
				orgIdentifierCondition: organizationRepo.IDCondition(organizationId),
				noOfDeletedRows:        noOfOrganizations,
			}
		}(),
		func() test {
			organizationName := gofakeit.Name()
			var noOfOrganizations int64 = 1
			return test{
				name: "happy path delete organization filter name",
				testFunc: func(t *testing.T) {
					for range noOfOrganizations {
						org := domain.Organization{
							ID:         gofakeit.Name(),
							Name:       organizationName,
							InstanceID: instanceId,
							State:      domain.OrgStateActive,
						}

						// create organization
						err := organizationRepo.Create(t.Context(), tx, &org)
						require.NoError(t, err)

					}
				},
				orgIdentifierCondition: organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
				noOfDeletedRows:        noOfOrganizations,
			}
		}(),
		func() test {
			non_existent_organization_name := gofakeit.Name()
			return test{
				name:                   "delete non existent organization",
				orgIdentifierCondition: organizationRepo.NameCondition(database.TextOperationEqual, non_existent_organization_name),
			}
		}(),
		func() test {
			organizationName := gofakeit.Name()
			return test{
				name: "deleted already deleted organization",
				testFunc: func(t *testing.T) {
					org := domain.Organization{
						ID:         gofakeit.Name(),
						Name:       organizationName,
						InstanceID: instanceId,
						State:      domain.OrgStateActive,
					}

					// create organization
					err := organizationRepo.Create(t.Context(), tx, &org)
					require.NoError(t, err)

					// delete organization
					affectedRows, err := organizationRepo.Delete(t.Context(), tx,
						database.And(
							organizationRepo.InstanceIDCondition(org.InstanceID),
							organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
						),
					)
					assert.Equal(t, int64(1), affectedRows)
					require.NoError(t, err)
				},
				orgIdentifierCondition: organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
				// this test should return 0 affected rows as the org was already deleted
				noOfDeletedRows: 0,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.testFunc != nil {
				tt.testFunc(t)
			}

			// delete organization
			noOfDeletedRows, err := organizationRepo.Delete(t.Context(), tx,
				database.And(
					organizationRepo.InstanceIDCondition(instanceId),
					tt.orgIdentifierCondition,
				),
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check organization was deleted
			organization, err := organizationRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						tt.orgIdentifierCondition,
						organizationRepo.InstanceIDCondition(instanceId),
					),
				),
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, organization)
		})
	}
}

func TestGetOrganizationWithSubResources(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(t.Context())
	})

	instanceRepo := repository.InstanceRepository()
	orgRepo := repository.OrganizationRepository()

	// create instance
	instanceId := gofakeit.Name()
	err = instanceRepo.Create(t.Context(), tx, &domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	})
	require.NoError(t, err)

	// create organization
	org := domain.Organization{
		ID:         "1",
		Name:       "org-name",
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	}
	err = orgRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	err = orgRepo.Create(t.Context(), tx, &domain.Organization{
		ID:         "without-sub-resources",
		Name:       "org-name-2",
		InstanceID: instanceId,
		State:      domain.OrgStateActive,
	})
	require.NoError(t, err)

	t.Run("domains", func(t *testing.T) {
		domainRepo := repository.OrganizationDomainRepository()

		domain1 := &domain.AddOrganizationDomain{
			InstanceID: org.InstanceID,
			OrgID:      org.ID,
			Domain:     "domain1.com",
			IsVerified: true,
			IsPrimary:  true,
		}
		err = domainRepo.Add(t.Context(), tx, domain1)
		require.NoError(t, err)

		domain2 := &domain.AddOrganizationDomain{
			InstanceID: org.InstanceID,
			OrgID:      org.ID,
			Domain:     "domain2.com",
			IsVerified: false,
			IsPrimary:  false,
		}
		err = domainRepo.Add(t.Context(), tx, domain2)
		require.NoError(t, err)

		t.Run("org by domain", func(t *testing.T) {
			orgRepo := orgRepo.LoadDomains()

			returnedOrg, err := orgRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						orgRepo.InstanceIDCondition(instanceId),
						orgRepo.ExistsDomain(domainRepo.DomainCondition(database.TextOperationEqual, domain1.Domain)),
					),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, org.ID, returnedOrg.ID)
			assert.Len(t, returnedOrg.Domains, 2)
		})

		t.Run("ensure org by domain works without LoadDomains", func(t *testing.T) {
			returnedOrg, err := orgRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						orgRepo.InstanceIDCondition(instanceId),
						orgRepo.ExistsDomain(domainRepo.DomainCondition(database.TextOperationEqual, domain1.Domain)),
					),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, org.ID, returnedOrg.ID)
			assert.Len(t, returnedOrg.Domains, 0)
		})
	})

	t.Run("metadata", func(t *testing.T) {
		metadataRepo := repository.OrganizationMetadataRepository()

		metadata := []*domain.OrganizationMetadata{
			{
				OrgID: org.ID,
				Metadata: domain.Metadata{
					InstanceID: org.InstanceID,
					Key:        "key1",
					Value:      []byte("value1"),
				},
			},
			{
				OrgID: org.ID,
				Metadata: domain.Metadata{
					InstanceID: org.InstanceID,
					Key:        "key2",
					Value:      []byte("value2"),
				},
			},
		}
		err = metadataRepo.Set(t.Context(), tx, metadata...)
		require.NoError(t, err)

		t.Run("org by metadata key", func(t *testing.T) {
			orgRepo := orgRepo.LoadMetadata()

			returnedOrg, err := orgRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						orgRepo.InstanceIDCondition(instanceId),
						orgRepo.ExistsMetadata(metadataRepo.KeyCondition(database.TextOperationEqual, "key1")),
					),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, org.ID, returnedOrg.ID)
			assert.Len(t, returnedOrg.Metadata, 2)
		})

		t.Run("org by metadata value", func(t *testing.T) {
			orgRepo := orgRepo.LoadMetadata()

			returnedOrg, err := orgRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						orgRepo.InstanceIDCondition(instanceId),
						orgRepo.ExistsMetadata(metadataRepo.ValueCondition(database.BytesOperationEqual, []byte("value1"))),
					),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, org.ID, returnedOrg.ID)
			assert.Len(t, returnedOrg.Metadata, 2)
		})

		t.Run("ensure org by metadata key works without LoadMetadata", func(t *testing.T) {
			returnedOrg, err := orgRepo.Get(t.Context(), tx,
				database.WithCondition(
					database.And(
						orgRepo.InstanceIDCondition(instanceId),
						orgRepo.ExistsMetadata(metadataRepo.KeyCondition(database.TextOperationEqual, "key1")),
					),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, org.ID, returnedOrg.ID)
			assert.Len(t, returnedOrg.Metadata, 0)
		})
	})
}

package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestCreateOrganization(t *testing.T) {
	tests := []struct {
		name         string
		testFunc     func() *domain.Organization
		organization domain.Organization
		err          error
	}{
		{
			name: "happy path",
			organization: func() domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				instanceId := gofakeit.Name()
				organization := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.Active,
				}
				return organization
			}(),
		},
		{
			name: "create organization wihtout name",
			organization: func() domain.Organization {
				organizationId := gofakeit.Name()
				// organizationName := gofakeit.Name()
				instanceId := gofakeit.Name()
				organization := domain.Organization{
					ID:         organizationId,
					Name:       "",
					InstanceID: instanceId,
					State:      domain.Active,
				}
				return organization
			}(),
			err: errors.New("organization name not provided"),
		},
		{
			name: "adding same organization twice",
			testFunc: func() *domain.Organization {
				organizationRepo := repository.OrgRepository(pool)
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				instanceId := gofakeit.Name()

				ctx := context.Background()
				inst := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.Active,
				}

				err := organizationRepo.Create(ctx, &inst)
				assert.NoError(t, err)
				return &inst
			},
			err: errors.New("organization id already exists"),
		},
		{
			name: "adding organization with no id",
			organization: func() domain.Organization {
				// organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				instanceId := gofakeit.Name()
				organization := domain.Organization{
					// ID:              organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.Active,
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
					ID:   organizationId,
					Name: organizationName,
					// InstanceID: instanceId,
					State: domain.Active,
				}
				return organization
			}(),
			err: errors.New("instance id not provided"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var organization *domain.Organization
			if tt.testFunc != nil {
				organization = tt.testFunc()
			} else {
				organization = &tt.organization
			}
			organizationRepo := repository.OrgRepository(pool)

			// create organization
			beforeCreate := time.Now()
			err := organizationRepo.Create(ctx, organization)
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check organization values
			organization, err = organizationRepo.Get(ctx,
				organizationRepo.NameCondition(database.TextOperationEqual, organization.Name),
			)
			assert.NoError(t, err)

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
	organizationRepo := repository.OrgRepository(pool)
	tests := []struct {
		name         string
		testFunc     func() *domain.Organization
		update       []database.Change
		rowsAffected int64
	}{
		{
			name: "happy path update name",
			testFunc: func() *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				instanceId := gofakeit.Name()

				ctx := context.Background()
				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.Active,
				}

				// create organization
				err := organizationRepo.Create(ctx, &org)
				assert.NoError(t, err)

				// update with updated value
				org.Name = "new_name"
				return &org
			},
			update:       []database.Change{organizationRepo.SetName("new_name")},
			rowsAffected: 1,
		},
		{
			name: "happy path change state",
			testFunc: func() *domain.Organization {
				organizationId := gofakeit.Name()
				organizationName := gofakeit.Name()
				instanceId := gofakeit.Name()

				ctx := context.Background()
				org := domain.Organization{
					ID:         organizationId,
					Name:       organizationName,
					InstanceID: instanceId,
					State:      domain.Active,
				}

				// create organization
				err := organizationRepo.Create(ctx, &org)
				assert.NoError(t, err)

				// update with updated value
				org.State = domain.Inactive
				return &org
			},
			update:       []database.Change{organizationRepo.SetState(domain.Inactive)},
			rowsAffected: 1,
		},
		// {
		// 	name: "update non existent organization",
		// 	testFunc: func() *domain.Organization {
		// 		organizationId := gofakeit.Name()

		// 		org := domain.Organization{
		// 			ID: organizationId,
		// 		}
		// 		return &org
		// 	},
		// 	rowsAffected: 0,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeUpdate := time.Now()

			ctx := context.Background()
			organizationRepo := repository.OrgRepository(pool)

			createdOrg := tt.testFunc()

			// update org
			rowsAffected, err := organizationRepo.Update(ctx,
				organizationRepo.IDCondition(createdOrg.ID),
				tt.update...,
			)
			afterUpdate := time.Now()
			assert.NoError(t, err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			if rowsAffected == 0 {
				return
			}

			// check organization values
			organization, err := organizationRepo.Get(ctx,
				organizationRepo.IDCondition(createdOrg.ID),
			)
			assert.NoError(t, err)

			assert.Equal(t, createdOrg.ID, organization.ID)
			assert.Equal(t, createdOrg.Name, organization.Name)
			assert.Equal(t, createdOrg.State, organization.State)
			assert.WithinRange(t, organization.UpdatedAt, beforeUpdate, afterUpdate)
			assert.Nil(t, organization.DeletedAt)
		})
	}
}

// func TestGetOrganization(t *testing.T) {
// 	organizationRepo := repository.OrgRepository(pool)
// 	type test struct {
// 		name                   string
// 		testFunc               func() *domain.Organization
// 		conditionClauses       []database.Condition
// 		noOrganizationReturned bool
// 	}

// 	tests := []test{
// 		func() test {
// 			organizationId := gofakeit.Name()
// 			return test{
// 				name: "happy path get using id",
// 				testFunc: func() *domain.Organization {
// 					organizationName := gofakeit.Name()

// 					ctx := context.Background()
// 					inst := domain.Organization{
// 						ID:              organizationId,
// 						Name:            organizationName,
// 						DefaultOrgID:    "defaultOrgId",
// 						IAMProjectID:    "iamProject",
// 						ConsoleClientID: "consoleCLient",
// 						ConsoleAppID:    "consoleApp",
// 						DefaultLanguage: "defaultLanguage",
// 					}

// 					// create organization
// 					err := organizationRepo.Create(ctx, &inst)
// 					assert.NoError(t, err)
// 					return &inst
// 				},
// 				conditionClauses: []database.Condition{organizationRepo.IDCondition(organizationId)},
// 			}
// 		}(),
// 		func() test {
// 			organizationName := gofakeit.Name()
// 			return test{
// 				name: "happy path get using name",
// 				testFunc: func() *domain.Organization {
// 					organizationId := gofakeit.Name()

// 					ctx := context.Background()
// 					inst := domain.Organization{
// 						ID:              organizationId,
// 						Name:            organizationName,
// 						DefaultOrgID:    "defaultOrgId",
// 						IAMProjectID:    "iamProject",
// 						ConsoleClientID: "consoleCLient",
// 						ConsoleAppID:    "consoleApp",
// 						DefaultLanguage: "defaultLanguage",
// 					}

// 					// create organization
// 					err := organizationRepo.Create(ctx, &inst)
// 					assert.NoError(t, err)
// 					return &inst
// 				},
// 				conditionClauses: []database.Condition{organizationRepo.NameCondition(database.TextOperationEqual, organizationName)},
// 			}
// 		}(),
// 		{
// 			name: "get non existent organization",
// 			testFunc: func() *domain.Organization {
// 				organizationId := gofakeit.Name()

// 				inst := domain.Organization{
// 					ID: organizationId,
// 				}
// 				return &inst
// 			},
// 			conditionClauses:       []database.Condition{organizationRepo.NameCondition(database.TextOperationEqual, "non-existent-organization-name")},
// 			noOrganizationReturned: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx := context.Background()
// 			organizationRepo := repository.OrgRepository(pool)

// 			var organization *domain.Organization
// 			if tt.testFunc != nil {
// 				organization = tt.testFunc()
// 			}

// 			// check organization values
// 			returnedOrganization, err := organizationRepo.Get(ctx,
// 				tt.conditionClauses...,
// 			)
// 			assert.NoError(t, err)
// 			if tt.noOrganizationReturned {
// 				assert.Nil(t, returnedOrganization)
// 				return
// 			}

// 			assert.Equal(t, returnedOrganization.ID, organization.ID)
// 			assert.Equal(t, returnedOrganization.Name, organization.Name)
// 			assert.Equal(t, returnedOrganization.DefaultOrgID, organization.DefaultOrgID)
// 			assert.Equal(t, returnedOrganization.IAMProjectID, organization.IAMProjectID)
// 			assert.Equal(t, returnedOrganization.ConsoleClientID, organization.ConsoleClientID)
// 			assert.Equal(t, returnedOrganization.ConsoleAppID, organization.ConsoleAppID)
// 			assert.Equal(t, returnedOrganization.DefaultLanguage, organization.DefaultLanguage)
// 			assert.NoError(t, err)
// 		})
// 	}
// }

// func TestListOrganization(t *testing.T) {
// 	type test struct {
// 		name                   string
// 		testFunc               func() ([]*domain.Organization, database.PoolTest, func())
// 		conditionClauses       []database.Condition
// 		noOrganizationReturned bool
// 	}
// 	tests := []test{
// 		{
// 			name: "happy path single organization no filter",
// 			testFunc: func() ([]*domain.Organization, database.PoolTest, func()) {
// 				ctx := context.Background()
// 				// create new db to make sure no organizations exist
// 				pool, stop, err := newEmbeededDB()
// 				assert.NoError(t, err)

// 				organizationRepo := repository.OrgRepository(pool)
// 				noOfOrganizations := 1
// 				organizations := make([]*domain.Organization, noOfOrganizations)
// 				for i := range noOfOrganizations {

// 					organizationId := gofakeit.Name()
// 					organizationName := gofakeit.Name()

// 					inst := domain.Organization{
// 						ID:              organizationId,
// 						Name:            organizationName,
// 						DefaultOrgID:    "defaultOrgId",
// 						IAMProjectID:    "iamProject",
// 						ConsoleClientID: "consoleCLient",
// 						ConsoleAppID:    "consoleApp",
// 						DefaultLanguage: "defaultLanguage",
// 					}

// 					// create organization
// 					err := organizationRepo.Create(ctx, &inst)
// 					assert.NoError(t, err)

// 					organizations[i] = &inst
// 				}

// 				return organizations, pool, stop
// 			},
// 		},
// 		{
// 			name: "happy path multiple organization no filter",
// 			testFunc: func() ([]*domain.Organization, database.PoolTest, func()) {
// 				ctx := context.Background()
// 				// create new db to make sure no organizations exist
// 				pool, stop, err := newEmbeededDB()
// 				assert.NoError(t, err)

// 				organizationRepo := repository.OrgRepository(pool)
// 				noOfOrganizations := 5
// 				organizations := make([]*domain.Organization, noOfOrganizations)
// 				for i := range noOfOrganizations {

// 					organizationId := gofakeit.Name()
// 					organizationName := gofakeit.Name()

// 					inst := domain.Organization{
// 						ID:              organizationId,
// 						Name:            organizationName,
// 						DefaultOrgID:    "defaultOrgId",
// 						IAMProjectID:    "iamProject",
// 						ConsoleClientID: "consoleCLient",
// 						ConsoleAppID:    "consoleApp",
// 						DefaultLanguage: "defaultLanguage",
// 					}

// 					// create organization
// 					err := organizationRepo.Create(ctx, &inst)
// 					assert.NoError(t, err)

// 					organizations[i] = &inst
// 				}

// 				return organizations, pool, stop
// 			},
// 		},
// 		func() test {
// 			organizationRepo := repository.OrgRepository(pool)
// 			organizationId := gofakeit.Name()
// 			return test{
// 				name: "organization filter on id",
// 				testFunc: func() ([]*domain.Organization, database.PoolTest, func()) {
// 					ctx := context.Background()

// 					noOfOrganizations := 1
// 					organizations := make([]*domain.Organization, noOfOrganizations)
// 					for i := range noOfOrganizations {

// 						organizationName := gofakeit.Name()

// 						inst := domain.Organization{
// 							ID:              organizationId,
// 							Name:            organizationName,
// 							DefaultOrgID:    "defaultOrgId",
// 							IAMProjectID:    "iamProject",
// 							ConsoleClientID: "consoleCLient",
// 							ConsoleAppID:    "consoleApp",
// 							DefaultLanguage: "defaultLanguage",
// 						}

// 						// create organization
// 						err := organizationRepo.Create(ctx, &inst)
// 						assert.NoError(t, err)

// 						organizations[i] = &inst
// 					}

// 					return organizations, nil, nil
// 				},
// 				conditionClauses: []database.Condition{organizationRepo.IDCondition(organizationId)},
// 			}
// 		}(),
// 		func() test {
// 			organizationRepo := repository.OrgRepository(pool)
// 			organizationName := gofakeit.Name()
// 			return test{
// 				name: "multiple organization filter on name",
// 				testFunc: func() ([]*domain.Organization, database.PoolTest, func()) {
// 					ctx := context.Background()

// 					noOfOrganizations := 5
// 					organizations := make([]*domain.Organization, noOfOrganizations)
// 					for i := range noOfOrganizations {

// 						organizationId := gofakeit.Name()

// 						inst := domain.Organization{
// 							ID:              organizationId,
// 							Name:            organizationName,
// 							DefaultOrgID:    "defaultOrgId",
// 							IAMProjectID:    "iamProject",
// 							ConsoleClientID: "consoleCLient",
// 							ConsoleAppID:    "consoleApp",
// 							DefaultLanguage: "defaultLanguage",
// 						}

// 						// create organization
// 						err := organizationRepo.Create(ctx, &inst)
// 						assert.NoError(t, err)

// 						organizations[i] = &inst
// 					}

// 					return organizations, nil, nil
// 				},
// 				conditionClauses: []database.Condition{organizationRepo.NameCondition(database.TextOperationEqual, organizationName)},
// 			}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx := context.Background()

// 			var organizations []*domain.Organization

// 			pool := pool
// 			if tt.testFunc != nil {
// 				var stop func()
// 				var pool_ database.PoolTest
// 				organizations, pool_, stop = tt.testFunc()
// 				if pool_ != nil {
// 					pool = pool_
// 					defer stop()
// 				}
// 			}
// 			organizationRepo := repository.OrgRepository(pool)

// 			// check organization values
// 			returnedOrganizations, err := organizationRepo.List(ctx,
// 				tt.conditionClauses...,
// 			)
// 			assert.NoError(t, err)
// 			if tt.noOrganizationReturned {
// 				assert.Nil(t, returnedOrganizations)
// 				return
// 			}

// 			assert.Equal(t, len(organizations), len(returnedOrganizations))
// 			for i, organization := range organizations {
// 				assert.Equal(t, returnedOrganizations[i].ID, organization.ID)
// 				assert.Equal(t, returnedOrganizations[i].Name, organization.Name)
// 				assert.Equal(t, returnedOrganizations[i].DefaultOrgID, organization.DefaultOrgID)
// 				assert.Equal(t, returnedOrganizations[i].IAMProjectID, organization.IAMProjectID)
// 				assert.Equal(t, returnedOrganizations[i].ConsoleClientID, organization.ConsoleClientID)
// 				assert.Equal(t, returnedOrganizations[i].ConsoleAppID, organization.ConsoleAppID)
// 				assert.Equal(t, returnedOrganizations[i].DefaultLanguage, organization.DefaultLanguage)
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestDeleteOrganization(t *testing.T) {
// 	type test struct {
// 		name             string
// 		testFunc         func()
// 		conditionClauses database.Condition
// 	}
// 	tests := []test{
// 		func() test {
// 			organizationRepo := repository.OrgRepository(pool)
// 			organizationId := gofakeit.Name()
// 			return test{
// 				name: "happy path delete single organization filter id",
// 				testFunc: func() {
// 					ctx := context.Background()

// 					noOfOrganizations := 1
// 					organizations := make([]*domain.Organization, noOfOrganizations)
// 					for i := range noOfOrganizations {

// 						organizationName := gofakeit.Name()

// 						inst := domain.Organization{
// 							ID:              organizationId,
// 							Name:            organizationName,
// 							DefaultOrgID:    "defaultOrgId",
// 							IAMProjectID:    "iamProject",
// 							ConsoleClientID: "consoleCLient",
// 							ConsoleAppID:    "consoleApp",
// 							DefaultLanguage: "defaultLanguage",
// 						}

// 						// create organization
// 						err := organizationRepo.Create(ctx, &inst)
// 						assert.NoError(t, err)

// 						organizations[i] = &inst
// 					}
// 				},
// 				conditionClauses: organizationRepo.IDCondition(organizationId),
// 			}
// 		}(),
// 		func() test {
// 			organizationRepo := repository.OrgRepository(pool)
// 			organizationName := gofakeit.Name()
// 			return test{
// 				name: "happy path delete single organization filter name",
// 				testFunc: func() {
// 					ctx := context.Background()

// 					noOfOrganizations := 1
// 					organizations := make([]*domain.Organization, noOfOrganizations)
// 					for i := range noOfOrganizations {

// 						organizationId := gofakeit.Name()

// 						inst := domain.Organization{
// 							ID:              organizationId,
// 							Name:            organizationName,
// 							DefaultOrgID:    "defaultOrgId",
// 							IAMProjectID:    "iamProject",
// 							ConsoleClientID: "consoleCLient",
// 							ConsoleAppID:    "consoleApp",
// 							DefaultLanguage: "defaultLanguage",
// 						}

// 						// create organization
// 						err := organizationRepo.Create(ctx, &inst)
// 						assert.NoError(t, err)

// 						organizations[i] = &inst
// 					}
// 				},
// 				conditionClauses: organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
// 			}
// 		}(),
// 		func() test {
// 			organizationRepo := repository.OrgRepository(pool)
// 			non_existent_organization_name := gofakeit.Name()
// 			return test{
// 				name:             "delete non existent organization",
// 				conditionClauses: organizationRepo.NameCondition(database.TextOperationEqual, non_existent_organization_name),
// 			}
// 		}(),
// 		func() test {
// 			organizationRepo := repository.OrgRepository(pool)
// 			organizationName := gofakeit.Name()
// 			return test{
// 				name: "multiple organization filter on name",
// 				testFunc: func() {
// 					ctx := context.Background()

// 					noOfOrganizations := 5
// 					organizations := make([]*domain.Organization, noOfOrganizations)
// 					for i := range noOfOrganizations {

// 						organizationId := gofakeit.Name()

// 						inst := domain.Organization{
// 							ID:              organizationId,
// 							Name:            organizationName,
// 							DefaultOrgID:    "defaultOrgId",
// 							IAMProjectID:    "iamProject",
// 							ConsoleClientID: "consoleCLient",
// 							ConsoleAppID:    "consoleApp",
// 							DefaultLanguage: "defaultLanguage",
// 						}

// 						// create organization
// 						err := organizationRepo.Create(ctx, &inst)
// 						assert.NoError(t, err)

// 						organizations[i] = &inst
// 					}
// 				},
// 				conditionClauses: organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
// 			}
// 		}(),
// 		func() test {
// 			organizationRepo := repository.OrgRepository(pool)
// 			organizationName := gofakeit.Name()
// 			return test{
// 				name: "deleted already deleted organization",
// 				testFunc: func() {
// 					ctx := context.Background()

// 					noOfOrganizations := 1
// 					organizations := make([]*domain.Organization, noOfOrganizations)
// 					for i := range noOfOrganizations {

// 						organizationId := gofakeit.Name()

// 						inst := domain.Organization{
// 							ID:              organizationId,
// 							Name:            organizationName,
// 							DefaultOrgID:    "defaultOrgId",
// 							IAMProjectID:    "iamProject",
// 							ConsoleClientID: "consoleCLient",
// 							ConsoleAppID:    "consoleApp",
// 							DefaultLanguage: "defaultLanguage",
// 						}

// 						// create organization
// 						err := organizationRepo.Create(ctx, &inst)
// 						assert.NoError(t, err)

// 						organizations[i] = &inst
// 					}

// 					// delete organization
// 					err := organizationRepo.Delete(ctx,
// 						organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
// 					)
// 					assert.NoError(t, err)
// 				},
// 				conditionClauses: organizationRepo.NameCondition(database.TextOperationEqual, organizationName),
// 			}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx := context.Background()
// 			organizationRepo := repository.OrgRepository(pool)

// 			if tt.testFunc != nil {
// 				tt.testFunc()
// 			}

// 			// delete organization
// 			err := organizationRepo.Delete(ctx,
// 				tt.conditionClauses,
// 			)
// 			assert.NoError(t, err)

// 			// check organization was deleted
// 			organization, err := organizationRepo.Get(ctx,
// 				tt.conditionClauses,
// 			)
// 			assert.NoError(t, err)
// 			assert.Nil(t, organization)
// 		})
// 	}
// }

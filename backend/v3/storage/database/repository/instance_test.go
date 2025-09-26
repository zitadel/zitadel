package repository_test

import (
	"context"
	"strconv"
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

func TestCreateInstance(t *testing.T) {
	beforeCreate := time.Now()
	tx, err := pool.Begin(context.Background(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil {
			t.Log("error during rollback:", err)
		}
	}()
	instanceRepo := repository.InstanceRepository()

	tests := []struct {
		name     string
		testFunc func(t *testing.T, tx database.QueryExecutor) *domain.Instance
		instance domain.Instance
		err      error
	}{
		{
			name: "happy path",
			instance: func() domain.Instance {
				instanceId := gofakeit.Name()
				instanceName := gofakeit.Name()
				instance := domain.Instance{
					ID:              instanceId,
					Name:            instanceName,
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}
				return instance
			}(),
		},
		{
			name: "create instance without name",
			instance: func() domain.Instance {
				instanceId := gofakeit.Name()
				// instanceName := gofakeit.Name()
				instance := domain.Instance{
					ID:              instanceId,
					Name:            "",
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}
				return instance
			}(),
			err: new(database.CheckError),
		},
		{
			name: "adding same instance twice",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Instance {
				inst := domain.Instance{
					ID:              gofakeit.UUID(),
					Name:            gofakeit.Name(),
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}

				err := instanceRepo.Create(t.Context(), tx, &inst)
				require.NoError(t, err)

				// change the name to make sure same only the id clashes
				inst.Name = gofakeit.Name()
				require.NoError(t, err)
				return &inst
			},
			err: new(database.UniqueError),
		},
		func() struct {
			name     string
			testFunc func(t *testing.T, tx database.QueryExecutor) *domain.Instance
			instance domain.Instance
			err      error
		} {
			instanceId := gofakeit.Name()
			instanceName := gofakeit.Name()
			return struct {
				name     string
				testFunc func(t *testing.T, tx database.QueryExecutor) *domain.Instance
				instance domain.Instance
				err      error
			}{
				name: "adding instance with same name twice",
				testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Instance {
					inst := domain.Instance{
						ID:              gofakeit.Name(),
						Name:            instanceName,
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					err := instanceRepo.Create(t.Context(), tx, &inst)
					require.NoError(t, err)

					// change the id
					inst.ID = instanceId
					inst.CreatedAt = time.Time{}
					inst.UpdatedAt = time.Time{}
					return &inst
				},
				instance: domain.Instance{
					ID:              instanceId,
					Name:            instanceName,
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				},
				// two instances can have the sane name
				err: nil,
			}
		}(),
		{
			name: "adding instance with no id",
			instance: func() domain.Instance {
				instance := domain.Instance{
					Name:            gofakeit.Name(),
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}
				return instance
			}(),
			err: new(database.CheckError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = savepoint.Rollback(t.Context())
				if err != nil {
					t.Log("error during rollback:", err)
				}
			}()

			var instance *domain.Instance
			if tt.testFunc != nil {
				instance = tt.testFunc(t, savepoint)
			} else {
				instance = &tt.instance
			}

			// create instance

			err = instanceRepo.Create(t.Context(), tx, instance)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check instance values
			instance, err = instanceRepo.Get(t.Context(), tx,
				database.WithCondition(
					instanceRepo.IDCondition(instance.ID),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.instance.ID, instance.ID)
			assert.Equal(t, tt.instance.Name, instance.Name)
			assert.Equal(t, tt.instance.DefaultOrgID, instance.DefaultOrgID)
			assert.Equal(t, tt.instance.IAMProjectID, instance.IAMProjectID)
			assert.Equal(t, tt.instance.ConsoleClientID, instance.ConsoleClientID)
			assert.Equal(t, tt.instance.ConsoleAppID, instance.ConsoleAppID)
			assert.Equal(t, tt.instance.DefaultLanguage, instance.DefaultLanguage)
			assert.WithinRange(t, instance.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, instance.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

func TestUpdateInstance(t *testing.T) {
	beforeUpdate := time.Now()
	tx, err := pool.Begin(context.Background(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil {
			t.Log("error during rollback:", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()

	tests := []struct {
		name         string
		testFunc     func(t *testing.T, tx database.QueryExecutor) *domain.Instance
		rowsAffected int64
		getErr       error
	}{
		{
			name: "happy path",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Instance {
				inst := domain.Instance{
					ID:              gofakeit.UUID(),
					Name:            gofakeit.Name(),
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}

				// create instance
				err := instanceRepo.Create(t.Context(), tx, &inst)
				require.NoError(t, err)
				return &inst
			},
			rowsAffected: 1,
		},
		{
			name: "update deleted instance",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Instance {
				inst := domain.Instance{
					ID:              gofakeit.UUID(),
					Name:            gofakeit.Name(),
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}

				// create instance
				err := instanceRepo.Create(t.Context(), tx, &inst)
				require.NoError(t, err)

				// delete instance
				affectedRows, err := instanceRepo.Delete(t.Context(), tx,
					inst.ID,
				)
				require.NoError(t, err)
				assert.Equal(t, int64(1), affectedRows)

				return &inst
			},
			rowsAffected: 0,
		},
		{
			name: "update non existent instance",
			testFunc: func(t *testing.T, tx database.QueryExecutor) *domain.Instance {
				inst := domain.Instance{
					ID: gofakeit.UUID(),
				}
				return &inst
			},
			rowsAffected: 0,
			getErr:       new(database.NoRowFoundError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := tt.testFunc(t, tx)

			// update name
			newName := "new_" + instance.Name
			rowsAffected, err := instanceRepo.Update(t.Context(), tx,
				instance.ID,
				instanceRepo.SetName(newName),
			)
			afterUpdate := time.Now()
			require.NoError(t, err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			if rowsAffected == 0 {
				return
			}

			// check instance values
			instance, err = instanceRepo.Get(t.Context(), tx,
				database.WithCondition(
					instanceRepo.IDCondition(instance.ID),
				),
			)
			require.Equal(t, tt.getErr, err)

			assert.Equal(t, newName, instance.Name)
			assert.WithinRange(t, instance.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

func TestGetInstance(t *testing.T) {
	tx, err := pool.Begin(context.Background(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil {
			t.Log("error during rollback:", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()
	domainRepo := repository.InstanceDomainRepository()

	type test struct {
		name     string
		testFunc func(t *testing.T) *domain.Instance
		err      error
	}
	tests := []test{
		func() test {
			return test{
				name: "happy path get using id",
				testFunc: func(t *testing.T) *domain.Instance {
					inst := domain.Instance{
						ID:              gofakeit.UUID(),
						Name:            gofakeit.BeerName(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					// create instance
					err := instanceRepo.Create(t.Context(), tx, &inst)
					require.NoError(t, err)
					return &inst
				},
			}
		}(),
		{
			name: "happy path including domains",
			testFunc: func(t *testing.T) *domain.Instance {
				inst := domain.Instance{
					ID:              gofakeit.NewCrypto().UUID(),
					Name:            gofakeit.BeerName(),
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}

				// create instance
				err := instanceRepo.Create(t.Context(), tx, &inst)
				require.NoError(t, err)

				d := &domain.AddInstanceDomain{
					InstanceID:  inst.ID,
					Domain:      gofakeit.DomainName(),
					IsPrimary:   gu.Ptr(true),
					IsGenerated: gu.Ptr(false),
					Type:        domain.DomainTypeCustom,
				}
				err = domainRepo.Add(t.Context(), tx, d)
				require.NoError(t, err)

				inst.Domains = append(inst.Domains, &domain.InstanceDomain{
					InstanceID:  d.InstanceID,
					Domain:      d.Domain,
					IsPrimary:   d.IsPrimary,
					IsGenerated: d.IsGenerated,
					Type:        d.Type,
					CreatedAt:   d.CreatedAt,
					UpdatedAt:   d.UpdatedAt,
				})

				return &inst
			},
		},
		{
			name: "get non existent instance",
			testFunc: func(t *testing.T) *domain.Instance {
				inst := domain.Instance{
					ID: "get non existent instance",
				}
				return &inst
			},
			err: new(database.NoRowFoundError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var instance *domain.Instance
			if tt.testFunc != nil {
				instance = tt.testFunc(t)
			}

			// check instance values
			returnedInstance, err := instanceRepo.Get(t.Context(), tx,
				database.WithCondition(
					instanceRepo.IDCondition(instance.ID),
				),
			)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}

			if instance.ID == "get non existent instance" {
				assert.Nil(t, returnedInstance)
				return
			}

			assert.Equal(t, returnedInstance.ID, instance.ID)
			assert.Equal(t, returnedInstance.Name, instance.Name)
			assert.Equal(t, returnedInstance.DefaultOrgID, instance.DefaultOrgID)
			assert.Equal(t, returnedInstance.IAMProjectID, instance.IAMProjectID)
			assert.Equal(t, returnedInstance.ConsoleClientID, instance.ConsoleClientID)
			assert.Equal(t, returnedInstance.ConsoleAppID, instance.ConsoleAppID)
			assert.Equal(t, returnedInstance.DefaultLanguage, instance.DefaultLanguage)
		})
	}
}

func TestListInstance(t *testing.T) {
	tx, err := pool.Begin(context.Background(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil {
			t.Log("error during rollback:", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()

	type test struct {
		name               string
		testFunc           func(t *testing.T, tx database.QueryExecutor) []*domain.Instance
		conditionClauses   []database.Condition
		noInstanceReturned bool
	}
	tests := []test{
		{
			name: "happy path single instance no filter",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Instance {
				noOfInstances := 1
				instances := make([]*domain.Instance, noOfInstances)
				for i := range noOfInstances {

					inst := domain.Instance{
						ID:              strconv.Itoa(i),
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					// create instance
					err := instanceRepo.Create(t.Context(), tx, &inst)
					require.NoError(t, err)

					instances[i] = &inst
				}

				return instances
			},
		},
		{
			name: "happy path multiple instance no filter",
			testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Instance {
				noOfInstances := 5
				instances := make([]*domain.Instance, noOfInstances)
				for i := range noOfInstances {

					inst := domain.Instance{
						ID:              strconv.Itoa(i),
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					// create instance
					err := instanceRepo.Create(t.Context(), tx, &inst)
					require.NoError(t, err)

					instances[i] = &inst
				}

				return instances
			},
		},
		func() test {
			instanceID := gofakeit.BeerName()
			return test{
				name: "instance filter on id",
				testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Instance {
					noOfInstances := 1
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						inst := domain.Instance{
							ID:              instanceID,
							Name:            gofakeit.Name(),
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(t.Context(), tx, &inst)
						require.NoError(t, err)

						instances[i] = &inst
					}

					return instances
				},
				conditionClauses: []database.Condition{instanceRepo.IDCondition(instanceID)},
			}
		}(),
		func() test {
			instanceName := gofakeit.BeerName()
			return test{
				name: "multiple instance filter on name",
				testFunc: func(t *testing.T, tx database.QueryExecutor) []*domain.Instance {
					noOfInstances := 5
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						inst := domain.Instance{
							ID:              strconv.Itoa(i),
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(t.Context(), tx, &inst)
						require.NoError(t, err)

						instances[i] = &inst
					}

					return instances
				},
				conditionClauses: []database.Condition{instanceRepo.NameCondition(database.TextOperationEqual, instanceName)},
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
					t.Log("error during rollback:", err)
				}
			}()
			instances := tt.testFunc(t, savepoint)

			var condition database.Condition
			if len(tt.conditionClauses) > 0 {
				condition = database.And(tt.conditionClauses...)
			}

			// check instance values
			returnedInstances, err := instanceRepo.List(t.Context(), tx,
				database.WithCondition(condition),
				database.WithOrderByAscending(instanceRepo.IDColumn()),
			)
			require.NoError(t, err)
			if tt.noInstanceReturned {
				assert.Len(t, returnedInstances, 0)
				return
			}

			assert.Equal(t, len(instances), len(returnedInstances))
			for i, instance := range instances {
				assert.Equal(t, returnedInstances[i].ID, instance.ID)
				assert.Equal(t, returnedInstances[i].Name, instance.Name)
				assert.Equal(t, returnedInstances[i].DefaultOrgID, instance.DefaultOrgID)
				assert.Equal(t, returnedInstances[i].IAMProjectID, instance.IAMProjectID)
				assert.Equal(t, returnedInstances[i].ConsoleClientID, instance.ConsoleClientID)
				assert.Equal(t, returnedInstances[i].ConsoleAppID, instance.ConsoleAppID)
				assert.Equal(t, returnedInstances[i].DefaultLanguage, instance.DefaultLanguage)
			}
		})
	}
}

func TestDeleteInstance(t *testing.T) {
	tx, err := pool.Begin(context.Background(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil {
			t.Log("error during rollback:", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()

	type test struct {
		name            string
		testFunc        func(t *testing.T, tx database.QueryExecutor)
		instanceID      string
		noOfDeletedRows int64
	}
	tests := []test{
		func() test {
			instanceID := gofakeit.NewCrypto().UUID()
			return test{
				name: "happy path delete single instance filter id",
				testFunc: func(t *testing.T, tx database.QueryExecutor) {
					inst := domain.Instance{
						ID:              instanceID,
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					// create instance
					err := instanceRepo.Create(t.Context(), tx, &inst)
					require.NoError(t, err)
				},
				instanceID:      instanceID,
				noOfDeletedRows: 1,
			}
		}(),
		func() test {
			non_existent_instance_name := gofakeit.Name()
			return test{
				name:       "delete non existent instance",
				instanceID: non_existent_instance_name,
			}
		}(),
		func() test {
			instanceID := gofakeit.Name()
			return test{
				name: "deleted already deleted instance",
				testFunc: func(t *testing.T, tx database.QueryExecutor) {

					inst := domain.Instance{
						ID:              instanceID,
						Name:            gofakeit.BeerName(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					// create instance
					err := instanceRepo.Create(t.Context(), tx, &inst)
					require.NoError(t, err)

					// delete instance
					affectedRows, err := instanceRepo.Delete(t.Context(), tx,
						inst.ID,
					)
					require.NoError(t, err)
					assert.Equal(t, int64(1), affectedRows)
				},
				instanceID: instanceID,
				// this test should return 0 affected rows as the instance was already deleted
				noOfDeletedRows: 0,
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
					t.Log("error during rollback:", err)
				}
			}()

			if tt.testFunc != nil {
				tt.testFunc(t, savepoint)
			}

			// delete instance
			noOfDeletedRows, err := instanceRepo.Delete(t.Context(), savepoint, tt.instanceID)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check instance was deleted
			instance, err := instanceRepo.Get(t.Context(), savepoint,
				database.WithCondition(
					instanceRepo.IDCondition(tt.instanceID),
				),
			)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, instance)
		})
	}
}

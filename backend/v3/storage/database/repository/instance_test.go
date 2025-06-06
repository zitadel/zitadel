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

func TestCreateInstance(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() *domain.Instance
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
			name: "create instance wihtout name",
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
			err: errors.New("instnace name not provided"),
		},
		{
			name: "adding same instance twice",
			testFunc: func() *domain.Instance {
				instanceRepo := repository.InstanceRepository(pool)
				instanceId := gofakeit.Name()
				instanceName := gofakeit.Name()

				ctx := context.Background()
				inst := domain.Instance{
					ID:              instanceId,
					Name:            instanceName,
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}

				err := instanceRepo.Create(ctx, &inst)
				assert.NoError(t, err)
				return &inst
			},
			err: errors.New("instnace id already exists"),
		},
		{
			name: "adding instance with no id",
			instance: func() domain.Instance {
				// instanceId := gofakeit.Name()
				instanceName := gofakeit.Name()
				instance := domain.Instance{
					// ID:              instanceId,
					Name:            instanceName,
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}
				return instance
			}(),
			err: errors.New("instnace id not provided"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var instance *domain.Instance
			if tt.testFunc != nil {
				instance = tt.testFunc()
			} else {
				instance = &tt.instance
			}
			instanceRepo := repository.InstanceRepository(pool)

			// create instance
			beforeCreate := time.Now()
			err := instanceRepo.Create(ctx, instance)
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check instance values
			instance, err = instanceRepo.Get(ctx,
				instanceRepo.NameCondition(database.TextOperationEqual, instance.Name),
			)
			assert.Equal(t, tt.instance.ID, instance.ID)
			assert.Equal(t, tt.instance.Name, instance.Name)
			assert.Equal(t, tt.instance.DefaultOrgID, instance.DefaultOrgID)
			assert.Equal(t, tt.instance.IAMProjectID, instance.IAMProjectID)
			assert.Equal(t, tt.instance.ConsoleClientID, instance.ConsoleClientID)
			assert.Equal(t, tt.instance.ConsoleAppID, instance.ConsoleAppID)
			assert.Equal(t, tt.instance.DefaultLanguage, instance.DefaultLanguage)
			assert.WithinRange(t, instance.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, instance.UpdatedAt, beforeCreate, afterCreate)
			assert.Nil(t, instance.DeletedAt)
			assert.NoError(t, err)
		})
	}
}

func TestUpdateInstance(t *testing.T) {
	tests := []struct {
		name         string
		testFunc     func() *domain.Instance
		rowsAffected int64
	}{
		{
			name: "happy path",
			testFunc: func() *domain.Instance {
				instanceRepo := repository.InstanceRepository(pool)
				instanceId := gofakeit.Name()
				instanceName := gofakeit.Name()

				ctx := context.Background()
				inst := domain.Instance{
					ID:              instanceId,
					Name:            instanceName,
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}

				// create instance
				err := instanceRepo.Create(ctx, &inst)
				assert.NoError(t, err)
				return &inst
			},
			rowsAffected: 1,
		},
		{
			name: "update non existent instance",
			testFunc: func() *domain.Instance {
				instanceId := gofakeit.Name()

				inst := domain.Instance{
					ID: instanceId,
				}
				return &inst
			},
			rowsAffected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeUpdate := time.Now()

			ctx := context.Background()
			instanceRepo := repository.InstanceRepository(pool)

			instance := tt.testFunc()

			// update name
			newName := "new_" + instance.Name
			rowsAffected, err := instanceRepo.Update(ctx,
				instanceRepo.IDCondition(instance.ID),
				instanceRepo.SetName(newName),
			)
			afterUpdate := time.Now()
			assert.NoError(t, err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			if rowsAffected == 0 {
				return
			}

			// check instance values
			instance, err = instanceRepo.Get(ctx,
				instanceRepo.IDCondition(instance.ID),
			)
			assert.NoError(t, err)

			assert.Equal(t, newName, instance.Name)
			assert.WithinRange(t, instance.UpdatedAt, beforeUpdate, afterUpdate)
			assert.Nil(t, instance.DeletedAt)
		})
	}
}

func TestGetInstance(t *testing.T) {
	tests := []struct {
		name               string
		testFunc           func() *domain.Instance
		noInstanceReturned bool
	}{
		{
			name: "happy path",
			testFunc: func() *domain.Instance {
				instanceRepo := repository.InstanceRepository(pool)
				instanceId := gofakeit.Name()
				instanceName := gofakeit.Name()

				ctx := context.Background()
				inst := domain.Instance{
					ID:              instanceId,
					Name:            instanceName,
					DefaultOrgID:    "defaultOrgId",
					IAMProjectID:    "iamProject",
					ConsoleClientID: "consoleCLient",
					ConsoleAppID:    "consoleApp",
					DefaultLanguage: "defaultLanguage",
				}

				// create instance
				err := instanceRepo.Create(ctx, &inst)
				assert.NoError(t, err)
				return &inst
			},
		},
		{
			name: "get non existent instance",
			testFunc: func() *domain.Instance {
				instanceId := gofakeit.Name()

				inst := domain.Instance{
					ID: instanceId,
				}
				return &inst
			},
			noInstanceReturned: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			instanceRepo := repository.InstanceRepository(pool)

			var instance *domain.Instance
			if tt.testFunc != nil {
				instance = tt.testFunc()
			}

			// check instance values
			returnedInstance, err := instanceRepo.Get(ctx,
				instanceRepo.IDCondition(instance.ID),
			)
			assert.NoError(t, err)
			if tt.noInstanceReturned {
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
			assert.NoError(t, err)
		})
	}
}

func TestListInstance(t *testing.T) {
	type test struct {
		name               string
		testFunc           func() ([]*domain.Instance, database.PoolTest, func())
		conditionClauses   []database.Condition
		noInstanceReturned bool
	}
	tests := []test{
		{
			name: "happy path single instance",
			testFunc: func() ([]*domain.Instance, database.PoolTest, func()) {
				ctx := context.Background()
				pool, stop, err := newEmbeededDB()
				assert.NoError(t, err)

				instanceRepo := repository.InstanceRepository(pool)
				noOfInstances := 1
				instances := make([]*domain.Instance, noOfInstances)
				for i := range noOfInstances {

					instanceId := gofakeit.Name()
					instanceName := gofakeit.Name()

					inst := domain.Instance{
						ID:              instanceId,
						Name:            instanceName,
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					// create instance
					err := instanceRepo.Create(ctx, &inst)
					assert.NoError(t, err)

					instances[i] = &inst
				}

				return instances, pool, stop
			},
		},
		{
			name: "happy path multiple instance",
			testFunc: func() ([]*domain.Instance, database.PoolTest, func()) {
				ctx := context.Background()
				pool, stop, err := newEmbeededDB()
				assert.NoError(t, err)

				instanceRepo := repository.InstanceRepository(pool)
				noOfInstances := 5
				instances := make([]*domain.Instance, noOfInstances)
				for i := range noOfInstances {

					instanceId := gofakeit.Name()
					instanceName := gofakeit.Name()

					inst := domain.Instance{
						ID:              instanceId,
						Name:            instanceName,
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					// create instance
					err := instanceRepo.Create(ctx, &inst)
					assert.NoError(t, err)

					instances[i] = &inst
				}

				return instances, pool, stop
			},
		},
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceId := gofakeit.Name()
			return test{
				name: "instance filter on id",
				testFunc: func() ([]*domain.Instance, database.PoolTest, func()) {
					ctx := context.Background()

					noOfInstances := 1
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						instanceName := gofakeit.Name()

						inst := domain.Instance{
							ID:              instanceId,
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						assert.NoError(t, err)

						instances[i] = &inst
					}

					return instances, nil, nil
				},
				conditionClauses: []database.Condition{instanceRepo.IDCondition(instanceId)},
			}
		}(),
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceName := gofakeit.Name()
			return test{
				name: "multiple instance filter on name",
				testFunc: func() ([]*domain.Instance, database.PoolTest, func()) {
					ctx := context.Background()

					noOfInstances := 5
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						instanceId := gofakeit.Name()

						inst := domain.Instance{
							ID:              instanceId,
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						assert.NoError(t, err)

						instances[i] = &inst
					}

					return instances, nil, nil
				},
				conditionClauses: []database.Condition{instanceRepo.NameCondition(database.TextOperationEqual, instanceName)},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var instances []*domain.Instance

			pool := pool
			if tt.testFunc != nil {
				var stop func()
				var pool_ database.PoolTest
				instances, pool_, stop = tt.testFunc()
				if pool_ != nil {
					pool = pool_
					defer stop()
				}
			}
			instanceRepo := repository.InstanceRepository(pool)

			// check instance values
			returnedInstances, err := instanceRepo.List(ctx,
				tt.conditionClauses...,
			)
			assert.NoError(t, err)
			if tt.noInstanceReturned {
				assert.Nil(t, returnedInstances)
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
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteInstance(t *testing.T) {
	type test struct {
		name             string
		testFunc         func()
		conditionClauses database.Condition
	}
	tests := []test{
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceName := gofakeit.Name()
			return test{
				name: "happy path delete single instance",
				testFunc: func() {
					ctx := context.Background()

					noOfInstances := 1
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						instanceId := gofakeit.Name()

						inst := domain.Instance{
							ID:              instanceId,
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						assert.NoError(t, err)

						instances[i] = &inst
					}
				},
				conditionClauses: instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
			}
		}(),
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			non_existent_instance_name := gofakeit.Name()
			return test{
				name:             "delete non existent instance",
				conditionClauses: instanceRepo.NameCondition(database.TextOperationEqual, non_existent_instance_name),
			}
		}(),
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceName := gofakeit.Name()
			return test{
				name: "multiple instance filter on name",
				testFunc: func() {
					ctx := context.Background()

					noOfInstances := 5
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						instanceId := gofakeit.Name()

						inst := domain.Instance{
							ID:              instanceId,
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						assert.NoError(t, err)

						instances[i] = &inst
					}
				},
				conditionClauses: instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
			}
		}(),
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceName := gofakeit.Name()
			return test{
				name: "deleted already deleted instance",
				testFunc: func() {
					ctx := context.Background()

					noOfInstances := 1
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						instanceId := gofakeit.Name()

						inst := domain.Instance{
							ID:              instanceId,
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						assert.NoError(t, err)

						instances[i] = &inst
					}

					// delete instance
					err := instanceRepo.Delete(ctx,
						instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
					)
					assert.NoError(t, err)
				},
				conditionClauses: instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			instanceRepo := repository.InstanceRepository(pool)

			if tt.testFunc != nil {
				tt.testFunc()
			}

			// delete instance
			err := instanceRepo.Delete(ctx,
				tt.conditionClauses,
			)
			assert.NoError(t, err)

			// check instance was deleted
			instance, err := instanceRepo.Get(ctx,
				tt.conditionClauses,
			)
			assert.NoError(t, err)
			assert.Nil(t, instance)
		})
	}
}

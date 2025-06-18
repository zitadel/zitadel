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

func TestCreateInstance(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.Instance
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
			err: errors.New("instance name not provided"),
		},
		{
			name: "adding same instance twice",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
				instanceRepo := repository.InstanceRepository(pool)
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

				err := instanceRepo.Create(ctx, &inst)
				require.NoError(t, err)
				return &inst
			},
			err: errors.New("instance id already exists"),
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
			err: errors.New("instance id not provided"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var instance *domain.Instance
			if tt.testFunc != nil {
				instance = tt.testFunc(ctx, t)
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
			assert.Nil(t, instance.DeletedAt)
		})
	}
}

func TestUpdateInstance(t *testing.T) {
	tests := []struct {
		name         string
		testFunc     func(ctx context.Context, t *testing.T) *domain.Instance
		rowsAffected int64
	}{
		{
			name: "happy path",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
				instanceRepo := repository.InstanceRepository(pool)
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
				require.NoError(t, err)
				return &inst
			},
			rowsAffected: 1,
		},
		{
			name: "update deleted instance",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
				instanceRepo := repository.InstanceRepository(pool)
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
				require.NoError(t, err)

				// delete instance
				err = instanceRepo.Delete(ctx,
					instanceRepo.IDCondition(inst.ID),
				)
				require.NoError(t, err)

				return &inst
			},
			rowsAffected: 0,
		},
		{
			name: "update non existent instance",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
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
			ctx := context.Background()
			instanceRepo := repository.InstanceRepository(pool)

			instance := tt.testFunc(ctx, t)

			beforeUpdate := time.Now()
			// update name
			newName := "new_" + instance.Name
			rowsAffected, err := instanceRepo.Update(ctx,
				instanceRepo.IDCondition(instance.ID),
				instanceRepo.SetName(newName),
			)
			afterUpdate := time.Now()
			require.NoError(t, err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			if rowsAffected == 0 {
				return
			}

			// check instance values
			instance, err = instanceRepo.Get(ctx,
				instanceRepo.IDCondition(instance.ID),
			)
			require.NoError(t, err)

			assert.Equal(t, newName, instance.Name)
			assert.WithinRange(t, instance.UpdatedAt, beforeUpdate, afterUpdate)
			assert.Nil(t, instance.DeletedAt)
		})
	}
}

func TestGetInstance(t *testing.T) {
	instanceRepo := repository.InstanceRepository(pool)
	type test struct {
		name             string
		testFunc         func(ctx context.Context, t *testing.T) *domain.Instance
		conditionClauses []database.Condition
	}

	tests := []test{
		func() test {
			instanceId := gofakeit.Name()
			return test{
				name: "happy path get using id",
				testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
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
					require.NoError(t, err)
					return &inst
				},
				conditionClauses: []database.Condition{instanceRepo.IDCondition(instanceId)},
			}
		}(),
		func() test {
			instanceName := gofakeit.Name()
			return test{
				name: "happy path get using name",
				testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
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
					require.NoError(t, err)
					return &inst
				},
				conditionClauses: []database.Condition{instanceRepo.NameCondition(database.TextOperationEqual, instanceName)},
			}
		}(),
		{
			name: "get non existent instance",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
				instanceId := gofakeit.Name()

				_ = domain.Instance{
					ID: instanceId,
				}
				return nil
			},
			conditionClauses: []database.Condition{instanceRepo.NameCondition(database.TextOperationEqual, "non-existent-instance-name")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			instanceRepo := repository.InstanceRepository(pool)

			var instance *domain.Instance
			if tt.testFunc != nil {
				instance = tt.testFunc(ctx, t)
			}

			// check instance values
			returnedInstance, err := instanceRepo.Get(ctx,
				tt.conditionClauses...,
			)
			require.NoError(t, err)
			if instance == nil {
				assert.Nil(t, instance, returnedInstance)
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
	ctx := context.Background()
	pool, stop, err := newEmbeddedDB(ctx)
	require.NoError(t, err)
	defer stop()

	type test struct {
		name               string
		testFunc           func(ctx context.Context, t *testing.T) []*domain.Instance
		conditionClauses   []database.Condition
		noInstanceReturned bool
	}
	tests := []test{
		{
			name: "happy path single instance no filter",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Instance {
				instanceRepo := repository.InstanceRepository(pool)
				noOfInstances := 1
				instances := make([]*domain.Instance, noOfInstances)
				for i := range noOfInstances {

					inst := domain.Instance{
						ID:              gofakeit.Name(),
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					// create instance
					err := instanceRepo.Create(ctx, &inst)
					require.NoError(t, err)

					instances[i] = &inst
				}

				return instances
			},
		},
		{
			name: "happy path multiple instance no filter",
			testFunc: func(ctx context.Context, t *testing.T) []*domain.Instance {
				instanceRepo := repository.InstanceRepository(pool)
				noOfInstances := 5
				instances := make([]*domain.Instance, noOfInstances)
				for i := range noOfInstances {

					inst := domain.Instance{
						ID:              gofakeit.Name(),
						Name:            gofakeit.Name(),
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					// create instance
					err := instanceRepo.Create(ctx, &inst)
					require.NoError(t, err)

					instances[i] = &inst
				}

				return instances
			},
		},
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceId := gofakeit.Name()
			return test{
				name: "instance filter on id",
				testFunc: func(ctx context.Context, t *testing.T) []*domain.Instance {
					noOfInstances := 1
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						inst := domain.Instance{
							ID:              instanceId,
							Name:            gofakeit.Name(),
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						require.NoError(t, err)

						instances[i] = &inst
					}

					return instances
				},
				conditionClauses: []database.Condition{instanceRepo.IDCondition(instanceId)},
			}
		}(),
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceName := gofakeit.Name()
			return test{
				name: "multiple instance filter on name",
				testFunc: func(ctx context.Context, t *testing.T) []*domain.Instance {
					noOfInstances := 5
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						inst := domain.Instance{
							ID:              gofakeit.Name(),
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
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
			t.Cleanup(func() {
				_, err := pool.Exec(ctx, "DELETE FROM zitadel.instances")
				require.NoError(t, err)
			})

			instances := tt.testFunc(ctx, t)

			instanceRepo := repository.InstanceRepository(pool)

			// check instance values
			returnedInstances, err := instanceRepo.List(ctx,
				tt.conditionClauses...,
			)
			require.NoError(t, err)
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
			}
		})
	}
}

func TestDeleteInstance(t *testing.T) {
	type test struct {
		name             string
		testFunc         func(ctx context.Context, t *testing.T)
		conditionClauses database.Condition
	}
	tests := []test{
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceId := gofakeit.Name()
			return test{
				name: "happy path delete single instance filter id",
				testFunc: func(ctx context.Context, t *testing.T) {
					noOfInstances := 1
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						inst := domain.Instance{
							ID:              instanceId,
							Name:            gofakeit.Name(),
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						require.NoError(t, err)

						instances[i] = &inst
					}
				},
				conditionClauses: instanceRepo.IDCondition(instanceId),
			}
		}(),
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceName := gofakeit.Name()
			return test{
				name: "happy path delete single instance filter name",
				testFunc: func(ctx context.Context, t *testing.T) {
					noOfInstances := 1
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						inst := domain.Instance{
							ID:              gofakeit.Name(),
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						require.NoError(t, err)

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
				testFunc: func(ctx context.Context, t *testing.T) {
					noOfInstances := 5
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						inst := domain.Instance{
							ID:              gofakeit.Name(),
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						require.NoError(t, err)

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
				testFunc: func(ctx context.Context, t *testing.T) {
					noOfInstances := 1
					instances := make([]*domain.Instance, noOfInstances)
					for i := range noOfInstances {

						inst := domain.Instance{
							ID:              gofakeit.Name(),
							Name:            instanceName,
							DefaultOrgID:    "defaultOrgId",
							IAMProjectID:    "iamProject",
							ConsoleClientID: "consoleCLient",
							ConsoleAppID:    "consoleApp",
							DefaultLanguage: "defaultLanguage",
						}

						// create instance
						err := instanceRepo.Create(ctx, &inst)
						require.NoError(t, err)

						instances[i] = &inst
					}

					// delete instance
					err := instanceRepo.Delete(ctx,
						instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
					)
					require.NoError(t, err)
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
				tt.testFunc(ctx, t)
			}

			// delete instance
			err := instanceRepo.Delete(ctx,
				tt.conditionClauses,
			)
			require.NoError(t, err)

			// check instance was deleted
			instance, err := instanceRepo.Get(ctx,
				tt.conditionClauses,
			)
			require.NoError(t, err)
			assert.Nil(t, instance)
		})
	}
}

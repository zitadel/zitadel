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
			err: new(database.CheckErr),
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
				// change the name to make sure same only the id clashes
				inst.Name = gofakeit.Name()
				require.NoError(t, err)
				return &inst
			},
			err: new(database.UniqueErr),
		},
		func() struct {
			name     string
			testFunc func(ctx context.Context, t *testing.T) *domain.Instance
			instance domain.Instance
			err      error
		} {
			instanceId := gofakeit.Name()
			instanceName := gofakeit.Name()
			return struct {
				name     string
				testFunc func(ctx context.Context, t *testing.T) *domain.Instance
				instance domain.Instance
				err      error
			}{
				name: "adding instance with same name twice",
				testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
					instanceRepo := repository.InstanceRepository(pool)

					inst := domain.Instance{
						ID:              gofakeit.Name(),
						Name:            instanceName,
						DefaultOrgID:    "defaultOrgId",
						IAMProjectID:    "iamProject",
						ConsoleClientID: "consoleCLient",
						ConsoleAppID:    "consoleApp",
						DefaultLanguage: "defaultLanguage",
					}

					err := instanceRepo.Create(ctx, &inst)
					require.NoError(t, err)

					// change the id
					inst.ID = instanceId
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
			err: new(database.CheckErr),
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
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check instance values
			instance, err = instanceRepo.Get(ctx,
				instance.ID,
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
	tests := []struct {
		name         string
		testFunc     func(ctx context.Context, t *testing.T) *domain.Instance
		rowsAffected int64
		getErr       error
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
				affectedRows, err := instanceRepo.Delete(ctx,
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
			testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
				instanceId := gofakeit.Name()

				inst := domain.Instance{
					ID: instanceId,
				}
				return &inst
			},
			rowsAffected: 0,
			getErr:       new(database.ErrNoRowFound),
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
			instance, err = instanceRepo.Get(ctx,
				instance.ID,
			)
			require.Equal(t, tt.getErr, err)

			assert.Equal(t, newName, instance.Name)
			assert.WithinRange(t, instance.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

func TestGetInstance(t *testing.T) {
	instanceRepo := repository.InstanceRepository(pool)
	type test struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T) *domain.Instance
		err      error
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
			}
		}(),
		{
			name: "get non existent instance",
			testFunc: func(ctx context.Context, t *testing.T) *domain.Instance {
				inst := domain.Instance{
					ID: "get non existent instance",
				}
				return &inst
			},
			err: new(database.ErrNoRowFound),
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
				instance.ID,
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
		name            string
		testFunc        func(ctx context.Context, t *testing.T)
		instanceID      string
		noOfDeletedRows int64
	}
	tests := []test{
		func() test {
			instanceRepo := repository.InstanceRepository(pool)
			instanceId := gofakeit.Name()
			var noOfInstances int64 = 1
			return test{
				name: "happy path delete single instance filter id",
				testFunc: func(ctx context.Context, t *testing.T) {
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
				instanceID:      instanceId,
				noOfDeletedRows: noOfInstances,
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
					affectedRows, err := instanceRepo.Delete(ctx,
						instances[0].ID,
					)
					require.NoError(t, err)
					assert.Equal(t, int64(1), affectedRows)
				},
				instanceID: instanceName,
				// this test should return 0 affected rows as the instance was already deleted
				noOfDeletedRows: 0,
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
			noOfDeletedRows, err := instanceRepo.Delete(ctx,
				tt.instanceID,
			)
			require.NoError(t, err)
			assert.Equal(t, noOfDeletedRows, tt.noOfDeletedRows)

			// check instance was deleted
			instance, err := instanceRepo.Get(ctx,
				tt.instanceID,
			)
			require.ErrorIs(t, err, new(database.ErrNoRowFound))
			assert.Nil(t, instance)
		})
	}
}

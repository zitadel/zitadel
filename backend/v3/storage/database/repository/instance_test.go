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
	instanceRepo := repository.InstanceRepository(pool)
	instanceId := gofakeit.Name()
	instanceName := gofakeit.Name()

	ctx := context.Background()
	inst := domain.Instance{
		ID:              instanceId,
		Name:            instanceName,
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientId: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}

	beforeCreate := time.Now()
	err := instanceRepo.Create(ctx, &inst)
	require.NoError(t, err)
	afterCreate := time.Now()

	instance, err := instanceRepo.Get(ctx,
		instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
	)
	require.Equal(t, inst.ID, instance.ID)
	require.Equal(t, inst.Name, instance.Name)
	require.Equal(t, inst.DefaultOrgID, instance.DefaultOrgID)
	require.Equal(t, inst.IAMProjectID, instance.IAMProjectID)
	require.Equal(t, inst.ConsoleClientId, instance.ConsoleClientId)
	require.Equal(t, inst.ConsoleAppID, instance.ConsoleAppID)
	require.Equal(t, inst.DefaultLanguage, instance.DefaultLanguage)
	assert.WithinRange(t, instance.CreatedAt, beforeCreate, afterCreate)
	assert.WithinRange(t, instance.UpdatedAt, beforeCreate, afterCreate)
	require.Nil(t, instance.DeletedAt)
	require.NoError(t, err)
}

func TestUpdateNameInstance(t *testing.T) {
	instanceRepo := repository.InstanceRepository(pool)
	instanceId := gofakeit.Name()
	instanceName := gofakeit.Name()

	ctx := context.Background()
	inst := domain.Instance{
		ID:              instanceId,
		Name:            instanceName,
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientId: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}

	err := instanceRepo.Create(ctx, &inst)
	require.NoError(t, err)

	_, err = instanceRepo.Get(ctx,
		instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
	)
	require.NoError(t, err)

	// update name
	err = instanceRepo.Update(ctx,
		instanceRepo.IDCondition(instanceId),
		instanceRepo.SetName("new_name"),
	)
	require.NoError(t, err)

	instance, err := instanceRepo.Get(ctx,
		instanceRepo.IDCondition(instanceId),
	)
	require.NoError(t, err)
	require.Equal(t, "new_name", instance.Name)
}

func TestUpdeDeleteInstance(t *testing.T) {
	instanceRepo := repository.InstanceRepository(pool)
	instanceId := gofakeit.Name()
	instanceName := gofakeit.Name()

	ctx := context.Background()
	inst := domain.Instance{
		ID:              instanceId,
		Name:            instanceName,
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientId: "consoleCLient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}

	err := instanceRepo.Create(ctx, &inst)
	require.NoError(t, err)

	instance, err := instanceRepo.Get(ctx,
		instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
	)
	require.NotNil(t, instance)
	require.NoError(t, err)

	// delete instance
	err = instanceRepo.Delete(ctx,
		instanceRepo.IDCondition(instanceId),
	)
	require.NoError(t, err)

	instance, err = instanceRepo.Get(ctx,
		instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
	)
	require.NoError(t, err)
	require.Nil(t, instance)
}

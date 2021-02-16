package migration

import (
	"errors"
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMigration_Cleanup1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	reason := "test"

	client.EXPECT().DeleteJob(namespace, getJobName(reason)).Times(1).Return(nil)
	cleanupFunc := GetCleanupFunc(monitor, namespace, reason)
	assert.NotNil(t, cleanupFunc)
	assert.NoError(t, cleanupFunc(client))
}

func TestMigration_Cleanup2(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test2"
	reason := "test2"

	client.EXPECT().DeleteJob(namespace, getJobName(reason)).Times(1).Return(nil)
	cleanupFunc := GetCleanupFunc(monitor, namespace, reason)
	assert.NotNil(t, cleanupFunc)
	assert.NoError(t, cleanupFunc(client))
}

func TestMigration_CleanupFailure1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	reason := "test"

	client.EXPECT().DeleteJob(namespace, getJobName(reason)).Times(1).Return(errors.New("failed"))
	cleanupFunc := GetCleanupFunc(monitor, namespace, reason)
	assert.NotNil(t, cleanupFunc)
	assert.Error(t, cleanupFunc(client))
}

func TestMigration_CleanupFailure2(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test2"
	reason := "test2"

	client.EXPECT().DeleteJob(namespace, getJobName(reason)).Times(1).Return(errors.New("failed"))
	cleanupFunc := GetCleanupFunc(monitor, namespace, reason)
	assert.NotNil(t, cleanupFunc)
	assert.Error(t, cleanupFunc(client))
}

package backup

import (
	"fmt"
	"github.com/caos/zitadel/operator/database/kinds/backups/core"
	"testing"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBackup_Cleanup1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}
	name := "test"
	namespace := "testNs"

	cleanupFunc := GetCleanupFunc(monitor, namespace, name)
	client.EXPECT().WaitUntilJobCompleted(namespace, core.GetBackupJobName(name), timeout).Times(1).Return(nil)
	client.EXPECT().DeleteJob(namespace, core.GetBackupJobName(name)).Times(1)
	assert.NoError(t, cleanupFunc(client))

	client.EXPECT().WaitUntilJobCompleted(namespace, core.GetBackupJobName(name), timeout).Times(1).Return(fmt.Errorf("fail"))
	assert.Error(t, cleanupFunc(client))
}

func TestBackup_Cleanup2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}
	name := "test2"
	namespace := "testNs2"

	cleanupFunc := GetCleanupFunc(monitor, namespace, name)
	client.EXPECT().WaitUntilJobCompleted(namespace, core.GetBackupJobName(name), timeout).Times(1).Return(nil)
	client.EXPECT().DeleteJob(namespace, core.GetBackupJobName(name)).Times(1)
	assert.NoError(t, cleanupFunc(client))

	client.EXPECT().WaitUntilJobCompleted(namespace, core.GetBackupJobName(name), timeout).Times(1).Return(fmt.Errorf("fail"))
	assert.Error(t, cleanupFunc(client))
}

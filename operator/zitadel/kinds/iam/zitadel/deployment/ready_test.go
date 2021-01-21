package deployment

import (
	"errors"
	"testing"

	"github.com/caos/orbos/pkg/labels/mocklabels"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeployment_Ready1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"

	client.EXPECT().WaitUntilDeploymentReady(namespace, mocklabels.NameVal, true, true, timeout).Times(1).Return(nil)
	readyFunc := GetReadyFunc(monitor, namespace, mocklabels.Name)
	assert.NotNil(t, readyFunc)
	assert.NoError(t, readyFunc(client))
}

func TestDeployment_Ready2(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test2"

	client.EXPECT().WaitUntilDeploymentReady(namespace, mocklabels.NameVal, true, true, timeout).Times(1).Return(nil)
	readyFunc := GetReadyFunc(monitor, namespace, mocklabels.Name)
	assert.NotNil(t, readyFunc)
	assert.NoError(t, readyFunc(client))
}

func TestDeployment_ReadyFailure1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"

	client.EXPECT().WaitUntilDeploymentReady(namespace, mocklabels.NameVal, true, true, timeout).Times(1).Return(errors.New("fail"))
	readyFunc := GetReadyFunc(monitor, namespace, mocklabels.Name)
	assert.NotNil(t, readyFunc)
	assert.Error(t, readyFunc(client))
}

func TestDeployment_ReadyFailure2(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test2"

	client.EXPECT().WaitUntilDeploymentReady(namespace, mocklabels.NameVal, true, true, timeout).Times(1).Return(errors.New("fail"))
	readyFunc := GetReadyFunc(monitor, namespace, mocklabels.Name)
	assert.NotNil(t, readyFunc)
	assert.Error(t, readyFunc(client))
}

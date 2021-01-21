package setup

import (
	"errors"
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetup_Done1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	reason := "test"

	client.EXPECT().WaitUntilJobCompleted(namespace, getJobName(reason), timeout).Times(1).Return(nil)
	cleanupFunc := GetDoneFunc(monitor, namespace, reason)
	assert.NotNil(t, cleanupFunc)
	assert.NoError(t, cleanupFunc(client))
}

func TestSetup_Done2(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test2"
	reason := "test2"

	client.EXPECT().WaitUntilJobCompleted(namespace, getJobName(reason), timeout).Times(1).Return(nil)
	cleanupFunc := GetDoneFunc(monitor, namespace, reason)
	assert.NotNil(t, cleanupFunc)
	assert.NoError(t, cleanupFunc(client))
}

func TestSetup_DoneFailure1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	reason := "test"

	client.EXPECT().WaitUntilJobCompleted(namespace, getJobName(reason), timeout).Times(1).Return(errors.New("failed"))
	doneFunc := GetDoneFunc(monitor, namespace, reason)
	assert.NotNil(t, doneFunc)
	assert.Error(t, doneFunc(client))
}

func TestSetup_DoneFailure2(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test2"
	reason := "test2"

	client.EXPECT().WaitUntilJobCompleted(namespace, getJobName(reason), timeout).Times(1).Return(errors.New("failed"))
	doneFunc := GetDoneFunc(monitor, namespace, reason)
	assert.NotNil(t, doneFunc)
	assert.Error(t, doneFunc(client))
}

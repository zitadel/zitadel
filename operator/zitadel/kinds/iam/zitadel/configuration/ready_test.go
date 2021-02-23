package configuration

import (
	"errors"
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfiguration_Ready1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	secretName := "testSecret"
	secretVarsName := "testVars"
	secretPasswordName := "testPasswords"
	cmName := "testCM"
	consoleCMName := "testConsole"

	client.EXPECT().WaitForConfigMap(namespace, cmName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForConfigMap(namespace, consoleCMName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretVarsName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretPasswordName, timeout).MinTimes(1).MaxTimes(1).Return(nil)

	readyFunc := GetReadyFunc(monitor, namespace, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName)
	assert.NotNil(t, readyFunc)
	assert.NoError(t, readyFunc(client))
}

func TestConfiguration_Ready2(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test2"
	secretName := "testSecret2"
	secretVarsName := "testVars2"
	secretPasswordName := "testPasswords2"
	cmName := "testCM2"
	consoleCMName := "testConsole2"

	client.EXPECT().WaitForConfigMap(namespace, cmName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForConfigMap(namespace, consoleCMName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretVarsName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretPasswordName, timeout).MinTimes(1).MaxTimes(1).Return(nil)

	readyFunc := GetReadyFunc(monitor, namespace, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName)
	assert.NotNil(t, readyFunc)
	assert.NoError(t, readyFunc(client))
}

func TestConfiguration_ReadyFail1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	secretName := "testSecret"
	secretVarsName := "testVars"
	secretPasswordName := "testPasswords"
	cmName := "testCM"
	consoleCMName := "testConsole"

	client.EXPECT().WaitForConfigMap(namespace, cmName, timeout).MinTimes(1).MaxTimes(1).Return(errors.New("fail"))
	client.EXPECT().WaitForConfigMap(namespace, consoleCMName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretVarsName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretPasswordName, timeout).MinTimes(1).MaxTimes(1).Return(nil)

	readyFunc := GetReadyFunc(monitor, namespace, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName)
	assert.NotNil(t, readyFunc)
	assert.Error(t, readyFunc(client))
}

func TestConfiguration_ReadyFail2(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	secretName := "testSecret"
	secretVarsName := "testVars"
	secretPasswordName := "testPasswords"
	cmName := "testCM"
	consoleCMName := "testConsole"

	client.EXPECT().WaitForConfigMap(namespace, cmName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForConfigMap(namespace, consoleCMName, timeout).MinTimes(1).MaxTimes(1).Return(errors.New("fail"))
	client.EXPECT().WaitForSecret(namespace, secretName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretVarsName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretPasswordName, timeout).MinTimes(1).MaxTimes(1).Return(nil)

	readyFunc := GetReadyFunc(monitor, namespace, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName)
	assert.NotNil(t, readyFunc)
	assert.Error(t, readyFunc(client))
}

func TestConfiguration_ReadyFail3(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	secretName := "testSecret"
	secretVarsName := "testVars"
	secretPasswordName := "testPasswords"
	cmName := "testCM"
	consoleCMName := "testConsole"

	client.EXPECT().WaitForConfigMap(namespace, cmName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForConfigMap(namespace, consoleCMName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretName, timeout).MinTimes(1).MaxTimes(1).Return(errors.New("fail"))
	client.EXPECT().WaitForSecret(namespace, secretVarsName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretPasswordName, timeout).MinTimes(1).MaxTimes(1).Return(nil)

	readyFunc := GetReadyFunc(monitor, namespace, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName)
	assert.NotNil(t, readyFunc)
	assert.Error(t, readyFunc(client))
}

func TestConfiguration_ReadyFail4(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	secretName := "testSecret"
	secretVarsName := "testVars"
	secretPasswordName := "testPasswords"
	cmName := "testCM"
	consoleCMName := "testConsole"

	client.EXPECT().WaitForConfigMap(namespace, cmName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForConfigMap(namespace, consoleCMName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretVarsName, timeout).MinTimes(1).MaxTimes(1).Return(errors.New("fail"))
	client.EXPECT().WaitForSecret(namespace, secretPasswordName, timeout).MinTimes(1).MaxTimes(1).Return(nil)

	readyFunc := GetReadyFunc(monitor, namespace, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName)
	assert.NotNil(t, readyFunc)
	assert.Error(t, readyFunc(client))
}

func TestConfiguration_ReadyFail5(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	secretName := "testSecret"
	secretVarsName := "testVars"
	secretPasswordName := "testPasswords"
	cmName := "testCM"
	consoleCMName := "testConsole"

	client.EXPECT().WaitForConfigMap(namespace, cmName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForConfigMap(namespace, consoleCMName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretVarsName, timeout).MinTimes(1).MaxTimes(1).Return(nil)
	client.EXPECT().WaitForSecret(namespace, secretPasswordName, timeout).MinTimes(1).MaxTimes(1).Return(errors.New("fail"))

	readyFunc := GetReadyFunc(monitor, namespace, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName)
	assert.NotNil(t, readyFunc)
	assert.Error(t, readyFunc(client))
}

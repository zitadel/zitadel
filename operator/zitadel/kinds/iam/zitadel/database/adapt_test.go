package database

import (
	"errors"
	"testing"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	databasemock "github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_Adapt(t *testing.T) {
	dbClient := databasemock.NewMockClient(gomock.NewController(t))
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	host := "host"
	port := "port"
	users := []string{"test"}

	monitor := mntr.Monitor{}
	queried := map[string]interface{}{}

	dbClient.EXPECT().GetConnectionInfo(monitor, k8sClient).Return(host, port, nil)
	dbClient.EXPECT().ListUsers(monitor, k8sClient).Return([]string{"test"}, nil)

	_, err := GetDatabaseInQueried(queried)
	assert.Error(t, err)

	query, err := AdaptFunc(monitor, dbClient)
	assert.NoError(t, err)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)

	assert.NoError(t, ensure(k8sClient))

	current, err := GetDatabaseInQueried(queried)
	assert.NoError(t, err)
	assert.Equal(t, host, current.Host)
	assert.Equal(t, port, current.Port)
	assert.Equal(t, users, current.Users)
}

func TestDatabase_Adapt2(t *testing.T) {
	dbClient := databasemock.NewMockClient(gomock.NewController(t))
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	host := "host2"
	port := "port2"
	users := []string{"test2"}

	monitor := mntr.Monitor{}
	queried := map[string]interface{}{}

	dbClient.EXPECT().GetConnectionInfo(monitor, k8sClient).Return(host, port, nil)
	dbClient.EXPECT().ListUsers(monitor, k8sClient).Return(users, nil)

	_, err := GetDatabaseInQueried(queried)
	assert.Error(t, err)

	query, err := AdaptFunc(monitor, dbClient)
	assert.NoError(t, err)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)

	assert.NoError(t, ensure(k8sClient))

	current, err := GetDatabaseInQueried(queried)
	assert.NoError(t, err)
	assert.Equal(t, host, current.Host)
	assert.Equal(t, port, current.Port)
	assert.Equal(t, users, current.Users)
}

func TestDatabase_AdaptFailConnection(t *testing.T) {
	dbClient := databasemock.NewMockClient(gomock.NewController(t))
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	monitor := mntr.Monitor{}
	queried := map[string]interface{}{}

	dbClient.EXPECT().GetConnectionInfo(monitor, k8sClient).Return("", "", errors.New("fail"))
	dbClient.EXPECT().ListUsers(monitor, k8sClient).Return([]string{"test"}, nil)

	_, err := GetDatabaseInQueried(queried)
	assert.Error(t, err)

	query, err := AdaptFunc(monitor, dbClient)
	assert.NoError(t, err)
	ensure, err := query(k8sClient, queried)
	assert.Error(t, err)
	assert.Nil(t, ensure)

	current, err := GetDatabaseInQueried(queried)
	assert.Error(t, err)
	assert.Nil(t, current)
}

func TestDatabase_AdaptFailUsers(t *testing.T) {
	dbClient := databasemock.NewMockClient(gomock.NewController(t))
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	host := "host"
	port := "port"

	monitor := mntr.Monitor{}
	queried := map[string]interface{}{}

	dbClient.EXPECT().GetConnectionInfo(monitor, k8sClient).Return(host, port, nil)
	dbClient.EXPECT().ListUsers(monitor, k8sClient).Return(nil, errors.New("fail"))

	_, err := GetDatabaseInQueried(queried)
	assert.Error(t, err)

	query, err := AdaptFunc(monitor, dbClient)
	assert.NoError(t, err)
	ensure, err := query(k8sClient, queried)
	assert.Error(t, err)
	assert.Nil(t, ensure)

	current, err := GetDatabaseInQueried(queried)
	assert.Error(t, err)
	assert.Nil(t, current)
}

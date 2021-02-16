package users

import (
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	databasemock "github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsers_CreateIfNecessary(t *testing.T) {
	users := []string{}
	monitor := mntr.Monitor{}
	user := "test"
	dbClient := databasemock.NewMockClientInt(gomock.NewController(t))
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	dbClient.EXPECT().AddUser(monitor, user, k8sClient)
	function := createIfNecessary(monitor, user, users, dbClient)
	assert.NotNil(t, function)
	err := function(k8sClient)
	assert.NoError(t, err)

	users = []string{"test"}
	function = createIfNecessary(monitor, user, users, dbClient)
	assert.Nil(t, function)

	user2 := "test2"
	dbClient.EXPECT().AddUser(monitor, user2, k8sClient)
	function = createIfNecessary(monitor, user2, users, dbClient)
	assert.NotNil(t, function)
	err = function(k8sClient)
	assert.NoError(t, err)
}

func TestUsers_DeleteIfNotRequired(t *testing.T) {
	users := []string{}
	monitor := mntr.Monitor{}
	user := "test"
	dbClient := databasemock.NewMockClientInt(gomock.NewController(t))
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	dbClient.EXPECT().DeleteUser(monitor, user, k8sClient)
	function := deleteIfNotRequired(monitor, user, users, dbClient)
	assert.NotNil(t, function)
	err := function(k8sClient)
	assert.NoError(t, err)

	users = []string{"test"}
	function = deleteIfNotRequired(monitor, user, users, dbClient)
	assert.Nil(t, function)

	user2 := "test2"
	dbClient.EXPECT().DeleteUser(monitor, user2, k8sClient)
	function = deleteIfNotRequired(monitor, user2, users, dbClient)
	assert.NotNil(t, function)
	err = function(k8sClient)
	assert.NoError(t, err)
}

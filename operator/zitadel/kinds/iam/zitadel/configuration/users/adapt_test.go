package users

import (
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
	databasemock "github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsers_Adapt_CreateFirst(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	users := map[string]string{"test": "testpw"}
	dbClient := databasemock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}

	queried := map[string]interface{}{}
	database.SetDatabaseInQueried(queried, &database.Current{
		Host:  "test",
		Port:  "test",
		Users: []string{},
	})
	dbClient.EXPECT().AddUser(monitor, "test", client)

	query, _, err := AdaptFunc(monitor, users, dbClient)
	assert.NoError(t, err)
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	err = ensure(client)
	assert.NoError(t, err)
}

func TestUsers_Adapt_DoNothing(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	users := map[string]string{"test": "testpw"}
	dbClient := databasemock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}

	queried := map[string]interface{}{}
	database.SetDatabaseInQueried(queried, &database.Current{
		Host:  "test",
		Port:  "test",
		Users: []string{"test"},
	})

	query, _, err := AdaptFunc(monitor, users, dbClient)
	assert.NoError(t, err)
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)
	assert.NoError(t, ensure(client))
}

func TestUsers_Adapt_Add(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	users := map[string]string{"test": "testpw", "test2": "testpw"}
	dbClient := databasemock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}

	queried := map[string]interface{}{}
	database.SetDatabaseInQueried(queried, &database.Current{
		Host:  "test",
		Port:  "test",
		Users: []string{"test"},
	})
	dbClient.EXPECT().AddUser(monitor, "test2", client)

	query, _, err := AdaptFunc(monitor, users, dbClient)
	assert.NoError(t, err)
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	err = ensure(client)
	assert.NoError(t, err)
}

func TestUsers_Adapt_Delete(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	users := map[string]string{"test": "testpw", "test2": "testpw"}
	dbClient := databasemock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}

	queried := map[string]interface{}{}
	database.SetDatabaseInQueried(queried, &database.Current{
		Host:  "test",
		Port:  "test",
		Users: []string{"test", "test2", "test3"},
	})

	dbClient.EXPECT().DeleteUser(monitor, "test3", client)

	query, _, err := AdaptFunc(monitor, users, dbClient)
	assert.NoError(t, err)
	ensure, err := query(client, queried)
	err = ensure(client)
	assert.NoError(t, err)
}

func TestUsers_Adapt_DeleteMultiple(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	users := map[string]string{}
	dbClient := databasemock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}

	queried := map[string]interface{}{}
	database.SetDatabaseInQueried(queried, &database.Current{
		Host:  "test",
		Port:  "test",
		Users: []string{"test", "test2", "test3"},
	})

	dbClient.EXPECT().DeleteUser(monitor, "test", client)
	dbClient.EXPECT().DeleteUser(monitor, "test2", client)
	dbClient.EXPECT().DeleteUser(monitor, "test3", client)

	query, _, err := AdaptFunc(monitor, users, dbClient)
	assert.NoError(t, err)
	ensure, err := query(client, queried)
	err = ensure(client)
	assert.NoError(t, err)
}

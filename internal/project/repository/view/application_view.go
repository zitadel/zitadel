package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func ApplicationByID(db *gorm.DB, table, appID string) (*model.ApplicationView, error) {
	app := new(model.ApplicationView)
	query := view.PrepareGetByKey(table, model.ApplicationSearchKey(proj_model.AppSearchKeyAppID), appID)
	err := query(db, app)
	return app, err
}

func ApplicationByOIDCClientID(db *gorm.DB, table, clientID string) (*model.ApplicationView, error) {
	app := new(model.ApplicationView)
	clientIDQuery := model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyOIDCClientID, Value: clientID, Method: global_model.SearchMethodEquals}
	query := view.PrepareGetByQuery(table, clientIDQuery)
	err := query(db, app)
	return app, err
}

func ApplicationByProjectIDAndAppName(db *gorm.DB, table, projectID, appName string) (*model.ApplicationView, error) {
	app := new(model.ApplicationView)
	projectIDQuery := model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals}
	appNameQuery := model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyName, Value: appName, Method: global_model.SearchMethodEquals}
	query := view.PrepareGetByQuery(table, projectIDQuery, appNameQuery)
	err := query(db, app)
	return app, err
}

func SearchApplications(db *gorm.DB, table string, req *proj_model.ApplicationSearchRequest) ([]*model.ApplicationView, int, error) {
	apps := make([]*model.ApplicationView, 0)
	query := view.PrepareSearchQuery(table, model.ApplicationSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &apps)
	return apps, count, err
}

func PutApplication(db *gorm.DB, table string, app *model.ApplicationView) error {
	save := view.PrepareSave(table)
	return save(db, app)
}

func DeleteApplication(db *gorm.DB, table, appID string) error {
	delete := view.PrepareDeleteByKey(table, model.ApplicationSearchKey(proj_model.AppSearchKeyAppID), appID)
	return delete(db)
}

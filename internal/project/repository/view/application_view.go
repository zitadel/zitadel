package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func ApplicationByID(db *gorm.DB, table, projectID, appID string) (*model.ApplicationView, error) {
	app := new(model.ApplicationView)
	projectIDQuery := &model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals}
	appIDQuery := &model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyAppID, Value: appID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, projectIDQuery, appIDQuery)
	err := query(db, app)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-DGdfx", "Errors.Application.NotFound")
	}
	return app, err
}

func ApplicationsByProjectID(db *gorm.DB, table, projectID string) ([]*model.ApplicationView, error) {
	applications := make([]*model.ApplicationView, 0)
	queries := []*proj_model.ApplicationSearchQuery{
		{Key: proj_model.AppSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.ApplicationSearchRequest{Queries: queries})
	_, err := query(db, &applications)
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func ApplicationByOIDCClientID(db *gorm.DB, table, clientID string) (*model.ApplicationView, error) {
	app := new(model.ApplicationView)
	clientIDQuery := model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyOIDCClientID, Value: clientID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, clientIDQuery)
	err := query(db, app)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-DG1qh", "Errors.Project.App.NotFound")
	}
	return app, err
}

func ApplicationByProjectIDAndAppName(db *gorm.DB, table, projectID, appName string) (*model.ApplicationView, error) {
	app := new(model.ApplicationView)
	projectIDQuery := model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals}
	appNameQuery := model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyName, Value: appName, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, projectIDQuery, appNameQuery)
	err := query(db, app)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Jqw1z", "Errors.Project.App.NotFound")
	}
	return app, err
}

func SearchApplications(db *gorm.DB, table string, req *proj_model.ApplicationSearchRequest) ([]*model.ApplicationView, uint64, error) {
	apps := make([]*model.ApplicationView, 0)
	query := repository.PrepareSearchQuery(table, model.ApplicationSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &apps)
	return apps, count, err
}

func PutApplication(db *gorm.DB, table string, app *model.ApplicationView) error {
	save := repository.PrepareSave(table)
	return save(db, app)
}

func PutApplications(db *gorm.DB, table string, apps ...*model.ApplicationView) error {
	save := repository.PrepareBulkSave(table)
	s := make([]interface{}, len(apps))
	for i, app := range apps {
		s[i] = app
	}
	return save(db, s...)
}

func DeleteApplication(db *gorm.DB, table, appID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ApplicationSearchKey(proj_model.AppSearchKeyAppID), appID)
	return delete(db)
}

func DeleteApplicationsByProjectID(db *gorm.DB, table, projectID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ApplicationSearchKey(proj_model.AppSearchKeyProjectID), projectID)
	return delete(db)
}

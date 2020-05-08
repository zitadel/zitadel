package view

import (
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func ApplicationByID(db *gorm.DB, table, appID string) (*model.ApplicationView, error) {
	app := new(model.ApplicationView)
	query := view.PrepareGetByKey(table, model.ApplicationSearchKey(proj_model.APPLICATIONSEARCHKEY_APP_ID), appID)
	err := query(db, app)
	return app, err
}

func SearchApplications(db *gorm.DB, table string, req *proj_model.ApplicationSearchRequest) ([]*model.ApplicationView, int, error) {
	apps := make([]*model.ApplicationView, 0)
	query := view.PrepareSearchQuery(table, model.ApplicationSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &apps)
	if err != nil {
		return nil, 0, err
	}
	return apps, count, nil
}

func PutApplication(db *gorm.DB, table string, app *model.ApplicationView) error {
	save := view.PrepareSave(table)
	return save(db, app)
}

func DeleteApplication(db *gorm.DB, table, appID string) error {
	delete := view.PrepareDeleteByKey(table, model.ApplicationSearchKey(proj_model.APPLICATIONSEARCHKEY_APP_ID), appID)
	return delete(db)
}

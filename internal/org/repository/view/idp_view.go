package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func IdpByID(db *gorm.DB, table, idpID string) (*model.IdpConfigView, error) {
	idp := new(model.IdpConfigView)
	userIDQuery := &model.IdpConfigSearchQuery{Key: org_model.IdpConfigSearchKeyIdpConfigID, Value: idpID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, userIDQuery)
	err := query(db, idp)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Ahq2s", "Errors.Iam.IdpNotExisting")
	}
	return idp, err
}

func SearchIdps(db *gorm.DB, table string, req *org_model.IdpConfigSearchRequest) ([]*model.IdpConfigView, uint64, error) {
	idps := make([]*model.IdpConfigView, 0)
	query := repository.PrepareSearchQuery(table, model.IdpConfigSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &idps)
	if err != nil {
		return nil, 0, err
	}
	return idps, count, nil
}

func PutIdp(db *gorm.DB, table string, role *model.IdpConfigView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func DeleteIdp(db *gorm.DB, table, idpID string) error {
	delete := repository.PrepareDeleteByKey(table, model.IdpConfigSearchKey(org_model.IdpConfigSearchKeyIdpConfigID), idpID)

	return delete(db)
}

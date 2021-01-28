package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func IDPByID(db *gorm.DB, table, idpID string) (*model.IDPConfigView, error) {
	idp := new(model.IDPConfigView)
	idpIDQuery := &model.IDPConfigSearchQuery{Key: iam_model.IDPConfigSearchKeyIdpConfigID, Value: idpID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, idpIDQuery)
	err := query(db, idp)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Ahq2s", "Errors.IAM.IdpNotExisting")
	}
	return idp, err
}

func GetIDPConfigsByAggregateID(db *gorm.DB, table string, aggregateID string) ([]*model.IDPConfigView, error) {
	idps := make([]*model.IDPConfigView, 0)
	queries := []*iam_model.IDPConfigSearchQuery{
		{
			Key:    iam_model.IDPConfigSearchKeyAggregateID,
			Value:  aggregateID,
			Method: global_model.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.IDPConfigSearchRequest{Queries: queries})
	_, err := query(db, &idps)
	if err != nil {
		return nil, err
	}
	return idps, nil
}

func SearchIDPs(db *gorm.DB, table string, req *iam_model.IDPConfigSearchRequest) ([]*model.IDPConfigView, uint64, error) {
	idps := make([]*model.IDPConfigView, 0)
	query := repository.PrepareSearchQuery(table, model.IDPConfigSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &idps)
	if err != nil {
		return nil, 0, err
	}
	return idps, count, nil
}

func PutIDP(db *gorm.DB, table string, idp *model.IDPConfigView) error {
	save := repository.PrepareSave(table)
	return save(db, idp)
}

func DeleteIDP(db *gorm.DB, table, idpID string) error {
	delete := repository.PrepareDeleteByKey(table, model.IDPConfigSearchKey(iam_model.IDPConfigSearchKeyIdpConfigID), idpID)

	return delete(db)
}

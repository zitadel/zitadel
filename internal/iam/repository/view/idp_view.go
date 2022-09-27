package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
)

func IDPByID(db *gorm.DB, table, idpID, instanceID string) (*model.IDPConfigView, error) {
	idp := new(model.IDPConfigView)
	idpIDQuery := &model.IDPConfigSearchQuery{Key: iam_model.IDPConfigSearchKeyIdpConfigID, Value: idpID, Method: domain.SearchMethodEquals}
	instanceIDQuery := &model.IDPConfigSearchQuery{Key: iam_model.IDPConfigSearchKeyInstanceID, Value: instanceID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, idpIDQuery, instanceIDQuery)
	err := query(db, idp)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Ahq2s", "Errors.IDP.NotExisting")
	}
	return idp, err
}

func GetIDPConfigsByAggregateID(db *gorm.DB, table string, aggregateID, instanceID string) ([]*model.IDPConfigView, error) {
	idps := make([]*model.IDPConfigView, 0)
	queries := []*iam_model.IDPConfigSearchQuery{
		{
			Key:    iam_model.IDPConfigSearchKeyAggregateID,
			Value:  aggregateID,
			Method: domain.SearchMethodEquals,
		}, {
			Key:    iam_model.IDPConfigSearchKeyInstanceID,
			Value:  instanceID,
			Method: domain.SearchMethodEquals,
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

func DeleteIDP(db *gorm.DB, table, idpID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{model.IDPConfigSearchKey(iam_model.IDPConfigSearchKeyIdpConfigID), idpID},
		repository.Key{model.IDPConfigSearchKey(iam_model.IDPConfigSearchKeyInstanceID), instanceID},
	)

	return delete(db)
}

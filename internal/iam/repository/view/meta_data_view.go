package view

import (
	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func GetMetaDataList(db *gorm.DB, table string, aggregateID string) ([]*model.MetaDataView, error) {
	texts := make([]*model.MetaDataView, 0)
	queries := []*domain.MetaDataSearchQuery{
		{
			Key:    domain.MetaDataSearchKeyAggregateID,
			Value:  aggregateID,
			Method: domain.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.MetaDataSearchRequest{Queries: queries})
	_, err := query(db, &texts)
	if err != nil {
		return nil, err
	}
	return texts, nil
}

func MetaDataByKey(db *gorm.DB, table, aggregateID, key string) (*model.MetaDataView, error) {
	customText := new(model.MetaDataView)
	aggregateIDQuery := &model.MetaDataSearchQuery{Key: domain.MetaDataSearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	keyQuery := &model.MetaDataSearchQuery{Key: domain.MetaDataSearchKeyKey, Value: key, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, keyQuery)
	err := query(db, customText)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-29kkd", "Errors.MetaData.NotExisting")
	}
	return customText, err
}

func PutMetaData(db *gorm.DB, table string, customText *model.MetaDataView) error {
	save := repository.PrepareSave(table)
	return save(db, customText)
}

func DeleteMetaData(db *gorm.DB, table, aggregateID, key string) error {
	aggregateIDQuery := repository.Key{Key: model.MetaDataSearchKey(domain.MetaDataSearchKeyAggregateID), Value: aggregateID}
	keyQuery := repository.Key{Key: model.MetaDataSearchKey(domain.MetaDataSearchKeyKey), Value: key}
	delete := repository.PrepareDeleteByKeys(table, aggregateIDQuery, keyQuery)
	return delete(db)
}

func DeleteMetaDataByAggregateID(db *gorm.DB, table, aggregateID string) error {
	aggregateIDQuery := repository.Key{Key: model.MetaDataSearchKey(domain.MetaDataSearchKeyAggregateID), Value: aggregateID}
	delete := repository.PrepareDeleteByKeys(table, aggregateIDQuery)
	return delete(db)
}

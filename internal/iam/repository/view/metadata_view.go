package view

import (
	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func GetMetadataList(db *gorm.DB, table string, aggregateID string) ([]*model.MetadataView, error) {
	texts := make([]*model.MetadataView, 0)
	queries := []*domain.MetadataSearchQuery{
		{
			Key:    domain.MetadataSearchKeyAggregateID,
			Value:  aggregateID,
			Method: domain.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.MetadataSearchRequest{Queries: queries})
	_, err := query(db, &texts)
	if err != nil {
		return nil, err
	}
	return texts, nil
}

func MetadataByKey(db *gorm.DB, table, aggregateID, key string) (*model.MetadataView, error) {
	customText := new(model.MetadataView)
	aggregateIDQuery := &model.MetadataSearchQuery{Key: domain.MetadataSearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	keyQuery := &model.MetadataSearchQuery{Key: domain.MetadataSearchKeyKey, Value: key, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, keyQuery)
	err := query(db, customText)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-29kkd", "Errors.Metadata.NotExisting")
	}
	return customText, err
}

func MetadataByKeyAndResourceOwner(db *gorm.DB, table, aggregateID, resourceOwner, key string) (*model.MetadataView, error) {
	customText := new(model.MetadataView)
	aggregateIDQuery := &model.MetadataSearchQuery{Key: domain.MetadataSearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	resourceOwnerQuery := &model.MetadataSearchQuery{Key: domain.MetadataSearchKeyResourceOwner, Value: resourceOwner, Method: domain.SearchMethodEquals}
	keyQuery := &model.MetadataSearchQuery{Key: domain.MetadataSearchKeyKey, Value: key, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, resourceOwnerQuery, keyQuery)
	err := query(db, customText)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-29kkd", "Errors.Metadata.NotExisting")
	}
	return customText, err
}

func SearchMetadata(db *gorm.DB, table string, req *domain.MetadataSearchRequest) ([]*model.MetadataView, uint64, error) {
	metadata := make([]*model.MetadataView, 0)
	query := repository.PrepareSearchQuery(table, model.MetadataSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &metadata)
	if err != nil {
		return nil, 0, err
	}
	return metadata, count, nil
}

func PutMetadata(db *gorm.DB, table string, customText *model.MetadataView) error {
	save := repository.PrepareSave(table)
	return save(db, customText)
}

func DeleteMetadata(db *gorm.DB, table, aggregateID, key string) error {
	aggregateIDQuery := repository.Key{Key: model.MetadataSearchKey(domain.MetadataSearchKeyAggregateID), Value: aggregateID}
	keyQuery := repository.Key{Key: model.MetadataSearchKey(domain.MetadataSearchKeyKey), Value: key}
	delete := repository.PrepareDeleteByKeys(table, aggregateIDQuery, keyQuery)
	return delete(db)
}

func DeleteMetadataByAggregateID(db *gorm.DB, table, aggregateID string) error {
	aggregateIDQuery := repository.Key{Key: model.MetadataSearchKey(domain.MetadataSearchKeyAggregateID), Value: aggregateID}
	delete := repository.PrepareDeleteByKeys(table, aggregateIDQuery)
	return delete(db)
}

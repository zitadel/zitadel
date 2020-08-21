package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetIdpProviderByAggregateIDAndConfigID(db *gorm.DB, table, aggregateID, idpConfigID string) (*model.IdpProviderView, error) {
	policy := new(model.IdpProviderView)
	aggIDQuery := &model.IdpProviderSearchQuery{Key: iam_model.IdpProviderSearchKeyAggregateID, Value: aggregateID, Method: global_model.SearchMethodEquals}
	idpConfigIDQuery := &model.IdpProviderSearchQuery{Key: iam_model.IdpProviderSearchKeyIdpConfigID, Value: idpConfigID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggIDQuery, idpConfigIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Skvi8", "Errors.Iam.LoginPolicy.IdpProviderNotExisting")
	}
	return policy, err
}

func IdpProvidersByIdpConfigID(db *gorm.DB, table string, idpConfigID string) ([]*model.IdpProviderView, error) {
	members := make([]*model.IdpProviderView, 0)
	queries := []*iam_model.IdpProviderSearchQuery{
		{
			Key:    iam_model.IdpProviderSearchKeyIdpConfigID,
			Value:  idpConfigID,
			Method: global_model.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.IdpProviderSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func SearchIdpProviders(db *gorm.DB, table string, req *iam_model.IdpProviderSearchRequest) ([]*model.IdpProviderView, uint64, error) {
	providers := make([]*model.IdpProviderView, 0)
	query := repository.PrepareSearchQuery(table, model.IdpProviderSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &providers)
	if err != nil {
		return nil, 0, err
	}
	return providers, count, nil
}

func PutIdpProvider(db *gorm.DB, table string, provider *model.IdpProviderView) error {
	save := repository.PrepareSave(table)
	return save(db, provider)
}

func PutIdpProviders(db *gorm.DB, table string, providers ...*model.IdpProviderView) error {
	save := repository.PrepareBulkSave(table)
	p := make([]interface{}, len(providers))
	for i, provider := range providers {
		p[i] = provider
	}
	return save(db, p...)
}

func DeleteIdpProvider(db *gorm.DB, table, aggregateID, idpConfigID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.IdpProviderSearchKey(iam_model.IdpProviderSearchKeyAggregateID), Value: aggregateID},
		repository.Key{Key: model.IdpProviderSearchKey(iam_model.IdpProviderSearchKeyIdpConfigID), Value: idpConfigID},
	)
	return delete(db)
}

func DeleteIdpProvidersByAggregateID(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.IdpProviderSearchKey(iam_model.IdpProviderSearchKeyAggregateID), Value: aggregateID},
	)
	return delete(db)
}

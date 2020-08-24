package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetIDPProviderByAggregateIDAndConfigID(db *gorm.DB, table, aggregateID, idpConfigID string) (*model.IDPProviderView, error) {
	policy := new(model.IDPProviderView)
	aggIDQuery := &model.IDPProviderSearchQuery{Key: iam_model.IDPProviderSearchKeyAggregateID, Value: aggregateID, Method: global_model.SearchMethodEquals}
	idpConfigIDQuery := &model.IDPProviderSearchQuery{Key: iam_model.IDPProviderSearchKeyIdpConfigID, Value: idpConfigID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggIDQuery, idpConfigIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Skvi8", "Errors.Iam.LoginPolicy.IdpProviderNotExisting")
	}
	return policy, err
}

func IDPProvidersByIdpConfigID(db *gorm.DB, table string, idpConfigID string) ([]*model.IDPProviderView, error) {
	members := make([]*model.IDPProviderView, 0)
	queries := []*iam_model.IDPProviderSearchQuery{
		{
			Key:    iam_model.IDPProviderSearchKeyIdpConfigID,
			Value:  idpConfigID,
			Method: global_model.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.IDPProviderSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func SearchIDPProviders(db *gorm.DB, table string, req *iam_model.IDPProviderSearchRequest) ([]*model.IDPProviderView, uint64, error) {
	providers := make([]*model.IDPProviderView, 0)
	query := repository.PrepareSearchQuery(table, model.IDPProviderSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &providers)
	if err != nil {
		return nil, 0, err
	}
	return providers, count, nil
}

func PutIDPProvider(db *gorm.DB, table string, provider *model.IDPProviderView) error {
	save := repository.PrepareSave(table)
	return save(db, provider)
}

func PutIDPProviders(db *gorm.DB, table string, providers ...*model.IDPProviderView) error {
	save := repository.PrepareBulkSave(table)
	p := make([]interface{}, len(providers))
	for i, provider := range providers {
		p[i] = provider
	}
	return save(db, p...)
}

func DeleteIDPProvider(db *gorm.DB, table, aggregateID, idpConfigID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyAggregateID), Value: aggregateID},
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyIdpConfigID), Value: idpConfigID},
	)
	return delete(db)
}

func DeleteIDPProvidersByAggregateID(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyAggregateID), Value: aggregateID},
	)
	return delete(db)
}

package view

import (
	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func GetIDPProviderByAggregateIDAndConfigID(db *gorm.DB, table, aggregateID, idpConfigID, instanceID string) (*model.IDPProviderView, error) {
	policy := new(model.IDPProviderView)
	aggIDQuery := &model.IDPProviderSearchQuery{Key: iam_model.IDPProviderSearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	idpConfigIDQuery := &model.IDPProviderSearchQuery{Key: iam_model.IDPProviderSearchKeyIdpConfigID, Value: idpConfigID, Method: domain.SearchMethodEquals}
	instanceIDQuery := &model.IDPProviderSearchQuery{Key: iam_model.IDPProviderSearchKeyInstanceID, Value: instanceID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggIDQuery, idpConfigIDQuery, instanceIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Skvi8", "Errors.IAM.LoginPolicy.IDP.NotExisting")
	}
	return policy, err
}

func IDPProvidersByIdpConfigID(db *gorm.DB, table, idpConfigID, instanceID string) ([]*model.IDPProviderView, error) {
	providers := make([]*model.IDPProviderView, 0)
	queries := []*iam_model.IDPProviderSearchQuery{
		{
			Key:    iam_model.IDPProviderSearchKeyIdpConfigID,
			Value:  idpConfigID,
			Method: domain.SearchMethodEquals,
		},
		{
			Key:    iam_model.IDPProviderSearchKeyInstanceID,
			Value:  instanceID,
			Method: domain.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.IDPProviderSearchRequest{Queries: queries})
	_, err := query(db, &providers)
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func IDPProvidersByAggregateIDAndState(db *gorm.DB, table string, aggregateID, instanceID string, idpConfigState iam_model.IDPConfigState) ([]*model.IDPProviderView, error) {
	providers := make([]*model.IDPProviderView, 0)
	queries := []*iam_model.IDPProviderSearchQuery{
		{
			Key:    iam_model.IDPProviderSearchKeyAggregateID,
			Value:  aggregateID,
			Method: domain.SearchMethodEquals,
		},
		{
			Key:    iam_model.IDPProviderSearchKeyState,
			Value:  int(idpConfigState),
			Method: domain.SearchMethodEquals,
		},
		{
			Key:    iam_model.IDPProviderSearchKeyInstanceID,
			Value:  instanceID,
			Method: domain.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.IDPProviderSearchRequest{Queries: queries})
	_, err := query(db, &providers)
	if err != nil {
		return nil, err
	}
	return providers, nil
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

func DeleteIDPProvider(db *gorm.DB, table, aggregateID, idpConfigID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyAggregateID), Value: aggregateID},
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyIdpConfigID), Value: idpConfigID},
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyInstanceID), Value: instanceID},
	)
	return delete(db)
}

func DeleteIDPProvidersByAggregateID(db *gorm.DB, table, aggregateID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyAggregateID), Value: aggregateID},
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyInstanceID), Value: instanceID},
	)
	return delete(db)
}

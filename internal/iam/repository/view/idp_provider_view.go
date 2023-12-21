package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func GetIDPProviderByAggregateIDAndConfigID(db *gorm.DB, table, aggregateID, idpConfigID, instanceID string) (*model.IDPProviderView, error) {
	policy := new(model.IDPProviderView)
	aggIDQuery := &model.IDPProviderSearchQuery{Key: iam_model.IDPProviderSearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	idpConfigIDQuery := &model.IDPProviderSearchQuery{Key: iam_model.IDPProviderSearchKeyIdpConfigID, Value: idpConfigID, Method: domain.SearchMethodEquals}
	instanceIDQuery := &model.IDPProviderSearchQuery{Key: iam_model.IDPProviderSearchKeyInstanceID, Value: instanceID, Method: domain.SearchMethodEquals}
	ownerRemovedQuery := &model.IDPProviderSearchQuery{Key: iam_model.IDPProviderSearchKeyOwnerRemoved, Value: false, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggIDQuery, idpConfigIDQuery, instanceIDQuery, ownerRemovedQuery)
	err := query(db, policy)
	if zerrors.IsNotFound(err) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-Skvi8", "Errors.IAM.LoginPolicy.IDP.NotExisting")
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
		{
			Key:    iam_model.IDPProviderSearchKeyOwnerRemoved,
			Value:  false,
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
		{
			Key:    iam_model.IDPProviderSearchKeyOwnerRemoved,
			Value:  false,
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

func DeleteInstanceIDPProviders(db *gorm.DB, table, instanceID string) error {
	delete := repository.PrepareDeleteByKey(table,
		model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyInstanceID),
		instanceID,
	)
	return delete(db)
}

func UpdateOrgOwnerRemovedIDPProviders(db *gorm.DB, table, instanceID, aggID string) error {
	update := repository.PrepareUpdateByKeys(table,
		model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyOwnerRemoved),
		true,
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyInstanceID), Value: instanceID},
		repository.Key{Key: model.IDPProviderSearchKey(iam_model.IDPProviderSearchKeyAggregateID), Value: aggID},
	)
	return update(db)
}

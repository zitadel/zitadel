package view

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"

	"github.com/jinzhu/gorm"

	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
)

func ExternalIDPByExternalUserIDAndIDPConfigID(db *gorm.DB, table, externalUserID, idpConfigID, instanceID string) (*model.ExternalIDPView, error) {
	user := new(model.ExternalIDPView)
	userIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyExternalUserID,
		Method: domain.SearchMethodEquals,
		Value:  externalUserID,
	}
	idpConfigIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyIdpConfigID,
		Method: domain.SearchMethodEquals,
		Value:  idpConfigID,
	}
	instanceIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	ownerRemovedQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyOwnerRemoved,
		Method: domain.SearchMethodEquals,
		Value:  false,
	}
	query := repository.PrepareGetByQuery(table, userIDQuery, idpConfigIDQuery, instanceIDQuery, ownerRemovedQuery)
	err := query(db, user)
	if zerrors.IsNotFound(err) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-Mso9f", "Errors.ExternalIDP.NotFound")
	}
	return user, err
}

func ExternalIDPByExternalUserIDAndIDPConfigIDAndResourceOwner(db *gorm.DB, table, externalUserID, idpConfigID, resourceOwner, instanceID string) (*model.ExternalIDPView, error) {
	user := new(model.ExternalIDPView)
	userIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyExternalUserID,
		Method: domain.SearchMethodEquals,
		Value:  externalUserID,
	}
	idpConfigIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyIdpConfigID,
		Method: domain.SearchMethodEquals,
		Value:  idpConfigID,
	}
	resourceOwnerQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyResourceOwner,
		Method: domain.SearchMethodEquals,
		Value:  resourceOwner,
	}
	instanceIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	ownerRemovedQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyOwnerRemoved,
		Method: domain.SearchMethodEquals,
		Value:  false,
	}
	query := repository.PrepareGetByQuery(table, userIDQuery, idpConfigIDQuery, resourceOwnerQuery, instanceIDQuery, ownerRemovedQuery)
	err := query(db, user)
	if zerrors.IsNotFound(err) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-Sf8sd", "Errors.ExternalIDP.NotFound")
	}
	return user, err
}

func ExternalIDPsByIDPConfigID(db *gorm.DB, table, idpConfigID, instanceID string) ([]*model.ExternalIDPView, error) {
	externalIDPs := make([]*model.ExternalIDPView, 0)
	orgIDQuery := &usr_model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyIdpConfigID,
		Method: domain.SearchMethodEquals,
		Value:  idpConfigID,
	}
	instanceIDQuery := &usr_model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	ownerRemovedQuery := &usr_model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyOwnerRemoved,
		Method: domain.SearchMethodEquals,
		Value:  false,
	}
	query := repository.PrepareSearchQuery(table, model.ExternalIDPSearchRequest{
		Queries: []*usr_model.ExternalIDPSearchQuery{orgIDQuery, instanceIDQuery, ownerRemovedQuery},
	})
	_, err := query(db, &externalIDPs)
	return externalIDPs, err
}

func PutExternalIDPs(db *gorm.DB, table string, externalIDPs ...*model.ExternalIDPView) error {
	save := repository.PrepareBulkSave(table)
	u := make([]interface{}, len(externalIDPs))
	for i, idp := range externalIDPs {
		u[i] = idp
	}
	return save(db, u...)
}

func PutExternalIDP(db *gorm.DB, table string, idp *model.ExternalIDPView) error {
	save := repository.PrepareSave(table)
	return save(db, idp)
}

func DeleteExternalIDP(db *gorm.DB, table, externalUserID, idpConfigID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyExternalUserID), Value: externalUserID},
		repository.Key{Key: model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyIdpConfigID), Value: idpConfigID},
		repository.Key{Key: model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyInstanceID), Value: instanceID},
	)
	return delete(db)
}

func DeleteExternalIDPsByUserID(db *gorm.DB, table, userID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyUserID), userID},
		repository.Key{model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}

func DeleteInstanceExternalIDPs(db *gorm.DB, table, instanceID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyInstanceID), instanceID)
	return delete(db)
}

func UpdateOrgOwnerRemovedExternalIDPs(db *gorm.DB, table, instanceID, aggID string) error {
	update := repository.PrepareUpdateByKeys(table,
		model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyOwnerRemoved),
		true,
		repository.Key{Key: model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyInstanceID), Value: instanceID},
		repository.Key{Key: model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyResourceOwner), Value: aggID},
	)
	return update(db)
}

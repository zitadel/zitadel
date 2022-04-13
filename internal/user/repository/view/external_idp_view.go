package view

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/view/repository"

	"github.com/jinzhu/gorm"

	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
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
	query := repository.PrepareGetByQuery(table, userIDQuery, idpConfigIDQuery, instanceIDQuery)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Mso9f", "Errors.ExternalIDP.NotFound")
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
	query := repository.PrepareGetByQuery(table, userIDQuery, idpConfigIDQuery, resourceOwnerQuery, instanceIDQuery)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Sf8sd", "Errors.ExternalIDP.NotFound")
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
	query := repository.PrepareSearchQuery(table, model.ExternalIDPSearchRequest{
		Queries: []*usr_model.ExternalIDPSearchQuery{orgIDQuery, instanceIDQuery},
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

package view

import (
	"github.com/caos/zitadel/internal/view/repository"

	"github.com/jinzhu/gorm"

	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

func ExternalIDPByExternalUserIDAndIDPConfigID(db *gorm.DB, table, externalUserID, idpConfigID string) (*model.ExternalIDPView, error) {
	user := new(model.ExternalIDPView)
	userIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyExternalUserID,
		Method: global_model.SearchMethodEquals,
		Value:  externalUserID,
	}
	idpConfigIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyIdpConfigID,
		Method: global_model.SearchMethodEquals,
		Value:  idpConfigID,
	}
	query := repository.PrepareGetByQuery(table, userIDQuery, idpConfigIDQuery)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Mso9f", "Errors.ExternalIDP.NotFound")
	}
	return user, err
}

func ExternalIDPByExternalUserIDAndIDPConfigIDAndResourceOwner(db *gorm.DB, table, externalUserID, idpConfigID, resourceOwner string) (*model.ExternalIDPView, error) {
	user := new(model.ExternalIDPView)
	userIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyExternalUserID,
		Method: global_model.SearchMethodEquals,
		Value:  externalUserID,
	}
	idpConfigIDQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyIdpConfigID,
		Method: global_model.SearchMethodEquals,
		Value:  idpConfigID,
	}
	resourceOwnerQuery := &model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyResourceOwner,
		Method: global_model.SearchMethodEquals,
		Value:  resourceOwner,
	}
	query := repository.PrepareGetByQuery(table, userIDQuery, idpConfigIDQuery, resourceOwnerQuery)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Sf8sd", "Errors.ExternalIDP.NotFound")
	}
	return user, err
}

func ExternalIDPsByIDPConfigID(db *gorm.DB, table, idpConfigID string) ([]*model.ExternalIDPView, error) {
	externalIDPs := make([]*model.ExternalIDPView, 0)
	orgIDQuery := &usr_model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyIdpConfigID,
		Method: global_model.SearchMethodEquals,
		Value:  idpConfigID,
	}
	query := repository.PrepareSearchQuery(table, model.ExternalIDPSearchRequest{
		Queries: []*usr_model.ExternalIDPSearchQuery{orgIDQuery},
	})
	_, err := query(db, &externalIDPs)
	return externalIDPs, err
}

func ExternalIDPsByIDPConfigIDAndResourceOwner(db *gorm.DB, table, idpConfigID, resourceOwner string) ([]*model.ExternalIDPView, error) {
	externalIDPs := make([]*model.ExternalIDPView, 0)
	idpConfigIDQuery := &usr_model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyIdpConfigID,
		Method: global_model.SearchMethodEquals,
		Value:  idpConfigID,
	}
	orgIDQuery := &usr_model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyResourceOwner,
		Method: global_model.SearchMethodEquals,
		Value:  resourceOwner,
	}
	query := repository.PrepareSearchQuery(table, model.ExternalIDPSearchRequest{
		Queries: []*usr_model.ExternalIDPSearchQuery{orgIDQuery, idpConfigIDQuery},
	})
	_, err := query(db, &externalIDPs)
	return externalIDPs, err
}

func ExternalIDPsByIDPConfigIDAndResourceOwners(db *gorm.DB, table, idpConfigID string, resourceOwners []string) ([]*model.ExternalIDPView, error) {
	externalIDPs := make([]*model.ExternalIDPView, 0)
	idpConfigIDQuery := &usr_model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyIdpConfigID,
		Method: global_model.SearchMethodEquals,
		Value:  idpConfigID,
	}
	orgIDQuery := &usr_model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyResourceOwner,
		Method: global_model.SearchMethodIsOneOf,
		Value:  resourceOwners,
	}
	query := repository.PrepareSearchQuery(table, model.ExternalIDPSearchRequest{
		Queries: []*usr_model.ExternalIDPSearchQuery{orgIDQuery, idpConfigIDQuery},
	})
	_, err := query(db, &externalIDPs)
	return externalIDPs, err
}

func ExternalIDPsByUserID(db *gorm.DB, table, userID string) ([]*model.ExternalIDPView, error) {
	externalIDPs := make([]*model.ExternalIDPView, 0)
	orgIDQuery := &usr_model.ExternalIDPSearchQuery{
		Key:    usr_model.ExternalIDPSearchKeyUserID,
		Method: global_model.SearchMethodEquals,
		Value:  userID,
	}
	query := repository.PrepareSearchQuery(table, model.ExternalIDPSearchRequest{
		Queries: []*usr_model.ExternalIDPSearchQuery{orgIDQuery},
	})
	_, err := query(db, &externalIDPs)
	return externalIDPs, err
}

func SearchExternalIDPs(db *gorm.DB, table string, req *usr_model.ExternalIDPSearchRequest) ([]*model.ExternalIDPView, uint64, error) {
	externalIDPs := make([]*model.ExternalIDPView, 0)
	query := repository.PrepareSearchQuery(table, model.ExternalIDPSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &externalIDPs)
	if err != nil {
		return nil, 0, err
	}
	return externalIDPs, count, nil
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

func DeleteExternalIDP(db *gorm.DB, table, externalUserID, idpConfigID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyExternalUserID), Value: externalUserID},
		repository.Key{Key: model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyIdpConfigID), Value: idpConfigID},
	)
	return delete(db)
}

func DeleteExternalIDPsByUserID(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ExternalIDPSearchKey(usr_model.ExternalIDPSearchKeyUserID), userID)
	return delete(db)
}

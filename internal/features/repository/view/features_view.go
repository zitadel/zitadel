package view

import (
	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/features/model"
	view_model "github.com/caos/zitadel/internal/features/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func GetDefaultFeatures(db *gorm.DB, table string) ([]*view_model.FeaturesView, error) {
	features := make([]*view_model.FeaturesView, 0)
	queries := []*model.FeaturesSearchQuery{
		{Key: model.FeaturesSearchKeyDefault, Value: true, Method: domain.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, view_model.FeaturesSearchRequest{Queries: queries})
	_, err := query(db, &features)
	if err != nil {
		return nil, err
	}
	return features, nil
}

func GetFeaturesByAggregateID(db *gorm.DB, table, aggregateID string) (*view_model.FeaturesView, error) {
	features := new(view_model.FeaturesView)
	query := repository.PrepareGetByKey(table, view_model.FeaturesSearchKey(model.FeaturesSearchKeyAggregateID), aggregateID)
	err := query(db, features)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Dbf3h", "Errors.Features.NotFound")
	}
	return features, err
}

func PutFeatures(db *gorm.DB, table string, features *view_model.FeaturesView) error {
	save := repository.PrepareSave(table)
	return save(db, features)
}

func PutFeaturesList(db *gorm.DB, table string, featuresList ...*view_model.FeaturesView) error {
	save := repository.PrepareBulkSave(table)
	f := make([]interface{}, len(featuresList))
	for i, features := range featuresList {
		f[i] = features
	}
	return save(db, f...)
}

//
//func DeleteOrg(db *gorm.DB, table, orgID string) error {
//	delete := repository.PrepareDeleteByKey(table, model.OrgSearchKey(org_model.OrgSearchKeyOrgID), orgID)
//	return delete(db)
//}

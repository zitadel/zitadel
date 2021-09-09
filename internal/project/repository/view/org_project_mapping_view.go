package view

import (
	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func OrgProjectMappingByIDs(db *gorm.DB, table, orgID, projectID string) (*model.OrgProjectMapping, error) {
	orgProjectMapping := new(model.OrgProjectMapping)

	projectIDQuery := model.OrgProjectMappingSearchQuery{Key: proj_model.OrgProjectMappingSearchKeyProjectID, Value: projectID, Method: domain.SearchMethodEquals}
	orgIDQuery := model.OrgProjectMappingSearchQuery{Key: proj_model.OrgProjectMappingSearchKeyOrgID, Value: orgID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, projectIDQuery, orgIDQuery)
	err := query(db, orgProjectMapping)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-fn9fs", "Errors.OrgProjectMapping.NotExisting")
	}
	return orgProjectMapping, err
}

func PutOrgProjectMapping(db *gorm.DB, table string, grant *model.OrgProjectMapping) error {
	save := repository.PrepareSave(table)
	return save(db, grant)
}

func DeleteOrgProjectMapping(db *gorm.DB, table, orgID, projectID string) error {
	projectIDSearch := repository.Key{Key: model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyProjectID), Value: projectID}
	orgIDSearch := repository.Key{Key: model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyOrgID), Value: orgID}
	delete := repository.PrepareDeleteByKeys(table, projectIDSearch, orgIDSearch)
	return delete(db)
}

func DeleteOrgProjectMappingsByProjectID(db *gorm.DB, table, projectID string) error {
	delete := repository.PrepareDeleteByKey(table, model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyProjectID), projectID)
	return delete(db)
}

func DeleteOrgProjectMappingsByProjectGrantID(db *gorm.DB, table, projectGrantID string) error {
	delete := repository.PrepareDeleteByKey(table, model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyProjectGrantID), projectGrantID)
	return delete(db)
}

func DeleteOrgProjectMappingsByOrgID(db *gorm.DB, table, orgID string) error {
	delete := repository.PrepareDeleteByKey(table, model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyOrgID), orgID)
	return delete(db)
}

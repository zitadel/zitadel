package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	proj_model "github.com/zitadel/zitadel/internal/project/model"
	"github.com/zitadel/zitadel/internal/project/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
)

func OrgProjectMappingByIDs(db *gorm.DB, table, orgID, projectID, instanceID string) (*model.OrgProjectMapping, error) {
	orgProjectMapping := new(model.OrgProjectMapping)

	projectIDQuery := model.OrgProjectMappingSearchQuery{Key: proj_model.OrgProjectMappingSearchKeyProjectID, Value: projectID, Method: domain.SearchMethodEquals}
	orgIDQuery := model.OrgProjectMappingSearchQuery{Key: proj_model.OrgProjectMappingSearchKeyOrgID, Value: orgID, Method: domain.SearchMethodEquals}
	instanceIDQuery := model.OrgProjectMappingSearchQuery{Key: proj_model.OrgProjectMappingSearchKeyInstanceID, Value: instanceID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, projectIDQuery, orgIDQuery, instanceIDQuery)
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

func DeleteOrgProjectMapping(db *gorm.DB, table, orgID, projectID, instanceID string) error {
	projectIDSearch := repository.Key{Key: model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyProjectID), Value: projectID}
	orgIDSearch := repository.Key{Key: model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyOrgID), Value: orgID}
	instanceIDSearch := repository.Key{Key: model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyInstanceID), Value: instanceID}
	delete := repository.PrepareDeleteByKeys(table, projectIDSearch, orgIDSearch, instanceIDSearch)
	return delete(db)
}

func DeleteOrgProjectMappingsByProjectID(db *gorm.DB, table, projectID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyProjectID), projectID},
		repository.Key{model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}

func DeleteOrgProjectMappingsByProjectGrantID(db *gorm.DB, table, projectGrantID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyProjectGrantID), projectGrantID},
		repository.Key{model.OrgProjectMappingSearchKey(proj_model.OrgProjectMappingSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}

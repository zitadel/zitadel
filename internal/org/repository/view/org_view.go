package view

import (
	"github.com/jinzhu/gorm"

	caos_errs "github.com/caos/zitadel/internal/errors"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func OrgByID(db *gorm.DB, table, orgID string) (*model.OrgView, error) {
	org := new(model.OrgView)
	query := repository.PrepareGetByKey(table, model.OrgSearchKey(org_model.OrgSearchKeyOrgID), orgID)
	err := query(db, org)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-GEwea", "Errors.Org.NotFound")
	}
	return org, err
}

func OrgByPrimaryDomain(db *gorm.DB, table, primaryDomain string) (*model.OrgView, error) {
	org := new(model.OrgView)
	query := repository.PrepareGetByKey(table, model.OrgSearchKey(org_model.OrgSearchKeyOrgDomain), primaryDomain)
	err := query(db, org)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-GEwea", "Errors.Org.NotFound")
	}
	return org, err
}

func SearchOrgs(db *gorm.DB, table string, req *org_model.OrgSearchRequest) ([]*model.OrgView, uint64, error) {
	orgs := make([]*model.OrgView, 0)
	query := repository.PrepareSearchQuery(table, model.OrgSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries, SortingColumn: req.SortingColumn})
	count, err := query(db, &orgs)
	if err != nil {
		return nil, 0, err
	}
	return orgs, count, nil
}

func PutOrg(db *gorm.DB, table string, org *model.OrgView) error {
	save := repository.PrepareSave(table)
	return save(db, org)
}

func DeleteOrg(db *gorm.DB, table, orgID string) error {
	delete := repository.PrepareDeleteByKey(table, model.OrgSearchKey(org_model.OrgSearchKeyOrgID), orgID)
	return delete(db)
}

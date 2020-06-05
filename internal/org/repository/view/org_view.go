package view

import (
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func OrgByID(db *gorm.DB, table, orgID string) (*model.OrgView, error) {
	org := new(model.OrgView)
	query := view.PrepareGetByKey(table, model.OrgSearchKey(org_model.ORGSEARCHKEY_ORG_ID), orgID)
	err := query(db, org)
	return org, err
}

func SearchOrgs(db *gorm.DB, table string, req *org_model.OrgSearchRequest) ([]*model.OrgView, int, error) {
	orgs := make([]*model.OrgView, 0)
	query := view.PrepareSearchQuery(table, model.OrgSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &orgs)
	if err != nil {
		return nil, 0, err
	}
	return orgs, count, nil
}

func GetGlobalOrgByDomain(db *gorm.DB, table, domain string) (*model.OrgView, error) {
	org := new(model.OrgView)
	query := view.PrepareGetByKey(table, model.OrgSearchKey(org_model.ORGSEARCHKEY_ORG_DOMAIN), domain)
	err := query(db, org)
	return org, err
}

func PutOrg(db *gorm.DB, table string, org *model.OrgView) error {
	save := view.PrepareSave(table)
	return save(db, org)
}

func DeleteOrg(db *gorm.DB, table, orgID string) error {
	delete := view.PrepareDeleteByKey(table, model.OrgSearchKey(org_model.ORGSEARCHKEY_ORG_ID), orgID)
	return delete(db)
}

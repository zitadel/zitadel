package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func OrgDomainByDomain(db *gorm.DB, table, domain string) (*model.OrgDomainView, error) {
	domainView := new(model.OrgDomainView)

	domainQuery := &model.OrgDomainSearchQuery{Key: org_model.ORGDOMAINSEARCHKEY_DOMAIN, Value: domain, Method: global_model.SEARCHMETHOD_EQUALS}
	query := view.PrepareGetByQuery(table, domainQuery)
	err := query(db, domainView)
	return domainView, err
}

func SearchOrgDomains(db *gorm.DB, table string, req *org_model.OrgDomainSearchRequest) ([]*model.OrgDomainView, int, error) {
	members := make([]*model.OrgDomainView, 0)
	query := view.PrepareSearchQuery(table, model.OrgDomainSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &members)
	if err != nil {
		return nil, 0, err
	}
	return members, count, nil
}

func OrgDomainsByOrgID(db *gorm.DB, table string, orgID string) ([]*model.OrgDomainView, error) {
	domains := make([]*model.OrgDomainView, 0)
	queries := []*org_model.OrgDomainSearchQuery{
		{
			Key:    org_model.ORGDOMAINSEARCHKEY_ORG_ID,
			Value:  orgID,
			Method: global_model.SEARCHMETHOD_EQUALS,
		},
	}
	query := view.PrepareSearchQuery(table, model.OrgDomainSearchRequest{Queries: queries})
	_, err := query(db, &domains)
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func PutOrgDomain(db *gorm.DB, table string, role *model.OrgDomainView) error {
	save := view.PrepareSave(table)
	return save(db, role)
}

func DeleteOrgDomain(db *gorm.DB, table, domain string) error {
	delete := view.PrepareDeleteByKey(table, model.OrgSearchKey(org_model.ORGDOMAINSEARCHKEY_DOMAIN), domain)
	return delete(db)
}

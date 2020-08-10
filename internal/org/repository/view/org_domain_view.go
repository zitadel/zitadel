package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func OrgDomainByOrgIDAndDomain(db *gorm.DB, table, orgID, domain string) (*model.OrgDomainView, error) {
	domainView := new(model.OrgDomainView)
	orgIDQuery := &model.OrgDomainSearchQuery{Key: org_model.OrgDomainSearchKeyOrgID, Value: orgID, Method: global_model.SearchMethodEquals}
	domainQuery := &model.OrgDomainSearchQuery{Key: org_model.OrgDomainSearchKeyDomain, Value: domain, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, orgIDQuery, domainQuery)
	err := query(db, domainView)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Gqwfq", "Errors.Org.DomainNotOnOrg")
	}
	return domainView, err
}

func VerifiedOrgDomain(db *gorm.DB, table, domain string) (*model.OrgDomainView, error) {
	domainView := new(model.OrgDomainView)
	domainQuery := &model.OrgDomainSearchQuery{Key: org_model.OrgDomainSearchKeyDomain, Value: domain, Method: global_model.SearchMethodEquals}
	verifiedQuery := &model.OrgDomainSearchQuery{Key: org_model.OrgDomainSearchKeyVerified, Value: true, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, domainQuery, verifiedQuery)
	err := query(db, domainView)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Tew2q", "Errors.Org.DomainNotFound")
	}
	return domainView, err
}

func SearchOrgDomains(db *gorm.DB, table string, req *org_model.OrgDomainSearchRequest) ([]*model.OrgDomainView, uint64, error) {
	members := make([]*model.OrgDomainView, 0)
	query := repository.PrepareSearchQuery(table, model.OrgDomainSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
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
			Key:    org_model.OrgDomainSearchKeyOrgID,
			Value:  orgID,
			Method: global_model.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.OrgDomainSearchRequest{Queries: queries})
	_, err := query(db, &domains)
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func PutOrgDomain(db *gorm.DB, table string, role *model.OrgDomainView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func DeleteOrgDomain(db *gorm.DB, table, orgID, domain string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.OrgDomainSearchKey(org_model.OrgDomainSearchKeyDomain), Value: domain},
		repository.Key{Key: model.OrgDomainSearchKey(org_model.OrgDomainSearchKeyOrgID), Value: orgID},
	)
	return delete(db)
}

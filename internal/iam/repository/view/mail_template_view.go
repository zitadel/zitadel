package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetMailTemplateByAggregateID(db *gorm.DB, table, aggregateID string) (*model.MailTemplateView, error) {
	template := new(model.MailTemplateView)
	aggregateIDQuery := &model.MailTemplateSearchQuery{Key: iam_model.MailTemplateSearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery)
	err := query(db, template)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-iPnmU", "Errors.IAM.MailTemplate.NotExisting")
	}
	return template, err
}

func PutMailTemplate(db *gorm.DB, table string, template *model.MailTemplateView) error {
	save := repository.PrepareSave(table)
	return save(db, template)
}

func DeleteMailTemplate(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.MailTemplateSearchKey(iam_model.MailTemplateSearchKeyAggregateID), aggregateID)

	return delete(db)
}

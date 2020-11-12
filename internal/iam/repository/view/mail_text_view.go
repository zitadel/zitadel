package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetMailTextByIDs(db *gorm.DB, table, aggregateID string, textType string, language string) (*model.MailTextView, error) {
	mailText := new(model.MailTextView)
	aggregateIDQuery := &model.MailTextSearchQuery{Key: iam_model.MailTextSearchKeyAggregateID, Value: aggregateID, Method: global_model.SearchMethodEquals}
	textTypeQuery := &model.MailTextSearchQuery{Key: iam_model.MailTextSearchKeyMailTextType, Value: textType, Method: global_model.SearchMethodEquals}
	languageQuery := &model.MailTextSearchQuery{Key: iam_model.MailTextSearchKeyLanguage, Value: language, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, textTypeQuery, languageQuery)
	err := query(db, mailText)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-IiJjm", "Errors.IAM.MailText.NotExisting")
	}
	return mailText, err
}

func PutMailText(db *gorm.DB, table string, mailText *model.MailTextView) error {
	save := repository.PrepareSave(table)
	return save(db, mailText)
}

func DeleteMailText(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.MailTextSearchKey(iam_model.MailTextSearchKeyAggregateID), aggregateID)

	return delete(db)
}

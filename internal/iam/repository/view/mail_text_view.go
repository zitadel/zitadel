package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
	"strings"
)

func GetMailTexts(db *gorm.DB, table string, aggregateID string) ([]*model.MailTextView, error) {
	texts := make([]*model.MailTextView, 0)
	queries := []*iam_model.MailTextSearchQuery{
		{
			Key:    iam_model.MailTextSearchKeyAggregateID,
			Value:  aggregateID,
			Method: domain.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.MailTextSearchRequest{Queries: queries})
	_, err := query(db, &texts)
	if err != nil {
		return nil, err
	}
	return texts, nil
}

func GetMailTextByIDs(db *gorm.DB, table, aggregateID string, textType string, language string) (*model.MailTextView, error) {
	mailText := new(model.MailTextView)
	aggregateIDQuery := &model.MailTextSearchQuery{Key: iam_model.MailTextSearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	textTypeQuery := &model.MailTextSearchQuery{Key: iam_model.MailTextSearchKeyMailTextType, Value: textType, Method: domain.SearchMethodEquals}
	languageQuery := &model.MailTextSearchQuery{Key: iam_model.MailTextSearchKeyLanguage, Value: strings.ToUpper(language), Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, textTypeQuery, languageQuery)
	err := query(db, mailText)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-IiJjm", "Errors.IAM.CustomMessageText.NotExisting")
	}
	return mailText, err
}

func PutMailText(db *gorm.DB, table string, mailText *model.MailTextView) error {
	save := repository.PrepareSave(table)
	return save(db, mailText)
}

func DeleteMailText(db *gorm.DB, table, aggregateID string, textType string, language string) error {
	aggregateIDSearch := repository.Key{Key: model.MailTextSearchKey(iam_model.MailTextSearchKeyAggregateID), Value: aggregateID}
	textTypeSearch := repository.Key{Key: model.MailTextSearchKey(iam_model.MailTextSearchKeyMailTextType), Value: textType}
	languageSearch := repository.Key{Key: model.MailTextSearchKey(iam_model.MailTextSearchKeyLanguage), Value: language}
	delete := repository.PrepareDeleteByKeys(table, aggregateIDSearch, textTypeSearch, languageSearch)
	return delete(db)
}

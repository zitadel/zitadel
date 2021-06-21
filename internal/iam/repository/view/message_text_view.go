package view

import (
	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func GetMessageTexts(db *gorm.DB, table string, aggregateID string) ([]*model.MessageTextView, error) {
	texts := make([]*model.MessageTextView, 0)
	queries := []*iam_model.MessageTextSearchQuery{
		{
			Key:    iam_model.MessageTextSearchKeyAggregateID,
			Value:  aggregateID,
			Method: domain.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.MessageTextSearchRequest{Queries: queries})
	_, err := query(db, &texts)
	if err != nil {
		return nil, err
	}
	return texts, nil
}

func GetMessageTextByIDs(db *gorm.DB, table, aggregateID, textType, lang string) (*model.MessageTextView, error) {
	mailText := new(model.MessageTextView)
	aggregateIDQuery := &model.MessageTextSearchQuery{Key: iam_model.MessageTextSearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	textTypeQuery := &model.MessageTextSearchQuery{Key: iam_model.MessageTextSearchKeyMessageTextType, Value: textType, Method: domain.SearchMethodEquals}
	languageQuery := &model.MessageTextSearchQuery{Key: iam_model.MessageTextSearchKeyLanguage, Value: lang, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, textTypeQuery, languageQuery)
	err := query(db, mailText)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-IiJjm", "Errors.IAM.CustomMessageText.NotExisting")
	}
	return mailText, err
}

func PutMessageText(db *gorm.DB, table string, mailText *model.MessageTextView) error {
	save := repository.PrepareSave(table)
	return save(db, mailText)
}

func DeleteMessageText(db *gorm.DB, table, aggregateID, textType, lang string) error {
	aggregateIDSearch := repository.Key{Key: model.MessageTextSearchKey(iam_model.MessageTextSearchKeyAggregateID), Value: aggregateID}
	textTypeSearch := repository.Key{Key: model.MessageTextSearchKey(iam_model.MessageTextSearchKeyMessageTextType), Value: textType}
	languageSearch := repository.Key{Key: model.MessageTextSearchKey(iam_model.MessageTextSearchKeyLanguage), Value: lang}
	delete := repository.PrepareDeleteByKeys(table, aggregateIDSearch, textTypeSearch, languageSearch)
	return delete(db)
}

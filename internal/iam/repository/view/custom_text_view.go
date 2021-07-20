package view

import (
	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func GetCustomTexts(db *gorm.DB, table string, aggregateID, template, lang string) ([]*model.CustomTextView, error) {
	texts := make([]*model.CustomTextView, 0)
	queries := []*iam_model.CustomTextSearchQuery{
		{
			Key:    iam_model.CustomTextSearchKeyAggregateID,
			Value:  aggregateID,
			Method: domain.SearchMethodEquals,
		},
		{
			Key:    iam_model.CustomTextSearchKeyTemplate,
			Value:  template,
			Method: domain.SearchMethodEquals,
		},
		{
			Key:    iam_model.CustomTextSearchKeyLanguage,
			Value:  lang,
			Method: domain.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.CustomTextSearchRequest{Queries: queries})
	_, err := query(db, &texts)
	if err != nil {
		return nil, err
	}
	return texts, nil
}

func GetCustomTextsByAggregateIDAndTemplate(db *gorm.DB, table string, aggregateID, template string) ([]*model.CustomTextView, error) {
	texts := make([]*model.CustomTextView, 0)
	queries := []*iam_model.CustomTextSearchQuery{
		{
			Key:    iam_model.CustomTextSearchKeyAggregateID,
			Value:  aggregateID,
			Method: domain.SearchMethodEquals,
		},
		{
			Key:    iam_model.CustomTextSearchKeyTemplate,
			Value:  template,
			Method: domain.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.CustomTextSearchRequest{Queries: queries})
	_, err := query(db, &texts)
	if err != nil {
		return nil, err
	}
	return texts, nil
}

func CustomTextByIDs(db *gorm.DB, table, aggregateID, template, lang, key string) (*model.CustomTextView, error) {
	customText := new(model.CustomTextView)
	aggregateIDQuery := &model.CustomTextSearchQuery{Key: iam_model.CustomTextSearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	textTypeQuery := &model.CustomTextSearchQuery{Key: iam_model.CustomTextSearchKeyTemplate, Value: template, Method: domain.SearchMethodEquals}
	languageQuery := &model.CustomTextSearchQuery{Key: iam_model.CustomTextSearchKeyLanguage, Value: lang, Method: domain.SearchMethodEquals}
	keyQuery := &model.CustomTextSearchQuery{Key: iam_model.CustomTextSearchKeyKey, Value: key, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, textTypeQuery, languageQuery, keyQuery)
	err := query(db, customText)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-8nUU3", "Errors.CustomCustomText.NotExisting")
	}
	return customText, err
}

func PutCustomText(db *gorm.DB, table string, customText *model.CustomTextView) error {
	save := repository.PrepareSave(table)
	return save(db, customText)
}

func DeleteCustomText(db *gorm.DB, table, aggregateID, template, lang, key string) error {
	aggregateIDSearch := repository.Key{Key: model.CustomTextSearchKey(iam_model.CustomTextSearchKeyAggregateID), Value: aggregateID}
	templateSearch := repository.Key{Key: model.CustomTextSearchKey(iam_model.CustomTextSearchKeyTemplate), Value: template}
	languageSearch := repository.Key{Key: model.CustomTextSearchKey(iam_model.CustomTextSearchKeyLanguage), Value: lang}
	keySearch := repository.Key{Key: model.CustomTextSearchKey(iam_model.CustomTextSearchKeyKey), Value: key}
	delete := repository.PrepareDeleteByKeys(table, aggregateIDSearch, templateSearch, keySearch, languageSearch)
	return delete(db)
}

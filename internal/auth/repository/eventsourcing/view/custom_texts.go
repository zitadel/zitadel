package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	customTextTable = "auth.custom_texts"
)

func (v *View) CustomTextByIDs(aggregateID, template, lang, key string) (*model.CustomTextView, error) {
	return view.CustomTextByIDs(v.Db, customTextTable, aggregateID, template, lang, key)
}

func (v *View) CustomTextsByAggregateIDAndTemplateAndLand(aggregateID, template, lang string) ([]*model.CustomTextView, error) {
	return view.GetCustomTexts(v.Db, customTextTable, aggregateID, template, lang)
}

func (v *View) CustomTextsByAggregateIDAndTemplate(aggregateID, template string) ([]*model.CustomTextView, error) {
	return view.GetCustomTextsByAggregateIDAndTemplate(v.Db, customTextTable, aggregateID, template)
}

func (v *View) PutCustomText(template *model.CustomTextView, event *models.Event) error {
	err := view.PutCustomText(v.Db, customTextTable, template)
	if err != nil {
		return err
	}
	return v.ProcessedCustomTextSequence(event)
}

func (v *View) DeleteCustomText(aggregateID, textType, lang, key string, event *models.Event) error {
	err := view.DeleteCustomText(v.Db, customTextTable, aggregateID, textType, lang, key)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedCustomTextSequence(event)
}

func (v *View) DeleteCustomTextTemplate(aggregateID, textType, lang string, event *models.Event) error {
	err := view.DeleteCustomTextTemplate(v.Db, customTextTable, aggregateID, textType, lang)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedCustomTextSequence(event)
}

func (v *View) GetLatestCustomTextSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(customTextTable)
}

func (v *View) ProcessedCustomTextSequence(event *models.Event) error {
	return v.saveCurrentSequence(customTextTable, event)
}

func (v *View) UpdateCustomTextSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(customTextTable)
}

func (v *View) GetLatestCustomTextFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(customTextTable, sequence)
}

func (v *View) ProcessedCustomTextFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

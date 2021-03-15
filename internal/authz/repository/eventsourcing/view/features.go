package view

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/features/repository/view"
	"github.com/caos/zitadel/internal/features/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	featuresTable = "authz.features"
)

func (v *View) AllDefaultFeatures() ([]*model.FeaturesView, error) {
	return view.GetDefaultFeatures(v.Db, featuresTable)
}

func (v *View) FeaturesByAggregateID(aggregateID string) (*model.FeaturesView, error) {
	return view.GetFeaturesByAggregateID(v.Db, featuresTable, aggregateID)
}

func (v *View) PutFeatures(features *model.FeaturesView, event *models.Event) error {
	err := view.PutFeatures(v.Db, featuresTable, features)
	if err != nil {
		return err
	}
	return v.ProcessedFeaturesSequence(event)
}

func (v *View) PutFeaturesList(features []*model.FeaturesView, event *models.Event) error {
	err := view.PutFeaturesList(v.Db, featuresTable, features...)
	if err != nil {
		return err
	}
	return v.ProcessedFeaturesSequence(event)
}

func (v *View) GetLatestFeaturesSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(featuresTable)
}

func (v *View) ProcessedFeaturesSequence(event *models.Event) error {
	return v.saveCurrentSequence(featuresTable, event)
}

func (v *View) UpdateFeaturesSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(featuresTable)
}

func (v *View) GetLatestFeaturesFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(featuresTable, sequence)
}

func (v *View) ProcessedFeaturesFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

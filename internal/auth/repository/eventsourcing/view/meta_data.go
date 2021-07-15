package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	metaDataTable = "auth.meta_data"
)

func (v *View) MetaDataByKey(aggregateID, key string) (*model.MetaDataView, error) {
	return view.MetaDataByKey(v.Db, metaDataTable, aggregateID, key)
}

func (v *View) MetaDataListByAggregateID(aggregateID string) ([]*model.MetaDataView, error) {
	return view.GetMetaDataList(v.Db, metaDataTable, aggregateID)
}

func (v *View) PutMetaData(template *model.MetaDataView, event *models.Event) error {
	err := view.PutMetaData(v.Db, metaDataTable, template)
	if err != nil {
		return err
	}
	return v.ProcessedMetaDataSequence(event)
}

func (v *View) DeleteMetaData(aggregateID, key string, event *models.Event) error {
	err := view.DeleteMetaData(v.Db, metaDataTable, aggregateID, key)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedMetaDataSequence(event)
}

func (v *View) DeleteMetaDataByAggregateID(aggregateID string, event *models.Event) error {
	err := view.DeleteMetaDataByAggregateID(v.Db, metaDataTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedMetaDataSequence(event)
}
func (v *View) GetLatestMetaDataSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(metaDataTable)
}

func (v *View) ProcessedMetaDataSequence(event *models.Event) error {
	return v.saveCurrentSequence(metaDataTable, event)
}

func (v *View) UpdateMetaDataSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(metaDataTable)
}

func (v *View) GetLatestMetaDataFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(metaDataTable, sequence)
}

func (v *View) ProcessedMetaDataFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

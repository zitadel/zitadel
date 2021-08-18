package view

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	metadataTable = "management.metadata"
)

func (v *View) MetadataByKey(aggregateID, key string) (*model.MetadataView, error) {
	return view.MetadataByKey(v.Db, metadataTable, aggregateID, key)
}

func (v *View) MetadataByKeyAndResourceOwner(aggregateID, resourceOwner, key string) (*model.MetadataView, error) {
	return view.MetadataByKeyAndResourceOwner(v.Db, metadataTable, aggregateID, resourceOwner, key)
}

func (v *View) MetadataListByAggregateID(aggregateID string) ([]*model.MetadataView, error) {
	return view.GetMetadataList(v.Db, metadataTable, aggregateID)
}

func (v *View) SearchMetadata(request *domain.MetadataSearchRequest) ([]*model.MetadataView, uint64, error) {
	return view.SearchMetadata(v.Db, metadataTable, request)
}

func (v *View) PutMetadata(template *model.MetadataView, event *models.Event) error {
	err := view.PutMetadata(v.Db, metadataTable, template)
	if err != nil {
		return err
	}
	return v.ProcessedMetadataSequence(event)
}

func (v *View) DeleteMetadata(aggregateID, key string, event *models.Event) error {
	err := view.DeleteMetadata(v.Db, metadataTable, aggregateID, key)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedMetadataSequence(event)
}

func (v *View) DeleteMetadataByAggregateID(aggregateID string, event *models.Event) error {
	err := view.DeleteMetadataByAggregateID(v.Db, metadataTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedMetadataSequence(event)
}
func (v *View) GetLatestMetadataSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(metadataTable)
}

func (v *View) ProcessedMetadataSequence(event *models.Event) error {
	return v.saveCurrentSequence(metadataTable, event)
}

func (v *View) UpdateMetadataSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(metadataTable)
}

func (v *View) GetLatestMetadataFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(metadataTable, sequence)
}

func (v *View) ProcessedMetadataFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

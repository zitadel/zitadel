package model

import (
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/domain"
	usr_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	MetadataKeyAggregateID   = "aggregate_id"
	MetadataKeyResourceOwner = "resource_owner"
	MetadataKeyKey           = "key"
	MetadataKeyValue         = "value"
)

type MetadataView struct {
	AggregateID   string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	ResourceOwner string    `json:"-" gorm:"column:resource_owner"`
	CreationDate  time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate    time.Time `json:"-" gorm:"column:change_date"`

	Key   string `json:"key" gorm:"column:key;primary_key"`
	Value []byte `json:"value" gorm:"column:value"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func MetadataViewsToDomain(texts []*MetadataView) []*domain.Metadata {
	result := make([]*domain.Metadata, len(texts))
	for i, text := range texts {
		result[i] = MetadataViewToDomain(text)
	}
	return result
}

func MetadataViewToDomain(data *MetadataView) *domain.Metadata {
	return &domain.Metadata{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  data.AggregateID,
			Sequence:     data.Sequence,
			CreationDate: data.CreationDate,
			ChangeDate:   data.ChangeDate,
		},
		Key:   data.Key,
		Value: data.Value,
	}
}

func (md *MetadataView) AppendEvent(event *models.Event) (err error) {
	md.Sequence = event.Sequence
	switch event.Type {
	case usr_model.UserMetadataSet:
		md.setRootData(event)
		err = md.SetData(event)
	}
	return err
}

func (md *MetadataView) setRootData(event *models.Event) {
	md.AggregateID = event.AggregateID
	md.ResourceOwner = event.ResourceOwner
	md.ChangeDate = event.CreationDate
	md.Sequence = event.Sequence
}

func (md *MetadataView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, md); err != nil {
		logging.Log("MODEL-3n9fs").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-5CVaR", "Could not unmarshal data")
	}
	return nil
}

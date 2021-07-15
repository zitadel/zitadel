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
	MetaDataKeyAggregateID   = "aggregate_id"
	MetaDataKeyResourceOwner = "resource_owner"
	MetaDataKeyKey           = "key"
	MetaDataKeyValue         = "value"
)

type MetaDataView struct {
	AggregateID   string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	ResourceOwner string    `json:"-" gorm:"column:resource_owner"`
	CreationDate  time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate    time.Time `json:"-" gorm:"column:change_date"`

	Key   string `json:"key" gorm:"column:key;primary_key"`
	Value string `json:"value" gorm:"column:value"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func MetaDataViewsToDomain(texts []*MetaDataView) []*domain.MetaData {
	result := make([]*domain.MetaData, len(texts))
	for i, text := range texts {
		result[i] = MetaDataViewToDomain(text)
	}
	return result
}

func MetaDataViewToDomain(data *MetaDataView) *domain.MetaData {
	return &domain.MetaData{
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

func (i *MetaDataView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	switch event.Type {
	case usr_model.UserMetaDataSet:
		i.setRootData(event)
		err = i.SetData(event)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *MetaDataView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
	r.ResourceOwner = event.ResourceOwner
	r.ChangeDate = event.CreationDate
	r.Sequence = event.Sequence
}

func (r *MetaDataView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("MODEL-3n9fs").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-5CVaR", "Could not unmarshal data")
	}
	return nil
}

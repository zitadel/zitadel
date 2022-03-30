package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	user_repo "github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/user/model"
)

const (
	ExternalIDPKeyExternalUserID = "external_user_id"
	ExternalIDPKeyUserID         = "user_id"
	ExternalIDPKeyIDPConfigID    = "idp_config_id"
	ExternalIDPKeyResourceOwner  = "resource_owner"
)

type ExternalIDPView struct {
	ExternalUserID  string    `json:"userID" gorm:"column:external_user_id;primary_key"`
	IDPConfigID     string    `json:"idpConfigID" gorm:"column:idp_config_id;primary_key"`
	UserID          string    `json:"-" gorm:"column:user_id"`
	IDPName         string    `json:"-" gorm:"column:idp_name"`
	UserDisplayName string    `json:"displayName" gorm:"column:user_display_name"`
	CreationDate    time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate      time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner   string    `json:"-" gorm:"column:resource_owner"`
	Sequence        uint64    `json:"-" gorm:"column:sequence"`
	InstanceID      string    `json:"instanceID" gorm:"column:instance_id"`
}

func ExternalIDPViewFromModel(externalIDP *model.ExternalIDPView) *ExternalIDPView {
	return &ExternalIDPView{
		UserID:          externalIDP.UserID,
		IDPConfigID:     externalIDP.IDPConfigID,
		ExternalUserID:  externalIDP.ExternalUserID,
		IDPName:         externalIDP.IDPName,
		UserDisplayName: externalIDP.UserDisplayName,
		Sequence:        externalIDP.Sequence,
		CreationDate:    externalIDP.CreationDate,
		ChangeDate:      externalIDP.ChangeDate,
		ResourceOwner:   externalIDP.ResourceOwner,
	}
}

func ExternalIDPViewToModel(externalIDP *ExternalIDPView) *model.ExternalIDPView {
	return &model.ExternalIDPView{
		UserID:          externalIDP.UserID,
		IDPConfigID:     externalIDP.IDPConfigID,
		ExternalUserID:  externalIDP.ExternalUserID,
		IDPName:         externalIDP.IDPName,
		UserDisplayName: externalIDP.UserDisplayName,
		Sequence:        externalIDP.Sequence,
		CreationDate:    externalIDP.CreationDate,
		ChangeDate:      externalIDP.ChangeDate,
		ResourceOwner:   externalIDP.ResourceOwner,
	}
}

func ExternalIDPViewsToModel(externalIDPs []*ExternalIDPView) []*model.ExternalIDPView {
	result := make([]*model.ExternalIDPView, len(externalIDPs))
	for i, r := range externalIDPs {
		result[i] = ExternalIDPViewToModel(r)
	}
	return result
}

func (i *ExternalIDPView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch eventstore.EventType(event.Type) {
	case user_repo.UserIDPLinkAddedType:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	}
	return err
}

func (r *ExternalIDPView) setRootData(event *models.Event) {
	r.UserID = event.AggregateID
	r.ResourceOwner = event.ResourceOwner
	r.InstanceID = event.InstanceID
}

func (r *ExternalIDPView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-48sfs").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}

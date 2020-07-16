package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/service_account/model"
)

type ServiceAccount struct {
	models.ObjectRoot
}

func (sa *ServiceAccount) AppendEvents(events ...*models.Event) error {
	for _, event := range events {
		if err := sa.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (sa *ServiceAccount) AppendEvent(event *models.Event) (err error) {
	return nil
}

func ServiceAccountFromModel(serviceAccount *model.ServiceAccount) *ServiceAccount {
	return &ServiceAccount{}
}

func ServiceAccountToModel(serviceAccount *ServiceAccount) *model.ServiceAccount {
	return &model.ServiceAccount{}
}

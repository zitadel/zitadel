package handler

import (
	"context"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	org_repo "github.com/caos/zitadel/internal/repository/org"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	idpConfigTable = "adminapi.idp_configs"
)

type IDPConfig struct {
	handler
	subscription *v1.Subscription
}

func newIDPConfig(handler handler) *IDPConfig {
	h := &IDPConfig{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (i *IDPConfig) subscribe() {
	i.subscription = i.es.Subscribe(i.AggregateTypes()...)
	go func() {
		for event := range i.subscription.Events {
			query.ReduceEvent(i, event)
		}
	}()
}

func (i *IDPConfig) Subscription() *v1.Subscription {
	return i.subscription
}

func (i *IDPConfig) ViewModel() string {
	return idpConfigTable
}

func (i *IDPConfig) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.IAMAggregate}
}

func (i *IDPConfig) CurrentSequence() (uint64, error) {
	sequence, err := i.view.GetLatestIDPConfigSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (i *IDPConfig) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := i.view.GetLatestIDPConfigSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(i.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (i *IDPConfig) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = i.processIDPConfig(event)
	}
	return err
}

func (i *IDPConfig) processIDPConfig(event *es_models.Event) (err error) {
	idp := new(iam_view_model.IDPConfigView)
	switch event.Type {
	case model.IDPConfigAdded:
		err = idp.AppendEvent(iam_model.IDPProviderTypeSystem, event)
	case model.IDPConfigChanged,
		model.OIDCIDPConfigAdded,
		model.OIDCIDPConfigChanged,
		es_models.EventType(iam_repo.IDPAuthConnectorConfigAddedEventType),
		es_models.EventType(org_repo.IDPAuthConnectorConfigAddedEventType),
		es_models.EventType(iam_repo.IDPAuthConnectorConfigChangedEventType),
		es_models.EventType(org_repo.IDPAuthConnectorConfigChangedEventType),
		es_models.EventType(iam_repo.IDPAuthConnectorMachineUserRemovedEventType),
		es_models.EventType(org_repo.IDPAuthConnectorMachineUserRemovedEventType):
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		idp, err = i.view.IDPConfigByID(idp.IDPConfigID)
		if err != nil {
			return err
		}
		err = idp.AppendEvent(iam_model.IDPProviderTypeSystem, event)
		if err != nil {
			err = i.fillData(idp)
		}
	case model.IDPConfigDeactivated,
		model.IDPConfigReactivated:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		idp, err = i.view.IDPConfigByID(idp.IDPConfigID)
		if err != nil {
			return err
		}
		err = idp.AppendEvent(iam_model.IDPProviderTypeSystem, event)
	case model.IDPConfigRemoved:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		return i.view.DeleteIDPConfig(idp.IDPConfigID, event)
	default:
		return i.view.ProcessedIDPConfigSequence(event)
	}
	if err != nil {
		return err
	}
	return i.view.PutIDPConfig(idp, event)
}

func (i *IDPConfig) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Mslo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp config handler")
	return spooler.HandleError(event, err, i.view.GetLatestIDPConfigFailedEvent, i.view.ProcessedIDPConfigFailedEvent, i.view.ProcessedIDPConfigSequence, i.errorCountUntilSkip)
}

func (i *IDPConfig) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateIDPConfigSpoolerRunTimestamp)
}

func (i *IDPConfig) fillData(idp *iam_view_model.IDPConfigView) error {
	if idp.AuthConnectorMachineID == "" {
		idp.AuthConnectorMachineName = ""
		return nil
	}

	user, err := i.getUserByID(idp.AuthConnectorMachineID)
	if err != nil {
		return err
	}
	if user.MachineView != nil {
		idp.AuthConnectorMachineName = user.MachineView.Name
	}
	return nil
}

func (i *IDPConfig) getUserByID(userID string) (*usr_view_model.UserView, error) {
	user, usrErr := i.view.UserByID(userID)
	if usrErr != nil && !caos_errs.IsNotFound(usrErr) {
		return nil, usrErr
	}
	if user == nil {
		user = &usr_view_model.UserView{}
	}
	events, err := i.getUserEvents(userID, user.Sequence)
	if err != nil {
		return user, usrErr
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return user, nil
		}
	}
	if userCopy.State == int32(usr_model.UserStateDeleted) {
		return nil, caos_errs.ThrowNotFound(nil, "HANDLER-GAdg2", "Errors.User.NotFound")
	}
	return &userCopy, nil
}

func (i *IDPConfig) getUserEvents(userID string, sequence uint64) ([]*es_models.Event, error) {
	query, err := view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}

	return i.es.FilterEvents(context.Background(), query)
}

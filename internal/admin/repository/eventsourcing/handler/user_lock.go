package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	userLockerTable = "adminapi.user_lock"
)

type UserLocker struct {
	handler
	subscription *v1.Subscription
	command      *command.Commands
}

func newUserLocker(
	handler handler,
	command *command.Commands,
) *UserLocker {
	h := &UserLocker{
		handler: handler,
		command: command,
	}

	h.subscribe()

	return h
}

func (u *UserLocker) subscribe() {
	u.subscription = u.es.Subscribe(u.AggregateTypes()...)
	go func() {
		for event := range u.subscription.Events {
			query.ReduceEvent(u, event)
		}
	}()
}

func (u *UserLocker) ViewModel() string {
	return userLockerTable
}

func (u *UserLocker) Subscription() *v1.Subscription {
	return u.subscription
}

func (u *UserLocker) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{es_model.UserAggregate}
}

func (u *UserLocker) CurrentSequence() (uint64, error) {
	sequence, err := u.view.GetLatestUserSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (u *UserLocker) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(u.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *UserLocker) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case es_model.UserAggregate:
		return u.ProcessUser(event)
	default:
		return nil
	}
}

func (u *UserLocker) ProcessUser(event *es_models.Event) (err error) {
	userLock := new(view_model.UserLockView)
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered,
		es_model.HumanRegistered,
		es_model.MachineAdded,
		es_model.HumanAdded:
		err = userLock.AppendEvent(event)
		if err != nil {
			return err
		}
	case es_model.UserPasswordCheckFailed:
		userLock, err = u.view.UserLockByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = userLock.AppendEvent(event)
	case es_model.UserPasswordCheckSucceeded:
		userLock, err = u.view.UserLockByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = userLock.AppendEvent(event)
	case es_model.UserRemoved:
		return u.view.DeleteUserLock(event.AggregateID, event)
	default:
		return u.view.ProcessedUserLockSequence(event)
	}
	if err != nil {
		return err
	}
	err = u.view.PutUserLock(userLock, event)
	if err != nil {
		return err
	}
	if userLock.State != int32(view_model.UserStateLocked) && userLock.PasswordCheckFailedCount > 0 {
		policy, err := u.getLockoutPolicy(context.Background(), userLock.ResourceOwner)
		if err != nil {
			return err
		}
		if policy.MaxPasswordAttempts == 0 || userLock.PasswordCheckFailedCount < policy.MaxPasswordAttempts {
			return nil
		}
		_, err = u.command.LockUser(context.Background(), userLock.UserID, userLock.ResourceOwner)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UserLocker) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-vLmwQ", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserFailedEvent, u.view.ProcessedUserFailedEvent, u.view.ProcessedUserSequence, u.errorCountUntilSkip)
}

func (u *UserLocker) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateUserSpoolerRunTimestamp)
}

func (u *UserLocker) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-kVLb2", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *UserLocker) getIAMByID(ctx context.Context) (*iam_model.IAM, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, 0)
	if err != nil {
		return nil, err
	}
	iam := &model.IAM{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: domain.IAMID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, iam.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	return model.IAMToModel(iam), nil
}

func (u *UserLocker) getLockoutPolicy(ctx context.Context, orgID string) (*iam_model.LockoutPolicy, error) {
	org, err := u.getOrgByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if org.PasswordLockoutPolicy != nil {
		return org.PasswordLockoutPolicy, nil
	}
	iam, err := u.getIAMByID(ctx)
	if err != nil {
		return nil, err
	}
	return iam.DefaultPasswordLockoutPolicy, nil
}

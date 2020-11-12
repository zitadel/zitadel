package policy

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	OrgIAMPolicyAddedEventType = "policy.org.iam.added"
)

type OrgIAMPolicyAggregate struct {
	eventstore.Aggregate

	UserLoginMustBeDomain bool
}

type OrgIAMPolicyReadModel struct {
	eventstore.ReadModel

	UserLoginMustBeDomain bool
}

func (rm *OrgIAMPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *OrgIAMPolicyAddedEvent:
			rm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		}
	}
	return rm.ReadModel.Reduce()
}

type OrgIAMPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain bool `json:"allowUsernamePassword"`
}

func (e *OrgIAMPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *OrgIAMPolicyAddedEvent) Data() interface{} {
	return e
}

func NewOrgIAMPolicyAddedEvent(
	ctx context.Context,
	userLoginMustBeDomain bool,
) *OrgIAMPolicyAddedEvent {

	return &OrgIAMPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgIAMPolicyAddedEventType,
		),
		UserLoginMustBeDomain: userLoginMustBeDomain,
	}
}

func OrgIAMPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OrgIAMPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-TvSmA", "unable to unmarshal policy")
	}

	return e, nil
}

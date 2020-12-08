package policy

import (
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

type OrgIAMPolicyWriteModel struct {
	eventstore.WriteModel

	UserLoginMustBeDomain bool
}

func (wm *OrgIAMPolicyWriteModel) Reduce() error {
	return errors.ThrowUnimplemented(nil, "POLIC-o0vMl", "reduce unimpelemnted")
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
	base *eventstore.BaseEvent,
	userLoginMustBeDomain bool,
) *OrgIAMPolicyAddedEvent {

	return &OrgIAMPolicyAddedEvent{
		BaseEvent:             *base,
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

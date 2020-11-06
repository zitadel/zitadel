package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	OrgIAMPolicyAddedEventType = "policy.org.iam.added"
)

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
	service string,
	userLoginMustBeDomain bool,
) *OrgIAMPolicyAddedEvent {

	return &OrgIAMPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			OrgIAMPolicyAddedEventType,
		),
		UserLoginMustBeDomain: userLoginMustBeDomain,
	}
}

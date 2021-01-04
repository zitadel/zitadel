package org

import (
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LoginPolicyAddedEventType   = orgEventTypePrefix + policy.LoginPolicyAddedEventType
	LoginPolicyChangedEventType = orgEventTypePrefix + policy.LoginPolicyChangedEventType
)

type LoginPolicyAddedEvent struct {
	policy.LoginPolicyAddedEvent
}

type LoginPolicyChangedEvent struct {
	policy.LoginPolicyChangedEvent
}

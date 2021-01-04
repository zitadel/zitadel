package org

import (
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LabelPolicyAddedEventType   = orgEventTypePrefix + policy.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = orgEventTypePrefix + policy.LabelPolicyChangedEventType
)

type LabelPolicyAddedEvent struct {
	policy.LabelPolicyAddedEvent
}

type LabelPolicyChangedEvent struct {
	policy.LabelPolicyChangedEvent
}

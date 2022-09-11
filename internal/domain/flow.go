package domain

import "strconv"

type FlowState int32

const (
	FlowStateActive FlowState = iota
	FlowStateInactive
	flowStateCount
)

func (s FlowState) Valid() bool {
	return s >= 0 && s < flowStateCount
}

type FlowType int32

const (
	FlowTypeUnspecified FlowType = iota
	FlowTypeExternalAuthentication
	FlowTypeCustomiseToken
	flowTypeCount
)

func (s FlowType) Valid() bool {
	return s > 0 && s < flowTypeCount
}

func (s FlowType) HasTrigger(triggerType TriggerType) bool {
	switch triggerType {
	case TriggerTypePostAuthentication:
		return s == FlowTypeExternalAuthentication
	case TriggerTypePreCreation:
		return s == FlowTypeExternalAuthentication
	case TriggerTypePostCreation:
		return s == FlowTypeExternalAuthentication
	default:
		return false
	}
}

func (s FlowType) ID() string {
	if s < 0 && s >= flowTypeCount {
		return FlowTypeUnspecified.ID()
	}
	return strconv.Itoa(int(s))
}

func (s FlowType) LocalizationKey() string {
	if s < 0 && s >= flowTypeCount {
		return FlowTypeUnspecified.LocalizationKey()
	}

	switch s {
	case FlowTypeExternalAuthentication:
		return "Action.Flow.Type.ExternalAuthentication"
	case FlowTypeCustomiseToken:
		return "Action.Flow.Type.CustomiseToken"
	default:
		return "Action.Flow.Type.Unspecified"
	}
}

type TriggerType int32

const (
	TriggerTypeUnspecified TriggerType = iota
	TriggerTypePostAuthentication
	TriggerTypePreCreation
	TriggerTypePostCreation
	TriggerTypePreUserinfoCreation
	TriggerTypePreAccessTokenCreation
	triggerTypeCount
)

func (s TriggerType) Valid() bool {
	return s >= 0 && s < triggerTypeCount
}

func (s TriggerType) ID() string {
	if !s.Valid() {
		return TriggerTypeUnspecified.ID()
	}
	return strconv.Itoa(int(s))
}

func (s TriggerType) LocalizationKey() string {
	if !s.Valid() {
		return FlowTypeUnspecified.LocalizationKey()
	}

	switch s {
	case TriggerTypePostAuthentication:
		return "Action.Flow.Type.ExternalAuthentication"
	case TriggerTypePreCreation:
		return "Action.Flow.Type.PreCreation"
	case TriggerTypePostCreation:
		return "Action.Flow.Type.PostCreation"
	case TriggerTypePreUserinfoCreation:
		return "Action.Flow.Type.PreUserinfoCreation"
	case TriggerTypePreAccessTokenCreation:
		return "Action.Flow.Type.PreAccessTokenCreation"
	default:
		return "Action.Flow.Type.Unspecified"
	}
}

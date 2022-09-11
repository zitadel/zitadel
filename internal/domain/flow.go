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
	for _, trigger := range s.TriggerTypess() {
		if trigger == triggerType {
			return true
		}
	}
	return false
}

func (s FlowType) TriggerTypess() []TriggerType {
	switch s {
	case FlowTypeExternalAuthentication:
		return []TriggerType{
			TriggerTypePostAuthentication,
			TriggerTypePreCreation,
			TriggerTypePostCreation,
		}
	case FlowTypeCustomiseToken:
		return []TriggerType{
			TriggerTypePreUserinfoCreation,
			TriggerTypePreAccessTokenCreation,
		}
	default:
		return nil
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
		return "Action.TriggerType.ExternalAuthentication"
	case TriggerTypePreCreation:
		return "Action.TriggerType.PreCreation"
	case TriggerTypePostCreation:
		return "Action.TriggerType.PostCreation"
	case TriggerTypePreUserinfoCreation:
		return "Action.TriggerType.PreUserinfoCreation"
	case TriggerTypePreAccessTokenCreation:
		return "Action.TriggerType.PreAccessTokenCreation"
	default:
		return "Action.TriggerType.Unspecified"
	}
}

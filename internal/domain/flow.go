package domain

type FlowState int32

const (
	FlowStateUnspecified FlowState = iota
	FlowStateActive
	FlowStateInactive
	FlowStateRemoved
	flowStateCount
)

func (s FlowState) Valid() bool {
	return s >= 0 && s < flowStateCount
}

func (s FlowState) Exists() bool {
	return s != FlowStateUnspecified && s != FlowStateRemoved
}

type FlowType int32

const (
	FlowTypeUnspecified FlowType = iota
	FlowTypeRegister
	flowTypeCount
)

func (s FlowType) Valid() bool {
	return s >= 0 && s < flowTypeCount
}

func (s FlowType) HasTrigger(triggerType TriggerType) bool {
	switch triggerType {
	case TriggerTypePreRegister:
		return s == FlowTypeRegister
	default:
		return false
	}
}

type TriggerType int32

const (
	TriggerTypeUnspecified TriggerType = iota
	TriggerTypePreRegister
	triggerTypeCount
)

func (s TriggerType) Valid() bool {
	return s >= 0 && s < triggerTypeCount
}

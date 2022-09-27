package domain

type SMSConfigState int32

const (
	SMSConfigStateUnspecified SMSConfigState = iota
	SMSConfigStateActive
	SMSConfigStateInactive
	SMSConfigStateRemoved
)

func (s SMSConfigState) Exists() bool {
	return s != SMSConfigStateUnspecified && s != SMSConfigStateRemoved
}

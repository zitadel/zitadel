package domain

type SMSConfigState int32

const (
	SMSConfigStateUnspecified SMSConfigState = iota
	SMSConfigStateActive
	SMSConfigStateRemoved
)

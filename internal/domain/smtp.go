package domain

type SMTPConfigState int32

const (
	SMTPConfigStateUnspecified SMTPConfigState = iota
	SMTPConfigStateActive
	SMTPConfigStateRemoved
	SMTPConfigStateInactive
)

func (s SMTPConfigState) Exists() bool {
	return s != SMTPConfigStateUnspecified && s != SMTPConfigStateRemoved
}

package domain

type SMTPConfigState int32

const (
	SMTPConfigStateUnspecified SMTPConfigState = iota
	SMTPConfigStateActive
)

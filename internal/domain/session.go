package domain

type SessionState int32

const (
	SessionStateUnspecified SessionState = iota
	SessionStateActive
	SessionStateTerminated
)

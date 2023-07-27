package domain

type OIDCSessionState int32

const (
	OIDCSessionStateUnspecified OIDCSessionState = iota
	OIDCSessionStateActive
	OIDCSessionStateTerminated
)

package domain

type SAMLSessionState int32

const (
	SAMLSessionStateUnspecified SAMLSessionState = iota
	SAMLSessionStateActive
	SAMLSessionStateTerminated
)
